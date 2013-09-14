package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
	"github.com/henyouqian/lwUtil"
	"net/http"
	"time"
)

type Match struct {
	Id     uint32
	Name   string
	GameId uint32
	Begin  int64
	End    int64
	Sort   string
}

func newMatch(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "POST")

	session, err := findSession(w, r)
	lwutil.CheckError("err_auth", err)
	checkAdmin(session)

	appid := session.Appid
	if appid == 0 {
		lwutil.SendError("err_auth", "Please login with app secret")
	}

	// input
	type Input struct {
		Name   string
		GameId uint32
		Begin  string
		End    string
	}
	input := Input{}
	lwutil.DecodeRequestBody(r, &input)

	if input.Name == "" || input.Begin == "" || input.End == "" || input.GameId == 0 {
		lwutil.SendError("err_input", "Missing Name || Begin || End || Gameid")
	}

	// game info
	game, err := findGame(input.GameId, appid)
	lwutil.CheckError("err_game", err)

	// times
	const timeform = "2006-01-02 15:04:05"
	begin, err := time.ParseInLocation(timeform, input.Begin, time.Local)
	lwutil.CheckError("", err)
	end, err := time.ParseInLocation(timeform, input.End, time.Local)
	lwutil.CheckError("", err)
	beginUnix := begin.Unix()
	endUnix := end.Unix()

	if endUnix-beginUnix <= 60 {
		lwutil.SendError("err_input", "endUnix - beginUnix must > 60 seconds")
	}
	if time.Now().Unix() > endUnix {
		lwutil.SendError("err_input", "end time before now")
	}

	//
	rc := redisPool.Get()
	defer rc.Close()

	matchId, err := redis.Int(rc.Do("incr", "idGen/match"))
	lwutil.CheckError("", err)

	match := Match{
		uint32(matchId),
		input.Name,
		input.GameId,
		beginUnix,
		endUnix,
		game.Sort,
	}

	matchJson, err := json.Marshal(match)
	lwutil.CheckError("", err)

	key := fmt.Sprintf("%d+%d", appid, matchId)
	rc.Send("hset", "matches", key, matchJson)
	key = fmt.Sprintf("matchesInApp/%d", appid)
	rc.Send("zadd", key, endUnix, matchId)
	rc.Flush()
	for i := 0; i < 2; i++ {
		_, err = rc.Receive()
		lwutil.CheckError("", err)
	}

	// reply
	lwutil.WriteResponse(w, match)
}

func delMatch(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "POST")

	session, err := findSession(w, r)
	lwutil.CheckError("err_auth", err)
	checkAdmin(session)

	appid := session.Appid
	if appid == 0 {
		lwutil.SendError("err_auth", "Please login with app secret")
	}

	// input
	matchIds := make([]int, 0, 8)
	lwutil.DecodeRequestBody(r, &matchIds)

	// redis
	rc := redisPool.Get()
	defer rc.Close()

	key := fmt.Sprintf("matchesInApp/%d", appid)
	params := make([]interface{}, 0, 8)
	params = append(params, key)
	matchIdsItf := make([]interface{}, len(matchIds))
	for i, v := range matchIds {
		matchIdsItf[i] = v
	}
	params = append(params, matchIdsItf...)
	rc.Send("zrem", params...)

	args := make([]interface{}, len(matchIds)+1)
	args[0] = "matches"
	for i, matchId := range matchIds {
		key = fmt.Sprintf("%d+%d", appid, matchId)
		args[i+1] = key
	}
	rc.Send("hdel", args...)
	rc.Flush()

	_, err = rc.Receive()
	lwutil.CheckError("", err)
	delNum, err := rc.Receive()
	lwutil.CheckError("", err)

	// reply
	lwutil.WriteResponse(w, delNum)
}

func listMatch(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "POST")

	session, err := findSession(w, r)
	lwutil.CheckError("err_auth", err)

	appid := session.Appid
	if appid == 0 {
		lwutil.SendError("err_auth", "Please login with app secret")
	}

	nowUnix := time.Now().Unix()

	rc := redisPool.Get()
	defer rc.Close()

	// get matchIds
	key := fmt.Sprintf("matchesInApp/%d", appid)
	matchIdValues, err := redis.Values(rc.Do("zrangebyscore", key, nowUnix, "+inf"))
	lwutil.CheckError("", err)

	args := make([]interface{}, len(matchIdValues)+1)
	args[0] = "matches"
	for i, v := range matchIdValues {
		var id int
		id, err := redis.Int(v, err)
		lwutil.CheckError("", err)
		matchkey := fmt.Sprintf("%d+%d", appid, id)
		args[i+1] = matchkey
	}

	// get match data
	matchesValues, err := redis.Values(rc.Do("hmget", args...))

	matches := make([]interface{}, len(matchesValues))

	for i, v := range matchesValues {
		var match interface{}
		err = json.Unmarshal(v.([]byte), &match)
		lwutil.CheckError("", err)
		matches[i] = match
	}

	lwutil.WriteResponse(w, matches)
}

func startMatch(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "POST")

	session, err := findSession(w, r)
	lwutil.CheckError("err_auth", err)

	appid := session.Appid
	if appid == 0 {
		lwutil.SendError("err_auth", "Please login with app secret")
	}

	// input
	type input struct {
		MatchId uint32
	}
	var in input
	lwutil.DecodeRequestBody(r, &in)

	// redis setup
	rc := redisPool.Get()
	defer rc.Close()

	// playing?
	secretRaw, err := rc.Do("get", fmt.Sprintf("trySecretsRev/%d", in.MatchId))
	lwutil.CheckError("", err)
	if secretRaw != nil {
		secret, err := redis.String(secretRaw, err)
		lwutil.CheckError("", err)
		lwutil.WriteResponse(w, secret)
		return
	}

	// get match info
	key := fmt.Sprintf("%d+%d", appid, in.MatchId)
	matchJson, err := redis.Bytes(rc.Do("hget", "matches", key))
	lwutil.CheckError("err_not_found", err)

	var match Match
	err = json.Unmarshal(matchJson, &match)
	lwutil.CheckError("", err)

	// check time
	now := time.Now().Unix()
	if now < match.Begin || now >= match.End-MATCH_TRY_DURATION_SEC {
		lwutil.SendError("err_time", "now < match.Begin || now >= match.End-MATCH_TRY_DURATION_SEC")
	}

	// played?
	keyFail := fmt.Sprintf("failboard/%d", in.MatchId)
	keyLeaderboard := fmt.Sprintf("leaderboard/%d", in.MatchId)
	rc.Send("zscore", keyLeaderboard, session.Userid)
	rc.Send("sismember", keyFail, session.Userid)
	rc.Flush()
	lbScore, err := rc.Receive()
	lwutil.CheckError("", err)
	inFail, err := redis.Int(rc.Receive())
	lwutil.CheckError("", err)

	glog.Infoln(lbScore, inFail)
	if lbScore != nil || inFail != 0 {
		lwutil.SendError("err_no_try", "no try left")
	}

	// add to failboard
	_, err = rc.Do("sadd", keyFail, session.Userid)
	lwutil.CheckError("", err)

	// new try secret
	trySecret := lwutil.GenUUID()
	rc.Send("setex", fmt.Sprintf("trySecrets/%s", trySecret), MATCH_TRY_DURATION_SEC, in.MatchId)
	rc.Send("setex", fmt.Sprintf("trySecretsRev/%d", in.MatchId), MATCH_TRY_DURATION_SEC, trySecret)
	err = rc.Flush()
	lwutil.CheckError("", err)

	// reply
	lwutil.WriteResponse(w, trySecret)
}

func addScore(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "POST")

	session, err := findSession(w, r)
	lwutil.CheckError("err_auth", err)

	//appid := session.Appid
	//if appid == 0 {
	//	lwutil.SendError("err_auth", "Please login with app secret")
	//}

	// input
	type input struct {
		TrySecret string
		Score     int64
	}
	var in input
	lwutil.DecodeRequestBody(r, &in)

	// redis setup
	rc := redisPool.Get()
	defer rc.Close()

	// use secret to get matchId
	matchIdRaw, err := rc.Do("get", fmt.Sprintf("trySecrets/%s", in.TrySecret))
	lwutil.CheckError("", err)
	if matchIdRaw == nil {
		lwutil.SendError("err_secret", "")
	}
	matchId, err := redis.Int(matchIdRaw, err)
	lwutil.CheckError("", err)

	// del from failboard and add to leaderboard and delete secret
	keyFail := fmt.Sprintf("failboard/%d", matchId)
	keyLeaderboard := fmt.Sprintf("leaderboard/%d", matchId)
	rc.Send("srem", keyFail, session.Userid)
	rc.Send("zadd", keyLeaderboard, in.Score, session.Userid)
	rc.Send("del", fmt.Sprintf("trySecrets/%s", in.TrySecret))
	rc.Flush()
	_, err = rc.Receive()
	lwutil.CheckError("", err)
	_, err = rc.Receive()
	lwutil.CheckError("", err)
	_, err = rc.Receive()
	lwutil.CheckError("", err)

	// reply
	lwutil.WriteResponse(w, matchId)
}

func regMatch() {
	http.Handle("/match/new", lwutil.ReqHandler(newMatch))
	http.Handle("/match/del", lwutil.ReqHandler(delMatch))
	http.Handle("/match/list", lwutil.ReqHandler(listMatch))
	http.Handle("/match/start", lwutil.ReqHandler(startMatch))
	http.Handle("/match/addscore", lwutil.ReqHandler(addScore))
}

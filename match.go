package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	//"github.com/golang/glog"
	"github.com/henyouqian/lwutil"
	"net/http"
	"time"
)

type Match struct {
	Id       uint32
	Name     string
	GameId   uint32
	Begin    int64
	End      int64
	Sort     string
	TryMax   uint32
	TryPrice uint32
}

func newMatch(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "POST")

	rc := redisPool.Get()
	defer rc.Close()

	session, err := findSession(w, r, rc)
	lwutil.CheckError(err, "err_auth")
	checkAdmin(session)

	appid := session.Appid
	if appid == 0 {
		lwutil.SendError("err_auth", "Please login with app secret")
	}

	//input
	var in struct {
		Name     string
		GameId   uint32
		Begin    string
		End      string
		TryMax   uint32
		TryPrice uint32
	}
	err = lwutil.DecodeRequestBody(r, &in)
	lwutil.CheckError(err, "err_decode_body")

	if in.Name == "" || in.Begin == "" || in.End == "" || in.GameId == 0 {
		lwutil.SendError("err_input", "Missing Name || Begin || End || Gameid")
	}
	if in.TryMax == 0 {
		in.TryMax = 1
	}

	//game info
	game, err := findGame(in.GameId, appid)
	lwutil.CheckError(err, "err_game")

	//times
	const timeform = "2006-01-02 15:04:05"
	begin, err := time.ParseInLocation(timeform, in.Begin, time.Local)
	lwutil.CheckError(err, "")
	end, err := time.ParseInLocation(timeform, in.End, time.Local)
	lwutil.CheckError(err, "")
	beginUnix := begin.Unix()
	endUnix := end.Unix()

	if endUnix-beginUnix <= 60 {
		lwutil.SendError("err_input", "endUnix - beginUnix must > 60 seconds")
	}
	if time.Now().Unix() > endUnix {
		lwutil.SendError("err_input", "end time before now")
	}

	//
	matchId, err := redis.Int(rc.Do("incr", "idGen/match"))
	lwutil.CheckError(err, "")

	match := Match{
		uint32(matchId),
		in.Name,
		in.GameId,
		beginUnix,
		endUnix,
		game.Sort,
		in.TryMax,
		in.TryPrice,
	}

	matchJson, err := json.Marshal(match)
	lwutil.CheckError(err, "")

	key := fmt.Sprintf("%d+%d", appid, matchId)
	rc.Send("hset", "matches", key, matchJson)
	key = fmt.Sprintf("matchesInApp/%d", appid)
	rc.Send("zadd", key, endUnix, matchId)
	rc.Flush()
	for i := 0; i < 2; i++ {
		_, err = rc.Receive()
		lwutil.CheckError(err, "")
	}

	//reply
	lwutil.WriteResponse(w, match)
}

func delMatch(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "POST")

	rc := redisPool.Get()
	defer rc.Close()

	session, err := findSession(w, r, rc)
	lwutil.CheckError(err, "err_auth")
	checkAdmin(session)

	appid := session.Appid
	if appid == 0 {
		lwutil.SendError("err_auth", "Please login with app secret")
	}

	//input
	matchIds := make([]int, 0, 8)
	err = lwutil.DecodeRequestBody(r, &matchIds)
	lwutil.CheckError(err, "err_decode_body")

	//redis
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
	lwutil.CheckError(err, "")
	delNum, err := rc.Receive()
	lwutil.CheckError(err, "")

	//reply
	lwutil.WriteResponse(w, delNum)
}

func listMatch(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "POST")

	rc := redisPool.Get()
	defer rc.Close()

	session, err := findSession(w, r, rc)
	lwutil.CheckError(err, "err_auth")

	appid := session.Appid
	if appid == 0 {
		lwutil.SendError("err_auth", "Please login with app secret")
	}

	nowUnix := time.Now().Unix()

	//get matchIds
	key := fmt.Sprintf("matchesInApp/%d", appid)
	matchIdValues, err := redis.Values(rc.Do("zrangebyscore", key, nowUnix, "+inf"))
	lwutil.CheckError(err, "")

	args := make([]interface{}, len(matchIdValues)+1)
	args[0] = "matches"
	for i, v := range matchIdValues {
		var id int
		id, err := redis.Int(v, err)
		lwutil.CheckError(err, "")
		matchkey := fmt.Sprintf("%d+%d", appid, id)
		args[i+1] = matchkey
	}

	//get match data
	matchesValues, err := redis.Values(rc.Do("hmget", args...))

	matches := make([]Match, len(matchesValues))

	for i, v := range matchesValues {
		var match Match
		err = json.Unmarshal(v.([]byte), &match)
		lwutil.CheckError(err, "")
		matches[i] = match
	}

	//out
	type OutMatch struct {
		Id       uint32
		Name     string
		GameId   uint32
		Begin    int64
		End      int64
		Sort     string
		TryMax   uint32
		TryPrice uint32
		TryNum   uint32
	}

	outMatches := make([]OutMatch, len(matches))

	// get try number
	for _, match := range matches {
		tryNumKey := makeTryNumKey(match.Id)
		rc.Send("hget", tryNumKey, session.Userid)
	}
	err = rc.Flush()
	lwutil.CheckError(err, "")

	for i, match := range matches {
		tryNum, err := redis.Int(rc.Receive())
		if err != nil && err != redis.ErrNil {
			lwutil.CheckError(err, "")
		}

		outMatches[i] = OutMatch{
			Id:       match.Id,
			Name:     match.Name,
			GameId:   match.GameId,
			Begin:    match.Begin,
			End:      match.End,
			Sort:     match.Sort,
			TryMax:   match.TryMax,
			TryPrice: match.TryPrice,
			TryNum:   uint32(tryNum),
		}
	}

	lwutil.WriteResponse(w, outMatches)
}

func startMatch(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "POST")

	rc := redisPool.Get()
	defer rc.Close()

	session, err := findSession(w, r, rc)
	lwutil.CheckError(err, "err_auth")

	appid := session.Appid
	if appid == 0 {
		lwutil.SendError("err_auth", "Please login with app secret")
	}

	//input
	var in struct {
		MatchId uint32
	}
	err = lwutil.DecodeRequestBody(r, &in)
	lwutil.CheckError(err, "err_decode_body")

	//get match info
	key := fmt.Sprintf("%d+%d", appid, in.MatchId)
	matchJson, err := redis.Bytes(rc.Do("hget", "matches", key))
	lwutil.CheckError(err, "err_not_found")

	var match Match
	err = json.Unmarshal(matchJson, &match)
	lwutil.CheckError(err, "")

	//check time
	now := time.Now().Unix()
	if now < match.Begin || now >= match.End-MATCH_TRY_DURATION_SEC {
		lwutil.SendError("err_time", "now < match.Begin || now >= match.End-MATCH_TRY_DURATION_SEC")
	}

	//incr and check try number
	tryNumKey := makeTryNumKey(in.MatchId)
	tryNum, err := redis.Int(rc.Do("hget", tryNumKey, session.Userid))
	if err != nil && err != redis.ErrNil {
		lwutil.CheckError(err, "")
	}
	if uint32(tryNum) >= match.TryMax {
		lwutil.SendError("err_no_try", "no try left")
	}
	_, err = rc.Do("hincrby", tryNumKey, session.Userid, 1)
	lwutil.CheckError(err, "")
	tryNum++

	//new try secret
	trySecret := lwutil.GenUUID()
	_, err = rc.Do("setex", fmt.Sprintf("trySecrets/%s", trySecret), MATCH_TRY_DURATION_SEC, in.MatchId)
	lwutil.CheckError(err, "")

	//out
	out := struct {
		Secret string
		TryNum uint32
	}{trySecret, uint32(tryNum)}
	lwutil.WriteResponse(w, out)
}

func addScore(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "POST")

	rc := redisPool.Get()
	defer rc.Close()

	session, err := findSession(w, r, rc)
	lwutil.CheckError(err, "err_auth")

	appid := session.Appid
	if appid == 0 {
		lwutil.SendError("err_auth", "Please login with app secret")
	}

	//input
	var in struct {
		TrySecret string
		Score     int64
	}
	err = lwutil.DecodeRequestBody(r, &in)
	lwutil.CheckError(err, "err_decode_body")

	//use secret to get matchId
	matchIdRaw, err := rc.Do("get", fmt.Sprintf("trySecrets/%s", in.TrySecret))
	lwutil.CheckError(err, "")
	if matchIdRaw == nil {
		lwutil.SendError("err_secret", "")
	}
	matchId64, err := redis.Int64(matchIdRaw, err)
	lwutil.CheckError(err, "")
	matchId := uint32(matchId64)

	//get match info and prev score
	keyLeaderboard := makeLeaderboardKey(matchId)
	matchKey := fmt.Sprintf("%d+%d", appid, matchId)
	rc.Send("hget", "matches", matchKey)
	rc.Send("zscore", keyLeaderboard, session.Userid)
	rc.Flush()
	matchJs, err := redis.Bytes(rc.Receive())
	lwutil.CheckError(err, "")

	var match Match
	err = json.Unmarshal(matchJs, &match)
	lwutil.CheckError(err, "")

	prevScore, err := redis.Int64(rc.Receive())
	needOverwrite := false
	if err == redis.ErrNil {
		needOverwrite = true
	} else {
		lwutil.CheckError(err, "")
		if match.Sort == SORT_ASC {
			if in.Score < prevScore {
				needOverwrite = true
			}
		} else if match.Sort == SORT_DESC {
			if in.Score > prevScore {
				needOverwrite = true
			}
		} else {
			lwutil.SendError("", "invalid match.Sort: "+match.Sort)
		}
	}

	//del from failboard and add to leaderboard and delete secret
	if needOverwrite {
		rc.Send("zadd", keyLeaderboard, in.Score, session.Userid)
	}
	rc.Send("zrank", keyLeaderboard, session.Userid)
	rc.Send("del", fmt.Sprintf("trySecrets/%s", in.TrySecret))

	err = rc.Flush()
	lwutil.CheckError(err, "")

	if needOverwrite {
		_, err := rc.Receive()
		lwutil.CheckError(err, "")
	}
	rank, err := redis.Int(rc.Receive())
	lwutil.CheckError(err, "")
	rank++

	//reply
	lwutil.WriteResponse(w, rank)
}

func makeLeaderboardKey(matchId uint32) string {
	return fmt.Sprintf("leaderboard/%d", matchId)
}

func makeTryNumKey(matchId uint32) string {
	return fmt.Sprintf("trynum/%d", matchId)
}

func regMatch() {
	http.Handle("/match/new", lwutil.ReqHandler(newMatch))
	http.Handle("/match/del", lwutil.ReqHandler(delMatch))
	http.Handle("/match/list", lwutil.ReqHandler(listMatch))
	http.Handle("/match/start", lwutil.ReqHandler(startMatch))
	http.Handle("/match/addscore", lwutil.ReqHandler(addScore))
}

/*
matches: Hash{field: [appid+matchId] int, value: match Match.json}
matchesInApp: SortedSet{score: matchEndTimeUnix, member: matchId}
trySecrets/<trySecret>: String(matchId int)
trynum/<matchId int>: Hash{field: userId int, value: tryNum int}
leaderboard/<matchId int>: SortedSet{score: matchScore int, member: userId int}
*/

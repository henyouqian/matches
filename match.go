package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"time"
)

type Match struct {
	Id        uint32
	Name      string
	GameId    uint32
	Begin     int64
	End       int64
	Sort      string
	TimeLimit uint32
}

func newMatch(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err, "err_auth")
	checkAdmin(session)

	appid := session.Appid
	if appid == 0 {
		sendError("err_auth", "Please login with app secret")
	}

	// input
	type Input struct {
		Name      string
		GameId    uint32
		Begin     string
		End       string
		Sort      string
		TimeLimit uint32
	}
	input := Input{}
	decodeRequestBody(r, &input)

	if input.Name == "" || input.Begin == "" || input.End == "" || input.GameId == 0 {
		sendError("err_input", "Missing Name || Begin || End || Gameid")
	}
	if input.Sort != "ASC" && input.Sort != "DESC" {
		sendError("err_input", "Invalid Sort, must be ASC or DESC")
	}
	if input.TimeLimit < 60 {
		sendError("err_input", "Time limit must > 60")
	}

	// times
	const timeform = "2006-01-02 15:04:05"
	begin, err := time.ParseInLocation(timeform, input.Begin, time.Local)
	checkError(err, "err_shit")
	end, err := time.ParseInLocation(timeform, input.End, time.Local)
	checkError(err, "")
	beginUnix := begin.Unix()
	endUnix := end.Unix()

	if endUnix-beginUnix <= 60 {
		sendError("err_input", "endUnix - beginUnix must > 60 seconds")
	}
	if time.Now().Unix() > endUnix {
		sendError("err_input", "end time before now")
	}

	//
	rc := redisPool.Get()
	defer rc.Close()

	matchId, err := redis.Int(rc.Do("incr", "matchIdAutoIncr"))
	checkError(err, "")

	match := Match{
		uint32(matchId),
		input.Name,
		input.GameId,
		beginUnix,
		endUnix,
		input.Sort,
		input.TimeLimit,
	}

	matchJson, err := json.Marshal(match)
	checkError(err, "")

	key := fmt.Sprintf("matches/%d+%d", appid, matchId)
	rc.Send("set", key, matchJson)
	key = fmt.Sprintf("matchesInApp/%d", appid)
	rc.Send("zadd", key, endUnix, matchId)
	rc.Flush()
	for i := 0; i < 2; i++ {
		_, err = rc.Receive()
		checkError(err, "")
	}

	// reply
	writeResponse(w, match)
}

func delMatch(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err, "err_auth")
	checkAdmin(session)

	appid := session.Appid
	if appid == 0 {
		sendError("err_auth", "Please login with app secret")
	}

	// input
	matchIds := make([]int, 0, 8)
	decodeRequestBody(r, &matchIds)

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

	keys := make([]interface{}, 0, 8)
	for _, matchId := range matchIds {
		key = fmt.Sprintf("matches/%d+%d", appid, matchId)
		keys = append(keys, key)
	}
	rc.Send("del", keys...)
	rc.Flush()

	_, err = rc.Receive()
	checkError(err, "")
	delNum, err := rc.Receive()
	checkError(err, "")

	// reply
	writeResponse(w, delNum)
}

func listMatch(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err, "err_auth")

	appid := session.Appid
	if appid == 0 {
		sendError("err_auth", "Please login with app secret")
	}

	nowUnix := time.Now().Unix()

	rc := redisPool.Get()
	defer rc.Close()

	// get matchIds
	key := fmt.Sprintf("matchesInApp/%d", appid)
	matchIdValues, err := redis.Values(rc.Do("zrangebyscore", key, nowUnix, "+inf"))
	checkError(err, "")

	matchKeys := make([]interface{}, 0, 10)
	for _, v := range matchIdValues {
		var id int
		id, err := redis.Int(v, err)
		checkError(err, "")
		matchkey := fmt.Sprintf("matches/%d+%d", appid, id)
		matchKeys = append(matchKeys, matchkey)
	}

	// get match data
	matchesValues, err := redis.Values(rc.Do("mget", matchKeys...))

	matches := make([]interface{}, 0, 10)
	for _, v := range matchesValues {
		var match interface{}
		err = json.Unmarshal(v.([]byte), &match)
		checkError(err, "")
		matches = append(matches, match)
	}

	writeResponse(w, matches)
}

func startMatch(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err, "err_auth")

	appid := session.Appid
	if appid == 0 {
		sendError("err_auth", "Please login with app secret")
	}

	writeResponse(w, 1)
}

func regMatch() {
	http.HandleFunc("/match/new", newMatch)
	http.HandleFunc("/match/del", delMatch)
	http.HandleFunc("/match/list", listMatch)
	http.HandleFunc("/match/start", startMatch)
}

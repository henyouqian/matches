package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
)

type Game struct {
	Id   uint32
	Name string
	Sort string
}

func newGame(w http.ResponseWriter, r *http.Request) {
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
		Name string
		Sort string
	}
	input := Input{}
	decodeRequestBody(r, &input)

	if input.Name == "" {
		sendError("err_input", "Missing Name")
	}
	if input.Sort != "ASC" && input.Sort != "DESC" {
		sendError("err_input", "Invalid Sort, must be ASC or DESC")
	}

	//
	rc := redisPool.Get()
	defer rc.Close()

	gameId, err := redis.Int(rc.Do("incr", "idGen/game"))
	checkError(err, "")

	game := Game{
		uint32(gameId),
		input.Name,
		input.Sort,
	}

	gameJson, err := json.Marshal(game)
	checkError(err, "")

	key := fmt.Sprintf("games/%d", appid)
	_, err = rc.Do("hset", key, gameId, gameJson)
	checkError(err, "")

	// reply
	writeResponse(w, game)
}

//func delGame(w http.ResponseWriter, r *http.Request) {
//	defer handleError(w)
//	checkMathod(r, "POST")

//	session, err := findSession(w, r)
//	checkError(err, "err_auth")
//	checkAdmin(session)

//	appid := session.Appid
//	if appid == 0 {
//		sendError("err_auth", "Please login with app secret")
//	}

//	// input
//	matchIds := make([]int, 0, 8)
//	decodeRequestBody(r, &matchIds)

//	// redis
//	rc := redisPool.Get()
//	defer rc.Close()

//	key := fmt.Sprintf("matchesInApp/%d", appid)
//	params := make([]interface{}, 0, 8)
//	params = append(params, key)
//	matchIdsItf := make([]interface{}, len(matchIds))
//	for i, v := range matchIds {
//		matchIdsItf[i] = v
//	}
//	params = append(params, matchIdsItf...)
//	rc.Send("zrem", params...)

//	keys := make([]interface{}, 0, 8)
//	for _, matchId := range matchIds {
//		key = fmt.Sprintf("matches/%d+%d", appid, matchId)
//		keys = append(keys, key)
//	}
//	rc.Send("del", keys...)
//	rc.Flush()

//	_, err = rc.Receive()
//	checkError(err, "")
//	delNum, err := rc.Receive()
//	checkError(err, "")

//	// reply
//	writeResponse(w, delNum)
//}

//func listMatch(w http.ResponseWriter, r *http.Request) {
//	defer handleError(w)
//	checkMathod(r, "POST")

//	session, err := findSession(w, r)
//	checkError(err, "err_auth")

//	appid := session.Appid
//	if appid == 0 {
//		sendError("err_auth", "Please login with app secret")
//	}

//	nowUnix := time.Now().Unix()

//	rc := redisPool.Get()
//	defer rc.Close()

//	// get matchIds
//	key := fmt.Sprintf("matchesInApp/%d", appid)
//	matchIdValues, err := redis.Values(rc.Do("zrangebyscore", key, nowUnix, "+inf"))
//	checkError(err, "")

//	matchKeys := make([]interface{}, 0, 10)
//	for _, v := range matchIdValues {
//		var id int
//		id, err := redis.Int(v, err)
//		checkError(err, "")
//		matchkey := fmt.Sprintf("matches/%d+%d", appid, id)
//		matchKeys = append(matchKeys, matchkey)
//	}

//	// get match data
//	matchesValues, err := redis.Values(rc.Do("mget", matchKeys...))

//	matches := make([]interface{}, 0, 10)
//	for _, v := range matchesValues {
//		var match interface{}
//		err = json.Unmarshal(v.([]byte), &match)
//		checkError(err, "")
//		matches = append(matches, match)
//	}

//	writeResponse(w, matches)
//}

func regGame() {
	http.HandleFunc("/game/new", newGame)
	//http.HandleFunc("/game/del", delGame)
	//http.HandleFunc("/game/list", listGame)
}

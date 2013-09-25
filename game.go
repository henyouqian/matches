package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/henyouqian/lwutil"
	"net/http"
)

type Game struct {
	Id   uint32
	Name string
	Sort string
}

func findGame(gameid, appid uint32) (*Game, error) {
	rc := redisPool.Get()
	defer rc.Close()

	key := fmt.Sprintf("games/%d", appid)
	gameJson, err := redis.Bytes(rc.Do("hget", key, gameid))
	if err != nil {
		return nil, err
	}

	var game Game
	err = json.Unmarshal(gameJson, &game)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func newGame(w http.ResponseWriter, r *http.Request) {
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

	// input
	input := struct {
		Id   uint32
		Name string
		Sort string
	}{}
	err = lwutil.DecodeRequestBody(r, &input)
	lwutil.CheckError(err, "err_decode_body")

	if input.Id == 0 || input.Name == "" {
		lwutil.SendError("err_input", "Missing Id or Name")
	}
	if input.Sort != SORT_ASC && input.Sort != SORT_DESC {
		lwutil.SendError("err_input", "Invalid Sort, must be ASC or DESC")
	}

	//
	game := Game{
		input.Id,
		input.Name,
		input.Sort,
	}

	gameJson, err := json.Marshal(game)
	lwutil.CheckError(err, "")

	key := fmt.Sprintf("games/%d", appid)
	_, err = rc.Do("hset", key, input.Id, gameJson)
	lwutil.CheckError(err, "")

	// reply
	lwutil.WriteResponse(w, game)
}

func delGame(w http.ResponseWriter, r *http.Request) {
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

	// input
	gameIds := make([]int, 0, 8)
	err = lwutil.DecodeRequestBody(r, &gameIds)
	lwutil.CheckError(err, "err_decode_body")

	// redis
	args := make([]interface{}, 1, 8)
	args[0] = fmt.Sprintf("games/%d", appid)
	for _, gameId := range gameIds {
		args = append(args, gameId)
	}

	delNum, err := redis.Int(rc.Do("hdel", args...))
	lwutil.CheckError(err, "")

	// reply
	lwutil.WriteResponse(w, delNum)
}

func listGame(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "POST")

	rc := redisPool.Get()
	defer rc.Close()

	session, err := findSession(w, r, rc)
	lwutil.CheckError(err, "err_auth")

	appid := session.Appid
	if appid == 0 {
		lwutil.SendError("err_auth", "Please login with app secret")
	}

	// get game data
	gameValues, err := redis.Values(rc.Do("hgetall", fmt.Sprintf("games/%d", appid)))
	lwutil.CheckError(err, "")

	games := make([]interface{}, 0, len(gameValues)/2)
	for i, v := range gameValues {
		if i%2 == 0 {
			continue
		}
		var game interface{}
		err = json.Unmarshal(v.([]byte), &game)
		lwutil.CheckError(err, "")
		games = append(games, game)
	}

	//reply
	lwutil.WriteResponse(w, games)
}

func regGame() {
	http.Handle("/game/new", lwutil.ReqHandler(newGame))
	http.Handle("/game/del", lwutil.ReqHandler(delGame))
	http.Handle("/game/list", lwutil.ReqHandler(listGame))
}

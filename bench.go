package main

import (
	"net/http"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

func benchLogin(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "GET")

	// input
	type Input struct {
		Username  string
		Password  string
		Appsecret string
	}
	input := Input{Username:"aa", Password:"aa"}

	if input.Username == "" || input.Password == "" {
		sendError("err_input", "")
	}

	pwsha := sha224(input.Password + passwordSalt)

	// get userid
	row := authDB.QueryRow("SELECT id FROM user_accounts WHERE username=? AND password=?", input.Username, pwsha)
	var userid uint64
	err := row.Scan(&userid)
	checkError(err, "")

	// get appid
	appid := uint32(0)
	if input.Appsecret != "" {
		row = authDB.QueryRow("SELECT id FROM apps WHERE secret=?", input.Appsecret)
		err = row.Scan(&appid)
		checkError(err, "")
	}

	// new session
	rc := redisPool.Get()
	defer rc.Close()

	usertoken, err := newSession(w, rc, userid, input.Username, appid)
	checkError(err, "")

	// reply
	type Reply struct {
		Usertoken string
		Appid     uint32
	}
	reply := Reply{usertoken, appid}
	writeResponse(w, reply)
}


func benchHello(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, "hello")
}

func benchDBSingleSelect(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "GET")

	row := authDB.QueryRow("SELECT id FROM user_accounts WHERE username=?", "admin")
	var userid uint32
	err := row.Scan(&userid)
	checkError(err, "")

	writeResponse(w, userid)
}

func benchRedisGet(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "GET")

	rc := redisPool.Get()
	defer rc.Close()

	rc.Send("set", "foo", "yes")
	rc.Send("get", "foo")
	rc.Flush()
	rc.Receive()
	foo, err := redis.String(rc.Receive())
	checkError(err, "")

	// foo, err := redis.String(rc.Do("get", "foo"))
	// checkError(err, "")

	writeResponse(w, foo)
}

func benchJson(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "GET")

	str := []byte(`
		{
			"Name": "aa",
			"Gameid": 885,
			"Begin": "2013/04/06 23:43:24",
			"End": "2013/05/06 23:43:24",
			"Sort": 0
		}
	`)

	type Input struct {
		Name   string
		Gameid uint32
		Begin  string
		End    string
		Sort   uint8
	}
	input := Input{}
	err := json.Unmarshal(str, &input)
	checkError(err, "")

	writeResponse(w, input)
}

func benchJson2(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "GET")

	str := []byte(`
		{
			"Name": "aa"
		}
	`)

	type Input struct {
		Name   string
	}
	input := Input{}
	err := json.Unmarshal(str, &input)
	checkError(err, "")

	writeResponse(w, input)
}

func regBench() {
	http.HandleFunc("/bench/login", benchLogin)
	http.HandleFunc("/bench/hello", benchHello)
	http.HandleFunc("/bench/dbsingleselect", benchDBSingleSelect)
	http.HandleFunc("/bench/redisget", benchRedisGet)
	http.HandleFunc("/bench/json", benchJson)
	http.HandleFunc("/bench/json2", benchJson2)
}

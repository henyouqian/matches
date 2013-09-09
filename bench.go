package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"net/http"
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
	input := Input{Username: "admin", Password: "admin"}

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

const insertCount = 10

func benchDBInsert(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "GET")

	//db
	stmt, err := matchDB.Prepare("INSERT INTO insertTest (a, b, c, d) VALUES (?, ?, ?, ?)")
	checkError(err, "")

	ids := make([]int64, insertCount)
	for i := 0; i < insertCount; i++ {
		res, err := stmt.Exec(1, 2, 3, 4)
		checkError(err, "err_account_exists")

		ids[i], err = res.LastInsertId()
		checkError(err, "")
	}

	writeResponse(w, ids)
}

func benchDBInsertTx(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "GET")

	//db
	tx, err := matchDB.Begin()
	defer endTx(tx, &err)

	checkError(err, "")
	stmt, err := tx.Prepare("INSERT INTO insertTest (a, b, c, d) VALUES (?, ?, ?, ?)")
	checkError(err, "")

	ids := make([]int64, insertCount)
	for i := 0; i < insertCount; i++ {
		res, err := stmt.Exec(1, 2, 3, 4)
		checkError(err, "err_account_exists")

		ids[i], err = res.LastInsertId()
		checkError(err, "")
	}

	writeResponse(w, ids)
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
		Name string
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
	http.HandleFunc("/bench/dbinsert", benchDBInsert)
	http.HandleFunc("/bench/dbinserttx", benchDBInsertTx)
	http.HandleFunc("/bench/redisget", benchRedisGet)
	http.HandleFunc("/bench/json", benchJson)
	http.HandleFunc("/bench/json2", benchJson2)
}

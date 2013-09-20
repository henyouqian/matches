package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	//"github.com/golang/glog"
	"github.com/henyouqian/lwutil"
	"net/http"
)

func benchLogin(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "GET")

	// input
	in := struct {
		Username  string
		Password  string
		Appsecret string
	}{Username: "admin", Password: "admin"}

	if in.Username == "" || in.Password == "" {
		lwutil.SendError("err_input", "")
	}

	pwsha := lwutil.Sha224(in.Password + passwordSalt)

	// get userid
	row := authDB.QueryRow("SELECT id FROM user_accounts WHERE username=? AND password=?", in.Username, pwsha)
	var userid uint64
	err := row.Scan(&userid)
	lwutil.CheckError(err, "")

	// get appid
	appid := uint32(0)
	if in.Appsecret != "" {
		row = authDB.QueryRow("SELECT id FROM apps WHERE secret=?", in.Appsecret)
		err = row.Scan(&appid)
		lwutil.CheckError(err, "")
	}

	// new session
	rc := redisPool.Get()
	defer rc.Close()

	usertoken, err := newSession(w, userid, in.Username, appid, rc)
	lwutil.CheckError(err, "")

	// reply
	reply := struct {
		Usertoken string
		Appid     uint32
	}{usertoken, appid}
	lwutil.WriteResponse(w, reply)
}

func benchHello(w http.ResponseWriter, r *http.Request) {
	lwutil.WriteResponse(w, "hello")
}

func benchDBSingleSelect(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "GET")

	row := authDB.QueryRow("SELECT id FROM user_accounts WHERE username=?", "admin")
	var userid uint32
	err := row.Scan(&userid)
	lwutil.CheckError(err, "")

	lwutil.WriteResponse(w, userid)
}

const insertCount = 10

func benchDBInsert(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "GET")

	//db
	stmt, err := matchDB.Prepare("INSERT INTO insertTest (a, b, c, d) VALUES (?, ?, ?, ?)")
	lwutil.CheckError(err, "")

	ids := make([]int64, insertCount)
	for i := 0; i < insertCount; i++ {
		res, err := stmt.Exec(1, 2, 3, 4)
		lwutil.CheckError(err, "err_account_exists")

		ids[i], err = res.LastInsertId()
		lwutil.CheckError(err, "")
	}

	lwutil.WriteResponse(w, ids)
}

func benchDBInsertTx(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "GET")

	//db
	tx, err := matchDB.Begin()
	defer lwutil.EndTx(tx, &err)

	lwutil.CheckError(err, "")
	stmt, err := tx.Prepare("INSERT INTO insertTest (a, b, c, d) VALUES (?, ?, ?, ?)")
	lwutil.CheckError(err, "")

	ids := make([]int64, insertCount)
	for i := 0; i < insertCount; i++ {
		res, err := stmt.Exec(1, 2, 3, 4)
		lwutil.CheckError(err, "err_account_exists")

		ids[i], err = res.LastInsertId()
		lwutil.CheckError(err, "")
	}

	lwutil.WriteResponse(w, ids)
}

func benchRedisGet(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "GET")

	rc := redisPool.Get()
	defer rc.Close()

	rc.Send("set", "foo", "yes")
	rc.Send("get", "foo")
	rc.Flush()
	rc.Receive()
	foo, err := redis.String(rc.Receive())
	lwutil.CheckError(err, "")

	// foo, err := redis.String(rc.Do("get", "foo"))
	// lwutil.CheckError(err, "")

	lwutil.WriteResponse(w, foo)
}

func benchJson(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "GET")

	str := []byte(`
		{
			"Name": "aa",
			"Gameid": 885,
			"Begin": "2013/04/06 23:43:24",
			"End": "2013/05/06 23:43:24",
			"Sort": 0
		}
	`)

	in := struct {
		Name   string
		Gameid uint32
		Begin  string
		End    string
		Sort   uint8
	}{}
	err := json.Unmarshal(str, &in)
	lwutil.CheckError(err, "")

	lwutil.WriteResponse(w, in)
}

func benchJson2(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "GET")

	str := []byte(`
		{
			"Name": "aa"
		}
	`)

	in := struct {
		Name string
	}{}
	err := json.Unmarshal(str, &in)
	lwutil.CheckError(err, "")

	lwutil.WriteResponse(w, in)
}

func regBench() {
	http.Handle("/bench/login", lwutil.ReqHandler(benchLogin))
	http.Handle("/bench/hello", lwutil.ReqHandler(benchHello))
	http.Handle("/bench/dbsingleselect", lwutil.ReqHandler(benchDBSingleSelect))
	http.Handle("/bench/dbinsert", lwutil.ReqHandler(benchDBInsert))
	http.Handle("/bench/dbinserttx", lwutil.ReqHandler(benchDBInsertTx))
	http.Handle("/bench/redisget", lwutil.ReqHandler(benchRedisGet))
	http.Handle("/bench/json", lwutil.ReqHandler(benchJson))
	http.Handle("/bench/json2", lwutil.ReqHandler(benchJson2))
}

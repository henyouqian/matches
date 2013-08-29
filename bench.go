package main

import (
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

func testdb(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)

	rows, err := authDB.Query("SELECT id FROM user_accounts")
	checkError(err, "")

	ids := make([]uint32, 0, 50)
	for rows.Next() {
		var userid uint32
		err = rows.Scan(&userid)
		checkError(err, "")
		ids = append(ids, userid)
	}

	rows, err = authDB.Query("SELECT id FROM user_accounts")
	checkError(err, "")

	// rows, err = auth_db.Query("SELECT id FROM user_accounts")
	// checkError(err, "")

	for rows.Next() {
		var userid uint32
		err = rows.Scan(&userid)
		checkError(err, "")
		ids = append(ids, userid)
	}

	writeResponse(w, ids)
}

func regBench() {
	http.HandleFunc("/bench/login", benchLogin)
	http.HandleFunc("/bench/hello", benchHello)
	http.HandleFunc("/bench/testdb", testdb)
}

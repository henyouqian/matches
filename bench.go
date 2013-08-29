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
	db := opendb("auth_db")
	defer db.Close()

	rows, err := db.Query("SELECT id FROM user_accounts WHERE username=? AND password=?", input.Username, pwsha)
	checkError(err, "")
	if rows.Next() == false {
		sendError("err_not_match", "")
	}
	var userid uint64
	err = rows.Scan(&userid)
	checkError(err, "")

	// get appid
	appid := uint32(0)
	if input.Appsecret != "" {
		rows, err = db.Query("SELECT id FROM apps WHERE secret=?", input.Appsecret)
		checkError(err, "")
		if rows.Next() == false {
			sendError("err_app_secret", "")
		}
		err = rows.Scan(&appid)
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
	rows, err := auth_db.Query("SELECT id FROM user_accounts")
	checkError(err, "")

	ids := make([]uint32, 0, 50)
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

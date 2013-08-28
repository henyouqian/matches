package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"time"
)

const passwordSalt = "liwei"
const sessionLifeSecond = 60 * 60 * 24 * 7
const sessionUpdateSecond = 60 * 60

type Session struct {
	Userid   uint64
	Username string
	Born     time.Time
	Appid    uint32
}

func init() {
	
}

func newSession(w http.ResponseWriter, rc redis.Conn, userid uint64, username string, appid uint32) (usertoken string, err error) {
	usertoken = ""
	usertokenRaw, err := rc.Do("get", fmt.Sprintf("userid+appid:usertoken/%d+%d", userid, appid))
	checkError(err, "")
	if usertokenRaw != nil {
		usertoken, err := redis.String(usertokenRaw, err)
		if err != nil {
			return usertoken, err
		}
		rc.Do("del", fmt.Sprintf("usertoken:session/%s", usertoken))
	}

	usertoken = genUUID()

	session := Session{userid, username, time.Now(), appid}
	jsonSession, err := json.Marshal(session)
	if err != nil {
		return usertoken, err
	}

	_, err = rc.Do("setex", fmt.Sprintf("usertoken:session/%s", usertoken), sessionLifeSecond, jsonSession)
	if err != nil {
		return usertoken, err
	}
	_, err = rc.Do("set", fmt.Sprintf("userid+appid:usertoken/%d+%d", userid, appid), usertoken)
	if err != nil {
		return usertoken, err
	}

	// cookie
	http.SetCookie(w, &http.Cookie{Name: "usertoken", Value: usertoken, MaxAge: sessionLifeSecond, Path: "/"})

	return usertoken, err
}

func findSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	session := new(Session)

	usertokenCookie, err := r.Cookie("usertoken")
	if err != nil {
		return session, errors.New("err_auth")
	}
	usertoken := usertokenCookie.Value

	//redis
	rc := redisPool.Get()
	defer rc.Close()

	sessionBytes, err := redis.Bytes(rc.Do("get", fmt.Sprintf("usertoken:session/%s", usertoken)))
	if err != nil {
		err = errors.New("err_auth")
		return session, errors.New("err_auth")
	}

	err = json.Unmarshal(sessionBytes, &session)
	checkError(err, "")

	//update session
	dt := time.Now().Sub(session.Born)
	if dt > sessionUpdateSecond*time.Second {
		newSession(w, rc, session.Userid, session.Username, session.Appid)
	}

	return session, nil
}

func register(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	// input
	type Input struct {
		Username string
		Password string
	}
	var input Input
	decodeRequestBody(r, &input)

	if input.Username == "" || input.Password == "" {
		sendError("err_input", "")
	}

	pwsha := sha224(input.Password + passwordSalt)

	// insert into db
	db := opendb("auth_db")
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO user_accounts (username, password) VALUES (?, ?)")
	checkError(err, "")

	res, err := stmt.Exec(input.Username, pwsha)
	checkError(err, "err_account_exists")

	id, err := res.LastInsertId()
	checkError(err, "")

	// reply
	type Reply struct {
		Userid int64
	}
	writeResponse(w, Reply{id})
}

func login(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	// input
	type Input struct {
		Username  string
		Password  string
		Appsecret string
	}
	var input Input
	decodeRequestBody(r, &input)

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

func newApp(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err, "")
	checkAdmin(session)

	// input
	type Input struct {
		Name string
	}
	var input Input
	decodeRequestBody(r, &input)

	if input.Name == "" {
		sendError("err_input", "input.Name empty")
	}

	// db
	db := opendb("auth_db")
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO apps (name, secret) VALUES (?, ?)")
	checkError(err, "")

	secret := genUUID()
	_, err = stmt.Exec(input.Name, secret)
	checkError(err, "err_name_exists")

	// reply
	type Reply struct {
		Name string
		Secret string
	}
	reply := Reply{input.Name, secret}
	writeResponse(w, reply)
}

func test(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err, "")

	writeResponse(w, session)
}

func regAuth() {
	http.HandleFunc("/auth/login", login)
	http.HandleFunc("/auth/register", register)
	http.HandleFunc("/auth/newapp", newApp)
	http.HandleFunc("/auth/test", test)
}

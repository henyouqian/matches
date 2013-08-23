package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/nu7hatch/gouuid"
	"net/http"
	"time"
	// "database/sql"
	// _ "github.com/go-sql-driver/mysql"
	// "time"
)

const passwordSalt = "liwei"
const sessionLifeSecond = 60 * 60 * 24 * 7
const sessionUpdateSecond = 60 * 60

type Session struct {
	Userid   int64
	Username string
	Born     time.Time
}

func init() {
	type Match struct {
		Begin time.Time
		End   time.Time
	}

}

func newSession(w http.ResponseWriter, rc redis.Conn, userid int64, username string) (usertoken string, err error) {
	usertoken = ""
	usertokenRaw, err := rc.Do("get", fmt.Sprintf("userid:usertoken/%d", userid))
	checkError(err)
	if usertokenRaw != nil {
		usertoken, err := redis.String(usertokenRaw, err)
		if err != nil {
			return usertoken, err
		}
		rc.Do("del", fmt.Sprintf("usertoken:session/%s", usertoken))
	}

	uuid, err := uuid.NewV4()
	if err != nil {
		return usertoken, err
	}
	usertoken = uuid.String()

	session := Session{userid, username, time.Now()}
	jsonSession, err := json.Marshal(session)
	if err != nil {
		return usertoken, err
	}

	_, err = rc.Do("setex", fmt.Sprintf("usertoken:session/%s", usertoken), sessionLifeSecond, jsonSession)
	if err != nil {
		return usertoken, err
	}
	_, err = rc.Do("set", fmt.Sprintf("userid:usertoken/%d", userid), usertoken)
	if err != nil {
		return usertoken, err
	}

	// cookie
	http.SetCookie(w, &http.Cookie{Name: "usertoken", Value: usertoken, MaxAge: sessionLifeSecond})

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
	checkError(err)

	//update session
	dt := time.Now().Sub(session.Born)
	if dt > sessionUpdateSecond*time.Second {
		newSession(w, rc, session.Userid, session.Username)
	}

	return session, nil
}

func register(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	// params
	type regParam struct {
		Username string
		Password string
	}
	var param regParam
	err := decodeRequestBody(r, &param)
	checkError(err)

	if param.Username == "" || param.Password == "" {
		panic("err_param")
	}

	pwsha := sha224(param.Password + passwordSalt)

	// insert into db
	db := opendb("account_db")
	defer db.Close()

	stmt, err := db.Prepare("INSERT user_account SET username=?,password=?")
	checkError(err)

	res, err := stmt.Exec(param.Username, pwsha)
	if err != nil {
		panic("err_account_exists")
	}

	id, err := res.LastInsertId()
	checkError(err)

	// reply
	type Reply struct {
		Userid int64
	}
	writeResponse(w, Reply{id})
}

func login(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	// params
	type regParam struct {
		Username string
		Password string
	}
	var param regParam
	err := decodeRequestBody(r, &param)
	checkError(err)

	if param.Username == "" || param.Password == "" {
		panic("err_param")
	}

	pwsha := sha224(param.Password + passwordSalt)

	// validate
	db := opendb("account_db")
	defer db.Close()

	rows, err := db.Query("SELECT id FROM user_account WHERE username=? AND password=?", param.Username, pwsha)
	checkError(err)
	if rows.Next() == false {
		panic("err_not_match")
	}

	var userid int64
	err = rows.Scan(&userid)
	checkError(err)

	// create session
	rc := redisPool.Get()
	defer rc.Close()

	usertoken, err := newSession(w, rc, userid, param.Username)
	checkError(err)

	// reply
	type Reply struct {
		Usertoken string
	}
	reply := Reply{usertoken}
	writeResponse(w, reply)
}

func test(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err)

	writeResponse(w, session)
}

func regAuth() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/test", test)
}

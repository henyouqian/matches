package main

import (
	"fmt"
	"net/http"
	"github.com/nu7hatch/gouuid"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"errors"
	"time"
	// "database/sql"
	// _ "github.com/go-sql-driver/mysql"
	// "time"
)

const passwordSalt = "liwei"
const sessionLifeSecond = 60 * 60 * 24 * 7


type Session struct {
	Userid int64
	Username string
}


func init() {
	type Session struct {
		Userid int64
		Username string
		Time time.Time
	}
	// db, err := sql.Open("mysql", "root@/wh_db?parseTime=true")
	// if err != nil {
	// 	panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	// }
	// defer db.Close()

	// // Prepare statement for reading data
	// rows, err := db.Query("SELECT id, time FROM wagons")
	// if err != nil {
	// 	panic(err.Error()) // proper error handling instead of panic in your app
	// }

	// var id int
	// var time time.Time

	// for rows.Next() {
	// 	err = rows.Scan(&id, &time) // WHERE number = 13
	// 	if err != nil {
	// 		panic(err.Error()) // proper error handling instead of panic in your app
	// 	}
	// 	fmt.Println(id, time)
	// }
}


func findSession(r *http.Request) (*Session, error){
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

	return session, nil
}


func register(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	// params
	type regParam struct{
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
	type Reply struct{
		Userid int64
	}
	writeResponse(w, Reply{id})
}


func login(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")
	
	// params
	type regParam struct{
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

	usertokenRaw, err := rc.Do("get", fmt.Sprintf("userid:usertoken/%d", userid))
	checkError(err)
	if usertokenRaw != nil {
		usertoken, err := redis.String(usertokenRaw, err)
		checkError(err)
		rc.Do("del", fmt.Sprintf("usertoken:session/%s", usertoken))
	}

	uuid, err := uuid.NewV4()
	checkError(err)
	usertoken := uuid.String()

	session := Session{userid, param.Username}
	jsonSession, err := json.Marshal(session)
	checkError(err)

	_, err = rc.Do("setex", fmt.Sprintf("usertoken:session/%s", usertoken), sessionLifeSecond, jsonSession)
	checkError(err)
	_, err = rc.Do("set", fmt.Sprintf("userid:usertoken/%d", userid), usertoken)
	checkError(err)

	// cookie
	http.SetCookie(w, &http.Cookie{Name:"usertoken", Value:usertoken, MaxAge:sessionLifeSecond})

	// reply
	type Reply struct{
		Usertoken string
	}
	reply := Reply{usertoken}
	writeResponse(w, reply)
}


func test(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(r)
	checkError(err)

	writeResponse(w, session)
}

func regAuth() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/test", test)
}
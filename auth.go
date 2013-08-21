package main

import (
	// "fmt"
	"net/http"
	// "database/sql"
	// _ "github.com/go-sql-driver/mysql"
	// "time"
	// "encoding/json"
)

func init() {
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

const passwordSalt = "liwei"

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
	
	
}

func regAuth() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
}
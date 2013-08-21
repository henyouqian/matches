package main

import (
	"fmt"
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

type regParam struct{
	Name string
	Password string
}

func register(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")
	
	var param regParam
	err := decodeRequestBody(r, &param)
	checkError(err)

	fmt.Fprintf(w, "%+v", param)
}

func login(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")
	
	r.ParseForm()

	a, b := r.Form["username"], r.Form["password"]
	if len(a) == 0 || len(b) == 0 {
		fmt.Fprintf(w, "{error=\"param\"}")
		return;
	}
	c, d := a[0], b[0]

	fmt.Fprintf(w, "%v, %v", c, d)
}

func regAuth() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
}
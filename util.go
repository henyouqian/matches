package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"crypto/sha256"
	"encoding/base64"
)

func handleError(w http.ResponseWriter) {
	if r := recover(); r != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(w)
		err := fmt.Sprintf("%v", r)

		type errorResponse struct {
			Error string
		}
		encoder.Encode(errorResponse{err})
    }
}

func checkMathod(r *http.Request, method string) {
	if r.Method != method {
		panic("err_method_not_allowed")
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func decodeRequestBody(r *http.Request, v interface{}) error{
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

func writeResponse(w http.ResponseWriter, v interface{}) {
	encoder := json.NewEncoder(w)
	encoder.Encode(v)
}

const (
	account_db = "account_db"
)

func opendb(dbname string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("root@/%s?parseTime=true", account_db))
	checkError(err)
	return db
}

func sha224(s string) string {
	hasher := sha256.New224()
	hasher.Write([]byte(s))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
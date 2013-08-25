package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"
	"github.com/nu7hatch/gouuid"
)

var redisPool *redis.Pool

func init() {
	redisPool = &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

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

func decodeRequestBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

func writeResponse(w http.ResponseWriter, v interface{}) {
	encoder := json.NewEncoder(w)
	encoder.Encode(v)
}


func opendb(dbname string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("root@/%s?parseTime=true", dbname))
	checkError(err)
	return db
}

func sha224(s string) string {
	hasher := sha256.New224()
	hasher.Write([]byte(s))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func checkAdmin(session *Session) {
	if session.Username != "admin" {
		panic("err_denied")
	}
}

func genUUID() string {
	uuid, err := uuid.NewV4()
	checkError(err)
	return base64.URLEncoding.EncodeToString(uuid[:])
}

package main

import (
	"database/sql"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/henyouqian/lwutil"
	"time"
)

var (
	redisPool *redis.Pool
	authDB    *sql.DB
	matchDB   *sql.DB
	kvDB      *sql.DB
)

func init() {
	redisPool = &redis.Pool{
		MaxIdle:     20,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	authDB = opendb("auth_db")
	authDB.SetMaxIdleConns(10)

	matchDB = opendb("match_db")
	matchDB.SetMaxIdleConns(10)

	kvDB = opendb("kv_db")
	kvDB.SetMaxIdleConns(10)

	lwutil.StartKV(kvDB, redisPool)
}

func opendb(dbname string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("root@/%s?parseTime=true", dbname))
	if err != nil {
		panic(err)
	}
	return db
}

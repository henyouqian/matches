package main

import (
	"database/sql"
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	KEY_KV_ZSETS = "matches/kv"
)

var (
	kvDB *sql.DB
)

func init() {
	kvDB = opendb("kv_db")
	kvDB.SetMaxIdleConns(10)
}

func setKV(key string, value string, expireSec uint, rc redis.Conn) {
	rc = redisPool.Get()
	defer rc.Close()

	if expireSec > 3600 {
		expireSec = 3600
	}

	rc.Send("SET", "kv/"+key, value)
	score := time.Now().Unix() + int64(expireSec)
	rc.Send("ZADD", KEY_KV_ZSETS, score, key)
	rc.Flush()
}

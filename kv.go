package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	KEY_KV_ZSETS = "matches/kv"
)

func setKV(key string, value string, expireSec uint, rc redis.Conn) {
	if rc == nil {
		rc = redisPool.Get()
		defer rc.Close()
	}

	if expireSec > 3600 {
		expireSec = 3600
	}

	rc.Send("SET", key, value)
	score := time.Now().Unix() + int64(expireSec)
	rc.Send("ZADD", KEY_KV_ZSETS, score, key)
	rc.Flush()
}
package main

import (
	"github.com/garyburd/redigo/redis"
)

func setKV(key string, value string, expireSec uint, rc redis.Conn) {
	if rc == nil {
		rc = redisPool.Get()
		defer rc.Close()
	}

	if expireSec < 60 {expireSec = 60}

	rc.Send("SET", key, value)
	ZADD myzset 1 "one"
	rc.Send("GET", "foo")
	rc.Flush()
}
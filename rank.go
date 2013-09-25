package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/henyouqian/lwutil"
	"net/http"
)

func getMyRank(w http.ResponseWriter, r *http.Request) {
	lwutil.CheckMathod(r, "POST")

	rc := redisPool.Get()
	defer rc.Close()

	session, err := findSession(w, r, rc)
	lwutil.CheckError(err, "err_auth")

	appid := session.Appid
	if appid == 0 {
		lwutil.SendError("err_auth", "Please login with app secret")
	}

	// in
	var in struct {
		MatchId uint32
	}
	err = lwutil.DecodeRequestBody(r, &in)
	lwutil.CheckError(err, "err_decode_body")

	keyLeaderboard := makeLeaderboardKey(in.MatchId)
	rc.Send("zrank", keyLeaderboard, session.Userid)
	rc.Send("zscore", keyLeaderboard, session.Userid)
	rc.Flush()
	rank, err := redis.Int64(rc.Receive())
	score := int64(0)
	if err == redis.ErrNil {
		rank = 0
	} else {
		lwutil.CheckError(err, "")
		if err == nil {
			rank += 1
		}
		score, err = redis.Int64(rc.Receive())
		lwutil.CheckError(err, "")
	}

	// out
	out := struct {
		Rank  int64
		Score int64
	}{rank, score}
	lwutil.WriteResponse(w, out)
}

func regRank() {
	http.Handle("/rank/mine", lwutil.ReqHandler(getMyRank))
}

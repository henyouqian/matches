package main

import (
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

	//keyFail := makeFailboardKey(in.MatchId)
	//keyLeaderboard := makeLeaderboardKey(in.MatchId)

	// out
	lwutil.WriteResponse(w, in)
}

func regRank() {
	http.Handle("/rank/mine", lwutil.ReqHandler(getMyRank))
}

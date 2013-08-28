package main

import (
	"net/http"
	"time"
)

type Match struct {
	Begin time.Time
	End   time.Time
	Order string
}

func newMatch(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err, "")
	checkAdmin(session)

	// input
	type Input struct {
		Begin string
		End   string
		Order string
	}
	var input Input
	decodeRequestBody(r, &input)
	
	if input.Begin == "" || input.End == "" || 
	  (input.Order != "ASC" && input.Order != "DESC") {
		sendError("err_input", "")
	}

	//
	match := Match{}
	const timeform = "2006-01-02 15:04:05"
	begin, err := time.Parse(timeform, input.Begin)
	checkError(err, "err_shit")
	end, err := time.Parse(timeform, input.End)
	checkError(err, "")
	match.Begin = begin.Local()
	match.End = end.Local()
	match.Order = input.Order

	// redis
	rc := redisPool.Get()
	defer rc.Close()

	// reply
	writeResponse(w, &match)
}


func regMatch() {
	http.HandleFunc("/match/newmatch", newMatch)
}
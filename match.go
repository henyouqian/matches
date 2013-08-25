package main

import (
	"net/http"
)

func newMatch(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err)
	checkAdmin(session)

	// db


	// reply
	writeResponse(w, "xxx")
}



func regMatch() {
	http.HandleFunc("/match/newmatch", newMatch)
}
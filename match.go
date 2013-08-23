package main

import (
	"net/http"
	"fmt"
)

func createMatch(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	fmt.Println(err)
	checkError(err)

	if session.Username != "admin" {
		panic("err_denied")
	}

	// reply
	writeResponse(w, "xxx")
}

func regMatch() {
	http.HandleFunc("/matchapi/create", createMatch)
}
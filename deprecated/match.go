package main

import (
	"net/http"
	"time"
	// "github.com/golang/glog"
)

type Match struct {
	Id uint32
	Name string
	Gameid uint32
	Begin uint64
	End uint64
	Sort uint8
}

func newMatch(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err, "")
	checkAdmin(session)

	appid := session.Appid
	if appid == 0 {
		sendError("err_auth", "Please login with app secret")
	}

	// input
	type Input struct {
		Name   string
		Gameid uint32
		Begin  string
		End    string
		Sort   uint8
	}
	input := Input{}
	decodeRequestBody(r, &input)
	
	if input.Name == "" || input.Begin == "" || input.End == "" || input.Gameid == 0 {
		sendError("err_input", "Missing Name || Begin || End || Gameid")
	}

	if input.Sort != 0 && input.Sort != 1 {
		sendError("err_input", "Invalid Sort, must be 0 or 1")
	}

	//
	const timeform = "2006-01-02 15:04:05"
	begin, err := time.ParseInLocation(timeform, input.Begin, time.Local)
	checkError(err, "err_shit")
	end, err := time.ParseInLocation(timeform, input.End, time.Local)
	checkError(err, "")
	beginUnix := begin.Unix()
	endUnix := end.Unix()

	if endUnix - beginUnix <= 60 {
		sendError("err_input", "endUnix - beginUnix must > 60 seconds")
	}
	if time.Now().Unix() > endUnix {
		sendError("err_input", "end time before now")
	}

	// db
	db := opendb("match_db")
	defer db.Close()

	stmt, err := db.Prepare(`INSERT INTO matches (name, appid, gameid, sort, begin, end)  VALUES (?, ?, ?, ?, ?, ?)`)
	checkError(err, "")

	res, err := stmt.Exec(input.Name, appid, input.Gameid, input.Sort, beginUnix, endUnix)
	checkError(err, "")

	id, err := res.LastInsertId()
	checkError(err, "")

	// reply
	type Reply struct {
		Matchid int64
	}
	reply := Reply{id}
	writeResponse(w, &reply)
}

func listOpening(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err, "")

	//
	now := time.Now().Unix()

	// db
	db := opendb("match_db")
	defer db.Close()

	rows, err := db.Query("SELECT id, name, gameid, sort, begin, end FROM matches WHERE begin < ? AND end > ? AND appid = ?", now, now, session.Appid)
	checkError(err, "")
	
	matches := make([]Match, 0, 16)
	var match Match
	for rows.Next() {
		err = rows.Scan(&match.Id, &match.Name, &match.Gameid, &match.Sort, &match.Begin, &match.End)
		checkError(err, "")
		matches = append(matches, match)
	}

	type Reply struct {
		Matches []Match
	}
	reply := Reply{
		matches,
	}
	writeResponse(w, reply)
}

func listComming(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err, "")

	//
	now := time.Now().Unix()

	// db
	db := opendb("match_db")
	defer db.Close()

	rows, err := db.Query("SELECT id, name, gameid, sort, begin, end FROM matches WHERE begin > ? AND appid = ?", now, session.Appid)
	checkError(err, "")
	
	matches := make([]Match, 0, 16)
	var match Match
	for rows.Next() {
		err = rows.Scan(&match.Id, &match.Name, &match.Gameid, &match.Sort, &match.Begin, &match.End)
		checkError(err, "")
		matches = append(matches, match)
	}

	type Reply struct {
		Matches []Match
	}
	reply := Reply{
		matches,
	}
	writeResponse(w, reply)
}

func listClosed(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err, "")

	//
	now := time.Now().Unix()

	// db
	db := opendb("match_db")
	defer db.Close()

	rows, err := db.Query("SELECT id, name, gameid, sort, begin, end FROM matches WHERE end < ? AND appid = ?", now, session.Appid)
	checkError(err, "")
	
	matches := make([]Match, 0, 16)
	var match Match
	for rows.Next() {
		err = rows.Scan(&match.Id, &match.Name, &match.Gameid, &match.Sort, &match.Begin, &match.End)
		checkError(err, "")
		matches = append(matches, match)
	}

	type Reply struct {
		Matches []Match
	}
	reply := Reply{
		matches,
	}
	writeResponse(w, reply)
}


func regMatch() {
	http.HandleFunc("/match/new", newMatch)
	http.HandleFunc("/match/listopening", listOpening)
	http.HandleFunc("/match/listcomming", listComming)
	http.HandleFunc("/match/listclosed", listClosed)
}
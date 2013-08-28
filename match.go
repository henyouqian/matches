package main

import (
	"net/http"
	"time"
	// "github.com/golang/glog"
)

// type Match struct {
// 	Appid uint32
// 	Gameid uint32
// 	Begin time.Time
// 	End   time.Time
// 	Order string
// }

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

	_, err := findSession(w, r)
	checkError(err, "")

	//db
	db := opendb("match_db")
	defer db.Close()

	rows, err := db.Query("SELECT id FROM matches")
	checkError(err, "")
	
	ids := make([]uint32, 0, 16)
	var id uint32
	for rows.Next() {
		err = rows.Scan(&id)
		checkError(err, "")
		ids = append(ids, id)
	}

	type Reply struct {
		Ids []uint32
	}
	reply := Reply{ids}

	writeResponse(w, reply)
}

func listComming(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err, "")

	writeResponse(w, session)
}

func listClosed(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)
	checkMathod(r, "POST")

	session, err := findSession(w, r)
	checkError(err, "")

	writeResponse(w, session)
}


func regMatch() {
	http.HandleFunc("/match/newmatch", newMatch)
	http.HandleFunc("/match/listopening", listOpening)
	http.HandleFunc("/match/listcomming", listComming)
	http.HandleFunc("/match/listclosed", listClosed)
}
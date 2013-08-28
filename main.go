package main

import (
	"github.com/golang/glog"
	"net/http"
	"flag"
	"fmt"
)

func staticFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	flag.Parse()

	http.HandleFunc("/static/", staticFile)

	regAuth()
	regMatch()
	
	glog.Info("Server running")
	glog.Fatal(http.ListenAndServe(":9999", nil))
}

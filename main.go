package main

import (
	"github.com/golang/glog"
	"net/http"
	"flag"
)

func main() {
	flag.Parse()

	regAuth()
	regMatch()
	
	glog.Info("Server running")
	glog.Fatal(http.ListenAndServe(":9999", nil))
}

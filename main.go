package main

import (
	"github.com/golang/glog"
	"net/http"
	"flag"
	"runtime"
)

func staticFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	flag.Parse()

	http.HandleFunc("/static/", staticFile)

	regAuth()
	regMatch()
	regBench()

	runtime.GOMAXPROCS(runtime.NumCPU())
	
	glog.Infof("Server running: cpu=%d", runtime.NumCPU())
	glog.Fatal(http.ListenAndServe(":9999", nil))
}

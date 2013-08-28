package main

import (
	"github.com/golang/glog"
	"net/http"
	"flag"
	"runtime"
)

func main() {
	flag.Parse()

	regAuth()
	regMatch()
	regBench()

	runtime.GOMAXPROCS(runtime.NumCPU())
	
	glog.Infof("Server running: cpu=%d", runtime.NumCPU())
	glog.Fatal(http.ListenAndServe(":9999", nil))
}

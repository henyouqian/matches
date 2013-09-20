package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/henyouqian/lwutil"
)

func a() error {
	//return errors.New("error in a")
	return lwutil.NewErrStr("err in a")
}

func b() error {
	err := a()
	if err != nil {
		return lwutil.NewErr(err)
	}
	return nil
}

func c() error {
	err := b()
	if err != nil {
		return lwutil.NewErr(err)
	}
	return nil
}

func g() {
	glog.Infoln("g")
}

func kv() {
	rc := redisPool.Get()
	defer rc.Close()

	glog.Infoln("begin")
	for i := 0; i < 100000; i++ {
		lwutil.SetKV(fmt.Sprintf("bbb/%d", i), []byte("ccc"), rc)
	}
	glog.Infoln("end")
}

func lab() {

}

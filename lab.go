package main

import (
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

func lab() {
	//v, err := getKV("aaa")
	//glog.Infoln(v, err)
	lwutil.SetKV("bbb", []byte("uuu"))
	aaa, err := lwutil.GetKV("aaa")
	glog.Infoln(string(aaa), err)
}

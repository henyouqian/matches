package main

import (
	// "github.com/garyburd/redigo/redis"
	// "time"
	// "fmt"
	//"errors"
	"github.com/golang/glog"
	"github.com/henyouqian/lwUtil"
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

func lab() {
	err := c()
	glog.Errorln(err)
}

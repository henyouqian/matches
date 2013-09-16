package main

import (
	"github.com/henyouqian/lwUtil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var httpTestData struct {
	cookie *http.Cookie
}

func httpTest(fn func(http.ResponseWriter, *http.Request), method, body string) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, "", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.AddCookie(httpTestData.cookie)

	w := httptest.NewRecorder()
	lwutil.HttpTest(fn, w, req)
	return w, nil
}

func init() {
	body := `{
		"Username": "admin",
		"Password": "admin",
		"Appsecret": "waZaxlO3S1do7FWpyz3sKQ=="
	}`
	req, err := http.NewRequest("POST", "http://localhost:9999/bench/hello", strings.NewReader(body))
	lwutil.PanicIfError(err)

	w := httptest.NewRecorder()
	lwutil.HttpTest(login, w, req)

	userToken := w.Body.String()
	httpTestData.cookie = &http.Cookie{Name: "usertoken", Value: userToken, MaxAge: sessionLifeSecond, Path: "/"}
}

func TestLoginInfo(t *testing.T) {
	w, err := httpTest(loginInfo, "POST", "")

	if err != nil || w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestListApp(t *testing.T) {
	w, err := httpTest(listApp, "POST", "")
	if err != nil || w.Code != http.StatusOK {
		t.Fail()
	}
}

package main

import (
	"fmt"
	"github.com/henyouqian/lwUtil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {

}

func TestLogin(t *testing.T) {
	body := `{
		"Username": "admin",
		"Password": "admin"..asdf
	}`
	req, err := http.NewRequest("POST", "http://localhost:9999/bench/hello", strings.NewReader(body))
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	lwutil.HttpTest(login, w, req)

	fmt.Printf("%d - %s", w.Code, w.Body.String())
	t.Fail()
}

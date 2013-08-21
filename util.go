package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)

type errorResponse struct {
	Error string
}

func handleError(w http.ResponseWriter) {
	if r := recover(); r != nil {
		w.WriteHeader(http.StatusBadRequest)
		encoder := json.NewEncoder(w)
		err := fmt.Sprintf("%v", r)
		encoder.Encode(errorResponse{err})
    }
}

func checkMathod(r *http.Request, method string) {
	if r.Method != method {
		panic("err_method_not_allowed")
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func decodeRequestBody(r *http.Request, v interface{}) error{
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}
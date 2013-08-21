package main

import (
	"log"
	"net/http"
)


func main() {
	regAuth()
	log.Println("Server running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
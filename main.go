package main

import (
	"log"
	"net/http"
)

func main() {
	regAuth()
	regMatch()
	log.Println("Server running")
	log.Fatal(http.ListenAndServe(":9999", nil))
}

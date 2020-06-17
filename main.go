package main

import (
	"log"
	"net/http"
)

func main() {
	handleRequests()
}

func handleRequests() {
	http.HandleFunc("/create", electionView)

	log.Println("Listening on port", ":8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

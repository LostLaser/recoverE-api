package main

import (
	"log"
	"net/http"

	"github.com/LostLaser/recoverE-api/handlers"
)

func main() {
	handleRequests()
}

func handleRequests() {
	http.HandleFunc("/election", handlers.ElectionView)

	port := ":8888"
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

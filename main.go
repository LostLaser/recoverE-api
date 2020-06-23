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

	log.Println("Listening on port", ":8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

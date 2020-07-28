package main

import (
	"log"
	"net/http"
	"os"

	"github.com/LostLaser/recoverE-api/handlers"
)

func main() {
	handleRequests()
}

func handleRequests() {
	http.HandleFunc("/election", handlers.ElectionView)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

package main

import (
	"log"
	"net/http"

	"github.com/LostLaser/recoverE-api/handlers"
	"github.com/LostLaser/recoverE-api/utils"
)

var (
	port = utils.Get("port").(string)
)

func main() {
	handleRequests()
}

func handleRequests() {
	http.HandleFunc("/election", handlers.ElectionView)

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

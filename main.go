package main

import (
	"log"
	"net/http"

	"github.com/LostLaser/recoverE-api/config"
	"github.com/LostLaser/recoverE-api/controller"
)

var (
	port = config.Get("port").(string)
)

func main() {
	handleRequests()
}

func handleRequests() {
	m := http.NewServeMux()
	withLogger := loggingHandlerFactory()

	m.Handle("/election", withLogger(controller.ElectionView))

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, m))
}

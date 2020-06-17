package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//func electionView(rw http.ResponseWriter, req *http.Request) {
//	// Make sure that the writer supports flushing.
//	//
//	flusher, ok := rw.(http.Flusher)
//
//	if !ok {
//		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
//		return
//	}
//
//	rw.Header().Set("Content-Type", "text/event-stream")
//	rw.Header().Set("Cache-Control", "no-cache")
//	rw.Header().Set("Connection", "keep-alive")
//	rw.Header().Set("Access-Control-Allow-Origin", "*")
//
//	c := cluster.New(2, time.Second*2)
//	defer c.Purge()
//
//	for {
//		// Write to the ResponseWriter
//		// Server Sent Events compatible
//		fmt.Fprintf(rw, "data: %s\n\n", c.ReadEvent())
//
//		// Flush the data immediatly instead of buffering it for later.
//		flusher.Flush()
//	}
//}

func electionView(rw http.ResponseWriter, req *http.Request) {
	u := websocket.Upgrader{}
	c, err := u.Upgrade(rw, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// receive message
	messageType, message, err := c.ReadMessage()
	if err != nil {
		log.Println(err)
	}
	log.Println(message)

	// send message
	err = c.WriteMessage(messageType, []byte("HEY"))
	if err != nil {
		log.Println(err)
	}
}

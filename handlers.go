package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	cluster "github.com/LostLaser/recoverE"
	"github.com/gorilla/websocket"
)

func electionView(w http.ResponseWriter, r *http.Request) {

	// set up websocket based off request
	u := websocket.Upgrader{}
	u.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	// create cluster
	keys := r.URL.Query()
	count, err := strconv.Atoi(keys.Get("count"))
	if err != nil {
		http.Error(w, "Query parameter 'count' missing or invalid", http.StatusBadRequest)
	}

	socketMessaging(conn, count)
}

func socketMessaging(conn *websocket.Conn, count int) {
	c := cluster.New(count, time.Second*4)
	defer c.Purge()
	defer conn.Close()

	ids := c.ServerIds()
	err := conn.WriteJSON(constructMessage("setup", ids))
	if err != nil {
		log.Println(err)
	}

	// receive messages
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(message)
		}
	}()

	// stream cluster events to client
	for {
		err := conn.WriteJSON(constructMessage("event", c.ReadEvent()))
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func constructMessage(operation string, payload interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	m["operation"] = operation
	m["payload"] = payload

	return m
}

package handlers

import (
	"log"
	"time"

	cluster "github.com/LostLaser/recoverE"
	"github.com/gorilla/websocket"
)

func socketMessaging(conn *websocket.Conn, count int) {
	c := cluster.New(count, time.Second*4)
	defer c.Purge()
	defer conn.Close()

	ids := c.ServerIds()
	err := conn.WriteJSON(map[string]interface{}{"action": "SETUP", "payload": ids})
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
		err := conn.WriteJSON(c.ReadEvent())
		if err != nil {
			log.Println(err)
			return
		}
	}
}

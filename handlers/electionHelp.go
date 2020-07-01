package handlers

import (
	"log"
	"time"

	election "github.com/LostLaser/election"
	"github.com/gorilla/websocket"
)

func socketMessaging(conn *websocket.Conn, count int) {
	c := election.New(count, time.Second*4)
	defer c.Purge()
	defer conn.Close()

	ids := c.ServerIds()
	err := conn.WriteJSON(map[string]interface{}{"action": "SETUP", "payload": ids})
	if err != nil {
		log.Println(err)
	}

	// receive messages
	go responseMessage(conn, c)

	// stream cluster events to client with a delay
	for {
		time.Sleep(time.Millisecond * 50)
		err := conn.WriteJSON(c.ReadEvent())
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func responseMessage(conn *websocket.Conn, c *election.Cluster) {
	type message struct {
		Action string
		ID     string
	}

	for {
		msg := message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}

		switch action := msg.Action; action {
		case "STOP":
			c.StopServer(msg.ID)
		case "START":
			c.StartServer(msg.ID)
		}

	}
}

package service

import (
	"log"
	"time"

<<<<<<< HEAD:service/electionService.go
	"github.com/LostLaser/election"
	"github.com/LostLaser/recoverE-api/config"
=======
	election "github.com/LostLaser/election"
	"github.com/LostLaser/recoverE-api/utils"
>>>>>>> 7d7d336a62fe86e74c588ffc98d19e6bb9bb4779:handlers/electionHelp.go
	"github.com/gorilla/websocket"
)

var (
<<<<<<< HEAD:service/electionService.go
	stopMessage  = config.Get("election.node.process.stop-message").(string)
	startMessage = config.Get("election.node.process.start-message").(string)
	setupMessage = config.Get("election.node.process.initial-node-setup").(string)
)

// SocketMessaging handles communication pertaining to created cluster
func SocketMessaging(conn *websocket.Conn, count int) {
=======
	stopMessage  = utils.Get("node.process.stop-message").(string)
	startMessage = utils.Get("node.process.start-message").(string)
	setupMessage = utils.Get("node.process.initial-node-setup").(string)
)

func socketMessaging(conn *websocket.Conn, count int) {
>>>>>>> 7d7d336a62fe86e74c588ffc98d19e6bb9bb4779:handlers/electionHelp.go
	c := election.New(count, time.Second*4)
	defer c.Purge()
	defer conn.Close()

	ids := c.ServerIds()
	err := conn.WriteJSON(map[string]interface{}{"action": setupMessage, "payload": ids})
	if err != nil {
		log.Println(err)
	}

	// receive messages
	go responseMessage(conn, c)

	// stream cluster events to client with a delay
	for {
		time.Sleep(time.Millisecond * 200)
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
		case stopMessage:
			c.StopServer(msg.ID)
		case startMessage:
			c.StartServer(msg.ID)
		}

	}
}

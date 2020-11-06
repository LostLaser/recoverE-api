package service

import (
	"time"

	"github.com/LostLaser/election"
	"github.com/LostLaser/election/server"
	"github.com/LostLaser/recoverE-api/config"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	stopMessage       = config.Get("election.node.process.stop-message").(string)
	startMessage      = config.Get("election.node.process.start-message").(string)
	setupMessage      = config.Get("election.node.process.initial-node-setup").(string)
	connectionTimeout = config.Get("election.timeout-minutes").(int)
)

// Messenger handles communication pertaining to created cluster
func Messenger(conn *websocket.Conn, count int, electionSetup server.Setup) {
	c := election.New(electionSetup, count, time.Second*4)
	defer c.Purge()
	defer conn.Close()

	ids := c.ServerIds()
	err := conn.WriteJSON(map[string]interface{}{
		"action":  setupMessage,
		"payload": ids,
	})
	if err != nil {
		log.Debug(err)
	}

	exp := make(chan (bool))

	// receive messages
	go responseMessage(conn, c, exp)

	go expireSocket(conn, exp)

	// stream cluster events to client with a delay
	for {
		err := conn.WriteJSON(c.ReadEvent())
		if err != nil {
			log.Debug(err)
			return
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func responseMessage(conn *websocket.Conn, c *election.Cluster, exp chan bool) {
	defer conn.Close()
	type message struct {
		Action string
		ID     string
	}

	for {
		msg := message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Debug(err)
			return
		}
		exp <- false
		switch action := msg.Action; action {
		case stopMessage:
			c.StopServer(msg.ID)
		case startMessage:
			c.StartServer(msg.ID)
		}

	}
}

func expireSocket(conn *websocket.Conn, exp chan bool) {
	defer conn.Close()
	duration := time.Minute * time.Duration(connectionTimeout)
	mw := time.Second * 2
	closeCode := 4001

	timer := time.NewTimer(duration)
	for {
		select {
		case <-exp:
			timer.Reset(duration)
		case <-timer.C:
			msg := websocket.FormatCloseMessage(closeCode, "session expired due to inactivity")
			conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(mw))
			time.Sleep(mw)
			return
		}
	}
}

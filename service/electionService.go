package service

import (
	"time"

	"github.com/LostLaser/election"
	"github.com/LostLaser/election/server"
	"github.com/LostLaser/recoverE-api/config"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	stopMessage   = config.Get("election.node.process.stop-message").(string)
	startMessage  = config.Get("election.node.process.start-message").(string)
	setupMessage  = config.Get("election.node.process.initial-node-setup").(string)
	timeoutExpire = config.Get("election.timeout.expired").(int)
	timeoutHard   = config.Get("election.timeout.hard").(int)
)

// Messenger handles communication pertaining to created cluster
func Messenger(conn *websocket.Conn, count int, electionSetup server.Setup, logger *zap.Logger) {
	c := election.New(electionSetup, count, time.Second*4, logger)
	defer c.Purge()
	defer conn.Close()

	ids := c.ServerIds()
	err := conn.WriteJSON(map[string]interface{}{
		"action":  setupMessage,
		"payload": ids,
	})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	exp := make(chan (bool))

	// receive messages
	go responseMessage(conn, c, exp, logger)

	go expireSocket(conn, exp, logger)

	// stream cluster events to client
	for {
		err := conn.WriteJSON(c.ReadEvent())
		if err != nil {
			logger.Debug(err.Error())
			return
		}
	}
}

func responseMessage(conn *websocket.Conn, c *election.Cluster, exp chan bool, logger *zap.Logger) {
	defer conn.Close()
	type message struct {
		Action string
		ID     string
	}

	for {
		msg := message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			logger.Debug(err.Error())
			return
		}
		exp <- false
		switch action := msg.Action; action {
		case stopMessage:
			logger.Debug("Request to stop simulated server", zap.String("Cluster ID", c.ID), zap.String("Server ID", msg.ID))
			c.StopServer(msg.ID)
		case startMessage:
			logger.Debug("Request to start simulated server", zap.String("Cluster ID", c.ID), zap.String("Server ID", msg.ID))
			c.StartServer(msg.ID)
		}
	}
}

func expireSocket(conn *websocket.Conn, exp chan bool, logger *zap.Logger) {
	defer conn.Close()
	expireTime := time.Minute * time.Duration(timeoutExpire)
	hardResetTime := time.Minute * time.Duration(timeoutHard)
	mw := time.Second * 2
	expireCode := 4001
	hardResetCode := 4002

	expirationTimer := time.NewTimer(expireTime)
	hardTimer := time.NewTimer(hardResetTime)

	for {
		select {
		case <-exp:
			expirationTimer.Reset(expireTime)
		case <-expirationTimer.C:
			expireMsg := "session expired due to inactivity"
			logger.Debug(expireMsg)
			msg := websocket.FormatCloseMessage(expireCode, expireMsg)
			conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(mw))
			time.Sleep(mw)
			return
		case <-hardTimer.C:
			logger.Debug("Session hit hard timeout limit", zap.Duration("Hard timeout", expireTime))
			msg := websocket.FormatCloseMessage(hardResetCode, "Maximum time hit for live connection")
			conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(mw))
			time.Sleep(mw)
			return
		}
	}
}

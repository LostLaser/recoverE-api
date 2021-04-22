package controller

import (
	"net/http"
	"strconv"

	"github.com/LostLaser/election/server"
	"github.com/LostLaser/election/server/bully"
	"github.com/LostLaser/election/server/ring"
	"github.com/LostLaser/recoverE-api/config"
	"github.com/LostLaser/recoverE-api/service"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	minNodes       = config.Get("election.node.min").(int)
	maxNodes       = config.Get("election.node.max").(int)
	allowedOrigins = config.Get("allowed-origins").([]interface{})
)

// ElectionView handles the full interaction
func ElectionView(w http.ResponseWriter, r *http.Request, logger *zap.Logger) {

	// input validation
	keys := r.URL.Query()
	count, err := strconv.Atoi(keys.Get("count"))
	if err != nil {
		http.Error(w, "Query parameter 'count' missing or invalid", http.StatusBadRequest)
		return
	}
	if count > maxNodes {
		count = maxNodes
	} else if count < minNodes {
		count = minNodes
	}

	electionType := keys.Get("election_type")
	if electionType == "" {
		http.Error(w, "Query parameter 'election_type' missing or invalid", http.StatusBadRequest)
		return
	}
	var electionSetup server.Setup
	switch electionType {
	case "bully":
		electionSetup = bully.Setup{}
	case "ring":
		electionSetup = ring.Setup{}
	default:
		http.Error(w, "Invalid election type", http.StatusBadRequest)
		return
	}

	// set up websocket based off request
	u := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("origin")
			for _, o := range allowedOrigins {
				if o == origin {
					logger.Debug("Request authorized with origin", zap.String("origin", origin))
					return true
				}
			}
			logger.Debug("Origin was not in the allow list", zap.String("origin", origin))
			return false
		},
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		logger.Debug(err.Error())
		return
	}

	service.Messenger(conn, count, electionSetup, logger)
}

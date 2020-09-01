package controller

import (
	"net/http"
	"strconv"

	"github.com/LostLaser/recoverE-api/config"
	"github.com/LostLaser/recoverE-api/service"
	"github.com/gorilla/websocket"
)

var (
	minNodes       = config.Get("election.node.min").(int)
	maxNodes       = config.Get("election.node.max").(int)
	allowedOrigins = config.Get("allowed-origins").([]interface{})
)

// ElectionView handles the full interaction
func ElectionView(w http.ResponseWriter, r *http.Request) {

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

	// set up websocket based off request
	u := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("origin")
			for _, o := range allowedOrigins {
				if o == origin {
					return true
				}
			}
			return false
		},
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	service.Messenger(conn, count)
}

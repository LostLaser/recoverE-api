package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/LostLaser/recoverE-api/config"
	"github.com/LostLaser/recoverE-api/service"
	"github.com/gorilla/websocket"
)

var (
	minNodes = config.Get("election.node.min").(int)
	maxNodes = config.Get("election.node.max").(int)
)

// ElectionView handles the full interaction
func ElectionView(w http.ResponseWriter, r *http.Request) {

	// set up websocket based off request
	u := websocket.Upgrader{}
	u.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	// input validation
	keys := r.URL.Query()
	count, err := strconv.Atoi(keys.Get("count"))
	if err != nil {
		http.Error(w, "Query parameter 'count' missing or invalid", http.StatusBadRequest)
	}
	electionType := keys.Get("election_type")
	if electionType == "" {
		http.Error(w, "Query parameter 'electionType' missing or invalid", http.StatusBadRequest)
	}

	log.Print(maxNodes, minNodes)
	if count > maxNodes {
		count = maxNodes
	} else if count < minNodes {
		count = minNodes
	}

	service.SocketMessaging(conn, count)
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
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

	max := 5
	min := 2

	if count > max {
		count = max
	} else if count < min {
		count = min
	}

	socketMessaging(conn, count)
}

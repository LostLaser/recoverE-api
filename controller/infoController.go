package controller

import (
	"net/http"

	"go.uber.org/zap"
)

func LandingView(w http.ResponseWriter, r *http.Request, log *zap.Logger) {
	w.Write([]byte("Welcome to the Recover-E API"))
}

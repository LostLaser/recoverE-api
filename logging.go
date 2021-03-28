package main

import (
	"net/http"
	"strings"

	"github.com/LostLaser/recoverE-api/config"
	"go.uber.org/zap"
)

var (
	logLevel = strings.ToUpper(config.Get("application.log-level").(string))
)

type loggingHandlerFunc = func(w http.ResponseWriter, r *http.Request, l *zap.Logger)

type loggingHandler struct {
	logger      *zap.Logger
	handlerFunc loggingHandlerFunc
}

func loggingHandlerFactory() func(loggingHandlerFunc) *loggingHandler {
	l, _ := zap.NewProduction()
	switch logLevel {
	case "DEVELOPMENT":
		l, _ = zap.NewDevelopment()
	case "PRODUCTION":
		l, _ = zap.NewProduction()
	case "NONE":
		l = zap.NewNop()
	}
	return func(hf loggingHandlerFunc) *loggingHandler {
		return &loggingHandler{l, hf}
	}
}

func (lh *loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lh.logger.Debug("Got a request!")
	lh.handlerFunc(w, r, lh.logger)
}

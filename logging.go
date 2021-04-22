package main

import (
	"net/http"
	"strings"

	"github.com/LostLaser/recoverE-api/config"
	"go.uber.org/zap"
)

var (
	logLevelKey = "application.log-level"
	logLevel    = strings.ToUpper(config.Get(logLevelKey).(string))
)

type loggingHandlerFunc = func(w http.ResponseWriter, r *http.Request, l *zap.Logger)

type loggingHandler struct {
	logger      *zap.Logger
	handlerFunc loggingHandlerFunc
}

func loggingHandlerFactory() func(loggingHandlerFunc) *loggingHandler {
	l, _ := zap.NewProduction()
	switch logLevel {
	case "DEV":
		l, _ = zap.NewDevelopment()
	case "PROD":
		l, _ = zap.NewProduction()
	case "NONE":
		l = zap.NewNop()
	default:
		l.Warn("Unknown log level supplied", zap.String(logLevelKey, logLevel))
	}
	return func(hf loggingHandlerFunc) *loggingHandler {
		return &loggingHandler{l, hf}
	}
}

func (lh *loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lh.handlerFunc(w, r, lh.logger)
}

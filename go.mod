module github.com/LostLaser/recoverE-api

go 1.15

require (
	github.com/LostLaser/election v1.0.0
	github.com/gorilla/websocket v1.4.2
	github.com/labstack/gommon v0.3.0
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.1
	go.uber.org/zap v1.16.0
	golang.org/x/sys v0.0.0-20201106081118-db71ae66460a // indirect
)

replace github.com/LostLaser/election => ../election

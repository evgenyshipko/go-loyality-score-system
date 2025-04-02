package main

import (
	"github.com/evgenyshipko/go-loyality-score-system/internal/logger"
	"github.com/evgenyshipko/go-loyality-score-system/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer func() {
		logger.Sync()
	}()

	customServer := server.Create()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)

	go customServer.Start()

	<-stopSignal

	customServer.ShutDown()
}

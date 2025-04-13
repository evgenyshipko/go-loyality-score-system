package main

import (
	"github.com/evgenyshipko/go-rag-chat-helper/internal/logger"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/server"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer func() {
		logger.Sync()
	}()

	err := godotenv.Load()
	if err != nil {
		logger.Instance.Errorf("Ошибка загрузки .env файла: %v", err)
		return
	}

	customServer := server.Create()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)

	go customServer.Start()

	<-stopSignal

	customServer.ShutDown()

	// РАБОТА ЗАПРОСОВ В LLM
	//res, err := llm.GetKeywords("Как работать с Images в Next.js?")
	//if err != nil {
	//	logger.Instance.Warnw("send gpt request failed", "error", err)
	//}
	//
	//logger.Instance.Info(res)
	//logger.Instance.Warnw("send gpt request success", "result", res)
}

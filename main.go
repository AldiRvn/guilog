package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/logdyhq/logdy-core/logdy"
)

func main() {
	appPort := "8080"

	logdyLogger := logdy.InitializeLogdy(logdy.Config{
		ServerIp:   "127.0.0.1",
		ServerPort: appPort,
	}, nil)

	fmt.Println("log by fmt")
	log.Println("log by log")
	slog.Info("log by slog")

	logdyLogger.LogString("This is a message")
	logdyLogger.LogString("This is a message 2")
	logdyLogger.Log(logdy.Fields{"msg": "supports structured logs too", "url": "some url here"})
	logdyLogger.Log(logdy.Fields{
		"correlation_id": "a", "msg": "Hello",
	})
	logdyLogger.Log(logdy.Fields{
		"correlation_id": "b", "msg": "Hi",
	})
	logdyLogger.Log(logdy.Fields{
		"correlation_id": "a", "msg": "World",
	})

	go func() {
		for i := 0; i < 100; i++ {
			log.Println(i)
			logdyLogger.Log(logdy.Fields{
				"correlation_id": "b", "msg": fmt.Sprintf("Hi: %d", i),
			})
			time.Sleep(1 * time.Second)
		}
	}()

	log.Printf("web: http://0.0.0.0:%s\n", appPort)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit
	os.Exit(0)
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"hq.0xa1.red/axdx/scheduler/internal/platform/api"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database/redis"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	database.SetBackend(string(database.KindRedis))

	//etcd.SetEtcdEndpoints("localhost:22379")
	redis.SetPassword("test")

	_, err := database.New()
	if err != nil {
		log.Panic(err)
	}

	// TODO make this close all clients
	defer func() {
		redis.Close()
	}()

	api.SetAddress("127.0.0.1:8080")
	httpServer := api.Start()

Loop:
	for {
		select {
		case <-signals:
			log.Println("shutting down http server")
			httpServer.Shutdown(context.TODO()) // nolint
			break Loop
		case err := <-httpServer.ErrorChannel:
			if err != nil {
				log.Fatalf("http server error: %v", err)
			}
			break Loop
		}
	}
}

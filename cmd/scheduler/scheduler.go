package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"hq.0xa1.red/axdx/scheduler/internal/platform/api"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

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

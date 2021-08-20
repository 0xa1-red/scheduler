package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"hq.0xa1.red/axdx/scheduler/internal/config"
	"hq.0xa1.red/axdx/scheduler/internal/logging"
	"hq.0xa1.red/axdx/scheduler/internal/platform/api"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database/redis"
)

var (
	configPath string
)

func main() {
	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")
	flag.Parse()

	if err := config.ConfigurePackages(configPath); err != nil {
		log.Panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	logger := logging.MustNew()

	database.SetBackend(string(database.KindRedis))

	redis.SetPassword("test")

	_, dbErr := database.New()
	if dbErr != nil {
		logger.Panicw("error connecting to the database", "kind", database.Backend(), "error", dbErr)
	}

	// TODO make this close all clients
	defer func() {
		database.Close()
	}()

	api.SetAddress("127.0.0.1:8080")
	httpServer := api.Start()

Loop:
	for {
		select {
		case <-signals:
			fmt.Printf("\r") // This just to remove the nasty ^C :shrug:
			logger.Info("shutting down http server")
			httpServer.Shutdown(context.TODO()) // nolint
			break Loop
		case err := <-httpServer.ErrorChannel:
			if err != nil {
				logger.Errorw("http server error", "error", err)
			}
			break Loop
		}
	}
}

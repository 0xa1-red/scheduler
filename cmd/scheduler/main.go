package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hq.0xa1.red/axdx/scheduler/internal/config"
	"hq.0xa1.red/axdx/scheduler/internal/logging"
	"hq.0xa1.red/axdx/scheduler/internal/platform/api"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database"
	"hq.0xa1.red/axdx/scheduler/internal/pusher/stdout"
)

var (
	commitHash string
	buildtime  string
	tag        string

	version    bool
	configPath string
)

func main() {
	flag.BoolVar(&version, "version", false, "version information")
	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")
	flag.Parse()

	if version {
		os.Exit(showVersion())
	}

	if err := config.ConfigurePackages(configPath); err != nil {
		log.Panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	logger := logging.MustNew()

	database.SetBackend(string(database.KindRedis))

	_, dbErr := database.New()
	if dbErr != nil {
		logger.Panicw("error connecting to the database", "kind", database.Backend(), "error", dbErr)
	}

	// TODO make this close all clients
	defer func() {
		database.Close()
	}()

	httpServer := api.Start()

	pusher := stdout.New()
	pusher.Start(time.Second)

Loop:
	for {
		select {
		case <-signals:
			fmt.Printf("\r") // This just to remove the nasty ^C :shrug:
			logger.Info("shutting down http server")
			httpServer.Shutdown(context.TODO()) // nolint
			pusher.Shutdown(context.TODO())     // nolint
			break Loop
		case httpError := <-httpServer.ErrorChannel:
			if httpError != nil {
				logger.Errorw("http server error", "error", httpError)
			}
			break Loop
		case pusherError := <-pusher.Errors():
			if pusherError != nil {
				logger.Errorw("pusher error", "error", pusherError)
			}
		}

	}
}

func showVersion() int {
	fmt.Printf(`
Scheduler (https://hq.0xa1.red/axdx/scheduler)	
===
Version: %s
Commit hash: %s
Build time: %s
`, tag, commitHash, buildtime)

	return 0
}

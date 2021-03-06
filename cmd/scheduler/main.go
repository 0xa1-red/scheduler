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

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/handlers"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/repository"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/db/pgdb"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/callback"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/logs"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	//Get config
	cfg := config.GetConfig()

	//Get logger
	logger := logs.Get()
	logger.Infof("config loaded successfully. Logger initialized with log level: %s", cfg.LogLevel)

	//Connect to database
	pgDB, err := pgdb.Dial()
	if err != nil {
		return fmt.Errorf("error while connecting to database %v", err)
	}

	//Migrations

	//Repository initializing
	recRepo := repository.New(pgDB)

	//Service initializing
	callbackHandler := handlers.NewCallbackHandler(recRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", callbackHandler.Post)

	//Create and run callback server
	cb := new(callback.Server)
	go func() {
		if err := cb.Run(mux); err != nil && err != http.ErrServerClosed {
			logger.Errorf("error while initializing callback: %v", err)
		}
	}()
	logger.Infof("callback server successfully loaded on port %s", cfg.CallbackPort)

	//Graceful shutdown callback server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logger.Info("shutting down callback server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := cb.Shutdown(ctx); err != nil {
		return fmt.Errorf("error while shutting down callback server")
	}

	return nil
}

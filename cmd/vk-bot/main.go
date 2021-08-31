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

	//Create and run callback server
	cb := new(callback.Server)
	go func() {
		if err := cb.Run(cfg.CallbackPort); err != nil && err != http.ErrServerClosed {
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

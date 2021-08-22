package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"projects/vk-stitch-bot/pkg/callback"
	"projects/vk-stitch-bot/pkg/config"
	"projects/vk-stitch-bot/pkg/logs"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("while getting config, caught error: %v", err)
	}
	logger := logs.Get(cfg.LogLevel)
	logger.Infof("Config loaded successfully. Logger initialized with log level: %s", cfg.LogLevel)

	cb := new(callback.Server)

	go func() {
		if err = cb.Run(cfg.CallbackPort); err != nil {
			logger.Errorf("error while initializing callback: %v", err)
		}
	}()
	logger.Infof("callback server successfully loaded on port %s", cfg.CallbackPort)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("Shutting down callback server")
	if err := cb.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("error while shutting down callback server")
	}

	return nil
}

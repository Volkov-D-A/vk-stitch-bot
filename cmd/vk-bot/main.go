package main

import (
	"log"
	"projects/vk-stitch-bot/pkg/config"
	"projects/vk-stitch-bot/pkg/logs"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.GetConfig()
	logger := logs.Get(cfg.LogLevel)
	logger.Infof("Config loaded successfully. Logger initialized with log level: %s", cfg.LogLevel)

	return nil
}

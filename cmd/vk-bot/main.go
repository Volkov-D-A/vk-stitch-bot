package main

import (
	"fmt"
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
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("while getting config, caught error: %v", err)
	}
	logger := logs.Get(cfg.LogLevel)
	logger.Infof("Config loaded successfully. Logger initialized with log level: %s", cfg.LogLevel)

	return nil
}

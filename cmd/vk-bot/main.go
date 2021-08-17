package main

import (
	"fmt"
	"log"
	"projects/vk-stitch-bot/pkg/config"
)

func main() {
	if err := run(); err != nil {log.Fatal(err) }
}

func run() error {
	cfg := config.GetConfig()
	fmt.Println(cfg)
	return nil
}
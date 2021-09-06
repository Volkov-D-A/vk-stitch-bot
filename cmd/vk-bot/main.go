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

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/db/pg"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/services"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/handlers"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/repository"

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
	DB, err := pg.Dial()
	if err != nil {
		return fmt.Errorf("error while connecting to database %v", err)
	}

	//Clean architecture repository - services - handlers
	//Repository initializing
	repos := repository.NewRepository(DB)
	//Service initializing
	service := services.NewService(repos, cfg)
	//Handler initializing
	callbackHandler := handlers.NewCallbackHandler(service, logger, cfg)

	//Init and setup VK callback server
	//if err = service.SetCallbackUrl(); err != nil {
	//	return err
	//}

	//FIXME: test
	service.Keyboard.SendProductKeyboard(&models.MessageRecipient{Id: 1111})

	//Create and run callback server
	cb := new(callback.Server)
	go func() {
		if err := cb.Run(callbackHandler.InitRoutes()); err != nil && err != http.ErrServerClosed {
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
		DB.Close()
		cancel()
	}()
	if err := cb.Shutdown(ctx); err != nil {
		return fmt.Errorf("error while shutting down callback server")
	}

	return nil
}

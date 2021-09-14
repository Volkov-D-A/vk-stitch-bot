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

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/db"

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
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("error getting config: %v", err)
	}

	//Get logger
	logger := logs.Get(cfg.LogLevel)
	logger.Infof("config loaded successfully. Logger initialized with log level: %s", cfg.LogLevel)

	//Connect to database
	pgUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.PG.User, cfg.PG.Password, cfg.PG.Host, cfg.PG.Database)
	DB, err := pg.Dial(pgUrl)
	if err != nil {
		return fmt.Errorf("error while connecting to database %v", err)
	}

	//Run migrations
	githubUrl := fmt.Sprintf("github://%s:%s@%s/%s/%s", cfg.Github.User, cfg.Github.Token, cfg.Github.User, cfg.Github.Repo, cfg.Github.Path)
	if err := db.PgMigrate(githubUrl, pgUrl); err != nil {
		return fmt.Errorf("error while migrating database %v", err)
	}

	//Clean architecture repository - services - handlers
	//Repository initializing
	repos := repository.NewRepository(DB, cfg)
	//Service initializing
	service := services.NewService(repos, cfg)
	//Handler initializing
	callbackHandler := handlers.NewCallbackHandler(service, logger, cfg)

	//Init and setup VK callback server
	configured, err := service.CheckCallbackServerInfo()
	if err != nil {
		return fmt.Errorf("error while checking callback server status: %v", err)
	}
	if !configured {
		if err = service.SetupCallback(); err != nil {
			return err
		}
	}

	//Create and run callback server
	cb := new(callback.Server)
	go func() {
		if err := cb.Run(callbackHandler.InitRoutes(), cfg.Callback.Port); err != nil && err != http.ErrServerClosed {
			logger.Errorf("error while initializing callback: %v", err)
		}
	}()
	logger.Infof("callback server successfully loaded on port %s", cfg.Callback.Port)

	//InitDatabase
	if err := service.InitDatabase(); err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}

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

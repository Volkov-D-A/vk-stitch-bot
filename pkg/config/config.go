package config

import (
	"errors"
	"os"
	"sync"
)

var (
	EnvParamNotExist = errors.New("environment parameter not exist")
	EnvParamIsEmpty = errors.New("environment parameter is empty")
)

type Config struct {
	Token string
}

var (
	once sync.Once
	config Config
)

func GetConfig() *Config {
	once.Do(func() {
		config.Token = os.Getenv("VK_TOKEN")
	})
	return &config
}
package config

import (
	"errors"
	"fmt"
	"log"
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
	var err error
	once.Do(func() {
		config.Token, err = getParam("VK_TOKEN")
		if err != nil {
			log.Fatal(err)
		}
	})
	return &config
}

func getParam(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("for key %s caught error: %v", key, EnvParamNotExist)
	} else {
		if val == "" {
			return "", fmt.Errorf("for key %s caught error: %v", key, EnvParamIsEmpty)
		}
	}
	return val, nil
}
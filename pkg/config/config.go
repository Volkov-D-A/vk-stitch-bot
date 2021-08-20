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
	EnvParamIsEmpty  = errors.New("environment parameter is empty")
	EnvKeyTooShort   = errors.New("environment key is too short")
)

type Config struct {
	Token    string
	LogLevel string
}

var (
	once   sync.Once
	config Config
)

func GetConfig() *Config {
	var err error
	once.Do(func() {
		// NOTE: can use external package to getting env parameter to config like viper or kelseyhightower/envconfig
		config.Token, err = getParam("VK_TOKEN")
		if err != nil {
			log.Fatalf("When getting params, caught error: %v", err)
		}
		config.LogLevel, err = getParam("LOG_LEVEL")
		if err != nil {
			log.Fatalf("When getting params, caught error: %v", err)
		}
	})
	return &config
}

func getParam(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("for key: '%s',  error: '%v'", key, EnvKeyTooShort)
	}
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("for key: '%s',  error: '%v'", key, EnvParamNotExist)
	} else {
		if val == "" {
			return "", fmt.Errorf("for key: '%s', error: '%v'", key, EnvParamIsEmpty)
		}
	}
	return val, nil
}

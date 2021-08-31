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
	Token        string //Token for accessing to VK api
	LogLevel     string //LogLevel for logger can be: panic, fatal, error, warn, info, debug, trace
	CallbackPort string //CallbackPort for listening callback server
}

var (
	config Config
	once   sync.Once
)

func GetConfig() *Config {
	// NOTE: can use external package to getting env parameter to config like viper or kelseyhightower/envconfig

	once.Do(func() {
		var err error
		//Getting config params from environment
		config.Token, err = getParam("VK_TOKEN")
		if err != nil {
			log.Fatal(fmt.Errorf("when getting params, caught error: %v", err))
		}
		config.LogLevel, err = getParam("LOG_LEVEL")
		if err != nil {
			log.Fatal(fmt.Errorf("when getting params, caught error: %v", err))
		}
		config.CallbackPort, err = getParam("CALLBACK_PORT")
		if err != nil {
			log.Fatal(fmt.Errorf("when getting params, caught error: %v", err))
		}
	})

	return &config
}

//Returns ENV parameter or error if not exist
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

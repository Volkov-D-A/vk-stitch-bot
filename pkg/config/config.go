package config

import (
	"errors"
	"fmt"
	"os"
)

var (
	EnvParamNotExist = errors.New("environment parameter not exist")
	EnvParamIsEmpty  = errors.New("environment parameter is empty")
	EnvKeyTooShort   = errors.New("environment key is too short")
)

type Config struct {
	Token    string
	LogLevel string
	CallbackPort string
}

var (
	config Config
)

func GetConfig() (*Config, error) {
	// NOTE: can use external package to getting env parameter to config like viper or kelseyhightower/envconfig
	var err error
	config.Token, err = getParam("VK_TOKEN")
	if err != nil {
		return nil, fmt.Errorf("when getting params, caught error: %v", err)
	}
	config.LogLevel, err = getParam("LOG_LEVEL")
	if err != nil {
		return nil, fmt.Errorf("when getting params, caught error: %v", err)
	}
	config.CallbackPort, err = getParam("CALLBACK_PORT")
	if err != nil {
		return nil, fmt.Errorf("when getting params, caught error: %v", err)
	}
	return &config, nil
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

package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel string `envconfig:"LOG_LEVEL"` //LogLevel for logger can be: panic, fatal, error, warn, info, debug, trace
	VK
	Callback
	Github
	PG
}

type VK struct {
	APIUrl     string `envconfig:"API_URL"`        //APIUrl - VK api url
	Token      string `envconfig:"VK_TOKEN"`       //Token for secure accessing to VK api
	Group      string `envconfig:"VK_GROUP"`       //Group - VK group id
	GroupOwner int    `envconfig:"VK_GROUP_OWNER"` //Group Owner - VK group owner id
}

type Callback struct {
	URL    string `envconfig:"CALLBACK_URL"`    //URL - callback url
	Port   string `envconfig:"CALLBACK_PORT"`   //Port - callback port
	Secret string `envconfig:"CALLBACK_SECRET"` //Secret - for secure handling callback requests
	Title  string `envconfig:"CALLBACK_TITLE"`  //Title - name of callback server
}

type Github struct {
	Token string `envconfig:"GH_TOKEN"`      //Token - Github token for secure access to migrations
	User  string `envconfig:"GH_USER"`       //User - Github user
	Repo  string `envconfig:"GH_REPOSITORY"` //Repo - Github repository
	Path  string `envconfig:"MIGRATE_PATH"`  //Path to migrations on Github
}

type PG struct {
	User     string `envconfig:"PG_USER"`     //User - Postgres user
	Password string `envconfig:"PG_PASSWORD"` //Password - Postgres password
	Database string `envconfig:"PG_DB"`       //Database - Postgres database
	Host     string `envconfig:"PG_HOST"`     //Host - Postgres host
}

func GetConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, fmt.Errorf("error creating config: %v", err)
	}
	return &config, nil
}

package pgdb

import (
	"errors"
	"fmt"
	"time"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"

	"github.com/go-pg/pg/v10"
)

const Timeout = 5

type DB struct {
	*pg.DB
}

func Dial() (*DB, error) {
	cfg := config.GetConfig()

	if cfg.PgURL == "" {
		return nil, errors.New("no postgres url specified")
	}
	pgOpts, err := pg.ParseURL(cfg.PgURL)
	if err != nil {
		return nil, err
	}
	pgDB := pg.Connect(pgOpts)
	_, err = pgDB.Exec("SELECT 1")
	if err != nil {
		return nil, fmt.Errorf("error while pinging database: %v", err)
	}
	pgDB.WithTimeout(time.Second * time.Duration(Timeout))
	return &DB{pgDB}, nil
}

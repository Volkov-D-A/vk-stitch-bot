package pg

import (
	"context"
	"fmt"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func Dial() (*DB, error) {
	cfg := config.GetConfig()
	conn, err := pgxpool.Connect(context.Background(), cfg.PgURL)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to database, %v", err)
	}
	return &DB{conn}, nil
}

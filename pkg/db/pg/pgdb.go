package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func Dial(url string) (*DB, error) {
	conn, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to database, %v", err)
	}
	return &DB{conn}, nil
}

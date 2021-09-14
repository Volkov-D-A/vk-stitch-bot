package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func PgMigrate(pg string) error {
	m, err := migrate.New("github://Volkov-D-A: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@Volkov-D-A/vk-stitch-bot/migrations", pg)
	if err != nil {
		return fmt.Errorf("error while creating migrations %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error while migrating %v", err)
	}
	return nil
}

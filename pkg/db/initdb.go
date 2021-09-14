package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func PgMigrate(githubUrl, pgUrl string) error {
	m, err := migrate.New(githubUrl, pgUrl)
	if err != nil {
		return fmt.Errorf("error while creating migrations %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error while migrating %v", err)
	}
	return nil
}

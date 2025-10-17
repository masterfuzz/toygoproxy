package migrations

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/masterfuzz/toygoproxy/pkg/database/migrations/fs"
)

func Run(pool *pgxpool.Pool) error {
	msrc, mErr := iofs.New(fs.Migrations, fs.MigrationsPath)
	if mErr != nil {
		return fmt.Errorf("Unable to create migrations source: %v", mErr)
	}

	driver, dErr := mpgx.WithInstance(stdlib.OpenDBFromPool(pool), &mpgx.Config{})
	if dErr != nil {
		return fmt.Errorf("Unable to create migrations driver: %v", dErr)
	}

	m, mErr := migrate.NewWithInstance("iofs", msrc, "pgx", driver)
	if mErr != nil {
		return fmt.Errorf("Unable to create migrations instance: %v", mErr)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("Unable to run migrations: %v", err)
	}

	version, _, _ := m.Version()
	log.Printf("migrations completed successfully. version=%v", version)

	return nil
}

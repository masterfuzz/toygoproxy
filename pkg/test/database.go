package test

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/masterfuzz/toygoproxy/pkg/database/migrations"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Database() (pool *pgxpool.Pool, cleanup func(), err error) {
	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14.6",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "postgres",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	cleanup = func() {
		container.Terminate(context.Background())
	}

	if err != nil {
		err = fmt.Errorf("failed to start postgres container: %w", err)
		return
	}

	port, err := container.MappedPort(context.Background(), "5432")
	if err != nil {
		err = fmt.Errorf("failed to get mapped port: %w", err)
		return
	}

	host, err := container.Host(context.Background())
	if err != nil {
		err = fmt.Errorf("failed to get container host: %w", err)
		return
	}

	pool, err = pgxpool.New(context.Background(), "postgres://postgres:postgres@"+host+":"+port.Port()+"/postgres")
	if err != nil {
		err = fmt.Errorf("failed to create pgxpool: %w", err)
		return
	}

	if err = migrations.Run(pool); err != nil {
		err = fmt.Errorf("failed to run migrations: %w", err)
		return
	}

	return
}

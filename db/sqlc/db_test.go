package db

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"os"
	"testing"
	"time"
)

func TestContainer(t *testing.T) {
	_, teardown := StartDBWithTestContainer(t)
	defer teardown()
}

func StartDBWithTestContainer(t *testing.T) (*Queries, func()) {
	t.Helper()

	dsn, terminateContainer, err := startTestContainer()
	require.NoError(t, err)

	db, teardown := initializeDatabase(t, dsn)
	return db, func() {
		teardown()
		terminateContainer()
	}
}

func startTestContainer() (string, func(), error) {
	_ = setColimaEnvVars()
	ctx := context.Background()

	// Disable ryuk until it is fixed in a future release of testcontainers-go
	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		panic(err)
	}
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		return "", nil, fmt.Errorf("starting PostgreSQL container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("getting PostrgreSQL container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return "", nil, fmt.Errorf("getting PostrgreSQL container port: %w", err)
	}

	host = fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		"postgres",
		"password",
		host,
		port.Int(),
		"test",
	)
	terminate := func() {
		if err := container.Terminate(context.Background()); err != nil {
			panic(err)
		}
	}
	return host, terminate, nil
}

func setColimaEnvVars() error {
	err := os.Setenv("TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE", "/var/run/docker .sock")
	if err != nil {
		return fmt.Errorf("setting TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE: %w", err)
	}
	err = os.Setenv("DOCKER_HOST", fmt.Sprintf("unix://%s/.colima/docker.sock", os.Getenv("HOME")))
	if err != nil {
		return fmt.Errorf("setting DOCKER_HOST: %w", err)
	}
	return err
}

func initializeDatabase(t *testing.T, host string) (*Queries, func()) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	conn, err := pgx.Connect(ctx, host)
	require.NoError(t, err)

	migrationPath := "file://../migration"
	m, err := migrate.New(migrationPath, host)

	require.NoError(t, err)
	require.NoError(t, m.Up())

	db := New(conn)

	teardown := func() {
		_ = conn.Close(ctx)
		cancel()
	}

	return db, func() {
		cancel()
		teardown()
	}
}

package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // driver migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // Postgres driver
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"os"
	"testing"
	"time"
)

func TestCreateAccount(t *testing.T) {
	db, teardown := StartDBWithTestContainer(t)
	defer teardown()

	ctx := context.Background()

	// Cria uma tabela simples para teste direto via sql.Exec
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			balance INT NOT NULL DEFAULT 0
		);
	`)
	require.NoError(t, err)

	// Insere um registro
	res, err := db.ExecContext(ctx, `INSERT INTO accounts (name, balance) VALUES ($1, $2)`, "Test User", 1000)
	require.NoError(t, err)

	affected, err := res.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), affected)
}

func TestContainer(t *testing.T) {
	db, teardown := StartDBWithTestContainer(t)
	defer teardown()
	db.QueryRow("SELECT NOW()")
}

func StartDBWithTestContainer(t *testing.T) (*sql.DB, func()) {
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

func initializeDatabase(t *testing.T, host string) (*sql.DB, func()) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	db, err := sql.Open("postgres", host)
	require.NoError(t, err)

	err = db.PingContext(ctx)
	require.NoError(t, err)

	migrationPath := "file://../migration"
	m, err := migrate.New(migrationPath, host)

	require.NoError(t, err)
	require.NoError(t, m.Up())

	teardown := func() {
		_ = db.Close()
		cancel()
	}

	return db, func() {
		cancel()
		teardown()
	}
}

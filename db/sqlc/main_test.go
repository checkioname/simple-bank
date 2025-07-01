package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

var (
	dsn         string
	testStore   Store
	testQueries *Queries
	teardown    func()
)

func TestMain(m *testing.M) {
	testQueries, dsn, teardown = StartDBWithTestContainer(nil)
	defer teardown()

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}

	testStore = NewStore(pool)
	code := m.Run()

	if pool != nil {
		pool.Close()
	}
	if teardown != nil {
		teardown()
	}
	os.Exit(code)
}

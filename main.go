package main

import (
	"context"
	"fmt"
	"github.com/checkioname/simple-bank/api"
	db "github.com/checkioname/simple-bank/db/sqlc"
	"github.com/checkioname/simple-bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connStr = "postgres://root:secret@localhost/simple_bank?sslmode=disable"
	addr    = "0.0.0.0:8080"
)

func main() {
	config, err := util.LoadConfig()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, config.ConnStr)
	if err != nil {
		fmt.Errorf("db conn: %v", err)
		return
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(addr)
	if err != nil {
		fmt.Errorf("server start: %v", err)
	}
}

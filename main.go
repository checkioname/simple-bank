package main

import (
	"context"
	"fmt"
	"github.com/checkioname/simple-bank/api"
	db "github.com/checkioname/simple-bank/db/sqlc"
	"github.com/checkioname/simple-bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
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

	err = server.Start(config.Addr)
	if err != nil {
		fmt.Errorf("server start: %v", err)
	}
}

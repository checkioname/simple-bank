package main

import (
	"context"
	"log"

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
		log.Printf("db conn: %v", err)
		return
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		panic(err)
	}

	err = server.Start(config.Addr)
	if err != nil {
		log.Printf("server start: %v", err)
	}
}

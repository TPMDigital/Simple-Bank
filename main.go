package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"

	"github.com/tpmdigital/simplebank/api"
	db "github.com/tpmdigital/simplebank/db/sqlc"
	"github.com/tpmdigital/simplebank/util"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configuration file:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

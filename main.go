package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/DenysBahachuk/Simple_Bank/api"
	db "github.com/DenysBahachuk/Simple_Bank/db/sqlc"
	"github.com/DenysBahachuk/Simple_Bank/utils"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(cfg.DBdriver, cfg.DBsource)
	if err != nil {
		log.Fatal("unable to connect to db:", err)
	}

	fmt.Println("successfully connected to db:", cfg.DBdriver)

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot start the server:", err)
	}
}

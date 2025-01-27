package main

import (
	"database/sql"

	"github.com/DenysBahachuk/Simple_Bank/api"
	db "github.com/DenysBahachuk/Simple_Bank/db/sqlc"
	"github.com/DenysBahachuk/Simple_Bank/utils"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	cfg, err := utils.LoadConfig(".")
	if err != nil {
		logger.Fatal("cannot load config:", err)
	}
	logger.Info("config successfully loaded")

	conn, err := sql.Open(cfg.DBdriver, cfg.DBsource)
	if err != nil {
		logger.Fatal("unable to connect to db:", err)
	}

	logger.Info("connection to db established:", cfg.DBdriver)

	store := db.NewStore(conn)
	server := api.NewServer(store, logger)

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		logger.Fatal("cannot start the server:", err)
	}
}

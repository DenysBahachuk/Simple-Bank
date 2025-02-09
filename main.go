package main

import (
	"database/sql"
	"net"

	"github.com/DenysBahachuk/Simple_Bank/api"
	db "github.com/DenysBahachuk/Simple_Bank/db/sqlc"
	"github.com/DenysBahachuk/Simple_Bank/gapi"
	"github.com/DenysBahachuk/Simple_Bank/pb"
	"github.com/DenysBahachuk/Simple_Bank/utils"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	//runGinServer(store, logger, cfg)

	runGrpcServer(store, logger, cfg)

}

func runGinServer(store db.Store, logger *zap.SugaredLogger, cfg utils.Config) {
	server, err := api.NewServer(store, logger, cfg)
	if err != nil {
		logger.Fatalf("failed to create gin server: %w", err)
	}

	err = server.Start(cfg.HTTPServerAddress)
	if err != nil {
		logger.Fatal("cannot start the gin server: %w", err)
	}
}

func runGrpcServer(store db.Store, logger *zap.SugaredLogger, cfg utils.Config) {
	grpcServer := grpc.NewServer()

	server, err := gapi.NewServer(store, cfg)
	if err != nil {
		logger.Fatalf("failed to create gRPC server: %w", err)
	}

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		logger.Fatalf("failed to create gRPC listener: %w", err)
	}

	logger.Info("start gRPC server at: ", cfg.GRPCServerAddress)
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Fatalf("failed to start gRPC server: %w", err)
	}
}

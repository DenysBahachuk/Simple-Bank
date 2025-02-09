package gapi

import (
	"fmt"

	db "github.com/DenysBahachuk/Simple_Bank/db/sqlc"
	"github.com/DenysBahachuk/Simple_Bank/pb"
	"github.com/DenysBahachuk/Simple_Bank/token"
	"github.com/DenysBahachuk/Simple_Bank/utils"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config utils.Config
	store  db.Store
	//logger     *zap.SugaredLogger
	tokenMaker token.Maker
}

func NewServer(store db.Store, config utils.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create a token maker: %w", err)
	}

	server := Server{
		config: config,
		store:  store,
		//logger:     logger,
		tokenMaker: tokenMaker,
	}

	return &server, nil
}

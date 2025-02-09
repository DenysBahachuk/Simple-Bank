package gapi

import (
	"context"

	db "github.com/DenysBahachuk/Simple_Bank/db/sqlc"
	"github.com/DenysBahachuk/Simple_Bank/pb"
	"github.com/DenysBahachuk/Simple_Bank/utils"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	args := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := s.store.CreateUser(ctx, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user already exists: %v", err)
			}
		}

		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	resp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return resp, nil
}

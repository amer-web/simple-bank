package gapi

import (
	"context"
	"github.com/amer-web/simple-bank/config"
	db "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/amer-web/simple-bank/helper"
	"github.com/amer-web/simple-bank/pb"
	tok "github.com/amer-web/simple-bank/token"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hash, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		FullName: req.FullName,
		Email:    req.Email,
		Password: hash,
	})
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				return nil, status.Error(codes.AlreadyExists, err.Error())
			}
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	handleToken := tok.NewMakerToken()
	token, _ := handleToken.CreateToken(user.Username, config.Source.TOKENDURATION)
	return &pb.CreateUserResponse{User: convertUser(user),
		AccessToken: token,
	}, nil
}

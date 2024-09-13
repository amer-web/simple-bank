package gapi

import (
	"context"
	"database/sql"
	"github.com/amer-web/simple-bank/config"
	"github.com/amer-web/simple-bank/helper"
	"github.com/amer-web/simple-bank/pb"
	tok "github.com/amer-web/simple-bank/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	// Validate request
	if err := req.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", err)
	}

	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = helper.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	handleToken := tok.NewMakerToken()
	token, _ := handleToken.CreateToken(user.Username, config.Source.TOKENDURATION)

	return &pb.LoginUserResponse{User: convertUser(user),
		AccessToken: token,
	}, nil
}

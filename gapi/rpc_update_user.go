package gapi

import (
	"context"
	"database/sql"
	db "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/amer-web/simple-bank/helper"
	"github.com/amer-web/simple-bank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	hashPassword := ""
	if req.GetPassword() != "" {
		hashPassword, _ = helper.HashPassword(req.GetPassword())
	}
	user, err := s.store.UpdateUser(ctx, db.UpdateUserParams{
		Username: req.Username,
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.GetFullName() != "",
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.GetEmail() != "",
		},
		Password: sql.NullString{
			String: hashPassword,
			Valid:  hashPassword != "",
		},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateUserResponse{
		User: convertUser(user),
	}, nil
}

package gapi

import (
	"context"
	"errors"
	db "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/amer-web/simple-bank/helper"
	"github.com/amer-web/simple-bank/pb"
	"github.com/jackc/pgx/v5/pgtype"
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
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid:  req.GetFullName() != "",
		},
		Email: pgtype.Text{
			String: req.GetEmail(),
			Valid:  req.GetEmail() != "",
		},
		Password: pgtype.Text{
			String: hashPassword,
			Valid:  hashPassword != "",
		},
	})
	if err != nil {
		if errors.Is(err, db.ErrorRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateUserResponse{
		User: convertUser(user),
	}, nil
}

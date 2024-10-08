package gapi

import (
	"context"
	"github.com/amer-web/simple-bank/pb"
	token2 "github.com/amer-web/simple-bank/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) VerifyEmailUser(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", err)
	}
	token := token2.NewMakerToken()
	payload, err := token.VerifyToken(req.Token)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid token: %v", err)
	}
	user, err := s.store.GetUser(ctx, payload.Sub)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &pb.VerifyEmailResponse{User: convertUser(user),
		Token: req.Token,
	}, nil
}

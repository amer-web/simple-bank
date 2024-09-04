package gapi

import (
	db "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/amer-web/simple-bank/pb"
)

type Server struct {
	store db.Store
	pb.UnimplementedSimpleBankServer
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	return server
}

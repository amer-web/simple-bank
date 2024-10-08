package gapi

import (
	db "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/amer-web/simple-bank/jobs"
	"github.com/amer-web/simple-bank/pb"
)

type Server struct {
	store db.Store
	tasks jobs.Distributor
	pb.UnimplementedSimpleBankServer
}

func NewServer(store db.Store, tasks jobs.Distributor) *Server {
	server := &Server{
		store: store,
		tasks: tasks,
	}
	return server
}

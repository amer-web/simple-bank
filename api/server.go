package api

import (
	db "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/amer-web/simple-bank/helper"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", helper.CurrencyValidate)
	}

	server.router = router
	return server
}

func (s *Server) Run() {
	s.routeAccount(s.router)
	s.routeTransfer(s.router)
	s.routeUser(s.router)
	s.router.Run(":8081")
}

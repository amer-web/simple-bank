package api

import (
	db "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/amer-web/simple-bank/helper"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

func (s *Server) routeUser(route *gin.Engine) {
	route.POST("/user", s.createUser)

}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Server) createUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, _ := helper.HashPassword(req.Password)
	user, err := s.store.CreateUser(c, db.CreateUserParams{
		Username: req.Username,
		FullName: req.FullName,
		Email:    req.Email,
		Password: hash,
	})
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

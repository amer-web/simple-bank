package api

import (
	"database/sql"
	"github.com/amer-web/simple-bank/config"
	db "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/amer-web/simple-bank/helper"
	tok "github.com/amer-web/simple-bank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"time"
)

func (s *Server) routeUser(route *gin.Engine) {
	route.POST("/user", s.createUser)
	route.POST("/login", s.loginUser)
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

type CreateUserLoginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type ResponseLoginUser struct {
	User        ResponseUser `json:"user" `
	AccessToken string       `json:"access_token"`
}
type ResponseUser struct {
	Username          string       `json:"username"`
	FullName          string       `json:"full_name"`
	Email             string       `json:"email"`
	PasswordChangedAt sql.NullTime `json:"password_changed_at"`
	CreatedAt         time.Time    `json:"created_at"`
}

func (s *Server) loginUser(c *gin.Context) {
	var req CreateUserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := s.store.GetUser(c, req.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = helper.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password not valid"})
		return
	}
	handleToken := tok.NewMakerToken()
	token, _ := handleToken.CreateToken(user.Username, config.Source.TOKENDURATION)
	response := ResponseLoginUser{
		AccessToken: token,
		User: ResponseUser{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: user.PasswordChangedAt,
			CreatedAt:         user.CreatedAt,
		},
	}
	c.JSON(http.StatusOK, response)
}

package api

import (
	"errors"
	db "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (s *Server) routeAccount(route *gin.Engine) {
	route.POST("/account", s.createAccount)
	route.GET("/account/:id", s.getAccount)
	route.GET("/accounts", s.getAccounts)
}

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := s.store.CreateAccount(c, db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
	})
	if err != nil {
		switch db.ErrorCode(err) {
		case "unique_violation", "foreign_key_violation":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"account": account})
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *Server) getAccount(c *gin.Context) {
	var req GetAccountRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := s.store.GetAccount(c, req.ID)
	if err != nil {
		if errors.Is(err, db.ErrorRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (s *Server) getAccounts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
		return
	}
	offset := (page - 1) * limit
	accounts, _ := s.store.ListAccounts(c, db.ListAccountsParams{
		Limit:  int64(limit),
		Offset: int64(offset),
	})
	c.JSON(http.StatusOK, gin.H{"accounts": accounts})

}

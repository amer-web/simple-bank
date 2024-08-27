package api

import (
	db "github.com/amer-web/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) routeTransfer(route *gin.Engine) {
	route.POST("/transfer", s.CreateTransferRequest)

}

type CreateTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) CreateTransferRequest(c *gin.Context) {
	var req CreateTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !s.validAccountCurrency(c, req.FromAccountID, req.Currency) {
		return
	}
	if !s.validAccountCurrency(c, req.ToAccountID, req.Currency) {
		return
	}
	result, err := s.store.TransferTx(c, db.ArrgTransfer{
		FromAcc: req.FromAccountID,
		ToAcc:   req.ToAccountID,
		Amount:  req.Amount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}
func (s *Server) validAccountCurrency(c *gin.Context, accId int64, currency string) bool {
	account, err := s.store.GetAccount(c, accId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	if account.Currency != currency {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account currency"})
		return false
	}
	return true
}

package api

import (
	"fmt"
	db "github.com/checkioname/simple-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccount(c *gin.Context) {
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Errorf("createAccount: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse body"})
		return
	}
	args := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
	}

	result, err := s.store.CreateAccount(c, args)
	if err != nil {
		fmt.Errorf("createAccount: %v", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(c *gin.Context) {
	var req getAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Errorf("getAccount: %v", err)
		badRequest(c, "Could not parse body")
	}
	acc, err := s.store.GetAccount(c, req.ID)
	if err != nil {
		fmt.Errorf("getAccount: %v", err)
		badRequest(c, "Account not found")
		return
	}

	c.JSON(http.StatusOK, acc)
	return
}

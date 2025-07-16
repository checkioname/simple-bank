package api

import (
	db "github.com/checkioname/simple-bank/db/sqlc"
	"github.com/checkioname/simple-bank/util"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccount(c *gin.Context) {
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("createAccount:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse body"})
		return
	}
	args := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
	}

	result, err := s.store.CreateAccount(c, args)
	if err != nil {
		slog.Error("createAccount:", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (s *Server) createUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("createUser:", err)
		return
	}

	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		slog.Error("createUser:", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	args := db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hashed,
	}

	result, err := s.store.CreateUser(c, args)
	if err != nil {
		slog.Error("createUser:", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)

}

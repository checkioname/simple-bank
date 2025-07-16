package api

import (
	"database/sql"
	"fmt"
	db "github.com/checkioname/simple-bank/db/sqlc"
	"github.com/checkioname/simple-bank/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	User        string `json:"user"`
}

func (server *Server) loginUser(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := server.store.GetUser(c, req.Username)
	if err != nil {
		fmt.Errorf("loginUser: %v", err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = util.VerifyPassword(req.Password, user.HashedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := server.token.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loginResponse{accessToken, user.Username})
	return
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
		fmt.Errorf("createUser: %v", err)
		return
	}

	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		fmt.Errorf("createUser: %v", err)
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
		fmt.Errorf("createUser: %v", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)

}

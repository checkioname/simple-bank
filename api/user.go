package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/checkioname/simple-bank/db/sqlc"
	"github.com/checkioname/simple-bank/util"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	SessionID             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User                  string    `json:"user"`
}

func (s *Server) loginUser(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, http.StatusBadRequest, err)
		return
	}

	user, err := s.store.GetUser(c, req.Username)
	if err != nil {
		log.Printf("loginUser: %v", err)
		if err == sql.ErrNoRows {
			errResponse(c, http.StatusNotFound, err)
			return
		}
		errResponse(c, http.StatusInternalServerError, err)
		return
	}

	err = util.VerifyPassword(req.Password, user.HashedPassword)
	if err != nil {
		errResponse(c, http.StatusUnauthorized, err)
		return
	}

	accessToken, payload, err := s.token.CreateToken(user.Username, s.config.AccessTokenDuration)
	fmt.Println(payload)
	if err != nil {
		errResponse(c, http.StatusInternalServerError, err)
		return
	}

	refreshToken, refreshPayload, err := s.token.CreateToken(user.Username, s.config.AccessTokenDuration)
	fmt.Println(payload)
	if err != nil {
		errResponse(c, http.StatusInternalServerError, err)
		return
	}

	sessionParams := db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	}
	session, err := s.store.CreateSession(c.Request.Context(), sessionParams)
	if err != nil {
		errResponse(c, http.StatusInternalServerError, err)
		return
	}

	rsp := loginResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  payload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  user.Username,
	}

	c.JSON(http.StatusOK, rsp)
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
		log.Printf("createUser: %v", err)
		return
	}

	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		log.Printf("createUser: %v", err)
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
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			errResponse(c, http.StatusBadRequest, ErrUserAlreadyExists)
			return
		}
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, result)

}

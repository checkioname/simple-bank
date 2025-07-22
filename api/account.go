package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	db "github.com/checkioname/simple-bank/db/sqlc"
	"github.com/checkioname/simple-bank/token"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrAccountNotFound         = fmt.Errorf("account not found")
	ErrAccountAlreadyExists    = fmt.Errorf("account already exists")
	ErrAccountBelongsToInvalid = fmt.Errorf("account belongs to invalid")
	ErrInvalidAccount          = fmt.Errorf("invalid account")
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccount(c *gin.Context) {
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("createAccount: %v", err)
		errResponse(c, http.StatusBadRequest, ErrInvalidAccount)
	}

	authPayload := c.MustGet(authPayloadKey).(*token.Payload)
	args := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
	}

	result, err := s.store.CreateAccount(c.Request.Context(), args)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		fmt.Println(pgErr)
		if ok && pgErr.Code == "23505" {
			errResponse(c, http.StatusBadRequest, ErrAccountAlreadyExists)
			return
		}
		errResponse(c, http.StatusInternalServerError, fmt.Errorf("internal error"))
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
		log.Printf("getAccount: %v", err)
		errResponse(c, http.StatusBadRequest, ErrCouldNotParse)
		return
	}
	acc, err := s.store.GetAccount(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("getAccount: %v", err)
			errResponse(c, http.StatusNotFound, ErrAccountNotFound)
			return
		}
		errResponse(c, http.StatusInternalServerError, err)
		return
	}

	authPayload := c.MustGet("authPayload").(*token.Payload)
	if authPayload.Username != acc.Owner {
		errResponse(c, http.StatusBadRequest, ErrAccountBelongsToInvalid)
		return
	}

	c.JSON(http.StatusOK, acc)
}

type listAccountsRequest struct {
	Owner  string `uri:"owner" binding:"required,min=1"`
	Limit  int32  `uri:"limit" binding:"required,min=1,max=100"`
	Offset int32  `uri:"offset" binding:"required,min=1,max=100"`
}

func (s *Server) listAccounts(c *gin.Context) {
	var req listAccountsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("listAccounts: %v", err)
		errResponse(c, http.StatusBadRequest, ErrInvalidAccount)
		return
	}

	authPayload := c.MustGet("authPayload").(*token.Payload)

	args := db.ListAccountsParams(req)
	args.Owner = authPayload.Username
	accounts, err := s.store.ListAccounts(c.Request.Context(), args)
	if err != nil {
		log.Printf("listAccounts: %v", err)
		errResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, accounts)
}

package api

import (
	"fmt"
	db "github.com/checkioname/simple-bank/db/sqlc"
	"github.com/checkioname/simple-bank/token"
	"github.com/checkioname/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config util.Config
	store  db.Store
	token  token.Maker
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("NewServer: %w", err)
	}

	server := &Server{config: config,
		store: store,
		token: tokenMaker,
	}

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.router = server.setupRoutes()
	return server, nil
}

func (server *Server) setupRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(authMiddleware(server.token))

	router.POST("/accounts", server.createAccount)
	router.GET("/account/:id", server.getAccount)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	return router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

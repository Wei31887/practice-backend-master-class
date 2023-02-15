package api

import (
	db "lesson/simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{ store: store }
	router := gin.Default()

	router.POST("users", server.CreateUser)
	router.GET("users", server.GetUser)

	router.POST("accounts", server.CreateAccount)
	router.GET("accounts/:id", server.GetAccount)
	router.GET("accounts", server.ListAccount)

	router.POST("transfers", server.CreateTransfer)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
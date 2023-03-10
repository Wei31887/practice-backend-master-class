package api

import (
	"lesson/simple-bank/config"
	db "lesson/simple-bank/db/sqlc"
	"lesson/simple-bank/token"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config     config.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

func NewServer(config config.Config, store db.Store) (*Server, error) {
	maker, err := token.NewJWTMaker(config.SecreteKey)
	if err != nil {
		log.Fatal("Can't make token maker ", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: maker,
	}

	server.setRouterGroup()

	return server, nil
}

func (server *Server) setRouterGroup() {
	router := gin.Default()

	router.POST("users", server.CreateUser)
	router.GET("users", server.GetUser)

	router.POST("accounts", server.CreateAccount)
	router.GET("accounts/:id", server.GetAccount)
	router.GET("accounts", server.ListAccount)

	router.POST("transfers", server.CreateTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

package gapi

import (
	"lesson/simple-bank/config"
	db "lesson/simple-bank/db/sqlc"
	pb "lesson/simple-bank/pb"
	"lesson/simple-bank/token"
	"log"
)

// server for gRPC service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     config.Config
	tokenMaker token.Maker
	store      db.Store
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

	return server, nil
}
package main

import (
	"database/sql"
	"lesson/simple-bank/api"
	"lesson/simple-bank/config"
	db "lesson/simple-bank/db/sqlc"
	"lesson/simple-bank/gapi"
	"lesson/simple-bank/initial"
	simple_bank "lesson/simple-bank/pb"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := initial.LoadingConfig(".")
	if err != nil {
		log.Fatal("Can't load config, ", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	go runGrpcServer(config, store)
	runGinServer(config, store)
}

func runGrpcServer(config config.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err!= nil {
        log.Fatal("Can't create server, ", err)
    }

	grpcServer := grpc.NewServer()
	simple_bank.RegisterSimpleBankServer(grpcServer, server)

	// reflection service
	reflection.Register(grpcServer)

	// listener
	lietener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
        log.Fatal("Can't listen ", err)
    }
	log.Println("Start grpc server, ", lietener.Addr())

	err = grpcServer.Serve(lietener)
	if err!= nil {
        log.Fatal("Can't serve Grpc ", err)
    }
}

func runGinServer(config config.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Can't create server, ", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Can't start server: ", err)
	}
}

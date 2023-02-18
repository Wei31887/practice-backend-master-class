package main

import (
	"database/sql"
	"lesson/simple-bank/api"
	db "lesson/simple-bank/db/sqlc"
	"lesson/simple-bank/initial"
	"log"

	_ "github.com/lib/pq"
)

func main() {
  config, err := initial.LoadingConfig(".")
  if err != nil {
    log.Fatal("Can't load config, ", err)
  }

  conn, err := sql.Open(config.DbDriver, config.DbSource)
  if err!= nil {
    panic(err)
  }
  defer conn.Close()

  store := db.NewStore(conn)
  server, err := api.NewServer(config, store)
  if err!= nil {
    log.Fatal("Can't create server, ", err)
  }

  err = server.Start(config.ServerAddress)
  if err!= nil {
    log.Fatal("Can't start server: ", err)
  }
}
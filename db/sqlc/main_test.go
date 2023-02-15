package db

import (
	"database/sql"
	"lesson/simple-bank/initial"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	path := "../../"
	config, err := initial.LoadingConfig(path)
	if err != nil {
		log.Fatal("Cann't load config", err)
	}

	testDB, err = sql.Open(config.DbDriver, config.DbSource)	
	if err != nil {
		log.Fatal("Cannot connect to db ,", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

const (
	driver = "postgres"
	source = "postgresql://amer:amer@127.0.0.1:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var dbTest *sql.DB

func TestMain(m *testing.M) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Fatal("error opening db:", err.Error())
	}
	defer db.Close()
	testQueries = New(db)
	dbTest = db
	os.Exit(m.Run())
}

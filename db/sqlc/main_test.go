package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

const (
	driver = "postgres"
	source = "postgresql://amer:amer@127.0.0.1:5432/simple_bank?sslmode=disable"
)

var testStore Store

func TestMain(m *testing.M) {
	connection, err := pgxpool.New(context.Background(), source)
	if err != nil {
		log.Fatal("error opening db:", err.Error())
	}
	testStore = NewStore(connection)
	os.Exit(m.Run())
}

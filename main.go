package main

import (
	"database/sql"
	"fmt"
	"github.com/amer-web/simple-bank/api"
	"github.com/amer-web/simple-bank/config"
	db2 "github.com/amer-web/simple-bank/db/sqlc"
	_ "github.com/lib/pq"
	"log"
)

func init() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	source := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.Source.DBUser, config.Source.DBPassword, config.Source.DBHost, config.Source.DBPort, config.Source.DBName)
	db, err := sql.Open(config.Source.DRIVER, source)
	if err != nil {
		log.Fatal("error opening db:", err.Error())
	}
	defer db.Close()
	store := db2.NewStore(db)
	server := api.NewServer(store)
	server.Run()
}

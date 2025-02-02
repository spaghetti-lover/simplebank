package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/spaghetti-lover/bank-system/api"
	db "github.com/spaghetti-lover/bank-system/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:root@localhost:5432/bank-system?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
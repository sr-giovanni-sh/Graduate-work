package main

import (
	"Graduate-work/pkg/db"
	"Graduate-work/pkg/server"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	dbFile := "scheduler.db"
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.CloseDB()
	server.StartServer()
}

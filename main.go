package main

import (
	"Graduate-work/pkg/api"
	"Graduate-work/pkg/config"
	"Graduate-work/pkg/db"
	"Graduate-work/pkg/server"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()

	if err := db.Init(cfg.DBFile); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	apiHandler := api.NewHandler(cfg.Password)
	if err := server.StartServer(cfg.Port, apiHandler); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

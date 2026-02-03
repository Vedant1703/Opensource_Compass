package main

import (
	"log"
	"net/http"

	"core-service/config"
	"core-service/internal/db"
	"core-service/routes"
)

func main() {
	cfg := config.Load()

	if cfg.DatabaseURL != "" {
		_, err := db.Connect(cfg.DatabaseURL)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
	}

	mux := http.NewServeMux()
	routes.RegisterHealthRoutes(mux)

	log.Printf("Core service running on :%s\n", cfg.ServerPort)
	err := http.ListenAndServe(":"+cfg.ServerPort, mux)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"net/http"

	"core-service/config"
	"core-service/internal/db"
	"core-service/internal/orchestration"
	"core-service/internal/preferences"
	"core-service/routes"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	var dbPool *pgxpool.Pool
	var err error

	if cfg.DatabaseURL != "" {
		dbPool, err = db.Connect(cfg.DatabaseURL)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
	}

	mux := http.NewServeMux()
	routes.RegisterHealthRoutes(mux)

	// initialize dependencies (even if dbPool is nil for now)
	prefRepo := preferences.NewRepository(dbPool)
	orchService := orchestration.NewService(prefRepo)

	routes.RegisterRoutes(mux, orchService)

	log.Printf("Core service running on :%s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, mux))
}

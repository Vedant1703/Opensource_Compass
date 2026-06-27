package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"core-service/config"
	"core-service/internal/auth"
	"core-service/internal/clients"
	"core-service/internal/db"
	"core-service/internal/orchestration"
	"core-service/internal/preferences"
	"core-service/internal/users"
	"core-service/internal/watchlist"
	"core-service/internal/middleware"
	"core-service/routes"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)



func main() {
	// Load the root .env file from the project root
	// Load the root .env file from the project root (optional, for local dev)
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load() // Fallback to local .env
	cfg := config.Load()
	aiClient := clients.NewAIClient(cfg.AISvcURL)
	githubClient := clients.NewGitHubClient(cfg.GitHubSvcURL)

	var dbPool *pgxpool.Pool
	var err error

	if cfg.DatabaseURL != "" {
		dbPool, err = db.Connect(cfg.DatabaseURL)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		// Run Migrations
		if err := db.RunMigrations(context.Background(), dbPool); err != nil {
			log.Printf("Migration warning: %v", err)
			// Don't fatal here in case of minor errors, or if user wants to run without it?
			// Actually, for missing column, we need it. But let's log and proceed.
		}
	}

	// Debug: Print JWT secret (first 8 chars only)
	log.Printf("🔑 Core Service JWT_SECRET: %s... (len: %d)", cfg.JWTSecret[:8], len(cfg.JWTSecret))

	mux := http.NewServeMux()
	routes.RegisterHealthRoutes(mux)

	prefRepo := preferences.NewRepository(dbPool)
	userRepo := users.NewRepository(dbPool)
	watchlistRepo := watchlist.NewRepository(dbPool)

	// Initialize Notification Service Notifier (WebSocket)
	wsNotifier := watchlist.NewNotifier(cfg.NotificationSvcURL)

	// Initialize Email Notifier
	emailNotifier := watchlist.NewEmailNotifier(userRepo)

	// Initialize Poller (runs every 2 minutes) with both notifiers
	poller := watchlist.NewPoller(watchlistRepo, githubClient, []watchlist.Notifier{wsNotifier, emailNotifier}, 2*time.Minute)

	// Start Poller in background
	ctx := context.Background() // TODO: Use proper context with cancellation on shutdown
	go poller.Start(ctx)

	orchService := orchestration.NewService(prefRepo, aiClient, githubClient)

	routes.RegisterRoutes(mux, orchService, cfg.JWTSecret, prefRepo, userRepo, githubClient, aiClient)
	routes.RegisterWatchlistRoutes(mux, watchlistRepo, cfg.JWTSecret)

	// Wrap mux with CORS middleware so browser preflight OPTIONS requests are handled.
	handler := auth.CORSMiddleware(mux)
	
	// Add RequestID as the outermost middleware
	handler = middleware.RequestID(handler)

	log.Printf("Core service running on :%s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, handler))
}

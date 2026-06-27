package main

import (
	"auth-service/routes"
	"crypto/rand"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// corsMiddleware adds CORS headers and handles preflight OPTIONS requests.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}

		c.Header("Access-Control-Allow-Origin", frontendURL)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// requestIDMiddleware adds an X-Request-ID header to every request.
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			b := make([]byte, 16)
			rand.Read(b)
			reqID = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
		}
		
		c.Set("request_id", reqID)
		
		// For standard context propagation
		type contextKey string
		ctx := context.WithValue(c.Request.Context(), contextKey("request_id"), reqID)
		c.Request = c.Request.WithContext(ctx)

		c.Header("X-Request-ID", reqID)
		c.Next()
	}
}

func main() {
	// Load the root .env file from the project root (optional, for local dev)
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load() // Fallback to local .env

	required := []string{
		"GITHUB_CLIENT_ID",
		"GITHUB_CLIENT_SECRET",
		"JWT_SECRET",
		"FRONTEND_URL",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			log.Fatalf("Missing required env variable: %s", key)
		}
	}

	// Debug: Print JWT secret (first 8 chars only)
	jwtSecret := os.Getenv("JWT_SECRET")
	log.Printf("🔑 Auth Service JWT_SECRET: %s... (len: %d)", jwtSecret[:8], len(jwtSecret))

	router := gin.Default()

	// Apply CORS middleware globally so preflight OPTIONS requests are handled
	router.Use(corsMiddleware())
	router.Use(requestIDMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "auth-service alive",
		})
	})

	routes.RegisterAuthRoutes(router)

	log.Println("Auth service running on :8080")
	router.Run(":8080")
}

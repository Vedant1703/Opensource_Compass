package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"notification-service/internal/websocket"

	"github.com/joho/godotenv"
)


// corsMiddleware adds CORS headers and handles preflight OPTIONS requests.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}

		w.Header().Set("Access-Control-Allow-Origin", frontendURL)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load the root .env file from the project root
	// Load the root .env file from the project root (optional, for local dev)
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load() // Fallback to local .env

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084" // Default port for notification service
	}

	hub := websocket.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	// Internal endpoint to send specific notifications
	// POST /notify?user_id=123
	// Body: JSON payload
	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Security check: Should be protected by some internal secret
		// secret := r.Header.Get("X-Internal-Secret")
		// if secret != os.Getenv("INTERNAL_SECRET") { ... }

		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "user_id query param required", http.StatusBadRequest)
			return
		}

		log.Printf("Received notify request for user: %s", userID)

		var payload interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Printf("Error decoding payload: %v", err)
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}

		msgBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshaling payload: %v", err)
			http.Error(w, "Failed to marshal payload", http.StatusInternalServerError)
			return
		}

		hub.SendToUser <- websocket.NotificationMessage{
			UserID:  userID,
			Message: msgBytes,
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Notification sent"))
	})

	log.Printf("Notification Service running on :%s", port)
	if err := http.ListenAndServe(":"+port, corsMiddleware(http.DefaultServeMux)); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

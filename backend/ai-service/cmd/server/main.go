package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/your-org/opensource-compass/ai-service/config"
	"github.com/your-org/opensource-compass/ai-service/internal/analysis"
	"github.com/your-org/opensource-compass/ai-service/internal/llm"
	"github.com/your-org/opensource-compass/ai-service/routes"
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

	cfg := config.Load()

	// Initialize LLM client based on provider
	var llmClient llm.Client
	if cfg.LLMProvider == "groq" {
		log.Println("Using Groq LLM provider")
		llmClient = llm.NewGroqClient(cfg.Model)
	} else {
		log.Println("Using Gemini LLM provider")
		llmClient = llm.NewGeminiClient(cfg.Model)
	}

	analyzer := analysis.NewAnalyzer(llmClient)
	setupGenerator := analysis.NewSetupGenerator(llmClient)

	http.HandleFunc("/analyze-issue", routes.AnalyzeIssue(analyzer))
	http.HandleFunc("/generate-setup-guide", routes.GenerateSetupGuide(setupGenerator))

	log.Println("AI Service running on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, corsMiddleware(http.DefaultServeMux)))
}

package main

import (
	"log"
	"net/http"

	"github.com/your-org/opensource-compass/ai-service/config"
	"github.com/your-org/opensource-compass/ai-service/internal/analysis"
	"github.com/your-org/opensource-compass/ai-service/internal/llm"
	"github.com/your-org/opensource-compass/ai-service/routes"
)

func main() {
	cfg := config.Load()

	llmClient := llm.NewGeminiClient(cfg.Model)
	analyzer := analysis.NewAnalyzer(llmClient)
	

	http.HandleFunc("/analyze-issue", routes.AnalyzeIssue(analyzer))

	log.Println("AI Service running on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}

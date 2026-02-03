package routes

import (
	"net/http"

	"core-service/internal/orchestration"
)

func RegisterRoutes(
	mux *http.ServeMux,
	orchService *orchestration.Service,
) {
	handler := NewRecommendationHandler(orchService)

	mux.HandleFunc("/recommendations", handler.GetRecommendations)
}

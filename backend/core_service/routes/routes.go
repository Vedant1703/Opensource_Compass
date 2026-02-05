package routes

import (
	"net/http"

	"core-service/internal/orchestration"
)

func RegisterRoutes(
	mux *http.ServeMux,
	orchService *orchestration.Service,
	jwtSecret string,
) {
	handler := NewRecommendationHandler(orchService, jwtSecret)

	mux.HandleFunc("/recommendations", handler.GetRecommendations)
}

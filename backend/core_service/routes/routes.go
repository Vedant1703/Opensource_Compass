package routes

import (
	"net/http"

	"core-service/internal/auth"
	"core-service/internal/orchestration"
)

func RegisterRoutes(
	mux *http.ServeMux,
	orchService *orchestration.Service,
	jwtSecret string,
) {
	handler := NewRecommendationHandler(orchService, jwtSecret)

	mux.Handle("/recommendations", auth.JWTMiddleware(jwtSecret, http.HandlerFunc(handler.GetRecommendations)))
}

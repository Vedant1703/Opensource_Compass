package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"core-service/internal/orchestration"
)

type RecommendationHandler struct {
	service *orchestration.Service
}

func NewRecommendationHandler(service *orchestration.Service) *RecommendationHandler {
	return &RecommendationHandler{service: service}
}

func (h *RecommendationHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id query parameter", http.StatusBadRequest)
		return
	}
	repoName := r.URL.Query().Get("repo")
	if repoName == "" {
		http.Error(w, "missing repo query parameter", http.StatusBadRequest)
		return
	}

	repoID := "repo-xyz"

	userCtx, err := h.service.BuildUserContext(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	signals := orchestration.BuildMockSignals(80)

	rec, err := h.service.ScoreRepositoryForUser(ctx, userCtx, repoID, signals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	issues := orchestration.BuildMockIssues()

	issues, err = h.service.FetchAndEnrichIssues(ctx, repoName, 3)
	if err != nil {
		log.Printf("failed to fetch issues: %v", err)
	}

	resp := orchestration.RecommendationResponse{
		RepoID:  rec.RepoID,
		Score:   rec.Score.TotalScore,
		Level:   rec.Score.Level,
		Reasons: rec.Score.Reasons,
		Issues:  issues,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

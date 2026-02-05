package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"core-service/internal/auth"
	"core-service/internal/orchestration"
)

type RecommendationHandler struct {
	service   *orchestration.Service
	jwtSecret string
}

func NewRecommendationHandler(service *orchestration.Service, jwtSecret string) *RecommendationHandler {
	return &RecommendationHandler{
		service:   service,
		jwtSecret: jwtSecret,
	}
}

func (h *RecommendationHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	token := ""
	if authHeader != "" {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	}

	userID, err := auth.ExtractUserID(token, h.jwtSecret)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userCtx, err := h.service.BuildUserContext(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// repo should be in "owner/repo" format (e.g., nodejs/node)
	repos, err := h.service.SearchReposForUser(ctx, userCtx, token)
	if err != nil || len(repos) == 0 {
		http.Error(w, "no repositories found", http.StatusInternalServerError)
		return
	}

	topRepo := "nodejs/" + repos[0]

	signals := h.service.BuildRepoSignals(ctx, topRepo, userCtx)

	rec, err := h.service.ScoreRepositoryForUser(ctx, userCtx, topRepo, signals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	issues, err := h.service.FetchAndEnrichIssues(ctx, topRepo, 3, token)
	if err != nil {
		http.Error(w, "failed to fetch issues from GitHub service", http.StatusInternalServerError)
		return
	}

	resp := orchestration.RecommendationResponse{
		RepoID:         rec.RepoID,
		Score:          rec.Score.TotalScore,
		Level:          rec.Score.Level,
		Reasons:        rec.Score.Reasons,
		SuggestedRepos: repos,
		Issues:         issues,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

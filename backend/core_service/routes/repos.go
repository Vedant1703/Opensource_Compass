package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"core-service/internal/auth"
	"core-service/internal/orchestration"
	"core-service/internal/users"
)

type RepoHandler struct {
	service   *orchestration.Service
	jwtSecret string
	userRepo  *users.Repository
}

func NewRepoHandler(service *orchestration.Service, jwtSecret string, userRepo *users.Repository) *RepoHandler {
	return &RepoHandler{
		service:   service,
		jwtSecret: jwtSecret,
		userRepo:  userRepo,
	}
}

func (h *RepoHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	token := ""
	if authHeader != "" {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	}

	userID, err := auth.ExtractUserID(token, h.jwtSecret)
	if err != nil {
		log.Printf("HandleSearch: Auth failed: %v", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch user's GitHub token from database
	githubToken, err := h.userRepo.GetGitHubToken(ctx, userID)
	if err != nil {
		log.Printf("HandleSearch: GetGitHubToken failed (non-fatal): %v", err)
		githubToken = ""
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page := 1
	limit := 30
	
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	repos, err := h.service.SearchRepositories(ctx, query, githubToken, page, limit)
	if err != nil {
		log.Printf("HandleSearch: SearchRepositories failed: %v", err)
		http.Error(w, "failed to search repos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

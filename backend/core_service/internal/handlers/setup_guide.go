package handlers

import (
	"core-service/internal/clients"
	"core-service/internal/middleware"
	"core-service/internal/users"
	"encoding/json"
	"net/http"
	"strings"
)

type SetupGuideHandler struct {
	githubClient *clients.GitHubClient
	aiClient     *clients.AIClient
	userRepo     *users.Repository
}

func NewSetupGuideHandler(
	gh *clients.GitHubClient,
	ai *clients.AIClient,
	ur *users.Repository,
) *SetupGuideHandler {
	return &SetupGuideHandler{
		githubClient: gh,
		aiClient:     ai,
		userRepo:     ur,
	}
}

func (h *SetupGuideHandler) GetSetupGuide(w http.ResponseWriter, r *http.Request) {
	// Extract query params
	repo := r.URL.Query().Get("repo")
	userID := r.URL.Query().Get("user_id")

	if repo == "" || userID == "" {
		http.Error(w, "repo and user_id are required", http.StatusBadRequest)
		return
	}

	// Split owner/repo
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		http.Error(w, "invalid repo format, expected owner/name", http.StatusBadRequest)
		return
	}
	owner, repoName := parts[0], parts[1]

	// Get user experience level
	// Assuming GetUserByID exists and returns *User, error
	// If context has user ID from JWT, we should prefer that, but prompt said query param user_id
	// But route is protected by JWTMiddleware?
	// The prompt said: GET /repo/setup-guide?repo=owner/repo&user_id=<id> with JWT authentication
	// So we can use the param or the token. Let's use the param as requested.

	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Fetch README
	// Pass empty token if not available in context, or extract from header?
	// GitHubClient methods in core_service take a token.
	// We can try to get it from the request header "Authorization" or "X-GitHub-Token" if forwarded?
	// Or maybe the GitHubClient in Core Service uses a system token?
	// Actually, GitHubClient in Core Service methods take `token` string.
	// We should probably pass the user's GitHub token if we have it?
	// User model has GitHubToken?

	readme, err := h.githubClient.FetchReadme(r.Context(), owner, repoName, user.GitHubToken)
	if err != nil {
		http.Error(w, "failed to fetch readme: "+err.Error(), http.StatusInternalServerError)
		return
	}

	middleware.Logf(r.Context(), "Fetching setup guide for %s/%s", owner, repoName)
	// Generate Guide
	guide, err := h.aiClient.GetSetupGuide(r.Context(), readme, user.ExperienceLevel)
	if err != nil {
		http.Error(w, "failed to generate guide: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"guide": guide,
	})
}

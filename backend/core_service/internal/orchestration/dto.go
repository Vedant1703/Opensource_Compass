package orchestration

type RecommendationResponse struct {
	RepoID string   `json:"repo_id"`
	Score  int      `json:"score"`
	Level  string   `json:"level"`
	Reasons []string `json:"reasons"`
	Issues []Issue 	`json:"issues"`
}

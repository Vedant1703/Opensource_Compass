package clients

type GitHubIssue struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

type GitHubRepo struct {
	Name        string `json:"name"`
	Language    string `json:"language"`
	Stars       int    `json:"stargazers_count"`
	UpdatedAt   string `json:"updated_at"`
	HasIssues   bool   `json:"has_issues"`
}
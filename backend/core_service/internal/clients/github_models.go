package clients

type GitHubIssue struct {
	ID      int      `json:"id"`
	Number  int      `json:"number"`
	Title   string   `json:"title"`
	Body    string   `json:"body"`
	HTMLURL string   `json:"html_url"`
	URL     string   `json:"url"`
	Labels  []string `json:"labels"`
}

type GitHubRepo struct {
	Owner           string   `json:"owner"`
	FullName        string   `json:"full_name"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	HTMLURL         string   `json:"html_url"`
	StargazersCount int      `json:"stargazers_count"`
	Language        string   `json:"language"`
	UpdatedAt       string   `json:"updated_at"`
	Topics          []string `json:"topics"`
	// Legacy fields for backward compatibility
	Stars        int    `json:"stars,omitempty"`
	URL          string `json:"url,omitempty"`
	LastPushedAt string `json:"last_pushed_at,omitempty"`
	OpenIssues   int    `json:"open_issues,omitempty"`
}

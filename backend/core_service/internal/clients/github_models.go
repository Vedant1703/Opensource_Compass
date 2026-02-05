package clients

type GitHubIssue struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	URL    string   `json:"url"`
	Labels []string `json:"labels"`
}

type GitHubRepo struct {
	FullName     string `json:"full_name"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Stars        int    `json:"stars"`
	URL          string `json:"url"`
	LastPushedAt string `json:"last_pushed_at"`
	OpenIssues   int    `json:"open_issues"`
}

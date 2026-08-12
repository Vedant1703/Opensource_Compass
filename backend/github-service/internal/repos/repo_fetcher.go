package repos

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type RepoDTO struct {
	FullName     string   `json:"full_name"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Stars        int      `json:"stars"`
	URL          string   `json:"url"`
	LastPushedAt string   `json:"last_pushed_at"`
	OpenIssues   int      `json:"open_issues"`
	Language     string   `json:"language"`
	Topics       []string `json:"topics"`
}

// Map user-friendly domain names to GitHub topics
// NOTE: Using only PRIMARY topic per domain because GitHub ANDs multiple topics together
func mapDomainToTopics(userDomains []string) []string {
	mapping := map[string]string{
		"web development":       "web",
		"mobile development":    "mobile",
		"ai & machine learning": "machine-learning",
		"devops & cloud":        "devops",
		"game development":      "game",
		"security":              "security",
		"data science":          "data-science",
		"tools & libraries":     "cli",
		"backend development":   "backend",
		"backend":               "backend", // robust fallback
	}

	topics := []string{}
	for _, domain := range userDomains {
		// Frontend converts to lowercase, but we lowercase again just to be safe
		lowerDomain := strings.ToLower(domain)
		if mapped, ok := mapping[lowerDomain]; ok {
			topics = append(topics, mapped)
		}
	}
	return topics
}

func buildQuery(prefix string, values []string) string {
	if len(values) == 0 {
		return ""
	}

	parts := make([]string, 0)
	for _, v := range values {
		parts = append(parts, prefix+":"+v)
	}

	return strings.Join(parts, " ")
}

// generateCacheKey creates a unique key from query parameters
func generateCacheKey(languages, frameworks, domains []string, nameQuery string, page, limit int) string {
	// Sort slices to ensure same key for same params regardless of order
	sortedLangs := make([]string, len(languages))
	copy(sortedLangs, languages)
	sort.Strings(sortedLangs)

	sortedFrameworks := make([]string, len(frameworks))
	copy(sortedFrameworks, frameworks)
	sort.Strings(sortedFrameworks)

	sortedDomains := make([]string, len(domains))
	copy(sortedDomains, domains)
	sort.Strings(sortedDomains)

	// Create string representation
	key := fmt.Sprintf("l:%s|f:%s|d:%s|n:%s|p:%d|lim:%d",
		strings.Join(sortedLangs, ","),
		strings.Join(sortedFrameworks, ","),
		strings.Join(sortedDomains, ","),
		nameQuery,
		page,
		limit,
	)

	// Hash the key to keep it short
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func FetchRepos(languages []string, frameworks []string, domains []string, nameQuery string, token string, cache *RepoCache, page int, limit int) ([]RepoDTO, error) {
	// Generate cache key
	cacheKey := generateCacheKey(languages, frameworks, domains, nameQuery, page, limit)

	// Check cache first
	if cache != nil {
		if cached, found := cache.Get(cacheKey); found {
			log.Printf("Cache hit for query, returning %d cached repositories", len(cached))
			return cached, nil
		}
		log.Println("Cache miss, fetching from GitHub API")
	}

	queryParts := []string{}

	if nameQuery != "" {
		queryParts = append(queryParts, nameQuery+" in:name")
	}

	// Add languages to query (GitHub searches with implicit OR)
	if q := buildQuery("language", languages); q != "" {
		queryParts = append(queryParts, q)
	}

	// Map user-friendly domain names to GitHub topics
	topics := mapDomainToTopics(domains)
	if len(topics) > 0 {
		topicQuery := buildQuery("topic", topics)
		queryParts = append(queryParts, topicQuery)
	}

	query := strings.Join(queryParts, " ")

	// Default query if nothing selected
	if query == "" {
		query = "stars:>1000"
	}

	log.Println("Github search query: ", query)

	url := fmt.Sprintf("https://api.github.com/search/repositories?q=%s&page=%d&per_page=%d&sort=stars&order=desc",
		url.QueryEscape(query), page, limit,
	)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("GitHub API request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("GitHub API returned status %d", resp.StatusCode)
		return nil, fmt.Errorf("github api returned status: %d", resp.StatusCode)
	}

	var raw struct {
		Items []struct {
			FullName        string   `json:"full_name"`
			Name            string   `json:"name"`
			Description     string   `json:"description"`
			Stars           int      `json:"stargazers_count"`
			HTMLURL         string   `json:"html_url"`
			PushedAt        string   `json:"pushed_at"`
			OpenIssuesCount int      `json:"open_issues_count"`
			Language        string   `json:"language"`
			Topics          []string `json:"topics"`
		} `json:"items"`
		TotalCount int `json:"total_count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		log.Printf("Failed to decode GitHub API response: %v", err)
		return nil, err
	}

	log.Printf("GitHub API returned %d repositories (total available: %d)", len(raw.Items), raw.TotalCount)

	repos := make([]RepoDTO, 0)
	for _, r := range raw.Items {
		repos = append(repos, RepoDTO{
			FullName:     r.FullName,
			Name:         strings.Split(r.FullName, "/")[1],
			Description:  r.Description,
			Stars:        r.Stars,
			URL:          r.HTMLURL,
			LastPushedAt: r.PushedAt,
			OpenIssues:   r.OpenIssuesCount,
			Language:     r.Language,
			Topics:       r.Topics,
		})
	}

	log.Printf("Returning %d repositories to client", len(repos))

	// Store in cache before returning
	if cache != nil {
		cache.Set(cacheKey, repos)
		log.Println("Stored results in cache")
	}

	return repos, nil
}

func FetchRepo(owner string, repo string, token string) (*RepoDTO, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status: %d", resp.StatusCode)
	}

	var raw struct {
		FullName        string   `json:"full_name"`
		Name            string   `json:"name"`
		Description     string   `json:"description"`
		Stars           int      `json:"stargazers_count"`
		HTMLURL         string   `json:"html_url"`
		PushedAt        string   `json:"pushed_at"`
		OpenIssuesCount int      `json:"open_issues_count"`
		Language        string   `json:"language"`
		Topics          []string `json:"topics"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	repoDto := &RepoDTO{
		FullName:     raw.FullName,
		Name:         raw.Name,
		Description:  raw.Description,
		Stars:        raw.Stars,
		URL:          raw.HTMLURL,
		LastPushedAt: raw.PushedAt,
		OpenIssues:   raw.OpenIssuesCount,
		Language:     raw.Language,
		Topics:       raw.Topics,
	}
	// Fixing Name assignment logic to be consistent with FetchRepos if needed, or just use raw.Name if it's the repo name
	// In FetchRepos: Name: strings.Split(r.FullName, "/")[1]
	// GitHub API /repos/:owner/:repo returns "name" as just the repo name, so raw.Name is correct if we want that.

	return repoDto, nil
}

func GetRepoReadme(owner string, repo string, token string) (string, error) {
	client := &http.Client{}
	// Use the raw.githubusercontent.com domain or the API with media type header
	// Using API with custom header is often more reliable for private repos if token is used
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/readme", owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	// Request raw content
	req.Header.Set("Accept", "application/vnd.github.v3.raw")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api returned status: %d", resp.StatusCode)
	}

	// Read body
	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

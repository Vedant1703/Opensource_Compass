package routes

import (
	"encoding/json"
	"fmt"
	"github-service/internal/issues"
	"github-service/internal/repos"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var repoCache *repos.RepoCache

func init() {
	// Initialize cache with 5-minute TTL
	repoCache = repos.NewRepoCache(5 * time.Minute)

	// Start background cleanup goroutine
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			repoCache.Clean()
		}
	}()
}

func RegisterGithubRoutes(router *gin.Engine) {
	router.GET("/repos/search", func(c *gin.Context) {
		languagesParam := c.Query("languages")
		frameworksParam := c.Query("frameworks")
		domainsParam := c.Query("domains")
		if domainsParam == "" {
			domainsParam = c.Query("domain")
		}
		// Support both q and name for search query
		nameQuery := c.Query("q")
		if nameQuery == "" {
			nameQuery = c.Query("name")
		}

		var languages []string
		if languagesParam != "" {
			languages = strings.Split(languagesParam, ",")
		}

		var frameworks []string
		if frameworksParam != "" {
			frameworks = strings.Split(frameworksParam, ",")
		}

		var domains []string
		if domainsParam != "" {
			domains = strings.Split(domainsParam, ",")
		}

		pageStr := c.Query("page")
		limitStr := c.Query("limit")
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

		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "GitHub token not configured",
			})
			return
		}

		repoList, err := repos.FetchRepos(languages, frameworks, domains, nameQuery, token, repoCache, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Transform to match frontend interface
		response := make([]gin.H, 0)
		for _, repo := range repoList {
			ownerAndName := strings.Split(repo.FullName, "/")
			owner := ownerAndName[0]

			response = append(response, gin.H{
				"owner":            owner,
				"name":             repo.Name,
				"full_name":        repo.FullName,
				"description":      repo.Description,
				"html_url":         repo.URL,
				"stargazers_count": repo.Stars,
				"language":         repo.Language,
				"updated_at":       repo.LastPushedAt,
				"topics":           repo.Topics,
			})
		}

		c.JSON(http.StatusOK, response)
	})

	router.GET("/repos/:owner/:repo", func(c *gin.Context) {
		owner := c.Param("owner")
		repo := c.Param("repo")
		token := os.Getenv("GITHUB_TOKEN")

		repository, err := repos.FetchRepo(owner, repo, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, repository)
	})

	router.GET("/repos/:owner/:repo/issues/latest", func(c *gin.Context) {
		owner := c.Param("owner")
		repo := c.Param("repo")
		token := os.Getenv("GITHUB_TOKEN")

		if token == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "GitHub token not configured",
			})
			return
		}

		// Call GitHub API to get latest issue
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?state=all&sort=created&direction=desc&per_page=1", owner, repo)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))
		req.Header.Set("Accept", "application/vnd.github+json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("github api returned status: %d", resp.StatusCode),
			})
			return
		}

		var issues []struct {
			Number  int    `json:"number"`
			Title   string `json:"title"`
			HTMLURL string `json:"html_url"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(issues) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"number": 0,
				"title":  "",
				"url":    "",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"number": issues[0].Number,
			"title":  issues[0].Title,
			"url":    issues[0].HTMLURL,
		})
	})

	router.GET("/issues/good-first", func(c *gin.Context) {
		owner := c.Query("owner")
		repo := c.Query("repo")

		if owner == "" || repo == "" {
			c.JSON(400, gin.H{
				"error": "owner and repo are required",
			})
			return
		}
		token := os.Getenv("GITHUB_TOKEN")

		issues, err := issues.FetchGoodFirstIssues(owner, repo, token)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, issues)
	})

	router.GET("/repos/:owner/:repo/readme", func(c *gin.Context) {
		owner := c.Param("owner")
		repo := c.Param("repo")
		token := os.Getenv("GITHUB_TOKEN")

		readme, err := repos.GetRepoReadme(owner, repo, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"content": readme})
	})
}

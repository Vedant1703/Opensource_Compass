package routes

import (
	"github-service/internal/issues"
	"github-service/internal/repos"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterGithubRoutes(router *gin.Engine){
	router.GET("/repos/search", func(c *gin.Context){
		languagesParam := c.Query("languages")
		frameworksParam := c.Query("frameworks")
		domainsParam := c.Query("domains")

		var languages []string
		if languagesParam != ""{
			languages = strings.Split(languagesParam, ",")
		}

		var frameworks []string
		if frameworksParam != ""{
			frameworks = strings.Split(frameworksParam, ",")
		}

		var domains []string
		if domainsParam != ""{
			domains = strings.Split(domainsParam, ",")
		}

		token := os.Getenv("GITHUB_TOKEN")
		if token == ""{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "GitHub token not configured",
			})
			return
		}

		repoList, err := repos.FetchRepos(languages, frameworks, domains, token)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, repoList)
	})

	router.GET("/issues/good-first", func(c *gin.Context){
		owner := c.Query("owner")
		repo := c.Query("repo")

		if owner == "" || repo == ""{
			c.JSON(400, gin.H{
				"error": "owner and repo are required",
			})
			return
		}
		token := os.Getenv("GITHUB_TOKEN")

		issues, err := issues.FetchGoodFirstIssues(owner, repo, token)
		if err != nil{
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, issues)
	})
}
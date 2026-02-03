package main

import (
	"github-service/internal/repos"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	_ = godotenv.Load()
	token := os.Getenv("GITHUB_TOKEN")

	if token == ""{
		log.Fatal("GITHUB_TOKEN not set in environment")
	}

	log.Println("Github Token Loaded")

	// Temporary tester
	repos, err := repos.FetchRepos(
		[]string{"python", "go"},
		nil,
		[]string{"backend"},
		token,
	)
	if err != nil{
		log.Fatalf("Error fetching repos: %v", err)
	}

	log.Printf("Fetched %d repositories\n", len(repos))

	if len(repos) > 0 {
		log.Printf(
			"🔎 Sample Repo → Name: %s | Stars: %d | URL: %s",
			repos[0].Name,
			repos[0].Stars,
			repos[0].URL,
		)
	}
	
	r := gin.Default()


	r.GET("/health", func(c *gin.Context){
		c.JSON(200, gin.H{
			"status": "github-service alive",
		})
	})

	log.Println("Github service running on :8081")
	r.Run(":8081")
}
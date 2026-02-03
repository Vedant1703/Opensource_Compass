package main

import (
	"github-service/routes"
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
	
	r := gin.Default()

	r.GET("/health", func(c *gin.Context){
		c.JSON(200, gin.H{
			"status": "github-service alive",
		})
	})

	routes.RegisterGithubRoutes(r)

	log.Println("Github service running on :8081")
	r.Run(":8081")
}
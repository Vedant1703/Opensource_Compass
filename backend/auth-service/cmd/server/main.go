package main

import (
	"auth-service/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	if err := godotenv.Load(); err != nil{
		fmt.Println("Failed to load .env file")
	}

	required := []string{
		"GITHUB_CLIENT_ID",
		"GITHUB_CLIENT_SECRET",
		"JWT_SECRET",
		"FRONTEND_URL",
	}

	for _, key := range required{
		if os.Getenv(key) == ""{
			log.Fatalf("Missing required env variable: %s", key)
		}
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context){
		c.JSON(200, gin.H{
			"status": "auth-service alive",
		})
	})

	routes.RegisterAuthRoutes(router)

	log.Println("Auth service running on :8080")
	router.Run(":8080")
}
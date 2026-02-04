package routes

import (
	"auth-service/internal/github"
	jwtutil "auth-service/internal/jwt"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine){
	r.GET("/auth/github", func(c *gin.Context){
		url := "https://github.com/login/oauth/authorize" +
		"?client_id=" + os.Getenv("GITHUB_CLIENT_ID") + 
		"&scope=read:user user:email"
		
		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	r.GET("/auth/github/callback", func(c *gin.Context){
		code := c.Query("code")

		accessToken, _ := github.ExchangeCode(code)
		user, _ := github.FetchUser(accessToken)

		jwtToken, _ := jwtutil.GenerateJWT(
			fmt.Sprintf("%v", user["id"]),
		)

		redirect := os.Getenv("FRONTEND_URL") + "/auth/success?token=" + jwtToken

		c.Redirect(http.StatusTemporaryRedirect, redirect)
	})
}
package github

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func ExchangeCode(code string) (string, error){
	data := url.Values{}
	data.Set("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GITHUB_CLIENT_SECRET"))
	data.Set("code", code)

	req, _ := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		strings.NewReader(data.Encode()),
	)

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		return "", err
	}
	defer resp.Body.Close()

	var respData struct{
		AccessToken string `json:"access_token"`
	}

	json.NewDecoder(resp.Body).Decode(&respData)
	return respData.AccessToken, nil
}

func FetchUser(token string) (map[string]interface{}, error){
	req, _ := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)

	req.Header.Set("Authorization", "token "+ token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()

	var user map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&user)

	return user, nil
}
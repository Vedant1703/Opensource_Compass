package repos

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type RepoDTO struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Stars int `json:"stars"`
	URL string `json:"url"`
}

func buildORQuery(prefix string, values []string) string{
	if len(values) == 0{
		return ""
	}

	parts := make([]string, 0)
	for _, v := range values{
		if prefix != ""{
			parts = append(parts, prefix+ ": "+v)
		}else{
			parts = append(parts, v)
		}
	}

	return "(" + strings.Join(parts, " OR ") + ")"
}

func FetchRepos(languages []string, frameworks []string, domains []string, token string) ([]RepoDTO, error){
	queryParts := []string{}

	if q:= buildORQuery("language", languages); q != ""{
		queryParts = append(queryParts, q)
	}

	if q:= buildORQuery("topic", domains); q != ""{
		queryParts = append(queryParts, q)
	}
	query := strings.Join(queryParts, " ")

	log.Println("Github search query: ", query)

	url := fmt.Sprintf("https://api.github.com/search/repositories?q=%s&sort=stars&order=desc", 
		url.QueryEscape(query),
	)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil{
		return nil, err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()

	var raw struct{
		Items []struct{
			FullName string `json:"full_name"`
			Description string `json:"description"`
			Stars int `json:"stargazers_count"`
			HTMLURL string `json:"html_url"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	repos := make([]RepoDTO, 0)
	for _, r := range raw.Items{
		repos = append(repos, RepoDTO{
			Name: r.FullName,
			Description: r.Description,
			Stars: r.Stars,
			URL: r.HTMLURL,
		})
	}
	return repos, nil
}

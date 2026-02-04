package issues

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func FetchGoodFirstIssues(owner, repo, token string) ([]IssueDTO, error){
	rawQuery := fmt.Sprintf(`repo:%s/%s label:"good first issue"`,
		owner, repo,
	)

	encodedQuery := url.QueryEscape(rawQuery)

	apiUrl := "https://api.github.com/search/issues?q=" + encodedQuery + "&per_page=50"


	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil{
		return nil, err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()

	var raw struct{
		Items []struct{
			ID int `json:"id"`
			Title string `json:"title"`
			Body string `json:"body"`
			HTMLURL string `json:"html_url"`
			Labels []struct{
				Name string `json:"name"`
			} `json:"labels"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil{
		return nil, err
	}

	issues := make([]IssueDTO, 0)

	for _, i := range raw.Items{
		labels := []string{}
		for _, l := range i.Labels{
			labels = append(labels, l.Name)
		}

		issues = append(issues, IssueDTO{
			ID: i.ID,
			Title: i.Title,
			Body: i.Body,
			URL: i.HTMLURL,
			Labels: labels,
		})
	}

	return issues, nil
}
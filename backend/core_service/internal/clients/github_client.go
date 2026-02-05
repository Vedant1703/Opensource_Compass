package clients

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

type GitHubClient struct {
	baseURL    string
	httpClient *http.Client
}

func addAuthHeader(req *http.Request, token string) {
	if token == "" {
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
}

func NewGitHubClient(baseURL string) *GitHubClient {
	return &GitHubClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *GitHubClient) FetchIssues(
	ctx context.Context,
	repo string,
) ([]GitHubIssue, error) {

	if c.baseURL == "" {
		return nil, errors.New("github service url not configured")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/repos/"+repo+"/issues",
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("github service returned non-200 response")
	}

	var issues []GitHubIssue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, err
	}

	return issues, nil
}

func (c *GitHubClient) FetchRepo(
	ctx context.Context,
	repo string,
) (*GitHubRepo, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/repos/"+repo,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("github service returned non-200 response")
	}

	var repoData GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repoData); err != nil {
		return nil, err
	}

	return &repoData, nil
}

func (c *GitHubClient) SearchRepos(
	ctx context.Context,
	languages string,
	domains string,
	token string,
) ([]GitHubRepo, error) {

	if c.baseURL == "" {
		return nil, errors.New("github service url not configured")
	}

	endpoint := c.baseURL + "/repos/search"

	params := url.Values{}
	params.Set("languages", languages)
	params.Set("domains", domains)


	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint+"?"+params.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	addAuthHeader(req, token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("github service returned non-200 response")
	}

	var repos []GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func (c *GitHubClient) FetchGoodFirstIssues(
	ctx context.Context,
	owner string,
	repo string,
	token string,
) ([]GitHubIssue, error) {

	if c.baseURL == "" {
		return nil, errors.New("github service url not configured")
	}

	endpoint := c.baseURL + "/issues/good-first"

	params := url.Values{}
	params.Set("owner", owner)
	params.Set("repo", repo)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint+"?"+params.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	addAuthHeader(req, token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("github service returned non-200 response")
	}

	var issues []GitHubIssue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, err
	}

	return issues, nil
}

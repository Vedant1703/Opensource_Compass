package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type AIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAIClient(baseURL string) *AIClient {
	return &AIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *AIClient) AnalyzeIssue(
	ctx context.Context,
	title string,
	body string,
	labels []string,
) (*AIAnalysis, error) {

	if c.baseURL == "" {
		return nil, errors.New("ai service url not configured")
	}

	reqBody := aiRequest{
		Title:  title,
		Body:   body,
		Labels: labels,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/analyze-issue",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("ai service returned non-200 response")
	}

	var analysis AIAnalysis
	if err := json.NewDecoder(resp.Body).Decode(&analysis); err != nil {
		return nil, err
	}

	return &analysis, nil
}

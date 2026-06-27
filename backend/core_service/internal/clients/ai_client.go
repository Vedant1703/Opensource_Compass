package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"core-service/internal/middleware"
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

	if reqID := middleware.GetRequestID(ctx); reqID != "" {
		req.Header.Set("X-Request-ID", reqID)
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

func (c *AIClient) GetSetupGuide(
	ctx context.Context,
	readme string,
	userLevel string,
) (string, error) {
	if c.baseURL == "" {
		return "", errors.New("ai service url not configured")
	}

	reqBody := struct {
		Readme    string `json:"readme"`
		UserLevel string `json:"user_level"`
	}{
		Readme:    readme,
		UserLevel: userLevel,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/generate-setup-guide",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return "", err
	}

	if reqID := middleware.GetRequestID(ctx); reqID != "" {
		req.Header.Set("X-Request-ID", reqID)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("ai service returned non-200 response")
	}

	var result struct {
		Guide string `json:"guide"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Guide, nil
}

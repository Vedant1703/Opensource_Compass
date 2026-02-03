package analysis

import (
	"time"

	"github.com/your-org/opensource-compass/ai-service/internal/formatter"
	"github.com/your-org/opensource-compass/ai-service/internal/llm"
	"github.com/your-org/opensource-compass/ai-service/internal/prompts"
)


type IssueRequest struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

type IssueAnalysis struct {
	Summary    string   `json:"summary"`
	Difficulty string   `json:"difficulty"`
	Skills     []string `json:"skills"`
	FirstSteps []string `json:"first_steps"`
}

type Analyzer struct {
	llm llm.Client
}

func NewAnalyzer(client llm.Client) *Analyzer {
	return &Analyzer{llm: client}
}

func (a *Analyzer) Analyze(req IssueRequest) (*IssueAnalysis, error) {
	var result *IssueAnalysis

	err := Retry(3, 500*time.Millisecond, func() error {
		var e error
		result, e = a.analyzeOnce(req)
		return e
	})

	return result, err
}

func (a *Analyzer) analyzeOnce(req IssueRequest) (*IssueAnalysis, error) {
	prompt := prompts.IssueAnalysisPrompt(
		req.Title,
		req.Body,
		req.Labels,
	)

	resp, err := a.llm.Chat([]llm.Message{
		{Role: "user", Content: prompt},
	})
	if err != nil {
		return nil, err
	}

	jsonBytes, err := formatter.ExtractJSON(resp)
	if err != nil {
		return nil, err
	}

	var result IssueAnalysis
	if err := formatter.ValidateJSON(jsonBytes, &result); err != nil {
		return nil, err
	}

	return &result, nil
}



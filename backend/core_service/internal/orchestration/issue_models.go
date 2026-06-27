package orchestration

import "core-service/internal/clients"

type Issue struct {
	Title       string   			`json:"title"`
	Body        string   			`json:"body"`
	Labels      []string            `json:"labels"`
	IssueNumber int          		`json:"issue_number,omitempty"`
	IssueURL    string       		`json:"issue_url,omitempty"`
	AI          *clients.AIAnalysis `json:"ai,omitempty"`
}
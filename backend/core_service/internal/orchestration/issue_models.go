package orchestration

import "core-service/internal/clients"

type Issue struct {
	Title   string   `json:"title"`
	Body    string   `json:"body"`
	Labels  []string `json:"labels"`

	AI *clients.AIAnalysis `json:"ai,omitempty"`
}

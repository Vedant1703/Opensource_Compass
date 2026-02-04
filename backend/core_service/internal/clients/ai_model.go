package clients

type AIAnalysis struct {
	Summary    string   `json:"summary"`
	Difficulty string   `json:"difficulty"`
	Skills     []string `json:"skills"`
	FirstSteps []string `json:"first_steps"`
}

type aiRequest struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

package routes

import (
	"encoding/json"
	"net/http"

	"github.com/your-org/opensource-compass/ai-service/internal/analysis"
)

func AnalyzeIssue(analyzer *analysis.Analyzer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req analysis.IssueRequest
		json.NewDecoder(r.Body).Decode(&req)

		result, err := analyzer.Analyze(req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}

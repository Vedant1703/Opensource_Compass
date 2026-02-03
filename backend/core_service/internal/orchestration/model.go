package orchestration

import "core-service/internal/scoring"

type RepoRecommendation struct {
	RepoID   string
	Score    scoring.ScoreResult
	Signals  scoring.RepoSignals
}

package orchestration

import (
	"context"
	"testing"

	"core-service/internal/preferences"
)

func TestScoreRepositoryForUser(t *testing.T) {
	ctx := context.Background()

	// fake user context
	userCtx := &UserContext{
		UserID: "user-123",
		Preferences: []preferences.Preference{
			{PreferenceType: "language", Value: "Go"},
			{PreferenceType: "domain", Value: "Backend"},
		},
	}

	signals := BuildMockSignals(80)

	svc := &Service{}

	rec, err := svc.ScoreRepositoryForUser(ctx, userCtx, "repo-xyz", signals)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Score.TotalScore <= 0 {
		t.Fatalf("expected positive score, got %d", rec.Score.TotalScore)
	}

	if rec.Score.Level == "" {
		t.Fatalf("expected difficulty level")
	}
}

package orchestration

import (
	"context"

	"core-service/internal/clients"
	"core-service/internal/preferences"
	"core-service/internal/scoring"
)

type Service struct {
	aiClient *clients.AIClient
	prefRepo *preferences.Repository
	githubClient  *clients.GitHubClient
}

func NewService(prefRepo *preferences.Repository, aiClient *clients.AIClient,githubClient  *clients.GitHubClient) *Service {
	return &Service{
		prefRepo: prefRepo,
		aiClient: aiClient,
		githubClient: githubClient,
	}
}

type UserContext struct {
	UserID      string
	Preferences []preferences.Preference
}

func (s *Service) BuildUserContext(
	ctx context.Context,
	userID string,
) (*UserContext, error) {

	prefs, err := s.prefRepo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserContext{
		UserID:      userID,
		Preferences: prefs,
	}, nil
}

func (s *Service) ScoreRepositoryForUser(
	ctx context.Context,
	userCtx *UserContext,
	repoID string,
	signals scoring.RepoSignals,
) (*RepoRecommendation, error) {

	scoreResult := scoring.ComputeScore(signals)

	return &RepoRecommendation{
		RepoID:  repoID,
		Score:   scoreResult,
		Signals: signals,
	}, nil
}

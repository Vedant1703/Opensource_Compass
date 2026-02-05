package orchestration

import (
	"context"
	"time"

	"core-service/internal/scoring"
)

func (s *Service) BuildRepoSignals(
	ctx context.Context,
	repo string,
	userCtx *UserContext,
) scoring.RepoSignals {

	signals := scoring.RepoSignals{}

	if s.githubClient == nil {
		return signals
	}

	repoData, err := s.githubClient.FetchRepo(ctx, repo)
	if err != nil {
		return signals
	}
	signals.TechStackMatchPct = 70

	updatedTime, _ := time.Parse(time.RFC3339, repoData.LastPushedAt)
	if time.Since(updatedTime).Hours() < 24*30 {
		signals.RecentActivity = true
	}

	if repoData.OpenIssues>0 {
		signals.HasGoodFirstIssues = true
	}

	signals.HasClearReadme = true

	return signals
}

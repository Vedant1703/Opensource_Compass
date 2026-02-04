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

	langs := ExtractLanguages(userCtx)
techMatch := ComputeTechMatch(langs, repoData.Language)

signals.TechStackMatchPct = techMatch

	updatedTime, _ := time.Parse(time.RFC3339, repoData.UpdatedAt)
	if time.Since(updatedTime).Hours() < 24*30 {
		signals.RecentActivity = true
	}

	if repoData.HasIssues {
		signals.HasGoodFirstIssues = true
	}

	signals.HasClearReadme = true

	return signals
}

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
	
	userLangs := ExtractLanguages(userCtx)
	userDomains := ExtractDomains(userCtx)
	signals.TechStackMatchPct = ComputeTechMatch(userLangs, userDomains, repoData.Language, repoData.Topics)

	updatedTime, _ := time.Parse(time.RFC3339, repoData.LastPushedAt)
	if time.Since(updatedTime).Hours() < 24*30 {
		signals.RecentActivity = true
	}

	// Mark MaintainerActive = true if the repo was updated in the last 7 days. Frequent recent pushes suggest active maintenance, which is rewarded in scoring (+20 points).
	pushedTime, err := time.Parse(time.RFC3339, repoData.LastPushedAt)
	if err == nil && time.Since(pushedTime).Hours() < 24*7 {
		signals.MaintainerActive = true
	}

	if repoData.OpenIssues>0 {
		signals.HasGoodFirstIssues = true
	}

	signals.HasClearReadme = true

	return signals
}

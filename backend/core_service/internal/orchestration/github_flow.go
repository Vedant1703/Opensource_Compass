package orchestration

import (
	"context"
	"strings"
)

func (s *Service) FetchAndEnrichIssues(
	ctx context.Context,
	fullRepo string,
	maxAI int,
	token string,
) ([]Issue, error) {

	parts := strings.Split(fullRepo, "/")
	if len(parts) != 2 {
		return BuildMockIssues(), nil
	}

	owner := parts[0]
	repo := parts[1]

	if s.githubClient == nil {
		return BuildMockIssues(), nil
	}

	rawIssues, err := s.githubClient.FetchGoodFirstIssues(ctx, owner, repo, token)
	if err != nil {
		// graceful fallback
		return BuildMockIssues(), nil
	}

	var issues []Issue
	for _, gi := range rawIssues {
		issues = append(issues, Issue{
			Title:  gi.Title,
			Body:   gi.Body,
			Labels: gi.Labels,
		})
	}

	issues = s.EnrichIssues(ctx, issues, maxAI)

	return issues, nil
}

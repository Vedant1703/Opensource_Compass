package orchestration

import (
	"context"
)

func (s *Service) FetchAndEnrichIssues(
	ctx context.Context,
	repo string,
	maxAI int,
) ([]Issue, error) {

	if s.githubClient == nil {
		return BuildMockIssues(), nil
	}

	rawIssues, err := s.githubClient.FetchIssues(ctx, repo)
	if err != nil {
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

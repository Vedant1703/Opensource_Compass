package orchestration

import "context"

func (s *Service) EnrichIssues(
	ctx context.Context,
	issues []Issue,
	maxAI int,
) []Issue {

	count := 0

	for i := range issues {
		if count >= maxAI {
			break
		}

		analysis := s.EnrichIssueWithAI(
			ctx,
			issues[i].Title,
			issues[i].Body,
			issues[i].Labels,
		)

		if analysis != nil {
			issues[i].AI = analysis
			count++
		}
	}

	return issues
}

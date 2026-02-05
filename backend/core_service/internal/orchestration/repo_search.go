package orchestration

import (
	"context"
	"errors"
)

func (s *Service) SearchReposForUser(
	ctx context.Context,
	userCtx *UserContext,
	token string,
) ([]string, error) {

	if s.githubClient == nil {
		return nil, errors.New("github client not configured")
	}

	langs := ExtractLanguages(userCtx)

	languageQuery := ""
	for i, l := range langs {
		if i > 0 {
			languageQuery += ","
		}
		languageQuery += l
	}

	domainQuery := "backend"

	repos, err := s.githubClient.SearchRepos(ctx, languageQuery, domainQuery, token)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, r := range repos {
		repoID := r.FullName
		if repoID == "" {
			repoID = r.Name
		}
		results = append(results, repoID)
	}

	return results, nil
}

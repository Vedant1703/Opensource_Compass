package orchestration

import (
	"context"
	"core-service/internal/clients"
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

	domains := ExtractDomains(userCtx)
	domainQuery := ""
	for i, d := range domains {
		if i > 0 {
			domainQuery += ","
		}
		domainQuery += d
	}

	if domainQuery == "" {
		domainQuery = "backend"
	}

	repos, err := s.githubClient.SearchRepos(ctx, languageQuery, domainQuery, "", token, 1, 30) // Hardcoded 1, 30 for SearchReposForUser recommendation internal tasks
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

func (s *Service) SearchRepositories(
	ctx context.Context,
	query string,
	token string,
	page int,
	limit int,
) ([]clients.GitHubRepo, error) {
	if s.githubClient == nil {
		return nil, errors.New("github client not configured")
	}

	return s.githubClient.SearchRepos(ctx, "", "", query, token, page, limit)
}

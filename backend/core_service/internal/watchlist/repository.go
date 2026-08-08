package watchlist

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

type UniqueRepo struct {
	RepoOwner string
	RepoName  string
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(ctx context.Context, userID, owner, name string) error {
	query := `
		INSERT INTO watched_repos (user_id, repo_owner, repo_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, repo_owner, repo_name) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, userID, owner, name)
	if err != nil {
		return fmt.Errorf("failed to add to watchlist: %w", err)
	}
	return nil
}

func (r *Repository) Remove(ctx context.Context, userID, owner, name string) error {
	query := `
		DELETE FROM watched_repos
		WHERE user_id = $1 AND repo_owner = $2 AND repo_name = $3
	`
	_, err := r.db.Exec(ctx, query, userID, owner, name)
	if err != nil {
		return fmt.Errorf("failed to remove from watchlist: %w", err)
	}
	return nil
}

func (r *Repository) List(ctx context.Context, userID string) ([]WatchedRepo, error) {
	query := `
		SELECT id, user_id, repo_owner, repo_name, created_at
		FROM watched_repos
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query watchlist: %w", err)
	}
	defer rows.Close()

	repos := []WatchedRepo{}
	for rows.Next() {
		var repo WatchedRepo
		if err := rows.Scan(&repo.ID, &repo.UserID, &repo.RepoOwner, &repo.RepoName, &repo.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

func (r *Repository) IsWatched(ctx context.Context, userID, owner, name string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM watched_repos 
			WHERE user_id = $1 AND repo_owner = $2 AND repo_name = $3
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, owner, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check watchlist status: %w", err)
	}
	return exists, nil
}

func (r *Repository) ListAllChunked(ctx context.Context, limit, offset int) ([]WatchedRepo, error) {
	query := `
		SELECT id, user_id, repo_owner, repo_name, last_checked_at, latest_issue_number, created_at
		FROM watched_repos
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list all watchlist items chunked: %w", err)
	}
	defer rows.Close()

	repos := []WatchedRepo{}
	for rows.Next() {
		var repo WatchedRepo
		var lastCheckedAt *time.Time // Handle NULL
		if err := rows.Scan(&repo.ID, &repo.UserID, &repo.RepoOwner, &repo.RepoName, &lastCheckedAt, &repo.LatestIssueNumber, &repo.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if lastCheckedAt != nil {
			repo.LastCheckedAt = *lastCheckedAt
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

func (r *Repository) UpdateLastChecked(ctx context.Context, id int, latestIssueNum int) error {
	query := `
		UPDATE watched_repos
		SET last_checked_at = NOW(), latest_issue_number = $2
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id, latestIssueNum)
	if err != nil {
		return fmt.Errorf("failed to update watchlist item: %w", err)
	}
	return nil
}

func (r *Repository) ListAllUniqueReposChunked(ctx context.Context, limit, offset int) ([]UniqueRepo, error) {
	query := `
		SELECT DISTINCT repo_owner, repo_name
		FROM watched_repos
		ORDER BY repo_owner, repo_name
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list unique repos: %w", err)
	}
	defer rows.Close()

	var repos []UniqueRepo
	for rows.Next() {
		var repo UniqueRepo
		if err := rows.Scan(&repo.RepoOwner, &repo.RepoName); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

func (r *Repository) GetOutdatedWatches(ctx context.Context, owner, name string, latestNum int) ([]WatchedRepo, error) {
	query := `
		SELECT id, user_id, repo_owner, repo_name, last_checked_at, latest_issue_number, created_at
		FROM watched_repos
		WHERE repo_owner = $1 AND repo_name = $2 AND latest_issue_number < $3
	`
	rows, err := r.db.Query(ctx, query, owner, name, latestNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get outdated watches: %w", err)
	}
	defer rows.Close()

	var repos []WatchedRepo
	for rows.Next() {
		var repo WatchedRepo
		var lastCheckedAt *time.Time // Handle NULL
		if err := rows.Scan(&repo.ID, &repo.UserID, &repo.RepoOwner, &repo.RepoName, &lastCheckedAt, &repo.LatestIssueNumber, &repo.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if lastCheckedAt != nil {
			repo.LastCheckedAt = *lastCheckedAt
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

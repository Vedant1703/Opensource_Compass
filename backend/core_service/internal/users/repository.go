package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, github_username, experience_level)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query, user.ID, user.GitHubUsername, user.ExperienceLevel)
	return err
}

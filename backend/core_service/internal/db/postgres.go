package db

import (
	"context"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	// Supabase recently deprecated the postgres.[project-ref] username format
	// for connection pooling (Supavisor). We sanitize the connection string
	// here to ensure seamless deployments with older environment variables.
	parsedURL, err := url.Parse(databaseURL)
	if err == nil && parsedURL.User != nil {
		username := parsedURL.User.Username()
		if strings.HasPrefix(username, "postgres.") {
			password, hasPassword := parsedURL.User.Password()
			if hasPassword {
				parsedURL.User = url.UserPassword("postgres", password)
			} else {
				parsedURL.User = url.User("postgres")
			}
			databaseURL = parsedURL.String()
			log.Println("Sanitized Supabase database URL: removed project-ref from username")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL successfully")
	return pool, nil
}

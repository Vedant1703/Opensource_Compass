package watchlist

import (
	"context"
	"log"
	"sync"
	"time"
)

type Poller struct {
	repo         *Repository
	githubClient GitHubClient
	notifiers    []Notifier // Support multiple notifiers
	interval     time.Duration
}

type GitHubClient interface {
	GetLatestIssueNumber(owner, name string) (int, error)
	GetLatestIssue(owner, name string) (number int, title string, url string, err error)
}

type Notifier interface {
	NotifyUser(userID string, payload interface{}) error
}

func NewPoller(repo *Repository, gh GitHubClient, notifiers []Notifier, interval time.Duration) *Poller {
	return &Poller{
		repo:         repo,
		githubClient: gh,
		notifiers:    notifiers,
		interval:     interval,
	}
}

func (p *Poller) Start(ctx context.Context) {
	log.Printf("Starting Poller with interval: %v", p.interval)
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	// Run immediately on start
	p.poll(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("Poller stopping...")
			return
		case <-ticker.C:
			p.poll(ctx)
		}
	}
}

func (p *Poller) poll(ctx context.Context) {
	log.Println("Poller: Starting poll cycle")
	// TODO: Fetch distinct repos to avoid duplicate checks if multiple users watch the same repo
	// For MVP, we iterate all entries. Optimization: implement "ListAllUniqueRepos" in repository.

	limit := 100
	offset := 0
	
	// Create a worker pool semaphore
	sem := make(chan struct{}, 10) // Concurrency bounded to 10
	var wg sync.WaitGroup

	for {
		entries, err := p.repo.ListAllChunked(ctx, limit, offset)
		if err != nil {
			log.Printf("Poller: Error listing watched repos chunked at offset %d: %v", offset, err)
			break
		}
		
		if len(entries) == 0 {
			break
		}
		
		log.Printf("Poller: Processing chunk of %d entries (offset %d)", len(entries), offset)

		for _, entry := range entries {
			wg.Add(1)
			sem <- struct{}{} // Acquire token
			
			go func(e WatchedRepo) {
				defer wg.Done()
				defer func() { <-sem }() // Release token
				
				log.Printf("Poller: Checking repo %s/%s (Last issue: %d)", e.RepoOwner, e.RepoName, e.LatestIssueNumber)
				latestNum, title, issueURL, err := p.githubClient.GetLatestIssue(e.RepoOwner, e.RepoName)
				if err != nil {
					log.Printf("Poller: Error fetching latest issue for %s/%s: %v", e.RepoOwner, e.RepoName, err)
					return
				}
				
				if latestNum > e.LatestIssueNumber {
					log.Printf("Poller: New issue detected for %s/%s! Updating DB and notifying user %s", e.RepoOwner, e.RepoName, e.UserID)

					// Updates DB
					if err := p.repo.UpdateLastChecked(ctx, e.ID, latestNum); err != nil {
						log.Printf("Poller: Error updating last checked for %s/%s: %v", e.RepoOwner, e.RepoName, err)
					}

					// Send notification with issue details via all notifiers
					payload := map[string]interface{}{
						"type":         "new_issue",
						"repo":         e.RepoOwner + "/" + e.RepoName,
						"issue_number": latestNum,
						"issue_title":  title,
						"issue_url":    issueURL,
						"message":      "New issue detected!",
					}

					for _, notifier := range p.notifiers {
						if err := notifier.NotifyUser(e.UserID, payload); err != nil {
							log.Printf("Poller: Error notifying user %s via notifier: %v", e.UserID, err)
						} else {
							log.Printf("Poller: Successfully notified user %s", e.UserID)
						}
					}
				} else {
					log.Printf("Poller: No new issues for %s/%s", e.RepoOwner, e.RepoName)
				}
			}(entry)
		}
		
		offset += limit
	}
	
	wg.Wait()
	log.Println("Poller: cycle finished")
}

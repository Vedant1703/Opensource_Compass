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

	limit := 100
	offset := 0
	
	// Create a worker pool semaphore
	sem := make(chan struct{}, 10) // Concurrency bounded to 10
	var wg sync.WaitGroup

	for {
		uniqueRepos, err := p.repo.ListAllUniqueReposChunked(ctx, limit, offset)
		if err != nil {
			log.Printf("Poller: Error listing unique watched repos at offset %d: %v", offset, err)
			break
		}
		
		if len(uniqueRepos) == 0 {
			break
		}
		
		log.Printf("Poller: Processing chunk of %d unique repos (offset %d)", len(uniqueRepos), offset)

		for _, repoEntry := range uniqueRepos {
			wg.Add(1)
			sem <- struct{}{} // Acquire token
			
			go func(e UniqueRepo) {
				defer wg.Done()
				defer func() { <-sem }() // Release token
				
				log.Printf("Poller: Checking unique repo %s/%s", e.RepoOwner, e.RepoName)
				latestNum, title, issueURL, err := p.githubClient.GetLatestIssue(e.RepoOwner, e.RepoName)
				if err != nil {
					log.Printf("Poller: Error fetching latest issue for %s/%s: %v", e.RepoOwner, e.RepoName, err)
					return
				}
				
				// Find all subscribed users who are outdated for this repo
				outdatedEntries, err := p.repo.GetOutdatedWatches(ctx, e.RepoOwner, e.RepoName, latestNum)
				if err != nil {
					log.Printf("Poller: Error fetching outdated watches for %s/%s: %v", e.RepoOwner, e.RepoName, err)
					return
				}

				if len(outdatedEntries) > 0 {
					log.Printf("Poller: New issue detected for %s/%s! Updating DB and notifying %d users", e.RepoOwner, e.RepoName, len(outdatedEntries))
					for _, watch := range outdatedEntries {
						// Updates DB
						if err := p.repo.UpdateLastChecked(ctx, watch.ID, latestNum); err != nil {
							log.Printf("Poller: Error updating last checked for watch ID %d: %v", watch.ID, err)
							continue
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
							if err := notifier.NotifyUser(watch.UserID, payload); err != nil {
								log.Printf("Poller: Error notifying user %s via notifier: %v", watch.UserID, err)
							} else {
								log.Printf("Poller: Successfully notified user %s", watch.UserID)
							}
						}
					}
				} else {
					log.Printf("Poller: No new issues for %s/%s", e.RepoOwner, e.RepoName)
				}
			}(repoEntry)
		}
		
		offset += limit
	}
	
	wg.Wait()
	log.Println("Poller: cycle finished")
}

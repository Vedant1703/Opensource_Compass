package watchlist

import (
	"context"
	"log"
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

		for _, e := range uniqueRepos {
			log.Printf("Poller: Checking unique repo %s/%s", e.RepoOwner, e.RepoName)
				latestNum, title, issueURL, err := p.githubClient.GetLatestIssue(e.RepoOwner, e.RepoName)
				if err != nil {
					log.Printf("Poller: Error fetching latest issue for %s/%s: %v", e.RepoOwner, e.RepoName, err)
					continue
				}
				
				// Find all subscribed users who are outdated for this repo
				outdatedEntries, err := p.repo.GetOutdatedWatches(ctx, e.RepoOwner, e.RepoName, latestNum)
				if err != nil {
					log.Printf("Poller: Error fetching outdated watches for %s/%s: %v", e.RepoOwner, e.RepoName, err)
					continue
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

				// Sleep to respect GitHub secondary rate limits (no concurrent requests)
				time.Sleep(1 * time.Second)
		}
		
		offset += limit
	}
	
	log.Println("Poller: cycle finished")
}

package orchestration

import (
	"context"
	"sync"
)

func (s *Service) EnrichIssues(ctx context.Context, issues []Issue, maxAI int) []Issue {

    // Collect the indices of issues eligible for AI enrichment
    var eligibleIdx []int
    for i, issue := range issues {
        if IsGoodFirstIssue(issue.Labels) {
            eligibleIdx = append(eligibleIdx, i)
        }
        if len(eligibleIdx) >= maxAI {
            break
        }
    }

    if len(eligibleIdx) == 0 {
        return issues
    }

    // Use a mutex to protect concurrent writes to the issues slice
    var mu sync.Mutex
    var wg sync.WaitGroup

    for _, idx := range eligibleIdx {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()

            analysis := s.EnrichIssueWithAI(
                ctx,
                issues[i].Title,
                issues[i].Body,
                issues[i].Labels,
            )

            if analysis != nil {
                mu.Lock()
                issues[i].AI = analysis
                mu.Unlock()
            }
        }(idx)
    }

    wg.Wait()  // Block until all goroutines complete
    return issues
}
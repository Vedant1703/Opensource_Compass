package orchestration

import (
	"context"
	"log"

	"core-service/internal/clients"
)

func (s *Service) EnrichIssueWithAI(
	ctx context.Context,
	title string,
	body string,
	labels []string,
) *clients.AIAnalysis {

	if s.aiClient == nil {
		return nil
	}

	analysis, err := s.aiClient.AnalyzeIssue(ctx, title, body, labels)
	if err != nil {
		log.Printf("AI analysis failed: %v", err)
		return nil
	}

	return analysis
}

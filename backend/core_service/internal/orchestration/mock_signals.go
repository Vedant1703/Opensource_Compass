package orchestration

import "core-service/internal/scoring"

func BuildMockSignals(techMatchPct int) scoring.RepoSignals {
	return scoring.RepoSignals{
		HasGoodFirstIssues: true,
		RecentActivity:     true,
		HasClearReadme:     true,
		MaintainerActive:   false,
		TechStackMatchPct:  techMatchPct,
	}
}

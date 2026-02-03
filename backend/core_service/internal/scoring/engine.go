package scoring

func ComputeScore(signals RepoSignals) ScoreResult {
	score := 0
	reasons := []string{}

	if signals.HasGoodFirstIssues {
		score += 25
		reasons = append(reasons, "Beginner-friendly issues available")
	}

	if signals.RecentActivity {
		score += 20
		reasons = append(reasons, "Repository is actively maintained")
	}

	if signals.HasClearReadme {
		score += 15
		reasons = append(reasons, "Clear documentation")
	}

	if signals.MaintainerActive {
		score += 20
		reasons = append(reasons, "Maintainer responds to issues/PRs")
	}

	techScore := (signals.TechStackMatchPct * 20) / 100
	score += techScore

	if signals.TechStackMatchPct >= 70 {
		reasons = append(reasons, "Strong match with your tech stack")
	}

	level := classify(score)

	return ScoreResult{
		TotalScore: score,
		Level:      level,
		Reasons:    reasons,
	}
}

func classify(score int) string {
	switch {
	case score >= 75:
		return "Beginner"
	case score >= 50:
		return "Intermediate"
	default:
		return "Advanced"
	}
}

package orchestration

func BuildMockIssues() []Issue {
	return []Issue{
		{
			Title:  "Fix broken README link",
			Body:   "The documentation contains an outdated link that needs updating.",
			Labels: []string{"good first issue", "documentation"},
		},
		{
			Title:  "Add unit tests for scoring engine",
			Body:   "We need test coverage for ComputeScore in the scoring module.",
			Labels: []string{"help wanted", "testing"},
		},
	}
}

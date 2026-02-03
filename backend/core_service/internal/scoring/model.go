package scoring

type RepoSignals struct {
	HasGoodFirstIssues bool
	RecentActivity     bool
	HasClearReadme     bool
	MaintainerActive   bool
	TechStackMatchPct  int 
}

type ScoreResult struct {
	TotalScore int
	Level      string 
	Reasons    []string
}

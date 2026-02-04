package orchestration

import "strings"

func IsGoodFirstIssue(labels []string) bool {
	for _, label := range labels {
		if strings.ToLower(label) == "good first issue" {
			return true
		}
	}
	return false
}

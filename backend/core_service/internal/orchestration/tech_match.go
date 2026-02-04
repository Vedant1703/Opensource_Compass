package orchestration

import "strings"

func ComputeTechMatch(
	userPrefs []string,
	repoLanguage string,
) int {

	if repoLanguage == "" {
		return 0
	}

	repoLang := strings.ToLower(repoLanguage)

	for _, pref := range userPrefs {
		if strings.ToLower(pref) == repoLang {
			return 100
		}
	}

	return 30
}

package orchestration

import "core-service/internal/preferences"

func GroupPreferencesByType(
	prefs []preferences.Preference,
) map[string][]string {

	result := make(map[string][]string)

	for _, p := range prefs {
		result[p.PreferenceType] = append(
			result[p.PreferenceType],
			p.Value,
		)
	}

	return result
}

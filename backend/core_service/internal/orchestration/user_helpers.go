package orchestration

func ExtractLanguages(userCtx *UserContext) []string {
	var langs []string

	for _, p := range userCtx.Preferences {
		if p.PreferenceType == "language" {
			langs = append(langs, p.Value)
		}
	}

	return langs
}

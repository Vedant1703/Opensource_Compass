package prompts

func IssueAnalysisPrompt(
	title string,
	body string,
	labels []string,
) string {
	return `
You are an experienced open-source mentor.

Analyze the GitHub issue below and respond in JSON.

Issue Title:
` + title + `

Issue Description:
` + body + `

Labels:
` + formatLabels(labels) + `

Return JSON with:
- summary (simple explanation)
- difficulty (Easy | Medium | Hard)
- skills (array)
- first_steps (array)
Return ONLY raw JSON.
No markdown.
No explanations.
`
}

func formatLabels(labels []string) string {
	out := ""
	for _, l := range labels {
		out += "- " + l + "\n"
	}
	return out
}

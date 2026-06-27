package orchestration

import "strings"

// Updated signature:
func ComputeTechMatch(userLanguages []string, userDomains []string, repoLanguage string, repoTopics []string) int {

    score := 0

    // Check primary language match (full weight: 100)
    if repoLanguage != "" {
        repoLang := strings.ToLower(repoLanguage)
        for _, pref := range userLanguages {
            if strings.ToLower(pref) == repoLang {
                score = 100
                break
            }
        }
    }

    // Check topic match (partial weight: 50 if at least one topic matches)
    if score == 0 && len(repoTopics) > 0 && len(userDomains) > 0 {
        domainTopicMap := map[string][]string{
            "Machine Learning":    {"machine-learning", "deep-learning", "ai"},
            "AI":                  {"ai", "machine-learning", "llm"},
            "Web Development":     {"web", "frontend", "react", "vue"},
            "Backend":             {"backend", "api", "rest"},
            "DevOps":              {"devops", "docker", "kubernetes", "ci-cd"},
            "Security":            {"security", "cryptography"},
            "Data Science":        {"data-science", "data-analysis", "pandas"},
            "Mobile":              {"android", "ios", "flutter", "swift"},
            "Game Development":    {"game", "unity", "godot"},
            "Blockchain":          {"blockchain", "web3", "ethereum"},
        }
        for _, domain := range userDomains {
            if topics, ok := domainTopicMap[domain]; ok {
                for _, userTopic := range topics {
                    for _, repoTopic := range repoTopics {
                        if strings.ToLower(repoTopic) == userTopic {
                            score = 50  // topic match → partial score
                            break
                        }
                    }
                    if score > 0 { break }
                }
            }
            if score > 0 { break }
        }
    }

    return score
}

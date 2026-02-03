package llm

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Client interface {
	Chat(messages []Message) (string, error)
}

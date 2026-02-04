package issues

type IssueDTO struct{
	ID int `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
	URL string `json:"url"`
	Labels []string `json:"labels"`
}
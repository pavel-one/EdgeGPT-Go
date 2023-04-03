package responses

type GptResponse interface {
	GetAnswer() string
	GetType() int
	GetMaxUnit() int
	GetUserUnit() int
	GetSuggestions() []*Suggestion
}

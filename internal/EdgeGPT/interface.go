package EdgeGPT

type GptResponse interface {
	GetAnswer() string
	GetType() int
}

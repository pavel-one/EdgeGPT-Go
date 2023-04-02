package EdgeGPT

type GptResponse interface {
	GetAnswer() string
	GetType() int
	GetMaxUnit() int
	GetUserUnit() int
}

type StorageGpt interface {
	GetOrSet(key string) (*GPT, error)
	Add(gpt *GPT, key string)
	Get(key string) (*GPT, error)
	Remove(key string) error
}

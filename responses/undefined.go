package responses

// Undefined response for unused messages
type Undefined struct {
	Type int `json:"type"`
}

// GetAnswer get answer text
func (r *Undefined) GetAnswer() string {
	return ""
}

// GetType get type
func (r *Undefined) GetType() int {
	return r.Type
}

// GetMaxUnit get max user questions for current session
func (r *Undefined) GetMaxUnit() int {
	return 0
}

// GetUserUnit get current question for current session
func (r *Undefined) GetUserUnit() int {
	return 0
}

func (r *Undefined) GetSuggestions() []*Suggestion {
	return nil
}

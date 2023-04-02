package EdgeGPT

// UndefinedResponse response for unused messages
type UndefinedResponse struct {
	Type int `json:"type"`
}

// GetAnswer get answer text
func (r *UndefinedResponse) GetAnswer() string {
	return ""
}

// GetType get type
func (r *UndefinedResponse) GetType() int {
	return r.Type
}

// GetMaxUnit get max user questions for current session
func (r *UndefinedResponse) GetMaxUnit() int {
	return 0
}

// GetUserUnit get current question for current session
func (r *UndefinedResponse) GetUserUnit() int {
	return 0
}

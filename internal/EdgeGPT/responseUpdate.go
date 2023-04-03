package EdgeGPT

// UpdateResponse response for work generate message
type UpdateResponse struct {
	Type      int    `json:"type"`
	Target    string `json:"target"`
	Arguments []struct {
		Cursor struct {
			J string `json:"j"`
			P int    `json:"p"`
		} `json:"cursor"`
		Messages  []*Message `json:"messages"`
		RequestId string     `json:"requestId"`
	} `json:"arguments"`
}

// GetAnswer get answer text
func (u *UpdateResponse) GetAnswer() string {
	arg := u.Arguments[0]
	if len(arg.Messages) == 0 {
		return ""
	}

	message := arg.Messages[len(arg.Messages)-1]

	return message.AdaptiveCards[0].Body[0].Text
}

// GetType get type
func (u *UpdateResponse) GetType() int {
	return u.Type
}

// GetMaxUnit get max user questions for current session
func (u *UpdateResponse) GetMaxUnit() int {
	return 0
}

// GetUserUnit get current question for current session
func (u *UpdateResponse) GetUserUnit() int {
	return 0
}

func (u *UpdateResponse) GetSuggestions() []*Suggestion {
	return nil
}

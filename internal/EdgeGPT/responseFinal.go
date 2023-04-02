package EdgeGPT

import "time"

// FinalResponse response for final generate message
type FinalResponse struct {
	Type         int    `json:"type"`
	InvocationId string `json:"invocationId,omitempty"`
	Item         struct {
		Messages               []*Message `json:"messages"`
		ConversationExpiryTime time.Time  `json:"conversationExpiryTime,omitempty"`
		Throttling             struct {
			MaxNumUserMessagesInConversation int `json:"maxNumUserMessagesInConversation"`
			NumUserMessagesInConversation    int `json:"numUserMessagesInConversation"`
		} `json:"throttling"`
	} `json:"item"`
}

// GetAnswer get answer text
func (r *FinalResponse) GetAnswer() string {
	item := r.Item
	if len(item.Messages) == 0 {
		return ""
	}

	message := item.Messages[len(item.Messages)-1]

	return message.AdaptiveCards[0].Body[0].Text
}

// GetType get type
func (r *FinalResponse) GetType() int {
	return r.Type
}

// GetMaxUnit get max user questions for current session
func (r *FinalResponse) GetMaxUnit() int {
	return r.Item.Throttling.MaxNumUserMessagesInConversation
}

// GetUserUnit get current question for current session
func (r *FinalResponse) GetUserUnit() int {
	return r.Item.Throttling.NumUserMessagesInConversation
}

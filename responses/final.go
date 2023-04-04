package responses

import (
	"time"
)

// Final response for final generate message
type Final struct {
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
func (r *Final) GetAnswer() string {
	item := r.Item
	if len(item.Messages) == 0 {
		return ""
	}

	message := item.Messages[len(item.Messages)-1]

	return message.AdaptiveCards[0].Body[0].Text
}

// GetType get type
func (r *Final) GetType() int {
	return r.Type
}

// GetMaxUnit get max user questions for current session
func (r *Final) GetMaxUnit() int {
	return r.Item.Throttling.MaxNumUserMessagesInConversation
}

// GetUserUnit get current question for current session
func (r *Final) GetUserUnit() int {
	return r.Item.Throttling.NumUserMessagesInConversation
}

func (r *Final) GetSuggestions() []*Suggestion {
	item := r.Item
	if len(item.Messages) == 0 {
		return nil
	}

	message := item.Messages[len(item.Messages)-1]

	return message.SuggestedResponses
}

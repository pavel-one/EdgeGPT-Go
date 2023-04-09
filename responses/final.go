package responses

import (
	"encoding/json"
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

	if message.AdaptiveCards == nil {
		return ""
	}
	if message.AdaptiveCards[0].Body == nil {
		return ""
	}

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

func (r *Final) GetSuggestions() []map[string]any {
	var out []map[string]any

	item := r.Item
	if len(item.Messages) == 0 {
		return nil
	}

	message := item.Messages[len(item.Messages)-1]

	for _, item := range message.SuggestedResponses {
		var m map[string]any

		b, err := json.Marshal(item)
		if err != nil {
			return nil
		}

		if err := json.Unmarshal(b, &m); err != nil {
			return nil
		}

		out = append(out, m)
	}

	return out
}

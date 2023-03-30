package EdgeGPT

import (
	"time"
)

type Suggestion struct {
	Text        string    `json:"text"`
	Author      string    `json:"author"`
	CreatedAt   time.Time `json:"createdAt"`
	Timestamp   time.Time `json:"timestamp"`
	MessageId   string    `json:"messageId"`
	MessageType string    `json:"messageType"`
	Offense     string    `json:"offense"`
	Feedback    struct {
		Tag       interface{} `json:"tag"`
		UpdatedOn interface{} `json:"updatedOn"`
		Type      string      `json:"type"`
	} `json:"feedback"`
	ContentOrigin string      `json:"contentOrigin"`
	Privacy       interface{} `json:"privacy"`
}

type MessageResponse struct {
	Text          string `json:"text"`
	Author        string `json:"author"`
	AdaptiveCards []struct {
		Type    string `json:"type"`
		Version string `json:"version"`
		Body    []struct {
			Type string `json:"type"`
			Text string `json:"text"`
			Wrap bool   `json:"wrap"`
		} `json:"body"`
	} `json:"adaptiveCards,omitempty"`
	SuggestedResponses []*Suggestion `json:"suggestedResponses,omitempty"`
}

type Response struct {
	Type         int    `json:"type"`
	InvocationId string `json:"invocationId"`
	Item         struct {
		Messages               []*MessageResponse `json:"messages"`
		ConversationExpiryTime time.Time          `json:"conversationExpiryTime,omitempty"`
		Throttling             struct {
			MaxNumUserMessagesInConversation int `json:"maxNumUserMessagesInConversation"`
			NumUserMessagesInConversation    int `json:"numUserMessagesInConversation"`
		} `json:"throttling"`
	} `json:"item"`
}

type MessageWrapper struct {
	Final    bool
	Question string
	Answer   *MessageResponse
	ch       chan *MessageResponse
}

func NewMessageWrapper(question string) *MessageWrapper {
	return &MessageWrapper{
		Question: question,
		ch:       make(chan *MessageResponse, 1),
	}
}

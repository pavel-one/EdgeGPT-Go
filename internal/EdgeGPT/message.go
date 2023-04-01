package EdgeGPT

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const (
	TypeUpdate float64 = 1
	TypeFinish float64 = 2
)

type UndefinedResponse struct {
	Type int `json:"type"`
}

func (r *UndefinedResponse) GetAnswer() string {
	return ""
}

func (r *UndefinedResponse) GetType() int {
	return r.Type
}

type UpdateResponse struct {
	Type      int    `json:"type"`
	Target    string `json:"target"`
	Arguments []struct {
		Cursor struct {
			J string `json:"j"`
			P int    `json:"p"`
		} `json:"cursor"`
		Messages  []MessageResponse `json:"messages"`
		RequestId string            `json:"requestId"`
	} `json:"arguments"`
}

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
	InvocationId string `json:"invocationId,omitempty"`
	Item         struct {
		Messages               []*MessageResponse `json:"messages"`
		ConversationExpiryTime time.Time          `json:"conversationExpiryTime,omitempty"`
		Throttling             struct {
			MaxNumUserMessagesInConversation int `json:"maxNumUserMessagesInConversation"`
			NumUserMessagesInConversation    int `json:"numUserMessagesInConversation"`
		} `json:"throttling"`
	} `json:"item"`
}

func (u *UpdateResponse) GetAnswer() string {
	arg := u.Arguments[0]
	if len(arg.Messages) == 0 {
		return ""
	}

	message := arg.Messages[len(arg.Messages)-1]

	return message.Text
}

func (u *UpdateResponse) GetType() int {
	return u.Type
}

func (r *Response) GetAnswer() string {
	item := r.Item
	if len(item.Messages) == 0 {
		return ""
	}

	message := item.Messages[len(item.Messages)-1]

	return message.Text
}

func (r *Response) GetType() int {
	return r.Type
}

type MessageWrapper struct {
	Final    bool
	Question string
	Answer   GptResponse
	Chan     chan []byte
	mu       *sync.Mutex
	conn     *websocket.Conn
}

func NewMessageWrapper(question string, mutex *sync.Mutex, conn *websocket.Conn) *MessageWrapper {
	return &MessageWrapper{
		Question: question,
		Chan:     make(chan []byte, 1),
		mu:       mutex,
		conn:     conn,
	}
}

func (m *MessageWrapper) Worker() error {
	defer m.mu.Unlock()

	var response map[string]any
	var updateResponse UpdateResponse
	var finishResponse Response
	var undefinedResponse UndefinedResponse

	for {
		var message []byte
		_, original, err := m.conn.ReadMessage()
		if err != nil {
			return err
		}

		// read to delimiter
		for _, b := range original {
			if b == DelimiterByte {
				break
			}

			message = append(message, b)
		}

		if err := json.Unmarshal(message, &response); err != nil {
			return err
		}

		switch response["type"] {
		case TypeUpdate:
			if err := json.Unmarshal(message, &updateResponse); err != nil {
				return err
			}

			m.Answer = &updateResponse
			break
		case TypeFinish:
			if err := json.Unmarshal(message, &finishResponse); err != nil {
				return err
			}

			m.Answer = &finishResponse
			m.Final = true
			m.Chan <- message
			close(m.Chan)
			return nil

		default:
			if err := json.Unmarshal(message, &undefinedResponse); err != nil {
				return err
			}

			m.Answer = &undefinedResponse
		}

		m.Chan <- message
	}

}

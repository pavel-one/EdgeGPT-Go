package EdgeGPT

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pavel-one/EdgeGPT-Go/internal/Helpers"
	"github.com/pavel-one/EdgeGPT-Go/responses"
	"sync"
)

type Hub struct {
	conversation *Conversation
	conn         *websocket.Conn
	InvocationId int
	mu           sync.Mutex
}

// initialHandshake request for initial session
func (c *Hub) initialHandshake() error {
	message := []byte("{\"protocol\": \"json\", \"version\": 1}" + Delimiter)

	if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}

	if _, _, err := c.conn.ReadMessage(); err != nil { //wait initial
		return err
	}

	c.InvocationId = 0

	return nil
}

// send new message to websocket
func (c *Hub) send(style, message string) (*responses.MessageWrapper, error) {
	c.mu.Lock()

	m, err := json.Marshal(c.getRequest(style, message))
	if err != nil {
		return nil, err
	}

	m = append(m, DelimiterByte)

	if err := c.conn.WriteMessage(websocket.TextMessage, m); err != nil {
		return nil, err
	}

	return responses.NewMessageWrapper(message, &c.mu, c.conn), nil
}

// Close hub and connection
// TODO: Use this!
func (c *Hub) Close() {
	c.conn.Close()
}

// getRequest generate struct for new request websocket
func (c *Hub) getRequest(style, message string) map[string]any {
	m := map[string]any{
		"invocationId": string(rune(c.InvocationId)),
		"target":       "chat",
		"type":         4,
		"arguments": []map[string]any{
			{
				"source": "cib",
				"optionsSets": []string{
					"nlu_direct_response_filter",
					"deepleo",
					"disable_emoji_spoken_text",
					"responsible_ai_policy_235",
					"enablemm",
					style,
					"dtappid",
					"cricinfo",
					"cricinfov2",
					"dv3sugg",
				},
				"sliceIds": []string{
					"222dtappid",
					"225cricinfo",
					"224locals0",
				},
				"traceId":          Helpers.RandomHex(32),
				"isStartOfSession": c.InvocationId == 0,
				"message": map[string]any{
					"author":      "user",
					"inputMethod": "Keyboard",
					"text":        message,
					"messageType": "Chat",
				},
				"conversationSignature": c.conversation.ConversationSignature,
				"participant": map[string]any{
					"id": c.conversation.ClientId,
				},
				"conversationId": c.conversation.ConversationId,
			},
		},
	}
	c.InvocationId++

	return m
}

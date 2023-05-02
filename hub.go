package EdgeGPT

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/pavel-one/EdgeGPT-Go/config"
	"github.com/pavel-one/EdgeGPT-Go/internal/Helpers"
	"github.com/pavel-one/EdgeGPT-Go/responses"
	"net/url"
	"sync"
)

type Hub struct {
	conversation *Conversation
	conn         *websocket.Conn
	wssUrl       *url.URL
	headers      map[string]string
	InvocationId int
	mu           sync.Mutex
}

func NewHub(conversation *Conversation, config *config.GPT) (*Hub, error) {
	if conversation == nil {
		return nil, errors.New("not set conversation")
	}

	h := &Hub{
		conversation: conversation,
		conn:         nil,
		wssUrl:       config.WssUrl,
		headers:      config.Headers,
	}

	conn, err := h.NewConnect()
	if err != nil {
		return nil, err
	}
	h.conn = conn

	log.Infoln("New hub for conversation:", conversation.ConversationId)

	return h, nil
}

// NewConnect create new websocket connection
func (h *Hub) NewConnect() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(h.wssUrl.String(), Helpers.GetHeaders(h.headers))
	if err != nil {
		return nil, err
	}

	message := []byte("{\"protocol\": \"json\", \"version\": 1}" + Delimiter)
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		return nil, err
	}
	if _, _, err := conn.ReadMessage(); err != nil { //wait initial
		return nil, err
	}

	return conn, nil
}

// CheckAndReconnect check active connection and reconnect
func (h *Hub) CheckAndReconnect() error {
	if h.conn == nil {
		return errors.New("not set connection")
	}

	if err := h.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
		log.Infoln("Reconnection")

		h.Close()
		h.conn = nil

		conn, err := h.NewConnect()
		if err != nil {
			return err
		}
		h.conn = conn
	}

	return nil
}

// send new message to websocket
func (h *Hub) send(style, message string) (*responses.MessageWrapper, error) {
	if h.conn == nil {
		return nil, errors.New("not set connection")
	}
	h.mu.Lock()

	if err := h.CheckAndReconnect(); err != nil {
		return nil, err
	}

	m, err := json.Marshal(h.getRequest(style, message))
	if err != nil {
		return nil, err
	}

	m = append(m, DelimiterByte)

	if err := h.conn.WriteMessage(websocket.TextMessage, m); err != nil {
		return nil, err
	}

	return responses.NewMessageWrapper(message, &h.mu, h.conn), nil
}

// Close hub and connection
// TODO: Use this!
func (h *Hub) Close() {
	if h.conn == nil {
		return
	}

	log.Infoln("Close connection")
	h.conn.Close()
}

// getRequest generate struct for new request websocket
func (h *Hub) getRequest(style, message string) map[string]any {
	switch style {
	case "creative":
		style = StyleCreative
		break
	case "balanced":
		style = StyleBalanced
		break
	case "precise":
		style = StylePrecise
		break
	case StyleCreative:
		style = StyleCreative
		break
	case StyleBalanced:
		style = StyleBalanced
		break
	case StylePrecise:
		style = StylePrecise
		break
	default:
		style = StyleBalanced
	}

	m := map[string]any{
		"invocationId": string(rune(h.InvocationId)),
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
				"isStartOfSession": h.InvocationId == 0,
				"message": map[string]any{
					"author":      "user",
					"inputMethod": "Keyboard",
					"text":        message,
					"messageType": "Chat",
				},
				"conversationSignature": h.conversation.ConversationSignature,
				"participant": map[string]any{
					"id": h.conversation.ClientId,
				},
				"conversationId": h.conversation.ConversationId,
			},
		},
	}
	h.InvocationId++

	return m
}

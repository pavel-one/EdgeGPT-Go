package EdgeGPT

import (
	"EdgeGPT-Go/internal/helpers"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Hub struct {
	conversation *Conversation
	conn         *websocket.Conn
	InvocationId int
}

func (c *Hub) initialHandshake() error {
	message := []byte("{\"protocol\": \"json\", \"version\": 1}" + Delimiter)

	if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}

	if _, _, err := c.conn.ReadMessage(); err != nil { //wait initial
		return err
	}

	c.InvocationId = 0

	go c.worker()

	return nil
}

func (c *Hub) send(message string) error {
	m, err := json.Marshal(c.getRequest(message))
	if err != nil {
		return err
	}

	m = append(m, DelimiterByte)

	if err := c.conn.WriteMessage(websocket.TextMessage, m); err != nil {
		return err
	}

	return nil
}

func (c *Hub) worker() {
	// читаем ответы от сервера
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			// обработка ошибки
			log.Println("Ошибка чтения сообщения:", err)
			return
		}
		fmt.Printf("Получено сообщение типа %d: %s\n", messageType, message)
	}
}

func (c *Hub) Close() {
	//c.InvocationId = 0 TODO: ?
	c.conn.Close()
}

func (c *Hub) getRequest(message string) map[string]any {
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
					StyleBalanced,
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
				"traceId":          helpers.RandomHex(32),
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

package EdgeGPT

import "gopkg.in/guregu/null.v4"

type Conversation struct {
	ConversationId        string `json:"conversationId,omitempty"`
	ClientId              string `json:"clientId,omitempty"`
	ConversationSignature string `json:"conversationSignature,omitempty"`
	Result                struct {
		Value   null.String `json:"value"`
		Message null.String `json:"message"`
	} `json:"result"`
}

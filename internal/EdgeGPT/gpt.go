package EdgeGPT

import (
	"EdgeGPT-Go/config"
	"EdgeGPT-Go/internal/helpers"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"os"
)

const (
	StyleCreative = "h3relaxedimg"
	StyleBalanced = "galileo"
	StylePrecise  = "h3precise"
	Delimiter     = "\x1e"
	DelimiterByte = uint8(30)
)

type GPT struct {
	Config       *config.GPT
	client       *http.Client
	cookies      []*http.Cookie
	Conversation *Conversation
}

func NewGPT(conf *config.GPT) (*GPT, error) {
	cookieFile, err := os.Open(conf.CookieFileName)
	if err != nil {
		return nil, err
	}
	defer cookieFile.Close()

	cookiesJSON, err := io.ReadAll(cookieFile)
	if err != nil {
		return nil, err
	}

	var parse []map[string]any
	err = json.Unmarshal(cookiesJSON, &parse)
	if err != nil {
		return nil, err
	}

	gpt := &GPT{
		Config:  conf,
		cookies: helpers.MapToCookies(parse),
		client: &http.Client{
			Timeout: conf.TimeoutRequest,
		},
	}

	if err := gpt.createConversation(); err != nil {
		return nil, err
	}

	return gpt, nil
}

// createConversation request for getting new dialog
func (g *GPT) createConversation() error {
	req, err := http.NewRequest("GET", g.Config.ConversationUrl.String(), nil)

	for k, v := range g.Config.Headers {
		req.Header.Set(k, v)
	}

	if err != nil {
		return err
	}

	for _, cookie := range g.cookies {
		req.AddCookie(cookie)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code not ok: %d, %s", resp.StatusCode, resp.Status)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	conversation := new(Conversation)
	if err := json.Unmarshal(r, conversation); err != nil {
		return err
	}

	if conversation.Result.Value.ValueOrZero() != "Success" {
		return nil
	}

	g.Conversation = conversation

	return nil
}

func (g *GPT) Ask(message string) error {
	if len(message) > 2000 {
		return fmt.Errorf("message very long, max: %d", 2000)
	}

	hub, err := g.createHub()
	if err != nil {
		return err
	}

	if err := hub.send(message); err != nil {
		return err
	}

	return nil
}

func (g *GPT) createHub() (*Hub, error) {
	if g.Conversation == nil {
		return nil, errors.New("not set conversation")
	}

	conn, _, err := websocket.DefaultDialer.Dial(g.Config.WssUrl.String(), helpers.GetHeaders(g.Config.Headers))
	if err != nil {
		return nil, err
	}

	h := &Hub{
		conversation: g.Conversation,
		conn:         conn,
	}

	if err := h.initialHandshake(); err != nil {
		return nil, err
	}

	return h, nil
}

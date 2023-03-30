package EdgeGPT

import (
	"EdgeGPT-Go/config"
	"EdgeGPT-Go/internal/helpers"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GPT struct {
	Config  *config.GPT
	client  *http.Client
	cookies []*http.Cookie
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

	return &GPT{
		Config:  conf,
		cookies: helpers.MapToCookies(parse),
		client: &http.Client{
			Timeout: conf.TimeoutRequest,
		},
	}, nil
}

func (g *GPT) NewConversation() (*Conversation, error) {
	req, err := http.NewRequest("GET", g.Config.ConversationUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	for _, cookie := range g.cookies {
		req.AddCookie(cookie)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code not ok: %d, %s", resp.StatusCode, resp.Status)
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	conversation := new(Conversation)
	if err := json.Unmarshal(r, conversation); err != nil {
		return nil, err
	}

	if conversation.Result.Value.ValueOrZero() != "Success" {
		return nil, fmt.Errorf("not valid cookies: %s", conversation.Result.Message.ValueOrZero())
	}

	return conversation, nil
}

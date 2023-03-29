package EdgeGPT

import (
	"EdgeGPT-Go/internal/helpers"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

type Conversation struct {
	client  *http.Client
	cookies []*http.Cookie
}

func NewConversation() (*Conversation, error) {
	cookieFile, err := os.Open("cookies.json")
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

	return &Conversation{
		cookies: helpers.MapToCookies(parse),
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}, nil
}

func (c *Conversation) Get() ([]byte, error) {
	req, err := http.NewRequest("GET", BING_PROXY_URL, nil)
	if err != nil {
		return nil, err
	}

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return r, nil
}

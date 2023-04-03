package EdgeGPT

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pavel-one/EdgeGPT-Go/config"
	"github.com/pavel-one/EdgeGPT-Go/internal/CookieManager"
	"github.com/pavel-one/EdgeGPT-Go/internal/Helpers"
	"github.com/pavel-one/EdgeGPT-Go/internal/Logger"
	"io"
	"net/http"
	"time"
)

var log = Logger.NewLogger("GPT Service")

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
	ExpiredAt    time.Time
	Hub          *Hub
}

// NewGPT create new service
func NewGPT(conf *config.GPT) (*GPT, error) {
	manager, err := CookieManager.NewManager()
	if err != nil {
		return nil, err
	}

	gpt := &GPT{
		Config:    conf,
		cookies:   Helpers.MapToCookies(manager.GetBestCookie()),
		ExpiredAt: time.Now().Add(time.Minute * 120),
		client: &http.Client{
			Timeout: conf.TimeoutRequest,
		},
	}

	if err := gpt.createConversation(); err != nil {
		return nil, err
	}

	hub, err := gpt.createHub()
	if err != nil {
		return nil, err
	}
	gpt.Hub = hub

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

	log.Infoln("New conversation", conversation)

	return nil
}

/*
AskAsync getting answer async:
Example:

	gpt, err := EdgeGPT.NewGPT(conf) //create service
	if err != nil {
		log.Fatalln(err)
	}

	mw, err := gpt.AskAsync("Привет, ты живой?") // send ask to gpt
	if err != nil {
		log.Fatalln(err)
	}

	go mw.Worker() // Run reading websocket messages

	for _ = range mw.Chan {
		// update answer
		log.Println(mw.Answer.GetAnswer())
	}
*/
func (g *GPT) AskAsync(message string) (*MessageWrapper, error) {

	if len(message) > 2000 {
		return nil, fmt.Errorf("message very long, max: %d", 2000)
	}

	log.Infoln("New ask:", message)
	return g.Hub.send(message)
}

// AskSync getting answer sync
func (g *GPT) AskSync(message string) (*MessageWrapper, error) {
	if len(message) > 2000 {
		return nil, fmt.Errorf("message very long, max: %d", 2000)
	}

	m, err := g.Hub.send(message)
	if err != nil {
		return nil, err
	}

	go m.Worker()

	for _ = range m.Chan {
		if m.Final {
			break
		}
	}

	log.Infoln("New ask:", message)
	return m, nil
}

// createHub create websocket hub
func (g *GPT) createHub() (*Hub, error) {
	if g.Conversation == nil {
		return nil, errors.New("not set conversation")
	}

	conn, _, err := websocket.DefaultDialer.Dial(g.Config.WssUrl.String(), Helpers.GetHeaders(g.Config.Headers))
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

	log.Infoln("New hub for conversation:", g.Conversation.ConversationId)
	return h, nil
}

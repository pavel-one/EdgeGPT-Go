package config

import (
	"EdgeGPT-Go/internal/helpers"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"time"
)

type GPT struct {
	ConversationUrl *url.URL
	WssUrl          *url.URL
	TimeoutRequest  time.Duration
	Headers         map[string]string
}

func NewGpt() (*GPT, error) {
	cu, err := url.Parse("https://edgeservices.bing.com/edgesvc/turing/conversation/create")
	if err != nil {
		return nil, err
	}

	wss, err := url.Parse("wss://sydney.bing.com/sydney/ChatHub")
	if err != nil {
		return nil, err
	}

	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	forwared := fmt.Sprintf(
		"13.%d.%d.%d",
		helpers.RandInt(104, 107),
		helpers.RandInt(0, 255),
		helpers.RandInt(0, 255))

	return &GPT{
		ConversationUrl: cu,
		WssUrl:          wss,
		TimeoutRequest:  time.Second * 30,
		Headers: map[string]string{
			"accept":                      "application/json",
			"accept-language":             "en-US,en;q=0.9",
			"content-type":                "application/json",
			"sec-ch-ua":                   "\"Not_A Brand\";v=\"99\", \"Microsoft Edge\";v=\"110\", \"Chromium\";v=\"110\"",
			"sec-ch-ua-arch":              "\"x86\"",
			"sec-ch-ua-bitness":           "\"64\"",
			"sec-ch-ua-full-version":      "\"109.0.1518.78\"",
			"sec-ch-ua-full-version-list": "\"Chromium\";v=\"110.0.5481.192\", \"Not A(Brand\";v=\"24.0.0.0\", \"Microsoft Edge\";v=\"110.0.1587.69\"",
			"sec-ch-ua-mobile":            "?0",
			"sec-ch-ua-model":             "",
			"sec-ch-ua-platform":          "\"Windows\"",
			"sec-ch-ua-platform-version":  "\"15.0.0\"",
			"sec-fetch-dest":              "empty",
			"sec-fetch-mode":              "cors",
			"sec-fetch-site":              "same-origin",
			"x-ms-client-request-id":      uid.String(),
			"x-ms-useragent":              "azsdk-js-api-client-factory/1.0.0-beta.1 core-rest-pipeline/1.10.0 OS/Win32",
			"Referer":                     "https://www.bing.com/search?q=Bing+AI&showconv=1&FORM=hpcodx",
			"Referrer-Policy":             "origin-when-cross-origin",
			"x-forwarded-for":             forwared,
		},
	}, nil
}

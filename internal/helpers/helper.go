package helpers

import (
	"math/rand"
	"net/http"
)

func RandInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func MapToCookies(m []map[string]any) []*http.Cookie {
	var cookies []*http.Cookie

	for _, item := range m {
		cookies = append(cookies, &http.Cookie{
			Name:  item["name"].(string),
			Value: item["value"].(string),
			//Expires:  time.Now().Add(time.Hour * 8766),
			//Path:     item["path"].(string),
			//Domain:   item["domain"].(string),
			//Secure:   item["secure"].(bool),
			//HttpOnly: item["httpOnly"].(bool),
		})
	}

	return cookies
}

package helpers

import (
	"crypto/rand"
	"encoding/hex"
	mrand "math/rand"
	"net/http"
)

func RandInt(min int, max int) int {
	return mrand.Intn(max-min+1) + min
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

func GetHeaders(m map[string]string) http.Header {
	headers := http.Header{}

	for k, v := range m {
		headers.Set(k, v)
	}

	return headers
}

//func RandomHex(length int) string {
//	bytes := make([]byte, length)
//	io.ReadFull(rand.Reader, bytes)
//	hexString := hex.EncodeToString(bytes)
//
//	return hexString
//}

func RandomHex(length int) string {
	buf := make([]byte, length/2)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}

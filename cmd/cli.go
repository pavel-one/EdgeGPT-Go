package main

import (
	"EdgeGPT-Go/internal/EdgeGPT"
	"log"
)

func main() {
	c, err := EdgeGPT.NewConversation()
	if err != nil {
		log.Fatalln(err)
	}

	get, err := c.Get()
	if err != nil {
		log.Fatalln(err)
	}
	resp := string(get)

	log.Println(resp)
}

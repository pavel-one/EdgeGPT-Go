package main

import (
	"EdgeGPT-Go/config"
	"EdgeGPT-Go/internal/EdgeGPT"
	"log"
)

func main() {
	conf, err := config.NewGpt()
	if err != nil {
		log.Fatalln(err)
	}

	gpt, err := EdgeGPT.NewGPT(conf)
	if err != nil {
		log.Fatalln(err)
	}

	c, err := gpt.NewConversation()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(c)
}

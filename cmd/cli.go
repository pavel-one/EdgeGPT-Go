package main

import (
	"EdgeGPT-Go/config"
	"EdgeGPT-Go/internal/EdgeGPT"
	"log"
	"time"
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

	err = gpt.Ask("Привет, ты живой?")
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(time.Minute * 5)
}

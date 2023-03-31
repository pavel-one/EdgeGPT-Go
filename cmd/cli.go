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

	mw, err := gpt.AskAsync("Привет, ты живой?")
	if err != nil {
		log.Fatalln(err)
	}

	go mw.Worker()

	for _ = range mw.Chan {
		log.Println(mw.Answer.GetAnswer())
	}

	as, err := gpt.AskSync("Какая погода в Ростове-на-Дону?")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(as.Answer.GetAnswer())

	time.Sleep(time.Minute * 5)
}

package main

import (
	"github.com/pavel-one/EdgeGPT-Go/config"
	"github.com/pavel-one/EdgeGPT-Go/internal/EdgeGPT"
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
		log.Infoln(mw.Answer.GetAnswer())
	}

	as, err := gpt.AskSync("Покажи пример сокетов на golang gorilla")
	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln(as.Answer.GetAnswer())

	time.Sleep(time.Minute * 5)
}

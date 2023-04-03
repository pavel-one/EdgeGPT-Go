package main

import (
	"github.com/pavel-one/EdgeGPT-Go/internal/EdgeGPT"
)

func main() {
	s := EdgeGPT.NewStorage()

	gpt, err := s.GetOrSet("any-key")
	if err != nil {
		log.Fatalln(err)
	}

	// send ask async
	mw, err := gpt.AskAsync("Hi, you're alive?")
	if err != nil {
		log.Fatalln(err)
	}

	go mw.Worker() // start worker

	for _ = range mw.Chan {
		// update answer
		log.Infoln(mw.Answer.GetAnswer())
		log.Infoln(mw.Answer.GetType())
		log.Infoln(mw.Answer.GetSuggestions())
		log.Infoln(mw.Answer.GetMaxUnit())
		log.Infoln(mw.Answer.GetUserUnit())
	}

	// send sync ask
	as, err := gpt.AskSync("Show an example of sockets in golang gorilla")
	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln(as.Answer.GetAnswer())
}

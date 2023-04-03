## EdgeGPT-Go

EdgeGPT-Go is a New Bing unofficial API developed using Golang.  
You can use it as a library, microservice or standalone cli application.  
The package supports multiple cookies.
The package supports multiple cookies. As well as rapid deployment as a microservice via docker.

## Feature:
- [x] GRPC interface
- [x] Library interface
- [x] Sync/Async request
- [ ] CLI interface

## How to use it:

### Getting authentication (Required)
- Install the cookie editor extension for [Chrome](https://chrome.google.com/webstore/detail/cookie-editor/hlkenndednhfkekhgcdicdfddnkalmdm) or [Firefox](https://addons.mozilla.org/en-US/firefox/addon/cookie-editor/)
- Go to `bing.com` and login
- Open the extension
- Click "Export" on the bottom right, then "Export as JSON" (This saves your cookies to clipboard)
- Create folder `cookies`
- Paste your cookies into a file `1.json`

If you have several accounts - repeat for each of them and save to the `cookies` folder

### Use as a library
`go get github.com/pavel-one/EdgeGPT-Go`   

```go
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

```
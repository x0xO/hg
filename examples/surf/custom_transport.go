package main

import (
	"log"
	"time"

	"github.com/x0xO/hg/surf"
)

func main() {
	cli := surf.NewClient()

	// transport custom settings
	cli.GetTransport().TLSHandshakeTimeout = time.Nanosecond

	_, err := cli.Get("https://google.com").Do()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/x0xO/hg/surf"
)

func main() {
	r, err := surf.NewClient().
		SetOptions(surf.NewOptions().DisableKeepAlive()).
		Get("http://www.keycdn.com").
		Do()
	if err != nil {
		log.Fatal(err)
	}

	r.Debug().Response().Print() // Connection: close
}

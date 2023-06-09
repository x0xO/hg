package main

import (
	"fmt"
	"log"
	"time"

	"github.com/x0xO/hg/surf"
)

func main() {
	r, err := surf.NewClient().
		SetOptions(surf.NewOptions().Timeout(time.Second)).
		Get("httpbingo.org/delay/2").
		Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.StatusCode)
}

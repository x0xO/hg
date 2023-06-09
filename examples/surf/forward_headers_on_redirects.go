package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hg/surf"
)

func main() {
	opt := surf.NewOptions().ForwardHeadersOnRedirect()

	r, err := surf.NewClient().
		SetOptions(opt).
		Get("google.com").
		AddHeaders(map[string]string{"Referer": "surf.xoxo"}).
		Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Referer())
}

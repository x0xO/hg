package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hg/surf"
)

func main() {
	r, err := surf.NewClient().
		SetOptions(surf.NewOptions().HTTP2()).
		Get("https://http2.pro/api/v1").
		Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Proto)

	r.Debug().Request().Response(true).Print()
}

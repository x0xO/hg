package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hg/surf"
)

func main() {
	r, err := surf.NewClient().
		SetOptions(surf.NewOptions().CacheBody()).
		Get("http://httpbingo.org/get").
		Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Body.Limit(10))
	fmt.Println(r.Body) // print cached body
}

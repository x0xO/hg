package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hg/surf"
)

func main() {
	r, err := surf.NewClient().Head("http://httpbingo.org/get").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Status)
}

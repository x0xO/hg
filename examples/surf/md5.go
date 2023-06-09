package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hg/surf"
)

func main() {
	r, err := surf.NewClient().Get("http://google.com").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Body.MD5())
}

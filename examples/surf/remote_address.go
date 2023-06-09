package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hg/surf"
)

func main() {
	// to get remote server ip address
	opt := surf.NewOptions().GetRemoteAddress()

	r, err := surf.NewClient().SetOptions(opt).Get("ya.ru").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.RemoteAddress())
}

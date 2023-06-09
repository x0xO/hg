package main

import (
	"fmt"

	"github.com/x0xO/hg/surf"
)

func main() {
	r, _ := surf.NewClient().Get("https://httpbingo.org/encoding/utf8").Do()
	fmt.Println(r.Body)

	r, _ = surf.NewClient().Get("http://vk.com").Do()
	fmt.Println(r.Body.UTF8())
}

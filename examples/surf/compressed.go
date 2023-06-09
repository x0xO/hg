package main

import (
	"fmt"

	"github.com/x0xO/hg/surf"
)

func main() {
	URL := "https://httpbin.org/gzip"
	r, _ := surf.NewClient().Get(URL).Do()
	fmt.Println(r.Body)

	URL = "https://httpbin.org/deflate"
	r, _ = surf.NewClient().Get(URL).Do()
	fmt.Println(r.Body)
}

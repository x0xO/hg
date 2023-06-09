package main

import (
	"fmt"

	"github.com/x0xO/hg/surf"
)

func main() {
	URL := "https://httpbingo.org/get"

	cli := surf.NewClient()
	req := cli.Get(URL)

	resp, _ := req.Do()

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Body)
	fmt.Println(resp.Cookies)
	fmt.Println(resp.Headers)
	fmt.Println(resp.URL)
}

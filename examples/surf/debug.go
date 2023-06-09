package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hg/surf"
)

func main() {
	opt := surf.NewOptions()
	cli := surf.NewClient().SetOptions(opt)
	req := cli.Get("https://httpbingo.org/get", map[string]string{"hg": "surf"})

	r, err := req.Do()
	if err != nil {
		log.Fatal(err)
	}

	d := r.Debug()
	d.Request(true) // true for verbose output with request body if set
	d.Response()    // true for verbose output with response body

	d.Print()

	fmt.Println(r.Time)
}

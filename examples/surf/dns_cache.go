package main

import (
	"log"
	"time"

	"github.com/x0xO/hg/surf"
)

func main() {
	opt := surf.NewOptions().DNSCache(time.Second*30, 10)

	cli := surf.NewClient().SetOptions(opt) // separate client to reuse client and DNS cache
	url := "httpbingo.org/get"

	r, err := cli.Get(url).Do() // cache the ip of the DNS response

	for i := 0; i < 10; i++ {
		r, err = cli.Get(url).Do() // use DNS cache
	}

	if err != nil {
		log.Fatal(err)
	}

	r.Debug().DNSStat().Print()
}

package main

import (
	"context"
	"log"
	"time"

	"github.com/x0xO/hg/surf"
)

func main() {
	URL := "https://httpbingo.org/get"

	cli := surf.NewClient()
	req := cli.Get(URL)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	resp, err := req.WithContext(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.Body)
}

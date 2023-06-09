package main

import (
	"fmt"
	"log"
	"time"

	"github.com/x0xO/hg/surf"
)

func main() {
	// opt := surf.NewOptions().Retry(3)
	opt := surf.NewOptions().Retry(2, time.Millisecond*50)

	for i := 0; i < 3; i++ {
		r, err := surf.NewClient().SetOptions(opt).Get("http://httpbingo.org/unstable").Do()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("StatusCode:", r.StatusCode, "Attempts:", r.Attempts)
		r.Debug().Request().Print()
	}
}

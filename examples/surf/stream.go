package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/x0xO/hg/surf"
)

func main() {
	r, err := surf.NewClient().Get("https://httpbingo.org/stream/10").Do()
	if err != nil {
		log.Fatal(err)
	}

	for {
		line, err := r.Body.Stream().ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
		}

		log.Println(line)
		time.Sleep(time.Second * 1)
	}

	// var bytesRead int
	// buffer := make([]byte, 4096)

	// for {
	// 	n, err := r.Body.Stream().Read(buffer)
	// 	bytesRead += n

	// 	if err == io.EOF {
	// 		break
	// 	}

	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	log.Println(string(buffer))
	// 	time.Sleep(time.Second * 1)
	// }
}

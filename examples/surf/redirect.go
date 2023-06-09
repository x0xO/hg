package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hg/surf"
)

func main() {
	opt := surf.NewOptions()
	opt.MaxRedirects(4)

	// or custom redirect policy

	// opt.RedirectPolicy(
	// 	func(req *http.Request, via []*http.Request) error {
	// 		if len(via) >= 4 {
	// 			return fmt.Errorf("stopped after %d redirects", 4)
	// 		}
	// 		return nil
	// 	},
	// )

	r, err := surf.NewClient().SetOptions(opt).Get("https://httpbingo.org/redirect/6").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.StatusCode)
}

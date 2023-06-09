package main

import (
	"fmt"

	"github.com/x0xO/hg"
)

func main() {
	d := hg.NewHDir(".").CopyDir("copy")

	if d.Error() != nil {
		fmt.Println(d.Error())
	}
}

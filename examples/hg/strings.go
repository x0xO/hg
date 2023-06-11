package main

import (
	"fmt"

	"github.com/x0xO/hg"
)

func main() {
	s := hg.NewHString("привет")
	fmt.Println(s.LeftJustify(20, "_"))
	fmt.Println(s.RightJustify(20, "_"))
	fmt.Println(s.Center(20, "_"))
}

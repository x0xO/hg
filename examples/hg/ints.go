package main

import (
	"fmt"

	"github.com/x0xO/hg"
)

func main() {
	fmt.Println(hg.NewHInt(1).Random())
	fmt.Println(hg.NewHInt().RandomRange(-10, -5))
	fmt.Println(hg.NewHInt().RandomRange(-10, 5))
}

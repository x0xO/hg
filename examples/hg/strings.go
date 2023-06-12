package main

import (
	"fmt"

	"github.com/x0xO/hg"
)

func main() {
	s := hg.NewHString("💛💚💙💜")

	fmt.Println(s.LeftJustify(10, "*"))
	fmt.Println(s.RightJustify(10, "*"))
	fmt.Println(s.Center(10, "*"))

	// 💛💚💙💜******
	// ******💛💚💙💜
	// ***💛💚💙💜***
}

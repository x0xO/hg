package main

import (
	"fmt"

	"github.com/x0xO/hg"
)

func main() {
	f := hg.NewHFloat(1.3339)

	md := f.Hash().MD5()
	fmt.Println(md)

	s := hg.HString("12.3348992")
	fmt.Println(s.HFloat().RoundDecimal(5))
}

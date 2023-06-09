package main

import (
	"fmt"

	"github.com/x0xO/hg"
)

func main() {
	sl := hg.HSliceOf(1, 2, 3, 4, 4, 2, 5)
	s := hg.HSetOf(sl...) // convert HSlice to HSet

	fmt.Println(s)

	s2 := hg.HSetOf(4, 5, 6, 7, 8)
	fmt.Println(s.SymmetricDifference(s2))

	set5 := hg.HSetOf(1, 2)
	set6 := hg.HSetOf(2, 3, 4)

	set7 := set5.Difference(set6)
	fmt.Println(set7)

	s = hg.HSetOf(1, 2, 3, 4, 5)
	even := s.Filter(func(val int) bool { return val%2 == 0 })
	fmt.Println(even)

	s = s.Remove(1)
	fmt.Println(s)
}

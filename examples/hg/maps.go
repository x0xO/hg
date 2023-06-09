package main

import (
	"fmt"

	"github.com/x0xO/hg"
	"github.com/x0xO/hg/pkg/iter"
)

func main() {
	md := hg.NewHMap[int, hg.HSlice[int]]()

	for i := range iter.N(5) {
		// md[i] = md.GetOrDefault(i, hg.NewHSlice[int]()).Append(i)
		md.Set(i, md.GetOrDefault(i, hg.NewHSlice[int]()).Append(i))
	}

	for i := range iter.N(10) {
		// md[i] = md.GetOrDefault(i, hg.NewHSlice[int]()).Append(i)
		md.Set(i, md.GetOrDefault(i, hg.NewHSlice[int]()).Append(i))
	}

	fmt.Println(md)

	mo := hg.HMapOrdFromHMap(md)
	fmt.Printf("mo: %v\n", mo)
}

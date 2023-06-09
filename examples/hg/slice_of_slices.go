package main

import (
	"fmt"

	"github.com/x0xO/hg"
)

func main() {
	ns1 := hg.NewHSlice[hg.HString]().Append("aaa")
	ns2 := hg.NewHSlice[hg.HString]().Append("bbb").Append("ccc")
	ns3 := hg.NewHSlice[hg.HString]().Append("ccc").Append("dddd").Append("wwwww")

	nx := hg.HSliceOf(ns3, ns2, ns1, ns2)

	fmt.Println(nx.Flatten())
	fmt.Println(nx.Flatten().Last().(hg.HString).ToUpper())

	nx = nx.Unique()

	fmt.Println(
		nx.SortBy(func(i, j int) bool { return nx[i].Get(0).String() < nx[j].Get(0).String() }),
	)
	fmt.Println(nx.SortBy(func(i, j int) bool { return nx[i].Len() < nx[j].Len() }))

	fmt.Println(nx.Reverse())

	fmt.Println(nx.Random())
	fmt.Println(nx.RandomSample(2))

	ch := nx.Chunks(2)           // return []HSlice[T]
	chunks := hg.HSliceOf(ch...) // make hslice chunks
	fmt.Println(chunks)

	pr := nx.Permutations()            // return []HSlice[T]
	permutations := hg.HSliceOf(pr...) // make hslice permutations
	fmt.Println(permutations)

	m := hg.NewHMap[string, hg.HSlice[hg.HSlice[hg.HString]]]()
	m.Set("one", nx)

	fmt.Println(m.Get("one").Last().Contains("aaa"))

	nestedSlice := hg.HSlice[any]{
		1,
		hg.HSliceOf(2, 3),
		"abc",
		hg.HSliceOf("def", "ghi"),
		hg.HSliceOf(4.5, 6.7),
	}

	fmt.Println(nestedSlice)           // Output: [1 [2 3] abc [def ghi] [4.5 6.7]]
	fmt.Println(nestedSlice.Flatten()) // Output: [1 2 3 abc def ghi 4.5 6.7]

	nestedSlice2 := hg.HSlice[any]{
		1,
		[]int{2, 3},
		"abc",
		hg.HSliceOf("awe", "som", "e"),
		[]string{"lol", "ov"},
		hg.HSliceOf(4.5, 6.7),
		[]float64{4.5, 6.7},
		map[string]string{"a": "ss"},
		hg.HSliceOf(hg.NewHMapOrd[int, int]().Set(1, 1), hg.NewHMapOrd[int, int]().Set(2, 2)),
	}

	// Output: [1 [2 3] abc [awe som e] [lol ov] [4.5 6.7] [4.5 6.7] map[a:ss]]
	fmt.Println(nestedSlice2)

	// Output: [1 2 3 abc awe som e lol ov 4.5 6.7 4.5 6.7 map[a:ss]]
	fmt.Println(nestedSlice2.Flatten())
}

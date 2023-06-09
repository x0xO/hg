package main

import (
	"fmt"

	"github.com/x0xO/hg"
	"github.com/x0xO/hg/pkg/deref"
	"github.com/x0xO/hg/pkg/iter"
)

func main() {
	// var mo *hg.HMapOrd[int string]
	// mo := &hg.HMapOrd[int, string]{}
	// mo := new(hg.HMapOrd[int, string])
	// mo := ref.Of(make(hg.HMapOrd[int, string], 0))

	md := hg.NewHMapOrd[int, hg.HSlice[int]]()

	for i := range iter.N(5) {
		md.Set(i, md.GetOrDefault(i, hg.NewHSlice[int]()).Append(i))
	}

	for i := range iter.N(10) {
		md.Set(i, md.GetOrDefault(i, hg.NewHSlice[int]()).Append(i))
	}

	fmt.Println(md)

	ms := hg.NewHMapOrd[hg.HInt, hg.HInt]()
	ms.Set(11, 99).Set(12, 2).Set(1, 22).Set(2, 32).Set(222, 2)

	ms1 := ms.Clone()

	ms1.Set(888, 000)
	ms1.Set(888, 300)

	if v, ok := ms1.Get(888); ok {
		fmt.Println(v)
	}

	if v, ok := ms1.Get(11); ok {
		fmt.Println(v)
	}

	ms1.Set(1, 223)
	fmt.Println(ms)
	fmt.Println(ms1)

	fmt.Println(ms.Eq(ms1))
	fmt.Println(ms.Contains(12))

	ms.ForEach(func(k, v hg.HInt) { fmt.Println(k, v) })

	ms = ms.Map(func(k, v hg.HInt) (hg.HInt, hg.HInt) { return k.Mul(2), v.Mul(2) })

	fmt.Println(ms)

	ms.Delete(12, 1, 222)
	fmt.Println(ms.Contains(12))

	msstr := hg.NewHMapOrd[hg.HString, hg.HString]()
	msstr.Set("aaa", "CCC").Set("ccc", "AAA").Set("bbb", "DDD").Set("ddd", "BBB")
	fmt.Println(msstr) // before sort

	msstr.SortBy(func(i, j int) bool { return deref.Of(msstr)[i].Key < deref.Of(msstr)[j].Key })
	fmt.Println(msstr) // after sort by key

	msstr.SortBy(
		func(i, j int) bool { return deref.Of(msstr)[i].Value < deref.Of(msstr)[j].Value },
	)

	fmt.Println(msstr) // after sort by value

	mss := hg.NewHMapOrd[hg.HInt, hg.HSlice[int]]()
	mss.Set(22, hg.HSlice[int]{4, 0, 9, 6, 7})
	mss.Set(11, hg.HSlice[int]{1, 2, 3, 4})
	fmt.Println(mss) // before sort

	mss.SortBy(func(i, j int) bool { return deref.Of(mss)[i].Key < deref.Of(mss)[j].Key })
	fmt.Println(mss) // after sort by key

	mss.SortBy(
		func(i, j int) bool { return deref.Of(mss)[i].Value.Get(1) < deref.Of(mss)[j].Value.Get(1) },
	)

	fmt.Println(mss) // after sort by value

	fmt.Println(hg.HMapOrdFromMap(mss.ToHMap().ToMap()))
}

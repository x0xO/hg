package main

import (
	"fmt"

	"github.com/x0xO/hg"
)

func main() {
	result := hg.HSlice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("%#v\n", result)

	result = result.RandomSample(5)

	fmt.Printf("%#v\n", result.Clone().Append(999))
	fmt.Printf("%#v\n", result)
	fmt.Printf("%#v\n", result.ToSlice())

	filled := hg.NewHSlice[int](10).Fill(88)

	fmt.Println(filled)

	slice := hg.HSlice[int]{1, 2, 3, 4, 5}.Cut(1, 4)
	fmt.Println(slice)

	// InPlace Methods
	sipl := hg.NewHSlice[int]()

	sipl.AppendInPlace(1)
	sipl.AppendInPlace(2)
	sipl.AppendInPlace(3)

	sipl.DeleteInPlace(1)
	sipl.Fill(999999)

	sipl.InsertInPlace(0, 22, 33, 44)
	sipl.AddUniqueInPlace(22, 22, 22, 33, 44, 55)

	fmt.Println(sipl)

	slice = hg.HSlice[int]{1, 2, 3}
	slice.MapInPlace(func(val int) int { return val * 2 })

	fmt.Println(slice)

	slice = hg.HSlice[int]{1, 2, 3, 4, 5}

	slice.FilterInPlace(func(val int) bool {
		return val%2 == 0
	})

	fmt.Println(slice)

	slicea := hg.HSlice[string]{"a", "b", "c", "d"}
	slicea.InsertInPlace(2, "e", "f")
	fmt.Println(slicea)

	slice = hg.HSlice[int]{1, 2, 0, 4, 0, 3, 0, 0, 0, 0}
	slice.FilterZeroValuesInPlace()
	slice.DeleteInPlace(0)
	fmt.Println(slice)

	sll := hg.NewHSlice[int](0, 100000)
	sll = sll.Append(1).Clip()

	fmt.Println(sll.Cap())
}

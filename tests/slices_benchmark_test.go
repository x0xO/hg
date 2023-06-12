package hg_test

import (
	"testing"

	"github.com/x0xO/hg"
	"github.com/x0xO/hg/pkg/iter"
)

// go test -bench=. -benchmem -count=4

func genSlice() hg.HSlice[hg.HString] {
	slice := hg.NewHSlice[hg.HString](0, 10000)
	for i := range iter.N(10000) {
		slice = slice.Append(hg.NewHInt(i).HString())
	}

	return slice
}

func BenchmarkMap(b *testing.B) {
	slice := genSlice()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Map(func(i hg.HString) hg.HString { return i.Enc().GzFlate() })
	}
}

func BenchmarkMapParallel(b *testing.B) {
	slice := genSlice()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.MapParallel(func(i hg.HString) hg.HString { return i.Enc().GzFlate() })
	}
}

func BenchmarkFilter(b *testing.B) {
	slice := genSlice()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.Filter(func(i hg.HString) bool { return i.Enc().GzFlate().Len()%2 == 0 })
	}
}

func BenchmarkFilterParallel(b *testing.B) {
	slice := genSlice()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		slice.FilterParallel(func(i hg.HString) bool { return i.Enc().GzFlate().Len()%2 == 0 })
	}
}

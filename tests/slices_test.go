package hg_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/x0xO/hg"
	"github.com/x0xO/hg/pkg/iter"
)

func TestPermutations(t *testing.T) {
	empty := hg.NewHSlice[int]()
	permsEmpty := empty.Permutations()
	expectedPermsEmpty := []hg.HSlice[int]{empty}

	if !reflect.DeepEqual(permsEmpty, expectedPermsEmpty) {
		t.Errorf("expected %v, but got %v", expectedPermsEmpty, permsEmpty)
	}

	slice1 := hg.HSliceOf(1)
	perms1 := slice1.Permutations()
	expectedPerms1 := []hg.HSlice[int]{slice1}

	if !reflect.DeepEqual(perms1, expectedPerms1) {
		t.Errorf("expected %v, but got %v", expectedPerms1, perms1)
	}

	slice2 := hg.HSliceOf("a", "b")
	perms2 := slice2.Permutations()
	expectedPerms2 := []hg.HSlice[string]{
		{"a", "b"},
		{"b", "a"},
	}

	if !reflect.DeepEqual(perms2, expectedPerms2) {
		t.Errorf("expected %v, but got %v", expectedPerms2, perms2)
	}

	slice3 := hg.HSliceOf(1.0, 2.0, 3.0)
	perms3 := slice3.Permutations()

	expectedPerms3 := []hg.HSlice[float64]{
		{1.0, 2.0, 3.0},
		{1.0, 3.0, 2.0},
		{2.0, 1.0, 3.0},
		{2.0, 3.0, 1.0},
		{3.0, 1.0, 2.0},
		{3.0, 2.0, 1.0},
	}

	if !reflect.DeepEqual(perms3, expectedPerms3) {
		t.Errorf("expected %v, but got %v", expectedPerms3, perms3)
	}
}

func TestHSliceInsert(t *testing.T) {
	empty := hg.NewHSlice[int]()
	empty = empty.Insert(0, 1, 2)

	expectedEmpty := hg.NewHSlice[int]().Append(1, 2)
	if !reflect.DeepEqual(empty, expectedEmpty) {
		t.Errorf("expected %v, but got %v", expectedEmpty, empty)
	}

	slice1 := hg.NewHSlice[int]().Append(3, 4)
	slice1 = slice1.Insert(0, 1, 2)

	expected1 := hg.NewHSlice[int]().Append(1, 2, 3, 4)
	if !reflect.DeepEqual(slice1, expected1) {
		t.Errorf("expected %v, but got %v", expected1, slice1)
	}

	slice2 := hg.NewHSlice[string]().Append("foo", "bar", "baz")
	slice2 = slice2.Insert(1, "qux", "quux")

	expected2 := hg.NewHSlice[string]().Append("foo", "qux", "quux", "bar", "baz")
	if !reflect.DeepEqual(slice2, expected2) {
		t.Errorf("expected %v, but got %v", expected2, slice2)
	}

	slice3 := hg.NewHSlice[float64]().Append(1.23, 4.56)
	slice3 = slice3.Insert(slice3.Len(), 7.89)

	expected3 := hg.NewHSlice[float64]().Append(1.23, 4.56, 7.89)
	if !reflect.DeepEqual(slice3, expected3) {
		t.Errorf("expected %v, but got %v", expected3, slice3)
	}
}

func TestHSliceToSlice(t *testing.T) {
	hsl := hg.NewHSlice[int]().Append(1, 2, 3, 4, 5)
	slice := hsl.ToSlice()

	if len(slice) != hsl.Len() {
		t.Errorf("Expected length %d, but got %d", hsl.Len(), len(slice))
	}

	for i, v := range hsl {
		if v != slice[i] {
			t.Errorf("Expected value %d at index %d, but got %d", v, i, slice[i])
		}
	}
}

func TestHSliceHMapHashedHInt(t *testing.T) {
	hsl := hg.HSlice[hg.HInt]{1, 2, 3, 4, 5}
	hmap := hsl.ToHMapHashed()

	if hmap.Len() != hsl.Len() {
		t.Errorf("Expected %d, got %d", hsl.Len(), hmap.Len())
	}

	for _, v := range hsl {
		if !hmap.Contains(v.Hash().MD5()) {
			t.Errorf("Expected %v, got %v", v, hmap[v.Hash().MD5()])
		}
	}
}

func TestHSliceHMapHashedHStrings(t *testing.T) {
	hsl := hg.HSlice[hg.HString]{"1", "2", "3", "4", "5"}
	hmap := hsl.ToHMapHashed()

	if hmap.Len() != hsl.Len() {
		t.Errorf("Expected %d, got %d", hsl.Len(), hmap.Len())
	}

	for _, v := range hsl {
		if !hmap.Contains(v.Hash().MD5()) {
			t.Errorf("Expected %v, got %v", v, hmap[v.Hash().MD5()])
		}
	}
}

func TestHSliceShuffle(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	hsl.Shuffle()

	if hsl.Len() != 10 {
		t.Error("Expected length of 10, got ", hsl.Len())
	}
}

func TestHSliceChunks(t *testing.T) {
	tests := []struct {
		name     string
		input    hg.HSlice[int]
		expected []hg.HSlice[int]
		size     int
	}{
		{
			name:     "empty slice",
			input:    hg.NewHSlice[int](),
			expected: []hg.HSlice[int]{},
			size:     2,
		},
		{
			name:     "single chunk",
			input:    hg.NewHSlice[int]().Append(1, 2, 3),
			expected: []hg.HSlice[int]{hg.NewHSlice[int]().Append(1, 2, 3)},
			size:     3,
		},
		{
			name:  "multiple chunks",
			input: hg.NewHSlice[int]().Append(1, 2, 3, 4, 5, 6),
			expected: []hg.HSlice[int]{
				hg.NewHSlice[int]().Append(1, 2),
				hg.NewHSlice[int]().Append(3, 4),
				hg.NewHSlice[int]().Append(5, 6),
			},
			size: 2,
		},
		{
			name:  "last chunk is smaller",
			input: hg.NewHSlice[int]().Append(1, 2, 3, 4, 5),
			expected: []hg.HSlice[int]{
				hg.NewHSlice[int]().Append(1, 2),
				hg.NewHSlice[int]().Append(3, 4),
				hg.NewHSlice[int]().Append(5),
			},
			size: 2,
		},
		{
			name:     "chunk size bigger than slice length",
			input:    hg.NewHSlice[int]().Append(1, 2, 3, 4),
			expected: []hg.HSlice[int]{hg.NewHSlice[int]().Append(1, 2, 3, 4)},
			size:     5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Chunks(tt.size)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d chunks, but got %d", len(tt.expected), len(result))
				return
			}

			for i, chunk := range result {
				if !chunk.Eq(tt.expected[i]) {
					t.Errorf("Chunk %d does not match expected result", i)
				}
			}
		})
	}
}

func TestHSliceReverse(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5}
	hsl = hsl.Reverse()

	if !reflect.DeepEqual(hsl, hg.HSlice[int]{5, 4, 3, 2, 1}) {
		t.Errorf("Expected %v, got %v", hg.HSlice[int]{5, 4, 3, 2, 1}, hsl)
	}
}

func TestHSliceAll(t *testing.T) {
	h1 := hg.NewHSlice[int]()
	h2 := hg.NewHSlice[int]().Append(1, 2, 3)
	h3 := hg.NewHSlice[int]().Append(2, 4, 6)

	testCases := []struct {
		f    func(int) bool
		name string
		h    hg.HSlice[int]
		want bool
	}{
		{
			name: "empty slice",
			f:    func(x int) bool { return x%2 == 0 },
			h:    h1,
			want: true,
		},
		{
			name: "all elements satisfy the condition",
			f:    func(x int) bool { return x%2 != 0 },
			h:    h2,
			want: false,
		},
		{
			name: "not all elements satisfy the condition",
			f:    func(x int) bool { return x%2 == 0 },
			h:    h3,
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.h.All(tc.f)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHSliceAny(t *testing.T) {
	h1 := hg.NewHSlice[int]()
	f1 := func(x int) bool { return x > 0 }

	if h1.Any(f1) {
		t.Errorf("Expected false for empty slice, got true")
	}

	h2 := hg.NewHSlice[int]().Append(1, 2, 3)
	f2 := func(x int) bool { return x < 1 }

	if h2.Any(f2) {
		t.Errorf("Expected false for slice with no matching elements, got true")
	}

	h3 := hg.NewHSlice[string]().Append("foo", "bar")
	f3 := func(x string) bool { return x == "bar" }

	if !h3.Any(f3) {
		t.Errorf("Expected true for slice with one matching element, got false")
	}

	h4 := hg.NewHSlice[int]().Append(1, 2, 3, 4, 5)
	f4 := func(x int) bool { return x%2 == 0 }

	if !h4.Any(f4) {
		t.Errorf("Expected true for slice with multiple matching elements, got false")
	}
}

func TestHSliceReduce(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5}
	sum := hsl.Reduce(func(index, value int) int { return index + value }, 0)

	if sum != 15 {
		t.Errorf("Expected %d, got %d", 15, sum)
	}
}

func TestHSliceFilter(t *testing.T) {
	var hsl hg.HSlice[int]

	hsl = hsl.Append(1, 2, 3, 4, 5)
	result := hsl.Filter(func(v int) bool { return v%2 == 0 })

	if result.Len() != 2 {
		t.Errorf("Expected 2, got %d", result.Len())
	}

	if result[0] != 2 {
		t.Errorf("Expected 2, got %d", result[0])
	}

	if result[1] != 4 {
		t.Errorf("Expected 4, got %d", result[1])
	}
}

func TestHSliceIndex(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5}

	if hsl.Index(1) != 0 {
		t.Error("Index of 1 should be 0")
	}

	if hsl.Index(2) != 1 {
		t.Error("Index of 2 should be 1")
	}

	if hsl.Index(3) != 2 {
		t.Error("Index of 3 should be 2")
	}

	if hsl.Index(4) != 3 {
		t.Error("Index of 4 should be 3")
	}

	if hsl.Index(5) != 4 {
		t.Error("Index of 5 should be 4")
	}

	if hsl.Index(6) != -1 {
		t.Error("Index of 6 should be -1")
	}
}

func TestHSliceRandomSample(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := hsl.RandomSample(5)

	if result.Len() != 5 {
		t.Errorf("Expected result length to be 5, got %d", result.Len())
	}

	for _, item := range result {
		if !hsl.Contains(item) {
			t.Errorf("Expected result to contain only items from the original slice, got %d", item)
		}
	}
}

func TestHSliceMap(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5}
	result := hsl.Map(func(i int) int { return i * 2 })

	if result.Len() != hsl.Len() {
		t.Errorf("Expected %d, got %d", hsl.Len(), result.Len())
	}

	for i := 0; i < result.Len(); i++ {
		if result[i] != hsl[i]*2 {
			t.Errorf("Expected %d, got %d", hsl[i]*2, result[i])
		}
	}
}

func TestHSliceAddUnique(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3}
	hsl = hsl.AddUnique(4, 5, 6)

	if !hsl.Contains(4) {
		t.Error("AddUnique failed")
	}

	hsl = hsl.AddUnique(4, 5, 6)
	if hsl.Len() != 6 {
		t.Error("AddUnique failed")
	}
}

func TestHSliceCount(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5, 6, 7}

	if hsl.Count(1) != 1 {
		t.Error("Expected 1, got ", hsl.Count(1))
	}

	if hsl.Count(2) != 1 {
		t.Error("Expected 1, got ", hsl.Count(2))
	}

	if hsl.Count(3) != 1 {
		t.Error("Expected 1, got ", hsl.Count(3))
	}

	if hsl.Count(4) != 1 {
		t.Error("Expected 1, got ", hsl.Count(4))
	}

	if hsl.Count(5) != 1 {
		t.Error("Expected 1, got ", hsl.Count(5))
	}

	if hsl.Count(6) != 1 {
		t.Error("Expected 1, got ", hsl.Count(6))
	}

	if hsl.Count(7) != 1 {
		t.Error("Expected 1, got ", hsl.Count(7))
	}
}

func TestHSliceSortBy(t *testing.T) {
	hsl1 := hg.NewHSlice[int]().Append(3, 1, 4, 1, 5)
	expected1 := hg.NewHSlice[int]().Append(1, 1, 3, 4, 5)

	actual1 := hsl1.SortBy(func(i, j int) bool { return hsl1[i] < hsl1[j] })

	if !actual1.Eq(expected1) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected1, actual1)
	}

	hsl2 := hg.NewHSlice[string]().Append("foo", "bar", "baz")
	expected2 := hg.NewHSlice[string]().Append("foo", "baz", "bar")

	actual2 := hsl2.SortBy(func(i, j int) bool { return hsl2[i] > hsl2[j] })

	if !actual2.Eq(expected2) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected2, actual2)
	}

	hsl3 := hg.NewHSlice[int]()
	expected3 := hg.NewHSlice[int]()

	actual3 := hsl3.SortBy(func(i, j int) bool { return hsl3[i] < hsl3[j] })

	if !actual3.Eq(expected3) {
		t.Errorf("SortBy failed: expected %v, but got %v", expected3, actual3)
	}
}

func TestHSliceJoin(t *testing.T) {
	hsl := hg.HSlice[string]{"1", "2", "3", "4", "5"}
	hstr := hsl.Join(",")

	if !strings.EqualFold("1,2,3,4,5", hstr.String()) {
		t.Errorf("Expected 1,2,3,4,5, got %s", hstr.String())
	}
}

func TestHSliceToStringSlice(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3}
	result := hsl.ToStringSlice()
	expected := []string{"1", "2", "3"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestHSliceAdd(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3}.Append(4, 5, 6)

	if !reflect.DeepEqual(hsl, hg.HSlice[int]{1, 2, 3, 4, 5, 6}) {
		t.Error("Add failed")
	}
}

func TestHSliceClone(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3}
	hslClone := hsl.Clone()

	if !hsl.Eq(hslClone) {
		t.Errorf("Clone() failed, expected %v, got %v", hsl, hslClone)
	}
}

func TestHSliceCut(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	hsl = hsl.Cut(0, 5)

	if !reflect.DeepEqual(hsl, hg.HSlice[int]{6, 7, 8, 9, 10}) {
		t.Errorf("Cut(0, 5) = %v, want %v", hsl, hg.HSlice[int]{6, 7, 8, 9, 10})
	}

	hsl = hsl.Cut(0, 5)
	if !reflect.DeepEqual(hsl, hg.HSlice[int]{}) {
		t.Errorf("Cut(0, 5) = %v, want %v", hsl, hg.HSlice[int]{})
	}
}

func TestHSliceLast(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5}
	if hsl.Last() != 5 {
		t.Error("Last() failed")
	}
}

func TestHSliceLastIndex(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5}
	if hsl.LastIndex() != 4 {
		t.Error("LastIndex() failed")
	}
}

func TestHSliceLen(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5}
	if hsl.Len() != 5 {
		t.Errorf("Expected 5, got %d", hsl.Len())
	}
}

func TestHSlicePop(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5}
	last, hsl := hsl.Pop()

	if last != 5 {
		t.Errorf("Expected 5, got %v", last)
	}

	if hsl.Len() != 4 {
		t.Errorf("Expected 4, got %v", hsl.Len())
	}
}

func TestHSliceRandom(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < 10; i++ {
		if hsl.Random() < 1 || hsl.Random() > 10 {
			t.Error("Random() should return a number between 1 and 10")
		}
	}
}

func TestHSliceMaxHInt(t *testing.T) {
	hsl := hg.HSlice[hg.HInt]{1, 2, 3, 4, 5}
	if max := hsl.Max(); max != 5 {
		t.Errorf("Max() = %d, want: %d.", max, 5)
	}
}

func TestHSliceMaxFloats(t *testing.T) {
	hsl := hg.HSlice[hg.HFloat]{2.2, 2.8, 2.1, 2.7}
	if max := hsl.Max(); max != 2.8 {
		t.Errorf("Max() = %f, want: %f.", max, 2.8)
	}
}

func TestHSliceMinHFloat(t *testing.T) {
	hsl := hg.HSlice[hg.HFloat]{2.2, 2.8, 2.1, 2.7}
	if min := hsl.Min(); min != 2.1 {
		t.Errorf("Min() = %f; want: %f", min, 2.1)
	}
}

func TestHSliceMinHInt(t *testing.T) {
	hsl := hg.HSlice[hg.HInt]{1, 2, 3, 4, 5}
	if min := hsl.Min(); min != 1 {
		t.Errorf("Min() = %d; want: %d", min, 1)
	}
}

func TestHSliceFilterZeroValues(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 0, 4, 0, 5, 0, 6, 0, 7, 0, 8, 0, 9, 0, 10}
	hsl = hsl.FilterZeroValues()

	if hsl.Len() != 10 {
		t.Errorf("Expected 10, got %d", hsl.Len())
	}

	for i := range iter.N(hsl.Len()) {
		if hsl[i] == 0 {
			t.Errorf("Expected non-zero value, got %d", hsl[i])
		}
	}
}

func TestHSliceRange(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	if hsl.Range(0, 0).Len() != 0 {
		t.Error("Expected 0, got", hsl.Range(0, 0).Len())
	}

	if hsl.Range(0, 1).Len() != 1 {
		t.Error("Expected 1, got", hsl.Range(0, 1).Len())
	}

	if hsl.Range(0, 2).Len() != 2 {
		t.Error("Expected 2, got", hsl.Range(0, 2).Len())
	}

	if hsl.Range(0, 3).Len() != 3 {
		t.Error("Expected 3, got", hsl.Range(0, 3).Len())
	}

	if hsl.Range(0, 4).Len() != 4 {
		t.Error("Expected 4, got", hsl.Range(0, 4).Len())
	}

	if !hsl.Range(0, -1).Eq(hg.HSlice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9}) {
		t.Error("Range(0, -1) failed")
	}

	if !hsl.Range(1, -2).Eq(hg.HSlice[int]{2, 3, 4, 5, 6, 7, 8}) {
		t.Error("Range(1, -2) failed")
	}

	if !hsl.Range(0, 10).Eq(hg.HSlice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Error("Range(0, 10) failed")
	}

	if !hsl.Range(0, 11).Eq(hg.HSlice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Error("Range(0, 11) failed")
	}

	if !hsl.Range(0, 100).Eq(hg.HSlice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Error("Range(0, 100) failed")
	}

	if !hsl.Range(0, 1).Eq(hg.HSlice[int]{1}) {
		t.Error("Range(0, 1) failed")
	}

	if !hsl.Range(-0, 1).Eq(hg.HSlice[int]{1}) {
		t.Error("Range(0, 1) failed")
	}

	if !hsl.Range(3, 6).Eq(hg.HSlice[int]{4, 5, 6}) {
		t.Error("Range(3, 6) failed")
	}

	if !hsl.Range(-1, -2).Eq(hg.HSlice[int]{}) {
		t.Error("Range(-1, -2) failed")
	}

	if !hsl.Range(0, -20).Eq(hg.HSlice[int]{}) {
		t.Error("Range(0, -20) failed")
	}

	if !hsl.Range(0, 0).Eq(hg.HSlice[int]{}) {
		t.Error("Range(0, 0) failed")
	}
}

func TestHSliceDelete(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5}
	hsl = hsl.Delete(2)

	if !reflect.DeepEqual(hsl, hg.HSlice[int]{1, 2, 4, 5}) {
		t.Errorf("Delete(2) = %v, want %v", hsl, hg.HSlice[int]{1, 2, 4, 5})
	}
}

func TestHSliceForEach(t *testing.T) {
	h1 := hg.NewHSlice[int]().Append(1, 2, 3, 4, 5)
	h2 := hg.NewHSlice[string]().Append("foo", "bar", "baz")
	h3 := hg.NewHSlice[float64]().Append(1.1, 2.2, 3.3, 4.4)

	var result1 []int

	h1.ForEach(func(i int) { result1 = append(result1, i) })

	if !reflect.DeepEqual(result1, []int{1, 2, 3, 4, 5}) {
		t.Errorf(
			"ForEach failed for %v, expected %v, but got %v",
			h1,
			[]int{1, 2, 3, 4, 5},
			result1,
		)
	}

	var result2 []string

	h2.ForEach(func(s string) { result2 = append(result2, s) })

	if !reflect.DeepEqual(result2, []string{"foo", "bar", "baz"}) {
		t.Errorf(
			"ForEach failed for %v, expected %v, but got %v",
			h2,
			[]string{"foo", "bar", "baz"},
			result2,
		)
	}

	var result3 []float64

	h3.ForEach(func(f float64) { result3 = append(result3, f) })

	if !reflect.DeepEqual(result3, []float64{1.1, 2.2, 3.3, 4.4}) {
		t.Errorf(
			"ForEach failed for %v, expected %v, but got %v",
			h3,
			[]float64{1.1, 2.2, 3.3, 4.4},
			result3,
		)
	}
}

func TestHSliceSFill(t *testing.T) {
	hsl := hg.HSlice[int]{1, 2, 3, 4, 5}
	hsl.Fill(0)

	for _, v := range hsl {
		if v != 0 {
			t.Errorf("Expected all elements to be 0, but found %d", v)
		}
	}
}

func TestHSliceSet(t *testing.T) {
	hsl := hg.NewHSlice[int](5)

	hsl.Set(0, 1)
	hsl.Set(0, 1)
	hsl.Set(2, 2)
	hsl.Set(4, 3)

	if !reflect.DeepEqual(hsl, hg.HSlice[int]{1, 0, 2, 0, 3}) {
		t.Errorf("Set() = %v, want %v", hsl, hg.HSlice[int]{1, 0, 2, 0, 3})
	}
}

func TestHSliceMapParallel(t *testing.T) {
	hsl := hg.NewHSlice[int](10).Fill(1)
	result := hsl.MapParallel(func(x int) int { return x * 2 })
	expected := hg.NewHSlice[int](10).Fill(2)

	if !result.Eq(expected) {
		t.Errorf("Unexpected result: got %v, expected %v", result, expected)
	}

	hsl = hg.NewHSlice[int](10000).Fill(1)
	result = hsl.MapParallel(func(x int) int { return x * 2 })
	expected = hg.NewHSlice[int](10000).Fill(2)

	if !result.Eq(expected) {
		t.Errorf("Unexpected result: got %v, expected %v", result, expected)
	}
}

func TestHSliceFilterParallel(t *testing.T) {
	hsl := hg.HSliceOf(1, 2, 3, 4, 5)
	expected := hg.HSliceOf(2, 4)
	actual := hsl.FilterParallel(func(x int) bool { return x%2 == 0 })

	if !actual.Eq(expected) {
		t.Errorf("FilterParallel failed. Expected %v, but got %v", expected, actual)
	}

	hsl = hg.HSliceOf(2, 4, 6, 8, 10)
	expected = hsl.Clone()
	actual = hsl.FilterParallel(func(x int) bool { return x%2 == 0 })

	if !actual.Eq(expected) {
		t.Errorf("FilterParallel failed. Expected %v, but got %v", expected, actual)
	}

	hsl = hg.HSliceOf(1, 3, 5, 7, 9)
	expected = hg.NewHSlice[int]()
	actual = hsl.FilterParallel(func(x int) bool { return x%2 == 0 })

	if !actual.Eq(expected) {
		t.Errorf("FilterParallel failed. Expected %v, but got %v", expected, actual)
	}
}

func TestHSliceReduceParallel(t *testing.T) {
	hsl := hg.NewHSlice[int](10).Fill(1)
	result := hsl.ReduceParallel(func(a, b int) int { return a + b }, 0)
	expected := hsl.Reduce(func(a, b int) int { return a + b }, 0)

	if result != expected {
		t.Errorf("Unexpected result: got %d, expected %d", result, expected)
	}

	hsl = hg.NewHSlice[int](10000).Fill(1)
	result = hsl.ReduceParallel(func(a, b int) int { return a + b }, 0)
	expected = hsl.Reduce(func(a, b int) int { return a + b }, 0)

	if result != expected {
		t.Errorf("Unexpected result: got %d, expected %d", result, expected)
	}
}

func TestHSliceZip(t *testing.T) {
	s1 := hg.HSliceOf(1, 2, 3, 4)
	s2 := hg.HSliceOf(5, 6, 7, 8)
	expected := []hg.HSlice[int]{{1, 5}, {2, 6}, {3, 7}, {4, 8}}
	result := s1.Zip(s2)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip(%v, %v) = %v, expected %v", s1, s2, result, expected)
	}

	s3 := hg.HSliceOf(1, 2, 3)
	s4 := hg.HSliceOf(4, 5)
	expected = []hg.HSlice[int]{{1, 4}, {2, 5}}
	result = s3.Zip(s4)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip(%v, %v) = %v, expected %v", s3, s4, result, expected)
	}

	s5 := hg.HSliceOf(1, 2, 3)
	s6 := hg.HSliceOf(4, 5, 6)
	s7 := hg.HSliceOf(7, 8, 9)
	expected = []hg.HSlice[int]{{1, 4, 7}, {2, 5, 8}, {3, 6, 9}}
	result = s5.Zip(s6, s7)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip(%v, %v, %v) = %v, expected %v", s5, s6, s7, result, expected)
	}

	s8 := hg.HSliceOf(1, 2, 3)
	s9 := hg.HSliceOf(4, 5)
	s10 := hg.HSliceOf(6)
	expected = []hg.HSlice[int]{{1, 4, 6}}
	result = s8.Zip(s9, s10)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zip(%v, %v, %v) = %v, expected %v", s8, s9, s10, result, expected)
	}
}

func TestHSliceFlatten(t *testing.T) {
	tests := []struct {
		name     string
		input    hg.HSlice[any]
		expected hg.HSlice[any]
	}{
		{
			name:     "Empty slice",
			input:    hg.HSlice[any]{},
			expected: hg.HSlice[any]{},
		},
		{
			name:     "Flat slice",
			input:    hg.HSlice[any]{1, "abc", 3.14},
			expected: hg.HSlice[any]{1, "abc", 3.14},
		},
		{
			name: "Nested slice",
			input: hg.HSlice[any]{
				1,
				hg.HSlice[int]{2, 3},
				"abc",
				hg.HSlice[string]{"def", "ghi"},
				hg.HSlice[float64]{4.5, 6.7},
			},
			expected: hg.HSlice[any]{1, 2, 3, "abc", "def", "ghi", 4.5, 6.7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Flatten()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Flatten() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHSliceCounter(t *testing.T) {
	h1 := hg.HSlice[int]{1, 2, 3, 2, 1, 4, 5, 4, 4}
	h2 := hg.HSlice[string]{"apple", "banana", "orange", "apple", "apple", "orange", "grape"}

	expected1 := hg.NewHMapOrd[any, int]()
	expected1.Set(3, 1)
	expected1.Set(5, 1)
	expected1.Set(1, 2)
	expected1.Set(2, 2)
	expected1.Set(4, 3)

	result1 := h1.Counter()
	if !result1.Eq(expected1) {
		t.Errorf("Counter() returned %v, expected %v", result1, expected1)
	}

	// Test with string values
	expected2 := hg.NewHMapOrd[any, int]()
	expected2.Set("banana", 1)
	expected2.Set("grape", 1)
	expected2.Set("orange", 2)
	expected2.Set("apple", 3)

	result2 := h2.Counter()
	if !result2.Eq(expected2) {
		t.Errorf("Counter() returned %v, expected %v", result2, expected2)
	}
}

func TestHSliceReplace(t *testing.T) {
	tests := []struct {
		name     string
		input    hg.HSlice[string]
		i, j     int
		values   []string
		expected hg.HSlice[string]
	}{
		{
			name:     "basic test",
			input:    hg.HSlice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        3,
			values:   []string{"e", "f"},
			expected: hg.HSlice[string]{"a", "e", "f", "d"},
		},
		{
			name:     "replace at start",
			input:    hg.HSlice[string]{"a", "b", "c", "d"},
			i:        0,
			j:        2,
			values:   []string{"e", "f"},
			expected: hg.HSlice[string]{"e", "f", "c", "d"},
		},
		{
			name:     "replace at end",
			input:    hg.HSlice[string]{"a", "b", "c", "d"},
			i:        2,
			j:        4,
			values:   []string{"e", "f"},
			expected: hg.HSlice[string]{"a", "b", "e", "f"},
		},
		{
			name:     "replace with more values",
			input:    hg.HSlice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        2,
			values:   []string{"e", "f", "g", "h"},
			expected: hg.HSlice[string]{"a", "e", "f", "g", "h", "c", "d"},
		},
		{
			name:     "replace with fewer values",
			input:    hg.HSlice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        3,
			values:   []string{"e"},
			expected: hg.HSlice[string]{"a", "e", "d"},
		},
		{
			name:     "replace entire slice",
			input:    hg.HSlice[string]{"a", "b", "c", "d"},
			i:        0,
			j:        4,
			values:   []string{"e", "f", "g", "h"},
			expected: hg.HSlice[string]{"e", "f", "g", "h"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Replace(tt.i, tt.j, tt.values...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHSliceReplaceInPlace(t *testing.T) {
	tests := []struct {
		name     string
		input    hg.HSlice[string]
		i, j     int
		values   []string
		expected hg.HSlice[string]
	}{
		{
			name:     "basic test",
			input:    hg.HSlice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        3,
			values:   []string{"e", "f"},
			expected: hg.HSlice[string]{"a", "e", "f", "d"},
		},
		{
			name:     "replace at start",
			input:    hg.HSlice[string]{"a", "b", "c", "d"},
			i:        0,
			j:        2,
			values:   []string{"e", "f"},
			expected: hg.HSlice[string]{"e", "f", "c", "d"},
		},
		{
			name:     "replace at end",
			input:    hg.HSlice[string]{"a", "b", "c", "d"},
			i:        2,
			j:        4,
			values:   []string{"e", "f"},
			expected: hg.HSlice[string]{"a", "b", "e", "f"},
		},
		{
			name:     "replace with more values",
			input:    hg.HSlice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        2,
			values:   []string{"e", "f", "g", "h"},
			expected: hg.HSlice[string]{"a", "e", "f", "g", "h", "c", "d"},
		},
		{
			name:     "replace with fewer values",
			input:    hg.HSlice[string]{"a", "b", "c", "d"},
			i:        1,
			j:        3,
			values:   []string{"e"},
			expected: hg.HSlice[string]{"a", "e", "d"},
		},
		{
			name:     "replace entire slice",
			input:    hg.HSlice[string]{"a", "b", "c", "d"},
			i:        0,
			j:        4,
			values:   []string{"e", "f", "g", "h"},
			expected: hg.HSlice[string]{"e", "f", "g", "h"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hsl := &tt.input
			hsl.ReplaceInPlace(tt.i, tt.j, tt.values...)
			if !reflect.DeepEqual(*hsl, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, *hsl)
			}
		})
	}
}

func TestHSliceContainsAny(t *testing.T) {
	testCases := []struct {
		hsl    hg.HSlice[int]
		other  hg.HSlice[int]
		expect bool
	}{
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{6, 7, 8, 9, 10}, false},
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{5, 6, 7, 8, 9}, true},
		{hg.HSlice[int]{}, hg.HSlice[int]{1, 2, 3, 4, 5}, false},
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{}, false},
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{1, 2, 3, 4, 5}, true},
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{6}, false},
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{1}, true},
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{6, 7, 8, 9, 0, 3}, true},
	}

	for _, tc := range testCases {
		if result := tc.hsl.ContainsAny(tc.other); result != tc.expect {
			t.Errorf("ContainsAny(%v, %v) = %v; want %v", tc.hsl, tc.other, result, tc.expect)
		}
	}
}

func TestHSliceContainsAll(t *testing.T) {
	testCases := []struct {
		hsl    hg.HSlice[int]
		other  hg.HSlice[int]
		expect bool
	}{
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{1, 2, 3}, true},
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{1, 2, 3, 6}, false},
		{hg.HSlice[int]{}, hg.HSlice[int]{1, 2, 3, 4, 5}, false},
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{}, true},
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{1, 2, 3, 4, 5}, true},
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{6}, false},
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{1}, true},
		{hg.HSlice[int]{1, 2, 3, 4, 5}, hg.HSlice[int]{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 5, 5}, true},
	}

	for _, tc := range testCases {
		if result := tc.hsl.ContainsAll(tc.other); result != tc.expect {
			t.Errorf("ContainsAll(%v, %v) = %v; want %v", tc.hsl, tc.other, result, tc.expect)
		}
	}
}

func TestHSliceUnique(t *testing.T) {
	testCases := []struct {
		input  hg.HSlice[int]
		output hg.HSlice[int]
	}{
		{
			input:  hg.NewHSlice[int]().Append(1, 2, 3, 4, 5),
			output: hg.NewHSlice[int]().Append(1, 2, 3, 4, 5),
		},
		{
			input:  hg.NewHSlice[int]().Append(1, 2, 3, 4, 5, 5, 4, 3, 2, 1),
			output: hg.NewHSlice[int]().Append(1, 2, 3, 4, 5),
		},
		{
			input:  hg.NewHSlice[int]().Append(1, 1, 1, 1, 1),
			output: hg.NewHSlice[int]().Append(1),
		},
		{
			input:  hg.NewHSlice[int](),
			output: hg.NewHSlice[int](),
		},
	}

	for _, tc := range testCases {
		actual := tc.input.Unique()
		if !reflect.DeepEqual(actual, tc.output) {
			t.Errorf("Unique(%v) returned %v, expected %v", tc.input, actual, tc.output)
		}
	}
}

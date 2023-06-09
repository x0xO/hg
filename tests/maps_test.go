package hg_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/x0xO/hg"
)

func TestHMapKeys(t *testing.T) {
	hmap := hg.NewHMap[string, int]()
	hmap.Set("a", 1)
	hmap.Set("b", 2)
	hmap.Set("c", 3)

	keys := hmap.Keys()
	if keys.Len() != 3 {
		t.Errorf("Expected 3 keys, got %d", keys.Len())
	}

	if !keys.Contains("a") {
		t.Errorf("Expected key 'a'")
	}

	if !keys.Contains("b") {
		t.Errorf("Expected key 'b'")
	}

	if !keys.Contains("c") {
		t.Errorf("Expected key 'c'")
	}
}

func TestHMapValues(t *testing.T) {
	hmap := hg.NewHMap[string, int]()

	hmap.Set("a", 1)
	hmap.Set("b", 2)
	hmap.Set("c", 3)

	values := hmap.Values()

	if values.Len() != 3 {
		t.Errorf("Expected 3 values, got %d", values.Len())
	}

	if !values.Contains(1) {
		t.Errorf("Expected value '1'")
	}

	if !values.Contains(2) {
		t.Errorf("Expected value '2'")
	}

	if !values.Contains(3) {
		t.Errorf("Expected value '3'")
	}
}

func TestHMapClone(t *testing.T) {
	hmap := hg.NewHMap[string, int]()
	hmap["a"] = 1
	hmap["b"] = 2
	hmap["c"] = 3

	nhmap := hmap.Clone()

	if hmap.Len() != nhmap.Len() {
		t.Errorf("Clone failed: expected %d, got %d", hmap.Len(), nhmap.Len())
	}

	for k, v := range hmap {
		if nhmap[k] != v {
			t.Errorf("Clone failed: expected %d, got %d", v, nhmap[k])
		}
	}
}

func TestHMapCopy(t *testing.T) {
	src := hg.HMap[string, int]{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	dst := hg.HMap[string, int]{
		"d": 4,
		"e": 5,
		"a": 6,
	}

	dst.Copy(src)

	if dst.Len() != 5 {
		t.Errorf("Expected len(dst) to be 5, got %d", len(dst))
	}

	if dst["a"] != 1 {
		t.Errorf("Expected dst[\"a\"] to be 1, got %d", dst["a"])
	}

	if dst["b"] != 2 {
		t.Errorf("Expected dst[\"b\"] to be 2, got %d", dst["b"])
	}

	if dst["c"] != 3 {
		t.Errorf("Expected dst[\"c\"] to be 3, got %d", dst["c"])
	}
}

func TestHMapAdd(t *testing.T) {
	hmap := hg.HMap[string, string]{}
	hmap = hmap.Set("key", "value")

	if hmap["key"] != "value" {
		t.Error("Expected value to be 'value'")
	}
}

func TestHMapDelete(t *testing.T) {
	hmap := hg.HMap[string, int]{"a": 1, "b": 2, "c": 3}

	hmap = hmap.Delete("a", "b")

	if hmap.Len() != 1 {
		t.Errorf("Expected length of 1, got %d", hmap.Len())
	}

	if _, ok := hmap["a"]; ok {
		t.Errorf("Expected key 'a' to be deleted")
	}

	if _, ok := hmap["b"]; ok {
		t.Errorf("Expected key 'b' to be deleted")
	}

	if _, ok := hmap["c"]; !ok {
		t.Errorf("Expected key 'c' to be present")
	}
}

func TestHMapEqual(t *testing.T) {
	hmap := hg.NewHMap[string, string]()
	hmap.Set("key", "value")

	other := hg.NewHMap[string, string]()
	other = other.Set("key", "value")

	if !hmap.Eq(other) {
		t.Error("hmap and other should be equal")
	}

	other = other.Set("key", "other value")

	if hmap.Eq(other) {
		t.Error("hmap and other should not be equal")
	}
}

func TestHMapToMap(t *testing.T) {
	hmap := hg.NewHMap[string, int]()
	hmap.Set("a", 1)
	hmap.Set("b", 2)
	hmap.Set("c", 3)

	nmap := hmap.ToMap()

	if len(nmap) != 3 {
		t.Errorf("Expected 3, got %d", len(nmap))
	}

	if nmap["a"] != 1 {
		t.Errorf("Expected 1, got %d", nmap["a"])
	}

	if nmap["b"] != 2 {
		t.Errorf("Expected 2, got %d", nmap["b"])
	}

	if nmap["c"] != 3 {
		t.Errorf("Expected 3, got %d", nmap["c"])
	}
}

func TestHMapLen(t *testing.T) {
	hmap := hg.HMap[int, int]{}
	if hmap.Len() != 0 {
		t.Errorf("Expected 0, got %d", hmap.Len())
	}

	hmap[1] = 1
	if hmap.Len() != 1 {
		t.Errorf("Expected 1, got %d", hmap.Len())
	}

	hmap[2] = 2
	if hmap.Len() != 2 {
		t.Errorf("Expected 2, got %d", hmap.Len())
	}
}

func TestHMapMap(t *testing.T) {
	hm := hg.NewHMap[int, string](3)
	hm.Set(1, "one")
	hm.Set(2, "two")
	hm.Set(3, "three")

	expected := hg.NewHMap[int, string](3)
	expected.Set(2, "one")
	expected.Set(4, "two")
	expected.Set(6, "three")

	mapped := hm.Map(func(k int, v string) (int, string) { return k * 2, v })

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("Map failed: expected %v, but got %v", expected, mapped)
	}

	expected = hg.NewHMap[int, string](3)
	expected.Set(1, "one_suffix")
	expected.Set(2, "two_suffix")
	expected.Set(3, "three_suffix")

	mapped = hm.Map(func(k int, v string) (int, string) { return k, v + "_suffix" })

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("Map failed: expected %v, but got %v", expected, mapped)
	}

	expected = hg.NewHMap[int, string](3)
	expected.Set(0, "")
	expected.Set(1, "one")
	expected.Set(3, "three")

	mapped = hm.Map(func(k int, v string) (int, string) {
		if k == 2 {
			return 0, ""
		}
		return k, v
	})

	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("Map failed: expected %v, but got %v", expected, mapped)
	}
}

func TestHMapFilter(t *testing.T) {
	hm := hg.NewHMap[string, int](3)
	hm.Set("one", 1)
	hm.Set("two", 2)
	hm.Set("three", 3)

	expected := hg.NewHMap[string, int](1)
	expected.Set("two", 2)

	filtered := hm.Filter(func(k string, v int) bool { return v%2 == 0 })

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Filter failed: expected %v, but got %v", expected, filtered)
	}

	expected = hg.NewHMap[string, int](2)
	expected.Set("one", 1)
	expected.Set("three", 3)

	filtered = hm.Filter(func(k string, v int) bool { return strings.Contains(k, "e") })

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Filter failed: expected %v, but got %v", expected, filtered)
	}

	expected = hg.NewHMap[string, int](3)
	expected.Set("one", 1)
	expected.Set("two", 2)
	expected.Set("three", 3)

	filtered = hm.Filter(func(k string, v int) bool { return true })

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Filter failed: expected %v, but got %v", expected, filtered)
	}

	expected = hg.NewHMap[string, int](0)

	filtered = hm.Filter(func(k string, v int) bool { return false })

	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Filter failed: expected %v, but got %v", expected, filtered)
	}
}

func TestHMapMapParallel(t *testing.T) {
	hmap1 := hg.NewHMap[string, int](0)
	fn1 := func(k string, v int) (string, int) { return k, v * 2 }
	expected1 := hg.NewHMap[string, int](0)
	actual1 := hmap1.MapParallel(fn1)

	if !actual1.Eq(expected1) {
		t.Errorf("Test case 1 failed: expected %v but got %v", expected1, actual1)
	}

	hmap2 := hg.NewHMap[string, int](4)
	hmap2.Set("a", 1)
	hmap2.Set("b", 2)
	hmap2.Set("c", 3)
	hmap2.Set("d", 4)

	fn2 := func(k string, v int) (string, int) { return k, v * 2 }

	expected2 := hg.NewHMap[string, int](4)
	expected2.Set("a", 2)
	expected2.Set("b", 4)
	expected2.Set("c", 6)
	expected2.Set("d", 8)

	actual2 := hmap2.MapParallel(fn2)
	if !actual2.Eq(expected2) {
		t.Errorf("Test case 2 failed: expected %v but got %v", expected2, actual2)
	}
}

func TestHMapOf(t *testing.T) {
	testCases := []struct {
		name          string
		entries       []any
		expectedHMap  hg.HMap[string, int]
		expectedError string
	}{
		{
			name:         "empty",
			entries:      []any{},
			expectedHMap: hg.HMap[string, int]{},
		},
		{
			name:         "single key-value pair",
			entries:      []any{"one", 1},
			expectedHMap: hg.HMap[string, int]{"one": 1},
		},
		{
			name:         "multiple key-value pairs",
			entries:      []any{"one", 1, "two", 2, "three", 3},
			expectedHMap: hg.HMap[string, int]{"one": 1, "two": 2, "three": 3},
		},
		{
			name:         "duplicate keys",
			entries:      []any{"one", 1, "two", 2, "one", 3},
			expectedHMap: hg.HMap[string, int]{"one": 3, "two": 2},
		},
		{
			name:          "odd number of arguments",
			entries:       []any{"one", 1, "two"},
			expectedError: "HMapOf requires an even number of arguments representing alternating keys and values",
		},
		{
			name:          "incorrect types",
			entries:       []any{"one", 1, 2, "two"},
			expectedError: "HMapOf requires alternating keys and values of the correct types",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if testCase.expectedError == "" {
						t.Errorf("Unexpected panic: %v", r)
					} else if r != testCase.expectedError {
						t.Errorf("Expected panic: %v, got: %v", testCase.expectedError, r)
					}
				} else if testCase.expectedError != "" {
					t.Errorf("Expected panic: %v, but no panic occurred", testCase.expectedError)
				}
			}()

			hmap := hg.HMapOf[string, int](testCase.entries...)
			if !reflect.DeepEqual(hmap, testCase.expectedHMap) {
				t.Errorf("Expected HMap: %v, got: %v", testCase.expectedHMap, hmap)
			}
		})
	}
}

func TestHMapInvertValues(t *testing.T) {
	hmap := hg.NewHMap[int, string](0)
	inv := hmap.Invert()

	if inv.Len() != 0 {
		t.Errorf("Expected inverted map to have length 0, but got length %d", inv.Len())
	}

	hmap2 := hg.NewHMap[string, int](3)
	hmap2.Set("one", 1)
	hmap2.Set("two", 2)
	hmap2.Set("three", 3)

	inv2 := hmap2.Invert()

	if inv2.Len() != 3 {
		t.Errorf("Expected inverted map to have length 3, but got length %d", inv2.Len())
	}

	if inv2.Get(1) != "one" {
		t.Errorf("Expected inverted map to map 1 to 'one', but got %s", inv2.Get(1))
	}

	if inv2.Get(2) != "two" {
		t.Errorf("Expected inverted map to map 2 to 'two', but got %s", inv2.Get(2))
	}

	if inv2.Get(3) != "three" {
		t.Errorf("Expected inverted map to map 3 to 'three', but got %s", inv2.Get(3))
	}
}

package hg

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/x0xO/hg/pkg/deref"
	"github.com/x0xO/hg/pkg/iter"
	"github.com/x0xO/hg/pkg/rand"
)

// NewHSlice creates a new HSlice of the given generic type T with the specified length and
// capacity.
// The size variadic parameter can have zero, one, or two integer values.
// If no values are provided, an empty HSlice with a length and capacity of 0 is returned.
// If one value is provided, it sets both the length and capacity of the HSlice.
// If two values are provided, the first value sets the length and the second value sets the
// capacity.
//
// Parameters:
//
// - size ...int: A variadic parameter specifying the length and/or capacity of the HSlice
//
// Returns:
//
// - HSlice[T]: A new HSlice of the specified generic type T with the given length and capacity
//
// Example usage:
//
//	s1 := hg.NewHSlice[int]()        // Creates an empty HSlice of type int
//	s2 := hg.NewHSlice[int](5)       // Creates an HSlice with length and capacity of 5
//	s3 := hg.NewHSlice[int](3, 10)   // Creates an HSlice with length of 3 and capacity of 10
func NewHSlice[T any](size ...int) HSlice[T] {
	length, capacity := 0, 0

	switch {
	case len(size) > 1:
		length, capacity = size[0], size[1]
	case len(size) == 1:
		length, capacity = size[0], size[0]
	}

	return make(HSlice[T], length, capacity)
}

// HSliceOf creates a new generic slice containing the provided elements.
func HSliceOf[T any](slice ...T) HSlice[T] { return slice }

// Counter returns an ordered map with the counts of each unique element in the slice.
// This function is useful when you want to count the occurrences of each unique element in an
// HSlice.
//
// Returns:
//
// - *hMapOrd[any, int]: An ordered HMap with keys representing the unique elements in the HSlice
// and values representing the counts of those elements.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 1, 2, 1}
//	counts := slice.Counter()
//	// The counts ordered HMap will contain:
//	// 1 -> 3 (since 1 appears three times)
//	// 2 -> 2 (since 2 appears two times)
//	// 3 -> 1 (since 3 appears once)
func (hsl HSlice[T]) Counter() *HMapOrd[any, int] {
	result := NewHMapOrd[any, int](hsl.Len())

	hsl.ForEach(func(t T) { result.Set(t, result.GetOrDefault(t, 0)+1) })
	result.SortBy(
		func(i, j int) bool { return deref.Of(result)[i].Value < deref.Of(result)[j].Value },
	)

	return result
}

// Enumerate returns a map with the index of each element as the key.
// This function is useful when you want to create an HMap where the keys are the indices of the
// elements in an HSlice, and the values are the corresponding elements.
//
// Returns:
//
// - HMap[int, T]: An HMap with keys representing the indices of the elements in the HSlice and
// values representing the corresponding elements.
//
// Example usage:
//
//	slice := hg.HSlice[int]{10, 20, 30}
//	indexedMap := slice.Enumerate()
//	// The indexedMap HMap will contain:
//	// 0 -> 10 (since 10 is at index 0)
//	// 1 -> 20 (since 20 is at index 1)
//	// 2 -> 30 (since 30 is at index 2)
func (hsl HSlice[T]) Enumerate() HMap[int, T] {
	result := NewHMap[int, T](hsl.Len())
	for k, v := range hsl {
		result.Set(k, v)
	}

	return result
}

// Fill fills the slice with the specified value.
// This function is useful when you want to create an HSlice with all elements having the same
// value.
// This method can be used in place, as it modifies the original slice.
//
// Parameters:
//
// - val T: The value to fill the HSlice with.
//
// Returns:
//
// - HSlice[T]: A reference to the original HSlice filled with the specified value.
//
// Example usage:
//
//	slice := hg.HSlice[int]{0, 0, 0}
//	slice.Fill(5)
//
// The modified slice will now contain: 5, 5, 5.
func (hsl HSlice[T]) Fill(val T) HSlice[T] {
	for i := range iter.N(hsl.Len()) {
		hsl.Set(i, val)
	}

	return hsl
}

// ToHMapHashed returns a map with the hashed version of each element as the key.
func (hsl HSlice[T]) ToHMapHashed() HMap[HString, T] {
	result := NewHMap[HString, T](hsl.Len())

	hsl.ForEach(func(t T) {
		switch val := any(t).(type) {
		case HInt:
			result.Set(val.Hash().MD5(), t)
		case int:
			result.Set(HInt(val).Hash().MD5(), t)
		case HString:
			result.Set(val.Hash().MD5(), t)
		case string:
			result.Set(HString(val).Hash().MD5(), t)
		case HBytes:
			result.Set(val.Hash().MD5().HString(), t)
		case []byte:
			result.Set(HBytes(val).Hash().MD5().HString(), t)
		case HFloat:
			result.Set(val.Hash().MD5(), t)
		case float64:
			result.Set(HFloat(val).Hash().MD5(), t)
		}
	})

	return result
}

// Chunks splits the HSlice into smaller chunks of the specified size.
// The function iterates through the HSlice, creating new HSlice[T] chunks of the specified size.
// If size is less than or equal to 0 or the HSlice is empty,
// it returns an empty slice of HSlice[T].
// If size is greater than or equal to the length of the HSlice,
// it returns a slice of HSlice[T] containing the original HSlice.
//
// Parameters:
//
// - size int: The size of each chunk.
//
// Returns:
//
// - []HSlice[T]: A slice of HSlice[T] containing the chunks of the original HSlice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 4, 5, 6}
//	chunks := slice.Chunks(2)
//
// The resulting chunks will be: [{1, 2}, {3, 4}, {5, 6}].
func (hsl HSlice[T]) Chunks(size int) []HSlice[T] {
	if size <= 0 || hsl.Empty() {
		return []HSlice[T]{}
	}

	chunks := (hsl.Len() + size - 1) / size // ceil(len/size)
	result := make([]HSlice[T], 0, chunks)

	for i := 0; i < hsl.Len(); i += size {
		end := i + size
		if end > hsl.Len() {
			end = hsl.Len()
		}

		result = append(result, hsl.Range(i, end))
	}

	return result
}

// All returns true if all elements in the slice satisfy the provided condition.
// This function is useful when you want to check if all elements in an HSlice meet a certain
// criteria.
//
// Parameters:
//
// - fn func(T) bool: A function that returns a boolean indicating whether the element satisfies
// the condition.
//
// Returns:
//
// - bool: True if all elements in the HSlice satisfy the condition, false otherwise.
//
// Example usage:
//
//	slice := hg.HSlice[int]{2, 4, 6, 8, 10}
//	isEven := func(num int) bool { return num%2 == 0 }
//	allEven := slice.All(isEven)
//
// The resulting allEven will be true since all elements in the slice are even.
func (hsl HSlice[T]) All(fn func(T) bool) bool {
	for _, val := range hsl {
		if !fn(val) {
			return false
		}
	}

	return true
}

// Any returns true if any element in the slice satisfies the provided condition.
// This function is useful when you want to check if at least one element in an HSlice meets a
// certain criteria.
//
// Parameters:
//
// - fn func(T) bool: A function that returns a boolean indicating whether the element satisfies
// the condition.
//
// Returns:
//
// - bool: True if at least one element in the HSlice satisfies the condition, false otherwise.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 3, 5, 7, 9}
//	isEven := func(num int) bool { return num%2 == 0 }
//	anyEven := slice.Any(isEven)
//
// The resulting anyEven will be false since none of the elements in the slice are even.
func (hsl HSlice[T]) Any(fn func(T) bool) bool {
	for _, val := range hsl {
		if fn(val) {
			return true
		}
	}

	return false
}

// Index returns the index of the first occurrence of the specified value in the slice, or -1 if
// not found.
func (hsl HSlice[T]) Index(val T) int {
	for i, v := range hsl {
		if reflect.DeepEqual(v, val) {
			return i
		}
	}

	return -1
}

// RandomSample returns a new slice containing a random sample of elements from the original slice.
// The sampling is done without replacement, meaning that each element can only appear once in the
// result.
//
// Parameters:
//
// - sequence int: The number of elements to include in the random sample.
//
// Returns:
//
// - HSlice[T]: A new HSlice containing the random sample of elements.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9}
//	sample := slice.RandomSample(3)
//
// The resulting sample will contain 3 unique elements randomly selected from the original
// slice.
func (hsl HSlice[T]) RandomSample(sequence int) HSlice[T] {
	if sequence >= hsl.Len() {
		return hsl
	}

	return hsl.Clone().Shuffle()[:sequence]
}

// Insert inserts values at the specified index in the slice and returns the resulting slice.
// The original slice remains unchanged.
//
// Parameters:
//
// - i int: The index at which to insert the new values.
//
// - values ...T: A variadic list of values to insert at the specified index.
//
// Returns:
//
// - HSlice[T]: A new HSlice containing the original elements and the inserted values.
//
// Example usage:
//
//	slice := hg.HSlice[string]{"a", "b", "c", "d"}
//	newSlice := slice.Insert(2, "e", "f")
//
// The resulting newSlice will be: ["a", "b", "e", "f", "c", "d"].
func (hsl HSlice[T]) Insert(i int, values ...T) HSlice[T] { return hsl.Replace(i, i, values...) }

// InsertInPlace inserts values at the specified index in the slice and modifies the original
// slice.
//
// Parameters:
//
// - i int: The index at which to insert the new values.
//
// - values ...T: A variadic list of values to insert at the specified index.
//
// Example usage:
//
//	slice := hg.HSlice[string]{"a", "b", "c", "d"}
//	slice.InsertInPlace(2, "e", "f")
//
// The resulting slice will be: ["a", "b", "e", "f", "c", "d"].
func (hsl *HSlice[T]) InsertInPlace(i int, values ...T) { hsl.ReplaceInPlace(i, i, values...) }

// Replace replaces the elements of hsl[i:j] with the given values, and returns
// a new slice with the modifications. The original slice remains unchanged.
// Replace panics if hsl[i:j] is not a valid slice of hsl.
//
// Parameters:
//
// - i int: The starting index of the slice to be replaced.
//
// - j int: The ending index of the slice to be replaced.
//
// - values ...T: A variadic list of values to replace the existing slice.
//
// Returns:
//
// - HSlice[T]: A new HSlice containing the original elements with the specified elements replaced.
//
// Example usage:
//
//	slice := hg.HSlice[string]{"a", "b", "c", "d"}
//	newSlice := slice.Replace(1, 3, "e", "f")
//
// The original slice remains ["a", "b", "c", "d"], and the newSlice will be: ["a", "e", "f", "d"].
func (hsl HSlice[T]) Replace(i, j int, values ...T) HSlice[T] {
	_ = hsl[i:j] // verify that i:j is a valid subslice

	total := hsl[:i].Len() + len(values) + hsl[j:].Len()
	slice := make(HSlice[T], total)

	copy(slice, hsl[:i])
	copy(slice[i:], values)
	copy(slice[i+len(values):], hsl[j:])

	return slice
}

// ReplaceInPlace replaces the elements of hsl[i:j] with the given values,
// and modifies the original slice in place. ReplaceInPlace panics if hsl[i:j]
// is not a valid slice of hsl.
//
// Parameters:
//
// - i int: The starting index of the slice to be replaced.
//
// - j int: The ending index of the slice to be replaced.
//
// - values ...T: A variadic list of values to replace the existing slice.
//
// Example usage:
//
//	slice := hg.HSlice[string]{"a", "b", "c", "d"}
//	slice.ReplaceInPlace(1, 3, "e", "f")
//
// After the ReplaceInPlace operation, the resulting slice will be: ["a", "e", "f", "d"].
func (hsl *HSlice[T]) ReplaceInPlace(i, j int, values ...T) {
	_ = deref.Of(hsl)[i:j] // verify that i:j is a valid subslice

	diff := len(values) - (j - i)
	if diff > 0 {
		*hsl = deref.Of(hsl).Append(NewHSlice[T](diff)...)
	}

	copy(deref.Of(hsl)[i+len(values):], deref.Of(hsl)[j:])
	copy(deref.Of(hsl)[i:], values)

	if diff < 0 {
		*hsl = deref.Of(hsl)[:hsl.Len()+diff]
	}
}

// Unique returns a new slice containing unique elements from the current slice.
//
// The order of elements in the returned slice is the same as the order in the original slice.
//
// Returns:
//
// - HSlice[T]: A new HSlice containing unique elements from the current slice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 2, 4, 5, 3}
//	unique := slice.Unique()
//
// The resulting unique slice will be: [1, 2, 3, 4, 5].
func (hsl HSlice[T]) Unique() HSlice[T] {
	result := NewHSlice[T](0, hsl.Len())

	hsl.ForEach(func(t T) {
		if !result.Contains(t) {
			result = result.Append(t)
		}
	})

	return result.Clip()
}

// ForEach applies a given function to each element in the slice.
//
// The function takes one parameter of type T (the same type as the elements of the slice).
// The function is applied to each element in the order they appear in the slice.
//
// Parameters:
//
// - fn (func(T)): The function to be applied to each element of the slice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3}
//	slice.ForEach(func(val int) {
//	    fmt.Println(val * 2)
//	})
//	// Output:
//	// 2
//	// 4
//	// 6
func (hsl HSlice[T]) ForEach(fn func(T)) {
	for _, val := range hsl {
		fn(val)
	}
}

// Map returns a new slice by applying a given function to each element in the current slice.
//
// The function takes one parameter of type T (the same type as the elements of the slice)
// and returns a value of type T. The returned value is added to a new slice,
// which is then returned as the result.
//
// Parameters:
//
// - fn (func(T) T): The function to be applied to each element of the slice.
//
// Returns:
//
// - HSlice[T]: A new slice containing the results of applying the function to each element
// of the current slice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3}
//	doubled := slice.Map(func(val int) int {
//	    return val * 2
//	})
//	fmt.Println(doubled)
//
// Output: [2 4 6].
func (hsl HSlice[T]) Map(fn func(T) T) HSlice[T] {
	result := NewHSlice[T](0, hsl.Len())
	hsl.ForEach(func(t T) { result = result.Append(fn(t)) })

	return result
}

// MapInPlace applies a given function to each element in the current slice,
// modifying the elements in place.
//
// The function takes one parameter of type T (the same type as the elements of the slice)
// and returns a value of type T. The returned value replaces the original element in the slice.
//
// Parameters:
//
// - fn (func(T) T): The function to be applied to each element of the slice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3}
//	slice.MapInPlace(func(val int) int {
//	    return val * 2
//	})
//	fmt.Println(slice)
//
// Output: [2 4 6].
func (hsl *HSlice[T]) MapInPlace(fn func(T) T) {
	for i := range iter.N(hsl.Len()) {
		hsl.Set(i, fn(hsl.Get(i)))
	}
}

// Filter returns a new slice containing elements that satisfy a given condition.
//
// The function takes one parameter of type T (the same type as the elements of the slice)
// and returns a boolean value. If the returned value is true, the element is added
// to a new slice, which is then returned as the result.
//
// Parameters:
//
// - fn (func(T) bool): The function to be applied to each element of the slice
// to determine if it should be included in the result.
//
// Returns:
//
// - HSlice[T]: A new slice containing the elements that satisfy the given condition.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 4, 5}
//	even := slice.Filter(func(val int) bool {
//	    return val%2 == 0
//	})
//	fmt.Println(even)
//
// Output: [2 4].
func (hsl HSlice[T]) Filter(fn func(T) bool) HSlice[T] {
	result := NewHSlice[T](0, hsl.Len())

	hsl.ForEach(func(t T) {
		if fn(t) {
			result = result.Append(t)
		}
	})

	return result.Clip()
}

// FilterInPlace removes elements from the current slice that do not satisfy a given condition.
//
// The function takes one parameter of type T (the same type as the elements of the slice)
// and returns a boolean value. If the returned value is false, the element is removed
// from the slice.
//
// Parameters:
//
// - fn (func(T) bool): The function to be applied to each element of the slice
// to determine if it should be kept in the slice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 4, 5}
//	slice.FilterInPlace(func(val int) bool {
//	    return val%2 == 0
//	})
//	fmt.Println(slice)
//
// Output: [2 4].
func (hsl *HSlice[T]) FilterInPlace(fn func(T) bool) {
	j := 0

	for i := range iter.N(hsl.Len()) {
		if fn(hsl.Get(i)) {
			hsl.Set(j, hsl.Get(i))
			j++
		}
	}

	*hsl = deref.Of(hsl)[:j]
}

// Reduce reduces the slice to a single value using a given function and an initial value.
//
// The function takes two parameters of type T (the same type as the elements of the slice):
// an accumulator and a value from the slice. The accumulator is initialized with the provided
// initial value, and the function is called for each element in the slice. The returned value
// from the function becomes the new accumulator value for the next iteration. After processing
// all the elements in the slice, the final accumulator value is returned as the result.
//
// Parameters:
//
// - fn (func(acc, val T) T): The function to be applied to each element of the slice
// and the accumulator. This function should return a new value for the accumulator.
//
// - initial (T): The initial value for the accumulator.
//
// Returns:
//
// - T: The final accumulator value after processing all the elements in the slice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 4, 5}
//	sum := slice.Reduce(func(acc, val int) int {
//	    return acc + val
//	}, 0)
//	fmt.Println(sum)
//
// Output: 15.
func (hsl HSlice[T]) Reduce(fn func(acc, val T) T, initial T) T {
	acc := initial

	hsl.ForEach(func(t T) { acc = fn(acc, t) })

	return acc
}

// MapParallel applies a given function to each element in the slice in parallel and returns a new
// slice.
//
// The function iterates over the elements of the slice and applies the provided function
// to each element. If the length of the slice is less than a predefined threshold (max),
// it falls back to the sequential Map function. Otherwise, the slice is divided into two
// halves and the function is applied to each half in parallel using goroutines. The
// resulting slices are then combined to form the final output slice.
//
// Note: The order of the elements in the output slice may not be the same as the input
// slice due to parallel processing.
//
// Parameters:
//
// - fn (func(T) T): The function to be applied to each element of the slice.
//
// Returns:
//
// - HSlice[T]: A new slice with the results of applying the given function to each element
// of the original slice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 4, 5}
//	squared := slice.MapParallel(func(val int) int {
//	    return val * val
//	})
//	fmt.Println(squared)
//
// Output: {1 4 9 16 25}.
func (hsl HSlice[T]) MapParallel(fn func(T) T) HSlice[T] {
	const max = 1 << 11
	if hsl.Len() < max {
		return hsl.Map(fn)
	}

	half := hsl.Len() / 2
	left := hsl.Range(0, half)
	right := hsl.Range(half, hsl.Len())

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		left = left.MapParallel(fn)

		wg.Done()
	}()

	right = right.MapParallel(fn)

	wg.Wait()

	return NewHSlice[T](0, hsl.Len()).Append(left...).Append(right...)
}

// FilterParallel returns a new slice containing elements that satisfy a given condition, computed
// in parallel.
//
// The function iterates over the elements of the slice and applies the provided predicate
// function to each element. If the length of the slice is less than a predefined threshold (max),
// it falls back to the sequential Filter function. Otherwise, the slice is divided into two
// halves and the predicate function is applied to each half in parallel using goroutines. The
// resulting slices are then combined to form the final output slice.
//
// Note: The order of the elements in the output slice may not be the same as the input
// slice due to parallel processing.
//
// Parameters:
//
// - fn (func(T) bool): The predicate function to be applied to each element of the slice.
//
// Returns:
//
// - HSlice[T]: A new slice containing the elements that satisfy the given condition.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 4, 5}
//	even := slice.FilterParallel(func(val int) bool {
//	    return val % 2 == 0
//	})
//	fmt.Println(even)
//
// Output: {2 4}.
func (hsl HSlice[T]) FilterParallel(fn func(T) bool) HSlice[T] {
	const max = 1 << 11
	if hsl.Len() < max {
		return hsl.Filter(fn)
	}

	half := hsl.Len() / 2
	left := hsl.Range(0, half)
	right := hsl.Range(half, hsl.Len())

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		left = left.FilterParallel(fn)

		wg.Done()
	}()

	right = right.FilterParallel(fn)

	wg.Wait()

	return NewHSlice[T](0, left.Len()+right.Len()).Append(left...).Append(right...)
}

// ReduceParallel reduces the slice to a single value using a given function and an initial value,
// computed in parallel.
//
// The function iterates over the elements of the slice and applies the provided reducer function
// to each element in a pairwise manner. If the length of the slice is less than a predefined
// threshold (max),
// it falls back to the sequential Reduce function. Otherwise, the slice is divided into two
// halves and the reducer function is applied to each half in parallel using goroutines. The
// resulting values are combined using the reducer function to produce the final output value.
//
// Note: Due to parallel processing, the order in which the reducer function is applied to the
// elements may not be the same as the input slice.
//
// Parameters:
//
// - fn (func(T, T) T): The reducer function to be applied to each element of the slice.
//
// - initial (T): The initial value to be used as the starting point for the reduction.
//
// Returns:
//
// - T: A single value obtained by applying the reducer function to the elements of the slice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 4, 5}
//	sum := slice.ReduceParallel(func(acc, val int) int {
//	    return acc + val
//	}, 0)
//	fmt.Println(sum)
//
// Output: 15.
func (hsl HSlice[T]) ReduceParallel(fn func(T, T) T, initial T) T {
	const max = 1 << 11
	if hsl.Len() < max {
		return hsl.Reduce(fn, initial)
	}

	half := hsl.Len() / 2
	left := hsl.Range(0, half)
	right := hsl.Range(half, hsl.Len())

	result := NewHSlice[T](2)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		result.Set(0, left.ReduceParallel(fn, initial))

		wg.Done()
	}()

	result.Set(1, right.ReduceParallel(fn, initial))

	wg.Wait()

	return result.Reduce(fn, initial)
}

// AddUnique appends unique elements from the provided arguments to the current slice.
//
// The function iterates over the provided elements and checks if they are already present
// in the slice. If an element is not already present, it is appended to the slice. The
// resulting slice is returned, containing the unique elements from both the original
// slice and the provided elements.
//
// Parameters:
//
// - elems (...T): A variadic list of elements to be appended to the slice.
//
// Returns:
//
// - HSlice[T]: A new slice containing the unique elements from both the original slice
// and the provided elements.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 4, 5}
//	slice = slice.AddUnique(3, 4, 5, 6, 7)
//	fmt.Println(slice)
//
// Output: [1 2 3 4 5 6 7].

func (hsl HSlice[T]) AddUnique(elems ...T) HSlice[T] {
	for _, elem := range elems {
		if !hsl.Contains(elem) {
			hsl = hsl.Append(elem)
		}
	}

	return hsl
}

// AddUniqueInPlace appends unique elements from the provided arguments to the current slice.
//
// The function iterates over the provided elements and checks if they are already present
// in the slice. If an element is not already present, it is appended to the slice.
//
// Parameters:
//
// - elems (...T): A variadic list of elements to be appended to the slice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 4, 5}
//	slice.AddUniqueInPlace(3, 4, 5, 6, 7)
//	fmt.Println(slice)
//
// Output: [1 2 3 4 5 6 7].
func (hsl *HSlice[T]) AddUniqueInPlace(elems ...T) {
	for _, elem := range elems {
		if !hsl.Contains(elem) {
			*hsl = hsl.Append(elem)
		}
	}
}

// Get returns the element at the given index, handling negative indices as counting from the end
// of the slice.
func (hsl HSlice[T]) Get(index int) T {
	if HInt(index).IsNegative() {
		index = hsl.Len() + index
	}

	if index > hsl.LastIndex() {
		index = hsl.LastIndex()
	}

	return hsl[index]
}

// Count returns the count of the given element in the slice.
func (hsl HSlice[T]) Count(elem T) int {
	if hsl.Empty() {
		return 0
	}

	var counter int

	hsl.ForEach(func(t T) {
		if reflect.DeepEqual(t, elem) {
			counter++
		}
	})

	return counter
}

// Max returns the maximum element in the slice, assuming elements are comparable.
func (hsl HSlice[T]) Max() T {
	if hsl.Empty() {
		return *new(T)
	}

	max := hsl.Get(0)

	var greater func(a, b any) bool

	switch any(max).(type) {
	case HInt:
		greater = func(a, b any) bool { return a.(HInt).Gt(b.(HInt)) }
	case int:
		greater = func(a, b any) bool { return a.(int) > b.(int) }
	case HString:
		greater = func(a, b any) bool { return a.(HString).Gt(b.(HString)) }
	case string:
		greater = func(a, b any) bool { return a.(string) > b.(string) }
	case HFloat:
		greater = func(a, b any) bool { return a.(HFloat).Gt(b.(HFloat)) }
	case float64:
		greater = func(a, b any) bool { return HFloat(a.(float64)).Gt(HFloat(b.(float64))) }
	}

	hsl.ForEach(func(t T) {
		if greater(t, max) {
			max = t
		}
	})

	return max
}

// Min returns the minimum element in the slice, assuming elements are comparable.
func (hsl HSlice[T]) Min() T {
	if hsl.Empty() {
		return *new(T)
	}

	min := hsl.Get(0)

	var less func(a, b any) bool

	switch any(min).(type) {
	case HInt:
		less = func(a, b any) bool { return a.(HInt).Lt(b.(HInt)) }
	case int:
		less = func(a, b any) bool { return a.(int) < b.(int) }
	case HString:
		less = func(a, b any) bool { return a.(HString).Lt(b.(HString)) }
	case string:
		less = func(a, b any) bool { return a.(string) < b.(string) }
	case HFloat:
		less = func(a, b any) bool { return a.(HFloat).Lt(b.(HFloat)) }
	case float64:
		less = func(a, b any) bool { return HFloat(a.(float64)).Lt(HFloat(b.(float64))) }
	}

	hsl.ForEach(func(t T) {
		if less(t, min) {
			min = t
		}
	})

	return min
}

// Shuffle shuffles the elements in the slice randomly. This method can be used in place, as it
// modifies the original slice.
//
// The function uses the crypto/rand package to generate random indices.
//
// Returns:
//
// - HSlice[T]: The modified slice with the elements shuffled randomly.
//
// Example usage:
//
// slice := hg.HSlice[int]{1, 2, 3, 4, 5}
// shuffled := slice.Shuffle()
// fmt.Println(shuffled)
//
// Output: A randomly shuffled version of the original slice, e.g., [4 1 5 2 3].
func (hsl HSlice[T]) Shuffle() HSlice[T] {
	n := hsl.Len()

	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		hsl.Swap(i, j)
	}

	return hsl
}

// Reverse reverses the order of the elements in the slice. This method can be used in place, as it
// modifies the original slice.
//
// Returns:
//
// - HSlice[T]: The modified slice with the elements reversed.
//
// Example usage:
//
// slice := hg.HSlice[int]{1, 2, 3, 4, 5}
// reversed := slice.Reverse()
// fmt.Println(reversed)
//
// Output: [5 4 3 2 1].
func (hsl HSlice[T]) Reverse() HSlice[T] {
	for i, j := 0, hsl.Len()-1; i < j; i, j = i+1, j-1 {
		hsl.Swap(i, j)
	}

	return hsl
}

// SortBy sorts the elements in the slice using the provided comparison function. This method can
// be used in place, as it modifies the original slice.
//
// The function takes a custom comparison function as an argument and sorts the elements
// of the slice using the provided logic. The comparison function should return true if
// the element at index i should come before the element at index j, and false otherwise.
//
// Parameters:
//
// - f func(i, j int) bool: A comparison function that takes two indices i and j.
//
// Returns:
//
// - HSlice[T]: The sorted HSlice.
//
// Example usage:
//
// hsl := NewHSlice[int](1, 5, 3, 2, 4)
// hsl.SortBy(func(i, j int) bool { return hsl[i] < hsl[j] }) // sorts in ascending order.
func (hsl HSlice[T]) SortBy(f func(i, j int) bool) HSlice[T] {
	sort.Slice(hsl, f)
	return hsl
}

// FilterZeroValues returns a new slice with all zero values removed.
//
// The function iterates over the elements in the slice and checks if they are
// zero values using the reflect.DeepEqual function. If an element is not a zero value,
// it is added to the resulting slice. The new slice, containing only non-zero values,
// is returned.
//
// Returns:
//
// - HSlice[T]: A new slice containing only non-zero values from the original slice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 0, 4, 0}
//	nonZeroSlice := slice.FilterZeroValues()
//	fmt.Println(nonZeroSlice)
//
// Output: [1 2 4].
func (hsl HSlice[T]) FilterZeroValues() HSlice[T] {
	return hsl.Filter(func(v T) bool { return !reflect.DeepEqual(v, *new(T)) })
}

// FilterZeroValuesInPlace removes all zero values from the current slice.
//
// The function iterates over the elements in the slice and checks if they are
// zero values using the reflect.DeepEqual function. If an element is a zero value,
// it is removed from the slice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 0, 4, 0}
//	slice.FilterZeroValuesInPlace()
//	fmt.Println(slice)
//
// Output: [1 2 4].
func (hsl *HSlice[T]) FilterZeroValuesInPlace() {
	hsl.FilterInPlace(func(v T) bool { return !reflect.DeepEqual(v, *new(T)) })
}

// ToStringSlice converts the slice into a slice of strings.
func (hsl HSlice[T]) ToStringSlice() []string {
	result := NewHSlice[string](0, hsl.Len())
	hsl.ForEach(func(t T) { result = result.Append(fmt.Sprint(t)) })

	return result
}

// Join joins the elements in the slice into a single HString, separated by the provided separator
// (if any).
func (hsl HSlice[T]) Join(sep ...T) HString {
	var separator string
	if len(sep) != 0 {
		separator = fmt.Sprint(sep[0])
	}

	return HString(strings.Join(hsl.ToStringSlice(), separator))
}

// Range returns a new slice containing elements from the current slice between the specified start
// and end indices.
//
// The function checks if the start and end indices are within the bounds of the original slice.
// If the end index is negative, it is added to the length of the slice to calculate the actual end
// index.
// If the start index is negative or greater than the end index, an empty slice is returned.
// If the end index is greater than the length of the slice, it is set to the length of the slice.
//
// Parameters:
//
// - start (int): The start index of the range.
//
// - end (int): The end index of the range.
//
// Returns:
//
// - HSlice[T]: A new slice containing elements from the current slice between the start and end
// indices.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 4, 5}
//	subSlice := slice.Range(1, 4)
//	fmt.Println(subSlice)
//
// Output: [2 3 4].
func (hsl HSlice[T]) Range(start, end int) HSlice[T] {
	if HInt(end).IsNegative() {
		end = hsl.Len() + end
	}

	if start > end || HInt(start).IsNegative() {
		return HSlice[T]{}
	}

	if end > hsl.Len() {
		end = hsl.Len()
	}

	return hsl.Clone()[start:end]
}

// Cut returns a new slice that is the current slice with the elements between the specified start
// and end indices removed.
//
// The function checks if the start and end indices are within the bounds of the original slice.
// If the end index is negative, it is added to the length of the slice to calculate the actual end
// index.
// If the start index is negative or greater than the end index, an empty slice is returned.
// If the end index is greater than the length of the slice, it is set to the length of the slice.
//
// Parameters:
//
// - start (int): The start index of the range to be removed.
//
// - end (int): The end index of the range to be removed.
//
// Returns:
//
// - HSlice[T]: A new slice containing elements from the current slice with the specified range
// removed.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 4, 5}
//	newSlice := slice.Cut(1, 4)
//	fmt.Println(newSlice)
//
// Output: [1 5].
func (hsl HSlice[T]) Cut(start, end int) HSlice[T] {
	if HInt(end).IsNegative() {
		end = hsl.Len() + end
	}

	if start > end || HInt(start).IsNegative() {
		return HSlice[T]{}
	}

	if end > hsl.Len() {
		end = hsl.Len()
	}

	return hsl.Range(0, start).Append(hsl.Range(end, hsl.Len())...)
}

// Random returns a random element from the slice.
//
// The function uses the crypto/rand package to generate a random index within the bounds of the
// slice. If the slice is empty, the zero value of type T is returned.
//
// Returns:
//
// - T: A random element from the slice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3, 4, 5}
//	randomElement := slice.Random()
//	fmt.Println(randomElement)
//
// Output: <any random element from the slice>.
func (hsl HSlice[T]) Random() T {
	if hsl.Empty() {
		return *new(T)
	}

	return hsl.Get(rand.Intn(hsl.Len()))
}

// Clone returns a copy of the slice.
func (hsl HSlice[T]) Clone() HSlice[T] {
	slice := NewHSlice[T](hsl.Len())
	copy(slice, hsl)

	return slice
}

// LastIndex returns the last index of the slice.
func (hsl HSlice[T]) LastIndex() int {
	if !hsl.Empty() {
		return hsl.Len() - 1
	}

	return 0
}

// factorial a utility function that calculates the factorial of a given number.
func factorial(n int) int {
	if n <= 1 {
		return 1
	}

	return n * factorial(n-1)
}

// Permutations returns all possible permutations of the elements in the slice.
//
// The function uses a recursive approach to generate all the permutations of the elements.
// If the slice has a length of 0 or 1, it returns the slice itself wrapped in a single-element
// slice.
//
// Returns:
//
// - []HSlice[T]: A slice of HSlice[T] containing all possible permutations of the elements in the
// slice.
//
// Example usage:
//
//	slice := hg.HSlice[int]{1, 2, 3}
//	perms := slice.Permutations()
//	for _, perm := range perms {
//	    fmt.Println(perm)
//	}
//	// Output:
//	// [1 2 3]
//	// [1 3 2]
//	// [2 1 3]
//	// [2 3 1]
//	// [3 1 2]
//	// [3 2 1]
func (hsl HSlice[T]) Permutations() []HSlice[T] {
	if hsl.Len() <= 1 {
		return []HSlice[T]{hsl}
	}

	perms := make([]HSlice[T], 0, factorial(hsl.Len()))

	for i, elem := range hsl {
		rest := NewHSlice[T](hsl.Len() - 1)

		copy(rest[:i], hsl[:i])
		copy(rest[i:], hsl[i+1:])

		subPerms := rest.Permutations()

		for j := range subPerms {
			subPerms[j] = append(HSlice[T]{elem}, subPerms[j]...)
		}

		perms = append(perms, subPerms...)
	}

	return perms
}

// Zip zips the elements of the given slices with the current slice into a new slice of HSlice[T]
// elements.
//
// The function combines the elements of the current slice with the elements of the given slices by
// index. The length of the resulting slice of HSlice[T] elements is determined by the shortest
// input slice.
//
// Params:
//
// - slices: The slices to be zipped with the current slice.
//
// Returns:
//
// - []HSlice[T]: A new slice of HSlice[T] elements containing the zipped elements of the input
// slices.
//
// Example usage:
//
//	slice1 := hg.HSlice[int]{1, 2, 3}
//	slice2 := hg.HSlice[int]{4, 5, 6}
//	slice3 := hg.HSlice[int]{7, 8, 9}
//	zipped := slice1.Zip(slice2, slice3)
//	for _, group := range zipped {
//	    fmt.Println(group)
//	}
//	// Output:
//	// [1 4 7]
//	// [2 5 8]
//	// [3 6 9]
func (hsl HSlice[T]) Zip(slices ...HSlice[T]) []HSlice[T] {
	minLen := hsl.Len()

	for _, slice := range slices {
		if slice.Len() < minLen {
			minLen = slice.Len()
		}
	}

	result := make([]HSlice[T], 0, minLen)

	for i := range iter.N(minLen) {
		values := NewHSlice[T](0, len(slices)+1).Append(hsl.Get(i))
		for _, j := range slices {
			values = values.Append(j.Get(i))
		}

		result = append(result, values)
	}

	return result
}

// Flatten flattens the nested slice structure into a single-level HSlice[any].
//
// It recursively traverses the nested slice structure and appends all non-slice elements to a new
// HSlice[any].
//
// Returns:
//
// - HSlice[any]: A new HSlice[any] containing the flattened elements.
//
// Example usage:
//
//	nested := hg.HSlice[any]{1, 2, hg.HSlice[int]{3, 4, 5}, []any{6, 7, []int{8, 9}}}
//	flattened := nested.Flatten()
//	fmt.Println(flattened)
//
// Output: [1 2 3 4 5 6 7 8 9].
func (hsl HSlice[T]) Flatten() HSlice[any] {
	flattened := NewHSlice[any]()
	flattenRecursive(reflect.ValueOf(hsl), &flattened)

	return flattened
}

// flattenRecursive a helper function for recursively flattening nested slices.
func flattenRecursive(val reflect.Value, flattened *HSlice[any]) {
	for i := range iter.N(val.Len()) {
		elem := val.Index(i)
		if elem.Kind() == reflect.Interface {
			elem = elem.Elem()
		}

		if elem.Kind() == reflect.Slice {
			flattenRecursive(elem, flattened)
		} else {
			*flattened = append(*flattened, elem.Interface())
		}
	}
}

// Eq returns true if the slice is equal to the provided other slice.
func (hsl HSlice[T]) Eq(other HSlice[T]) bool {
	if hsl.Len() != other.Len() {
		return false
	}

	for index, val := range hsl {
		if !reflect.DeepEqual(val, other.Get(index)) {
			return false
		}
	}

	return true
}

// String returns a string representation of the slice.
func (hsl HSlice[T]) String() string {
	var builder strings.Builder

	hsl.ForEach(func(v T) { builder.WriteString(fmt.Sprintf("%v, ", v)) })

	return HString(builder.String()).AddPrefix("HSlice[").TrimRight(", ").Add("]").String()
}

// Append appends the provided elements to the slice and returns the modified slice.
func (hsl HSlice[T]) Append(elems ...T) HSlice[T] { return append(hsl, elems...) }

// AppendInPlace appends the provided elements to the slice and modifies the original slice.
func (hsl *HSlice[T]) AppendInPlace(elems ...T) { *hsl = hsl.Append(elems...) }

// Cap returns the capacity of the HSlice.
func (hsl HSlice[T]) Cap() int { return cap(hsl) }

// Contains returns true if the slice contains the provided value.
func (hsl HSlice[T]) Contains(val T) bool { return hsl.Index(val) >= 0 }

// ContainsAny checks if the HSlice contains any element from another HSlice.
func (hsl HSlice[T]) ContainsAny(other HSlice[T]) bool {
	for _, val := range other {
		if hsl.Contains(val) {
			return true
		}
	}

	return false
}

// ContainsAll checks if the HSlice contains all elements from another HSlice.
func (hsl HSlice[T]) ContainsAll(other HSlice[T]) bool {
	for _, val := range other {
		if !hsl.Contains(val) {
			return false
		}
	}

	return true
}

// Delete removes the element at the specified index from the slice and returns the modified slice.
func (hsl HSlice[T]) Delete(i int) HSlice[T] { return hsl.Cut(i, i+1) }

// DeleteInPlace removes the element at the specified index from the slice and modifies the
// original slice.
func (hsl *HSlice[T]) DeleteInPlace(i int) {
	if i < 0 || i >= hsl.Len() {
		return
	}

	copy(deref.Of(hsl)[i:], deref.Of(hsl)[i+1:])
	*hsl = deref.Of(hsl)[:hsl.Len()-1]
}

// Empty returns true if the slice is empty.
func (hsl HSlice[T]) Empty() bool { return hsl.Len() == 0 }

// Last returns the last element of the slice.
func (hsl HSlice[T]) Last() T { return hsl.Get(-1) }

// Len returns the length of the slice.
func (hsl HSlice[T]) Len() int { return len(hsl) }

// Ne returns true if the slice is not equal to the provided other slice.
func (hsl HSlice[T]) Ne(other HSlice[T]) bool { return !hsl.Eq(other) }

// NotEmpty returns true if the slice is not empty.
func (hsl HSlice[T]) NotEmpty() bool { return hsl.Len() != 0 }

// Pop returns the last element of the slice and a new slice without the last element.
func (hsl HSlice[T]) Pop() (T, HSlice[T]) { return hsl.Last(), hsl.Range(0, -1) }

// Set sets the value at the specified index in the slice and returns the modified slice.
// This method can be used in place, as it modifies the original slice.
//
// Parameters:
//
// - i (int): The index at which to set the new value.
//
// - val (T): The new value to be set at the specified index.
//
// Returns:
//
// - HSlice[T]: The modified slice with the new value set at the specified index.
//
// Example usage:
//
// slice := hg.HSlice[int]{1, 2, 3, 4, 5}
// slice.Set(2, 99)
// fmt.Println(slice)
//
// Output: [1 2 99 4 5].
func (hsl HSlice[T]) Set(i int, val T) HSlice[T] { hsl[i] = val; return hsl }

// Swap swaps the elements at the specified indices in the slice and returns the modified slice.
// This method can be used in place, as it modifies the original slice.
//
// Parameters:
//
// - i (int): The index of the first element to be swapped.
//
// - j (int): The index of the second element to be swapped.
//
// Returns:
//
// - HSlice[T]: The modified slice with the elements at the specified indices swapped.
//
// Example usage:
//
// slice := hg.HSlice[int]{1, 2, 3, 4, 5}
// slice.Swap(1, 3)
// fmt.Println(slice)
//
// Output: [1 4 3 2 5].
func (hsl HSlice[T]) Swap(i, j int) HSlice[T] { hsl[i], hsl[j] = hsl[j], hsl[i]; return hsl }

// Clip removes unused capacity from the slice.
func (hsl HSlice[T]) Clip() HSlice[T] { return hsl[:hsl.Len():hsl.Len()] }

// ToSlice returns a new slice with the same elements as the HSlice[T].
func (hsl HSlice[T]) ToSlice() []T { return hsl }

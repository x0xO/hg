package hg

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// NewHMap creates a new HMap of the specified size or an empty HMap if no size is provided.
func NewHMap[K comparable, V any](size ...int) HMap[K, V] {
	if len(size) == 0 {
		return make(HMap[K, V], 0)
	}

	return make(HMap[K, V], size[0])
}

// HMapFromMap creates an HMap from a given Go map.
func HMapFromMap[K comparable, V any](m map[K]V) HMap[K, V] {
	hmap := NewHMap[K, V](len(m))
	for k, v := range m {
		hmap.Set(k, v)
	}

	return hmap
}

// HMapOf creates an HMap from a list of alternating keys and values.
// The function takes a variadic parameter representing a list of alternating keys and values.
// The keys must be of a comparable type, while the values can be of any type.
// The function returns a newly created HMap containing the provided key-value pairs.
//
// Parameters:
//
// - entries ...any: A variadic parameter representing a list of alternating keys and values.
// The number of elements in this list must be even, as it should contain pairs of keys and values.
//
// Returns:
//
// - HMap[K, V]: A new HMap containing the provided key-value pairs.
//
// Panics:
//
// - If the number of elements in 'entries' is not even, as it must contain pairs of keys and
// values. If the provided keys and values are not of the correct types (K and V, respectively).
//
// Example usage:
//
//	hmap := hg.HMapOf["string", int]("key1", 1, "key2", 2, "key3", 3)
func HMapOf[K comparable, V any](entries ...any) HMap[K, V] {
	if len(entries)%2 != 0 {
		panic(
			"HMapOf requires an even number of arguments representing alternating keys and values",
		)
	}

	hmap := NewHMap[K, V](len(entries) / 2)

	for i := 0; i < len(entries); i += 2 {
		key, keyOk := entries[i].(K)
		value, valueOk := entries[i+1].(V)

		if !keyOk || !valueOk {
			panic("HMapOf requires alternating keys and values of the correct types")
		}

		hmap.Set(key, value)
	}

	return hmap
}

// Invert inverts the keys and values of the HMap, returning a new HMap with values as keys and
// keys as values. Note that the inverted HMap will have 'any' as the key type, since not all value
// types are guaranteed to be comparable.
func (hmap HMap[K, V]) Invert() HMap[any, K] {
	result := NewHMap[any, K](hmap.Len())
	hmap.ForEach(func(k K, v V) { result.Set(v, k) })

	return result
}

// Keys returns a slice of the HMap's keys.
func (hmap HMap[K, V]) Keys() HSlice[K] {
	keys := NewHSlice[K](0, hmap.Len())
	hmap.ForEach(func(k K, _ V) { keys = keys.Append(k) })

	return keys
}

// Values returns a slice of the HMap's values.
func (hmap HMap[K, V]) Values() HSlice[V] {
	values := NewHSlice[V](0, hmap.Len())
	hmap.ForEach(func(_ K, v V) { values = values.Append(v) })

	return values
}

// Contains checks if the HMap contains the specified key.
func (hmap HMap[K, V]) Contains(key K) bool {
	_, ok := hmap[key]
	return ok
}

// Clone creates a new HMap that is a copy of the original HMap.
func (hmap HMap[K, V]) Clone() HMap[K, V] {
	result := NewHMap[K, V](hmap.Len())
	hmap.ForEach(func(k K, v V) { result.Set(k, v) })

	return result
}

// Copy copies the source HMap's key-value pairs to the target HMap.
func (hmap HMap[K, V]) Copy(src HMap[K, V]) HMap[K, V] {
	src.ForEach(func(k K, v V) { hmap.Set(k, v) })
	return hmap
}

// Delete removes the specified keys from the HMap.
func (hmap HMap[K, V]) Delete(keys ...K) HMap[K, V] {
	for _, key := range keys {
		delete(hmap, key)
	}

	return hmap
}

// ToMap converts the HMap to a regular Go map.
func (hmap HMap[K, V]) ToMap() map[K]V { return hmap }

// Map applies a function to each key-value pair in the HMap and returns a new HMap with the
// results. The provided function 'fn' should take a key and a value as input parameters and return
// a new key-value pair.
//
// Parameters:
//
// - fn func(K, V) (K, V): A function that takes a key and a value as input parameters and returns
// a new key-value pair.
//
// Returns:
//
// - HMap[K, V]: A new HMap containing the key-value pairs resulting from applying the provided
// function to each key-value pair in the original HMap.
//
// Example usage:
//
//	mappedHMap := originalHMap.Map(func(key K, value V) (K, V) {
//		return key, value * 2
//	})
func (hmap HMap[K, V]) Map(fn func(K, V) (K, V)) HMap[K, V] {
	result := NewHMap[K, V](hmap.Len())
	hmap.ForEach(func(k K, v V) { result.Set(fn(k, v)) })

	return result
}

// Filter filters the HMap based on a given function and returns a new HMap containing the matching
// key-value pairs. The provided function 'fn' should take a key and a value as input parameters
// and return a boolean value.
// If the function returns true, the key-value pair will be included in the resulting HMap.
//
// Parameters:
//
// - fn func(K, V) bool: A function that takes a key and a value as input parameters and returns a
// boolean value.
//
// Returns:
//
// - HMap[K, V]: A new HMap containing the key-value pairs for which the provided function returned
// true.
//
// Example usage:
//
//	filteredHMap := originalHMap.Filter(func(key K, value V) bool {
//		return value >= 10
//	})
func (hmap HMap[K, V]) Filter(fn func(K, V) bool) HMap[K, V] {
	result := NewHMap[K, V]()

	hmap.ForEach(func(k K, v V) {
		if fn(k, v) {
			result.Set(k, v)
		}
	})

	return result
}

// ForEach applies a function to each key-value pair in the HMap.
// The provided function 'fn' should take a key and a value as input parameters and perform an
// operation.
// This function is useful for side effects, as it does not return a new HMap.
//
// Parameters:
//
// - fn func(K, V): A function that takes a key and a value as input parameters and performs an
// operation.
//
// Example usage:
//
//	originalHMap.ForEach(func(key K, value V) {
//		fmt.Printf("Key: %v, Value: %v\n", key, value)
//	})
func (hmap HMap[K, V]) ForEach(fn func(K, V)) {
	for key, val := range hmap {
		fn(key, val)
	}
}

// MapParallel applies a function to each key-value pair in the HMap in parallel and returns a new
// HMap with the results.
// The provided function 'fn' should take a key and a value as input parameters and return a new
// key-value pair.
// This function is designed for better performance on large HMaps by utilizing
// parallel processing.
//
// Parameters:
//
// - fn func(K, V) (K, V): A function that takes a key and a value as input parameters and returns
// a new key-value pair.
//
// Returns:
//
// - HMap[K, V]: A new HMap containing the key-value pairs resulting from applying the provided
// function to each key-value pair in the original HMap.
//
// Example usage:
//
//	mappedHMap := originalHMap.MapParallel(func(key K, value V) (K, V) {
//		return key, value * 2
//	})
func (hmap HMap[K, V]) MapParallel(fn func(K, V) (K, V)) HMap[K, V] {
	const max = 1 << 11
	if hmap.Len() < max {
		return hmap.Map(fn)
	}

	half := hmap.Len() / 2

	left := NewHMap[K, V](half)
	right := NewHMap[K, V](half - left.Len())

	i := 0
	for k, v := range hmap {
		if i < half {
			left.Set(k, v)
		} else {
			right.Set(k, v)
		}
		i++
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		left = left.MapParallel(fn)

		wg.Done()
	}()

	right = right.MapParallel(fn)

	wg.Wait()

	return hmap.Clear().Copy(left).Copy(right)
}

// FilterParallel filters the HMap based on a given function in parallel and returns a new HMap
// containing the matching key-value pairs. The provided function 'fn' should take a key and a
// value as input parameters and return a boolean value.
// If the function returns true, the key-value pair will be included in the resulting HMap.
// This function is designed for better performance on large HMaps by utilizing parallel
// processing.
//
// Parameters:
//
// - fn func(K, V) bool: A function that takes a key and a value as input parameters and returns a
// boolean value.
//
// Returns:
//
// - HMap[K, V]: A new HMap containing the key-value pairs for which the provided function returned
// true.
//
// Example usage:
//
//	filteredHMap := originalHMap.FilterParallel(func(key K, value V) bool {
//		return value >= 10
//	})
//
// TODO: написать тесты.
func (hmap HMap[K, V]) FilterParallel(fn func(K, V) bool) HMap[K, V] {
	const max = 1 << 11
	if hmap.Len() < max {
		return hmap.Filter(fn)
	}

	half := hmap.Len() / 2

	left := NewHMap[K, V](half)
	right := NewHMap[K, V](half - left.Len())

	i := 0
	for k, v := range hmap {
		if i < half {
			left.Set(k, v)
		} else {
			right.Set(k, v)
		}
		i++
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		left = left.FilterParallel(fn)

		wg.Done()
	}()

	right = right.FilterParallel(fn)

	wg.Wait()

	return NewHMap[K, V](left.Len() + right.Len()).Copy(left).Copy(right)
}

// Eq checks if two HMaps are equal.
func (hmap HMap[K, V]) Eq(other HMap[K, V]) bool {
	if hmap.Len() != other.Len() {
		return false
	}

	for key, value := range hmap {
		if !other.Contains(key) || !reflect.DeepEqual(other[key], value) {
			return false
		}
	}

	return true
}

// String returns a string representation of the HMap.
func (hmap HMap[K, V]) String() string {
	var builder strings.Builder

	hmap.ForEach(func(k K, v V) { builder.WriteString(fmt.Sprintf("%v:%v, ", k, v)) })

	return HString(builder.String()).TrimRight(", ").Format("HMap{%s}").String()
}

// GetOrDefault returns the value for a key. If the key does not exist, returns the default value
// instead. This function is useful when you want to provide a fallback value for keys that may not
// be present in the HMap.
//
// Parameters:
//
// - key K: The key for which to retrieve the value.
//
// - defaultValue V: The default value to return if the key does not exist in the HMap.
//
// Returns:
//
// - V: The value associated with the key if it exists in the HMap, or the default value if the key
// is not found.
//
// Example usage:
//
//	value := hmap.GetOrDefault("someKey", "defaultValue")
func (hmap HMap[K, V]) GetOrDefault(key K, defaultValue V) V {
	if value, ok := hmap[key]; ok {
		return value
	}

	return defaultValue
}

// Clear removes all key-value pairs from the HMap.
func (hmap HMap[K, V]) Clear() HMap[K, V] { return hmap.Delete(hmap.Keys()...) }

// Empty checks if the HMap is empty.
func (hmap HMap[K, V]) Empty() bool { return hmap.Len() == 0 }

// Get retrieves the value associated with the given key.
func (hmap HMap[K, V]) Get(k K) V { return hmap[k] }

// Len returns the number of key-value pairs in the HMap.
func (hmap HMap[K, V]) Len() int { return len(hmap) }

// Ne checks if two HMaps are not equal.
func (hmap HMap[K, V]) Ne(other HMap[K, V]) bool { return !hmap.Eq(other) }

// NotEmpty checks if the HMap is not empty.
func (hmap HMap[K, V]) NotEmpty() bool { return !hmap.Empty() }

// Set sets the value for the given key in the HMap.
func (hmap HMap[K, V]) Set(k K, v V) HMap[K, V] { hmap[k] = v; return hmap }

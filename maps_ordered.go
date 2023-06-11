package hg

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/x0xO/hg/pkg/deref"
	"github.com/x0xO/hg/pkg/ref"
)

// NewHMapOrd creates a new ordered HMap with the specified size (if provided).
// An ordered HMap is an HMap that maintains the order of its key-value pairs based on the
// insertion order. If no size is provided, the default size will be used.
//
// Parameters:
//
// - size ...int: (Optional) The initial size of the ordered HMap. If not provided, a default size
// will be used.
//
// Returns:
//
// - *hMapOrd[K, V]: A pointer to a new ordered HMap with the specified initial size (or default
// size if not provided).
//
// Example usage:
//
//	hmapOrd := hg.NewHMapOrd[string, int](10)
//
// Creates a new ordered HMap with an initial size of 10.
func NewHMapOrd[K comparable, V any](size ...int) *HMapOrd[K, V] {
	if len(size) == 0 {
		return ref.Of(make(HMapOrd[K, V], 0))
	}

	return ref.Of(make(HMapOrd[K, V], 0, size[0]))
}

// HMapOrdFromHMap converts a standard HMap to an ordered HMap.
// The resulting ordered HMap will maintain the order of its key-value pairs based on the order of
// insertion.
// This function is useful when you want to create an ordered HMap from an existing HMap.
//
// Parameters:
//
// - hmap HMap[K, V]: The input HMap to be converted to an ordered HMap.
//
// Returns:
//
// - *hMapOrd[K, V]: A pointer to a new ordered HMap containing the same key-value pairs as the
// input HMap.
//
// Example usage:
//
//	hmapOrd := hg.HMapOrdFromHMap[string, int](hmap)
//
// Converts the standard HMap 'hmap' to an ordered HMap.
func HMapOrdFromHMap[K comparable, V any](hmap HMap[K, V]) *HMapOrd[K, V] {
	hmaps := NewHMapOrd[K, V](hmap.Len())
	hmap.ForEach(func(k K, v V) { hmaps.Set(k, v) })

	return hmaps
}

// HMapOrdFromMap converts a standard Go map to an ordered HMap.
// The resulting ordered HMap will maintain the order of its key-value pairs based on the order of
// insertion.
// This function is useful when you want to create an ordered HMap from an existing Go map.
//
// Parameters:
//
// - m map[K]V: The input Go map to be converted to an ordered HMap.
//
// Returns:
//
// - *hMapOrd[K, V]: A pointer to a new ordered HMap containing the same key-value pairs as the
// input Go map.
//
// Example usage:
//
//	hmapOrd := hg.HMapOrdFromMap[string, int](goMap)
//
// Converts the standard Go map 'map[K]V' to an ordered HMap.
func HMapOrdFromMap[K comparable, V any](m map[K]V) *HMapOrd[K, V] {
	return HMapOrdFromHMap(HMapFromMap(m))
}

// SortBy sorts the ordered HMap by a custom comparison function.
// The comparison function should return true if the element at index i should be sorted before the
// element at index j. This function is useful when you want to sort the key-value pairs in an
// ordered HMap based on a custom comparison logic.
//
// Parameters:
//
// - f func(i, j int) bool: The custom comparison function used for sorting the ordered HMap.
//
// Returns:
//
// - *hMapOrd[K, V]: A pointer to the same ordered HMap, sorted according to the custom comparison
// function.
//
// Example usage:
//
//	hmapo.SortBy(func(i, j int) bool { return (*hmapo)[i].Key < (*hmapo)[j].Key })
//	hmapo.SortBy(func(i, j int) bool { return (*hmapo)[i].Value < (*hmapo)[j].Value })
func (hmapo *HMapOrd[K, V]) SortBy(f func(i, j int) bool) *HMapOrd[K, V] {
	sort.Slice(*hmapo, f)
	return hmapo
}

// Clone creates a new ordered HMap with the same key-value pairs.
func (hmapo *HMapOrd[K, V]) Clone() *HMapOrd[K, V] {
	result := NewHMapOrd[K, V](hmapo.Len())
	hmapo.ForEach(func(k K, v V) { result.Set(k, v) })

	return result
}

// Copy copies key-value pairs from the source ordered HMap to the current ordered HMap.
func (hmapo *HMapOrd[K, V]) Copy(src *HMapOrd[K, V]) *HMapOrd[K, V] {
	src.ForEach(func(k K, v V) { hmapo.Set(k, v) })
	return hmapo
}

// ToHMap converts the ordered HMap to a standard HMap.
func (hmapo *HMapOrd[K, V]) ToHMap() HMap[K, V] {
	hmap := NewHMap[K, V](hmapo.Len())
	hmapo.ForEach(func(k K, v V) { hmap.Set(k, v) })

	return hmap
}

// Map applies a provided function to all key-value pairs in the ordered HMap and returns a new
// ordered HMap with the results. The provided function should take the key and value as input and
// return a new key-value pair as output. This function is useful when you want to transform the
// key-value pairs of an ordered HMap using a custom function.
//
// Parameters:
//
// - fn func(K, V) (K, V): The custom function that takes the key and value as input and returns a
// new key-value pair.
//
// Returns:
//
// - *hMapOrd[K, V]: A pointer to a new ordered HMap containing the key-value pairs after applying
// the custom function.
//
// Example usage:
//
//	hmapo.Map(func(k string, v int) (string, int) {
//		return strings.ToUpper(k), v * 2
//	}) // Transforms the keys to uppercase and doubles the values in the ordered HMap.
func (hmapo *HMapOrd[K, V]) Map(fn func(K, V) (K, V)) *HMapOrd[K, V] {
	result := NewHMapOrd[K, V](hmapo.Len())
	hmapo.ForEach(func(k K, v V) { result.Set(fn(k, v)) })

	return result
}

// Filter filters the ordered HMap based on a provided predicate function,
// returning a new ordered HMap with only the key-value pairs that satisfy the predicate.
// The predicate function should take the key and value as input and return a boolean value.
// This function is useful when you want to create a new ordered HMap containing only the key-value
// pairs that meet certain criteria.
//
// Parameters:
//
// - fn func(K, V) bool: The predicate function that takes the key and value as input and returns a
// boolean value.
//
// Returns:
//
// - *hMapOrd[K, V]: A pointer to a new ordered HMap containing only the key-value pairs that
// satisfy the predicate.
//
// Example usage:
//
//	hmapo.Filter(func(k string, v int) bool {
//		return v > 10
//	})
//
// Filters the ordered HMap to include only the key-value pairs where the value is greater
// than 10.
func (hmapo *HMapOrd[K, V]) Filter(fn func(K, V) bool) *HMapOrd[K, V] {
	result := NewHMapOrd[K, V](hmapo.Len())

	hmapo.ForEach(func(k K, v V) {
		if fn(k, v) {
			result.Set(k, v)
		}
	})

	return result
}

// Set sets the value for the specified key in the ordered HMap.
func (hmapo *HMapOrd[K, V]) Set(key K, value V) *HMapOrd[K, V] {
	if i := hmapo.index(key); i != -1 {
		deref.Of(hmapo)[i].Value = value
		return hmapo
	}

	hmp := hMapPair[K, V]{key, value}
	*hmapo = append(*hmapo, hmp)

	return hmapo
}

// Get retrieves the value for the specified key, along with a boolean indicating whether the key
// was found in the ordered HMap. This function is useful when you want to access the value
// associated with a key in the ordered HMap, and also check if the key exists in the map.
//
// Parameters:
//
// - key K: The key to search for in the ordered HMap.
//
// Returns:
//
// - V: The value associated with the specified key if found, or the zero value for the value type
// if the key is not found.
//
// - bool: A boolean value indicating whether the key was found in the ordered HMap.
//
// Example usage:
//
//	value, found := hmapo.Get("some_key")
//
// Retrieves the value associated with the key "some_key" and checks if the key exists in the
// ordered HMap.
func (hmapo *HMapOrd[K, V]) Get(key K) (V, bool) {
	if i := hmapo.index(key); i != -1 {
		return deref.Of(hmapo)[i].Value, true
	}

	return *new(V), false // Returns the zero value for type V and false (not found)
}

// GetOrDefault returns the value for a key. If the key does not exist, returns the default value
// instead. This function is useful when you want to access the value associated with a key in the
// ordered HMap, but if the key does not exist, you want to return a specified default value.
//
// Parameters:
//
// - key K: The key to search for in the ordered HMap.
//
// - defaultValue V: The default value to return if the key is not found in the ordered HMap.
//
// Returns:
//
// - V: The value associated with the specified key if found, or the provided default value if the
// key is not found.
//
// Example usage:
//
//	value := hmapo.GetOrDefault("some_key", "default_value")
//
// Retrieves the value associated with the key "some_key" or returns "default_value" if the key is
// not found.
func (hmapo *HMapOrd[K, V]) GetOrDefault(key K, defaultValue V) V {
	if i := hmapo.index(key); i != -1 {
		return deref.Of(hmapo)[i].Value
	}

	return defaultValue
}

// Invert inverts the key-value pairs in the ordered HMap, creating a new ordered HMap with the
// values as keys and the original keys as values.
func (hmapo *HMapOrd[K, V]) Invert() *HMapOrd[any, K] {
	result := NewHMapOrd[any, K](hmapo.Len())
	hmapo.ForEach(func(k K, v V) { result.Set(v, k) })

	return result
}

func (hmapo *HMapOrd[K, V]) index(key K) int {
	for i, hmap := range *hmapo {
		if hmap.Key == key {
			return i
		}
	}

	return -1
}

// Keys returns an HSlice containing all the keys in the ordered HMap.
func (hmapo *HMapOrd[K, V]) Keys() HSlice[K] {
	keys := NewHSlice[K](0, hmapo.Len())
	hmapo.ForEach(func(k K, _ V) { keys = keys.Append(k) })

	return keys
}

// Values returns an HSlice containing all the values in the ordered HMap.
func (hmapo *HMapOrd[K, V]) Values() HSlice[V] {
	values := NewHSlice[V](0, hmapo.Len())
	hmapo.ForEach(func(_ K, v V) { values = values.Append(v) })

	return values
}

// Delete removes the specified keys from the ordered HMap.
func (hmapo *HMapOrd[K, V]) Delete(keys ...K) *HMapOrd[K, V] {
	for _, key := range keys {
		if i := hmapo.index(key); i != -1 {
			*hmapo = append(deref.Of(hmapo)[:i], deref.Of(hmapo)[i+1:]...)
		}
	}

	return hmapo
}

// ForEach executes a provided function for each key-value pair in the ordered HMap.
// This function is useful when you want to perform an operation or side effect for each key-value
// pair in the ordered HMap.
//
// Parameters:
//
// - fn func(K, V): The function to execute for each key-value pair in the ordered HMap. It takes a
// key (K) and a value (V) as arguments.
//
// Example usage:
//
//	hmapo.ForEach(func(key K, value V) { fmt.Printf("Key: %v, Value: %v\n", key, value) }).
//
// Prints each key-value pair in the ordered HMap.
func (hmapo *HMapOrd[K, V]) ForEach(fn func(K, V)) {
	for _, hmap := range *hmapo {
		fn(hmap.Key, hmap.Value)
	}
}

// Eq compares the current ordered HMap to another ordered HMap and returns true if they are equal.
func (hmapo *HMapOrd[K, V]) Eq(other *HMapOrd[K, V]) bool {
	if hmapo.Len() != other.Len() {
		return false
	}

	for _, hmap := range *hmapo {
		value, ok := other.Get(hmap.Key)
		if !ok || !reflect.DeepEqual(value, hmap.Value) {
			return false
		}
	}

	return true
}

// String returns a string representation of the ordered HMap.
func (hmapo *HMapOrd[K, V]) String() string {
	var builder strings.Builder

	hmapo.ForEach(func(k K, v V) { builder.WriteString(fmt.Sprintf("%v:%v, ", k, v)) })

	return HString(builder.String()).TrimRight(", ").Format("HMapOrd{%s}").String()
}

// Clear removes all key-value pairs from the ordered HMap.
func (hmapo *HMapOrd[K, V]) Clear() *HMapOrd[K, V] { return hmapo.Delete(hmapo.Keys()...) }

// Contains checks if the ordered HMap contains the specified key.
func (hmapo *HMapOrd[K, V]) Contains(key K) bool { return hmapo.index(key) >= 0 }

// Empty checks if the ordered HMap is empty.
func (hmapo *HMapOrd[K, V]) Empty() bool { return hmapo.Len() == 0 }

// Len returns the number of key-value pairs in the ordered HMap.
func (hmapo *HMapOrd[K, V]) Len() int { return len(*hmapo) }

// Ne compares the current ordered HMap to another ordered HMap and returns true if they are not
// equal.
func (hmapo *HMapOrd[K, V]) Ne(other *HMapOrd[K, V]) bool { return !hmapo.Eq(other) }

// NotEmpty checks if the ordered HMap is not empty.
func (hmapo *HMapOrd[K, V]) NotEmpty() bool { return !hmapo.Empty() }

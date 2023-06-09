package hg

import (
	"fmt"
	"strings"
)

// NewHSet creates a new HSet of the specified size or an empty HSet if no size is provided.
func NewHSet[T comparable](size ...int) HSet[T] {
	if len(size) == 0 {
		return make(HSet[T], 0)
	}

	return make(HSet[T], size[0])
}

// HSetOf creates a new generic set containing the provided elements.
func HSetOf[T comparable](values ...T) HSet[T] {
	hset := NewHSet[T](len(values))
	for _, v := range values {
		hset.Add(v)
	}

	return hset
}

// Add adds the provided elements to the set and returns the modified set.
func (s HSet[T]) Add(values ...T) HSet[T] {
	for _, v := range values {
		s[v] = struct{}{}
	}

	return s
}

// Remove removes the specified values from the HSet.
func (s HSet[T]) Remove(values ...T) HSet[T] {
	for _, v := range values {
		delete(s, v)
	}

	return s
}

// Len returns the number of values in the HSet.
func (s HSet[T]) Len() int { return len(s) }

// Contains checks if the HSet contains the specified value.
func (s HSet[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

// ContainsAny checks if the HSet contains any element from another HSet.
func (s HSet[T]) ContainsAny(other HSet[T]) bool {
	for v := range other {
		if s.Contains(v) {
			return true
		}
	}

	return false
}

// ContainsAll checks if the HSet contains all elements from another HSet.
func (s HSet[T]) ContainsAll(other HSet[T]) bool {
	if s.Len() < other.Len() {
		return false
	}

	for v := range other {
		if !s.Contains(v) {
			return false
		}
	}

	return true
}

// Clone creates a new HSet that is a copy of the original HSet.
func (s HSet[T]) Clone() HSet[T] {
	result := NewHSet[T](s.Len())
	s.ForEach(func(t T) { result.Add(t) })

	return result
}

// ForEach applies a function to each value in the HSet.
// The provided function 'fn' should take a value as input parameter and perform an
// operation.
// This function is useful for side effects, as it does not return a new HSet.
//
// Parameters:
//
// - fn func(T): A function that takes a value as input parameter and performs an
// operation.
//
// Example usage:
//
//	originalHSet.ForEach(func(value T) {
//		fmt.Printf("Value: %v\n", value)
//	})
func (s HSet[T]) ForEach(fn func(T)) {
	for value := range s {
		fn(value)
	}
}

// Map returns a new set by applying a given function to each element in the current set.
//
// The function takes one parameter of type T (the same type as the elements of the set)
// and returns a value of type T. The returned value is added to a new set,
// which is then returned as the result.
//
// Parameters:
//
// - fn (func(T) T): The function to be applied to each element of the set.
//
// Returns:
//
// - HSet[T]: A new set containing the results of applying the function to each element
// of the current set.
//
// Example usage:
//
//	s := hg.HSetOf(1, 2, 3)
//	doubled := s.Map(func(val int) int {
//	    return val * 2
//	})
//	fmt.Println(doubled)
//
// Output: [2 4 6].
func (s HSet[T]) Map(fn func(T) T) HSet[T] {
	result := NewHSet[T](s.Len())
	s.ForEach(func(t T) { result.Add(fn(t)) })

	return result
}

// Filter returns a new set containing elements that satisfy a given condition.
//
// The function takes one parameter of type T (the same type as the elements of the set)
// and returns a boolean value. If the returned value is true, the element is added
// to a new set, which is then returned as the result.
//
// Parameters:
//
// - fn (func(T) bool): The function to be applied to each element of the set
// to determine if it should be included in the result.
//
// Returns:
//
// - HSet[T]: A new set containing the elements that satisfy the given condition.
//
// Example usage:
//
//	s := hg.HSetOf(1, 2, 3, 4, 5)
//	even := s.Filter(func(val int) bool {
//	    return val%2 == 0
//	})
//	fmt.Println(even)
//
// Output: [2 4].
func (s HSet[T]) Filter(fn func(T) bool) HSet[T] {
	result := NewHSet[T]()

	s.ForEach(func(t T) {
		if fn(t) {
			result.Add(t)
		}
	})

	return result
}

// HSlice returns a new HSlice with the same elements as the HSet[T].
func (s HSet[T]) HSlice() HSlice[T] {
	hsl := NewHSlice[T](0, s.Len())
	s.ForEach(func(v T) { hsl = hsl.Append(v) })

	return hsl
}

// ToSlice returns a new slice with the same elements as the HSet[T].
func (s HSet[T]) ToSlice() []T { return s.HSlice() }

// Intersection returns the intersection of the current set and another set, i.e., elements
// present in both sets.
//
// Parameters:
//
// - other HSet[T]: The other set to calculate the intersection with.
//
// Returns:
//
// - HSet[T]: A new HSet containing the intersection of the two sets.
//
// Example usage:
//
//	s1 := hg.HSetOf(1, 2, 3, 4, 5)
//	s2 := hg.HSetOf(4, 5, 6, 7, 8)
//	intersection := s1.Intersection(s2)
//
// The resulting intersection will be: [4, 5].
func (s HSet[T]) Intersection(other HSet[T]) HSet[T] {
	result := NewHSet[T]()

	s.ForEach(func(t T) {
		if other.Contains(t) {
			result.Add(t)
		}
	})

	return result
}

// Difference returns the difference between the current set and another set,
// i.e., elements present in the current set but not in the other set.
//
// Parameters:
//
// - other HSet[T]: The other set to calculate the difference with.
//
// Returns:
//
// - HSet[T]: A new HSet containing the difference between the two sets.
//
// Example usage:
//
//	s1 := hg.HSetOf(1, 2, 3, 4, 5)
//	s2 := hg.HSetOf(4, 5, 6, 7, 8)
//	diff := s1.Difference(s2)
//
// The resulting diff will be: [1, 2, 3].
func (s HSet[T]) Difference(other HSet[T]) HSet[T] {
	result := NewHSet[T]()

	s.ForEach(func(t T) {
		if !other.Contains(t) {
			result.Add(t)
		}
	})

	return result
}

// Union returns a new set containing the unique elements of the current set and the provided
// other set.
//
// Parameters:
//
// - other HSet[T]: The other set to create the union with.
//
// Returns:
//
// - HSet[T]: A new HSet containing the unique elements of the current set and the provided
// other set.
//
// Example usage:
//
//	s1 := hg.HSetOf(1, 2, 3)
//	s2 := hg.HSetOf(3, 4, 5)
//	union := s1.Union(s2)
//
// The resulting union set will be: [1, 2, 3, 4, 5].
func (s HSet[T]) Union(other HSet[T]) HSet[T] {
	result := NewHSet[T](s.Len() + other.Len())
	return result.Add(s.HSlice()...).Add(other.HSlice()...)
}

// SymmetricDifference returns the symmetric difference between the current set and another
// set, i.e., elements present in either the current set or the other set but not in both.
//
// Parameters:
//
// - other HSet[T]: The other set to calculate the symmetric difference with.
//
// Returns:
//
// - HSet[T]: A new HSet containing the symmetric difference between the two sets.
//
// Example usage:
//
//	s1 := hg.HSetOf(1, 2, 3, 4, 5)
//	s2 := hg.HSetOf(4, 5, 6, 7, 8)
//	symDiff := s1.SymmetricDifference(s2)
//
// The resulting symDiff will be: [1, 2, 3, 6, 7, 8].
func (s HSet[T]) SymmetricDifference(other HSet[T]) HSet[T] {
	return s.Difference(other).Union(other.Difference(s))
}

// Subset checks if the current set 's' is a subset of the provided 'other' set.
// A set 's' is a subset of 'other' if all elements of 's' are also elements of 'other'.
//
// Parameters:
//
// - other HSet[T]: The other set to compare with.
//
// Returns:
//
// - bool: true if 's' is a subset of 'other', false otherwise.
//
// Example usage:
//
//	s1 := hg.HSetOf(1, 2, 3)
//	s2 := hg.HSetOf(1, 2, 3, 4, 5)
//	isSubset := s1.Subset(s2) // Returns true
func (s HSet[T]) Subset(other HSet[T]) bool { return other.ContainsAll(s) }

// Superset checks if the current set 's' is a superset of the provided 'other' set.
// A set 's' is a superset of 'other' if all elements of 'other' are also elements of 's'.
//
// Parameters:
//
// - other HSet[T]: The other set to compare with.
//
// Returns:
//
// - bool: true if 's' is a superset of 'other', false otherwise.
//
// Example usage:
//
//	s1 := hg.HSetOf(1, 2, 3, 4, 5)
//	s2 := hg.HSetOf(1, 2, 3)
//	isSuperset := s1.Superset(s2) // Returns true
func (s HSet[T]) Superset(other HSet[T]) bool { return s.ContainsAll(other) }

// Eq checks if two HSets are equal.
func (s HSet[T]) Eq(other HSet[T]) bool {
	if s.Len() != other.Len() {
		return false
	}

	for v := range other {
		if !s.Contains(v) {
			return false
		}
	}

	return true
}

// Ne checks if two HSets are not equal.
func (s HSet[T]) Ne(other HSet[T]) bool { return !s.Eq(other) }

// Clear removes all values from the HSet.
func (s HSet[T]) Clear() HSet[T] { return s.Remove(s.HSlice()...) }

// Empty checks if the HSet is empty.
func (s HSet[T]) Empty() bool { return s.Len() == 0 }

// String returns a string representation of the HSet.
func (s HSet[T]) String() string {
	var builder strings.Builder

	s.ForEach(func(v T) { builder.WriteString(fmt.Sprintf("%v, ", v)) })

	return HString(builder.String()).AddPrefix("HSet[").TrimRight(", ").Add("]").String()
}

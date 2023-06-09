package hg_test

import (
	"testing"

	"github.com/x0xO/hg"
)

func TestHSetDifference(t *testing.T) {
	set1 := hg.HSetOf(1, 2, 3, 4)
	set2 := hg.HSetOf(3, 4, 5, 6)
	set5 := hg.HSetOf(1, 2)
	set6 := hg.HSetOf(2, 3, 4)

	set3 := set1.Difference(set2)
	set4 := set2.Difference(set1)
	set7 := set5.Difference(set6)
	set8 := set6.Difference(set5)

	if set3.Len() != 2 || set3.Ne(hg.HSetOf(1, 2)) {
		t.Errorf("Unexpected result: %v", set3)
	}

	if set4.Len() != 2 || set4.Ne(hg.HSetOf(5, 6)) {
		t.Errorf("Unexpected result: %v", set4)
	}

	if set7.Len() != 1 || set7.Ne(hg.HSetOf(1)) {
		t.Errorf("Unexpected result: %v", set7)
	}

	if set8.Len() != 2 || set8.Ne(hg.HSetOf(3, 4)) {
		t.Errorf("Unexpected result: %v", set8)
	}
}

func TestHSetSymmetricDifference(t *testing.T) {
	set1 := hg.NewHSet[int](10)
	set2 := set1.Clone()
	result := set1.SymmetricDifference(set2)

	if !result.Empty() {
		t.Errorf("SymmetricDifference between equal sets should be empty, got %v", result)
	}

	set1 = hg.HSetOf(0, 1, 2, 3, 4)
	set2 = hg.HSetOf(5, 6, 7, 8, 9)
	result = set1.SymmetricDifference(set2)
	expected := set1.Union(set2)

	if !result.Eq(expected) {
		t.Errorf(
			"SymmetricDifference between disjoint sets should be their union, expected %v but got %v",
			expected,
			result,
		)
	}

	set1 = hg.HSetOf(0, 1, 2, 3, 4, 5)
	set2 = hg.HSetOf(4, 5, 6, 7, 8)
	result = set1.SymmetricDifference(set2)
	expected = hg.HSetOf(0, 1, 2, 3, 6, 7, 8)

	if !result.Eq(expected) {
		t.Errorf(
			"SymmetricDifference between sets with common elements should be correct, expected %v but got %v",
			expected,
			result,
		)
	}
}

func TestHSetIntersection(t *testing.T) {
	set1 := hg.HSet[int]{}
	set2 := hg.HSet[int]{}

	set1 = set1.Add(1, 2, 3)
	set2 = set2.Add(2, 3, 4)

	set3 := set1.Intersection(set2)

	if !set3.Contains(2) || !set3.Contains(3) {
		t.Error("Intersection failed")
	}
}

func TestHSetUnion(t *testing.T) {
	set1 := hg.NewHSet[int]().Add(1, 2, 3)
	set2 := hg.NewHSet[int]().Add(2, 3, 4)
	set3 := hg.NewHSet[int]().Add(1, 2, 3, 4)

	result := set1.Union(set2)

	if result.Len() != 4 {
		t.Errorf("Union(%v, %v) returned %v; expected %v", set1, set2, result, set3)
	}

	for v := range set3 {
		if !result.Contains(v) {
			t.Errorf("Union(%v, %v) missing element %v", set1, set2, v)
		}
	}
}

func TestHSetSubset(t *testing.T) {
	tests := []struct {
		name  string
		s     hg.HSet[int]
		other hg.HSet[int]
		want  bool
	}{
		{
			name:  "test_subset_1",
			s:     hg.HSetOf(1, 2, 3),
			other: hg.HSetOf(1, 2, 3, 4, 5),
			want:  true,
		},
		{
			name:  "test_subset_2",
			s:     hg.HSetOf(1, 2, 3, 4),
			other: hg.HSetOf(1, 2, 3),
			want:  false,
		},
		{
			name:  "test_subset_3",
			s:     hg.HSetOf(5, 4, 3, 2, 1),
			other: hg.HSetOf(1, 2, 3, 4, 5),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Subset(tt.other); got != tt.want {
				t.Errorf("Subset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHSetSuperset(t *testing.T) {
	tests := []struct {
		name  string
		s     hg.HSet[int]
		other hg.HSet[int]
		want  bool
	}{
		{
			name:  "test_superset_1",
			s:     hg.HSetOf(1, 2, 3, 4, 5),
			other: hg.HSetOf(1, 2, 3),
			want:  true,
		},
		{
			name:  "test_superset_2",
			s:     hg.HSetOf(1, 2, 3),
			other: hg.HSetOf(1, 2, 3, 4),
			want:  false,
		},
		{
			name:  "test_superset_3",
			s:     hg.HSetOf(1, 2, 3, 4, 5),
			other: hg.HSetOf(5, 4, 3, 2, 1),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Superset(tt.other); got != tt.want {
				t.Errorf("Superset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHSetEq(t *testing.T) {
	tests := []struct {
		name  string
		s     hg.HSet[int]
		other hg.HSet[int]
		want  bool
	}{
		{
			name:  "test_eq_1",
			s:     hg.HSetOf(1, 2, 3),
			other: hg.HSetOf(1, 2, 3),
			want:  true,
		},
		{
			name:  "test_eq_2",
			s:     hg.HSetOf(1, 2, 3),
			other: hg.HSetOf(1, 2, 4),
			want:  false,
		},
		{
			name:  "test_eq_3",
			s:     hg.HSetOf(1, 2, 3),
			other: hg.HSetOf(3, 2, 1),
			want:  true,
		},
		{
			name:  "test_eq_4",
			s:     hg.HSetOf(1, 2, 3, 4),
			other: hg.HSetOf(1, 2, 3),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Eq(tt.other); got != tt.want {
				t.Errorf("Eq() = %v, want %v", got, tt.want)
			}
		})
	}
}

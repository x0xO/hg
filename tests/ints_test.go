package hg_test

import (
	"testing"

	"github.com/x0xO/hg"
	"github.com/x0xO/hg/pkg/iter"
)

func TestHIntIsPositive(t *testing.T) {
	tests := []struct {
		name string
		hi   hg.HInt
		want bool
	}{
		{"positive", 1, true},
		{"negative", -1, false},
		{"zero", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hi.IsPositive(); got != tt.want {
				t.Errorf("HInt.IsPositive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHIntNegative(t *testing.T) {
	tests := []struct {
		name string
		hi   hg.HInt
		want bool
	}{
		{"positive", 1, false},
		{"negative", -1, true},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hi.IsNegative(); got != tt.want {
				t.Errorf("HInt.IsNegative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomRange(t *testing.T) {
	for range iter.N(100) {
		min := hg.NewHInt(100).Random()
		max := hg.NewHInt(100).Random().Add(min)

		r := hg.NewHInt().RandomRange(min, max)
		if r.Lt(min) || r.Gt(max) {
			t.Errorf("RandomRange(%d, %d) = %d, want in range [%d, %d]", min, max, r, min, max)
		}
	}
}

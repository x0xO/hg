package hg_test

import (
	"testing"

	"github.com/x0xO/hg"
)

func TestHFloatCompare(t *testing.T) {
	testCases := []struct {
		hf1      hg.HFloat
		hf2      hg.HFloat
		expected hg.HInt
	}{
		{3.14, 6.28, -1},
		{6.28, 3.14, 1},
		{1.23, 1.23, 0},
		{-2.5, 2.5, -1},
	}

	for _, tc := range testCases {
		result := tc.hf1.Compare(tc.hf2)
		if !result.Eq(tc.expected) {
			t.Errorf("Compare(%f, %f): expected %d, got %d", tc.hf1, tc.hf2, tc.expected, result)
		}
	}
}

func TestHFloatEq(t *testing.T) {
	testCases := []struct {
		hf1      hg.HFloat
		hf2      hg.HFloat
		expected bool
	}{
		{3.14, 6.28, false},
		{1.23, 1.23, true},
		{0.0, 0.0, true},
		{-2.5, 2.5, false},
	}

	for _, tc := range testCases {
		result := tc.hf1.Eq(tc.hf2)
		if result != tc.expected {
			t.Errorf("Eq(%f, %f): expected %t, got %t", tc.hf1, tc.hf2, tc.expected, result)
		}
	}
}

func TestHFloatNe(t *testing.T) {
	testCases := []struct {
		hf1      hg.HFloat
		hf2      hg.HFloat
		expected bool
	}{
		{3.14, 6.28, true},
		{1.23, 1.23, false},
		{0.0, 0.0, false},
		{-2.5, 2.5, true},
	}

	for _, tc := range testCases {
		result := tc.hf1.Ne(tc.hf2)
		if result != tc.expected {
			t.Errorf("Ne(%f, %f): expected %t, got %t", tc.hf1, tc.hf2, tc.expected, result)
		}
	}
}

func TestHFloatGt(t *testing.T) {
	testCases := []struct {
		hf1      hg.HFloat
		hf2      hg.HFloat
		expected bool
	}{
		{3.14, 6.28, false},
		{6.28, 3.14, true},
		{1.23, 1.23, false},
		{-2.5, 2.5, false},
	}

	for _, tc := range testCases {
		result := tc.hf1.Gt(tc.hf2)
		if result != tc.expected {
			t.Errorf("Gt(%f, %f): expected %t, got %t", tc.hf1, tc.hf2, tc.expected, result)
		}
	}
}

func TestHFloatLt(t *testing.T) {
	testCases := []struct {
		hf1      hg.HFloat
		hf2      hg.HFloat
		expected bool
	}{
		{3.14, 6.28, true},
		{6.28, 3.14, false},
		{1.23, 1.23, false},
		{-2.5, 2.5, true},
	}
	for _, tc := range testCases {
		result := tc.hf1.Lt(tc.hf2)
		if result != tc.expected {
			t.Errorf("Lt(%f, %f): expected %t, got %t", tc.hf1, tc.hf2, tc.expected, result)
		}
	}
}

func TestHFloatRoundDecimal(t *testing.T) {
	testCases := []struct {
		value    hg.HFloat
		decimals int
		expected hg.HFloat
	}{
		{3.1415926535, 2, 3.14},
		{3.1415926535, 3, 3.142},
		{100.123456789, 4, 100.1235},
		{-5.6789, 1, -5.7},
		{12345.6789, 0, 12346},
	}

	for _, testCase := range testCases {
		result := testCase.value.RoundDecimal(testCase.decimals)
		if result != testCase.expected {
			t.Errorf(
				"Failed: value=%.10f decimals=%d, expected=%.10f, got=%.10f\n",
				testCase.value,
				testCase.decimals,
				testCase.expected,
				result,
			)
		}
	}
}

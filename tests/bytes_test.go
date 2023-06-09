package hg_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/x0xO/hg"
)

func TestHBytesCompare(t *testing.T) {
	testCases := []struct {
		hbs1     hg.HBytes
		hbs2     hg.HBytes
		expected hg.HInt
	}{
		{[]byte("apple"), []byte("banana"), -1},
		{[]byte("banana"), []byte("apple"), 1},
		{[]byte("banana"), []byte("banana"), 0},
		{[]byte("apple"), []byte("Apple"), 1},
		{[]byte(""), []byte(""), 0},
	}

	for _, tc := range testCases {
		result := tc.hbs1.Compare(tc.hbs2)
		if !result.Eq(tc.expected) {
			t.Errorf(
				"HBytes.Compare(%q, %q): expected %d, got %d",
				tc.hbs1,
				tc.hbs2,
				tc.expected,
				result,
			)
		}
	}
}

func TestHBytesEq(t *testing.T) {
	testCases := []struct {
		hbs1     hg.HBytes
		hbs2     hg.HBytes
		expected bool
	}{
		{[]byte("apple"), []byte("banana"), false},
		{[]byte("banana"), []byte("banana"), true},
		{[]byte("Apple"), []byte("apple"), false},
		{[]byte(""), []byte(""), true},
	}

	for _, tc := range testCases {
		result := tc.hbs1.Eq(tc.hbs2)
		if result != tc.expected {
			t.Errorf(
				"HBytes.Eq(%q, %q): expected %t, got %t",
				tc.hbs1,
				tc.hbs2,
				tc.expected,
				result,
			)
		}
	}
}

func TestHBytesNe(t *testing.T) {
	testCases := []struct {
		hbs1     hg.HBytes
		hbs2     hg.HBytes
		expected bool
	}{
		{[]byte("apple"), []byte("banana"), true},
		{[]byte("banana"), []byte("banana"), false},
		{[]byte("Apple"), []byte("apple"), true},
		{[]byte(""), []byte(""), false},
	}

	for _, tc := range testCases {
		result := tc.hbs1.Ne(tc.hbs2)
		if result != tc.expected {
			t.Errorf(
				"HBytes.Ne(%q, %q): expected %t, got %t",
				tc.hbs1,
				tc.hbs2,
				tc.expected,
				result,
			)
		}
	}
}

func TestHBytesGt(t *testing.T) {
	testCases := []struct {
		hbs1     hg.HBytes
		hbs2     hg.HBytes
		expected bool
	}{
		{[]byte("apple"), []byte("banana"), false},
		{[]byte("banana"), []byte("apple"), true},
		{[]byte("Apple"), []byte("apple"), false},
		{[]byte("banana"), []byte("banana"), false},
		{[]byte(""), []byte(""), false},
	}

	for _, tc := range testCases {
		result := tc.hbs1.Gt(tc.hbs2)
		if result != tc.expected {
			t.Errorf(
				"HBytes.Gt(%q, %q): expected %t, got %t",
				tc.hbs1,
				tc.hbs2,
				tc.expected,
				result,
			)
		}
	}
}

func TestHBytesLt(t *testing.T) {
	testCases := []struct {
		hbs1     hg.HBytes
		hbs2     hg.HBytes
		expected bool
	}{
		{[]byte("apple"), []byte("banana"), true},
		{[]byte("banana"), []byte("apple"), false},
		{[]byte("Apple"), []byte("apple"), true},
		{[]byte("banana"), []byte("banana"), false},
		{[]byte(""), []byte(""), false},
	}

	for _, tc := range testCases {
		result := tc.hbs1.Lt(tc.hbs2)
		if result != tc.expected {
			t.Errorf(
				"HBytes.Lt(%q, %q): expected %t, got %t",
				tc.hbs1,
				tc.hbs2,
				tc.expected,
				result,
			)
		}
	}
}

func TestHBytesNormalizeNFC(t *testing.T) {
	testCases := []struct {
		input    hg.HBytes
		expected hg.HBytes
	}{
		{[]byte("Mëtàl Hëàd"), []byte("Mëtàl Hëàd")},
		{[]byte("Café"), []byte("Café")},
		{[]byte("Ĵūņě"), []byte("Ĵūņě")},
		{[]byte("A\u0308"), []byte("Ä")},
		{[]byte("o\u0308"), []byte("ö")},
		{[]byte("u\u0308"), []byte("ü")},
		{[]byte("O\u0308"), []byte("Ö")},
		{[]byte("U\u0308"), []byte("Ü")},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			output := tc.input.NormalizeNFC()
			if string(output) != string(tc.expected) {
				t.Errorf("HBytes.NormalizeNFC(%q) = %q; want %q", tc.input, output, tc.expected)
			}
		})
	}
}

func TestHBytesReader(t *testing.T) {
	tests := []struct {
		name     string
		hbytes   hg.HBytes
		expected []byte
	}{
		{"Empty HBytes", hg.HBytes{}, []byte{}},
		{"Single byte HBytes", hg.HBytes{0x41}, []byte{0x41}},
		{
			"Multiple bytes HBytes",
			hg.HBytes{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			[]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
		},
		{
			"HBytes with various values",
			hg.HBytes{0x00, 0xff, 0x80, 0x7f},
			[]byte{0x00, 0xff, 0x80, 0x7f},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := test.hbytes.Reader()
			resultBytes, err := io.ReadAll(reader)
			if err != nil {
				t.Fatalf("Error reading from *bytes.Reader: %v", err)
			}

			if !bytes.Equal(resultBytes, test.expected) {
				t.Errorf("HBytes.Reader() content = %v, expected %v", resultBytes, test.expected)
			}
		})
	}
}

func TestHBytesContainsAny(t *testing.T) {
	testCases := []struct {
		hbs      hg.HBytes
		bss      []hg.HBytes
		expected bool
	}{
		{
			hbs:      hg.HBytes("Hello, world!"),
			bss:      []hg.HBytes{hg.HBytes("world"), hg.HBytes("Go")},
			expected: true,
		},
		{
			hbs:      hg.HBytes("Welcome to the HumanGo-1!"),
			bss:      []hg.HBytes{hg.HBytes("Go-3"), hg.HBytes("Go-4")},
			expected: false,
		},
		{
			hbs:      hg.HBytes("Have a great day!"),
			bss:      []hg.HBytes{hg.HBytes(""), hg.HBytes(" ")},
			expected: true,
		},
		{
			hbs:      hg.HBytes(""),
			bss:      []hg.HBytes{hg.HBytes("Hello"), hg.HBytes("world")},
			expected: false,
		},
		{
			hbs:      hg.HBytes(""),
			bss:      []hg.HBytes{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		result := tc.hbs.ContainsAny(tc.bss...)
		if result != tc.expected {
			t.Errorf(
				"HBytes.ContainsAny(%v, %v) = %v; want %v",
				tc.hbs,
				tc.bss,
				result,
				tc.expected,
			)
		}
	}
}

func TestHBytesContainsAll(t *testing.T) {
	testCases := []struct {
		hbs      hg.HBytes
		bss      []hg.HBytes
		expected bool
	}{
		{
			hbs:      hg.HBytes("Hello, world!"),
			bss:      []hg.HBytes{hg.HBytes("Hello"), hg.HBytes("world")},
			expected: true,
		},
		{
			hbs:      hg.HBytes("Welcome to the HumanGo-1!"),
			bss:      []hg.HBytes{hg.HBytes("Go-3"), hg.HBytes("Go-4")},
			expected: false,
		},
		{
			hbs:      hg.HBytes("Have a great day!"),
			bss:      []hg.HBytes{hg.HBytes("Have"), hg.HBytes("a")},
			expected: true,
		},
		{
			hbs:      hg.HBytes(""),
			bss:      []hg.HBytes{hg.HBytes("Hello"), hg.HBytes("world")},
			expected: false,
		},
		{
			hbs:      hg.HBytes("Hello, world!"),
			bss:      []hg.HBytes{},
			expected: true,
		},
	}

	for _, tc := range testCases {
		result := tc.hbs.ContainsAll(tc.bss...)
		if result != tc.expected {
			t.Errorf(
				"HBytes.ContainsAll(%v, %v) = %v; want %v",
				tc.hbs,
				tc.bss,
				result,
				tc.expected,
			)
		}
	}
}

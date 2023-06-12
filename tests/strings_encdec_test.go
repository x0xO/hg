package hg_test

import (
	"testing"

	"github.com/x0xO/hg"
	"github.com/x0xO/hg/pkg/iter"
)

func TestHStringBase64Encode(t *testing.T) {
	tests := []struct {
		name string
		e    hg.HString
		want hg.HString
	}{
		{"empty", "", ""},
		{"hello", "hello", "aGVsbG8="},
		{"hello world", "hello world", "aGVsbG8gd29ybGQ="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Enc().Base64(); got != tt.want {
				t.Errorf("enc.Base64Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHStringBase64Decode(t *testing.T) {
	tests := []struct {
		name string
		d    hg.HString
		want hg.HString
	}{
		{"base64 decode", "aGVsbG8gd29ybGQ=", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Dec().Base64(); got != tt.want {
				t.Errorf("dec.Base64Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHStringRot13(t *testing.T) {
	input := hg.HString("hello world")
	expected := hg.HString("uryyb jbeyq")
	actual := input.Enc().Rot13()

	if actual != expected {
		t.Errorf("Rot13Encode(%q) = %q; expected %q", input, actual, expected)
	}

	input = hg.HString("uryyb jbeyq")
	expected = hg.HString("hello world")
	actual = input.Dec().Rot13()

	if actual != expected {
		t.Errorf("Rot13Decode(%q) = %q; expected %q", input, actual, expected)
	}
}

func TestHStringXOR(t *testing.T) {
	for range iter.N(100) {
		input := hg.NewHString().Random(hg.NewHInt().RandomRange(30, 100).Int())
		key := hg.NewHString().Random(10)
		obfuscated := input.Enc().XOR(key)
		deobfuscated := obfuscated.Dec().XOR(key)

		if input != deobfuscated {
			t.Errorf("expected %s, but got %s", input, deobfuscated)
		}
	}
}

func TestXOR(t *testing.T) {
	tests := []struct {
		input string
		key   string
		want  string
	}{
		{"01", "qsCDE", "AB"},
		{"123", "ABCDE", "ppp"},
		{"12345", "98765", "\x08\x0a\x04\x02\x00"},
		{"Hello", "wORLD", "?*> +"},
		// {"Hello,", "World!", "\x0f\x0a\x1e\x00\x0b\x0d"},
		// {"`c345", "QQ", "12345"},
		{"abcde", "01234", "QSQWQ"},
		{"lowercase", "9?'      ", "UPPERCASE"},
		{"test", "", "test"},
		{"test", "test", "\x00\x00\x00\x00"},
	}

	for _, tt := range tests {
		got := hg.NewHString(tt.input).Enc().XOR(hg.HString(tt.key))
		if got != hg.HString(tt.want) {
			t.Errorf("XOR(%q, %q) = %q; want %q", tt.input, tt.key, got, tt.want)
		}
	}
}

func TestGzFlateDecode(t *testing.T) {
	testCases := []struct {
		name     string
		input    hg.HString
		expected hg.HString
	}{
		{"Empty input", "", ""},
		{"Valid compressed data", "8kjNycnXUXCvcstJLElVBAAAAP//AQAA//8=", "Hello, GzFlate!"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.Dec().GzFlate()
			if result.Ne(tc.expected) {
				t.Errorf("GzFlateDecode, expected: %s, got: %s", tc.expected, result)
			}
		})
	}
}

func TestGzFlateEncode(t *testing.T) {
	testCases := []struct {
		name     string
		input    hg.HString
		expected hg.HString
	}{
		{"Empty input", "", "AAAA//8BAAD//w=="},
		{"Valid input", "Hello, GzFlate!", "8kjNycnXUXCvcstJLElVBAAAAP//AQAA//8="},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.Enc().GzFlate()
			if result.Ne(tc.expected) {
				t.Errorf("GzFlateEncode, expected: %s, got: %s", tc.expected, result)
			}
		})
	}
}

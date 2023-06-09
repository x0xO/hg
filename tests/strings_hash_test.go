package hg_test

import (
	"testing"

	"github.com/x0xO/hg"
)

func TestHStringMD5(t *testing.T) {
	tests := []struct {
		name string
		h    hg.HString
		want hg.HString
	}{
		{
			name: "empty",
			h:    hg.NewHString().Hash().MD5(),
			want: hg.HString("d41d8cd98f00b204e9800998ecf8427e"),
		},
		{
			name: "hello",
			h:    hg.NewHString("hello").Hash().MD5(),
			want: hg.HString("5d41402abc4b2a76b9719d911017c592"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h; got != tt.want {
				t.Errorf("hg.HString.MD5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHStringSHA1(t *testing.T) {
	h := hg.NewHString("Hello, world!")
	expected := "943a702d06f34599aee1f8da8ef9f7296031d699"

	actual := h.Hash().SHA1().String()
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestHStringSHA256(t *testing.T) {
	tests := []struct {
		name string
		h    hg.HString
		want hg.HString
	}{
		{
			"empty",
			hg.HString(""),
			"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{"a", hg.HString("a"), "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb"},
		{
			"abc",
			hg.HString("abc"),
			"ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
		},
		{
			"message digest",
			hg.HString("message digest"),
			"f7846f55cf23e14eebeab5b4e1550cad5b509e3348fbc4efa3a1413d393cb650",
		},
		{
			"secure hash algorithm",
			hg.HString("secure hash algorithm"),
			"f30ceb2bb2829e79e4ca9753d35a8ecc00262d164cc077080295381cbd643f0d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Hash().SHA256(); got != tt.want {
				t.Errorf("hg.HString.SHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHStringSHA512(t *testing.T) {
	tests := []struct {
		name string
		h    hg.HString
		want hg.HString
	}{
		{
			"empty",
			hg.HString(""),
			"cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e",
		},
		{
			"hello",
			hg.HString("hello"),
			"9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043",
		},
		{
			"hello world",
			hg.HString("hello world"),
			"309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.Hash().SHA512(); got != tt.want {
				t.Errorf("hg.HString.SHA512() = %v, want %v", got, tt.want)
			}
		})
	}
}

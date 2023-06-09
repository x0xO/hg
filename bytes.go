package hg

import (
	"bytes"
	"regexp"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// NewHBytes creates a new HBytes value.
func NewHBytes(bs ...[]byte) HBytes {
	if len(bs) != 0 {
		return bs[0]
	}

	return HBytes("")
}

// Reverse returns a new HBytes with the order of its runes reversed.
func (hbs HBytes) Reverse() HBytes {
	reversed := make(HBytes, hbs.Len())
	i := 0

	for hbs.Len() > 0 {
		r, size := utf8.DecodeLastRune(hbs)
		hbs = hbs[:hbs.Len()-size]
		i += utf8.EncodeRune(reversed[i:], r)
	}

	return reversed
}

// Replace replaces the first 'n' occurrences of 'oldB' with 'newB' in the HBytes.
func (hbs HBytes) Replace(oldB, newB HBytes, n int) HBytes {
	return bytes.Replace(hbs, oldB, newB, n)
}

// ReplaceAll replaces all occurrences of 'oldB' with 'newB' in the HBytes.
func (hbs HBytes) ReplaceAll(oldB, newB HBytes) HBytes { return bytes.ReplaceAll(hbs, oldB, newB) }

// Trim trims the specified characters from the beginning and end of the HBytes.
func (hbs HBytes) Trim(cutset HString) HBytes { return bytes.Trim(hbs, cutset.String()) }

// TrimLeft trims the specified characters from the beginning of the HBytes.
func (hbs HBytes) TrimLeft(cutset HString) HBytes { return bytes.TrimLeft(hbs, cutset.String()) }

// TrimRight trims the specified characters from the end of the HBytes.
func (hbs HBytes) TrimRight(cutset HString) HBytes { return bytes.TrimRight(hbs, cutset.String()) }

// TrimPrefix trims the specified HBytes prefix from the HBytes.
func (hbs HBytes) TrimPrefix(cutset HBytes) HBytes { return bytes.TrimPrefix(hbs, cutset) }

// TrimSuffix trims the specified HBytes suffix from the HBytes.
func (hbs HBytes) TrimSuffix(cutset HBytes) HBytes { return bytes.TrimSuffix(hbs, cutset) }

// Split splits the HBytes at each occurrence of the specified HBytes separator.
func (hbs HBytes) Split(sep ...HBytes) HSlice[HBytes] {
	var separator []byte
	if len(sep) != 0 {
		separator = sep[0]
	}

	return hSliceHBytesFromSlice(bytes.Split(hbs, separator))
}

func hSliceHBytesFromSlice(bb [][]byte) HSlice[HBytes] {
	result := NewHSlice[HBytes](0, len(bb))
	for _, v := range bb {
		result = result.Append(NewHBytes(v))
	}

	return result
}

// Add appends the given HBytes to the current HBytes.
func (hbs HBytes) Add(bs HBytes) HBytes { return append(hbs, bs...) }

// AddPrefix prepends the given HBytes to the current HBytes.
func (hbs HBytes) AddPrefix(bs HBytes) HBytes { return bs.Add(hbs) }

// Bytes returns the HBytes as a byte slice.
func (hbs HBytes) Bytes() []byte { return hbs }

// Clone creates a new HBytes instance with the same content as the current HBytes.
func (hbs HBytes) Clone() HBytes { return bytes.Clone(hbs) }

// Compare compares the HBytes with another HBytes and returns an HInt.
func (hbs HBytes) Compare(bs HBytes) HInt { return HInt(bytes.Compare(hbs, bs)) }

// Contains checks if the HBytes contains the specified HBytes.
func (hbs HBytes) Contains(bs HBytes) bool { return bytes.Contains(hbs, bs) }

// ContainsRegexp checks if the HBytes contains a match for the specified regular expression
// pattern.
func (hbs HBytes) ContainsRegexp(pattern *regexp.Regexp) bool { return pattern.Match(hbs) }

// ContainsAny checks if the HBytes contains any of the specified HBytes.
func (hbs HBytes) ContainsAny(bss ...HBytes) bool {
	for _, bs := range bss {
		if hbs.Contains(bs) {
			return true
		}
	}

	return false
}

// ContainsAll checks if the HBytes contains all of the specified HBytes.
func (hbs HBytes) ContainsAll(bss ...HBytes) bool {
	for _, bs := range bss {
		if !hbs.Contains(bs) {
			return false
		}
	}

	return true
}

// ContainsAnyChart checks if the given HBytes contains any characters from the input HString.
func (hbs HBytes) ContainsAnyChars(chars HString) bool {
	return bytes.ContainsAny(hbs, chars.String())
}

// ContainsRune checks if the HBytes contains the specified rune.
func (hbs HBytes) ContainsRune(r rune) bool { return bytes.ContainsRune(hbs, r) }

// Count counts the number of occurrences of the specified HBytes in the HBytes.
func (hbs HBytes) Count(bs HBytes) int { return bytes.Count(hbs, bs) }

// Empty checks if the HBytes is empty.
func (hbs HBytes) Empty() bool { return hbs.Len() == 0 }

// Eq checks if the HBytes is equal to another HBytes.
func (hbs HBytes) Eq(bs HBytes) bool { return hbs.Compare(bs).Eq(0) }

// EqFold compares two HBytes slices case-insensitively.
func (hbs HBytes) EqFold(bs HBytes) bool { return bytes.EqualFold(hbs, bs) }

// Gt checks if the HBytes is greater than another HBytes.
func (hbs HBytes) Gt(bs HBytes) bool { return hbs.Compare(bs).Gt(0) }

// HString returns the HBytes as an HString.
func (hbs HBytes) HString() HString { return HString(hbs) }

// Index returns the index of the first instance of bs in hbs, or -1 if bs is not present in hbs.
func (hbs HBytes) Index(bs HBytes) int { return bytes.Index(hbs, bs) }

// IndexByte returns the index of the first instance of the byte b in hbs, or -1 if b is not
// present in hbs.
func (hbs HBytes) IndexByte(b byte) int { return bytes.IndexByte(hbs, b) }

// IndexRune returns the index of the first instance of the rune r in hbs, or -1 if r is not
// present in hbs.
func (hbs HBytes) IndexRune(r rune) int { return bytes.IndexRune(hbs, r) }

// Len returns the length of the HBytes.
func (hbs HBytes) Len() int { return len(hbs) }

// LenRunes returns the number of runes in the HBytes.
func (hbs HBytes) LenRunes() int { return utf8.RuneCount(hbs) }

// Lt checks if the HBytes is less than another HBytes.
func (hbs HBytes) Lt(bs HBytes) bool { return hbs.Compare(bs).Lt(0) }

// Map applies a function to each rune in the HBytes and returns the modified HBytes.
func (hbs HBytes) Map(fn func(rune) rune) HBytes { return bytes.Map(fn, hbs) }

// NormalizeNFC returns a new HBytes with its Unicode characters normalized using the NFC form.
func (hbs HBytes) NormalizeNFC() HBytes { return norm.NFC.Bytes(hbs) }

// Ne checks if the HBytes is not equal to another HBytes.
func (hbs HBytes) Ne(bs HBytes) bool { return !hbs.Eq(bs) }

// NotEmpty checks if the HBytes is not empty.
func (hbs HBytes) NotEmpty() bool { return hbs.Len() != 0 }

// Reader returns a *bytes.Reader initialized with the content of HBytes.
func (hbs HBytes) Reader() *bytes.Reader { return bytes.NewReader(hbs) }

// Repeat returns a new HBytes consisting of the current HBytes repeated 'count' times.
func (hbs HBytes) Repeat(count int) HBytes { return bytes.Repeat(hbs, count) }

// Runes returns the HBytes as a slice of runes.
func (hbs HBytes) Runes() []rune { return bytes.Runes(hbs) }

// String returns the HBytes as a string.
func (hbs HBytes) String() string { return string(hbs) }

// ToTitle converts the HBytes to title case.
func (hbs HBytes) ToTitle() HBytes { return cases.Title(language.English).Bytes(hbs) }

// ToLower converts the HBytes to lowercase.
func (hbs HBytes) ToLower() HBytes { return cases.Lower(language.English).Bytes(hbs) }

// ToUpper converts the HBytes to uppercase.
func (hbs HBytes) ToUpper() HBytes { return cases.Upper(language.English).Bytes(hbs) }

// TrimSpace trims white space characters from the beginning and end of the HBytes.
func (hbs HBytes) TrimSpace() HBytes { return bytes.TrimSpace(hbs) }

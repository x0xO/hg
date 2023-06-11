package hg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/x0xO/hg/pkg/iter"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// NewHString creates a new HString from the provided string (optional).
func NewHString(strs ...string) HString {
	if len(strs) != 0 {
		return HString(strs[0])
	}

	return ""
}

// Random generates a random HString with the specified length.
//
// This function uses a predefined set of characters (ASCII_LETTERS and DIGITS) and iterates
// 'count' times, appending a random character from the set to the result HString.
//
// Parameters:
//
// - count (int): The length of the random HString to be generated.
//
// Returns:
//
// - HString: A random HString with the specified length.
//
// Example usage:
//
//	randomString := hg.HString.Random(10)
//
// randomString contains a random HString with 10 characters.
func (HString) Random(count int) HString {
	letters := HString(ASCII_LETTERS + DIGITS).Split()

	var result HString

	for range iter.N(count) {
		result += letters.Random()
	}

	return result
}

// IsASCII checks if all characters in the HString are ASCII bytes.
func (hs HString) IsASCII() bool {
	for _, r := range hs {
		if r > unicode.MaxASCII {
			return false
		}
	}

	return true
}

// IsDigit checks if all characters in the HString are digits.
func (hs HString) IsDigit() bool {
	if hs.Empty() {
		return false
	}

	for _, c := range hs {
		if !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}

// HInt tries to parse the HString as an int and returns an HInt.
func (hs HString) HInt() HInt {
	if hint, err := strconv.ParseInt(hs.String(), 0, 32); err == nil {
		return HInt(hint)
	}

	return 0
}

// HFloat tries to parse the HString as a float64 and returns an HFloat.
func (hs HString) HFloat() HFloat {
	if hfloat, err := strconv.ParseFloat(hs.String(), 64); err == nil {
		return HFloat(hfloat)
	}

	return 0
}

// ToTitle converts the HString to title case.
func (hs HString) ToTitle() HString {
	return HString(cases.Title(language.English).String(hs.String()))
}

// ToLower returns the HString in lowercase.
func (hs HString) ToLower() HString {
	return HString(cases.Lower(language.English).String(hs.String()))
}

// ToUpper returns the HString in uppercase.
func (hs HString) ToUpper() HString {
	return HString(cases.Upper(language.English).String(hs.String()))
}

// Trim trims characters in the cutset from the beginning and end of the HString.
func (hs HString) Trim(cutset HString) HString {
	return HString(strings.Trim(hs.String(), cutset.String()))
}

// TrimLeft trims characters in the cutset from the beginning of the HString.
func (hs HString) TrimLeft(cutset HString) HString {
	return HString(strings.TrimLeft(hs.String(), cutset.String()))
}

// TrimRight trims characters in the cutset from the end of the HString.
func (hs HString) TrimRight(cutset HString) HString {
	return HString(strings.TrimRight(hs.String(), cutset.String()))
}

// TrimPrefix trims the specified prefix from the HString.
func (hs HString) TrimPrefix(cutset HString) HString {
	return HString(strings.TrimPrefix(hs.String(), cutset.String()))
}

// TrimSuffix trims the specified suffix from the HString.
func (hs HString) TrimSuffix(cutset HString) HString {
	return HString(strings.TrimSuffix(hs.String(), cutset.String()))
}

// Replace replaces the 'oldS' HString with the 'newS' HString for the specified number of
// occurrences.
func (hs HString) Replace(oldS, newS HString, n int) HString {
	return HString(strings.Replace(hs.String(), oldS.String(), newS.String(), n))
}

// ReplaceAll replaces all occurrences of the 'oldS' HString with the 'newS' HString.
func (hs HString) ReplaceAll(oldS, newS HString) HString {
	return HString(strings.ReplaceAll(hs.String(), oldS.String(), newS.String()))
}

// ReplaceNth returns a new HString instance with the nth occurrence of oldS
// replaced with newS. If there aren't enough occurrences of oldS, the
// original HString is returned. If n is less than -1, the original HString
// is also returned. If n is -1, the last occurrence of oldS is replaced with newS.
//
// Returns:
//
// - A new HString instance with the nth occurrence of oldS replaced with newS.
//
// Example usage:
//
//	hs := hg.HString("The quick brown dog jumped over the lazy dog.")
//	result := hs.ReplaceNth("dog", "fox", 2)
//	fmt.Println(result)
//
// Output: "The quick brown dog jumped over the lazy fox.".
func (hs HString) ReplaceNth(oldS, newS HString, n int) HString {
	if n < -1 || oldS.Len() == 0 {
		return hs
	}

	count, i := 0, 0

	for {
		pos := hs[i:].Index(oldS)
		if pos == -1 {
			break
		}

		pos += i
		count++

		if count == n || (n == -1 && hs[pos+oldS.Len():].Index(oldS) == -1) {
			return hs[:pos].Add(newS).Add(hs[pos+oldS.Len():])
		}

		i = pos + oldS.Len()
	}

	return hs
}

// Contains checks if the HString contains the specified substring.
func (hs HString) Contains(substr HString) bool {
	return strings.Contains(hs.String(), substr.String())
}

// ContainsRegexp checks if the HString contains a match for the specified regular expression
// pattern.
func (hs HString) ContainsRegexp(pattern *regexp.Regexp) bool { return pattern.Match(hs.Bytes()) }

// ContainsAny checks if the HString contains any of the specified substrings.
func (hs HString) ContainsAny(substrs ...HString) bool {
	for _, substr := range substrs {
		if hs.Contains(substr) {
			return true
		}
	}

	return false
}

// ContainsAll checks if the given HString contains all the specified substrings.
func (hs HString) ContainsAll(substrs ...HString) bool {
	for _, substr := range substrs {
		if !hs.Contains(substr) {
			return false
		}
	}

	return true
}

// ContainsAnyChars checks if the HString contains any characters from the specified HString.
func (hs HString) ContainsAnyChars(chars HString) bool {
	return strings.ContainsAny(hs.String(), chars.String())
}

// StartsWith checks if the HString starts with any of the provided prefixes.
// The method accepts a variable number of arguments, allowing for checking against multiple
// prefixes at once. It iterates over the provided prefixes and uses the HasPrefix function from
// the strings package to check if
// the HString starts with each prefix.
// The function returns true if the HString starts with any of the prefixes, and false otherwise.
//
// Usage:
//
//	hs := hg.HString("http://example.com")
//	if hs.StartsWith("http://", "https://") {
//	   // do something
//	}
func (hs HString) StartsWith(prefixes ...HString) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(string(hs), prefix.String()) {
			return true
		}
	}

	return false
}

// EndsWith checks if the HString ends with any of the provided suffixes.
// The method accepts a variable number of arguments, allowing for checking against multiple
// suffixes at once. It iterates over the provided suffixes and uses the HasSuffix function from
// the strings package to check if
// the HString ends with each suffix.
// The function returns true if the HString ends with any of the suffixes, and false otherwise.
//
// Usage:
//
//	hs := hg.HString("example.com")
//	if hs.EndsWith(".com", ".net") {
//	   // do something
//	}
func (hs HString) EndsWith(suffixes ...HString) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(string(hs), suffix.String()) {
			return true
		}
	}

	return false
}

// Split splits the HString by the specified separator.
func (hs HString) Split(sep ...HString) HSlice[HString] {
	var separator string
	if len(sep) != 0 {
		separator = sep[0].String()
	}

	return hSliceHStringFromSlice(strings.Split(hs.String(), separator))
}

// Fields splits the HString into a slice of substrings, removing any whitespace.
func (hs HString) Fields() HSlice[HString] {
	return hSliceHStringFromSlice(strings.Fields(hs.String()))
}

func hSliceHStringFromSlice(ss []string) HSlice[HString] {
	result := NewHSlice[HString](0, len(ss))
	for _, v := range ss {
		result = result.Append(NewHString(v))
	}

	return result
}

// Chunks splits the HString into chunks of the specified size.
//
// This function iterates through the HString, creating new HString chunks of the specified size.
// If size is less than or equal to 0 or the HString is empty,
// it returns an empty HSlice[HString].
// If size is greater than or equal to the length of the HString,
// it returns an HSlice[HString] containing the original HString.
//
// Parameters:
//
// - size (int): The size of the chunks to split the HString into.
//
// Returns:
//
// - HSlice[HString]: A slice of HString chunks of the specified size.
//
// Example usage:
//
//	text := hg.HString("Hello, World!")
//	chunks := text.Chunks(4)
//
// chunks contains {"Hell", "o, W", "orld", "!"}.
func (hs HString) Chunks(size int) HSlice[HString] {
	if size <= 0 || hs.Empty() {
		return HSlice[HString]{}
	}

	if size >= hs.Len() {
		return HSlice[HString]{hs}
	}

	hr, hrLen := hs.Runes(), hs.LenRunes()
	chunks := NewHSlice[HString](0, (hrLen+size-1)/size)

	for i := 0; i < hrLen; i += size {
		end := i + size
		if end > hrLen {
			end = hrLen
		}

		chunks = chunks.Append(HString(hr[i:end]))
	}

	return chunks
}

// Cut returns a new HString that contains the text between the first
// occurrences of the 'start' and 'end' strings.
//
// The function searches for the 'start' and 'end' strings within the HString,
// and if both are found,
// it returns a new HString containing the text between the first occurrences
// of 'start' and 'end'.
// If either 'start' or 'end' is empty or not found in the HString,
// it returns the original HString.
//
// Parameters:
//
// - start (HString): The HString marking the beginning of the text to be cut.
//
// - end (HString): The HString marking the end of the text to be cut.
//
// Returns:
//
// - HString: A new HString containing the text between the first occurrences
// of 'start' and 'end', or the original HString if 'start' or 'end' is empty or not found.
//
// Example usage:
//
//	hs := hg.HString("Hello, [world]! How are you?")
//	cut := hs.Cut("[", "]") // "world"
func (hs HString) Cut(start, end HString) HString {
	if start.Empty() || end.Empty() {
		return hs
	}

	startIndex := hs.Index(start)
	if startIndex == -1 {
		return hs
	}

	endIndex := hs[startIndex+start.Len():].Index(end)
	if endIndex == -1 {
		return hs
	}

	return hs[startIndex+start.Len() : startIndex+start.Len()+endIndex]
}

// Similarity calculates the similarity between two HStrings using the
// Levenshtein distance algorithm and returns the similarity percentage as an HFloat.
//
// The function compares two HStrings using the Levenshtein distance,
// which measures the difference between two sequences by counting the number
// of single-character edits required to change one sequence into the other.
// The similarity is then calculated by normalizing the distance by the maximum
// length of the two input HStrings.
//
// Parameters:
//
// - hstr (HString): The HString to compare with hs.
//
// Returns:
//
// - HFloat: The similarity percentage between the two HStrings as a value between 0 and 100.
//
// Example usage:
//
//	hs1 := hg.HString("kitten")
//	hs2 := hg.HString("sitting")
//	similarity := hs1.Similarity(hs2) // 57.14285714285714
func (hs HString) Similarity(hstr HString) HFloat {
	if hs.Eq(hstr) {
		return 100
	}

	if hs.Len() == 0 || hstr.Len() == 0 {
		return 0
	}

	s1 := hs.Runes()
	s2 := hstr.Runes()

	lenS1 := hs.LenRunes()
	lenS2 := hstr.LenRunes()

	if lenS1 > lenS2 {
		s1, s2, lenS1, lenS2 = s2, s1, lenS2, lenS1
	}

	distance := NewHSlice[HInt](lenS1 + 1)

	for i, r2 := range s2 {
		prev := HInt(i).Add(1)

		for j, r1 := range s1 {
			current := distance[j]
			if r2 != r1 {
				current = distance[j].Add(1).Min(prev.Add(1)).Min(distance[j+1].Add(1))
			}

			distance[j], prev = prev, current
		}

		distance[lenS1] = prev
	}

	return HFloat(1).
		Sub(distance[lenS1].HFloat().Div(HInt(lenS1).Max(HInt(lenS2)).HFloat())).Mul(100)
}

// Compare compares two HStrings and returns an HInt indicating their relative order.
// The result will be 0 if hs==hstr, -1 if hs < hstr, and +1 if hs > hstr.
func (hs HString) Compare(hstr HString) HInt {
	return HInt(strings.Compare(hs.String(), hstr.String()))
}

// Add appends the specified HString to the current HString.
func (hs HString) Add(hstr HString) HString { return hs + hstr }

// AddPrefix prepends the specified HString to the current HString.
func (hs HString) AddPrefix(hstr HString) HString { return hstr.Add(hs) }

// Bytes returns the HString as a byte slice.
func (hs HString) Bytes() []byte { return []byte(hs) }

// ContainsRune checks if the HString contains the specified rune.
func (hs HString) ContainsRune(r rune) bool { return strings.ContainsRune(hs.String(), r) }

// Count returns the number of non-overlapping instances of the substring in the HString.
func (hs HString) Count(substr HString) int { return strings.Count(hs.String(), substr.String()) }

// Empty checks if the HString is empty.
func (hs HString) Empty() bool { return hs.Len() == 0 }

// Eq checks if two HStrings are equal.
func (hs HString) Eq(hstr HString) bool { return hs.Compare(hstr).Eq(0) }

// EqFold compares two HString strings case-insensitively.
func (hs HString) EqFold(hstr HString) bool {
	return strings.EqualFold(hs.String(), hstr.String())
}

// Gt checks if the HString is greater than the specified HString.
func (hs HString) Gt(hstr HString) bool { return hs.Compare(hstr).Gt(0) }

// HBytes returns the HString as an HBytes.
func (hs HString) HBytes() HBytes { return HBytes(hs) }

// Index returns the index of the first instance of the specified substring in the HString, or -1
// if substr is not present in hs.
func (hs HString) Index(substr HString) int { return strings.Index(hs.String(), substr.String()) }

// IndexRune returns the index of the first instance of the specified rune in the HString.
func (hs HString) IndexRune(r rune) int { return strings.IndexRune(hs.String(), r) }

// Len returns the length of the HString.
func (hs HString) Len() int { return len(hs) }

// LenRunes returns the number of runes in the HString.
func (hs HString) LenRunes() int { return utf8.RuneCountInString(hs.String()) }

// Lt checks if the HString is less than the specified HString.
func (hs HString) Lt(hstr HString) bool { return hs.Compare(hstr).Lt(0) }

// Map applies the provided function to all runes in the HString and returns the resulting HString.
func (hs HString) Map(fn func(rune) rune) HString { return HString(strings.Map(fn, hs.String())) }

// NormalizeNFC returns a new HString with its Unicode characters normalized using the NFC form.
func (hs HString) NormalizeNFC() HString { return HString(norm.NFC.String(hs.String())) }

// Ne checks if two HStrings are not equal.
func (hs HString) Ne(hstr HString) bool { return !hs.Eq(hstr) }

// NotEmpty checks if the HString is not empty.
func (hs HString) NotEmpty() bool { return hs.Len() != 0 }

// Reader returns a *strings.Reader initialized with the content of HString.
func (hs HString) Reader() *strings.Reader { return strings.NewReader(hs.String()) }

// Repeat returns a new HString consisting of the specified count of the original HString.
func (hs HString) Repeat(count int) HString { return HString(strings.Repeat(hs.String(), count)) }

// Reverse reverses the HString.
func (hs HString) Reverse() HString { return hs.HBytes().Reverse().HString() }

// Runes returns the HString as a slice of runes.
func (hs HString) Runes() []rune { return []rune(hs) }

// String returns the HString as a string.
func (hs HString) String() string { return string(hs) }

// TrimSpace trims whitespace from the beginning and end of the HString.
func (hs HString) TrimSpace() HString { return HString(strings.TrimSpace(hs.String())) }

// Format applies a specified format to the HString object.
func (hs HString) Format(format HString) HString {
	return HString(fmt.Sprintf(format.String(), hs))
}

// LeftJustify justifies the HString to the left by adding padding to the right, up to the
// specified length. If the length of the HString is already greater than or equal to the specified
// length, or the pad is empty, the original HString is returned.
//
// The padding HString is repeated as necessary to fill the remaining length.
// The padding is added to the right of the HString.
//
// Parameters:
//   - length: The desired length of the resulting justified HString.
//   - pad: The HString used as padding.
//
// Example usage:
//
//	hs := hg.HString("Hello")
//	result := hs.LeftJustify(10, "...")
//	// result: "Hello....."
func (hs HString) LeftJustify(length int, pad HString) HString {
	if hs.LenRunes() >= length || pad.Eq("") {
		return hs
	}

	var output strings.Builder

	output.WriteString(hs.String())
	writePadding(&output, pad, pad.LenRunes(), length-hs.LenRunes())

	return HString(output.String())
}

// RightJustify justifies the HString to the right by adding padding to the left, up to the
// specified length. If the length of the HString is already greater than or equal to the specified
// length, or the pad is empty, the original HString is returned.
//
// The padding HString is repeated as necessary to fill the remaining length.
// The padding is added to the left of the HString.
//
// Parameters:
//   - length: The desired length of the resulting justified HString.
//   - pad: The HString used as padding.
//
// Example usage:
//
//	hs := hg.HString("Hello")
//	result := hs.RightJustify(9, "...")
//	// result: "....Hello"
func (hs HString) RightJustify(length int, pad HString) HString {
	if hs.LenRunes() >= length || pad.Eq("") {
		return hs
	}

	var output strings.Builder

	writePadding(&output, pad, pad.LenRunes(), length-hs.LenRunes())
	output.WriteString(hs.String())

	return HString(output.String())
}

// Center justifies the HString by adding padding on both sides, up to the specified length.
// If the length of the HString is already greater than or equal to the specified length, or the
// pad is empty, the original HString is returned.
//
// The padding HString is repeated as necessary to evenly distribute the remaining length on both
// sides.
// The padding is added to the left and right of the HString.
//
// Parameters:
//   - length: The desired length of the resulting justified HString.
//   - pad: The HString used as padding.
//
// Example usage:
//
//	hs := hg.HString("Hello")
//	result := hs.Center(10, "...")
//	// result: "..Hello..."
func (hs HString) Center(length int, pad HString) HString {
	if hs.LenRunes() >= length || pad.Eq("") {
		return hs
	}

	var output strings.Builder

	remains := length - hs.LenRunes()
	writePadding(&output, pad, pad.LenRunes(), remains/2)
	output.WriteString(hs.String())
	writePadding(&output, pad, pad.LenRunes(), (remains+1)/2)

	return HString(output.String())
}

// writePadding writes the padding HString to the output Builder to fill the remaining length.
// It repeats the padding HString as necessary and appends any remaining runes from the padding
// HString.
func writePadding(output *strings.Builder, pad HString, padlen, remains int) {
	if repeats := remains / padlen; repeats > 0 {
		output.WriteString(pad.Repeat(repeats).String())
	}

	padrunes := pad.Runes()
	for i := range iter.N(remains % padlen) {
		output.WriteRune(padrunes[i])
	}
}

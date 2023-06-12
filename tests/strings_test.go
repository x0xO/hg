package hg_test

import (
	"io"
	"testing"

	"github.com/x0xO/hg"
	"github.com/x0xO/hg/pkg/iter"
)

func TestHStringIsDigit(t *testing.T) {
	tests := []struct {
		name string
		hs   hg.HString
		want bool
	}{
		{"empty", hg.HString(""), false},
		{"one", hg.HString("1"), true},
		{"nine", hg.HString("99999"), true},
		{"non-digit", hg.HString("1111a"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hs.IsDigit(); got != tt.want {
				t.Errorf("HString.IsDigit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHStringToHInt(t *testing.T) {
	tests := []struct {
		name string
		hs   hg.HString
		want hg.HInt
	}{
		{
			name: "empty",
			hs:   hg.HString(""),
			want: 0,
		},
		{
			name: "one digit",
			hs:   hg.HString("1"),
			want: 1,
		},
		{
			name: "two digits",
			hs:   hg.HString("12"),
			want: 12,
		},
		{
			name: "one letter",
			hs:   hg.HString("a"),
			want: 0,
		},
		{
			name: "one digit and one letter",
			hs:   hg.HString("1a"),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hs.HInt(); got != tt.want {
				t.Errorf("HString.ToHInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHStringToTitle(t *testing.T) {
	tests := []struct {
		name string
		hs   hg.HString
		want hg.HString
	}{
		{"empty", "", ""},
		{"one word", "hello", "Hello"},
		{"two words", "hello world", "Hello World"},
		{"three words", "hello world, how are you?", "Hello World, How Are You?"},
		{"multiple hyphens", "foo-bar-baz", "Foo-Bar-Baz"},
		{"non-ascii letters", "こんにちは, 世界!", "こんにちは, 世界!"},
		{"all whitespace", "   \t\n   ", "   \t\n   "},
		{"numbers", "12345 67890", "12345 67890"},
		{"arabic", "مرحبا بالعالم", "مرحبا بالعالم"},
		{"chinese", "你好世界", "你好世界"},
		{"czech", "ahoj světe", "Ahoj Světe"},
		{"danish", "hej verden", "Hej Verden"},
		{"dutch", "hallo wereld", "Hallo Wereld"},
		{"french", "bonjour tout le monde", "Bonjour Tout Le Monde"},
		{"german", "hallo welt", "Hallo Welt"},
		{"hebrew", "שלום עולם", "שלום עולם"},
		{"hindi", "नमस्ते दुनिया", "नमस्ते दुनिया"},
		{"hungarian", "szia világ", "Szia Világ"},
		{"italian", "ciao mondo", "Ciao Mondo"},
		{"japanese", "こんにちは世界", "こんにちは世界"},
		{"korean", "안녕하세요 세상", "안녕하세요 세상"},
		{"norwegian", "hei verden", "Hei Verden"},
		{"polish", "witaj świecie", "Witaj Świecie"},
		{"portuguese", "olá mundo", "Olá Mundo"},
		{"russian", "привет мир", "Привет Мир"},
		{"spanish", "hola mundo", "Hola Mundo"},
		{"swedish", "hej världen", "Hej Världen"},
		{"turkish", "merhaba dünya", "Merhaba Dünya"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hs.ToTitle(); got.Ne(tt.want) {
				t.Errorf("HString.ToTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHStringAdd(t *testing.T) {
	tests := []struct {
		name string
		hs   hg.HString
		s    hg.HString
		want hg.HString
	}{
		{
			name: "empty",
			hs:   hg.HString(""),
			s:    hg.HString(""),
			want: hg.HString(""),
		},
		{
			name: "empty_hs",
			hs:   hg.HString(""),
			s:    hg.HString("test"),
			want: hg.HString("test"),
		},
		{
			name: "empty_s",
			hs:   hg.HString("test"),
			s:    hg.HString(""),
			want: hg.HString("test"),
		},
		{
			name: "not_empty",
			hs:   hg.HString("test"),
			s:    hg.HString("test"),
			want: hg.HString("testtest"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hs.Add(tt.s); got != tt.want {
				t.Errorf("HString.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHStringAddPrefix(t *testing.T) {
	tests := []struct {
		name string
		hs   hg.HString
		s    hg.HString
		want hg.HString
	}{
		{
			name: "empty",
			hs:   hg.HString(""),
			s:    hg.HString(""),
			want: hg.HString(""),
		},
		{
			name: "empty_hs",
			hs:   hg.HString(""),
			s:    hg.HString("test"),
			want: hg.HString("test"),
		},
		{
			name: "empty_s",
			hs:   hg.HString("test"),
			s:    hg.HString(""),
			want: hg.HString("test"),
		},
		{
			name: "not_empty",
			hs:   hg.HString("rest"),
			s:    hg.HString("test"),
			want: hg.HString("testrest"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hs.AddPrefix(tt.s); got != tt.want {
				t.Errorf("HString.AddPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHStringRandom(t *testing.T) {
	for i := range iter.N(100) {
		random := hg.NewHString().Random(i)

		if random.Len() != i {
			t.Errorf("Random string length %d is not equal to %d", random.Len(), i)
		}
	}
}

func TestHStringChunks(t *testing.T) {
	h := hg.HString("")
	chunks := h.Chunks(3)

	if chunks.Len() != 0 {
		t.Errorf("Expected empty slice, but got %v", chunks)
	}

	h = hg.HString("hello")
	chunks = h.Chunks(10)

	if chunks.Len() != 1 {
		t.Errorf("Expected 1 chunk, but got %v", chunks.Len())
	}

	if chunks[0] != h {
		t.Errorf("Expected chunk to be %v, but got %v", h, chunks.Get(0))
	}

	h = hg.HString("hello")
	chunks = h.Chunks(2)

	if chunks.Len() != 3 {
		t.Errorf("Expected 3 chunks, but got %v", chunks.Len())
	}

	expectedChunks := hg.HSlice[hg.HString]{"he", "ll", "o"}

	for i, c := range chunks {
		if c != expectedChunks.Get(i) {
			t.Errorf("Expected chunk %v to be %v, but got %v", i, expectedChunks.Get(i), c)
		}
	}

	h = hg.HString("hello world")
	chunks = h.Chunks(3)

	if chunks.Len() != 4 {
		t.Errorf("Expected 4 chunks, but got %v", chunks.Len())
	}

	expectedChunks = hg.HSlice[hg.HString]{"hel", "lo ", "wor", "ld"}

	for i, c := range chunks {
		if c != expectedChunks.Get(i) {
			t.Errorf("Expected chunk %v to be %v, but got %v", i, expectedChunks.Get(i), c)
		}
	}

	h = hg.HString("hello")
	chunks = h.Chunks(5)

	if chunks.Len() != 1 {
		t.Errorf("Expected 1 chunk, but got %v", chunks.Len())
	}

	if chunks.Get(0) != h {
		t.Errorf("Expected chunk to be %v, but got %v", h, chunks.Get(0))
	}

	h = hg.HString("hello")
	chunks = h.Chunks(-1)

	if chunks.Len() != 0 {
		t.Errorf("Expected empty slice, but got %v", chunks)
	}
}

func TestHStringCut(t *testing.T) {
	tests := []struct {
		name   string
		input  hg.HString
		start  hg.HString
		end    hg.HString
		output hg.HString
	}{
		{"Basic", "Hello [start]world[end]!", "[start]", "[end]", "world"},
		{"No start", "Hello world!", "[start]", "[end]", "Hello world!"},
		{"No end", "Hello [start]world!", "[start]", "[end]", "Hello [start]world!"},
		{"Start equals end", "Hello [tag]world[tag]!", "[tag]", "[tag]", "world"},
		{
			"Multiple instances",
			"A [start]first[end] B [start]second[end] C",
			"[start]",
			"[end]",
			"first",
		},
		{"Empty input", "", "[start]", "[end]", ""},
		{"Empty start and end", "Hello world!", "", "", "Hello world!"},
		{
			"Nested tags",
			"A [start]first [start]nested[end] value[end] B",
			"[start]",
			"[end]",
			"first [start]nested",
		},
		{
			"Overlapping tags",
			"A [start]first[end][start]second[end] B",
			"[start]",
			"[end]",
			"first",
		},
		{"Unicode characters", "Привет [начало]мир[конец]!", "[начало]", "[конец]", "мир"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.input.Cut(test.start, test.end)
			if result != test.output {
				t.Errorf("Expected '%s', got '%s'", test.output, result)
			}
		})
	}
}

func TestHStringEscape(t *testing.T) {
	tests := []struct {
		input    hg.HString
		expected hg.HString
	}{
		{
			input:    hg.HString("hello world"),
			expected: hg.HString("hello+world"),
		},
		{
			input:    hg.HString("a+b=c/d"),
			expected: hg.HString("a%2Bb%3Dc%2Fd"),
		},
		{
			input:    hg.HString("foo?bar=baz&abc=123"),
			expected: hg.HString("foo%3Fbar%3Dbaz%26abc%3D123"),
		},
		{
			input:    hg.HString(""),
			expected: hg.HString(""),
		},
	}

	for _, test := range tests {
		actual := test.input.Enc().URL()
		if actual != test.expected {
			t.Errorf("Escape(%s): expected %s, but got %s", test.input, test.expected, actual)
		}
	}
}

func TestHStringUnEscape(t *testing.T) {
	tests := []struct {
		input    hg.HString
		expected hg.HString
	}{
		{
			input:    hg.HString("hello+world"),
			expected: hg.HString("hello world"),
		},
		{
			input:    hg.HString("hello%20world"),
			expected: hg.HString("hello world"),
		},
		{
			input:    hg.HString("a%2Bb%3Dc%2Fd"),
			expected: hg.HString("a+b=c/d"),
		},
		{
			input:    hg.HString("foo%3Fbar%3Dbaz%26abc%3D123"),
			expected: hg.HString("foo?bar=baz&abc=123"),
		},
		{
			input:    hg.HString(""),
			expected: hg.HString(""),
		},
	}

	for _, test := range tests {
		actual := test.input.Dec().URL()
		if actual != test.expected {
			t.Errorf("UnEscape(%s): expected %s, but got %s", test.input, test.expected, actual)
		}
	}
}

func TestHStringCompare(t *testing.T) {
	testCases := []struct {
		hs1      hg.HString
		hs2      hg.HString
		expected hg.HInt
	}{
		{"apple", "banana", -1},
		{"banana", "apple", 1},
		{"banana", "banana", 0},
		{"apple", "Apple", 1},
		{"", "", 0},
	}

	for _, tc := range testCases {
		result := tc.hs1.Compare(tc.hs2)
		if !result.Eq(tc.expected) {
			t.Errorf("Compare(%q, %q): expected %d, got %d", tc.hs1, tc.hs2, tc.expected, result)
		}
	}
}

func TestHStringEq(t *testing.T) {
	testCases := []struct {
		hs1      hg.HString
		hs2      hg.HString
		expected bool
	}{
		{"apple", "banana", false},
		{"banana", "banana", true},
		{"Apple", "apple", false},
		{"", "", true},
	}

	for _, tc := range testCases {
		result := tc.hs1.Eq(tc.hs2)
		if result != tc.expected {
			t.Errorf("Eq(%q, %q): expected %t, got %t", tc.hs1, tc.hs2, tc.expected, result)
		}
	}
}

func TestHStringNe(t *testing.T) {
	testCases := []struct {
		hs1      hg.HString
		hs2      hg.HString
		expected bool
	}{
		{"apple", "banana", true},
		{"banana", "banana", false},
		{"Apple", "apple", true},
		{"", "", false},
	}

	for _, tc := range testCases {
		result := tc.hs1.Ne(tc.hs2)
		if result != tc.expected {
			t.Errorf("Ne(%q, %q): expected %t, got %t", tc.hs1, tc.hs2, tc.expected, result)
		}
	}
}

func TestHStringGt(t *testing.T) {
	testCases := []struct {
		hs1      hg.HString
		hs2      hg.HString
		expected bool
	}{
		{"apple", "banana", false},
		{"banana", "apple", true},
		{"Apple", "apple", false},
		{"banana", "banana", false},
		{"", "", false},
	}

	for _, tc := range testCases {
		result := tc.hs1.Gt(tc.hs2)
		if result != tc.expected {
			t.Errorf("Gt(%q, %q): expected %t, got %t", tc.hs1, tc.hs2, tc.expected, result)
		}
	}
}

func TestHStringLt(t *testing.T) {
	testCases := []struct {
		hs1      hg.HString
		hs2      hg.HString
		expected bool
	}{
		{"apple", "banana", true},
		{"banana", "apple", false},
		{"Apple", "apple", true},
		{"banana", "banana", false},
		{"", "", false},
	}

	for _, tc := range testCases {
		result := tc.hs1.Lt(tc.hs2)
		if result != tc.expected {
			t.Errorf("Lt(%q, %q): expected %t, got %t", tc.hs1, tc.hs2, tc.expected, result)
		}
	}
}

func TestHStringReverse(t *testing.T) {
	testCases := []struct {
		in      hg.HString
		wantOut hg.HString
	}{
		{in: "", wantOut: ""},
		{in: " ", wantOut: " "},
		{in: "a", wantOut: "a"},
		{in: "ab", wantOut: "ba"},
		{in: "abc", wantOut: "cba"},
		{in: "abcdefg", wantOut: "gfedcba"},
		{in: "ab丂d", wantOut: "d丂ba"},
		{in: "abåd", wantOut: "dåba"},

		{in: "世界", wantOut: "界世"},
		{in: "🙂🙃", wantOut: "🙃🙂"},
		{in: "こんにちは", wantOut: "はちにんこ"},

		// Punctuation and whitespace
		{in: "Hello, world!", wantOut: "!dlrow ,olleH"},
		{in: "Hello\tworld!", wantOut: "!dlrow\tolleH"},
		{in: "Hello\nworld!", wantOut: "!dlrow\nolleH"},

		// Mixed languages and scripts
		{in: "Hello, 世界!", wantOut: "!界世 ,olleH"},
		{in: "Привет, мир!", wantOut: "!рим ,тевирП"},
		{in: "안녕하세요, 세계!", wantOut: "!계세 ,요세하녕안"},

		// Palindromes
		{in: "racecar", wantOut: "racecar"},
		{in: "A man, a plan, a canal: Panama", wantOut: "amanaP :lanac a ,nalp a ,nam A"},

		{
			in:      "The quick brown fox jumps over the lazy dog.",
			wantOut: ".god yzal eht revo spmuj xof nworb kciuq ehT",
		},
		{in: "A man a plan a canal panama", wantOut: "amanap lanac a nalp a nam A"},
		{in: "Was it a car or a cat I saw?", wantOut: "?was I tac a ro rac a ti saW"},
		{in: "Never odd or even", wantOut: "neve ro ddo reveN"},
		{in: "Do geese see God?", wantOut: "?doG ees eseeg oD"},
		{in: "A Santa at NASA", wantOut: "ASAN ta atnaS A"},
		{in: "Yo, Banana Boy!", wantOut: "!yoB ananaB ,oY"},
		{in: "Madam, in Eden I'm Adam", wantOut: "madA m'I nedE ni ,madaM"},
		{in: "Never odd or even", wantOut: "neve ro ddo reveN"},
		{in: "Was it a car or a cat I saw?", wantOut: "?was I tac a ro rac a ti saW"},
		{in: "Do geese see God?", wantOut: "?doG ees eseeg oD"},
		{in: "No 'x' in Nixon", wantOut: "noxiN ni 'x' oN"},
		{in: "A Santa at NASA", wantOut: "ASAN ta atnaS A"},
		{in: "Yo, Banana Boy!", wantOut: "!yoB ananaB ,oY"},
	}

	for _, tc := range testCases {
		result := tc.in.Reverse()
		if result.Ne(tc.wantOut) {
			t.Errorf("Reverse(%s): expected %s, got %s", tc.in, result, tc.wantOut)
		}
	}
}

func TestHStringNormalizeNFC(t *testing.T) {
	testCases := []struct {
		input    hg.HString
		expected hg.HString
	}{
		{input: "Mëtàl Hëàd", expected: "Mëtàl Hëàd"},
		{input: "Café", expected: "Café"},
		{input: "Ĵūņě", expected: "Ĵūņě"},
		{input: "𝓽𝓮𝓼𝓽 𝓬𝓪𝓼𝓮", expected: "𝓽𝓮𝓼𝓽 𝓬𝓪𝓼𝓮"},
		{input: "ḀϊṍẀṙṧ", expected: "ḀϊṍẀṙṧ"},
		{input: "えもじ れんしゅう", expected: "えもじ れんしゅう"},
		{input: "Научные исследования", expected: "Научные исследования"},
		{input: "🌟Unicode✨", expected: "🌟Unicode✨"},
		{input: "A\u0308", expected: "Ä"},
		{input: "o\u0308", expected: "ö"},
		{input: "u\u0308", expected: "ü"},
		{input: "O\u0308", expected: "Ö"},
		{input: "U\u0308", expected: "Ü"},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			output := tc.input.NormalizeNFC()
			if output != tc.expected {
				t.Errorf("Normalize(%q) = %q; want %q", tc.input, output, tc.expected)
			}
		})
	}
}

func TestHStringSimilarity(t *testing.T) {
	testCases := []struct {
		str1     hg.HString
		str2     hg.HString
		expected hg.HFloat
	}{
		{"hello", "hello", 100},
		{"hello", "world", 20},
		{"hello", "", 0},
		{"", "", 100},
		{"cat", "cats", 75},
		{"kitten", "sitting", 57.14},
		{"good", "bad", 25},
		{"book", "back", 50},
		{"abcdef", "azced", 50},
		{"tree", "three", 80},
		{"house", "horse", 80},
		{"language", "languish", 62.50},
		{"programming", "programmer", 72.73},
		{"algorithm", "logarithm", 77.78},
		{"software", "hardware", 50},
		{"tea", "ate", 33.33},
		{"pencil", "pen", 50},
		{"information", "informant", 63.64},
		{"coffee", "toffee", 83.33},
		{"developer", "develop", 77.78},
		{"distance", "difference", 50},
		{"similar", "similarity", 70},
		{"apple", "apples", 83.33},
		{"internet", "internets", 88.89},
		{"education", "dedication", 80},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			output := tc.str1.Similarity(tc.str2)
			if output.RoundDecimal(2).Ne(tc.expected) {
				t.Errorf(
					"hg.HString(\"%s\").SimilarText(\"%s\") = %.2f%% but want %.2f%%\n",
					tc.str1,
					tc.str2,
					output,
					tc.expected,
				)
			}
		})
	}
}

func TestHStringReader(t *testing.T) {
	tests := []struct {
		name     string
		hstring  hg.HString
		expected string
	}{
		{"Empty HString", "", ""},
		{"Single character HString", "a", "a"},
		{"Multiple characters HString", "hello world", "hello world"},
		{"HString with special characters", "こんにちは、世界！", "こんにちは、世界！"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := test.hstring.Reader()
			resultBytes, err := io.ReadAll(reader)
			if err != nil {
				t.Fatalf("Error reading from *strings.Reader: %v", err)
			}

			result := string(resultBytes)

			if result != test.expected {
				t.Errorf("Reader() content = %s, expected %s", result, test.expected)
			}
		})
	}
}

func TestHStringContainsAny(t *testing.T) {
	testCases := []struct {
		name    string
		input   hg.HString
		substrs hg.HSlice[hg.HString]
		want    bool
	}{
		{
			name:    "ContainsAny_OneSubstringMatch",
			input:   "This is an example",
			substrs: []hg.HString{"This", "missing"},
			want:    true,
		},
		{
			name:    "ContainsAny_NoSubstringMatch",
			input:   "This is an example",
			substrs: []hg.HString{"notfound", "missing"},
			want:    false,
		},
		{
			name:    "ContainsAny_EmptySubstrings",
			input:   "This is an example",
			substrs: []hg.HString{},
			want:    false,
		},
		{
			name:    "ContainsAny_EmptyInput",
			input:   "",
			substrs: []hg.HString{"notfound", "missing"},
			want:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.ContainsAny(tc.substrs...)
			if got != tc.want {
				t.Errorf("ContainsAny() = %v; want %v", got, tc.want)
			}
		})
	}
}

func TestHStringContainsAll(t *testing.T) {
	testCases := []struct {
		name    string
		input   hg.HString
		substrs hg.HSlice[hg.HString]
		want    bool
	}{
		{
			name:    "ContainsAll_AllSubstringsMatch",
			input:   "This is an example",
			substrs: []hg.HString{"This", "example"},
			want:    true,
		},
		{
			name:    "ContainsAll_NotAllSubstringsMatch",
			input:   "This is an example",
			substrs: []hg.HString{"This", "missing"},
			want:    false,
		},
		{
			name:    "ContainsAll_EmptySubstrings",
			input:   "This is an example",
			substrs: []hg.HString{},
			want:    true,
		},
		{
			name:    "ContainsAll_EmptyInput",
			input:   "",
			substrs: []hg.HString{"notfound", "missing"},
			want:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.ContainsAll(tc.substrs...)
			if got != tc.want {
				t.Errorf("ContainsAll() = %v; want %v", got, tc.want)
			}
		})
	}
}

func TestHStringReplaceNth(t *testing.T) {
	tests := []struct {
		name     string
		hs       hg.HString
		oldS     hg.HString
		newS     hg.HString
		n        int
		expected hg.HString
	}{
		{
			"First occurrence",
			"The quick brown dog jumped over the lazy dog.",
			"dog",
			"fox",
			1,
			"The quick brown fox jumped over the lazy dog.",
		},
		{
			"Second occurrence",
			"The quick brown dog jumped over the lazy dog.",
			"dog",
			"fox",
			2,
			"The quick brown dog jumped over the lazy fox.",
		},
		{
			"Last occurrence",
			"The quick brown dog jumped over the lazy dog.",
			"dog",
			"fox",
			-1,
			"The quick brown dog jumped over the lazy fox.",
		},
		{
			"Negative n (except -1)",
			"The quick brown dog jumped over the lazy dog.",
			"dog",
			"fox",
			-2,
			"The quick brown dog jumped over the lazy dog.",
		},
		{
			"Zero n",
			"The quick brown dog jumped over the lazy dog.",
			"dog",
			"fox",
			0,
			"The quick brown dog jumped over the lazy dog.",
		},
		{
			"Longer replacement",
			"Hello, world!",
			"world",
			"beautiful world",
			1,
			"Hello, beautiful world!",
		},
		{
			"Shorter replacement",
			"A wonderful day, isn't it?",
			"wonderful",
			"nice",
			1,
			"A nice day, isn't it?",
		},
		{
			"Replace entire string",
			"Hello, world!",
			"Hello, world!",
			"Greetings, world!",
			1,
			"Greetings, world!",
		},
		{"No replacement", "Hello, world!", "x", "y", 1, "Hello, world!"},
		{"Nonexistent substring", "Hello, world!", "foobar", "test", 1, "Hello, world!"},
		{"Replace empty string", "Hello, world!", "", "x", 1, "Hello, world!"},
		{"Multiple identical substrings", "banana", "na", "xy", 1, "baxyna"},
		{"Multiple identical substrings, last", "banana", "na", "xy", -1, "banaxy"},
		{"Replace with empty string", "Hello, world!", "world", "", 1, "Hello, !"},
		{"Empty input string", "", "world", "test", 1, ""},
		{"Empty input, empty oldS, empty newS", "", "", "", 1, ""},
		{"Replace multiple spaces", "Hello    world!", "    ", " ", 1, "Hello world!"},
		{"Unicode characters", "こんにちは世界！", "世界", "World", 1, "こんにちはWorld！"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.hs.ReplaceNth(test.oldS, test.newS, test.n)
			if result != test.expected {
				t.Errorf("ReplaceNth() got %q, want %q", result, test.expected)
			}
		})
	}
}

func TestHStringIsASCII(t *testing.T) {
	testCases := []struct {
		input    hg.HString
		expected bool
	}{
		{"Hello, world!", true},
		{"こんにちは", false},
		{"", true},
		{"1234567890", true},
		{"ABCabc", true},
		{"~`!@#$%^&*()-_+={[}]|\\:;\"'<,>.?/", true},
		{"áéíóú", false},
		{"Привет", false},
	}

	for _, tc := range testCases {
		result := tc.input.IsASCII()
		if result != tc.expected {
			t.Errorf("IsASCII(%q) returned %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

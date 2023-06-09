package hg

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html"
	"io"
	"log"
	"net/url"
	"strconv"

	"github.com/x0xO/hg/pkg/iter"
)

type (
	// A struct that wraps an HString for encoding.
	enc struct{ str HString }
	// A struct that wraps an HString for decoding.
	dec struct{ str HString }
)

// Encode returns an enc struct wrapping the given HString.
func (hs HString) Encode() enc { return enc{hs} }

// Decode returns a dec struct wrapping the given HString.
func (hs HString) Decode() dec { return dec{hs} }

// GzFlate compresses the wrapped HString using GzDeflate and returns the compressed data as a
// Base64-encoded HString.
func (e enc) GzFlate() HString {
	// GzDeflate
	buffer := new(bytes.Buffer)

	writer, err := flate.NewWriter(buffer, 7)
	if err != nil {
		log.Printf("gzdeflate error: %s\n", err)
		return ""
	}

	_, _ = writer.Write(e.str.Bytes())
	_ = writer.Flush()
	_ = writer.Close()

	return HString(buffer.String()).Encode().Base64()
}

// GzFlate decompresses the Base64-encoded wrapped HString using GzInflate and returns the
// decompressed data as an HString.
func (d dec) GzFlate() HString {
	// GzInflate
	decoded := d.str.Decode().Base64()
	if decoded == "" {
		return decoded
	}

	reader := flate.NewReader(decoded.Reader())
	buffer := new(bytes.Buffer)
	_, _ = io.Copy(buffer, reader)
	_ = reader.Close()

	return HString(buffer.String())
}

// Base64 encodes the wrapped HString using Base64 and returns the encoded result as an HString.
func (e enc) Base64() HString {
	return HString(base64.StdEncoding.EncodeToString(e.str.Bytes()))
}

// Base64 decodes the wrapped HString using Base64 and returns the decoded result as an HString.
func (d dec) Base64() HString {
	decoded, err := base64.StdEncoding.DecodeString(d.str.String())
	if err != nil {
		log.Printf("base64decode error: %s\n", err)
		return ""
	}

	return HString(decoded)
}

// URL URL-encodes the wrapped HString and returns the encoded result as an HString.
func (e enc) URL() HString { return HString(url.QueryEscape(e.str.String())) }

// URL URL-decodes the wrapped HString and returns the decoded result as an HString.
func (d dec) URL() HString {
	result, _ := url.QueryUnescape(d.str.String())
	return HString(result)
}

// HTML HTML-encodes the wrapped HString and returns the encoded result as an HString.
func (e enc) HTML() HString { return HString(html.EscapeString(e.str.String())) }

// HTML HTML-decodes the wrapped HString and returns the decoded result as an HString.
func (d dec) HTML() HString { return HString(html.UnescapeString(d.str.String())) }

// Rot13 encodes the wrapped HString using ROT13 cipher and returns the encoded result as an
// HString.
func (e enc) Rot13() HString {
	rot := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+13)%26
		default:
			return r
		}
	}

	return e.str.Map(rot)
}

// Rot13 decodes the wrapped HString using ROT13 cipher and returns the decoded result as an
// HString.
func (d dec) Rot13() HString { return d.str.Encode().Rot13() }

// XOR encodes the wrapped HString using XOR cipher with the given key and returns the encoded
// result as an HString.
func (e enc) XOR(key HString) HString {
	if key.Empty() {
		return e.str
	}

	encrypted := e.str.Bytes()

	for i := range iter.N(e.str.Len()) {
		encrypted[i] ^= key[i%key.Len()]
	}

	return HString(encrypted)
}

// XOR decodes the wrapped HString using XOR cipher with the given key and returns the decoded
// result as an HString.
func (d dec) XOR(key HString) HString { return d.str.Encode().XOR(key) }

// Hex hex-encodes the wrapped HString and returns the encoded result as an HString.
func (e enc) Hex() HString { return HString(hex.EncodeToString(e.str.Bytes())) }

// Hex hex-decodes the wrapped HString and returns the decoded result as an HString.
func (d dec) Hex() HString {
	result, _ := hex.DecodeString(d.str.String())
	return HString(result)
}

// Binary converts the wrapped HString to its binary representation as an HString.
func (e enc) Binary() HString {
	var buf bytes.Buffer
	for i := range iter.N(e.str.Len()) {
		fmt.Fprintf(&buf, "%08b", e.str[i])
	}

	return HString(buf.String())
}

// Binary converts the wrapped binary HString back to its original HString representation.
func (d dec) Binary() HString {
	var result HBytes

	for i := 0; i+8 <= d.str.Len(); i += 8 {
		b, _ := strconv.ParseUint(d.str[i:i+8].String(), 2, 8)
		result = append(result, byte(b))
	}

	return result.HString()
}

package surf

import (
	"bufio"
	"compress/zlib"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"math"
	"regexp"

	"github.com/x0xO/hg"
	"golang.org/x/net/html/charset"
)

type body struct {
	body        io.ReadCloser
	contentType string
	content     hg.HBytes
	limit       int64
	deflate     bool
	cache       bool
}

// MD5 returns the MD5 hash of the body's content as a HString.
func (b *body) MD5() hg.HString { return b.HString().Hash().MD5() }

// XML decodes the body's content as XML into the provided data structure.
func (b *body) XML(data any) error { return xml.Unmarshal(b.HBytes(), data) }

// JSON decodes the body's content as JSON into the provided data structure.
func (b *body) JSON(data any) error { return json.Unmarshal(b.HBytes(), data) }

// Stream returns the body's bufio.Reader for streaming the content.
func (b *body) Stream() *bufio.Reader { return bufio.NewReader(b.body) }

// HString returns the body's content as a HString.
func (b *body) HString() hg.HString { return b.HBytes().HString() }

// String returns the body's content as a string.
func (b *body) String() string { return b.HBytes().String() }

// Limit sets the body's size limit and returns the modified body.
func (b *body) Limit(limit int64) *body { b.limit = limit; return b }

// Close closes the body and returns any error encountered.
func (b *body) Close() error {
	if b.body == nil {
		return errors.New("empty body error")
	}

	if _, err := io.Copy(io.Discard, b.body); err != nil {
		return err
	}

	return b.body.Close()
}

// UTF8 converts the body's content to UTF-8 encoding and returns it as a string.
func (b *body) UTF8() hg.HString {
	reader, err := charset.NewReader(b.HBytes().Reader(), b.contentType)
	if err != nil {
		return b.HString()
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return b.HString()
	}

	return hg.HString(content)
}

// Bytes returns the body's content as a byte slice.
func (b *body) HBytes() hg.HBytes {
	if b.cache && b.content != nil {
		return b.content
	}

	if _, err := b.body.Read(nil); err != nil {
		if err.Error() == "http: read on closed response body" {
			return nil
		}
	}

	defer b.Close()

	var err error
	if b.deflate {
		if b.body, err = zlib.NewReader(b.body); err != nil {
			return nil
		}
	}

	if b.limit == -1 {
		b.limit = math.MaxInt64
	}

	var content hg.HBytes

	content, err = io.ReadAll(io.LimitReader(b.body, b.limit))
	if err != nil {
		return nil
	}

	if b.cache {
		b.content = content
	}

	return content
}

// Dump dumps the body's content to a file with the given filename.
func (b *body) Dump(filename string) error {
	defer b.Close()
	f := hg.NewHFile(hg.HString(filename)).WriteFromReader(b.body)

	return f.Error()
}

// Contains checks if the body's content contains the provided pattern (byte slice, string, or
// *regexp.Regexp) and returns a boolean.
func (b *body) Contains(pattern any) bool {
	switch p := pattern.(type) {
	case []byte:
		return b.HBytes().ToLower().Contains(hg.HBytes(p).ToLower())
	case string:
		return b.HString().ToLower().Contains(hg.HString(p).ToLower())
	case *regexp.Regexp:
		return b.HBytes().ContainsRegexp(p)
	}

	return false
}

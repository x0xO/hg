package surf

import (
	"net/http"
	"net/textproto"
	"regexp"
	"strings"
)

type headers http.Header

// Contains checks if the header contains any of the specified patterns.
// It accepts a header name and a pattern (or list of patterns) and returns a boolean value
// indicating whether any of the patterns are found in the header values.
// The patterns can be a string, a slice of strings, or a slice of *regexp.Regexp.
func (h headers) Contains(header string, patterns any) bool {
	if h.Values(header) != nil {
		for _, value := range h.Values(header) {
			switch ps := patterns.(type) {
			case string:
				if strings.Contains(strings.ToLower(value), strings.ToLower(ps)) {
					return true
				}
			case []string:
				for _, pattern := range ps {
					if strings.Contains(strings.ToLower(value), strings.ToLower(pattern)) {
						return true
					}
				}
			case []*regexp.Regexp:
				for _, pattern := range ps {
					if pattern.Match([]byte(value)) {
						return true
					}
				}
			}
		}
	}

	return false
}

// Values returns the values associated with a specified header key.
// It wraps the Values method from the textproto.MIMEHeader type.
func (h headers) Values(key string) []string {
	return textproto.MIMEHeader(h).Values(key)
}

// Get returns the first value associated with a specified header key.
// It wraps the Get method from the textproto.MIMEHeader type.
func (h headers) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}

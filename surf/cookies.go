package surf

import (
	"net/http"
	"regexp"

	"github.com/x0xO/hg"
)

type cookies []*http.Cookie

// Contains checks if the cookies collection contains a cookie that matches the provided pattern.
// The pattern parameter can be either a string or a pointer to a regexp.Regexp object.
// The method returns true if a matching cookie is found and false otherwise.
func (c *cookies) Contains(pattern any) bool {
	for _, cookie := range *c {
		switch p := pattern.(type) {
		case string:
			return hg.HString(cookie.String()).ToLower().Contains(hg.HString(p).ToLower())
		case *regexp.Regexp:
			return hg.HBytes(cookie.String()).ContainsRegexp(p)
		}
	}

	return false
}

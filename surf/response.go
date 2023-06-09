package surf

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

type Response struct {
	*Client
	remoteAddr    net.Addr
	URL           *url.URL
	response      *http.Response
	Body          *body
	request       *Request
	Headers       headers
	Status        string
	UserAgent     string
	Proto         string
	History       history
	Cookies       cookies
	Time          time.Duration
	ContentLength int64
	StatusCode    int
	Attempts      int
}

func (resp Response) GetResponse() *http.Response { return resp.response }

// Referer returns the referer of the response.
func (resp Response) Referer() string { return resp.response.Request.Referer() }

// GetCookies returns the cookies from the response for the given URL.
func (resp Response) GetCookies(rawURL string) []*http.Cookie { return resp.getCookies(rawURL) }

// RemoteAddress returns the remote address of the response.
func (resp Response) RemoteAddress() net.Addr { return resp.remoteAddr }

// SetCookies sets cookies for the given URL in the response.
func (resp *Response) SetCookies(rawURL string, cookies []*http.Cookie) error {
	return resp.setCookies(rawURL, cookies)
}

// TLSGrabber returns a tlsData struct containing information about the TLS connection if it
// exists.
func (resp Response) TLSGrabber() *tlsData {
	if resp.response.TLS != nil {
		return tlsGrabber(resp.response.TLS)
	}

	return nil
}

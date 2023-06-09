package surf

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// Options a struct that holds options for configuring the HTTP client.
type Options struct {
	proxy         any
	dialer        *net.Dialer
	dnsCacheStats *cacheDialerStats
	cliMW         []clientMiddleware
	reqMW         []requestMiddleware
	maxRedirects  int
	retryWait     time.Duration
	retryMax      int
	http2         bool
	history       bool
	cacheBody     bool
}

// NewOptions creates a new Options instance with default values.
func NewOptions() *Options { return new(Options) }

func (opt *Options) addreqMW(m requestMiddleware) *Options {
	opt.reqMW = append(opt.reqMW, m)
	return opt
}

func (opt *Options) addcliMW(m clientMiddleware) *Options {
	opt.cliMW = append(opt.cliMW, m)
	return opt
}

// UnixDomainSocket sets the path for a Unix domain socket in the Options.
// This allows the HTTP client to connect to the server using a Unix domain
// socket instead of a traditional TCP/IP connection.
func (opt *Options) UnixDomainSocket(socketPath string) *Options {
	return opt.addcliMW(func(client *Client) { unixDomainSocketMW(client, socketPath) })
}

// DNSCache configures the DNS cache settings of the HTTP client.
//
// DNS caching can improve the performance of HTTP clients by caching the DNS
// lookup results for the specified Time-To-Live (TTL) duration and limiting the usage of the
// cached DNS result.
//
// Parameters:
//
// - ttl: the TTL duration for the DNS cache. After this duration, the DNS cache for a host will be
// invalidated.
//
// - maxUsage: the maximum number of times a cached DNS lookup result can be used. After this
// number is reached, the DNS cache for a host will be invalidated.
//
// Returns the same Options object, allowing for chaining of configuration calls.
//
// Example:
//
//	opt := surf.NewOptions().DNSCache(time.Second*30, 10)
//	cli := surf.NewClient().SetOptions(opt)
//
// The client will now use a DNS cache with a 30-second TTL and a maximum usage count of 10.
func (opt *Options) DNSCache(ttl time.Duration, maxUsage int64) *Options {
	return opt.addcliMW(func(client *Client) { dnsCacheMW(client, ttl, maxUsage) })
}

// DNS sets the custom DNS resolver address.
func (opt *Options) DNS(dns string) *Options {
	return opt.addcliMW(func(client *Client) { dnsMW(client, dns) })
}

// DNSOverTLS configures the client to use DNS over TLS.
func (opt *Options) DNSOverTLS() *dnsOverTLS { return &dnsOverTLS{opt: opt} }

// Timeout sets the timeout duration for the client.
func (opt *Options) Timeout(timeout time.Duration) *Options {
	return opt.addcliMW(func(client *Client) { timeoutMW(client, timeout) })
}

// InterfaceAddr sets the network interface address for the client.
func (opt *Options) InterfaceAddr(address string) *Options {
	return opt.addcliMW(func(client *Client) { interfaceAddrMW(client, address) })
}

// Proxy sets the proxy settings for the client.
func (opt *Options) Proxy(proxy any) *Options {
	opt.proxy = proxy
	return opt.addcliMW(func(client *Client) { proxyMW(client, proxy) })
}

// BasicAuth sets the basic authentication credentials for the client.
func (opt *Options) BasicAuth(authentication any) *Options {
	return opt.addreqMW(func(req *Request) error { return basicAuthMW(req, authentication) })
}

// BearerAuth sets the bearer token for the client.
func (opt *Options) BearerAuth(authentication string) *Options {
	return opt.addreqMW(func(req *Request) error { return bearerAuthMW(req, authentication) })
}

// UserAgent sets the user agent for the client.
func (opt *Options) UserAgent(userAgent any) *Options {
	return opt.addreqMW(func(req *Request) error { return userAgentMW(req, userAgent) })
}

// ContentType sets the content type for the client.
func (opt *Options) ContentType(contentType string) *Options {
	return opt.addreqMW(func(req *Request) error { return contentTypeMW(req, contentType) })
}

// CacheBody configures whether the client should cache the body of the response.
func (opt *Options) CacheBody(enable ...bool) *Options {
	if len(enable) != 0 {
		opt.cacheBody = enable[0]
	} else {
		opt.cacheBody = true
	}

	return opt
}

// GetRemoteAddress configures whether the client should get the remote address.
func (opt *Options) GetRemoteAddress() *Options { return opt.addreqMW(remoteAddrMW) }

// DisableKeepAlive disable keep-alive connections.
func (opt *Options) DisableKeepAlive() *Options { return opt.addcliMW(disableKeepAliveMW) }

// Retry configures the retry behavior of the client.
func (opt *Options) Retry(retryMax int, retryWait ...time.Duration) *Options {
	opt.retryMax = retryMax

	if len(retryWait) != 0 {
		opt.retryWait = retryWait[0]
	} else {
		opt.retryWait = time.Second * 1
	}

	return opt
}

// History configures whether the client should keep a history of requests (for debugging purposes
// only).
// WARNING: use only for debugging, not in async mode, no concurrency safe!!!
func (opt *Options) History() *Options {
	opt.history = true
	return opt
}

// HTTP2 configures whether the client should use HTTP/2.
func (opt *Options) HTTP2(enable ...bool) *Options {
	force := true

	if len(enable) != 0 {
		force = enable[0]
	}

	opt.http2 = force

	return opt.addcliMW(http2MW)
}

// Session configures whether the client should maintain a session.
func (opt *Options) Session() *Options { return opt.addcliMW(sessionMW) }

// MaxRedirects sets the maximum number of redirects the client should follow.
func (opt *Options) MaxRedirects(maxRedirects int) *Options {
	opt.maxRedirects = maxRedirects
	return opt.addcliMW(func(client *Client) { redirectPolicyMW(client, false, false, nil) })
}

// FollowOnlyHostRedirects configures whether the client should only follow redirects within the
// same host.
func (opt *Options) FollowOnlyHostRedirects() *Options {
	return opt.addcliMW(func(client *Client) { redirectPolicyMW(client, true, false, nil) })
}

// ForwardHeadersOnRedirect adds a middleware to the Options object that ensures HTTP headers are
// forwarded during a redirect.
func (opt *Options) ForwardHeadersOnRedirect() *Options {
	return opt.addcliMW(func(client *Client) { redirectPolicyMW(client, false, true, nil) })
}

// RedirectPolicy sets a custom redirect policy for the client.
func (opt *Options) RedirectPolicy(f func(*http.Request, []*http.Request) error) *Options {
	return opt.addcliMW(func(client *Client) { redirectPolicyMW(client, false, false, f) })
}

// String generate a string representation of the Options instance.
func (opt Options) String() string { return fmt.Sprintf("%#v", opt) }

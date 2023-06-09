package surf

import (
	"context"
	"crypto/tls"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// default dialer for surf.
func surfDialerMW(client *Client) {
	client.dialer = &net.Dialer{Timeout: _dialerTimeout, KeepAlive: _TCPKeepAlive}
}

// default tlsConfig for surf.
func surfTLSConfigMW(client *Client) { client.tlsConfig = &tls.Config{InsecureSkipVerify: true} }

// default transport for surf.
func surfTransportMW(client *Client) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = client.dialer.DialContext
	transport.TLSClientConfig = client.tlsConfig
	transport.MaxIdleConns = _maxIdleConns
	transport.MaxConnsPerHost = _maxConnsPerHost
	transport.MaxIdleConnsPerHost = _maxIdleConnsPerHost
	transport.IdleConnTimeout = _idleConnTimeout

	client.transport = transport
}

// default client for surf.
func surfClientMW(client *Client) {
	client.cli = &http.Client{Transport: client.transport, Timeout: _clientTimeout}
}

// http2MW configures the client to use HTTP2 or HTTP1.1, depending on the provided options.
func http2MW(client *Client) {
	if client.opt.http2 {
		client.GetTransport().ForceAttemptHTTP2 = true
	} else {
		client.GetTransport().TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
	}
}

// sessionMW configures the client's cookie jar to enable session handling.
func sessionMW(client *Client) { client.GetClient().Jar, _ = cookiejar.New(nil) }

// disableKeepAliveMW disable keep-alive setting.
func disableKeepAliveMW(client *Client) { client.GetTransport().DisableKeepAlives = true }

// interfaceAddrMW configures the client's local address for dialing based on the provided
// options.
func interfaceAddrMW(client *Client, address string) error {
	if address != "" {
		ip, err := net.ResolveTCPAddr("tcp", address+":0")
		if err != nil {
			return err
		}

		client.GetDialer().LocalAddr = ip
	}

	return nil
}

// timeoutMW configures the client's timeout setting based on the provided options.
func timeoutMW(client *Client, timeout time.Duration) error {
	client.GetClient().Timeout = timeout
	return nil
}

// redirectPolicyMW configures the client's redirect policy based on the
// provided options.
func redirectPolicyMW(
	client *Client,
	followOnlyHostRedirects, forwardHeadersOnRedirect bool,
	f func(*http.Request, []*http.Request) error,
) {
	maxRedirects := _maxRedirects
	if client.opt != nil && client.opt.maxRedirects != 0 {
		maxRedirects = client.opt.maxRedirects
	}

	redirectPolicy := func(req *http.Request, via []*http.Request) error {
		if len(via) >= maxRedirects {
			return http.ErrUseLastResponse
		}

		if followOnlyHostRedirects {
			newHost := req.URL.Host
			oldHost := via[0].Host

			if oldHost == "" {
				oldHost = via[0].URL.Host
			}

			if newHost != oldHost {
				return http.ErrUseLastResponse
			}
		}

		if forwardHeadersOnRedirect {
			for key, val := range via[0].Header {
				req.Header[key] = val
			}
		}

		if client.opt != nil && client.opt.history {
			client.history = append(client.history, req.Response)
		}

		return nil
	}

	if f != nil {
		redirectPolicy = f
	}

	client.GetClient().CheckRedirect = redirectPolicy
}

// dnsMW sets the DNS for client.
func dnsMW(client *Client, dns string) {
	client.GetDialer().Resolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var dialer net.Dialer
			return dialer.DialContext(ctx, "udp", dns)
		},
	}
}

// dnsTLSMW sets up a DNS over TLS for client.
func dnsTLSMW(client *Client, resolver *net.Resolver) { client.GetDialer().Resolver = resolver }

// dnsCacheMW sets up a DNS cache for client.
func dnsCacheMW(client *Client, ttl time.Duration, maxUsage int64) {
	if ttl != 0 && maxUsage != 0 {
		client.cacheDialer(ttl, maxUsage)
	}
}

// configureUnixSocket sets the DialContext function for the client's HTTP transport to use
// a Unix domain socket if the unixDomainSocket option is set.
func unixDomainSocketMW(client *Client, unixDomainSocket string) {
	if unixDomainSocket == "" {
		return
	}

	client.GetTransport().DialContext = func(_ context.Context, _, addr string) (net.Conn, error) {
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}

		unixaddr, err := net.ResolveUnixAddr(host, unixDomainSocket)
		if err != nil {
			return nil, err
		}

		return net.DialUnix(host, nil, unixaddr)
	}
}

// proxyMW configures the request's proxy settings based on the provided
// proxy options. It supports single or multiple proxy options.
func proxyMW(client *Client, proxys any) {
	var proxy string

	switch pr := proxys.(type) {
	case string:
		proxy = pr
	case []string:
		proxy = pr[rand.Intn(len(pr))]
	}

	if proxy == "" {
		return
	}

	proxyURL, _ := url.Parse(proxy)
	client.GetTransport().Proxy = http.ProxyURL(proxyURL)
}

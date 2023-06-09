package surf

import "time"

const (
	_userAgent = "hg-http-client/6.6.6 (+https://github.com/x0xO/hg)"

	_maxRedirects = 10
	_maxWorkers   = 10

	_dialerTimeout   = 30 * time.Second
	_clientTimeout   = 30 * time.Second
	_TCPKeepAlive    = 15 * time.Second
	_idleConnTimeout = 20 * time.Second

	_maxIdleConns        = 512
	_maxConnsPerHost     = 128
	_maxIdleConnsPerHost = 128
)

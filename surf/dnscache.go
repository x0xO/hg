package surf

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/x0xO/hg"
)

// cacheDialerStats contains stats of a dialer.
type cacheDialerStats struct {
	totalConn          int64
	cacheMiss          int64
	cacheHit           int64
	dnsQuery           int64
	successfulDNSQuery int64
}

// cacheItem describes a cached dns result.
type cacheItem struct {
	host           string
	ips            []net.IPAddr
	expirationTime time.Time
	usageCount     int64
	maxUsageCount  int64
}

func newCacheItem(host string, ips []net.IPAddr, ttl time.Duration, maxUsageCount int64) *cacheItem {
	return &cacheItem{
		host:           host,
		ips:            ips,
		expirationTime: time.Now().Add(ttl),
		maxUsageCount:  maxUsageCount,
	}
}

// ip returns an ip and a bool value which indicates whether the cache is valid.
func (i *cacheItem) ip() (net.IPAddr, bool) {
	n := len(i.ips)
	if n == 0 {
		return net.IPAddr{}, false
	}

	count := atomic.AddInt64(&i.usageCount, 1)
	index := int(count-1) % n

	return i.ips[index], i.maxUsageCount >= count && time.Now().Before(i.expirationTime)
}

type dialer struct {
	chanLock          sync.Mutex
	lock              sync.RWMutex
	cache             hg.HMap[string, *cacheItem]
	dialer            *net.Dialer
	resolveChannels   hg.HMap[string, chan error]
	resolver          *net.Resolver
	stats             cacheDialerStats
	cacheDuration     time.Duration
	forceRefreshTimes int64
}

// cacheDialer creates a dialer with dns cache.
func (c *Client) cacheDialer(ttl time.Duration, maxUsage int64) *dialer {
	cd := &dialer{
		dialer:            c.dialer,
		resolver:          c.dialer.Resolver,
		cache:             hg.NewHMap[string, *cacheItem](),
		resolveChannels:   hg.NewHMap[string, chan error](),
		cacheDuration:     ttl,
		forceRefreshTimes: maxUsage,
	}

	c.transport.DialContext = cd.DialContext
	c.opt.dnsCacheStats = &cd.stats

	return cd
}

func (d *dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	atomic.AddInt64(&d.stats.totalConn, 1)

	if (network == "tcp" || network == "tcp4") && address != "" {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}

		if host != "" {
			ip, err := d.resolveHost(ctx, host)
			if err != nil {
				atomic.AddInt64(&d.stats.cacheMiss, 1)
				return nil, err
			}

			atomic.AddInt64(&d.stats.cacheHit, 1)

			address = net.JoinHostPort(ip.String(), port)
		}
	}

	return d.dialer.DialContext(ctx, network, address)
}

func (d *dialer) resolveHost(ctx context.Context, host string) (net.IPAddr, error) {
	ip, exist := d.getIPFromCache(ctx, host)
	if exist {
		return ip, nil
	}

	d.chanLock.Lock()

	ch := d.resolveChannels.Get(host)
	if ch == nil {
		ch = make(chan error, 1)
		d.resolveChannels.Set(host, ch)

		go d.resolveAndCache(ctx, host, ch)
	}

	d.chanLock.Unlock()

	select {
	case err := <-ch:
		ch <- err

		if err != nil {
			return net.IPAddr{}, err
		}

		ip, _ := d.getIPFromCache(ctx, host)

		return ip, nil
	case <-ctx.Done():
		return net.IPAddr{}, ctx.Err()
	}
}

func (d *dialer) resolveAndCache(ctx context.Context, host string, ch chan<- error) {
	atomic.AddInt64(&d.stats.dnsQuery, 1)

	var (
		item            *cacheItem
		noDNSRecordsErr = fmt.Errorf("no dns records for host %s", host)
	)

	defer func() {
		if item != nil {
			atomic.AddInt64(&d.stats.successfulDNSQuery, 1)
		}

		d.lock.Lock()
		defer d.lock.Unlock()

		if item == nil {
			d.cache.Delete(host)
		} else {
			d.cache.Set(host, item)
		}

		d.chanLock.Lock()
		defer d.chanLock.Unlock()

		d.resolveChannels.Delete(host)
	}()

	ips, err := d.resolver.LookupIPAddr(ctx, host)
	if err != nil || len(ips) == 0 {
		ch <- noDNSRecordsErr
		return
	}

	var convertedIPs []net.IPAddr

	for _, ip := range ips {
		if ip4 := ip.IP.To4(); ip4 != nil {
			convertedIPs = append(convertedIPs, net.IPAddr{IP: ip4})
		}
	}

	if len(convertedIPs) == 0 {
		ch <- noDNSRecordsErr
		return
	}

	item = newCacheItem(host, convertedIPs, d.cacheDuration, d.forceRefreshTimes)
	ch <- nil
}

func (d *dialer) getIPFromCache(_ context.Context, host string) (net.IPAddr, bool) {
	d.lock.RLock()
	item := d.cache.Get(host)
	d.lock.RUnlock()

	if item == nil {
		return net.IPAddr{}, false
	}

	ip, valid := item.ip()
	if !valid {
		d.invalidateCache(host)
	}

	return ip, valid
}

func (d *dialer) invalidateCache(host string) {
	d.lock.Lock()
	d.cache.Delete(host)
	d.lock.Unlock()
}

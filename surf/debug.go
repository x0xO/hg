package surf

import (
	"fmt"
	"io"
	"net/http/httputil"
	"strings"

	"github.com/x0xO/hg"
)

type debug struct {
	print strings.Builder
	resp  Response
}

func (resp Response) Debug() *debug { return &debug{resp: resp} }

func (d *debug) Print() { fmt.Println(d.print.String()) }

func (d *debug) DNSStat() *debug {
	if d.resp.opt == nil {
		return d
	}

	if d.resp.opt.dnsCacheStats == nil {
		return d
	}

	stats := d.resp.opt.dnsCacheStats
	fmt.Fprint(&d.print, hg.NewHString(" DNS ").Center(79, "=").Add("\n").String())
	fmt.Fprintf(&d.print, "Total Connections: %d\n", stats.totalConn)
	fmt.Fprintf(&d.print, "Total DNS Queries: %d\n", stats.dnsQuery)
	fmt.Fprintf(&d.print, "Successful DNS Queries: %d\n", stats.successfulDNSQuery)
	fmt.Fprintf(&d.print, "Cache Hit: %d\n", stats.cacheHit)
	fmt.Fprintf(&d.print, "Cache Miss: %d\n", stats.cacheMiss)
	fmt.Fprint(&d.print, hg.NewHString().Center(79, "=").Add("\n").String())

	return d
}

func (d *debug) Request(verbos ...bool) *debug {
	body, err := httputil.DumpRequestOut(d.resp.request.request, false)
	if err != nil {
		return d
	}

	fmt.Fprint(&d.print, hg.NewHString(" Request ").Center(79, "=").Add("\n").String())
	fmt.Fprintf(&d.print, "%s\n", hg.HBytes(body).TrimSpace())

	cookies := d.resp.getCookies(d.resp.request.request.URL.String())
	if len(cookies) != 0 {
		fmt.Fprint(&d.print, hg.NewHString(" Cookie ").Center(79, "=").Add("\n").String())

		for _, cookie := range cookies {
			fmt.Fprintf(&d.print, "%s\n", cookie.String())
		}
	}

	if len(verbos) != 0 && verbos[0] && d.resp.request.body != nil {
		if bytes, err := io.ReadAll(d.resp.request.body); err == nil {
			reqBody := hg.NewHBytes(bytes).TrimSpace()
			fmt.Fprint(&d.print, hg.NewHString(" ReqBody ").Center(79, "=").String())
			fmt.Fprint(&d.print, reqBody.HString().Format("\n%s\n").String())
		}
	}

	fmt.Fprint(&d.print, hg.NewHString().Center(79, "=").Add("\n").String())

	return d
}

func (d *debug) Response(verbos ...bool) *debug {
	body, err := httputil.DumpResponse(d.resp.response, false)
	if err != nil {
		return d
	}

	fmt.Fprint(&d.print, hg.NewHString(" Response ").Center(79, "=").Add("\n").String())
	fmt.Fprint(&d.print, hg.HBytes(body).TrimSpace())

	if len(verbos) != 0 && verbos[0] {
		fmt.Fprint(&d.print, hg.NewHString(" ResBody ").Center(79, "=").Format("\n%s\n").String())
		fmt.Fprint(&d.print, d.resp.Body.HString().TrimSpace())
	}

	fmt.Fprint(&d.print, hg.NewHString().Center(79, "=").Format("\n%s\n").String())

	return d
}

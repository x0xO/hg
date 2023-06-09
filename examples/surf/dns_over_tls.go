package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hg/surf"
)

func main() {
	opt := surf.NewOptions()

	// opt.DNSOverTLS().Google()
	// opt.DNSOverTLS().Switch()
	// opt.DNSOverTLS().Cloudflare()
	opt.DNSOverTLS().LibreDNS()
	// opt.DNSOverTLS().Quad9()
	// opt.DNSOverTLS().AdGuard()
	// opt.DNSOverTLS().CIRAShield()
	// opt.DNSOverTLS().Ali()
	// opt.DNSOverTLS().Quad101()
	// opt.DNSOverTLS().SB()
	// opt.DNSOverTLS().Forge()

	// custom dns provider
	// opt.DNSOverTLS().AddProvider("dns.provider.com", "0.0.0.0:853", "2.2.2.2:853")

	r, err := surf.NewClient().SetOptions(opt).Get("http://httpbingo.org/get").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.Body)
}

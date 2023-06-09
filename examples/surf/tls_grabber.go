package main

import (
	"fmt"
	"log"

	"github.com/x0xO/hg/surf"
)

func main() {
	r, err := surf.NewClient().Get("https://google.com").Do()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r.TLSGrabber().CommonName)
	fmt.Println(r.TLSGrabber().DNSNames)
	fmt.Println(r.TLSGrabber().Emails)
	fmt.Println(r.TLSGrabber().ExtensionServerName)
	fmt.Println(r.TLSGrabber().FingerprintSHA256)
	fmt.Println(r.TLSGrabber().FingerprintSHA256OpenSSL)
	fmt.Println(r.TLSGrabber().IssuerCommonName)
	fmt.Println(r.TLSGrabber().IssuerOrg)
	fmt.Println(r.TLSGrabber().Organization)
	fmt.Println(r.TLSGrabber().TLSVersion)
}

package surf

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
)

type tlsData struct {
	ExtensionServerName      string
	FingerprintSHA256        string
	FingerprintSHA256OpenSSL string
	TLSVersion               string
	CommonName               []string
	DNSNames                 []string
	Emails                   []string
	IssuerCommonName         []string
	IssuerOrg                []string
	Organization             []string
}

// tlsGrabber takes a TLS connection state and returns a tlsData struct containing information
// about the TLS connection.
func tlsGrabber(cs *tls.ConnectionState) *tlsData {
	var td tlsData

	if cs != nil {
		cert := cs.PeerCertificates[0]
		td.DNSNames = append(td.DNSNames, cert.DNSNames...)
		td.Emails = append(td.Emails, cert.EmailAddresses...)
		td.CommonName = append(td.CommonName, cert.Subject.CommonName)
		td.Organization = append(td.Organization, cert.Subject.Organization...)
		td.IssuerOrg = append(td.IssuerOrg, cert.Issuer.Organization...)
		td.IssuerCommonName = append(td.IssuerCommonName, cert.Issuer.CommonName)
		td.ExtensionServerName = cs.ServerName

		tlsVersionStringMap := map[uint16]string{
			0x0300: "SSL30",
			0x0301: "TLS10",
			0x0302: "TLS11",
			0x0303: "TLS12",
			0x0304: "TLS13",
		}

		if version, ok := tlsVersionStringMap[cs.Version]; ok {
			td.TLSVersion = version
		}

		if fingerprintSHA256, err := calculateFingerprints(cs); err == nil {
			td.FingerprintSHA256 = hex.EncodeToString(fingerprintSHA256)
			td.FingerprintSHA256OpenSSL = openSSL(fingerprintSHA256)
		}
	}

	return &td
}

// calculateFingerprints takes a TLS connection state and returns the SHA256 fingerprint
// of the first certificate in the chain or an error if no certificates are found.
func calculateFingerprints(c *tls.ConnectionState) ([]byte, error) {
	if len(c.PeerCertificates) == 0 {
		err := errors.New("no certificates found")
		return nil, err
	}

	cert := c.PeerCertificates[0]
	dataSHA256 := sha256.Sum256(cert.Raw)
	fingerprintSHA256 := dataSHA256[:]

	return fingerprintSHA256, nil
}

// openSSL takes a byte slice of a fingerprint and returns a string representation in the OpenSSL
// format.
func openSSL(fpBytes []byte) string {
	var buf bytes.Buffer

	for i, f := range fpBytes {
		if i > 0 {
			fmt.Fprintf(&buf, ":")
		}

		fmt.Fprintf(&buf, "%02X", f)
	}

	return buf.String()
}

package proxy

import (
	"crypto/tls"
	"log"
)

type CertificateProvider struct {
	certificates map[string]*tls.Certificate
	fallback *tls.Certificate
}

func NewCertificateProvider() *CertificateProvider {
	fallback, err := tls.LoadX509KeyPair("cert.crt", "cert.key")
	if err != nil {
		log.Fatalf("couldn't load fallback certificate, %v", err)
	}

	return &CertificateProvider{
		certificates: make(map[string]*tls.Certificate),
		fallback: &fallback,
	}
}

func (c *CertificateProvider) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if cert, ok := c.certificates[hello.ServerName]; ok {
		return cert, nil
	}
	log.Printf("Could not find certificate for %q", hello.ServerName)

	return c.fallback, nil
	
}

func (c *CertificateProvider) SetCertificate(name string, cert *tls.Certificate) {
	c.certificates[name] = cert
}

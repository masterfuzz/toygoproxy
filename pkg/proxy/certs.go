package proxy

import (
	"crypto/tls"
	"fmt"
)

type CertificateProvider struct {
	certificates map[string]*tls.Certificate
}

func NewCertificateProvider() *CertificateProvider {
	return &CertificateProvider{
		certificates: make(map[string]*tls.Certificate),
	}
}

func (c *CertificateProvider) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if cert, ok := c.certificates[hello.ServerName]; ok {
		return cert, nil
	}
	return nil, fmt.Errorf("Could not find certificate for %q", hello.ServerName)
	
}

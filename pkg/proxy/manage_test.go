package proxy_test

import (
	"crypto/tls"
	"testing"
)

func TestEnsureCertificate(t *testing.T) {
	manage.EnsureCertificate("hello.com")

	cert, err := certs.GetCertificate(&tls.ClientHelloInfo{ServerName: "hello.com"})
	if err != nil {
		t.Fatal(err)
	}
	if cert == nil {
		t.Fatal("certificate should not be nil")
	}
}

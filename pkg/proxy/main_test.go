package proxy_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/masterfuzz/toygoproxy/pkg/issuer"
	. "github.com/masterfuzz/toygoproxy/pkg/proxy"
	"github.com/masterfuzz/toygoproxy/pkg/test"
)

var (
	certs  *CertificateProvider
	proxy  *ProxyServer
	manage *ManagementServer
)

func TestMain(m *testing.M) {
	conn, cleanup, err := test.Database()
	defer cleanup()
	if err != nil {
		fmt.Println("Failed to set up database:", err)
		os.Exit(1)
	}
	ctx := context.Background()

	selfSigned := &issuer.SelfSignedIssuer{}
	fallbackCert, err := selfSigned.RequestCertificate("fallback")
	if err != nil {
		fmt.Println("Failed to create fallback certificate", err)
		os.Exit(1)
	}

	certs = NewCertificateProvider(ctx, conn, fallbackCert)
	proxy = NewProxyServer(conn, certs)
	manage = NewManagementServer(conn, certs, selfSigned)

	m.Run()
}

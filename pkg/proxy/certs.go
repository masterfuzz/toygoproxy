package proxy

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/masterfuzz/toygoproxy/pkg/database/postgres"
)

type CertificateProvider struct {
	certificates map[string]*tls.Certificate
	fallback     *tls.Certificate
	q            *db.Queries
}

func NewCertificateProvider(ctx context.Context, conn *pgxpool.Pool, fallback *tls.Certificate) *CertificateProvider {

	q := db.New(conn)
	dbcerts, err := q.GetCertificates(ctx)
	if err != nil {
		panic(err)
	}

	certificates := make(map[string]*tls.Certificate)
	for _, dbc := range dbcerts {
		cert, err := pemToCert(dbc.Certificate, dbc.PrivateKey)
		if err != nil {
			panic(err)
		}
		certificates[dbc.Hostname] = cert
		log.Printf("loaded certificate for %q", dbc.Hostname)
	}

	return &CertificateProvider{
		certificates: certificates,
		fallback:     fallback,
		q:            q,
	}
}

func (c *CertificateProvider) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if cert, ok := c.certificates[hello.ServerName]; ok {
		return cert, nil
	}
	log.Printf("Could not find certificate for %q", hello.ServerName)

	return c.fallback, nil

}

func (c *CertificateProvider) SetCertificate(ctx context.Context, name string, cert *tls.Certificate) error {
	certPem, keyPem, err := certToPem(cert)
	if err != nil {
		return err
	}

	_, err = c.q.InsertCertificate(ctx, db.InsertCertificateParams{
		Hostname:    name,
		Certificate: certPem,
		PrivateKey:  keyPem,
	})
	if err != nil {
		return err
	}

	c.certificates[name] = cert

	return nil
}

func pemToCert(cert string, key string) (*tls.Certificate, error) {
	tlsCert, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate and key: %w", err)
	}
	return &tlsCert, nil
}

func certToPem(cert *tls.Certificate) (string, string, error) {
	// Encode certificate to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Certificate[0], // First cert in chain
	})
	if certPEM == nil {
		return "", "", fmt.Errorf("failed to encode certificate to PEM")
	}

	// Encode private key to PEM
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(cert.PrivateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal private key: %w", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKeyBytes,
	})
	if keyPEM == nil {
		return "", "", fmt.Errorf("failed to encode private key to PEM")
	}

	return string(certPEM), string(keyPEM), nil
}

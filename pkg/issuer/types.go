package issuer

import "crypto/tls"

type Issuer interface {
	RequestCertificate(hostname string) (*tls.Certificate, error)
}

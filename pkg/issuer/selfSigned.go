package issuer

import "crypto/tls"

var _ Issuer = &SelfSignedIssuer{}

type SelfSignedIssuer struct {
}

// RequestCertificate implements Issuer.
func (s *SelfSignedIssuer) RequestCertificate(hostname string) (*tls.Certificate, error) {
	panic("unimplemented")
}


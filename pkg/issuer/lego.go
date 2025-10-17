package issuer

import (
	"crypto"
	"crypto/tls"
	"log"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

var _ Issuer = &AcmeHttp01Issuer{}

type AcmeHttp01Issuer struct {
	client *lego.Client
}

func NewAcmeHttp01Issuer(user *LegoUser) *AcmeHttp01Issuer {
	client, err := lego.NewClient(lego.NewConfig(user))
	if err != nil {
		log.Fatalf("Failed to create the lego client %v", err)
	}

	// TODO: mux this with the other http server
	client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "5000"))

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Fatal(err)
	}
	user.Registration = reg

	return &AcmeHttp01Issuer{
		client: client,
	}
}

// RequestCertificate implements Issuer.
func (a *AcmeHttp01Issuer) RequestCertificate(hostname string) (*tls.Certificate, error) {
	request := certificate.ObtainRequest{
		Domains: []string{"mydomain.com"},
		Bundle:  true,
	}
	certificates, err := a.client.Certificate.Obtain(request)
	if err != nil {
		return nil, err
	}

	tlsCert, err := tls.X509KeyPair(certificates.Certificate, certificates.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &tlsCert, nil
}

type LegoUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *LegoUser) GetEmail() string {
	return u.Email
}
func (u LegoUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *LegoUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func NewLegoUser(email string, key crypto.PrivateKey) *LegoUser {
	return &LegoUser{
		Email: email,
		key:   key,
	}
}

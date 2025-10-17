package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/masterfuzz/toygoproxy/pkg/database/postgres"
	"github.com/masterfuzz/toygoproxy/pkg/issuer"
)

var _ http.Handler = &ManagementServer{}

type ManagementServer struct {
	q *db.Queries
	certs *CertificateProvider
	issuer issuer.Issuer
}

func NewManagementServer(conn *pgxpool.Pool, certificateProvider *CertificateProvider, certificateIssuer issuer.Issuer) *ManagementServer {
	return &ManagementServer{
		q: db.New(conn),
		certs: certificateProvider,
		issuer: certificateIssuer,
	}
}


func (p *ManagementServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	hostname := r.FormValue("hostname")
	pageDataUrl := r.FormValue("page_data_url")

	if hostname == "" || pageDataUrl == "" {
		http.Error(w, "Missing required fields: hostname and page_data_url", http.StatusBadRequest)
		return
	}

	id, err := p.q.RegisterStatusPage(context.Background(), db.RegisterStatusPageParams{
		Hostname:    hostname,
		PageDataUrl: pageDataUrl,
	})
	if err != nil {
		log.Printf("Failed to register status page: %v", err)
		http.Error(w, "Failed to register status page", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Status page registered with ID: %d\n", id)

	go p.EnsureCertificate(hostname)
}

func (p *ManagementServer) EnsureCertificate(hostname string) {
	_, ok := p.certs.certificates[hostname]
	if !ok {
		log.Printf("issuing certificate for %q", hostname)
		cert, err := p.issuer.RequestCertificate(hostname)
		if err != nil {
			// TODO: what do we do here if this fails? Unregister?
			log.Printf("Error registering certificate for %q: %v", hostname, err)
			return
		}
		p.certs.SetCertificate(context.TODO(), hostname, cert)
		log.Printf("Certificate for %q set", hostname)
		return
	}
	log.Printf("%q already has a certificate", hostname)
}

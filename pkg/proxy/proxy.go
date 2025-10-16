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

var _ http.Handler = &ProxyServer{}

type ProxyServer struct {
	conn *pgxpool.Pool
	q *db.Queries
	certs *CertificateProvider
	issuer issuer.Issuer
	magic string
}

func NewProxyServer(conn *pgxpool.Pool, certificateProvider *CertificateProvider, certificateIssuer issuer.Issuer, magic string) *ProxyServer {
	return &ProxyServer{
		conn: conn,
		q: db.New(conn),
		magic: magic,
		certs: certificateProvider,
		issuer: certificateIssuer,
	}
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hostname := r.Host
	if hostname == p.magic {
		p.handleApi(w, r)
		return
	}

	log.Printf("%v: %v %v", hostname, r.Method, r.URL)
	resp, err := p.q.GetStatusPage(context.Background(), hostname)
	if err != nil {
		// TODO check if row just not found or if database error
		log.Printf("Hostname %q not found or error: %v", hostname, err)
		http.NotFound(w, r)
		return 
	}

	fmt.Fprint(w, resp)
}

func (p *ProxyServer) handleApi(w http.ResponseWriter, r *http.Request) {
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

func (p *ProxyServer) EnsureCertificate(hostname string) {
	_, ok := p.certs.certificates[hostname]
	if !ok {
		log.Printf("issuing certificate for %q", hostname)
		cert, err := p.issuer.RequestCertificate(hostname)
		if err != nil {
			// TODO: what do we do here if this fails? Unregister?
			log.Printf("Error registering certificate for %q: %v", hostname, err)
			return
		}
		p.certs.SetCertificate(hostname, cert)
		log.Printf("Certificate for %q set", hostname)
		return
	}
	log.Printf("%q already has a certificate", hostname)
}

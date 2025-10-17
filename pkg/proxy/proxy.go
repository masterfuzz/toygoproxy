package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/masterfuzz/toygoproxy/pkg/database/postgres"
)

var _ http.Handler = &ProxyServer{}

type ProxyServer struct {
	conn *pgxpool.Pool
	q *db.Queries
	certs *CertificateProvider
	magic string
}

func NewProxyServer(conn *pgxpool.Pool, certificateProvider *CertificateProvider) *ProxyServer {
	return &ProxyServer{
		conn: conn,
		q: db.New(conn),
		certs: certificateProvider,
	}
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hostname := strings.Split(r.Host, ":")[0]

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



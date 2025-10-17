package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/masterfuzz/toygoproxy/pkg/database/postgres"
	"github.com/masterfuzz/toygoproxy/pkg/metrics"
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
	start := time.Now()
	hostname := strings.Split(r.Host, ":")[0]

	// Track active connections
	metrics.ActiveConnections.Inc()
	defer metrics.ActiveConnections.Dec()

	log.Printf("%v: %v %v", hostname, r.Method, r.URL)
	resp, err := p.q.GetStatusPage(context.Background(), hostname)
	if err != nil {
		// TODO check if row just not found or if database error
		log.Printf("Hostname %q not found or error: %v", hostname, err)

		// Record metrics for failed request
		metrics.HTTPRequestsTotal.WithLabelValues(hostname, r.Method, "404").Inc()
		metrics.HTTPRequestDuration.WithLabelValues(hostname, r.Method).Observe(time.Since(start).Seconds())

		http.NotFound(w, r)
		return
	}

	// Record successful request metrics
	metrics.HTTPRequestsTotal.WithLabelValues(hostname, r.Method, "200").Inc()
	metrics.HTTPRequestDuration.WithLabelValues(hostname, r.Method).Observe(time.Since(start).Seconds())

	fmt.Fprint(w, resp)
}



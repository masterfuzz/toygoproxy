package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/masterfuzz/toygoproxy/pkg/database/postgres"
)

var _ http.Handler = &ProxyServer{}

type ProxyServer struct {
	conn *pgxpool.Pool
	q *db.Queries
}

func NewProxyServer(conn *pgxpool.Pool) *ProxyServer {
	return &ProxyServer{
		conn: conn,
		q: db.New(conn),
	}
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL)
	resp, err := p.q.GetStatusPage(context.Background(), r.URL.Hostname())
	if err != nil {
		// TODO check if row just not found or if database error
		log.Printf("Hostname not found or error: %v", err)
		http.NotFound(w, r)
		return 
	}

	fmt.Fprint(w, resp)
}

package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/masterfuzz/toygoproxy/pkg/database/postgres"
)

type ProxyServer struct {
	http.Handler
	conn *pgxpool.Pool
	q *db.Queries
}

func NewProxyServer(conn *pgxpool.Pool) *ProxyServer {
	return &ProxyServer{
		conn: conn,
		q: db.New(conn),
	}
}

func (p *ProxyServer) ServeHttp(w http.ResponseWriter, r *http.Request) {
	resp, err := p.q.GetStatusPage(context.Background(), r.URL.Hostname())
	if err != nil {
		// TODO check if row just not found or if database error
		log.Printf("Hostname not found or error: %v", err)
		http.NotFound(w, r)
		return 
	}

	fmt.Fprint(w, resp)
}

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/masterfuzz/toygoproxy/pkg/proxy"
	"github.com/masterfuzz/toygoproxy/pkg/database/migrations"
)

func main() {

	pool, err := pgxpool.New(context.Background(), "") // uses pg env vars
	if err != nil {
		fmt.Printf("Unable to create database connection pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()
	if err := migrations.Run(pool); err != nil {
		fmt.Printf("Unable to run migrations: %v\n", err)
		os.Exit(1)
	}

	prox := proxy.NewProxyServer(pool)

	mux := http.NewServeMux()
	mux.Handle("/", prox)


	certs := proxy.NewCertificateProvider()

	tlsConfig := &tls.Config{
		GetCertificate: certs.GetCertificate,
	}

	server := &http.Server{
		Addr: ":8443",
		Handler: mux,
		TLSConfig: tlsConfig,
	}

	log.Printf("Starting server on :8443")
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("Error statring https server %v", err)
	}
}

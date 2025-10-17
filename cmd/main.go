package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/masterfuzz/toygoproxy/pkg/database/migrations"
	"github.com/masterfuzz/toygoproxy/pkg/issuer"
	"github.com/masterfuzz/toygoproxy/pkg/proxy"
)

var (
	httpsPort = envOrDefault("HTTPS_PORT", "8443")
	httpPort = envOrDefault("HTTP_PORT", "8080")
	managementPort = envOrDefault("MANAGMENT_PORT", "9080")
)

func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, "") // uses pg env vars
	if err != nil {
		fmt.Printf("Unable to create database connection pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()
	if err := migrations.Run(pool); err != nil {
		fmt.Printf("Unable to run migrations: %v\n", err)
		os.Exit(1)
	}

	certs := proxy.NewCertificateProvider(ctx, pool)
	prox := proxy.NewProxyServer(pool, certs)

	// with lego issuer
	// lego := issuer.NewAcmeHttp01Issuer(issuer.NewLegoUser("user@example.com", privateKey))
	// manage := proxy.NewManagementServer(pool, certs, lego)

	manage := proxy.NewManagementServer(pool, certs, &issuer.SelfSignedIssuer{})

	mux := http.NewServeMux()
	mux.Handle("/", prox)

	tlsConfig := &tls.Config{
		GetCertificate: certs.GetCertificate,
	}

	server := &http.Server{
		Addr: ":" + httpsPort,
		Handler: mux,
		TLSConfig: tlsConfig,
	}

	// Redirect to HTTPS
	// TODO: we would need to know the external port
	go func() {
		if err := http.ListenAndServe(":" + httpPort, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, fmt.Sprintf("https://%v:%v/%v", strings.Split(r.Host, ":")[0], httpsPort, r.RequestURI), http.StatusMovedPermanently)
		})); err != nil {
			log.Fatalf("HTTP listen error: %v", err)
		}
	}()

	// Start management server
	go func() {
		if err := http.ListenAndServe(":" + managementPort, manage); err != nil {
			log.Fatalf("Management server error: %v", err)
		}
	}()

	log.Printf("Starting server on :%v", httpsPort)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("Error statring https server %v", err)
	}
}

func envOrDefault(key string, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

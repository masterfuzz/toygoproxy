package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "toygoproxy_http_requests_total",
			Help: "Total number of HTTP requests received",
		},
		[]string{"hostname", "method", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "toygoproxy_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"hostname", "method"},
	)

	// Certificate metrics
	CertificatesTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "toygoproxy_certificates_total",
			Help: "Total number of certificates stored",
		},
	)

	CertificateExpiryTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "toygoproxy_certificate_expiry_timestamp_seconds",
			Help: "Certificate expiration time in unix timestamp",
		},
		[]string{"hostname"},
	)

	CertificateIssuesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "toygoproxy_certificate_issues_total",
			Help: "Total number of certificate issuance attempts",
		},
		[]string{"hostname", "status"},
	)

	// Database metrics
	DatabaseQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "toygoproxy_database_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"query", "status"},
	)

	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "toygoproxy_database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query"},
	)

	// Proxy metrics
	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "toygoproxy_active_connections",
			Help: "Current number of active connections",
		},
	)

	BytesTransferred = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "toygoproxy_bytes_transferred_total",
			Help: "Total bytes transferred",
		},
		[]string{"hostname", "direction"},
	)
)

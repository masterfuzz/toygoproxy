# Prometheus Metrics

The toygoproxy application exposes Prometheus metrics at `http://localhost:9090/metrics` (or the port configured via `METRICS_PORT` environment variable).

## Available Metrics

### HTTP Request Metrics

- **`toygoproxy_http_requests_total`** (Counter)
  - Total number of HTTP requests received
  - Labels: `hostname`, `method`, `status`
  - Example: Track request counts per hostname and response status

- **`toygoproxy_http_request_duration_seconds`** (Histogram)
  - HTTP request latency in seconds
  - Labels: `hostname`, `method`
  - Example: Measure request processing time

### Certificate Metrics

- **`toygoproxy_certificates_total`** (Gauge)
  - Total number of certificates stored
  - Example: Monitor certificate inventory

- **`toygoproxy_certificate_expiry_timestamp_seconds`** (Gauge)
  - Certificate expiration time in unix timestamp
  - Labels: `hostname`
  - Example: Alert on upcoming certificate expirations

- **`toygoproxy_certificate_issues_total`** (Counter)
  - Total number of certificate issuance attempts
  - Labels: `hostname`, `status` (success/failed)
  - Example: Track certificate issuance success rate

### Database Metrics

- **`toygoproxy_database_queries_total`** (Counter)
  - Total number of database queries
  - Labels: `query`, `status`
  - Example: Monitor database query patterns

- **`toygoproxy_database_query_duration_seconds`** (Histogram)
  - Database query duration in seconds
  - Labels: `query`
  - Example: Identify slow queries

### Proxy Metrics

- **`toygoproxy_active_connections`** (Gauge)
  - Current number of active connections
  - Example: Monitor concurrent connection load

- **`toygoproxy_bytes_transferred_total`** (Counter)
  - Total bytes transferred
  - Labels: `hostname`, `direction` (sent/received)
  - Example: Track bandwidth usage

## Usage Example

### Scrape Configuration

Add the following to your Prometheus `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'toygoproxy'
    static_configs:
      - targets: ['localhost:9090']
```

### Example Queries

```promql
# Request rate by hostname
rate(toygoproxy_http_requests_total[5m])

# 95th percentile request latency
histogram_quantile(0.95, rate(toygoproxy_http_request_duration_seconds_bucket[5m]))

# Active connections
toygoproxy_active_connections

# Certificate expiration alert (expires in less than 30 days)
(toygoproxy_certificate_expiry_timestamp_seconds - time()) < (30 * 24 * 3600)

# Certificate issuance failure rate
rate(toygoproxy_certificate_issues_total{status="failed"}[5m])
```

## Accessing Metrics

The metrics endpoint is exposed on a dedicated metrics server:

```bash
curl http://localhost:9090/metrics
```

This will return metrics in Prometheus text format.

## Configuration

You can customize the metrics port using the `METRICS_PORT` environment variable:

```bash
export METRICS_PORT=9091
```

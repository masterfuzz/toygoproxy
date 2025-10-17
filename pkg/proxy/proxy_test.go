package proxy_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetStatusPage(t *testing.T) {
	// Test data
	hostname := "test.example.com"
	pageDataURL := "https://status.example.com/api/v1/status"

	// Register a status page through the management server
	form := url.Values{}
	form.Add("hostname", hostname)
	form.Add("page_data_url", pageDataURL)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	manage.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d. Body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	// Test the proxy server's HTTP handler
	proxyReq := httptest.NewRequest(http.MethodGet, "/", nil)
	proxyReq.Host = hostname
	proxyRR := httptest.NewRecorder()

	proxy.ServeHTTP(proxyRR, proxyReq)

	if proxyRR.Code != http.StatusOK {
		t.Fatalf("Expected proxy status code %d, got %d. Body: %s", http.StatusOK, proxyRR.Code, proxyRR.Body.String())
	}

	if proxyRR.Body.String() != pageDataURL {
		t.Errorf("Expected proxy response %q, got %q", pageDataURL, proxyRR.Body.String())
	}

	// Test 404 for non-existent hostname
	notFoundReq := httptest.NewRequest(http.MethodGet, "/", nil)
	notFoundReq.Host = "nonexistent.example.com"
	notFoundRR := httptest.NewRecorder()

	proxy.ServeHTTP(notFoundRR, notFoundReq)

	if notFoundRR.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for non-existent hostname, got %d", http.StatusNotFound, notFoundRR.Code)
	}
}

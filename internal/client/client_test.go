package client

import (
	"testing"
)

func TestNewClientSetsCookieHeader(t *testing.T) {
	cfg := &Config{
		Endpoint: "http://example.local",
		Cookie:   "session=abc123",
		Timeout:  5,
	}

	c := NewClient(cfg)
	httpClient := c.GetHTTPClient()
	// resty client stores headers in client.Header
	headers := httpClient.Header
	if headers.Get("Cookie") != "session=abc123" {
		t.Fatalf("expected Cookie header to be set, got %q", headers.Get("Cookie"))
	}
}

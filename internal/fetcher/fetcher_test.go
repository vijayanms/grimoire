package fetcher

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestFetcher(bodyByPath map[string]string, statusByPath map[string]int) *Fetcher {
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if user, pass, ok := req.BasicAuth(); !ok || user != "key" || pass != "secret" {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader("unauthorized")),
				Header:     make(http.Header),
				Request:    req,
			}, nil
		}
		status := http.StatusOK
		if statusByPath != nil {
			if s, ok := statusByPath[req.URL.Path]; ok {
				status = s
			}
		}
		body := ""
		if bodyByPath != nil {
			body = bodyByPath[req.URL.Path]
		}
		return &http.Response{
			StatusCode: status,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	})
	return &Fetcher{
		uri:       "http://opnsense.local",
		apiKey:    "key",
		apiSecret: "secret",
		client:    &http.Client{Transport: transport},
	}
}

func TestFetcherListUnderKey(t *testing.T) {
	payload := `{"alias":{"uuid-1":{"name":"test1"},"uuid-2":{"name":"test2"}}}`
	f := newTestFetcher(map[string]string{
		"/api/firewall/alias/getItem": payload,
	}, nil)
	entries, err := f.ListUnderKey(context.Background(), "/firewall/alias/getItem", "alias")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	var item struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(entries["uuid-1"], &item); err != nil {
		t.Fatalf("unmarshal uuid-1: %v", err)
	}
	if item.Name != "test1" {
		t.Errorf("expected test1, got %q", item.Name)
	}
}

func TestFetcherListUnderKeyMissingMonad(t *testing.T) {
	f := newTestFetcher(map[string]string{
		"/api/some/endpoint": `{"other_key":{}}`,
	}, nil)
	entries, err := f.ListUnderKey(context.Background(), "/some/endpoint", "alias")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty map, got %d entries", len(entries))
	}
}

func TestFetcherHTTPError(t *testing.T) {
	f := newTestFetcher(nil, map[string]int{
		"/api/any": http.StatusInternalServerError,
	})
	_, err := f.ListUnderKey(context.Background(), "/any", "monad")
	if err == nil {
		t.Fatal("expected error for non-200 response")
	}
}

package fetcher

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
)

type Fetcher struct {
	uri       string
	apiKey    string
	apiSecret string
	client    *http.Client
	opnsense  opnsense.Client
}

func New(uri, apiKey, apiSecret string, insecure bool, logger *log.Logger) *Fetcher {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}, //nolint:gosec
	}
	apiClient := api.NewClient(api.Options{
		Uri:           uri,
		APIKey:        apiKey,
		APISecret:     apiSecret,
		AllowInsecure: insecure,
		Logger:        logger,
	})
	return &Fetcher{
		uri:       uri,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client:    &http.Client{Transport: transport},
		opnsense:  opnsense.NewClient(apiClient),
	}
}

func (f *Fetcher) Client() opnsense.Client {
	return f.opnsense
}

// List fetches {uri}{endpoint} with Basic auth and returns the raw JSON map
// one level deep — the keys are UUIDs.
func (f *Fetcher) List(ctx context.Context, endpoint string) (map[string]json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.uri+"/api"+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(f.apiKey, f.apiSecret)

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: status %d: %s", endpoint, resp.StatusCode, body)
	}

	var outer map[string]json.RawMessage
	if err := json.Unmarshal(body, &outer); err != nil {
		return nil, fmt.Errorf("GET %s: unmarshal: %w", endpoint, err)
	}
	return outer, nil
}

// ListRows fetches a searchItem/searchRule-style endpoint and returns UUID→raw entries.
// The response must be {"rows": [{uuid: ..., ...}, ...]}.
func (f *Fetcher) ListRows(ctx context.Context, endpoint string) (map[string]json.RawMessage, error) {
	outer, err := f.List(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	raw, ok := outer["rows"]
	if !ok {
		return map[string]json.RawMessage{}, nil
	}
	var rows []json.RawMessage
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, fmt.Errorf("GET %s rows: unmarshal: %w", endpoint, err)
	}
	entries := make(map[string]json.RawMessage, len(rows))
	for _, row := range rows {
		var id struct {
			UUID string `json:"uuid"`
		}
		if err := json.Unmarshal(row, &id); err != nil || id.UUID == "" {
			continue
		}
		entries[id.UUID] = row
	}
	return entries, nil
}

// ListUnderKey fetches and unwraps one monad level, returning UUID→raw entries.
func (f *Fetcher) ListUnderKey(ctx context.Context, endpoint, monad string) (map[string]json.RawMessage, error) {
	outer, err := f.List(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	raw, ok := outer[monad]
	if !ok {
		return map[string]json.RawMessage{}, nil
	}
	var entries map[string]json.RawMessage
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, fmt.Errorf("GET %s[%s]: unmarshal entries: %w", endpoint, monad, err)
	}
	for k, v := range entries {
		if len(v) == 0 || v[0] != '{' {
			delete(entries, k)
		}
	}
	return entries, nil
}

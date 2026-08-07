package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
)

// httpTestFetcher is a Fetcher backed by an httptest.Server, so resource
// Fetch functions run their real HTTP + JSON-unmarshal path (including
// opnsense-go's own struct unmarshaling) against canned fixture bodies keyed
// by exact request path. Mirrors the List/ListRows/ListUnderKey behavior in
// ../fetcher.go, which resource fetch functions call through the Fetcher
// interface — duplicated here because the resources package cannot import
// package main.
type httpTestFetcher struct {
	uri    string
	client *http.Client
	opn    opnsense.Client
}

// newHTTPTestFetcher starts an httptest.Server that serves bodyByPath[path]
// for any GET request whose path matches exactly (e.g.
// "/api/firewall/alias/searchItem"). Unknown paths get a 404 naming the path,
// so a missing fixture shows up as a clear error from the Fetch call under
// test rather than a hang or an unsafe Fatal from the server goroutine.
func newHTTPTestFetcher(t *testing.T, bodyByPath map[string]string) *httpTestFetcher {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, ok := bodyByPath[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "no fixture for path %s", r.URL.Path)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	apiClient := api.NewClient(api.Options{
		Uri:       server.URL,
		APIKey:    "key",
		APISecret: "secret",
	})
	return &httpTestFetcher{
		uri:    server.URL,
		client: server.Client(),
		opn:    opnsense.NewClient(apiClient),
	}
}

func (f *httpTestFetcher) Client() opnsense.Client { return f.opn }

func (f *httpTestFetcher) list(ctx context.Context, endpoint string) (map[string]json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.uri+"/api"+endpoint, nil)
	if err != nil {
		return nil, err
	}
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

func (f *httpTestFetcher) ListRows(ctx context.Context, endpoint string) (map[string]json.RawMessage, error) {
	outer, err := f.list(ctx, endpoint)
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

func (f *httpTestFetcher) ListUnderKey(ctx context.Context, endpoint, monad string) (map[string]json.RawMessage, error) {
	outer, err := f.list(ctx, endpoint)
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

// testTracker is a no-op LabelTracker returning the uuid unchanged, keeping
// generated HCL deterministic across test runs.
type testTracker struct{}

func (testTracker) Derive(name, description, uuid string) string { return uuid }

// selectedMap renders OPNsense's "selected map" shape for api.SelectedMap /
// api.SelectedMapList / api.SelectedMapListNL fields:
// {"key1":{"value":"...","selected":1},"key2":{"value":"...","selected":0}}
// selected lists which keys are selected (order doesn't matter, unmarshal sorts).
func selectedMap(all []string, selected ...string) string {
	sel := make(map[string]bool, len(selected))
	for _, s := range selected {
		sel[s] = true
	}
	m := make(map[string]any, len(all))
	for _, k := range all {
		flag := 0
		if sel[k] {
			flag = 1
		}
		m[k] = map[string]any{"value": k, "selected": flag}
	}
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// rowsFixture builds a ListRows/searchItem-style {"rows":[...]} body from
// UUIDs; row contents beyond "uuid" are irrelevant since ListRows only reads it.
func rowsFixture(uuids ...string) string {
	rows := make([]string, len(uuids))
	for i, u := range uuids {
		rows[i] = fmt.Sprintf(`{"uuid":%q}`, u)
	}
	return fmt.Sprintf(`{"rows":[%s]}`, joinJSON(rows))
}

func joinJSON(parts []string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += ","
		}
		out += p
	}
	return out
}

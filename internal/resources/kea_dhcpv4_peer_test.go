package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchKeaDHCPv4Peer(t *testing.T) {
	const uuid = "11111111-1111-1111-1111-111111111111"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv4/search_peer": rowsFixture(uuid),
		"/api/kea/dhcpv4/get_peer/" + uuid: `{"peer":{
			"name": "peer1",
			"url": "http://10.0.0.2:8000/",
			"role": ` + selectedMap([]string{"primary", "secondary"}, "primary") + `
		}}`,
	})

	entries, err := fetchKeaDHCPv4Peer(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`name = "peer1"`,
		`url  = "http://10.0.0.2:8000/"`,
		`role = "primary"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchKeaDHCPv4PeerEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv4/search_peer": `{"rows":[]}`,
	})
	entries, err := fetchKeaDHCPv4Peer(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

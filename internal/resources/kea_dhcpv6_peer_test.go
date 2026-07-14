package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchKeaDHCPv6Peer(t *testing.T) {
	const uuid = "22222222-2222-2222-2222-222222222222"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv6/search_peer": rowsFixture(uuid),
		"/api/kea/dhcpv6/get_peer/" + uuid: `{"peer":{
			"name": "peer6",
			"url": "http://[fd00::2]:8000/",
			"role": ` + selectedMap([]string{"primary", "secondary"}, "secondary") + `
		}}`,
	})

	entries, err := fetchKeaDHCPv6Peer(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`name = "peer6"`,
		`url  = "http://[fd00::2]:8000/"`,
		`role = "secondary"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchKeaDHCPv6PeerEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv6/search_peer": `{"rows":[]}`,
	})
	entries, err := fetchKeaDHCPv6Peer(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

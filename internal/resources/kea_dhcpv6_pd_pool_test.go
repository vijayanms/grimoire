package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchKeaDHCPv6PDPool(t *testing.T) {
	const uuid = "55555555-5555-5555-5555-555555555555"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv6/search_pd_pool": rowsFixture(uuid),
		"/api/kea/dhcpv6/get_pd_pool/" + uuid: `{"pd_pool":{
			"subnet": ` + selectedMap([]string{"1", "2"}, "1") + `,
			"prefix": "fd00:1::",
			"prefix_len": "64",
			"delegated_len": "62",
			"description": "delegation pool"
		}}`,
	})

	entries, err := fetchKeaDHCPv6PDPool(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`subnet_id     = "1"`,
		`prefix        = "fd00:1::"`,
		`prefix_len    = "64"`,
		`delegated_len = "62"`,
		`description   = "delegation pool"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchKeaDHCPv6PDPoolEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv6/search_pd_pool": `{"rows":[]}`,
	})
	entries, err := fetchKeaDHCPv6PDPool(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

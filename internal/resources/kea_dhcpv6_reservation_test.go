package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchKeaDHCPv6Reservation(t *testing.T) {
	const uuid = "44444444-4444-4444-4444-444444444444"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv6/search_reservation": rowsFixture(uuid),
		"/api/kea/dhcpv6/get_reservation/" + uuid: `{"reservation":{
			"subnet": ` + selectedMap([]string{"1", "2"}, "2") + `,
			"ip_address": "fd00::50",
			"duid": "00:01:00:01:aa:bb:cc:dd:ee:ff",
			"hostname": "myhost6",
			"domain_search": ` + selectedMap([]string{"example.com", "internal.example.com"}, "example.com", "internal.example.com") + `,
			"description": "static host v6"
		}}`,
	})

	entries, err := fetchKeaDHCPv6Reservation(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`subnet_id     = "2"`,
		`ip_address    = "fd00::50"`,
		`duid          = "00:01:00:01:aa:bb:cc:dd:ee:ff"`,
		`hostname      = "myhost6"`,
		`domain_search = ["example.com", "internal.example.com"]`,
		`description   = "static host v6"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchKeaDHCPv6ReservationEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv6/search_reservation": `{"rows":[]}`,
	})
	entries, err := fetchKeaDHCPv6Reservation(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

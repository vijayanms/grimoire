package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchKeaDHCPv4Reservation(t *testing.T) {
	const uuid = "33333333-3333-3333-3333-333333333333"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv4/search_reservation": rowsFixture(uuid),
		"/api/kea/dhcpv4/get_reservation/" + uuid: `{"reservation":{
			"subnet": ` + selectedMap([]string{"1", "2"}, "1") + `,
			"ip_address": "10.0.0.50",
			"hw_address": "aa:bb:cc:dd:ee:ff",
			"hostname": "myhost",
			"description": "static host"
		}}`,
	})

	entries, err := fetchKeaDHCPv4Reservation(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`subnet_id   = "1"`,
		`ip_address  = "10.0.0.50"`,
		`mac_address = "aa:bb:cc:dd:ee:ff"`,
		`hostname    = "myhost"`,
		`description = "static host"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchKeaDHCPv4ReservationEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv4/search_reservation": `{"rows":[]}`,
	})
	entries, err := fetchKeaDHCPv4Reservation(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

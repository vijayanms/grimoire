package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchQuaggaBGPPrefixList(t *testing.T) {
	const uuid = "99999999-9999-9999-9999-999999999999"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/quagga/bgp/searchPrefixlist": rowsFixture(uuid),
		"/api/quagga/bgp/getPrefixlist/" + uuid: `{"prefixlist":{
			"enabled": "1",
			"description": "allow home net",
			"name": "HOME-NET",
			"version": ` + selectedMap([]string{"4", "6"}, "4") + `,
			"seqnumber": "15",
			"action": ` + selectedMap([]string{"permit", "deny"}, "permit") + `,
			"network": "192.168.0.0/16"
		}}`,
	})

	entries, err := fetchQuaggaBGPPrefixList(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`description = "allow home net"`,
		`name        = "HOME-NET"`,
		`ip_version  = "4"`,
		`number      = 15`,
		`action      = "permit"`,
		`network     = "192.168.0.0/16"`,
		`enabled     = true`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchQuaggaBGPPrefixListEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/quagga/bgp/searchPrefixlist": `{"rows":[]}`,
	})
	entries, err := fetchQuaggaBGPPrefixList(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

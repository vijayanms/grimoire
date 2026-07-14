package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchInterfacesVIP(t *testing.T) {
	const uuid = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/interfaces/vip_settings/searchItem": rowsFixture(uuid),
		"/api/interfaces/vip_settings/getItem/" + uuid: `{"vip":{
			"interface": ` + selectedMap([]string{"wan", "lan"}, "wan") + `,
			"mode": ` + selectedMap([]string{"proxyarp", "carp", "other"}, "proxyarp") + `,
			"network": "192.168.0.195/32",
			"descr": "Test VIP",
			"gateway": ""
		}}`,
	})

	entries, err := fetchInterfacesVIP(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`mode        = "proxyarp"`,
		`interface   = "wan"`,
		`network     = "192.168.0.195/32"`,
		`gateway     = null`,
		`description = "Test VIP"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchInterfacesVIPEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/interfaces/vip_settings/searchItem": `{"rows":[]}`,
	})
	entries, err := fetchInterfacesVIP(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

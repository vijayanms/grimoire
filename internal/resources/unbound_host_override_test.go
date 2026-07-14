package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchUnboundHostOverride(t *testing.T) {
	const uuid = "66666666-6666-6666-6666-666666666666"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/unbound/settings/searchHostOverride": rowsFixture(uuid),
		"/api/unbound/settings/getHostOverride/" + uuid: `{"host":{
			"enabled": "1",
			"hostname": "nas",
			"domain": "example.com",
			"rr": ` + selectedMap([]string{"A", "AAAA", "MX"}, "A") + `,
			"server": "10.0.0.10",
			"mxprio": "10",
			"mx": "mail.example.com",
			"description": "NAS host override"
		}}`,
	})

	entries, err := fetchUnboundHostOverride(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`hostname    = "nas"`,
		`domain      = "example.com"`,
		`type        = "A"`,
		`server      = "10.0.0.10"`,
		`mx_priority = 10`,
		`mx_host     = "mail.example.com"`,
		`enabled     = true`,
		`description = "NAS host override"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchUnboundHostOverrideEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/unbound/settings/searchHostOverride": `{"rows":[]}`,
	})
	entries, err := fetchUnboundHostOverride(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

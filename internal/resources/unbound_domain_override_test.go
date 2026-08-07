package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchUnboundDomainOverride(t *testing.T) {
	const uuid = "22222222-2222-2222-2222-222222222222"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/unbound/settings/searchDomainOverride": rowsFixture(uuid),
		"/api/unbound/settings/getDomainOverride/" + uuid: `{"domain":{
			"enabled": "1",
			"domain": "internal.example.com",
			"server": "10.0.0.53",
			"description": "Internal DNS"
		}}`,
	})

	entries, err := fetchUnboundDomainOverride(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`domain      = "internal.example.com"`,
		`server      = "10.0.0.53"`,
		`enabled     = true`,
		`description = "Internal DNS"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchUnboundDomainOverrideEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/unbound/settings/searchDomainOverride": `{"rows":[]}`,
	})
	entries, err := fetchUnboundDomainOverride(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

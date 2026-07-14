package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchUnboundHostAlias(t *testing.T) {
	const uuid = "44444444-4444-4444-4444-444444444444"
	const hostUUID = "55555555-5555-5555-5555-555555555555"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/unbound/settings/searchHostAlias": rowsFixture(uuid),
		"/api/unbound/settings/getHostAlias/" + uuid: `{"alias":{
			"enabled": "1",
			"host": ` + selectedMap([]string{hostUUID}, hostUUID) + `,
			"hostname": "alias-host",
			"domain": "example.com",
			"description": "Alias for host"
		}}`,
	})

	entries, err := fetchUnboundHostAlias(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`override    = "` + hostUUID + `"`,
		`hostname    = "alias-host"`,
		`domain      = "example.com"`,
		`enabled     = true`,
		`description = "Alias for host"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchUnboundHostAliasEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/unbound/settings/searchHostAlias": `{"rows":[]}`,
	})
	entries, err := fetchUnboundHostAlias(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

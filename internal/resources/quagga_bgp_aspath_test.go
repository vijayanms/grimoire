package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchQuaggaBGPASPath(t *testing.T) {
	const uuid = "77777777-7777-7777-7777-777777777777"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/quagga/bgp/searchAspath": rowsFixture(uuid),
		"/api/quagga/bgp/getAspath/" + uuid: `{"aspath":{
			"enabled": "1",
			"description": "block AS 65000",
			"number": "10",
			"action": ` + selectedMap([]string{"permit", "deny"}, "deny") + `,
			"as": "_65000"
		}}`,
	})

	entries, err := fetchQuaggaBGPASPath(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`description = "block AS 65000"`,
		`number      = 10`,
		`action      = "deny"`,
		`as          = "_65000"`,
		`enabled     = true`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchQuaggaBGPASPathEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/quagga/bgp/searchAspath": `{"rows":[]}`,
	})
	entries, err := fetchQuaggaBGPASPath(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

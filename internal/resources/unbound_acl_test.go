package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchUnboundACL(t *testing.T) {
	const uuid = "11111111-1111-1111-1111-111111111111"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/unbound/settings/searchAcl": rowsFixture(uuid),
		"/api/unbound/settings/get_acl/" + uuid: `{"acl":{
			"enabled": "1",
			"name": "trusted",
			"action": ` + selectedMap([]string{"allow", "deny"}, "allow") + `,
			"networks": ` + selectedMap([]string{"10.0.0.0/24", "10.0.1.0/24"}, "10.0.0.0/24", "10.0.1.0/24") + `,
			"description": "Trusted networks"
		}}`,
	})

	entries, err := fetchUnboundACL(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`name        = "trusted"`,
		`action      = "allow"`,
		`networks    = ["10.0.0.0/24", "10.0.1.0/24"]`,
		`enabled     = true`,
		`description = "Trusted networks"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchUnboundACLEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/unbound/settings/searchAcl": `{"rows":[]}`,
	})
	entries, err := fetchUnboundACL(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

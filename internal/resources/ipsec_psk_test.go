package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchIPsecPSK(t *testing.T) {
	const uuid = "55555555-5555-5555-5555-555555555555"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/ipsec/pre_shared_keys/searchItem": rowsFixture(uuid),
		"/api/ipsec/pre_shared_keys/get_item/" + uuid: `{"preSharedKey":{
			"ident": "local@example.com",
			"remote_ident": "remote@example.com",
			"Key": "supersecretkey",
			"keyType": ` + selectedMap([]string{"PSK", "EAP"}, "PSK") + `,
			"description": "psk"
		}}`,
	})

	entries, err := fetchIPsecPSK(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`identity_local  = "local@example.com"`,
		`identity_remote = "remote@example.com"`,
		`pre_shared_key  = "supersecretkey"`,
		`type            = "PSK"`,
		`description     = "psk"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchIPsecPSKEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/ipsec/pre_shared_keys/searchItem": `{"rows":[]}`,
	})
	entries, err := fetchIPsecPSK(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchOpenVPNStaticKey(t *testing.T) {
	const uuid = "99999999-9999-9999-9999-999999999999"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/openvpn/instances/search_static_key": rowsFixture(uuid),
		"/api/openvpn/instances/get_static_key/" + uuid: `{"statickey":{
			"mode": ` + selectedMap([]string{"crypt", "auth"}, "crypt") + `,
			"key": "dummy-static-key-data",
			"description": "site to site key"
		}}`,
	})

	entries, err := fetchOpenVPNStaticKey(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`mode        = "crypt"`,
		`key         = "dummy-static-key-data"`,
		`description = "site to site key"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchOpenVPNStaticKeyEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/openvpn/instances/search_static_key": `{"rows":[]}`,
	})
	entries, err := fetchOpenVPNStaticKey(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

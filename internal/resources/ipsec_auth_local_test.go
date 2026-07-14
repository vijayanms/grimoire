package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchIPsecAuthLocal(t *testing.T) {
	const uuid = "11111111-1111-1111-1111-111111111111"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/ipsec/connections/searchLocal": rowsFixture(uuid),
		"/api/ipsec/connections/get_local/" + uuid: `{"local":{
			"enabled": "1",
			"connection": ` + selectedMap([]string{"con1", "con2"}, "con1") + `,
			"round": "1",
			"auth": ` + selectedMap([]string{"pubkey", "psk"}, "pubkey") + `,
			"id": "local@example.com",
			"eap_id": "eapid",
			"certs": ` + selectedMap([]string{"certA", "certB"}, "certA", "certB") + `,
			"public_keys": ` + selectedMap([]string{"pkA"}, "pkA") + `,
			"description": "local auth"
		}}`,
	})

	entries, err := fetchIPsecAuthLocal(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`enabled          = "1"`,
		`ipsec_connection = "con1"`,
		`round            = "1"`,
		`authentication   = "pubkey"`,
		`auth_id          = "local@example.com"`,
		`eap_id           = "eapid"`,
		`certificates     = ["certA", "certB"]`,
		`public_keys      = ["pkA"]`,
		`description      = "local auth"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchIPsecAuthLocalEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/ipsec/connections/searchLocal": `{"rows":[]}`,
	})
	entries, err := fetchIPsecAuthLocal(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchIPsecAuthRemote(t *testing.T) {
	const uuid = "22222222-2222-2222-2222-222222222222"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/ipsec/connections/searchRemote": rowsFixture(uuid),
		"/api/ipsec/connections/get_remote/" + uuid: `{"remote":{
			"enabled": "1",
			"connection": ` + selectedMap([]string{"con1", "con2"}, "con2") + `,
			"round": "2",
			"auth": ` + selectedMap([]string{"pubkey", "eap-mschapv2"}, "eap-mschapv2") + `,
			"id": "remote@example.com",
			"eap_id": "eapid-remote",
			"certs": ` + selectedMap([]string{"certC"}, "certC") + `,
			"public_keys": ` + selectedMap([]string{"pkB", "pkC"}, "pkB", "pkC") + `,
			"description": "remote auth"
		}}`,
	})

	entries, err := fetchIPsecAuthRemote(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`enabled          = "1"`,
		`ipsec_connection = "con2"`,
		`round            = "2"`,
		`authentication   = "eap-mschapv2"`,
		`auth_id          = "remote@example.com"`,
		`eap_id           = "eapid-remote"`,
		`certificates     = ["certC"]`,
		`public_keys      = ["pkB", "pkC"]`,
		`description      = "remote auth"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchIPsecAuthRemoteEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/ipsec/connections/searchRemote": `{"rows":[]}`,
	})
	entries, err := fetchIPsecAuthRemote(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

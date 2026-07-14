package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchIPsecChild(t *testing.T) {
	const uuid = "33333333-3333-3333-3333-333333333333"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/ipsec/connections/searchChild": rowsFixture(uuid),
		"/api/ipsec/connections/get_child/" + uuid: `{"child":{
			"enabled": "1",
			"connection": ` + selectedMap([]string{"con1"}, "con1") + `,
			"esp_proposals": ` + selectedMap([]string{"aes256-sha256", "aes128-sha1"}, "aes256-sha256", "aes128-sha1") + `,
			"sha256_96": "0",
			"start_action": ` + selectedMap([]string{"none", "trap", "start"}, "trap") + `,
			"close_action": ` + selectedMap([]string{"none", "restart"}, "none") + `,
			"dpd_action": ` + selectedMap([]string{"none", "clear", "restart"}, "clear") + `,
			"mode": ` + selectedMap([]string{"tunnel", "transport"}, "tunnel") + `,
			"policies": "1",
			"local_ts": ` + selectedMap([]string{"10.0.0.0/24"}, "10.0.0.0/24") + `,
			"remote_ts": ` + selectedMap([]string{"10.1.0.0/24", "10.2.0.0/24"}, "10.1.0.0/24", "10.2.0.0/24") + `,
			"reqid": "1",
			"rekey_time": "3600",
			"description": "child sa"
		}}`,
	})

	entries, err := fetchIPsecChild(context.Background(), f, testTracker{})
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
		`proposals        = ["aes128-sha1", "aes256-sha256"]`,
		`sha256_96        = "0"`,
		`start_action     = "trap"`,
		`close_action     = "none"`,
		`dpd_action       = "clear"`,
		`mode             = "tunnel"`,
		`install_policies = "1"`,
		`local_networks   = ["10.0.0.0/24"]`,
		`remote_networks  = ["10.1.0.0/24", "10.2.0.0/24"]`,
		`request_id       = "1"`,
		`rekey_time       = "3600"`,
		`description      = "child sa"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchIPsecChildEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/ipsec/connections/searchChild": `{"rows":[]}`,
	})
	entries, err := fetchIPsecChild(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

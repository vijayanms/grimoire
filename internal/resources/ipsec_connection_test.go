package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchIPsecConnection(t *testing.T) {
	const uuid = "44444444-4444-4444-4444-444444444444"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/ipsec/connections/searchConnection": rowsFixture(uuid),
		"/api/ipsec/connections/get_connection/" + uuid: `{"connection":{
			"enabled": "1",
			"proposals": ` + selectedMap([]string{"aes256-sha256-modp2048"}, "aes256-sha256-modp2048") + `,
			"unique": ` + selectedMap([]string{"no", "replace", "keep"}, "replace") + `,
			"aggressive": "0",
			"version": ` + selectedMap([]string{"0", "1", "2"}, "2") + `,
			"mobike": "1",
			"local_addrs": ` + selectedMap([]string{"192.0.2.1"}, "192.0.2.1") + `,
			"remote_addrs": ` + selectedMap([]string{"198.51.100.1", "198.51.100.2"}, "198.51.100.1", "198.51.100.2") + `,
			"local_port": ` + selectedMap([]string{"500"}, "500") + `,
			"remote_port": ` + selectedMap([]string{"500"}, "500") + `,
			"encap": "0",
			"reauth_time": "0",
			"rekey_time": "7200",
			"over_time": "10800",
			"dpd_delay": "30",
			"dpd_timeout": "150",
			"pools": ` + selectedMap([]string{"pool1"}, "pool1") + `,
			"send_certreq": "1",
			"send_cert": ` + selectedMap([]string{"never", "always", "ifasked"}, "ifasked") + `,
			"keyingtries": "3",
			"description": "main conn"
		}}`,
	})

	entries, err := fetchIPsecConnection(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`enabled                 = "1"`,
		`proposals               = ["aes256-sha256-modp2048"]`,
		`unique                  = "replace"`,
		`aggressive              = "0"`,
		`version                 = "2"`,
		`mobike                  = "1"`,
		`local_addresses         = ["192.0.2.1"]`,
		`remote_addresses        = ["198.51.100.1", "198.51.100.2"]`,
		`local_port              = "500"`,
		`remote_port             = "500"`,
		`udp_encapsulation       = "0"`,
		`reauthentication_time   = "0"`,
		`rekey_time              = "7200"`,
		`ike_lifetime            = "10800"`,
		`dpd_delay               = "30"`,
		`dpd_timeout             = "150"`,
		`ip_pools                = ["pool1"]`,
		`send_certificate_request = "1"`,
		`send_certificate        = "ifasked"`,
		`keying_tries            = "3"`,
		`description             = "main conn"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchIPsecConnectionEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/ipsec/connections/searchConnection": `{"rows":[]}`,
	})
	entries, err := fetchIPsecConnection(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchWireguardServer(t *testing.T) {
	const uuid = "99999999-9999-9999-9999-999999999999"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/wireguard/server/searchServer": rowsFixture(uuid),
		"/api/wireguard/server/getServer/" + uuid: `{"server":{
			"enabled": "1",
			"name": "home-vpn",
			"instance": "0",
			"pubkey": "SERVERPUBKEY==",
			"privkey": "SERVERPRIVKEY==",
			"port": "51820",
			"mtu": "1420",
			"dns": ` + selectedMap([]string{"10.10.0.1", "10.10.0.2"}, "10.10.0.1") + `,
			"tunneladdress": ` + selectedMap([]string{"10.10.0.1/24"}, "10.10.0.1/24") + `,
			"disableroutes": "0",
			"gateway": "10.10.0.1",
			"peers": ` + selectedMap([]string{"peer1", "peer2"}, "peer1", "peer2") + `
		}}`,
	})

	entries, err := fetchWireguardServer(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`enabled        = true`,
		`name           = "home-vpn"`,
		`public_key     = "SERVERPUBKEY=="`,
		`private_key    = "SERVERPRIVKEY=="`,
		`port           = 51820`,
		`mtu            = 1420`,
		`dns            = ["10.10.0.1"]`,
		`tunnel_address = ["10.10.0.1/24"]`,
		`disable_routes = false`,
		`gateway        = "10.10.0.1"`,
		`instance       = "0"`,
		`peers          = ["peer1", "peer2"]`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchWireguardServerEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/wireguard/server/searchServer": `{"rows":[]}`,
	})
	entries, err := fetchWireguardServer(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

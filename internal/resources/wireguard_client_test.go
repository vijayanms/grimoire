package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchWireguardClient(t *testing.T) {
	const uuid = "88888888-8888-8888-8888-888888888888"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/wireguard/client/searchClient": rowsFixture(uuid),
		"/api/wireguard/client/getClient/" + uuid: `{"client":{
			"enabled": "1",
			"name": "phone",
			"pubkey": "PUBKEYBASE64==",
			"psk": "PSKBASE64==",
			"tunneladdress": ` + selectedMap([]string{"10.10.0.2/32", "fd10::2/128"}, "10.10.0.2/32", "fd10::2/128") + `,
			"serveraddress": "vpn.example.com",
			"serverport": "51820",
			"keepalive": "25"
		}}`,
	})

	entries, err := fetchWireguardClient(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`enabled        = true`,
		`name           = "phone"`,
		`public_key     = "PUBKEYBASE64=="`,
		`psk            = "PSKBASE64=="`,
		`tunnel_address = ["10.10.0.2/32", "fd10::2/128"]`,
		`server_address = "vpn.example.com"`,
		`server_port    = 51820`,
		`keepalive      = 25`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchWireguardClientEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/wireguard/client/searchClient": `{"rows":[]}`,
	})
	entries, err := fetchWireguardClient(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

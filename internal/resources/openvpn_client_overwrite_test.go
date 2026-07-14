package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchOpenVPNClientOverwrite(t *testing.T) {
	const uuid = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/openvpn/client_overwrites/search": rowsFixture(uuid),
		"/api/openvpn/client_overwrites/get/" + uuid: `{"cso":{
			"enabled": "1",
			"servers": ` + selectedMap([]string{"server1"}, "server1") + `,
			"common_name": "client1.example.com",
			"block": "0",
			"push_reset": "1",
			"tunnel_network": "10.20.0.0/24",
			"tunnel_networkv6": "",
			"local_networks": ` + selectedMap([]string{"192.168.1.0/24"}, "192.168.1.0/24") + `,
			"remote_networks": ` + selectedMap([]string{"192.168.2.0/24"}, "192.168.2.0/24") + `,
			"route_gateway": "10.20.0.1",
			"redirect_gateway": ` + selectedMap([]string{"ipv4", "ipv6"}, "ipv4") + `,
			"register_dns": "1",
			"dns_domain": ` + selectedMap([]string{"example.com"}, "example.com") + `,
			"dns_domain_search": ` + selectedMap([]string{"example.com"}, "example.com") + `,
			"dns_servers": ` + selectedMap([]string{"10.20.0.53"}, "10.20.0.53") + `,
			"ntp_servers": ` + selectedMap([]string{"10.20.0.123"}, "10.20.0.123") + `,
			"wins_servers": ` + selectedMap([]string{}) + `,
			"description": "client1 overwrite"
		}}`,
	})

	entries, err := fetchOpenVPNClientOverwrite(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`enabled           = true`,
		`common_name       = "client1.example.com"`,
		`description       = "client1 overwrite"`,
		`servers           = ["server1"]`,
		`block             = false`,
		`push_reset        = true`,
		`tunnel_network    = "10.20.0.0/24"`,
		`local_network     = ["192.168.1.0/24"]`,
		`remote_network    = ["192.168.2.0/24"]`,
		`route_gateway     = "10.20.0.1"`,
		`redirect_gateway  = ["ipv4"]`,
		`register_dns      = true`,
		`dns_domain        = ["example.com"]`,
		`dns_domain_search = ["example.com"]`,
		`dns_server        = ["10.20.0.53"]`,
		`ntp_server        = ["10.20.0.123"]`,
		`wins_server       = []`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchOpenVPNClientOverwriteEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/openvpn/client_overwrites/search": `{"rows":[]}`,
	})
	entries, err := fetchOpenVPNClientOverwrite(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchOpenVPNInstance(t *testing.T) {
	const uuid = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/openvpn/instances/search": rowsFixture(uuid),
		"/api/openvpn/instances/get/" + uuid: `{"instance":{
			"enabled": "1",
			"role": ` + selectedMap([]string{"server", "client"}, "server") + `,
			"description": "test-vpn",
			"dev_type": ` + selectedMap([]string{"tun", "tap"}, "tun") + `,
			"proto": ` + selectedMap([]string{"UDP4", "TCP4"}, "UDP4") + `,
			"port": "1194",
			"local": "0.0.0.0",
			"remote": ` + selectedMap([]string{"remote1.example.com"}, "remote1.example.com") + `,
			"topology": ` + selectedMap([]string{"subnet", "net30"}, "subnet") + `,
			"server": "10.8.0.0/24",
			"cert": ` + selectedMap([]string{"cert1"}, "cert1") + `,
			"data-ciphers": ` + selectedMap([]string{"AES-256-GCM"}, "AES-256-GCM") + `,
			"maxclients": "10",
			"keepalive_interval": "10",
			"tun_mtu": "1500",
			"dns_servers": ` + selectedMap([]string{"10.0.0.1"}, "10.0.0.1") + `,
			"ifconfig-pool-persist": "1"
		}}`,
	})

	entries, err := fetchOpenVPNInstance(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`enabled                = true`,
		`role                   = "server"`,
		`description            = "test-vpn"`,
		`dev_type               = "tun"`,
		`protocol               = "UDP4"`,
		`port                   = 1194`,
		`local                  = "0.0.0.0"`,
		`remote                 = ["remote1.example.com"]`,
		`topology               = "subnet"`,
		`server                 = "10.8.0.0/24"`,
		`certificate            = "cert1"`,
		`data_ciphers           = ["AES-256-GCM"]`,
		`max_clients            = 10`,
		`keepalive_interval     = 10`,
		`tun_mtu                = 1500`,
		`dns_servers            = ["10.0.0.1"]`,
		`ifconfig_pool_persist  = true`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchOpenVPNInstanceEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/openvpn/instances/search": `{"rows":[]}`,
	})
	entries, err := fetchOpenVPNInstance(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

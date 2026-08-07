package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchKeaDHCPv4Subnet(t *testing.T) {
	const uuid = "66666666-6666-6666-6666-666666666666"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv4/search_subnet": rowsFixture(uuid),
		"/api/kea/dhcpv4/get_subnet/" + uuid: `{"subnet4":{
			"subnet": "10.0.0.0/24",
			"next_server": "10.0.0.1",
			"pools": "10.0.0.100-10.0.0.200\n10.0.0.210-10.0.0.220",
			"match-client-id": "1",
			"option_data_autocollect": "0",
			"description": "main subnet",
			"option_data": {
				"routers": ` + selectedMap([]string{"10.0.0.1"}, "10.0.0.1") + `,
				"domain_name_servers": ` + selectedMap([]string{"10.0.0.1", "10.0.0.2"}, "10.0.0.1", "10.0.0.2") + `,
				"domain_name": "example.com",
				"domain_search": ` + selectedMap([]string{"example.com"}, "example.com") + `,
				"ntp_servers": ` + selectedMap([]string{"10.0.0.1"}, "10.0.0.1") + `,
				"time_servers": ` + selectedMap([]string{"10.0.0.1"}, "10.0.0.1") + `,
				"tftp_server_name": "tftp.example.com",
				"boot_file_name": "pxelinux.0",
				"static_routes": "192.168.1.0/24,10.0.0.254;192.168.2.0/24,10.0.0.253"
			}
		}}`,
	})

	entries, err := fetchKeaDHCPv4Subnet(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`subnet          = "10.0.0.0/24"`,
		`next_server     = "10.0.0.1"`,
		`pools           = ["10.0.0.100-10.0.0.200", "10.0.0.210-10.0.0.220"]`,
		`match_client_id = true`,
		`auto_collect    = false`,
		`description     = "main subnet"`,
		`routers         = ["10.0.0.1"]`,
		`dns_servers     = ["10.0.0.1", "10.0.0.2"]`,
		`domain_name     = "example.com"`,
		`tftp_server     = "tftp.example.com"`,
		`tftp_bootfile   = "pxelinux.0"`,
		`{ destination_ip = "192.168.1.0/24", router_ip = "10.0.0.254" }`,
		`{ destination_ip = "192.168.2.0/24", router_ip = "10.0.0.253" }`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchKeaDHCPv4SubnetEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv4/search_subnet": `{"rows":[]}`,
	})
	entries, err := fetchKeaDHCPv4Subnet(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchKeaDHCPv6Subnet(t *testing.T) {
	const uuid = "77777777-7777-7777-7777-777777777777"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv6/search_subnet": rowsFixture(uuid),
		"/api/kea/dhcpv6/get_subnet/" + uuid: `{"subnet6":{
			"subnet": "fd00:1::/64",
			"allocator": ` + selectedMap([]string{"iterative", "random"}, "iterative") + `,
			"pd-allocator": ` + selectedMap([]string{"iterative", "random"}, "random") + `,
			"pools": "fd00:1::100-fd00:1::200",
			"interface": ` + selectedMap([]string{"lan", "opt1"}, "lan") + `,
			"description": "main subnet v6",
			"option_data": {
				"dns_servers": ` + selectedMap([]string{"fd00::1"}, "fd00::1") + `,
				"domain_search": ` + selectedMap([]string{"example.com"}, "example.com") + `
			}
		}}`,
	})

	entries, err := fetchKeaDHCPv6Subnet(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`subnet       = "fd00:1::/64"`,
		`allocator    = "iterative"`,
		`pd_allocator = "random"`,
		`pools        = ["fd00:1::100-fd00:1::200"]`,
		`interface    = "lan"`,
		`description  = "main subnet v6"`,
		`domain_name_servers = ["fd00::1"]`,
		`domain_search       = ["example.com"]`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchKeaDHCPv6SubnetEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/kea/dhcpv6/search_subnet": `{"rows":[]}`,
	})
	entries, err := fetchKeaDHCPv6Subnet(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

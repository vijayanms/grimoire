package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchFirewallAlias(t *testing.T) {
	const uuid = "11111111-1111-1111-1111-111111111111"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/alias/searchItem": rowsFixture(uuid),
		"/api/firewall/alias/getItem/" + uuid: `{"alias":{
			"enabled": "1",
			"name": "LAN_Hosts",
			"type": ` + selectedMap([]string{"host", "network"}, "host") + `,
			"proto": ` + selectedMap([]string{"IPv4", "IPv6"}, "IPv4") + `,
			"interface": ` + selectedMap([]string{"lan", "wan"}, "lan") + `,
			"content": ` + selectedMap([]string{"10.0.0.1", "10.0.0.2"}, "10.0.0.1", "10.0.0.2") + `,
			"categories": ` + selectedMap([]string{"internal"}, "internal") + `,
			"updatefreq": "1.5",
			"path_expression": "",
			"counters": "0",
			"description": "LAN hosts alias"
		}}`,
	})

	entries, err := fetchFirewallAlias(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`name            = "LAN_Hosts"`,
		`type            = "host"`,
		`ip_protocol     = ["IPv4"]`,
		`interface       = "lan"`,
		`content         = ["10.0.0.1", "10.0.0.2"]`,
		`categories      = ["internal"]`,
		`update_freq     = 1.5`,
		`enabled         = true`,
		`description     = "LAN hosts alias"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchFirewallAliasEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/alias/searchItem": `{"rows":[]}`,
	})
	entries, err := fetchFirewallAlias(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchFirewallFilter(t *testing.T) {
	const uuid = "33333333-3333-3333-3333-333333333333"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/filter/searchRule": rowsFixture(uuid),
		"/api/firewall/filter/getRule/" + uuid: `{"rule":{
			"enabled": "1",
			"sequence": "10",
			"nosync": "0",
			"description": "Allow LAN to WAN",
			"categories": ` + selectedMap([]string{"internal"}, "internal") + `,
			"interfacenot": "0",
			"interface": ` + selectedMap([]string{"lan", "wan"}, "lan") + `,
			"quick": "1",
			"action": ` + selectedMap([]string{"pass", "block"}, "pass") + `,
			"allowopts": "0",
			"direction": ` + selectedMap([]string{"in", "out"}, "in") + `,
			"ipprotocol": ` + selectedMap([]string{"inet", "inet6"}, "inet") + `,
			"protocol": ` + selectedMap([]string{"TCP", "UDP"}, "TCP") + `,
			"log": "0",
			"source_net": "192.168.1.0/24",
			"source_not": "0",
			"destination_net": "any",
			"destination_port": "443",
			"destination_not": "0",
			"statetype": ` + selectedMap([]string{"keep", "none"}, "keep") + `,
			"statetimeout": "3600",
			"max": "100",
			"prio": ` + selectedMap([]string{"1", "2"}, "1") + `
		}}`,
	})

	entries, err := fetchFirewallFilter(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`description      = "Allow LAN to WAN"`,
		`interface = ["lan"]`,
		`action        = "pass"`,
		`net    = "192.168.1.0/24"`,
		`net    = "any"`,
		`type           = "keep"`,
		`states             = 100`,
		`match         = 1`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchFirewallFilterEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/filter/searchRule": `{"rows":[]}`,
	})
	entries, err := fetchFirewallFilter(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

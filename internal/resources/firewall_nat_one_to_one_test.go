package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchFirewallNATOneToOne(t *testing.T) {
	const uuid = "66666666-6666-6666-6666-666666666666"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/one_to_one/searchRule": rowsFixture(uuid),
		"/api/firewall/one_to_one/getRule/" + uuid: `{"rule":{
			"enabled": "1",
			"log": "0",
			"sequence": "1",
			"interface": ` + selectedMap([]string{"wan", "lan"}, "wan") + `,
			"type": ` + selectedMap([]string{"binat", "nat"}, "binat") + `,
			"source_net": "192.168.2.100",
			"source_not": "0",
			"destination_net": "any",
			"destination_not": "0",
			"external": "192.168.1.100",
			"natreflection": ` + selectedMap([]string{"enable", "disable"}, "enable") + `,
			"categories": ` + selectedMap([]string{"internal"}, "internal") + `,
			"description": "Test 1:1 NAT rule"
		}}`,
	})

	entries, err := fetchFirewallNATOneToOne(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`enabled        = true`,
		`sequence       = 1`,
		`interface      = "wan"`,
		`type           = "binat"`,
		`external_net   = "192.168.1.100"`,
		`nat_reflection = "enable"`,
		`categories     = ["internal"]`,
		`description    = "Test 1:1 NAT rule"`,
		`net    = "192.168.2.100"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchFirewallNATOneToOneDefaultsWhenUnset(t *testing.T) {
	const uuid = "77777777-7777-7777-7777-777777777777"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/one_to_one/searchRule": rowsFixture(uuid),
		"/api/firewall/one_to_one/getRule/" + uuid: `{"rule":{
			"enabled": "0",
			"log": "0",
			"sequence": "",
			"interface": ` + selectedMap([]string{"wan"}, "wan") + `,
			"type": ` + selectedMap([]string{"binat"}, "binat") + `,
			"source_net": "any",
			"destination_net": "any",
			"external": "192.168.1.100",
			"natreflection": ` + selectedMap([]string{"enable", "disable"}) + `,
			"description": ""
		}}`,
	})

	entries, err := fetchFirewallNATOneToOne(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	hcl := entries[0].HCL
	if !strings.Contains(hcl, "sequence       = null") {
		t.Errorf("expected null sequence when unset, got:\n%s", hcl)
	}
	if !strings.Contains(hcl, `nat_reflection = "default"`) {
		t.Errorf("expected default nat_reflection when unset, got:\n%s", hcl)
	}
}

func TestFetchFirewallNATOneToOneEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/one_to_one/searchRule": `{"rows":[]}`,
	})
	entries, err := fetchFirewallNATOneToOne(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

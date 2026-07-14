package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchFirewallNAT(t *testing.T) {
	const uuid = "44444444-4444-4444-4444-444444444444"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/source_nat/searchRule": rowsFixture(uuid),
		"/api/firewall/source_nat/getRule/" + uuid: `{"rule":{
			"enabled": "1",
			"nonat": "0",
			"sequence": "5",
			"interface": ` + selectedMap([]string{"wan", "lan"}, "wan") + `,
			"ipprotocol": ` + selectedMap([]string{"inet", "inet6"}, "inet") + `,
			"protocol": ` + selectedMap([]string{"TCP", "UDP"}, "TCP") + `,
			"source_net": "any",
			"source_not": "0",
			"destination_net": "192.168.1.0/24",
			"destination_port": "80",
			"destination_not": "0",
			"target": "10.0.0.5",
			"target_port": "8080",
			"log": "1",
			"description": "Test NAT"
		}}`,
	})

	entries, err := fetchFirewallNAT(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`enabled     = true`,
		`sequence    = 5`,
		`interface   = "wan"`,
		`ip_protocol = "inet"`,
		`protocol    = "TCP"`,
		`description = "Test NAT"`,
		`net    = "192.168.1.0/24"`,
		`ip   = "10.0.0.5"`,
		`port = "8080"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchFirewallNATEmptySequence(t *testing.T) {
	const uuid = "55555555-5555-5555-5555-555555555555"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/source_nat/searchRule": rowsFixture(uuid),
		"/api/firewall/source_nat/getRule/" + uuid: `{"rule":{
			"enabled": "0",
			"nonat": "0",
			"sequence": "",
			"interface": ` + selectedMap([]string{"wan"}, "wan") + `,
			"ipprotocol": ` + selectedMap([]string{"inet"}, "inet") + `,
			"protocol": ` + selectedMap([]string{"TCP"}, "TCP") + `,
			"source_net": "any",
			"destination_net": "any",
			"target": "10.0.0.5",
			"target_port": "80",
			"description": ""
		}}`,
	})

	entries, err := fetchFirewallNAT(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if !strings.Contains(entries[0].HCL, "sequence    = null") {
		t.Errorf("expected null sequence when unset, got:\n%s", entries[0].HCL)
	}
}

func TestFetchFirewallNATEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/source_nat/searchRule": `{"rows":[]}`,
	})
	entries, err := fetchFirewallNAT(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

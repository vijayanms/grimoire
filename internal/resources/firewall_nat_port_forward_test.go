package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchFirewallNATPortForward(t *testing.T) {
	const uuid = "88888888-8888-8888-8888-888888888888"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/d_nat/searchRule": rowsFixture(uuid),
		"/api/firewall/d_nat/getRule/" + uuid: `{"rule":{
			"disabled": "0",
			"sequence": "1",
			"interface": ` + selectedMap([]string{"wan"}, "wan") + `,
			"ipprotocol": ` + selectedMap([]string{"inet", "inet6"}, "inet") + `,
			"protocol": ` + selectedMap([]string{"TCP", "UDP"}, "TCP") + `,
			"log": "1",
			"natreflection": ` + selectedMap([]string{"enable", "purenat", "disable"}, "enable") + `,
			"descr": "Port forward web",
			"source": {"network": "any", "address": "", "port": "", "not": "0"},
			"destination": {"network": "wanip", "address": "", "port": "80", "not": "0"},
			"target": "192.168.1.50",
			"local-port": "8080"
		}}`,
	})

	entries, err := fetchFirewallNATPortForward(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`enabled        = true`,
		`interface      = ["wan"]`,
		`ip_protocol    = "inet"`,
		`protocol       = "TCP"`,
		`nat_reflection = "enable"`,
		`description    = "Port forward web"`,
		`net    = "wanip"`,
		`port   = "80"`,
		`ip   = "192.168.1.50"`,
		`port = "8080"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchFirewallNATPortForwardMultiInterface(t *testing.T) {
	const uuid = "99999999-9999-9999-9999-999999999999"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/d_nat/searchRule": rowsFixture(uuid),
		"/api/firewall/d_nat/getRule/" + uuid: `{"rule":{
			"disabled": "1",
			"sequence": "2",
			"interface": ` + selectedMap([]string{"wan", "openvpn", "lan"}, "wan", "openvpn") + `,
			"ipprotocol": ` + selectedMap([]string{"inet"}, "inet") + `,
			"protocol": ` + selectedMap([]string{"TCP"}, "TCP") + `,
			"natreflection": ` + selectedMap([]string{"enable"}) + `,
			"descr": "",
			"source": {"network": "any", "port": "", "not": "0"},
			"destination": {"network": "any", "port": "", "not": "0"},
			"target": "192.168.1.51",
			"local-port": "80"
		}}`,
	})

	entries, err := fetchFirewallNATPortForward(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	hcl := entries[0].HCL
	if !strings.Contains(hcl, `enabled        = false`) {
		t.Errorf("expected disabled rule to render enabled = false, got:\n%s", hcl)
	}
	if !strings.Contains(hcl, `interface      = ["openvpn", "wan"]`) {
		t.Errorf("expected interfaces rendered as a sorted HCL set, got:\n%s", hcl)
	}
	if !strings.Contains(hcl, `nat_reflection = "disable"`) {
		t.Errorf("expected default disable nat_reflection for unrecognized/unset value, got:\n%s", hcl)
	}
}

func TestFetchFirewallNATPortForwardEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/d_nat/searchRule": `{"rows":[]}`,
	})
	entries, err := fetchFirewallNATPortForward(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

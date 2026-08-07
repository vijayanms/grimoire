package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchDnsmasqHost(t *testing.T) {
	const uuid = "10101010-1010-1010-1010-101010101010"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/dnsmasq/settings/search_host": rowsFixture(uuid),
		"/api/dnsmasq/settings/get_host/" + uuid: `{"host":{
			"host": "printer",
			"domain": "lan",
			"local": "0",
			"ip": ` + selectedMap([]string{"192.168.1.10", "192.168.1.11"}, "192.168.1.10", "192.168.1.11") + `,
			"aliases": ` + selectedMap([]string{"alias1"}, "alias1") + `,
			"cnames": ` + selectedMap([]string{"cname1"}) + `,
			"client_id": "",
			"hwaddr": ` + selectedMap([]string{"AA:BB:CC:DD:EE:FF"}, "AA:BB:CC:DD:EE:FF") + `,
			"set_tag": ` + selectedMap([]string{"tag1", "tag2"}, "tag1") + `,
			"ignore": "0",
			"descr": "Printer host",
			"comments": "office printer"
		}}`,
	})

	entries, err := fetchDnsmasqHost(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`hostname           = "printer"`,
		`domain             = "lan"`,
		`ip_addresses       = ["192.168.1.10", "192.168.1.11"]`,
		`alias_records      = ["alias1"]`,
		`cname_records      = []`,
		`client_id          = null`,
		`hardware_addresses = ["AA:BB:CC:DD:EE:FF"]`,
		`tag                = "tag1"`,
		`description        = "Printer host"`,
		`comment            = "office printer"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchDnsmasqHostEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/dnsmasq/settings/search_host": `{"rows":[]}`,
	})
	entries, err := fetchDnsmasqHost(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

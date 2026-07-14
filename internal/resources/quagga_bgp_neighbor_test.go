package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchQuaggaBGPNeighbor(t *testing.T) {
	const uuid = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/quagga/bgp/searchNeighbor": rowsFixture(uuid),
		"/api/quagga/bgp/getNeighbor/" + uuid: `{"neighbor":{
			"enabled": "1",
			"description": "upstream peer",
			"address": "203.0.113.1",
			"remoteas": "65001",
			"password": "s3cr3t",
			"weight": "100",
			"localip": "203.0.113.2",
			"updatesource": ` + selectedMap([]string{"lan", "wan"}, "wan") + `,
			"linklocalinterface": ` + selectedMap([]string{"lan", "wan"}, "wan") + `,
			"nexthopself": "1",
			"nexthopselfall": "0",
			"multihop": "0",
			"multiprotocol": "0",
			"rrclient": "0",
			"bfd": "1",
			"keepalive": "10",
			"holddown": "30",
			"connecttimer": "5",
			"defaultoriginate": "0",
			"asoverride": "0",
			"disable_connected_check": "0",
			"attributeunchanged": ` + selectedMap([]string{"as-path", "med"}, "as-path") + `,
			"linkedPrefixlistIn": ` + selectedMap([]string{"pfx-in"}, "pfx-in") + `,
			"linkedPrefixlistOut": ` + selectedMap([]string{"pfx-out"}, "pfx-out") + `,
			"linkedRoutemapIn": ` + selectedMap([]string{"rm-in"}, "rm-in") + `,
			"linkedRoutemapOut": ` + selectedMap([]string{"rm-out"}, "rm-out") + `
		}}`,
	})

	entries, err := fetchQuaggaBGPNeighbor(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`description              = "upstream peer"`,
		`peer_ip                  = "203.0.113.1"`,
		`remote_as                = 65001`,
		`password                 = "s3cr3t"`,
		`weight                   = 100`,
		`local_ip                 = "203.0.113.2"`,
		`update_source            = "wan"`,
		`link_local_interface     = "wan"`,
		`next_hop_self            = true`,
		`next_hop_self_all        = false`,
		`bfd                      = true`,
		`keepalive                = 10`,
		`hold_down                = 30`,
		`connect_timer            = 5`,
		`attribute_unchanged      = "as-path"`,
		`prefix_list_in           = "pfx-in"`,
		`prefix_list_out          = "pfx-out"`,
		`route_map_in             = "rm-in"`,
		`route_map_out            = "rm-out"`,
		`enabled                  = true`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchQuaggaBGPNeighborEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/quagga/bgp/searchNeighbor": `{"rows":[]}`,
	})
	entries, err := fetchQuaggaBGPNeighbor(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

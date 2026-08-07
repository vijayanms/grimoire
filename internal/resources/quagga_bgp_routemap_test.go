package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchQuaggaBGPRouteMap(t *testing.T) {
	const uuid = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/quagga/bgp/searchRoutemap": rowsFixture(uuid),
		"/api/quagga/bgp/getRoutemap/" + uuid: `{"routemap":{
			"enabled": "1",
			"description": "outbound policy",
			"name": "OUT-POLICY",
			"action": ` + selectedMap([]string{"permit", "deny"}, "permit") + `,
			"id": "10",
			"set": "local-preference 200",
			"match": ` + selectedMap([]string{"aspath1", "aspath2"}, "aspath1") + `,
			"match2": ` + selectedMap([]string{"prefix1", "prefix2"}, "prefix1", "prefix2") + `,
			"match3": ` + selectedMap([]string{"comm1"}, "comm1") + `
		}}`,
	})

	entries, err := fetchQuaggaBGPRouteMap(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`description    = "outbound policy"`,
		`name           = "OUT-POLICY"`,
		`action         = "permit"`,
		`route_map_id   = 10`,
		`set            = "local-preference 200"`,
		`aspath_list    = ["aspath1"]`,
		`prefix_list    = ["prefix1", "prefix2"]`,
		`community_list = ["comm1"]`,
		`enabled        = true`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchQuaggaBGPRouteMapEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/quagga/bgp/searchRoutemap": `{"rows":[]}`,
	})
	entries, err := fetchQuaggaBGPRouteMap(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

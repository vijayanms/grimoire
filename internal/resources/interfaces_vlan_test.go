package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchInterfacesVLAN(t *testing.T) {
	const uuid = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/interfaces/vlan_settings/searchItem": rowsFixture(uuid),
		"/api/interfaces/vlan_settings/getItem/" + uuid: `{"vlan":{
			"descr": "Test VLAN",
			"tag": "100",
			"pcp": ` + selectedMap([]string{"0", "1", "2"}, "1") + `,
			"if": ` + selectedMap([]string{"vtnet0", "vtnet1"}, "vtnet0") + `,
			"vlanif": "vtnet0_vlan100"
		}}`,
	})

	entries, err := fetchInterfacesVLAN(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`description = "Test VLAN"`,
		`tag         = 100`,
		`priority    = "1"`,
		`parent      = "vtnet0"`,
		`device      = "vtnet0_vlan100"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchInterfacesVLANNullTagWhenUnset(t *testing.T) {
	const uuid = "cccccccc-cccc-cccc-cccc-cccccccccccc"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/interfaces/vlan_settings/searchItem": rowsFixture(uuid),
		"/api/interfaces/vlan_settings/getItem/" + uuid: `{"vlan":{
			"descr": "Auto-assigned VLAN",
			"tag": "",
			"pcp": ` + selectedMap([]string{"0"}, "0") + `,
			"if": ` + selectedMap([]string{"vtnet0"}, "vtnet0") + `,
			"vlanif": ""
		}}`,
	})

	entries, err := fetchInterfacesVLAN(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if !strings.Contains(entries[0].HCL, "tag         = null") {
		t.Errorf("expected null tag when unset, got:\n%s", entries[0].HCL)
	}
}

func TestFetchInterfacesVLANEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/interfaces/vlan_settings/searchItem": `{"rows":[]}`,
	})
	entries, err := fetchInterfacesVLAN(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

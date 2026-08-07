package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchIPsecVTI(t *testing.T) {
	const uuid = "66666666-6666-6666-6666-666666666666"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/ipsec/vti/search": rowsFixture(uuid),
		"/api/ipsec/vti/get/" + uuid: `{"vti":{
			"enabled": "1",
			"reqid": "1",
			"local": "192.0.2.1",
			"remote": "198.51.100.1",
			"tunnel_local": "10.10.10.1",
			"tunnel_remote": "10.10.10.2",
			"tunnel_local2": "",
			"tunnel_remote2": "",
			"description": "vti0"
		}}`,
	})

	entries, err := fetchIPsecVTI(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`enabled           = "1"`,
		`request_id        = "1"`,
		`local_ip          = "192.0.2.1"`,
		`remote_ip         = "198.51.100.1"`,
		`tunnel_local_ip   = "10.10.10.1"`,
		`tunnel_remote_ip  = "10.10.10.2"`,
		`description       = "vti0"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchIPsecVTIEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/ipsec/vti/search": `{"rows":[]}`,
	})
	entries, err := fetchIPsecVTI(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

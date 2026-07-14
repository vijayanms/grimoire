package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchUnboundForward(t *testing.T) {
	const uuid = "33333333-3333-3333-3333-333333333333"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/unbound/settings/searchDot": rowsFixture(uuid),
		"/api/unbound/settings/getDot/" + uuid: `{"dot":{
			"enabled": "1",
			"domain": "example.com",
			"type": ` + selectedMap([]string{"forward", "static"}, "forward") + `,
			"server": "1.1.1.1",
			"port": "853",
			"verify": "cloudflare-dns.com"
		}}`,
	})

	entries, err := fetchUnboundForward(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`domain      = "example.com"`,
		`server_ip   = "1.1.1.1"`,
		`server_port = 853`,
		`verify_cn   = "cloudflare-dns.com"`,
		`enabled     = true`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchUnboundForwardEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/unbound/settings/searchDot": `{"rows":[]}`,
	})
	entries, err := fetchUnboundForward(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchRoutesRoute(t *testing.T) {
	const uuid = "dddddddd-dddd-dddd-dddd-dddddddddddd"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/routes/routes/searchroute": rowsFixture(uuid),
		"/api/routes/routes/getroute/" + uuid: `{"route":{
			"disabled": "0",
			"gateway": ` + selectedMap([]string{"GW_WAN", "GW_LAN"}, "GW_WAN") + `,
			"network": "10.10.10.0/24",
			"descr": "Test route"
		}}`,
	})

	entries, err := fetchRoutesRoute(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`enabled     = true`,
		`gateway     = "GW_WAN"`,
		`network     = "10.10.10.0/24"`,
		`description = "Test route"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchRoutesRouteDisabled(t *testing.T) {
	const uuid = "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/routes/routes/searchroute": rowsFixture(uuid),
		"/api/routes/routes/getroute/" + uuid: `{"route":{
			"disabled": "1",
			"gateway": ` + selectedMap([]string{"GW_WAN"}, "GW_WAN") + `,
			"network": "10.10.20.0/24",
			"descr": ""
		}}`,
	})

	entries, err := fetchRoutesRoute(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if !strings.Contains(entries[0].HCL, "enabled     = false") {
		t.Errorf("expected disabled route to render enabled = false, got:\n%s", entries[0].HCL)
	}
}

func TestFetchRoutesRouteEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/routes/routes/searchroute": `{"rows":[]}`,
	})
	entries, err := fetchRoutesRoute(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchCronJob(t *testing.T) {
	const uuid = "ffffffff-ffff-ffff-ffff-ffffffffffff"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/cron/settings/searchJobs": rowsFixture(uuid),
		"/api/cron/settings/getJob/" + uuid: `{"job":{
			"enabled": "1",
			"minutes": "0",
			"hours": "4",
			"days": "*",
			"months": "*",
			"weekdays": "*",
			"who": "root",
			"command": ` + selectedMap([]string{"firmware poll", "backup"}, "firmware poll") + `,
			"parameters": "",
			"description": "Nightly poll"
		}}`,
	})

	entries, err := fetchCronJob(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`enabled     = true`,
		`hours       = "4"`,
		`who         = "root"`,
		`command     = "firmware poll"`,
		`parameters  = null`,
		`description = "Nightly poll"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchCronJobEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/cron/settings/searchJobs": `{"rows":[]}`,
	})
	entries, err := fetchCronJob(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

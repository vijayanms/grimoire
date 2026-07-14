package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchFirewallCategory(t *testing.T) {
	const uuid = "22222222-2222-2222-2222-222222222222"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/category/searchItem": rowsFixture(uuid),
		"/api/firewall/category/getItem/" + uuid: `{"category":{
			"auto": "1",
			"name": "Cat1",
			"color": "FF0000"
		}}`,
	})

	entries, err := fetchFirewallCategory(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`auto  = true`,
		`name  = "Cat1"`,
		`color = "FF0000"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchFirewallCategoryEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/firewall/category/searchItem": `{"rows":[]}`,
	})
	entries, err := fetchFirewallCategory(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

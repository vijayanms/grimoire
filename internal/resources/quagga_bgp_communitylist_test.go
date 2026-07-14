package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchQuaggaBGPCommunityList(t *testing.T) {
	const uuid = "88888888-8888-8888-8888-888888888888"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/quagga/bgp/searchCommunitylist": rowsFixture(uuid),
		"/api/quagga/bgp/getCommunitylist/" + uuid: `{"communitylist":{
			"enabled": "1",
			"description": "customer routes",
			"number": "20",
			"seqnumber": "5",
			"action": ` + selectedMap([]string{"permit", "deny"}, "permit") + `,
			"community": "65000:100"
		}}`,
	})

	entries, err := fetchQuaggaBGPCommunityList(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`description     = "customer routes"`,
		`number          = 20`,
		`sequence_number = 5`,
		`action          = "permit"`,
		`community       = "65000:100"`,
		`enabled         = true`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchQuaggaBGPCommunityListEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/quagga/bgp/searchCommunitylist": `{"rows":[]}`,
	})
	entries, err := fetchQuaggaBGPCommunityList(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

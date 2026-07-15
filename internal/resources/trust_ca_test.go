package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchTrustCA(t *testing.T) {
	const uuid = "77777777-7777-7777-7777-777777777777"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/trust/ca/search": rowsFixture(uuid),
		"/api/trust/ca/get/" + uuid: `{"ca":{
			"refid": "5f1a2b3c",
			"descr": "internal ca",
			"action": ` + selectedMap([]string{"internal", "import"}, "internal") + `,
			"crt": "dummy-cert-data",
			"prv": "dummy-key-data",
			"serial": "1",
			"caref": ` + selectedMap([]string{"someref"}, "someref") + `,
			"key_type": ` + selectedMap([]string{"RSA", "ECDSA"}, "RSA") + `,
			"lifetime": "3650",
			"digest": ` + selectedMap([]string{"sha256", "sha512"}, "sha256") + `,
			"country": ` + selectedMap([]string{"US", "GB"}, "US") + `,
			"state": "California",
			"city": "San Francisco",
			"organization": "Acme",
			"organizationalunit": "IT",
			"email": "admin@example.com",
			"commonname": "Acme Root CA",
			"ocsp_uri": "http://ocsp.example.com",
			"crt_payload": "",
			"prv_payload": "",
			"name": "",
			"valid_from": "",
			"valid_to": ""
		}}`,
	})

	entries, err := fetchTrustCA(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`description          = "internal ca"`,
		`action               = "internal"`,
		`crt                  = "dummy-cert-data"`,
		`prv                  = "dummy-key-data"`,
		`serial               = "1"`,
		`caref                = "someref"`,
		`key_type             = "RSA"`,
		`lifetime             = "3650"`,
		`digest               = "sha256"`,
		`country              = "US"`,
		`state                = "California"`,
		`city                 = "San Francisco"`,
		`organization         = "Acme"`,
		`organizational_unit  = "IT"`,
		`email                = "admin@example.com"`,
		`common_name          = "Acme Root CA"`,
		`ocsp_uri             = "http://ocsp.example.com"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchTrustCAEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/trust/ca/search": `{"rows":[]}`,
	})
	entries, err := fetchTrustCA(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

package resources

import (
	"context"
	"strings"
	"testing"
)

func TestFetchTrustCert(t *testing.T) {
	const uuid = "88888888-8888-8888-8888-888888888888"
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/trust/cert/search": rowsFixture(uuid),
		"/api/trust/cert/get/" + uuid: `{"cert":{
			"refid": "6a2b3c4d",
			"descr": "server cert",
			"caref": ` + selectedMap([]string{"internal-ca"}, "internal-ca") + `,
			"crt": "dummy-cert-data",
			"csr": "",
			"prv": "dummy-key-data",
			"action": ` + selectedMap([]string{"internal", "import"}, "internal") + `,
			"key_type": ` + selectedMap([]string{"RSA", "ECDSA"}, "ECDSA") + `,
			"digest": ` + selectedMap([]string{"sha256", "sha512"}, "sha256") + `,
			"cert_type": ` + selectedMap([]string{"server", "client"}, "server") + `,
			"lifetime": "825",
			"private_key_location": ` + selectedMap([]string{"local", "remote"}, "local") + `,
			"country": ` + selectedMap([]string{"US", "GB"}, "GB") + `,
			"state": "London",
			"city": "London",
			"organization": "Acme",
			"organizationalunit": "Ops",
			"email": "ops@example.com",
			"commonname": "server.example.com",
			"ocsp_uri": "",
			"altnames_dns": "server.example.com",
			"altnames_ip": "",
			"altnames_uri": "",
			"altnames_email": "",
			"in_use": "1",
			"is_user": "0",
			"crt_payload": "",
			"csr_payload": "",
			"prv_payload": "",
			"name": "",
			"valid_from": "",
			"valid_to": ""
		}}`,
	})

	entries, err := fetchTrustCert(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	hcl := entries[0].HCL
	for _, want := range []string{
		`ref_id               = "6a2b3c4d"`,
		`description          = "server cert"`,
		`caref                = "internal-ca"`,
		`crt                  = "dummy-cert-data"`,
		`csr                  = null`,
		`prv                  = "dummy-key-data"`,
		`action               = "internal"`,
		`key_type             = "ECDSA"`,
		`digest               = "sha256"`,
		`cert_type            = "server"`,
		`lifetime             = "825"`,
		`private_key_location = "local"`,
		`country              = "GB"`,
		`state                = "London"`,
		`city                 = "London"`,
		`organization         = "Acme"`,
		`organizational_unit  = "Ops"`,
		`email                = "ops@example.com"`,
		`common_name          = "server.example.com"`,
		`altnames_dns         = "server.example.com"`,
		`in_use               = "1"`,
		`is_user              = "0"`,
	} {
		if !strings.Contains(hcl, want) {
			t.Errorf("HCL missing %q\ngot:\n%s", want, hcl)
		}
	}
}

func TestFetchTrustCertEmpty(t *testing.T) {
	f := newHTTPTestFetcher(t, map[string]string{
		"/api/trust/cert/search": `{"rows":[]}`,
	})
	entries, err := fetchTrustCert(context.Background(), f, testTracker{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

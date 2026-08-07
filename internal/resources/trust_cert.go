package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_trust_cert",
		Filename: "trust_cert.tf",
		Fetch:    fetchTrustCert,
	})
}

func fetchTrustCert(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/trust/cert/search")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Trust().GetCert(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("trust_cert %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Description, d.CommonName, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_trust_cert" %s {
  description          = %s
  caref                = %s
  crt                  = %s
  csr                  = %s
  prv                  = %s
  action               = %s
  key_type             = %s
  digest               = %s
  cert_type            = %s
  lifetime             = %s
  private_key_location = %s
  country              = %s
  state                = %s
  city                 = %s
  organization         = %s
  organizational_unit  = %s
  email                = %s
  common_name          = %s
  ocsp_uri             = %s
  altnames_dns         = %s
  altnames_ip          = %s
  altnames_uri         = %s
  altnames_email       = %s
}
`, hclString(label),
			hclString(d.Description),
			hclString(d.CaRef.String()),
			hclString(d.Crt),
			hclStringOrNull(d.Csr),
			hclString(d.Prv),
			hclString(d.Action.String()),
			hclString(d.KeyType.String()),
			hclString(d.Digest.String()),
			hclString(d.CertType.String()),
			hclString(d.Lifetime),
			hclString(d.PrivateKeyLocation.String()),
			hclString(d.Country.String()),
			hclString(d.State),
			hclString(d.City),
			hclString(d.Organization),
			hclString(d.OrganizationalUnit),
			hclString(d.Email),
			hclString(d.CommonName),
			hclStringOrNull(d.OcspUri),
			hclStringOrNull(d.AltnamesDns),
			hclStringOrNull(d.AltnamesIp),
			hclStringOrNull(d.AltnamesUri),
			hclStringOrNull(d.AltnamesEmail))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

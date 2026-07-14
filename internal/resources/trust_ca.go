package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_trust_ca",
		Filename: "trust_ca.tf",
		Fetch:    fetchTrustCA,
	})
}

func fetchTrustCA(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/trust/ca/search")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Trust().GetCa(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("trust_ca %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Description, d.CommonName, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_trust_ca" %s {
  ref_id               = %s
  description          = %s
  action               = %s
  crt                  = %s
  prv                  = %s
  serial               = %s
  ca_ref               = %s
  key_type             = %s
  lifetime             = %s
  digest               = %s
  country              = %s
  state                = %s
  city                 = %s
  organization         = %s
  organizational_unit  = %s
  email                = %s
  common_name          = %s
  ocsp_uri             = %s
  crt_payload          = %s
  prv_payload          = %s
  name                 = %s
  valid_from           = %s
  valid_to             = %s
}
`, hclString(label),
			hclStringOrNull(d.RefId),
			hclString(d.Description),
			hclString(d.Action.String()),
			hclString(d.Crt),
			hclString(d.Prv),
			hclStringOrNull(d.Serial),
			hclString(d.CaRef.String()),
			hclString(d.KeyType.String()),
			hclString(d.Lifetime),
			hclString(d.Digest.String()),
			hclString(d.Country.String()),
			hclString(d.State),
			hclString(d.City),
			hclString(d.Organization),
			hclString(d.OrganizationalUnit),
			hclString(d.Email),
			hclString(d.CommonName),
			hclString(d.OcspUri),
			hclStringOrNull(d.CrtPayload),
			hclStringOrNull(d.PrvPayload),
			hclStringOrNull(d.Name),
			hclStringOrNull(d.ValidFrom),
			hclStringOrNull(d.ValidTo))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

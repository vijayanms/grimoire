package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_ipsec_auth_remote",
		Filename: "ipsec_auth_remote.tf",
		Fetch:    fetchIPsecAuthRemote,
	})
}

func fetchIPsecAuthRemote(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/ipsec/connections/searchRemote")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Ipsec().GetIPsecAuthRemote(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("ipsec_auth_remote %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_ipsec_auth_remote" %s {
  enabled          = %s
  ipsec_connection = %s
  round            = %s
  authentication   = %s
  auth_id          = %s
  eap_id           = %s
  certificates     = %s
  public_keys      = %s
  description      = %s
}
`, hclString(label),
			hclString(d.Enabled),
			hclString(d.Connection.String()),
			hclString(d.Round),
			hclString(d.Authentication.String()),
			hclString(d.Id),
			hclString(d.EAPId),
			hclSet([]string(d.Certificates)),
			hclSet([]string(d.PublicKeys)),
			hclString(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

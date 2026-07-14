package resources

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/ipsec"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_ipsec_auth_local",
		Filename: "ipsec_auth_local.tf",
		Fetch:    fetchIPsecAuthLocal,
	})
}

type ipsecAuthStruct = ipsec.IPsecAuthLocal

func fetchIPsecAuthLocal(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/ipsec/connections/searchLocal")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Ipsec().GetIPsecAuthLocal(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("ipsec_auth_local %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_ipsec_auth_local" %s {
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

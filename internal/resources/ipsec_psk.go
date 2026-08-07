package resources

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/ipsec"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_ipsec_psk",
		Filename: "ipsec_psk.tf",
		Fetch:    fetchIPsecPSK,
	})
}

type ipsecPSKStruct = ipsec.IPsecPSK

func fetchIPsecPSK(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/ipsec/pre_shared_keys/searchItem")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Ipsec().GetIPsecPSK(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("ipsec_psk %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_ipsec_psk" %s {
  identity_local  = %s
  identity_remote = %s
  pre_shared_key  = %s
  type            = %s
  description     = %s
}
`, hclString(label),
			hclString(d.IdentityLocal),
			hclString(d.IdentityRemote),
			hclString(d.PreSharedKey),
			hclString(d.Type.String()),
			hclString(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

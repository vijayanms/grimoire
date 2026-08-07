package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_openvpn_static_key",
		Filename: "openvpn_static_key.tf",
		Fetch:    fetchOpenVPNStaticKey,
	})
}

func fetchOpenVPNStaticKey(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/openvpn/instances/search_static_key")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Openvpn().GetStaticKey(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("openvpn_static_key %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_openvpn_static_key" %s {
  mode        = %s
  key         = %s
  description = %s
}
`, hclString(label),
			hclString(d.Mode.String()),
			hclString(d.Key),
			hclString(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

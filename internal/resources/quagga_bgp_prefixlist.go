package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_quagga_bgp_prefixlist",
		Filename: "quagga_bgp_prefixlist.tf",
		Fetch:    fetchQuaggaBGPPrefixList,
	})
}

func fetchQuaggaBGPPrefixList(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/quagga/bgp/searchPrefixlist")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Quagga().GetBGPPrefixList(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("quagga_bgp_prefixlist %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Name, d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_quagga_bgp_prefixlist" %s {
  enabled     = %s
  description = %s
  name        = %s
  ip_version  = %s
  number      = %s
  action      = %s
  network     = %s
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclString(d.Description),
			hclString(d.Name),
			hclString(d.IPVersion.String()),
			hclInt(stringToInt64(d.SequenceNumber)),
			hclString(d.Action.String()),
			hclString(d.Network))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_quagga_bgp_aspath",
		Filename: "quagga_bgp_aspath.tf",
		Fetch:    fetchQuaggaBGPASPath,
	})
}

func fetchQuaggaBGPASPath(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/quagga/bgp/searchAspath")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Quagga().GetBGPASPath(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("quagga_bgp_aspath %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_quagga_bgp_aspath" %s {
  enabled     = %s
  description = %s
  number      = %s
  action      = %s
  as          = %s
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclString(d.Description),
			hclInt(stringToInt64(d.Number)),
			hclString(d.Action.String()),
			hclString(d.AS))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

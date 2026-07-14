package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_quagga_bgp_communitylist",
		Filename: "quagga_bgp_communitylist.tf",
		Fetch:    fetchQuaggaBGPCommunityList,
	})
}

func fetchQuaggaBGPCommunityList(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/quagga/bgp/searchCommunitylist")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Quagga().GetBGPCommunityList(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("quagga_bgp_communitylist %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_quagga_bgp_communitylist" %s {
  enabled         = %s
  description     = %s
  number          = %s
  sequence_number = %s
  action          = %s
  community       = %s
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclString(d.Description),
			hclInt(stringToInt64(d.Number)),
			hclInt(stringToInt64(d.SequenceNumber)),
			hclString(d.Action.String()),
			hclString(d.Community))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

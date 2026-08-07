package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_kea_dhcpv6_pd_pool",
		Filename: "kea_dhcpv6_pd_pool.tf",
		Fetch:    fetchKeaDHCPv6PDPool,
	})
}

func fetchKeaDHCPv6PDPool(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/kea/dhcpv6/search_pd_pool")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Kea().GetPDPool(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("kea_dhcpv6_pd_pool %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Prefix, d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_kea_dhcpv6_pd_pool" %s {
  subnet_id     = %s
  prefix        = %s
  prefix_len    = %s
  delegated_len = %s
  description   = %s
}
`, hclString(label),
			hclString(d.Subnet.String()),
			hclString(d.Prefix),
			hclString(d.PrefixLen),
			hclString(d.DelegatedLen),
			hclString(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

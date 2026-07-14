package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_kea_dhcpv6_subnet",
		Filename: "kea_dhcpv6_subnet.tf",
		Fetch:    fetchKeaDHCPv6Subnet,
	})
}

func fetchKeaDHCPv6Subnet(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/kea/dhcpv6/search_subnet")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Kea().GetSubnetV6(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("kea_dhcpv6_subnet %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Subnet, d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_kea_dhcpv6_subnet" %s {
  subnet       = %s
  allocator    = %s
  pd_allocator = %s
  pools        = %s
  interface    = %s
  description  = %s

  option_data {
    domain_name_servers = %s
    domain_search       = %s
  }
}
`, hclString(label),
			hclString(d.Subnet),
			hclString(d.Allocator.String()),
			hclString(d.PDAllocator.String()),
			hclSet(splitNL(d.Pools)),
			hclString(d.Interface.String()),
			hclString(d.Description),
			hclSet([]string(d.OptionData.DomainNameServers)),
			hclSet([]string(d.OptionData.DomainSearch)))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_kea_dhcpv6_peer",
		Filename: "kea_dhcpv6_peer.tf",
		Fetch:    fetchKeaDHCPv6Peer,
	})
}

func fetchKeaDHCPv6Peer(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/kea/dhcpv6/search_peer")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Kea().GetPeerV6(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("kea_dhcpv6_peer %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Name, "", uuid)
		hcl := fmt.Sprintf(`resource "opnsense_kea_dhcpv6_peer" %s {
  name = %s
  url  = %s
  role = %s
}
`, hclString(label), hclString(d.Name), hclString(d.Url), hclString(d.Role.String()))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

package resources

import (
	"context"
	"encoding/json"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_kea_dhcpv6_reservation",
		Filename: "kea_dhcpv6_reservation.tf",
		Fetch:    fetchKeaDHCPv6Reservation,
	})
}

func fetchKeaDHCPv6Reservation(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/kea/dhcpv6/search_reservation")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid, row := range raw {
		d, err := f.Client().Kea().GetReservationV6(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("kea_dhcpv6_reservation %s: %w", uuid, err)
		}
		if d.Description == "" {
			var rd struct {
				Description string `json:"description"`
			}
			_ = json.Unmarshal(row, &rd)
			d.Description = rd.Description
		}
		label := tracker.Derive(d.Hostname, d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_kea_dhcpv6_reservation" %s {
  subnet_id     = %s
  ip_address    = %s
  duid          = %s
  hostname      = %s
  domain_search = %s
  description   = %s
}
`, hclString(label),
			hclString(d.Subnet.String()),
			hclString(d.IpAddress),
			hclString(d.DUID),
			hclString(d.Hostname),
			hclSet([]string(d.DomainSearch)),
			hclString(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

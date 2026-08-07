package resources

import (
	"context"
	"encoding/json"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_kea_dhcpv4_reservation",
		Filename: "kea_dhcpv4_reservation.tf",
		Fetch:    fetchKeaDHCPv4Reservation,
	})
}

func fetchKeaDHCPv4Reservation(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/kea/dhcpv4/search_reservation")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid, row := range raw {
		d, err := f.Client().Kea().GetReservationV4(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("kea_dhcpv4_reservation %s: %w", uuid, err)
		}
		// get_reservation omits description; fall back to search row
		if d.Description == "" {
			var rd struct {
				Description string `json:"description"`
			}
			_ = json.Unmarshal(row, &rd)
			d.Description = rd.Description
		}
		label := tracker.Derive(d.Hostname, d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_kea_dhcpv4_reservation" %s {
  subnet_id   = %s
  ip_address  = %s
  mac_address = %s
  hostname    = %s
  description = %s
}
`, hclString(label),
			hclString(d.Subnet.String()),
			hclString(d.IpAddress),
			hclString(d.HwAddress),
			hclString(d.Hostname),
			hclString(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

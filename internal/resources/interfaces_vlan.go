package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_interfaces_vlan",
		Filename: "interfaces_vlan.tf",
		Fetch:    fetchInterfacesVLAN,
	})
}

func fetchInterfacesVLAN(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/interfaces/vlan_settings/searchItem")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Interfaces().GetVlan(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("interfaces_vlan %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Device, d.Description, uuid)
		tagVal := "null"
		if d.Tag != "" {
			tagVal = hclInt(stringToInt64(d.Tag))
		}
		hcl := fmt.Sprintf(`resource "opnsense_interfaces_vlan" %s {
  description = %s
  tag         = %s
  priority    = %s
  parent      = %s
  device      = %s
}
`, hclString(label),
			hclStringOrNull(d.Description),
			tagVal,
			hclString(d.Priority.String()),
			hclString(d.Parent.String()),
			hclString(d.Device))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

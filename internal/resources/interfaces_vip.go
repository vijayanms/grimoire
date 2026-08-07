package resources

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/interfaces"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_interfaces_vip",
		Filename: "interfaces_vip.tf",
		Fetch:    fetchInterfacesVIP,
	})
}

type interfacesVIPStruct = interfaces.Vip

func fetchInterfacesVIP(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/interfaces/vip_settings/searchItem")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Interfaces().GetVip(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("interfaces_vip %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_interfaces_vip" %s {
  mode        = %s
  interface   = %s
  network     = %s
  gateway     = %s
  description = %s
}
`, hclString(label),
			hclString(d.Mode.String()),
			hclString(d.Interface.String()),
			hclString(d.Network),
			hclStringOrNull(d.Gateway),
			hclStringOrNull(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

package resources

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/firewall"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_firewall_category",
		Filename: "firewall_category.tf",
		Fetch:    fetchFirewallCategory,
	})
}

type firewallCategoryStruct = firewall.Category

func fetchFirewallCategory(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/firewall/category/searchItem")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Firewall().GetCategory(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("firewall_category %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Name, "", uuid)
		hcl := fmt.Sprintf(`resource "opnsense_firewall_category" %s {
  auto  = %s
  name  = %s
  color = %s
}
`, hclString(label), hclBool(stringToBool(d.Automatic)),
			hclString(d.Name), hclString(d.Color))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

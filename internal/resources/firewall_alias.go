package resources

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/firewall"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_firewall_alias",
		Filename: "firewall_alias.tf",
		Fetch:    fetchFirewallAlias,
	})
}

type firewallAliasStruct = firewall.Alias

func fetchFirewallAlias(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/firewall/alias/searchItem")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Firewall().GetAlias(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("firewall_alias %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Name, d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_firewall_alias" %s {
  enabled         = %s
  name            = %s
  type            = %s
  ip_protocol     = %s
  interface       = %s
  content         = %s
  categories      = %s
  update_freq     = %s
  path_expression = %s
  stats           = %s
  description     = %s
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclString(d.Name),
			hclString(d.Type.String()),
			hclSet([]string(d.IPProtocol)),
			hclString(d.Interface.String()),
			hclSet([]string(d.Content)),
			hclSet([]string(d.Categories)),
			hclFloat(stringToFloat64Default(d.UpdateFreq, -1)),
			hclString(d.PathExpression),
			hclBool(stringToBool(d.Statistics)),
			hclStringOrNull(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

package resources

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/firewall"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_firewall_nat_one_to_one",
		Filename: "firewall_nat_one_to_one.tf",
		Fetch:    fetchFirewallNATOneToOne,
	})
}

type firewallNATOneToOneStruct = firewall.NatOneToOne

func fetchFirewallNATOneToOne(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/firewall/one_to_one/searchRule")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Firewall().GetNatOneToOne(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("firewall_nat_one_to_one %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		seq := "null"
		if d.Sequence != "" {
			seq = hclInt(stringToInt64(d.Sequence))
		}
		natRefl := d.NatReflection.String()
		if natRefl == "" {
			natRefl = "default"
		}
		hcl := fmt.Sprintf(`resource "opnsense_firewall_nat_one_to_one" %s {
  enabled        = %s
  log            = %s
  sequence       = %s
  interface      = %s
  type           = %s
  external_net   = %s
  nat_reflection = %s
  categories     = %s
  description    = %s

  source {
    net    = %s
    invert = %s
  }

  destination {
    net    = %s
    invert = %s
  }
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclBool(stringToBool(d.Log)),
			seq,
			hclString(d.Interface.String()),
			hclString(d.Type.String()),
			hclString(d.ExternalNet),
			hclString(natRefl),
			hclSet([]string(d.Categories)),
			hclStringOrNull(d.Description),
			hclString(d.SourceNet), hclBool(stringToBool(d.SourceInvert)),
			hclString(d.DestinationNet), hclBool(stringToBool(d.DestinationInvert)))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

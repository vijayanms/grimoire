package resources

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/firewall"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_firewall_nat",
		Filename: "firewall_nat.tf",
		Fetch:    fetchFirewallNAT,
	})
}

type firewallNATStruct = firewall.NAT

func fetchFirewallNAT(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/firewall/source_nat/searchRule")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Firewall().GetNAT(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("firewall_nat %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		seq := "null"
		if d.Sequence != "" {
			seq = hclInt(stringToInt64(d.Sequence))
		}
		hcl := fmt.Sprintf(`resource "opnsense_firewall_nat" %s {
  enabled     = %s
  disable_nat = %s
  sequence    = %s
  interface   = %s
  ip_protocol = %s
  protocol    = %s
  log         = %s
  description = %s

  source {
    net    = %s
    port   = %s
    invert = %s
  }

  destination {
    net    = %s
    port   = %s
    invert = %s
  }

  target {
    ip   = %s
    port = %s
  }
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclBool(stringToBool(d.DisableNAT)),
			seq,
			hclString(d.Interface.String()),
			hclString(d.IPProtocol.String()),
			hclString(d.Protocol.String()),
			hclBool(stringToBool(d.Log)),
			hclStringOrNull(d.Description),
			hclString(d.SourceNet), hclStringOrNull(d.SourcePort), hclBool(stringToBool(d.SourceInvert)),
			hclString(d.DestinationNet), hclStringOrNull(d.DestinationPort), hclBool(stringToBool(d.DestinationInvert)),
			hclString(d.Target), hclString(d.TargetPort))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

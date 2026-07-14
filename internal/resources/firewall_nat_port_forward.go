package resources

import (
	"context"
	"fmt"
	"strings"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_firewall_nat_port_forward",
		Filename: "firewall_nat_port_forward.tf",
		Fetch:    fetchFirewallNATPortForward,
	})
}

func natReflectionAPIToSchema(s string) string {
	switch strings.ToLower(s) {
	case "enable":
		return "enable"
	case "purenat":
		return "purenat"
	default:
		return "disable"
	}
}

func fetchFirewallNATPortForward(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/firewall/d_nat/searchRule")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Firewall().GetNatPortForward(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("firewall_nat_port_forward %s: %w", uuid, err)
		}

		label := tracker.Derive("", d.Description, uuid)
		seq := "null"
		if d.Sequence != "" {
			seq = hclInt(stringToInt64(d.Sequence))
		}
		hcl := fmt.Sprintf(`resource "opnsense_firewall_nat_port_forward" %s {
  enabled        = %s
  sequence       = %s
  interface      = %s
  ip_protocol    = %s
  protocol       = %s
  log            = %s
  nat_reflection = %s
  description    = %s

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
			hclBool(!stringToBool(d.Disabled)),
			seq,
			hclSet([]string(d.Interface)),
			hclString(d.IPProtocol.String()),
			hclString(d.Protocol.String()),
			hclBool(stringToBool(d.Log)),
			hclString(natReflectionAPIToSchema(d.NatReflection.String())),
			hclStringOrNull(d.Description),
			hclString(d.Source.Network), hclStringOrNull(d.Source.Port), hclBool(stringToBool(d.Source.Invert)),
			hclString(d.Destination.Network), hclStringOrNull(d.Destination.Port), hclBool(stringToBool(d.Destination.Invert)),
			hclString(d.Target), hclString(d.TargetPort))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

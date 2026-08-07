package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_unbound_forward",
		Filename: "unbound_forward.tf",
		Fetch:    fetchUnboundForward,
	})
}

func fetchUnboundForward(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/unbound/settings/searchDot")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Unbound().GetForward(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("unbound_forward %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Domain, d.Server, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_unbound_forward" %s {
  enabled     = %s
  domain      = %s
  server_ip   = %s
  server_port = %s
  verify_cn   = %s
}
`, hclString(label), hclBool(stringToBool(d.Enabled)),
			hclString(d.Domain), hclString(d.Server),
			hclInt(stringToInt64(d.Port)),
			hclString(d.VerifyCN))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_unbound_host_override",
		Filename: "unbound_host_override.tf",
		Fetch:    fetchUnboundHostOverride,
	})
}

func fetchUnboundHostOverride(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/unbound/settings/searchHostOverride")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Unbound().GetHostOverride(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("unbound_host_override %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Hostname, d.Description, uuid)
		mxprio := stringToInt64(d.MXPriority)
		if mxprio == 0 && d.MXPriority == "" {
			mxprio = -1
		}
		hcl := fmt.Sprintf(`resource "opnsense_unbound_host_override" %s {
  enabled     = %s
  hostname    = %s
  domain      = %s
  type        = %s
  server      = %s
  mx_priority = %s
  mx_host     = %s
  description = %s
}
`, hclString(label), hclBool(stringToBool(d.Enabled)),
			hclString(d.Hostname), hclString(d.Domain),
			hclString(d.Type.String()), hclString(d.Server),
			hclInt(mxprio), hclString(d.MXDomain),
			hclStringOrNull(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_unbound_host_alias",
		Filename: "unbound_host_alias.tf",
		Fetch:    fetchUnboundHostAlias,
	})
}

func fetchUnboundHostAlias(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/unbound/settings/searchHostAlias")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Unbound().GetHostAlias(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("unbound_host_alias %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Hostname, d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_unbound_host_alias" %s {
  enabled     = %s
  override    = %s
  hostname    = %s
  domain      = %s
  description = %s
}
`, hclString(label), hclBool(stringToBool(d.Enabled)),
			hclString(d.Host.String()),
			hclString(d.Hostname), hclString(d.Domain),
			hclStringOrNull(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

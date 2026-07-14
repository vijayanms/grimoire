package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_unbound_domain_override",
		Filename: "unbound_domain_override.tf",
		Fetch:    fetchUnboundDomainOverride,
	})
}

func fetchUnboundDomainOverride(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/unbound/settings/searchDomainOverride")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Unbound().GetDomainOverride(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("unbound_domain_override %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Domain, d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_unbound_domain_override" %s {
  enabled     = %s
  domain      = %s
  server      = %s
  description = %s
}
`, hclString(label), hclBool(stringToBool(d.Enabled)), hclString(d.Domain), hclString(d.Server), hclStringOrNull(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

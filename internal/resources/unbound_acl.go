package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_unbound_acl",
		Filename: "unbound_acl.tf",
		Fetch:    fetchUnboundACL,
	})
}

func fetchUnboundACL(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/unbound/settings/searchAcl")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Unbound().GetAcl(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("unbound_acl %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Name, d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_unbound_acl" %s {
  enabled     = %s
  name        = %s
  action      = %s
  networks    = %s
  description = %s
}
`, hclString(label), hclBool(stringToBool(d.Enabled)),
			hclString(d.Name), hclString(d.Action.String()),
			hclSet([]string(d.Networks)),
			hclStringOrNull(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

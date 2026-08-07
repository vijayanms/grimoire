package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_route",
		Filename: "routes_route.tf",
		Fetch:    fetchRoutesRoute,
	})
}

func fetchRoutesRoute(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/routes/routes/searchroute")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Routes().GetRoute(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("route %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		enabled := !stringToBool(d.Disabled)
		hcl := fmt.Sprintf(`resource "opnsense_route" %s {
  enabled     = %s
  gateway     = %s
  network     = %s
  description = %s
}
`, hclString(label), hclBool(enabled), hclString(d.Gateway.String()), hclString(d.Network), hclStringOrNull(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

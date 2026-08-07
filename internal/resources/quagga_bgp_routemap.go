package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_quagga_bgp_routemap",
		Filename: "quagga_bgp_routemap.tf",
		Fetch:    fetchQuaggaBGPRouteMap,
	})
}

func fetchQuaggaBGPRouteMap(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/quagga/bgp/searchRoutemap")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Quagga().GetBGPRouteMap(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("quagga_bgp_routemap %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Name, d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_quagga_bgp_routemap" %s {
  enabled        = %s
  description    = %s
  name           = %s
  action         = %s
  route_map_id   = %s
  set            = %s
  aspath_list    = %s
  prefix_list    = %s
  community_list = %s
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclString(d.Description),
			hclString(d.Name),
			hclString(d.Action.String()),
			hclInt(stringToInt64(d.RouteMapID)),
			hclString(d.Set),
			hclSet([]string(d.ASPathList)),
			hclSet([]string(d.PrefixList)),
			hclSet([]string(d.CommunityList)))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

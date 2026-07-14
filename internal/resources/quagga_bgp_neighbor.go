package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_quagga_bgp_neighbor",
		Filename: "quagga_bgp_neighbor.tf",
		Fetch:    fetchQuaggaBGPNeighbor,
	})
}

func fetchQuaggaBGPNeighbor(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/quagga/bgp/searchNeighbor")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Quagga().GetBGPNeighbor(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("quagga_bgp_neighbor %s: %w", uuid, err)
		}
		label := tracker.Derive(d.PeerIP, d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_quagga_bgp_neighbor" %s {
  enabled                  = %s
  description              = %s
  peer_ip                  = %s
  remote_as                = %s
  password                 = %s
  weight                   = %s
  local_ip                 = %s
  update_source            = %s
  link_local_interface     = %s
  next_hop_self            = %s
  next_hop_self_all        = %s
  multi_hop                = %s
  multi_protocol           = %s
  rr_client                = %s
  bfd                      = %s
  keepalive                = %s
  hold_down                = %s
  connect_timer            = %s
  default_route            = %s
  as_override              = %s
  disable_connected_check  = %s
  attribute_unchanged      = %s
  prefix_list_in           = %s
  prefix_list_out          = %s
  route_map_in             = %s
  route_map_out            = %s
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclString(d.Description),
			hclString(d.PeerIP),
			hclInt(stringToInt64(d.RemoteAS)),
			hclString(d.Password),
			hclInt(stringToInt64(d.Weight)),
			hclString(d.LocalIP),
			hclString(d.UpdateSource.String()),
			hclString(d.LinkLocalInterface.String()),
			hclBool(stringToBool(d.NextHopSelf)),
			hclBool(stringToBool(d.NextHopSelfAll)),
			hclBool(stringToBool(d.MultiHop)),
			hclBool(stringToBool(d.MultiProtocol)),
			hclBool(stringToBool(d.RRClient)),
			hclBool(stringToBool(d.BFD)),
			hclInt(stringToInt64(d.KeepAlive)),
			hclInt(stringToInt64(d.HoldDown)),
			hclInt(stringToInt64(d.ConnectTimer)),
			hclBool(stringToBool(d.DefaultRoute)),
			hclBool(stringToBool(d.ASOverride)),
			hclBool(stringToBool(d.DisableConnectedCheck)),
			hclString(d.AttributeUnchanged.String()),
			hclString(d.PrefixListIn.String()),
			hclString(d.PrefixListOut.String()),
			hclString(d.RouteMapIn.String()),
			hclString(d.RouteMapOut.String()))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

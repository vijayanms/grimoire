package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_wireguard_server",
		Filename: "wireguard_server.tf",
		Fetch:    fetchWireguardServer,
	})
}

func fetchWireguardServer(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/wireguard/server/searchServer")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Wireguard().GetServer(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("wireguard_server %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Name, "", uuid)
		hcl := fmt.Sprintf(`resource "opnsense_wireguard_server" %s {
  enabled        = %s
  name           = %s
  public_key     = %s
  private_key    = %s
  port           = %s
  mtu            = %s
  dns            = %s
  tunnel_address = %s
  disable_routes = %s
  gateway        = %s
  instance       = %s
  peers          = %s
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclString(d.Name),
			hclString(d.PublicKey),
			hclString(d.PrivateKey),
			hclInt(stringToInt64(d.Port)),
			hclInt(stringToInt64(d.MTU)),
			hclSet([]string(d.DNS)),
			hclSet([]string(d.TunnelAddress)),
			hclBool(stringToBool(d.DisableRoutes)),
			hclString(d.Gateway),
			hclString(d.Instance),
			hclSet([]string(d.Peers)))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_wireguard_client",
		Filename: "wireguard_client.tf",
		Fetch:    fetchWireguardClient,
	})
}

func fetchWireguardClient(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/wireguard/client/searchClient")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Wireguard().GetClient(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("wireguard_client %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Name, "", uuid)
		hcl := fmt.Sprintf(`resource "opnsense_wireguard_client" %s {
  enabled        = %s
  name           = %s
  public_key     = %s
  psk            = %s
  tunnel_address = %s
  server_address = %s
  server_port    = %s
  keepalive      = %s
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclString(d.Name),
			hclString(d.PublicKey),
			hclString(d.PSK),
			hclSet([]string(d.TunnelAddress)),
			hclString(d.ServerAddress),
			hclInt(stringToInt64(d.ServerPort)),
			hclInt(stringToInt64(d.KeepAlive)))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

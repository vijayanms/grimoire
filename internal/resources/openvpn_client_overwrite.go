package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_openvpn_client_overwrite",
		Filename: "openvpn_client_overwrite.tf",
		Fetch:    fetchOpenVPNClientOverwrite,
	})
}

func fetchOpenVPNClientOverwrite(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/openvpn/client_overwrites/search")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Openvpn().GetClientOverwrite(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("openvpn_client_overwrite %s: %w", uuid, err)
		}
		label := tracker.Derive(d.CommonName, d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_openvpn_client_overwrite" %s {
  enabled           = %s
  common_name       = %s
  description       = %s
  servers           = %s
  block             = %s
  push_reset        = %s
  tunnel_network    = %s
  tunnel_networkv6  = %s
  local_network     = %s
  remote_network    = %s
  route_gateway     = %s
  redirect_gateway  = %s
  register_dns      = %s
  dns_domain        = %s
  dns_domain_search = %s
  dns_server        = %s
  ntp_server        = %s
  wins_server       = %s
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclString(d.CommonName),
			hclString(d.Description),
			hclSet([]string(d.Servers)),
			hclBool(stringToBool(d.Block)),
			hclBool(stringToBool(d.PushReset)),
			hclString(d.TunnelNetwork),
			hclString(d.TunnelNetworkV6),
			hclSet([]string(d.LocalNetworks)),
			hclSet([]string(d.RemoteNetworks)),
			hclString(d.RouteGateway),
			hclSet([]string(d.RedirectGateway)),
			hclBool(stringToBool(d.RegisterDNS)),
			hclSet([]string(d.DNSDomain)),
			hclSet([]string(d.DNSDomainSearch)),
			hclSet([]string(d.DNSServers)),
			hclSet([]string(d.NTPServers)),
			hclSet([]string(d.WINSServers)))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

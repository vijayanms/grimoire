package resources

import (
	"context"
	"fmt"
	"strings"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_kea_dhcpv4_subnet",
		Filename: "kea_dhcpv4_subnet.tf",
		Fetch:    fetchKeaDHCPv4Subnet,
	})
}

func fetchKeaDHCPv4Subnet(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/kea/dhcpv4/search_subnet")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Kea().GetSubnetV4(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("kea_dhcpv4_subnet %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Subnet, d.Description, uuid)

		// static routes: "dest,router;dest2,router2"
		staticRoutes := ""
		if d.OptionData.StaticRoutes != "" {
			pairs := strings.Split(d.OptionData.StaticRoutes, ";")
			lines := make([]string, 0, len(pairs))
			for _, p := range pairs {
				parts := strings.SplitN(strings.TrimSpace(p), ",", 2)
				if len(parts) == 2 && parts[0] != "" {
					lines = append(lines, fmt.Sprintf(`
    static_route {
      destination = %s
      router      = %s
    }`, hclString(parts[0]), hclString(parts[1])))
				}
			}
			staticRoutes = strings.Join(lines, "")
		}

		hcl := fmt.Sprintf(`resource "opnsense_kea_dhcpv4_subnet" %s {
  subnet      = %s
  next_server = %s
  pools       = %s
  match_client_id = %s
  option_data_auto_collect = %s
  description = %s

  option_data {
    routers             = %s
    domain_name_servers = %s
    domain_name         = %s
    domain_search       = %s
    ntp_servers         = %s
    time_servers        = %s
    tftp_server_name    = %s
    boot_file_name      = %s%s
  }
}
`, hclString(label),
			hclString(d.Subnet),
			hclString(d.NextServer),
			hclSet(splitNL(d.Pools)),
			hclBool(stringToBool(d.MatchClientId)),
			hclBool(stringToBool(d.OptionDataAutoCollect)),
			hclString(d.Description),
			hclSet([]string(d.OptionData.Routers)),
			hclSet([]string(d.OptionData.DomainNameServers)),
			hclString(d.OptionData.DomainName),
			hclSet([]string(d.OptionData.DomainSearch)),
			hclSet([]string(d.OptionData.NtpServers)),
			hclSet([]string(d.OptionData.TimeServers)),
			hclString(d.OptionData.TftpServerName),
			hclString(d.OptionData.BootFileName),
			staticRoutes)
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

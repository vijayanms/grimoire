package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_dnsmasq_host",
		Filename: "dnsmasq_host.tf",
		Fetch:    fetchDnsmasqHost,
	})
}

func fetchDnsmasqHost(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/dnsmasq/settings/search_host")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Dnsmasq().GetHost(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("dnsmasq_host %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Hostname, d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_dnsmasq_host" %s {
  hostname           = %s
  domain             = %s
  is_local_domain    = %s
  ip_addresses       = %s
  alias_records      = %s
  cname_records      = %s
  client_id          = %s
  hardware_addresses = %s
  tag                = %s
  is_ignored         = %s
  description        = %s
  comment            = %s
}
`, hclString(label), hclString(d.Hostname), hclStringOrNull(d.Domain),
			hclBool(stringToBool(d.IsLocalDomain)), hclSet([]string(d.IpAddresses)),
			hclSet([]string(d.AliasRecords)), hclSet([]string(d.CnameRecords)),
			hclStringOrNull(d.ClientId), hclSet([]string(d.HardwareAddresses)),
			hclStringOrNull(d.Tag.String()), hclBool(stringToBool(d.IsIgnored)),
			hclStringOrNull(d.Description), hclStringOrNull(d.Comments))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

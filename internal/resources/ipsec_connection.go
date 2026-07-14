package resources

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/ipsec"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_ipsec_connection",
		Filename: "ipsec_connection.tf",
		Fetch:    fetchIPsecConnection,
	})
}

type ipsecConnectionStruct = ipsec.IPsecConnection

func fetchIPsecConnection(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/ipsec/connections/searchConnection")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Ipsec().GetIPsecConnection(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("ipsec_connection %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_ipsec_connection" %s {
  enabled                 = %s
  proposals               = %s
  unique                  = %s
  aggressive              = %s
  version                 = %s
  mobike                  = %s
  local_addresses         = %s
  remote_addresses        = %s
  local_port              = %s
  remote_port             = %s
  udp_encapsulation       = %s
  reauthentication_time   = %s
  rekey_time              = %s
  ike_lifetime            = %s
  dpd_delay               = %s
  dpd_timeout             = %s
  ip_pools                = %s
  send_certificate_request = %s
  send_certificate        = %s
  keying_tries            = %s
  description             = %s
}
`, hclString(label),
			hclString(d.Enabled),
			hclSet([]string(d.Proposals)),
			hclString(d.Unique.String()),
			hclString(d.Aggressive),
			hclString(d.Version.String()),
			hclString(d.Mobike),
			hclSet([]string(d.LocalAddresses)),
			hclSet([]string(d.RemoteAddresses)),
			hclString(d.LocalPort.String()),
			hclString(d.RemotePort.String()),
			hclString(d.UDPEncapsulation),
			hclString(d.ReauthenticationTime),
			hclString(d.RekeyTime),
			hclString(d.IKELifetime),
			hclString(d.DPDDelay),
			hclString(d.DPDTimeout),
			hclSet([]string(d.IPPools)),
			hclString(d.SendCertificateRequest),
			hclString(d.SendCertificate.String()),
			hclString(d.KeyingTries),
			hclString(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

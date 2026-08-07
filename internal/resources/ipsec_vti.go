package resources

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/ipsec"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_ipsec_vti",
		Filename: "ipsec_vti.tf",
		Fetch:    fetchIPsecVTI,
	})
}

type ipsecVTIStruct = ipsec.IPsecVTI

func fetchIPsecVTI(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/ipsec/vti/search")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Ipsec().GetIPsecVTI(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("ipsec_vti %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_ipsec_vti" %s {
  enabled           = %s
  request_id        = %s
  local_ip          = %s
  remote_ip         = %s
  tunnel_local_ip   = %s
  tunnel_remote_ip  = %s
  tunnel_local_ip2  = %s
  tunnel_remote_ip2 = %s
  description       = %s
}
`, hclString(label),
			hclString(d.Enabled),
			hclString(d.RequestID),
			hclString(d.LocalIP),
			hclString(d.RemoteIP),
			hclString(d.TunnelLocalIP),
			hclString(d.TunnelRemoteIP),
			hclString(d.TunnelLocalIP2),
			hclString(d.TunnelRemoteIP2),
			hclString(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

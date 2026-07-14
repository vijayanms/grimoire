package resources

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/ipsec"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_ipsec_child",
		Filename: "ipsec_child.tf",
		Fetch:    fetchIPsecChild,
	})
}

type ipsecChildStruct = ipsec.IPsecChild

func fetchIPsecChild(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/ipsec/connections/searchChild")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Ipsec().GetIPsecChild(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("ipsec_child %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_ipsec_child" %s {
  enabled          = %s
  ipsec_connection = %s
  proposals        = %s
  sha256_96        = %s
  start_action     = %s
  close_action     = %s
  dpd_action       = %s
  mode             = %s
  install_policies = %s
  local_networks   = %s
  remote_networks  = %s
  request_id       = %s
  rekey_time       = %s
  description      = %s
}
`, hclString(label),
			hclString(d.Enabled),
			hclString(d.Connection.String()),
			hclSet([]string(d.Proposals)),
			hclString(d.SHA256_96),
			hclString(d.StartAction.String()),
			hclString(d.CloseAction.String()),
			hclString(d.DPDAction.String()),
			hclString(d.Mode.String()),
			hclString(d.InstallPolicies),
			hclSet([]string(d.LocalNetworks)),
			hclSet([]string(d.RemoteNetworks)),
			hclString(d.RequestID),
			hclString(d.RekeyTime),
			hclString(d.Description))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

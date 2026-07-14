package resources

import (
	"context"
	"fmt"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_openvpn_instance",
		Filename: "openvpn_instance.tf",
		Fetch:    fetchOpenVPNInstance,
	})
}

func fetchOpenVPNInstance(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/openvpn/instances/search")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Openvpn().GetInstance(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("openvpn_instance %s: %w", uuid, err)
		}
		label := tracker.Derive(d.Description, "", uuid)
		hcl := fmt.Sprintf(`resource "opnsense_openvpn_instance" %s {
  enabled                = %s
  role                   = %s
  description            = %s
  dev_type               = %s
  protocol               = %s
  port                   = %s
  port_share             = %s
  local                  = %s
  remote                 = %s
  topology               = %s
  server                 = %s
  server_ipv6            = %s
  nopool                 = %s
  bridge_gateway         = %s
  bridge_pool            = %s
  route                  = %s
  push_route             = %s
  push_excluded_routes   = %s
  certificate            = %s
  crl                    = %s
  ca                     = %s
  cert_depth             = %s
  remote_cert_tls        = %s
  verify_client_cert     = %s
  use_ocsp               = %s
  auth_digest            = %s
  data_ciphers           = %s
  data_ciphers_fallback  = %s
  tls_key                = %s
  auth_mode              = %s
  local_group            = %s
  various_flags          = %s
  various_push_flags     = %s
  push_inactive          = %s
  username_as_common_name = %s
  strict_user_cn          = %s
  username               = %s
  password               = %s
  max_clients            = %s
  keepalive_interval     = %s
  keepalive_timeout      = %s
  reneg_sec              = %s
  auth_gen_token         = %s
  auth_gen_token_renewal = %s
  auth_gen_token_secret  = %s
  provision_exclusive    = %s
  redirect_gateway       = %s
  route_metric           = %s
  register_dns           = %s
  dns_domain             = %s
  dns_domain_search      = %s
  dns_servers            = %s
  ntp_servers            = %s
  tun_mtu                = %s
  fragment               = %s
  mssfix                 = %s
  carp_depend_on         = %s
  compress_migrate       = %s
  ifconfig_pool_persist  = %s
  http_proxy             = %s
  verify_x509_name       = %s
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclString(d.Role.String()),
			hclString(d.Description),
			hclString(d.DevType.String()),
			hclString(d.Protocol.String()),
			hclInt(stringToInt64(d.Port)),
			hclString(d.PortShare),
			hclString(d.Local),
			hclSet([]string(d.Remote)),
			hclString(d.Topology.String()),
			hclString(d.Server),
			hclString(d.ServerIPv6),
			hclBool(stringToBool(d.NoPool)),
			hclString(d.BridgeGateway),
			hclString(d.BridgePool),
			hclSet([]string(d.Route)),
			hclSet([]string(d.PushRoute)),
			hclSet([]string(d.PushExcludedRoutes)),
			hclString(d.Certificate.String()),
			hclString(d.CRL.String()),
			hclString(d.CertificateAuthority.String()),
			hclString(d.CertDepth.String()),
			hclBool(stringToBool(d.RemoteCertTLS)),
			hclString(d.VerifyClientCert.String()),
			hclBool(stringToBool(d.UseOCSP)),
			hclString(d.AuthDigest.String()),
			hclSet([]string(d.DataCiphers)),
			hclString(d.DataCiphersFallback.String()),
			hclString(d.TLSKey.String()),
			hclSet([]string(d.AuthMode)),
			hclString(d.LocalGroup.String()),
			hclSet([]string(d.VariousFlags)),
			hclSet([]string(d.VariousPushFlags)),
			hclBool(stringToBool(d.PushInactive)),
			hclBool(stringToBool(d.UsernameAsCommonName)),
			hclString(d.StrictUserCN.String()),
			hclString(d.Username),
			hclString(d.Password),
			hclInt(stringToInt64(d.MaxClients)),
			hclInt(stringToInt64(d.KeepaliveInterval)),
			hclInt(stringToInt64(d.KeepaliveTimeout)),
			hclInt(stringToInt64(d.RenegSec)),
			hclBool(stringToBool(d.AuthGenToken)),
			hclInt(stringToInt64(d.AuthGenTokenRenewal)),
			hclString(d.AuthGenTokenSecret),
			hclBool(stringToBool(d.ProvisionExclusive)),
			hclSet([]string(d.RedirectGateway)),
			hclInt(stringToInt64(d.RouteMetric)),
			hclBool(stringToBool(d.RegisterDNS)),
			hclSet([]string(d.DNSDomain)),
			hclSet([]string(d.DNSDomainSearch)),
			hclSet([]string(d.DNSServers)),
			hclSet([]string(d.NTPServers)),
			hclInt(stringToInt64(d.TunMTU)),
			hclInt(stringToInt64(d.Fragment)),
			hclInt(stringToInt64(d.MSSFix)),
			hclString(d.CARPDependOn.String()),
			hclBool(stringToBool(d.CompressMigrate)),
			hclBool(stringToBool(d.IfConfigPoolPersist)),
			hclString(d.HTTPProxy),
			hclString(d.VerifyX509Name))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

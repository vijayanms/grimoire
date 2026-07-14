package resources

import (
	"context"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/firewall"
)

func init() {
	Registry = append(Registry, ResourceDef{
		TFType:   "opnsense_firewall_filter",
		Filename: "firewall_filter.tf",
		Fetch:    fetchFirewallFilter,
	})
}

type firewallFilterStruct = firewall.Filter

func fetchFirewallFilter(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error) {
	raw, err := f.ListRows(ctx, "/firewall/filter/searchRule")
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(raw))
	for uuid := range raw {
		d, err := f.Client().Firewall().GetFilter(ctx, uuid)
		if err != nil {
			return nil, fmt.Errorf("firewall_filter %s: %w", uuid, err)
		}
		label := tracker.Derive("", d.Description, uuid)
		hcl := fmt.Sprintf(`resource "opnsense_firewall_filter" %s {
  enabled          = %s
  sequence         = %s
  no_xmlrpc_sync   = %s
  description      = %s
  categories       = %s

  interface {
    invert    = %s
    interface = %s
  }

  filter {
    quick         = %s
    action        = %s
    allow_options = %s
    direction     = %s
    ip_protocol   = %s
    protocol      = %s
    log           = %s
    schedule      = %s
    icmp_type     = %s
    tcp_flags        = %s
    tcp_flags_out_of = %s

    source {
      net    = %s
      port   = %s
      invert = %s
    }

    destination {
      net    = %s
      port   = %s
      invert = %s
    }
  }

  stateful_firewall {
    type           = %s
    timeout        = %s
    overload_table = %s
    no_pfsync      = %s

    adaptive_timeouts {
      start = %s
      end   = %s
    }

    max {
      states             = %s
      source_nodes       = %s
      source_states      = %s
      source_connections  = %s
      new_connections {
        count   = %s
        seconds = %s
      }
    }
  }

  traffic_shaping {
    shaper         = %s
    reverse_shaper = %s
  }

  source_routing {
    gateway          = %s
    disable_reply_to = %s
    reply_to         = %s
  }

  priority {
    match         = %s
    set           = %s
    low_delay_set = %s
    match_tos     = %s
  }

  internal_tagging {
    set_local   = %s
    match_local = %s
  }
}
`, hclString(label),
			hclBool(stringToBool(d.Enabled)),
			hclInt(stringToInt64(d.Sequence)),
			hclBool(stringToBool(d.NoXMLRPCSync)),
			hclStringOrNull(d.Description),
			hclSet([]string(d.Categories)),
			hclBool(stringToBool(d.InvertInterface)),
			hclSet([]string(d.Interface)),
			hclBool(stringToBool(d.Quick)),
			hclString(d.Action.String()),
			hclBool(stringToBool(d.AllowOptions)),
			hclString(d.Direction.String()),
			hclString(d.IPProtocol.String()),
			hclString(d.Protocol.String()),
			hclBool(stringToBool(d.Log)),
			hclStringOrNull(d.Schedule.String()),
			hclSet([]string(d.ICMPType)),
			hclSet([]string(d.TCPFlags)),
			hclSet([]string(d.TCPFlagsOutOf)),
			hclString(d.SourceNet),
			hclStringOrNull(d.SourcePort),
			hclBool(stringToBool(d.SourceInvert)),
			hclString(d.DestinationNet),
			hclStringOrNull(d.DestinationPort),
			hclBool(stringToBool(d.DestinationInvert)),
			hclString(d.StateType.String()),
			hclInt(stringToInt64(d.StateTimeout)),
			hclStringOrNull(d.OverloadTable.String()),
			hclBool(stringToBool(d.NoPfsync)),
			hclInt(stringToInt64(d.AdaptiveTimeoutsStart)),
			hclInt(stringToInt64(d.AdaptiveTimeoutsEnd)),
			hclInt(stringToInt64(d.MaxStates)),
			hclInt(stringToInt64(d.MaxSourceNodes)),
			hclInt(stringToInt64(d.MaxSourceStates)),
			hclInt(stringToInt64(d.MaxSourceConnections)),
			hclInt(stringToInt64(d.MaxNewConnectionsCount)),
			hclInt(stringToInt64(d.MaxNewConnectionsSeconds)),
			hclStringOrNull(d.TrafficShaper.String()),
			hclStringOrNull(d.TrafficShaperReverse.String()),
			hclStringOrNull(d.Gateway.String()),
			hclBool(stringToBool(d.DisableReplyTo)),
			hclStringOrNull(d.ReplyTo.String()),
			hclInt(stringToInt64(d.MatchPriority.String())),
			hclInt(stringToInt64(d.SetPriority.String())),
			hclInt(stringToInt64(d.SetPriorityLowDelay.String())),
			hclStringOrNull(d.MatchTOS.String()),
			hclStringOrNull(d.SetLocalTag),
			hclStringOrNull(d.MatchLocalTag))
		entries = append(entries, Entry{UUID: uuid, Label: label, HCL: hcl})
	}
	return entries, nil
}

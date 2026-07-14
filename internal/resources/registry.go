package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/browningluke/opnsense-go/pkg/opnsense"
)

// Entry is one rendered resource block.
type Entry struct {
	UUID  string
	Label string
	HCL   string
}

// ResourceDef describes one resource type.
type ResourceDef struct {
	TFType   string // e.g. "opnsense_firewall_alias"
	Filename string // e.g. "firewall_alias.tf"
	Fetch    func(ctx context.Context, f Fetcher, tracker LabelTracker) ([]Entry, error)
}

// Fetcher is the interface our HTTP fetcher satisfies (internal/fetcher, passed here as interface).
type Fetcher interface {
	Client() opnsense.Client
	ListUnderKey(ctx context.Context, endpoint, monad string) (map[string]json.RawMessage, error)
	ListRows(ctx context.Context, endpoint string) (map[string]json.RawMessage, error)
}

// LabelTracker deduplicates labels within a single resource type.
type LabelTracker interface {
	Derive(name, description, uuid string) string
}

// Registry holds all resource definitions, populated by init() in each resource file.
var Registry []ResourceDef

// ---- shared conversion helpers ----

func stringToBool(s string) bool   { return s == "1" }
func stringToInt64(s string) int64 { n, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64); return n }
func stringToFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

func hclBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func hclString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}

func hclStringOrNull(s string) string {
	if s == "" {
		return "null"
	}
	return hclString(s)
}

func hclSet(items []string) string {
	filtered := make([]string, 0, len(items))
	for _, v := range items {
		if v != "" {
			filtered = append(filtered, v)
		}
	}
	sort.Strings(filtered)
	parts := make([]string, len(filtered))
	for i, v := range filtered {
		parts[i] = hclString(v)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func hclSetUnsorted(items []string) string {
	filtered := make([]string, 0, len(items))
	for _, v := range items {
		if v != "" {
			filtered = append(filtered, v)
		}
	}
	parts := make([]string, len(filtered))
	for i, v := range filtered {
		parts[i] = hclString(v)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func hclInt(n int64) string { return fmt.Sprintf("%d", n) }
func hclFloat(f float64) string {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return s
}

// splitCSV splits a comma-delimited string, filtering empties.
func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := parts[:0]
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// splitNL splits a newline-delimited string, filtering empties.
func splitNL(s string) []string {
	parts := strings.Split(s, "\n")
	out := parts[:0]
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

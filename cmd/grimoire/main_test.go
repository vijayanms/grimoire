package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteImportsTF(t *testing.T) {
	dir := t.TempDir()
	blocks := []string{
		"import {\n  to = opnsense_firewall_alias.lan_net\n  id = \"__lan_network\"\n}\n",
		"import {\n  to = opnsense_firewall_alias.jellyfin\n  id = \"c8e9f400-8d61-4728-bdc8-b2b5034fe4a9\"\n}\n",
	}

	if err := writeImportsTF(dir, blocks); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dir, "imports.tf"))
	if err != nil {
		t.Fatalf("failed to read imports.tf: %v", err)
	}

	for _, want := range []string{
		`to = opnsense_firewall_alias.lan_net`,
		`id = "__lan_network"`,
		`to = opnsense_firewall_alias.jellyfin`,
		`id = "c8e9f400-8d61-4728-bdc8-b2b5034fe4a9"`,
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("imports.tf missing %q\ngot:\n%s", want, got)
		}
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vijayanms/grimoire/internal/fetcher"
	"github.com/vijayanms/grimoire/internal/label"
	"github.com/vijayanms/grimoire/internal/resources"
)

func main() {
	uri := flag.String("uri", os.Getenv("OPNSENSE_URI"), "OPNsense base URI (e.g. https://opnsense.local)")
	apiKey := flag.String("api-key", os.Getenv("OPNSENSE_API_KEY"), "OPNsense API key")
	apiSecret := flag.String("api-secret", os.Getenv("OPNSENSE_API_SECRET"), "OPNsense API secret")
	insecure := flag.Bool("insecure", false, "Skip TLS certificate verification")
	outDir := flag.String("out-dir", "./out", "Output directory for generated files")
	resourceFilter := flag.String("resources", "", "Comma-separated list of TF resource types to generate (default: all)")
	flag.Parse()

	if *uri == "" || *apiKey == "" || *apiSecret == "" {
		fmt.Fprintln(os.Stderr, "error: --uri, --api-key, and --api-secret are required")
		flag.Usage()
		os.Exit(1)
	}

	filter := map[string]bool{}
	if *resourceFilter != "" {
		for _, r := range strings.Split(*resourceFilter, ",") {
			filter[strings.TrimSpace(r)] = true
		}
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating output directory: %v\n", err)
		os.Exit(1)
	}

	logFile, err := os.Create(filepath.Join(*outDir, "opnsense.log"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags|log.LUTC)

	f := fetcher.New(*uri, *apiKey, *apiSecret, *insecure, logger)
	ctx := context.Background()

	var importBlocks []string

	for _, def := range resources.Registry {
		if len(filter) > 0 && !filter[def.TFType] {
			continue
		}

		tracker := label.New()
		fmt.Printf("fetching %s...\n", def.TFType)

		entries, err := def.Fetch(ctx, f, tracker)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: %s: %v\n", def.TFType, err)
			continue
		}
		if len(entries) == 0 {
			continue
		}

		var sb strings.Builder
		for _, e := range entries {
			sb.WriteString(e.HCL)
			sb.WriteString("\n")
			importBlocks = append(importBlocks, fmt.Sprintf("import {\n  to = %s.%s\n  id = %q\n}\n", def.TFType, e.Label, e.UUID))
		}

		tfPath := filepath.Join(*outDir, def.Filename)
		if err := os.WriteFile(tfPath, []byte(sb.String()), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", tfPath, err)
			os.Exit(1)
		}
		fmt.Printf("  wrote %d resources to %s\n", len(entries), def.Filename)
	}

	if err := writeProviderTF(*outDir, *uri, *apiKey, *apiSecret); err != nil {
		fmt.Fprintf(os.Stderr, "error writing provider.tf: %v\n", err)
		os.Exit(1)
	}

	if err := writeImportsTF(*outDir, importBlocks); err != nil {
		fmt.Fprintf(os.Stderr, "error writing imports.tf: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\ndone. %d import blocks written to %s/imports.tf\n", len(importBlocks), *outDir)
}

const providerTF = `terraform {
  required_providers {
    opnsense = {
      source = "browningluke/opnsense"
    }
  }
}

provider "opnsense" {
  uri        = var.opnsense_uri
  api_key    = var.opnsense_api_key
  api_secret = var.opnsense_api_secret
}
`

const variablesTF = `variable "opnsense_uri" {
  description = "OPNsense base URI"
  type        = string
}

variable "opnsense_api_key" {
  description = "OPNsense API key"
  type        = string
  sensitive   = true
}

variable "opnsense_api_secret" {
  description = "OPNsense API secret"
  type        = string
  sensitive   = true
}
`

func writeProviderTF(outDir, uri, apiKey, apiSecret string) error {
	if err := os.WriteFile(filepath.Join(outDir, "provider.tf"), []byte(providerTF), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outDir, "variables.tf"), []byte(variablesTF), 0o644); err != nil {
		return err
	}
	tfvars := fmt.Sprintf("opnsense_uri = %q\nopnsense_api_key = %q\nopnsense_api_secret = %q\n", uri, apiKey, apiSecret)
	return os.WriteFile(filepath.Join(outDir, "terraform.tfvars"), []byte(tfvars), 0o644)
}

func writeImportsTF(outDir string, blocks []string) error {
	path := filepath.Join(outDir, "imports.tf")
	return os.WriteFile(path, []byte(strings.Join(blocks, "\n")), 0o644)
}

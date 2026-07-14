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

	var importLines []string

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
			importLines = append(importLines, fmt.Sprintf("tofu import '%s.%s' '%s'", def.TFType, e.Label, e.UUID))
		}

		tfPath := filepath.Join(*outDir, def.Filename)
		if err := os.WriteFile(tfPath, []byte(sb.String()), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing %s: %v\n", tfPath, err)
			os.Exit(1)
		}
		fmt.Printf("  wrote %d resources to %s\n", len(entries), def.Filename)
	}

	if err := writeProviderTF(*outDir, *uri); err != nil {
		fmt.Fprintf(os.Stderr, "error writing provider.tf: %v\n", err)
		os.Exit(1)
	}

	if err := writeImportSh(*outDir, importLines); err != nil {
		fmt.Fprintf(os.Stderr, "error writing import.sh: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\ndone. %d import commands written to %s/import.sh\n", len(importLines), *outDir)
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
  api_key    = "YOUR_API_KEY"
  api_secret = "YOUR_API_SECRET"
}
`

const variablesTF = `variable "opnsense_uri" {
  description = "OPNsense base URI"
  type        = string
}
`

func writeProviderTF(outDir string, uri string) error {
	if err := os.WriteFile(filepath.Join(outDir, "provider.tf"), []byte(providerTF), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outDir, "variables.tf"), []byte(variablesTF), 0o644); err != nil {
		return err
	}
	tfvars := fmt.Sprintf("opnsense_uri = %q\n", uri)
	return os.WriteFile(filepath.Join(outDir, "terraform.tfvars"), []byte(tfvars), 0o644)
}

func writeImportSh(outDir string, lines []string) error {
	var sb strings.Builder
	sb.WriteString("#!/bin/sh\nset -e\n\n")
	for _, l := range lines {
		sb.WriteString(l)
		sb.WriteString("\n")
	}
	path := filepath.Join(outDir, "import.sh")
	if err := os.WriteFile(path, []byte(sb.String()), 0o755); err != nil {
		return err
	}
	return nil
}

# grimoire — Claude Code Instructions

## Communication Style
Respond like a caveman. No articles, no filler words, no pleasantries. Short. Direct. Code speaks for itself.

## Build & test

```bash
task build     # compile
task check     # fmt + vet + test
go test ./...  # tests only
```

`go build ./...` must pass before any commit. Run it after every file change.

## Architecture

See `AGENTS.md` for full architecture. Key points:

- `internal/resources/registry.go` — shared helpers and interfaces; never import outside the `resources` package
- Each resource file uses `init()` for self-registration; adding a resource = new file only
- `internal/fetcher/fetcher.go` is `package fetcher`; the `resources` package depends on the `Fetcher` interface, not the concrete type
- entry point is `cmd/grimoire/main.go`
- HCL is rendered with `fmt.Sprintf` + helpers — no external HCL library

## Adding a resource

HCL attribute names must match the provider's schema **exactly** — a mismatch causes `tofu plan` to show spurious diffs or fail outright.

Before writing HCL output for a resource, check the canonical attribute names in the provider's source on GitHub:

`https://github.com/browningluke/terraform-provider-opnsense/tree/master/internal/service/<service>/<resource>_schema.go`

The `Schema` map keys in that file are the exact attribute names to use.

## Conventions

- No comments unless the WHY is non-obvious
- No opportunistic cleanup beyond the task scope
- Minimal diffs — don't parametrize or abstract unless the task requires it
- `hclStringOrNull` for optional string fields; `hclInt` for numeric; `hclSet` for sets (filters empty strings, sorts)

## opnsense-go

Used for struct types only (`api.SelectedMap`, `api.SelectedMapList`, etc.). Do not modify it. `go.mod` always points at the published module — never commit a local `replace` directive; if you need to test against a local checkout, use an uncommitted `go.work`.

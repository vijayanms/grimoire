# Agent Instructions

## Build, Test, and Lint

```bash
# Build binary to dist/
task build
# or
go build -o dist/opnsense-tf-import ./cmd/grimoire

# Run unit tests
task test
# or
go test ./...

# Format code
task fmt

# Format + vet + test
task check

# Run with arguments
task run -- --uri https://opnsense.example.com --api-key K --api-secret S --insecure
```

## Architecture

Standalone CLI tool. No Terraform plugin framework. Three layers:

```
cmd/grimoire/main.go        # flag parsing, orchestration loop, writes output files
internal/fetcher/fetcher.go # HTTP Basic auth client; UUID-preserving list calls
internal/label/label.go     # label derivation and collision tracking
internal/resources/
├── registry.go             # ResourceDef, Entry, Fetcher/LabelTracker interfaces, shared HCL helpers
└── <service>_<name>.go     # one file per resource type; registers via init()
```

### Registration pattern

Each resource file calls `init()` to append a `ResourceDef` to `resources.Registry`. No manual registration in main. Adding a new resource = new file only.

```go
func init() {
    Registry = append(Registry, ResourceDef{
        TFType:   "opnsense_firewall_alias",
        Filename: "firewall_alias.tf",
        Fetch:    fetchFirewallAlias,
    })
}
```

### Fetcher

`internal/fetcher/fetcher.go` (`package fetcher`) holds a `*Fetcher` that makes raw `net/http` GET requests with Basic auth. Satisfies the `resources.Fetcher` interface:

```go
ListUnderKey(ctx context.Context, endpoint, monad string) (map[string]json.RawMessage, error)
```

The monad key is the wrapper object key in the OPNsense API response (e.g. `"alias"` in `{"alias": {"uuid": {...}}}`).

### Adding a new resource

1. Create `internal/resources/<service>_<name>.go`
2. Define a struct with JSON tags matching the opnsense-go struct for that resource
3. Implement `fetch<Name>(ctx, Fetcher, LabelTracker) ([]Entry, error)`:
   - Call `f.ListUnderKey(ctx, "<endpoint>", "<monad>")`
   - Unmarshal each entry
   - Derive a label via `tracker.Derive(name, description, uuid)`
   - Render HCL via `fmt.Sprintf` using the shared helpers
4. Register in `init()`
5. `go build ./...` — done

Match the API endpoint and monad to what opnsense-go uses in its `ReqOpts` for that resource type.

## Key Conventions

### Field conversion helpers (in `internal/resources/registry.go`)

| OPNsense API type | TF type | Helper |
|---|---|---|
| `"0"` / `"1"` → bool | bool | `hclBool(stringToBool(s))` |
| Disabled (inverted) | bool | `hclBool(!stringToBool(s))` |
| `api.SelectedMap` → string | string | `hclString(m.String())` |
| `api.SelectedMapList` → set | set | `hclSet([]string(m))` |
| `api.SelectedMapListNL` → set | set | `hclSet([]string(m))` |
| String int → int64 | number | `hclInt(stringToInt64(s))` |
| Optional string (empty = null) | string/null | `hclStringOrNull(s)` |
| Newline-delimited string → slice | set | `hclSet(splitNL(s))` |

### Label derivation

`tracker.Derive(name, description, uuid)` — priority: name → description (first 20 chars) → uuid prefix (8 chars). Sanitized to `[a-z0-9_]`, max 63 chars. Collisions get `_N` suffix. Use a fresh `label.New()` per resource type in the main loop.

### HCL output

No HCL library. Pure `fmt.Sprintf` with the helper functions above. One resource block per entry, blank line between blocks. File named by `ResourceDef.Filename`.

### opnsense-go dependency

`github.com/browningluke/opnsense-go` is the API client library. It provides the struct types used here (`api.SelectedMap`, `api.SelectedMapList`, etc.). Do not modify it in this repo.

If a new resource needs a struct not yet in opnsense-go, define a local struct in the resource file with matching JSON tags — that's sufficient and adds no new dependency.

**Do not import `github.com/browningluke/terraform-provider-opnsense/...`**. The provider's internal packages are not a public API and importing them would couple this tool to the provider's internal structure.

## Keeping This File Up-to-Date

Update `AGENTS.md` when:
- A new resource type is added (if it introduces a non-obvious pattern)
- A new shared helper is added to `registry.go`
- Build or test commands change
- A new convention is established

# Contributing

## Setup

```sh
git clone https://github.com/vijayanms/grimoire
cd grimoire
go mod download
```

The module depends on [browningluke/opnsense-go](https://github.com/browningluke/opnsense-go) as its API struct source. For most contributions you don't need to touch it — the published module version is used automatically via Go modules.

If a contribution requires adding or changing a struct in `opnsense-go`, that change must be merged and released upstream first. Open a PR there, wait for a release, then update the version in `go.mod` here.

## Build & run

```sh
task build                          # compile to dist/
task run -- --uri https://... ...   # go run with args
task check                          # fmt + vet + test
```

## Adding a resource type

1. Create `internal/resources/<service>_<name>.go` following the existing pattern:
   - `init()` registers a `ResourceDef` in `resources.Registry`
   - Define a struct mirroring the opnsense-go type (same JSON tags)
   - Implement `fetch<Name>(ctx, Fetcher, LabelTracker) ([]Entry, error)`
   - Use the shared helpers in `internal/resources/registry.go` for field conversion

2. Run `go build ./...` — no other registration step needed.

## Field conversion quick reference

| OPNsense API type | Go field | HCL helper |
|---|---|---|
| `"0"` / `"1"` → bool | `string` | `hclBool(stringToBool(s))` |
| Disabled (inverted) | `string` | `hclBool(!stringToBool(s))` |
| `api.SelectedMap` → string | `api.SelectedMap` | `hclString(m.String())` |
| `api.SelectedMapList` → set | `api.SelectedMapList` | `hclSet([]string(m))` |
| String int → int64 | `string` | `hclInt(stringToInt64(s))` |
| Optional string | `string` | `hclStringOrNull(s)` |

## Pull requests

**Title format:** `<resource or area> — brief description`  
Example: `openvpn_instance — fix port field type`

**PR template:**

```
## What

Short description of the change.

## Why

Motivation — bug, missing field, wrong type, etc.

## Type

- [ ] Bug fix
- [ ] New resource type
- [ ] Breaking change
- [ ] Docs / tooling

## Testing

- [ ] `go build ./...` passes
- [ ] `task check` passes
- [ ] Manually verified against OPNsense (version: )
```

Keep PRs focused — one resource type or one fix per PR.

## Code style

- `gofmt` — run `task fmt` before committing
- No comments unless the WHY is non-obvious
- Match field naming in the generated HCL to the provider's schema attribute names exactly

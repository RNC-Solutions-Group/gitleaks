# Dependency matrix — full upgrade experiment

**Status:** Experimental result; not a production dependency decision
**Canonical application SHA:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`
**Experiment branch:** `test/deps-upgrade-all`
**Inventory basis:** `go list -m -u -json all` at the experiment start, then `go list -m -json all` after `go get -u -t ./...`, `go mod tidy`, and `go mod verify`.
**Toolchain used for resolution:** Go 1.27.3 linux/amd64; final module `go` directive is **1.26.0**.

The candidate column is the latest version visible to `go list -m -u` at experiment time. “Selected” is the version resolved in this experiment after module-graph pruning. A selected version can be older than a candidate when the parent module pins it or when the update would require an import-path change. `go mod tidy -diff` was clean and `go mod verify` reported `all modules verified`.

## Summary

| Metric | Baseline | Upgraded |
| --- | ---: | ---: |
| Modules in graph (excluding main) | 186 | 143 |
| Direct modules | 16 | 16 |
| Indirect modules | 170 | 127 |
| Updated modules | — | 69 (9 direct, 60 indirect) |
| Modules added to graph | — | 34 |
| Modules removed from graph | — | 77 |
| Unchanged modules | — | 40 |
| `go.mod` Go directive | 1.24.11 | 1.26.0 |
| `go.sum` lines | 423 | 172 |

The module count fell mainly because `go-re2` v1.13 no longer requires `wazero`, and because Go 1.26 graph pruning removed old test-only dependencies. The add/remove churn is not all runtime code: many added modules (cloud/OpenTelemetry, `golang.org/x/tools` test packages) are graph-only and are not linked into the scanner binary.

## Direct dependencies

| Module | Current | Candidate | Selected | Class | Compatibility | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| `github.com/BobuSumisu/aho-corasick` | v1.0.3 | v1.0.3 | v1.0.3 | Direct | Supported |  |
| `github.com/Masterminds/sprig/v3` | v3.3.0 | v3.3.0 | v3.3.0 | Direct | Supported |  |
| `github.com/charmbracelet/lipgloss` | v0.5.0 | v1.1.0 | v1.1.0 | Direct | Supported | Major v0→v1; simple API usage still compiles |
| `github.com/fatih/semgroup` | v1.2.0 | v1.3.0 | v1.3.0 | Direct | Supported |  |
| `github.com/gitleaks/go-gitdiff` | v0.9.1 | v0.9.1 | v0.9.1 | Direct | Supported |  |
| `github.com/google/go-cmp` | v0.7.0 | v0.7.0 | v0.7.0 | Direct | Supported |  |
| `github.com/h2non/filetype` | v1.1.3 | v1.1.3 | v1.1.3 | Direct | Supported |  |
| `github.com/hashicorp/go-version` | v1.7.0 | v1.9.0 | v1.9.0 | Direct | Supported |  |
| `github.com/lucasjones/reggen` | v0.0.0-20200904144131-37ba4fa293bb | v0.0.0-20200904144131-37ba4fa293bb | v0.0.0-20200904144131-37ba4fa293bb | Direct | Supported |  |
| `github.com/mholt/archives` | v0.1.5 | v0.1.5 | v0.1.5 | Direct | Supported |  |
| `github.com/rs/zerolog` | v1.33.0 | v1.35.1 | v1.35.1 | Direct | Supported |  |
| `github.com/spf13/cobra` | v1.9.1 | v1.10.2 | v1.10.2 | Direct | Supported |  |
| `github.com/spf13/viper` | v1.19.0 | v1.21.0 | v1.21.0 | Direct | Supported | Uses mapstructure/v2 and go.yaml.in/yaml/v3; precedence tested |
| `github.com/stretchr/testify` | v1.10.0 | v1.12.1 | v1.12.1 | Direct | Supported |  |
| `github.com/wasilibs/go-re2` | v1.9.0 | v1.13.0 | v1.13.0 | Direct | Supported (Go >=1.26) | wasm2go replaces wazero; **requires Go 1.26** |
| `golang.org/x/exp` | v0.0.0-20250218142911-aa4b98e5adaa | v0.0.0-20261109185158-a52e3aef5920 | v0.0.0-20261109185158-a52e3aef5920 | Direct | Supported | Latest timestamped pseudo-version |

## Regex and WebAssembly dependencies

| Module | Current | Candidate | Selected | Class | Compatibility | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| `github.com/wasilibs/go-re2` | v1.9.0 | v1.13.0 | v1.13.0 | Regex / Wasm | Supported (Go >=1.26) | wasm2go replaces wazero; **requires Go 1.26** |
| `github.com/tetratelabs/wazero` | v1.9.0 | v1.12.0 | removed | Regex / Wasm | Removed | Removed by graph pruning/tidy |
| `github.com/wasilibs/wazero-helpers` | v0.0.0-20240620070341-3dff1577cd52 | v0.0.0-20250123031827-cd30c44769bb | removed | Regex / Wasm | Removed | Removed by graph pruning/tidy |

`go-re2` v1.13.0 is the substantive regex-engine change. It moved from the wazero WebAssembly runtime to `wasm2go`, removes `wazero` and `wazero-helpers` from the graph, and requires **Go 1.26.0**. The public `regexp` API used by Gitleaks did not require a source change. The optional `re2_cgo` variant remains a separate native build and requires system `re2` plus `pkg-config`; it is not selected by `CGO_ENABLED=1` alone.

## Updated indirect dependencies

| Module | Current | Candidate | Selected | Class | Compatibility | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| `cloud.google.com/go` | v0.112.1 | v0.123.0 | v0.121.6 | Indirect | Uncertain (graph-only) |  |
| `cloud.google.com/go/compute/metadata` | v0.2.3 | v0.10.0 | v0.9.0 | Indirect | Uncertain (graph-only) |  |
| `cloud.google.com/go/storage` | v1.5.0 | v1.70.0 | v1.56.0 | Indirect | Uncertain (graph-only) | Graph-only via go4.org; not linked |
| `dario.cat/mergo` | v1.0.1 | v1.0.2 | v1.0.2 | Indirect | Supported |  |
| `github.com/andybalholm/brotli` | v1.2.0 | v1.2.8 | v1.2.8 | Indirect | Supported | Fixes reachable GO-2026-6888 |
| `github.com/bodgit/sevenzip` | v1.6.1 | v1.6.5 | v1.6.5 | Indirect | Supported |  |
| `github.com/coreos/go-systemd/v22` | v22.5.0 | v22.7.0 | v22.7.0 | Indirect | Supported |  |
| `github.com/davecgh/go-spew` | v1.1.2-0.20180830191138-d8f796af33cc | v1.1.2-0.20180830191138-d8f796af33cc | v1.1.1 | Indirect | Supported |  |
| `github.com/envoyproxy/protoc-gen-validate` | v0.1.0 | v1.3.3 | v1.2.1 | Indirect | Supported |  |
| `github.com/fsnotify/fsnotify` | v1.8.0 | v1.10.1 | v1.10.1 | Indirect | Supported |  |
| `github.com/go-logr/logr` | v1.4.1 | v1.4.4 | v1.4.3 | Indirect | Supported |  |
| `github.com/google/s2a-go` | v0.1.7 | v0.1.11 | v0.1.9 | Indirect | Supported |  |
| `github.com/googleapis/enterprise-certificate-proxy` | v0.3.2 | v0.3.25 | v0.3.7 | Indirect | Supported |  |
| `github.com/googleapis/gax-go/v2` | v2.12.3 | v2.27.0 | v2.16.0 | Indirect | Supported |  |
| `github.com/huandu/xstrings` | v1.5.0 | v1.6.3 | v1.6.3 | Indirect | Supported |  |
| `github.com/klauspost/compress` | v1.18.0 | v1.20.4 | v1.20.4 | Indirect | Supported | Fixes reachable GO-2026-6865 |
| `github.com/klauspost/pgzip` | v1.2.6 | v1.3.0 | v1.3.0 | Indirect | Supported |  |
| `github.com/lucasb-eyer/go-colorful` | v1.2.0 | v1.4.1 | v1.4.1 | Indirect | Supported |  |
| `github.com/Masterminds/semver/v3` | v3.3.0 | v3.5.0 | v3.5.0 | Indirect | Supported |  |
| `github.com/mattn/go-colorable` | v0.1.14 | v0.1.16 | v0.1.16 | Indirect | Supported |  |
| `github.com/mattn/go-isatty` | v0.0.20 | v0.0.24 | v0.0.24 | Indirect | Supported |  |
| `github.com/mattn/go-runewidth` | v0.0.14 | v0.0.31 | v0.0.31 | Indirect | Supported |  |
| `github.com/minio/minlz` | v1.0.1 | v1.2.4 | v1.2.4 | Indirect | Supported |  |
| `github.com/muesli/termenv` | v0.15.1 | v0.16.0 | v0.16.0 | Indirect | Supported |  |
| `github.com/nwaples/rardecode/v2` | v2.2.0 | v2.4.6 | v2.4.6 | Indirect | Supported |  |
| `github.com/pelletier/go-toml/v2` | v2.2.3 | v2.4.3 | v2.4.3 | Indirect | Supported |  |
| `github.com/pierrec/lz4/v4` | v4.1.22 | v4.1.33 | v4.1.33 | Indirect | Supported |  |
| `github.com/pmezard/go-difflib` | v1.0.1-0.20181226105442-5d4384ee4fb2 | v1.0.1-0.20181226105442-5d4384ee4fb2 | v1.0.0 | Indirect | Supported |  |
| `github.com/rivo/uniseg` | v0.2.0 | v0.4.7 | v0.4.7 | Indirect | Supported |  |
| `github.com/rs/xid` | v1.5.0 | v1.6.0 | v1.6.0 | Indirect | Supported |  |
| `github.com/sagikazarmark/locafero` | v0.7.0 | v0.12.0 | v0.12.0 | Indirect | Supported |  |
| `github.com/shopspring/decimal` | v1.4.0 | v1.6.0 | v1.6.0 | Indirect | Supported |  |
| `github.com/sourcegraph/conc` | v0.3.0 | v0.3.0 | v0.3.1-0.20240121214520-5f936abd7ae8 | Indirect | Supported |  |
| `github.com/spf13/cast` | v1.7.1 | v1.10.0 | v1.10.0 | Indirect | Supported |  |
| `github.com/spf13/pflag` | v1.0.6 | v1.1.0 | v1.1.0 | Indirect | Supported |  |
| `github.com/stretchr/objx` | v0.5.2 | v0.8.0 | v0.5.3 | Indirect | Supported |  |
| `github.com/ulikunitz/xz` | v0.5.15 | v0.5.17 | v0.5.17 | Indirect | Supported |  |
| `go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc` | v0.49.0 | v0.72.0 | v0.61.0 | Indirect | Uncertain (graph-only) |  |
| `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp` | v0.49.0 | v0.72.0 | v0.61.0 | Indirect | Uncertain (graph-only) |  |
| `go.opentelemetry.io/otel` | v1.24.0 | v1.47.0 | v1.38.0 | Indirect | Uncertain (graph-only) |  |
| `go.opentelemetry.io/otel/metric` | v1.24.0 | v1.47.0 | v1.38.0 | Indirect | Uncertain (graph-only) |  |
| `go.opentelemetry.io/otel/trace` | v1.24.0 | v1.47.0 | v1.38.0 | Indirect | Uncertain (graph-only) |  |
| `go4.org` | v0.0.0-20230225012048-214862532bf5 | v0.0.0-20261028031214-ed6929a230fe | v0.0.0-20261028031214-ed6929a230fe | Indirect | Supported | Latest pseudo-version pulls cloud/OpenTelemetry graph-only modules |
| `golang.org/x/crypto` | v0.35.0 | v0.58.0 | v0.58.0 | Indirect | Supported | Module-level GO-2026-5932 (openpgp), not imported |
| `golang.org/x/mod` | v0.27.0 | v0.43.0 | v0.42.0 | Indirect | Supported |  |
| `golang.org/x/net` | v0.23.0 | v0.62.0 | v0.60.0 | Indirect | Supported |  |
| `golang.org/x/oauth2` | v0.18.0 | v0.38.0 | v0.34.0 | Indirect | Supported |  |
| `golang.org/x/sync` | v0.17.0 | v0.25.0 | v0.25.0 | Indirect | Supported |  |
| `golang.org/x/sys` | v0.30.0 | v0.50.0 | v0.50.0 | Indirect | Supported |  |
| `golang.org/x/term` | v0.29.0 | v0.48.0 | v0.47.0 | Indirect | Supported |  |
| `golang.org/x/text` | v0.29.0 | v0.44.0 | v0.44.0 | Indirect | Supported |  |
| `golang.org/x/time` | v0.5.0 | v0.16.0 | v0.14.0 | Indirect | Supported |  |
| `golang.org/x/tools` | v0.36.0 | v0.52.0 | v0.52.1-0.20261106145633-d407941a339a | Indirect | Supported |  |
| `google.golang.org/api` | v0.171.0 | v0.311.0 | v0.259.0 | Indirect | Uncertain (graph-only) | Graph-only via go4.org; not linked |
| `google.golang.org/genproto` | v0.0.0-20240213162025-012b6fc9bca9 | v0.0.0-20261104224025-d1e00e9f33ab | v0.0.0-20251202230838-ff82c1b0f217 | Indirect | Uncertain (graph-only) |  |
| `google.golang.org/genproto/googleapis/api` | v0.0.0-20240311132316-a219d84964c2 | v0.0.0-20261104224025-d1e00e9f33ab | v0.0.0-20251202230838-ff82c1b0f217 | Indirect | Uncertain (graph-only) |  |
| `google.golang.org/genproto/googleapis/rpc` | v0.0.0-20240314234333-6e1732d8331c | v0.0.0-20261104224025-d1e00e9f33ab | v0.0.0-20251222181119-0a764e51fe1b | Indirect | Uncertain (graph-only) |  |
| `google.golang.org/grpc` | v1.62.1 | v1.85.0 | v1.78.0 | Indirect | Uncertain (graph-only) |  |
| `google.golang.org/protobuf` | v1.33.0 | v1.36.12 | v1.36.11 | Indirect | Uncertain (graph-only) |  |
| `gopkg.in/check.v1` | v1.0.0-20190902080502-41f04d3bba15 | v1.0.0-20201130134442-10cb98267c6c | v0.0.0-20161208181325-20d25e280405 | Indirect | Supported |  |

## Added modules

| Module | Current | Candidate | Selected | Class | Compatibility | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| `cel.dev/expr` | — | v0.26.0 | v0.24.0 | Indirect | Supported |  |
| `cloud.google.com/go/auth` | — | v0.24.3 | v0.18.0 | Indirect | Uncertain (graph-only) |  |
| `cloud.google.com/go/auth/oauth2adapt` | — | v0.3.1 | v0.2.8 | Indirect | Uncertain (graph-only) |  |
| `cloud.google.com/go/iam` | — | v1.15.1 | v1.5.3 | Indirect | Uncertain (graph-only) |  |
| `cloud.google.com/go/monitoring` | — | v1.32.1 | v1.24.3 | Indirect | Uncertain (graph-only) |  |
| `github.com/aymanbagabas/go-udiff` | — | v0.4.1 | v0.2.0 | Tooling | Supported |  |
| `github.com/bits-and-blooms/bitset` | — | v1.29.0 | v1.24.6 | Indirect | Supported |  |
| `github.com/cespare/xxhash/v2` | — | v2.3.0 | v2.3.0 | Indirect | Supported |  |
| `github.com/charmbracelet/colorprofile` | — | v0.4.3 | v0.4.3 | Indirect | Supported |  |
| `github.com/charmbracelet/x/ansi` | — | v0.11.9 | v0.11.9 | Indirect | Supported |  |
| `github.com/charmbracelet/x/cellbuf` | — | v0.0.15 | v0.0.15 | Indirect | Supported |  |
| `github.com/charmbracelet/x/exp/golden` | — | v0.1.0 | v0.0.0-20240806155701-69247e0abc2a | Tooling | Supported |  |
| `github.com/charmbracelet/x/term` | — | v0.2.2 | v0.2.2 | Indirect | Supported |  |
| `github.com/clipperhouse/displaywidth` | — | v0.11.0 | v0.11.0 | Indirect | Supported |  |
| `github.com/clipperhouse/stringish` | — | v0.1.1 | v0.1.1 | Indirect | Supported |  |
| `github.com/clipperhouse/uax29/v2` | — | v2.7.0 | v2.7.0 | Indirect | Supported |  |
| `github.com/cncf/xds/go` | — | v0.0.0-20260202195803-dba9d589def2 | v0.0.0-20251022180443-0feb69152e9f | Indirect | Supported |  |
| `github.com/envoyproxy/go-control-plane/envoy` | — | v1.40.0 | v1.35.0 | Indirect | Supported |  |
| `github.com/go-jose/go-jose/v4` | — | v4.1.9 | v4.1.3 | Indirect | Supported |  |
| `github.com/go-viper/mapstructure/v2` | — | v2.5.0 | v2.5.0 | Indirect | Supported |  |
| `github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp` | — | v1.38.0 | v1.30.0 | Indirect | Uncertain (graph-only) |  |
| `github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric` | — | v0.62.0 | v0.53.0 | Indirect | Uncertain (graph-only) |  |
| `github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping` | — | v0.62.0 | v0.53.0 | Indirect | Uncertain (graph-only) |  |
| `github.com/planetscale/vtprotobuf` | — | v0.6.1-0.20240319094008-0393e58bdf10 | v0.6.1-0.20240319094008-0393e58bdf10 | Indirect | Supported |  |
| `github.com/spiffe/go-spiffe/v2` | — | v2.9.0 | v2.6.0 | Indirect | Supported |  |
| `github.com/stangelandcl/ppmd` | — | v0.1.1 | v0.1.1 | Indirect | Supported |  |
| `github.com/xo/terminfo` | — | v1.3.1 | v1.3.1 | Indirect | Supported |  |
| `go.opentelemetry.io/auto/sdk` | — | v1.2.1 | v1.2.1 | Indirect | Uncertain (graph-only) |  |
| `go.opentelemetry.io/contrib/detectors/gcp` | — | v1.47.0 | v1.38.0 | Indirect | Uncertain (graph-only) |  |
| `go.opentelemetry.io/otel/sdk` | — | v1.47.0 | v1.38.0 | Indirect | Uncertain (graph-only) |  |
| `go.opentelemetry.io/otel/sdk/metric` | — | v1.47.0 | v1.38.0 | Indirect | Uncertain (graph-only) |  |
| `go.yaml.in/yaml/v3` | — | v3.0.6 | v3.0.6 | Indirect | Supported |  |
| `golang.org/x/tools/go/expect` | — | v0.1.1-deprecated | v0.1.1-deprecated | Tooling | Supported |  |
| `golang.org/x/tools/go/packages/packagestest` | — | v0.1.1-deprecated | v0.1.1-deprecated | Tooling | Supported |  |

## Unchanged modules

| Module | Current | Candidate | Selected | Class | Compatibility | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| `github.com/aymanbagabas/go-osc52/v2` | v2.0.1 | v2.0.1 | v2.0.1 | Indirect | Supported |  |
| `github.com/BobuSumisu/aho-corasick` | v1.0.3 | v1.0.3 | v1.0.3 | Direct | Supported |  |
| `github.com/bodgit/plumbing` | v1.3.0 | v1.3.0 | v1.3.0 | Indirect | Supported |  |
| `github.com/bodgit/windows` | v1.0.1 | v1.0.1 | v1.0.1 | Indirect | Supported |  |
| `github.com/cpuguy83/go-md2man/v2` | v2.0.6 | v2.0.7 | v2.0.6 | Indirect | Supported |  |
| `github.com/dsnet/compress` | v0.0.2-0.20230904184137-39efe44ab707 | v0.0.2-0.20230904184137-39efe44ab707 | v0.0.2-0.20230904184137-39efe44ab707 | Indirect | Supported |  |
| `github.com/dsnet/golib` | v0.0.0-20171103203638-1ea166775780 | v1.0.2 | v0.0.0-20171103203638-1ea166775780 | Indirect | Supported |  |
| `github.com/felixge/httpsnoop` | v1.0.4 | v1.1.0 | v1.0.4 | Indirect | Supported |  |
| `github.com/frankban/quicktest` | v1.14.6 | v1.14.6 | v1.14.6 | Indirect | Supported |  |
| `github.com/gitleaks/go-gitdiff` | v0.9.1 | v0.9.1 | v0.9.1 | Direct | Supported |  |
| `github.com/go-logr/stdr` | v1.2.2 | v1.2.2 | v1.2.2 | Indirect | Supported |  |
| `github.com/google/go-cmdtest` | v0.4.0 | v0.4.0 | v0.4.0 | Indirect | Supported |  |
| `github.com/google/go-cmp` | v0.7.0 | v0.7.0 | v0.7.0 | Direct | Supported |  |
| `github.com/google/renameio` | v0.1.0 | v1.0.1 | v0.1.0 | Indirect | Supported |  |
| `github.com/google/uuid` | v1.6.0 | v1.6.0 | v1.6.0 | Indirect | Supported |  |
| `github.com/h2non/filetype` | v1.1.3 | v1.1.3 | v1.1.3 | Direct | Supported |  |
| `github.com/hashicorp/golang-lru/v2` | v2.0.7 | v2.0.7 | v2.0.7 | Indirect | Supported |  |
| `github.com/inconshreveable/mousetrap` | v1.1.0 | v1.1.0 | v1.1.0 | Indirect | Supported |  |
| `github.com/klauspost/cpuid` | v1.2.0 | v1.3.1 | v1.2.0 | Indirect | Supported |  |
| `github.com/kr/pretty` | v0.3.1 | v0.3.1 | v0.3.1 | Indirect | Supported |  |
| `github.com/kr/text` | v0.2.0 | v0.2.0 | v0.2.0 | Indirect | Supported |  |
| `github.com/lucasjones/reggen` | v0.0.0-20200904144131-37ba4fa293bb | v0.0.0-20200904144131-37ba4fa293bb | v0.0.0-20200904144131-37ba4fa293bb | Direct | Supported |  |
| `github.com/Masterminds/goutils` | v1.1.1 | v1.1.1 | v1.1.1 | Indirect | Supported |  |
| `github.com/Masterminds/sprig/v3` | v3.3.0 | v3.3.0 | v3.3.0 | Direct | Supported |  |
| `github.com/mholt/archives` | v0.1.5 | v0.1.5 | v0.1.5 | Direct | Supported |  |
| `github.com/mikelolasagasti/xz` | v1.0.1 | v1.0.1 | v1.0.1 | Indirect | Supported |  |
| `github.com/mitchellh/copystructure` | v1.2.0 | v1.2.0 | v1.2.0 | Indirect | Supported |  |
| `github.com/mitchellh/reflectwalk` | v1.0.2 | v1.0.2 | v1.0.2 | Indirect | Supported |  |
| `github.com/pkg/errors` | v0.9.1 | v0.9.1 | v0.9.1 | Indirect | Supported |  |
| `github.com/rogpeppe/go-internal` | v1.9.0 | v1.16.0 | v1.9.0 | Tooling | Supported |  |
| `github.com/russross/blackfriday/v2` | v2.1.0 | v2.1.0 | v2.1.0 | Indirect | Supported |  |
| `github.com/rwcarlsen/goexif` | v0.0.0-20190401172101-9e8deecbddbd | v0.0.0-20190401172101-9e8deecbddbd | v0.0.0-20190401172101-9e8deecbddbd | Indirect | Supported |  |
| `github.com/sorairolake/lzip-go` | v0.3.8 | v0.3.8 | v0.3.8 | Indirect | Supported |  |
| `github.com/spf13/afero` | v1.15.0 | v1.15.0 | v1.15.0 | Indirect | Supported |  |
| `github.com/STARRY-S/zip` | v0.2.3 | v0.2.3 | v0.2.3 | Indirect | Supported |  |
| `github.com/subosito/gotenv` | v1.6.0 | v1.6.0 | v1.6.0 | Indirect | Supported |  |
| `github.com/xyproto/randomstring` | v1.0.5 | v1.2.0 | v1.0.5 | Tooling | Supported |  |
| `golang.org/x/xerrors` | v0.0.0-20191204190536-9bdfabe68543 | v0.0.0-20240903120638-7835f813f4da | v0.0.0-20191204190536-9bdfabe68543 | Indirect | Supported |  |
| `gopkg.in/yaml.v2` | v2.3.0 | v2.4.0 | v2.3.0 | Indirect | Supported |  |
| `gopkg.in/yaml.v3` | v3.0.1 | v3.0.1 | v3.0.1 | Indirect | Supported |  |

## Removed modules

These modules were present in the baseline graph and absent after the upgrade/tidy. Removal is not by itself evidence of improved security; it reduces graph surface but may also remove old test-only tooling.

| Module | Baseline version | Disposition | Class |
| --- | --- | --- | --- |
| `cloud.google.com/go/bigquery` | v1.3.0 | Removed by graph pruning/tidy | Indirect |
| `cloud.google.com/go/compute` | v1.24.0 | Removed by graph pruning/tidy | Indirect |
| `cloud.google.com/go/datastore` | v1.0.0 | Removed by graph pruning/tidy | Indirect |
| `cloud.google.com/go/firestore` | v1.15.0 | Removed by graph pruning/tidy | Indirect |
| `cloud.google.com/go/longrunning` | v0.5.5 | Removed by graph pruning/tidy | Indirect |
| `cloud.google.com/go/pubsub` | v1.1.0 | Removed by graph pruning/tidy | Indirect |
| `dmitri.shuralyov.com/gpu/mtl` | v0.0.0-20190408044501-666a987793e9 | Removed by graph pruning/tidy | Indirect |
| `github.com/armon/go-metrics` | v0.4.1 | Removed by graph pruning/tidy | Indirect |
| `github.com/BurntSushi/toml` | v0.3.1 | Removed by graph pruning/tidy | Indirect |
| `github.com/BurntSushi/xgb` | v0.0.0-20160522181843-27f122750802 | Removed by graph pruning/tidy | Indirect |
| `github.com/census-instrumentation/opencensus-proto` | v0.2.1 | Removed by graph pruning/tidy | Indirect |
| `github.com/chzyer/logex` | v1.1.10 | Removed by graph pruning/tidy | Indirect |
| `github.com/chzyer/readline` | v0.0.0-20180603132655-2972be24d48e | Removed by graph pruning/tidy | Indirect |
| `github.com/chzyer/test` | v0.0.0-20180213035817-a1ea475d72b1 | Removed by graph pruning/tidy | Indirect |
| `github.com/client9/misspell` | v0.3.4 | Removed by graph pruning/tidy | Indirect |
| `github.com/coreos/go-semver` | v0.3.0 | Removed by graph pruning/tidy | Indirect |
| `github.com/envoyproxy/go-control-plane` | v0.9.1-0.20191026205805-5f8ba28d4473 | Removed by graph pruning/tidy | Indirect |
| `github.com/fatih/color` | v1.14.1 | Removed by graph pruning/tidy | Indirect |
| `github.com/go-gl/glfw` | v0.0.0-20190409004039-e6da0acd62b1 | Removed by graph pruning/tidy | Indirect |
| `github.com/go-gl/glfw/v3.3/glfw` | v0.0.0-20191125211704-12ad95a8df72 | Removed by graph pruning/tidy | Indirect |
| `github.com/godbus/dbus/v5` | v5.0.4 | Removed by graph pruning/tidy | Indirect |
| `github.com/gogo/protobuf` | v1.3.2 | Removed by graph pruning/tidy | Indirect |
| `github.com/golang/glog` | v0.0.0-20160126235308-23def4e6c14b | Removed by graph pruning/tidy | Indirect |
| `github.com/golang/groupcache` | v0.0.0-20210331224755-41bb18bfe9da | Removed by graph pruning/tidy | Indirect |
| `github.com/golang/mock` | v1.4.0 | Removed by graph pruning/tidy | Indirect |
| `github.com/golang/protobuf` | v1.5.3 | Removed by graph pruning/tidy | Indirect |
| `github.com/google/btree` | v1.0.0 | Removed by graph pruning/tidy | Indirect |
| `github.com/google/martian` | v2.1.0+incompatible | Removed by graph pruning/tidy | Indirect |
| `github.com/google/pprof` | v0.0.0-20200212024743-f11f1df84d12 | Removed by graph pruning/tidy | Indirect |
| `github.com/hashicorp/consul/api` | v1.28.2 | Removed by graph pruning/tidy | Indirect |
| `github.com/hashicorp/errwrap` | v1.1.0 | Removed by graph pruning/tidy | Indirect |
| `github.com/hashicorp/go-cleanhttp` | v0.5.2 | Removed by graph pruning/tidy | Indirect |
| `github.com/hashicorp/go-hclog` | v1.5.0 | Removed by graph pruning/tidy | Indirect |
| `github.com/hashicorp/go-immutable-radix` | v1.3.1 | Removed by graph pruning/tidy | Indirect |
| `github.com/hashicorp/go-multierror` | v1.1.1 | Removed by graph pruning/tidy | Indirect |
| `github.com/hashicorp/go-rootcerts` | v1.0.2 | Removed by graph pruning/tidy | Indirect |
| `github.com/hashicorp/golang-lru` | v0.5.4 | Removed by graph pruning/tidy | Indirect |
| `github.com/hashicorp/hcl` | v1.0.0 | Removed by graph pruning/tidy | Indirect |
| `github.com/hashicorp/serf` | v0.10.1 | Removed by graph pruning/tidy | Indirect |
| `github.com/ianlancetaylor/demangle` | v0.0.0-20181102032728-5e5cf60278f6 | Removed by graph pruning/tidy | Indirect |
| `github.com/json-iterator/go` | v1.1.12 | Removed by graph pruning/tidy | Indirect |
| `github.com/jstemmer/go-junit-report` | v0.9.1 | Removed by graph pruning/tidy | Indirect |
| `github.com/kisielk/gotool` | v1.0.0 | Removed by graph pruning/tidy | Indirect |
| `github.com/kr/pty` | v1.1.1 | Removed by graph pruning/tidy | Indirect |
| `github.com/magiconair/properties` | v1.8.9 | Removed by graph pruning/tidy | Indirect |
| `github.com/mitchellh/go-homedir` | v1.1.0 | Removed by graph pruning/tidy | Indirect |
| `github.com/mitchellh/mapstructure` | v1.5.0 | Removed by graph pruning/tidy | Indirect |
| `github.com/modern-go/concurrent` | v0.0.0-20180306012644-bacd9c7ef1dd | Removed by graph pruning/tidy | Indirect |
| `github.com/modern-go/reflect2` | v1.0.2 | Removed by graph pruning/tidy | Indirect |
| `github.com/muesli/reflow` | v0.2.1-0.20210115123740-9e1d0d53df68 | Removed by graph pruning/tidy | Indirect |
| `github.com/nats-io/nats.go` | v1.34.0 | Removed by graph pruning/tidy | Indirect |
| `github.com/nats-io/nkeys` | v0.4.7 | Removed by graph pruning/tidy | Indirect |
| `github.com/nats-io/nuid` | v1.0.1 | Removed by graph pruning/tidy | Indirect |
| `github.com/prometheus/client_model` | v0.0.0-20190812154241-14fe0d1b01d4 | Removed by graph pruning/tidy | Indirect |
| `github.com/sagikazarmark/crypt` | v0.19.0 | Removed by graph pruning/tidy | Indirect |
| `github.com/sagikazarmark/slog-shim` | v0.1.0 | Removed by graph pruning/tidy | Indirect |
| `github.com/tetratelabs/wazero` | v1.9.0 | Removed by graph pruning/tidy | Regex / Wasm |
| `github.com/wasilibs/wazero-helpers` | v0.0.0-20240620070341-3dff1577cd52 | Removed by graph pruning/tidy | Regex / Wasm |
| `github.com/yuin/goldmark` | v1.4.13 | Removed by graph pruning/tidy | Indirect |
| `go.etcd.io/etcd/api/v3` | v3.5.12 | Removed by graph pruning/tidy | Indirect |
| `go.etcd.io/etcd/client/pkg/v3` | v3.5.12 | Removed by graph pruning/tidy | Indirect |
| `go.etcd.io/etcd/client/v2` | v2.305.12 | Removed by graph pruning/tidy | Indirect |
| `go.etcd.io/etcd/client/v3` | v3.5.12 | Removed by graph pruning/tidy | Indirect |
| `go.opencensus.io` | v0.24.0 | Removed by graph pruning/tidy | Indirect |
| `go.uber.org/atomic` | v1.9.0 | Removed by graph pruning/tidy | Indirect |
| `go.uber.org/multierr` | v1.11.0 | Removed by graph pruning/tidy | Indirect |
| `go.uber.org/zap` | v1.21.0 | Removed by graph pruning/tidy | Indirect |
| `golang.org/x/image` | v0.0.0-20190802002840-cff245a6509b | Removed by graph pruning/tidy | Indirect |
| `golang.org/x/lint` | v0.0.0-20200130185559-910be7a94367 | Removed by graph pruning/tidy | Indirect |
| `golang.org/x/mobile` | v0.0.0-20190719004257-d2bd2a29d028 | Removed by graph pruning/tidy | Indirect |
| `google.golang.org/appengine` | v1.6.8 | Removed by graph pruning/tidy | Indirect |
| `gopkg.in/errgo.v2` | v2.1.0 | Removed by graph pruning/tidy | Indirect |
| `gopkg.in/ini.v1` | v1.67.0 | Removed by graph pruning/tidy | Indirect |
| `honnef.co/go/tools` | v0.0.1-2019.2.3 | Removed by graph pruning/tidy | Indirect |
| `rsc.io/binaryregexp` | v0.2.0 | Removed by graph pruning/tidy | Indirect |
| `rsc.io/quote/v3` | v3.1.0 | Removed by graph pruning/tidy | Indirect |
| `rsc.io/sampler` | v1.3.0 | Removed by graph pruning/tidy | Indirect |

## Remaining available updates

At the end of the experiment, 53 indirect graph modules still reported a newer candidate. None is imported by the main module's build graph; they are retained at the versions selected by parent modules, predominantly the `go4.org`/Google Cloud test graph. No direct dependency retained a newer candidate. Forcing those graph-only modules upward would expand the dependency surface without changing the scanner binary, so they were left at parent-selected versions and are not represented as failed upgrades.

## Toolchain and deprecation notes

- **Go 1.26.0 minimum:** imposed by `github.com/wasilibs/go-re2@v1.13.0`. The repository's CI uses Go 1.25 and its Dockerfile uses Go 1.24; both fail before compilation with `go.mod requires go >= 1.26.0`.
- **`go4.org` pseudo-version:** the latest timestamped revision adds Google Cloud/OpenTelemetry requirements (for `go4.org/readerutil` via `mholt/archives`/`sevenzip`). These modules are graph-only in this build.
- **Deprecated modules:** no selected module reports `Deprecated` in `go list -m -json all`. The removed baseline graph contained deprecated modules such as `github.com/golang/protobuf`; they are no longer selected.
- **Retractions:** no selected module exposes a `Retracted` field. `go list -m -retracted all` was not completed because it attempted to add go.sum entries under `-mod=readonly`; this is a verification limitation, not evidence of retractions.
- **`golang.org/x/crypto@v0.58.0`:** `govulncheck` reports module-level `GO-2026-5932` for the unmaintained `openpgp` package; the scanner does not import or call it. A fixed pseudo-version exists, but no stable release was selected by `go get -u`.

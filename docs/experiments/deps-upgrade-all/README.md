# Experimental full dependency and RE2 upgrade

**Status:** Experiment complete — owner review pending; **not merged and not production-approved**
**Source baseline:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (canonical application source)
**Branch point / documentation HEAD:** `afa2c3c8280b7dd581c6bf3b2bc4c7d7b224a20b`
**Experiment branch:** `test/deps-upgrade-all`
**Canonical CP5 baseline remains unchanged.** This experiment does not supersede CP5 and does not mark CP5 owner acceptance.

## Objective

Determine whether the fork can move to the latest stable Go dependency stack, including `github.com/wasilibs/go-re2` v1.13.0, and characterize the build, test, detection, performance, security, and migration consequences before any remediation work is authorized.

The experiment was intentionally aggressive: dependency breakage was an acceptable outcome. The result is recorded even where it is inconvenient.

## Outcome in one page

| Area | Result |
| --- | --- |
| Dependency resolution | **Pass** (all direct modules latest; 53 graph-only indirect modules remain parent-pinned) |
| Upgrade completeness | **Partial** — 69 modules updated (9 direct, 60 indirect); 34 added; 77 removed |
| Standard build | **Pass** |
| Standard tests | **Pass** — 275/275, 0 skip (Linux Go 1.27.3) |
| RE2 application build | **Pass** |
| RE2 full package graph | **Fail** — same pre-existing generator-rule errors as baseline |
| RE2 tests | **Fail (full graph)** — same pre-existing errors; targeted RE2 package tests **Pass** |
| Detection parity | **Equivalent** — 40/40 normalized findings in the exercised matrix, Linux and Windows |
| Performance | **Improved** for RE2 v1.13; standard engine neutral/slightly improved; see methodology caveat |
| Windows compatibility | **Pass** for build and CLI parity; full native test suite **not run** |
| Security advisory review | **Complete** — two reachable baseline vulnerabilities fixed; no reachable upgraded vulnerabilities |

**Migration gate:** the upgrade raises the module Go requirement from **1.24.11 to 1.26.0**. Current CI uses Go 1.25 and the Dockerfile uses Go 1.24; both fail immediately with `go.mod requires go >= 1.26.0`. This is a required toolchain migration, not a source incompatibility.

**Detection caveat:** the RE2 engine was not assumed equivalent. Identical rule/config/source inputs were compared after normalizing nondeterministic finding order. No detection differences were found in the exercised corpus, including decoding depth 2, archive depth 1, Unicode, entropy, allowlists, and generic/specific deduplication. This is bounded evidence, not proof of universal rule equivalence.

**Performance caveat:** the largest apparent speed-up is a **Go runtime default change**, not a dependency effect. The baseline module directive (Go 1.24.11) disables automatic container-aware `GOMAXPROCS`; the upgraded directive (Go 1.26.0) enables it. Measurements were therefore repeated with `GOMAXPROCS=2` for both stacks. The matched results show standard-engine changes within roughly ±10% (small workloads noisy), while RE2 v1.13 is materially faster and lower-memory than RE2 v1.9. The CP3 `28.7×` figure is not used as an optimization claim.

## Tested environments

| Environment | Use |
| --- | --- |
| Linux/amd64, Go 1.27.3, `CGO_ENABLED=1` | Builds, tests, race, targeted RE2 tests |
| Linux/amd64, Go 1.27.3, `CGO_ENABLED=0` | Application binaries and detection/performance runs |
| Docker 29.9.2, `golang:1.27.3` | Primary reproducible environment |
| Windows/amd64, Go 1.27.3 | Native stdlib and RE2 builds; CLI detection parity |
| `golang:1.25` (Go 1.25.14) | Compatibility check — **fails at go.mod gate** |
| `golang:1.24` (Go 1.24.13) | Compatibility check — **fails at go.mod gate** |

Performance runs used 2 vCPU, 1 GiB, no swap, `--network none`, warm module/build caches. Detection runs used the same synthetic corpus and configs on both stacks.

## Dependency outcome

- `go 1.24.11` → `go 1.26.0`.
- `github.com/wasilibs/go-re2` `v1.9.0` → `v1.13.0`.
  - Moves from the wazero WebAssembly runtime to `wasm2go`.
  - Removes `github.com/tetratelabs/wazero` and `github.com/wasilibs/wazero-helpers` from the graph.
  - Requires **Go 1.26.0**.
  - The Gitleaks `regexp` wrapper required no source change.
- Major-version direct move: `github.com/charmbracelet/lipgloss` `v0.5.0` → `v1.1.0`; the limited API used by the scanner compiles.
- `github.com/spf13/viper` `v1.19.0` → `v1.21.0`; config precedence and unmarshalling were re-exercised successfully.
- Security-relevant upgrades: `github.com/andybalholm/brotli` `v1.2.0` → `v1.2.8` and `github.com/klauspost/compress` `v1.18.0` → `v1.20.4` remove the two reachable baseline vulnerabilities.
- The latest `go4.org` pseudo-version adds graph-only Google Cloud/OpenTelemetry modules. They are not linked into the scanner binary, but they enlarge the module graph and supply-chain surface.
- `go mod verify`: all modules verified. `go mod tidy -diff`: clean.
- Two test-only modules resolved to stable releases below baseline pseudo-versions (`davecgh/go-spew`, `pmezard/go-difflib`); this is a deliberate “latest stable” result, not a build regression.

## Existing failures preserved

These failures were present at the CP5 baseline and remain; they are **not** introduced regressions:

1. `go vet ./...` reports the two test-only `lostcancel` diagnostics at `detect/detect_test.go:2083` and `:2122`.
2. `go build -tags gore2regex ./...` and `go test -tags gore2regex ./...` fail in `cmd/generate/config/rules/airtable.go:30` and `perplexity.go:16` (stdlib `*regexp.Regexp` used where the tagged wrapper type is expected).
3. `TestTranslateExtend` is order-sensitive. A fixed shuffle seed (`1794264498871797264`) fails on both baseline and upgraded stacks; the upgraded random `-shuffle=on` run happened to expose it. This is the CP5-F07 global-state/test-order issue, not a dependency regression.

## Introduced compatibility changes

1. **Go 1.26 minimum** (go-re2 v1.13). Go 1.25.14 and Go 1.24.13 both fail before compilation.
2. **CI and Docker toolchains must be raised** to Go 1.26+ (tested at Go 1.27.3) before this branch can be a CI candidate.
3. **Optional `re2_cgo` remains separate.** `CGO_ENABLED=1` alone does not select it. The `re2_cgo` build requires system `re2` and `pkg-config` (`libre2-dev` on Debian/Ubuntu); without them the build fails with `Package re2 was not found`.
4. **Repository graph expansion** from the latest `go4.org` pseudo-version. This does not affect linked code but is a supply-chain review item.

No application source, tests, fixtures, scripts, config, CI, Dockerfile, Makefile, or Git metadata were changed by this experiment.

## Deliverables

- [`dependency-matrix.md`](dependency-matrix.md) — complete baseline/candidate/selected inventory, updates, additions, removals, and remaining candidates.
- [`compatibility-report.md`](compatibility-report.md) — build/test/toolchain/regex-engine/CLI/report comparison.
- [`regression-register.md`](regression-register.md) — introduced changes, existing failures, and non-regressions.
- [`performance-report.md`](performance-report.md) — matched workloads, memory, and the `GOMAXPROCS` confound.
- [`upgrade-assessment.md`](upgrade-assessment.md) — final recommendation against the assignment’s ten questions.

## Recommendation

**B. Adopt with targeted compatibility changes.**

The full stack is technically viable and passes the local Go 1.27.3 test suite, both application builds, targeted RE2 tests, detection parity, and Windows build/CLI checks. It fixes two reachable dependency vulnerabilities and materially improves RE2 performance. Adoption should be conditional on:

1. Raising CI and the Docker build image to **Go 1.26+** (Go 1.27.3 is the tested version).
2. Accepting or explicitly reducing the graph expansion from the latest `go4.org`/charmbracelet transitive updates.
3. Keeping CP5 F01 (scan-completeness/false-success) and F02 (destructive report paths) as separate P0 remediation work; this upgrade does not fix them.
4. Treating `re2_cgo` as an optional, separately provisioned native variant.

If Go 1.25 must remain the supported toolchain, this is no longer a full-stack recommendation: fall back to **C** and pin `go-re2` to v1.12.0 (Go 1.25) or defer the RE2 move. That fallback was not tested in this experiment.

No merge, production release, CP5 remediation, or external PR adoption was performed. Stop for owner review.

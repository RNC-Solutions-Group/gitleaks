# Regression register — dependency and RE2 upgrade

**Status:** Experimental; no production remediation
**Branch:** `test/deps-upgrade-all`
**Comparison basis:** canonical source `5f6ec56…`, baseline docs HEAD `afa2c3c…`, upgraded dependency commit `70aaf01`.

Classification follows the assignment: existing baseline failure, new dependency incompatibility, toolchain incompatibility, regex-engine behavior change, environment limitation, or unrelated application defect.

## Introduced failures / mandatory changes

### REG-01 — Go 1.26 minimum (introduced toolchain incompatibility)

| Field | Value |
| --- | --- |
| Status | **Introduced, mandatory** |
| Evidence | Upgraded `go.mod` has `go 1.26.0`; Go 1.25.14 and Go 1.24.13 both fail before compilation |
| Exact error | `go: go.mod requires go >= 1.26.0 (running go 1.25.14; GOTOOLCHAIN=local)`; same with 1.24.13 |
| Cause | `github.com/wasilibs/go-re2@v1.13.0` declares `go 1.26.0` |
| Reproduction | Extract `src-upgraded-lf.tar`; run `GOFLAGS=-mod=readonly GOTOOLCHAIN=local go build ./...` with Go 1.25 or 1.24 |
| Blocks further tests? | No, for Go 1.27.3. Yes, for the repository’s current CI (1.25) and Dockerfile (1.24) |
| Disposition | Required migration: raise CI and Docker toolchains to Go 1.26+ (1.27.3 tested). No source change. |

### REG-02 — CI and Docker toolchain no longer match (introduced configuration consequence)

| Field | Value |
| --- | --- |
| Status | **Introduced (derived from REG-01)** |
| Evidence | `.github/workflows/test.yml:23` uses `go-version: 1.25`; `Dockerfile:1` uses `golang:1.24` |
| Effect | CI build/test and container build fail at the `go.mod` gate before downloading modules |
| Cause | Same as REG-01 |
| Reproduction | Same as REG-01 with the exact workflow/Dockerfile versions |
| Disposition | Update both to Go 1.26+ as part of adoption. This experiment did not edit CI or Docker files. |

## Behavior differences without test/build failure

### REG-03 — Module graph expansion and churn

| Field | Value |
| --- | --- |
| Status | **Introduced behavior difference; not a failure** |
| Evidence | 143 final graph modules vs 186 baseline; 69 updated, 34 added, 77 removed, 40 unchanged |
| Cause | `go get -u -t ./...` + `go mod tidy`, including latest `go4.org`, lipgloss v1.1, viper v1.21, and go-re2 v1.13 graph changes |
| Notable additions | Google Cloud/OpenTelemetry modules via `go4.org`; charmbracelet x/*; graph-only `golang.org/x/tools` test packages |
| Impact | Larger dependency-review surface; most new modules are not linked into the scanner binary |
| Reproduction | `go list -m all`; compare with baseline inventory |
| Disposition | Accept or reduce deliberately; no runtime detection/test failure observed. |

### REG-04 — Two test-only modules resolved below baseline pseudo-versions

| Field | Value |
| --- | --- |
| Status | **Introduced version change; not a test failure** |
| Evidence | `github.com/davecgh/go-spew` `v1.1.2-0.20180830191138-d8f796af33cc` → `v1.1.1`; `github.com/pmezard/go-difflib` `v1.0.1-0.20181226105442-5d4384ee4fb2` → `v1.0.0` |
| Cause | Latest-stable selection and graph tidy; baseline used development pseudo-versions |
| Impact | Test-only; no build/test failure. Not security-analyzed as a downgrade beyond `govulncheck` |
| Disposition | Accept as latest stable, or pin baseline pseudo-versions if preserving exact test graph is required. |

## Existing baseline failures (not introduced)

### REG-05 — `go vet` lostcancel diagnostics

| Field | Value |
| --- | --- |
| Status | **Existing baseline failure** |
| Evidence | `detect/detect_test.go:2083`, `:2122`; identical before/after |
| Cause | Test-only context cancellation analysis |
| Reproduction | `go vet ./...` |
| Disposition | Unchanged CP5-F14. Not part of this upgrade. |

### REG-06 — `gore2regex` full package graph fails in generator packages

| Field | Value |
| --- | --- |
| Status | **Existing baseline failure** |
| Evidence | `airtable.go:30`, `perplexity.go:16`; identical before/after |
| Cause | Generator packages use stdlib `*regexp.Regexp` with the tagged wrapper type |
| Reproduction | `go build -tags gore2regex ./...`; `go test -tags gore2regex ./...` |
| Disposition | Tagged application binary still builds; targeted RE2 package tests pass. Unchanged CP5-F14. |

### REG-07 — `TestTranslateExtend` shuffle-order failure

| Field | Value |
| --- | --- |
| Status | **Existing baseline failure / test-order sensitivity** |
| Evidence | Upgraded random shuffle failed at seed `1794264498871797264`; re-running the same seed on baseline also fails identically |
| Cause | `extendDepth` package-global state and test ordering (CP5-F07) |
| Reproduction | `go test -shuffle=1794264498871797264 ./config -count=1 -v` on either stack |
| Disposition | Not a dependency regression. A single random shuffle pass is not proof of order independence. |

## Environment / optional-variant limitations

### REG-08 — Native `re2_cgo` requires system RE2

| Field | Value |
| --- | --- |
| Status | **Existing optional-variant requirement; not a regression** |
| Evidence | `CGO_ENABLED=1 go build -tags "gore2regex re2_cgo" .` fails with `Package re2 was not found in the pkg-config search path` |
| Cause | `re2_cgo` requires the native RE2 library and `pkg-config` metadata (`libre2-dev` on Debian/Ubuntu) |
| Impact | Default WebAssembly-based RE2 works; the native variant was not built/tested in this experiment |
| Reproduction | Same build command in `golang:1.27.3` |
| Disposition | Keep as a separate, provisioned variant; do not infer it from `CGO_ENABLED=1`. |

### REG-09 — Small-workload standard-engine timing deltas are inconclusive

| Field | Value |
| --- | --- |
| Status | **Inconclusive; not confirmed regression** |
| Evidence | Matched `GOMAXPROCS=2` medians: std encoded depth 0 `87 → 99 ms` (+13.8%), archive depth 1 `89 → 98 ms` (+10.1%); other standard workloads improved |
| Cause | Three repetitions, sub-100 ms workloads, warm-cache noise, process startup and scheduling |
| Impact | No detection difference; no coverage loss |
| Disposition | Do not report as a performance regression without more repetitions; larger workloads and RE2 show improvement. |

## Confirmed non-regressions

- Detection parity: 40/40 normalized findings for baseline stdlib vs upgraded stdlib, baseline RE2 vs upgraded RE2, upgraded stdlib vs upgraded RE2, and baseline stdlib vs baseline RE2 on the exercised corpus; 38/38 at decode depth 0.
- Default config: 3/3 expected findings across all four builds; same rule IDs and values.
- Windows native: 40/40 parity for all four builds.
- Report formats and `--redact`: valid output, no raw synthetic token in redacted JSON.
- Config precedence: source-local, `GITLEAKS_CONFIG`, and explicit `--config` all behaved as documented.
- Application binaries: both engines build and run; default configuration loads.

## Summary

| Class | Count |
| --- | ---: |
| New source/build/test failures on Go 1.27.3 | **0** |
| New mandatory toolchain incompatibilities | **1** (Go 1.26 minimum) |
| New configuration consequences | **1** (CI/Docker must move) |
| New detection regressions | **0** |
| Existing persistent failures | **3 categories** (vet, tagged full graph, shuffle-order test) |
| Unverified/environment limitations | Windows full tests, `re2_cgo`, CI/release execution |

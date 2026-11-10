# Upgrade assessment — latest dependency stack and RE2

**Status:** Experimental assessment complete; owner decision required
**Branch:** `test/deps-upgrade-all`
**Canonical source:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`
**CP5 baseline acceptance status:** unchanged; not accepted on the owner’s behalf

## Decision summary

| Area | Result |
| --- | --- |
| Dependency resolution | **Pass** |
| Upgrade completeness | **Partial** — 69 updated (9 direct, 60 indirect); 53 graph-only indirect candidates remain |
| Standard build | **Pass** |
| Standard tests | **Pass** — 275/275, 0 skip on Go 1.27.3 |
| RE2 application build | **Pass** |
| RE2 full graph | **Fail** — pre-existing generator-rule errors |
| RE2 tests | **Fail (full graph)**; targeted RE2 packages **Pass** |
| Detection parity | **Equivalent** in the exercised matrix |
| Performance | **Improved** for RE2; standard engine neutral/slightly improved |
| Windows compatibility | **Pass** for native build and CLI parity; full tests **Not run** |
| Security advisory review | **Complete** — 2 reachable baseline vulnerabilities fixed; 0 reachable upgraded |

## 1. Is the latest dependency stack viable?

**Yes on Go 1.26+, no on the repository’s current CI/Docker toolchains.**

On Go 1.27.3 the upgraded stack:

- resolves and verifies (`go mod verify`: all modules verified; `go mod tidy -diff`: clean);
- builds the full default package graph;
- passes **275/275** tests and **275/275** race tests;
- builds both the standard and `gore2regex` application binaries;
- passes targeted RE2 package tests;
- builds and runs natively on Windows;
- preserves the baseline’s known `vet` and tagged-full-graph failures.

The blocker is the module directive: `go 1.26.0` is required by `go-re2 v1.13.0`. Go 1.25.14 (CI) and Go 1.24.13 (Dockerfile) fail before compilation with `go.mod requires go >= 1.26.0`. This is a toolchain migration, not an application-source incompatibility.

The upgrade is **whole-stack** for direct modules: all 16 are at their latest eligible versions. It is **partial** for the module graph: 53 indirect modules still have newer candidates, but they are graph-only (predominantly the `go4.org`/Google Cloud/OpenTelemetry and test-tooling graph) and are not imported by the scanner. Forcing them would enlarge the graph without changing linked code. This distinction is recorded rather than hiding it under a “Pass” label.

## 2. Is the latest RE2 implementation viable?

**Yes for the default WebAssembly-based engine.** `github.com/wasilibs/go-re2` v1.13.0:

- builds the `gore2regex` application;
- loads the embedded default config and rules;
- produces the same normalized findings as the standard engine and as RE2 v1.9 on the exercised corpus;
- is materially faster and lower-memory than v1.9;
- requires no Gitleaks source change.

The full `go test -tags gore2regex ./...` command still fails in `cmd/generate/config/rules/airtable.go` and `perplexity.go`, exactly as at baseline. Targeted tests for the packages that do compile under the tag pass.

The optional native `re2_cgo` variant was **not** tested successfully. It requires system `re2` and `pkg-config` (`libre2-dev` on Debian/Ubuntu), and `CGO_ENABLED=1` alone does not select it. This is a separate deployment variant, not part of the default result.

## 3. Which existing failures remain?

1. `go vet ./...`: two test-only `lostcancel` diagnostics at `detect/detect_test.go:2083` and `:2122`.
2. `go build -tags gore2regex ./...` / `go test -tags gore2regex ./...`: generator-rule type errors in `airtable.go:30` and `perplexity.go:16`.
3. `TestTranslateExtend` order sensitivity under shuffle; the fixed seed fails on both stacks (CP5-F07/global `extendDepth` state).
4. CP5 F01 (false-success/partial scan status) and F02 (destructive report paths), plus all other CP5 application findings. The dependency upgrade does not address them.
5. Windows full tests, CI/release execution, advisory database completeness, and the optional native RE2 variant remain unverified in this experiment.

## 4. Which failures were introduced?

- **One mandatory toolchain incompatibility:** Go 1.26 minimum (REG-01), which directly causes the CI Go 1.25 and Docker Go 1.24 build failures (REG-02).
- **One module-graph behavior change:** 34 added modules and 77 removed modules, largely graph-only; this is a review-surface change, not a runtime failure (REG-03).
- **Two test-only version downgrades** to stable releases (`go-spew`, `go-difflib`); no test failure (REG-04).

No new source, build, or test failure was observed on Go 1.27.3. No detection regression was measured.

## 5. What compatibility work is required?

Minimum adoption work:

1. Raise CI from Go 1.25 to **Go 1.26+**; Go 1.27.3 is the tested version.
2. Raise the Docker build image from Go 1.24 to **Go 1.26+**.
3. Update contributor/build documentation for the new minimum.
4. Decide whether the graph expansion from `go4.org`/charmbracelet is acceptable, or explicitly constrain it.
5. If native `re2_cgo` is a supported artifact, provision `libre2-dev`/`pkg-config` and test that variant separately.
6. Preserve the existing application behavior and do not couple this migration to CP5 F01/F02 remediation.

No application logic change is required by the upgrade itself.

## 6. What detection differences were measured?

**None in the exercised matrix.**

- 40/40 normalized findings for baseline stdlib vs upgraded stdlib, baseline RE2 vs upgraded RE2, upgraded stdlib vs upgraded RE2, and baseline stdlib vs baseline RE2.
- 38/38 at decode depth 0; the two encoded findings are absent by design.
- 30/30 in a 30-file ordering corpus.
- 3/3 default-config findings (AWS access token, GitHub PAT, private RSA key).
- Covers rule allowlist, path allowlist, entropy threshold, Unicode `\p{L}`, one- and two-level Base64, archive depth 1, and generic/specific deduplication.
- Windows native CLI repeated the 40/40 parity for all four builds.

Raw finding order differed across runs, consistent with CP5-F12; comparisons were normalized by stable finding fields. This is bounded evidence, not proof that every rule or input matches. No real credentials were used; the corpus is synthetic.

## 7. Is performance better, worse, or inconclusive?

**RE2: better. Standard engine: neutral/slightly better. The large “2×” default-speed difference is a toolchain runtime-default effect, not a dependency effect.**

Matched `GOMAXPROCS=2`, 2 vCPU/1 GiB, 3 repetitions:

- Standard engine: 32 MiB/2,000 files default **3,447 → 3,272 ms** (−5.1%); depth 0 **240 → 226 ms**; encoded depth 2 **619 → 563 ms**; 1,000 findings **245 → 197 ms**; archive depth 1 **89 → 98 ms** and encoded depth 0 **87 → 99 ms** (small-workload regressions, within noise).
- RE2: 32 MiB default **594 → 299 ms** (−49.7%); encoded depth 0/1/2 **433/361/381 → 101/110/101 ms**; archive **415 → 107 ms**; 1,000 findings **766 → 257 ms**; startup **358 → 68 ms**.
- Peak RSS: standard engine roughly unchanged (+9% files, +13% encoded, −7% findings); RE2 materially lower (−36% to −43% on tested workloads).

The baseline module’s Go 1.24.11 default disabled container-aware `GOMAXPROCS`; the upgraded Go 1.26.0 default enables it. In an isolated check, baseline default was 6.34 s, baseline with `GODEBUG=containermaxprocs=1` was 3.41 s, and upgraded default was 3.34 s on the same 32 MiB corpus. The dependency delta must therefore be measured with `GOMAXPROCS` matched, as above. The CP3 8,196 ms baseline and 28.7× depth-0 ratio are not used as an optimization claim here.

## 8. Estimated migration complexity

**Low source complexity, medium build/CI complexity, medium supply-chain review.**

- Application source changes: **none required**.
- Build/test migration: bump CI/Docker to Go 1.26+; update documented minimum. Small.
- Runtime behavior: `GOMAXPROCS` default changes under cgroups; this is a benefit in the tested 2-vCPU container but should be understood as a runtime-default change, not silently attributed to dependencies.
- Detection validation: the bounded parity matrix is clean; broader real-rule regression coverage should still be run before production adoption.
- Dependency review: 34 added/77 removed modules, mostly graph-only; the latest `go4.org` pull-in of cloud/OpenTelemetry modules is the main review item.
- Optional native RE2: requires a separate system-library build and test.

## 9. Should we adopt the upgraded stack before CP5 remediation?

**Yes as an implementation candidate, after the toolchain migration — not as a substitute for CP5 remediation.**

The upgrade is technically viable and improves dependency security and RE2 performance without changing detection in the exercised corpus. But it does not fix the two P0 application risks: incomplete scans can still exit 0, and report paths can still destroy files before scanning. Those remain independent CP5 work.

If the organization cannot move CI and containers off Go 1.25/1.24, the full-stack upgrade should not be adopted now. In that case, use the selective path (§10).

## 10. Which dependency versions should remain pinned?

| Dependency | Version | Reason |
| --- | --- | --- |
| `github.com/wasilibs/go-re2` | **v1.13.0** | Latest stable; `wasm2go`; requires Go 1.26.0. This is the version that provides the tested RE2 performance and graph simplification. |
| Go toolchain | **1.26+** (1.27.3 tested) | Required by go-re2 v1.13; CI and Docker must match. |
| `github.com/wasilibs/go-re2` fallback | **v1.12.0** | Go 1.25-compatible fallback if the toolchain cannot move. Not tested in this experiment; it retains the wazero path. |
| `go4.org` | Latest pseudo or baseline pin | Latest pulls the cloud/OpenTelemetry graph; pinning the baseline `v0.0.0-20230225012048-214862532bf5` may reduce graph surface, but that selective combination was not tested here. |
| `golang.org/x/crypto` | v0.58.0 (as resolved) | Module-level `GO-2026-5932` is in unimported `openpgp`; no reachable call. Do not treat the module-level advisory as a scanner vulnerability. |

No other direct dependency requires a special pin. `aho-corasick`, `sprig`, `go-gitdiff`, `go-cmp`, `filetype`, `archives`, and `reggen` were already at their latest available versions.

## Final recommendation

**B. Adopt with targeted compatibility changes.**

Reasoning: on Go 1.26+ the full stack builds, passes the default and race suites, preserves known failures without adding new ones on Go 1.27.3, maintains detection parity, builds and runs on Windows, fixes two reachable dependency vulnerabilities, and materially improves RE2. The only mandatory migration is the Go 1.26 toolchain floor plus CI/Docker updates. That is a bounded, targeted compatibility change rather than a source rewrite. Adoption should still be an implementation-candidate decision, not a production release, and must not be conflated with closing CP5 F01/F02.

If the Go 1.25 toolchain is a hard constraint, the correct fallback is **C. Upgrade selectively; retain specific pinned dependencies** — pin `go-re2` at v1.12.0 (or defer the RE2 change) and keep the toolchain at 1.25. That path is not validated by this experiment.

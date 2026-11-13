# Dependencies — Gitleaks candidate register

**Current active queue (2026-11-13):** Actioned issues/PRs have moved to the [RNC disposition register](exclusions.md); implementation, test evidence and contributor attribution are maintained in [fixes/QOL](https://github.com/RNC-Solutions-Group/gitleaks/blob/fix/validation/fixes-qol-tracking.md). Do not reassign finished work from historical review notes.

Snapshot: 2026-11-08 upstream open inventory · RNC fork `master` at `5f6ec56cfa6f` · Review branch `feat/issues-and-pr`.

**6 active issues · 2 active PRs · 8 work-package accounting slots.** All retained PR diffs have been inspected against the pinned fork; **source/patch review is not a successful build, dependency resolution, runtime acceptance, or adoption**. PRs remain upstream-open. [Security CVE changes](security-hardening.md) have a separate owner queue; [RNC dispositions](exclusions.md) now also hold independently matched or superseded dependency proposals.

## Current active dependency scope — 2026-11-13

Nine historical dependency proposals were evaluated against the **current `fix/validation`** Go module graph and moved to [dispositions](exclusions.md): matching version selections and independent Lip Gloss v2 migration were not PR integrations; v1 Lip Gloss was superseded. The remaining active PRs are **#2323** (unmodified release login-action pin) and **#2115** (Go 1.26.3 minimum/toolchain security claim still not fully met: `go.mod` says 1.26.0 despite the Go 1.27.4 Docker builder). Go module-path/install/old architecture issues remain active separately.

## Pinned baseline and integration constraints

- `master/go.mod`: module path `github.com/zricethezav/gitleaks/v8`, Go `1.24.11`; `go-re2 v1.9.0`, `testify v1.10.0`, `cobra v1.9.1`, `pflag v1.0.6`, `viper v1.19.0`, `lipgloss v0.5.0`, `zerolog v1.33.0`, `semgroup v1.2.0`, `go-version v1.7.0`.
- `master/Dockerfile` uses `golang:1.24` builder; `.github/workflows/test.yml` uses Go `1.25`; the *separate, unmerged* RNC `fix/docker-security-bases` Dockerfile uses builder Go `1.27.3`. **A newer builder does not upgrade module pins or the test workflow.**
- Dependency PRs are generated from **different upstream base commits**. Do not concatenate their `go.mod`/`go.sum` patches blindly. Re-resolve one graph on the chosen base, preserve authorship/commit references and test the *combined* outcome.
- **Confirmed graph collisions:** [#2250](https://github.com/gitleaks/gitleaks/pull/2250) requests `testify v1.12.1` while [#2139](https://github.com/gitleaks/gitleaks/pull/2139) requests `v1.11.1`; [#2149](https://github.com/gitleaks/gitleaks/pull/2149) requests `pflag v1.0.9` while #2139 requests `v1.0.10`; [#2289](https://github.com/gitleaks/gitleaks/pull/2289) requests `x/sys v0.48.0` while Security [#2308](https://github.com/gitleaks/gitleaks/pull/2308) requests `v0.47.0`. Lipgloss [#2145](https://github.com/gitleaks/gitleaks/pull/2145) is a **v1 package** upgrade; [#2100](https://github.com/gitleaks/gitleaks/pull/2100) is a **v2 module-path migration**. These are not interchangeable patches.

## Remaining open dependency and build PRs — by-item diff review

| Upstream PR / contributor | Actual patch compared to pinned master | RNC disposition | Remaining acceptance gate |
| --- | --- | --- | --- |
| [#2323](https://github.com/gitleaks/gitleaks/pull/2323) — @dependabot[bot] | `.github/workflows/release.yml`: both Docker Hub and GHCR `docker/login-action` uses move from pinned `dbcb813…` v4.6.0 to `81c6f93…` v4.7.0. Master still uses 4.6.0 twice. | **PATCH-REVIEWED — RELEASE WORKFLOW TEST** | Verify full upstream action SHA and permissions, Docker Hub/GHCR authentication using safe CI, and both publish targets. No other workflow is changed in the patch. |
| [#2115](https://github.com/gitleaks/gitleaks/pull/2115) — AlonBilman (historical author, not tagged; proposed approach remains unfulfilled) | Proposes Go module minimum `1.26.3`, CI `1.26.3` and floating `golang:1.26` builder. **Current RNC validation** instead declares `go 1.26.0`, with CI and Docker builder explicitly pinned to `1.27.4`; these are different constraints. | **RETAIN — MODULE MINIMUM / SECURITY SUPPORT POLICY NOT RESOLVED** | The current pinned builder is newer but does **not** establish the requested minimum for downstream consumers building with Go 1.26.0. Decide minimum supported toolchain and verify its CVE/security policy; retain pinned CI/builder and compatibility tests. Do not retire this proposal just because the build container is newer. |

## Open build / dependency issues — by-item source review

**The historical #1588 build/Makefile side-effect defect is corrected by CP-V10** and migrated to [dispositions](exclusions.md). Other module path, legacy install and release issues listed below remain separate active work.

| Upstream issue | Pinned-source finding | RNC disposition and gate |
| --- | --- | --- |
| [#1839](https://github.com/gitleaks/gitleaks/issues/1839) — module path mismatch | `master/go.mod` still declares `github.com/zricethezav/gitleaks/v8`, while repository and some README badges use `github.com/gitleaks/gitleaks/v8`. Dockerfile/Makefile ldflags and imports still reference the old namespace. | **ACTIVE — MIGRATION DESIGN REQUIRED.** Reproduce modern `go get` on versioned `/v8`, assess existing importers and internal build paths before any rename. |
| [#1305](https://github.com/gitleaks/gitleaks/issues/1305) — Git 1.8 rejects `-C` | `sources/git.go` still executes multiple `git -C` commands. The issue's `-c` workaround is **not** syntax-equivalent to `-C`. | **ACTIVE — MINIMUM-GIT POLICY.** Define supported Git floor and provide either a tested fallback or an explicit unsupported-platform message; coordinate #833 documentation. |
| [#1226](https://github.com/gitleaks/gitleaks/issues/1226) — unversioned `go run` | Reported `go run github.com/gitleaks/gitleaks@latest` resolves legacy `v1.25.1`, not the v8 module. | **SOURCE-CHECKED — INSTALL COMMAND MATRIX.** Same naming family as #1209/#1839, but preserve the `go run` acceptance case: fixing only `/v8` imports does not prove unversioned `@latest` works. |
| [#1209](https://github.com/gitleaks/gitleaks/issues/1209) — unversioned `go install` | Reported `go install github.com/gitleaks/gitleaks@latest` hits the same legacy module-path mismatch. | **SOURCE-CHECKED — INSTALL COMMAND MATRIX.** Test this distinct CLI path alongside #1226 with documented versioned install commands; no claim that a module rename fixes historical unversioned paths. |
| [#833](https://github.com/gitleaks/gitleaks/issues/833) — build prerequisites | `go.mod` demands Go `1.24.11`; README lacks a clear Git floor or Go-version prerequisite. A 2022 RHEL 7/Go 1.16 reproducer is historical, not a current support commitment. | **ACTIVE — DOC / SUPPORT POLICY.** State tested minimum Go/Git and platform limits; validate installation paths and link to #1305. |
| [#1944](https://github.com/gitleaks/gitleaks/issues/1944) — release binaries | Reporter describes `SIGILL` on old x86 from a published standalone binary, while a binary extracted from the image works. `.github/workflows/release.yml` currently publishes containers only. | **ACTIVE — BINARY PROVENANCE GAP.** Require reproducible release-binary workflow, explicit `GOAMD64` target/CGO behavior, checksums, builder metadata, old-x86 verification and platform matrix. Stack-trace paths are **not** conclusive proof of where a binary was built. |

## Acceptance and promotion gate

1. Select one baseline, resolve a **single** combined module graph, and show a `go.mod/go.sum` diff. Do not merge multiple Dependabot heads by overwriting each other's resolved checksums.
2. Pin Go toolchain, Docker builder/image digest, Actions commit SHAs and OS/architecture matrix. Validate minimum toolchain versus all selected modules.
3. Run `go mod tidy -diff` or an equivalent no-mutation check, `go build ./...`, `go test ./... -race`, config generation/diff, CLI snapshots, Git and dir scan fixtures, and platform-specific acceptance (when toolchain is available).
4. Exercise regex engine, config parsing, archive inputs, ANSI/redaction, CLI flags, signal/cancel handling and performance under the integrated graph; compare expected detections rather than only exit status.
5. Verify binary provenance, container SBOM/CVE inventory and release workflows independently; those are not established by a clean module resolution.
6. Record exact upstream contributor/PR SHA, RNC integration branch/commit/PR, test evidence and owner approval **before** claiming adopted, corrected or released. Only the user promotes changes; no upstream state changes from triage.

**Scope:** Historical comparison to `master@5f6ec56` is preserved; resolved/superseded status above uses the live `fix/validation` dependency graph. No upstream PRs or commits were integrated by this reconciliation.

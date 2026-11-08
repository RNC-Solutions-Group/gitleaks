# Dependencies — Gitleaks candidate register

Snapshot source: [Open PR inventory](../02-open-pull-requests.md), [issue inventory](../01-issue-backlog.md) (2026-11-08). This is a **title/search-based first pass**, enriched with pinned version checks from the fork. It is not an SBOM comparison against every image.

**Ownership rule:** [closed.md](closed.md) is the canonical RNC exclusion list. Items recorded there are omitted from active tables here. Security CVE PRs and their linked issues are **excluded from this register** and tracked exclusively in [security-hardening.md](security-hardening.md).

## Summary

| PR disposition | Count | Basis |
| --- | ---: | --- |
| Partially superseded — further review | 2 | Docker builder or some modules are newer, but each PR still contains unresolved source/CI or dependency changes |
| Unreviewed / relevance check | 10 | Candidate changes still require diff, version, safety, and test review |
| **Total active dependency/build PR candidates** | **12** | Six superseded/obsolete PRs moved to [closed.md](closed.md); excludes #2326 |

## Open PRs to review

| Upstream PR | Title | Pinned version or preliminary difference | RNC disposition |
| --- | --- | --- | --- |
| [#2323](https://github.com/gitleaks/gitleaks/pull/2323) | build(deps): bump docker/login-action from 4.6.0 to 4.7.0 in the actions group | Docker login action 4.6.0 → 4.7.0 | UNREVIEWED / RELEVANCE CHECK |
| [#2289](https://github.com/gitleaks/gitleaks/pull/2289) | build(deps): bump github.com/wasilibs/go-re2 from 1.9.0 to 1.13.0 | go-re2 v1.9.0 → v1.13.0 | UNREVIEWED / RELEVANCE CHECK |
| [#2250](https://github.com/gitleaks/gitleaks/pull/2250) | build(deps): bump github.com/stretchr/testify from 1.10.0 to 1.12.1 | testify v1.10.0 → v1.12.1 | UNREVIEWED / RELEVANCE CHECK |
| [#2149](https://github.com/gitleaks/gitleaks/pull/2149) | build(deps): bump github.com/spf13/cobra from 1.9.1 to 1.10.2 | cobra v1.9.1 → v1.10.2 | UNREVIEWED / RELEVANCE CHECK |
| [#2147](https://github.com/gitleaks/gitleaks/pull/2147) | build(deps): bump github.com/hashicorp/go-version from 1.7.0 to 1.9.0 | hashicorp/go-version v1.7.0 → v1.9.0 | UNREVIEWED / RELEVANCE CHECK |
| [#2145](https://github.com/gitleaks/gitleaks/pull/2145) | build(deps): bump github.com/charmbracelet/lipgloss from 0.5.0 to 1.1.0 | lipgloss v0.5.0 → v1.1.0 | UNREVIEWED / RELEVANCE CHECK |
| [#2144](https://github.com/gitleaks/gitleaks/pull/2144) | build(deps): bump github.com/rs/zerolog from 1.33.0 to 1.35.1 | zerolog v1.33.0 → v1.35.1 | UNREVIEWED / RELEVANCE CHECK |
| [#2141](https://github.com/gitleaks/gitleaks/pull/2141) | build(deps): bump github.com/fatih/semgroup from 1.2.0 to 1.3.0 | semgroup v1.2.0 → v1.3.0 | UNREVIEWED / RELEVANCE CHECK |
| [#2139](https://github.com/gitleaks/gitleaks/pull/2139) | build(deps): bump github.com/spf13/viper from 1.19.0 to 1.21.0 | viper v1.19.0 → v1.21.0 | UNREVIEWED / RELEVANCE CHECK |
| [#2115](https://github.com/gitleaks/gitleaks/pull/2115) | chore: bump Go to 1.26.3 to fix vulnerabilities | RNC hardened Docker builder uses Go 1.27.3, **but** `master/go.mod` still says 1.24.11 and `.github/workflows/test.yml` still requests Go 1.25 | **PARTIAL — REVIEW CI / SOURCE-BUILD SCOPE** |

| [#2100](https://github.com/gitleaks/gitleaks/pull/2100) | chore: Upgrade lipgloss to v2.0.3 | lipgloss v0.5.0 → proposed v2.0.3 (major API migration) | UNREVIEWED / RELEVANCE CHECK |
| [#2035](https://github.com/gitleaks/gitleaks/pull/2035) | chore(go.mod): update go version and dependencies | Go builder, `mholt/archives v0.1.5`, and many module updates already supersede the proposal, **but** the PR proposes `go-re2 v1.10.0` while master retains v1.9.0; related PR #2289 proposes v1.13.0 | **PARTIAL — REVIEW REGEX BACKEND / MODULE GRAPH** |





**Critical boundary:** Go toolchain updates do **not** automatically upgrade dependencies in `go.mod` or `go.sum`. The fork's `master` still pins `klauspost/compress v1.18.0`, `go-re2 v1.9.0`, `cobra v1.9.1`, `viper v1.19.0`, and other older modules. A new Go builder does not prove these PRs are redundant.

### Verified RNC supersession — moved out of this active register

- [PR #1960](closed.md) and [PR #1991](closed.md) targeted older Go and/or Alpine Docker bases. The **separate** RNC `fix/docker-security-bases` Dockerfile now builds with `golang:1.27.3` and runs on `alpine:3.24`; both targets are surpassed. Their Docker changes are not yet on fork `master`.
- [PR #1959](closed.md) requested `go 1.24.9`; fork `master/go.mod` already specifies `go 1.24.11`. This supersedes the exact requested minimum version, independently of Docker.
- [PR #1103](closed.md) requested `yaml.v3 v3.0.0`; fork `master/go.mod` already requires `v3.0.1`.
- [PR #2103](closed.md) aimed to keep CVE-2025-68121 out of **future Docker-built binaries** by raising the minimum Go version. The RNC hardened Docker builder uses Go 1.27.3, superseding that *image-build remediation*. The exact `go.mod` directive change was **not** adopted; non-Docker binaries and source-build minimums require separate verification.
- [PR #2115](https://github.com/gitleaks/gitleaks/pull/2115) remains active because its CI Go version and `go.mod` minimum are not mirrored; [PR #2035](https://github.com/gitleaks/gitleaks/pull/2035) remains active because its Go-RE2 and module-graph changes are not fully matched.
- **Evidence:** [`master/go.mod`](https://github.com/RNC-Solutions-Group/gitleaks/blob/master/go.mod) · [`fix/docker-security-bases/Dockerfile`](https://github.com/RNC-Solutions-Group/gitleaks/blob/fix/docker-security-bases/Dockerfile) · [`master/.github/workflows/test.yml`](https://github.com/RNC-Solutions-Group/gitleaks/blob/master/.github/workflows/test.yml) · original upstream PR patches (reviewed 2026-11-08). No image SBOM or runtime regression was newly generated for this routing decision.

A proposal superseded for **RNC's hardened Docker build** is not necessarily fixed in upstream releases, or in the fork's untouched `master`/non-Docker release pipeline. The closed ledger preserves reopen conditions.

## Open issues requiring dependency/build triage

These came from the open-issue inventory plus targeted build and dependency searches. CVE work owned by [security-hardening](security-hardening.md) is excluded from this table.

| Issue | Description | Routing |
| --- | --- | --- |
| [#1120](https://github.com/gitleaks/gitleaks/issues/1120) | Please bump golang.org/x/sys dependency | Check x/sys platform compatibility against the fork’s pinned version |
| [#2031](https://github.com/gitleaks/gitleaks/issues/2031) | Rename in the go.mod file? | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#1839](https://github.com/gitleaks/gitleaks/issues/1839) | Module path discrepancy after repository migration | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#1305](https://github.com/gitleaks/gitleaks/issues/1305) | git version 1.8.3.1 does not support git -C which the latest version of gitleaks uses. date(2023-11-30).   | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#1226](https://github.com/gitleaks/gitleaks/issues/1226) | `go run ...` module declares its path as: github.com/zricethezav/gitleaks | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#1209](https://github.com/gitleaks/gitleaks/issues/1209) | Cannot build from source: incorrect go.mod | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#833](https://github.com/gitleaks/gitleaks/issues/833) | Undocumented golang and git version requirements | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#1944](https://github.com/gitleaks/gitleaks/issues/1944) | Build release binary in GitHub Actions | Release reproducibility / build pipeline |

### Title-search exclusions, not dependency remediation

The search also returned issues about `--version` output, config rule versioning, `useDefault`, and JSON reporting, plus #2326. These are not package-upgrade tasks; leave them for feature/CLI categories rather than inflate a dependency closure rate.

## Recommended filter before closing anything else

1. **Scope:** build toolchain, Go module, container OS package, GitHub Action, or unrelated functionality.
2. **Exact pin:** current `go.mod/go.sum`, Dockerfile, pinned Actions SHA, and binary `go version -m`.
3. **Actual image:** digest, platform, SBOM/package version, Trivy/Scout database timestamp, and any CVEs still detected.
4. **Behavior:** Go tests and Docker acceptance, including archive/chunking regressions for compression libraries.
5. **Disposition:** move into `closed.md` **only after an explicit RNC closure decision**, preserving the upstream issue/PR as a link; do not auto-close GitHub issues.

The dependency closure rate should be calculated from verified package and image evidence, not inferred from the Go builder version alone.

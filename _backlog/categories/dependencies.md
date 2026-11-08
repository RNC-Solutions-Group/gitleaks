# Dependencies — Gitleaks candidate register

Snapshot source: [Open PR inventory](../02-open-pull-requests.md), [issue inventory](../01-issue-backlog.md) (2026-11-08). This is a **title/search-based first pass**, enriched with pinned version checks from the fork. It is not an SBOM comparison against every image.

**Ownership rule:** [closed.md](closed.md) is the canonical RNC exclusion list. Items recorded there are omitted from active tables here. Security PRs and their linked issues remain visible as references below, but their active review and test disposition is owned by [security-hardening.md](security-hardening.md).

## Summary

| PR disposition | Count | Basis |
| --- | ---: | --- |
| Reviewed — testing in security-hardening | 2 | [#2308](https://github.com/gitleaks/gitleaks/pull/2308), [#2336](https://github.com/gitleaks/gitleaks/pull/2336) |
| Locally superseded candidate | 6 | RNC builder / current pinned module already newer than the proposal; **not equivalent to an upstream merge** |
| Unreviewed / relevance check | 13 | Candidate changes still require diff, version, safety, and test review |
| **Total title-based dependency/build PR candidates** | **21** | Excludes [#2326](https://github.com/gitleaks/gitleaks/pull/2326) (documentation mentioning a report option dependency, not package updates) |

## Open PRs to review

| Upstream PR | Title | Pinned version or preliminary difference | RNC disposition |
| --- | --- | --- | --- |
| [#2336](https://github.com/gitleaks/gitleaks/pull/2336) | build(deps): bump klauspost/compress to v1.18.7 | master: klauspost/compress v1.18.0 → proposed v1.18.7 | REVIEWED — TESTING · security-hardening |
| [#2323](https://github.com/gitleaks/gitleaks/pull/2323) | build(deps): bump docker/login-action from 4.6.0 to 4.7.0 in the actions group | Docker login action 4.6.0 → 4.7.0 | UNREVIEWED / RELEVANCE CHECK |
| [#2308](https://github.com/gitleaks/gitleaks/pull/2308) | build(deps): update dependencies to resolve CVEs (#2305, #2291) | master: x/crypto v0.35.0, x/sys v0.30.0, x/text v0.29.0; PR raises Go minimum | REVIEWED — TESTING · security-hardening |
| [#2289](https://github.com/gitleaks/gitleaks/pull/2289) | build(deps): bump github.com/wasilibs/go-re2 from 1.9.0 to 1.13.0 | go-re2 v1.9.0 → v1.13.0 | UNREVIEWED / RELEVANCE CHECK |
| [#2250](https://github.com/gitleaks/gitleaks/pull/2250) | build(deps): bump github.com/stretchr/testify from 1.10.0 to 1.12.1 | testify v1.10.0 → v1.12.1 | UNREVIEWED / RELEVANCE CHECK |
| [#2149](https://github.com/gitleaks/gitleaks/pull/2149) | build(deps): bump github.com/spf13/cobra from 1.9.1 to 1.10.2 | cobra v1.9.1 → v1.10.2 | UNREVIEWED / RELEVANCE CHECK |
| [#2147](https://github.com/gitleaks/gitleaks/pull/2147) | build(deps): bump github.com/hashicorp/go-version from 1.7.0 to 1.9.0 | hashicorp/go-version v1.7.0 → v1.9.0 | UNREVIEWED / RELEVANCE CHECK |
| [#2145](https://github.com/gitleaks/gitleaks/pull/2145) | build(deps): bump github.com/charmbracelet/lipgloss from 0.5.0 to 1.1.0 | lipgloss v0.5.0 → v1.1.0 | UNREVIEWED / RELEVANCE CHECK |
| [#2144](https://github.com/gitleaks/gitleaks/pull/2144) | build(deps): bump github.com/rs/zerolog from 1.33.0 to 1.35.1 | zerolog v1.33.0 → v1.35.1 | UNREVIEWED / RELEVANCE CHECK |
| [#2141](https://github.com/gitleaks/gitleaks/pull/2141) | build(deps): bump github.com/fatih/semgroup from 1.2.0 to 1.3.0 | semgroup v1.2.0 → v1.3.0 | UNREVIEWED / RELEVANCE CHECK |
| [#2139](https://github.com/gitleaks/gitleaks/pull/2139) | build(deps): bump github.com/spf13/viper from 1.19.0 to 1.21.0 | viper v1.19.0 → v1.21.0 | UNREVIEWED / RELEVANCE CHECK |
| [#2115](https://github.com/gitleaks/gitleaks/pull/2115) | chore: bump Go to 1.26.3 to fix vulnerabilities | RNC Docker builder already Go 1.27.3; upstream go.mod still 1.24.11 | LOCAL SUPERSEDED — verify scope |
| [#2103](https://github.com/gitleaks/gitleaks/pull/2103) | build: update Go toolchain for CVE-2025-68121 | RNC Docker builder already Go 1.27.3; upstream go.mod still 1.24.11 | LOCAL SUPERSEDED — verify scope |
| [#2100](https://github.com/gitleaks/gitleaks/pull/2100) | chore: Upgrade lipgloss to v2.0.3 | lipgloss v0.5.0 → proposed v2.0.3 (major API migration) | UNREVIEWED / RELEVANCE CHECK |
| [#2035](https://github.com/gitleaks/gitleaks/pull/2035) | chore(go.mod): update go version and dependencies | master archives v0.1.5 already; RNC builder Go 1.27.3 | LOCAL SUPERSEDED — verify scope |
| [#1991](https://github.com/gitleaks/gitleaks/pull/1991) | Update Dockerfile | Generic Dockerfile change; manual diff needed | UNREVIEWED / RELEVANCE CHECK |
| [#1960](https://github.com/gitleaks/gitleaks/pull/1960) | Upgrade Go and Alpine versions in Dockerfile | RNC Dockerfile already Go 1.27.3 / Alpine 3.24 | LOCAL SUPERSEDED — verify scope |
| [#1959](https://github.com/gitleaks/gitleaks/pull/1959) | Update go to 1.24.9 in go.mod | master go.mod already 1.24.11, newer than proposed 1.24.9 | LOCAL SUPERSEDED — verify scope |
| [#1103](https://github.com/gitleaks/gitleaks/pull/1103) | fix(sec): upgrade gopkg.in/yaml.v3 to 3.0.0 | master yaml.v3 already v3.0.1, newer than proposed v3.0.0 | LOCAL SUPERSEDED — verify scope |
| [#1099](https://github.com/gitleaks/gitleaks/pull/1099) | Fix: go.mod version changed back to 1.18 | Version-compatibility change, not a vulnerability dependency upgrade | UNREVIEWED / RELEVANCE CHECK |
| [#652](https://github.com/gitleaks/gitleaks/pull/652) | Minimum go version should be enforced when installing gitleaks with go get command | Installation Go-version enforcement, not a module upgrade | UNREVIEWED / RELEVANCE CHECK |

**Critical boundary:** Go toolchain updates do **not** automatically upgrade dependencies in `go.mod` or `go.sum`. The fork's `master` still pins `klauspost/compress v1.18.0`, `go-re2 v1.9.0`, `cobra v1.9.1`, `viper v1.19.0`, and other older modules. A new Go builder does not prove these PRs are redundant.

### Evidence used to mark six locally superseded candidates

- [`master/go.mod`](https://github.com/RNC-Solutions-Group/gitleaks/blob/master/go.mod): `go 1.24.11`; `mholt/archives v0.1.5`; `gopkg.in/yaml.v3 v3.0.1`.
- [`fix/docker-security-bases/Dockerfile`](https://github.com/RNC-Solutions-Group/gitleaks/blob/fix/docker-security-bases/Dockerfile): build image `golang:1.27.3`, runtime `alpine:3.24`.
- Toolchain-only PRs proposing Go 1.24.9, 1.24.13, 1.25.7, or 1.26.3 are locally superseded for **that builder**. However, source `go.mod` changes and release pipeline impacts require independent review.
- The YAML v3 proposal targeted v3.0.0; the current source already pins v3.0.1. No need to reapply the older update.
- A proposal being superseded in **RNC's hardened build** does not mean it is fixed in upstream releases or the fork's untouched `master`.

## Open issues requiring dependency/build triage

These came from the open-issue inventory plus targeted searches for CVE and build terms. The first three are owned by [security-hardening](security-hardening.md), not separately evaluated here.

| Issue | Description | Routing |
| --- | --- | --- |
| [#2305](https://github.com/gitleaks/gitleaks/issues/2305) | Please Resolve CVEs in 8.30.1 | REVIEWED — TESTING via security-hardening |
| [#2291](https://github.com/gitleaks/gitleaks/issues/2291) | Go CVE Vulnerabilities | REVIEWED — TESTING via security-hardening |
| [#2330](https://github.com/gitleaks/gitleaks/issues/2330) | fix below CVEs | REVIEWED — TESTING via security-hardening |
| [#1120](https://github.com/gitleaks/gitleaks/issues/1120) | Please bump golang.org/x/sys dependency | Check x/sys platform compatibility against #2308 |
| [#2031](https://github.com/gitleaks/gitleaks/issues/2031) | Rename in the go.mod file? | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#1839](https://github.com/gitleaks/gitleaks/issues/1839) | Module path discrepancy after repository migration | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#1544](https://github.com/gitleaks/gitleaks/issues/1544) | invalid go version '1.22.0': must match format 1.23 | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#1305](https://github.com/gitleaks/gitleaks/issues/1305) | git version 1.8.3.1 does not support git -C which the latest version of gitleaks uses. date(2023-11-30).   | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#1226](https://github.com/gitleaks/gitleaks/issues/1226) | `go run ...` module declares its path as: github.com/zricethezav/gitleaks | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#1209](https://github.com/gitleaks/gitleaks/issues/1209) | Cannot build from source: incorrect go.mod | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#1057](https://github.com/gitleaks/gitleaks/issues/1057) | Go version 1.19 is unsuported by tidy ( requires 1.19 ) | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#1029](https://github.com/gitleaks/gitleaks/issues/1029) | Documentation: installation of dependencies when installing from source. | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#833](https://github.com/gitleaks/gitleaks/issues/833) | Undocumented golang and git version requirements | Build, module-path, compatibility, or documentation; not established as an outstanding CVE |
| [#1944](https://github.com/gitleaks/gitleaks/issues/1944) | Build release binary in GitHub Actions | Release reproducibility / build pipeline |

### Title-search exclusions, not dependency remediation

The search also returned issues about `--version` output, config rule versioning, `useDefault`, and JSON reporting, plus [#2326](https://github.com/gitleaks/gitleaks/pull/2326). These are not package-upgrade tasks; leave them for feature/CLI categories rather than inflate a dependency closure rate.

## Recommended filter before closing anything else

1. **Scope:** build toolchain, Go module, container OS package, GitHub Action, or unrelated functionality.
2. **Exact pin:** current `go.mod/go.sum`, Dockerfile, pinned Actions SHA, and binary `go version -m`.
3. **Actual image:** digest, platform, SBOM/package version, Trivy/Scout database timestamp, and any CVEs still detected.
4. **Behavior:** Go tests and Docker acceptance, including archive/chunking regressions for compression libraries.
5. **Disposition:** move into `closed.md` **only after an explicit RNC closure decision**, preserving the upstream issue/PR as a link; do not auto-close GitHub issues.

This is where the easy wins should emerge. It would be premature to claim 90% of these PRs are closed from a newer Go version alone.

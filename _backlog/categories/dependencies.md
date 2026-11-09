# Dependencies — Gitleaks candidate register

Snapshot source: [Open PR inventory](../02-open-pull-requests.md), [issue inventory](../01-issue-backlog.md) (2026-11-08). **8 active issues · 12 active PRs · 20 work packages.** Candidate dependency and build-pin changes, enriched with fork-pinned version checks; no current-image SBOM comparison is implied.

Security CVE work belongs to [security-hardening](security-hardening.md). This active register excludes items recorded in [RNC exclusions](exclusions.md).

## Summary

| PR disposition | Count | Basis |
| --- | ---: | --- |
| Partially superseded — further review | 2 | Docker builder or some modules are newer, but each PR still contains unresolved source/CI or dependency changes |
| Unreviewed / relevance check | 10 | Candidate changes still require diff, version, safety, and test review |
| **Total active dependency/build PR candidates** | **12** | Open-snapshot PRs requiring build/dependency review |

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

## Open issues requiring dependency/build triage

These came from the open-issue inventory plus targeted build and dependency searches. CVE work owned by [security-hardening](security-hardening.md) is excluded from this table.

| Issue | Description | Routing |
| --- | --- | --- |
| [#1120](https://github.com/gitleaks/gitleaks/issues/1120) | Update `golang.org/x/sys` to support FreeBSD/RISC-V. | **PENDING — PLATFORM MATRIX.** Master pins `x/sys v0.30.0` (not the 2023 version). No FreeBSD/RISC-V build, supported-architecture matrix, or target-specific proof. Do not close based only on the version being newer. |
| [#1839](https://github.com/gitleaks/gitleaks/issues/1839) | Canonical module-path discrepancy after repository migration. | **ACTIVE — REPRO SUPPLIED.** Master `go.mod` still declares `github.com/zricethezav/gitleaks/v8`; report demonstrates `go get github.com/gitleaks/gitleaks/v8@...` fails. Assess compatibility consequences before proposing a rename. |
| [#1305](https://github.com/gitleaks/gitleaks/issues/1305) | Git 1.8.3.1 rejects `git -C`. | **PENDING — SUPPORT POLICY.** `master/sources/git.go` still invokes `git -C` in scan paths. Historic Git 1.8 is not proof of support, but the precise incompatibility remains; do not claim corrected until minimum Git support is documented or tested. Related #833. |
| [#1226](https://github.com/gitleaks/gitleaks/issues/1226) | 2023 `go run github.com/gitleaks/gitleaks@latest` failure. | **PENDING — VERSIONED COMMAND.** Shares repository/module naming root with #1839, but uses **unversioned** `/gitleaks@latest` and resolves legacy `v1.25.1`; changing a `/v8` module path alone would not validate that exact command. Keep a distinct install/migration test case. |
| [#1209](https://github.com/gitleaks/gitleaks/issues/1209) | 2023 `go install github.com/gitleaks/gitleaks@latest` failure. | **PENDING — VERSIONED COMMAND.** Same historical unversioned-module problem as #1226, with a different command; assess current `/v8` installation advice. Do not misrepresent a module rename as proof the old command is supported. |
| [#833](https://github.com/gitleaks/gitleaks/issues/833) | Publish minimum Git and Go build prerequisites. | **ACTIVE — DOC REQUIREMENT.** `master/go.mod` explicitly requires Go 1.24.11, but root README build instructions only say “make sure go is installed” and do not state the Git floor. 2022 RHEL/Go1.16 case is historical; documentation gap survives. |
| [#1944](https://github.com/gitleaks/gitleaks/issues/1944) | Reproducible release binaries built by GitHub Actions. | **ACTIVE — SEPARATE PIPELINE GAP.** Reporter supplies a `SIGILL` on old x86 and observations indicating a different binary build environment than Docker. Current `.github/workflows/release.yml` publishes **containers only**, not standalone release binaries. Do not treat a Docker build as a fix. |
| [#1588](https://github.com/gitleaks/gitleaks/issues/1588) | Normalize source builds and installation, reduce ad-hoc linker flags. | **PENDING REVIEW.** This broader build request mentions module rename like #1839, **plus** simplifying ldflags and `go install`. Keep the distinct normalization scope here; it is not wholly superseded by the module-path defect. |


## Review gate

1. **Scope:** build toolchain, Go module, container OS package, GitHub Action, or unrelated functionality.
2. **Exact pin:** current `go.mod/go.sum`, Dockerfile, pinned Actions SHA, and binary `go version -m`.
3. **Actual image:** digest, platform, SBOM/package version, Trivy/Scout database timestamp, and any CVEs still detected.
4. **Behavior:** Go tests and Docker acceptance, including archive/chunking regressions for compression libraries.
5. **Disposition:** move into `exclusions.md` **only after an explicit RNC closure decision**, preserving the upstream issue/PR as a link; do not auto-close GitHub issues.

Use verified package/image evidence for any dependency disposition. A newer builder tag alone does not establish remediation.

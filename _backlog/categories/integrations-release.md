# Integrations Release — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**3 issues · 6 PRs · 8 work packages.**

CI, installations, releases, platforms, images, and packaging, outside the separately managed dependency pin review.

**Review states:** `SOURCE-CHECKED` records a targeted source/patch comparison, **not** a built or tested correction. `PENDING` records an unresolved issue or competing proposal. No PR is adopted or upstream item closed by categorization. See [RNC exclusions](exclusions.md) for retired review slots.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| [#1984](https://github.com/gitleaks/gitleaks/issues/1984) | [#2081](https://github.com/gitleaks/gitleaks/pull/2081) (@Smeet23) | Docker `--version` placeholder; PR injects the release ref as a Docker build argument. | **SOURCE-CHECKED — RETAIN, TEST PENDING** | The v8.29.0 regression received an upstream Makefile fix/release v8.29.1, but RNC master Dockerfile still uses `git describe` and its release workflow does not supply a version argument. #2081 modifies both; compare with the separate hardened Dockerfile and test final image metadata. | — |
| — | [#2220](https://github.com/gitleaks/gitleaks/pull/2220) (@kobihikri) | BuildKit provenance and SBOM attestations. | **SOURCE-CHECKED — RETAIN** | Master release workflow lacks explicit `provenance`/`sbom` options; PR adds both. No published image attestation was checked. Distinct from Cosign signing in #1217. | — |
| — | [#2209](https://github.com/gitleaks/gitleaks/pull/2209) (@trinathone) | Pull-request scan trigger with fork-secret skip. | **SOURCE-CHECKED — RETAIN** | Master workflow only triggers on push/dispatch and restricts the job to `gitleaks/gitleaks`; PR adds `pull_request` but remains upstream-repository-specific. No proven RNC-fork CI fix or adoption. Older #1915 remains excluded. | — |
| — | [#1661](https://github.com/gitleaks/gitleaks/pull/1661) (@ianzone) | npm binary distribution and publication workflow. | **SOURCE-CHECKED — RETAIN** | Eight-file workflow/JS/package addition with `NPM_TOKEN` publishing; no npm release path currently exists in master. New distribution/security support surface, with narrower platform coverage than GoReleaser. | — |
| — | [#1485](https://github.com/gitleaks/gitleaks/pull/1485) (@suhovius) | Bash installation downloader and README `curl`/`wget` pipes. | **SOURCE-CHECKED — RETAIN** | Adds executable installer; master has no `install.sh`. Requires checksum/signature verification and shell/platform policy review before any implementation decision. | — |
| [#1127](https://github.com/gitleaks/gitleaks/issues/1127) | — | PowerShell `Start-Process` invocation of Gitleaks 8.14.1 in a Windows Server 2016/Azure DevOps build reports an error after an apparently successful scan. | **CI/Windows INTEGRATION — REPRO PENDING** | Issue supplies PowerShell commands and two historical screenshots (error text unavailable for independent verification). Master still exposes deprecated `detect --no-git`, whereas modern `dir` scanning differs; neither proves the observed failure was in scanner code rather than process handling or CI output. Keep open; require transcribed error, exit code, captured stderr/stdout, and current `dir` reproduction. | — |
| — | [#1217](https://github.com/gitleaks/gitleaks/pull/1217) (@becojo) | Cosign container image signatures. | **SOURCE-CHECKED — RETAIN** | PR introduces OIDC `id-token` and Cosign installation/signing to release workflow; neither is in master. Not interchangeable with BuildKit attestation PR #2220. | — |
| [#885](https://github.com/gitleaks/gitleaks/issues/885) | — | Moving major/minor Docker tags (`vX`, `vX.Y`). | **SOURCE-CHECKED — RETAIN** | Release workflow delegates tags to `docker/metadata-action` without explicit major/minor semver patterns; no registry inventory validated the requested tags. | — |

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

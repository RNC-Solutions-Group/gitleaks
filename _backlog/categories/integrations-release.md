# Integrations Release — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**3 active upstream issues · 6 upstream PRs · 8 work packages.** Additional action/permission reports #2063, #1913, #1343, host Podman report #1385, and competing PR #1915 have evidence-backed [RNC exclusions](exclusions.md). Cross-platform baseline issue #1893 moved to [scan integrity](scan-integrity.md); broad build-normalization request #1588 moved to [dependencies](dependencies.md). Docker version issue #1984 / PR #2081 moved here from Git & Hooks; the upstream-only orphaned-tag report #2086 is now a documented [RNC exclusion](exclusions.md). Former Needs Review issue #1127 is now routed here as a Windows/Azure DevOps scan-invocation problem; its cause remains open for reproduction. This focused scrub excluded four issue slots plus the issue-linked SDK-template PR, and transferred configuration policy issue #1754 to [Configuration](configuration-allowlists.md). Upstream state and scanner code are unchanged.

CI, installations, releases, platforms, images and packaging outside the already separated dependency pin register. Current master and hardened Docker branch must not be conflated.

**Evidence rules:** “SOURCE-CHECKED” means a targeted patch/current-source comparison, **not** that a fix was built or tested. Most records are title/crosswalk-only and must stay pending until reproduced, diffed, and assessed on the fork. “Obsolete,” “superseded,” “adopted,” or “corrected” require separate evidence; no PR is adopted by routing it here.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| [#1984](https://github.com/gitleaks/gitleaks/issues/1984) | [#2081](https://github.com/gitleaks/gitleaks/pull/2081) (@Smeet23) | Docker `--version` placeholder; PR injects the release ref as a Docker build argument. | **SOURCE-CHECKED — RETAIN, TEST PENDING** | The v8.29.0 regression received an upstream Makefile fix/release v8.29.1, but RNC master Dockerfile still uses `git describe` and its release workflow does not supply a version argument. #2081 modifies both; compare with the separate hardened Dockerfile and test final image metadata. | — |
| — | [#2220](https://github.com/gitleaks/gitleaks/pull/2220) (@kobihikri) | BuildKit provenance and SBOM attestations. | **SOURCE-CHECKED — RETAIN** | Master release workflow lacks explicit `provenance`/`sbom` options; PR adds both. No published image attestation was checked. Distinct from Cosign signing in #1217. | — |
| — | [#2209](https://github.com/gitleaks/gitleaks/pull/2209) (@trinathone) | Pull-request scan trigger with fork-secret skip. | **SOURCE-CHECKED — RETAIN** | Master workflow only triggers on push/dispatch and restricts the job to `gitleaks/gitleaks`; PR adds `pull_request` but remains upstream-repository-specific. No proven RNC-fork CI fix or adoption. Older #1915 remains excluded. | — |
| — | [#1661](https://github.com/gitleaks/gitleaks/pull/1661) (@ianzone) | npm binary distribution and publication workflow. | **SOURCE-CHECKED — RETAIN** | Eight-file workflow/JS/package addition with `NPM_TOKEN` publishing; no npm release path currently exists in master. New distribution/security support surface, with narrower platform coverage than GoReleaser. | — |
| — | [#1485](https://github.com/gitleaks/gitleaks/pull/1485) (@suhovius) | Bash installation downloader and README `curl`/`wget` pipes. | **SOURCE-CHECKED — RETAIN** | Adds executable installer; master has no `install.sh`. Requires checksum/signature verification and shell/platform policy review before any implementation decision. | — |
| [#1127](https://github.com/gitleaks/gitleaks/issues/1127) | — | PowerShell `Start-Process` invocation of Gitleaks 8.14.1 in a Windows Server 2016/Azure DevOps build reports an error after an apparently successful scan. | **RECLASSIFIED — CI/Windows INTEGRATION; REPRO PENDING** | Issue supplies PowerShell commands and two historical screenshots (error text unavailable for independent verification). Master still exposes deprecated `detect --no-git`, whereas modern `dir` scanning differs; neither proves the observed failure was in scanner code rather than process handling or CI output. Keep open; require transcribed error, exit code, captured stderr/stdout, and current `dir` reproduction. | — |
| — | [#1217](https://github.com/gitleaks/gitleaks/pull/1217) (@becojo) | Cosign container image signatures. | **SOURCE-CHECKED — RETAIN** | PR introduces OIDC `id-token` and Cosign installation/signing to release workflow; neither is in master. Not interchangeable with BuildKit attestation PR #2220. | — |
| [#885](https://github.com/gitleaks/gitleaks/issues/885) | — | Moving major/minor Docker tags (`vX`, `vX.Y`). | **SOURCE-CHECKED — RETAIN** | Release workflow delegates tags to `docker/metadata-action` without explicit major/minor semver patterns; no registry inventory validated the requested tags. | — |

## Needs Review retirement — category ownership clarification (2026-11-08)

- **Issue #1127** is a concrete Windows CI/PowerShell scanner-invocation report rather than an unclassifiable request. Retained **here** without claiming the 2023 error is fixed or Azure DevOps caused it. Its screenshots could not be retrieved; no runtime reproduction was performed.
- Former Needs Review issue [#1864](exclusions.md) is a GitLab-owned secret-detection ruleset concern, now recorded in RNC exclusions. The catch-all category has no remaining owner entries.

## Focused Integrations & Releases scrub — 2026-11-08

**Baseline:** RNC fork `master` at `5f6ec56c`. Eight upstream issue bodies, seven open PR diffs, issue comments, and current release/build/config sources were reviewed. **No builds, releases, tag mutations, or code integrations.**

- **Excluded:** issues #1365 (Go prerequisite), #1210 (old Debian downstream packaging), #2086 (upstream orphaned tag outside the fork's tag refs), and issue #2155 with PR #2156 (upstream SDK issue-intake template). The canonical [exclusions ledger](exclusions.md) records evidence and reopen triggers.
- **Moved:** #1754 to [Configuration](configuration-allowlists.md). `minVersion` exists at the top level of the current config, but `validateMinVersion` merely **warns**, not enforces a minimum.
- **Retained:** three issues and six PRs in eight packages. None is adopted or presumed fixed by a newer release name.

## Focused review: preserved distinctions

- **#2220 vs #1217:** the release PRs change different supply-chain controls. #2220 enables provenance/SBOM attestations in Buildx; #1217 adds Cosign signing. Neither independently supersedes the other, and neither exists in `master/.github/workflows/release.yml` today. **Both pending.**
- **#2209 vs #1915:** master already has a release-pinned Gitleaks scan workflow, but no PR event trigger. #1915's older unpinned-action rewrite is redundant/unsafe relative to #2209's retained approach. **#2209 pending**, not merged or proven safe for forks.
- **#2155/#2156:** excluded as upstream-only SDK issue intake; no runtime scanner or release change. Reopen if RNC launches an SDK request channel.
- **#1661, #1485:** npm distribution and an executable `curl|bash` installer expand release/distribution surfaces; neither is a doc-only change. **Pending separate security/support review.**
- **#1210, #1365:** excluded as downstream Debian packaging and missing-Go installation-prerequisite cases, **not** as proven scanner fixes. **#1944** release-binary provenance/`SIGILL` stays with [dependencies](dependencies.md).
- **#885:** current release workflow delegates tags to `docker/metadata-action` without explicit major/minor semver tags; no proof the requested stable major/minor tags exist. **Pending.**

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

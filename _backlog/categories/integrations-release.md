# Integrations Release — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**5 active upstream issues · 6 upstream PRs · 10 work packages.** Additional action/permission reports #2063, #1913, #1343, host Podman report #1385, and competing PR #1915 have evidence-backed [RNC exclusions](exclusions.md). Cross-platform baseline issue #1893 moved to [scan integrity](scan-integrity.md); broad build-normalization request #1588 moved to [dependencies](dependencies.md). Existing upstream state is unchanged.

CI, installations, releases, platforms, images and packaging outside the already separated dependency pin register. Current master and hardened Docker branch must not be conflated.

**Evidence rules:** “SOURCE-CHECKED” means a targeted patch/current-source comparison, **not** that a fix was built or tested. Most records are title/crosswalk-only and must stay pending until reproduced, diffed, and assessed on the fork. “Obsolete,” “superseded,” “adopted,” or “corrected” require separate evidence; no PR is adopted by routing it here.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| — | [#2220](https://github.com/gitleaks/gitleaks/pull/2220) (@kobihikri) | ci: attach provenance and SBOM attestations to the published images | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#2209](https://github.com/gitleaks/gitleaks/pull/2209) (@trinathone) | Add PR triggering with a fork-secret skip condition to this repository's `gitleaks.yml` workflow. | **SOURCE-CHECKED — REVIEW PENDING** | Current master workflow runs on `push` / manual dispatch only. PR adds `pull_request` but skips forks outright; this does **not** fix the general reusable Action license behavior in excluded issue #2063. Supersedes unsafe older workflow rewrite #1915 as a separate review slot; test event/permissions behavior before any adoption. | — |
| [#2155](https://github.com/gitleaks/gitleaks/issues/2155) | [#2156](https://github.com/gitleaks/gitleaks/pull/2156) (@MChorfa) | Add dedicated SDK embedding issue template; proposal: chore(github): add SDK embedding issue template | **LINKED — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1754](https://github.com/gitleaks/gitleaks/issues/1754) | — | Setting a minimum version of GitLeaks when useDefault = true | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1661](https://github.com/gitleaks/gitleaks/pull/1661) (@ianzone) | chore: release npm | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1485](https://github.com/gitleaks/gitleaks/pull/1485) (@suhovius) | Add install.sh bash script installation option | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1365](https://github.com/gitleaks/gitleaks/issues/1365) | — | Installing gitleaks not working | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1217](https://github.com/gitleaks/gitleaks/pull/1217) (@becojo) | Sign container using Cosign | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1210](https://github.com/gitleaks/gitleaks/issues/1210) | — | Debian and GitLeak build failures | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#885](https://github.com/gitleaks/gitleaks/issues/885) | — | Create docker tags also for major and minor versions | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |

## Focused review: preserved distinctions

- **#2220 vs #1217:** the release PRs change different supply-chain controls. #2220 enables provenance/SBOM attestations in Buildx; #1217 adds Cosign signing. Neither independently supersedes the other, and neither exists in `master/.github/workflows/release.yml` today. **Both pending.**
- **#2209 vs #1915:** master already has a release-pinned Gitleaks scan workflow, but no PR event trigger. #1915's older unpinned-action rewrite is redundant/unsafe relative to #2209's retained approach. **#2209 pending**, not merged or proven safe for forks.
- **#2155/#2156:** the proposed SDK issue template is new and absent from `.github/ISSUE_TEMPLATE/`; treat it as an optional governance decision, not an already-implemented feature. **Pending.**
- **#1661, #1485:** npm distribution and an executable `curl|bash` installer expand release/distribution surfaces; neither is a doc-only change. **Pending separate security/support review.**
- **#1210, #1365:** older Debian architecture and WSL build failures have no validated current-toolchain reproduction; keep pending rather than assuming new Go fixes them. **#1944** release-binary provenance and `SIGILL` belongs in [dependencies](dependencies.md), not this workflow-only package set.
- **#885:** current release workflow delegates tags to `docker/metadata-action` without explicit major/minor semver tags; no proof the requested stable major/minor tags exist. **Pending.**

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

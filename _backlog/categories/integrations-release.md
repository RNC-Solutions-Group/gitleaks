# Integrations Release — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**14 upstream issues · 7 upstream PRs · 19 work packages.** Each ID is owned by this category once; linked PRs and issues share a row. Existing upstream state is unchanged.

CI, installations, releases, platforms, images and packaging outside the already separated dependency pin register. Current master and hardened Docker branch must not be conflated.

**Evidence rules:** “SOURCE-CHECKED” means a targeted patch/current-source comparison, **not** that a fix was built or tested. Most records are title/crosswalk-only and must stay pending until reproduced, diffed, and assessed on the fork. “Obsolete,” “superseded,” “adopted,” or “corrected” require separate evidence; no PR is adopted by routing it here.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| — | [#2220](https://github.com/gitleaks/gitleaks/pull/2220) (@kobihikri) | ci: attach provenance and SBOM attestations to the published images | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#2063](https://github.com/gitleaks/gitleaks/issues/2063) | [#2209](https://github.com/gitleaks/gitleaks/pull/2209) (@trinathone) | Org. repo with licence + PR from fork = missing license key error; proposal: fix: Org. repo with licence + PR from fork = missing license key error | **LINKED — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#2188](https://github.com/gitleaks/gitleaks/issues/2188) | — | [Notice] FYI I submitted a schema for gitleaks to schema store | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#2155](https://github.com/gitleaks/gitleaks/issues/2155) | [#2156](https://github.com/gitleaks/gitleaks/pull/2156) (@MChorfa) | Add dedicated SDK embedding issue template; proposal: chore(github): add SDK embedding issue template | **LINKED — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1915](https://github.com/gitleaks/gitleaks/pull/1915) (@SGreeshma) | ci: add GitHub Actions workflow to run Gitleaks on push/PR | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1913](https://github.com/gitleaks/gitleaks/issues/1913) | — | GitLeaks to support on merge_group | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1893](https://github.com/gitleaks/gitleaks/issues/1893) | — | The results generated on the Windows platform and the Linux platform are inconsistent. | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1754](https://github.com/gitleaks/gitleaks/issues/1754) | — | Setting a minimum version of GitLeaks when useDefault = true | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1661](https://github.com/gitleaks/gitleaks/pull/1661) (@ianzone) | chore: release npm | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1588](https://github.com/gitleaks/gitleaks/issues/1588) | — | Request: Normalize build process | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1485](https://github.com/gitleaks/gitleaks/pull/1485) (@suhovius) | Add install.sh bash script installation option | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1465](https://github.com/gitleaks/gitleaks/issues/1465) | — | ubi8/python 3.12 emits RuntimeWarning for  tarfile | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1385](https://github.com/gitleaks/gitleaks/issues/1385) | — | lsetxattr: permission denied with podman | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1365](https://github.com/gitleaks/gitleaks/issues/1365) | — | Installing gitleaks not working | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1343](https://github.com/gitleaks/gitleaks/issues/1343) | — | Resource not accessible by integration | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1257](https://github.com/gitleaks/gitleaks/issues/1257) | — | Failed with Self hosted GitHub Runner | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1217](https://github.com/gitleaks/gitleaks/pull/1217) (@becojo) | Sign container using Cosign | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1210](https://github.com/gitleaks/gitleaks/issues/1210) | — | Debian and GitLeak build failures | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#885](https://github.com/gitleaks/gitleaks/issues/885) | — | Create docker tags also for major and minor versions | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

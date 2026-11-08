# Documentation — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**7 upstream issues · 17 upstream PRs · 22 work packages.** Each ID is owned by this category once; linked PRs and issues share a row. Existing upstream state is unchanged.

Documentation, setup notes, labels, and explanatory examples. Older PRs require line-by-line comparison; a modern README does not prove every old doc fix is present.

**Evidence rules:** “SOURCE-CHECKED” means a targeted patch/current-source comparison, **not** that a fix was built or tested. Most records are title/crosswalk-only and must stay pending until reproduced, diffed, and assessed on the fork. “Obsolete,” “superseded,” “adopted,” or “corrected” require separate evidence; no PR is adopted by routing it here.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| [#1862](https://github.com/gitleaks/gitleaks/issues/1862) | [#2298](https://github.com/gitleaks/gitleaks/pull/2298) (@bharathvemula) | The READMe.md mentions to run a deprecated command; proposal: docs: replace deprecated detect wording in baseline guide | **LINKED — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#2287](https://github.com/gitleaks/gitleaks/pull/2287) (@ump45nose) | docs: fix anchor for official Docker images link in README | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#2282](https://github.com/gitleaks/gitleaks/pull/2282) (@haimingZZ) | docs: fix typos in CONTRIBUTING.md and README.md | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#2256](https://github.com/gitleaks/gitleaks/issues/2256) | [#2258](https://github.com/gitleaks/gitleaks/pull/2258) (@jankesec) | Typo in README: identiifer should be identifier; proposal: docs: fix identifier typo in README | **SOURCE-CHECKED — TEST PENDING** | Patch inspected: the typo `identiifer` is still present in `master/README.md`. | — |
| — | [#2247](https://github.com/gitleaks/gitleaks/pull/2247) (@CAOShurong) | docs: replace dead blog.gitleaks.io link in README | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#2202](https://github.com/gitleaks/gitleaks/pull/2202) (@rylena) | docs: fix baseline typo | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#2134](https://github.com/gitleaks/gitleaks/pull/2134) (@santoshyadavdev) | chore: update logo | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#2127](https://github.com/gitleaks/gitleaks/pull/2127) (@StephenJarso) | docs: improve README phrasing and fix typographical errors | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#2084](https://github.com/gitleaks/gitleaks/pull/2084) (@vincent067) | docs: add enhanced GitHub Actions example with security best practices | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#2008](https://github.com/gitleaks/gitleaks/pull/2008) (@jylenhof) | docs(README): add mise alternative method installation | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1900](https://github.com/gitleaks/gitleaks/pull/1900) (@vignesh1507) | Update profile.sh | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1855](https://github.com/gitleaks/gitleaks/issues/1855) | — | Blog unavailable blocked by cloudflair | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1852](https://github.com/gitleaks/gitleaks/pull/1852) (@janderssonse) | docs(security): add initial security policy | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1753](https://github.com/gitleaks/gitleaks/issues/1753) | — | Unable to receive license from email for github organization account | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1745](https://github.com/gitleaks/gitleaks/pull/1745) (@jessp01) | Fix typos | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1314](https://github.com/gitleaks/gitleaks/pull/1314) (@riley-kohler) | chore: Fix Contributing Link in README | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1296](https://github.com/gitleaks/gitleaks/pull/1296) (@jdanyow) | zricethezav -> gitleaks in README | **SOURCE-CHECKED — TEST PENDING** | Patch inspected: some URLs have already changed on `master`; compare every link before treating it as redundant. | — |
| [#1080](https://github.com/gitleaks/gitleaks/issues/1080) | — | Mistake in the description for aws-access-token | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1067](https://github.com/gitleaks/gitleaks/issues/1067) | — | Configuration examples gone | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1060](https://github.com/gitleaks/gitleaks/pull/1060) (@it-praktyk) | Add expansion of the SAST abbreviation | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1027](https://github.com/gitleaks/gitleaks/issues/1027) | — | Additional documentation for "extending a central config" | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#958](https://github.com/gitleaks/gitleaks/pull/958) (@YogevBokobza) | better explanation for [extend] feature | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

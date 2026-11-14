# RNC Gitleaks issue/PR category register

Upstream snapshot: **2026-11-08**, **285 open issues + 252 open PRs**. Historical source comparisons reference [`master@5f6ec56`](https://github.com/RNC-Solutions-Group/gitleaks/tree/5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c); current RNC integrated baseline is [`rnc-mod@3a4635f`](https://github.com/RNC-Solutions-Group/gitleaks/commit/3a4635f82dcd0a216f9d2b88ebbd72341b3b60cc). Triage branch: `feat/issues-and-pr`. Source comparisons against historical master are not evidence of a defect in current RNC.

## Routed coverage

| Canonical category | Issues | PRs | Work packages |
| --- | ---: | ---: | ---: |
| [Security and hardening](security-hardening.md) | 4 | 2 | 6 |
| [Scan integrity and file readers](scan-integrity.md) | 15 | 8 | 17 |
| [Detection accuracy](detection-accuracy.md) | 49 | 26 | 51 |
| [Provider rules](detection-rules.md) | 35 | 60 | 61 |
| [Configuration and allowlists](configuration-allowlists.md) | 32 | 18 | 30 |
| [Git and hooks](git-hooks-sources.md) | 12 | 5 | 13 |
| [Reporting and output](reporting-output.md) | 10 | 5 | 11 |
| [Performance and stability](performance-stability.md) | 4 | 4 | 6 |
| [Dependency and build pins](dependencies.md) | 6 | 2 | 8 |
| [Integrations and releases](integrations-release.md) | 4 | 3 | 6 |
| [Optional features](features.md) | 20 | 13 | 22 |
| [Documentation](documentation.md) | 3 | 5 | 4 |
| [RNC dispositions — corrected, superseded, excluded](exclusions.md) | 91 | 101 | 153 |
| **Total assigned from open snapshot** | **285** | **252** | **388** |

**Current RNC queue (post-CP-V12, 2026-11-13):** The original 285 issues and 252 PRs each retain exactly one assigned disposition: **186 active issues / 144 active PRs** in the technical categories and **99 removed issues / 108 removed PRs** in the [resolved, superseded and excluded register](exclusions.md). The 388 original work-package slots consist of **226 active and 162 dispositioned/retired slots**. Unassigned: **0 issues / 0 PRs**. Upstream issue/PR states remain unchanged. The 388 current work-package slots are accounting groups, not 388 distinct bugs or fixes. Related reports and PRs can share a review package while retaining their own acceptance tests and contributor attribution.

The inventories and `_backlog/data/` remain the original dated snapshot. The 2026-11-13 reconciliation moved 31 actioned packages (21 issue IDs, 33 PR IDs), including three additional CP-V5/V7/V10 packages discovered in the final source audit to [dispositions](exclusions.md), with RNC build/CI proof and original contributor attribution cross-referenced in the [validation implementation ledger](https://github.com/RNC-Solutions-Group/gitleaks/blob/rnc-mod/RNC-MOD.md). Strikethrough in an inventory means *routed*, never closed on upstream GitHub. Classification and source-checking do not constitute an implemented fix.

### CP-V12 merged work retired from active queue (2026-11-13)

[RNC PR #5](https://github.com/RNC-Solutions-Group/gitleaks/pull/5) merged as [`3a4635f`](https://github.com/RNC-Solutions-Group/gitleaks/commit/3a4635f82dcd0a216f9d2b88ebbd72341b3b60cc) corrected the V12 Git/reporting/baseline/hook packages. Nine packages (8 issues, 7 PRs) moved from Git/Hooks, Reporting, Performance, Security and Scan Integrity into [RNC dispositions](exclusions.md). Four targets remain active with their precise residuals: #1651/#2324, #1243, #2019 and #1028/#2301. See the [dated CP-V12 reconciliation report](../reconciliations/2026-11-13-v12-reconciliation.md). No upstream issue or PR was closed. Category headline counts were already non-additive before this pass; the global accounting above and package-level exclusion rows are authoritative.

### CP-V11 implemented work retired from active queue (2026-11-13)

[RNC PR #4](https://github.com/RNC-Solutions-Group/gitleaks/pull/4), merged as [`197a7e5`](https://github.com/RNC-Solutions-Group/gitleaks/commit/197a7e5516b0029d0b7bb2554552fc9ed54efb68), preserved all six original upstream contributor commits from PRs #2166, #2294, and #2346. Their three linked issue/PR packages have moved from [Configuration and allowlists](configuration-allowlists.md) into [RNC dispositions](exclusions.md). The older #637 CI-permissions proposal has also been retired from Security because `permissions: contents: read` was already present in V10. Four packages moved total: **3 issues, 4 PRs**, with no upstream issue/PR closure. Historical source comparisons, original inventories and `_backlog/data/` are unchanged. The QOL implementation ledger will be comprehensively refreshed only after V12.

### Scan-integrity second-pass reconciliation (2026-11-09)

One issue (#1464) and two PRs (#2253/#2252) moved from Scan Integrity to the existing Git/Hooks #1981 **same-root-cause** package; no proposal was excluded. PR #1992 moved from Scan Integrity to Optional Features as a new URL/JSON reader. Issue #1651 and PR #2324 now share a package because the PR explicitly describes a partial correction for that issue. Together these reassignments reduce accounting packages from 391 to **388**, without adding or losing an upstream issue or open PR. The prior exclusions ledger remains unchanged.

### Optional Features final-pass reconciliation (2026-11-09)

- **Issue #1397 + open PR #2216** now belong to [Configuration and allowlists](configuration-allowlists.md): signed inline suppression changes finding eligibility, not scanner interface capability.
- **Issue #2017** now belongs to [Integrations and releases](integrations-release.md): the `gore2regex` virtual-address-space failure is tied to **published binary build flags**, not a generic optional feature. Its **closed unmerged PR #2018** is context only and adds no open PR.
- **Issue #1511** now belongs to [RNC exclusions](exclusions.md): source-file rewriting is outside current Gitleaks scan/report semantics; a separate report-only masking request can reopen review.

Optional Features is now **6 issues / 3 open PRs / 6 packages**. Those historical 285 issues, 252 open PRs, and 388 accounting package slots remain complete and unchanged in total. Source/PR inspection was performed, but **no new runtime or regression tests** were run; no implementation or upstream disposition was changed.

### Wave 2 routing update (2026-11-12)

**16 elective work packages (14 issues and 10 open PRs)** were parked in [Optional Features](features.md): Configuration (7), Scan Integrity (4), Reporting (2), Integrations (2), Git/Hooks (1). Optional Features now has **20 issues / 13 PRs / 22 work packages**. All original 285 issues, 252 PRs and 388 work-package accounting slots retain a single owner: **221 active issues + 199 active PRs** in 282 active package slots, and 106 excluded slots. Historical inventories are unchanged.

Three near-term informative QOL candidate packages: [#2130 / PR #2244](features.md#informative-reporting-extensions--high-value-qol-shortlist), [#983 / PR #1654](features.md#informative-reporting-extensions--high-value-qol-shortlist), [#1054](features.md#informative-reporting-extensions--high-value-qol-shortlist). Routing does not authorize code changes. See the [Wave 2 joint reconciliation plan](../04-rnc-reconciliation-plan.md) for post-V8 harness and release priorities.

### Dependency and RE2 source-vs-release reconciliation — 2026-11-13

The current `fix/validation` graph independently matches or supersedes nine dependency PRs and two security upgrade PRs (three linked CVE issues): **12 more accounting packages retired**, without integrating upstream code. Exact-version matches and the independently completed Lip Gloss v2 migration qualify for profile-only contributor thanks; superseded recommendations do not. **The RE2 release address-space defect #2017 is still active**; so are the Go 1.26.3 minimum and Docker login-action proposals. The prior branch's 2026-11-08 baseline columns are historical, not current validation pins.

## Queue housekeeping — 2026-11-13

Only **open or partially corrected work** stays in the active category tables. A corrected original issue and its competing historical PRs move as one package to [dispositions](exclusions.md). An implemented scope still awaiting final owner acceptance is labeled as such instead of being falsely declared released. Distinct residuals retain their own active item/release gate; the original dated snapshots and issue/PR crosswalk are unchanged.

## Triage confidence

- **RNC CORRECTED / IMPLEMENTED / DECLINED / SUPERSEDED:** only with explicit scoped code/test evidence, accepted operator exception, verified replacement or nonissue; [dispositions](exclusions.md) is the removal register, not an upstream GitHub closure. No upstream issue or PR is closed by these documents.
- **REVIEWED — TESTING:** existing CVE PRs #2308/#2336 and their linked issues; not yet adopted.
- **SOURCE-CHECKED — TEST PENDING:** select open PR patches compared directly to fork `master`. **PATCH-REVIEWED** marks upstream diffs inspected without a complete baseline comparison. Neither state implies a build/run result.
- **RNC EXCLUDED:** the [exclusion ledger](exclusions.md) records patch/version evidence, replacement PR or issue when known, and reopen conditions. A PR superseded by another **open** PR is a duplicate review slot, not a proven fix. Nonidentical provider and scanner proposals remain pending.
- **PENDING REVIEW:** remaining records were routed by issue/PR titles and direct crosswalk references. The right category helps evaluate them; it does not establish root cause, patch correctness, priority, or obsolescence.

## Single-owner rule

Every upstream issue and PR has exactly **one primary category**. Verified related issues and the PRs that address them form **one work package** in the same category. Category documents may mention other IDs for context, but that does not transfer ownership. We do **not** turn weak mentions or third-party release-note references into confirmed relationships. A shared numeric reference is never sufficient to prove duplicate root causes.

## Branch and attribution workflow

1. Keep `master` distinct from the unmerged `fix/docker-security-bases` hardened image branch. A newer Go image does not update source dependencies, tests, or CLI behavior.
2. On review, record root cause, applicable PR and author, original commit SHA, code comparison, reproducer, regression tests, and functional tradeoffs. Prefer the best implementation, not merely the first or smallest PR.
3. **Adopted:** only after upstream code is actually integrated; record upstream PR/commit, RNC integration commit, RNC PR, branch, and verification. **Corrected:** reserved for an RNC-origin fix, with RNC commit/PR and tests.
4. **Superseded:** record the known replacement PR/commit or concrete current code; an open replacement does not close the underlying defect. **Closed:** presently a nonissue or RNC scope exception with precise reason. **Obsolete:** a validated outdated/inapplicable approach, with evidence. Do not use any of these merely to reduce counts.
5. Keep only one canonical category for a work package; move all tied IDs together. Once categorized, strike each source link in the three Markdown inventories. Leave `_backlog/data/` (JSON and CSV) immutable.
6. Only the user promotes or merges code between security PR, shared-security branch, experiment branch, and fork `master`.

## Review evidence checklist

| Evidence | Required before implementation/disposition |
| --- | --- |
| Baseline | Exact master SHA, module pins, container image digest if applicable |
| Cause | Reproducer and actual source behavior |
| Proposal | Upstream author, PR and head commit, competing PRs |
| Safety | Negative tests, compatibility, regression and performance where relevant |
| Integration | Which RNC PR/commit contains which upstream commit, or RNC-only correction |
| Outcome | Adopted / corrected / superseded / closed / obsolete / deferred, with rationale and reopen trigger |

See [issue backlog](../01-issue-backlog.md), [open PRs](../02-open-pull-requests.md), and [crosswalk](../03-pr-issue-crosswalk.md).

# RNC Gitleaks issue/PR category register

Upstream snapshot: **2026-11-08**, **285 open issues + 252 open PRs**. Current fork comparison target: [`RNC-Solutions-Group/gitleaks@master`](https://github.com/RNC-Solutions-Group/gitleaks/tree/5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c); triage work branch: `feat/issues-and-pr`. This is a review/ownership index, **not** a claim that each fix was verified, merged, or released.

## Routed coverage

| Canonical category | Issues | PRs | Work packages |
| --- | ---: | ---: | ---: |
| [Security and hardening](security-hardening.md) | 10 | 7 | 13 |
| [Scan integrity and file readers](scan-integrity.md) | 24 | 19 | 30 |
| [Detection accuracy](detection-accuracy.md) | 49 | 26 | 51 |
| [Provider rules](detection-rules.md) | 35 | 60 | 61 |
| [Configuration and allowlists](configuration-allowlists.md) | 43 | 29 | 42 |
| [Git and hooks](git-hooks-sources.md) | 18 | 14 | 19 |
| [Reporting and output](reporting-output.md) | 13 | 7 | 14 |
| [Performance and stability](performance-stability.md) | 4 | 4 | 6 |
| [Dependency and build pins](dependencies.md) | 7 | 11 | 18 |
| [Integrations and releases](integrations-release.md) | 4 | 6 | 9 |
| [Optional features](features.md) | 6 | 3 | 6 |
| [Documentation](documentation.md) | 8 | 13 | 13 |
| [RNC exclusions — obsolete, superseded, exceptions](exclusions.md) | 64 | 53 | 106 |
| **Total assigned from open snapshot** | **285** | **252** | **388** |

**Coverage:** All 285 issues and 252 PRs from the snapshot have exactly one RNC disposition: **221 active issues / 199 active PRs** in the technical categories, and **64 excluded issues / 53 excluded PRs** documented in [exclusions](exclusions.md). Unassigned: **0 issues / 0 PRs**. The 388 current work-package slots are accounting groups, not 388 distinct bugs or fixes. Related reports and PRs can share a review package while retaining their own acceptance tests and contributor attribution.

The inventories and `_backlog/data/` remain the original dated snapshot. Strikethrough in an inventory means *routed*, never closed on upstream GitHub. Classification and source-checking do not constitute an implemented fix.

### Scan-integrity second-pass reconciliation (2026-11-09)

One issue (#1464) and two PRs (#2253/#2252) moved from Scan Integrity to the existing Git/Hooks #1981 **same-root-cause** package; no proposal was excluded. PR #1992 moved from Scan Integrity to Optional Features as a new URL/JSON reader. Issue #1651 and PR #2324 now share a package because the PR explicitly describes a partial correction for that issue. Together these reassignments reduce accounting packages from 391 to **388**, without adding or losing an upstream issue or open PR. The prior exclusions ledger remains unchanged.

### Optional Features final-pass reconciliation (2026-11-09)

- **Issue #1397 + open PR #2216** now belong to [Configuration and allowlists](configuration-allowlists.md): signed inline suppression changes finding eligibility, not scanner interface capability.
- **Issue #2017** now belongs to [Integrations and releases](integrations-release.md): the `gore2regex` virtual-address-space failure is tied to **published binary build flags**, not a generic optional feature. Its **closed unmerged PR #2018** is context only and adds no open PR.
- **Issue #1511** now belongs to [RNC exclusions](exclusions.md): source-file rewriting is outside current Gitleaks scan/report semantics; a separate report-only masking request can reopen review.

Optional Features is now **6 issues / 3 open PRs / 6 packages**. The 285 issues, 252 open PRs, and 388 accounting packages remain complete and unchanged in total. Source/PR inspection was performed, but **no new runtime or regression tests** were run; no implementation or upstream disposition was changed.

## Triage confidence

- **RNC CLOSED / SUPERSEDED:** only when a known replacement or a verified nonissue has supporting evidence, with [exclusions.md](exclusions.md) as the exclusion ledger. No upstream issue or PR is closed by these documents.
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

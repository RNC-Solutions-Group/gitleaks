# RNC Gitleaks issue/PR category register

Upstream snapshot: **2026-11-08**, **285 open issues + 252 open PRs**. Current fork comparison target: [`RNC-Solutions-Group/gitleaks@master`](https://github.com/RNC-Solutions-Group/gitleaks/tree/5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c); triage work branch: `feat/issues-and-pr`. This is a review/ownership index, **not** a claim that each fix was verified, merged, or released.

## Routed coverage

| Canonical category | Issues | PRs | Work packages |
| --- | ---: | ---: | ---: |
| [Security and hardening](security-hardening.md) | 12 | 6 | 14 |
| [Scan integrity and file readers](scan-integrity.md) | 23 | 21 | 32 |
| [Detection accuracy](detection-accuracy.md) | 56 | 37 | 79 |
| [Provider rules](detection-rules.md) | 33 | 52 | 56 |
| [Configuration and allowlists](configuration-allowlists.md) | 47 | 31 | 61 |
| [Git and hooks](git-hooks-sources.md) | 19 | 13 | 22 |
| [Reporting and output](reporting-output.md) | 15 | 8 | 17 |
| [Performance and stability](performance-stability.md) | 4 | 3 | 5 |
| [Dependency and build pins](dependencies.md) | 8 | 12 | 20 |
| [Integrations and releases](integrations-release.md) | 3 | 6 | 8 |
| [Optional features](features.md) | 9 | 3 | 8 |
| [Documentation](documentation.md) | 6 | 12 | 11 |
| [RNC exclusions — obsolete, superseded, exceptions](exclusions.md) | 50 | 48 | 87 |
| **Total assigned from open snapshot** | **285** | **252** | **420** |

**Coverage:** All 285 issues and 252 PRs from the snapshot have exactly one RNC disposition: **235 active issues / 204 active PRs** in the technical categories, and **50 excluded issues / 48 excluded PRs** documented in [exclusions](exclusions.md). Unassigned: **0 issues / 0 PRs**. The 420 work-package slots are accounting groups, not 420 distinct bugs or fixes.

The inventories and `_backlog/data/` remain the original dated snapshot. Strikethrough in an inventory means *routed*, never closed on upstream GitHub. Classification and source-checking do not constitute an implemented fix.

## Triage confidence

- **RNC CLOSED / SUPERSEDED:** only when a known replacement or a verified nonissue has supporting evidence, with [exclusions.md](exclusions.md) as the exclusion ledger. No upstream issue or PR is closed by these documents.
- **REVIEWED — TESTING:** existing CVE PRs #2308/#2336 and their linked issues; not yet adopted.
- **SOURCE-CHECKED — TEST PENDING:** select open PR patches compared directly to fork `master`. No build/run proof is implied.
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

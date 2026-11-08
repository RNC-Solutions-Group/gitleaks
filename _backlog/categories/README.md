# RNC Gitleaks issue/PR category register

Upstream snapshot: **2026-11-08**, **285 open issues + 252 open PRs**. Current fork comparison target: [`RNC-Solutions-Group/gitleaks@master`](https://github.com/RNC-Solutions-Group/gitleaks/tree/5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c); triage work branch: `feat/issues-and-pr`. This is a review/ownership index, **not** a claim that each fix was verified, merged, or released.

## Routed coverage

| Canonical category | Issues | PRs | Work packages |
| --- | ---: | ---: | ---: |
| [Security and hardening](security-hardening.md) | 7 | 5 | 9 |
| [Scan integrity and file readers](scan-integrity.md) | 34 | 17 | 46 |
| [Detection accuracy](detection-accuracy.md) | 46 | 28 | 66 |
| [Provider rules](detection-rules.md) | 39 | 45 | 72 |
| [Configuration and allowlists](configuration-allowlists.md) | 34 | 28 | 49 |
| [Git and hooks](git-hooks-sources.md) | 21 | 15 | 30 |
| [Reporting and output](reporting-output.md) | 25 | 10 | 28 |
| [Performance and stability](performance-stability.md) | 4 | 2 | 4 |
| [Dependency and build pins](dependencies.md) | 11 | 13 | 24 |
| [Integrations and releases](integrations-release.md) | 14 | 7 | 19 |
| [Optional features](features.md) | 6 | 2 | 7 |
| [Candidate duplicates / competing fixes](crosswalk-duplicates.md) | 32 | 57 | 36 |
| [Documentation](documentation.md) | 7 | 17 | 22 |
| [Unclassified review](needs-review.md) | 4 | 0 | 4 |
| [Closed and superseded (RNC)](closed.md) | 1 | 6 | 7 |
| **Total assigned from open snapshot** | **285** | **252** | **423** |

**Unassigned from snapshot: 0 issues, 0 PRs.** This is **100% routing coverage**, not 100% technical review or remediation. A routed candidate remains active unless documented otherwise. The original three inventories and raw API exports keep their historical totals. A struck-through link in the three Markdown inventories means *categorized*, never *closed on GitHub*.

Two duplicate-category issue targets may be referenced outside the current open-issue snapshot; those historical references are not included in the 285 count. The original duplicate report's 174 candidate-link tally is unchanged. Counts in this table use **unique open-snapshot IDs**, so historic counts and a document's incidental hyperlinks must not be summed.

## Triage confidence

- **RNC CLOSED / SUPERSEDED:** only when a known replacement or a verified nonissue has supporting evidence, with [closed.md](closed.md) as the exclusion ledger. No upstream issue or PR is closed by these documents.
- **REVIEWED — TESTING:** existing CVE PRs #2308/#2336 and their linked issues; not yet adopted.
- **SOURCE-CHECKED — TEST PENDING:** select open PR patches compared directly to fork `master` (including #2329, #2317, #2319, #2294, #2237, #2206, #2307, #2348 and #2258). No build/run proof is implied.
- **PENDING REVIEW:** remaining records were routed by issue/PR titles and direct crosswalk references. The right category helps evaluate them; it does not establish root cause, patch correctness, priority, or obsolescence.

## Single-owner rule

Every upstream issue and PR has exactly **one primary category**. Verified related issues and the PRs that address them form **one work package** in the same category. Category documents may mention other IDs for context, but that does not transfer ownership. We do **not** turn weak mentions or third-party release-note references into confirmed relationships. For example, issue #2304 merely mentioning PR #2300 does not make a report-deletion bug the same fix as a report-permission issue.

## Branch and attribution workflow

1. Keep `master` distinct from the unmerged `fix/docker-security-bases` hardened image branch. A newer Go image does not update source dependencies, tests, or CLI behavior.
2. On review, record root cause, applicable PR and author, original commit SHA, code comparison, reproducer, regression tests, and functional tradeoffs. Prefer the best implementation, not merely the first or smallest PR.
3. **Adopted:** only after upstream code is actually integrated; record upstream PR/commit, RNC integration commit, RNC PR, branch, and verification. **Corrected:** reserved for an RNC-origin fix, with RNC commit/PR and tests.
4. **Superseded:** record the known replacement PR/commit or concrete current code. **Closed:** presently a nonissue but provenance unknown. **Obsolete:** a validated outdated/inapplicable approach, with evidence. Do not use any of these merely to reduce counts.
5. Keep only one canonical category for a work package; move all tied IDs together. Once categorized, strike each source link in the three Markdown inventories. Leave `_backlog/data/` (JSON and CSV) immutable.
6. Only the user promotes or merges code between security PR, shared-security branch, experiment branch, and fork `main`.

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

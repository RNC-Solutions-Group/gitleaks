# Gitleaks backlog — RNC assessment workspace

Collection timestamp: 2026-11-08T07:14:05.557Z

This branch contains a dated, reproducible review snapshot of **upstream `gitleaks/gitleaks`**, not changes to the scanner. Active categories identify what remains to assess; separate records retain exclusion decisions and duplicate-history evidence. Linked false-positive/false-negative reports and competing patches share a review package only when their overlap is established, with separate fixtures retained where behavior differs.

## Contents

- [01 — Open issue backlog](01-issue-backlog.md): all 285 open issues returned by the date-partitioned searches.
- [02 — Open pull requests](02-open-pull-requests.md): all 252 open PRs, with authors, refs, and review state.
- [03 — PR/issue crosswalk](03-pr-issue-crosswalk.md): links, confidence, and exclusions.
- [data/](data/): structured public-safe JSON snapshots, filterable CSV, and collection manifest.
- [categories/](categories/README.md): current technical work queues, ownership and decision gates.

## Current RNC review status

| Disposition | Issues | PRs |
| --- | ---: | ---: |
| Active technical-category review | **235** | **203** |
| [RNC exclusions](categories/exclusions.md) | **50** | **49** |
| **Original snapshot assigned** | **285** | **252** |

All upstream snapshot items have one disposition in the [active category index](categories/README.md), with **0 unassigned**. The register records 404 work-package accounting slots, which are not necessarily unique root causes. Source review is not implementation, regression testing, merging, or an upstream closure.

## Historical evidence (outside active categories)

- [Duplicate / competing-fix crosswalk](crosswalk-duplicates.md) — the 174 original candidate relationships, subsequent patch comparisons, retained contributors and disposition reasoning. This is a **historical audit artifact**, not an active work queue.
- [PR/issue crosswalk inventory](03-pr-issue-crosswalk.md) — original relationship extraction and confidence.
- [Exclusions ledger](categories/exclusions.md) — canonical reasons, evidence and reopen gates for retired RNC review slots.

The source inventories and JSON/CSV exports are retained unchanged. Category files focus on open work and the tests needed to evaluate it; retirement history stays in the historical crosswalk or exclusion ledger.

## Review lifecycle

1. **Unreviewed:** collected, automatically linked/classified, no technical judgment yet.
2. **Candidate:** review the issue, PR, diff, upstream discussion, threat model, and compatibility.
3. **Tested:** reproduce on a pinned baseline, cherry-pick or implement on an isolated experiment branch, execute focused and full regressions.
4. **Disposition:** adopt, defer, reject, superseded, duplicate, blocked, or failed test. Record **why**, including failed commands and expected versus observed behavior.
5. **Promotion:** only the owner authorizes promotion into security-only, combined-security, or general-development branches. A textual relationship is not a passing test.

This branch is a backlog workspace, not a patch-integration branch. No PRs were cherry-picked or merged and no upstream issues were changed.

## Data safety and interpretation

The GitHub API issue and PR descriptions were used in memory to derive references. **Body text is excluded from committed exports** because issue descriptions may contain token-shaped values, test credentials, private logs, or other unsafe content. `data/issues.json` and `data/open-prs.json` are structured metadata exports, not byte-for-byte archival dumps. Each record retains the upstream URL for direct inspection. Full code diffs, review threads, and linked-issue GraphQL nodes have **not** been exported.

GitHub's `linked:issue` search indicator, explicit closing keywords, issue URLs, and ordinary `#number` mentions are recorded separately. A mere mention is not proof a PR fixes an issue. The category suggestions are keyword-based and require human review.

Run a fresh collection before making claims about the *current* upstream backlog. Source state changes continuously.

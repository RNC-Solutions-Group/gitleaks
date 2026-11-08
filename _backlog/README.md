# Gitleaks issue and PR assessment workspace

Collection timestamp: 2026-11-08T07:14:05.557Z

This branch contains an initial, reproducible review snapshot of **upstream `gitleaks/gitleaks`**, not a change to the Gitleaks scanner. It is a workspace for deciding which community fixes and enhancements merit independent tests in RNC's downstream branches.

## Contents

- [01 — Open issue backlog](01-issue-backlog.md): all 285 open issues returned by the date-partitioned searches.
- [02 — Open pull requests](02-open-pull-requests.md): all 252 open PRs, with authors, refs, and review state.
- [03 — PR/issue crosswalk](03-pr-issue-crosswalk.md): links, confidence, and exclusions.
- [data/](data/): structured public-safe JSON snapshots, filterable CSV, and collection manifest.
- [categories/](categories/README.md): where reviewed candidates and experiment outcomes will be documented.

## RNC second-pass junk-drawer assessment (2026-11-08)

The [category ownership register](categories/README.md) now reflects code/diff and targeted issue-body elimination passes against fork `master` at `5f6ec56c`. **25 issues and 31 PRs** have defensible RNC exclusions (obsolete, superseded, or closed by exception); **260 issues and 221 PRs** remain active for independent review. Details, original contributor links, rationale, and reopen criteria are in [exclusions.md](categories/exclusions.md). Additional cross-category moves were made from actual changed files; competing nonidentical fixes remain pending.

The focused categories are [documentation](categories/documentation.md) (including an **attribution-preserving, separate-branch seven-PR consolidation recommendation**), [dependencies](categories/dependencies.md) (open issue triage and a repaired PR table), [integrations/release](categories/integrations-release.md), and [needs review](categories/needs-review.md). No proposed documentation patches have been adopted here.

These are **RNC review dispositions only**. They do not close upstream GitHub items, endorse or adopt a surviving PR, or prove that any unresolved defect is fixed. The three full Markdown inventories still represent their original upstream snapshot, with strikethrough meaning *routed*; the JSON/CSV data exports remain unchanged.

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

# RNC Gitleaks — issue and PR review workspace

> [!IMPORTANT]
> **This is the `feat/issues-and-pr` triage branch, not the published Gitleaks project README, a feature release, or an implementation branch.** Its root README is intentionally branch-specific. The normal Gitleaks README remains on RNC [`master`](https://github.com/RNC-Solutions-Group/gitleaks/blob/master/README.md) and [upstream](https://github.com/gitleaks/gitleaks/blob/master/README.md).

## What is happening here?

RNC Solutions Group forked Gitleaks to maintain its own security-scanning workflow and evaluate existing community contributions. This branch is the **review desk** for the upstream GitHub issue and pull-request backlog: a place to identify duplicate work, compare proposed changes with the fork's existing code, and decide what deserves a separate technical assessment.

We are cleaning out the backlog before committing engineering time. That means identifying:

- **Obsolete:** the proposed code, API, platform, or execution path no longer applies.
- **Superseded:** another identified PR, issue, or implemented code path already covers the same work. An open replacement PR remains **pending**; it is not an adopted fix.
- **RNC exclusion:** the report is outside the fork's scope, is not an actionable scanner defect, or does not justify independent RNC review.
- **Still under review:** anything requiring a reproducer, more evidence, tradeoff analysis, tests, or an adoption decision.

An RNC exclusion **does not close an upstream item**. It removes an item from *our independent review queue* with a documented reason and reopening condition.

## Where is the backlog?

Start with [`_backlog/README.md`](_backlog/README.md), then use:

| File | Purpose |
| --- | --- |
| [`categories/README.md`](_backlog/categories/README.md) | Category ownership, counts, confidence, and workflow rules |
| [`categories/exclusions.md`](_backlog/categories/exclusions.md) | Evidence-backed obsolete, superseded, and exception dispositions; contributor references; reopen triggers |
| [`01-issue-backlog.md`](_backlog/01-issue-backlog.md) | Full upstream issue inventory |
| [`02-open-pull-requests.md`](_backlog/02-open-pull-requests.md) | Full upstream PR inventory with authors and links |
| [`03-pr-issue-crosswalk.md`](_backlog/03-pr-issue-crosswalk.md) | Issue/PR relationships and confidence |
| [`data/`](_backlog/data/) | Immutable snapshot exports, not a live issue tracker |

The initial snapshot on **2026-11-08** contained **285 open upstream issues and 252 upstream PRs**. Dispositions are evaluated against fork `master` at [`5f6ec56c`](https://github.com/RNC-Solutions-Group/gitleaks/commit/5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c), with the separately hardened Docker branch tracked as a different implementation baseline. These counts and comparisons are a dated snapshot, not a claim about today's upstream state.

## Branch boundaries and contributor credit

1. **No upstream PR is adopted merely because it was reviewed here.** This branch contains reports and dispositions, not cherry-picked scanner changes.
2. Every retained PR keeps its upstream link and original contributor. A future consolidation must identify the source PR(s), source commits, RNC integration commit, and any adaptation made.
3. A contributor's original commit can be preserved in an authorized implementation branch, but any integration, correction, or merge requires independent review and owner approval.
4. **Do not merge this triage README as the product README.** If the review branch is ever merged, the root README is intentionally disposable; keep the normal upstream-derived README from the destination branch.
5. Original issue and PR exports remain unmodified. Strikethrough in the Markdown inventories indicates *routed to a category*, not closed upstream.

## Upstream project attribution

**Gitleaks is an upstream open-source project, not an RNC-origin scanner.** The original software and documentation belong to the [Gitleaks project](https://github.com/gitleaks/gitleaks), founded by [@zricethezav](https://github.com/zricethezav), and its [contributors](https://github.com/gitleaks/gitleaks/graphs/contributors). This fork retains the upstream [MIT license](LICENSE), notices, and Git history. RNC's branch-specific review documents are independent analysis; they do not transfer authorship or represent upstream maintainer decisions.

Original Gitleaks identity and project mark:

```
┌─○───┐
│ │╲  │
│ │ ○ │
│ ○ ░ │
└─░───┘
```

Gitleaks detects secrets, including passwords, API keys, and tokens, in Git repositories, ordinary files, and standard input. That scanner and its original user documentation remain the upstream project's work.

The upstream maintainer states that Gitleaks is feature complete, with future Gitleaks releases limited to security patches, and directs feature development toward [Betterleaks](https://github.com/betterleaks/betterleaks). For normal usage, installation, CLI flags, and configuration, consult the [original upstream README](https://github.com/gitleaks/gitleaks/blob/master/README.md) or [RNC master README](https://github.com/RNC-Solutions-Group/gitleaks/blob/master/README.md).

The original README also credits its sponsor, [CodeRabbit](https://coderabbit.ai/?utm_source=oss&utm_medium=sponsorship&utm_campaign=gitleaks), and the [Gitleaks community Discord](https://discord.gg/8Hzbrnkr7E). Those credits are preserved here; this review branch makes no separate sponsorship claim.

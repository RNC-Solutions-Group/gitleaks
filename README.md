# RNC Gitleaks — assessment and issue/PR review workspace

> [!IMPORTANT]
> **This is the `feat/issues-and-pr` assessment branch, not the published Gitleaks project README, a feature release, or an implementation branch.** Its root README is intentionally branch-specific. The normal Gitleaks README remains on RNC [`master`](https://github.com/RNC-Solutions-Group/gitleaks/blob/master/README.md) and [upstream](https://github.com/gitleaks/gitleaks/blob/master/README.md).

## What is happening here?

RNC Solutions Group forked Gitleaks to maintain its own security-scanning workflow and evaluate existing community contributions. This branch is an **evidence and assessment workspace** with three related purposes:

1. Curate the upstream GitHub issue/PR backlog and preserve the original contributor references and disposition reasoning.
2. Establish an executable, source-grounded **current-state technical baseline** for the fork: architecture, capabilities, correctness, regression coverage, performance, stability, security, and build/integration behavior.
3. Later, compare external PR implementations against that verified baseline to determine applicability, compatibility, benefits, regressions, and unresolved questions. A mismatch is not proof that the external implementation is wrong.

The original upstream inventory and assessment remain preserved, while active development has reached [RNC CP-V10](https://github.com/RNC-Solutions-Group/gitleaks/commit/bd42a958f7476e9a6d8e17bd2dbb2abe3ff060ea) on **`fix/validation`** ([12/12 cross-platform CI](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43321594496)). This branch continues to own *documentation-only* issue/PR dispositions and their research history. Queue dispositions now distinguish:

- **Obsolete:** the proposed code, API, platform, or execution path no longer applies.
- **Superseded:** another identified PR, issue, or implemented code path already covers the same work. An open replacement PR remains **pending**; it is not an adopted fix.
- **Corrected / implemented:** RNC code or documentation demonstrably addresses the scoped package; final operator or harness acceptance is recorded separately when still pending.
- **Declined RNC:** the operator explicitly rejected the proposed integration or behavior.
- **RNC scope exclusion:** the report does not justify independent work, is unrelated to this fork, or is an intentionally excluded proposal.
- **Still under review:** anything requiring a reproducer, more evidence, tradeoff analysis, tests, or an adoption decision.

An entry in the RNC disposition register **does not close an upstream item**. It removes it from *our active review queue* with scoped source/CI evidence (where corrected), original contributor credit, and an explicit reopening condition. Distinct unresolved residuals retain their own active issue or named release gate.

## Where is the backlog?

Start with [`_backlog/README.md`](_backlog/README.md), then use:

| File | Purpose |
| --- | --- |
| [`categories/README.md`](_backlog/categories/README.md) | Category ownership, counts, confidence, and workflow rules |
| [`categories/exclusions.md`](_backlog/categories/exclusions.md) | Canonical **corrected, implemented, declined, superseded and excluded** dispositions, contributor credit, scope and reopen triggers |
| [`01-issue-backlog.md`](_backlog/01-issue-backlog.md) | Full upstream issue inventory |
| [`02-open-pull-requests.md`](_backlog/02-open-pull-requests.md) | Full upstream PR inventory with authors and links |
| [`03-pr-issue-crosswalk.md`](_backlog/03-pr-issue-crosswalk.md) | Issue/PR relationships and confidence |
| [`crosswalk-duplicates.md`](_backlog/crosswalk-duplicates.md) | Historical duplicate/competing-fix assessment, contributor links, and decisions; **not an active category** |
| [`data/`](_backlog/data/) | Immutable snapshot exports, not a live issue tracker |

The initial snapshot on **2026-11-08** contained **285 open upstream issues and 252 upstream PRs**. Dispositions are evaluated against fork `master` at [`5f6ec56c`](https://github.com/RNC-Solutions-Group/gitleaks/commit/5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c), with the separately hardened Docker branch tracked as a different implementation baseline. These counts and comparisons are a dated snapshot, not a claim about today's upstream state.

**Current reconciled queue (2026-11-13):** **200 active issue references / 166 active PR references / 251 active package slots**. The [resolved/superseded/excluded disposition register](_backlog/categories/exclusions.md) contains **85 issue references / 86 PR references / 137 retired slots**, including 31 newly moved completed or declined packages. Original counts remain **285 issues / 252 PRs / 388 package slots**, with zero unassigned; `_backlog/data/` is unchanged. See the [canonical register](_backlog/README.md), [category index](_backlog/categories/README.md), and [implementation and attribution ledger on `fix/validation`](https://github.com/RNC-Solutions-Group/gitleaks/blob/fix/validation/fixes-qol-tracking.md).

## Historical current-state technical assessment

The [assessment workspace](_assessment/README.md) preserves the original directory map, evidence requirements and CP0–CP5 baseline review framework, which **predates** the subsequent CP-V1–CP-V10 scanner implementation. Those assignments are historical assessment references, not the current development roadmap. For present source/checkpoint decisions use the [Wave 2 roadmap](_backlog/04-rnc-reconciliation-plan.md) and [validation engineering ledger](https://github.com/RNC-Solutions-Group/gitleaks/blob/fix/validation/fixes-qol-tracking.md).

- **CP0:** repository, Git lineage, environment, dependencies, test infrastructure, and safe command inventory.
- **CP1:** architecture, interfaces, execution paths, and component relationships.
- **CP2:** functional capabilities, correctness, failure modes, and regression-test coverage.
- **CP3:** reproducible performance, resource use, and stability measurements.
- **CP4:** security, configuration, integrations, developer tooling, Makefile side effects, and build/release reproducibility.
- **CP5:** reconcile findings into a dated baseline and a prioritized investigation roadmap.

The current-state checkpoints are stored under `_assessment/current-state/CPX-*/`. Later cross-repository PR comparisons will be documented separately under `_assessment/` after the baseline is accepted. External PRs and issues in `_backlog/` are leads to investigate, not evidence that a defect exists in this fork.

**Execution boundary:** source, tests, Makefiles, dependencies, configuration, and CI are read-only on this branch. Authorized agents may inspect the code, run safe Go diagnostics, and exercise the sibling `../gitleaks-test-repo` read-only when available. `_assessment/` is the agent's documentation workspace; `_backlog/` is managed by Dave and GPT and is read-only to other agents absent explicit authorization. The [root agent instructions](AGENTS.md) and [assessment instructions](_assessment/AGENTS.md) govern these boundaries. Existing Makefile targets with implicit formatting, dependency changes, regeneration, cleanup, or network fetches are not authorized diagnostic entry points.

A baseline report must pin the **assessed code SHA** separately from subsequent documentation commits, preserve failed/unavailable test outcomes, and distinguish verified facts, measured results, inferred risks, and unverified hypotheses. The assessment is meant to characterize the application, not fix every discovered issue or declare a PR suitable for adoption.

## Branch boundaries and contributor credit

1. **No scanner source changes or upstream PR integrations are authorized on this branch.** Source is examined and executed for assessment, not modified; PR review is not adoption.
2. Only `_backlog/` (owner/GPT-managed), `_assessment/` (authorized assessment work), and **this branch-specific README when explicitly requested by the owner** are documentation write areas. General agents may write only within their specifically assigned scope.
3. Every retained PR keeps its upstream link and original contributor. A future consolidation must identify the source PR(s), source commits, RNC integration commit, and any adaptation made.
4. A contributor's original commit can be preserved in an authorized implementation branch, but any integration, correction, or merge requires independent review and owner approval.
5. **Do not merge this branch-specific README as the product README.** If the assessment branch is ever merged, retain the normal upstream-derived README from the destination branch.
6. Original issue and PR exports remain unmodified. Strikethrough in the Markdown inventories indicates *routed to a category*, not closed upstream.

## AI-assisted assessment and owner-directed development

RNC uses AI-assisted and AI-managed workflows to organize backlog evidence, examine source code, run authorized diagnostics, compare upstream proposals, and prepare technical documentation. ChatGPT, OpenCode, DeepSeek-based agents, and other tools may contribute within explicitly assigned scopes. **Dave defines the objectives, controls acceptance and integration decisions, and authorizes any source modifications on separate implementation branches.**

AI-produced findings and recommendations are **working analysis, not automatic proof**. Claims must be checked against pinned source code, reproducible tests, measurements, and original PR discussions; uncertainty and failed verification must remain visible. Agent instructions and checkpoint assignments constrain work on this branch. No use of AI here implies that upstream code, external PRs, or original Gitleaks contributions were authored by RNC or endorsed by upstream maintainers.


**Current branch/source split:** `feat/issues-and-pr` = historical inventories, active category queues, resolved dispositions and assessment documents; `fix/validation` = implemented CP-V8–V10 scanner and [fixes/QOL ledger](https://github.com/RNC-Solutions-Group/gitleaks/blob/fix/validation/fixes-qol-tracking.md); independent harness = separate testing repository. The former branch's source tree is not the new scanner implementation. Do not merge this review README over the actual product README.

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

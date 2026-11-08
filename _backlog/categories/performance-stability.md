# Performance Stability — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**4 upstream issues · 2 upstream PRs · 4 work packages.** Each ID is owned by this category once; linked PRs and issues share a row. Existing upstream state is unchanged.

Memory, hangs, race-prone parsers, decoding cost, and crashes. First establish a measurable reproduction and a bounded test workload.

**Evidence rules:** “SOURCE-CHECKED” means a targeted patch/current-source comparison, **not** that a fix was built or tested. Most records are title/crosswalk-only and must stay pending until reproduced, diffed, and assessed on the fork. “Obsolete,” “superseded,” “adopted,” or “corrected” require separate evidence; no PR is adopted by routing it here.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| [#1338](https://github.com/gitleaks/gitleaks/issues/1338) | [#2246](https://github.com/gitleaks/gitleaks/pull/2246) (@aron-intframe) | gitdiff report "parser" can cause an infinite wait on parser error; proposal: fix(sources): stop a failed diff parse from hanging or silently truncating a scan | **SOURCE-CHECKED — TEST PENDING** | Patch inspected: `master/sources/git.go` still uses the unguarded diff-parser flow; needs hang tests. | — |
| [#2019](https://github.com/gitleaks/gitleaks/issues/2019) | [#2212](https://github.com/gitleaks/gitleaks/pull/2212) (@RajeshRajendiran) | Memory issue in GitLab when scanning large repos with Gitleaks (--max-decode-depth behavior); proposal: docs: fix stale --max-decode-depth default in README | **LINKED — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1437](https://github.com/gitleaks/gitleaks/issues/1437) | — | hang on futex | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#752](https://github.com/gitleaks/gitleaks/issues/752) | — | invalid memory address or nil pointer dereference | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

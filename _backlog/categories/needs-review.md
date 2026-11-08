# Needs Review — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**4 upstream issues · 0 upstream PRs · 4 work packages.** Each ID is owned by this category once; linked PRs and issues share a row. Existing upstream state is unchanged.

Ambiguous, insufficiently evidenced, or out-of-domain reports. These are parked for clarification; none is presumed bogus or resolved.

**Evidence rules:** “SOURCE-CHECKED” means a targeted patch/current-source comparison, **not** that a fix was built or tested. Most records are title/crosswalk-only and must stay pending until reproduced, diffed, and assessed on the fork. “Obsolete,” “superseded,” “adopted,” or “corrected” require separate evidence; no PR is adopted by routing it here.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| [#2322](https://github.com/gitleaks/gitleaks/issues/2322) | — | ghghj | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1864](https://github.com/gitleaks/gitleaks/issues/1864) | — | Sample AKIA credentials cause SAST/Security Scanners to false flag | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1624](https://github.com/gitleaks/gitleaks/issues/1624) | — | Organisation key email never received | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1127](https://github.com/gitleaks/gitleaks/issues/1127) | — | Gitleaks 8.14.1 Strange bug!? | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

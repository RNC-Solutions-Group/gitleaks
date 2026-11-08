# Needs Review — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**2 active upstream issues · 0 upstream PRs · 2 work packages.** Issue #2322 (empty template) and the external licensing email complaint #1624 have been excluded in [exclusions.md](exclusions.md); no upstream state was changed.

Ambiguous, insufficiently evidenced, or out-of-domain reports. These are parked for clarification; none is presumed bogus or resolved.

**Evidence rules:** “SOURCE-CHECKED” means a targeted patch/current-source comparison, **not** that a fix was built or tested. Most records are title/crosswalk-only and must stay pending until reproduced, diffed, and assessed on the fork. “Obsolete,” “superseded,” “adopted,” or “corrected” require separate evidence; no PR is adopted by routing it here.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| [#1864](https://github.com/gitleaks/gitleaks/issues/1864) | — | Vendored AWS SDK sample AKIA credentials repeatedly trigger configured SAST/secret-scanning exclusions. | **SOURCE-READ — SCOPE UNCLEAR, RETAINED** | Reporter links [AWS SDK upstream #3087](https://github.com/aws/aws-sdk-go-v2/issues/3087) but gives no Gitleaks version, exact finding/rule, or controlled reproduction. Might be a third-party scanner or legitimate fixture detection. Keep parked; only move to detection accuracy if proven to involve this scanner's matching rules. | — |
| [#1127](https://github.com/gitleaks/gitleaks/issues/1127) | — | Unexplained PowerShell/Azure DevOps error after an apparently successful Gitleaks 8.14.1 scan. | **SOURCE-READ — HISTORICAL REPRO, RETAINED** | Script uses old `detect --no-git` command and `Start-Process`; exact error is embedded in screenshots, not transcribed. Current `dir` mode differs, but a newer command alone does not explain whether the observed output was scanner or CI-host generated. Keep until the screenshot/error text and a current reproduction are assessed. | — |

**Focused disposition:** no exclusion is justified for these two from the issue bodies alone. Both stay parked rather than inferring that old version numbers automatically resolve the reported symptoms.

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

# Documentation — RNC review queue

**Current active queue (2026-11-13):** Actioned issues/PRs have moved to the [RNC disposition register](exclusions.md); implementation, test evidence and contributor attribution are maintained in [fixes/QOL](https://github.com/RNC-Solutions-Group/gitleaks/blob/fix/validation/fixes-qol-tracking.md). Do not reassign finished work from historical review notes.

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**8 issues · 13 PRs · 13 work packages.**

Documentation, setup guidance, labels, and examples. Compare proposed text with current code and README behavior before consolidating patches.

**Review states:** `SOURCE-CHECKED` records a targeted source/patch comparison, **not** a built or tested correction. `PENDING` records an unresolved issue or competing proposal. No PR is adopted or upstream item closed by categorization. See [RNC exclusions](exclusions.md) for retired review slots.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| [#2173](https://github.com/gitleaks/gitleaks/issues/2173) | [#2174](https://github.com/gitleaks/gitleaks/pull/2174) (@moutons), [#2213](https://github.com/gitleaks/gitleaks/pull/2213) (@rylena), [#2344](https://github.com/gitleaks/gitleaks/pull/2344) (@delaysystem) | Document supported `--platform` choices in CLI help and README. | **SOURCE-CHECKED — TEST PENDING** | #2174 derives values from SCM enum but misses `detect` help; #2213 documents five values; #2344 adds `none` in help but omits README. Consolidate semantics and test both commands before adoption. | — |
| — | [#2084](https://github.com/gitleaks/gitleaks/pull/2084) (@vincent067) | Add `examples/README.md` and a new reusable GitHub Actions example. | **HOLD — NOT SAFE AS-IS** | The proposed `gitleaks git --baseline-path ...` claims to **generate** a baseline (it actually consumes one); default workflow runs an Action rather than the shown manual SARIF upload; suggested `--exit-code 0` conflicts with blocking deployments; workflow summary has malformed shell quoting around `.gitleaksignore`. Requires functional/security review, not a typo bundle. | — |
| [#1067](https://github.com/gitleaks/gitleaks/issues/1067) | — | Legacy `examples/leaky-repo.toml` benchmark config removed. | **ISSUE — REVIEW PENDING** | Original request refers to a removed 2022-era example; current master does not include it. Whether a maintained benchmark configuration is useful remains a separate decision, not solved by the docs typo bundle. | — |
| [#2295](https://github.com/gitleaks/gitleaks/issues/2295) | [#2306](https://github.com/gitleaks/gitleaks/pull/2306) (@rotirotirafa) | Document Firebase example-key provenance in generated GCP API-key allowlist. | **PATCH-REVIEWED — DOC/GENERATOR GATE** | PR #2306 adds explanatory comments in GCP rule generator/bundled output; verify both remain synchronized when regenerated. No claim that detection or allowlisting semantics change. | — |

### Completed CP-V9 documentation package

The eight upstream documentation PRs and baseline-security issue #1263 are implemented, with original sources, authors and scope recorded in [dispositions](exclusions.md) and the [engineering implementation ledger](https://github.com/RNC-Solutions-Group/gitleaks/blob/fix/validation/fixes-qol-tracking.md). They are not active implementation candidates. Do not mistake historic patch discussions for pending fixes.

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

**Completed 2026-11-12:** Baseline-security guidance #1263, corrected Docker anchor #2287, and verified mise instructions #2008 are implemented. They remain referenced in [dispositions](exclusions.md), not the active queue. Unsafe #2084 Actions recipe has not been adopted.

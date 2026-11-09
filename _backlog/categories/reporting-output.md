# Reporting Output — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**14 issues · 7 PRs · 15 work packages.**

SARIF/CSV/JSON/XML output, finding coordinates, report content, redaction, and output contracts. Verify escaping, permissions, and error paths.

**Review states:** `SOURCE-CHECKED` records a targeted source/patch comparison, **not** a built or tested correction. `PENDING` records an unresolved issue or competing proposal. No PR is adopted or upstream item closed by categorization. See [RNC exclusions](exclusions.md) for retired review slots.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| [#1626](https://github.com/gitleaks/gitleaks/issues/1626) | — | Include a copyable `git show <commit> -L ...:<path>` context command or equivalent historic-line link alongside findings. | **FEATURE REVIEW PENDING** | Reporter asks for additional finding output to inspect deleted/moved source lines, not a Git scanning fix. Current report has commit/file/line fields; choose between optional CLI decoration and short README example. No PR or implementation verified. | — |
| — | [#2205](https://github.com/gitleaks/gitleaks/pull/2205) (@trinathone) | Normalize multiline commit messages in CSV output, with a new CSV fixture. | **PENDING REVIEW** | Diff edits only `report/csv.go`, its tests, and an expected CSV; it does **not** touch `--log-opts` parsing for issue #2003. Assess CSV output semantics separately. | — |
| [#1866](https://github.com/gitleaks/gitleaks/issues/1866) | [#2326](https://github.com/gitleaks/gitleaks/pull/2326) (@bharathvemula) | Document `--report-format=json --report-path=-` for stdout and explain why `--report-format` alone produces no report. | **SOURCE-CHECKED — DOC REVIEW PENDING** | Current `report.StdoutReportPath = "-"` and `cmd/root.go` support this; the README **lacks the explicit example** added by PR #2326. This PR stays pending; no documentation adopted. | — |
| [#2304](https://github.com/gitleaks/gitleaks/issues/2304) | [#2307](https://github.com/gitleaks/gitleaks/pull/2307) (@jdotc1) | Existing report is removed before invalid report format is rejected; proposal: fix: preserve existing reports during report validation | **SOURCE-CHECKED — TEST PENDING** | Patch inspected: `cmd/root.go` on `master` creates the destination during report-option validation. | — |
| [#1962](https://github.com/gitleaks/gitleaks/issues/1962) | [#2238](https://github.com/gitleaks/gitleaks/pull/2238) (@hamodywe) | SARIF columns should count Unicode codepoints, not bytes. | **SOURCE-CHECKED — TEST PENDING** | Current `detect/location.go` subtracts byte offsets; PR #2238 uses `utf8.RuneCountInString` and adds multilingual fixtures. Distinct from PR #2215's line-start offset changes; assess interaction together. | — |
| [#1827](https://github.com/gitleaks/gitleaks/issues/1827) | [#2233](https://github.com/gitleaks/gitleaks/pull/2233) (@CBOSSX) | Invalid URI in SARIF file report; proposal: fix: escape SARIF artifact location URIs | **SOURCE-CHECKED — TEST PENDING** | Patch inspected: `master/report/sarif.go` has no proposed URI escaping. | — |
| [#1424](https://github.com/gitleaks/gitleaks/issues/1424) | [#2215](https://github.com/gitleaks/gitleaks/pull/2215) (@Eljees) | First-line vs later-line start/end column offsets. | **SOURCE-CHECKED — TEST PENDING** | Current `detect/location.go` advances `prevNewLine` from `pair[0]`; PR #2215 switches to actual post-newline line start `pair[1]` and revises test expectations. Distinct from #1962 Unicode counting. | — |
| [#2068](https://github.com/gitleaks/gitleaks/issues/2068) | [#2096](https://github.com/gitleaks/gitleaks/pull/2096) (@ChrisJr404) | Support --report-format gitlab-code-quality; proposal: feat(report): add gitlab-code-quality report format (#2068) | **LINKED — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1858](https://github.com/gitleaks/gitleaks/issues/1858) | — | The “level” field is not displayed in the sarif report | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1766](https://github.com/gitleaks/gitleaks/issues/1766) | — | Add "context" to findings | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1714](https://github.com/gitleaks/gitleaks/issues/1714) | — | Generate links for directory scans? | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1708](https://github.com/gitleaks/gitleaks/issues/1708) | — | Base64-decoded findings aren't printed with colour | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1357](https://github.com/gitleaks/gitleaks/issues/1357) | — | Incorrect Start and End Column Values in Gitleaks Report for Findings | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1048](https://github.com/gitleaks/gitleaks/issues/1048) | — | Generate multiple report formats (e.g. JSON and JUnit) from one scan. | **ISSUE — FEATURE REVIEW PENDING** | The **multi-format requirement remains open**; design against current `report.Reporter` semantics and test writer lifecycle and exit behavior. | — |
| [#988](https://github.com/gitleaks/gitleaks/issues/988) | — | No json report file generated in docker version v8.13.0 | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

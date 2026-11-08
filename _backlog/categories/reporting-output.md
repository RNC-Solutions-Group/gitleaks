# Reporting Output — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**25 upstream issues · 10 upstream PRs · 28 work packages.** Each ID is owned by this category once; linked PRs and issues share a row. Existing upstream state is unchanged.

SARIF/CSV/JSON/XML output, finding coordinates, reporting data preservation, and redaction. Tests should cover file permissions, escaping, error paths, and stable output contracts.

**Evidence rules:** “SOURCE-CHECKED” means a targeted patch/current-source comparison, **not** that a fix was built or tested. Most records are title/crosswalk-only and must stay pending until reproduced, diffed, and assessed on the fork. “Obsolete,” “superseded,” “adopted,” or “corrected” require separate evidence; no PR is adopted by routing it here.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| [#1866](https://github.com/gitleaks/gitleaks/issues/1866) | [#2326](https://github.com/gitleaks/gitleaks/pull/2326) (@bharathvemula) | Reporting: Include explicitly example for writing output as JSON to stdout; proposal: docs: document JSON report to stdout and report-format/report-path dependency | **LINKED — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#2304](https://github.com/gitleaks/gitleaks/issues/2304) | [#2307](https://github.com/gitleaks/gitleaks/pull/2307) (@jdotc1) | Existing report is removed before invalid report format is rejected; proposal: fix: preserve existing reports during report validation | **SOURCE-CHECKED — TEST PENDING** | Patch inspected: `cmd/root.go` on `master` creates the destination during report-option validation. | — |
| [#1962](https://github.com/gitleaks/gitleaks/issues/1962) | [#2238](https://github.com/gitleaks/gitleaks/pull/2238) (@hamodywe) | Column numbers should count unicode codepoints, not bytes; proposal: fix(detect): count columns in codepoints, not bytes | **LINKED — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1827](https://github.com/gitleaks/gitleaks/issues/1827) | [#2233](https://github.com/gitleaks/gitleaks/pull/2233) (@CBOSSX) | Invalid URI in SARIF file report; proposal: fix: escape SARIF artifact location URIs | **SOURCE-CHECKED — TEST PENDING** | Patch inspected: `master/report/sarif.go` has no proposed URI escaping. | — |
| [#1424](https://github.com/gitleaks/gitleaks/issues/1424) | [#2215](https://github.com/gitleaks/gitleaks/pull/2215) (@Eljees) | Off-by-one error in startColumn and endColumn calculation; proposal: fix: calculate columns from line starts | **LINKED — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#2068](https://github.com/gitleaks/gitleaks/issues/2068) | [#2096](https://github.com/gitleaks/gitleaks/pull/2096) (@ChrisJr404) | Support --report-format gitlab-code-quality; proposal: feat(report): add gitlab-code-quality report format (#2068) | **LINKED — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1992](https://github.com/gitleaks/gitleaks/pull/1992) (@bplaxco) | Add URL & JSON sources | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1858](https://github.com/gitleaks/gitleaks/issues/1858) | — | The “level” field is not displayed in the sarif report | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1766](https://github.com/gitleaks/gitleaks/issues/1766) | — | Add "context" to findings | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1714](https://github.com/gitleaks/gitleaks/issues/1714) | — | Generate links for directory scans? | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1708](https://github.com/gitleaks/gitleaks/issues/1708) | — | Base64-decoded findings aren't printed with colour | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1689](https://github.com/gitleaks/gitleaks/issues/1689) | — | Missing XML report format file | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1405](https://github.com/gitleaks/gitleaks/issues/1405) | — | Please update the deprecated way how to set output | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1357](https://github.com/gitleaks/gitleaks/issues/1357) | — | Incorrect Start and End Column Values in Gitleaks Report for Findings | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1348](https://github.com/gitleaks/gitleaks/issues/1348) | — | Alibaba secret in CSV format | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1312](https://github.com/gitleaks/gitleaks/issues/1312) | — | adding support for writing formatted reports to stdout when --report-path = "" and --report-format is != "" | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1275](https://github.com/gitleaks/gitleaks/issues/1275) | — | Help on finding right flag to export report in json | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1265](https://github.com/gitleaks/gitleaks/pull/1265) (@baruchiro) | feat: catch secret in XML child | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1243](https://github.com/gitleaks/gitleaks/issues/1243) | — | Output shows redacted secret if final new line not present | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1048](https://github.com/gitleaks/gitleaks/issues/1048) | [#1232](https://github.com/gitleaks/gitleaks/pull/1232) (@Coruscant11) | Support multiple report formats in a single execution; proposal: Support multiple report formats in a single execution | **LINKED — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1122](https://github.com/gitleaks/gitleaks/issues/1122) | — | Add option to completely disable color codes | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1110](https://github.com/gitleaks/gitleaks/issues/1110) | — | Unclear Sarif Short Description | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1065](https://github.com/gitleaks/gitleaks/issues/1065) | — | Validate report-path before scan | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| — | [#1043](https://github.com/gitleaks/gitleaks/pull/1043) (@icecream-monster) | + non-encrypted prv key in non-pem format(no -----begin header) check pt2 | **PR — REVIEW PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#1004](https://github.com/gitleaks/gitleaks/issues/1004) | — | Add logging around failure to write report | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#988](https://github.com/gitleaks/gitleaks/issues/988) | — | No json report file generated in docker version v8.13.0 | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#983](https://github.com/gitleaks/gitleaks/issues/983) | — | print resulting configuration for debug purpose if extend directive is used | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |
| [#908](https://github.com/gitleaks/gitleaks/issues/908) | — | Default Generic API rule - does not scan secrets in XML tags | **ISSUE — REPRO PENDING** | Title/crosswalk metadata only; master equivalence and runtime behavior **not verified**. | — |

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

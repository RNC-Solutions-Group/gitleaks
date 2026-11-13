# Security Hardening — RNC review queue

**Current active queue (2026-11-13):** Completed/declined original issue and PR packages were removed from this category. See the [RNC disposition register](exclusions.md) and [fixes/QOL implementation ledger](https://github.com/RNC-Solutions-Group/gitleaks/blob/fix/validation/fixes-qol-tracking.md). Remaining original candidates below are not closed by overlap with those corrections.

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**10 issues · 7 PRs · 13 accounting work packages (12 table rows; #2305/#2291 retain separate CVE acceptance tracks).** This category covers dependency CVEs, secrets exposure, report-file permissions, execution hardening, and supply-chain controls.

**Review states:** CVE proposals #2308 and #2336 have operator-reported test evidence but remain unadopted. Other PR diffs have been inspected against pinned fork source; behavior remains untested unless explicitly stated. Do not reproduce potential credential values in this ledger. Out-of-scope and superseded entries reside only in [RNC exclusions](exclusions.md).

## Tracked work

| PR and contributor | Linked issue(s) | Proposed security change | RNC state | Evidence and remaining gate |
| --- | --- | --- | --- | --- |
| [#2308](https://github.com/gitleaks/gitleaks/pull/2308) — @sanjitpawar2007-hue; head `71c6e141c640d5241b77745164fab8279c6e78f1` | [#2305](https://github.com/gitleaks/gitleaks/issues/2305), [#2291](https://github.com/gitleaks/gitleaks/issues/2291) | Upgrade Go dependency modules (notably `x/sys`, `x/text`, `x/crypto`, `rardecode`) for reported CVEs; raises minimum Go toolchain. | **REVIEWED — TESTING** | PR diff changes `go.mod`/`go.sum`: Go `1.24.11` → `1.26.0`, plus `rardecode`, `x/crypto`, `x/sync`, `x/sys` and `x/text`. User-reported Go tests passed on 1.27.3. **Unmodified PR Docker build fails** against the 1.24 builder; RNC's `gitleaks:pr2308mod` includes a Dockerfile-only adjustment. Preserve tested image/source provenance, CVE inventory and both issue-specific acceptance tracks; **not adopted**. |
| [#2336](https://github.com/gitleaks/gitleaks/pull/2336) — @VISHNU0906; head `817e1ef908e02a72460c2cf087c9cfcdd5b0b2e0` | [#2330](https://github.com/gitleaks/gitleaks/issues/2330) | Upgrade `github.com/klauspost/compress` from v1.18.0 to v1.18.7 for reported CVE-2026-63209; consider interplay with archive handling. | **REVIEWED — TESTING** | Diff changes `go.mod`/`go.sum` to `klauspost/compress v1.18.7` (from `v1.18.0`); it complements #2308 rather than duplicating it. User-reported Go suite passed. Pinned `xz v0.5.15` and `rardecode v2.2.0` already match the issue's other claimed fixed-version targets. No digest-pinned combined image, current vulnerability scan or full archive acceptance is recorded; **not adopted**. |

## Current technical baselines

- Fork `master/go.mod`: `go 1.24.11`; `klauspost/compress v1.18.0`; `nwaples/rardecode/v2 v2.2.0`; `ulikunitz/xz v0.5.15`; `x/crypto v0.35.0`; `x/sys v0.30.0`; `x/text v0.29.0`.
- RNC `fix/docker-security-bases` Dockerfile: `golang:1.27.3` builder and `alpine:3.24` runtime. **This branch does not itself modify the Go module pins.**
- [RNC test harness](https://github.com/RNC-Solutions-Group/gitleaks/tree/test/security-integration-2308-2336/_rnc-build) reports 19/19 scenarios and 121/121 assertions passing for its tested images; 34/38 Docker-applicable fixture files validated, four exercised-only, and 75 upstream Go-only contracts not imported. This establishes parity over the tested scope **only**, not complete release acceptance.
- Go test evidence for the upstream PRs is operator-reported from earlier local work; it should be preserved alongside exact tested commit hashes and command outputs.

## Additional security and supply-chain candidates — source and patch review

All source comparisons below target fork `master` at `5f6ec56cfa6f`. **SOURCE-CHECKED** identifies baseline behavior; **PATCH-REVIEWED** identifies inspected upstream diffs. Neither means the issue was reproduced, the patch was tested, or any code was integrated.

| Upstream issue(s) | Upstream PR(s) / contributor(s) | Scope | RNC disposition | Gate |
| --- | --- | --- | --- | --- |
| [#2110](https://github.com/gitleaks/gitleaks/issues/2110) | — | Potential upstream Slack **bot** token exposure. | **INCIDENT TRIAGE — UNVERIFIED** | Separate exposure report; maintainer-only validation, immediate revoke/rotate if authentic, then assess history and downstream artifacts. **No token values or credential-derived examples in docs or tests.** Not a scanner-rule fix. |
| [#2111](https://github.com/gitleaks/gitleaks/issues/2111) | — | Potential upstream Slack **user** token exposure. | **INCIDENT TRIAGE — UNVERIFIED** | Distinct credential class; verify affected identity, revoke/rotate and review access if authentic. Do not infer duplication with bot or app-level report from similar wording. No values copied to RNC. |
| [#2112](https://github.com/gitleaks/gitleaks/issues/2112) | — | Potential upstream Slack **app-level** token exposure. | **INCIDENT TRIAGE — UNVERIFIED** | Separate validation/rotation/history-remediation track. Needs upstream maintainer or secret owner authorization; do not attempt to validate by using reported credentials. |
| — | [#2237](https://github.com/gitleaks/gitleaks/pull/2237) (@bunlongheng) | Spreadsheet formula injection through untrusted CSV report cells. | **PATCH-REVIEWED — CSV MATRIX PENDING** | Master `report/csv.go` writes finding fields verbatim. PR adds a leading-apostrophe sanitizer for first-byte `=`, `+`, `-`, `@`, tab and CR; **no tests included**. Test leading whitespace/control variants, malicious scanner-controlled fields, Excel/LibreOffice/Sheets interpretation, and parser/machine-consumer compatibility before adoption. Distinct from plaintext file permissions. |
| — | [#637](https://github.com/gitleaks/gitleaks/pull/637) (@varunsh-coder) | Constrain CI `GITHUB_TOKEN` scope. | **PATCH-REVIEWED — ADAPT TO CURRENT WORKFLOW** | 2022 patch adds `permissions: contents: read` to two legacy jobs; current `.github/workflows/test.yml` has a **single** matrix `test` job and no explicit permissions. Old diff does not apply as written; least-privilege requirement still needs a current-workflow edit and Actions validation. |
| [#1243](https://github.com/gitleaks/gitleaks/issues/1243) | — | Reported redaction escape in `Finding:` context without a trailing source newline. | **SOURCE-CHECKED — HIGH PRIORITY OUTPUT REPRO** | Current `report/finding.go:Redact` replaces only the finding's selected `Secret` inside `Line`/`Match`; a different secret on the same line may remain visible. Original no-terminal-newline claim is **not yet reproduced** on pinned master. Use synthetic multi-secret fixtures, `--redact`, newline/no-newline, Git/dir/pre-commit and verbose output; retain no raw-secret logs. Separate from #2299 permission control. |
| — | [#2036](https://github.com/gitleaks/gitleaks/pull/2036) (@Tanmaykaturi) | Replace realistic-looking Facebook token fixture in generator source. | **PATCH-REVIEWED — FIXTURE/SAST GATE** | One-line replacement in `cmd/generate/config/rules/facebook.go` preserves a token-shaped positive case. Confirm new fixture is safely synthetic, detector-positive, and accepted by configured SAST rules; avoid embedding either token-shaped value in review artifacts. |

**Priority checks (current 2026-11-13):** #2299/#2312 F02 permissions and #1572 F05 bounded ignore loader were corrected and moved to [dispositions](exclusions.md); retain separately scoped publication, namespace replacement and Windows-permission residuals, 3) #1243 redaction leakage, 4) #2237 CSV injection, 5) #2176/#2206 container least privilege. The Slack reports are confidential incident-response tracks, not scanner PRs. CVE patch tests and image baselines remain independently gated. No scanner changes have been adopted by this documentation pass.

## Promotion and completion gate

1. Pin a known source commit, PR head SHA, Dockerfile changes, toolchain, and image digest for **each** candidate.
2. Run a CVE inventory with Trivy/Scout against actual final images and record exact package versions, severity, database timestamp, and remaining findings.
3. Re-run the complete applicable Docker regression and upstream Go tests, including archive/rar/xz/compress combinations, negative fixtures, symlinks, and file-chunk boundaries.
4. Compare to independent expected findings—not blindly to release-image behavior—and report all fixture exclusions.
5. Record the user-approved RNC branch/PR/commit receiving each contribution. Until then, status remains **REVIEWED — TESTING**.
6. Only after evidence and owner approval promote security changes downstream according to the established branch policy; do not merge community PRs or alter upstream state from these docs.

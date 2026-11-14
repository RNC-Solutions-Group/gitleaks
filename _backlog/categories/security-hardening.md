# Security Hardening — RNC review queue

**Current active queue (post-CP-V11, 2026-11-13):** Completed/declined original issue and PR packages were removed from this category. See the [RNC disposition register](exclusions.md) and [RNC-MOD engineering reference](https://github.com/RNC-Solutions-Group/gitleaks/blob/rnc-mod/RNC-MOD.md). Remaining original candidates below are not closed by overlap with those corrections.

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**4 active issues · 1 active PR · 5 accounting work packages.** This category covers dependency CVEs, secrets exposure, report-file permissions, execution hardening, and supply-chain controls.

**Review states:** CVE dependency-version proposals #2308/#2336 and their three linked issues were superseded by newer pins in `fix/validation` and moved to [dispositions](exclusions.md). Final released-artifact vulnerability and SBOM verification remains an independent gate. Other PR diffs were reviewed against the pinned upstream baseline; do not claim untested patches were adopted. Do not reproduce potential credential values in this ledger. Out-of-scope and superseded entries reside only in [RNC exclusions](exclusions.md).

## Tracked work

| PR and contributor | Linked issue(s) | Proposed security change | RNC state | Evidence and remaining gate |
| --- | --- | --- | --- | --- |

## CVE proposals reconciled against current source — 2026-11-13

The original #2305/#2291/#2330 issues and #2308/#2336 proposals no longer need *those particular dependency bumps*: current RNC `go.mod` exceeds all cited fixed-version minima; Go module minimum `1.26.0`, Docker builder `1.27.4`. RNC CP-V10 [CI 12/12](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43321594496) includes `govulncheck` and container vulnerability jobs. This does **not** certify any unbuilt/published release artifact or that every future vulnerability database finds zero issues. See [RNC dispositions](exclusions.md) for exact version crosswalk; current module/version release-gate work remains open.

## Historical technical baselines (master, not current validation)

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
| [#1243](https://github.com/gitleaks/gitleaks/issues/1243) | — | Reported redaction escape in `Finding:` context without a trailing source newline. | **SOURCE-CHECKED — HIGH PRIORITY OUTPUT REPRO** | Current `report/finding.go:Redact` replaces only the finding's selected `Secret` inside `Line`/`Match`; a different secret on the same line may remain visible. Original no-terminal-newline claim is **not yet reproduced** on pinned master. Use synthetic multi-secret fixtures, `--redact`, newline/no-newline, Git/dir/pre-commit and verbose output; retain no raw-secret logs. Separate from #2299 permission control. |
| — | [#2036](https://github.com/gitleaks/gitleaks/pull/2036) (@Tanmaykaturi) | Replace realistic-looking Facebook token fixture in generator source. | **PATCH-REVIEWED — FIXTURE/SAST GATE** | One-line replacement in `cmd/generate/config/rules/facebook.go` preserves a token-shaped positive case. Confirm new fixture is safely synthetic, detector-positive, and accepted by configured SAST rules; avoid embedding either token-shaped value in review artifacts. |

**Priority checks (post-V11, 2026-11-13):** The original F02 permissions, F05 ignore loader, three CVE dependency-version reports, and obsolete CI token-permission proposal #637 have moved to [dispositions](exclusions.md). Retain separately scoped report publication, Windows ACL, #1243 redaction, and fixture/SAST review; CSV formula injection was corrected in V12 and is dispositioned. Slack reports are confidential incident-response tracks, not scanner PRs. Release CVE/SBOM and image provenance remain independently gated. No scanner changes were made by this reconciliation.

## Promotion and completion gate

1. Pin a known source commit, PR head SHA, Dockerfile changes, toolchain, and image digest for **each** candidate.
2. Run a CVE inventory with Trivy/Scout against actual final images and record exact package versions, severity, database timestamp, and remaining findings.
3. Re-run the complete applicable Docker regression and upstream Go tests, including archive/rar/xz/compress combinations, negative fixtures, symlinks, and file-chunk boundaries.
4. Compare to independent expected findings—not blindly to release-image behavior—and report all fixture exclusions.
5. Record the user-approved RNC branch/PR/commit receiving each contribution. Until then, status remains **REVIEWED — TESTING**.
6. Only after evidence and owner approval promote security changes downstream according to the established branch policy; do not merge community PRs or alter upstream state from these docs.

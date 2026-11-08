# Security hardening — reviewed and testing

Primary active category for upstream dependency-security PRs [#2308](https://github.com/gitleaks/gitleaks/pull/2308) and [#2336](https://github.com/gitleaks/gitleaks/pull/2336) and associated upstream issues. The records below reflect *review and testing*, **not** integrated or released fixes.

The [closed ledger](closed.md) is authoritative: excluded RNC-closed items are not duplicated in this active table.

## Tracked work

| PR and contributor | Linked issue(s) | Proposed security change | RNC state | Evidence and remaining gate |
| --- | --- | --- | --- | --- |
| [#2308](https://github.com/gitleaks/gitleaks/pull/2308) — @sanjitpawar2007-hue; head `71c6e141c640d5241b77745164fab8279c6e78f1` | [#2305](https://github.com/gitleaks/gitleaks/issues/2305), [#2291](https://github.com/gitleaks/gitleaks/issues/2291) | Upgrade Go dependency modules (notably `x/sys`, `x/text`, `x/crypto`, `rardecode`) for reported CVEs; raises minimum Go toolchain. | **REVIEWED — TESTING** | User-reported Go regression tests passed on Go 1.27.3. **Original upstream PR Dockerfile does not build** because its Go requirement exceeds the Docker builder version. RNC's `gitleaks:pr2308mod` uses a Dockerfile-only correction; it is not an unmodified-PR image. |
| [#2336](https://github.com/gitleaks/gitleaks/pull/2336) — @VISHNU0906; head `817e1ef908e02a72460c2cf087c9cfcdd5b0b2e0` | [#2330](https://github.com/gitleaks/gitleaks/issues/2330) | Upgrade `github.com/klauspost/compress` from v1.18.0 to v1.18.7 for reported CVE-2026-63209; consider interplay with archive handling. | **REVIEWED — TESTING** | User-reported Go regression suite passed. No standalone / integrated Docker image with a verified image-level CVE clearance and complete archive behavior acceptance has been established in this ledger. |

## Current technical baselines

- Fork `master/go.mod`: `go 1.24.11`; `klauspost/compress v1.18.0`; `nwaples/rardecode/v2 v2.2.0`; `ulikunitz/xz v0.5.15`; `x/crypto v0.35.0`; `x/sys v0.30.0`; `x/text v0.29.0`.
- RNC `fix/docker-security-bases` Dockerfile: `golang:1.27.3` builder and `alpine:3.24` runtime. **This branch does not itself modify the Go module pins.**
- [RNC test harness](https://github.com/RNC-Solutions-Group/gitleaks/tree/test/security-integration-2308-2336/_rnc-build) reports 19/19 scenarios and 121/121 assertions passing for its tested images; 34/38 Docker-applicable fixture files validated, four exercised-only, and 75 upstream Go-only contracts not imported. This establishes parity over the tested scope **only**, not complete release acceptance.
- Go test evidence for the upstream PRs is operator-reported from earlier local work; it should be preserved alongside exact tested commit hashes and command outputs.

## Issue-to-PR mapping

| Upstream issue | Upstream state | RNC assessment |
| --- | --- | --- |
| [#2305](https://github.com/gitleaks/gitleaks/issues/2305) — CVEs in v8.30.1 | Open | Linked to [#2308](https://github.com/gitleaks/gitleaks/pull/2308); **reviewed — testing**; confirm each CVE against final image SBOM. |
| [#2291](https://github.com/gitleaks/gitleaks/issues/2291) — Go CVE vulnerabilities | Open | Linked to [#2308](https://github.com/gitleaks/gitleaks/pull/2308); **reviewed — testing**; distinguish Go stdlib from individually pinned modules. |
| [#2330](https://github.com/gitleaks/gitleaks/issues/2330) — compression/archive module CVEs | Open | Linked to [#2336](https://github.com/gitleaks/gitleaks/pull/2336); **reviewed — testing**; package-specific findings require image verification. |

## Promotion and completion gate

1. Pin a known source commit, PR head SHA, Dockerfile changes, toolchain, and image digest for **each** candidate.
2. Run a CVE inventory with Trivy/Scout against actual final images and record exact package versions, severity, database timestamp, and remaining findings.
3. Re-run the complete applicable Docker regression and upstream Go tests, including archive/rar/xz/compress combinations, negative fixtures, symlinks, and file-chunk boundaries.
4. Compare to independent expected findings—not blindly to release-image behavior—and report all fixture exclusions.
5. Record the user-approved RNC branch/PR/commit receiving each contribution. Until then, status remains **REVIEWED — TESTING**.
6. Only after evidence and owner approval promote security changes downstream according to the established branch policy; do not merge community PRs or alter upstream state from these docs.

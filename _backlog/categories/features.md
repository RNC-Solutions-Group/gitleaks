# Optional Features — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**6 issues · 3 active PRs · 6 work packages.** Optional interfaces and developer-experience capabilities not required for core scan correctness.

**Review states:** All six packages remain unadopted; direct upstream patch inspection is not runtime verification. Three previously routed issues were reassigned after the 2026-11-09 by-item pass: [#1397 / PR #2216](configuration-allowlists.md) to Configuration & Allowlists; [#2017 / closed PR #2018](integrations-release.md) to Integrations & Release; and [#1511](exclusions.md) to the RNC exclusions ledger. The original open-PR snapshot and contributor links remain preserved.

## Candidates

These candidates have not been adopted or corrected by RNC. PR #1992 was re-homed from Scan Integrity because it adds **new optional URL/JSON input modes** rather than repairing coverage of existing file sources.

| Upstream issue(s) | Upstream PR(s) / contributor(s) | Scope | RNC disposition | Gate |
| --- | --- | --- | --- | --- |
| — | [#1992](https://github.com/gitleaks/gitleaks/pull/1992) (@bplaxco) | New `url` command and `json` structured-source reader; fetch URLs, scan HTTP bodies and JSON string leaves, preserve JSON-path provenance. | **PATCH-REVIEWED — REJECT PATCH AS-IS; RETAIN FEATURE** | New `cmd/url.go`, `cmd/json.go`, `sources/url.go`, `sources/json.go`: no HTTP timeout or response/body-size ceiling, unrestricted redirects/IP targets, unbounded stdin read, a nested non-200 body leak, and errors logged with nil-return paths that can disguise incomplete scans as clean. Review SSRF, explicit network authorization, URL credential logging, TLS, JSON provenance/limits and fail-closed status. Moved from Scan Integrity 2026-11-09. |
| [#2152](https://github.com/gitleaks/gitleaks/issues/2152), [#2153](https://github.com/gitleaks/gitleaks/issues/2153) | [#2154](https://github.com/gitleaks/gitleaks/pull/2154) (@MChorfa) | Embeddable scanner SDK plus canary exemptions and vector state. | **PATCH-REVIEWED — REWORK BEFORE TEST/ADOPTION** | SDK `Scanner` defaults decoding to zero depth while CLI defaults to five: parity risk. PR uses process-global `viper.Reset()`, unsynchronized exemption maps, shallow config copying, and loses partial findings on scan errors. Test API lifecycle, concurrent use, failure propagation, and actual canary/vector semantics before considering an adapter. |
| [#1869](https://github.com/gitleaks/gitleaks/issues/1869) | [#2056](https://github.com/gitleaks/gitleaks/pull/2056) (@raajheshkannaa) | MCP-server command for coding-agent initiated local scans. | **PATCH-REVIEWED — HARDEN BEFORE ADOPTION** | Local stdio MCP is a useful agent-facing adapter, but `redact=0` is caller-configurable, path roots are unrestricted, payload budget is unlimited by default, and Git-option parsing/ignored-file load failures need stronger validation. Require permitted roots, output redaction, bounded scans, and no accidental secret disclosure to AI clients. |
| [#1054](https://github.com/gitleaks/gitleaks/issues/1054) | — | Deduplicate identical secrets across Git commits and optionally share ignore fingerprints. | **SOURCE-CHECKED — REPORT-ONLY GROUPING PREFERRED** | Current finding fingerprints include commit/path/rule/line, not secret identity. Group repeated occurrences by a protected normalized-value digest for reporting while preserving per-commit findings/provenance and ignore semantics. Do not silently suppress renewed exposures. |
| [#1013](https://github.com/gitleaks/gitleaks/issues/1013) | — | Optional validation of discovered credentials against their issuing services. | **ISSUE — DESIGN/SECURITY REVIEW PENDING** | Live provider checks must be explicitly authorized, network-bound, rate-limited and redacted, with `verified`, `invalid`, `unknown`, and `error` distinguished. Expired credentials are not automatically safe and no unapproved verification egress should occur. |
| [#1456](https://github.com/gitleaks/gitleaks/issues/1456) | — | Checksum-aware declarative rules for provider tokens (e.g. Adler-32 suffix validation). | **ISSUE — DESIGN/TEST PENDING** | Deterministic **offline** checksum validation extends TOML rule schema/evaluator (e.g. Adler-32) rather than one provider pattern. Validate capture groups, checksum variants, false-positive tradeoffs, generator parity and backward compatibility. |

## Evaluation gate

1. Confirm the actual gap on current `master`, with a safe reproducer; **PENDING REVIEW** does not mean the proposal is endorsed.
2. Establish SDK parity with CLI defaults (especially decode depth), concurrency safety, lifecycle behavior, and error/partial-result semantics; compare MCP as an **external** adapter before embedding the SDK.
3. Compare surviving proposals with the current scanner behavior, security implications, and tests before any future adoption decision.
4. PR #1992 introduces active outbound HTTP fetching. Before any experimental run, require an explicit opt-in network target, redirect/IP protections, timeout and byte budgets and fail-closed errors; do not use production secrets as test data. PR #2056 likewise requires allowed filesystem roots, redacted outputs and bounded MCP requests.
5. Preserve original contributor attribution and separate RNC implementation/verification records if adoption occurs in another branch.

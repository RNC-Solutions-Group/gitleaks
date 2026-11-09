# Optional Features — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**7 issues · 3 active PRs · 6 work packages.** Optional interfaces and developer-experience capabilities not required for core scan correctness.

**Review states:** All proposals remain unadopted. Issue #2017 links a **closed, unmerged** PR #2018 for context, not as an active PR. [RNC exclusions](exclusions.md) are recorded separately.

## Candidates

These candidates have not been adopted or corrected by RNC.

| Upstream issue(s) | Upstream PR(s) / contributor(s) | Scope | RNC disposition | Gate |
| --- | --- | --- | --- | --- |
| [#2152](https://github.com/gitleaks/gitleaks/issues/2152), [#2153](https://github.com/gitleaks/gitleaks/issues/2153) | [#2154](https://github.com/gitleaks/gitleaks/pull/2154) (@MChorfa) | Embeddable scanner SDK plus canary exemptions and vector state. | **SOURCE-CHECKED — TEST PENDING** | One SDK PR addresses two **distinct** feature requests; neither issue is a duplicate. Review API stability, side effects, serialization and security model as one implementation package. |
| [#1869](https://github.com/gitleaks/gitleaks/issues/1869) | [#2056](https://github.com/gitleaks/gitleaks/pull/2056) (@raajheshkannaa) | MCP-server command for coding-agent initiated local scans. | **PENDING REVIEW** — distinct active feature | PR adds `cmd/mcp.go` plus `go.mod`/`go.sum`; no `cmd/mcp.go` on master. Requires separate MCP access, output-handling, and dependency review. |
| [#2017](https://github.com/gitleaks/gitleaks/issues/2017) | [#2018](https://github.com/gitleaks/gitleaks/pull/2018) (@prashanthpai), **closed unmerged** — contextual, not active PR | Prefer standard Go regex for **published binaries** under constrained address-space limits. | **PENDING REVIEW** — release build remains distinct | Current `.goreleaser.yml` explicitly uses `gore2regex`, unlike default untagged stdlib builds; linked one-line removal PR #2018 was closed unmerged. Reassess build/release ownership later without changing another category in this pass. |
| [#1511](https://github.com/gitleaks/gitleaks/issues/1511) | — | Configurable secret replacement token (Drupal/YAML-safe output). | **PENDING REVIEW** — no equivalent CLI option | `cmd/root.go` exposes percentage `--redact`; `report/finding.go` hardcodes `REDACTED` at 100%. This does not implement a selectable replacement string. |
| [#1397](https://github.com/gitleaks/gitleaks/issues/1397) | [#2216](https://github.com/gitleaks/gitleaks/pull/2216) (@hantaee22) | Signed line-offset `gitleaks:allow` directives. | **PENDING REVIEW** — unique open PR | Current `detect/detect.go` only checks the finding's own line; #2216 changes detection, adds offset parsing/tests, and documents semantics. Nothing equivalent on master. |
| [#1054](https://github.com/gitleaks/gitleaks/issues/1054) | — | Deduplicate identical secrets across Git commits and optionally share ignore fingerprints. | **PENDING REVIEW** — distinct behavior | Current `detect/baseline.go` compares commit metadata; default fingerprints are commit-scoped. Existing per-commit processing is not the requested cross-commit secret deduplication. |

## Evaluation gate

1. Confirm the actual gap on current `master`, with a safe reproducer; **PENDING REVIEW** does not mean the proposal is endorsed.
2. Keep issue #2017's release-binary build tags separate from ordinary source builds, and treat PR #2018 as closed **without merge**, not as a shipped correction.
3. Compare surviving proposals with the current scanner behavior, security implications, and tests before any future adoption decision.
4. Preserve original contributor attribution and separate RNC implementation/verification records if adoption occurs in another branch.

# Features — Gitleaks candidate register

This is the **single active category** for feature and developer-experience proposals selected for RNC evaluation. Issues and PRs for the same work package belong together here, even when one has a different title-only classification in the original snapshot. The [exclusions ledger](exclusions.md) remains a separate exclusion register.

**Duplicate-category scrub (2026-11-08):** Issues #2152 and #2153 plus PR #2154 were added as one SDK package. These two issues request different capabilities; neither was excluded as a duplicate.

## Focused Optional Features scrub — 2026-11-08

**Source baseline:** RNC fork `master` at [`5f6ec56c`](https://github.com/RNC-Solutions-Group/gitleaks/tree/5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c). Seven upstream issue bodies, two active upstream PR patches, and the issue-linked closed PR #2018 were compared with current fork source. This is **source review only**, not a fresh build, runtime reproduction, adoption, or independent patch correctness certification.

**Disposition:** two issue-only review slots excluded; five issues and two associated open PRs remain in five active work packages. The older v7/Go 1.16 [PR #652](exclusions.md) was already excluded as obsolete; it is **not** a new elimination in this pass.

- [Issue #1246](exclusions.md): **CLOSED_RNC — already covered by master**. The `--ignore-gitleaks-allow` flag disables inline suppression in the detector when enabled.
- [Issue #829](exclusions.md): **CLOSED_RNC — scope exception / existing workflow**. Users can include the complete matched line using `--report-format=template` and `.Line`; master includes a tested full-line JSON template example. This is **not** a native `-add_full_row` switch or an added field in ordinary JSON/CSV reports. Full-line export may disclose more secrets and remains an explicit template choice.

These decisions, links to supporting source/tests, and reopen conditions are canonical in [exclusions.md](exclusions.md). All upstream issue/PR states and original JSON/CSV exports remain unchanged.

## Evaluation gate

1. Confirm the actual gap on current `master`, with a safe reproducer; **PENDING REVIEW** does not mean the proposal is endorsed.
2. Keep issue #2017's release-binary build tags separate from ordinary source builds, and treat PR #2018 as closed **without merge**, not as a shipped correction.
3. Compare surviving proposals with the current scanner behavior, security implications, and tests before any future adoption decision.
4. Preserve original contributor attribution and separate RNC implementation/verification records if adoption occurs in another branch.

## Remaining Optional Features — pending independent review

These candidates have not been adopted or corrected by RNC.

| Upstream issue(s) | Upstream PR(s) / contributor(s) | Scope | RNC disposition | Gate |
| --- | --- | --- | --- | --- |
| [#2152](https://github.com/gitleaks/gitleaks/issues/2152), [#2153](https://github.com/gitleaks/gitleaks/issues/2153) | [#2154](https://github.com/gitleaks/gitleaks/pull/2154) (@MChorfa) | Embeddable scanner SDK plus canary exemptions and vector state. | **RECLASSIFIED — SOURCE-CHECKED, TEST PENDING** | One SDK PR addresses two **distinct** feature requests; neither issue is a duplicate. Review API stability, side effects, serialization and security model as one implementation package. |
| [#1869](https://github.com/gitleaks/gitleaks/issues/1869) | [#2056](https://github.com/gitleaks/gitleaks/pull/2056) (@raajheshkannaa) | MCP-server command for coding-agent initiated local scans. | **PENDING REVIEW** — distinct active feature | PR adds `cmd/mcp.go` plus `go.mod`/`go.sum`; no `cmd/mcp.go` on master. Requires separate MCP access, output-handling, and dependency review. |
| [#2017](https://github.com/gitleaks/gitleaks/issues/2017) | [#2018](https://github.com/gitleaks/gitleaks/pull/2018) (@prashanthpai), **closed unmerged** — contextual, not active PR | Prefer standard Go regex for **published binaries** under constrained address-space limits. | **PENDING REVIEW** — release build remains distinct | Current `.goreleaser.yml` explicitly uses `gore2regex`, unlike default untagged stdlib builds; linked one-line removal PR #2018 was closed unmerged. Reassess build/release ownership later without changing another category in this pass. |
| [#1511](https://github.com/gitleaks/gitleaks/issues/1511) | — | Configurable secret replacement token (Drupal/YAML-safe output). | **PENDING REVIEW** — no equivalent CLI option | `cmd/root.go` exposes percentage `--redact`; `report/finding.go` hardcodes `REDACTED` at 100%. This does not implement a selectable replacement string. |
| [#1397](https://github.com/gitleaks/gitleaks/issues/1397) | [#2216](https://github.com/gitleaks/gitleaks/pull/2216) (@hantaee22) | Signed line-offset `gitleaks:allow` directives. | **PENDING REVIEW** — unique open PR | Current `detect/detect.go` only checks the finding's own line; #2216 changes detection, adds offset parsing/tests, and documents semantics. Nothing equivalent on master. |
| [#1054](https://github.com/gitleaks/gitleaks/issues/1054) | — | Deduplicate identical secrets across Git commits and optionally share ignore fingerprints. | **PENDING REVIEW** — distinct behavior | Current `detect/baseline.go` compares commit metadata; default fingerprints are commit-scoped. Existing per-commit processing is not the requested cross-commit secret deduplication. |



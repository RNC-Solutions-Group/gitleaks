# CP-V12 Backlog Reconciliation — 2026-11-13

**Baseline:** `rnc-mod` at `3a4635f82dcd0a216f9d2b88ebbd72341b3b60cc` (merged V12 PR #5)
**Assessment branch:** `feat/issues-and-pr`
**Original snapshot:** 2026-11-08; 285 issues, 252 PRs, 388 accounting slots
**Authority:** current source and tests at the merged V12 baseline, not historical patch summaries.

## 1. Method

Read the full category register, exclusions register, `RNC-MOD.md`, V10/V11/V12 provenance records, current `rnc-mod` source and tests. For each V12 target, compared the final behavior rather than repeating the pre-V12 review state.

A package was retired only when the merged source/test evidence established one of:

- full correction under the same or a different RNC implementation;
- adopted contributor implementation with preserved ancestry;
- explicit policy decision that intentionally retains the old default while adding an opt-in control.

Partial fixes remain active with their residual stated. No upstream issue or PR was closed or represented as merged upstream.

## 2. Dispositioned V12 packages

Nine packages (8 issues, 7 PRs) moved to the exclusion register:

| Upstream item(s) | Actual V12 outcome | Final evidence |
|---|---|---|
| #1316 / PR #2311 | Fixed: pretty/date/color formatting determinism; contributor commit retained. | `2790da3`, merged V12; hostile-format tests. |
| #1814 | Fixed: textconv/ext-diff blocked; no repository helper execution. | V12 Git flags; textconv marker and staged tests. |
| #1338 / PR #2246 | Corrected through V11/V12 Git lifecycle and parser completeness; unsupported formats fail closed. | `5bed958`, `9092703`; malformed/truncated/section tests. |
| PR #2237 | Adopted and expanded CSV formula neutralization. | `904871b`; Unicode/control/fullwidth/RFC 4180 tests. |
| #1962 / PR #2238 | Corrected through RNC SARIF-only codepoint fields; historical JSON/CSV columns retained. PR not cherry-picked. | V12 location/SARIF tests; design influence only. |
| #1827 / PR #2233 | Adopted and expanded SARIF artifact URI escaping, including Windows drive/UNC. | `7021b7d`; URI matrix. |
| #2121 / PR #2271 | Adopted regex match reuse. | `8e787c2`; benchmark and parity suite. |
| #1893 | Fixed CRLF detector/baseline identity including legacy baseline migration. | V12 detector/baseline tests. |
| #1409 / PR #1918 | Corrected through the new `gitleaks files` command and explicit-file hooks; broad directory hook not adopted. | V12 source/CLI/real pre-commit tests. |

## 3. Targets retained active

| Upstream item(s) | Status | Reason retained |
|---|---|---|
| #1651 / PR #2324 | **PARTIALLY CORRECTED — RESIDUAL CONFIRMED** | 100–300 KB boundary splits are corrected, but lines above the 1 MiB finish budget fail closed rather than being fully scanned. |
| #1243 | **PARTIALLY CORRECTED — RESIDUAL CONFIRMED** | Known adjacent/multiline/empty-capture redaction improved; unknown non-finding text and caller-provided template surfaces remain. |
| #2019 | **NOT REPRODUCED — EVIDENCE INSUFFICIENT** | Finite probes show depth-linear memory amplification; no reproducible 24 GB case. Default decode depth 5 retained; V13-15 must measure independently. |
| #1028 / PR #2301 | **POLICY DECISION REQUIRED** | Merge-only default remains unchanged; first-parent is opt-in. Default change alters duplicates/fingerprints. |

## 4. Accounting

| Measure | Before V12 reconciliation | V12 moves | After |
|---|---:|---:|---:|
| Active issues | 194 | -8 | **186** |
| Active PRs | 151 | -7 | **144** |
| Active work-package slots | 235 | -9 | **226** |
| Dispositioned issues | 91 | +8 | **99** |
| Dispositioned PRs | 101 | +7 | **108** |
| Dispositioned/retired slots | 153 | +9 | **162** |
| Assigned total | 285 / 252 / 388 | 0 | **285 / 252 / 388** |
| Unassigned | 0 | 0 | **0** |

These nine are the only V12 packages moved in this reconciliation. No package was retired merely because its upstream PR existed or was reviewed.

### Category-header reconciliation note

The pre-existing category headline counts did not add up to the prior global 194/151/235 accounting; they were stale routing summaries, not a second accounting system. This pass updated the five affected category headers by the exact package moves (Security, Scan Integrity, Git/Hooks, Reporting, Performance). It does **not** invent a new per-category distribution for the remaining legacy discrepancy. The global totals above are the authoritative one-owner accounting; package-level rows and the exclusion register are the audit trail. A future full category recount should reconcile every headline against the immutable snapshot plus exclusions, not merely adjust arithmetic.

## 5. Conflicts / corrections to historical classifications

- Historical rows for #1316/#2311, #1814, #1338/#2246, #2237, #1962/#2238, #1827/#2233, #2121/#2271, #1893 and #1409/#1918 described pre-V12 gaps. They are superseded by merged behavior.
- PR #2238's implementation was **not** adopted; the issue is corrected by a different SARIF-only design. It is not commit-level attribution.
- PR #1918's broad `dir` hook was **not** adopted; the issue is corrected by a bounded explicit-file command/hook design.
- #1651/#2324 and #1243 remain active because an opt-in/partial correction is not a complete fix.
- No historical `fixes-qol-tracking.md` link remains in the assessment tree; canonical references now point to `RNC-MOD.md`.

## 6. V13 implementation plan

| Workstream | Packages | Initial approach |
|---|---|---|
| Reporting/QOL | V13-01..07 | Unused-ignore diagnostics, safe effective-config summary, multi-format architecture assessment, bounded context design, SARIF level policy, keyed grouping policy, console location/decoded UX. |
| Git/config reliability | V13-08..12 | Parsed log-opts vector and fail-closed validation, shallow-clone reporting, subdirectory/bare-repo scope validation, remote normalization, actionable config errors. |
| Hooks/performance/stability | V13-13..16 | Pre-commit environment reproducer, trustworthy profile script, decode/hang measurements, explicit-file scalability. |
| Release engineering | V13-17..20 | RNC GHCR-only workflow, version/install policy, release artifact portability, attestations/documentation. Plan only; no publication. |

Implementation evidence will be recorded in the V13 remediation matrix. CP0 does not itself change scanner source.

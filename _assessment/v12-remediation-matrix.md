# CP-V12 Remediation Matrix

**Baseline:** `197a7e5516b0029d0b7bb2554552fc9ed54efb68` (V11 merge on `rnc-mod`)
**Development branch:** `fix/v12-hardening`
**Scope:** 15 designated correctness, security, reporting, performance, portability, and hook packages.
**Evidence:** disposable synthetic repositories, files, reports, and hooks under the local temporary directory; no production repositories or real credentials.

## Status summary

| ID | Package | Status | Primary correction / decision |
|---|---|---|---|
| V12-01 | Merge-only Git history | POLICY DECISION REQUIRED | Default remains unchanged; explicit `--log-opts=--diff-merges=first-parent` is tested and documented. Default change would alter duplicates/fingerprints. |
| V12-02 | Git textconv authenticity | FIXED — LOCALLY VALIDATED | `--no-textconv` and `--no-ext-diff` pinned for log and diff commands. |
| V12-03 | Host Git formatting | FIXED — LOCALLY VALIDATED | `--pretty=medium`, `--date=default`, `--no-color` pinned; hostile formats tested. |
| V12-04 | Stream termination/completeness | ALREADY CORRECTED — REGRESSION VERIFIED | V11 drain/section/error collector retained; patch-suppressing/combined log formats now fail closed. |
| V12-05 | Source-local config trust | POLICY DECISION REQUIRED | Default auto-discovery preserved; additive `--no-local-config` trust control added. |
| V12-06 | Contextual redaction | FIXED — LOCALLY VALIDATED | Redaction groups known exposures by file/commit and covers empty-capture Match. |
| V12-07 | CSV formula injection | FIXED — LOCALLY VALIDATED | All string cells neutralized, including whitespace/control/zero-width/fullwidth triggers. |
| V12-08 | SARIF Unicode columns | FIXED — LOCALLY VALIDATED | SARIF-only codepoint columns; historical JSON/CSV column contract unchanged. |
| V12-09 | SARIF artifact URI | FIXED — LOCALLY VALIDATED | Single-escaped relative/path URIs; Windows drive/UNC handling. |
| V12-10 | Paired report publication | PARTIALLY CORRECTED — RESIDUAL CONFIRMED | Handled overwrite primary failure restores/removes companion; crash atomicity remains impossible. |
| V12-11 | Duplicate regex scan | FIXED — LOCALLY VALIDATED | Reuses the first match slice; benchmark improved. |
| V12-12 | Recursive decode resources | NOT REPRODUCED — EVIDENCE INSUFFICIENT | Finite probes show depth-linear allocation but no correctness defect; no default depth change. |
| V12-13 | CRLF baseline identity | FIXED — LOCALLY VALIDATED | Detector normalizes CRLF findings; legacy baseline entries migrated. |
| V12-14 | Long-line boundaries | FIXED — LOCALLY VALIDATED | 100-300 KB cases complete; over-cap lines fail closed rather than silently splitting. |
| V12-15 | Pre-commit explicit files | FIXED — LOCALLY VALIDATED | New `gitleaks files` plus explicit-file hook IDs; staged hooks unchanged. |

## Findings and decisions

### V12-01 — merge-only changes

**Confirmed:** default `git log -p` omits merge-only content. `--log-opts=--diff-merges=first-parent` exposes it. `--remerge-diff` and combined diffs produce unusable or duplicate output. Default adoption would change findings, duplicates and fingerprints. The established default remains; explicit opt-in is tested by `TestGitMergeOnlyFindingRequiresExplicitMergeDiff` and documented in README.

### V12-02 — textconv and ext-diff

**Confirmed:** a repository-local `diff.<driver>.textconv` command executed during a synthetic scan. The scanner now pins `--no-textconv` and `--no-ext-diff` for history and worktree diff commands. `TestGitLogDoesNotInvokeRepositoryTextconv` uses a missing helper and a synthetic binary; the scan succeeds without executing it.

### V12-03 — host Git presentation

**Confirmed:** `format.pretty` and `log.date=relative` destroyed commit identity/date parsing. Git log now pins medium/default output after user log options; `TestGitLogNeutralizesRepositoryFormatting` asserts SHA and UTC date. `--no-color` remains.

### V12-04 — parser completeness and unsupported output

V11's drain/section/error collector remains and existing malformed/truncated/cancellation tests pass. Silent-clean residuals for `--name-only`, `--name-status`, raw/numstat/shortstat, no-patch, combined and remerge formats are now rejected by `validateGitLogOpts`; `TestValidateGitLogOptsRejectsSilentCleanFormats` covers them. Ordinary patch and first-parent opt-in remain supported.

### V12-05 — local configuration trust

**Confirmed:** an automatic source-local `.gitleaks.toml` can replace rules or suppress everything. Default behavior is retained for compatibility. `--no-local-config` skips only auto-discovery; explicit `--config` and environment sources still apply. `TestCLINoLocalConfigUsesDefaultRules` verifies the boundary. Operator adoption of a secure default is still required if the project wants untrusted scans protected by default.

### V12-06 — contextual redaction

**Confirmed:** two synthetic cases could expose an adjacent credential: findings with different line ranges in one file, and empty-capture findings whose Match carried a neighbor. Redaction grouping now uses file+commit and includes non-path empty-capture Match values. Unit tests cover multiline ranges and empty captures; existing CLI redaction tests pass.

### V12-07 — CSV formula injection

All CSV string cells are neutralized with a leading apostrophe when the first significant rune is `=`, `+`, `-`, `@`, or a fullwidth variant. Leading whitespace/control/format characters, BOM and zero-width characters are skipped before classification. Tests cover idempotence, RFC 4180 round-trip, multiline cells and normal dates/paths.

### V12-08 — SARIF Unicode columns

V11's internal byte columns and line-boundary semantics are preserved for JSON/CSV/console compatibility. New `StartColumnCodepoint`/`EndColumnCodepoint` fields are populated from the actual V11 line start and consumed only by SARIF. Tests cover ASCII, Arabic, CJK, emoji, combining marks, later lines and CRLF. Existing location expectations remain stable.

### V12-09 — SARIF artifact URIs

Artifact paths are single-escaped via `net/url`; relative and POSIX absolute paths remain path references, Windows drives become `file:///C:/...`, UNC paths become `file://server/share/...`. Percent, hash, query, spaces and Unicode are covered. Existing SARIF fixtures remain unchanged for ordinary ASCII paths.

### V12-10 — paired reports

Overwrite-mode publication now backs up an existing companion (hardlink or copy), publishes the new companion, and restores the old one if the primary rename fails. If no companion existed, the new one is removed. Published identity is captured to avoid deleting a concurrent replacement. Crash atomicity across two files remains impossible and is documented; tests cover both restore and remove paths.

### V12-11 — duplicate regex scan

`detectRule` already computed `matches`; the loop now reuses it. Benchmark on a positive-heavy rule:

- Before: approximately 3.3-3.9 µs/op, 1,917-1,920 B/op, 14 allocs/op.
- After: approximately 3.0-3.3 µs/op, 1,894-1,896 B/op, 11 allocs/op.

A separate 8 MiB probe measured approximately 2.8-2.9 ms/op before and 0.8-1.0 ms/op after. Findings, locations, allowlists and decoded-pass semantics are unchanged.

### V12-12 — recursive decoding resources

Finite synthetic probe (6.67 MB encoded payload, empty rule set) measured 13.3 MB total allocation and ~20.7 MB heap after scan at depths 0, 1, 3 and 5; no depth-correlated leak was observed in that control. A nested-token probe showed depth-linear allocation (roughly 1 MB live per pass for a 0.4 MB payload), so amplification is real but was not established as a correctness defect in V12. No default decode-depth change and no unproven cap were introduced. This remains a policy/benchmark candidate for the independent harness.

### V12-13 — CRLF baseline identity

The real detector now trims CRLF from match spans and canonicalizes `Match`/`Secret` before entropy; LF and CRLF scans of the same content produce equal fields. Legacy baselines with a trailing bare CR get `EndColumn` corrected and entropy recomputed during `LoadBaseline`. Tests use real detector output and a legacy JSON baseline, not hand-built invariant values.

### V12-14 — long lines

Fragments now finish the current line under a 1 MiB extra-read budget after the normal 25 KB safe-boundary lookahead. Tests at 100 KB, 125 KB, 125 KB+1, 150 KB and 250 KB keep a boundary-straddling secret intact. A line that cannot be finished within the cap returns an operational error rather than silently splitting and reporting a clean scan.

### V12-15 — explicit-file hooks

`gitleaks files` scans exactly the named files, deduplicates paths, rejects missing files and directories, and never walks a parent directory. New `gitleaks-files` and `gitleaks-files-system` hooks enable `pass_filenames`; the original three hooks remain staged-index scanners. Go, CLI and real pre-commit `--files`/`--all-files` tests cover scope, secrets, clean runs, missing paths and directory rejection. `--all-files` processes files sequentially; this is a performance limitation, not a scope fallback.

## Validation gates

The final development and integration candidates are validated with:

```text
go build ./...
go test -count=1 ./...
go test -count=1 -tags gore2regex ./...
go vet ./...
go mod verify
go generate ./...
git diff --check
```

Plus focused tests, real pre-commit integration, and a non-root Linux race run from an exact-tree Git bundle. Results are recorded in the final checkpoint report and PR.

# CP-V14A Corrective Review

**Parent PR:** #7, head at review `f4327bcf0947ba7fa04da43d1eb523694b46babc`
**Parent baseline:** PR #6 / `audit/v13-integration@71be6f3` (open, not merged)
**Correction branch:** `feat/v14-policy-implementation` (no new PR, no retarget)
**Policy authority:** `feat/issues-and-pr:_backlog/operator-decisions/2026-11-14-v13-policy.md`

## Confirmed defects and dispositions

| Finding | Disposition | Root cause | Corrective implementation | Regression evidence | Residual |
|---|---|---|---|---|---|
| V14A-01 grouped report can overwrite its key | **Confirmed / corrected** | Key path loaded but never registered in `reportOutput.protectedInputs`; grouped output constructed with unconditional `reportOverwrite`; only raw primary/grouped equality checked. | Key loaded before report setup; registered as protected input on primary and grouped outputs; grouped uses resolved `mode`; key outside scan target; full destination-set collision validation; existing grouped sidecar excluded from scan. | `TestGroupedReportRejectsKeyAsDestination`, `...KeyAsPrimaryDestination`, `...KeyInsideScanTarget`, `...SymlinkKey`, `...HardlinkKeyDestination`, `...ObeysOutputPolicy`, aggregate/informative collision tests; key bytes unchanged after rejection. | Multi-file publication is not crash-atomic. Preflight is shared, but a late grouped failure after primary publication can leave a complete primary and no grouped sidecar; exits nonzero and reports the primary path. |
| V14A-02 grouped output ignores policy | **Confirmed / corrected** | Grouped output always used `reportOverwrite`. | Grouped output uses the same resolved reject/overwrite/aggregate mode as the primary; protected destination set includes aggregate-numbered candidates. | Existing grouped file under reject is preserved; `--overwrite` replaces; `--aggregate` writes numbered sidecar; primary/grouped aggregate collision rejected. | Aggregate collision scan is bounded to 1,000 candidate numbers per path; beyond that, operators should use distinct names. |
| V14A-03 cross-report collision/atomicity | **Confirmed / corrected (preflight), explicitly bounded** | Only direct primary/grouped and informative/grouped raw path equality was checked; key/config/ignore/baseline/template and numbered destinations were incomplete. | `validateReportDestinationSet` pairwise-checks primary, informative and grouped paths; `samePath` resolves parent symlinks and file identity; key/config/baseline/ignore/template/primary/companion are protected inputs; grouped destination excluded from scan. | Collision tests for primary, informative, hardlink key, aggregate numbering; all destination files absent/unchanged on rejection. | No cross-file rollback or crash-atomicity claim. Primary publication precedes grouped publication; partial failure is reported nonzero. Primary/informative V12/V13 paired rollback behavior is preserved. |
| V14A-04 inline strict extension gap | **Confirmed / corrected** | Inline TOML was validated only when `len(ConfigFiles)==0`; an inline config extending a file made `ConfigFiles` nonempty, skipping the inline root. | Inline effective source is validated whenever selected by precedence, independently of extension files; every `ConfigFiles` entry is still validated. | `TestCLIStrictConfigValidatesInlineRootWithExtension` rejects unknown inline root and accepts valid inline+extension. | If `GITLEAKS_CONFIG`/`--config` selects a file, ignored inline TOML is not validated; it is not the effective source. |
| V14A-05 key entropy/lifecycle | **Partially corrected by design control** | Length check cannot prove entropy; file identity race not previously checked. | Secure generation guidance; symlink/regular/identity/access checks; key outside scan target; no logging; key rotation/domain separation documented. | Generated-key workflow, symlink, hardlink, permissive-mode (POSIX), short/empty key tests. | No cryptographic randomness test is claimed. Windows ACLs are operator-controlled; key is held in memory for the scan and Go cannot guarantee zeroization. Hardlinked aliases inside a scan target are an operator-guidance residual. |

## Corrective architecture

- `cmd/root.go`: loads/validates the key before report construction; protects key input; uses resolved report mode for grouped output; validates complete destination set; tracks published primary path for partial-publication diagnostics.
- `cmd/report.go`: `protectInput` revalidates after adding protected inputs; `samePath` resolves parent-directory symlink aliases.
- `cmd/report_destinations.go`: primary/informative/grouped set validation, including aggregate candidates.
- `detect/detect.go`: excludes configured grouped output from source scanning; records published report path.
- `config/strict.go` unchanged schema; CLI flow fixed for inline roots.
- `KEY-MANAGEMENT.md`: generation/rotation/Windows ACL/limitations.

## Validation

| Check | Result |
|---|---|
| `go build ./...` | PASS |
| `go test -count=1 ./...` | PASS (final run after corrections) |
| `go test -count=1 -tags gore2regex ./...` | PASS |
| `go vet ./...` | PASS |
| `go mod verify` | PASS |
| Non-root Linux `go test -race -count=1 ./...` | PASS at final corrective head |
| Windows grouped/key/hardlink/symlink tests | PASS where supported; POSIX-permission test skipped on Windows with explicit platform reason |
| V13/V14 primary finding/fingerprint parity with opt-ins disabled | Covered by unchanged primary path and existing regression suites; no default scanner change |
| Grouped key unchanged after rejection | PASS |
| Existing grouped sidecar excluded | PASS |

GitHub Actions is **not triggered** for PR #7 while it targets stacked base `audit/v13-integration`; this is a documented limitation, not a pass claim. PR #6's CI remains the V13 gate.

## Architectural experiments

- `/v2` module experiment: isolated, not committed; all local build/test/RE2/vet/install/linker/cross-build checks pass. Public module-proxy/`go install` is unverified. Revised recommendation: **COA2 conditional** on operator approval and a `minVersion` compatibility floor; COA1 fallback.
- First-parent: expanded ordinary/evil/conflict/octopus/repeated/pathspec/ignore/large-history matrix; opt-in retained. Default change still requires duplicate/fingerprint/baseline migration and runtime evidence.

## Remaining gates

- Operator decision on COA1 vs COA2 and the `minVersion` compatibility policy.
- Operator decision on first-parent default (not recommended now).
- Operator decision on fingerprint v2 (design only).
- Key-management service/ACL acceptance.
- Normal GitHub Actions matrix after V13 merges and PR #7 is retargeted.

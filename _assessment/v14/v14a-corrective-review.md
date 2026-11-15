# CP-V14A Corrective Review

**Parent PR:** #7, head at review `f4327bcf0947ba7fa04da43d1eb523694b46babc`
**Parent baseline:** PR #6 / `audit/v13-integration@71be6f3` (open, not merged)
**Correction branch:** `feat/v14-policy-implementation` (no new PR, no retarget)
**Policy authority:** `feat/issues-and-pr:_backlog/operator-decisions/2026-11-14-v13-policy.md`

> **V14B superseding note:** V14B removed V14A-08's 1,000-candidate allocation ceiling and replaced the nested candidate enumeration with structural naming-family checks. The V14A text below is preserved as the historical corrective record; current behavior is documented in `v14b-report-architecture.md`.

## Confirmed defects and dispositions

| Finding | Disposition | Root cause | Corrective implementation | Regression evidence | Residual |
|---|---|---|---|---|---|
| V14A-01 grouped report can overwrite its key | **Confirmed / corrected** | Key path loaded but never registered in `reportOutput.protectedInputs`; grouped output constructed with unconditional `reportOverwrite`; only raw primary/grouped equality checked. | Key loaded before report setup; registered as protected input on primary and grouped outputs; grouped uses resolved `mode`; key outside scan target; full destination-set collision validation; only the exact configured grouped destination (or its filesystem alias) is excluded from scan. | `TestGroupedReportRejectsKeyAsDestination`, `...KeyAsPrimaryDestination`, `...KeyInsideScanTarget`, `...SymlinkKey`, `...HardlinkKeyDestination`, `...ObeysOutputPolicy`, aggregate/informative/parent-symlink collision tests; key bytes unchanged after rejection. | Multi-file publication is not crash-atomic. Preflight is shared, but a late grouped failure after primary publication can leave a complete primary and no grouped sidecar; exits nonzero and reports the primary path. |
| V14A-02 grouped output ignores policy | **Confirmed / corrected** | Grouped output always used `reportOverwrite`. | Grouped output uses the same resolved reject/overwrite/aggregate mode as the primary; protected destination set includes allocatable family members. | Existing grouped file under reject is preserved; `--overwrite` replaces; `--aggregate` writes numbered sidecar; primary/grouped aggregate collision rejected at `_001` and `_1000`; V14B extends this beyond `_1000`. | V14B removed the former fixed 1,000-candidate allocation bound; see `v14b-report-architecture.md`. |
| V14A-03 cross-report collision/atomicity | **Confirmed / corrected (preflight), explicitly bounded** | Only direct primary/grouped and informative/grouped raw path equality was checked; key/config/ignore/baseline/template and numbered destinations were incomplete. | `validateReportDestinationSet` pairwise-checks primary, informative and grouped paths; `samePath` resolves parent symlinks and file identity; key/config/baseline/ignore/template/primary/companion are protected inputs; only the exact grouped destination is excluded from scan. | Collision tests for primary, informative, hardlink key, parent symlink, aggregate numbering; all destination files absent/unchanged on rejection. | No cross-file rollback or crash-atomicity claim. Primary publication precedes grouped publication; partial failure is reported nonzero. Primary/informative V12/V13 paired rollback behavior is preserved. |
| V14A-04 inline strict extension gap | **Confirmed / corrected** | Inline TOML was validated only when `len(ConfigFiles)==0`; an inline config extending a file made `ConfigFiles` nonempty, skipping the inline root. | Inline effective source is validated whenever selected by precedence, independently of extension files; every `ConfigFiles` entry is still validated. | `TestCLIStrictConfigValidatesInlineRootWithExtension` rejects unknown inline root and accepts valid inline+extension. | If `GITLEAKS_CONFIG`/`--config` selects a file, ignored inline TOML is not validated; it is not the effective source. |
| V14A-05 key entropy/lifecycle | **Partially corrected by design control** | Length check cannot prove entropy; file identity race not previously checked. | Secure generation guidance; symlink/regular/identity/access checks; key outside scan target; no logging; key rotation/domain separation documented. | Generated-key workflow, symlink, hardlink, permissive-mode (POSIX), short/empty key tests; repeat/key-change determinism unit tests. | No cryptographic randomness test is claimed. Windows ACLs are operator-controlled; key is held in memory for the scan and Go cannot guarantee zeroization. Hardlinked aliases inside a scan target are an operator-guidance residual. |
| V14A-06 incomplete architecture evidence | **Corrected evaluation** | Prior module report and first-parent matrix were too narrow. | Isolated `/v2` migration experiment; expanded Git matrix including conflict, octopus, repeated, pathspec, ignore, baseline and 100-commit runtime; `minVersion` policy proposal. | Module experiment standard/RE2/vet/install/linker/cross-build; measured Git table. | No production module/default change. Public proxy install, older Git, and large-scale runtime remain unverified. Current `minVersion` behavior is a warning, not a rejection; the compatibility-floor proposal is separate and unimplemented. |
| V14A-07 overbroad grouped-report exclusion | **Confirmed / corrected** | Any same-directory name matching `<configured-stem>_<digits><ext>` was treated as a generated grouped sidecar, so ordinary source files could be silently skipped. | Exclusion now requires exact path or filesystem identity with the configured grouped destination. Similar names are ordinary scan inputs. | `TestGroupedReportScansLookalikeSourceFiles` covers `grouped_001.json`, `grouped_123.json`, `grouped_999.json`, `grouped_1000.json`, `grouped_123.txt`; exact configured path remains excluded. | Corrected. The scanner does not infer ownership of prior numbered files without a configured path or explicit tracking; no persistent generated-path registry is introduced. |
| V14A-08 aggregate validation/allocation mismatch | **Confirmed / corrected; superseded by V14B** | Preflight inspected 1,000 aggregate candidates, while allocation could continue indefinitely. A protected input or cross-report collision at `_1001` or later could not be validated. | V14A shared `maxAggregateReportCandidates=1000`. V14B removes the cap, checks naming families structurally and rechecks the selected path. | V14A `TestAggregateAllocationUsesValidationBoundary` is historical. Current evidence: `v14b-report-architecture.md`, high-index compatibility tests and 10,000-candidate benchmark. | V14B restores unbounded aggregate numbering without nested enumeration; no 1,000-candidate exhaustion is imposed. |
| V14A-09 `minVersion` wording | **Confirmed / corrected documentation** | The module assessment said an older RNC product version would "reject" otherwise compatible upstream configs. `validateMinVersion` actually warns and continues; only invalid SemVer/scanner-version parsing errors. | Assessment now distinguishes warning behavior from the proposed future compatibility-floor policy; no production change. | Source inspection of `config/config.go:validateMinVersion`; doc review. | Proposal remains unimplemented; release tag-base validation is a separate release-policy check. |

## Corrective architecture

- `cmd/root.go`: loads/validates the key before report construction; protects key input; uses resolved report mode for grouped output; validates complete destination set; tracks published primary path for partial-publication diagnostics.
- `cmd/report.go`: `protectInput` revalidates after adding protected inputs; `samePath` resolves parent-directory symlink aliases.
- `cmd/report_destinations.go`: primary/informative/grouped set validation, including aggregate candidates.
- `detect/detect.go`: excludes only the exact configured grouped output (or a filesystem alias); lookalike numeric filenames are scanned. Records the published report path.
- `config/strict.go`: schema unchanged; CLI flow fixed for inline roots.
- `KEY-MANAGEMENT.md`: generation/rotation/Windows ACL/limitations.

## Validation

Final corrective head is recorded after the validation commit. Actual commands/results:

| Check | Result |
|---|---|
| `go build ./...` | PASS |
| `go test -count=1 ./...` | PASS |
| `go test -count=1 -tags gore2regex ./...` | PASS |
| `go vet ./...` | PASS |
| `go mod verify` | PASS |
| `go generate ./...` | PASS; only generated `config/gitleaks.toml` line endings, restored |
| Non-root Linux `go test -race -count=1 ./...` | PASS at `ef55493`; PASS again at closeout code `ca606da` from a clean Linux clone |
| Windows grouped/key/hardlink/symlink/parent-alias tests | PASS where supported; POSIX-permission test skipped on Windows with explicit platform reason |
| Grouped preflight/render/publication failure tests | PASS |
| Lookalike source filenames (`grouped_001/123/999/1000.json`, `grouped_123.txt`) scanned | PASS |
| Exact configured grouped destination excluded | PASS |
| Aggregate `_1000` allocation, protected-input rejection and `_1001` exhaustion | PASS at V14A head; **superseded by V14B** (uncapped, structural validation) |
| Key unchanged after rejection | PASS |
| Default non-strict configuration behavior | PASS |
| V13/V14 primary finding/fingerprint parity with opt-ins disabled | Existing regression suites pass; no primary detection/fingerprint change |

Late grouped publication after primary publication and concurrent destination replacement are not injected in an end-to-end CLI test. The grouped writer uses the established `reportOutput` publication path; its render/temp/rename failure behavior is unit-tested. No crash-atomicity claim is made.

GitHub Actions is **not triggered** for PR #7 while it targets stacked base `audit/v13-integration`; this is a documented limitation, not a pass claim. PR #6 has 12/12 applicable checks on its V13A head.

## Architectural experiments

- `/v2` module experiment: isolated, not committed; local build/test/RE2/vet/install/linker/cross-build checks pass. Public module-proxy/`go install` is unverified. Revised recommendation: **COA2 conditional** on operator approval and a `minVersion` compatibility floor; COA1 fallback.
- First-parent: expanded ordinary/evil/conflict/octopus/repeated/pathspec/ignore/baseline/large-history matrix; opt-in retained. Default change still requires duplicate/fingerprint/baseline migration and runtime evidence.

## Remaining gates

- Operator decision on COA1 vs COA2 and the `minVersion` compatibility policy.
- Operator decision on first-parent default (not recommended now).
- Operator decision on fingerprint v2 (design only).
- Key-management service/ACL acceptance.
- Normal GitHub Actions matrix after V13 merges and PR #7 is retargeted.

# CP-V14A Final Closeout

**Status:** FINAL CORRECTIVE PACKAGE — pending operator acceptance.
**Repository:** `RNC-Solutions-Group/gitleaks`
**Branch / PR:** `feat/v14-policy-implementation` / PR #7
**Base:** PR #6, `audit/v13-integration@71be6f3fdf99ac711e7d1dbb5ee1683147ff2883` (open, not merged)
**Code baseline before this closeout:** `ef55493f7063ade3948826f620865544062f83e8`
**Final PR head:** recorded on PR #7 after the closeout commit; no new PR was opened.
**Constraints observed:** no retarget, rebase, force-push, merge, tag, release, image, package, Go module change, detection-rule change or scanner-harness change.

> **V14B superseding note (2026-11-15):** this document records the V14A state at `1f845e0`. V14B removed the V14A 1,000-candidate aggregate allocation ceiling and the quadratic candidate-enumeration validation described below. The V14A record is preserved as history; see `v14b-report-architecture.md` for current behavior.

## 1. Finding-by-finding disposition

| Finding | Disposition | Correction | Evidence |
|---|---|---|---|
| **A — overbroad report-file exclusion** | **Confirmed / corrected** | Removed the filename-pattern exclusion. Only the exact configured grouped destination, or a filesystem alias of it, is excluded from scanning. A source file named `grouped_001.json`, `grouped_123.json`, `grouped_999.json`, `grouped_1000.json` or `grouped_123.txt` is ordinary input. | `detect/detect.go`, `detect/utils.go`; `TestGroupedReportScansLookalikeSourceFiles` (five positive cases), `TestGroupedReportExcludesExactConfiguredDestinationFromScan` (negative exact-match case). |
| **B — aggregate destination boundary** | **Superseded by V14B** | V14A temporarily shared a 1,000-candidate cap between validation and allocation. V14B removed that cap: allocation again supports `_1001` and higher, while preflight checks naming families structurally and rechecks the selected path. | Historical V14A tests are preserved; current evidence is in `v14b-report-architecture.md`. |
| **C — `minVersion` characterization** | **Confirmed / corrected documentation** | Current `validateMinVersion` warns when the scanner is older; it does not reject the configuration solely from that comparison. The module assessment now separates this behavior from the proposed future compatibility floor. | `config/config.go:validateMinVersion`; `_assessment/v14/module-identity-comparison.md`. No production code change. |
| V14A-01 grouped key/destination | Corrected (prior V14A) | Key is a protected input; grouped output obeys resolved policy; destination set is validated. | `cmd/root.go`, `cmd/report.go`, `cmd/report_destinations.go`; grouped security tests. |
| V14A-02 grouped policy | Corrected (prior V14A) | Reject/overwrite/aggregate semantics apply to the sidecar. | `TestGroupedReportObeysOutputPolicy`. |
| V14A-03 cross-report collisions | Corrected at preflight (prior V14A) | Primary, informative and grouped paths checked as one set, including parent symlinks and file identity. | Collision tests; not crash-atomic across files. |
| V14A-04 inline strict root | Corrected (prior V14A) | Inline root is validated even when an extension populates `ConfigFiles`. | `TestCLIStrictConfigValidatesInlineRootWithExtension`. |
| V14A-05 key lifecycle | Partially corrected by design control | CSPRNG/ACL/rotation guidance; regular/non-symlink/identity checks; no entropy claim. | `KEY-MANAGEMENT.md`; key tests. Windows ACLs and memory zeroization remain operator-controlled. |
| V14A-06 architecture evidence | Corrected evaluation only | Isolated `/v2` and expanded first-parent matrix; no production module/default change. | `module-identity-comparison.md`, `first-parent-default-evaluation.md`. |

## 2. Validation

Executed on the closeout worktree at the final code state, Go `1.27.3 windows/amd64` unless noted:

| Check | Result |
|---|---|
| `go build ./...` | **PASS** |
| `go test -count=1 ./...` | **PASS** (all packages; `cmd` 389.822s) |
| `go test -count=1 -tags gore2regex ./...` | **PASS** (all packages; `cmd` 381.769s) |
| `go vet ./...` | **PASS** |
| `go mod verify` | **PASS** |
| `go generate ./...` | **PASS**; `config/gitleaks.toml` line endings changed, restored; no other generated tracked diff |
| `git diff --check` | **PASS** |
| Focused V14A tests | **PASS** at the V14A head: lookalike filenames, exact destination exclusion, `_1000` allocation, protected-input boundary, `_1000` collision, exhaustion. The cap/exhaustion portion is superseded by V14B. |
| Focused strict-config tests | **PASS**, including inline root plus extension and malformed/misplaced TOML |
| Non-root Linux `go test -race -count=1 ./...` | **PASS** — `golang:1.27.4`, UID/GID 10001, clean Linux clone at closeout head `cb7e2c0` (code commit `ca606da`); symlink fixture materialized; all packages |
| Windows-specific hardlink/symlink/parent-alias tests | **PASS** where supported; POSIX permission test skipped on Windows with an explicit platform reason |

**Not run / not claimed:**

- GitHub Actions for PR #7: not triggered while the PR targets stacked base `audit/v13-integration`. PR #6 has 12/12 applicable checks. PR #7 must be retargeted after V13 lands.
- Docker build/smoke for V14A: not run. V13 CI Docker results are not a V14 pass.
- Independent scanner harness/accuracy certification: not run; outside this assignment.
- Published-release/SBOM/signature verification: nothing published.

## 3. Documentation review

Reconciled:

- `README.md` — current V14 candidate status, CI limitation, opt-in controls.
- `RELEASING.md` — removed the contradictory moving-tag claim; only immutable full-version container tags; release-tag base distinguished from runtime `minVersion`; no V14 publication claim.
- `RNC-MOD.md` — V14 document status and PR #7 metadata; exact-destination scan rule; V14B aggregate-family validation; grouped/strict limitations updated from "policy decision" to implemented opt-in; validation claims separated into local, V13 CI and not-run; corrected V14A IDs and artifacts; three H1 sections and attribution register preserved.
- `KEY-MANAGEMENT.md` — exact destination distinction and V14B uncapped allocation/selected-path check.
- `_assessment/v14/module-identity-comparison.md` — warning-versus-reject correction.
- `_assessment/v13-remediation-matrix.md`, `_assessment/v13-reporting-options.md` — V13 historical statuses now point to the V14 opt-in implementation where applicable.
- `_assessment/contributor-provenance/v13-contributor-manifest.md` — V13 policy-gated wording updated without changing authorship.
- `_assessment/v14/v14a-corrective-review.md`, `implementation-and-evaluation-report.md`, `operator-policy-matrix.md` — V14A-07/08/09 and final evidence recorded, with V14B supersession noted.
- `_assessment/v14/v14b-report-architecture.md` — V14B compatibility/collision redesign, benchmark and security evidence.
- `_assessment/v14/backlog-reconciliation-proposal.md` — post-merge proposal described below.

Attribution structure remains three H1 sections in `RNC-MOD.md`; §3 retains all original upstream contributor credit. No contributor commit was adopted or rewritten in V14.

## 4. Proposed post-merge backlog delta

Canonical baseline: 285 issues / 252 PRs / 388 slots; 186 active issues / 144 active PRs / 226 active slots; 99 / 108 / 162 dispositioned. Full detail, exact sources, tests, commits and relationships are in `backlog-reconciliation-proposal.md`.

| Measure | Canonical | Proposed | After merge |
|---|---:|---:|---:|
| Active issues | 186 | -4 | 182 |
| Active PRs | 144 | -2 | 142 |
| Active slots | 226 | -4 | 222 |
| Dispositioned issues | 99 | +4 | 103 |
| Dispositioned PRs | 108 | +2 | 110 |
| Dispositioned slots | 162 | +4 | 166 |
| Totals | 285 / 252 / 388 | 0 | 285 / 252 / 388 |

Proposed full issue records: #1858 (SARIF level), #1714 (file links), #811 (strict unknown fields), #827 (malformed/strict diagnostics). Proposed PR-only records: #1900 (profile script), #1809 (Azure SSH). #811/#827 share the strict package with #1279, so they do **not** reduce the package-slot count; four standalone packages account for the four slot retirements.

Kept active with residuals: #2130/#2244 (no fail-on-unused), #983/#1654 (summary, not full config dump), #1054/#1162 (sidecar, not primary dedup; unkeyed hash rejected), #1279, #2003/#2128, #1310, #1163, #1293, #1754/#764. Deferred/operator-dependent: #1048, #1766, #2019, #2017, #1217/#2220/#885, #1028/#2301, #1839/#1209/#1226, #2239/#2286/#1566.

No canonical category, exclusion ledger or historical export was modified. The proposal becomes actionable only after both PRs merge; the operator must re-read the then-current backlog.

## 5. Remaining operator decisions and limitations

1. Merge PR #6 first, then retarget PR #7 to `rnc-mod` and run the full CI matrix.
2. Accept or reject the Windows ACL/key-lifecycle, memory-zeroization and hardlink-alias residuals.
3. Select COA1 or COA2 for module identity; no `/v2` migration is implemented.
4. Retain the first-parent default, or approve duplicate/fingerprint/baseline migration.
5. Keep fingerprint v2 design-only or authorize a separate implementation.
6. Decide whether `minVersion` should gain a compatibility-floor policy; current behavior warns.
7. Multi-file publication (primary/informative/grouped) is not crash-atomic. **V14B supersedes the former 1,000-candidate cap:** aggregate allocation is unbounded in numbering, with linear reservation attempts and structural preflight checks.
8. No GitHub Actions result exists for PR #7 until it is retargeted; no Docker, harness, release or attestation result is claimed.

## 6. Recommendation

**CONDITIONAL GO for PR #7**, contingent on:

- V13/PR #6 landing or otherwise being reconciled first;
- PR #7 being retargeted to `rnc-mod` and passing the full configured CI matrix;
- operator acceptance of the key-management/ACL, COA2, first-parent and bounded-aggregate decisions;
- backlog application only after both PRs are merged and re-verified.

No merge, release, tag, image, package or attestation was produced.

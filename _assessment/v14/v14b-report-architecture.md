# CP-V14B — Aggregate Report Compatibility and Collision Algorithm

**Status:** CORRECTIVE IMPLEMENTATION ON PR #7 — pending operator review.
**Starting head:** `1f845e043b01862c7fef2b17a168e2be4af05092`.
**Stacked base:** PR #6, `audit/v13-integration@71be6f3` (open, unmerged).
**Scope:** aggregate allocation compatibility and destination-collision validation. No new PR, merge, module, detection-rule, default-history or release change.

## 1. Finding V14B-01 — arbitrary aggregate allocation ceiling

### Root cause

V14A-08 introduced `maxAggregateReportCandidates = 1000` in `cmd/report_destinations.go` and reused it in `publishAggregate`. The constant was introduced to make the exhaustive preflight match the allocator, not because of a threat model or a V13 compatibility requirement. It applied to every aggregate report, including ordinary primary reports with no V14 grouped reporting.

### Reproduction

The V14A tests filled the requested path and `_001` through `_1000`, then expected publication to fail. That was a regression from V13, where `publishAggregate` had no numbered candidate cap. A primary report-only aggregate invocation could previously allocate `report_1001.json`; V14A returned a publication error.

### Correction

- Removed the constant and all references.
- `publishAggregate` again searches the requested path and numbered alternatives without an arbitrary ceiling.
- Existing candidates are skipped through exclusive `O_EXCL` reservation; a candidate selected for publication is rechecked against protected inputs/peers before the temporary report is renamed.
- Preflight no longer depends on enumerating the allocation range, so removing the ceiling does not make preflight unbounded.

### Compatibility result

Ordinary aggregation can now allocate beyond `_1000`, including `_1001` and substantially higher indices. Existing destination preservation, exclusive allocation, reject/overwrite semantics and nonzero publication failures are unchanged. No intentional V13 behavior change remains in this path.

## 2. Finding V14B-02 — quadratic collision validation

### Root cause

The V14A preflight combined:

- `validateReportDestinationSet()` pairwise checks;
- `aggregateDestinationsOverlap()`, which looped over 1,000 possible numbers; and
- a special primary/companion/grouped loop that, for each of 1,000 primary candidates, called `aggregateDestinationsOverlap()` again.

The all-three case therefore performed roughly O(1,000²) comparisons. Each comparison could invoke `samePath`, which canonicalizes parent paths and calls `os.Stat`, producing millions of filesystem metadata operations before scanning.

`checkAggregateProtectedInputs()` independently looped over 1,000 candidates for every protected input, even when the candidate paths were unrelated.

### Correction

Preflight now compares **naming families structurally** rather than enumerating candidate numbers:

- Exact configured aliases are still checked with `samePath`, including hardlink identity.
- Aggregate family overlap is determined from the canonical directory and the base/stem/extension naming rule produced by `numberedReportPath`.
- The primary informative companion family is derived from the selected primary name; a grouped destination is rejected if it equals the companion base or a numbered companion candidate.
- Protected inputs are rejected if they lie in the primary aggregate family or the informative companion family. This is O(protected inputs), not O(1000 × protected inputs).
- The allocator validates the **actual selected destination** after exclusive reservation and before publication. This is the dynamic check for a protected input, filesystem alias or peer destination that was not predictable at preflight.
- Paired primary/informative aggregate allocation retains complete-pair reservation but no longer revalidates the full static set for every occupied candidate.

No general-purpose path abstraction was introduced. The helpers are local to `cmd/report_destinations.go` and mirror the existing `numberedReportPath` naming rule.

### Complexity

Let `D` be the three configured report destinations, `P` the number of protected inputs, and `N` the number of occupied numbered candidates.

- Preflight naming checks: **O(D² + P)**, with D fixed at 3.
- Aggregate allocation: **O(N)** exclusive reservation attempts; each attempt is O(1), plus one O(P) selected-destination check before publication.
- Paired aggregate allocation: **O(N)** pair reservation attempts; one selected-pair safety check before publication.
- No nested candidate loops or O(N²) validation.

The allocator is intentionally not capped. A directory containing N occupied candidates costs N reservation attempts; this is the inherent cost of finding the first free name and is not quadratic validation.

## 3. Measured benchmark

Environment: Go `1.27.3 windows/amd64`, Intel i7-11375H, NTFS, Windows host, standard test configuration. Absolute values are not portable. Benchmark source: `cmd/report_benchmark_test.go`.

Commands:

```text
go test ./cmd -run '^$' -bench 'BenchmarkReportDestinationValidation|BenchmarkAggregateAllocation' -benchtime=1x -count=5
```

V14A comparison used the same benchmark file against the exported `1f845e0` tree, with `-count=3` for the bounded cases. The all-three V14A case was terminated after more than 90 seconds without a result and is not reported as a measurement.

### Preflight validation, median

| Case | V14A | V14B | Observation |
|---|---:|---:|---|
| Primary only | ~0.3–0.4 µs | ~0.3–0.4 µs | No aggregate naming comparison needed. |
| Primary + grouped | **3.180 s**, ~22.3 MB/op, 286,843 allocs/op | **5.51 ms**, ~31.4 KB/op, 411 allocs/op | ~577× lower measured median; allocations fall by ~700×. |
| Primary + companion + grouped | **>90 s, no result** | **12.18 ms**, ~75.7 KB/op, 968 allocs/op | Nested V14A path was not practically completable. |

### Aggregate allocation, median including report write

| Occupied numbered reports | V14A | V14B | V14B allocations |
|---:|---:|---:|---:|
| 0 | 5.788 s | 31.6 ms | 1,520 allocs/op, ~102 KB/op |
| 10 | 6.447 s | 31.1 ms | 1,688 allocs/op, ~120 KB/op |
| 100 | 4.750 s | 43.8 ms | 2,964 allocs/op, ~252 KB/op |
| 1,000 | **publication failure** | 151.1 ms | 16,313 allocs/op, ~1.54 MB/op |
| 10,000 | not run (ceiling) | 1.304 s | 151,369 allocs/op, ~14.6 MB/op |

The V14A times are dominated by its fixed 1,000-candidate validation and do not improve at low occupancy. The V14B allocation grows approximately linearly with occupied candidates, consistent with one exclusive reservation attempt per occupied name. The 10,000-entry case is a synthetic directory of hardlinked placeholders; it is not a claim about scan performance or typical report directories.

## 4. Security and compatibility preservation

| V14A guarantee | V14B result |
|---|---|
| Grouping-key protection | Key remains a protected input for primary and grouped outputs; key bytes preserved on rejection. |
| Config, baseline, ignore, template protection | Protected inputs are checked structurally at preflight and again on the selected destination. |
| Primary/informative/grouped collision detection | Exact aliases plus structural aggregate families and companion families are checked. |
| Symlink, hardlink, parent-directory aliases | `samePath` identity/canonical-parent checks remain; selected destination is rechecked after reservation. |
| Non-destructive rejection | Validation remains pre-scan; publication uses temp render and reservation. |
| Reject/overwrite/aggregate | Unchanged; aggregate allocation is no longer capped. |
| Report-destination exclusion | Only the exact configured grouped destination or alias is excluded; lookalike source files remain scanned. |
| Findings/fingerprints/redaction/exit behavior | No detection changes; existing focused and full suites pass. |
| Fail-closed publication | Reservation/render/rename failures remain nonzero; no successful partial result is reported. |
| Partial multi-file publication | Still explicitly non-atomic; the primary publication path is reported on a later grouped failure. |

The static-versus-dynamic split is deliberate:

- **Static hazards** — exact paths, predictable aggregate family overlap, companion-family overlap, protected input in a family — are rejected before scanning.
- **Dynamic hazards** — a selected numbered path that is now a protected input, alias or peer — are checked after reservation and before publication.
- **Allocation races** — exclusive create decides the candidate; a competing creator cannot claim the same name.

## 5. Tests added or updated

- `cmd/report_destinations_test.go`: aggregate family overlap and companion-family tables, including `_1000`, `_1001`, high indices, noncanonical padding/width and disjoint naming.
- `cmd/report_test.go`: ordinary primary aggregate allocation through `_1001` and `_10001`; existing `_001`, `_999`, `_1000` and base preserved; protected high candidate rejected; destination created after preflight skipped.
- `cmd/report_pair_test.go`: informative pair allocation through `_1002` after 1,001 occupied pairs.
- `cmd/grouped_report_security_test.go`: grouped allocation through `_1002`; primary/grouped collision at `_10001`; key at `_10001`; companion collision at `_001`/`_1001`; config protected at `_10001`.
- `cmd/report_benchmark_test.go`: isolated preflight/allocation benchmark at 0/10/100/1,000/10,000 occupied destinations.

## 6. Validation

Code/test commit: `abf91b307d089ea36fc5fb71d1625a7de47f1fa7`.

| Check | Result |
|---|---|
| `go build ./...` | **PASS** |
| `go test -count=1 ./...` | **PASS**; `cmd` 99.403s (down from 389.822s at V14A) |
| `go test -count=1 -tags gore2regex ./...` | **PASS**; `cmd` 98.274s |
| `go vet ./...` | **PASS** |
| `go mod verify` | **PASS** |
| `go generate ./...` | **PASS**; generated config line endings restored |
| `git diff --check` | **PASS** |
| Focused V14B compatibility/collision tests | **PASS**: `_1001`/`_10001`, grouped `_1002`, paired `_1002`, family/companion tables, high protected inputs |
| V14A lookalike/source-exclusion regressions | **PASS** |
| Non-root Linux `go test -race -count=1 ./...` | **PASS** at code commit `abf91b3`; `golang:1.27.4`, UID/GID 10001, clean clone, symlink fixture materialized |
| Windows filesystem tests | **PASS** where supported; POSIX permission test skipped on Windows |

No new dependency, module-path, detection-rule, default-history, release or harness change was made. Full GitHub Actions still requires PR #7 to be retargeted after PR #6 lands.

## 7. Residual limitations

- Allocation is unbounded in numbering, but a directory with an extremely large number of occupied names will take linearly longer and perform one metadata reservation per occupied candidate. No allocation-time resource cap is imposed; an operator can choose a different destination.
- Multi-file publication remains non-crash-atomic.
- The structural naming comparison mirrors the scanner's generated spelling; a manually named `report_0001.json` is not treated as a candidate of `report.json`.
- The selected-destination check is not a transaction against a hostile process that replaces the reserved path between final validation and rename; that existing TOCTOU residual is unchanged.
- Benchmarks are host-specific and do not measure scanner startup, regex matching or detection accuracy.

## 8. Recommendation

**CONDITIONAL GO for V14B corrective review**, contingent on full configured CI after PR #7 is retargeted and operator acceptance of the remaining V14A key/ACL, module, first-parent and fingerprint decisions. The V14B code is locally validated; PR #7 remains stacked on PR #6, and no merge or release is performed.

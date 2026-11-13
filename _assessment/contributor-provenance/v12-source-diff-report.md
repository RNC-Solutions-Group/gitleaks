# CP-V12 Source-Diff and Verification Report

**V10 authority:** `48c76717ccbbe2159185edb4510fb1924d7dbcc6` (tree `39b13207713d7102d7ac813a0599b2eaabf5d29f`)
**V12 development HEAD:** `1ff69b7836b32e4f738883a09bc2a3daeccbd69b`
**V12 integration product:** `9092703a25f4c10edc162699955a6993ec66ce37` (tree `a348eea61670da139b89b410efa5156325091e87`)

## Product equality

`git diff --exit-code fix/v12-hardening 9092703` returned zero for the entire tree before the documentation-only addendum. The integration product and the validated development HEAD are the same object:

```text
V12 development tree: a348eea61670da139b89b410efa5156325091e87
integration tree:     a348eea61670da139b89b410efa5156325091e87
```

The V12 provenance manifest and this report are the only additional files on the final audit tip.

## Difference from V10

Including the V12 remediation matrix, the integrated tree changes 40 files relative to V10:

```text
40 files changed, 1873 insertions(+), 180 deletions(-)
```

The changes are limited to the three approved V11 line/allowlist corrections already in the baseline, the V12 remediation areas below, their regression tests, the public documentation, and the remediation/provenance reports. No provider rules or benchmark harness were changed.

### V12-01 — merge-only Git changes

Default `git log -p` remains merge-diff-free. The tested, documented opt-in is `--log-opts=--diff-merges=first-parent`. Patch-suppressing and combined/remerge formats are rejected rather than reported clean. This is a deliberate policy decision because a default change alters findings, duplicates, fingerprints and baselines.

### V12-02/03 — Git authenticity and deterministic metadata

History and worktree commands pin:

```text
--no-textconv
--no-ext-diff
--no-color
--date=default
--pretty=medium
```

Synthetic textconv execution and hostile `format.pretty`/`log.date` configurations are covered. The flags are forced after user log options; users cannot re-enable repository-controlled textconv through `--log-opts`.

### V12-04 — stream and parser completeness

V10/V11's drain, raw-section tracker, error collector, cancellation wait and partial-scan report protection remain. V12 adds fail-closed rejection of `--name-only`, `--name-status`, raw/numstat/shortstat, no-patch, combined and remerge log formats.

### V12-05 — local configuration trust

Default automatic `(source)/.gitleaks.toml` discovery is unchanged for compatibility. `--no-local-config` is additive and skips only automatic discovery; `--config` and environment configuration retain precedence. Default-policy adoption remains an operator decision.

### V12-06 — redaction context

Known redaction exposures are grouped by file and commit instead of exact line range; empty-capture Match text is included unless it is a filename-only marker. Real multiline/adjacent and empty-capture tests pass.

### V12-07/08/09 — reporting

- CSV: all string cells are formula-neutralized, including leading whitespace/control/format characters, BOM, zero-width characters and fullwidth variants. RFC 4180 round-trip and idempotence are tested.
- SARIF columns: SARIF-only codepoint fields are populated from the V11 line start. JSON/CSV/console byte columns and baselines remain compatible.
- SARIF URIs: relative and POSIX paths are percent-escaped once; Windows drive paths become `file:///C:/...`; UNC paths become `file://server/share/...`. Percent/hash/query/space/Unicode cases are covered.

### V12-10 — paired reports

Overwrite mode backs up an existing informative companion, restores it if the primary rename fails, and removes only the newly published companion when none existed. The published file identity is captured for cleanup safety. Crash-window/atomicity limitations remain explicit and are not misrepresented as transactional.

### V12-11 — duplicate regex work

`detectRule` reuses the already-computed `matches` slice.

Measured on a positive-heavy rule:

```text
before: ~3.3-3.9 us/op, ~1918 B/op, 14 allocs/op
after:  ~3.0-3.3 us/op, ~1895 B/op, 11 allocs/op
```

An 8 MiB probe measured ~2.8-2.9 ms/op before and ~0.8-1.0 ms/op after. Behavior is unchanged.

### V12-12 — recursive decode resources

The finite 6.67 MB encoded-payload probe did not show depth-correlated memory growth in the empty-rule control; a nested-token probe showed depth-linear allocations. No correctness defect or leak was established, so no default decode-depth change or unproven segment cap was introduced. This remains a benchmark candidate for the independent harness.

### V12-13 — CRLF baseline identity

Real detector output now trims CRLF from match spans and canonicalizes Match/Secret before entropy. Legacy baselines are migrated in `LoadBaseline`: trailing bare CR is removed, EndColumn is adjusted, and entropy is recomputed. Real LF-to-CRLF baseline suppression is tested.

### V12-14 — long lines

The fragment reader finishes the current line under a 1 MiB extra-read budget. Boundary cases at 100 KB, 125 KB, 125 KB+1, 150 KB and 250 KB retain a complete boundary-straddling secret. Lines that cannot be finished within the cap return an operational error instead of silently splitting and reporting clean.

### V12-15 — pre-commit explicit files

`gitleaks files` scans exactly the named files, deduplicates paths, rejects missing paths and directories, and never walks a parent tree. New `gitleaks-files`/`gitleaks-files-system` hooks enable `pass_filenames`; the three staged hooks and their semantics are unchanged. Real pre-commit `--files` and `--all-files` behavior is tested.

## Validation

| Gate | Result |
|---|---|
| `go build ./...` | PASS |
| `go test -count=1 ./...` | PASS |
| `go test -count=1 -tags gore2regex ./...` | PASS |
| `go vet ./...` | PASS |
| `go mod verify` | PASS |
| `go generate ./...` + clean tree | PASS |
| `git diff --check` | PASS |
| Real pre-commit hook integration | PASS |
| Non-root Linux `go test -race -count=1 ./...` | PASS |

## Residual limitations

- Merge-only default coverage remains opt-in by policy.
- Automatic local configuration remains trusted by default; `--no-local-config` is the additive control.
- Paired report publication cannot be crash-atomic across two files.
- Lines above the 1 MiB finish budget fail closed rather than being partially scanned.
- `gitleaks files` scans sequentially; large `--all-files` sets are slower than an optimized directory walk.

## Merge recommendation

Do not squash the contributor ancestry. Use a merge commit to retain the five exact integration commit SHAs and trailers. Rebase and merge is acceptable only if linear history is required and the resulting SHA rewrite is acceptable; verify authors/trailers afterward. Squash and merge is prohibited for this PR.

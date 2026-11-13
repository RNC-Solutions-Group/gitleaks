# CP-V11 Source-Diff and Verification Report

**Baseline:** `V10-backup` at `48c76717ccbbe2159185edb4510fb1924d7dbcc6` (tree `39b13207713d7102d7ac813a0599b2eaabf5d29f`)
**V11 development HEAD:** `1b06b86a2ad815917b8e3f6265afb304f0eedd7e` (tree `73cfdd5c02943d97a731e0c7012380f9c4d00958`)
**Integration product commit:** `45f6155a56b2c8c8fad1272cf1f38f5602da8e45` (tree `73cfdd5c02943d97a731e0c7012380f9c4d00958`)
**Integration base:** `rnc-mod` at `ee1e59e38231b84a8277131827a09fc929383ecc`

## Product-tree equality

`git diff --exit-code fix/security 45f6155` returned zero before the documentation-only commit. The integration product tree and the validated V11 development tree are the same object:

```text
V11 development tree:   73cfdd5c02943d97a731e0c7012380f9c4d00958
integration tree:       73cfdd5c02943d97a731e0c7012380f9c4d00958
```

This documentation commit is the only approved difference on the report-bearing integration tip. The product source, tests, build configuration, and CI configuration remain identical to `fix/security`.

## Differences from V10

Excluding `_assessment/contributor-provenance/**`, V11 changes 14 files:

```text
 .github/workflows/test.yml                         |  14 +-
 cmd/generate/config/base/config.go                 |   2 +-
 cmd/generate/config/base/config_test.go            |  11 +
 config/config_test.go                              |  27 +++
 config/gitleaks.toml                               |   2 +-
 detect/codec/segment.go                            |   8 +-
 detect/codec/segment_test.go                       |  26 ++
 detect/detect_test.go                              | 270 ++++++++++++++-------
 detect/git_commit_allowlist_test.go                |  83 +++++++
 detect/location.go                                 |  12 +-
 detect/location_test.go                            |   4 +-
 sources/git.go                                     |  21 +-
 sources/git_test.go                                |  64 +++++
 testdata/config/valid/allowlist_global_commit.toml |  10 +
 14 files changed, 453 insertions(+), 101 deletions(-)
```

The six V10 contributor-provenance reports are inherited unchanged from `rnc-mod`. The two V11 provenance reports are additional documentation-only paths.

## Intentional behavior corrections

These differences are the purpose of V11 and must not be reverted to V10 equivalence.

### V11-A: line-boundary correctness

- `Finding.Line` and decoded `CurrentLine` now start at the matched line, not at the newline ending the previous line.
- Anchored `regexTarget="line"` allowlists work on lines after the first.
- Empty decoded ranges and empty/newline matches are clamped instead of panicking.
- Line and column numbers are unchanged, and `Finding.Line` remains excluded from JSON reports.

### V11-B: anchored placeholder allowlist

- The global allowlist regex is now `(?i)^(?:true|false|null)$`.
- Values such as `trueSuffix`, `prefixnull`, and `false_key` are no longer globally suppressed.
- `config/gitleaks.toml` was regenerated from the generator source; no unrelated generated changes appeared.
- A default-config regression asserts exact placeholders are allowed and placeholder-containing values are not.

### V11-C: global commit allowlist

- A matching global `[[allowlists]]` commit now skips the entire diff file and therefore the commit.
- The decision is centralized in `isCommitAllowed`, which stops at the first match.
- Staged scans with an empty commit SHA remain unaffected.
- A Windows-runnable temporary-repository test proves the allowed commit is suppressed while an unaffected commit still produces a finding.

### CI trigger correction

`.github/workflows/test.yml` now runs only for pull requests targeting `rnc-mod` and pushes to `rnc-mod`, matching the protected-branch policy. Job definitions and the test matrix are unchanged.

## V10 compatibility

All non-V11 behavior is inherited from the V10 product tree. The full standard and RE2 test suites, V10 report lifecycle tests, Git lifecycle tests, staged-index tests, archive tests, root-relative tests, and Windows-specific tests pass. No V10 correction was restored or removed outside the three approved V11 areas.

## Local validation

Executed on the integration candidate and the V11 development tree:

| Check | Result |
|---|---|
| `go build ./...` | PASS |
| `go test -count=1 ./...` | PASS |
| `go test -count=1 -tags gore2regex ./...` | PASS |
| `go vet ./...` | PASS |
| `go mod verify` | PASS |
| `go generate ./...` then clean tracked tree | PASS |
| `gofmt` on changed files | PASS |
| `git diff --check` | PASS |
| `go test -race -count=1 ./...` in non-root Linux `golang:1.27.4` | PASS |

The race run was executed at V11 development HEAD `1b06b86`, whose tree is the same `73cfdd5` object as the integration product. The integration candidate additionally passed the full Windows local suite. GitHub Actions is expected only after the final PR is opened.

## Contributor and history checks

- Six upstream source commits retained with original authors and source trailers.
- One RNC reconciliation commit on top.
- RNC V11 checkpoint commits are not present in the integration history.
- No unrelated V10 history was rewritten.
- The only cherry-pick conflict was the test-import/test-body conflict in `sources/git_test.go`; it was resolved without dropping either V10 tests or the upstream helper test.

## Merge recommendation

Preserve the contributor commits. Use a **merge commit** for the final PR if the operator wants the exact pre-merge contributor SHAs retained; it adds a merge commit but does not rewrite authors or source trailers.

If linear history is mandatory, **Rebase and merge** is the only acceptable alternative: it preserves author identity and trailers but rewrites the integration commit SHAs and must be followed by an author/trailer verification. Do not use **Squash and merge**, which would collapse the six contributor commits into the RNC reconciliation.

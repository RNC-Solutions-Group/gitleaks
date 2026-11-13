# Source-Diff Analysis: Integrated Upstream Branch vs CP-V10

**Integrated tree:** `audit/rnc-mod-upstream-integration` (pre-documentation HEAD
`100ae90`)  
**Reference:** CP-V10 `bd42a958f7476e9a6d8e17bd2dbb2abe3ff060ea`

## Comparison context

This is not a branch-versus-branch regression comparison. The integration branch
is based on frozen upstream `5f6ec56` and contains 18 selected upstream commits
plus two RNC integration-adjustment commits. It deliberately does **not** contain
the RNC V1–V10 implementation, informative reporting, build contract, CI, or
Docker work. The large deletion side of the diff is therefore mostly V10-only
functionality absent from the pre-V10 foundation.

Measured tree comparison:

```text
123 files changed, 2704 insertions(+), 14122 deletions(-)
```

This count is not a quality metric. The meaningful comparison is the
function-level overlap below.

## Substantive differences by area

### 1. Finding redaction

| Item | Detail |
| --- | --- |
| Upstream | #2261, `report/finding.go`, `report/finding_test.go` |
| Integration | `maskSecret` is the upstream rune implementation, `report/finding.go:88` |
| CP-V10 | `report/finding.go:231–254`; adds `percent == 0`, lower/upper clamps, `redactEmptySecret`, `uniqueSecrets`, `redactSecrets`, `maskedSecret`, `truncateRunes` |
| Behavioral difference | V10 protects empty captures, adjacent known secrets, required-rule output and more redaction surfaces; integration only fixes multibyte masking |
| Risk | Applying V10 on top will conflict at `maskSecret`; a naive squash could drop V10 safeguards |
| Resolution | Keep #2261 as its own commit; V10 squash should merge its rune logic into the broader V10 `maskSecret` and retain the V10 functions |
| Tests | Keep upstream UTF-8 tests; retain V10 empty/adjacent/required-rule regressions |

### 2. Configuration regex compilation

| Item | Detail |
| --- | --- |
| Upstream | #2082, `config/config.go`, `regexp/*_regex.go`, invalid fixture |
| Integration | `regexp.Compile` wrappers and compile-at-translate/allowlist sites; `config/config.go:125–130`, `:348–356` |
| CP-V10 | `config/config.go:150–159`, `:383–393`; same wrappers, expanded contextual errors and additional invalid-format tests |
| Behavioral difference | Both fail instead of panic. V10 reports richer paths and covers extended/legacy configurations; integration keeps upstream messages |
| Risk | Low; same abstraction and error direction |
| Resolution | Preserve #2082 commit, then let V10's error-message/test expansion win |
| Tests | Full config and RE2 backend tests |

### 3. Archive identification

| Item | Detail |
| --- | --- |
| Upstream | #2317, `sources/file.go`, `sources/file_2316_test.go` |
| Integration | `archivePeekSize`, `contentMatchesArchive` returning `bool`, `verifyArchiveStream`; `sources/file.go:28`, `:117–147` |
| CP-V10 | `sources/file.go:29`, `:170–205`; same bounded 64 KB peek and format dispatch, returns `(bool,error)`, shared `restorePrefix`, explicit read-error propagation |
| Behavioral difference | Integration treats read/verify errors as “not an archive” and falls back to plain content; V10 reports the error. V10 restores the prefix through one helper and handles non-seekable readers |
| Risk | Integration can silently scan a corrupt archive as text; V10 is stricter |
| Resolution | Preserve #2317 for ancestry; V10 should supersede the error path |
| Tests | Keep plaintext-with-archive-name, real ZIP/gzip tests; add V10 read-error/prefix-replay tests |

### 4. File and archive size propagation

| Item | Detail |
| --- | --- |
| Upstream | #2099, 10 files across `cmd`, `detect`, `sources` |
| Integration | `MaxFileSize` propagated; archive member precheck; decompressor uses `bytes.Buffer` and `io.CopyN`; `sources/file.go:52`, `:201–206`, `:236–251` |
| CP-V10 | `MaxFileSize` plus `Size`; bounded temporary spooling (`spoolLimitedReader`), exact byte limits, informative exclusion records; `sources/file.go:209–238`, `:265–329`, `:361–368` |
| Behavioral difference | Integration buffers an entire decompressed stream in memory after checking its limit; V10 spools and tracks size/unknown-size exclusion detail. V10 covers Git/stdin/reader and member paths more broadly |
| Risk | Memory pressure on large compressed input in integration; no V10 informative records |
| Resolution | Preserve #2099; V10 supersedes memory handling and should keep field names where possible |
| Tests | Keep size-limit tests; V10 spool/unknown-size tests remain authoritative |

### 5. Staged binary-marked text

| Item | Detail |
| --- | --- |
| Upstream | #2338, `sources/git.go`, `detect/git_staged_binary_test.go` |
| Integration | `staged` field and `yieldFromBlob`; reads `commitSHA:path` via `NewBlobReaderContext`; `sources/git.go:36–39`, `:374–385`, `:420` |
| CP-V10 | Same staged guard, but uses `NewIndexBlobReaderContext` (`:0:path`) and MIME classification; `sources/git.go:629–640`, `:702–703` |
| Behavioral difference | The upstream form relies on an empty commit SHA producing a `:path` index lookup; V10 makes stage-zero explicit, propagates errors, and tracks blob exclusion |
| Risk | Integration has implicit/fragile index behavior and can return nil on blob error; V10 fails operationally |
| Resolution | Preserve #2338; V10's explicit index path is the accepted correction |
| Tests | Upstream big-file-threshold/-diff fixtures plus V10 staged-index errors |

### 6. Git parser integrity, stderr and process lifecycle

| Item | Detail |
| --- | --- |
| Upstream | #2246, #2284, #2339, #2080 |
| Integration | `newGitCmd`, `drain`, `UnparsedBytes`; `--no-color` insertion; `waitCmd`; direct stderr return; `sources/git.go:84–137`, `:187–217`, `:374–505` |
| CP-V10 | `forceNoColor`, `gitOutputTracker`, `gitStderr`, `commandError`, `scanerr`, raw section-count comparison, duplicate-safe `Wait`; `sources/git.go:30–51`, `:120`, `:400–432`, `:577–583` |
| Behavioral difference | Integration detects leftover parser bytes and a nonzero exit; V10 also validates raw section counts, classifies fatal vs warning stderr, sanitizes diagnostics, bounds errors via `scanerr`, and avoids duplicate wait |
| Risk | Integration may still miss a truncated/parse-inconsistent stream or misclassify warnings; V10 is the accepted completeness layer |
| Resolution | Preserve all four selected commits; V10 reconciliation must not drop their tests or author ancestry |
| Tests | Parser overflow, blocking pipe, non-repo stderr, forced color, exit status, V10 section-count tests |

### 7. Root-relative directory paths

| Item | Detail |
| --- | --- |
| Upstream | #2248, seven files |
| Integration | `Fragment.Root`, `File.Root`, `RootRelativePath` with explicit `!` split, `Files.root`, `shouldSkip`, `isConfigOrBaseline`, `knownRelativeToRoot`, `pathAllowedRelativeToRoot`; `sources/fragment.go:29–62`, `sources/files.go:61–160`, `detect/detect.go:724–853` |
| CP-V10 | Same concepts, but `RootRelativePath` relies on lexical full path; `fragmentPaths`, `configOrBaselinePath`, `samePathAtRoot`, absolute baseline normalization and archive-root tests |
| Behavioral difference | Integration splits on `!` explicitly; V10 treats the archive member path lexically. V10 covers baseline normalization, absolute/dot/relative roots, config self-exclusion and root-relative allowlists more completely |
| Risk | Archive-member edge cases and cross-volume relative paths; V10's F16 corrections must win |
| Resolution | Preserve #2248; V10 F16 implementation is the reconciler |
| Tests | Keep `dir_root_test.go`/`files_test.go`; retain CP-V8 cross-volume additions |

### 8. Directory CLI arguments

| Item | Detail |
| --- | --- |
| Upstream | #2334, `cmd/directory.go` |
| Integration | `cobra.MaximumNArgs(1)` and `resolveDirSource`; empty arg defaults to `.` |
| CP-V10 | Same maximum and function; V10 rejects an explicitly empty argument and reports it |
| Behavioral difference | Integration silently treats an empty positional argument as `.`; V10 fails closed |
| Risk | Low; V10's behavior is stricter and already validated |
| Resolution | Preserve #2334; V10 behavior wins in reconciliation |
| Tests | Extra args, missing path, empty arg |

### 9. AND-allowlist path sufficiency

| Item | Detail |
| --- | --- |
| Upstream | #2227, two commits, `sources/common.go` |
| Integration | `pathIsSufficientToSkip` and `shouldSkipPath`; `sources/common.go:78–105` |
| CP-V10 | Same function plus `shouldSkipPathAtRoot` for F16; `sources/common.go:45`, `:85–110` |
| Behavioral difference | Integration checks path sufficiency but does not evaluate root-relative paths during enumeration; V10's `shouldSkipPathAtRoot` does |
| Risk | Absolute/relative directory roots can differ in allowlist skip behavior on integration |
| Resolution | Preserve both #2227 commits; V10 root-aware function is additive |
| Tests | OR/AND/nil/commit/stopword cases plus V10 root tests |

### 10. File and fragment error propagation

| Item | Detail |
| --- | --- |
| Upstream | #2249 (files), #2253 (fragments) |
| Integration | local `skipped`/`fragmentErrs` slices and `errors.Join`; `sources/files.go:143–210`, `detect/detect.go:216–285` |
| CP-V10 | `internal/scanerr` bounded collector, `Files.scanFile`, collector integration; `sources/files.go:195`, `detect/detect.go:324` |
| Behavioral difference | Integration returns all errors joined (unbounded list); V10 bounds and deduplicates error detail and keeps scan results. V10 also has broader walker/read/worker coverage |
| Risk | Large unreadable trees can produce large error lists; V10's collector is the accepted behavior |
| Resolution | Preserve both commits; V10 `scanerr` implementation supersedes |
| Tests | Unreadable file, partial scan, readable-fragment continuation |

### 11. Report validation and permissions

| Item | Detail |
| --- | --- |
| Upstream | #2307 (preflight removal), #2312 (0600 final write) |
| Integration | No preflight truncation; `openReportFile` final write; `cmd/root.go:347–352`, `:467`, `:502` |
| CP-V10 | Entire publication model: `cmd/report.go` `newReportOutput`, temporary rendering, pre-publication validation, aggregate/reject reservations, rename retry; temp files `Chmod(0600)` |
| Behavioral difference | Integration avoids preflight destruction and writes 0600, but still writes directly to the final destination with no temp/rename. V10 preserves old reports during scan failures and supports aggregate/reject |
| Risk | Integration has the old direct-write failure mode; V10 supersedes |
| Resolution | Preserve #2307 and #2312 for provenance; V10 report implementation replaces the direct write |
| Tests | Upstream preservation/permission tests; V10 report pair/TOCTOU tests |

## V10 behavior absent from the integration branch

The following are expected pre-V10 gaps, not regressions:

- F16 full root-relative baseline/ignore/allowlist and cross-volume corrections
- RNC V1 AND-allowlist detector integration and root-aware skip
- Git lifecycle error classification, section counting and `scanerr`
- Staged index blob explicit resolution and MIME/error handling
- Bounded size spooling and informative reporting
- Report preflight/temp/atomic/aggregate/reject architecture and Windows rename retry
- Informative reporting, build contract, RE2 generator fix, Docker and CI work

## Build/test results

| Check | Result |
| --- | --- |
| `go build ./...` | Pass |
| `go test -count=1 ./...` | Pass |
| `go vet ./...` | Fail in pre-existing base test `detect/detect_test.go:2084` (`lostcancel`) |
| `go test -count=1 -tags gore2regex ./...` | Fail at base generator type mismatch (`airtable.go`, `perplexity.go`); fixed only by CP-V10 |

The vet and tagged failures are limitations of the pre-V10 base, not of the
selected upstream commits. They are recorded rather than hidden.

## Recommended reconciliation

1. Keep the 18 upstream commits and the two RNC adjustments intact.
2. Apply the CP-V10 implementation as a separate RNC-authored reconciliation
   commit (or one squash) after the contributor commits, never before.
3. Resolve each conflict in favor of the validated V10 behavior, while retaining
   contributor tests and the specific upstream logic that V10 did not supersede.
4. Re-run full standard/tagged tests and the F16/report/Git regressions.
5. Compare the resulting tree to `bd42a95`; treat only unexplained behavior
   differences as defects.

The integration branch is **not** a release candidate and must not be promoted
as `rnc-mod`.

# Contributor Integration Manifest

**Branch:** `audit/rnc-mod-upstream-integration`
**Base:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (upstream `master`)
**Final integration HEAD before this documentation:** `100ae90`
**CP-V10 comparison source:** `bd42a958f7476e9a6d8e17bd2dbb2abe3ff060ea`

This branch is an experimental provenance and integration foundation. It is not
`rnc-mod`, is not a production baseline, and does not contain the RNC V10
implementation or any V10 squash.

## Integration summary

- 17 upstream PRs attempted and integrated.
- 18 original upstream commits preserved with their original author names and
  email addresses.
- 2 RNC integration-adjustment commits.
- 7 conflicting cherry-picks resolved; every resolution is listed below.
- 5 secondary proposals were screened but not integrated because they were
  alternative or superseded implementations; see
  `contributor-audit.md` and `attribution-summary.md`.
- No protected branch was modified.

## Sequence and SHA crosswalk

Cherry-picks were made with `git cherry-pick -x`; each commit message retains
the original upstream SHA. Sequence is the actual first-parent order.

| # | PR | Original SHA | Fork SHA | Original author | Result |
| ---: | ---: | --- | --- | --- | --- |
| 1 | #2227 | `894b4a3b04dfd7c10e9fb887cdaaad35e4b79c31` | `3004bda` | BasitS-hash `<basitsherazi319@gmail.com>` | Clean |
| 2 | #2227 | `49b26eece1ea83ea58da0c6e0bf92326fb82a25c` | `c8d7070` | BasitS-hash `<basitsherazi319@gmail.com>` | Clean |
| 3 | #2334 | `ad321b7d7aa9b1a1642bb56813fa7443d5412971` | `bdc880a` | Shanmukh-l33t `<Shanmukh-l33t@users.noreply.github.com>` | Clean |
| 4 | #2261 | `b1ce07f27d8251ef22e5b47c80387e9b1590fc2f` | `1f3533b` | Arpit Jain `<arpitjain099@gmail.com>` | Clean |
| 5 | #2082 | `df8223e0008b6d5bf82d0433592b8f1557c97311` | `7a4c18f` | Smeet Agrawal `<smeetagrawal23@gmail.com>` | Clean |
| 6 | #2307 | `0d1ae3ea1087bebcf5ce45982a770d0f4925af6c` | `d68ecb2` | Jeff `<138259519+jdotc1@users.noreply.github.com>` | Clean |
| 7 | #2329 | `21028d8858e31fd088029f2fdb6634c99a1d3563` | `2cf0ade` | cardene777 `<cardene777@gmail.com>` | Clean |
| 8 | #2284 | `1894e6f2cdaef94fe66388e577e0f59f14bfcff2` | `f146e23` | Sakariyah Abdulhazeem `<sakariyahabdulhazeem@gmail.com>` | Clean |
| 9 | #2246 | `9b7b2243fcdeb4279cccd9df9348f9e537d4fb95` | `5bed958` | INTFRAME `<hello@intframe.com>` | Conflict in `sources/git_test.go`; resolved |
| 10 | #2338 | `7fed4140d87f62090c87b9f09c5c9b6073a84a27` | `16ba49e` | KVR Vishnu `<46535652+VISHNU0906@users.noreply.github.com>` | Conflict in `sources/git.go`; resolved |
| 11 | #2339 | `7e5132faa1ccea4bfd0b4b006ca896b5dd03df3f` | `b603fa0` | DelaySystem `<DelaySystem@tubot.lat>` | Conflict in `sources/git.go` and `sources/git_test.go`; resolved |
| 12 | #2080 | `ff77f9f58b1ec5fbbece852e0bfe9c1f050f24e5` | `2bd5ad3` | Smeet Agrawal `<smeetagrawal23@gmail.com>` | Conflict in `sources/git_test.go`; resolved |
| 13 | #2099 | `1b3896fd7df9a667a09568b751dce6143613319e` | `62a89d4` | Nikhil `<tad.areas_0y@icloud.com>` | Add/add conflict in `sources/file_test.go`; resolved |
| 14 | #2317 | `71e91200afb2415b283bb1721709b19c2fddd8aa` | `7ce04bb` | Varshith Puli `<pulivarshit@gmail.com>` | Clean |
| 15 | #2248 | `5c9ca86f8f4d1db1e17004b6a1b598137203face` | `8ee08ef` | hamodywe `<iosapk.org@gmail.com>` | Conflict in `sources/file.go`; resolved |
| 16 | #2249 | `509fce15ad2eb2536bc8d840a3ed3e374d73bb2d` | `5e35358` | Vaibhav Srivastava `<vaibhavsri1712@gmail.com>` | Clean |
| 17 | #2253 | `4f58f98a1424203af39cc00a26f9672e9b45b794` | `260d7e1` | Arpit Jain `<arpitjain099@gmail.com>` | Clean |
| â€” | integration fix | â€” | `22f63b4` | Dave Woods `<209306207+davidwoods15@users.noreply.github.com>` | RNC adjustment after #2246/#2338 |
| 18 | #2312 | `2009273e0029f91aeda4732552ae5489634d0749` | `6445d1e` | Varshith Puli `<pulivarshit@gmail.com>` | Conflict in `cmd/root.go`; resolved |
| â€” | integration fix | â€” | `100ae90` | Dave Woods `<209306207+davidwoods15@users.noreply.github.com>` | RNC Windows test adjustment after #2312 |

## Conflict resolutions

### #2246 â€” `sources/git_test.go`

PR #2246 replaced `GitCmd` construction with `newGitCmd` and added parser-drain
tests. #2284 had already added a color-forced staged-diff regression test in the
same file. The conflict was in the import block and the adjacent test blocks.

Resolution: kept both test suites, unioned the imports, and ran `gofmt`. No
production conflict.

### #2338 â€” `sources/git.go`

#2246 moved `GitCmd` construction into `newGitCmd`; #2338 added a `staged` field
at the older construction site.

Resolution: added `staged bool` alongside `drainOnce`/`unparsed`, retained the
new constructor, and set `gitCmd.staged = staged` in
`NewGitDiffCmdContext` after `newGitCmd` returns. This required the separate
`22f63b4` adjustment because the first conflict resolution left the field
assignment in a helper without a `staged` parameter.

### #2339 â€” `sources/git.go` and test file

#2246 returned an â€œunparsed bytesâ€ error; #2339 added a `waitCmd` wrapper that
turns a non-zero Git exit into an error. Both belong on the normal completion
path.

Resolution: after `wg.Wait()`, return the parser-integrity error first when
`unparsed > 0`, otherwise return the `waitCmd()` error. The test file keeps the
#2246 parser tests and the #2339 exit-code test. Imports were unioned.

### #2080 â€” `sources/git_test.go`

The production change auto-merged: the stderr channel error is returned
directly instead of being routed through `yield`. The test file conflicted on
imports and the placement of the new stderr-propagation test.

Resolution: unioned imports, preserved both the #2284 color test and the #2080
stderr test, and retained the later parser tests.

### #2099 â€” `sources/file_test.go`

#2329 added `sources/file_test.go` with first-chunk MIME tests; #2099 added the
same path with archive-size tests. Git reported an add/add conflict.

Resolution: combined both test sets in one file and retained all helpers.

### #2248 â€” `sources/file.go`

#2099 moved archive-member size checks before opening the member; #2248 added
`Root` propagation to the member `File`.

Resolution: kept the size precheck before `d.Open()`, added `Root: s.Root` to
the member struct, and preserved `MaxFileSize` propagation.

### #2312 â€” `cmd/root.go`

#2307 removed the destructive preflight `os.Create` from report validation.
#2312 added `openReportFile` at both the preflight and final publication sites.

Resolution: kept #2307's non-destructive validation, dropped only the #2312
preflight hunk, and retained `openReportFile` for the final report write. This
avoids reintroducing destination truncation while preserving the contributor's
0600-permission implementation. The Windows-specific permission assertion in
the new test was skipped in `100ae90`; Windows does not expose Unix mode bits
through `FileMode.Chmod`.

## Dependency order

The integration follows this dependency-aware order:

1. Independent CLI, allowlist, redaction, config, report-validation and
   first-chunk changes (#2227, #2334, #2261, #2082, #2307, #2329).
2. Git command construction (#2284, #2246).
3. Staged and parser-error behavior (#2338, #2339, #2080).
4. File size and archive verification (#2099, #2317).
5. Root-relative path handling (#2248).
6. File/fragment error propagation (#2249, #2253).
7. Report permission handling (#2312).

The only non-upstream commits are the two documented integration adjustments.

## Resulting file changes

The integrated branch changes 34 files relative to `5f6ec56` (before audit
documentation), principally:

- `sources/git.go`, `sources/git_test.go`
- `sources/file.go`, `sources/file_test.go`, `sources/file_2316_test.go`
- `sources/files.go`, `sources/files_test.go`, `sources/files_unreadable_test.go`
- `sources/common.go`, `sources/common_test.go`
- `sources/fragment.go`
- `detect/detect.go`, `detect/detect_test.go`, `detect/baseline.go`,
  `detect/dir_root_test.go`, `detect/git_staged_binary_test.go`
- `cmd/root.go`, `cmd/root_test.go`, `cmd/directory.go`,
  `cmd/directory_test.go`, `cmd/report_preservation_test.go`
- `config/config.go`, `config/config_test.go`,
  `regexp/stdlib_regex.go`, `regexp/wasilibs_regex.go`
- `report/finding.go`, `report/finding_test.go`
- `testdata/config/invalid/allowlist_bad_regex.toml`

The exact final file list is:
`git diff --name-only 5f6ec56..HEAD`.

## Not integrated

| PR | Contributor | Reason |
| ---: | --- | --- |
| #2300 | ToniAdreal | Alternative 0600 report implementation; still truncates during preflight. #2312 was selected for the permission code, and #2307 removed destructive preflight. |
| #2292 | Saadanjum0 | Alternative stderr-return implementation; conflicts with #2080, which was selected. |
| #2252 | CAOShurong | Immediate fragment-error return would abandon readable findings; V10 and the integrated branch preserve continued collection. |
| #2313 | Varshith-Kali | Alternative unreadable-file error handling; #2249 was selected for the same behavior. |
| #2203 | RajeshRajendiran | Drain-fix alternative overlapping #2246/#2299; #2246 was selected. |

These are documentary/design acknowledgments, not commit-level adoption. No
cherry-pick was attempted because each is a duplicate or incompatible
implementation of a behavior already represented by a selected commit.

## Verification

On `100ae90` plus this documentation's source tree:

```text
go build ./...                         PASS
go test -count=1 ./...                 PASS
go vet ./...                           FAIL: pre-existing base lostcancel test at detect/detect_test.go:2084
go test -count=1 -tags gore2regex ./... FAIL: known base generator type mismatch (fixed only in CP-V10)
```

The failure contexts are recorded in `source-diff-analysis.md`. Full V10
equivalence is neither expected nor claimed; this branch deliberately stops
before the RNC implementation squash.

## Attribution rules applied

- Original author name/email preserved by `git cherry-pick -x`.
- No `Co-authored-by` trailers were fabricated.
- No contributor commit was squashed into another contributor commit.
- No V10 commit was applied.
- The two RNC adjustment commits are clearly identified as integration
  maintenance, not upstream authorship.
- The scanner tree remains a pre-V10 experimental branch.

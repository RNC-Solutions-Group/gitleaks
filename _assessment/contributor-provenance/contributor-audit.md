# Contributor Attribution Audit

**Repository:** `RNC-Solutions-Group/gitleaks`
**Audit branch:** `audit/rnc-mod-upstream-integration`
**Upstream base:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`
**CP-V10 comparison:** `bd42a958f7476e9a6d8e17bd2dbb2abe3ff060ea`
**Current validation branch:** local `fix/validation` at `bd42a95`; remote
`origin/fix/validation` at `48c76717ccbbe2159185edb4510fb1924d7dbcc6`

## Scope and method

All 22 requested PRs were retrieved from the upstream repository through the
GitHub API and their actual head commits were fetched with Git. Every claimed
author and email below comes from the fetched commit object, not from the PR
display name. Every PR in this audit is **open and unmerged upstream** at the
time of retrieval.

For each primary PR I compared:

- the original patch against its actual merge base;
- the production source and test files;
- the corresponding CP-V10 file and symbol;
- the integrated branch result;
- the surviving code, adapted code, and superseded behavior.

A raw shared-line count is not used as proof of copyrightable copying. The
confidence labels describe implementation correspondence, not legal advice.

## Reference discrepancy

`fix/validation` is locally at `bd42a95`, while `origin/fix/validation` is at
`48c7671`. The three intervening commits are documentation-only:

- `5ef695e` reconcile CP-V8â€“V10 implementations and completed backlog
- `43785c1` complete 31-package closure crosswalk and queue accounting
- `48c7671` acknowledge independently matched dependency upgrades

No scanner source changed. This audit uses `bd42a95` as the CP-V10 source
reference and reads the newer ledger through `origin/fix/validation`. The local
protected branch was not moved or rewritten.

## Primary contributors

### #2261 â€” @arpitjain099 â€” rune-safe masking

- **Head:** `b1ce07f27d8251ef22e5b47c80387e9b1590fc2f`
- **Commit author:** Arpit Jain `<arpitjain099@gmail.com>`
- **Files:** `report/finding.go`, `report/finding_test.go`
- **State:** open, unmerged

**Production overlap.** The PR converts `maskSecret` from byte length/slicing to
`[]rune` length/slicing. CP-V10 `report/finding.go:231â€“254` contains the same
`runes := []rune(secret)`, `length := float64(len(runes))`,
`RoundToEven(length * prc / 100)`, and `string(runes[:lth]) + "..."` algorithm.
The source comment explicitly cites `gitleaks/gitleaks#2261
(@arpitjain099)`. V10 adds `percent == 0` and bounds clamps.

**Test overlap.** The PR adds multibyte and mixed-width cases plus
`utf8.ValidString`. V10 has broader redaction tests, including empty and
adjacent-secret cases; the original test expressions are substantially
reflected in the V10 suite.

**Survival/adaptation.** Exact algorithm and adapted tests survive. V10 fixes
additional escaping cases the original PR did not cover.

**Recommendation:** Category 1; preserve original commit. Confidence: high.

### #2082 â€” @Smeet23 â€” invalid regex error handling

- **Head:** `df8223e0008b6d5bf82d0433592b8f1557c97311`
- **Commit author:** Smeet Agrawal `<smeetagrawal23@gmail.com>`
- **Files:** `config/config.go`, `config/config_test.go`,
  `regexp/stdlib_regex.go`, `regexp/wasilibs_regex.go`,
  `testdata/config/invalid/allowlist_bad_regex.toml`
- **State:** open, unmerged

**Production overlap.** The PR adds `regexp.Compile` wrappers to both backends
and replaces `MustCompile` at rule path/regex and allowlist regex/path sites.
CP-V10 has the same wrappers in `regexp/stdlib_regex.go` and
`regexp/wasilibs_regex.go`, and the same call-site structure in
`config/config.go:150â€“159` and `:383â€“393`, with expanded contextual messages.

**Test overlap.** The PR adds an invalid allowlist-path fixture and a
table-driven error assertion. V10 retains the invalid-pattern test family and
adds formats/extensions and CLI coverage.

**Survival/adaptation.** The wrapper functions and call-site approach survive;
V10 improved diagnostics and extended validation. This is substantive
implementation reuse, not merely an idea.

**Recommendation:** Category 1; preserve original commit. Confidence: high.

### #2317 â€” @Varshith-Kali â€” archive content verification

- **Head:** `71e91200afb2415b283bb1721709b19c2fddd8aa`
- **Commit author:** Varshith Puli `<pulivarshit@gmail.com>`
- **Files:** `sources/file.go`, `sources/file_2316_test.go`
- **State:** open, unmerged

**Production overlap.** The PR introduces `archivePeekSize = 64 * 1_000`,
`contentMatchesArchive`, `verifyArchiveStream`, compressed-format dispatch,
seek/rewind, and non-seekable prefix replay. CP-V10
`sources/file.go:26â€“29`, `:112â€“180`, and `:199â€“205` retains the same constant,
function names, bounded peek, format dispatch, and replay concept. V10 returns
an error and uses a shared `restorePrefix` helper.

**Test overlap.** The PR adds archive-like plaintext, real ZIP, and real gzip
tests. V10 retains expanded archive-identification tests and error paths.

**Survival/adaptation.** Strong direct/adapted overlap. The original code's
error suppression was corrected in V10.

**Recommendation:** Category 1/2; preserve original commit. Confidence: high.

### #2338 â€” @VISHNU0906 â€” staged binary fallback

- **Head:** `7fed4140d87f62090c87b9f09c5c9b6073a84a27`
- **Commit author:** KVR Vishnu `<46535652+VISHNU0906@users.noreply.github.com>`
- **Files:** `sources/git.go`, `detect/git_staged_binary_test.go`
- **State:** open, unmerged

**Production overlap.** The PR adds `GitCmd.staged`, changes the binary-file
branch to `yieldFromBlob`, and reads a staged index blob. V10
`sources/git.go:629â€“640` and `:702â€“703` contains the same staged/binary guard
and index fallback, with `NewIndexBlobReaderContext(":0:"+path)` and MIME
handling. The integration branch has the PR's `NewBlobReaderContext` path; V10
corrects it to the explicit stage-zero index path.

**Test overlap.** The PR's `core.bigFileThreshold` and `-diff` fixture is the
direct ancestor of V10's staged binary tests.

**Survival/adaptation.** Design and staged flag survive; V10 corrected index
resolution, errors and lifecycle. This is adapted implementation, not mere
research.

**Recommendation:** Category 2; preserve original commit, then reconcile V10's
index correction. Confidence: high.

### #2099 â€” @alephnull-sh â€” archive/member size limits

- **Head:** `1b3896fd7df9a667a09568b751dce6143613319e`
- **Commit author:** Nikhil `<tad.areas_0y@icloud.com>`
- **Files:** 10 files under `cmd/`, `detect/`, `sources/`
- **State:** open, unmerged

**Production overlap.** The PR adds `MaxFileSize` propagation through `cmd`,
`detect`, `Files`, `Git`, archive members, decompressors, and stdin/reader.
CP-V10 `sources/file.go:52â€“60`, `:209â€“238`, `:317â€“329`, and `:361â€“368` carries
the same field and member/decompressor prechecks. V10 replaces the PR's
`bytes.Buffer` whole-content copy with bounded temporary spooling and adds
`Size`/informative records.

**Test overlap.** The PR adds archive-member and compressed-file over-limit
tests. V10 expands these into size-policy, spooling, Git and CLI tests.

**Survival/adaptation.** Field propagation and precheck behavior survive;
memory strategy was deliberately replaced. This is partial/adapted reuse.

**Recommendation:** Category 2; preserve original commit. Confidence: high.

### #2246 â€” @aron-intframe â€” Git parser drain and truncation detection

- **Head:** `9b7b2243fcdeb4279cccd9df9348f9e537d4fb95`
- **Commit author:** INTFRAME `<hello@intframe.com>`
- **Files:** `sources/git.go`, `sources/git_test.go`
- **State:** open, unmerged

**Production overlap.** The PR introduces `newGitCmd`, buffered stdout,
`drainOnce`, `unparsed`, `UnparsedBytes`, a `Wait` drain, and an unparsed-output
error. CP-V10 `sources/git.go:344â€“424` retains `newGitCmd`, drain/unparsed
mechanics and the same incomplete-results concept; V10 adds a raw
`gitOutputTracker` section count and error/collector integration.

**Test overlap.** The PR's overflowed hunk header fixture and blocking-pipe
tests are reflected in V10's parser-integrity tests.

**Survival/adaptation.** Substantial direct/adapted reuse. V10's section-count
check is an additional correction.

**Recommendation:** Category 1/2; preserve original commit. Confidence: high.

### #2248 â€” @hamodywe â€” root-relative paths

- **Head:** `5c9ca86f8f4d1db1e17004b6a1b598137203face`
- **Commit author:** hamodywe `<iosapk.org@gmail.com>`
- **Files:** `detect/baseline.go`, `detect/detect.go`, `detect/dir_root_test.go`,
  `sources/file.go`, `sources/files.go`, `sources/files_test.go`,
  `sources/fragment.go`
- **State:** open, unmerged

**Production overlap.** The PR introduces `Fragment.Root`, `File.Root`,
`RootRelativePath`, `Files.root`, root-relative allowlist scanning,
`isConfigOrBaseline`, `knownRelativeToRoot`, and
`pathAllowedRelativeToRoot`. CP-V10 has the same `Root` field and
`RootRelativePath` in `sources/fragment.go:33â€“56`, and retains the root-relative
matching design across `detect/detect.go:402â€“459`, `detect/baseline.go:109â€“157`
and `sources/files.go:151â€“160`. V10 refactors toward `fragmentPaths`,
`configOrBaselinePath`, and `samePathAtRoot`.

**Test overlap.** The PR's `detect/dir_root_test.go` and
`sources/files_test.go` are the origin of the V10 F16 test family, which later
expanded to baseline/config/ignore/allowlist cases.

**Survival/adaptation.** Strong adaptation. V10 fixed archive-member and
absolute-baseline/normalization cases and changed the helper structure.

**Recommendation:** Category 1/2; preserve original commit. Confidence: high.

### #2334 â€” @Shanmukh-l33t â€” explicit directory target

- **Head:** `ad321b7d7aa9b1a1642bb56813fa7443d5412971`
- **Commit author:** Shanmukh-l33t `<Shanmukh-l33t@users.noreply.github.com>`
- **Files:** `cmd/directory.go`, `cmd/directory_test.go`
- **State:** open, unmerged

**Production overlap.** The PR adds `cobra.MaximumNArgs(1)`,
`resolveDirSource`, extra-argument rejection and existence validation. CP-V10
`cmd/directory.go:23â€“45` contains the same `MaximumNArgs(1)` and
`resolveDirSource` contract; V10 additionally rejects an explicitly empty
argument instead of defaulting it to `.`.

**Test overlap.** The PR's test cases are substantially retained in V10's
directory tests.

**Survival/adaptation.** Direct implementation reuse with one V10 behavior
correction.

**Recommendation:** Category 1; preserve original commit. Confidence: high.

### #2227 â€” @BasitS-hash â€” AND-allowlist path sufficiency

- **Head:** `49b26eece1ea83ea58da0c6e0bf92326fb82a25c`
- **Commits:** `894b4a3b04dfd7c10e9fb887cdaaad35e4b79c31`,
  `49b26eece1ea83ea58da0c6e0bf92326fb82a25c`
- **Commit author:** BasitS-hash `<basitsherazi319@gmail.com>`
- **Files:** `sources/common.go`, `sources/common_test.go`
- **State:** open, unmerged

**Production overlap.** The PR introduces `pathIsSufficientToSkip` and the
check in `shouldSkipPath`; the review commit adds nil safety, commit handling,
and RE2 backend tests. CP-V10 `sources/common.go:85â€“110` contains the same
function and nearly identical OR/AND logic. V10 adds the root-aware
`shouldSkipPathAtRoot` used by F16.

**Test overlap.** The PR's table-driven OR/AND/stopword/commit/nil cases are
substantially retained and extended in V10.

**Survival/adaptation.** Direct reuse with V10 root-awareness extension. Both
original commits are preserved separately.

**Recommendation:** Category 1; preserve both original commits. Confidence:
high.

### #2284 â€” @lekanjava â€” forced Git color

- **Head:** `1894e6f2cdaef94fe66388e577e0f59f14bfcff2`
- **Commit author:** Sakariyah Abdulhazeem `<sakariyahabdulhazeem@gmail.com>`
- **Files:** `sources/git.go`, `sources/git_test.go`
- **State:** open, unmerged

**Production overlap.** The PR adds `--no-color` to log and diff command
construction and a forced-color fixture. CP-V10 `sources/git.go:30â€“51`,
`:303`, and `:324â€“330` implements `forceNoColor` with argument-order handling
and warns about user-supplied color options. The behavior is the same; the V10
implementation is broader than the PR's straight insertion.

**Test overlap.** The PR's forced-color staged test is reflected in V10's
`git_color_test.go` family.

**Survival/adaptation.** Adapted implementation; direct behavior and tests
survive.

**Recommendation:** Category 2; preserve original commit. Confidence: high.

### #2329 â€” @cardene777 â€” first-chunk MIME guard

- **Head:** `21028d8858e31fd088029f2fdb6634c99a1d3563`
- **Commit author:** cardene777 `<cardene777@gmail.com>`
- **Files:** `sources/file.go`, `sources/file_test.go`
- **State:** open, unmerged

**Production overlap.** The PR replaces `totalLines == 0` with an explicit
`isFirstChunk` flag and resets it after the first classification. CP-V10
`sources/file.go:457` and `:484â€“485` contains the same flag and assignment.

**Test overlap.** The PR's long-first-line and binary-signature tests are
substantially retained in V10's chunk-boundary tests.

**Survival/adaptation.** Direct implementation reuse.

**Recommendation:** Category 1; preserve original commit. Confidence: high.

## Secondary contributors

| PR | Contributor | Head | Initial finding | Disposition |
| ---: | --- | --- | --- | --- |
| #2307 | @jdotc1 | `0d1ae3e` | Removes destructive preflight `os.Create`; adds preservation test | Integrated (`d68ecb2`); adapted by V10 report architecture |
| #2300 | @ToniAdreal | `2b08eb9` | Alternative 0600 `createReportFile`; preflight still truncates | Not integrated; #2312 permission code selected, #2307 preserved non-destructive validation |
| #2312 | @Varshith-Kali | `2009273` | `openReportFile` 0600 + tests; preflight still truncates | Integrated (`6445d1e`); preflight hunk dropped, permission code retained |
| #2339 | @delaysystem | `7e5132f` | `waitCmd` reaps Git and returns nonzero exit | Integrated (`b603fa0`); merged with parser drain |
| #2080 | @Smeet23 | `ff77f9f` | Return stderr error directly instead of yielding it | Integrated (`2bd5ad3`); V10 expands stderr classification |
| #2292 | @Saadanjum0 | `de6f59a` | Alternative stderr-error accumulation | Not integrated; conflicts with selected #2080 |
| #2252 | @CAOShurong | `424b4a9` | Immediate fragment-error return | Not integrated; would abandon readable findings; V10 uses collector |
| #2253 | @arpitjain099 | `4f58f98` | Collect fragment errors and return joined error | Integrated (`260d7e1`); V10 uses `internal/scanerr` |
| #2249 | @vaibhav8a | `509fce1` | Collect unreadable-file errors and propagate | Integrated (`5e35358`); V10 uses `internal/scanerr` |
| #2313 | @Varshith-Kali | `acfeacd` | Alternative open-error aggregation | Not integrated; #2249 selected for the same behavior |
| #2203 | @RajeshRajendiran | `221ca62` | Drain errCh without dropping queued diff files | Not integrated; #2246/#2299 selected; V10 drain/section logic supersedes |

The secondary patches were inspected in full because each contains production
code and regression tests. The non-integrated items are alternative
implementations, not research-only comments. They receive documentary/design
recognition but are not assigned fabricated commits.

Detailed code correspondence:

- **#2307:** V10 has no preflight destination creation; publication is owned by
  `cmd/report.go` (`newReportOutput`, temp rendering, final validate/rename).
  The upstream deletion is semantically subsumed but was retained as an
  independent commit in the pre-V10 integration branch.
- **#2312:** V10 uses temp-file `Chmod(0o600)` in `cmd/report.go`; the
  integrated branch retains the contributor's `openReportFile` final-write
  helper. V10's stronger temp publication supersedes the preflight/truncation
  portions.
- **#2339:** V10's `GitCmd.Wait`/error collector covers exit errors; the
  integration adds the contributor's `waitCmd` on top of the selected parser
  drain.
- **#2080 / #2292:** Both return stderr failures; only one can be selected.
  V10's stderr classifier (`gitStderr`, `isFatalGitStderr`, collector) is the
  later design.
- **#2253 / #2252:** V10's `scanerr` collector preserves readable findings while
  returning a bounded operational error; #2252's immediate return is rejected.
- **#2249 / #2313:** V10's `sources/files.go` uses `scanerr` rather than a local
  error slice; #2249 is the closer ancestor and was integrated.
- **#2203:** V10 drains the parser and validates raw section counts; #2246 is
  the selected ancestor. #2203's contributor already appears in CP-V9
  documentation.

## Documentation and dependency acknowledgments

The CP-V9 documentation bundle and the 2026-11-13 validation ledger establish
the following documentary credits. Their work is minor prose, corrected
instructions, or independently selected dependency versions; no scanner commit
is claimed.

| Source | Contributor | Recognition |
| --- | --- | --- |
| #2236, #2212 | @RajeshRajendiran | Git added-line and decode-depth documentation corrections |
| #2298 | @bharathvemula | Baseline prose/examples |
| #2287 | @ump45nose | Pre-commit Docker anchor correction (adapted to local anchor) |
| #2282 | @haimingZZ | README/CONTRIBUTING typos |
| #2258 | @jankesec | Configuration example typo |
| #2247 | @CAOShurong | Configuration guide URL |
| #2008 | @jylenhof | Optional mise installation instructions |
| #1263 | @noblehelm | Baseline credential-storage/security guidance |
| #2319 | @abhinavrastogi2627 | One-line pre-commit manifest correction; CP-V9 source commit was RNC-owned |
| #2206 | @trinathone | Non-root Docker proposal reviewed; CP-V10 optional target is RNC implementation and root default retained |
| Dependency graph | @corani, Dependabot | Independently matched upgrades; profile-only appreciation, no PR cherry-pick |

The MIT license permits reuse and requires preservation of the copyright and
permission notice. The repository retains the upstream `LICENSE` with
`Copyright (c) 2019 Zachary Rice`; no contributor-specific copyright notice was
added or removed. Attribution in Git history is voluntary policy, not a
substitute for the license notice.

## Contributor-level count

- Distinct PR contributors reviewed: **18** (11 primary plus 7 new secondary
  logins; Varshith-Kali, Smeet23 and Arpit Jain appear in both groups).
- Distinct documentation/dependency contributors reviewed: **7 individuals**
  plus Dependabot, with overlap to the PR group for Rajesh Rajendiran and
  CAOShurong.
- Repeated contributors: Varshith Puli (#2317, #2312, #2313), Smeet Agrawal
  (#2082, #2080), Arpit Jain (#2261, #2253), Rajesh Rajendiran (#2203 and docs),
  CAOShurong (#2252 and docs).

## Outstanding questions

- Exact legal significance of the three-line screening rule is not evaluated;
  this is an engineering and voluntary attribution policy.
- The final V10 squash must preserve the selected upstream commits as separate
  ancestry; otherwise the commit-level credit is lost.
- V10's later rewrites are substantial in several areas; the integration report
  recommends preserving contributor commits and reconciling V10 changes in a
  separate RNC commit, not claiming the final code is byte-identical to the PRs.

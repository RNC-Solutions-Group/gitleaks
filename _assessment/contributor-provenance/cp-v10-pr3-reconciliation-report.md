# CP-V10 PR #3 Reconciliation Checkpoint

**Checkpoint:** CP-V10-PR3-RECONCILIATION
**Date:** 2026-11-13
**Repository:** `RNC-Solutions-Group/gitleaks`
**Candidate branch:** `audit/v10-pr3-reconciliation`
**Long-lived audit branch:** `audit/rnc-mod-upstream-integration`
**Product reference:** `origin/V10-backup` (`48c76717ccbbe2159185edb4510fb1924d7dbcc6`)
**PR reference:** [RNC-Solutions-Group/gitleaks#3](https://github.com/RNC-Solutions-Group/gitleaks/pull/3)

## A. Integration summary

### Reference verification

| Reference | Verified commit | Result |
|---|---|---|
| `master` | `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` | Matches |
| `audit/rnc-mod-upstream-integration` | `78b5c7f30aa1c86e491fc947637016f7ae5ca59c` | Matches |
| `fix/validation` | `48c76717ccbbe2159185edb4510fb1924d7dbcc6` | Matches |
| `V10-backup` | `48c76717ccbbe2159185edb4510fb1924d7dbcc6` | Matches |
| `rnc-mod` | `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` | Matches; not modified |

Local branches `V10-backup` and `rnc-mod` are not present in this checkout; their remote-tracking refs were verified. This has no effect on the comparison or promotion.

PR #3 was verified as open, draft, mergeable state `CONFLICTING`, head `fix/validation` at `48c7671`, base `audit/rnc-mod-upstream-integration` at `78b5c7f`. GitHub reports 64 commits in the PR. Those 64 commits are the V10/validation commits on the `fix/validation` side of merge base `5f6ec56`; they are not the 18 contributor commits, which are on the audit side. The two sides diverged at `5f6ec56`, which explains the PR conflict and is why this checkpoint uses a local squash reconciliation rather than the GitHub Merge button.

### Reconciliation commits

| Item | Commit | Meaning |
|---|---|---|
| Original audit HEAD | `78b5c7f30aa1c86e491fc947637016f7ae5ca59c` | Provenance-integrated foundation |
| Provenance-document encoding correction | `baacad503202894ac9d70ebc27d2ce416e993e37` | Five Markdown reports only |
| RNC V10 integration commit | `9c82fe4d789f486ccdeb0106c6d543721d198b32` | Authoritative product tree |
| Reconciled product HEAD | `9c82fe4d789f486ccdeb0106c6d543721d198b32` | Exact V10 product tree |
| Report-bearing audit tip | docs-only child of `9c82fe4` | This report; same product tree |

The integration commit is authored and committed by Dave Woods, representing the RNC reconciliation commit. The 18 upstream contributor commits retain their original authors and are ancestral.

### Squash and conflicts

`git merge --squash origin/fix/validation` produced **21 conflicted paths**:

`cmd/directory.go`, `cmd/directory_test.go`, `cmd/root.go`, `config/config.go`, `detect/baseline.go`, `detect/detect.go`, `detect/dir_root_test.go`, `detect/files.go`, `regexp/stdlib_regex.go`, `regexp/wasilibs_regex.go`, `report/finding.go`, `report/finding_test.go`, `sources/common.go`, `sources/common_test.go`, `sources/file.go`, `sources/file_test.go`, `sources/files.go`, `sources/files_test.go`, `sources/fragment.go`, `sources/git.go`, and `sources/git_test.go`.

Resolution policy was applied literally: all conflicted and non-conflicted product paths were replaced by the authoritative `origin/V10-backup` blobs using a reset of the index and working tree to the V10 tree. No upstream implementation was preferred for cleanliness, and no hybrid conflict resolution was introduced. The five provenance reports were then restored from the corrected audit-side commit.

The integration commit changes **123 paths** relative to the audit HEAD: **14,203 insertions** and **2,704 deletions**. This includes V10 additions, replacements, and removal of obsolete audit-branch files that are not part of V10. The separate encoding-correction commit changes five documentation paths only (108 insertions, 108 deletions).

## B. Exact-diff report

### Authoritative comparison

Command, run from the repository root:

```text
git diff --exit-code origin/V10-backup HEAD -- . \
  ':(exclude)_assessment/contributor-provenance/**'
```

Result: **exit 0**, no output. There are zero product-tree differences outside the approved provenance documentation.

### File-count and path comparison

| Check | Result |
|---|---|
| Tracked files in `origin/V10-backup` | 524 |
| Tracked files at product commit `9c82fe4` | 529 |
| Tracked files at report-bearing audit tip | 530 |
| Approved provenance reports at `9c82fe4` | 5 |
| Approved provenance paths at report-bearing audit tip | 6 (the five reports plus this report) |
| V10 files missing from candidate | 0 |
| Candidate files absent from V10 at `9c82fe4` | 5, all under `_assessment/contributor-provenance/` |
| Candidate files absent from V10 at report-bearing tip | 6, all documentation-only |
| Unexpected extra product paths | 0 |
| V10 product files accidentally deleted | 0 |
| Obsolete upstream/audit product code remaining | 0 |

The five pre-existing permitted differences are:

- `_assessment/contributor-provenance/attribution-summary.md`
- `_assessment/contributor-provenance/contributor-audit.md`
- `_assessment/contributor-provenance/integration-manifest.md`
- `_assessment/contributor-provenance/rnc-mod-integration-plan.md`
- `_assessment/contributor-provenance/source-diff-analysis.md`

This report is the sixth permitted provenance path and is added by the report-only commit on top of `9c82fe4`.

### Tree-hash proof

A temporary index was populated from the candidate, the five provenance paths were removed, and the resulting product tree was written:

```text
origin/V10-backup tree:       39b13207713d7102d7ac813a0599b2eaabf5d29f
candidate filtered tree:      39b13207713d7102d7ac813a0599b2eaabf5d29f
trees_equal:                  true
```

This is an object-level equality check over every tracked product path and blob, not a behavioral approximation. It includes files present on only one side because the path sets were compared independently as well. Adding this report does not change the filtered tree: the report-only child of `9c82fe4` has the same `39b13207713d7102d7ac813a0599b2eaabf5d29f` product tree.

### Key blob comparison

| Product path | Blob (identical on both sides) |
|---|---|
| `report/finding.go` | `95e5ffe48b6813a41fa9768843e0c9f58c8a9f20` |
| `config/config.go` | `1874e341398092c84ac953300d7b5356859265a3` |
| `sources/file.go` | `963a6e223c03f8249578f2630b779cbaf683fb3f` |
| `sources/git.go` | `ec594dafe7914d182a22d7cfb48742a7888f1184` |
| `sources/fragment.go` | `2acfaa7db17d23536dac45f25df7123358e1217a` |
| `sources/common.go` | `937ea941d6d91ef8ec3e8221bab8c887763da885` |
| `cmd/directory.go` | `58727cd6da21384a4ae90ffbd028ad8d10dd8779` |
| `cmd/root.go` | `fbaf344046c91e793fd1cacf964b2d78823fd39d` |
| `cmd/report.go` | `b548f8579517fd76974a068622ad1137d7ba9ce6` |
| `regexp/stdlib_regex.go` | `ba3daf7f55da1428b26fe4b9543d1fda5ffc7d75` |
| `regexp/wasilibs_regex.go` | `e564ab75f14e92e93611b4ff2247fe5beb501e4a` |
| `Dockerfile` | `65e2918facf9c8a76dde86424beadc2d4cd3728c` |
| `.github/workflows/test.yml` | `222a88a78638119b991557522603b9ca8ddb7ff8` |
| `go.mod` | `78c6dc6fffe813300007d5e8e0b661e67f27fd02` |
| `go.sum` | `4773d926cc682fb8008b0ab7abf2dba4e6eb9456` |
| `Makefile` | `7155aeeb4ee42f9e468e726be17126cacd01a7c9` |

### Verdict

**Exact product equivalence with `origin/V10-backup` is confirmed.** The only differences are the intentionally retained provenance reports.

## C. Contributor preservation

### Preserved upstream commits

All 18 selected upstream commits remain in the candidate's ancestry. The source object for each is present locally, each local commit has one parent, each is an ancestor of `HEAD`, and each local author name/email exactly matches the corresponding source commit.

| # | PR | Integration commit | Upstream source commit | Author | Email |
|---|---|---|---|---|---|
| 1 | #2227 | `3004bda1fc3833c68218f2b8589581fdc99715c9` | `894b4a3b04dfd7c10e9fb887cdaaad35e4b79c31` | BasitS-hash | `basitsherazi319@gmail.com` |
| 2 | #2227 | `c8d707016d51eca7d5e02d1f5066b48ecc668788` | `49b26eece1ea83ea58da0c6e0bf92326fb82a25c` | BasitS-hash | `basitsherazi319@gmail.com` |
| 3 | #2334 | `bdc880a71479a3e3c4e14a86e02858b1831752cb` | `ad321b7d7aa9b1a1642bb56813fa7443d5412971` | Shanmukh-l33t | `Shanmukh-l33t@users.noreply.github.com` |
| 4 | #2261 | `1f3533b1a1e8498ee5e99aba54ee557dda9c0367` | `b1ce07f27d8251ef22e5b47c80387e9b1590fc2f` | Arpit Jain | `arpitjain099@gmail.com` |
| 5 | #2082 | `7a4c18ff97b9771695485672b409582a37b20aac` | `df8223e0008b6d5bf82d0433592b8f1557c97311` | Smeet Agrawal | `smeetagrawal23@gmail.com` |
| 6 | #2307 | `d68ecb2c2cc6978673351f98b971c77827922204` | `0d1ae3ea1087bebcf5ce45982a770d0f4925af6c` | Jeff | `138259519+jdotc1@users.noreply.github.com` |
| 7 | #2329 | `2cf0ade305f4a60049b0c57af466dee4744e04ad` | `21028d8858e31fd088029f2fdb6634c99a1d3563` | cardene777 | `cardene777@gmail.com` |
| 8 | #2284 | `f146e235258cf21d695f74ba1791c2033859f721` | `1894e6f2cdaef94fe66388e577e0f59f14bfcff2` | Sakariyah Abdulhazeem | `sakariyahabdulhazeem@gmail.com` |
| 9 | #2246 | `5bed9587521b429aecb3b223acf0d20dbeef754f` | `9b7b2243fcdeb4279cccd9df9348f9e537d4fb95` | INTFRAME | `hello@intframe.com` |
| 10 | #2338 | `16ba49e717e24e5bea47dc7a2419033ebbe6bcde` | `7fed4140d87f62090c87b9f09c5c9b6073a84a27` | KVR Vishnu | `46535652+VISHNU0906@users.noreply.github.com` |
| 11 | #2339 | `b603fa0bd8a563909c062d341f2f1ea5781e2142` | `7e5132faa1ccea4bfd0b4b006ca896b5dd03df3f` | DelaySystem | `DelaySystem@tubot.lat` |
| 12 | #2080 | `2bd5ad38561fd54ff92975f2c38344d56f935fc4` | `ff77f9f58b1ec5fbbece852e0bfe9c1f050f24e5` | Smeet Agrawal | `smeetagrawal23@gmail.com` |
| 13 | #2099 | `62a89d4339a344e52e209a7c6d260de711b34518` | `1b3896fd7df9a667a09568b751dce6143613319e` | Nikhil | `tad.areas_0y@icloud.com` |
| 14 | #2317 | `7ce04bb344f4ddb7882167f9dcfce73598900c77` | `71e91200afb2415b283bb1721709b19c2fddd8aa` | Varshith Puli | `pulivarshit@gmail.com` |
| 15 | #2248 | `8ee08eff98e602ea2a1e301bd2e6ab575b0da529` | `5c9ca86f8f4d1db1e17004b6a1b598137203face` | hamodywe | `iosapk.org@gmail.com` |
| 16 | #2249 | `5e353588a898b56abae78018c4f13c10f59122f0` | `509fce15ad2eb2536bc8d840a3ed3e374d73bb2d` | Vaibhav Srivastava | `vaibhavsri1712@gmail.com` |
| 17 | #2253 | `260d7e18947668fa47b660f8d9fa242c62191454` | `4f58f98a1424203af39cc00a26f9672e9b45b794` | Arpit Jain | `arpitjain099@gmail.com` |
| 18 | #2312 | `6445d1e949d1ff49caa7a6db86340d66ea4c8973` | `2009273e0029f91aeda4732552ae5489634d0749` | Varshith Puli | `pulivarshit@gmail.com` |

### Ancestry verification

- `5f6ec56` is an ancestor of `HEAD`.
- `78b5c7f` (the original audit branch) is an ancestor of `HEAD`.
- All 18 integration commits are ancestors of `HEAD`.
- Each upstream commit contains the exact `(cherry picked from commit <40-hex>)` reference shown above.
- All 18 referenced source commits are present in the object database.
- Author names and emails match their source commits exactly.
- No merge commits were introduced among the contributor commits.
- The 18 commits cover all 17 selected PRs: #2227 supplies two commits; the other 16 PRs supply one each.

The RNC integration adjustments `22f63b4273566228012126072806f86530436456` and `100ae9064fa98728480a2818c8a157889b9e1c2e` also remain ancestral and unrewritten. The three audit-documentation commits (`110d6fd`, `78b5c7f`, `baacad5`) and the report commit follow them.

Ancestry is therefore:

```text
master (5f6ec56)
  -> 18 contributor commits
  -> 2 RNC integration adjustments
  -> contributor audit documentation
  -> encoding correction (baacad5)
  -> RNC V10 integration (9c82fe4)
  -> this report
```

No upstream contributor commit was amended, squashed, replayed, or dropped.

## D. Test results

### Environment

- Windows 11 host, Go `go1.27.3 windows/amd64`
- Docker Engine `29.9.3`
- Linux race container: `golang:1.27.4` (`linux/amd64`)
- No native C compiler (`gcc`) on the Windows host, so local `-race` was unavailable.
- Test outputs and scratch fixtures were written under `%TEMP%`; no generated scan output was committed.

### Required gates

| Command | Result |
|---|---|
| `go build ./...` | PASS |
| `go test -count=1 ./...` | PASS |
| `go test -count=1 -tags gore2regex ./...` | PASS |
| `go vet ./...` | PASS |
| `go mod verify` | PASS (`all modules verified`) |
| `git diff --check` | PASS |
| `go test -race -count=1 ./...` (Linux, non-root) | PASS |

All standard and RE2 package tests passed on Windows, including the `_windows_test.go` cases. The prior CP-V10 base limitations (`detect/detect_test.go:2084` vet failure and the pre-V10 generator type mismatch) are absent because this tree is the validated V10 tree.

### Targeted regressions

| Target | Command | Result |
|---|---|---|
| Redaction and UTF-8 masking | `go test -count=1 ./report -run 'Test(Redact\|Mask)'` | PASS |
| Regex and AND allowlists | `go test -count=1 ./config -run 'Test(InvalidUserRegexReturnsError\|TranslateAllowlists\|Translate$)'` | PASS |
| Git, archive, file, root handling | `go test -count=1 ./sources -run 'Test(ShouldSkipPath\|Git\|Staged\|File\|Archive\|Files\|Root)'` | PASS |
| Error collection and staged/root detection | `go test -count=1 ./detect -run 'Test(DetectSource\|Staged\|DirRoot)'` | PASS |
| CLI, report publication, status, rename | `go test -count=1 ./cmd -run 'Test(CLI\|Report\|Resolve\|NewReport\|Aggregate\|Numbered\|Partial\|Informative\|Rename\|IsTransient\|Scan)'` | PASS |
| Error collector | `go test -count=1 ./internal/scanerr` | PASS |

These cover redaction and UTF-8, AND allowlists, Git parsing/stderr/exit codes/lifecycle, staged index blobs, archive verification and size limits, root-relative paths and fingerprints, scan error collection, report preservation/publication/permissions, and Windows-specific tests through the native suite.

### Race tests

`go test -race -count=1 ./...` was executed as uid 10001 in a Linux `golang:1.27.4` container from a Git bundle clone of candidate `9c82fe4`, with a temporary Git trust configuration for the fixture repositories. Result: **PASS**. All test-bearing packages passed, including `cmd`, `config`, `detect`, `detect/codec`, `internal/informative`, `internal/scanerr`, `report`, and `sources`.

A preliminary attempt used a `git archive` extracted on Windows. It failed only in the `detect` fixture tests because the extraction lost the fixture repositories' Git metadata/symlink behavior (`fatal: bad object` and a symlink assertion mismatch). Re-running from a proper Linux clone eliminated the environment artifact and passed the full suite. The failed attempt is not counted as a product failure.

Native Windows race tests were **not run** because the host has no C compiler. Darwin, ARM64, and Docker image builds were **not rerun** in this checkpoint; their configuration is blob-identical to V10, and no claim is made that those environments were independently executed here.

### V10 reference runtime comparison

Two independently built binaries were produced from the reconciled candidate and from `fix/validation` at V10-backup:

| Artifact | SHA-256 |
|---|---|
| Candidate binary | `1b4b2a43e3e1383a2289b45c104c91fe1f4c2c553d3559ced613d15b124e9c19` |
| V10 reference binary | `1b4b2a43e3e1383a2289b45c104c91fe1f4c2c553d3559ced613d15b124e9c19` |

The binaries are byte-identical. Both were run against the same synthetic fixture with `gitleaks dir`:

| Format | Exit code | Bytes | SHA-256 | Byte-identical |
|---|---|---|---|---|
| JSON | 1 | 705 | `e610693ec47d5dcbcb65a9578e30c58f45dcf588bda723a00c10f54dae96d3d7` | Yes |
| SARIF | 1 | 49,830 | `cb649df5670e2037d512a014a1c0bb4152c9bb4cec7fb4a7c476395661d8392d` | Yes |
| CSV | 1 | 428 | `cb4bd892385006b61eea357651038397590530d34d73353260e82c3ff74b4245` | Yes |

Both binaries reported the same finding: `slack-bot-token` at `secrets.txt:3`, fingerprint `.../secrets.txt:slack-bot-token:3`, with the same secret hash and finding fields. Report structures and exit codes are identical.

### Scanner harness

The existing scanner harness lives on the separate `origin/test/scan-harness` branch. Its local worktree is not a usable checkout (all tracked files are currently deleted in that worktree), and no pinned fixture cache is present. It was therefore not run in this checkpoint. The binary and report comparison above is the direct V10 reference comparison required by the checkpoint; no harness infrastructure or oracle was created or modified.

### Remaining exceptions

- No source/product exception remains.
- GitHub PR #3 itself remains `CONFLICTING` because its head (`fix/validation`) was not rewritten and its base was reconciled by squash rather than by merging the PR. This is expected and is addressed in Section E.
- Platform coverage not executed here (native Windows race, Darwin, ARM64, Docker image build) is explicitly not claimed.

## E. Promotion assessment

The reconciled candidate is ready to become the long-lived `rnc-mod` development baseline **after operator review**. It has:

- preserved all 18 contributor commits and all 17 PR provenance references;
- an exact V10 product tree, object-for-object, outside the approved provenance reports;
- passing build, standard tests, RE2 tests, vet, module verification, whitespace check, targeted regressions, and the non-root Linux race suite;
- no changes to `master`, `fix/validation`, `V10-backup`, or `rnc-mod`.

Promotion in this checkpoint is limited to fast-forwarding `audit/rnc-mod-upstream-integration` to the report-bearing reconciliation tip. `rnc-mod` is not promoted.

PR #3 should be **closed as superseded**, not merged through the GitHub UI. Its head branch still contains the 64 unsquashed V10 commits; using the Merge button now would create a merge/squash commit, duplicate the reconciled content, and obscure the approved ancestry. The reconciliation commit and this report record PR #3 as the source reference.

## F. Documentation correction

Commit `baacad5` repaired the five provenance reports only. The previous normalization commit, `78b5c7f`, had interpreted UTF-8 bytes as Windows-1252 while normalizing Markdown, producing visible sequences such as `a-circumflex euro quote` in place of em dashes and smart quotes. The repair restores the intended Unicode and retains the trailing-whitespace normalization, without changing findings, evidence, conclusions, headings, or tables.

Verification: each restored report exactly matches the pre-corruption `110d6fd` text after the same trailing-whitespace normalization; all five contain zero U+FFFD and zero U+009D characters, and no mojibake sequences remain.

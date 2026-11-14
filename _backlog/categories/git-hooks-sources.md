# Git Hooks Sources — RNC review queue

**Current active queue (2026-11-13):** Completed/declined original issue and PR packages were removed from this category. See the [RNC disposition register](exclusions.md) and [fixes/QOL implementation ledger](https://github.com/RNC-Solutions-Group/gitleaks/blob/rnc-mod/RNC-MOD.md). Remaining original candidates below are not closed by overlap with those corrections.

Second-pass evidence review: 2026-11-09 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**14 issues · 12 PRs · 15 work packages.**

**2026-11-12 routing:** Elective work packages moved to [Optional Features](features.md), preserving issues, PRs, contributor credit and original review evidence. No source changes or scanner-defect closures are implied.

Git history, hooks, staged changes, diff/log parsing, source integration and failure propagation. The original Git/Hooks review inspected 20 issue bodies and 13 PR diffs; this Scan Integrity follow-up incorporated issue #1464 and PRs #2253/#2252 after additional source/patch comparison. Three older issues and one PR remain in [evidence-backed exclusions](exclusions.md). **No Go tests or live Git-hook reproductions were executed.**

**Review states:** `SOURCE-CHECKED` includes a targeted fork-source comparison; `PATCH-REVIEWED` means an upstream diff was inspected, **not** that the patch was built or tested against the fork. `PENDING` records unresolved issues or competing proposals. No PR is adopted or upstream item closed by categorization. See [RNC exclusions](exclusions.md) for retired review slots.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| [#2003](https://github.com/gitleaks/gitleaks/issues/2003) | [#2128](https://github.com/gitleaks/gitleaks/pull/2128) (@puneetdixit200) | Quoted Git `--log-opts` argument handling versus claimed notes/CSV symptom. | **PATCH RELEVANT ONLY TO QUOTING — ISSUE CLAIM UNPROVEN** | Master `strings.Split(logOpts," ")` breaks quoted multi-word pathspecs; #2128 adds a quote/escape parser and unit+Git-note tests. But user's `--log-opts "--no-notes"` shell command already supplies plain `--no-notes`; Git `--no-notes` hides annotations, **not** commit messages, and CSV permits quoted line breaks. PR test explicitly preserves commit messages. Verify Windows quoting, empty values, backslashes and unsupported flags; do not claim issue #2003 or CSV 'fix' resolved. | — |
| [#1028](https://github.com/gitleaks/gitleaks/issues/1028) | [#2301](https://github.com/gitleaks/gitleaks/pull/2301) (@bhaskargurram-ai) | Merge-commit-only secrets omitted by `git log -p`. | **F11 REPRODUCED — COA 3 RECOMMENDED, NOT APPROVED/IMPLEMENTED** | #2301 adds `--diff-merges=first-parent` before user log opts and an actual evil-merge fixture; `--cc` emits combined diffs the current `gitdiff.Parse` cannot read. First-parent also re-reports branch additions on merges, changing fingerprint/ignore semantics; test octopus/fast-forward/non-FF merges, first-parent paths, duplicates, `--no-diff-merges` override, Git-version and Windows compatibility (upstream test skips Windows). | —  **2026-11-12 current-RNC disposition:** Default 5/4 misses merge-only secrets; first-parent 10/6 plus 3 duplicate findings; remerge 7/6 has conflict path and octopus defects. Recommend explicit `--log-opts` first-parent and separate harness lane; PR #2301 @bhaskargurram-ai reviewed, not adopted. |
| [#1993](https://github.com/gitleaks/gitleaks/issues/1993) | — | Staged scan hangs in PTY/pipe-driven macOS developer terminals. | **REPRO BLOCKED — KEEP** | Reporter observed `git diff --staged \| head` hang but `bash -c` and file redirection work; pinned master retains deprecated `protect` and pipes Git stdout/stderr through goroutines. A 0%-CPU observation does not establish the exact deadlock. Instrument child process, stdout/stderr draining, cancellation and timeout in macOS VS Code/Cursor, then distinguish Git process hang from scanner hang; do not recommend bypassing hooks as a fix. | — |
| [#1647](https://github.com/gitleaks/gitleaks/issues/1647) | — | Kubernetes Secret keys with changed diff-context/base64 scope. | **OPEN — DETECTION CROSS-CHECK** | Reproducer adds base64 `c29tZXRoaW5nZ2VoZWltZXMK`, which decodes to ordinary short text, not a verifiable credential. Master `git log -U0` plus `TextFragment.Raw(OpAdd)` discards unchanged multiline YAML context. Upstream discussion proposes requiring finding span overlap added lines. Reproduce with a valid K8s rule/secret fixture, single-line versus multi-line changes, decoding depth and `git` versus `dir`; avoid claiming generic Git reader failure. | — |
| [#1364](https://github.com/gitleaks/gitleaks/issues/1364) | — | Docker pre-commit Git `--staged` reported unknown on macOS. | **OPEN — CONTAINER TRUST/WORKDIR** | Reporter's error is Git's `--no-index` usage, not evidence `--staged` is unsupported; comments reproduce `safe.directory` dubious ownership and container UID/mount mismatch, with Mac/WSL variance. Current Docker hook still runs staged diff. Test Docker Desktop/OrbStack/WSL Linux with `git status` inside identical mounted container, `.git` visibility, ownership and Git version; avoid globally trusting arbitrary paths. | — |
| [#1310](https://github.com/gitleaks/gitleaks/issues/1310) | — | Shallow clone gives misleading commit/fingerprint identity. | **OPEN — HISTORY-BOUNDARY TEST** | v8.18.1 `--depth 1` reproducer reports new clone-tip commit SHA against old README secrets, breaking `.gitleaksignore`. Current parser uses `git log -p -U0` and SHA-based fingerprints; no baseline proof of repair. Compare full and shallow scans, root snapshot diff semantics, merge commits and `--unshallow` behavior; distinguish missing provenance from actual newly introduced secret. | — |
| [#821](https://github.com/gitleaks/gitleaks/issues/821) | — | GitLab pre-receive sees existing refs, not incoming quarantine commits. | **OPEN — INTEGRATION DESIGN** | Upstream maintainer and users identify revision selection/quarantine as cause; a working external hook computes incoming revisions and passes `--log-opts`. Master defaults to `git log --all`, which need not include incoming objects. Test GitLab bare/quarantine hook with stdin old/new/ref triples, forced pushes, new branches, deletion, ancestry and safely failing closed; no upstream plugin assumed. | — |
| [#1293](https://github.com/gitleaks/gitleaks/issues/1293) | — | Auto-load `.gitleaksignore` stored in a bare repository tree. | **OPEN — EXPECTATION/DESIGN** | `cmd/root.go` finds ignore files on the **filesystem** at `source/.gitleaksignore`; a bare repo exposes Git objects/refs, not a checked-out root file. An explicit `--gitleaks-ignore-path` or `-c` workaround is not auto-discovery. Decide whether hook should read `.gitleaksignore` at a specified revision, respect trust boundaries and handle missing/changed refs; test cloned versus bare. | — |
| [#1163](https://github.com/gitleaks/gitleaks/issues/1163) | — | Git `--source=<subfolder>` scans repository-wide history. | **OPEN — SCOPE CONTRACT** | Master `git -C <source> log -p ...` selects Git working directory, **not a path filter**; issue's subdir leaks are expected under current Git behavior, while `dir` scopes filesystem files. Upstream documents `--log-opts='--full-history --all -- :/modules/module30/*'` workaround. Decide explicit subdir-scoped flag versus current repository-wide contract; test nested repos, relative/absolute paths and user pathspec overrides. | — |
| — | [#1809](https://github.com/gitleaks/gitleaks/pull/1809) (@JoostVoskuil) | Azure DevOps SSH remotes fail finding-link platform detection. | **REAL GAP — ORIGINAL PATCH OBSOLETE** | #1809 edits removed `cmd/detect.go` and `detect/git.go`; current path is `sources/git.go` + `cmd/scm`. Current SSH normalization can yield `ssh.dev.azure.com` but `platformFromHost` only recognizes `dev.azure.com`/`visualstudio.com`. Old PR used substring host matching (spoof risk), no tests and reviewer requested rebase after #1872. Re-implement narrowly with exact trusted hosts and Azure SSH v3 path mapping; preserve GitHub/GitLab/explicit platform; never cherry-pick old diff. | — |

## Cross-cutting implementation conclusions

1. **P0 original CLI fail-open paths corrected in RNC CP-V3.** The #1981/#2129 packages and their alternative PRs are now in [dispositions](exclusions.md). Continue separate malformed/custom diff and blocking-I/O regression tests; a new false-clean must be demonstrated before reopening.
2. **P0: keep allowlist Boolean semantics.** [PR #2318](https://github.com/gitleaks/gitleaks/pull/2318) is **RNC rejected**, not selected. `detect.checkCommitOrPathAllowed` already evaluates global commit/path allowlists, including `AND`. The PR's new unconditional `continue` would discard a file if its commit matches but its required path does **not**, creating false negatives. The original inner-loop `continue` in `sources/git.go` was corrected in RNC CP-V11 (PR #4); preserve detector and Git-history regressions. The rejected PR #2318 remains separate. See the exclusions ledger.
3. **P1 remaining Git output hazards.** Forced color, pretty/date, and textconv hazards were corrected in V12. Keep #2301 merge diffs separately opt-out capable, with duplicate-finding policy and regression tests; #2003 argument parsing remains active for V13.
4. **System-hook filename forwarding corrected CP-V9 and exercised in CP-V10 CI.** Hook remains staged-index scanning, even when the pre-commit runner was asked for `--files` or `--all-files`. The #1918 directory-scan alternative changes scope, speed and index/worktree behavior and remains separate.
5. **No wholesale import of old PRs.** #1899 preserves a non-1-error fail-open hook, #1809 targets deleted files and changes trusted-host equality to unsafe substring matching, and #2128 does not prove the actual #2003 CSV complaint. Preserve upstream author credit, but require a corrected integrated design.

## Test gates for implementation

| Gate | Mandatory fixtures / assertions |
| --- | --- |
| Git process completeness | Nonrepo, invalid `--log-opts`, corrupted index, fatal `safe.directory`, stderr warning + valid fragments, nonzero exit with **empty** stderr, timeout and cancellation. No partial success returned as clean; no child process leaked. |
| Diff fidelity | Forced ANSI color, `format.pretty`, `log.date=relative`, external diff, SOPS textconv, merge-only secret, shallow clone, added/deleted lines and pathspec-scoped directories. Confirm actual finding + SHA/line, not only parser file count. |
| Hook contract | Go, Docker and system hook; multiple staged filenames; `--files`, `--all-files`, unstaged edits, deleted files, no staged changes, missing executable, every nonzero scanner exit, container UID and safe.directory. |
| Release safety | Linux, Windows, macOS as applicable; Go tests and race tests; performance delta on large repository; secret redaction, findings/report consistency, timeout cleanup; no new default trust-bypass or network scan behavior. |

## Provenance — inspected upstream heads

- [#2319](https://github.com/gitleaks/gitleaks/pull/2319) @abhinavrastogi2627: `70cec8c274b7`
- [#2080](https://github.com/gitleaks/gitleaks/pull/2080) @Smeet23: `ff77f9f58b1e`
- [#2292](https://github.com/gitleaks/gitleaks/pull/2292) @Saadanjum0: `de6f59a0154a`
- [#2339](https://github.com/gitleaks/gitleaks/pull/2339) @delaysystem: `7e5132faa1cc`
- [#2128](https://github.com/gitleaks/gitleaks/pull/2128) @puneetdixit200: `9190d351d019`
- [#2311](https://github.com/gitleaks/gitleaks/pull/2311) @tonymet: `ae702afb4b1c`
- [#2301](https://github.com/gitleaks/gitleaks/pull/2301) @bhaskargurram-ai: `6bcf2f91c1b8`
- [#2284](https://github.com/gitleaks/gitleaks/pull/2284) @lekanjava: `1894e6f2cdae`
- [#2203](https://github.com/gitleaks/gitleaks/pull/2203) @RajeshRajendiran: `221ca6296bba`
- [#1918](https://github.com/gitleaks/gitleaks/pull/1918) @Holzhaus: `928e0f1239a2`
- [#1899](https://github.com/gitleaks/gitleaks/pull/1899) @vignesh1507: `b90d5518137e`
- [#1809](https://github.com/gitleaks/gitleaks/pull/1809) @JoostVoskuil: `10704c9a086f`

This list identifies **reviewed upstream heads**, not RNC integrations. [Rejected #2318](https://github.com/gitleaks/gitleaks/pull/2318) @MAHIR-BABBAR (`2f4877db1a86`) is tracked only in [exclusions](exclusions.md).

**Scan Integrity crosswalk (2026-11-09):** Issue [#1464](https://github.com/gitleaks/gitleaks/issues/1464) and PRs [#2253](https://github.com/gitleaks/gitleaks/pull/2253), [#2252](https://github.com/gitleaks/gitleaks/pull/2252) moved here from the scan-reader queue. They do **not** count as separate work packages because the core fault—fragment errors swallowed before CLI exit—is already owned by #1981; all original authors and candidate strategies remain in this row, not excluded or adopted. The original #2129 queued-fragment truncation acceptance remains a separate package.

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

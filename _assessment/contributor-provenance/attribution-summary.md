# Attribution Summary

**Branch:** `audit/rnc-mod-upstream-integration`
**Base:** `5f6ec56`
**CP-V10 reference:** `bd42a95`
**Audited PRs:** 22 (11 primary, 11 secondary)
**Distinct PR contributors:** 18

## Category 1 — Direct/substantive implementation reuse

Recommend preserving the original commit-level authorship and ancestry for these
contributions.

| PR | Contributor | Commit(s) | What survives |
| ---: | --- | --- | --- |
| #2261 | @arpitjain099 | `b1ce07f` | Rune-based masking algorithm and UTF-8 tests |
| #2082 | @Smeet23 | `df8223e` | Error-returning regex wrappers and compile-at-config sites |
| #2317 | @Varshith-Kali | `71e9120` | 64 KB archive content verification and prefix replay |
| #2334 | @Shanmukh-l33t | `ad321b7` | `resolveDirSource`, `MaximumNArgs(1)` and regression cases |
| #2227 | @BasitS-hash | `894b4a3`, `49b26ee` | `pathIsSufficientToSkip` and AND handling; both commits preserved |
| #2329 | @cardene777 | `21028d8` | First-chunk-only MIME classification |
| #2246 | @aron-intframe | `9b7b224` | `newGitCmd`, parser drain, unparsed-byte detection |
| #2248 | @hamodywe | `5c9ca86` | `Root`/`RootRelativePath` and root-relative matching design |

These are not claims that the final V10 file is byte-identical. They are
substantive code/design contributions with identifiable surviving expression.

## Category 2 — Adapted implementation

Substantial influence or reused expression requiring reconciliation and
contributor recognition. The original commit is still preserved in the
experimental integration branch because the implementation is genuinely
represented there.

| PR | Contributor | Commit(s) | Adaptation |
| ---: | --- | --- | --- |
| #2338 | @VISHNU0906 | `7fed414` | Staged binary fallback; V10 uses explicit `:0:` index lookup and MIME/error handling |
| #2099 | @alephnull-sh | `1b3896f` | Size propagation survives; whole-buffer decompression replaced by bounded spooling |
| #2284 | @lekanjava | `1894e6f` | `--no-color` behavior survives; V10 adds `forceNoColor` argument handling |
| #2307 | @jdotc1 | `0d1ae3e` | Destructive preflight removed; V10 replaces the publication architecture |
| #2339 | @delaysystem | `7e5132f` | Git exit error handling merged with parser drain and V10 collector |
| #2080 | @Smeet23 | `ff77f9f` | Stderr error return; V10 classifies fatal/warning diagnostics |
| #2253 | @arpitjain099 | `4f58f98` | Fragment error collection; V10 uses bounded `scanerr` |
| #2249 | @vaibhav8a | `509fce1` | Unreadable-file accumulation; V10 uses `scanerr` and scanFile lifecycle |
| #2312 | @Varshith-Kali | `2009273` | 0600 `openReportFile` retained; destructive preflight hunk dropped in integration; V10 uses temp 0600 |

## Category 3 — Research/design influence

Documentary acknowledgment recommended; no fabricated code authorship.

| PR | Contributor | Why |
| ---: | --- | --- |
| #2300 | @ToniAdreal | Alternative 0600 report design; preflight still destructive; superseded by selected #2312/#2307 and V10 temp publication |
| #2292 | @Saadanjum0 | Alternative stderr propagation; conflicts with selected #2080 |
| #2252 | @CAOShurong | Immediate fragment-error return; rejected because it abandons readable findings |
| #2313 | @Varshith-Kali | Alternative unreadable-file error handling; #2249 selected |
| #2203 | @RajeshRajendiran | Alternative errCh/drain fix; #2246 selected; V10 drain/section logic supersedes |
| #2206 | @trinathone | Non-root Docker proposal reviewed; V10 optional target is RNC implementation; root default retained |

## Category 4 — Documentation/dependencies

Appropriate acknowledgment in documentation and the ledger; no scanner-history
commit.

| Source | Contributor | Recognition |
| --- | --- | --- |
| #2236, #2212 | @RajeshRajendiran | Git added-line and decode-depth documentation |
| #2298 | @bharathvemula | Baseline prose/examples |
| #2287 | @ump45nose | Pre-commit Docker anchor correction |
| #2282 | @haimingZZ | README/CONTRIBUTING typo corrections |
| #2258 | @jankesec | Configuration example correction |
| #2247 | @CAOShurong | Configuration guide URL |
| #2008 | @jylenhof | mise installation instructions |
| #1263 | @noblehelm | Baseline credential-storage guidance |
| #2319 | @abhinavrastogi2627 | One-line system-hook manifest correction; RNC source commit |
| Dependency graph | @corani, Dependabot | Profile-only recognition for independently matched versions; no cherry-pick |

## Category 5 — Superseded or unrelated

No additional recognition beyond the historical records and the documentary
credits above:

- Alternative implementations listed in Category 3.
- Dependency/version proposals superseded by newer independently selected pins.
- The upstream Gitleaks Action integration, which RNC explicitly declines.
- Any proposal whose production behavior is fully replaced by V10 and which
  has no identifiable surviving expression beyond the problem statement.

## Contributor-level roll-up

| Contributor | Primary | Secondary | Documentation/deps | Recognition |
| --- | --- | --- | --- | --- |
| @arpitjain099 | #2261 | #2253 | — | Commit-level for #2261; commit-level adapted for #2253 |
| @Smeet23 | #2082 | #2080 | — | Commit-level for both; #2080 adapted |
| @Varshith-Kali | #2317 | #2312, #2313 | — | Commit-level for #2317 and #2312; #2313 documentary |
| @VISHNU0906 | #2338 | — | — | Commit-level adapted |
| @alephnull-sh | #2099 | — | — | Commit-level adapted |
| @aron-intframe | #2246 | — | — | Commit-level |
| @hamodywe | #2248 | — | — | Commit-level adapted |
| @Shanmukh-l33t | #2334 | — | — | Commit-level |
| @BasitS-hash | #2227 | — | — | Commit-level, two commits |
| @lekanjava | #2284 | — | — | Commit-level adapted |
| @cardene777 | #2329 | — | — | Commit-level |
| @jdotc1 | — | #2307 | — | Commit-level adapted |
| @ToniAdreal | — | #2300 | — | Documentary |
| @delaysystem | — | #2339 | — | Commit-level adapted |
| @Saadanjum0 | — | #2292 | — | Documentary |
| @CAOShurong | — | #2252 | #2247 | Documentary/adapted acknowledgment |
| @vaibhav8a | — | #2249 | — | Commit-level adapted |
| @RajeshRajendiran | — | #2203 | #2236, #2212 | Documentary/code influence; existing CP-V9 credit |
| @bharathvemula | — | — | #2298 | Documentary |
| @haimingZZ | — | — | #2282 | Documentary |
| @jankesec | — | — | #2258 | Documentary |
| @ump45nose | — | — | #2287 | Documentary |
| @jylenhof | — | — | #2008 | Documentary |
| @noblehelm | — | — | #1263 | Documentary |
| @abhinavrastogi2627 | — | — | #2319 | Documentary |
| @trinathone | — | #2206 | — | Documentary |
| @corani / Dependabot | — | — | Dependency graph | Profile-only |

## Recommended disposition

Preserve the 18 upstream commits from 17 PRs as separate ancestry. Do not
squash them into the RNC V10 commit. Apply the RNC V10 implementation afterward
as one RNC-authored reconciliation commit (or an operator-approved squash),
resolving conflicts in favor of validated V10 behavior. This preserves both the
contributors' original Git authorship and the RNC implementation history.

The final history should therefore contain, after `5f6ec56`:

1. the selected upstream commits (18 commits, original authors preserved);
2. any justified RNC integration adjustments;
3. the RNC V10 reconciliation commit(s);
4. no V10 squash performed in this audit.

This is the only structure that preserves commit-level credit while allowing
the validated V10 tree to become `rnc-mod`.

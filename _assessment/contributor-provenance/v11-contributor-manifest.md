# CP-V11 Contributor Manifest

**Status:** Final integration candidate
**Integration branch:** `audit/v11-integration`
**Integration base:** `rnc-mod` at `ee1e59e38231b84a8277131827a09fc929383ecc`
**V11 development HEAD:** `1b06b86a2ad815917b8e3f6265afb304f0eedd7e`
**RNC reconciliation commit:** `45f6155a56b2c8c8fad1272cf1f38f5602da8e45`
**V11 product tree:** `73cfdd5c02943d97a731e0c7012380f9c4d00958`

## Selected upstream contributions

All six original commits were retained. No upstream implementation was omitted.

| Task | PR | Source commit | Integration commit | Author | Scope |
|---|---|---|---|---|---|
| V11-C | #2166 | `55099ffe9d039c6fb0500a5f65cca1abfb398826` | `c13eb34551fb4c2fe2917a0a866dd222e0ac1ab8` | ashvinctrl `<sharmaashvin27@gmail.com>` | Global commit allowlist outer-loop correction |
| V11-C | #2166 | `fcd84a613a099e698aa811b32e6d3494c6af0326` | `9c60717293ad113f09e24dbd29e7fcee7af18333` | ashvinctrl `<sharmaashvin27@gmail.com>` | `isCommitAllowed` helper and unit tests |
| V11-B | #2294 | `b04e70b06a62e95ae09b68eafeb9b7a2050153e8` | `079ee1cb39aa50bb39f82c27eb7f2f7f8d0fbddc` | Paradox `<159635503+i-am-paradox@users.noreply.github.com>` | Anchored boolean/null allowlist regex, generator and config |
| V11-A | #2346 | `126bab3a4497026961f7e78d373b098fb845dc92` | `f6c111c8e701b2129381535b638bf130a0de9002` | DelaySystem `<DelaySystem@tubot.lat>` | Exclude previous newline from `Finding.Line` and decoded lines |
| V11-A | #2346 | `5379957823eab9f6d4ed31827cc044ba00797e35` | `05d58c560e84b29552ae4810f0a930e05e6266c5` | DelaySystem `<DelaySystem@tubot.lat>` | Clamp empty decoded line ranges |
| V11-A | #2346 | `760252fd84604b2a6179b060241237548fca6804` | `29a6f766c051a5cbd60aba2d96f3e1410b8a8e22` | DelaySystem `<DelaySystem@tubot.lat>` | Clamp `startLineIndex` to `endLineIndex` |

The source PRs were open upstream at review time; the listed source commits are the PR head commits, not upstream merge commits.

## Provenance checks

- All six local integration commits retain their source author name and email.
- All six carry the exact `(cherry picked from commit <40-hex>)` trailer.
- The commits are applied in dependency order on top of the approved `rnc-mod` baseline.
- The commits were not amended, rebased, squashed, or replayed during review.
- The RNC development checkpoint commits were not transferred into the integration branch.
- One cherry-pick conflict occurred in `sources/git_test.go` for `fcd84a613a099e698aa811b32e6d3494c6af0326`. It was resolved by retaining the existing V10 Git lifecycle tests and appending the upstream `TestIsCommitAllowed` unit test. The production helper auto-merged.
- One reconciliation conflict occurred in `detect/detect_test.go`; it was resolved by restoring the authoritative V11 development tree, not by preserving an intermediate merged test file.

## RNC reconciliation

The single RNC-authored commit `45f6155a56b2c8c8fad1272cf1f38f5602da8e45` supplies the final V11 product tree. The six upstream commits remain ancestral and individually identifiable; their implementations may be adapted or superseded by the final RNC tree.

The final product tree after reconciliation is identical to V11 development HEAD `1b06b86`: tree `73cfdd5c02943d97a731e0c7012380f9c4d00958`.

## Licensing and attribution

The repository retains the upstream MIT license and copyright notice. Source SHAs and `-x` trailers preserve traceability. No co-author trailers were added, and no contributor implementation is presented as original RNC work.

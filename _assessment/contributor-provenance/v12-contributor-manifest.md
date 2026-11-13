# CP-V12 Contributor Manifest

**Integration branch:** `audit/v12-integration`
**Integration base:** `rnc-mod` at `197a7e5516b0029d0b7bb2554552fc9ed54efb68`
**V12 development HEAD:** `1ff69b7836b32e4f738883a09bc2a3daeccbd69b`
**RNC reconciliation commit:** `9092703a25f4c10edc162699955a6993ec66ce37`
**V12 product tree:** `a348eea61670da139b89b410efa5156325091e87`

## Retained upstream commits

| Package | PR | Source commit | Integration commit | Author |
|---|---|---|---|---|
| V12-03 | #2311 | `ae702afb4b1cf41d560baf89a9f454eadfa251f2` | `2790da33f1250426c8e8901ea0cde6476410eee3` | Anthony Metzidis `<anthony.metzidis@gmail.com>` |
| V12-07 | #2237 | `981e18e2ba6dc3d5419047d026352fb4a85ab5b1` | `904871b791ab40fd92fc59917d7640123e1448d5` | Bunlong Heng `<bheng.code@gmail.com>` |
| V12-09 | #2233 | `b8fc92f16740f33715985a2373e05b8efa663fa0` | `7021b7d242c1089d20359420245de8def434b246` | CaoBochun `<caobochun@bytedance.com>` |
| V12-11 | #2271 | `158d1e0bc38c49322593401dabff5738afb6dc4b` | `8e787c2e8fd8adddf38d5109e1a943c1b74614a6` | Dhanush Satyavolu `<satyavoludhanush@gmail.com>` |
| V12-14 | #2324 | `a77d865f9eebbe38903c5fc365d1f0accfbff8d8` | `d1fe7bca8f34608358f4de2f7cb0f416fd979e57` | cardene777 `<cardene777@gmail.com>` |

All five integration commits retain their original author and the exact `(cherry picked from commit <40-hex>)` trailer. The source PRs were open upstream at review time.

## Considered but not cherry-picked

| Package/PR | Reason |
|---|---|
| V12-01 / #2301 | Default `--diff-merges=first-parent` changes findings, duplicates and fingerprints; V12 keeps the default and documents/tests the explicit opt-in. |
| V12-08 / #2238 | Shared byte-column replacement would change JSON/CSV and baselines. V12 uses SARIF-only codepoint fields; the upstream design is credited here, not represented as adopted code. |
| V12-15 / #1918 | The proposed directory scan changes scope and fails with multiple filenames. V12 adds a bounded `gitleaks files` command and explicit-file hooks instead. |
| V12-04 / #2246 | The stream-lifecycle defect was already corrected in V10/V11 ancestry; V12 extends validation of unsupported output formats. |

## Provenance checks

- Five original contributor commits remain ancestral and identifiable.
- Authors are preserved; no co-author trailer was fabricated.
- RNC V12 development checkpoint commits are not present in the integration history.
- Five squash conflicts (`report/csv.go`, `report/sarif.go`, `sources/common.go`, `sources/common_test.go`, `sources/git.go`) were resolved by restoring the accepted V12 development tree.
- The final product tree matches V12 development HEAD exactly: `a348eea61670da139b89b410efa5156325091e87`.
- Only the V11/V12 provenance documentation differs from the raw product tree.

## Reconciliation

The single RNC-authored commit `9092703a25f4c10edc162699955a6993ec66ce37` supplies the final V12 implementation. Contributor code may be adapted or superseded by that tree while authorship and source references remain preserved.

MIT license and copyright notices remain in place.

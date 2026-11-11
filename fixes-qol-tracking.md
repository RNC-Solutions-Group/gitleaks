# Fixes, quality-of-life improvements, and upstream attribution

> **Living engineering and release-preparation ledger.** This is not a claim that every item below is complete or that an upstream author's code has been incorporated. Update it as implementation, verification, and attribution decisions change.

**Project:** RNC Gitleaks fork (`RNC-Solutions-Group/gitleaks`)  
**Working branch:** `fix/validation`  
**Snapshot:** 2026-11-11; branch head before this document: `46488522cfa9c58dab9f31d7eb541443e56a642e`  
**Comparison baseline:** `dev/gitleaks-upgrade` at `c0516b4bd66c0dc3878d80283a267166f7b4fd09`  
**Upstream reference reviewed for CP-V0:** `gitleaks/gitleaks` `master` at `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`  
**Primary evidence:** CP-V0 assessment, CP-V1 completion report, branch diff/tests, and linked upstream issues/PRs. PR state and head commits are snapshots, not permanent integration pins.

## 1. Project contract: improve Gitleaks without changing how people use it

1. **Backward compatibility is a release gate.** Preserve upstream commands, flags, configuration interpretation, default scanning behavior, exit codes, finding content, output formats, and automation contracts unless correcting a demonstrated defect requires a change.
2. **Accuracy and completeness come first.** Never trade real detection coverage for speed, a superficially green test, or a nicer report.
3. **Safety fixes should be transparent.** Prevent unintended data loss, misleading clean results, partial reports, and unsafe permissions without imposing a new workflow on established users.
4. **New features are opt-in whenever practical.** Keep documented defaults compatible with original Gitleaks. Feature flags override configuration only when explicitly supplied.
5. **No unnecessary product expansion.** Gitleaks is a scanner, not a retention, archive, or report-management service. Do not mix in the separate benchmark harness or security orchestration toolkit.
6. **Preserve contributor credit.** Research influence, adapted code, and incorporated commits are different forms of contribution; record them honestly and distinctly.
7. **Keep the public history clean.** Iterative development commits are valuable on working branches but should not become hundreds of unrelated commits on the final fork baseline.

**Compatibility evidence required before release:** Run representative synthetic fixtures with upstream and the final fork using identical default settings. Compare detection results, report schemas, CLI contracts, exit codes, and basic runtime behavior. Document each intentional difference; do not assume compatibility merely because unit tests pass.

## 2. Improvement and change register

Statuses: **Implemented** = code exists on `fix/validation`; **Validated** = specific tests/evidence exist; **Planned** = approved design, not yet on branch; **Open** = not fixed; **Deferred** = outside current cycle. A local pass is not a claim that all CI platforms pass.

| ID | Area | Original defect or behavior | Change / expected contract | State at snapshot | Relevant evidence / attribution |
| --- | --- | --- | --- | --- | --- |
| V1-A / CPV0-F02 | Report integrity | Early `os.Create`/remove could destroy old output; final write could truncate or leave partial output. | Validate first; write temporary restrictive-permission report; publish atomically; preserve old output on failed validation/publication; reject unsafe targets. | **Implemented; local validation reported**, wider cross-platform acceptance pending. | [#2304](https://github.com/gitleaks/gitleaks/issues/2304), [#2307](https://github.com/gitleaks/gitleaks/pull/2307), [#2299](https://github.com/gitleaks/gitleaks/issues/2299), [#2300](https://github.com/gitleaks/gitleaks/pull/2300), [#2312](https://github.com/gitleaks/gitleaks/pull/2312) |
| V1-B / CPV0-F06 | AND allowlists | A path match could suppress a file before other required AND criteria were evaluated. | Only path-skip if OR or path-only AND; evaluate content/commit-dependent AND conditions during detection. | **Implemented; tests reported passing.** | [#2227](https://github.com/gitleaks/gitleaks/pull/2227) |
| V1-C / CPV0-F08 | Git diff integrity | Forced terminal coloring could break machine-readable Git log/diff parsing and silently miss findings. | Use Git `--no-color` for supported history/staged/unstaged acquisition; validate forced-color fixtures. | **Implemented; tests reported passing.** | [#2284](https://github.com/gitleaks/gitleaks/pull/2284) |
| V1-D / CPV0-F10 | File-chunk classification | Later chunks could reclassify an initially text file as binary and drop findings. | Classify MIME on the initial chunk only; preserve binary-at-start behavior. | **Implemented; tests reported passing.** | [#2329](https://github.com/gitleaks/gitleaks/pull/2329) |
| V1-E / CPV0-F12 | `dir` CLI | Extra/missing source arguments could silently fall back to scanning CWD. | Accept zero (default `.`) or one valid source; reject ambiguous/missing paths. | **Implemented; tests reported passing.** | [#1979](https://github.com/gitleaks/gitleaks/issues/1979), [#2334](https://github.com/gitleaks/gitleaks/pull/2334) |
| V1-F / CPV0-F13 | Pre-commit integration | Wrapper could interpret nonzero scanner exits as success, particularly on Windows. | Use portable subprocess status; fail on every nonzero code, with useful diagnostics. | **Implemented; tests reported passing.** | No matching upstream PR identified in CP-V0. |
| V1G | Path/collision hardening | Current development-only in-tree report policy adds path-resolution edge cases and cross-platform CI failures. | Correct relevant real safety bugs, then retire the superseded in-tree policy under V2. | **Partially implemented; CI not accepted.** | `cmd/report.go` and test changes through `4648852`; CI evidence must be refreshed. |
| V2-A | Output compatibility | Default upstream behavior overwrites the requested report. | `auto-overwrite = true` by default; keep stable original output path and normal overwrite when aggregation is disabled. | **Planned / approved contract.** | [#338](https://github.com/gitleaks/gitleaks/issues/338), [#2304](https://github.com/gitleaks/gitleaks/issues/2304) |
| V2-B | Optional aggregation | Repeated output paths collide. | `aggregation = false` by default; when enabled, create `report.json`, then `report_001.json`, etc. No content merge or retention. | **Planned / approved contract.** | [#338](https://github.com/gitleaks/gitleaks/issues/338) |
| V2-C | Explicit overrides | Users need one-run exceptions without editing configuration. | Add `--overwrite` and `--aggregate`; CLI wins over settings; both flags together error. | **Planned / approved contract.** | New RNC design; not an upstream flag adoption. |
| V2-D | Collision rejection | Users may prohibit both overwrite and aggregation. | If both settings false and destination exists, error **before scanning**; absent destination writes normally. | **Planned / approved contract.** | New RNC design; [#1065](https://github.com/gitleaks/gitleaks/issues/1065) discussed early path validation. |
| V2-E | Output simplicity | Experimental `--allow-in-tree-report-overwrite` adds special handling incompatible with the simplified policy. | Retire this RNC-only experimental flag while retaining scan-input/symlink/special-file protection. | **Planned.** | Current `cmd/root.go`, `cmd/report.go`; not part of original upstream interface. |
| V2-F | Publication correctness | Competing runs and OS-specific rename behavior can cause filename races. | Exclusive numbered destination allocation; atomic complete writes; real Windows/Linux/macOS validation; report actual output path. | **Planned.** | CP-V2 acceptance tests. |

**Output decision table (when the requested destination exists):**

| `auto-overwrite` | `aggregation` | Default outcome | One-run override |
| --- | --- | --- | --- |
| true | true | Aggregate (numbered file) | `--overwrite` |
| true | false | Overwrite requested file | `--aggregate` |
| false | true | Aggregate (numbered file) | `--overwrite` |
| false | false | Reject before scan | Either explicit flag |

Neither `--overwrite` nor `--no-overwrite` was found in the reviewed upstream v8 CLI. **Do not call either an original upstream flag.** Explicit overwrite must never authorize destroying a scan input. No report destination means no file-output policy change; stdout must remain compatible.

## 3. Known issue disposition

This section distinguishes an implemented fix from full release acceptance. Only the six V1 items below have been reported complete; **do not label the rest fixed based on design intent**.

| CP-V0 finding | Severity | Description | Disposition |
| --- | --- | --- | --- |
| F01 | P0 | Source/Git/worker errors can be swallowed, creating false clean scans or reports. | **OPEN**; separate error-propagation project. |
| F02 | P0 | Report validation/truncation, symlink, partial output, and input-collision data loss. | **V1 implemented**, V2 output-policy redesign pending; cross-platform acceptance outstanding. |
| F03 | P1 | Secret-bearing logs and incomplete redaction across multiple/empty/Unicode secrets. | **OPEN**; evaluate #2350/#2261 and broader log boundary. |
| F04 | P1 | Untrusted local configuration can replace rules; malformed regex can panic. | **OPEN**; evaluate #2082/#2241. |
| F05 | P1 | Unsafe `.gitleaksignore` file types, blocking FIFO, scan errors, lack of limits. | **OPEN**; issue #1572. |
| F06 | P0 | Global AND allowlist path short-circuit suppresses findings. | **V1 implemented;** release regression proof still required. |
| F07 | P0 | Staged text classified as binary may be skipped. | **OPEN**; candidate #2338. |
| F08 | P1 | Forced Git color hides findings. | **V1 implemented;** release regression proof still required. |
| F09 | P1 | Archive extension false positives can skip ordinary text. | **OPEN**; candidate #2317. |
| F10 | P1 | MIME reclassification on later chunks misses secrets. | **V1 implemented;** release regression proof still required. |
| F11 | P1 | Merge-only Git changes can be missed in default history mode. | **OPEN / compatibility decision required**; #2301 changes behavior. |
| F12 | P1 | `dir` silently scans the wrong directory on invalid positional args. | **V1 implemented;** release regression proof still required. |
| F13 | P1 | Pre-commit helper lets nonzero scanner statuses pass. | **V1 implemented;** release regression proof still required. |
| F14–F25 | P2/P3 | Archive identification, size limits, path identity, Unicode redaction, stdout contracts, cancellation, build/CI, and other lower-priority findings. | **Deferred / individually unverified**; retain CP-V0 details rather than claiming closure. |

At CP-V1 completion, `go build`, normal tests, generation, targeted Linux race suites, RE2 build, and Docker smoke tests were reported passing; `go vet` had pre-existing test-only `lostcancel` warnings and some comprehensive tagged/race lanes had known limitations. Subsequent V1G branch CI had Windows/macOS report-path test failures. **Re-run and link current CI before upgrading any item to release-validated.**

## 4. External contributor recognition and code provenance

### Attribution states (must be used accurately)

- **Reviewed:** Studied an issue, PR, test, or discussion; no claim that code was copied.
- **Informed / independently implemented:** Contributor exposed a failure mode or approach, but RNC wrote its own implementation. Credit the source in documentation and the final upstream PR.
- **Adapted implementation:** RNC incorporated identifiable logic or code from another person's work and modified it. Record exactly what was adapted, the original PR/commit, and appropriate author credit.
- **Adopted commit:** A contributor's commit is deliberately incorporated as a **separate authored commit** (prefer `git cherry-pick -x`), not hidden inside the RNC squash. Preserve the original author's metadata and record the resulting fork SHA.
- **Not selected:** Reviewed but deliberately not integrated. Still credit relevant research when it shaped a decision; do not add a false co-author trailer.

**Current adoption ledger: No external PR is yet recorded here as a verified separately adopted commit.** Similarity to an upstream PR does not establish copying or permission to claim its code has been integrated. Determine provenance from actual diffs and development records before changing a status.

### 4A. Upstream contributions directly relevant to implemented V1 work

All entries below are **reviewed/informed**, not confirmed cherry-picks. Head SHAs were rechecked on 2026-11-11; recheck before actual integration.

| RNC area | Upstream work | Contributor | PR head SHA | What it contributed to our assessment | Code-adoption status |
| --- | --- | --- | --- | --- | --- |
| F02 | [#2307](https://github.com/gitleaks/gitleaks/pull/2307), issue [#2304](https://github.com/gitleaks/gitleaks/issues/2304) | [@jdotc1](https://github.com/jdotc1) | `0d1ae3ea1087bebcf5ce45982a770d0f4925af6c` | Proved preflight deletion; proposed non-destructive validation and regression. Atomic final publication is additional RNC scope. | **Reviewed; no adopted commit verified.** |
| F02 / permissions | [#2300](https://github.com/gitleaks/gitleaks/pull/2300) | [@ToniAdreal](https://github.com/ToniAdreal) | `2b08eb90383f6f221496d984bc3d9a575d002f8b` | Restrictive `0600` report-permission proposal and discussion of preflight deletion. | **Reviewed; no adopted commit verified.** |
| F02 / permissions | [#2312](https://github.com/gitleaks/gitleaks/pull/2312) | [@Varshith-Kali](https://github.com/Varshith-Kali) | `2009273e0029f91aeda4732552ae5489634d0749` | Alternative `0600` permissions proposal; still truncates existing destinations. | **Reviewed; no adopted commit verified.** |
| F06 | [#2227](https://github.com/gitleaks/gitleaks/pull/2227) | [@BasitS-hash](https://github.com/BasitS-hash) | `49b26eece1ea83ea58da0c6e0bf92326fb82a25c` | AND-allowlist path sufficiency and regression cases. | **Reviewed; inspect for code overlap before release.** |
| F08 | [#2284](https://github.com/gitleaks/gitleaks/pull/2284) | [@lekanjava](https://github.com/lekanjava) | `1894e6f2cdaef94fe66388e577e0f59f14bfcff2` | Explicit `--no-color` in Git acquisition to protect parsers. | **Reviewed; inspect for code overlap before release.** |
| F10 | [#2329](https://github.com/gitleaks/gitleaks/pull/2329) | [@cardene777](https://github.com/cardene777) | `21028d8858e31fd088029f2fdb6634c99a1d3563` | First-chunk-only MIME classification. | **Reviewed; inspect for code overlap before release.** |
| F12 | [#2334](https://github.com/gitleaks/gitleaks/pull/2334), issue [#1979](https://github.com/gitleaks/gitleaks/issues/1979) | [@Shanmukh-l33t](https://github.com/Shanmukh-l33t) | `ad321b7d7aa9b1a1642bb56813fa7443d5412971` | `dir` positional argument validation and wrong-directory regression. | **Reviewed; inspect for code overlap before release.** |

### 4B. Other inspected open PRs and future integration candidates

These are **research/deferred**, not delivered by CP-V1. SHAs marked here were recorded in CP-V0 (2026-11-10), except those explicitly rechecked above or below. Revalidate PR state, diff, author metadata, and final tip before integration.

| Scope | PR | Contributor | Historical PR head / note | Assessment / integration disposition |
| --- | --- | --- | --- | --- |
| F01 Git exit | [#2339](https://github.com/gitleaks/gitleaks/pull/2339) | [@delaysystem](https://github.com/delaysystem) | `7e5132faa1ccea4bfd0b4b006ca896b5dd03df3f` | Partial exit-status fix; add stderr/parse/worker handling. **Rechecked 11 Nov.** |
| F01 diff parsing | [#2246](https://github.com/gitleaks/gitleaks/pull/2246) | [@aron-intframe](https://github.com/aron-intframe) | `9b7b2243fcdeb4279cccd9df9348f9e537d4fb95` | Partial parser-drain approach; source callback still swallows error. **Rechecked 11 Nov.** |
| F01 Git stderr | [#2080](https://github.com/gitleaks/gitleaks/pull/2080) | [@Smeet23](https://github.com/Smeet23) | `ff77f9f58b1ec5fbbece852e0bfe9c1f050f24e5` | Partial stderr propagation. |
| F01 Git stderr | [#2292](https://github.com/gitleaks/gitleaks/pull/2292) | [@Saadanjum0](https://github.com/Saadanjum0) | **SHA to recheck** | Similar partial stderr propagation. |
| F01 error propagation | [#2252](https://github.com/gitleaks/gitleaks/pull/2252) | [@CAOShurong](https://github.com/CAOShurong) | `424b4a92ffbcf4aa4b16ef98b53ecce45e1e3e33` | Fragment error flow; not complete source coverage. **Rechecked 11 Nov.** |
| F01 fragment errors | [#2253](https://github.com/gitleaks/gitleaks/pull/2253) | [@arpitjain099](https://github.com/arpitjain099) | `4f58f98a1424203af39cc00a26f9672e9b45b794` | Bounded/aggregate error ideas; verify actual resource limits. |
| F01 file open errors | [#2249](https://github.com/gitleaks/gitleaks/pull/2249) | [@vaibhav8a](https://github.com/vaibhav8a) | `509fce15ad2eb2536bc8d840a3ed3e374d73bb2d` | Path-aware open errors, incomplete walk/read handling. |
| F01 file open errors | [#2313](https://github.com/gitleaks/gitleaks/pull/2313) | [@Varshith-Kali](https://github.com/Varshith-Kali) | `acfeacda3adc17b8659d46b82180edd9e5adc7c7` | Alternative open-error reporting. |
| F01 file open errors | [#2235](https://github.com/gitleaks/gitleaks/pull/2235) | [@hamodywe](https://github.com/hamodywe) | `d86cc860f3c5f7085d8d68c24bc09424a73aabf7` | Count-only alternative, less useful than path-aware reports. |
| F07 staged binary | [#2338](https://github.com/gitleaks/gitleaks/pull/2338) | [@VISHNU0906](https://github.com/VISHNU0906) | `7fed4140d87f62090c87b9f09c5c9b6073a84a27` | Strong staged blob-coverage candidate; check error propagation. **Rechecked 11 Nov.** |
| F09 archive identity | [#2317](https://github.com/gitleaks/gitleaks/pull/2317) | [@Varshith-Kali](https://github.com/Varshith-Kali) | `71e91200afb2415b283bb1721709b19c2fddd8aa` | Content-first sniffing candidate; test opposite extension case. **Rechecked 11 Nov.** |
| F11 merge-only secrets | [#2301](https://github.com/gitleaks/gitleaks/pull/2301) | [@bhaskargurram-ai](https://github.com/bhaskargurram-ai) | `6bcf2f91c1b8765bbf5387507b4f0e1258b01dcc` | Behavior-changing; requires operator decision and compatibility tests. |
| F03 empty secret redaction | [#2350](https://github.com/gitleaks/gitleaks/pull/2350) | [@kevin9327](https://github.com/kevin9327) | `f19096fae4ba7a3016476c7b9d5528a6f31bbd9c` | Narrow fix; full log boundary remains. |
| F03/F17 UTF-8 masking | [#2261](https://github.com/gitleaks/gitleaks/pull/2261) | [@arpitjain099](https://github.com/arpitjain099) | `b1ce07f27d8251ef22e5b47c80387e9b1590fc2f` | Rune-safe masking candidate. |
| F04 invalid regex | [#2082](https://github.com/gitleaks/gitleaks/pull/2082) | [@Smeet23](https://github.com/Smeet23) | `df8223e0008b6d5bf82d0433592b8f1557c97311` | Controlled regex error candidate. |
| F04 invalid regex | [#2241](https://github.com/gitleaks/gitleaks/pull/2241) | [@harriiinnii](https://github.com/harriiinnii) | **SHA to recheck** | Broader alternative with unrelated expected-regex modifications. |
| F15 size limits | [#2099](https://github.com/gitleaks/gitleaks/pull/2099) | [@alephnull-sh](https://github.com/alephnull-sh) | `1b3896fd7df9a667a09568b751dce6143613319e` | Partial size budget; streaming/memory issues remain. |
| F16 path identity | [#2248](https://github.com/gitleaks/gitleaks/pull/2248) | [@hamodywe](https://github.com/hamodywe) | `5c9ca86f8f4d1db1e17004b6a1b598137203face` | Broad root-relative path proposal; sequence after AND-allowlist safety. |

### 4C. Issue reporters and research contributors

Issue authors deserve recognition even when their proposed implementation is not chosen. This is an initial verified/CP-V0-informed list; extend it as more research is performed:

| Issue | Reporter | Contribution |
| --- | --- | --- |
| [#2299](https://github.com/gitleaks/gitleaks/issues/2299) | [@carfeii](https://github.com/carfeii) | Report files could expose discovered plaintext secrets through ambient umask. |
| [#2304](https://github.com/gitleaks/gitleaks/issues/2304) | [@jdotc1](https://github.com/jdotc1) | Destructive report-validation preflight and reproduction. |
| [#338](https://github.com/gitleaks/gitleaks/issues/338) | [@lochlainnw](https://github.com/lochlainnw) | Overwrites when scanning many repositories; prompted investigation of aggregation; upstream maintainer preferred external scripts over merging reports. |
| [#1065](https://github.com/gitleaks/gitleaks/issues/1065) | [@mesaglio](https://github.com/mesaglio) | Fail-fast validation of invalid report destinations before expensive scans. |
| [#1572](https://github.com/gitleaks/gitleaks/issues/1572) | [@joerg1985](https://github.com/joerg1985) | Ignore-file safety and blocking concerns. |
| [#2232](https://github.com/gitleaks/gitleaks/issues/2232) | [@siliroid](https://github.com/siliroid) | Source-file error visibility and false-clean behavior. |
| [#2337](https://github.com/gitleaks/gitleaks/issues/2337) | [@farrrr](https://github.com/farrrr) | Staged binary-classified text can be missed. |
| [#2316](https://github.com/gitleaks/gitleaks/issues/2316) | [@bhandrigan](https://github.com/bhandrigan) | Archive filename classification false positives. |
| [#1979](https://github.com/gitleaks/gitleaks/issues/1979) | [@caiodelgadonew](https://github.com/caiodelgadonew) | Wrong-directory `dir` behavior. |
| [#1651](https://github.com/gitleaks/gitleaks/issues/1651), [#1757](https://github.com/gitleaks/gitleaks/issues/1757) | [@rgmz](https://github.com/rgmz) | Decode/chunk boundary and input-coverage concerns. |
| [#1243](https://github.com/gitleaks/gitleaks/issues/1243) | Reporter to verify | Secret exposure in verbose redaction. |

**License:** Preserve upstream Gitleaks MIT `LICENSE` and notices. License permission to reuse code is not a substitute for accurate author credit. Do not invent email addresses or apply `Co-authored-by` merely because someone reported a problem.

## 5. Adopt-vs-adapt decision log (required before clean release)

For every outside PR potentially overlapping RNC work, add a decision row **after examining the exact patch and our corresponding code**:

| PR / commit | RNC files or behavior overlapped | Decision: reviewed / independent / adapted / adopt | Original author and email from Git object | Original SHA -> resulting fork SHA | Tests / compatibility evidence | Operator approval |
| --- | --- | --- | --- | --- | --- | --- |
| #2227 / `49b26ee` | F06 AND allowlist | **Undecided; reviewed, overlapping design** | Resolve before integration | — | V1 tests; revalidate | Pending |
| #2284 / `1894e6f` | F08 Git color | **Undecided; reviewed, overlapping design** | Resolve before integration | — | V1 tests; revalidate | Pending |
| #2329 / `21028d8` | F10 first-chunk MIME | **Undecided; reviewed, overlapping design** | Resolve before integration | — | V1 tests; revalidate | Pending |
| #2334 / `ad321b7` | F12 CLI arguments | **Undecided; reviewed, overlapping design** | Resolve before integration | — | V1 tests; revalidate | Pending |
| #2307 / `0d1ae3e` | F02 report validation | **Undecided; reviewed, partial design overlap** | Resolve before integration | — | V1 tests; V2 pending | Pending |
| #2300 / `2b08eb9`, #2312 / `2009273` | F02 permissions | **Alternative proposals; no adopted commit** | Resolve if borrowing specific patch | — | Permission tests; revalidate | Pending |

Do **not** call a PR "adopted" simply because a similar fix was independently implemented. If code was actually carried over, record the evidence and correct attribution before final integration.

## 6. Clean-history and upstream-delivery strategy

**Operator intent:** Iterate freely on feature/development branches, then deliver **one clean RNC-owned squash commit** for original work while preserving any outside contributors' code as separate authored commits. The final public diff and attribution should be intelligible to upstream reviewers.

### Release reconciliation procedure (operator-controlled; not authorized now)

1. **Freeze and audit:** Confirm final upstream/main base, selected fixes, feature flags, test evidence, licenses, and every overlapping community PR. Compare the fork against the actual chosen upstream SHA, not just a stale historical pin.
2. **Classify provenance:** For each reviewed PR, choose *independently implemented*, *adapted*, *adopted*, or *not selected*. Record what changed and who authored the original idea and code.
3. **Construct a clean staging branch** (working name `dev/clean-up`) from the approved release base. This should be performed only on explicit operator authorization; do not rewrite or disturb `fix/validation`.
4. **For adopted commits:** Fetch the contributor PR/commit, verify the actual author and patch, and cherry-pick with `-x` into the clean staging sequence. Preserve the contributor's author metadata; document original and new fork SHAs. If cherry-picking requires adjustments, use a separate RNC follow-up commit or an expressly documented adaptation, preserving credit.
5. **For RNC-owned work:** Rebuild one clean cumulative delta **excluding code already contained in adopted contributor commits**, and record it as one RNC-authored squash commit. It may be ordered before or after the adopted commits based on dependencies. Never duplicate somebody else's patch inside the squash and then attempt to cherry-pick it again.
6. **Validate the reconstructed tree:** Compare final code/output against the fully accepted working branch; run all test/CI/benchmark compatibility checks; inspect the complete diff and Git authorship graph.
7. **Preserve the final authorship graph:** Integrate the clean staging history into RNC `master` without squashing the separate contributor commits again. **No merge, rebase, branch deletion, or upstream PR in the current task.**

### Honest commit-count expectation

- If **no external code** is adopted, the fork can finish **one RNC implementation commit ahead** of the chosen upstream baseline (subject to the actual baseline state).
- If **N independently authored external commits** are deliberately adopted, the final history necessarily has **one RNC implementation commit plus N contributor commits** (and any justified integration adjustments). You cannot simultaneously preserve their separate authored commits and have exactly one total commit ahead.
- Merely mentioning a contributor in a squash commit, adding a `Co-authored-by` trailer, or linking an issue **does not** preserve their original commit as a separate history entry. Use each mechanism only for its accurate purpose.

### Final upstream PR attribution requirements

The future upstream PR description should contain:

1. Clear, bounded problems solved and reproducible before/after evidence.
2. Backward-compatibility statement and all intentional behavior changes/defaults.
3. Test matrix with Windows/Linux/macOS, standard/RE2, Docker, and relevant failure paths.
4. **Code incorporated:** each contributor, PR, original commit, resulting fork commit, and any adjustments.
5. **Reviewed/influential work not incorporated:** relevant authors and issues/PRs with a sentence explaining how they informed the solution.
6. Open limitations, independent follow-up issues, and a clean diff against upstream.

**Upstream caveat:** The checked `README.md` currently states Gitleaks is feature-complete and upstream accepts security patches only, with active product focus moving to Betterleaks. A future upstream PR is an operator decision, not an assumed acceptance path; keep the patch scoped and distinguish safety fixes from new optional QoL features.

## 7. Updating this file

Whenever a fix or QoL feature changes:

1. Update its row and status (**Planned**, **Implemented**, **Validated**, **Open**, **Deferred**).
2. Add the applicable CP finding, direct issue/PR links, and evidence/test links.
3. Record default/CLI/report/exit-code changes explicitly under the compatibility contract.
4. Review the upstream attribution table for missing contributors, conflicting proposals, and code overlap.
5. If outside code is incorporated, update the adopted-commit ledger immediately; preserve actual author metadata and SHA crosswalk.
6. Reconcile this ledger with release notes and the eventual upstream PR. Never silently remove someone from the acknowledgments solely because their exact patch was not selected.

This file tracks engineering evidence and attribution; it is **not** a request to adopt upstream work or a license to merge anything automatically.

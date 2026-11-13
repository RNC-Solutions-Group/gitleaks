# rnc-mod Integration Plan

**Status:** proposed; not executed  
**Base for integration:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`  
**Experimental branch:** `audit/rnc-mod-upstream-integration`  
**Final target branch:** `rnc-mod` (not created or modified by this audit)  
**Source to reconcile:** CP-V10 `bd42a958f7476e9a6d8e17bd2dbb2abe3ff060ea`

## 1. Final eligible contributor commits

Preserve these 18 upstream commits, in this order, from the experimental branch:

1. `3004bda` #2227 — BasitS-hash
2. `c8d7070` #2227 — BasitS-hash
3. `bdc880a` #2334 — Shanmukh-l33t
4. `1f3533b` #2261 — Arpit Jain
5. `7a4c18f` #2082 — Smeet Agrawal
6. `d68ecb2` #2307 — Jeff / jdotc1
7. `2cf0ade` #2329 — cardene777
8. `f146e23` #2284 — Sakariyah Abdulhazeem
9. `5bed958` #2246 — INTFRAME
10. `16ba49e` #2338 — KVR Vishnu
11. `b603fa0` #2339 — DelaySystem
12. `2bd5ad3` #2080 — Smeet Agrawal
13. `62a89d4` #2099 — Nikhil
14. `7ce04bb` #2317 — Varshith Puli
15. `8ee08ef` #2248 — hamodywe
16. `5e35358` #2249 — Vaibhav Srivastava
17. `260d7e1` #2253 — Arpit Jain
18. `6445d1e` #2312 — Varshith Puli

The two integration-adjustment commits (`22f63b4`, `100ae90`) are RNC-authored
and should remain only if the final reconciliation still needs them. They are
not contributor commits.

Do not squash any of the 18 together. Do not rewrite their author fields.

## 2. Integration order and dependencies

The experimental branch order is the recommended order because it minimizes
conflicts:

1. CLI/allowlist/redaction/config/report preflight/first-chunk:
   #2227, #2334, #2261, #2082, #2307, #2329.
2. Git command and parser:
   #2284, #2246, #2338, #2339, #2080.
3. Size/archive:
   #2099, #2317.
4. Root-relative matching:
   #2248.
5. File/fragment errors:
   #2249, #2253.
6. Report permissions:
   #2312.

Dependencies to watch:

- #2246 precedes #2338 because it moves `GitCmd` construction into `newGitCmd`;
  #2338 sets `staged` in the constructor path.
- #2246/#2339/#2080 modify the same parser error paths.
- #2099/#2248/#2317/#2329 modify `sources/file.go`.
- #2099/#2248/#2249 modify `sources/files.go`.
- #2248/#2253 modify `detect/detect.go`.
- #2307/#2312 modify `cmd/root.go`.

## 3. Preparing the final branch

Proposed operator commands, after approval:

```bash
git fetch origin
git switch --create rnc-mod 5f6ec56
git merge --no-ff audit/rnc-mod-upstream-integration   # or cherry-pick the
                                                       # documented sequence
```

If `rnc-mod` already exists by the time of execution, stop and reconcile its
actual state rather than force-moving it. The audit branch and `rnc-mod` must
not be created by rewriting protected history.

## 4. V10 reconciliation strategy

Do **not** cherry-pick the individual CP-V1–V10 commits into the contributor
history. The approved release unit is the validated CP-V10 tree. Apply it as a
single RNC-authored reconciliation commit (or an operator-approved squash of
the RNC work) on top of the contributor commits.

The reconciliation commit must:

- resolve conflicts in favor of validated V10 behavior, except where an upstream
  contribution is demonstrably better and the operator approves the deviation;
- preserve contributor tests where they still express valid behavior;
- remove duplicated upstream code that V10 replaced, rather than retaining dead
  alternates;
- retain the original upstream commit ancestry;
- not claim that V10 is byte-identical to any upstream PR.

Suggested commit message:

```text
feat(rnc): reconcile validated V10 scanner implementation

Apply the accepted CP-V10 tree on top of the preserved upstream contributor
commits. Conflicts resolved in favor of validated V10 behavior; contributor
regression coverage retained where applicable.
```

## 5. Anticipated conflict map

| Area | Contributor commits | V10 side | Resolution |
| --- | --- | --- | --- |
| Finding redaction | #2261 | `report/finding.go` expanded redaction | Keep V10 helpers; merge rune mask and upstream UTF-8 tests |
| Config regex | #2082 | V10 compile/error expansion | Keep V10 errors; retain wrapper API |
| Archive verification | #2317 | V10 error-returning `contentMatchesArchive` and `restorePrefix` | Keep V10; retain upstream tests |
| Size limits | #2099 | V10 spooling and `Size` | Keep V10 spooling; preserve field propagation and tests |
| Staged binary | #2338 | V10 `:0:` index path | Keep V10 index lookup; preserve staged fixture |
| Git parser | #2246/#2339/#2080/#2284 | V10 tracker, stderr classifier, scanerr, forceNoColor | Keep V10 lifecycle; preserve upstream parser tests |
| Root-relative | #2248 | V10 `fragmentPaths`, absolute baseline normalization, F16 fixes | Keep V10; preserve upstream root tests |
| Directory args | #2334 | V10 empty-arg rejection | Keep V10 behavior |
| AND allowlist | #2227 | V10 root-aware `shouldSkipPathAtRoot` | Keep V10 and contributor helper |
| File/fragment errors | #2249/#2253 | V10 `scanerr` | Keep V10 collector; preserve upstream continuation tests |
| Report validation | #2307 | V10 `reportOutput` temp publication | Keep V10; retain no-preflight principle |
| Report permissions | #2312 | V10 temp `0600` + publication | Keep V10 temp model; the upstream `openReportFile` may be dropped if no longer used |
| Build/RE2/Docker/CI | none | V10-specific | Apply directly; no contributor conflict |

The exact function-level analysis is in `source-diff-analysis.md`.

## 6. Validation commands

Run from the reconciled tree:

```bash
gofmt -l .
go build ./...
go test -count=1 ./...
go test -count=1 -tags gore2regex ./...
go vet ./...
go mod verify
```

Targeted regressions:

```bash
go test ./report -run 'Mask|Redact' -count=1
go test ./config -run 'Translate|Regex|Allowlist' -count=1
go test ./sources -run 'Archive|File|Git|Root|Size|Unreadable' -count=1
go test ./detect -run 'DirRoot|Git|Fragment|Allowlist|Baseline' -count=1
go test ./cmd -run 'Directory|Report|PreCommit' -count=1
```

If the local Windows host still lacks a C compiler, run race tests in Linux CI
or a Linux container; do not report unavailable local race tests as passing.
The final acceptance must also compare the reconciled tree and representative
scanner outputs against `bd42a95` using the CP-V10 comparison fixtures.

## 7. Source-equivalence acceptance criteria

The reconciled tree is acceptable only when:

1. `go test -count=1 ./...` passes.
2. Full tagged RE2 tests pass (the CP-V10 generator fix is present).
3. `go vet ./...` passes with no new failures.
4. `git diff bd42a95..rnc-mod` contains only the intended contributor ancestry,
   the V10 reconciliation commit, audit documentation, and explicitly approved
   differences.
5. Findings, paths, fingerprints, exit codes and report structure match
   `bd42a95` on the CP-V10 fixture set.
6. Git history shows the original contributor names/emails for all 18 commits.
7. No contributor commit is squashed into the RNC commit.
8. No protected branch was rewritten.

## 8. Rollback strategy

- Keep `audit/rnc-mod-upstream-integration` untouched as the provenance record.
- Create `rnc-mod` only from `5f6ec56`; if reconciliation fails, delete the
  local `rnc-mod` candidate and restart from the audit branch (branch deletion
  requires operator authorization).
- Never force-push a shared `rnc-mod`; publish only after validation.
- Preserve the integration manifest and SHA crosswalk so any failed experiment
  can be reconstructed.

## 9. Required decisions / approvals

1. Approve preservation of the 18 upstream commits as separate ancestry.
2. Approve the two RNC integration-adjustment commits or require them to be
   folded into the reconciliation commit.
3. Approve V10 as the behavioral source of truth for every conflict.
4. Decide whether #2300/#2292/#2252/#2313/#2203 remain documentary only (this
   audit recommends yes, because each is an alternative to a selected commit).
5. Confirm that the final `rnc-mod` branch name and merge mechanism are
   authorized; this audit does not create or push `rnc-mod`.
6. Confirm MIT notice preservation and whether an additional `NOTICE`/
   acknowledgment section is desired; the license already permits the use.

## 10. Explicit answers to the required recommendation questions

1. **Which upstream commits should be preserved?** The 18 commits from 17 PRs
   listed in §1.
2. **Which PRs fail the substantive threshold?** None of the 11 primary PRs.
   The secondary alternatives #2300, #2292, #2252, #2313 and #2203 are not
   separate implementations once the selected commits are applied; they remain
   documentary. #2312 passes through its 0600 implementation and is integrated.
3. **Which contributions were underestimated?** #2246, #2248, #2334, #2227,
   #2317, #2338, #2099, #2261 and #2082 have direct/adapted V10 overlap. The
   existing ledger already flagged several; this audit raises #2317, #2338 and
   #2099 from “review” to explicit adapted implementation.
4. **Which upstream commits conflict with one another?** #2246/#2338/#2339/
   #2080 in Git; #2099/#2317/#2329/#2248 in `sources/file.go`; #2249/#2313 in
   files; #2307/#2312 in report; #2252/#2253 in detection errors.
5. **Which conflict with V10?** All areas in §5, especially Git lifecycle,
   report publication, redaction breadth, RE2/build, and F16.
6. **Can conflicts be resolved preserving V10 behavior?** Yes; the selected
   upstream commits represent narrow precursors and V10 contains corrections.
   The experimental branch compiles and tests cleanly before reconciliation.
7. **Are any upstream changes worth retaining instead of V10?** The contributor
   tests and the specific `RootRelativePath`, `pathIsSufficientToSkip`,
   `contentMatchesArchive`, parser-drain and staged-field concepts should be
   retained where V10 already preserves them. No wholesale replacement of V10
   is recommended.
8. **Licensing concerns?** MIT permits reuse; keep the upstream copyright and
   permission notice. Preserve original Git authorship; do not invent notices.
9. **Exact merge/squash strategy?** Build `rnc-mod` from the audit branch and
   apply one RNC V10 reconciliation commit; do not squash contributor commits;
   do not squash V10 into them. The final branch is one RNC commit on top of
   preserved upstream ancestry, not one total commit.
10. **Unverified/approval-required?** The final `rnc-mod` branch creation,
    final reconciliation, race tests on the target environment, and any
    decision to retain an upstream behavior over V10 require operator approval.

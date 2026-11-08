# Performance Stability — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**3 active upstream issues · 2 upstream PRs · 4 work packages.** **Focused performance scrub completed 2026-11-08:** all four work packages remain active for independent review; **no new obsolete, superseded, or exclusion-by-exception dispositions** were supported. Historic v8.2.1 crash #752 is already [obsolete](exclusions.md). The README-only #2212 patch remains in [documentation](documentation.md); issue #2019 stays here for its memory/OOM behavior. Profiling-script PR #1900 belongs here because it changes executed commands. Upstream state and scanner code are unchanged.

Memory, hangs, race-prone parsers, decoding cost, and crashes. First establish a measurable reproduction and a bounded test workload.

**Evidence rules:** “SOURCE-CHECKED” means a targeted patch/current-source comparison, **not** that a fix was built or tested. Most records are title/crosswalk-only and must stay pending until reproduced, diffed, and assessed on the fork. “Obsolete,” “superseded,” “adopted,” or “corrected” require separate evidence; no PR is adopted by routing it here.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| [#1338](https://github.com/gitleaks/gitleaks/issues/1338) | [#2246](https://github.com/gitleaks/gitleaks/pull/2246) (@aron-intframe) | A malformed Git diff can make the streaming parser stop early, silently report an incomplete clean scan, or deadlock on a filled pipe. | **SOURCE-CHECKED — RETAIN; TEST PENDING** | Master `sources/git.go` still calls `gitdiff.Parse(stdout)` without detecting unread output after parser termination; its `blobReader` drain is a *different* path. #2246 adds source changes and regression fixtures. The PR author's synthetic Git shim shows both failure modes; that is not an independently reproduced native-Git emitter. Must verify incomplete-scan signaling, hangs, and normal Git throughput before deciding. | — |
| [#2019](https://github.com/gitleaks/gitleaks/issues/2019) | — | Large GitLab repository scan allegedly reaches ~24 GB/OOM with recursive decoding enabled. | **ISSUE — RETAIN; REPRO PENDING** | `cmd/root.go` still defaults `--max-decode-depth` to **5**, intentionally according to maintainer comments; setting depth 0 is a workaround with different detection coverage, not a correction. The supplied pprof excerpt does not reproduce the reported 24 GB peak, and no public fixture exists. README-only #2212 is separate. Profile bounded input, decode depth, peak RSS and exit status before assessment. | — |
| — | [#1900](https://github.com/gitleaks/gitleaks/pull/1900) (@vignesh1507) | Change `scripts/profile.sh` benchmarking/profiling output, shell quoting, `find`, optional `yq`, and diagnostics invocation. | **SOURCE-CHECKED — RETAIN; SCRIPT TEST PENDING** | The original script remains on master. PR is not a cosmetic-only diff: its diagnostics scan no longer inherits baseline `--exit-code=0 --max-decode-depth 8` options used by its benchmark; command equivalence and comparability need tests. Existing master diagnostic flags still exist; no evidence yet to call the PR obsolete or adopted. | — |
| [#1437](https://github.com/gitleaks/gitleaks/issues/1437) | — | Kubernetes/Linux hang with a futex wait observed in Gitleaks v8.18.4. | **ISSUE — RETAIN; EVIDENCE INSUFFICIENT** | Original report supplies a `strace` futex wait and SIGTERM, but no stack dump, public repository, or reproducer. Age alone does not prove the underlying hang is fixed; no responsible superseding PR was identified in this focused pass. Require current-master reproduction with bounded timeout and goroutine stack capture. | — |

## Focused scrub evidence and limits

- **Baseline:** [RNC `master` `5f6ec56c`](https://github.com/RNC-Solutions-Group/gitleaks/tree/5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c), upstream issue bodies/comments #1338, #2019, #1437, and complete upstream patches #2246 and #1900.
- **Elimination result:** 0 issues and 0 PRs excluded. No deduplication proved: #1338 is a diff-parser failure with deterministic *synthetic* reproduction; #1437 is an unidentified runtime futex hang; #2019 is memory pressure associated with decode depth; #1900 is profiling tooling. Superficially similar symptoms do not establish one root cause.
- **Safety note:** #1338 also has scan-integrity consequences (incomplete results reported as clean). Keep its single primary owner here during this one-category pass; don't silently transfer it or treat the source-checked proposal as a tested fix.
- **Runtime evidence:** no builds, benchmarks, new runtime tests, external deployments, upstream closures, or PR adoptions performed in this triage branch.

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

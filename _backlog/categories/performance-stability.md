# Performance Stability — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**4 issues · 4 PRs · 6 work packages.**

Memory pressure, hangs, parsing stability, decoding cost, runtime efficiency, and profiling. Establish bounded reproducible workloads.

**Review states:** `SOURCE-CHECKED` includes a targeted fork-source comparison; `PATCH-REVIEWED` means an upstream diff was inspected, **not** that the patch was built or tested against the fork. `PENDING` records unresolved issues or competing proposals. No PR is adopted or upstream item closed by categorization. See [RNC exclusions](exclusions.md) for retired review slots.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| [#2121](https://github.com/gitleaks/gitleaks/issues/2121) | [#2271](https://github.com/gitleaks/gitleaks/pull/2271) (@dhanushs11) | Detection rule regex runs twice on each fragment. | **SOURCE-CHECKED — TEST PENDING** | PR #2271 reuses existing match indices and adds a benchmark. The included benchmark measures regex work, not scanner correctness; add a regression before adoption. | — |
| [#1338](https://github.com/gitleaks/gitleaks/issues/1338) | [#2246](https://github.com/gitleaks/gitleaks/pull/2246) (@aron-intframe) | A malformed Git diff can make the streaming parser stop early, silently report an incomplete clean scan, or deadlock on a filled pipe. | **SOURCE-CHECKED — RETAIN; TEST PENDING** | Master `sources/git.go` still calls `gitdiff.Parse(stdout)` without detecting unread output after parser termination; its `blobReader` drain is a *different* path. #2246 adds source changes and regression fixtures. The PR author's synthetic Git shim shows both failure modes; that is not an independently reproduced native-Git emitter. Must verify incomplete-scan signaling, hangs, and normal Git throughput before deciding. | — |
| [#2019](https://github.com/gitleaks/gitleaks/issues/2019) | — | Large GitLab repository scan allegedly reaches ~24 GB/OOM with recursive decoding enabled. | **ISSUE — RETAIN; REPRO PENDING** | `cmd/root.go` still defaults `--max-decode-depth` to **5**, intentionally according to maintainer comments; setting depth 0 is a workaround with different detection coverage, not a correction. The supplied pprof excerpt does not reproduce the reported 24 GB peak, and no public fixture exists. README-only #2212 is separate. Profile bounded input, decode depth, peak RSS and exit status before assessment. | — |
| — | [#1900](https://github.com/gitleaks/gitleaks/pull/1900) (@vignesh1507) | Change `scripts/profile.sh` benchmarking/profiling output, shell quoting, `find`, optional `yq`, and diagnostics invocation. | **SOURCE-CHECKED — RETAIN; SCRIPT TEST PENDING** | The original script remains on master. PR is not a cosmetic-only diff: its diagnostics scan no longer inherits baseline `--exit-code=0 --max-decode-depth 8` options used by its benchmark; command equivalence and comparability need tests. Existing master diagnostic flags still exist; no evidence yet to call the PR obsolete or adopted. | — |
| [#1437](https://github.com/gitleaks/gitleaks/issues/1437) | — | Kubernetes/Linux hang with a futex wait observed in Gitleaks v8.18.4. | **ISSUE — RETAIN; EVIDENCE INSUFFICIENT** | Original report supplies a `strace` futex wait and SIGTERM, but no stack dump, public repository, or reproducer. Age alone does not prove the underlying hang is fixed; no responsible superseding PR was identified in this focused pass. Require current-master reproduction with bounded timeout and goroutine stack capture. | — |
| — | [#2023](https://github.com/gitleaks/gitleaks/pull/2023) (@sissl0) | Reduce repeated regular-expression work by caching newline offsets and reusing match indexes in `detect/detect.go`. | **PATCH-REVIEWED — PERF TEST PENDING** | Changed files are detector runtime and tests, not regex rule changes. Benchmark allocations/CPU and verify location/line correctness against current master; compare performance approach to #2271 without assuming equivalence. | — |

## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

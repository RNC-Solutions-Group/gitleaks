# CP-V17 C2 disposition matrix

The V15 C2 roster has twelve packages. Outcomes below are based on the accepted V16
tree plus the V17 candidate, not on stale issue labels.

| Package | Upstream references | V17 outcome | Evidence |
| --- | --- | --- | --- |
| `GIT-01` | #2003, #2128 | **Verified / regression retained** | The V13/V16 shell-style parser already turns quoted `'--no-notes'` into one argument. New test `TestQuotedNoNotesReachesGitCommand` proves the option reaches `git log`; the CSV claim was not reproducible on Git 2.45.2 because `--pretty=medium` does not include notes without an explicit notes option. No product behavior change. |
| `GIT-02` | #1028, #2301 | **Opt-in implemented; default unchanged; OPERATOR DECISION REQUIRED (OD-1)** | `--scan-merges` adds `--diff-merges=first-parent`. Tests cover ordinary merge, evil merge, conflict resolution, and octopus. Default scans still omit merge patches. |
| `GIT-03` | #1993 | **Not reproduced; residual** | Windows/Linux pipe tests and the existing drain-before-wait path complete. No macOS PTY host was available, so the reported 0%-CPU hang is not confirmed. No arbitrary process-group or timeout change was added. |
| `GIT-04` | #1647 | **Opt-in implemented** | `--git-context-lines N` includes bounded unchanged context for multi-line rules. Findings are credited only when their line span overlaps an added line. Kubernetes Secret fixture: 0 default findings, 2 with 10 context lines; unchanged-context control produces no second finding. |
| `GIT-05` | #1364 | **Not reproduced; environment/version-specific** | Current `git diff --staged` works with Git 2.45.2; the Docker build uses Git 2.54.1. The issue was filed against 8.18.2-era behavior. No code change. |
| `GIT-06` | #1310 | **Warning retained; reconstruction impossible** | Shallow state is still warned as incomplete. Missing pre-boundary objects cannot be reconstructed locally. No false-clean suppression or history fabrication. |
| `GIT-07` | #821 | **Deferred feature design** | Incoming pre-receive objects are not refs until the hook accepts them; scanning is a range/refs problem, not a golden-history guarantee. Document `--log-opts <old>..<new>`; no hidden ref mutation. |
| `GIT-08` | #1293 | **Deferred policy** | Bare repositories have no working tree for automatic local config/ignore discovery. Existing warning tells the operator to pass explicit paths; a bare-local trust policy needs an operator decision. |
| `GIT-09` | #1163 | **Deferred policy** | Git history is repository-wide; the existing warning tells the operator to use a pathspec. Silently narrowing a subdirectory argument would change scope and fingerprints. |
| `PERF-01` | #2019 | **Measured; residual** | Default decode depth 5 on 1.25 MB/14,641-finding corpus peaked ~136 MB; depth 0 ~41 MB. No 24 GB reproduction. Finding/segment retention is the dominant scale risk; no unsafe cap added. |
| `PERF-03` | #1437 | **Not reproduced; residual** | No futex stack or hang was obtained. Current Git and file sources drain both pipes and wait before returning. The report remains unconfirmed rather than being closed by assertion. |
| `PERF-04` | #2023 | **Partially adopted; wholesale patch rejected** | Adopted the safe newline precompute. Rejected the patch's broader rewrite and any size/newline semantic change. Added a conservative decoder prefilter and measured it; no duplicate regex search remains in the current path. |

## GIT-02 / OD-1 evidence

Controlled fixtures, `--exit-code 0`, JSON output:

| Scenario | Default | `--scan-merges` |
| --- | ---: | ---: |
| Evil merge only | 0 | 1 at merge SHA |
| Conflict-resolution merge only | 0 | 1 at merge SHA |
| Octopus merge only (4 parents) | 0 | 1 at merge SHA |
| Branch secret merged normally | 1 at branch SHA | 2: branch + merge SHA |
| Same branch secret with SHA-less `.gitleaksignore` | 0 | 0 |

The merge copy has a distinct fingerprint (`<mergeSHA>:<file>:<rule>:<line>`). A
SHA-less ignore entry suppresses both copies; a SHA-scoped ignore or an existing
baseline can surface the merge copy once.

### OD-1 alternatives

| Strategy | Coverage | Cost / compatibility | Recommendation |
| --- | --- | --- | --- |
| Keep default `git log -p` | Misses merge-only secrets | No new fingerprints | Current accepted default |
| `--diff-merges=first-parent` default | Covers evil/conflict/octopus merge-only content | Duplicates branch content at merge; new fingerprints; larger output | Not without operator approval and C5 migration |
| Combined `--cc` / remerge | Parser cannot consume `diff --cc`; remerge not universally available | Broken/incomplete scan risk | Reject |
| Opt-in `--scan-merges` | Same as first-parent, operator-selected | Explicit duplicate/fingerprint tradeoff | **Recommended V17 posture** |
| Cross-commit dedup | Avoids duplicates | Changes attribution, identity, and baseline behavior | Defer to C5 |

**Technical recommendation:** retain the default and ship the discoverable opt-in. The
security benefit is real, but making it default would silently add findings and change
fingerprints for every merge in every deployment before C5 has a migration contract.
RNC CI templates should adopt `--scan-merges` where complete merge coverage is required.

## GIT-04 evidence

`--git-context-lines` is a separate opt-in. It does not change merge behavior. The
implementation records added-line ranges and rejects findings wholly inside context, so
the pre-existing secret in unchanged context is not re-reported at a later commit.

## Validation boundary

All C2 tests run in-process or against temporary Git repositories. No independent
harness, frozen oracle, or benchmark manifest was modified. The macOS PTY case and a
true pre-receive environment remain unexecuted and are reported as residual, not passed.

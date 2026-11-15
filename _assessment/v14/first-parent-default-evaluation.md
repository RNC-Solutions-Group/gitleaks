# First-Parent Merge-History Evaluation (OP-03, revised V14A)

**Decision state:** EVALUATION ONLY. The current default remains unsupported merge-patch omission plus explicit `--log-opts=--diff-merges=first-parent` opt-in. No default or fingerprint change was made.

## Workloads and commands

Controlled synthetic repositories under a temporary root. The same V14A scanner/config was used for every mode:

```bash
gitleaks git <repo> --log-opts="--all" ...
gitleaks git <repo> --log-opts="--all --diff-merges=first-parent" ...
```

Git: 2.45.2.windows.1. A non-root Linux Docker race/tests run covered the scanner, not the Git-history experiment.

## Measured results

| Scenario | Default findings | First-parent findings | Unique locations (first) | Duplicates (first) | Merge-only recovered | Notes |
| --- | ---: | ---: | ---: | ---: | ---: | --- |
| Ordinary merge | 1 | 2 | 1 | 1 | 0 | feature commit + merge duplicate |
| Evil merge-only amendment | 1 | 3 | 2 | 1 | 1 | feature commit + merge duplicate + merge-only |
| Conflict-resolution secret | 0 | 1 | 1 | 0 | 1 | secret exists only in merge resolution |
| Octopus merge-only amendment | 0 | 1 | 1 | 0 | 1 | merge commit carries the secret |
| Repeated value on two branches | 2 | 3 | 2 | 1 | 0 | one merge duplicate of feature occurrence |
| Pathspec `-- feature.txt` | 1 | 2 | 1 | 1 | 0 | pathspec does not remove merge duplicate |
| Feature fingerprint ignored | 1 | 1 | 1 | 0 | 0 | ignore covers contributing commit; merge duplicate remains |
| 100-commit history + one merge | 1 | 2 | 1 | 1 | 0 | median 315.8 ms vs 306.6 ms (3 runs; startup-dominated) |

The extra first-parent count is not accuracy: in ordinary/repeated/pathspec/large histories it is a duplicate occurrence. Genuine new coverage appears only for merge-only/conflict/octopus changes.

## Fingerprint and baseline effects

- The same content/line in a merge commit and a contributing commit has different commit-prefixed fingerprints.
- A v1 `.gitleaksignore` entry for the contributing feature fingerprint suppressed that finding but did **not** suppress the merge-commit duplicate.
- A baseline generated under the default likewise cannot suppress merge-commit findings until regenerated or explicitly ignored.
- Commit attribution of the duplicate points at the merge commit; the original occurrence remains attributed to its branch commit.

## Git-version and operational compatibility

- `--diff-merges=first-parent` requires Git 2.31+; older Git fails operationally rather than silently clean.
- The existing parser accepts the option; no new shell execution.
- Octopus/conflict cases produced no operational errors in this test.
- No pathspec failure was observed; the pathspec contract was preserved.
- Runtime impact was not separable from process startup on this workload. Larger histories and additional merge patches need dedicated benchmarking before a default claim.

## Alternatives

1. Retain the current default and explicit opt-in (recommended).
2. Change the default to first-parent, accepting duplicate findings and new fingerprints.
3. Add a "comprehensive history" mode, which would require an explicit duplicate/identity policy and is not authorized.
4. Deduplicate findings automatically — rejected: it would change primary finding identity and fingerprint semantics.

## Recommendation

**Retain the current default for V14A and keep first-parent explicit opt-in.** The experiment recovers real merge-only/conflict/octopus secrets, but it also creates duplicate occurrences, new fingerprints, ignore/baseline migration work and a Git-version floor. There is insufficient evidence for a default change.

A future default change requires:
- an approved duplicate/fingerprint/baseline migration;
- large-repository timing and memory measurements;
- octopus/conflict/pathspec fixtures retained as regression tests;
- operator approval and a Git-version support decision.

## Residual limitations

- This is a synthetic, small-history experiment; it does not establish detection accuracy or production runtime behavior.
- Older Git versions and non-Windows Git behavior were not executed.
- No independent scanner-harness accuracy run was performed.

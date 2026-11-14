# First-Parent Merge-History Evaluation (OP-03)

**Decision state:** EVALUATION ONLY. The current default remains unsupported merge-patch omission plus explicit `--log-opts=--diff-merges=first-parent` opt-in.

## Workload

Synthetic repository with:

- base commit;
- feature branch adding `FEATURE_MERGE_SECRET_123456` in `feature.txt`;
- master commit adding `master.txt`;
- non-fast-forward merge commit;
- evil-merge amendment adding `MERGE_ONLY_SECRET_987654` to `base.txt`.

Commands (same binary/config):

```bash
gitleaks git <repo> --log-opts="--all" ...
gitleaks git <repo> --log-opts="--all --diff-merges=first-parent" ...
```

## Measured results

| Mode | Findings | Commits | Unique secret occurrences | Duplicate |
| --- | ---: | --- | ---: | ---: |
| Default | 1 | `c152670` feature | 1 | 0 |
| First-parent | 3 | `73c776f` merge (feature duplicate), `c152670` feature, `73c776f` merge-only | 2 | 1 |

First-parent output:

- `feature.txt` from the merge commit `73c776f` and from `c152670` — the same content/line under two commit identities;
- `base.txt` from `73c776f` — the merge-only evil-edit secret.

Default output:

- only `feature.txt` from `c152670`;
- merge-only secret absent;
- merge duplicate absent.

The extra finding count is **not** accuracy: one of the three first-parent findings duplicates the feature commit. The genuine coverage gain is the `base.txt` merge-only finding.

## Compatibility consequences of a default change

- New findings at merge commits and duplicate findings for content already attributable to contributing commits.
- New commit-prefixed fingerprints; existing `.gitleaksignore` entries for contributing commits do not suppress merge-commit duplicates.
- Baselines created under the default would not suppress merge findings until regenerated or explicitly ignored.
- Merge commit attribution may be less intuitive for code-review workflows.
- Git requirement rises to Git 2.31+ for the opt-in form; older Git fails operationally rather than silently clean.
- Runtime increases with additional merge patches; this small workload does not quantify large-repository cost.
- `--no-diff-merges` remains an override.

## Alternatives

- Keep default; document the opt-in (current).
- Default first-parent with duplicate findings; no safe automatic dedupe exists under the current report contract.
- Default first-parent plus source-level dedupe would change primary finding identity and is not authorized.
- Combined/remerge output is rejected; the parser cannot complete it.

## Recommendation

**Retain the current default for V14.** The evaluation confirms the known coverage defect, but changing the default would alter fingerprint/baseline/ignore behavior and increase findings without a complete duplicate-policy migration. Operators needing merge-only coverage should use the explicit first-parent option and regenerate/ignore findings deliberately.

## Reopen criteria

Operator approval of a default change with a migration plan covering duplicates, fingerprints, baselines, ignore files, Git-version floor, runtime and CI policy. Octopus/conflict-resolution/pathspec matrix and a large-history benchmark remain required before that decision.

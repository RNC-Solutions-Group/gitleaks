# Post-V14 backlog reconciliation — 2026-11-15

**Status:** APPLIED to the canonical active category and exclusion ledgers. This is a
documentation reconciliation, not a scanner change, release or upstream closure.
**Accepted product:** `RNC-Solutions-Group/gitleaks@rnc-mod` =
`7d74921841114be7677a423cf60620210e832bf7` (tree
`8eb327dd5860a4665daf38ce2f2fbc0aa5c07963`).
**Canonical backlog before this pass:** `feat/issues-and-pr@7f7e4eada3bdf3b594a3c0ec1fd101b89f2f08e2`.
**Post-merge CI:** run `43763741641`; 13 successful checks, 1 intentionally skipped
`gitleaks` action, 0 failures.

## 1. Authority and method

1. Confirm the V13/V14 merge commits, PR head trees and product trees.
2. Re-read the active category rows and the V14A reconciliation proposal.
3. Verify each proposed retirement against accepted source, tests and, where practical,
   an isolated runtime fixture.
4. Independently test two packages not in the V14A proposal (#1243, #1919) and the
   conditional #1900 retirement.
5. Remove only packages whose current accepted behavior satisfies the recorded request.
   A shared package remains active while any issue in it remains unresolved.
6. Add the retirement to `categories/exclusions.md` with source/test evidence, exact
   limitation and reopen condition.
7. Preserve `_backlog/data/issues.json`, `open-prs.json`, `crosswalk.csv`,
   `collection-manifest.json`, the original inventories and all historical snapshots.

An item removed here is **not** closed on GitHub. No upstream PR was merged or
cherry-picked. No RNC product source was changed.

## 2. Verified baseline

| Fact | Value |
|---|---|
| PR #6 head | `71be6f3fdf99ac711e7d1dbb5ee1683147ff2883` |
| PR #6 merge | `1e88ac8d3a299c6670a57bb229fb9b5ba811e607` |
| PR #6 head/merge tree | `c1686124229b7751be58851211cb379341d45e27` (identical) |
| PR #7 head | `cb7b0582159ac386fc0883ab47a63e2851a1b1ed` |
| PR #7 merge | `7d74921841114be7677a423cf60620210e832bf7` |
| PR #7 head/merge tree | `8eb327dd5860a4665daf38ce2f2fbc0aa5c07963` (identical) |
| PR #6 ancestor of #7 merge | yes |

The squashed integration commits preserve the reviewed trees exactly. The previous V14A
proposal's “not merged” premise is superseded; it did not change any ledger and remains a
historical proposal.

## 3. Final disposition delta

| Measure | Before | Retired | After |
|---|---:|---:|---:|
| Active issues | 186 | **-6** | **180** |
| Active PRs | 144 | **-2** | **142** |
| Active work packages | 226 | **-6** | **220** |
| Dispositioned issues | 99 | **+6** | **105** |
| Dispositioned PRs | 108 | **+2** | **110** |
| Dispositioned packages | 162 | **+6** | **168** |
| Original accounting | 285 / 252 / 388 | 0 | **285 / 252 / 388** |

Checks:

```text
180 + 105 = 285 active/dispositioned issues
142 + 110 = 252 active/dispositioned PRs
220 + 168 = 388 active/dispositioned packages
```

This differs from the V14A proposal's 182/142/222 and 103/110/166 because two additional
corrections are present in the accepted tree: #1243 and #1919.

## 4. Dispositions applied

### 4.1 #1858 — SARIF result level (@bdovaz)

**Original request:** SARIF results omitted `level`.

**Accepted behavior:** opt-in `--sarif-level=none|note|warning|error`; empty preserves
historical omission. No per-rule severity model is invented.

**Evidence:** `report/sarif.go` (`Level`, `getResults`), `cmd/root.go` flag validation,
`cmd/sarif_level_test.go:TestCLISARIFLevelIsOptIn`, `report/sarif_test.go`. V13 source
`74a042e522c3d8646d39fd41e8a0ed88d00e5e83`; accepted tree `7d74921`.

**Limitation:** the level is constant across results; default remains omission. This is
an opt-in capability, not a per-rule severity system.

**Reopen:** a consumer requires a default level or a per-result severity model.

### 4.2 #1714 — local file links (@rgmz)

**Original request:** clickable local file locations for directory scans.

**Accepted behavior:** opt-in `--file-links` emits a console-only `File-Link:` URI for
non-Git findings; `--no-color` suppresses terminal hyperlinks; JSON/CSV/SARIF unchanged.

**Evidence:** `detect/utils.go:310-348`, `cmd/root.go` flag, `cmd/file_links_test.go`
(URI, line anchor, JSON isolation, no ANSI). V13 source `74a042e`.

**Limitation:** absolute host paths are printed only when opted in; not a machine-report
field; no Git remote link change.

**Reopen:** a link appears in a machine report or by default, or leaks a path in a
console context that should be redacted.

### 4.3 #811 and #827 — strict configuration validation (@svyotov, @sergiomarotco)

**Original requests:** unknown fields silently ignored (#811); malformed configuration
validation (#827).

**Accepted behavior:** default loading remains permissive. Opt-in `--strict-config`
recursively rejects unknown root/rule/allowlist/extend/required fields and validates
extension chains and inline TOML. Malformed TOML already fails in the base parser; strict
mode adds contextual schema validation without echoing values.

**Evidence:** `config/strict.go`, `cmd/root.go:validateStrictConfiguration`,
`cmd/strict_config_test.go` (unknown root/rule, extension chain, inline root with
extension, malformed/misplaced allowlist), V14 source `b41ebe9` and `75e9c0a`.

**Accounting treatment:** **issue records only.** The #811/#827/#1279 shared
strict-configuration package remains active because #1279's default permissiveness,
allowlist-placement diagnostics and reported false positive are not fully resolved.

**Reopen:** strict mode fails to reject an unknown field, or default loading changes
without a compatibility decision.

### 4.4 #1900 — profile script measurement validity (@vignesh1507)

**Original proposal:** the profile script changed command flags and could abort on a
finding.

**Accepted behavior:** the script uses the same command array for benchmark and profile:
`--exit-code=0 --max-decode-depth "$decode_depth" --no-banner --log-level error`, quoted
paths, explicit decode depth and metadata. The `profile.sh` rewrite is independent of the
upstream patch; no upstream commit was adopted.

**V15 evidence (manual, 2026-11-15):** Git Bash, scanner built from the accepted HEAD,
Git target with a synthetic finding, target path containing a space, `DECODE_DEPTH=1`.
The script exited 0, did not stop on the finding, and produced CPU/memory/trace artifacts
for both `dir` and `git`; metadata recorded the target, Git HEAD, shallow=false, Go
version and decode depth. `hyperfine` and `yq` were absent, so the benchmark wrapper and
pretty-printer were not exercised.

**Limitation:** there is **no automated shell regression test**. The functional claim is
verified for the executed path; the test gap remains an internal verification obligation.
Do not cite the V13 smoke alone as proof.

**Reopen:** a shell regression, changed profile/benchmark flags, or a finding that aborts
the run despite `--exit-code=0`.

### 4.5 #1809 — Azure DevOps SSH remote identity (@JoostVoskuil)

**Original request:** Azure DevOps SSH remotes did not produce finding links.

**Accepted behavior:** exact `git@ssh.dev.azure.com:v3/org/project/repo` maps to
`https://dev.azure.com/org/project/_git/repo`; credentials are stripped; platform
detection recognizes `dev.azure.com`/`visualstudio.com`. No substring host matching.

**Evidence:** `sources/git.go:1062-1100` (`normalizeGitRemoteURL`),
`sources/git_test.go:389:TestNormalizeGitRemoteURL` (Azure SSH v3, GitHub SSH, credential
strip, encoded GitLab path). V13 source `74a042e`.

**Limitation:** only the exact Azure v3 SSH shape; no arbitrary SSH host mapping.

**Reopen:** Azure link regression, credential retention, or a host-spoofing case.

### 4.6 #1243 — no-final-newline redaction (@SyedMonzAli)

**Original request:** redaction exposed the preceding line's password when the file had
no final newline.

**Accepted behavior:** V11 line-start correction excludes the previous line's newline
from `Finding.Line`; known-exposure grouping masks other reported secrets for the same
file/commit. The V15 direct fixture produced `Finding: asana=REDACTED` at line 2, with no
preceding password.

**Evidence:** `f6c111c8` (V11 contributor line-start fix), `detect/location.go:59,73`,
`detect/utils.go:166-217`, `report/finding.go:116-221`; manual no-final-newline fixture.

**Limitation:** only secrets known to the detector are masked; unknown adjacent text can
remain. That is a separate redaction residual, not this no-newline defect.

**Reopen:** an adjacent-line exposure is reproduced with a controlled multi-secret
fixture.

### 4.7 #1919 — `targetRules` on global allowlists (@maksym-dilanian)

**Original request:** `targetRules` in `[[allowlists]]` was ignored.

**Accepted behavior:** the loader stores rule-targeted allowlists and attaches them after
extension, validating that the target rule exists; non-targeted rules are unaffected.
This is upstream code already present in the accepted baseline (`4451b455`), not an RNC
V14 change.

**Evidence:** `config/config.go:273-282`; V15 two-rule path/regex fixtures suppressed only
the targeted `alpha` finding and preserved `beta`.

**Limitation:** malformed/mis-scoped allowlist tables still need #1279's strict
diagnostics; this retirement covers the target-rule routing defect only.

**Reopen:** targetRules is ignored or over-applies under extend/path/inline config.

## 5. Category totals after reconciliation

| Category | Active issues | Active PRs | Active packages |
|---|---:|---:|---:|
| Security and hardening | 3 | 1 | 4 |
| Scan integrity and file readers | 14 | 8 | 16 |
| Detection accuracy | 49 | 26 | 51 |
| Provider rules | 35 | 60 | 61 |
| Configuration and allowlists | 29 | 18 | 29 |
| Git and hooks | 9 | 2 | 9 |
| Reporting and output | 6 | 3 | 7 |
| Performance and stability | 2 | 1 | 3 |
| Dependencies and build | 6 | 2 | 8 |
| Integrations and releases | 4 | 3 | 6 |
| Optional features | 20 | 13 | 22 |
| Documentation | 3 | 5 | 4 |
| **Active total** | **180** | **142** | **220** |
| **Dispositioned total** | **105** | **110** | **168** |
| **Original** | **285** | **252** | **388** |

## 6. Contributor and evidence preservation

- #1900: @vignesh1507 — design influence; independent RNC rewrite, no upstream shell
  code adopted.
- #1809: @JoostVoskuil — Azure SSH research; independent current-architecture
  implementation.
- #1243: @SyedMonzAli — defect research; V11 line-start correction.
- #1919: @maksym-dilanian — targetRules report; upstream implementation present.
- #1858 and #1714: issue reporters' requests remain credited through the upstream links;
  the RNC implementation is opt-in and independently authored.
- #811/#827: strict-validation requests; RNC opt-in strict mode, shared package retained
  for #1279.

Canonical code-attribution remains `RNC-MOD.md` §3. This reconciliation adds no
code-authorship claim and closes no upstream item.

## 7. Preservation check

Unchanged:

- `_backlog/data/issues.json`
- `_backlog/data/open-prs.json`
- `_backlog/data/crosswalk.csv`
- `_backlog/data/collection-manifest.json`
- `_backlog/01-issue-backlog.md`
- `_backlog/02-open-pull-requests.md`
- `_backlog/03-pr-issue-crosswalk.md`
- `_backlog/crosswalk-duplicates.md`
- historical reconciliation records

Edited:

- active category rows and headers for Security, Reporting, Performance, Git and
  Configuration;
- `_backlog/categories/README.md`;
- `_backlog/categories/exclusions.md`;
- `_backlog/README.md`;
- this record.

No scanner source, rule, default, test, module, CI workflow or release artifact was
changed by this reconciliation.

## 8. Remaining uncertainty

- #1900 has no permanent automated shell test; add one in the verification gate.
- The strict-config package remains open for #1279; #811/#827 issue retirement is not a
  full close of default permissiveness.
- #1858's SARIF level is constant and opt-in; no per-rule severity exists.
- #1714 prints absolute local paths only when explicitly requested.
- #1243's underlying “unknown adjacent text” residual remains separate.
- #1919 was an upstream implementation, not an RNC-originated fix; do not describe it as
  V14 work.

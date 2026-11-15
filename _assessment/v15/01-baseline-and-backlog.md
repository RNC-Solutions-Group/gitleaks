# CP-V15 — Accepted baseline and reconciled backlog

## 1. Repository authority

**Repository:** `RNC-Solutions-Group/gitleaks`
**Accepted integration branch:** `rnc-mod`
**Accepted HEAD:** `7d74921841114be7677a423cf60620210e832bf7`
**Accepted tree:** `8eb327dd5860a4665daf38ce2f2fbc0aa5c07963`
**Post-merge CI run:** `43763741641` (2026-11-15T21:26–21:41Z)

PR #6 (`CP-V13`) and PR #7 (`CP-V14`) are both merged. The ancestry and tree checks are
exact, not inferred from the PR titles:

| Fact | Verified value |
|---|---|
| PR #6 head | `71be6f3fdf99ac711e7d1dbb5ee1683147ff2883` |
| PR #6 merge commit | `1e88ac8d3a299c6670a57bb229fb9b5ba811e607` |
| PR #6 head tree | `c1686124229b7751be58851211cb379341d45e27` |
| PR #6 merged-commit tree | `c1686124229b7751be58851211cb379341d45e27` |
| PR #7 head | `cb7b0582159ac386fc0883ab47a63e2851a1b1ed` |
| PR #7 merge commit | `7d74921841114be7677a423cf60620210e832bf7` |
| PR #7 head tree | `8eb327dd5860a4665daf38ce2f2fbc0aa5c07963` |
| PR #7 merged-commit tree | `8eb327dd5860a4665daf38ce2f2fbc0aa5c07963` |
| `71be6f3` tree vs `1e88ac8` tree | identical |
| `cb7b058` tree vs `7d74921` tree | identical |
| PR #6 ancestor of PR #7 merge | yes |

The squash merges preserved the reviewed trees exactly. V13's merge parent is
`3a4635f`; V14's merge parent is `1e88ac8`. No unrelated source change was introduced by
either merge.

### CI status at the accepted HEAD

GitHub check runs at `7d74921`:

| Check | Result |
|---|---|
| Linux build and tests (race + real pre-commit + generator diff) | SUCCESS |
| Windows build and tests (race) | SUCCESS |
| macOS Intel and Apple Silicon build/race | SUCCESS |
| Standard regexp application build | SUCCESS |
| gore2regex application build | SUCCESS |
| Docker amd64/arm64 root and optional non-root build/smoke | SUCCESS (4 jobs) |
| Dependency, module and `govulncheck` verification | SUCCESS |
| Container vulnerability assessment | SUCCESS |
| Module graph update | SUCCESS |
| `gitleaks` action workflow | SKIPPED — intentionally not an RNC CI gate |

**13 successful check runs; 1 skipped check; 0 failed.** This is a post-merge success,
not a PR-only claim. The skipped `gitleaks` action is an intentional policy exclusion,
not a passing security scan.

The V14A statement that PR #7's CI had not run is now historical. The stacked PR did
receive a full Test-workflow run; the authoritative result is the post-merge run above.
`RNC-MOD.md` still describes PRs #6/#7 as open and the product as an unmerged candidate
(`RNC-MOD.md:19-24`, `2.6.6`/`2.7`). That metadata is stale and should be corrected in a
separate documentation checkpoint; it does not change the accepted source tree.

## 2. What was read

Current code and operations:

- `cmd/`, `detect/`, `sources/`, `config/`, `report/`, `regexp/`, `internal/`,
  `scripts/`, `Dockerfile`, `.goreleaser.yml`, `.github/workflows/`.
- `RNC-MOD.md`, `RELEASING.md`, `CI-SECURITY.md`, `KEY-MANAGEMENT.md`, `README.md`.
- V8–V14 assessment records, V13 remediation matrix and V14 policy/closeout documents.
- Backlog inventories, category ledgers, crosswalk, exclusions, V12 reconciliation and
  the V14A post-merge proposal.
- Operator policy record `2026-11-14-v13-policy.md`.

Historical proposals and future plans were not treated as accepted behavior. In
particular, `/v2` module identity, fingerprint v2, first-parent default, live credential
validation, remote configuration, multi-format reporting, and all release attestation
plans remain **not implemented**.

## 3. Backlog reconciliation

### 3.1 Independent result

The V14A proposal is correct in direction but incomplete. Two additional packages are
already corrected in the accepted product and were missed by that proposal:

- **#1243** — no-terminal-newline redaction exposure. The V11 line-start correction at
  `f6c111c8` excludes the previous line's newline from `Finding.Line`; accepted
  `detect/location.go:59,73` and `detect/utils.go:166-217`. A direct 2026-11-15
  reproduction produced `Finding: asana=REDACTED`, not the reported preceding-line
  password. **Retire the issue and package slot.**
- **#1919** — `targetRules` on a global allowlist. The accepted loader implements
  rule-targeted allowlists at `config/config.go:273-282` (upstream code from
  `4451b455`). A two-rule path/regex fixture suppresses only the targeted rule; the
  non-targeted rule still reports. **Retire the issue and package slot.**

The #1900 conditional retirement was independently tested, not accepted from the prior
smoke claim. `scripts/profile.sh` was run with the accepted binary against a Git target
containing a synthetic finding, under a path with spaces, with `DECODE_DEPTH=1`:

- exit 0 despite the finding;
- `dir` and `git` both produced CPU, memory and trace artifacts;
- metadata recorded target, Git head, shallow state, Go version and decode depth;
- benchmark and profile use the same command array (`--exit-code=0`,
  `--max-decode-depth`, `--no-banner`, `--log-level error`);
- `hyperfine`/`yq` were absent, so the benchmark wrapper and pretty-printer were not
  exercised.

There is still **no automated shell regression test**. That is a test-coverage residual,
not an open functional defect. Retire #1900's review slot and add an internal test
obligation (see the roadmap's verification gate).

### 3.2 Retirement delta

| Upstream item | Kind | Why retired | Evidence | Reopen trigger |
|---|---|---|---|---|
| #1243 (@SyedMonzAli) | issue + package | V11 line-start correction; manual no-final-newline repro | `f6c111c8`; `detect/location.go`; `detect/utils.go`; direct CLI run | New adjacent-line redaction exposure with a verified fixture |
| #1858 (@bdovaz) | issue + package | Opt-in `--sarif-level` implements the requested capability; empty preserves historical omission | `74a042e`; `report/sarif.go`; `cmd/sarif_level_test.go` | Consumer requires a per-rule severity model or default level |
| #1714 (@rgmz) | issue + package | Opt-in console-only local `file://` link; machine reports unchanged | `74a042e`; `detect/utils.go`; `cmd/file_links_test.go` | Link is emitted in machine output or leaks by default |
| #1900 (@vignesh1507) | PR + package | Independent profile rewrite; manual finding/decode-depth/artifact test | `74a042e`; `scripts/profile.sh`; 2026-11-15 manual smoke | Shell regression or benchmark/profile command drift |
| #1809 (@JoostVoskuil) | PR + package | Independent exact Azure SSH v3 normalization | `74a042e`; `sources/git.go:1062-1100`; `sources/git_test.go:389` | Azure link regression or unsafe substring matching |
| #1919 (@maksym-dilanian) | issue + package | Accepted loader implements `targetRules`; direct test | `4451b455` baseline; `config/config.go:273-282`; 2026-11-15 CLI test | `targetRules` regresses under extend/path semantics |
| #811 (@svyotov) | issue record only | Opt-in strict recursive field validation | `b41ebe9`; `config/strict.go`; `cmd/strict_config_test.go` | Strict mode fails to reject an unknown field |
| #827 (@sergiomarotco) | issue record only | Malformed TOML still fails; strict mode adds contextual validation | `b41ebe9`; strict tests | Default parser accepts malformed syntax |

The strict-configuration **package slot stays active** because #1279 remains unresolved.
Retiring #811/#827 reduces issue count but not the package count.

### 3.3 Accounting

| Measure | V14A proposal | Verified V15 | Delta vs V14A |
|---|---:|---:|---:|
| Active issues | 182 | **180** | -2 (#1243, #1919) |
| Active PRs | 142 | **142** | 0 |
| Active work packages | 222 | **220** | -2 |
| Dispositioned issues | 103 | **105** | +2 |
| Dispositioned PRs | 110 | **110** | 0 |
| Dispositioned/retired packages | 166 | **168** | +2 |
| Original issues/PRs/packages | 285/252/388 | **285/252/388** | 0 |

Checks:

- active issues + dispositioned issues = **180 + 105 = 285**
- active PRs + dispositioned PRs = **142 + 110 = 252**
- active packages + dispositioned packages = **220 + 168 = 388**

### 3.4 Category totals after reconciliation

| Category | Issues | PRs | Packages |
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
| **Total active** | **180** | **142** | **220** |
| **Dispositioned** | **105** | **110** | **168** |

The original dates, upstream identities, contributor links, JSON/CSV exports and
historical inventories are unchanged. Backlog changes are limited to the active category
ledgers, the exclusion register, the backlog indexes/plan note and the new V15
reconciliation record; the V15 assessment documents are additive.

## 4. Interpretation limits

- A package retirement means **removed from RNC's active review queue**, not fixed or
  closed upstream.
- The V15 counts are **review obligations**, not 220 code changes. Planning triage
  recommends 43 deferrals, 35 rejections and 52 targeted reproductions; see the matrix.
- `#811/#827` are corrected only in the sense that their acceptance condition is
  implemented as an opt-in strict mode. Default loading remains permissive; the shared
  strict package stays open for #1279.
- No V15 change closes a GitHub issue or PR.
- The current RNC product is still not a published release, and #2017's constrained-RE2
  release behavior remains unproven.

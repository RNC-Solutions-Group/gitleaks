# CP-V14A Backlog Reconciliation Proposal

**Status:** PROPOSAL ONLY — not applied to the canonical backlog.
**Canonical backlog:** `origin/feat/issues-and-pr` at `7f7e4eada3bdf3b594a3c0ec1fd101b89f2f08e2`.
**Code under review:** PR #7, `feat/v14-policy-implementation`, stacked on open PR #6 / `audit/v13-integration@71be6f3`. Neither PR is merged.
**Accounting baseline:** 285 issues / 252 PRs / 388 work-package slots; 186 active issues / 144 active PRs / 226 active slots; 99 / 108 / 162 dispositioned.
**Rule:** no V13 or V14 package is an accepted product baseline until the operator merges PR #6 and PR #7. No historical `_backlog/data/` export, inventory or canonical category file was changed.

## 1. Method and guardrails

- Read the canonical category register, exclusion ledger, V12 reconciliation, V13 matrix and V14 evidence at the commits above.
- Match an active upstream package only where the RNC behavior, source path and regression test address the recorded use case.
- A **package-slot retirement** requires all original issue/PR records in that package to be resolved. An issue-only retirement does not reduce the slot count while another issue keeps the package active.
- An upstream PR superseded by an independent RNC implementation is retired as a **review slot**, not represented as adopted code. Original contributor credit remains in `RNC-MOD.md` §3.
- Partial controls stay active. An opt-in diagnostic is not a full correction of a default-behavior request.
- Nothing below closes an upstream issue or PR on GitHub.

## 2. Proposed disposition delta

| Measure | Canonical | Proposed delta | Proposed after merge |
|---|---:|---:|---:|
| Active issues | 186 | **-4** | **182** |
| Active PRs | 144 | **-2** | **142** |
| Active work-package slots | 226 | **-4** | **222** |
| Dispositioned issues | 99 | **+4** | **103** |
| Dispositioned PRs | 108 | **+2** | **110** |
| Dispositioned/retired slots | 162 | **+4** | **166** |
| Original totals | 285 / 252 / 388 | 0 | **285 / 252 / 388** |

The difference between four retired issues and four retired slots is intentional: #811 and #827 share the strict-configuration package with #1279, which remains active. Their issue records can be retired, but the package slot cannot.

## 3. Fully corrected — proposed issue retirements

| Upstream item | V13/V14 change | Exact source / test evidence | Commit(s) | Relationship / accounting |
|---|---|---|---|---|
| [#1858](https://github.com/gitleaks/gitleaks/issues/1858) — SARIF `level` omitted | Additive opt-in `--sarif-level=none|note|warning|error`; empty preserves historical omission; no per-rule severity invented. | `report/sarif.go`, `cmd/root.go`; `cmd/sarif_level_test.go:TestCLISARIFLevelIsOptIn`; `report/sarif_test.go`. | `74a042e522c3d8646d39fd41e8a0ed88d00e5e83` | Standalone Reporting package. Retire issue and one slot. |
| [#1714](https://github.com/gitleaks/gitleaks/issues/1714) — local `file://` links | Opt-in console-only `File-Link:` URI for non-Git findings; `--no-color` suppresses terminal hyperlinks; machine reports unchanged. | `detect/utils.go`, `cmd/root.go`; `cmd/file_links_test.go:TestCLIFileLinksConsoleOnly`; URI/JSON-isolation assertions. | `74a042e522c3d8646d39fd41e8a0ed88d00e5e83` | Standalone Reporting package. Retire issue and one slot. |
| [#811](https://github.com/gitleaks/gitleaks/issues/811) — unknown fields silently ignored | Opt-in recursive `--strict-config` rejects unknown root/rule/allowlist/extend fields with file/table/field context; default loading remains permissive. | `config/strict.go`, `cmd/root.go`; `cmd/strict_config_test.go:TestCLIStrictConfigRejectsUnknownFields`, `TestCLIStrictConfigRejectsUnknownRuleField`, `TestCLIStrictConfigValidatesExtensionChain`. | `b41ebe9` (schema/CLI); `75e9c0a` (inline-root correction) | Shares the three-issue strict package with #827/#1279. Retire the issue record; **do not reduce the package slot** while #1279 remains active. |
| [#827](https://github.com/gitleaks/gitleaks/issues/827) — malformed configuration validation | Malformed TOML already fails; strict mode adds recursive validation and contextual diagnostics without echoing values. Verify against the current parser rather than treating this as a new parser feature. | `config/strict.go`; `cmd/strict_config_test.go:TestCLIStrictConfigNestedAndSyntaxFailures` (`malformed` case); default parser tests. | `b41ebe9` | Same shared strict package. Retire the issue record as verified/superseded; no slot reduction while #1279 remains active. |

## 4. Already implemented or superseded — proposed PR retirements requiring post-merge verification

| Upstream PR | Relationship to RNC work | Exact source / validation evidence | Commit(s) | Accounting |
|---|---|---|---|---|
| [#1900](https://github.com/gitleaks/gitleaks/pull/1900) — `profile.sh` measurement validity | Independent RNC rewrite, not an adopted patch. Same command shape for benchmark/profile, explicit decode depth, `--exit-code=0`, safe quoting, metadata and artifact checks. | `scripts/profile.sh`; V13 Linux smoke produced CPU/memory/trace artifacts for `dir` and `git`. No automated Go test covers the shell script; retain as a post-merge verification item. | `74a042e522c3d8646d39fd41e8a0ed88d00e5e83` | PR-only Performance package. Retire PR and one slot. Contributor credit remains documentary. |
| [#1809](https://github.com/gitleaks/gitleaks/pull/1809) — Azure DevOps SSH identity | Independent current-architecture normalization; old patch targeted deleted files and used unsafe substring matching. RNC maps the exact `ssh.dev.azure.com:v3` form to HTTPS and strips credentials. | `sources/git.go`; `sources/git_test.go:TestNormalizeGitRemoteURL` (`azure ssh v3` case); URL-escaping/link tests. | `74a042e522c3d8646d39fd41e8a0ed88d00e5e83` | PR-only Git package. Retire PR and one slot. No upstream code adopted. |

`#1900` is the only proposed retirement without an automated regression test. Its source and smoke evidence are recorded in `_assessment/v13-remediation-matrix.md`; the operator should either accept that evidence or keep the PR slot open until a script test exists.

## 5. Partially corrected — retain active, with residual

| Upstream item(s) | Implemented portion | Retained residual / reason not retired |
|---|---|---|
| [#2130](https://github.com/gitleaks/gitleaks/issues/2130) + [#2244](https://github.com/gitleaks/gitleaks/pull/2244) | `--report-unused-ignore` reports loaded entries that suppressed no finding after a completed scan; incomplete scans never certify entries unused. `74a042e`; detector/CLI tests. | No `--fail-on-unused-ignores`/strict CI failure, which #2130 and #2244 request. The RNC diagnostic is stderr-only, not the PR's informative output. Keep both records active; PR is not adopted. |
| [#983](https://github.com/gitleaks/gitleaks/issues/983) + [#1654](https://github.com/gitleaks/gitleaks/pull/1654) | `gitleaks config` prints precedence, loaded files, rule/allowlist counts, report policy and limits without regex export. `74a042e`; summary/precedence tests. | This is a bounded summary, not the full resolved/merged TOML or a round-trippable config dump requested by #983 and implemented by #1654. Keep active. |
| [#1054](https://github.com/gitleaks/gitleaks/issues/1054) + [#1162](https://github.com/gitleaks/gitleaks/pull/1162) | Opt-in keyed HMAC grouped sidecar groups repeated values without suppressing or rewriting primary findings. `report/grouped.go`; `b41ebe9` + `75e9c0a`; repetition/redaction/key/destination tests. | It does not deduplicate the primary `detect -v` output or share one `.gitleaksignore` fingerprint across occurrences. #1162's unkeyed SHA-256 is explicitly rejected as an unsafe default; the PR is not adopted. Keep #1054 and the package active. |
| [#1279](https://github.com/gitleaks/gitleaks/issues/1279) | The loader rejects `[rules.allowlist]` without a rule ID (`rule |id| is missing or empty`); strict mode provides additional unknown-field/allowlist-placement diagnostics elsewhere. `config/strict.go`; `b41ebe9`; nested/misplaced allowlist tests. | Default non-strict loading remains permissive, the diagnosed error is not specific to allowlist placement, and the reported false positive still requires a corrected allowlist. Keep active and keep the shared strict package slot. |
| [#2003](https://github.com/gitleaks/gitleaks/issues/2003) + [#2128](https://github.com/gitleaks/gitleaks/pull/2128) | Shell-free quoted `--log-opts` tokenizer and parsed-vector validation fix quoted pathspecs and preserve fail-closed format checks. `sources/git.go`; `74a042e`; parser/vector tests. | The issue's notes/CSV symptom was not reproduced; the PR was not adopted. Do not retire #2003 on the quoting correction alone. |
| [#1310](https://github.com/gitleaks/gitleaks/issues/1310) | Shallow-repository warning clarifies missing pre-boundary history; SHAs and fingerprints are not rewritten. `sources/git.go`; shallow fixture test; `74a042e`. | Missing history cannot be reconstructed; misleading identity remains possible. Retain. |
| [#1163](https://github.com/gitleaks/gitleaks/issues/1163) | Subdirectory history-scope warning; repository-wide Git behavior is explicit. `74a042e`. | No subdirectory-history feature or default scope change. Retain. |
| [#1293](https://github.com/gitleaks/gitleaks/issues/1293) | Bare-repository warning and explicit config/ignore guidance. `74a042e`. | Automatic `.gitleaksignore` discovery from a bare Git tree is not implemented. Retain. |
| [#1754](https://github.com/gitleaks/gitleaks/issues/1754), [#764](https://github.com/gitleaks/gitleaks/issues/764) | No behavior change. Current `validateMinVersion` warns when the scanner is older; it does not reject solely from that comparison. | The requested fail-closed compatibility policy is an operator decision and remains unimplemented. See §7. |

## 6. Deferred or rejected proposals

| Item | Status | Reason |
|---|---|---|
| [#1048](https://github.com/gitleaks/gitleaks/issues/1048) multiple report formats | Deferred | Requires a validated multi-destination report plan and explicit partial-failure semantics; V13 documented the design but shipped no unsafe implementation. |
| [#1766](https://github.com/gitleaks/gitleaks/issues/1766) bounded finding context | Deferred | Requires a separate opt-in schema, byte/line budgets and redaction of every nearby secret. |
| [#2019](https://github.com/gitleaks/gitleaks/issues/2019) recursive-decoding OOM | Deferred / not reproduced | V13's finite 6.67 MB probe found no depth-correlated peak-RSS growth in that control; the reported 24 GB case is not established. Keep the resource residual. |
| [#2017](https://github.com/gitleaks/gitleaks/issues/2017) RE2 address-space under release constraints | Deferred release proof | No constrained release binary was reproduced; compiling the full tagged graph is not runtime proof. |
| [#1217](https://github.com/gitleaks/gitleaks/pull/1217), [#2220](https://github.com/gitleaks/gitleaks/pull/2220), [#885](https://github.com/gitleaks/gitleaks/issues/885) signing/SBOM/tags | Deferred release engineering | No image, attestation, SBOM or registry digest exists to verify. |
| [#1028](https://github.com/gitleaks/gitleaks/issues/1028) + [#2301](https://github.com/gitleaks/gitleaks/pull/2301) first-parent default | Operator-policy dependent | V14 measured real merge-only recovery but also duplicates, new fingerprints and baseline/ignore migration. Default unchanged; PR not adopted. |
| [#1839](https://github.com/gitleaks/gitleaks/issues/1839), [#1209](https://github.com/gitleaks/gitleaks/issues/1209), [#1226](https://github.com/gitleaks/gitleaks/issues/1226) module identity | Operator-policy dependent | `/v2` local experiment passed build/test/install/linker/cross-build, but no module was published and no `go.mod` change was made. |
| [#2239](https://github.com/gitleaks/gitleaks/issues/2239) + [#2286](https://github.com/gitleaks/gitleaks/pull/2286)/[#1566](https://github.com/gitleaks/gitleaks/pull/1566) fingerprint v2 | Design only | Keyed/versioned grammar and migration are documented; no v2 code or default switch. |
| [#1162](https://github.com/gitleaks/gitleaks/pull/1162) unkeyed `SecretSHA` | Rejected in favor of keyed grouping | A bare SHA-256 of a secret is an enumerable correlation identifier. It is not adopted; #1054 remains partial and active. |

## 7. Operator-policy-dependent items

1. **Module identity:** choose COA1 (retain upstream path) or COA2 (RNC `/v2` product `v2.0.0`). COA2 requires a separate `minVersion` compatibility-floor policy, module-proxy/checksum proof and downstream migration. No production change is authorized.
2. **First-parent default:** retain opt-in unless the operator accepts duplicate findings, new fingerprints and baseline/ignore migration, plus larger-history evidence.
3. **Fingerprint v2:** approve grammar, key source/rotation, dual-acceptance versus v2-only enforcement and migration before implementation.
4. **`minVersion`:** current comparison warns, not rejects. Decide whether upstream `v8.x` becomes a configuration-schema floor against an embedded compatibility version; do not conflate it with release-tag-base validation.
5. **Key lifecycle:** accept OS/CI secret storage, Windows ACL limits, in-memory lifetime and key-rotation behavior, or keep grouping disabled.
6. **Release evidence:** authorize and execute signed/SBOM/provenance publication and independent consumer verification; no release is claimed now.

## 8. Post-merge application procedure

1. Merge PR #6, then retarget PR #7 to `rnc-mod` and run the full GitHub Actions matrix.
2. Re-read the canonical category rows at that merged HEAD; do not apply this table mechanically.
3. For each retirement, add the RNC commit, source/test paths and reopen trigger to `_backlog/categories/exclusions.md`; remove the active row and reconcile category headers.
4. Retire #811/#827 issue records **without** moving the shared strict-config package until #1279 is separately resolved.
5. Do not edit `_backlog/data/issues.json`, `open-prs.json`, `crosswalk.csv` or historical inventory exports.
6. Record the final counts and verify `active + dispositioned = 285 / 252 / 388` before committing the canonical reconciliation.

## 9. Attribution and exclusions

No upstream patch is adopted by this proposal. #1809 and #1900 are superseded review slots, with design credit retained for @JoostVoskuil and @vignesh1507. #1162 is rejected for privacy reasons and credited to @mn3m0nic. #2286 (@breken-ai) and #1566 (@rgmz) remain design references. The canonical attribution register remains `RNC-MOD.md` §3; this document creates no new code-authorship claim.

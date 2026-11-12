# Wave 2 — RNC issue/PR reconciliation, harness and release roadmap

**Initiated:** 2026-11-12 · **Workspace:** `feat/issues-and-pr` (documentation-only triage)  
**Historical upstream snapshot:** 285 issues / 252 open PRs / 388 work-package slots, dated 2026-11-08–09.  
**Historical scanner baseline:** RNC `master@5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`.  
**Modern scanner comparator:** `fix/validation` at CP-V7A `57931abf68eb9eff799732962562ddab2156ced8`; re-pin after current CP-V8.  
**Existing owner/GPT register:** `fixes-qol-tracking.md` on `fix/validation` (not edited here).

## Reconciliation method

Review current-code status *with the operator*, one complete issue/PR work package at a time. Retain historic upstream author, issue/PR links, original bug or requested feature, original baseline findings, RNC implementation commit, current SHA, actual reproduction and regression evidence, compatibility decision, and any residual limitation.

Use explicit states: **implemented and verified**, **partial or conditional**, **reproduced open**, **not yet verified**, **deferred/optional**, **accepted intentional limitation**, **superseded/duplicate with evidence**. Do not change historical `_backlog/data/` JSON/CSV. A matching PR title, green CI, or already-reviewed old snapshot does not prove closure on the modern scanner. Changes on this triage branch do not authorize scanner modifications or merging.

## Joint P0–P3 package review

| Priority | Proposed package clusters | Required evidence |
| --- | --- | --- |
| **P0** | Original fail-open scan errors (CP-V0 F01); destructive report paths (F02); global AND allowlist bypass (F06); staged binary-marked text (F07); fragment-boundary losses (issue #1651 / PR #2324); Git parser failure edges | Reproduce old behavior and demonstrate modern correction, correct exit semantics and report preservation; no unmitigated false-clean or destructive write within supported paths |
| **P1** | CP-V8 merge-only history F11; source-path identity F16; repo-local config trust; residual sensitive logging/Git diagnostics; ignore-file mutation; Git pretty/date/textconv; Unicode/invalid-byte ingestion | Positive/negative, Windows/Linux/macOS and standard/RE2 fixtures; no silent default/fingerprint change without operator approval |
| **P2** | SARIF column/Unicode/URI; anchored line allowlists; cross-platform baselines; current dependency/RE2 compatibility; provider accuracy batches; paired-report late publication failures | Current-code reproduction, separate source/report invariants, bounded resource/error tests; avoid conflating optional format requests with active P2 defects |
| **P3** | Decoder throughput, per-file buffer allocation, redundant regex calls, deterministic order, documentation/CLI QOL | Controlled benchmarks plus byte/findings/exit parity; do not disable detection to improve speed |
| **Optional** | 22 parked packages, including 3 informative-report follow-ups (unused-ignore, effective config, report grouping) | Opt-in interface design, confidentiality, complexity and owner approval; not part of core P0/P1 defect closure |

## Post-V8 independent harness acceptance

1. **Pin artifacts:** upstream reference, original RNC baseline, and post-V8 candidate commit; executable/container digests, default config hash, standard/RE2 versions and Git version. Do not merge candidate to master merely to test.
2. **Add candidate lane to dedicated harness:** compare against the existing baseline lane, retaining all original fixtures and expected data; do not replace harness baseline. Include directory, staged Git, history/merges, archives, symlinks, and failures.
3. **Differential datasets:** known-positive/negative controls, Levain, RNC synthetic secrets lab (1,600 adjudicated items), trap negatives, multi-occurrence groupings, evidence spans, rule IDs, path/fingerprint parity, statuses, and actual exit codes.
4. **Informative reporting:** check disabled mode is original behavior, console opt-in prints stderr, file opt-in produces paired JSON, and skipped large/binary/allowlisted targets do not become scan errors. Verify output privacy and report collision protection.
5. **Release readiness:** compare modern `fix/validation/go.mod` versions with old CVE PRs #2308/#2336 and current advisories; correct fork `.github/workflows/gitleaks.yml` skipped-only upstream-ownership gate; test non-root Docker on read-only/writeable mounts and reconsider global Git `safe.directory '*'`; separate mutating/network Makefile targets from reproducible builds; validate RE2 full tagged graph limitation, macOS/Windows/Linux artifacts, SBOM/provenance and image signing where accepted.
6. **Disposition and approval:** document accepted regression baseline, remaining risks, contributor attribution (adopted vs adapted vs informed), operator behavior choices, and clean integration/release gate. No PR, merge, rebase, or history rewrite authorized here.

## Reconciliation decisions — reviewed and operator accepted

### F01 — Historical P0 false-clean source/Git/worker failures (2026-11-12)

**Decision: ACCEPT core CLI correction.** The original CP5-F01 and CP-V0-F01 behavior was reproduced against the original pinned baseline and addressed through CP-V3 commits `105c0a8`, `56458e7`, `fd996ec`, `b556701`. Current source as of CP-V7A `57931ab` routes supported `git`, `dir`, `stdin`, and hidden `protect` through corrected `DetectSource` / `DetectGitContext` handling. Directory read/walk errors, Git nonzero process exits/fatal stderr/parse failures, worker errors and cancellation propagate; incomplete scans return nonzero regardless of `--exit-code 0`, and existing reports are preserved.

**Evidence:** `sources/files_errors_test.go`, `sources/git_test.go`, `detect/scan_integrity_test.go`, `cmd/scan_status_test.go`, `cmd/scan_timeout_unix_test.go`, `cmd/staged_blob_failure_unix_test.go`, and [CP-V3 cross-platform CI 10/10](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43094102089) plus [CP-V7A CI 10/10](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43278900115). Source/CI evidence was independently inspected; the latest candidate was not rerun locally in this review.

**Related packages:** Scan Integrity [#2232, PRs #2235/#2249/#2313](categories/scan-integrity.md); Git/Hooks [#1981/#1464 and PRs #2080/#2292/#2339/#2252/#2253](categories/git-hooks-sources.md), and [#2129 / PR #2203](categories/git-hooks-sources.md) within tested stderr/drain scenarios. Prior upstream PRs were researched/informed rather than cherry-picked; preserve original contributor attribution.

**Residuals are separate, not an open core P0:** deprecated `DetectFiles` and `DirectoryTargets` still lose errors when used directly; pathological custom/truncated Git diff parser coverage is heuristic; blocking OS reads can outlast cancellation; 50-iteration adversarial Git stderr/queued-output test was not independently reproduced. Track the first as non-CLI API correctness, the second as parser-integrity regression, and the third under F19 cancellation. Final harness must rerun the accepted controls before release.

**Next:** F02 report-output safety has been reviewed and accepted with explicit residuals; see below. F01 and F02 are independent decisions.

### F02 — Historical P0 report destruction / output integrity (2026-11-12)

**Decision: ACCEPT original P0 correction; keep two publication limitations separately.** RNC V1/V2/V2A corrected the original #2304 / PR #2307 preflight deletion/truncation issue. Incomplete scans, bad destinations and report-writer failures preserve prior report contents; reject-existing, overwrite and aggregate policies are explicit, and config/baseline/ignore/source collisions are protected. Related #2299 / PRs #2300/#2312 plaintext-report permissions are addressed with restrictive temporary publication. Normal upstream primary overwrite semantics remain. `cmd/report.go`, `cmd/report_test.go`, `cmd/report_config_test.go`, `cmd/report_pair_test.go`; [CP-V2 10/10](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43075289894), [CP-V2A 10/10](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43084497923), [CP-V7A 10/10](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43278900115) independently verified **10/10 CI each**. Historical upstream PRs reviewed; no patch directly adopted in this reconciliation.

**Release residuals:** (a) a narrow final-publication TOCTOU interval can remain when another process mutates destinations/parent directories after validation, with no current exploit reproduced; (b) V7A companion-first overwrite may leave a newer informative file with the older primary report if the second rename fails, with explicit error and primary preservation. This does **not** imply incomplete source scanning was accepted as successful; it is consistency of independently published outputs. Retain adversarial filesystem tests and consumer mismatch handling, Windows ACLs under real identities. Avoid a two-file transaction redesign absent evidence.

**Informative timestamp — operator approved requirement:** Add only to opt-in informative output (JSON plus console summary): propose `scan.generated_at` ISO 8601/RFC3339 **UTC** once the scan completes and document is assembled. No change to the primary report, findings, flags, console-only behavior or exit codes. The date helps distinguish file ages but cannot prove paired-file atomicity. Queue for next QOL implementation after CP-V8/harness; test injected-clock formatting and additive JSON schema expectations.

### F06 — Historical P0 global AND-allowlist path bypass (operator accepted 2026-11-12)

**Source/code assessment:** RNC V1-B implements `sources/common.go:pathIsSufficientToSkip` to forbid path-only early exclusion when global AND lists also require regex, stopword or commit criteria. `sources/common_test.go` includes AND path+regex, AND path+stopword, AND path+commit, OR, path-only AND and nil config cases. `sources/files_test.go` confirms targets are enumerated; `detect/detect_test.go` covers AND/OR detector combinations and commit/path. Current CP-V7A CI is green 10/10, but this review did not independently rerun final binary end-to-end CLI cases. The scan-root change in ongoing V8 requires follow-up parity checks; upstream #2227 attribution/patch overlap still to be confirmed before clean release.

**Decision: ACCEPT F06 original bypass correction within supported core scan paths.** Operator approved 2026-11-12. Retain full end-to-end CLI AND/OR negative-condition, Windows/Unix, standard/RE2, and post-CP-V8 F16 path-identity checks as release evidence. **Separate concern:** #2165/#2166 global Git commit-allowlist loop behavior in `sources/git.go` still requires independent P1 review; it is not the closed F06 defect.

### F07 — Historical P0 Git-staged binary-marked text bypass (technical review 2026-11-12; owner acceptance pending)

**Original issue:** [#2337](https://github.com/gitleaks/gitleaks/issues/2337) (@farrrr), proposed [PR #2338](https://github.com/gitleaks/gitleaks/pull/2338) (@VISHNU0906). Git can label readable staged text as binary through `.gitattributes -diff` or `core.bigFileThreshold` so upstream diff parsing lacks text hunks. The original P0 omitted staged secrets without signaling a scan failure.

**Current RNC implementation:** CP-V4 fallback in `sources/git.go` detects staged binary-marked diffs and reads **stage-zero Git index content** (`git cat-file blob :0:<path>`), not the worktree. Existing `sources/file.go` MIME classification, archive-depth/size policy and true binary exclusions remain. Failing to open/read/cancel index blobs produces an operational scan error via the CP-V3 collector and does not overwrite existing reports. CP-V5 `CloseAfterSkip` stops the producer for intentionally skipped large/binary blobs rather than draining the entire object. History and unstaged binary policies are unchanged and should **not** be conflated with the staged P0 correction.

**Positive/negative evidence:** `detect/git_staged_test.go`: normal text, `-diff`, `core.bigFileThreshold`, index/worktree divergence both directions, real binary, deletion, rename, spaces/Unicode, staged archives and size limits. `sources/git_index_test.go`: missing index object, cancellation, special path, read/close failures, intentional early skip. `cmd/staged_status_test.go`: real staged CLI findings, JSON report and exit-code overrides; `cmd/staged_blob_failure_unix_test.go`: missing blob fails nonzero even with `--exit-code 0`, previous report sentinel untouched. Independently verified workflow results: [CP-V4 CI 10/10](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43099917379), [CP-V5 CI 10/10](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43109905591), [CP-V7A CI 10/10](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43278900115). Modern `fix/validation` at `cfb2815` retains the reviewed staged Git implementation. No independent external harness run in this reconciliation.

**Proposed disposition:** ACCEPT the historical staged-binary-text P0 as corrected, conditional on routine post-V8 harness regression and documented exclusions; no new scanning mode or default history behavior. Verify CI/native/Docker standard and RE2 builds, real Git staged index/worktree divergence, bigFileThreshold and `-diff`, true binary negative, missing blob failure, max-file-size intentional skip and no default report/exit regressions. Upstream PR #2338's **design was adapted**, not cherry-picked; preserve @VISHNU0906 and reporter @farrrr attribution, with final exact patch overlap review for clean release. Performance of large skipped blobs should be monitored; no claim of universal binary/archived-history coverage.

**Next proposed joint review:** [#1651](https://github.com/gitleaks/gitleaks/issues/1651) / [PR #2324](https://github.com/gitleaks/gitleaks/pull/2324) chunk-boundary false negatives, a potentially still-active P0/P1 package requiring current-branch reproduction. No implementation change or backlog defect closure is claimed for that package.

## First review session

Begin with original **P0** F01/F02/F06/F07 and the still-unreconciled #1651/#2324 chunk-boundary package, then resolve report publication TOCTOU and CP-V8 F11/F16 decisions. Update category statuses only after the operator reviews the current evidence. This file is a kickoff plan, not a completed reconciliation.

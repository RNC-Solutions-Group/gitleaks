# Fixes, quality-of-life improvements, and upstream attribution

> **Living engineering and release-preparation ledger.** This is not a claim that every item below is complete or that an upstream author's code has been incorporated. Update it as implementation, verification, and attribution decisions change.

**Project:** RNC Gitleaks fork (`RNC-Solutions-Group/gitleaks`)  
**Working branch:** `fix/validation`  
**Snapshot:** 2026-11-12; CP-V7 implementation at `cccdef68731e2c59fa4af2d246a5927a111ffacd` reviewed and CI-validated. F09 accepted within tested scope; F15 remains partially open for nonseekable ZIP/7z outer-container behavior. Previous ledger-only commit `96e5f0991011ee64c234805a93998fc9881a1831`.
**Comparison baseline:** `dev/gitleaks-upgrade` at `c0516b4bd66c0dc3878d80283a267166f7b4fd09`  
**Upstream reference reviewed for CP-V0:** `gitleaks/gitleaks` `master` at `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`  
**Primary evidence:** CP-V0 and CP-V1–V3 evidence; now additionally **full operator-uploaded CP-V4, CP-V5, CP-V6, and CP-V7 reports**, GitHub source and commit comparisons, CI job verification, and upstream PRs/issues. Agent-reported local tests/performance remain distinct from independently verified source/CI results; historical PR SHAs are not release pins.

**CP-V2 snapshot:** implementation at `452d1224afe57b68888b4e581ad70878758426ad`; CI [43075289894](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43075289894) passed all jobs.

**CP-V2A snapshot:** implementation at `cf3007301441eaf1f7b449707c87c2e87cbfebc8`; CI [43084497923](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43084497923) passed all jobs.

**CP-V3 snapshot:** implementation at `b556701f46cf1955d9b18b107f8a771a28beebb5` (four commits from `f56fd00`); CI [43094102089](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43094102089) passed all 10 jobs: Windows, Linux, both macOS architectures, standard/RE2 builds, module/vulnerability verification, Docker amd64/arm64, and container vulnerability assessment.


**CP-V4 snapshot:** staged-index/binary classification fix at [`266415418f1872ba77da7025789c199d37e9d8e0`](https://github.com/RNC-Solutions-Group/gitleaks/commit/266415418f1872ba77da7025789c199d37e9d8e0); [CI 43099917379](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43099917379), **10/10 jobs successful**.

**CP-V5 snapshot:** redaction/logging/large-blob lifecycle fix at [`877900749d422302b15a51db69f318bce0cc70ff`](https://github.com/RNC-Solutions-Group/gitleaks/commit/877900749d422302b15a51db69f318bce0cc70ff); four commits from CP-V4; [CI 43109905591](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43109905591), **10/10 jobs successful**.

**CP-V6 snapshot:** configuration/ignore-file fix at [`d0795f59af188284c6e146d97cd0bd1452389689`](https://github.com/RNC-Solutions-Group/gitleaks/commit/d0795f59af188284c6e146d97cd0bd1452389689); two commits from CP-V5; [CI 43114548814](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43114548814), **10/10 jobs successful**. The configuration-trust policy decision remains open.

**CP-V7 snapshot:** [implementation `cccdef6`](https://github.com/RNC-Solutions-Group/gitleaks/commit/cccdef68731e2c59fa4af2d246a5927a111ffacd) (13 changed files, +862/−41); [CI 43121770371](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43121770371) independently verified **10/10 successful jobs**. **F09 accepted for tested false-negative behavior; F15 not fully closed** pending the oversized nonseekable ZIP/7z policy and regression decision.

## 1. Project contract: improve Gitleaks without changing how people use it

1. **Backward compatibility is a release gate.** Preserve upstream commands, flags, configuration interpretation, default scanning behavior, exit codes, finding content, output formats, and automation contracts unless correcting a demonstrated defect requires a change.
2. **Accuracy and completeness come first.** Never trade real detection coverage for speed, a superficially green test, or a nicer report.
3. **Safety fixes should be transparent.** Prevent unintended data loss, misleading clean results, partial reports, and unsafe permissions without imposing a new workflow on established users.
4. **New features are opt-in whenever practical.** Keep documented defaults compatible with original Gitleaks. Feature flags override configuration only when explicitly supplied.
5. **No unnecessary product expansion.** Gitleaks is a scanner, not a retention, archive, or report-management service. Do not mix in the separate benchmark harness or security orchestration toolkit.
6. **Preserve contributor credit.** Research influence, adapted code, and incorporated commits are different forms of contribution; record them honestly and distinctly.
7. **Keep the public history clean.** Iterative development commits are valuable on working branches but should not become hundreds of unrelated commits on the final fork baseline.

**Compatibility evidence required before release:** Run representative synthetic fixtures with upstream and the final fork using identical default settings. Compare detection results, report schemas, CLI contracts, exit codes, and basic runtime behavior. Document each intentional difference; do not assume compatibility merely because unit tests pass.

## 2. Improvement and change register

Statuses: **Implemented** = code exists on `fix/validation`; **Validated** = specific tests/evidence exist; **Planned** = approved design, not yet on branch; **Open** = not fixed; **Deferred** = outside current cycle. A local pass is not a claim that all CI platforms pass.

| ID | Area | Original defect or behavior | Change / expected contract | State at snapshot | Relevant evidence / attribution |
| --- | --- | --- | --- | --- | --- |
| V1-A / CPV0-F02 | Report integrity | Early `os.Create`/remove could destroy old output; final write could truncate or leave partial output. | Validate first; write temporary restrictive-permission report; publish atomically; preserve old output on failed validation/publication; reject unsafe targets. | **Implemented and validated through CP-V2 CI.** | [#2304](https://github.com/gitleaks/gitleaks/issues/2304), [#2307](https://github.com/gitleaks/gitleaks/pull/2307), [#2299](https://github.com/gitleaks/gitleaks/issues/2299), [#2300](https://github.com/gitleaks/gitleaks/pull/2300), [#2312](https://github.com/gitleaks/gitleaks/pull/2312) |
| V1-B / CPV0-F06 | AND allowlists | A path match could suppress a file before other required AND criteria were evaluated. | Only path-skip if OR or path-only AND; evaluate content/commit-dependent AND conditions during detection. | **Implemented; tests reported passing.** | [#2227](https://github.com/gitleaks/gitleaks/pull/2227) |
| V1-C / CPV0-F08 | Git diff integrity | Forced terminal coloring could break machine-readable Git log/diff parsing and silently miss findings. | Use Git `--no-color` for supported history/staged/unstaged acquisition; validate forced-color fixtures. | **Implemented; tests reported passing.** | [#2284](https://github.com/gitleaks/gitleaks/pull/2284) |
| V1-D / CPV0-F10 | File-chunk classification | Later chunks could reclassify an initially text file as binary and drop findings. | Classify MIME on the initial chunk only; preserve binary-at-start behavior. | **Implemented; tests reported passing.** | [#2329](https://github.com/gitleaks/gitleaks/pull/2329) |
| V1-E / CPV0-F12 | `dir` CLI | Extra/missing source arguments could silently fall back to scanning CWD. | Accept zero (default `.`) or one valid source; reject ambiguous/missing paths. | **Implemented; tests reported passing.** | [#1979](https://github.com/gitleaks/gitleaks/issues/1979), [#2334](https://github.com/gitleaks/gitleaks/pull/2334) |
| V1-F / CPV0-F13 | Pre-commit integration | Wrapper could interpret nonzero scanner exits as success, particularly on Windows. | Use portable subprocess status; fail on every nonzero code, with useful diagnostics. | **Implemented; tests reported passing.** | No matching upstream PR identified in CP-V0. |
| V1G | Path/collision hardening | Development-only in-tree report policy added path-resolution and cross-platform failures. | Retired by the V2 output policy. The replaced policy, content recognition, sole-file guard, and `--allow-in-tree-report-overwrite` are gone. | **Retired by V2; superseded.** | CP-V1G commits through `4648852`; V2 changes in `452d122`. |
| V2-A | Output compatibility | Default upstream behavior overwrites the requested report. | `auto-overwrite = true` by default; keep the original output path and normal overwrite when aggregation is disabled. | **Implemented and validated.** | `config/config.go`, `cmd/report.go`, `cmd/report_test.go`; CI run [43075289894](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43075289894). |
| V2-B | Optional aggregation | Repeated output paths collide. | `aggregation = false` by default; when enabled, create `report.json`, then `report_001.json`, etc. No content merge or retention. | **Implemented and validated.** | Numbering, concurrent allocation, extension, unicode, and template tests in `cmd/report_test.go`; CI run above. |
| V2-C | Explicit overrides | Users need one-run exceptions without editing configuration. | Add `--overwrite` and `--aggregate`; CLI wins over settings; both flags together error. | **Implemented and validated.** | Matrix tests cover all four configurations and both flags; CI run above. |
| V2-D | Collision rejection | Users may prohibit both overwrite and aggregation. | If both settings are false and the destination exists, error **before scanning**; absent destination writes normally. | **Implemented and validated.** | CLI tests assert no scan summary, preserved report, and both overrides named; CI run above. |
| V2-E | Output simplicity | Experimental `--allow-in-tree-report-overwrite` added special handling incompatible with the simplified policy. | Retire this RNC-only experimental flag while retaining scan-input, symlink, and special-file protection. | **Implemented and validated.** | Flag removed; CLI test asserts exit 126 for the old flag; CI run above. |
| V2-F | Publication correctness | Competing runs and OS-specific rename behavior can cause filename races. | Exclusive numbered destination allocation; atomic complete writes; real Windows/Linux/macOS validation; report actual output path. | **Implemented and validated.** | Exclusive-create reservation and temporary-file publication; Windows, macOS amd64/arm64, Linux, and Docker jobs passed in CI run above. |
| V2A | Configuration collision protection | `Config(cmd)` reset `cfg.Path` to the possibly-empty `--config` value after loading. Report validation therefore missed active `GITLEAKS_CONFIG` and extended configuration files; auto-discovered configs relied on the source-path fallback. | Track every file actually loaded, including the main file and all `extend.path` files, in `Config.ConfigFiles`. Pass the full list to report collision validation. Inline `GITLEAKS_CONFIG_TOML` data remains pathless. | **Implemented and validated.** | `config/config_test.go` and `cmd/report_config_test.go` cover explicit, auto-discovered, environment, one-level, and nested extensions, relative/absolute paths, pre-scan rejection, unchanged bytes, and override attempts; CI run below. |
| V3-A / CPV0-F01 | Source acquisition and worker errors | Directory walk, metadata, symlink, open, archive, and read failures could be silently dropped; semgroup errors were not joined. | Propagate path-aware failures, retain successfully read fragments, join workers, and aggregate errors with bounded `internal/scanerr` collector (default limit 10). Intentional exclusions remain non-errors. | **Implemented; cross-platform CI validated.** | `sources/files.go`, `sources/file.go`, `detect/detect.go`, `internal/scanerr`; Windows share-lock, partial-read, archive, and worker tests; [CI 43094102089](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43094102089). |
| V3-B / CPV0-F01 | Git execution integrity | Nonzero Git exits, fatal stderr, malformed/truncated diff streams, and incomplete process cleanup could yield false-clean results. | Wait/reap Git processes, classify fatal stderr without treating warnings as failure, drain pipes, and compare raw diff-section counts with parser output to detect tested incomplete streams. | **Implemented; cross-platform CI validated.** | `sources/git.go`, `sources/git_test.go`; PR research #2339/#2246/#2080/#2292; [CI 43094102089](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43094102089). |
| V3-C / CPV0-F19 | Worker lifecycle and cancellation | File and Git sources could return on cancellation before workers finished; deprecated `protect` lost caller context. | Wait for workers before returning; propagate command context through `DetectGitContext` and `protect`. Blocking OS reads remain a documented limitation. | **Worker lifecycle corrected; residual I/O limitation open.** | `detect/git.go`, `cmd/protect.go`, source/detector lifecycle tests; [CI 43094102089](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43094102089). |
| V3-D / CPV0-F01 | CLI scan completeness and report preservation | Failed/partial scans could exit 0, honor `--exit-code 0`, and replace a valid report with `[]` or partial findings. | Preserve default statuses on completed scans; operational/incomplete scans exit nonzero regardless of findings override; failed scans preserve existing reports and do not publish partial reports. | **Implemented; CLI regressions and cross-platform CI validated.** | `cmd/scan_status_test.go`, `cmd/report_test.go`, `README.md`; upstream compatible valid scans; [CI 43094102089](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43094102089). |

| V4-A / CPV0-F07 | Staged Git coverage | A text secret in a Git index blob marked binary by `.gitattributes -diff` could be silently omitted. | Read the staged index object with `git cat-file` when the diff omits text; continue MIME-based exclusion for genuinely binary content; preserve history/unstaged semantics. | **Implemented; CP-V4 CI 10/10 validated.** | `sources/git.go`, staged-index regressions; [#2338](https://github.com/gitleaks/gitleaks/pull/2338) (@VISHNU0906, **adapted design**, no verbatim code); [CI 43099917379](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43099917379). |
| V5-A / CPV0-F03 | Finding redaction | Empty captures corrupted strings or exposed matches; Unicode masking split UTF-8; adjacent and auxiliary required-rule secrets could leak in verbose/output fields. | Rune-safe partial masking; safe full-match fallback for empty captures; batch known secrets by file/commit/line; redact auxiliary fields; retain path-only findings. | **Implemented; CP-V5 CI 10/10 validated; unknown adjacent values remain residual.** | `report/finding.go`, `detect/utils.go`, CLI/report tests; [#2350](https://github.com/gitleaks/gitleaks/pull/2350) reviewed but unsafe early return rejected; [#2261](https://github.com/gitleaks/gitleaks/pull/2261) rune approach **adapted/credited in source**; [CI 43109905591](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43109905591). |
| V5-B / CPV0-F03 / F21 | Diagnostic credential leakage | Debug/trace, `GITLEAKS_CONFIG_TOML`, `AddFinding`, and Git stderr could disclose secrets independently of finding-report redaction. | Strip secret-bearing diagnostic fields, stop raw environment-config logging, sanitize bounded Git stderr/error strings, protect required-rule output under `--redact`; preserve default raw finding reports. | **Implemented; CP-V5 CI 10/10 validated; sanitizer is not comprehensive for unknown credentials.** | `detect/detect.go`, `detect/codec/decoder.go`, `cmd/root.go`, `sources/git.go`; CLI stdout/stderr and config tests; [CI 43109905591](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43109905591). |
| V5-C / CPV0-F07 | Staged blob lifecycle/performance | CP-V4's `blobReader.Close()` could drain the entire remainder of an intentionally excluded large blob. | Track EOF; `CloseAfterSkip` terminates and reaps the producer for intentional early exclusions; preserve error/cancellation/report behavior. | **Implemented; CP-V5 CI 10/10 validated; performance figure is agent-reported.** | `sources/git.go`; 100 MB staged-binary fixture reported ~530 ms CP-V4 vs ~305 ms CP-V5 (~286 ms pre-V4); blocked-read timeout test; [CI 43109905591](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43109905591). |
| V6-A / CPV0-F04 | Regex configuration | Malformed user-controlled rule/allowlist patterns could panic via `MustCompile`. | Error-returning `regexp.Compile` for user rule/allowlist patterns, legacy formats and extended config, on standard and RE2 backends; fail before scan/report publication. | **Implemented; CP-V6 CI 10/10 validated.** | `config/config.go`, `regexp/stdlib_regex.go`, `regexp/wasilibs_regex.go`; [#2082](https://github.com/gitleaks/gitleaks/pull/2082) (@Smeet23, adapted wrapper approach); [#2241](https://github.com/gitleaks/gitleaks/pull/2241) reviewed/not adopted; [CI 43114548814](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43114548814). |
| V6-B / CPV0-F04 / F22–F25 | Config extension state | Package-global `extendDepth` leaked across successive config translations. | Pass explicit depth through translation/extension; repeated config loads retain correct inheritance; preserve config-file tracking and valid policy semantics. | **Implemented; CP-V6 CI 10/10 validated; package-global Viper concurrency remains.** | `config/config.go`, sequential-extension tests; [CI 43114548814](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43114548814). |
| V6-C / CPV0-F04 | Source-local config trust | An automatically discovered repository `.gitleaks.toml` can replace default rules, relevant when scanning untrusted repos. | Trust boundary assessed; **do not force default rules** into intentionally custom configurations; use explicit trusted config in untrusted-scan orchestration. | **Assessed; policy decision OPEN, no new CLI flag.** | CP-V6 config precedence tests; possible opt-in policy deferred for operator review. |
| V6-D / CPV0-F05 | Ignore-file input safety | A FIFO/special file could hang `.gitleaksignore` processing, and incomplete reads could appear successful. | Reject non-regular inputs; Unix nonblocking open + post-open verification; Windows regular-file checks; `Scanner.Err()` and close/read failures propagated; merge entries only after full parse. | **Implemented; CP-V6 CI 10/10 validated; TOCTOU not fully eliminated.** | `detect/detect.go`, `detect/ignore_file_unix.go`, `detect/ignore_file_windows.go`, CLI/unit tests; issue [#1572](https://github.com/gitleaks/gitleaks/issues/1572); [CI 43114548814](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43114548814). |
| V6-E / CPV0-F05 | Ignore-file bounded parsing/discovery | Unbounded lines/files/entry counts and silent missing explicit paths could undermine reliable scans. | Enforce 16 MiB/file, 1 MiB/line, 100,000 non-comment entries; explicit missing paths error; optional absent files remain optional; regular-file symlinks retained; malformed fingerprint warning semantics unchanged. | **Implemented; CP-V6 CI 10/10 validated; extreme legitimate inputs may be rejected.** | `cmd/root.go`, `detect/detect.go`, ignore security tests; [CI 43114548814](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43114548814). |

| V7-A / CPV0-F09 | Archive identification | Filename-only archive classification silently missed ordinary plaintext with archive-like filenames. | Require content signature within 64 KB before extraction, restore consumed bytes on seekable/nonseekable readers; propagate inspection/read failures. | **Implemented; CI 10/10 validated; F09 accepted within tested scope.** | `sources/file.go`, source/CLI regressions; [#2317](https://github.com/gitleaks/gitleaks/pull/2317) (@Varshith-Kali; adapted), [#2316](https://github.com/gitleaks/gitleaks/issues/2316); [CI 43121770371](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43121770371). |
| V7-B / CPV0-F15 | Source/archive size limits | Size policy did not cover many compressed/reader/Git paths and had a rounded-limit legacy check. | Propagate `File.MaxFileSize` and `Size` in exact bytes; bound unknown/decompressed streams by temp spooling; zero limit remains unlimited; preserve Git-history textual diffs. | **Implemented; CI 10/10 validated; F15 partially OPEN.** | `sources/file.go`, `sources/files.go`, `sources/git.go`, CLI/source/reader tests; [#2099](https://github.com/gitleaks/gitleaks/pull/2099) (@alephnull-sh; partially adapted/informed); [CI 43121770371](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43121770371). |
| V7-C / CPV0-F09 / F15 | Coverage compatibility / resource budget | Nonseekable ZIP/7z temp seekable copy is capped by the member limit, potentially excluding a large compressed container with small valid members. Leading-only signatures can miss self-extracting/prefixed archives. | Reproduce Git/nonseekable oversized outer ZIP with ≤limit member; decide independent spool budget versus explicit incomplete-scan error. Keep unknown-extension archive discovery (#2010) separate. | **OPEN operator decision; do not close F15 fully.** | CP-V7 report §§6.2/12; `sources/file.go` temp-copy branch; [#2010](https://github.com/gitleaks/gitleaks/issues/2010). |

**Output decision table (when the requested destination exists):**

| `auto-overwrite` | `aggregation` | Default outcome | One-run override |
| --- | --- | --- | --- |
| true | true | Aggregate (numbered file) | `--overwrite` |
| true | false | Overwrite requested file | `--aggregate` |
| false | true | Aggregate (numbered file) | `--overwrite` |
| false | false | Reject before scan | Either explicit flag |

Neither `--overwrite` nor `--no-overwrite` was found in the reviewed upstream v8 CLI. **Do not call either an original upstream flag.** Explicit overwrite must never authorize destroying a scan input. No report destination means no file-output policy change; stdout must remain compatible.

## 3. Known issue disposition

This section distinguishes an implemented correction from full release acceptance. Status reflects delivered code and verification through CP-V7 implementation/CI (2026-11-12); it does not imply universal coverage or release acceptance. F09 accepted within tested scope; F15 conditional.

| CP-V0 finding | Severity | Description | Disposition |
| --- | --- | --- | --- |
| F01 | P0 | Source/Git/worker errors can be swallowed, creating false clean scans or reports. | **CP-V3 corrected and CI-validated for documented CLI-used scan paths:** source/worker/Git/parser errors return nonzero, `--exit-code 0` cannot suppress operational failure, and prior reports survive failed scans. Deprecated `DetectFiles`/`DirectoryTargets`, non-cancelable reads, and pathological parser cases remain limitations. |
| F02 | P0 | Report validation/truncation, symlink, partial output, and input-collision data loss. | **V1, V2, and V2A implemented; CP-V2/V2A cross-platform CI passed.** Residual final publication TOCTOU is documented in the CP-V2 report. |
| F03 | P1 | Secret-bearing logs and incomplete redaction across multiple/empty/Unicode secrets. | **CP-V5 implemented and CI-validated:** sanitized diagnostic paths, Unicode/empty/multiple-secret redaction, adjacent known findings, and required-rule output; #2261 adapted, #2350 not adopted. **Partially open:** unknown adjacent credentials, fragment-boundary grouping, partial mask disclosures, and caller-provided fragment/template surfaces. |
| F04 | P1 | Untrusted local configuration can replace rules; malformed regex can panic. | **CP-V6 implemented and CI-validated for malformed user regex and extension-depth leakage.** Standard/RE2 rule and allowlist compile errors are controlled, and repeated translation depth is isolated. **OPEN policy:** untrusted source-local `.gitleaks.toml` can intentionally replace default rules; concurrent global Viper state and pattern-bearing diagnostics remain residual. |
| F05 | P1 | Unsafe `.gitleaksignore` file types, blocking FIFO, scan errors, lack of limits. | **CP-V6 implemented and CI-validated:** reject FIFO/socket/device inputs, check scanner/read/close errors, merge only complete parses, and enforce 16 MiB/1 MiB/100,000-entry bounds; regular symlinks preserved. **Residual:** Windows and cross-platform file mutation/TOCTOU, compatibility with exceptionally large legitimate ignore files. |
| F06 | P0 | Global AND allowlist path short-circuit suppresses findings. | **V1 implemented;** release regression proof still required. |
| F07 | P0 | Staged text classified as binary may be skipped. | **CP-V4 implemented and CI-validated:** staged binary-marked *text* acquired from the index blob and scanned; genuine binary data still excluded. CP-V5 corrected the expensive blob-drain exit path. **Release regression proof and limits-policy review still required.** |
| F08 | P1 | Forced Git color hides findings. | **V1 implemented;** release regression proof still required. |
| F09 | P1 | Archive extension false positives can skip ordinary text. | **CP-V7 implemented and 10/10 CI validated:** ordinary text with misleading archive-like names is scanned; archive-verification read errors are operational. **Residual:** prefixed/self-extracting archive signature limitations, unknown-extension archive discovery (#2010) separate/open. |
| F10 | P1 | MIME reclassification on later chunks misses secrets. | **V1 implemented;** release regression proof still required. |
| F11 | P1 | Merge-only Git changes can be missed in default history mode. | **OPEN / compatibility decision required**; #2301 changes behavior. |
| F12 | P1 | `dir` silently scans the wrong directory on invalid positional args. | **V1 implemented;** release regression proof still required. |
| F13 | P1 | Pre-commit helper lets nonzero scanner statuses pass. | **V1 implemented;** release regression proof still required. |
| F14–F18 | P2 | Archive identification, size budgets, path identity, Unicode redaction, stdout contracts. | **Partially addressed:** F17 Unicode masking corrected in CP-V5; F15 byte-accurate limits improved in CP-V7 but the nonseekable ZIP/7z outer-container skip remains unresolved. F14/F16/F18 and residual path/stdout items remain open/deferred. |
| F19 | P2 | Returning before worker completion on cancellation. | **CP-V3 worker-join issue corrected and tested;** blocking OS reads can still outlast configured timeout. |
| F20 | P2 | Report permissions and close/flush weaknesses. | **CP-V2 implemented major safeguards;** final-publication TOCTOU remains documented. |
| F21 | P2 | Git stderr can contain and expose credentials through logs/error text. | **CP-V5 implemented and CI-validated for known Git stderr credential patterns:** sanitized URL userinfo, sensitive query values, and token-shaped runs in logs/error strings. **PARTIALLY OPEN:** unknown/short formats may escape heuristic sanitization and benign SHA/path values can be over-redacted; consider structured diagnostic boundary. |
| F22–F25 | P2/P3 | Output contracts, extension state, performance, build/CI hygiene. | **Partially addressed:** CP-V6 removed global extension-depth leakage; CP-V5 improved intentional blob-skip performance; CP-V3 fixed test-only `lostcancel`. Residual global Viper concurrency, deprecated APIs, build/generation graph, and other CP-V0 items require separate acceptance. |

At CP-V1 completion, `go build`, normal tests, generation, targeted Linux race suites, RE2 build, and Docker smoke tests were reported passing; `go vet` had pre-existing test-only `lostcancel` warnings and some comprehensive tagged/race lanes had known limitations. Subsequent V1G branch CI had Windows/macOS report-path test failures. CP-V2 removed that path policy and the failures no longer reproduce.

CP-V2 code at `452d1224afe57b68888b4e581ad70878758426ad` passed Windows, Linux, macOS amd64, macOS arm64, standard-build, RE2-build, module/vulnerability, and Docker amd64/arm64 CI jobs in [run 43075289894](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43075289894). Local validation also passed `go build ./...`, `go test -count=1 ./...` with Go 1.27.4, `go generate ./...` with unchanged generated configuration, Linux race tests for `./cmd`, `./config`, and `./report`, plus Docker CLI/report smoke tests. `go vet ./...` still reports only the two pre-existing `lostcancel` diagnostics in `detect/detect_test.go`.

CP-V2A code at `cf3007301441eaf1f7b449707c87c2e87cbfebc8` preserves every configuration file loaded from disk in `Config.ConfigFiles` and passes the complete list into report collision validation. Regression coverage proves explicit, auto-discovered, `GITLEAKS_CONFIG`, one-level, and nested extended configuration collisions are rejected before scanning, remain byte-for-byte unchanged, and cannot be bypassed with `--overwrite` or `--aggregate`. [Run 43084497923](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43084497923) passed all Windows, Linux, macOS amd64/arm64, module/vulnerability, RE2, standard-build, and Docker amd64/arm64 jobs.

CP-V3 delivered four commits: `105c0a8` (source errors), `56458e7` (Git process/diff integrity), `fd996ec` (detector error propagation and worker completion), and `b556701` (CLI regressions). [Final CI run 43094102089](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43094102089) passed all 10 jobs at `b556701`; independent review confirmed the four-commit diff, 20 changed files, and matching branch HEAD. Agent-reported local validation included Go 1.27.4 Windows build/tests/race, Linux race tests, generator hash unchanged (`256f6479…`), RE2 application build, Docker smoke, and clean `go vet`. The historical `lostcancel` test diagnostics were corrected incidentally in the updated archive-test path using `defer cancel()`.


### CP-V4–V6 reconciliation evidence (2026-11-11)

- **CP-V4 / F07:** `2664154` corrects staged `-diff` text-blob omission using index acquisition and MIME verification. [Run 43099917379](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43099917379) independently verified **10 successful jobs**. Flash's local validation and the detailed CP-V4 report are operator-reported; no direct local-report file was inspected.
- **CP-V5 / F03, F17, F21:** commits `920b3e2`, `51d3e2a`, `d3789d0`, `8779007`. Reports default to unredacted findings when `--redact=0`; no metadata-only mode was introduced. Operational logs omit diagnosed secret fields, `GITLEAKS_CONFIG_TOML` text is not logged, and Git diagnostics are sanitized heuristically. Redaction tests cover multiple/required/empty/Unicode findings and output formats. Flash reported identical seven-finding default fixture reports against upstream/pre-V5; 100 MB staged-binary averages ~286 ms pre-V4, ~530 ms V4, ~305 ms V5. These fixture/performance claims are **agent-reported**, while [run 43109905591](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43109905591) is independently verified **10/10 successful**.
- **CP-V6 / F04, F05, extension state:** commits `ed74f7e`, `d0795f5`. Invalid user regex now returns an error rather than a panic; extension depth is explicit; unsafe ignore inputs are rejected with bounded, atomic-in-memory parsing. Flash reports valid default/custom/extended/ignore scans maintain expected findings; config trust boundary deliberately unchanged; generation hash `256f6479…` unchanged. [Run 43114548814](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43114548814) independently verified **10/10 successful**.
- **Full CP-V4, CP-V5, CP-V6 and CP-V7 reports are now uploaded and reviewed.** Agent-reported local results are distinguished from independently verified GitHub source/CI evidence. No claim of universal secret-safe logs, all scan paths being cancelable, or release certification is made.

**Residual technical gates after CP-V7:** F02 final-publication TOCTOU; F03 unknown adjacent secret/template exposure; F04 untrusted local config policy, Viper concurrency, and sensitive regex diagnostics; F05 file-replacement race; F11 merge-history default behavior; F15 nonseekable ZIP/7z container budget (CP-V7 conditional); F16 path identity; F19 blocking reads/deprecated APIs; F21 imperfect Git stderr sanitization. Existing experimental/default report compatibility and contributor-provenance reviews remain release gates.


### CP-V7 acceptance evidence and residual decisions (2026-11-12)

**F09 (accepted within tested scope).** Archive-like plaintext now gets ordinary scanning because filename-derived format selection requires a content signature from a bounded 64,000-byte prefix. Read errors and corrupt genuine archive errors propagate; prefix replay preserves seekable/nonseekable input. Agent-reported fixtures: `notes.7z.txt` 0→1 finding and 100 misleading archive-name files 0→100 findings. #2010 (unknown-extension archives) remains out of scope.

**F15 (partially open).** Exact byte-based `File.MaxFileSize`/`Size` semantics now propagate to filesystem, stdin, direct readers, Git staged/index/archive content, archive members, and decompression. Zero means no explicit limit, size equal to limit is scanned, and Git-history textual-diff scanning remains unchanged. A nonseekable ZIP/7z outer file over the per-file limit is currently skipped even when all members are smaller. Require a focused staged/Git/nonseekable regression and operator choice on a separate temporary-copy budget or explicit incomplete-scan error rather than a misleading clean result. Known-size metadata races and temporary-spool disk usage remain limitations. Self-extracting/prefixed archives may be downgraded to text because only the initial signature is checked.

**Performance (agent-reported, not controlled release benchmarks):** 100 small text files median 532→653 ms (+121 ms, ~23% relative), peak RSS 48.2→47.5 MB. Twenty-member ZIP median 413→485 ms; 50 staged binary-marked files 1636→1775 ms. Re-run under controlled conditions before publishing performance claims.

**Validation:** `cccdef6`, 13 files +862/−41; [CI 43121770371](https://github.com/RNC-Solutions-Group/gitleaks/actions/runs/43121770371), **10/10 successful jobs independently verified**. Agent's full uploaded CP-V7 report and its local tests/benchmark details were reviewed; local execution was not independently rerun. RE2 app and changed-package tests pass; full tagged generator packages retain the pre-existing type mismatch.

**Provenance:** [PR #2317](https://github.com/gitleaks/gitleaks/pull/2317) @Varshith-Kali SHA `71e91200afb2415b283bb1721709b19c2fddd8aa` **adapted** 64 KB prefix verification/replay; [PR #2099](https://github.com/gitleaks/gitleaks/pull/2099) @alephnull-sh SHA `1b3896fd7df9a667a09568b751dce6143613319e` **partially adapted/informed** size propagation, not its memory-heavy buffer; [issue #2316](https://github.com/gitleaks/gitleaks/issues/2316) @bhandrigan and [issue #2010](https://github.com/gitleaks/gitleaks/issues/2010) @grodriguezl credited. No external commit cherry-picked; release code-overlap audit pending. Full CP-V4 report additionally establishes that #2338 was **adapted**, not only reviewed/informed.

**Remaining scan-integrity limitations:** deprecated APIs still have error-dropping paths outside the CLI; blocking OS reads cannot always be promptly canceled; raw `diff --git` header counting is a conservative parser-integrity heuristic and can false-fail unusual custom Git formatting or miss a pathological complete-prefix truncation. Do not claim universal completeness or cancellation. CP-V2 publication TOCTOU and scan-included numbered reports also remain documented.

## 4. External contributor recognition and code provenance

### Attribution states (must be used accurately)

- **Reviewed:** Studied an issue, PR, test, or discussion; no claim that code was copied.
- **Informed / independently implemented:** Contributor exposed a failure mode or approach, but RNC wrote its own implementation. Credit the source in documentation and the final upstream PR.
- **Adapted implementation:** RNC incorporated identifiable logic or code from another person's work and modified it. Record exactly what was adapted, the original PR/commit, and appropriate author credit.
- **Adopted commit:** A contributor's commit is deliberately incorporated as a **separate authored commit** (prefer `git cherry-pick -x`), not hidden inside the RNC squash. Preserve the original author's metadata and record the resulting fork SHA.
- **Not selected:** Reviewed but deliberately not integrated. Still credit relevant research when it shaped a decision; do not add a false co-author trailer.

**Current adoption ledger: No external PR is yet recorded here as a verified separately adopted commit.** Similarity to an upstream PR does not establish copying or permission to claim its code has been integrated. Determine provenance from actual diffs and development records before changing a status.

### 4A. Upstream contributions directly relevant to implemented V1 work

All entries below are **reviewed/informed**, not confirmed cherry-picks. Head SHAs were rechecked on 2026-11-11; recheck before actual integration.

| RNC area | Upstream work | Contributor | PR head SHA | What it contributed to our assessment | Code-adoption status |
| --- | --- | --- | --- | --- | --- |
| F02 | [#2307](https://github.com/gitleaks/gitleaks/pull/2307), issue [#2304](https://github.com/gitleaks/gitleaks/issues/2304) | [@jdotc1](https://github.com/jdotc1) | `0d1ae3ea1087bebcf5ce45982a770d0f4925af6c` | Proved preflight deletion; proposed non-destructive validation and regression. Atomic final publication is additional RNC scope. | **Reviewed; no adopted commit verified.** |
| F02 / permissions | [#2300](https://github.com/gitleaks/gitleaks/pull/2300) | [@ToniAdreal](https://github.com/ToniAdreal) | `2b08eb90383f6f221496d984bc3d9a575d002f8b` | Restrictive `0600` report-permission proposal and discussion of preflight deletion. | **Reviewed; no adopted commit verified.** |
| F02 / permissions | [#2312](https://github.com/gitleaks/gitleaks/pull/2312) | [@Varshith-Kali](https://github.com/Varshith-Kali) | `2009273e0029f91aeda4732552ae5489634d0749` | Alternative `0600` permissions proposal; still truncates existing destinations. | **Reviewed; no adopted commit verified.** |
| F06 | [#2227](https://github.com/gitleaks/gitleaks/pull/2227) | [@BasitS-hash](https://github.com/BasitS-hash) | `49b26eece1ea83ea58da0c6e0bf92326fb82a25c` | AND-allowlist path sufficiency and regression cases. | **Reviewed; inspect for code overlap before release.** |
| F08 | [#2284](https://github.com/gitleaks/gitleaks/pull/2284) | [@lekanjava](https://github.com/lekanjava) | `1894e6f2cdaef94fe66388e577e0f59f14bfcff2` | Explicit `--no-color` in Git acquisition to protect parsers. | **Reviewed; inspect for code overlap before release.** |
| F10 | [#2329](https://github.com/gitleaks/gitleaks/pull/2329) | [@cardene777](https://github.com/cardene777) | `21028d8858e31fd088029f2fdb6634c99a1d3563` | First-chunk-only MIME classification. | **Reviewed; inspect for code overlap before release.** |
| F12 | [#2334](https://github.com/gitleaks/gitleaks/pull/2334), issue [#1979](https://github.com/gitleaks/gitleaks/issues/1979) | [@Shanmukh-l33t](https://github.com/Shanmukh-l33t) | `ad321b7d7aa9b1a1642bb56813fa7443d5412971` | `dir` positional argument validation and wrong-directory regression. | **Reviewed; inspect for code overlap before release.** |

### 4B. Previously inspected PRs and future integration candidates

This is a historical research/candidate index, not an adoption record. CP-V3 decisions for F01 are summarized in §4D; entries below retain CP-V0-era SHAs where available and must be revalidated before actual code integration.

| Scope | PR | Contributor | Historical PR head / note | Assessment / integration disposition |
| --- | --- | --- | --- | --- |
| F01 Git exit | [#2339](https://github.com/gitleaks/gitleaks/pull/2339) | [@delaysystem](https://github.com/delaysystem) | `7e5132faa1ccea4bfd0b4b006ca896b5dd03df3f` | Partial exit-status fix; add stderr/parse/worker handling. **Rechecked 11 Nov.** |
| F01 diff parsing | [#2246](https://github.com/gitleaks/gitleaks/pull/2246) | [@aron-intframe](https://github.com/aron-intframe) | `9b7b2243fcdeb4279cccd9df9348f9e537d4fb95` | Partial parser-drain approach; source callback still swallows error. **Rechecked 11 Nov.** |
| F01 Git stderr | [#2080](https://github.com/gitleaks/gitleaks/pull/2080) | [@Smeet23](https://github.com/Smeet23) | `ff77f9f58b1ec5fbbece852e0bfe9c1f050f24e5` | Partial stderr propagation. |
| F01 Git stderr | [#2292](https://github.com/gitleaks/gitleaks/pull/2292) | [@Saadanjum0](https://github.com/Saadanjum0) | **SHA to recheck** | Similar partial stderr propagation. |
| F01 error propagation | [#2252](https://github.com/gitleaks/gitleaks/pull/2252) | [@CAOShurong](https://github.com/CAOShurong) | `424b4a92ffbcf4aa4b16ef98b53ecce45e1e3e33` | Fragment error flow; not complete source coverage. **Rechecked 11 Nov.** |
| F01 fragment errors | [#2253](https://github.com/gitleaks/gitleaks/pull/2253) | [@arpitjain099](https://github.com/arpitjain099) | `4f58f98a1424203af39cc00a26f9672e9b45b794` | Bounded/aggregate error ideas; verify actual resource limits. |
| F01 file open errors | [#2249](https://github.com/gitleaks/gitleaks/pull/2249) | [@vaibhav8a](https://github.com/vaibhav8a) | `509fce15ad2eb2536bc8d840a3ed3e374d73bb2d` | Path-aware open errors, incomplete walk/read handling. |
| F01 file open errors | [#2313](https://github.com/gitleaks/gitleaks/pull/2313) | [@Varshith-Kali](https://github.com/Varshith-Kali) | `acfeacda3adc17b8659d46b82180edd9e5adc7c7` | Alternative open-error reporting. |
| F01 file open errors | [#2235](https://github.com/gitleaks/gitleaks/pull/2235) | [@hamodywe](https://github.com/hamodywe) | `d86cc860f3c5f7085d8d68c24bc09424a73aabf7` | Count-only alternative, less useful than path-aware reports. |
| F07 staged binary | [#2338](https://github.com/gitleaks/gitleaks/pull/2338) | [@VISHNU0906](https://github.com/VISHNU0906) | `7fed4140d87f62090c87b9f09c5c9b6073a84a27` | **CP-V4 adapted staged-flag/binary-fallback design, not cherry-picked.** RNC corrected staged-index blob retrieval with broader MIME/error/lifecycle handling; final exact patch-overlap review still required. |
| F09 archive identity | [#2317](https://github.com/gitleaks/gitleaks/pull/2317) | [@Varshith-Kali](https://github.com/Varshith-Kali) | `71e91200afb2415b283bb1721709b19c2fddd8aa` | **CP-V7 adapted** bounded content-signature verification and prefix replay with expanded error tests; not cherry-picked. |
| F11 merge-only secrets | [#2301](https://github.com/gitleaks/gitleaks/pull/2301) | [@bhaskargurram-ai](https://github.com/bhaskargurram-ai) | `6bcf2f91c1b8765bbf5387507b4f0e1258b01dcc` | Behavior-changing; requires operator decision and compatibility tests. |
| F03 empty secret redaction | [#2350](https://github.com/gitleaks/gitleaks/pull/2350) | [@kevin9327](https://github.com/kevin9327) | `f19096fae4ba7a3016476c7b9d5528a6f31bbd9c` | **CP-V5 reviewed/informed; early-return implementation rejected** because empty selected captures can leave `Match` exposed. |
| F03/F17 UTF-8 masking | [#2261](https://github.com/gitleaks/gitleaks/pull/2261) | [@arpitjain099](https://github.com/arpitjain099) | `b1ce07f27d8251ef22e5b47c80387e9b1590fc2f` | **CP-V5 adapted** rune-based mask with author credit in `report/finding.go`; no cherry-picked authored commit. |
| F04 invalid regex | [#2082](https://github.com/gitleaks/gitleaks/pull/2082) | [@Smeet23](https://github.com/Smeet23) | `df8223e0008b6d5bf82d0433592b8f1557c97311` | **CP-V6 adapted** error-returning compile-wrapper direction, expanded contextual errors/tests; no cherry-pick. |
| F04 invalid regex | [#2241](https://github.com/gitleaks/gitleaks/pull/2241) | [@harriiinnii](https://github.com/harriiinnii) | `0c2402a9f71703c660c9ec197e75c2a452f02ac7` | **CP-V6 reviewed/informed; not adopted** due to unrelated expected-regex edits and incomplete allowlist error context. |
| F15 size limits | [#2099](https://github.com/gitleaks/gitleaks/pull/2099) | [@alephnull-sh](https://github.com/alephnull-sh) | `1b3896fd7df9a667a09568b751dce6143613319e` | **CP-V7 partially adapted/informed:** size propagation/member precheck; used bounded temporary spooling instead of whole-content RAM; ZIP/7z exception remains. |
| F16 path identity | [#2248](https://github.com/gitleaks/gitleaks/pull/2248) | [@hamodywe](https://github.com/hamodywe) | `5c9ca86f8f4d1db1e17004b6a1b598137203face` | Broad root-relative path proposal; sequence after AND-allowlist safety. |

### 4C. Issue reporters and research contributors

Issue authors deserve recognition even when their proposed implementation is not chosen. This is an initial verified/CP-V0-informed list; extend it as more research is performed:

| Issue | Reporter | Contribution |
| --- | --- | --- |
| [#2299](https://github.com/gitleaks/gitleaks/issues/2299) | [@carfeii](https://github.com/carfeii) | Report files could expose discovered plaintext secrets through ambient umask. |
| [#2304](https://github.com/gitleaks/gitleaks/issues/2304) | [@jdotc1](https://github.com/jdotc1) | Destructive report-validation preflight and reproduction. |
| [#338](https://github.com/gitleaks/gitleaks/issues/338) | [@lochlainnw](https://github.com/lochlainnw) | Overwrites when scanning many repositories; prompted investigation of aggregation; upstream maintainer preferred external scripts over merging reports. |
| [#1065](https://github.com/gitleaks/gitleaks/issues/1065) | [@mesaglio](https://github.com/mesaglio) | Fail-fast validation of invalid report destinations before expensive scans. |
| [#1572](https://github.com/gitleaks/gitleaks/issues/1572) | [@joerg1985](https://github.com/joerg1985) | Ignore-file safety and blocking concerns; CP-V6 implemented bounded special-file/read handling. |
| [#2232](https://github.com/gitleaks/gitleaks/issues/2232) | [@siliroid](https://github.com/siliroid) | Source-file error visibility and false-clean behavior. |
| [#2337](https://github.com/gitleaks/gitleaks/issues/2337) | [@farrrr](https://github.com/farrrr) | Staged binary-classified text can be missed. |
| [#2316](https://github.com/gitleaks/gitleaks/issues/2316) | [@bhandrigan](https://github.com/bhandrigan) | Archive filename classification false positives. |
| [#1979](https://github.com/gitleaks/gitleaks/issues/1979) | [@caiodelgadonew](https://github.com/caiodelgadonew) | Wrong-directory `dir` behavior. |
| [#1651](https://github.com/gitleaks/gitleaks/issues/1651), [#1757](https://github.com/gitleaks/gitleaks/issues/1757) | [@rgmz](https://github.com/rgmz) | Decode/chunk boundary and input-coverage concerns. |
| [#1243](https://github.com/gitleaks/gitleaks/issues/1243) | [@SyedMonzAli](https://github.com/SyedMonzAli) | Secret exposure in verbose redaction; partially addressed by CP-V5. |
| [#1846](https://github.com/gitleaks/gitleaks/issues/1846) | [@AnthonyPoschen](https://github.com/AnthonyPoschen) | Malformed user config regex causes a panic; corrected by CP-V6. |
| [#2296](https://github.com/gitleaks/gitleaks/issues/2296) | [@sporteka2](https://github.com/sporteka2) | Request for redaction/report configuration; reviewed and deferred, no new option in CP-V5. |

### 4D. CP-V3 scan-completeness contribution review (2026-11-11)

Flash reports reviewing each priority PR's implementation and tests at its then-current head, plus pre-existing merged #1711. **No external code was adopted verbatim or recorded as a separately authored commit.** The assessments below are research/influence classifications, not proof of direct code incorporation. The comprehensive CP-V3 report is maintained separately from this ledger; source SHAs in §4B are historical snapshots and should be checked again before clean-history reconstruction.

| PR | Contributor | CP-V3 decision | Relation to RNC correction / remaining limitation |
| --- | --- | --- | --- |
| [#2339](https://github.com/gitleaks/gitleaks/pull/2339) | @delaysystem | **Reviewed/informed; independent implementation** | Git nonzero exit/Wait handling; RNC also handles stderr, parser integrity, and workers. |
| [#2246](https://github.com/gitleaks/gitleaks/pull/2246) | @aron-intframe | **Informed; independent implementation** | Parser-drain/unparsed-output approach; RNC adds raw diff-section count verification. |
| [#2080](https://github.com/gitleaks/gitleaks/pull/2080) | @Smeet23 | **Reviewed/informed** | Stderr error return; RNC distinguishes fatal diagnostics from warnings and checks process exit. |
| [#2292](https://github.com/gitleaks/gitleaks/pull/2292) | @Saadanjum0 | **Reviewed/informed** | Alternative stderr reporting; no direct code adoption. |
| [#2252](https://github.com/gitleaks/gitleaks/pull/2252) | @CAOShurong | **Reviewed/not selected as-is** | Immediate fragment-error return would abandon continued collection of readable findings. |
| [#2253](https://github.com/gitleaks/gitleaks/pull/2253) | @arpitjain099 | **Informed; independent implementation** | Error accumulation informed bounded `internal/scanerr` collector; no direct patch adoption established. |
| [#2249](https://github.com/gitleaks/gitleaks/pull/2249) | @vaibhav8a | **Informed; independent implementation** | Path-aware open errors; RNC covers broader walker/read/worker failures. |
| [#2313](https://github.com/gitleaks/gitleaks/pull/2313) | @Varshith-Kali | **Reviewed/not selected as-is** | Alternative open-error reporting; RNC chose a broader integrated pipeline. |
| [#2235](https://github.com/gitleaks/gitleaks/pull/2235) | @hamodywe | **Rejected as implementation** | Count-only error state loses useful file paths and lacks broader failure coverage. |
| [#1711](https://github.com/gitleaks/gitleaks/pull/1711) | @dgn | **Merged upstream/pre-existing baseline** | Historical Git warning behavior contribution, reported merged as `3188ad6`; not a new RNC adoption. |

Related research/issues #1464, #1981, #1338, and [#2232](https://github.com/gitleaks/gitleaks/issues/2232) were reviewed in CP-V3. Credit these sources when relevant in any future upstream submission; do not invent code authorship from research influence.

**Current F01 provenance decision:** `internal/scanerr`, Git lifecycle/parser integrity, source error handling, and CLI regressions are recorded as RNC-implemented, informed by the above research. Reassess exact code overlap during the clean-release audit before asserting definitive independent authorship.

**License:** Preserve upstream Gitleaks MIT `LICENSE` and notices. License permission to reuse code is not a substitute for accurate author credit. Do not invent email addresses or apply `Co-authored-by` merely because someone reported a problem.


### 4E. CP-V4–V6 provenance and contributor disposition (2026-11-11)

These are checkpoint decisions, **not** automatic permission to adopt any upstream commit. Final release construction still requires original-Git-object author metadata and an exact patch-overlap review.

| Checkpoint | Upstream work / contributor | Source SHA or issue | Disposition and RNC relationship | Attribution / unresolved decision |
| --- | --- | --- | --- | --- |
| CP-V4 / F07 | [#2338](https://github.com/gitleaks/gitleaks/pull/2338) / @VISHNU0906 | `7fed4140d87f62090c87b9f09c5c9b6073a84a27` (historical) | **Adapted** staged-flag and binary-fallback design; RNC implementation reads index blobs and handles MIME/cleanup. | No upstream commit cherry-picked; final code-overlap audit. |
| CP-V5 / F03 | [#2350](https://github.com/gitleaks/gitleaks/pull/2350) / @kevin9327 | `f19096fae4ba7a3016476c7b9d5528a6f31bbd9c` | **Reviewed/informed; not adopted.** Early return could leave raw `Match` visible. | Credit defect research, not authoring of RNC fallback. |
| CP-V5 / F03/F17 | [#2261](https://github.com/gitleaks/gitleaks/pull/2261) / @arpitjain099 | `b1ce07f27d8251ef22e5b47c80387e9b1590fc2f` | **Adapted rune-mask implementation.** Source includes explicit attribution; no external commit cherry-picked. | Ensure adapted logic and original author clearly credited in eventual clean release. |
| CP-V5 / F03 | [#1243](https://github.com/gitleaks/gitleaks/issues/1243) / @SyedMonzAli | Issue, still open on review | **Partially addressed:** known adjacent and trailing-line redaction; unknown/non-finding text remains. | Credit issue reporter; no copied patch claimed. |
| CP-V5 / UX | [#2296](https://github.com/gitleaks/gitleaks/issues/2296) / @sporteka2 | Issue, still open on review | **Deferred:** no new report/redaction configuration surface. | No code adoption. |
| CP-V6 / F04 | [#2082](https://github.com/gitleaks/gitleaks/pull/2082) / @Smeet23 | `df8223e0008b6d5bf82d0433592b8f1557c97311` | **Adapted** error-returning regex-wrapper approach; extended validation/diagnostics across rule/allowlist formats. | Credit adapted implementation; exact patch-overlap review before release. |
| CP-V6 / F04 | [#2241](https://github.com/gitleaks/gitleaks/pull/2241) / @harriiinnii | `0c2402a9f71703c660c9ec197e75c2a452f02ac7` | **Reviewed/informed; not adopted.** Contains unrelated test expectation changes. | Credit scenario research only. |
| CP-V6 / F04 | [#1846](https://github.com/gitleaks/gitleaks/issues/1846) / @AnthonyPoschen | Issue, closed on review | **Issue reference:** invalid regex panic reproduced/corrected. | Reporter credited; no implementation adoption. |
| CP-V6 / F04 | [#2228](https://github.com/gitleaks/gitleaks/pull/2228) / @ankit090701 | `637df0f…` (historical short SHA; verify full Git object) | **Reviewed duplicate / not adopted.** | No authored commit incorporated. |
| CP-V6 / F05 | [#1572](https://github.com/gitleaks/gitleaks/issues/1572) / @joerg1985 | Issue, open on review | **Informed:** trust boundary/special-file handling; RNC implemented its own guarded reader. | Reporter credited; no copied patch. |

**Checkpoint provenance summary:** CP-V4, CP-V5, and CP-V6 were committed as RNC branch changes. No external PR was cherry-picked as a separately authored commit in those checkpoints. **Adapted** (#2261 and #2082), **informed** (#2338 and others), **reviewed-not-adopted**, and **issue-reported** are deliberately distinct. Do not invent author emails or `Co-authored-by` trailers. Do not collapse adapted contributor code into a future RNC-only squash without clear credit and a deliberate operator integration choice.

### 4F. CP-V7 contributor provenance

| Scope | Source | Contributor | SHA | Relation | Clean-release requirement |
| --- | --- | --- | --- | --- | --- |
| F09 signature verification | [#2317](https://github.com/gitleaks/gitleaks/pull/2317) | @Varshith-Kali | `71e91200afb2415b283bb1721709b19c2fddd8aa` | **Adapted:** bounded 64 KB signature and prefix replay; RNC added explicit error path and broader tests. | Credit code/design; verify overlap before clean history. |
| F15 source/member size | [#2099](https://github.com/gitleaks/gitleaks/pull/2099) | @alephnull-sh | `1b3896fd7df9a667a09568b751dce6143613319e` | **Partially adapted/informed:** size field propagation and prechecks; RAM-buffer approach rejected. | Credit adopted/influenced elements; review exact overlap. |
| F09 false-clean report | [#2316](https://github.com/gitleaks/gitleaks/issues/2316) | @bhandrigan | Issue | Reported misleading archive filenames and false clean. | Credit reporter; no code adopted. |
| Unknown-extension archive discovery | [#2010](https://github.com/gitleaks/gitleaks/issues/2010) | @grodriguezl | Issue | Separate, deferred problem; not corrected by CP-V7. | Do not claim closure. |

**CP-V4 provenance correction:** The uploaded V4 report classifies [#2338](https://github.com/gitleaks/gitleaks/pull/2338) (@VISHNU0906) as **adapted**, specifically the staged-flag/binary-fallback design. RNC independently implemented stage-zero index resolution, errors, lifecycle, and regression tests. No external PR was cherry-picked in CP-V4–V7.

## 5. Adopt-vs-adapt decision log (required before clean release)

For every outside PR potentially overlapping RNC work, add a decision row **after examining the exact patch and our corresponding code**:

| PR / commit | RNC files or behavior overlapped | Decision: reviewed / independent / adapted / adopt | Original author and email from Git object | Original SHA -> resulting fork SHA | Tests / compatibility evidence | Operator approval |
| --- | --- | --- | --- | --- | --- | --- |
| #2227 / `49b26ee` | F06 AND allowlist | **Undecided; reviewed, overlapping design** | Resolve before integration | — | V1 tests; revalidate | Pending |
| #2284 / `1894e6f` | F08 Git color | **Undecided; reviewed, overlapping design** | Resolve before integration | — | V1 tests; revalidate | Pending |
| #2329 / `21028d8` | F10 first-chunk MIME | **Undecided; reviewed, overlapping design** | Resolve before integration | — | V1 tests; revalidate | Pending |
| #2334 / `ad321b7` | F12 CLI arguments | **Undecided; reviewed, overlapping design** | Resolve before integration | — | V1 tests; revalidate | Pending |
| #2307 / `0d1ae3e` | F02 report validation | **Undecided; reviewed, partial design overlap** | Resolve before integration | — | V1 tests; V2 pending | Pending |
| #2300 / `2b08eb9`, #2312 / `2009273` | F02 permissions | **Alternative proposals; no adopted commit** | Resolve if borrowing specific patch | — | Permission tests; revalidate | Pending |
| #2339, #2246, #2080, #2292 | F01 Git lifecycle/parser/stderr | **CP-V3 reviewed/informed; no external patch adopted** | Refer to §4D; preserve original author records during final reconciliation | — | Git process/parser regression tests; CI 43094102089 | Re-audit code overlap before clean release |
| #2252, #2253 | F01 fragment/error propagation | **CP-V3 reviewed/informed; RNC implemented bounded collector** | Refer to §4D; no separate external commit recorded | — | `internal/scanerr`, detector tests; CI 43094102089 | Re-audit code overlap before clean release |
| #2249, #2313, #2235 | F01 file source errors | **CP-V3 reviewed/informed or rejected as-is** | Refer to §4D; no separate external commit recorded | — | Source/lock/read tests; CI 43094102089 | Re-audit code overlap before clean release |

| #2338 / `7fed414` | F07 staged-index acquisition and MIME behavior | **CP-V4 adapted staged-fallback design; no cherry-pick** | @VISHNU0906; original Git metadata not yet recorded | — | Staged-index tests; CI 43099917379 | Final code overlap/provenance review pending |
| #2350 / `f19096f` | F03 empty-capture redaction | **Reviewed/informed; early-return patch not adopted** | @kevin9327; original Git metadata not yet recorded | — | Empty-capture and verbose-output tests; CI 43109905591 | Credit researched problem; no direct adoption |
| #2261 / `b1ce07f` | F03/F17 rune-safe masking | **CP-V5 adapted with in-source credit** | @arpitjain099; original Git metadata not yet recorded | No cherry-pick; RNC `920b3e2` contains adaptation | Unicode/report tests; CI 43109905591 | **Attribution required** in clean release |
| #2082 / `df8223e` | F04 rule/allowlist error-returning regex compilation | **CP-V6 adapted wrapper direction** | @Smeet23; original Git metadata not yet recorded | No cherry-pick; RNC `ed74f7e` includes correction | Standard/RE2 and CLI config tests; CI 43114548814 | **Attribution/code-overlap audit required** |
| #2241 / `0c2402a` | F04 invalid-regex fixtures | **Reviewed/informed; patch not adopted** | @harriiinnii; original Git metadata not yet recorded | — | Validation fixtures; CI 43114548814 | Credit scenario review; no direct adoption |
| #2228 / `637df0f` | F04 duplicate correction proposal | **Reviewed/not adopted** | @ankit090701; original Git metadata not yet recorded | — | CP-V6 CI 43114548814 | No code-adoption claim |

| #2317 / `71e9120` | F09 bounded archive signature verification/replay | **CP-V7 adapted, no cherry-pick** | @Varshith-Kali; Git metadata to verify | RNC `cccdef6`; no adopted external commit | CI 43121770371 | Exact code-overlap and credit audit pending |
| #2099 / `1b3896f` | F15 `MaxFileSize` propagation and known member checks | **CP-V7 partially adapted/informed** | @alephnull-sh; Git metadata to verify | RNC `cccdef6`; no adopted external commit | CI 43121770371 | Attribution and ZIP/7z semantics pending |

Do **not** call a PR "adopted" simply because a similar fix was independently implemented. If code was actually carried over, record the evidence and correct attribution before final integration.

## 6. Clean-history and upstream-delivery strategy

**Operator intent:** Iterate freely on feature/development branches, then deliver **one clean RNC-owned squash commit** for original work while preserving any outside contributors' code as separate authored commits. The final public diff and attribution should be intelligible to upstream reviewers.

### Release reconciliation procedure (operator-controlled; not authorized now)

1. **Freeze and audit:** Confirm final upstream/main base, selected fixes, feature flags, test evidence, licenses, and every overlapping community PR. Compare the fork against the actual chosen upstream SHA, not just a stale historical pin.
2. **Classify provenance:** For each reviewed PR, choose *independently implemented*, *adapted*, *adopted*, or *not selected*. Record what changed and who authored the original idea and code.
3. **Construct a clean staging branch** (working name `dev/clean-up`) from the approved release base. This should be performed only on explicit operator authorization; do not rewrite or disturb `fix/validation`.
4. **For adopted commits:** Fetch the contributor PR/commit, verify the actual author and patch, and cherry-pick with `-x` into the clean staging sequence. Preserve the contributor's author metadata; document original and new fork SHAs. If cherry-picking requires adjustments, use a separate RNC follow-up commit or an expressly documented adaptation, preserving credit.
5. **For RNC-owned work:** Rebuild one clean cumulative delta **excluding code already contained in adopted contributor commits**, and record it as one RNC-authored squash commit. It may be ordered before or after the adopted commits based on dependencies. Never duplicate somebody else's patch inside the squash and then attempt to cherry-pick it again.
6. **Validate the reconstructed tree:** Compare final code/output against the fully accepted working branch; run all test/CI/benchmark compatibility checks; inspect the complete diff and Git authorship graph.
7. **Preserve the final authorship graph:** Integrate the clean staging history into RNC `master` without squashing the separate contributor commits again. **No merge, rebase, branch deletion, or upstream PR in the current task.**

### Honest commit-count expectation

- If **no external code** is adopted, the fork can finish **one RNC implementation commit ahead** of the chosen upstream baseline (subject to the actual baseline state).
- If **N independently authored external commits** are deliberately adopted, the final history necessarily has **one RNC implementation commit plus N contributor commits** (and any justified integration adjustments). You cannot simultaneously preserve their separate authored commits and have exactly one total commit ahead.
- Merely mentioning a contributor in a squash commit, adding a `Co-authored-by` trailer, or linking an issue **does not** preserve their original commit as a separate history entry. Use each mechanism only for its accurate purpose.

### Final upstream PR attribution requirements

The future upstream PR description should contain:

1. Clear, bounded problems solved and reproducible before/after evidence.
2. Backward-compatibility statement and all intentional behavior changes/defaults.
3. Test matrix with Windows/Linux/macOS, standard/RE2, Docker, and relevant failure paths.
4. **Code incorporated:** each contributor, PR, original commit, resulting fork commit, and any adjustments.
5. **Reviewed/influential work not incorporated:** relevant authors and issues/PRs with a sentence explaining how they informed the solution.
6. Open limitations, independent follow-up issues, and a clean diff against upstream.

**Upstream caveat:** The checked `README.md` currently states Gitleaks is feature-complete and upstream accepts security patches only, with active product focus moving to Betterleaks. A future upstream PR is an operator decision, not an assumed acceptance path; keep the patch scoped and distinguish safety fixes from new optional QoL features.

## 7. Updating this file

Whenever a fix or QoL feature changes:

1. Update its row and status (**Planned**, **Implemented**, **Validated**, **Open**, **Deferred**).
2. Add the applicable CP finding, direct issue/PR links, and evidence/test links.
3. Record default/CLI/report/exit-code changes explicitly under the compatibility contract.
4. Review the upstream attribution table for missing contributors, conflicting proposals, and code overlap.
5. If outside code is incorporated, update the adopted-commit ledger immediately; preserve actual author metadata and SHA crosswalk.
6. Reconcile this ledger with release notes and the eventual upstream PR. Never silently remove someone from the acknowledgments solely because their exact patch was not selected.

This file tracks engineering evidence and attribution; it is **not** a request to adopt upstream work or a license to merge anything automatically.

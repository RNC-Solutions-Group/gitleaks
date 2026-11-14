# RNC-Mod Engineering Reference

> **Scope.** This document is the canonical reader-facing technical record for the RNC-Mod fork of Gitleaks. It describes the operational difference between upstream Gitleaks `master` at `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` and the V13 integration candidate, the engineering corrections beneath that behavior, the evidence supporting them, and the external contributions recognized by the project. Historical assessment reports remain available under `_assessment/`; they are evidence archives, not competing current specifications.

## Baseline and document metadata

| Item | Value |
|---|---|
| Repository | `RNC-Solutions-Group/gitleaks` |
| Original upstream baseline | `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (upstream `master`) |
| V10 product authority | `48c76717ccbbe2159185edb4510fb1924d7dbcc6`; tree `39b13207713d7102d7ac813a0599b2eaabf5d29f` |
| V11 integration baseline | `rnc-mod` at `197a7e5516b0029d0b7bb2554552fc9ed54efb68` |
| V12 development HEAD (pre-report-publication correction) | `1ff69b7836b32e4f738883a09bc2a3daeccbd69b`; tree `a348eea61670da139b89b410efa5156325091e87` |
| V12 RNC reconciliation | `9092703a25f4c10edc162699955a6993ec66ce37` |
| V12 report-publication correction | Final commit on PR #5; see the PR head for the exact commit and tree |
| V12 correction branch | `audit/v12-integration` |
| V13 starting baseline | `3a4635f82dcd0a216f9d2b88ebbd72341b3b60cc` (merged V12, PR #5) |
| V13 development branch | `fix/v13-operational-hardening` |
| V13 integration | `audit/v13-integration`, final PR to `rnc-mod` (pending operator review) |
| Default branch contract | MIT-licensed upstream-derived scanner; no default detection-rule expansion |
| Document status | Authoritative for V12 candidate behavior; release acceptance remains an operator decision |

**Contents**

- [1. Operational and Usability Changes](#1-operational-and-usability-changes)
- [2. Security, Reporting, Quality and Internal Engineering Changes](#2-security-reporting-quality-and-internal-engineering-changes)
- [3. Attribution Register](#3-attribution-register)

**Reading conventions.** “Upstream” means the original Gitleaks behavior at the baseline above. “RNC-Mod” means the integrated fork. “Correction” means a verified behavioral or implementation change. “Retained” means intentionally preserved upstream behavior, even where an opt-in extension exists. “Residual” means a known limitation that is not claimed as fixed. Historical checkpoint identifiers (`CP-V*`, `F*`, `V12-*`) are retained for traceability but are not current release labels.

# 1. Operational and Usability Changes

## 1.1 Operational change matrix

| Component | Upstream behavior | RNC-Mod behavior | Implementation | Compatibility | Verification |
|---|---|---|---|---|---|
| Report destination validation | Destination validation could mutate or truncate an existing report before scanning. | The requested destination is validated without modification before scanning; reports are rendered to restrictive temporary files and published only after a complete scan and successful render. | `cmd/report.go` (`newReportOutput`, temp render, final validation/rename) | Existing valid invocations still replace the requested report by default. | Report-preservation, TOCTOU, permission and CLI tests. |
| Default overwrite policy | Existing report files were overwritten. | `auto-overwrite` defaults to `true`; auto-overwrite remains the default. | `config/config.go`, `cmd/report.go` | Drop-in compatible for ordinary report paths. | Configuration/output-policy matrix; CI. |
| Optional aggregation | Repeated runs could collide or overwrite. | `aggregation = false` by default. When enabled, a numbered destination (`report.json`, `report_001.json`, …) is exclusively reserved and published; no content merge. | `cmd/report.go` allocation and reservation | Opt-in; default unchanged. | Numbering, extension, Unicode, concurrent-allocation, Windows/macOS/Linux CLI tests. |
| Explicit output overrides | No RNC `--overwrite` or `--aggregate` controls. | One-run overrides of the configured output policy; supplying both is an error; explicit overwrite never authorizes a destination that is also a scan input. | Root persistent flags; `resolveReportMode`; collision checks | New additive flags; configuration remains valid. | Four-mode matrix, collision and override tests. |
| Collision rejection | No complete policy for disabling both overwrite and aggregation. | When both configured policies are false, an existing destination is rejected before scanning; an absent destination writes normally. | `cmd/report.go`, `cmd/report_pair.go` | Opt-in via configuration; default remains overwrite. | Pre-scan rejection, preserved-report and CLI tests. |
| Informative coverage report | Findings reports did not show intentionally skipped targets. | `informative-report = true` enables a versioned skipped-target report. With a primary file it writes a `*-informative.json` companion; without a file, or with `console-only = true`, it reports to stderr. | `internal/informative`, `cmd/report.go`, `cmd/report_pair.go` | Disabled by default; primary findings/report bytes are unchanged when disabled or under console-only. | Schema, exclusion-reason, companion/console, aggregation and differential-fixture tests. |
| Informative timestamp | No scan-generation timestamp. | `scan.generated_at` is emitted as a UTC RFC3339 value only when informative reporting is enabled; the primary report schema is unchanged. | `internal/informative`, scan wiring | Additive JSON field within informative schema version 1. | Deterministic clock and schema tests. |
| Paired report publication | No companion publication model. | The companion and primary are reserved, validated, rendered and published as a logical pair. Overwrite failure after companion publication triggers identity-checked rollback. | `cmd/report_pair.go` | Companion is opt-in; single-report path is unchanged. | Pair overwrite/aggregate/reject, collision, symlink, rollback, concurrent-replacement and failed-restore tests. |
| Operational exit status | A scan could exit successfully after a source/Git/parser failure, and `--exit-code 0` could mask that. | Completed scans preserve the requested finding exit code; operational/incomplete scans return nonzero regardless of `--exit-code`, and failed scans do not publish partial output. | `cmd/scan_status_test.go`, `internal/scanerr`, detector/source propagation | Intentional correction; only failing scans change status. | Failure-path CLI tests on Windows, Linux and macOS CI. |
| Source and worker errors | Some directory, metadata, symlink, open, archive and read errors were silently dropped. | Errors are collected with paths through a bounded collector; successfully read fragments are retained but the scan is reported incomplete and exits nonzero. | `sources/files.go`, `sources/file.go`, `internal/scanerr` | Intentional correction; clean scans unchanged. | Unreadable files, partial reads, worker errors, archive errors and CI. |
| Git process and diff completeness | Nonzero Git exits, fatal stderr, malformed/truncated diffs and incomplete process cleanup could yield false-clean results. | Git is waited/reaped; fatal stderr is classified separately from warnings; parser leftovers and raw `diff --git` section counts are checked; errors are joined and partial reports suppressed. | `sources/git.go`, `gitStderr`, `gitOutputTracker`, `scanerr` | Intentional correction for failing/partial streams. | Malformed, truncated, overflowed-header, blocking-pipe, cancellation, exit-code and stderr tests. |
| Git presentation determinism | Repository `color`, `format.pretty`, `log.date`, textconv and external-diff settings could alter parser input or execute helpers. | History and worktree Git commands force `--no-color`, `--no-textconv`, `--no-ext-diff`, `--date=default` and `--pretty=medium`. | `sources/git.go` (`forceGitOutputFlags`) | Fingerprints and metadata become stable across hostile host configuration; no detection loss for ordinary repositories. | Forced-color, hostile-format/date, textconv-marker and staged-binary tests. |
| Merge-only history | Default `git log -p` omits merge patches, so merge-only resolutions can be missed. | Default remains unchanged for compatibility. Explicit `--log-opts=--diff-merges=first-parent` opts into merge attribution. Combined/remerge and patch-suppressing formats are rejected. | `sources/git.go`, `cmd/git.go`, README | Retained default; opt-in behavior may duplicate content already reported from a contributing commit. | Evil-merge, normal-merge, no-diff-merges, first-parent and unsupported-format tests. |
| Staged binary-marked text | Text staged under `.gitattributes -diff` could be silently omitted. | When the staged diff omits text, the stage-zero index blob is read and MIME-checked; genuine binary content remains excluded. | `sources/git.go` index-blob path | Intentional coverage correction for staged scans. | `core.bigFileThreshold`, `-diff`, index blob, MIME and error tests. |
| Archive identification | Filename extension alone could classify ordinary text as an archive and skip it. | Archive dispatch requires a bounded content signature; the consumed prefix is replayed for both seekable and non-seekable readers; verification errors are operational. | `sources/file.go` (`archivePeekSize`, content verification, prefix restore) | Fewer false skips; genuine archives still scanned according to depth/size policy. | Archive-like plaintext, real ZIP/gzip/tar, prefix replay and corrupt-archive tests. |
| File and archive size limits | Limits were uneven across readers, decompression and archive members, and some paths used rounded units. | Exact-byte `MaxFileSize` propagation to files, archives, members, decompression, Git and stdin; known-size prechecks; unknown/decompressed streams use bounded temporary spooling; zero means no explicit limit. | `sources/file.go`, `sources/files.go`, `sources/git.go`, `detect` | Intentional exactness correction; a value equal to the limit is scanned. | Size-boundary, archive-member, decompressor, spool, stdin and CLI tests. |
| First-chunk classification | Later chunks could reclassify an initially text file as binary and drop findings. | MIME classification is performed only on the first chunk; binary-at-start remains excluded. | `sources/file.go` (`isFirstChunk`) | Intentional correction of chunk reclassification. | Long first-line, binary signature, chunk-boundary tests. |
| Long-line fragmentation | A long line could be split at a chunk boundary, placing a secret in separate fragments. | The normal 25 KB blank-line lookahead remains; a bounded 1 MiB finish-current-line pass prevents splits through 100–300 KB and beyond. A line that cannot be finished within the cap is an operational error. | `sources/common.go` | Intentional correction; over-cap lines fail closed instead of reporting clean. | 100 KB, 125 KB, 125 KB+1, 150 KB, 250 KB, 300 KB and over-cap tests. |
| Directory positional arguments | Extra or empty arguments could silently fall back to scanning the current directory. | Zero arguments default to `.`; one argument is accepted; extra, empty or invalid arguments are rejected. | `cmd/directory.go`, `cmd/root.go` | Intentional fail-closed correction; no multi-root feature. | CLI argument matrix and cross-platform tests. |
| Explicit-file command | No command that scans exactly a filename list without directory walking. | `gitleaks files` requires one or more files, deduplicates cleaned paths, rejects missing files and directories, scans symlinks only when requested, and never expands a parent directory. | `cmd/files.go`, `sources/explicit_files.go` | New additive command. | Source, CLI and real pre-commit tests. |
| Pre-commit hooks | Staged hooks ignore forwarded filenames, so `--all-files` still scanned only staged changes. | Existing staged hooks remain unchanged; new explicit-file hooks forward `--files`/`--all-files` to `gitleaks files`. | `.pre-commit-hooks.yaml`, scripts/tests | Additive; existing hook IDs retain their contract. | Manifest, Python helper, real pre-commit 4.7.0 integration and CI. |
| Source-local configuration | Automatic `(target)/.gitleaks.toml` discovery could replace rules or suppress findings in untrusted repositories. | Discovery and precedence are preserved; `--no-local-config` disables only automatic discovery. Explicit `--config` and environment configuration still apply. | `cmd/root.go` (`initConfig`) | Additive control; default remains compatible. | Default-rule and precedence tests; CLI integration. |
| Ignore-file input safety | FIFO/special inputs, unbounded sizes and incomplete reads could hang or partially suppress. | Non-regular ignore inputs are rejected; Unix uses nonblocking open plus post-open verification; Windows checks regular files; a complete parse is merged only on success. Limits: 16 MiB file, 1 MiB line, 100,000 entries. Explicit missing path errors; optional discovery remains optional. | `detect/detect.go`, `detect/ignore_file_*.go` | Intentional safety correction; legitimate ordinary files unchanged. | FIFO/special-file, oversize, over-entry, partial-read and explicit-path tests. |
| Contextual redaction | Adjacent, multiline, empty-capture and auxiliary-field secrets could remain visible under redaction. | Rune-safe masking, empty-capture fallback, known-exposure grouping by file/commit, auxiliary required-rule redaction and limited disclosure paths. | `report/finding.go`, `detect/utils.go`, diagnostics wiring | `--redact=0` remains raw by default; redacted output may mask more known values. | Unicode, empty, adjacent, multiline, required-rule, verbose and template tests. |
| CSV report | Formula-like cells were written verbatim. | Every string cell is neutralized with a leading apostrophe when the first significant rune is `=`, `+`, `-`, `@` or a supported fullwidth variant; leading whitespace/control/format characters, BOM and zero-width characters are skipped for classification. | `report/csv.go` | Intentional output change for formula-like values; ordinary cells unchanged. | RFC 4180 round-trip, trigger table, idempotence, negative/normal-cell tests. |
| SARIF columns | Columns were byte offsets; non-ASCII prefixes produced incorrect character columns for SARIF consumers. | SARIF receives 1-based Unicode codepoint columns computed from the actual line start; JSON/CSV/console retain historical byte columns. | `detect/location.go`, `report/finding.go`, `report/sarif.go` | SARIF-only correction; baseline/JSON/CSV compatibility preserved. | ASCII, Arabic, CJK, emoji, combining marks, later-line, CRLF and multiline tests. |
| SARIF artifact URIs | Raw paths were emitted in `artifactLocation.uri`. | Relative/POSIX paths are single-escaped; Windows drive paths become `file:///C:/...`; UNC paths become `file://server/share/...`. | `report/sarif.go` | Output changes only for paths needing escaping or absolute URI representation. | Space, `%`, `#`, `?`, Unicode, drive, UNC and relative tests. |
| Baseline identity | LF and CRLF checkouts could produce different baseline `Match`/`Secret` values, entropy and end columns, so an LF baseline failed to suppress the same CRLF finding. | Detector spans trim CRLF, Match/Secret are canonicalized before entropy, and legacy baselines have trailing CR removed, EndColumn adjusted and entropy recomputed at load. | `detect/detect.go`, `detect/baseline.go` | Intentional cross-line-ending correction; path/fingerprint/line identity unchanged. | Real detector LF/CRLF equality, legacy baseline migration and negative tests. |
| Build contract | Build targets could fetch remote tags or mutate tracked module/config files. | `build`, `test`, `test-race`, `test-cover`, `failfast`, `verify`, `lint`, `version` and `profile` are non-mutating. `format`, `tidy`, `generate` and `clean` are explicit maintenance operations. Versions resolve from local state unless supplied. | `Makefile`, `cmd/makefile_contract_test.go` | Developer/CI build behavior change; release version override remains. | Makefile contract tests, offline/cached build checks, generated-config verification. |
| RE2 build support | Tagged full-repository tests could fail on generator `regexp` type mismatches. | Generator code imports the project regexp abstraction; standard and full `gore2regex` builds/tests are supported. | `cmd/generate/config/base`, `cmd/generate/config/rules`, `regexp/*` | Build-only correction; generated rules unchanged after normalization. | Full tagged suite, generator tests, application build and CI. |
| Container runtime | Historical root image only; `safe.directory=*` retained. | Default remains UID 0/root. `runtime-nonroot` is an additional tested target at UID/GID 10001 with `HOME=/home/gitleaks`. | `Dockerfile`, CI | Opt-in; default behavior unchanged. | amd64/arm64 root and non-root Docker builds and smoke tests. |
| CI authority | Broader or inconsistent trigger policy. | GitHub Actions runs only for pull requests to `rnc-mod` and pushes to `rnc-mod`; content read-only; actions/tool versions pinned. | `.github/workflows/test.yml`, `release.yml` | Governance change; routine development uses local CI. | Trigger inspection and 12/12 PR run. |
| Archive feature scope | General archive-engine expansion could continue indefinitely. | No additional archive-engine expansion in this fork; nonseekable ZIP/7z outer-size exclusions remain intentional, informative skips. Unknown-extension discovery remains out of scope. | `sources/file.go`, operator policy | Deliberate scope limit. | Archive size/identification tests and documented residuals. |
| Legacy `protect` command | Hidden compatibility command existed. | Remains hidden and compatibility-preserving, with command context propagated through its scan lifecycle. | `cmd/protect.go` | Unchanged CLI surface; not advertised. | Legacy command tests/inspection. |
| Effective configuration inspection | No safe way to see the resolved configuration. | `gitleaks config [path]` prints precedence, loaded files, rule/allowlist counts, report policy and limits without regex/allowlist patterns. | `cmd/config_summary.go` | New command; no scan behavior change. | CLI summary, local precedence and invalid-config tests. |
| Unused-ignore diagnostic | No visibility into loaded ignore entries that did not suppress a finding. | `--report-unused-ignore` prints a sorted, opt-in completed-scan diagnostic; suppressed entries are tracked; incomplete scans never certify entries unused. | `detect` ignore tracking, `cmd/root.go` | Off by default; suppression behavior unchanged. | Detector and CLI tests; incomplete-scan guard. |
| Local file links | Console showed only text paths. | `--file-links` adds an opt-in `file://` URI with line anchor for non-Git findings; `--no-color` suppresses terminal hyperlinks; machine reports are unchanged. | `detect/utils.go`, `cmd/root.go` | Off by default; absolute local paths remain console-only. | URI/Windows/no-color/JSON-isolation tests. |
| SARIF result level | Results omitted `level`, leaving severity to consumers. | `--sarif-level=none,note,warning,error` is opt-in; empty preserves historical output; no per-rule severity model is invented. | `report/sarif.go`, `cmd/root.go` | Additive flag; SARIF structure otherwise unchanged. | Report and CLI tests. |
| Quoted Git log arguments | `--log-opts` used `strings.Split` and could not represent quoted pathspecs. | A shell-free parser handles quotes, escapes, Windows backslashes, empty/unmatched input; disallowed formats are validated on the parsed vector. | `sources/git.go` | More options work; unsafe formats remain rejected. | Parser/vector/CLI-adjacent tests. |
| Shallow/bare/subdirectory Git provenance | Shallow history and source-scope expectations were implicit. | Shallow scans warn that pre-boundary history is unavailable; bare repositories warn that working-tree config/ignore discovery is unavailable; subdirectory history scans warn that Git history is repository-wide. | `sources/git.go` | Warnings only; no fingerprint or scope change. | Git fixture tests. |
| Azure SSH remote identity | Azure DevOps v3 SSH remotes were not normalized. | `git@ssh.dev.azure.com:v3/org/project/repo` maps to canonical `https://dev.azure.com/org/project/_git/repo`; credentials are stripped; link paths are URL-escaped. | `sources/git.go`, `detect/utils.go` | Existing GitHub/GitLab/explicit-platform behavior retained. | Remote normalization and link tests. |
| Explicit-file scalability | `gitleaks files` processed all files sequentially. | Bounded concurrency uses the detector semaphore, preserves exact-file scope, deduplicates paths, sorts findings deterministically and propagates per-file errors. | `sources/explicit_files.go`, `cmd/files.go` | Same file set and policy; order is now deterministic. | Bounded concurrency, error, dedupe and CLI tests. |
| RNC release identity | Release workflow published upstream Docker Hub image and used upstream references. | Release workflow is GHCR-only, RNC image identity, tag/format and base-version validation, protected `rnc-mod` ancestry, least privilege, no Docker Hub credentials. | `.github/workflows/release.yml`, `.pre-commit-hooks.yaml`, `RELEASING.md` | No artifact is published by V13; RNC image only exists after an approved release. | Static workflow scope test; YAML parse; tag-policy probe. |
| Installation and version policy | Upstream install commands could be mistaken for RNC distribution. | README distinguishes upstream tools from RNC source/release paths; `RELEASING.md` defines `vX.Y.Z-rnc.N`, base > `minVersion`, container tags and module-path consequences. | `README.md`, `RELEASING.md` | No module-path rename; no release created. | Documentation/link review; version validation probe. |
| Upstream Gitleaks Action workflow | Could be activated as a scanner workflow. | Not part of the RNC test pipeline; no Action implementation was adopted. | CI/workflow policy | Deliberate exclusion. | Workflow inspection. |

## 1.2 CLI command and option reference

Only RNC additions and RNC-modified semantics are listed. Upstream flags not discussed here retain their original names, defaults and values.

### 1.2.1 `gitleaks files`

| Field | Behavior |
|---|---|
| Purpose | Scan exactly the files named on the command line without walking parent directories. |
| Command form | `gitleaks files [flags] <file>...`; at least one argument is required. |
| Default behavior | No implicit target; omission is an argument error. |
| Accepted values | One or more file paths; `--follow-symlinks` optionally follows named symlinks. |
| Execution semantics | Cleans and deduplicates paths, rejects missing paths and directories, skips empty files, scans files in argument order, and never expands a directory. |
| Output | Normal findings summary and selected report format; `File` uses the cleaned supplied path. |
| Exit behavior | Findings use `--exit-code`; argument, missing-file, directory, read or scan-completeness failures return nonzero. |
| Platform considerations | Windows path cleaning and symlink handling follow the host filesystem; POSIX permissions are not interpreted as content. |
| Upstream compatibility | New command; no upstream command changes. |
| Limitations | Bounded concurrency is used; very large `--all-files` sets still incur per-file open/read/close overhead. Findings are sorted by file, line, column and rule for deterministic output. |

```bash
gitleaks files --redact --verbose src/config.yaml scripts/deploy.sh
gitleaks files --follow-symlinks --report-path findings.json --report-format json current-config
```

### 1.2.2 `--no-local-config`

| Field | Behavior |
|---|---|
| Purpose | Refuse to trust an automatically discovered source-local `.gitleaks.toml`. |
| Form | `--no-local-config[=true\|false]`, persistent on all scan commands. |
| Default | `false`; automatic discovery remains enabled. |
| Accepted values | Boolean. |
| Execution semantics | Skips only precedence step 4 (the target-local file). `--config`, `GITLEAKS_CONFIG` and `GITLEAKS_CONFIG_TOML` still apply in order. |
| Output | No output change; rule/allowlist selection changes only when a local file would otherwise have been loaded. |
| Exit behavior | Invalid explicit/environment configuration still fails before scanning. |
| Platform considerations | None beyond path handling for the target. |
| Upstream compatibility | Additive; default behavior is preserved. |
| Limitations | Operators must opt in for untrusted scans; the default still trusts a repository-local config. |

```bash
gitleaks dir --no-local-config --redact untrusted-checkout
```

### 1.2.3 `--overwrite` and `--aggregate`

| Field | Behavior |
|---|---|
| Purpose | Apply a one-run report output policy, overriding configuration. |
| Form | `--overwrite[=true\|false]`, `--aggregate[=true\|false]`. |
| Default | Neither flag set; resolved configuration applies (`auto-overwrite=true`, `aggregation=false` by default). |
| Accepted values | Boolean. |
| Execution semantics | `--aggregate` selects numbered output; `--overwrite` selects replacement. Supplying both is an error. `--aggregate=false`/`--overwrite=false` do not force the opposite policy; unresolved configuration applies. |
| Output | Overwrite keeps the requested name; aggregate writes `name.json`, `name_001.json`, and so on with the same extension. |
| Exit behavior | Both flags together or an unsafe destination fails before scanning; existing destination policy is enforced before any scan. |
| Platform considerations | Windows rename contention is retried within a bounded window. |
| Upstream compatibility | Additive flags; default overwrite behavior matches historical operation. |
| Limitations | Aggregation numbers files; it does not merge report contents. Paired report publication is not crash-atomic. |

```bash
gitleaks git --overwrite --report-path findings.json .
gitleaks git --aggregate --report-path findings.json .
```

### 1.2.4 `--log-opts`

| Field | Behavior |
|---|---|
| Purpose | Restrict or extend `git log` history selection. |
| Form | `--log-opts="<git log options>"` on `git` and legacy `protect`. |
| Default | Empty; normal full-history options are supplied by the scanner. |
| Accepted values | Ordinary patch-producing Git log options, including ranges and `--diff-merges=first-parent`. |
| Execution semantics | A shell-free tokenizer parses quotes, escapes, empty/unmatched input and Windows backslashes; the parsed vector is used for execution and format validation. RNC safety/presentation flags are inserted before a `--` separator. Patch-suppressing and combined/remerge formats are rejected. |
| Output | Findings correspond to added lines in ordinary patch output; merge opt-in may report the same introduced content from multiple commits. |
| Exit behavior | Unsupported formats fail before Git starts; Git/parse failures are operational errors. |
| Platform considerations | Git version must support the options supplied; `--diff-merges` requires Git 2.31 or later. |
| Upstream compatibility | Default remains compatible; the rejection of silent-clean formats is intentional. |
| Limitations | Merge-only default coverage is not complete; opt-in may duplicate findings. The parser is not a shell and does not expand variables or globs. |

```bash
gitleaks git --log-opts="--all --diff-merges=first-parent" .
gitleaks git --log-opts="v1.0.0..v1.1.0" .
```

### 1.2.5 `dir` positional target

| Field | Behavior |
|---|---|
| Purpose | Scan one directory or one file. |
| Form | `gitleaks dir [path]`; aliases `file` and `directory`. |
| Default | No argument scans `.`. |
| Accepted values | Zero or one path. |
| Execution semantics | Empty explicit argument, extra arguments, missing paths and invalid targets are rejected rather than silently defaulting. |
| Output | Reported paths use the supplied root spelling; root-relative matching additionally evaluates a normalized relative path. |
| Exit behavior | Argument/metadata errors are nonzero operational failures. |
| Platform considerations | Absolute/dot/relative roots and Windows cross-volume paths are supported for matching; reported identity is not rewritten. |
| Upstream compatibility | Intentional fail-closed correction; ordinary single-root scans unchanged. |
| Limitations | No multi-root scanning feature; rule-level `path` regex still sees reported paths rather than root-relative paths. |

```bash
gitleaks dir --redact .
gitleaks dir --redact ./checkout
```

### 1.2.6 Report path, format and template

| Field | Behavior |
|---|---|
| Purpose | Select file/stdout output and serialization. |
| Form | `--report-path/-r`, `--report-format/-f`, `--report-template`. |
| Default | No report path; console summary only. With a file path and no explicit format, `.csv`/`.json`/`.sarif` are inferred; other extensions fail closed. Template rendering requires `--report-format=template` together with `--report-template`; mismatches fail before destination inspection. |
| Accepted values | Path or `-`; `json`, `csv`, `junit`, `sarif`, `template`. |
| Execution semantics | File output is validated, rendered to a temporary file, revalidated, then published according to output policy; stdout is not paired with a file companion. |
| Output | Findings report; with informative reporting and a file destination, a sidecar is also published unless console-only. |
| Exit behavior | Invalid path, format/template, collision or publication failure is nonzero; finding exit code remains separate. |
| Platform considerations | Restrictive temporary-file permissions on POSIX; Windows rename/sharing retries; symlink/special-file destinations rejected. |
| Upstream compatibility | Ordinary report bytes/format remain compatible; hostile destinations fail earlier. |
| Limitations | Two-file publication is not crash-atomic; template content is caller-controlled and may itself expose fields. |

```bash
gitleaks git -f sarif -r findings.sarif .
gitleaks git -r - -f json . | jq .
```

### 1.2.7 `--exit-code`

| Field | Behavior |
|---|---|
| Purpose | Select the status used when leaks are found. |
| Default | `1`. |
| Accepted values | Integer exit code. |
| Execution semantics | Applies to completed scans with findings. Operational/incomplete scans return their own nonzero status and are not converted to 0. |
| Output | No change. |
| Exit behavior | `0` suppresses only the findings status, never an operational failure. |
| Platform considerations | Standard process exit-code range. |
| Upstream compatibility | Intentional correction for incomplete scans; completed clean/finding behavior matches the contract. |
| Limitations | Shell callers that only test zero/nonzero cannot distinguish finding status from operational status except through stderr/messages. |

### 1.2.8 `--gitleaks-ignore-path`

| Field | Behavior |
|---|---|
| Purpose | Select a `.gitleaksignore` file or the directory containing one. |
| Default | `.`; automatic discovery. |
| Accepted values | Existing regular file or directory; explicit missing path is an error. |
| Execution semantics | Complete parse required; CRLF and backslashes normalized for fingerprints; entry limits enforced; non-regular inputs rejected. |
| Output | Suppressed findings disappear as before; malformed entries warn where supported. |
| Exit behavior | Explicit missing/unsafe/oversized ignore input is a nonzero operational failure. |
| Platform considerations | Unix nonblocking FIFO defense; Windows regular-file checks. |
| Upstream compatibility | Intentional safety correction; valid ignore files unchanged. |
| Limitations | Metadata race/TOCTOU and extremely large legitimate ignore files remain constraints. |

### 1.2.9 Resource controls

| Option | Default | Semantics and limitations |
|---|---|---|
| `--max-target-megabytes` | `0` (no explicit limit) | Exact byte limit for files/streams/members; equal-to-limit is scanned; oversized targets are intentional exclusions and may be reported informatively. |
| `--max-decode-depth` | `5` | Recursive decoding passes; depth-linear allocation is possible for encoded input. No unproven cap was added. |
| `--max-archive-depth` | `0` (no traversal) | Nested archive traversal limit; depth exclusions are intentional skips. |
| `--timeout` | `0` (no timeout) | Command context timeout; blocking OS reads may outlast it on some platforms. |
| `--redact` | `0`, but bare `--redact` means `100` | Redacts findings and known diagnostic surfaces; unknown adjacent text can remain. |
| `--follow-symlinks` | `false` | Applies to `dir` and `files`; symlink targets are resolved, reported with the link path where applicable, and directory targets are not followed as a tree. |

### 1.2.10 Legacy `protect` command

`protect` remains a hidden compatibility command with `--source`, `--staged` and `--log-opts`. Its scanner context propagation and exit handling follow the current `git` implementation. It is not advertised for new integrations; use `gitleaks git` instead.

### 1.2.11 `gitleaks config [path]`

| Field | Behavior |
|---|---|
| Purpose | Print a safe summary of the effective configuration. |
| Form | `gitleaks config [flags] [path]`; the optional path is the discovery/precedence target (default `.`). |
| Default | No scan; summarizes the configuration that a scan of the target would use. |
| Accepted values | One optional path. |
| Execution semantics | Resolves `--config`, environment sources, target-local `.gitleaks.toml` and extension files; reports loaded files, rule/disabled/ordered counts, global/rule allowlist counts, minVersion/title, report policy and resource limits. |
| Output | Human-readable summary to stdout. Rule regexes, allowlist patterns and secret values are intentionally omitted. |
| Exit behavior | Invalid/malformed/unreadable configuration exits nonzero; no partial scan is started. |
| Platform considerations | Loaded paths retain host spelling; no absolute path is synthesized beyond the configuration actually loaded. |
| Upstream compatibility | New command; no scan/report schema change. |
| Limitations | It is a summary, not a lossless configuration export. Unknown TOML fields may still be ignored by the loader. |

### 1.2.12 `--report-unused-ignore`

| Field | Behavior |
|---|---|
| Purpose | Diagnose loaded ignore entries that did not suppress a finding. |
| Default | `false`. |
| Accepted values | Boolean. |
| Execution semantics | Tracks entries that actually suppress a finding. After a completed scan only, prints a sorted list of entries that were not triggered. Incomplete scans print no certification. |
| Output | stderr diagnostics; findings report schemas are unchanged. |
| Exit behavior | Diagnostic only; does not change finding/operational exit status. |
| Platform considerations | Fingerprint text uses the scanner's normalized `/` path spelling. |
| Upstream compatibility | Additive opt-in. |
| Limitations | “Not triggered” is not proof of obsolescence: an unexamined path, commit or revision can make an entry appear unused. Treat the list as confidential operational metadata. |

### 1.2.13 `--file-links`

| Field | Behavior |
|---|---|
| Purpose | Offer local clickable file URIs in verbose console output. |
| Default | `false`. |
| Accepted values | Boolean. |
| Execution semantics | For non-Git findings with a file path, prints a separate `File-Link:` `file://` URI with `#Lline`. No finding/report field is changed. |
| Output | Console only. Machine-readable reports remain unchanged. |
| Exit behavior | Diagnostic only. |
| Platform considerations | Windows drive paths use `file:///C:/...`; POSIX paths use ordinary file URIs; UNC-like paths follow URL rules. |
| Upstream compatibility | Additive opt-in; `--no-color` suppresses OSC 8 terminal escapes. |
| Limitations | Local absolute paths may be sensitive; do not enable in shared CI logs. Decoded-value highlighting remains a separate residual. |

### 1.2.14 `--sarif-level`

| Field | Behavior |
|---|---|
| Purpose | Emit an explicit SARIF `level` for findings. |
| Default | Empty; historical SARIF omits `level`. |
| Accepted values | `none`, `note`, `warning`, `error` (empty allowed). |
| Execution semantics | Only valid with `--report-format=sarif`; sets the same level on each result. No per-rule severity classification is inferred. |
| Output | SARIF result `level` field. |
| Exit behavior | Invalid level or use with a non-SARIF report fails before scanning. |
| Platform considerations | None. |
| Upstream compatibility | Additive opt-in; locations, rule IDs and fingerprints unchanged. |
| Limitations | A constant level does not model differing credential risk. |

### 1.2.15 Configuration-only output controls

The following are configuration keys, not CLI flags:

```toml
auto-overwrite = true
aggregation = false
informative-report = false
console-only = false
```

`auto-overwrite` and `aggregation` affect file reports only. `informative-report` enables skipped-target accounting; `console-only` routes that accounting to stderr and suppresses the companion file. No report path means console-only informative output regardless of `console-only` when the feature is enabled.

## 1.3 Git scanning behavior

### 1.3.1 History and merge handling

A `git log -p` scan examines lines added by each commit in range, not unchanged pre-existing file content. A credential committed before a range can remain unseen if a later in-range commit changes only another line. Periodic full-history scans remain necessary for complete repository coverage.

Default history mode does not emit merge patches. This is retained because adopting first-parent output by default changes finding multiplicity, fingerprints, baselines and ignore behavior. The supported opt-in is:

```bash
gitleaks git --log-opts="--all --diff-merges=first-parent" .
```

First-parent output can report the same introduced content from the merge and from a contributing commit. The scanner does not deduplicate those cross-commit findings. `--no-diff-merges` remains available. `--diff-merges=cc`, `combined`, `dense-combined`, `remerge`, `--remerge-diff`, `--cc`, `-c`, `--name-only`, `--name-status`, `--raw`, `--numstat`, `--shortstat`, `--no-patch` and `-s` are rejected because they suppress or alter the patch stream beyond what the parser can scan completely.

### 1.3.2 Argument parsing and safe boundaries

`--log-opts` is tokenized without a shell. Quoted multiword pathspecs, escaped whitespace, literal quotes, empty arguments, unmatched quoting and Windows backslashes have defined behavior; invalid input fails before Git starts. Disallowed output modes are checked against the parsed vector, so quoting cannot bypass the fail-closed format policy. The parser does not expand variables, globs or command substitutions.

### 1.3.3 Deterministic metadata and repository-controlled helpers

History and diff command construction pins:

```text
--no-color
--no-textconv
--no-ext-diff
--date=default
--pretty=medium
```

Safety flags are inserted before a user `--` pathspec separator. This prevents repository `format.pretty`, `log.date`, color settings and textconv drivers from corrupting metadata or executing repository-controlled commands on the scanner host. User `--log-opts` cannot re-enable textconv through the supported option parser; this is a deliberate security boundary.

### 1.3.4 Staged and unstaged scanning

`git --staged`/`--pre-commit` and unstaged history/worktree paths use the same parser-integrity and error model. When a staged text file is marked binary by Git attributes, the scanner reads the stage-zero index blob and verifies MIME classification before deciding whether to scan or exclude it. Genuine binary content remains excluded. `--pre-commit` retains the staged-diff contract; it is not a directory scan.

### 1.3.5 Failure, cancellation and completeness

The scanner waits for Git and its parser goroutines, drains remaining stdout, records unparsed bytes and raw diff-section counts, classifies fatal stderr, and joins operational errors. Cancellation waits for workers before returning. On a failed/incomplete scan:

- no partial report is published;
- an existing report remains untouched;
- `--exit-code 0` does not turn the failure into success.

Blocking OS reads may outlast a context timeout; this is a residual platform limitation rather than a universal cancellation guarantee.

### 1.3.6 Fingerprints and baselines

Commit identity, path, rule ID and line number remain the fingerprint basis. Deterministic pretty/date output restores stable commit metadata under hostile host configuration. CRLF normalization affects baseline content comparison only: path, line and fingerprint identity are unchanged. Merge opt-in can produce multiple findings for the same content because the contributing commit and merge are distinct commit identities.

### 1.3.7 Shallow, bare and subdirectory provenance

Git history is repository-wide even when the source is a subdirectory. A warning identifies that scope and recommends `--log-opts` pathspec filtering; directory scans remain bounded. A shallow repository warning states that history before the shallow boundary is unavailable and is not a complete full-history claim. A bare repository warning states that automatic working-tree configuration/ignore discovery is unavailable; explicit paths must be supplied. None of these warnings rewrites commit identities or fingerprints.

### 1.3.8 Remote identity and links

SCP-style remotes are normalized to HTTPS, credentials are stripped, and Azure DevOps `ssh.dev.azure.com:v3/org/project/repo` remotes map to the canonical `https://dev.azure.com/org/project/_git/repo` form. Link paths use URL escaping rather than ad-hoc substitutions; GitHub, GitLab, Azure DevOps, Gitea and Bitbucket behavior uses the existing platform templates. Invalid/absent remotes skip links without fabricating a URL.

### 1.3.9 Path identity

Directory scans report the file spelling derived from the supplied root, while allowlist, ignore and baseline matching also evaluates the path relative to that root. This makes `dir .`, `dir checkout` and an absolute-root scan apply the same relative policy without rewriting reported `File` or `Fingerprint` values. Git and single-file paths retain their existing semantics. Rule-level `path` regex evaluation remains based on the reported path; that limitation is tracked separately.

## 1.4 Pre-commit and integration behavior

### 1.4.1 Hook compatibility matrix

| Hook ID | Language | Filename forwarding | Scan scope | Default use |
|---|---|---|---|---|
| `gitleaks` | Go | `pass_filenames: false` | Staged index via `gitleaks git --pre-commit --staged` | Historical default hook. |
| `gitleaks-docker` | RNC GHCR image | `false` | Staged index | Containerized hook; image exists only after an RNC release. |
| `gitleaks-system` | System | `false` | Staged index | Historical hook using an installed `gitleaks`. |
| `gitleaks-files` | Go | `true` | Exactly forwarded files via `gitleaks files` | New explicit-file mode. |
| `gitleaks-files-system` | System | `true` | Exactly forwarded files | New explicit-file mode with installed binary. |

The three staged hooks deliberately do **not** forward pre-commit filenames: forwarding them as positional Git targets would change the command contract. Consequently, `pre-commit run --all-files` on a staged hook still scans only staged changes. The new explicit-file hooks exist for callers who want pre-commit's file selection honored.

### 1.4.2 Explicit-file behavior

`gitleaks files` receives the pre-commit file list as positional arguments. It:

- deduplicates cleaned paths;
- rejects directories and missing paths;
- skips empty files;
- follows symlinks only with `--follow-symlinks`;
- scans each file once;
- never expands a directory or falls back to `.`.

Exit codes follow normal findings/operational rules. Real pre-commit integration has been exercised with a clean multi-file set, a committed secret under `--all-files`, a named secret, missing/deleted-file filtering, and directory rejection.

### 1.4.3 Performance constraint

Explicit-file scanning is sequential. It avoids a repository-wide expansion and is correct for the filename contract, but a very large `--all-files` invocation is slower than a tuned directory walk. This is a documented constraint, not a scope fallback.

## 1.5 Configuration behavior

### 1.5.1 Precedence

1. `--config/-c`
2. `GITLEAKS_CONFIG` (path)
3. `GITLEAKS_CONFIG_TOML` (inline content)
4. `(target)/.gitleaks.toml`
5. Embedded default configuration

`--no-local-config` skips only step 4. Explicit and environment sources retain their precedence. Invalid explicit/environment configuration fails before scanning; it never falls back silently to defaults.

### 1.5.2 Extension and validation

Configuration extension is bounded to two levels. Every file loaded from disk is recorded in `Config.ConfigFiles`, including extended files, so report destination collision validation covers the entire active configuration set. Inline TOML remains pathless. User rule and allowlist regexes compile through error-returning wrappers on both standard and RE2 backends; invalid patterns fail with context instead of panicking. `minVersion` is checked when present.

### 1.5.3 Allowlists, ignores and baselines

- Global and rule-level allowlists support path/commit/regex/stopword conditions; AND path matches no longer short-circuit other required conditions.
- Boolean/null placeholder allowlists are anchored to the whole value, so values merely containing `true`, `false` or `null` are not suppressed.
- Global commit allowlists skip the matching diff file/commit, including the first commit in a multi-commit stream.
- `.gitleaksignore` contains fingerprints, not secrets; it is parsed completely and bounded before use.
- Baseline JSON is a findings report and may contain plaintext. It is sensitive data. Baseline path matching supports equivalent directory roots; CRLF content is canonicalized.

### 1.5.4 Effective configuration inspection

`gitleaks config` resolves the same precedence and extension graph as a scan and reports counts/metadata without exporting regexes or allowlist patterns. It is a safe summary, not a lossless export. The command performs no scan and does not initialize a Viper reset or process-wide race.

### 1.5.5 Trust boundary

The default remains compatible: a repository-local `.gitleaks.toml` is trusted and can replace rules or suppress findings. This is an intentional compatibility decision, not a security guarantee. Untrusted-scan orchestration should pass `--no-local-config` or an explicit trusted configuration. The flag is additive; it does not change resolution order when an explicit/environment source is present.

## 1.6 Filesystem and archive behavior

### 1.6.1 Directory traversal and target selection

`gitleaks dir` accepts zero or one target. A directory is walked; one file is scanned as a file. Symlinks are not followed unless `--follow-symlinks` is set. Metadata, open, read, archive and close errors are operational failures rather than silently skipped targets. Intentional exclusions (empty file, binary content, size limit, archive depth, path allowlist, symlink not followed, excluded input) are distinguishable in informative accounting.

### 1.6.2 File identity and paths

Reported paths retain the supplied root spelling. Root-relative matching evaluates a normalized path for allowlists, ignore fingerprints and baselines. A baseline/configuration file that is also a scan target is excluded from the scan to avoid reading its own sensitive content. Input/report collisions — including symlink and hardlink identity — are rejected before scanning.

### 1.6.3 Binary classification

MIME classification is performed on the first chunk only. A file that starts as text remains text for later chunks; a binary signature at the start excludes the file. This avoids reclassifying later chunks as binary and dropping findings, while preserving genuine binary exclusion.

### 1.6.4 Archives

Archive handling requires content verification within a 64 KB prefix. Genuine ZIP/gzip/tar and other supported formats are traversed up to `--max-archive-depth`; member and decompressed sizes obey exact-byte limits. Corruption, read failure and verification failure are operational errors. Ordinary text with an archive-like name is scanned as text.

The fork deliberately does not expand the archive engine further. Nonseekable ZIP/7z containers whose outer size exceeds the per-file limit can be excluded even when all members are small; this is recorded as a `size_limit` intentional exclusion. Self-extracting/prefixed archives may be treated as text because only the leading signature is checked. Unknown-extension archive discovery is out of scope.

### 1.6.5 Fragments and long lines

Files are read in approximately 100 KB fragments with a 25 KB blank-line lookahead. When the lookahead ends inside a line, a 1 MiB finish-current-line pass prevents a secret from being split across fragments. If the line cannot be finished within the cap, the scan returns an operational error rather than reporting a clean partial scan.

### 1.6.6 Skipped versus failed targets

Intentional exclusions are successful scan decisions with accounting. Read/open/metadata/parse failures are operational failures. Informative output distinguishes reasons such as `size_limit`, `binary_content`, `archive_depth`, `path_allowlist`, `commit_allowlist`, `symlink_not_followed`, `empty_file` and `excluded_input`. Observed/inspected counts are scoped to encountered targets; they are not a census of unvisited subtree content.

## 1.7 Operational compatibility summary

| Category | Behavior relative to upstream |
|---|---|
| Unchanged defaults | Default `dir` single-root scanning, default Git history selection, default report overwrite, default rules, default exit code for completed scans, standard JSON/CSV field names, standard console output. |
| Intentional corrections | Operational failures now return nonzero; failed reports are not published; directory arguments fail closed; Git metadata is deterministic; textconv is not executed; staged binary-marked text is scanned; archive/MIME false skips are corrected; root-relative matching is stable; CRLF baselines work; redaction/CSV/SARIF defects are corrected. |
| Newly added opt-in capabilities | `--overwrite`, `--aggregate`, informative report/console-only, `--no-local-config`, `files` command, explicit-file hooks, optional non-root container. |
| Intentional default differences | Formula-like CSV cells are apostrophe-prefixed; CRLF findings normalize Match/Secret/entropy; hostile Git settings are overridden; over-cap long lines fail closed; unsupported Git log formats fail closed. |
| Retained limitations | Merge-only default coverage; source-local config trust; two-file crash atomicity; recursive decode resources; archive outer-size/self-extracting limits; sequential explicit-file scanning; blocking reads; rule-level path regex semantics; Windows rename locks beyond retry. |
| Detection parity boundary | V10 established object-level product equivalence with its accepted tree and differential report parity on supplied fixtures. V11 and V12 intentionally change line boundaries, allowlists, baselines, long-line handling and merge opt-in behavior. Do not claim universal historical detection parity after V11/V12. |
| Release boundary | Source/build/test validation is not a released-binary attestation, a completed scanner-harness accuracy study, or an upstream acceptance commitment. |

# 2. Security, Reporting, Quality and Internal Engineering Changes

## 2.1 Engineering correction matrix

| ID / Area | Original defect or limitation | Root cause | RNC correction | Source components | Validation | Final disposition |
|---|---|---|---|---|---|---|
| F01 / V3 | Source, walker, Git, parser and worker errors could be swallowed, producing false-clean scans. | Error returns were dropped at multiple layers; workers were not always joined. | Path-aware bounded error collector; worker joins; fatal Git stderr classification; raw-section completeness check; operational status cannot be overridden to 0. | `internal/scanerr`, `sources/files.go`, `sources/git.go`, `detect/detect.go`, CLI | Unreadable files, malformed/truncated diff, blocking pipe, cancellation, exit-code tests; CI | Corrected for CLI-used paths; deprecated APIs/blocking reads residual. |
| F02 / V1-A / V2 | Report preflight could delete/truncate old output; partial final writes; unsafe destinations. | `os.Create` before validation; direct final write. | Validate without mutation; restrictive temp render; revalidate; rename publication; collision/symlink/special-file checks; Windows retry. | `cmd/report.go`, `cmd/report_rename*.go`, `cmd/root.go` | Preservation, permissions, TOCTOU, hardlink, concurrent writer, platform tests | Corrected for handled failures; final-publication race/orphan temp residual. |
| F02 / V7A / V12-10 | Paired report could leave a new companion beside an old primary; rollback could overwrite a concurrent writer. | Two independent paths; backup restoration did not re-check current identity. | Reserve/validate/publish pair; backup existing companion; identity + size/modtime snapshot check before remove/restore; retain backup on failed restore or concurrent replacement. | `cmd/report_pair.go` | Overwrite/aggregate/reject, concurrent replacement, failed restore, no-prior-companion tests | Handled-failure correction; crash atomicity impossible. |
| F03 / V5 / V12-06 | Secret-bearing diagnostics and incomplete redaction across multiple, empty, Unicode, adjacent, multiline and auxiliary fields. | Byte slicing; exact-line grouping; empty captures skipped; unbounded diagnostic fields. | Rune-safe mask; empty-capture fallback; file/commit exposure grouping; auxiliary required-rule redaction; bounded sanitized diagnostics. | `report/finding.go`, `detect/utils.go`, `sources/git.go`, diagnostics | Unicode/empty/adjacent/multiline/template/verbose and Git stderr tests | Known-reporting surfaces corrected; unknown adjacent text/partial masks residual. |
| F04 / V6 | Malformed user regex could panic; extension depth leaked globally; local config trust unclear. | `MustCompile`; package-global extension state; implicit discovery. | Error-returning compile wrappers; explicit extension depth; config-file tracking; additive `--no-local-config`. | `config/config.go`, `regexp/*.go`, `cmd/root.go` | Invalid regex/extension/precedence and CLI tests | Corrected defaults; local-config trust remains a policy decision. |
| F05 / V6 | Unsafe `.gitleaksignore` types, blocking FIFO, incomplete reads, unbounded input. | Direct file read/scan without type/bounds. | Regular-file checks, Unix nonblocking open + fstat, complete parse before merge, 16 MiB/1 MiB/100k limits, explicit missing error. | `detect/detect.go`, `detect/ignore_file_*.go`, `cmd/root.go` | FIFO/special/oversize/partial/explicit-path tests | Corrected; TOCTOU/extreme legitimate files residual. |
| F06 / V1-B | Global AND allowlist path match suppressed content before other required conditions. | Path pre-skip ignored AND sufficiency. | Path skip only for OR/path-only AND; evaluate content/commit-dependent conditions during detection. | `sources/common.go`, `detect/detect.go` | OR/AND/nil/commit/stopword/root tests | Corrected. |
| F07 / V4 | Text secret staged under binary attributes omitted. | Parser skipped binary diff without index fallback. | Stage-zero index blob read plus MIME verification; genuine binary still excluded. | `sources/git.go` | Big-file-threshold, `-diff`, index, MIME, error tests | Corrected. |
| F08 / V1-C / V12-02/03 | Forced Git color, pretty format, date, textconv or ext-diff could corrupt parsing or execute helpers. | Only color was partially handled; host config trusted. | Force `--no-color`, `--no-textconv`, `--no-ext-diff`, `--pretty=medium`, `--date=default`. | `sources/git.go` | Forced color, hostile format/date, textconv marker, diff/log tests | Corrected. |
| F09 / V7 | Filename-only archive classification skipped ordinary text. | Extension trusted without content check. | 64 KB content signature, prefix replay, explicit verification error. | `sources/file.go` | Archive-like plaintext, real archive, corrupt/prefix tests | Corrected within scope. |
| F10 / V1-D | Later chunks could classify an initially text file as binary. | Classification repeated on every chunk. | First-chunk-only MIME guard. | `sources/file.go` | Long-first-line/binary/chunk tests | Corrected. |
| F11 / V12-01 | Merge-only changes missed by default history. | Git emits no patch for merges by default. | Retain default; test/document first-parent opt-in; reject unusable merge formats. | `sources/git.go`, `cmd/git.go`, README | Evil-merge/normal-merge/no-diff-merges tests | Policy decision; default remains incomplete by design. |
| F12 / V1-E | Invalid/multiple `dir` arguments silently scanned `.`. | Cobra argument contract/default fallback. | Zero/one argument, empty/missing/extra rejection. | `cmd/directory.go` | CLI argument matrix | Corrected. |
| F13 / V1-F / V9 | Pre-commit helper/system hook could pass on nonzero or forward filenames to Git. | Subprocess status/pass_filenames contract. | Portable subprocess status; staged hooks `pass_filenames: false`; explicit-file hooks added separately. | `scripts/pre-commit.py`, `.pre-commit-hooks.yaml`, `cmd/precommit_test.go` | Python helper, manifest and real pre-commit tests; CI | Corrected. |
| F14/F18 | Report/output contracts and Unicode handling residuals. | Mixed formatting/redaction/column assumptions. | JSON/CSV schema preserved; SARIF codepoint columns; URI escaping; redaction fixes. | `report/*.go`, `detect/location.go` | Format matrices and CLI tests | Scoped corrections; template/caller surfaces residual. |
| F15 / V7 | Incomplete size propagation and broad in-memory decompression. | Inconsistent fields and buffers. | Exact-byte propagation; prechecks; bounded temp spooling; zero unlimited; member/decompressor coverage. | `sources/file.go`, `sources/files.go`, `sources/git.go` | Boundary/member/decompressor/spool tests | Corrected; nonseekable outer cap accepted. |
| F16 / V8 | Equivalent roots produced different ignore/allowlist/baseline matches. | Reported-path-only comparisons. | Root-relative matching for policy while preserving reported path/fingerprint. | `sources/fragment.go`, `sources/files.go`, `detect/baseline.go`, `detect/detect.go` | Dot/relative/absolute/cross-volume/archive tests | Corrected; rule-level path regex residual. |
| F17 / V5 | Unicode redaction split UTF-8. | Byte-length slicing. | Rune-based masking and bounds/percent handling. | `report/finding.go` | Multibyte/mixed-width/UTF-8 validity tests | Corrected. |
| F19 / V3 | Workers could outlive cancellation; legacy protect lost context. | Early return/deprecated API context. | Wait for workers; propagate command context. | `detect/git.go`, `cmd/protect.go`, sources | Cancellation/lifecycle tests | Corrected; blocking OS reads residual. |
| F20 / V2 | Report permissions/close/flush weaknesses. | Direct create/write without restricted temp. | Temp files `0600` on POSIX; explicit close/flush errors; atomic rename. | `cmd/report.go` | Permission/close/publication tests | Corrected for publication path. |
| F21 / V5 | Git stderr could expose secrets in logs/errors. | Raw stderr propagated. | Bounded sanitizer for URL userinfo/query/token shapes; fatal/warning classification. | `sources/git.go` | Git stderr credential tests | Heuristic correction; unknown short formats residual. |
| F22 / V6 | Extension state leaked across loads. | Package-global `extendDepth`. | Explicit depth through translation/extension. | `config/config.go` | Sequential extension tests | Corrected; global Viper concurrency residual. |
| F23 / V2A | Config collision validation missed environment/extended files. | `cfg.Path` reset after load. | `ConfigFiles` tracking and full collision validation. | `config/config.go`, `cmd/report.go` | Config collision matrix | Corrected. |
| V2 | Report destination races and unsafe overwrite. | Direct writes and shared names. | Temp render, exclusive reservation, aggregate pairs, collision checks, retry. | `cmd/report.go`, `cmd/report_pair.go` | Concurrency/TOCTOU/platform tests | Corrected for handled failures; pair crash window residual. |
| V5-C | Large excluded staged blob was fully drained. | Reader had no early-close path. | EOF tracking and `CloseAfterSkip` producer termination/reap. | `sources/git.go` | Blocked-read/lifecycle/performance probe | Corrected. |
| V7A | No skipped-target accounting; no optional companion. | Findings-only model. | Versioned informative collector, exclusion reasons, deterministic JSON/console, companion publication. | `internal/informative`, `cmd/report*.go`, `sources/*` | Schema/reason/console/pair tests | Added opt-in; counts are observed-target scoped. |
| V9-A | Windows report rename sharing violations. | Immediate rename under transient locks. | Bounded Windows-only retry for access/sharing errors. | `cmd/report_rename_windows.go` | Stress/retry tests | Mitigated; locks beyond window residual. |
| V9-B | Informative output lacked generation time. | No provenance field. | UTC RFC3339 `scan.generated_at`, informative only. | `internal/informative` | Deterministic clock tests | Added. |
| V10-A | Build mutated tree/fetched tags. | Build recipe side effects. | Non-mutating build/test/verify; explicit maintenance; local SemVer resolution. | `Makefile` | Makefile contract/offline/cached builds | Corrected. |
| V10-B | Full RE2 tagged generator build failed. | Generator imported stdlib regexp types. | Project regexp abstraction in generator. | `cmd/generate/config/base`, `rules` | Full tagged suite/CI | Corrected. |
| V10-C | Real pre-commit integration skipped in CI; Windows cross-volume fixtures. | Runner setup/path assumptions. | Pin pre-commit, require it, fix fixtures, expand matrix. | `.github/workflows/test.yml`, tests | CI 12/12 | Corrected. |
| V10-D | No tested non-root container. | Root-only runtime. | Optional UID/GID 10001 target; root default retained. | `Dockerfile`, CI | amd64/arm64 smoke | Added opt-in. |
| V11-A / V12-13 | Findings line started at previous newline; CRLF baseline mismatch. | Line-slice/location offset; raw CRLF content. | V11 line-start correction; V12 CRLF/entropy/legacy baseline canonicalization. | `detect/location.go`, `detect/detect.go`, `detect/baseline.go` | Location/CRLF/baseline tests | Corrected. |
| V11-B | Boolean/null allowlist regex not fully anchored. | Alternation not grouped. | `(?i)^(?:true\|false\|null)$` and regenerated config. | `cmd/generate/config/base`, `config/gitleaks.toml` | Placeholder positive/negative tests | Corrected. |
| V11-C | Global commit allowlist could be bypassed. | Loop/continue placement. | Central `isCommitAllowed`, skip entire diff file/commit. | `sources/git.go`, `detect/detect.go` | Multi-commit/staged tests | Corrected. |
| V12-04 | Patch-suppressing/combined Git formats returned silent clean. | Parser saw no supported sections but no error. | Reject unsupported log formats; retain V11 drain/section checks. | `sources/git.go` | Unsupported-format matrix | Corrected. |
| V12-06 | Cross-line/empty-capture redaction leaks. | Exact line-span grouping; empty Secret skipped. | File/commit grouping; include non-path empty-capture Match. | `detect/utils.go` | Multiline/empty tests | Corrected. |
| V12-07 | CSV formula injection. | Raw string cells. | Rune-aware apostrophe neutralization including Unicode/control/fullwidth. | `report/csv.go` | Trigger/idempotence/round-trip tests | Corrected. |
| V12-08 | SARIF byte columns wrong for Unicode/non-first lines. | Historical byte columns reused. | SARIF-only codepoint columns from line start. | `detect/location.go`, `report/finding.go`, `report/sarif.go` | Multi-script/later-line/CRLF tests | Corrected without JSON/CSV break. |
| V12-09 | SARIF artifact URIs unescaped. | Raw path. | URL escaping; drive/UNC handling. | `report/sarif.go` | Special-character/Windows/UNC tests | Corrected. |
| V12-10 | Paired overwrite rollback destroyed concurrent companion or deleted backup. | Unconditional remove/restore; deferred cleanup. | Identity/snapshot check; preserve concurrent output; retain failed-restore backup. | `cmd/report_pair.go` | Concurrent replacement/failed restore/regression tests | Handled-failure correction; crash window residual. |
| V12-11 | Regex matches computed twice. | Precheck and loop both called `FindAllStringIndex`. | Reuse first match slice. | `detect/detect.go` | Benchmark and full tests | Corrected. |
| V12-12 | Recursive decode resource amplification. | Whole-fragment decode passes; no segment cap. | No correctness change; finite measurements documented; default retained. | `detect/codec`, `detect/detect.go` | 6.67 MB control and nested-token probe | Not reproduced as V12 correctness defect; resource residual. |
| V12-14 | Long lines split at chunk boundary. | Blank-line lookahead ended inside long line. | 1 MiB finish-line pass; over-cap operational error. | `sources/common.go`, `sources/file.go` | 100–300 KB + over-cap tests | Corrected within cap. |
| V12-15 | Pre-commit `--all-files` false pass. | Filenames not forwarded to a safe scanner command. | New `files` command and explicit-file hooks. | `cmd/files.go`, `sources/explicit_files.go`, hooks | Real pre-commit and scope tests | Corrected; sequential performance limit. |
| V13-01 | No unused-ignore visibility. | Loaded entries had no usage state; incomplete scans could falsely imply obsolescence. | Track suppressing entries; opt-in completed-scan sorted diagnostic; no JSON schema change. | `detect/detect.go`, `cmd/root.go` | Detector/CLI tests and incomplete guard | Corrected, opt-in. |
| V13-02 | Effective configuration opaque. | No safe post-precedence summary. | `gitleaks config` safe summary with source/files/counts/policy/limits; no regex export. | `cmd/config_summary.go` | Summary/precedence/invalid tests | Corrected, opt-in command. |
| V13-03 | Multiple report formats one scan. | Single reportOutput model. | Architecture design for a validated report plan; no unsafe implementation forced. | `_assessment/v13-reporting-options.md` | Design review | Deferred — outside safe V13 scope. |
| V13-04 | Bounded finding context. | Context would change schemas and redaction scope. | Opt-in context design with byte/line budgets and redaction requirements; not enabled. | `_assessment/v13-reporting-options.md` | Design review | Deferred — outside safe V13 scope. |
| V13-05 | SARIF result severity omitted. | No per-rule severity model. | Opt-in constant `--sarif-level`; empty preserves output; no invented classification. | `report/sarif.go`, `cmd/root.go` | SARIF/CLI tests | Corrected, opt-in. |
| V13-06 | Repeated-finding grouping. | Unkeyed secret hashes are enumerable; primary dedup would alter findings. | Keyed-HMAC policy design; unkeyed SHA rejected; no grouping enabled. | `_assessment/v13-reporting-options.md` | Design review | Policy decision required. |
| V13-07 | Local path UX / decoded highlighting. | Absolute paths leak; decoded match has no source coordinate. | Opt-in local `file://` links with no-color handling; decoded highlighting remains separate. | `detect/utils.go`, `cmd/root.go` | File-link/JSON isolation tests | Partially corrected — decoded highlight residual. |
| V13-08 | `--log-opts` split on spaces. | No shell-aware tokenization; validation on raw strings. | Shell-free tokenizer; parsed-vector validation; quoting cannot bypass fail-closed formats. | `sources/git.go` | Parser/vector tests | Corrected. |
| V13-09 | Shallow history provenance. | Missing history is unavoidable; no warning. | Detect shallow repository and warn; never rewrite SHAs/fingerprints. | `sources/git.go` | Shallow fixture test | Partially corrected — warning only. |
| V13-10 | Subdirectory/bare repository scope. | Git history is repo-wide; bare trees lack worktree ignore discovery. | Explicit warnings; no default scope change; explicit config/ignore paths documented. | `sources/git.go` | Bare/subdir tests | Partially corrected — explicit behavior, no new scope feature. |
| V13-11 | Azure SSH remote links. | SSH v3 host/path not normalized; ad-hoc escaping. | Canonical Azure URL mapping; URL path escaping; credential stripping. | `sources/git.go`, `detect/utils.go` | Remote/link tests | Corrected. |
| V13-12 | Configuration diagnostics. | Empty rules/malformed TOML mostly fail, unknown fields remain ignored. | Verified fail-closed empty/malformed cases; safe summary provides context; strict unknown-field mode not forced. | `config`, `cmd/config_summary.go` | Empty/malformed tests | Partially corrected — unknown-field/allowlist-placement residual. |
| V13-13 | Pre-commit environment reliability. | Docker/macOS/PTY reports lack a common proven cause. | Local real pre-commit passes; no unsafe safe.directory/scope change. | `.pre-commit-hooks.yaml`, tests | Local integration | Not reproduced — evidence insufficient for Docker/macOS. |
| V13-14 | Profiling script measurement validity. | `set -e`, decode-depth mismatch, unquoted paths, early findings. | Same command for benchmark/profile, explicit decode depth, safe quoting, exit-code 0, metadata, artifact checks. | `scripts/profile.sh` | Linux smoke run with dir/git + all profiles | Corrected. |
| V13-15 | Decode/hang stability. | V12 evidence did not establish a correctness defect. | Finite 6.67 MB line-bounded probe at depths 0/1/3/5; no depth-correlated peak-RSS growth in that control; hangs not reproduced. | measurement evidence | Linux `/usr/bin/time` probe | Not reproduced — evidence insufficient; defaults retained. |
| V13-16 | Explicit-file scalability. | Sequential file loop. | Bounded semaphore concurrency, deterministic findings sort, per-file errors and dedupe retained. | `sources/explicit_files.go`, `cmd/files.go` | Concurrency/dedupe/error tests; 200-file probe | Corrected; measured ~2.2x on synthetic Windows workload. |
| V13-17 | Upstream release identity. | Docker Hub/upstream image publication and no protected ancestry/tag validation. | GHCR-only RNC workflow, tag/base/ancestry checks, least privilege, no Docker Hub credentials. | `.github/workflows/release.yml`, hooks, `RELEASING.md` | Static workflow/tag probes | Corrected in source; no publication. |
| V13-18 | Version/install contract. | Module path vs RNC distribution ambiguity; prerelease SemVer trap. | `vX.Y.Z-rnc.N`, base strictly > minVersion, container tags, installation matrix; module path retained. | `RELEASING.md`, README | Tag-policy probe | Policy decision documented; module rename not performed. |
| V13-19 | Release binary portability. | GoReleaser still targets gore2regex and legacy architectures; no artifact test. | Assessed portability/RE2/architecture policy; no release config rewrite without artifact-level proof. | `RELEASING.md` | Documentation review; goreleaser unavailable locally | Partially corrected — packaging plan only. |
| V13-20 | Attestations/release proof. | No SBOM/provenance/signature/verification design. | Exact digest/SBOM/provenance/checksum/signature verification targets and release checklist; not enabled. | `RELEASING.md` | Design review | Policy/deferred — no attestation claimed. |

## 2.2 Security hardening

### 2.2.1 Fail-closed execution

The dominant security correction is the separation of “clean scan” from “scan that failed to execute completely.” Source walkers, Git acquisition, parser completion, workers, ignore parsing, report publication and argument parsing all feed operational error paths. Completed scans retain normal finding/exit behavior; incomplete scans cannot be converted to success by `--exit-code 0`.

### 2.2.2 Git trust boundary

Repository configuration is attacker-controlled when scanning untrusted content. Pinning presentation flags prevents `format.pretty`, `log.date` and color from changing parsed metadata, and `--no-textconv`/`--no-ext-diff` prevent repository-selected helpers from executing. The parser also rejects Git output modes it cannot interpret completely. Git processes are drained, waited and reaped; stderr is classified and sanitized.

### 2.2.3 Input and filesystem boundaries

- The directory command rejects ambiguous targets instead of scanning the wrong tree.
- Report destinations are validated without modification; symlinks, directories, special files and scan-input collisions are rejected.
- `.gitleaksignore` must be a regular file; FIFOs/special inputs are rejected before content is consumed.
- Archive members and decompressed streams are bounded by exact byte limits and temporary spooling.
- Long lines above the finish budget fail closed.
- Explicit-file mode never expands a directory.

### 2.2.4 Diagnostic and report confidentiality

Finding redaction covers known secrets across fields and auxiliary required-rule records, including Unicode, empty captures, adjacent findings and cross-line ranges. Diagnostic paths avoid raw secret-bearing fields; Git stderr is sanitized heuristically. Reports remain raw by default when `--redact=0`; this is a compatibility contract. CSV formula neutralization prevents report cells from becoming spreadsheet formulas. Report files use restrictive temporary permissions on POSIX.

### 2.2.5 Container and supply chain

Default container runtime remains root for compatibility and mount access, with `safe.directory=*` retained. The explicit `runtime-nonroot` stage provides UID/GID 10001 and a controlled home directory for operators who can provide matching mount permissions. CI pins Go, action SHAs, test reporter, vulnerability scanner, pre-commit and Trivy versions; workflow permissions are read-only except the release image workflow’s package write. Dependencies are module-verified and vulnerability-scanned in CI.

### 2.2.6 Residual security risks

Do not claim universal secret-safe logging, complete redaction of unknown text, atomic two-file publication, perfect `.gitleaksignore` TOCTOU protection, comprehensive archive handling, or release-artifact certification. Section 2.8 records each residual.

## 2.3 Reporting architecture

### 2.3.1 Report lifecycle

1. Resolve report mode from CLI/config.
2. Validate destination, parent, collision and input overlap without modifying the destination.
3. Scan and collect findings.
4. Render the primary report to a restrictive temporary file.
5. Render the optional informative companion to a separate temporary file.
6. Reserve or select the final primary/companion names according to mode.
7. Revalidate destinations immediately before publication.
8. Publish companion, then primary; on handled primary failure, roll back safely.
9. Remove unused temporary files; report actual primary path.

A failed scan or render never publishes a partial report and does not replace existing output.

### 2.3.2 Output modes

| Mode | Trigger | Primary naming | Existing companion | Failure behavior |
|---|---|---|---|---|
| Overwrite | `auto-overwrite=true` or `--overwrite` | Requested path | Backed up before publication; restored only if still safe | Primary failure leaves old primary; companion rollback identity-checked |
| Aggregate | `aggregation=true` or `--aggregate` | First free numbered path | Same numbered pair | Complete pair reservation; new pair removed on handled primary failure |
| Reject | Both policies false, no override | Requested path must be absent | Both must be absent | Rejected before scanning |
| Stdout | `--report-path -` | stdout | No file companion; informative goes to stderr | Normal stdout exit behavior |

### 2.3.3 Informative report

The informative document is schema version 1 and additive within the major version. It contains source type, optional generation time, observed/inspected target counts, exact exclusion totals by reason, and bounded per-target detail (10,000 records; console detail 20). Reasons cover size, binary content, archive depth, path/commit allowlists, symlink policy, empty files and excluded inputs. Counts describe encountered targets, not unvisited subtree contents. Findings-level suppression is not misrepresented as a skipped file.

### 2.3.4 Paired publication and V12-10 correction

The original companion implementation published the companion before the primary and, on primary failure, removed the published companion and restored a backup. That was insufficient under concurrency:

- an unrelated writer could replace the published companion after its identity was captured;
- unconditional restoration could then overwrite that writer’s file;
- a failed restoration could be followed by deferred backup deletion, destroying the only recoverable prior companion.

The corrected algorithm:

1. Back up an existing companion before publication (hardlink preferred, copy fallback).
2. Capture the published companion identity plus size/modtime after its rename.
3. On primary failure, inspect the current path.
4. If it is still the file this run published, remove it and restore the prior companion.
5. If it is absent, restore the prior companion.
6. If it has been replaced or cannot be inspected, do not touch it; retain the backup and report its path.
7. If removal or restoration fails, retain the backup and report the failure; do not let deferred cleanup delete it.
8. If no prior companion existed, remove only this run’s companion; a concurrent replacement is preserved.

This closes the reproduced capture/unconditional-restore and failed-cleanup windows. A simultaneous replacement between the final identity check and removal remains a narrow filesystem TOCTOU residual. It also cannot make two independent files crash-atomic; a process death between companion and primary publication can still leave a mismatched pair. Both limitations remain documented.

### 2.3.5 SARIF severity policy

SARIF result-level severity is deliberately opt-in. `--sarif-level` accepts `none`, `note`, `warning` or `error`, applies one constant level to all results, and is rejected for non-SARIF output. Empty remains historical omission. No per-rule risk classification is invented; consumers that require a level can opt in without changing rule IDs, locations or fingerprints.

### 2.3.6 Formats and compatibility

- JSON: historical field names/byte columns; `Finding.Fragment` remains excluded from JSON.
- CSV: RFC 4180 writer; formula-like string cells apostrophe-neutralized; historical schema retained.
- JUnit/template: existing contracts; template fields remain caller-controlled.
- SARIF: codepoint columns and escaped artifact URIs; ASCII relative paths produce the same ordinary URIs as before.
- Console: redacted only when requested; informative output goes to stderr when appropriate.
- Timestamp: informative-only; primary findings reports do not gain a timestamp.

## 2.4 Detection-engine correctness without rule expansion

No provider detection rules were added. Engine corrections affect how existing rules see content and policy:

- **Allowlists:** AND path sufficiency, anchored boolean/null placeholders, global commit allowlist semantics, root-relative path matching.
- **Line boundaries:** `Finding.Line` begins at the matched line; anchored line allowlists work after the first line; empty/newline decoded ranges are clamped.
- **Unicode:** rune-safe masking; SARIF codepoint columns; UTF-8-preserving redaction.
- **CRLF:** detector match spans and entropy canonicalize line endings; legacy baselines migrate columns/entropy.
- **Baselines:** content comparison tolerates line-ending conversion without weakening path/fingerprint/line identity.
- **Fragmentation:** 25 KB lookahead remains; 1 MiB finish-line pass fixes 100–300 KB boundary splits; over-cap lines fail closed.
- **Required rules:** auxiliary findings are redacted consistently with primary findings.
- **Decoding:** recursive depth remains bounded by the configured maximum; encoded content can still amplify resources.
- **Regex evaluation:** the match slice is computed once and reused per rule/pass.
- **Source completeness:** MIME first-chunk, archive content verification, staged index fallback, Git parser counts and operational errors.
- **Merge semantics:** default retained; first-parent is explicit; unsupported formats rejected.

These are correctness changes, not an expansion of detection signatures. Accuracy against a controlled oracle remains the separate scanner-harness responsibility.

## 2.5 Performance and resource management

| Workload / change | Before | After | Conditions and interpretation |
|---|---|---|---|
| Duplicate regex scan, positive-heavy rule | ~3.3–3.9 µs/op; ~1,917–1,920 B/op; 14 allocs | ~3.0–3.3 µs/op; ~1,894–1,896 B/op; 11 allocs | Windows 11, i7-11375H, Go 1.27.3, `go test -bench`; synthetic single-match rule. |
| 8 MiB positive regex probe | ~2.8–2.9 ms/op | ~0.8–1.0 ms/op | Same host; large-input probe, not a general throughput claim. |
| Staged binary-marked blob drain | ~286 ms pre-V4; ~530 ms V4; ~305 ms V5 | ~305 ms | Agent-reported 100 MB fixture; not an independently controlled release benchmark. |
| Long-line fragmentation | Split possible at 125 KB+; missed boundary-straddling secrets | Complete through 100–300 KB; over-cap fails closed | Memory bounded by a 1 MiB extra-read budget per file; larger lines are an operational error. |
| Recursive decode control | 6.67 MB encoded payload | 13.3 MB total allocation, ~20.7 MB heap at depths 0/1/3/5 | Empty-rule control showed no depth-correlated leak; nested-token probe showed linear allocation. Not a correctness fix. |
| Explicit-file hook | Staged hook false-passes `--all-files` | Exact-file scan passes/fails correctly | Sequential; large `--all-files` sets slower than a tuned directory walk. |
| Report publication | Direct write | Temp render + rename + validation | Adds one temp write/rename; required for data safety. |
| Git acquisition | Single parser pass | Parser pass + raw section tracker/drain | Small CPU/memory overhead for completeness guarantees. |
| Explicit-file scan, 200 x 50 KB | ~0.577 s median sequential V12 | ~0.267 s median bounded-concurrency V13 | Windows 11, local disk, 4 timed runs after one warm-up; synthetic no-finding workload. Not a general throughput claim. |

Performance evidence is platform- and workload-specific. No universal production speedup, memory bound, or detection-accuracy claim is made from these probes.

## 2.6 Build, dependency and platform engineering

### 2.6.1 Toolchain and modules

`go.mod` declares Go `1.26.0`; CI and the container builder pin Go `1.27.4` for vulnerability and platform coverage. Local validation in this cycle used Go `1.27.3` on Windows. The independently resolved module graph includes `github.com/wasilibs/go-re2 v1.13.0`, Lip Gloss v2, `go-gitdiff v0.9.1`, Viper, Cobra and `mholt/archives`. Both standard and `gore2regex` builds/tests pass. Full repository generation uses the project regexp abstraction.

### 2.6.2 Build contract

The Makefile distinguishes commands that do not mutate tracked source, configuration or module files (`build`, `test`, `test-race`, `test-cover`, `failfast`, `verify`, `lint`, `version`, `profile`) from explicit maintenance (`format`, `tidy`, `generate`, `clean`). Coverage and profile commands still write their requested local artifacts; they do not rewrite the source tree. Version resolution never fetches remote tags; it prefers explicit `VERSION`, then a local exact tag, nearest tag plus commit metadata, or a valid source fallback. `make verify` checks vet, generated configuration and module consistency where supported.

### 2.6.3 Container and CI

- Default image: Alpine 3.24 runtime, UID 0/root, Git/SSH/bash, `safe.directory=*`, version-stamped build.
- Optional `runtime-nonroot`: UID/GID 10001, `HOME=/home/gitleaks`, restrictive home permissions.
- Build stage cross-compiles using BuildKit target args; multiarch images do not run the compiler under QEMU.
- CI triggers only PRs to `rnc-mod` and pushes to `rnc-mod`; content permission is read-only.
- Twelve jobs: Linux race + real pre-commit + generator, Windows race, two macOS architectures, stdlib/RE2 builds, four Docker build/smoke jobs (amd64/arm64 root/non-root), dependency/vulnerability verification, container scan.
- The existing Gitleaks Action scanner workflow is deliberately not part of RNC CI.

### 2.6.4 RNC release and installation controls

The proposed RNC release identity is GHCR-only: `ghcr.io/rnc-solutions-group/gitleaks`. The release workflow validates `vMAJOR.MINOR.PATCH-rnc.N` tags, requires the numeric base to be strictly newer than the embedded `minVersion`, verifies the release commit is reachable from protected `rnc-mod`, and requests only `contents: read`/`packages: write`. It does not use Docker Hub credentials or publish an upstream image. No V13 action publishes a release.

The repository keeps the upstream Go module path for compatibility; `go install` of that path obtains upstream Gitleaks, not RNC-Mod. Source builds and future RNC release artifacts are the supported RNC paths. `RELEASING.md` records the full installation matrix, container-tag aliases, artifact portability constraints and proposed checksum/SBOM/provenance/signature verification. An attestation is not claimed merely because the workflow contains a build option.

### 2.6.5 Platform support

Code and tests cover Windows, Linux and macOS. Windows-specific paths include rename retry, cross-volume fixtures, regular-file checks and permission-test limitations. POSIX-specific paths include nonblocking ignore-file open and mode-bit assertions. Cross-architecture Docker builds cover amd64/arm64. Native race execution is performed in CI; local Windows race requires a C toolchain and was not available on this host.

### 2.6.6 Validation versus release

Source validation, local builds/tests, GitHub Actions, image smoke tests and a published release artifact are distinct. The final V12 candidate has local validation and PR CI; no published-release attestation, signature verification, SBOM attestation or independent scanner-harness certification is claimed by this document.

## 2.7 Validation and regression evidence

| Capability | Validation method | Coverage | Result | Limitations |
|---|---|---|---|---|
| Standard build/tests | `go build ./...`, `go test -count=1 ./...` | All packages | PASS locally and in CI | Does not prove detection accuracy. |
| RE2 build/tests | `go test -count=1 -tags gore2regex ./...` | All packages plus generator graph | PASS locally and CI build jobs | RE2 constrained-memory release issue remains. |
| Static/module checks | `go vet`, `go mod verify`, `go generate` clean diff | Whole repo | PASS | Windows CRLF `mod tidy` check deferred to Linux. |
| Race | `go test -race ./...` non-root Linux; Windows/macOS CI race jobs | All test-bearing packages | PASS | Local Windows race unavailable without C toolchain. |
| Git integrity | Synthetic evil/normal merges, malformed/truncated/blocked streams, hostile format/date/textconv, staged binary | `sources`, `detect`, CLI | PASS | Pathological custom formats and blocking reads residual. |
| Report publication | Temp/rename, permissions, aggregate/reject, hardlink/symlink, concurrent writer, pair rollback | `cmd/report*`, CLI | PASS | Crash atomicity not testable as transactional. |
| Config and ignore safety | Invalid regex, extension depth, precedence, missing local config, FIFO/oversize/partial ignore | `config`, `detect`, CLI | PASS | TOCTOU/extreme legitimate files residual. |
| Archive/size | Signature verification, prefix replay, member/decompressor limits, spool, nonseekable cases | `sources` | PASS | Outer-size/self-extracting limitations accepted. |
| Unicode/locations | ASCII/Arabic/CJK/emoji/combining, later lines, CRLF, multiline | `detect`, `report` | PASS | Rule-level path regex remains reported-path-only. |
| Redaction | Empty, Unicode, adjacent, multiline, auxiliary, verbose/template | `report`, `detect`, CLI | PASS | Unknown text/partial masks residual. |
| CSV/SARIF | Injection/idempotence, URI escaping, schema compatibility | `report` | PASS | Spreadsheet save/reopen durability not exhaustively tested. |
| Pre-commit | Manifest, Python helper, real pre-commit 4.7.0 `--files`/`--all-files` | Local and Linux CI | PASS | Explicit-file mode sequential. |
| Docker | amd64/arm64 root/non-root build and smoke | CI | PASS on PR head (pre-correction) | Final correction CI recorded in PR. |
| V10 equivalence | Object-level tree equality with accepted V10 authority; binary/report comparison | Product tree, fixtures | PASS | Applies to V10; V11/V12 intentionally differ. |
| Detection harness | Independent oracle/accuracy evaluation | Separate harness | NOT RUN in this cycle | Required release/accuracy gate remains external. |
| Published release attestation | Signed artifact/SBOM verification | Release process | NOT RUN | Source/CI validation is not release certification. |
| V13 config/ignore/file-link/SARIF diagnostics | CLI and unit tests | New opt-in controls and JSON isolation | PASS | Docker/macOS consumer GUI not exercised. |
| V13 Git argument parser | Unit vector tests plus existing Git suites | Quotes, escapes, Windows paths, pathspec boundary and disallowed-format bypass | PASS | No exhaustive Git-version matrix. |
| V13 explicit-file concurrency | Source/CLI tests and 200-file timing | Dedupe, scope, errors, deterministic order | PASS | Synthetic Windows workload; no huge descriptor-limit test. |
| V13 profile script | Linux dirty/Git smoke run | dir/git CPU, memory, trace artifacts | PASS | hyperfine/yq optional tools absent. |
| V13 release workflow | YAML parse, static scope test, tag-policy probe | GHCR-only, ancestry/tag/base validation | PASS in source | Workflow not executed with registry credentials; no image/chart/SBOM published. |

## 2.8 Remaining engineering limitations

| Residual / policy decision | Status | Impact | Mitigation | Reconsider when |
|---|---|---|---|---|
| Default merge-only Git coverage | Open policy decision | Merge-only resolutions can be missed in default scans | Use `--log-opts=--diff-merges=first-parent`; pair with full-history scans | Operator accepts duplicate/fingerprint changes and approves default change. |
| Source-local config trust | Open policy decision | A repository can replace rules/suppress findings | Use `--no-local-config` or explicit trusted config for untrusted scans | Secure-by-default adoption is approved with migration plan. |
| Two-file crash atomicity | Accepted residual | Process death can leave mismatched primary/companion | Handled-failure rollback; operator can re-run/reconcile | A journal/transactional publication design is authorized. |
| Paired rollback concurrent writer | Corrected for handled failures | A detected replacement is preserved; prior companion backup may remain | Identity/size/modtime checks; error reports retained backup path | Crash-safe recovery workflow is implemented. |
| Companion rollback check-to-remove window | Narrow TOCTOU residual | A writer replacing the companion between the final check and removal could still be displaced | Narrow window; identity/snapshot checks; publication failure is reported nonzero | Filesystem transaction or handle-relative removal is designed. |
| Report final-publication TOCTOU / parent-rename orphan temp | Residual | A late filesystem race can leave an owner-restricted temp file or fail publication | Pre-publication revalidation, identity checks, temp cleanup | Stronger directory-lease/atomic protocol is justified. |
| Windows rename locks beyond retry | Residual | Persistent sharing violation can fail publication | Bounded retry; nonzero failure; old report preserved | Longer/backoff policy or platform-specific redesign is approved. |
| Long lines above 1 MiB finish budget | Accepted fail-closed limit | Very long minified/pathological lines cannot be completely scanned | Operational error rather than false clean | Larger bounded reader or streaming engine is designed. |
| Recursive decode resources | Unresolved resource behavior | Encoded input can amplify memory/CPU with depth 5 | Bound input size/depth; monitor; no correctness claim | Controlled benchmark/segment cap is approved. |
| Archive nonseekable outer-size cap | Accepted intentional exclusion | Small members in large nonseekable ZIP/7z may be skipped | Informative `size_limit` reporting | Archive engine expansion is separately authorized. |
| Self-extracting/prefixed archive detection | Residual | May be treated as text | Scan text or use supported archive format | Signature discovery is expanded. |
| Rule-level `path` regex semantics | Residual | Rule `path` sees reported spelling, not root-relative | Use global/root-aware allowlist paths | Rule path normalization is separately designed. |
| Ignore/config TOCTOU | Residual | File replacement between checks and read | Regular-file checks, nonblocking open, complete parse | Directory-handle/descriptor-relative design is adopted. |
| Blocking OS reads | Residual | Timeout/cancellation may not interrupt promptly | Timeout plus process/worker lifecycle | Context-aware readers are available. |
| Deprecated library APIs | Residual | Some non-CLI callers may retain older error behavior | CLI uses corrected paths | Public API cleanup is authorized. |
| Global Viper state | Residual | Concurrent config loads in one process could interfere | CLI is single-config per run | Config API becomes instance-scoped. |
| Diagnostic redaction heuristics | Residual | Unknown/short secret formats or over-redaction possible | `--redact`; avoid trace logs on sensitive content | Structured diagnostic boundary is designed. |
| Redaction unknown adjacent text | Residual | Text not itself reported can remain | Known-exposure grouping; controlled output | Full document taint tracking is designed. |
| CSV apostrophe durability | Residual | Spreadsheet may strip mitigation on save/reopen; data changes for formula-like cells | Apostrophe mitigation; consume in parser/trusted tooling | A safer encoding contract is standardized. |
| Explicit-file performance | Accepted limitation | Large `--all-files` sets slower | Use directory scans where scope permits | Bounded parallel explicit-file reader is implemented. |
| Release RE2 address space | External release issue | Constrained runtimes may fail RE2 initialization | Use standard backend or adequate memory | Upstream/runtime fix is available. |
| Go minimum/toolchain policy | Open release policy | `go.mod` 1.26.0 vs CI 1.27.4 | CI pins current patched toolchain | Operator sets support floor. |
| Multi-format report plan | Deferred | One scan still produces one primary report (plus informative companion) | Use separate scans/report paths; no unsafe partial plan | Versioned report-plan contract with pre-scan validation is authorized. |
| Bounded finding context | Deferred | No surrounding context in machine/console reports | Use existing `Match`/Line console view; no raw context export | Opt-in context schema and redaction window are designed and tested. |
| Repeated-finding grouping | Policy decision | No secret-derived grouping identifier | Primary findings remain intact; use external grouping only with a scoped key | Keyed-HMAC threat model and key lifecycle are approved. |
| Decoded-finding highlighting | Residual | Decoded matches may render without styled source coordinates | Plain match/secret output remains | Provenance-aware decoded presentation is designed. |
| Config unknown fields / allowlist placement | Residual | Loader may ignore unknown TOML fields; strict rejection not forced | `gitleaks config` shows effective result; fix configuration on validation error | Strict validation mode with migration documentation is approved. |
| Shallow Git history | Warning only | Pre-boundary history is absent and cannot be reconstructed | Treat warning as incomplete provenance; unshallow for full history | A formal coverage manifest for shallow scans is authorized. |
| Subdirectory/bare Git scope | Warning/explicit configuration | History remains repository-wide; bare trees require explicit config/ignore | Use `--log-opts` pathspec or explicit paths | A bounded subdirectory-history feature is designed without changing defaults. |
| Pre-commit Docker/macOS/PTY environments | Not reproduced | Local pre-commit passes; host-specific failures remain unproven | Use documented mount/ownership/safe-directory guidance; no broad trust bypass | A reproducible Docker/macOS trace and root cause are supplied. |
| Recursive decoding/hang resources | Not reproduced in V13 probe | Encoded input can still amplify resources; historical futex report unproven | Bound inputs/depth; retain default depth 5 | Controlled harness finds a reproducible defect with stack/RSS evidence. |
| Release artifact/attestation | Not published | No RNC binary/image/tag/SBOM/signature exists yet | `RELEASING.md` verification checklist; no release claim | Operator authorizes a release and a separate consumer verifies all artifacts. |
| Published-release attestation | Not performed | No claim of signed release/SBOM certification | Build/test evidence only | Release pipeline attestation is executed. |
| Independent scanner harness | Not run | No independent accuracy/oracle certification | Existing tests + separate harness required | Harness baseline is published and executed. |

# 3. Attribution Register

## 3.1 Attribution methodology

The register distinguishes evidence classes:

- **Direct adoption:** an original contributor commit is retained and its implementation materially survives; source SHA, author identity and `-x` trailer are preserved.
- **Adapted implementation:** an original contributor commit is retained, but the final RNC tree modifies or supersedes part of its implementation. The commit remains ancestral; the adaptation is recorded, not hidden inside an RNC-only claim.
- **Independent implementation / design influence:** prior work informed the RNC design, but no contributor commit was cherry-picked and no code authorship is claimed.
- **Research or defect discovery:** an issue reporter or reviewer supplied a reproduction, failure mode or analysis used to guide a correction.
- **Documentation or dependency contribution:** prose, build metadata, install instructions or independently matched dependency work; no scanner commit is claimed.
- **Reviewed but not adopted:** a proposal was evaluated and excluded for a stated technical reason; recognition is documentary only.
- **Pre-existing baseline:** work already merged upstream before RNC's contributor integration; it is not counted as an RNC adoption.

Git ancestry and authorship claims below are derived from actual commit objects: full 40-hex source and integration SHAs, author name/email, and `(cherry picked from commit …)` trailers. No `Co-authored-by` trailer was fabricated. Existing upstream co-author trailers are retained as part of the original commits. Attribution is voluntary engineering policy and does not replace the upstream MIT license notice (`Copyright (c) 2019 Zachary Rice`).

## 3.2 Complete recognized contributor register

### 3.2.1 Commit-level contributors

All rows are retained ancestral commits. “Adapted” means the final tree reconciles or expands the original implementation; it does not erase authorship.

| Contributor | Upstream PR | Original commit SHA(s) | RNC integration SHA(s) | Contribution | Attribution type | Final disposition |
|---|---|---|---|---|---|---|
| BasitS-hash `<basitsherazi319@gmail.com>` | [#2227](https://github.com/gitleaks/gitleaks/pull/2227) | `894b4a3b04dfd7c10e9fb887cdaaad35e4b79c31`<br>`49b26eece1ea83ea58da0c6e0bf92326fb82a25c` | `3004bda1fc3833c68218f2b8589581fdc99715c9`<br>`c8d707016d51eca7d5e02d1f5066b48ecc668788` | AND-allowlist path sufficiency; nil/commit/RE2 review cases | Direct adoption | Ancestral; final tree retains path-sufficiency logic with root-aware extensions. |
| Shanmukh-l33t `<Shanmukh-l33t@users.noreply.github.com>` | [#2334](https://github.com/gitleaks/gitleaks/pull/2334) | `ad321b7d7aa9b1a1642bb56813fa7443d5412971` | `bdc880a71479a3e3c4e14a86e02858b1831752cb` | Explicit `dir` target contract and wrong-directory regression | Direct adoption | Ancestral; one-target validation and empty-argument correction retained. |
| Arpit Jain `<arpitjain099@gmail.com>` | [#2261](https://github.com/gitleaks/gitleaks/pull/2261), [#2253](https://github.com/gitleaks/gitleaks/pull/2253) | `b1ce07f27d8251ef22e5b47c80387e9b1590fc2f`<br>`4f58f98a1424203af39cc00a26f9672e9b45b794` | `1f3533b1a1e8498ee5e99aba54ee557dda9c0367`<br>`260d7e18947668fa47b660f8d9fa242c62191454` | Rune-based masking/UTF-8 tests; fragment error preservation | Direct adoption (#2261); adapted (#2253) | Ancestral; redaction expanded and fragment errors moved to bounded collector. |
| Smeet Agrawal `<smeetagrawal23@gmail.com>` | [#2082](https://github.com/gitleaks/gitleaks/pull/2082), [#2080](https://github.com/gitleaks/gitleaks/pull/2080) | `df8223e0008b6d5bf82d0433592b8f1557c97311`<br>`ff77f9f58b1ec5fbbece852e0bfe9c1f050f24e5` | `7a4c18ff97b9771695485672b409582a37b20aac`<br>`2bd5ad38561fd54ff92975f2c38344d56f935fc4` | Error-returning regex wrappers; Git stderr propagation | Direct adoption (#2082); adapted (#2080) | Ancestral; validation expanded; stderr fatal/warning classification and sanitization supersede. |
| Jeff `<138259519+jdotc1@users.noreply.github.com>` | [#2307](https://github.com/gitleaks/gitleaks/pull/2307), issue [#2304](https://github.com/gitleaks/gitleaks/issues/2304) | `0d1ae3ea1087bebcf5ce45982a770d0f4925af6c` | `d68ecb2c2cc6978673351f98b971c77827922204` | Removed destructive report preflight; report preservation | Adapted implementation | Ancestral; replaced by temp-render/final-publication architecture. |
| cardene777 `<cardene777@gmail.com>` | [#2329](https://github.com/gitleaks/gitleaks/pull/2329), [#2324](https://github.com/gitleaks/gitleaks/pull/2324) | `21028d8858e31fd088029f2fdb6634c99a1d3563`<br>`a77d865f9eebbe38903c5fc365d1f0accfbff8d8` | `2cf0ade305f4a60049b0c57af466dee4744e04ad`<br>`d1fe7bca8f34608358f4de2f7cb0f416fd979e57` | First-chunk MIME classification; finish current line at chunk boundary | Direct adoption (#2329); adapted (#2324) | Ancestral; RNC added long-line budget and fail-closed over-cap behavior. |
| Sakariyah Abdulhazeem `<sakariyahabdulhazeem@gmail.com>` | [#2284](https://github.com/gitleaks/gitleaks/pull/2284) | `1894e6f2cdaef94fe66388e577e0f59f14bfcff2` | `f146e235258cf21d695f74ba1791c2033859f721` | Forced `--no-color` for Git parsing | Adapted implementation | Ancestral; expanded into full deterministic output flag set. |
| INTFRAME `<hello@intframe.com>` | [#2246](https://github.com/gitleaks/gitleaks/pull/2246) | `9b7b2243fcdeb4279cccd9df9348f9e537d4fb95` | `5bed9587521b429aecb3b223acf0d20dbeef754f` | Git parser drain, unparsed-byte detection, process lifecycle | Adapted implementation | Ancestral; expanded with raw-section tracking, stderr classification and collector. |
| KVR Vishnu `<46535652+VISHNU0906@users.noreply.github.com>` | [#2338](https://github.com/gitleaks/gitleaks/pull/2338) | `7fed4140d87f62090c87b9f09c5c9b6073a84a27` | `16ba49e717e24e5bea47dc7a2419033ebbe6bcde` | Staged binary-marked text fallback | Adapted implementation | Ancestral; explicit stage-zero index/MIME/error path supersedes fragile lookup. |
| DelaySystem `<DelaySystem@tubot.lat>` | [#2339](https://github.com/gitleaks/gitleaks/pull/2339), [#2346](https://github.com/gitleaks/gitleaks/pull/2346) | `7e5132faa1ccea4bfd0b4b006ca896b5dd03df3f`<br>`126bab3a4497026961f7e78d373b098fb845dc92`<br>`5379957823eab9f6d4ed31827cc044ba00797e35`<br>`760252fd84604b2a6179b060241237548fca6804` | `b603fa0bd8a563909c062d341f2f1ea5781e2142`<br>`f6c111c8e701b2129381535b638bf130a0de9002`<br>`05d58c560e84b29552ae4810f0a930e05e6266c5`<br>`29a6f766c051a5cbd60aba2d96f3e1410b8a8e22` | Git nonzero-exit propagation/reaping; line-start/newline boundary and empty-range clamping | Adapted implementation | Ancestral; merged into collector; V12 CRLF/Unicode location handling extends the corrections. |
| Nikhil `<tad.areas_0y@icloud.com>` | [#2099](https://github.com/gitleaks/gitleaks/pull/2099) | `1b3896fd7df9a667a09568b751dce6143613319e` | `62a89d4339a344e52e209a7c6d260de711b34518` | File/archive/member size-limit propagation | Adapted implementation | Ancestral; whole-buffer decompression replaced by bounded spooling. |
| Varshith Puli `<pulivarshit@gmail.com>` | [#2317](https://github.com/gitleaks/gitleaks/pull/2317), [#2312](https://github.com/gitleaks/gitleaks/pull/2312), alternative [#2313](https://github.com/gitleaks/gitleaks/pull/2313) | `71e91200afb2415b283bb1721709b19c2fddd8aa`<br>`2009273e0029f91aeda4732552ae5489634d0749` | `7ce04bb344f4ddb7882167f9dcfce73598900c77`<br>`6445d1e949d1ff49caa7a6db86340d66ea4c8973` | Archive content verification/prefix replay; restrictive report permissions | Direct adoption (#2317); adapted (#2312) | Ancestral; error propagation expanded; preflight hunk dropped; #2313 documentary only. |
| hamodywe `<iosapk.org@gmail.com>` | [#2248](https://github.com/gitleaks/gitleaks/pull/2248), design influence from [#2238](https://github.com/gitleaks/gitleaks/pull/2238) | `5c9ca86f8f4d1db1e17004b6a1b598137203face` | `8ee08eff98e602ea2a1e301bd2e6ab575b0da529` | Root-relative ignore/allowlist/baseline design | Adapted implementation | Ancestral; V8/V12 root-aware and CRLF corrections supersede helpers; codepoint design credited, not cherry-picked. |
| Vaibhav Srivastava `<vaibhavsri1712@gmail.com>` | [#2249](https://github.com/gitleaks/gitleaks/pull/2249) | `509fce15ad2eb2536bc8d840a3ed3e374d73bb2d` | `5e353588a898b56abae78018c4f13c10f59122f0` | Unreadable-file error collection | Adapted implementation | Ancestral; bounded `scanerr` collector supersedes local slices. |
| ashvinctrl `<sharmaashvin27@gmail.com>` | [#2166](https://github.com/gitleaks/gitleaks/pull/2166) | `55099ffe9d039c6fb0500a5f65cca1abfb398826`<br>`fcd84a613a099e698aa811b32e6d3494c6af0326` | `c13eb34551fb4c2fe2917a0a866dd222e0ac1ab8`<br>`9c60717293ad113f09e24dbd29e7fcee7af18333` | Global commit-allowlist outer-loop correction; helper/tests | Adapted implementation | Ancestral; centralized helper retained in final tree. |
| Paradox `<159635503+i-am-paradox@users.noreply.github.com>` | [#2294](https://github.com/gitleaks/gitleaks/pull/2294) | `b04e70b06a62e95ae09b68eafeb9b7a2050153e8` | `079ee1cb39aa50bb39f82c27eb7f2f7f8d0fbddc` | Anchored boolean/null placeholder allowlist and regenerated config | Adapted implementation | Ancestral; exact `^(?:true\|false\|null)$` behavior retained. |
| Anthony Metzidis `<anthony.metzidis@gmail.com>` | [#2311](https://github.com/gitleaks/gitleaks/pull/2311) | `ae702afb4b1cf41d560baf89a9f454eadfa251f2` | `2790da33f1250426c8e8901ea0cde6476410eee3` | Standard Git pretty format | Adapted implementation | Ancestral; extended with date/textconv/ext-diff pinning. |
| Bunlong Heng `<bheng.code@gmail.com>` | [#2237](https://github.com/gitleaks/gitleaks/pull/2237) | `981e18e2ba6dc3d5419047d026352fb4a85ab5b1` | `904871b791ab40fd92fc59917d7640123e1448d5` | CSV formula-injection mitigation | Adapted implementation | Ancestral; expanded to Unicode/control/fullwidth and all string cells. |
| CaoBochun `<caobochun@bytedance.com>` | [#2233](https://github.com/gitleaks/gitleaks/pull/2233) | `b8fc92f16740f33715985a2373e05b8efa663fa0` | `7021b7d242c1089d20359420245de8def434b246` | Escaped SARIF artifact URIs | Adapted implementation | Ancestral; expanded with drive/UNC handling and tests. |
| Dhanush Satyavolu `<satyavoludhanush@gmail.com>` | [#2271](https://github.com/gitleaks/gitleaks/pull/2271) | `158d1e0bc38c49322593401dabff5738afb6dc4b` | `8e787c2e8fd8adddf38d5109e1a943c1b74614a6` | Reuse regex match slice | Direct adoption | Ancestral; final benchmark and code retain single scan. |

### 3.2.2 Documentary, design, research and dependency acknowledgments

These identities are not claimed as commit authors unless they also appear in 3.2.1. Recognition does not imply code incorporation.

| Contributor | Upstream source(s) | Contribution | Attribution type | Final disposition |
|---|---|---|---|---|
| ToniAdreal | [#2300](https://github.com/gitleaks/gitleaks/pull/2300) | Alternative restrictive report permissions | Reviewed but not adopted | Superseded by selected permission work and temp publication. |
| Saadanjum0 | [#2292](https://github.com/gitleaks/gitleaks/pull/2292) | Alternative Git stderr propagation | Reviewed but not adopted | Conflicted with the selected stderr path; documentary credit. |
| CAOShurong | [#2252](https://github.com/gitleaks/gitleaks/pull/2252), [#2247](https://github.com/gitleaks/gitleaks/pull/2247) | Fragment-error alternative; configuration-guide URL correction | Design influence / documentation | Immediate-return implementation rejected; documentation corrected. |
| RajeshRajendiran | [#2203](https://github.com/gitleaks/gitleaks/pull/2203), [#2236](https://github.com/gitleaks/gitleaks/pull/2236), [#2212](https://github.com/gitleaks/gitleaks/pull/2212) | Git drain alternative; added-line and decode-depth documentation | Design influence / documentation | Alternative superseded; documentation credited. |
| trinathone | [#2206](https://github.com/gitleaks/gitleaks/pull/2206), [#2209](https://github.com/gitleaks/gitleaks/pull/2209) | Non-root Docker proposal; upstream Action proposal | Reviewed but not adopted | Optional non-root target is RNC implementation; Action explicitly declined. |
| harriiinnii | [#2241](https://github.com/gitleaks/gitleaks/pull/2241) | Invalid-regex fixture/approach | Reviewed but not adopted | Unrelated expected-regex edits and incomplete context; documentary. |
| kevin9327 | [#2350](https://github.com/gitleaks/gitleaks/pull/2350) | Empty-capture redaction failure mode | Research/defect discovery | Early-return patch rejected; risk incorporated into RNC redaction. |
| ankit090701 | [#2228](https://github.com/gitleaks/gitleaks/pull/2228) | Duplicate invalid-regex proposal | Reviewed but not adopted | No authored commit incorporated. |
| Bhaskar Gurram `<gurrambhaskar.ai@gmail.com>` | [#2301](https://github.com/gitleaks/gitleaks/pull/2301) | Merge-only Git coverage proposal | Design influence; commit not adopted | Policy keeps default; first-parent opt-in documented/tested. |
| Jan Holthuis (ETAS-SEC/XPC-Bo1) `<jan.holthuis@etas.com>` | [#1918](https://github.com/gitleaks/gitleaks/pull/1918) | Pre-commit filename-forwarding hook proposal | Design influence; commit not adopted | Broad `dir` approach rejected; separate `files` command implemented. |
| bharathvemula | [#2298](https://github.com/gitleaks/gitleaks/pull/2298) | Baseline prose/examples | Documentation | Credited in baseline guidance. |
| haimingZZ | [#2282](https://github.com/gitleaks/gitleaks/pull/2282) | README/CONTRIBUTING typo corrections | Documentation | Credited. |
| jankesec | [#2258](https://github.com/gitleaks/gitleaks/pull/2258) | Configuration example correction | Documentation | Credited. |
| ump45nose | [#2287](https://github.com/gitleaks/gitleaks/pull/2287) | Pre-commit Docker anchor correction | Documentation/adapted | Corrected anchor adopted in documentation. |
| jylenhof | [#2008](https://github.com/gitleaks/gitleaks/pull/2008) | Optional mise installation instructions | Documentation | Credited; command verified. |
| noblehelm | [#1263](https://github.com/gitleaks/gitleaks/issues/1263) | Baseline credential-storage/security guidance | Research/documentation | Restricted storage/retention guidance implemented. |
| abhinavrastogi2627 | [#2319](https://github.com/gitleaks/gitleaks/pull/2319), issue [#2302](https://github.com/gitleaks/gitleaks/issues/2302) | System-hook filename-forwarding metadata correction | Documentation/design influence | Staged hook `pass_filenames: false` retained; explicit-file hooks added later. |
| carfeii | [#2299](https://github.com/gitleaks/gitleaks/issues/2299) | Report permissions/umask exposure | Research/defect discovery | Credited; publication permissions corrected. |
| lochlainnw | [#338](https://github.com/gitleaks/gitleaks/issues/338) | Repeated report overwrite/aggregation problem | Research/defect discovery | Aggregation is an opt-in RNC capability. |
| mesaglio | [#1065](https://github.com/gitleaks/gitleaks/issues/1065) | Fail-fast report-destination validation | Research/defect discovery | Pre-scan validation implemented. |
| joerg1985 | [#1572](https://github.com/gitleaks/gitleaks/issues/1572) | Ignore-file special-file/blocking safety | Research/defect discovery | Guarded reader/limits implemented. |
| siliroid | [#2232](https://github.com/gitleaks/gitleaks/issues/2232) | Source-error visibility/false-clean report | Research/defect discovery | Error collector implemented. |
| farrrr | [#2337](https://github.com/gitleaks/gitleaks/issues/2337) | Staged binary-classified text omission | Research/defect discovery | Staged index fallback implemented. |
| bhandrigan | [#2316](https://github.com/gitleaks/gitleaks/issues/2316) | Archive filename false-classification | Research/defect discovery | Content verification implemented. |
| caiodelgadonew | [#1979](https://github.com/gitleaks/gitleaks/issues/1979) | Wrong-directory `dir` behavior | Research/defect discovery | Target validation implemented. |
| rgmz | [#1651](https://github.com/gitleaks/gitleaks/issues/1651), [#1757](https://github.com/gitleaks/gitleaks/issues/1757) | Chunk/decode boundary and input coverage | Research/defect discovery | Long-line boundary correction; decode resource limitation documented. |
| SyedMonzAli | [#1243](https://github.com/gitleaks/gitleaks/issues/1243) | Verbose redaction exposure | Research/defect discovery | Partially addressed; unknown adjacent text residual. |
| AnthonyPoschen | [#1846](https://github.com/gitleaks/gitleaks/issues/1846) | Malformed user config regex panic | Research/defect discovery | Error-returning validation implemented. |
| sporteka2 | [#2296](https://github.com/gitleaks/gitleaks/issues/2296) | Redaction/report configuration request | Research/design influence | Deferred; no new option in this cycle. |
| grodriguezl | [#2010](https://github.com/gitleaks/gitleaks/issues/2010) | Unknown-extension archive discovery | Research/defect discovery | Deferred by archive scope policy. |
| corani | Dependency/RE2 migration research | Lip Gloss v2 migration identification | Dependency/documentation | Profile-only appreciation; no PR adopted. |
| Dependabot | Dependency update proposals | Matched-version recommendations | Dependency/documentation | Profile-only; RNC independently resolved module graph. |
| dgn | [#1711](https://github.com/gitleaks/gitleaks/pull/1711) | Pre-existing upstream Git warning behavior | Pre-existing baseline | Already merged upstream; not an RNC adoption. |
| puneetdixit200 | [#2128](https://github.com/gitleaks/gitleaks/pull/2128) | Quoted `--log-opts` handling and regression research | Design influence; no commit adopted | V13 implemented an independent parser; no upstream patch copied. |
| JoostVoskuil | [#1809](https://github.com/gitleaks/gitleaks/pull/1809) | Azure DevOps SSH remote/link research | Design influence; no commit adopted | V13 implemented current-architecture normalization; old patch targets removed files. |
| mn3m0nic | [#1162](https://github.com/gitleaks/gitleaks/pull/1162) | Repeated-finding/secret-hash proposal | Reviewed/design influence; rejected unsafe default | V13 rejects unkeyed secret hashes; keyed grouping remains policy-gated. |
| Coruscant11 | [#1232](https://github.com/gitleaks/gitleaks/pull/1232) | Multiple report writers proposal | Reviewed/design influence; not adopted | V13 documented a safe plan but did not force the legacy architecture. |

## 3.3 Attribution completeness and reconciliation

### 3.3.1 Totals

| Measure | Count |
|---|---|
| Retained upstream commits | **29** |
| Retained upstream PRs | **25** |
| Unique commit-level author identities | **20** |
| V10 retained commits / PRs | 18 / 17 |
| V11 retained commits / PRs | 6 / 3 |
| V12 retained commits / PRs | 5 / 5 |
| RNC reconciliation commits | 3 product reconciliations (V10, V11, V12) plus this corrective/docs commit |
| Fabricated co-author trailers | 0 |
| Retained commits missing source trailer | 0 |
| Retained commits with author mismatch | 0 |

These totals are different measurements: commits are not PRs, and authors are not commits. Multiple commits from one PR or author are not double-counted as unique contributors.

### 3.3.2 Ancestry and evidence checks

1. Every adopted upstream commit mapped above is ancestral to the final integration branch.
2. Original author names/emails match the source commit objects.
3. Every retained commit carries `(cherry picked from commit <40-hex>)` for the source SHA.
4. `5f6ec56` is ancestral to `rnc-mod` and to the V12 candidate.
5. V10 contributor ancestry is preserved; no V10/V11 contributor commit was squashed into the RNC product commits.
6. V11 contributor ancestry is preserved on the `rnc-mod` line.
7. V12 contributor commits are the last five ancestral commits before the V12 reconciliation.
8. No V10/V11/V12 development checkpoint commit was replayed into integration history.
9. No contributor appears twice in the commit-level register; multi-PR contributors are consolidated in one row where practical, with all PRs named.
10. Adapted implementations are labeled as such; no adapted contributor is presented as an RNC-original invention.
11. Reviewed-not-adopted proposals are not represented as integrated code.
12. Documentary/dependency recognition is kept separate from commit authorship.

### 3.3.3 Historic attribution reconciliation

The pre-V10 ledger recorded that no external PR was yet a verified separately adopted commit. That statement was accurate before the V10 contributor audit but is obsolete. The V10 reconciliation established 18 ancestry-preserved commits from 17 PRs; V11 added six from three PRs; V12 added five from five PRs. The final register supersedes the earlier “no adopted commits” statement.

V13 adopted no upstream commit: its parser, configuration summary, ignore-usage diagnostic, shallow/bare scope warnings, Azure normalization, SARIF level and explicit-file concurrency are RNC-authored. Design reviews of #2128, #1809, #1162 and #1232 are credited above without asserting code adoption.

Similarly, early checkpoint labels such as “implemented” or “validated” described a working branch at that time. They are reconciled here to final behavior: some early designs were replaced, some contributor implementations were adapted, and some defects described in historical reports were subsequently corrected. The historical labels are not republished as current status.

## 3.4 Contribution summary

RNC-Mod recognizes **20 commit-level contributors across 29 retained upstream commits from 25 PRs**, plus documentary, design, research and dependency acknowledgments recorded above. All 29 commits retain original author identity and source trailers; none was fabricated or collapsed into an RNC-only authorship claim. The three RNC reconciliation commits represent integration work, not upstream authorship.

This attribution register is the single authoritative recognition location for RNC-Mod. It records direct adoption, adapted implementation, independent design influence, defect discovery, documentation/dependency work, reviewed-not-adopted proposals and pre-existing upstream baseline work as distinct categories. The project preserves the upstream MIT license and copyright notice; Git ancestry, source trailers and this register together provide traceability without claiming that every contributor’s code survives byte-for-byte in the final tree.

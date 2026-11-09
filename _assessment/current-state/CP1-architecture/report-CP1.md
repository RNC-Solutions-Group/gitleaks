# CP1 — Architecture and execution-flow assessment

**Status:** Draft — ready for owner review  
**Assessment date:** 2026-11-09 (America/New_York)  
**Repository / branch:** `RNC-Solutions-Group/gitleaks` / `feat/issues-and-pr`  
**Assessed code SHA:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (canonical code baseline)  
**Branch HEAD at assessment:** `9fb7c4e6bcf6af0b05b281fc46c4a33da8314624` (documentation-only commits; source tree identical to the code baseline)  
**Environment:** CP0 environment carried forward — Windows 11 AMD64; native Go/module cache/`make`/C compiler absent; Docker `golang:1.27.3` available but dependencies unresolved. No executable Go diagnostics were run in CP1.  
**Assignment:** [`assignment-CP1.md`](assignment-CP1.md)

---

## Executive assessment

The fork is a single Go binary assembled from a Cobra CLI (`cmd`), a Viper/TOML configuration layer (`config`), three fragment sources (`sources`), a detector with a recursive decoder (`detect`, `detect/codec`), a regex abstraction (`regexp`), and five report writers (`report`). Every scanning mode ultimately reaches `Detector.DetectSource`. Most call it directly; the deprecated `protect` command reaches it indirectly through the deprecated `DetectGit` wrapper, which replaces the command context with `context.Background()`.

The control path is: command-specific bootstrap → `initConfig` → `initDiagnostics` → `config.Translate` → `Detector` construction → `DetectSource` → source-specific fragment production → per-fragment `DetectContext` (global commit/path allowlists, keyword prefilter, per-rule `detectRule`, up to `MaxDecodeDepth` decode iterations) → per-finding `AddFinding` (fingerprint ignore, baseline, mutex append) → summary/report/exit. Findings are filtered and redacted in `DetectContext` before they reach storage or reporting.

CP1 corrects four CP0 assumptions: initialization order is config → diagnostics → config translation → detector (not config translation first); `protect` is a separate execution path with a background context; `MaxDecodeDepth` produces up to `N+1` detection passes (one raw pass plus one per decode iteration); and CLI scan errors do force exit 1 when they surface — the fail-open behavior comes from error absorption inside the sources and `DetectSource`, not from the summary function.

The most consequential source-level findings are: empty-secret path-only findings are corrupted by `--redact` (`REDACTED` inserted between every character); Git commit allowlists are logged as skipped but not skipped in `Git.Fragments` (the detector-level guard still filters findings); `Git.Fragments` and `sources.Files.Fragments` discard worker errors and `DetectSource` absorbs fragment errors, so partial scans can return success; `protect` bypasses `--timeout`/cancellation; report-path validation deletes an existing target before scanning; and `extendDepth`, `viper`, and the logger are package globals that make repeated library configuration unsafe. All are source-verified; none was executed.

---

## Scope and methodology

**In scope:** the CP1 assignment’s ten required investigations — topology/interfaces, CLI paths, source lifecycle, configuration/allowlist precedence, detection/decoding order, concurrency/cancellation/error flow, output contracts, build/platform variants, and handoff. **Out of scope:** runtime behavior, capability/regression testing, performance measurement, security exploitation, PR adoption, and any source or backlog edit.

**Method:** read-only Git verification and full-file source reading at the pinned SHA; cross-check of CP0; static call-graph tracing with `git grep`; no build, test, scanner, benchmark, `make`, `go generate`, or network operation. CP0’s Windows/CRLF/symlink/fixture caveats are carried forward unchanged.

**Evidence labels:** **Verified (source)** = directly read in the pinned tree; **Verified (Git)** = produced by a read-only Git command; **Measured** would mean executed (not used in CP1); **Inferred** = reasoned consequence of verified code with no execution; **Unverified** = requires runtime/other evidence.

**CP0 cross-check:** the CP0 report, roadmap, and shared template were re-read. CP0’s opening state, baseline, and source equality were re-verified; its fixture and toolchain limitations were re-checked against the unchanged source tree and the current environment. Its initial data-flow map is confirmed with the refinements recorded in §7 of this report.

---

## Findings register

| ID | Claim | Status | Evidence | Impact | CP owner |
| --- | --- | --- | --- | --- | --- |
| CP1-001 | Canonical source tree is identical at `5f6ec56` and branch HEAD `9fb7c4e` | Verified (Git) | `git diff --name-only 5f6ec56..HEAD` contains only `AGENTS.md`, `README.md`, `_assessment/**`, `_backlog/**`; `git diff --stat` over Go/build paths is empty | Findings can cite `5f6ec56` and remain valid at HEAD | All |
| CP1-002 | All scanning modes ultimately reach `DetectSource`; `protect` uses an indirect `DetectGit` path that loses the command context | Verified (source) | `runDetect`/`runGit`/`runDirectory`/`runStdIn` call `Detector.DetectSource`; `runProtect` calls `DetectGit` ([`cmd/protect.go:61`](../../../cmd/protect.go#L61)), which wraps `DetectSource` with `context.Background()` ([`detect/git.go:19-29`](../../../detect/git.go#L19-L29)) | `protect` loses CLI context/timeout/cancellation; divergent error and shutdown behavior | CP2 |
| CP1-003 | Empty-secret path-only findings are corrupted when redaction is 100% | Verified (source) / Inferred (runtime) | Path-only finding has no `Secret` ([`detect/detect.go:397-405`](../../../detect/detect.go#L397-L405)); `Redact` sets `secret = "REDACTED"` then `strings.ReplaceAll(f.Match, "", "REDACTED")` ([`report/finding.go:78-86`](../../../report/finding.go#L78-L86)); `filter` applies it to every finding ([`detect/utils.go:158-160`](../../../detect/utils.go#L158-L160)) | `--redact` can turn `file detected: …` into repeated `REDACTED` text; no test covers path-only + redact | CP2 |
| CP1-004 | Git commit allowlists are logged as skipped but not skipped at the source | Verified (source) | `continue` at [`sources/git.go:341-346`](../../../sources/git.go#L341-L346) continues the allowlist loop, not the diff loop; commit info is then built and the diff is scanned | Redundant scanning and misleading logs; findings are still filtered by the detector-level allowlist check ([`detect/detect.go:301-305`](../../../detect/detect.go#L301-L305)), so no confirmed false negative — runtime verification needed | CP2 |
| CP1-005 | Fragment/worker errors are absorbed; a partial scan can return success | Verified (source) | `DetectSource` logs fragment errors and returns `nil` ([`detect/detect.go:228-232`](../../../detect/detect.go#L228-L232)); `Git.Fragments` returns `yield(Fragment{}, err)` and ignores `wg.Wait()` ([`sources/git.go:414-430`](../../../sources/git.go#L414-L430)); `Files.Fragments` returns only the walk error ([`sources/files.go:142-190`](../../../sources/files.go#L142-L190)) | CLI may report success on incomplete scans; scanner integrity risk | CP2 |
| CP1-006 | `protect` bypasses the `--timeout` context | Verified (source) | Timeout is installed on the command context ([`cmd/root.go:48-58`](../../../cmd/root.go#L48-L58)); `DetectGit` uses `context.Background()` ([`detect/git.go:21`](../../../detect/git.go#L21)) | A hung `protect` scan can outlive `--timeout`; cancellation does not propagate | CP2 |
| CP1-007 | Report-path validation truncates/deletes an existing target before scanning | Verified (source) | `os.Create(reportPath)` then `os.Remove(reportPath)` for any non-`-` path ([`cmd/root.go:347-358`](../../../cmd/root.go#L347-L358)) | A pre-existing report (or a symlink target) is destroyed before findings exist; data-loss hazard | CP2/CP4 |
| CP1-008 | `--report-template` does not imply `--report-format=template` despite the help text | Verified (source) | Help says “implies --report-format=template” ([`cmd/root.go:79`](../../../cmd/root.go#L79)); code fatals unless the format is exactly `template` ([`cmd/root.go:399-402`](../../../cmd/root.go#L399-L402)) | CLI contract/help mismatch; automation may fail | CP2 |
| CP1-009 | Config/baseline path comparisons are form-sensitive and can miss the file they intend to skip | Verified (source) / Inferred | `Config.Path` is set from `--config` or `<source>/.gitleaks.toml` ([`cmd/root.go:274-283`](../../../cmd/root.go#L274-L283)); git fragments use repo-relative paths ([`sources/git.go:399-405`](../../../sources/git.go#L399-L405)); skip is strict equality ([`detect/detect.go:294-299`](../../../detect/detect.go#L294-L299)); baseline path is made relative to source ([`detect/baseline.go:52-79`](../../../detect/baseline.go#L52-L79)) | Git scans by path can miss the config-path skip (masked by the default `gitleaks\.toml` allowlist); directory scans by path can miss the baseline skip, so baseline JSON may be re-scanned | CP2 |
| CP1-010 | `extendDepth` is a package global that is never reset | Verified (source) | Declared at [`config/config.go:24`](../../../config/config.go#L24), incremented at [`config/config.go:376`](../../../config/config.go#L376) and [`:396`](../../../config/config.go#L396); no production reset; tests reset manually ([`config/config_test.go:582`](../../../config/config_test.go#L582)) | Repeated config translation in one process (library/tests) can skip validation/extension and set the wrong config path | CP2 |
| CP1-011 | Invalid regexes panic rather than returning a configuration error | Verified (source) | `regexp.MustCompile` in rule/allowlist translation ([`config/config.go:123-128`](../../../config/config.go#L123-L128), [`:341-348`](../../../config/config.go#L341-L348)); `joinRegexOr` also uses `MustCompile` ([`config/utils.go:29-40`](../../../config/utils.go#L29-L40)) | A malformed user config crashes the process instead of producing a controlled error | CP2 |
| CP1-012 | Redaction does not cover all early diagnostic logs | Verified (source) | Trace/debug logs print the unredacted secret before `filter` runs ([`detect/detect.go:510-515`](../../../detect/detect.go#L510-L515), [`:545-553`](../../../detect/detect.go#L545-L553)); decoder debug logs print encoded/decoded values ([`detect/codec/decoder.go:93-101`](../../../detect/codec/decoder.go#L93-L101)) | `--redact` does not guarantee secret-free debug/trace logs; logging-level dependent | CP4 |
| CP1-013 | Worker errors are discarded and the semaphore is never waited on | Verified (source) | `s.Sema.Go(...)` results are not collected in [`sources/files.go:151-178`](../../../sources/files.go#L151-L178) or [`sources/git.go:362-413`](../../../sources/git.go#L362-L413); only local `wg` is awaited | Per-file failures can vanish; error reporting depends on logs | CP2/CP3 |
| CP1-014 | The 40-worker concurrency bound is hard-coded and not configurable | Verified (source) | `semgroup.NewGroup(ctx, 40)` ([`detect/detect.go:130`](../../../detect/detect.go#L130)) | Concurrency/backpressure is an architectural constant; performance tuning is not exposed | CP3 |
| CP1-015 | `--enable-rule` is applied after the keyword prefilter is built | Verified (source) | Prefilter built in `NewDetectorContext` ([`detect/detect.go:121-131`](../../../detect/detect.go#L121-L131)); rules replaced later ([`cmd/root.go:332-345`](../../../cmd/root.go#L332-L345)) | Prefilter remains a superset; correctness preserved, but unnecessary keyword work remains | CP3 |
| CP1-016 | Default `max-target-megabytes=0` means no size guard, and the guard is evaluated after path checks and decoding | Verified (source) | Flag default [`cmd/root.go:84`](../../../cmd/root.go#L84); detector guard [`detect/detect.go:430-440`](../../../detect/detect.go#L430-L440); path-only handling precedes it ([`:393-428`](../../../detect/detect.go#L393-L428)); dir source pre-check exists only when the flag is nonzero ([`sources/files.go:92-99`](../../../sources/files.go#L92-L99)); Git has no pre-read cap | Supports #2019 plausibility; bounds are post-decode and opt-in | CP3 |
| CP1-017 | The `.gitleaksignore` loader follows links, has no aggregate size/entry cap, and ignores `scanner.Err()` | Verified (source) | `fileExists` uses `os.Stat` ([`cmd/root.go:502-515`](../../../cmd/root.go#L502-L515)); `AddGitleaksIgnore` uses `os.Open` + `bufio.Scanner` and returns `nil` without checking `scanner.Err()` ([`detect/detect.go:153-193`](../../../detect/detect.go#L153-L193)) | #1572 untrusted-input path remains structurally present; scanner’s default line limit applies but its error is unchecked; no runtime reproduction attempted | CP4 |
| CP1-018 | Findings retain an unused `Fragment` field, and JSON can expose whatever a library embeds | Verified (source) | Field at [`report/finding.go:49-52`](../../../report/finding.go#L49-L52); no first-party assignment (`git grep 'Fragment:'` finds none) | Dead field in CLI; library users can unintentionally serialize a fragment | CP2 |
| CP1-019 | Report writers have distinct empty-findings and secret-exposure contracts | Verified (source) | CSV returns without writing for zero findings ([`report/csv.go:15-18`](../../../report/csv.go#L15-L18)); SARIF places `Secret` in the snippet ([`report/sarif.go:110-133`](../../../report/sarif.go#L110-L133)); template removes `env`, `expandenv`, `getHostByName` ([`report/template.go:33-37`](../../../report/template.go#L33-L37)); JSON omits `Line` ([`report/finding.go:24`](../../../report/finding.go#L24)) | Empty CSV files, redaction-dependent SARIF exposure, and a constrained template environment | CP2/CP4 |
| CP1-020 | Report file descriptors are not closed on the exit paths | Verified (source) | `defer file.Close()` ([`cmd/root.go:477-479`](../../../cmd/root.go#L477-L479)) is followed by `os.Exit` at [`:493-499`](../../../cmd/root.go#L493-L499), which skips defers | Writers are mostly unbuffered, but this is a lifecycle inconsistency and a risk if buffering is added | CP2 |
| CP1-021 | Tests rename `testdata/repos/*/dotGit` ↔ `.git` in place | Verified (source) | [`detect/detect_test.go:1282-1283`](../../../detect/detect_test.go#L1282-L1283), [`:1368-1369`](../../../detect/detect_test.go#L1368-L1369), [`:2432-2455`](../../../detect/detect_test.go#L2432-L2455) | CP1 did not run tests; a writable disposable copy is required for CP2 | CP2 |
| CP1-022 | Release builds use a different regex engine (`gore2regex`) than default/CI builds | Verified (source) | Build tags at [`regexp/stdlib_regex.go:1`](../../../regexp/stdlib_regex.go#L1), [`regexp/wasilibs_regex.go:1`](../../../regexp/wasilibs_regex.go#L1); release tag at [`.goreleaser.yml:18-19`](../../../.goreleaser.yml#L18-L19) | Engine-specific compile/runtime behavior is not exercised by the default test path | CP2/CP4 |
| CP1-023 | A repository-root `gitleaks.toml` is allowlisted by default, but baseline JSON is not | Verified (source) | Default path allowlist `gitleaks\.toml` ([`cmd/generate/config/base/config.go:63-64`](../../../cmd/generate/config/base/config.go#L63-L64)); baseline is a JSON report ([`detect/baseline.go:37-49`](../../../detect/baseline.go#L37-L49)) | Baseline files can be scanned and their embedded secrets re-reported when path matching misses | CP2/CP4 |
| CP1-024 | Rule traversal uses Go map order, so finding order is not deterministic | Verified (source) | `for _, rule := range d.Config.Rules` ([`detect/detect.go:327`](../../../detect/detect.go#L327)); findings are appended in rule/match order | Report ordering and diffs vary across runs; tests must compare order-insensitively (they use `ElementsMatch`, e.g. [`detect/detect_test.go:2121`](../../../detect/detect_test.go#L2121)) | CP2 |

---

## Execution evidence

No build, test, scanner, benchmark, `make` target, generation, or network operation was executed. Commands below are the read-only investigation set; file contents were read with the file-reading tool.

| Command | Purpose | Result / exit |
| --- | --- | --- |
| `git rev-parse --show-toplevel`, `rev-parse --abbrev-ref HEAD`, `rev-parse HEAD` | Identity/HEAD | Root, `feat/issues-and-pr`, `9fb7c4e…`; exit 0 |
| `git status --short --branch`, `git status --porcelain=v1` | Opening integrity | Clean; exit 0 |
| `git merge-base HEAD master`; `git rev-list --count 5f6ec56..HEAD` | Baseline relationship | Merge base `5f6ec56`; 62 documentation commits; exit 0 |
| `git diff --name-only 5f6ec56..HEAD`; `git diff --stat` over Go/build paths | Source equality | Documentation paths only; Go/build diff empty; exit 0 |
| `git grep -n <symbol>` (multiple) | Call-graph/ownership tracing | Exit 0 on matches; exit 1 only where a pattern had no match (noted in analysis) |
| `git ls-files`, `Get-ChildItem`, `Get-Content`/file reads | Package and asset inventory | Exit 0 |
| `go build/test/vet`, `make`, `go generate` | Executable diagnostics | **NOT RUN** — CP1 does not authorize or require them; CP0 dependency/toolchain limits still apply |

---

## 1. Source-pinned component map

```
main.go
  └─ cmd.Execute ── Cobra root + subcommands
       ├─ config bootstrap (viper) ── config.ViperConfig.Translate ── config.Config
       ├─ diagnostics manager (global)
       ├─ Detector construction (flags, .gitleaksignore, baseline, reporter)
       └─ DetectSource ── sources.Source
            ├─ sources.Files ── sources.File (chunking, archives, symlinks)
            ├─ sources.File  ── stdin / single reader
            └─ sources.Git   ── git log/diff → gitdiff.TextFragments
                 └─ detect.DetectContext
                      ├─ global commit/path allowlists
                      ├─ keyword prefilter (aho-corasick)
                      ├─ per-rule detectRule
                      │    ├─ path/commit allowlists, path-only findings
                      │    ├─ size guard
                      │    ├─ regex matches + location construction
                      │    ├─ gitleaks:allow, secret group, entropy
                      │    ├─ global/rule finding allowlists
                      │    └─ required-rule recursion
                      ├─ decoder.Decode (percent/unicode/hex/base64) up to MaxDecodeDepth
                      └─ filter (generic dedupe + redaction)
                           └─ AddFinding (fingerprint ignore, baseline, mutex)
                                └─ Reporter.Write (json/csv/junit/sarif/template)
```

| Component | Responsibility | Entry / exit symbols | Exchanged types | Owned state | Depends on | Evidence |
| --- | --- | --- | --- | --- | --- | --- |
| `main` | Process entry; SIGINT handler | `main` → `cmd.Execute`; `listenForInterrupt` | — | signal channel | `cmd`, `logging` | [`main.go:11-24`](../../../main.go#L11-L24) |
| `cmd` | CLI wiring, flags, config bootstrap, detector assembly, report selection, summary/exit | `rootCmd`, `Execute`, `Config`, `Detector`, `findingSummaryAndExit`; `runDetect`/`runGit`/`runDirectory`/`runProtect`/`runStdIn`/`runVersion` | `config.Config`, `*detect.Detector`, `[]report.Finding`, errors | global `rootCmd`, `diagnosticsManager`, `logLevel` | cobra, viper, all lower packages | [`cmd/root.go:43-103`](../../../cmd/root.go#L43-L103), [`cmd/root.go:224-500`](../../../cmd/root.go#L224-L500) |
| `config` | TOML→`ViperConfig`→`Config`; rules, allowlists, extension, validation, embedded default | `ViperConfig.Translate`, `Rule.Validate`, `Allowlist.Validate`, `GetOrderedRules` | `ViperConfig`, `Config`, `Rule`, `Allowlist`, `Extend` | global `viper`, `extendDepth`; embedded `DefaultConfig` | viper, regexp, go-version | [`config/config.go:18-281`](../../../config/config.go#L18-L281), [`config/rule.go:12-107`](../../../config/rule.go#L12-L107), [`config/allowlist.go:29-179`](../../../config/allowlist.go#L29-L179) |
| `sources` | Produce `Fragment`s from files, directories, stdin, or Git | `Source.Fragments`; `Files`, `File`, `Git` | `Fragment`, `ScanTarget`, `CommitInfo`, `RemoteInfo` | per-source buffers/depths; detector `Sema` | archives, filetype, gitdiff, semgroup | [`sources/source.go:9-15`](../../../sources/source.go#L9-L15), [`sources/file.go:52-264`](../../../sources/file.go#L52-L264), [`sources/files.go:142-191`](../../../sources/files.go#L142-L191), [`sources/git.go:299-431`](../../../sources/git.go#L299-L431) |
| `detect` | Detection loop, rule evaluation, decoding orchestration, ignore/baseline, finding storage | `NewDetectorContext`, `DetectSource`, `DetectContext`, `detectRule`, `AddFinding`, `filter` | `Fragment`, `config.Rule`, `codec.EncodedSegment`, `report.Finding` | findings slice + mutex, commit map + mutex, baseline, ignore map, prefilter, `Sema`, `TotalBytes` | config, codec, regexp, report, sources | [`detect/detect.go:39-132`](../../../detect/detect.go#L39-L132), [`detect/detect.go:208-368`](../../../detect/detect.go#L208-L368), [`detect/detect.go:371-900`](../../../detect/detect.go#L371-L900) |
| `detect/codec` | Detect/decode encoded segments across iterations | `NewDecoder`, `Decode`, `findEncodedSegments`, `AdjustMatchIndex`, `CurrentLine` | `EncodedSegment` tree, `encodingKind` | per-`DetectContext` decoder map and segment slices | regexp | [`detect/codec/decoder.go:10-105`](../../../detect/codec/decoder.go#L10-L105), [`detect/codec/encodings.go:22-161`](../../../detect/codec/encodings.go#L22-L161), [`detect/codec/segment.go:35-173`](../../../detect/codec/segment.go#L35-L173) |
| `regexp` | Compile-time regex engine selection | `MustCompile`, `Version` | `*regexp.Regexp` alias | — | stdlib or go-re2 by build tag | [`regexp/stdlib_regex.go:1-14`](../../../regexp/stdlib_regex.go#L1-L14), [`regexp/wasilibs_regex.go:1-14`](../../../regexp/wasilibs_regex.go#L1-L14) |
| `report` | Finding model; output writers | `Finding`, `Redact`, `Reporter`, five writers | `Finding`, `RequiredFinding`, SARIF/JUnit DTOs | none (writers are stateless except template) | encoding/json/xml/csv, sprig, lipgloss | [`report/finding.go:14-126`](../../../report/finding.go#L14-L126), [`report/report.go:7-16`](../../../report/report.go#L7-L16), [`report/json.go`](../../../report/json.go#L1), [`report/csv.go`](../../../report/csv.go#L1), [`report/junit.go`](../../../report/junit.go#L1), [`report/sarif.go`](../../../report/sarif.go#L1), [`report/template.go`](../../../report/template.go#L1) |
| `logging` | Global zerolog logger to stderr | `Logger`, `With`, level helpers | zerolog events | global `Logger` | zerolog | [`logging/log.go:10-48`](../../../logging/log.go#L10-L48) |
| `version` | Build-time version string | `Version`, `DefaultMsg` | string | ldflags-provided | — | [`version/version.go:4-5`](../../../version/version.go#L4-L5) |
| `cmd/scm` | SCM platform parsing for links | `PlatformFromString`, `Platform` | enum | — | — | [`cmd/scm/scm.go:8-52`](../../../cmd/scm/scm.go#L8-L52) |
| `cmd/generate` | Generate tracked/embedded default rule config | `//go:generate` main; 222-rule list | `config.Config` → TOML template | none | rules/base packages | [`cmd/generate/config/main.go:14-301`](../../../cmd/generate/config/main.go#L14-L301) |

**Interfaces and ownership:** `sources.Source` has one method, `Fragments(ctx, yield)`; the callback signature returns an error, but the CLI path (`DetectSource`) returns `nil` on fragment errors, while the deprecated `DetectReader`/`StreamDetectReader` paths propagate them. `report.Reporter` has one method, `Write(io.WriteCloser, []Finding)`. `config.Config` is a value copied into the detector; its maps/regexes are then read-only by contract, though nothing enforces that.

---

## 2. CLI execution-path matrix

Startup is Cobra-driven: `main` installs an interrupt listener and calls `cmd.Execute` ([`main.go:11-19`](../../../main.go#L11-L19)); `init()` registers root flags and binds only `config` to Viper ([`cmd/root.go:73-103`](../../../cmd/root.go#L73-L103)); `cobra.OnInitialize(initLog)` sets the global log level before each command; `PersistentPreRunE` installs an optional timeout on the command context ([`cmd/root.go:48-58`](../../../cmd/root.go#L48-L58)). `Execute` maps an “unknown flag” error to exit 126 and all other Cobra errors to a fatal log ([`cmd/root.go:224-232`](../../../cmd/root.go#L224-L232)).

| CLI mode | Source selection | Bootstrap | Detector entry | Error policy | Report / exit |
| --- | --- | --- | --- | --- | --- |
| `git [repo]` | positional path, default `.` | `initConfig(source)` → diagnostics → `Config` → `Detector` | `DetectSource(&sources.Git{...})` | logs `DetectSource` error, then summary exits 1 if error survived | report + `--exit-code`/1 |
| `git --pre-commit` / `git --staged` | positional path | same | `DetectSource(&sources.Git{Cmd: NewGitDiffCmdContext(...)})` | same | same |
| `dir [path]` (aliases `file`, `directory`) | positional path, default `.` | same | `DetectSource(&sources.Files{...})` | logs walk error, then summary | same |
| `detect` (hidden) | `--source`, default `.`; history unless `--no-git`/`--pipe` | same | `DetectSource(Git/Files/File)` | Git/Files: log then summary; pipe: fatal on error | same |
| `detect --no-git` | `--source` | same | `DetectSource(&sources.Files{...})` | log then summary | same |
| `detect --pipe` | `os.Stdin` | same | `DetectSource(&sources.File{Content: os.Stdin})` | fatal on error (no report) | summary otherwise |
| `stdin` | `os.Stdin`; config from `.`; detector source `""` | `initConfig(".")` → diagnostics → `Config` → `Detector(cmd, cfg, "")` | `DetectSource(&sources.File{Content: os.Stdin})` | fatal on error | summary otherwise |
| `protect` (hidden) | `--source`, default `.` | `initConfig(source)` → diagnostics → `Config` → `Detector` | `DetectGit(gitCmd, &detect.RemoteInfo{Platform: NoPlatform})` | logs error then summary | same, but background context |
| `version` | none | none | none | none | prints `version.Version`, exit 0 |
| root without subcommand | — | root has no `Run` | — | Cobra help/usage path; exact exit code unverified | **Unverified** |

**Source-level differences that matter:** `detect` is deprecated and hidden but still wired ([`cmd/detect.go:1-18`](../../../cmd/detect.go#L1-L18), [`:43-48`](../../../cmd/detect.go#L43-L48)); `protect` is deprecated/hidden ([`cmd/protect.go:22-27`](../../../cmd/protect.go#L22-L27)) and is the only CLI path through `DetectGit`; `git --pre-commit` explicitly sets `scm.NoPlatform` so links are not generated ([`cmd/git.go:64-70`](../../../cmd/git.go#L64-L70)); `stdin` and `detect --pipe` differ in config source (`"."` vs `--source`) and detector source argument.

**Flags that change downstream behavior**

| Flag(s) | Applied at | Downstream effect |
| --- | --- | --- |
| `--config`, `GITLEAKS_CONFIG`, `GITLEAKS_CONFIG_TOML` | `initConfig` | rule/allowlist set; see §3 |
| `--max-decode-depth` (default 5) | `Detector` | decode iterations; `N` decodes, `N+1` detection passes |
| `--max-archive-depth` (default 0) | `Detector` → source structs | archive traversal depth; 0 disables |
| `--max-target-megabytes` (default 0) | `Detector`; `Files.MaxFileSize` | opt-in size guard; dir pre-check only |
| `--follow-symlinks` | dir/detect only | `Files` symlink policy |
| `--gitleaks-ignore-path` (default `.`) | `Detector` | three `.gitleaksignore` load attempts |
| `--baseline-path` | `Detector` | baseline suppression and baseline-file skip |
| `--enable-rule` | `Detector` | post-construction rule-map replacement |
| `--redact` | `filter` | redaction of findings (see CP1-003/012) |
| `--report-path/-format/-template` | `Detector` | writer selection and target validation |
| `--exit-code` (default 1) | summary | exit status when findings exist |
| `--verbose`, `--no-color`, `--log-level` | logging/detector | stdout findings; logger configuration |
| `--timeout` | root `PersistentPreRunE` | context deadline for all modes except `protect`/deprecated APIs |
| `--diagnostics`, `--diagnostics-dir` | `initDiagnostics` | pprof CPU/mem/trace or localhost HTTP server |
| `--platform` | git/detect | SCM link generation |
| `--staged`, `--pre-commit`, `--log-opts` | git/detect/protect | diff vs history, Git arguments |
| `--source`, `--no-git`, `--pipe` | detect/protect | source selection |

---

## 3. Data-flow and state-lifetime map

| Stage | Operation | State created / mutated | Lifetime / ownership | Guards and ordering | Evidence |
| --- | --- | --- | --- | --- | --- |
| Config load | `initConfig` selects `--config` / env / source-local / embedded default | global Viper state; banner to stderr | process (global) | precedence chain; file-source fallback to default; fatal on read error | [`cmd/root.go:132-195`](../../../cmd/root.go#L132-L195) |
| Config translation | `Config` unmarshals Viper and calls `Translate` | `Config` value: rules map, keywords, ordered rules, allowlists, extension | returned by value to detector | rule validation, required-rule existence, duplicate IDs last-wins, extension depth ≤2 | [`config/config.go:109-281`](../../../config/config.go#L109-L281), [`cmd/root.go:234-247`](../../../cmd/root.go#L234-L247) |
| Detector assembly | `Detector` reads flags, loads ignores/baseline, selects reporter | `Detector` with prefilter, `Sema` (40), findings/commit maps, baseline, ignore map | one per command invocation | `.gitleaksignore` (3 attempts), baseline error logged, report path validated/deleted, `--enable-rule` override | [`cmd/root.go:249-409`](../../../cmd/root.go#L249-L409) |
| Source production | `Files`/`File`/`Git` yield `Fragment`s | per-source buffers, archive temp files, worker goroutines; `Fragment` strings | per source call; worker lifetime bounded by semaphore/wg | empty-file, size (dir only), symlink, archive depth, binary MIME, Git delete/binary/commit, context | [`sources/file.go:52-251`](../../../sources/file.go#L52-L251), [`sources/files.go:60-190`](../../../sources/files.go#L60-L190), [`sources/git.go:299-431`](../../../sources/git.go#L299-L431) |
| Fragment detection | `DetectContext` applies global allowlists, prefilter, rules, decode loop | local `currentRaw`, `encodedSegments`, `decoder` (per fragment), finding slice | per fragment | config/baseline path skip; commit/path allowlist; keyword prefilter; depth limit; decode termination | [`detect/detect.go:276-368`](../../../detect/detect.go#L276-L368) |
| Rule evaluation | `detectRule` builds findings | matches slice, newline indices, local findings, required-rule findings | per rule/fragment call | skip-report, rule allowlist, path check, size guard, regex, line/column, inline allow, secret group, entropy, allowlists | [`detect/detect.go:371-656`](../../../detect/detect.go#L371-L656) |
| Filtering | `filter` dedupes generic findings and redacts | new findings slice; mutates finding fields | per fragment | generic rule ID heuristic, same line/commit, secret containment; redaction after dedupe | [`detect/utils.go:136-166`](../../../detect/utils.go#L136-L166) |
| Storage | `AddFinding` computes fingerprint, applies ignore/baseline, appends | detector findings slice (mutex), commit map (mutex) | detector lifetime (whole scan) | global fingerprint then commit fingerprint; `IsNew` baseline comparison; verbose print | [`detect/detect.go:710-761`](../../../detect/detect.go#L710-L761), [`detect/baseline.go:12-35`](../../../detect/baseline.go#L12-L35) |
| Reporting | summary stops diagnostics, writes report, exits | report file/stdout; findings read | one per command | reporter selection; `os.Create` truncates; exit precedence error > findings > success | [`cmd/root.go:438-500`](../../../cmd/root.go#L438-L500) |

### Detection and decoding order (verified)

1. `DetectContext` records bytes, skips config/baseline paths, then applies **global commit/path allowlists** before any content work ([`detect/detect.go:294-305`](../../../detect/detect.go#L294-L305)).
2. It creates a **new `codec.Decoder` per fragment** ([`:311`](../../../detect/detect.go#L311)), so the encoded→decoded cache does not outlive the fragment.
3. Each iteration lowercases the current text for keyword matching, builds the keyword set, then iterates `d.Config.Rules` (a Go map, so order is not deterministic); rules with no keywords always run; keyword rules run only if their lowercased keyword is present ([`:319-347`](../../../detect/detect.go#L319-L347)).
4. After rule evaluation, depth increments; if `currentDecodeDepth > MaxDecodeDepth` the loop ends; otherwise `decoder.Decode` produces the next text and segment metadata ([`:349-363`](../../../detect/detect.go#L349-L363)). Therefore **depth 5 performs five decodes and up to six detection passes** (raw plus five decoded).
5. `detectRule` checks `SkipReport`, commit/path allowlists, path-only/path+regex rules, then the optional size guard, then regex matches ([`:383-445`](../../../detect/detect.go#L383-L445)). The size guard sees the **current (possibly decoded)** text and only runs when `MaxTargetMegaBytes > 0`; the default is 0 (no guard).
6. Two `FindAllStringIndex` calls occur: one for the empty check, one for iteration ([`:442`](../../../detect/detect.go#L442), [`:452`](../../../detect/detect.go#L452)). Newline indices are computed once per rule/fragment when at least one match exists ([`:448`](../../../detect/detect.go#L448)).
7. For a decoded pass, matches not overlapping a decoded segment are skipped as already reported ([`:462-471`](../../../detect/detect.go#L462-L471)); match indices are mapped back through the segment tree; `Line` remains the original fragment line while `Match`/`Secret` are decoded.
8. Inline `gitleaks:allow`, secret-group extraction, entropy, and global/rule finding allowlists are applied in that order ([`:510-566`](../../../detect/detect.go#L510-L566)). Required-rule processing runs after primary findings and recurses once with `InheritedFromFinding=true` ([`:570-656`](../../../detect/detect.go#L570-L656)).
9. `filter` dedupes generic findings and applies redaction; `AddFinding` computes `file:rule:line` or `commit:file:rule:line`, checks ignore map, baseline, then appends under a mutex ([`detect/utils.go:137-166`](../../../detect/utils.go#L137-L166), [`detect/detect.go:711-749`](../../../detect/detect.go#L711-L749)).

**#2019 / #1572 cross-check.** The separate pathways described in the backlog remain structurally distinct. #2019 reaches the decode loop through ordinary fragments, with an opt-in post-decode size guard and default depth 5 ([`cmd/root.go:91`](../../../cmd/root.go#L91), [`detect/detect.go:353-358`](../../../detect/detect.go#L353-L358)); the decoder allocates a new buffer per pass and retains per-fragment encoded/decoded strings ([`detect/codec/decoder.go:23-39`](../../../detect/codec/decoder.go#L23-L39)). #1572 reaches the ignore loader through `fileExists`/`os.Stat` and `AddGitleaksIgnore` ([`cmd/root.go:502-515`](../../../cmd/root.go#L502-L515), [`detect/detect.go:153-193`](../../../detect/detect.go#L153-L193)) and is not governed by `--follow-symlinks` (which only affects `sources.Files`). Neither was reproduced. CP1 adds no runtime evidence for either; see [`performance-stability.md:20-42`](../../../_backlog/categories/performance-stability.md#L20-L42) and the security record for canonical ownership.

### Configuration and policy precedence

| Input | Precedence / condition | Loader | Effective value | Consumer | Failure mode |
| --- | --- | --- | --- | --- | --- |
| `--config` path | 1 | `viper.SetConfigFile` + `ReadInConfig` | rules/allowlists; `cfg.Path` flag value | `Translate`, detector skip | fatal |
| `GITLEAKS_CONFIG` path | 2 | same | same; `cfg.Path` becomes `<source>/.gitleaks.toml` (flag empty) | same | fatal |
| `GITLEAKS_CONFIG_TOML` content | 3 | `viper.ReadConfig(bytes)` then early return | same; `cfg.Path` becomes `<source>/.gitleaks.toml` | same | fatal |
| `<source>/.gitleaks.toml` | 4 (source must be a directory and file exists) | `AddConfigPath` + `SetConfigName(".gitleaks")` | same; `cfg.Path` source-joined | same | fatal |
| Embedded default | 5 (no config found) or `[extend] useDefault` | `config.DefaultConfig` (`go:embed`) | default 222-rule config | same | fatal |
| `[extend] path` | during translation, max depth 2 | `viper.SetConfigFile(path)` | merged rules/allowlists; `cfg.Path` set to extension path | same | fatal; relative paths follow Viper/process-CWD semantics (Inferred), not the parent config directory |
| `[extend] url` | not implemented | `extendURL` TODO | no effect | — | silently ignored |
| `minVersion` | translation | semver compare with `version.Version` | warning only | operator | invalid semver → fatal |
| `.gitleaksignore` | after detector creation | three loads: `--gitleaks-ignore-path`, `<path>/.gitleaksignore`, `<source>/.gitleaksignore` | fingerprint set | `AddFinding` | fatal on open error; no size/entry cap; `scanner.Err()` unchecked |
| Baseline | after ignores | JSON unmarshal of a report | baseline findings + baseline path | `AddFinding`/`DetectContext` | logged, scan continues |
| `--enable-rule` | after detector/prefilter | replacement map built from `cfg.Rules` | detector rule set | `DetectContext` | fatal if a requested rule is unknown |
| Report flags | after rules | path validation + writer selection | reporter, report path | summary | fatal on unwritable/unknown format; existing target deleted |
| Diagnostics | before `Config` | CPU/mem/trace or HTTP | profiler state | summary stop | fatal on setup |

**Policy stages are distinct:** (a) config *loading* (Viper), (b) *translation/validation* (`Translate`, rule/allowlist `Validate`), (c) *detector construction* (flags, ignore, baseline, reporter, prefilter), and (d) *per-fragment/per-rule filtering* (global allowlists, keyword prefilter, rule allowlists, finding allowlists, ignore fingerprints, baseline). A failure in (a)/(b) is fatal; failures in (c) are mostly fatal, except baseline load; filtering never returns errors.

---

## 4. Concurrency, cancellation, and error-flow register

| ID | Originating operation | Mechanism | Propagation path | Handling outcome | Evidence | Certainty |
| --- | --- | --- | --- | --- | --- | --- |
| C-01 | File read | `sources.File.fileFragments` yields `(fragment, err)` on read error | `DetectSource` callback | error logged; callback returns `nil`; loop may continue/end | [`sources/file.go:183-248`](../../../sources/file.go#L183-L248), [`detect/detect.go:228-232`](../../../detect/detect.go#L228-L232) | Verified |
| C-02 | Git stderr | `listenForStdErr` sends one generic “stderr is not empty” error | `Git.Fragments` returns `yield(Fragment{}, err)` immediately → `DetectSource` | error logged; callback returns `nil`; Git source stops early with partial findings and no caller-visible error | [`sources/git.go:226-271`](../../../sources/git.go#L226-L271), [`:414-420`](../../../sources/git.go#L414-L420) | Verified |
| C-03 | Git process exit | `defer s.Cmd.Wait()` | local defer | debug log only; not returned | [`sources/git.go:299-304`](../../../sources/git.go#L299-L304) | Verified |
| C-04 | Per-file/worker failure | `s.Sema.Go(...)` result, local `wg` | discarded / only `wg.Wait` | not surfaced; `Files.Fragments` returns walk error only; `Git.Fragments` returns nil | [`sources/files.go:142-190`](../../../sources/files.go#L142-L190), [`sources/git.go:361-430`](../../../sources/git.go#L361-L430) | Verified |
| C-05 | Context cancellation | `PersistentPreRunE` deadline; source `ctx.Done()` checks | `Fragments` returns `ctx.Err()` | `DetectSource` returns it; summary exits 1 (except `protect` background context) | [`cmd/root.go:48-58`](../../../cmd/root.go#L48-L58), [`sources/files.go:184-190`](../../../sources/files.go#L184-L190), [`sources/git.go:313-316`](../../../sources/git.go#L313-L316) | Verified |
| C-06 | Archive extraction | extractor/decompressor errors logged; inner errors returned | worker → semgroup discarded | partial archive scan, no failed status | [`sources/file.go:90-164`](../../../sources/file.go#L90-L164) | Verified |
| C-07 | Config/report failure | `logging.Fatal` | immediate process exit 1 | dedicated exit path, no report | [`cmd/root.go:224-232`](../../../cmd/root.go#L224-L232) and 59 `Fatal` call sites | Verified |
| C-08 | Report write failure | `Reporter.Write` error | summary `Fatal` | exit 1 | [`cmd/root.go:483-490`](../../../cmd/root.go#L483-L490) | Verified |
| C-09 | SIGINT | signal goroutine calls `logging.Fatal` | process exit 1 | no graceful scan cancellation or report finalization | [`main.go:14-24`](../../../main.go#L14-L24) | Verified |
| C-10 | Goroutine leaks | stderr sender may block on an unbuffered channel after early return; HTTP pprof server never shuts down | — | possible leak; `protect` path has no cancellation | [`sources/git.go:111-112`](../../../sources/git.go#L111-L112), [`:226-270`](../../../sources/git.go#L226-L270), [`cmd/diagnostics.go:136-147`](../../../cmd/diagnostics.go#L136-L147) | Inferred |
| C-11 | Findings/commit state | mutexes; `TotalBytes` atomic | shared detector | safe appends; verbose printing occurs while `findingMutex` is held, so output is serialized but slow stdout blocks other workers | [`detect/detect.go:69-107`](../../../detect/detect.go#L69-L107), [`:710-761`](../../../detect/detect.go#L710-L761) | Verified |
| C-12 | Concurrency bound | `semgroup.NewGroup(ctx, 40)` per detector; `Go` is expected to block when the bound is saturated | walk/parse loop backpressure | 40 is a hard constant; no CLI control | [`detect/detect.go:130`](../../../detect/detect.go#L130), [`sources/files.go:150-178`](../../../sources/files.go#L150-L178) | Verified (bound) / Inferred (blocking) |

**Producer/consumer summary:** directory/Git parsing is the producer; each file/diff becomes a worker task; workers produce fragments synchronously into `DetectSource`, which mutates shared findings under a mutex. There is no findings queue or streaming consumer; findings accumulate for the entire scan. Backpressure is expected through the semaphore and, for Git, through the `gitdiff` output channel. Cancellation reaches sources through the command context, but deprecated entry points and `protect` replace it with `context.Background()` ([`detect/git.go:21`](../../../detect/git.go#L21), [`detect/files.go:75`](../../../detect/files.go#L75), [`detect/reader.go:22`](../../../detect/reader.go#L22), [`:89`](../../../detect/reader.go#L89)).

---

## 5. Reporting contracts and build variants

| Format | Writer | Selection | Empty findings | Redaction | Notable contract |
| --- | --- | --- | --- | --- | --- |
| JSON | `JsonReporter` | `--report-format=json` or `.json` inference | `[]` (empty array) | fields already redacted by `filter` | `Line` excluded; `Fragment` omitted if nil; `Secret`/`Match` included |
| CSV | `CsvReporter` | `csv` or `.csv` | **nothing written** (empty file) | inherited | Link column only if the first finding has a link |
| JUnit | `JunitReporter` | `junit` only | suite with zero failures/tests | inherited | failure data is the JSON-encoded finding |
| SARIF | `SarifReporter` | `sarif` or `.sarif` | results `[]`; rules from `OrderedRules` | inherited | snippet is `f.Secret`; symlink URI preferred |
| Template | `TemplateReporter` | `template` + `--report-template` | template renders empty slice | inherited | sprig `env`, `expandenv`, `getHostByName` removed |

**Report path and exit contracts:** only a non-empty `--report-path` enables a reporter; `-` writes to stdout; any other path is created/truncated at validation and again at write; an existing file is deleted during validation ([`cmd/root.go:347-358`](../../../cmd/root.go#L347-L358), [`:463-491`](../../../cmd/root.go#L463-L491)). Exit precedence is: scan error → 1; else findings → `--exit-code` (default 1); else 0 ([`:493-499`](../../../cmd/root.go#L493-L499)). Cobra “unknown flag” is 126. Logs go to stderr; verbose findings and `version` go to stdout; the banner goes to stderr.

| Build variant | Trigger | Code consequence | Verification status |
| --- | --- | --- | --- |
| stdlib regex | default / CI | `regexp.Regexp` = stdlib | default path; tests assume it |
| `gore2regex` | release build tag ([`.goreleaser.yml:18-19`](../../../.goreleaser.yml#L18-L19)) | `regexp.Regexp` = `wasilibs/go-re2`; rule compilation and matching semantics follow RE2/WASM | not exercised in CP1; default tests do not run this tag |
| Generated/embedded config | `//go:generate` + `go:embed` | `config/gitleaks.toml` is tracked, generated, and embedded; 222 rules | generator code read; not executed |
| Version/minVersion | ldflags → `version.Version` | config `minVersion` is skipped for the default dev string and warning-only otherwise | source-verified; release flow not run |
| Windows | runtime | path slash normalization, `WindowsFilePath`, symlink materialization disabled, CRLF | CP0-measured environment; two tests skip |
| Linux/container | CI/Docker | Linux path/symlink behavior; archive extraction | not measured |

---

## 6. Coupling, invariants, and verification backlog

| ID | Coupling / invariant | Why it matters | Verification experiment (not executed) | CP |
| --- | --- | --- | --- | --- |
| AR-01 | `DetectSource` absorbs fragment errors; `Git.Fragments`/`Files.Fragments` discard worker errors | Partial scans may look successful | Induce unreadable file, malformed diff tail, non-empty Git stderr; assert CLI output, findings, exit code, report presence | CP2 |
| AR-02 | `protect` uses `DetectGit` with `context.Background()` | `--timeout` and cancellation do not apply | Run `protect` against a controlled slow/hung diff with a short timeout; verify process exit | CP2 |
| AR-03 | Redaction occurs after detection logs | `--redact` may not cover debug/trace output | Capture debug/trace logs with a seeded synthetic secret and redaction on; scan for the raw value | CP4 |
| AR-04 | Path-only findings have empty `Secret`; `Redact(100)` uses `ReplaceAll` with an empty old string | `--redact` output corruption for path-only rules | Unit-level fixture: path rule + `--redact`; assert Match/Line/report text | CP2 |
| AR-05 | `Config.Path` and baseline paths are compared as strings across different path forms | Config/baseline files may be scanned or skipped inconsistently | Scan a repo by name vs `.`; place baseline inside/outside source; assert self-finding behavior | CP2 |
| AR-06 | `extendDepth`, Viper, and the logger are package globals | Repeated library configuration can retain extension depth/config state | In one process, translate an extended config then a plain config; assert rules/validation/path | CP2 |
| AR-07 | `MustCompile` panics on invalid regex | Malformed config yields a panic, not a CLI error | Feed a config with an invalid regex; assert process behavior/exit | CP2 |
| AR-08 | Git commit allowlist loop `continue` does not skip the diff | Redundant scanning and misleading logs | Git fixture with an allowed commit and a secret; assert finding suppressed and whether diff was still scanned (log/benchmark) | CP2 |
| AR-09 | `max-target-megabytes` is post-decode, opt-in, and absent for Git | A decoded fragment can allocate before the guard | Bounded synthetic encoded workload at depths 0/1/3/5 with/without the flag; measure peak RSS and allocations | CP3 |
| AR-10 | Decoder cache and buffers are per fragment, at most `MaxDecodeDepth` passes | Per-fragment allocation can still be large; cache lifetime is bounded | Profile bounded encoded corpora; attribute allocations to decode passes | CP3 |
| AR-11 | `.gitleaksignore` loader follows links and has no caps | #1572 input-boundary risk remains | Finite synthetic newline-delimited ignore files, controlled symlink, malformed entries, scanner error; assert fail-closed | CP4 |
| AR-12 | `gore2regex` release variant is untested here | Release behavior may differ in compile/runtime | Build with `-tags gore2regex` on a pinned toolchain; run rule compilation and representative fixtures | CP2/CP4 |
| AR-13 | Reporter path validation deletes an existing target | Data-loss hazard before scan | Pre-create a report file, run with `--report-path`, observe deletion; then test `-`/stdout | CP2 |
| AR-14 | `--report-template` help/code mismatch | Automation contract mismatch | Run with only `--report-template`; assert error | CP2 |
| AR-15 | Generic-rule dedupe is O(n²) per fragment and heuristic by rule-ID substring | Finding-heavy inputs may be slow; semantics depend on naming | Benchmark many generic/non-generic findings on one fragment; compare rule-ID variants | CP3 |

**Interfaces/invariants external PRs must preserve**

1. `sources.Source.Fragments(ctx, yield)` semantics: fragment ordering is not guaranteed; `yield` may be concurrent; nil error from the callback means “continue”.
2. `DetectSource` return contract: findings are stored in the detector, not only returned; callers use `Findings()`.
3. One decoder per fragment; decoded text is scanned only through overlapping segments, and finding locations must map back to original text.
4. Finding identity: `file:rule:line` (global) or `commit:file:rule:line` (commit); `.gitleaksignore` depends on both forms; fingerprint changes break suppressions.
5. Rule/allowlist evaluation order: commit/path first, then content; `gitleaks:allow` is checked on the finding line before secret-group/entropy/allowlist filtering.
6. Reporter interface: `Write` receives all findings once; report format and documented exit semantics are compatibility considerations. **Do not preserve destructive report-path validation merely because it is existing behavior.**
7. Build-tag contract: default tests use stdlib regex; release binaries use `gore2regex`.
8. Finding order is not guaranteed because rules are traversed from a map; report consumers must sort or compare order-insensitively.

**Observed behaviors to challenge, not preserve as invariants:** absorbed fragment/worker errors (CP1-005/013), discarded `protect` cancellation context (CP1-002/006), premature report-target deletion (CP1-007), and empty-secret redaction corruption (CP1-003). A later PR may intentionally correct these behaviors if tests establish the expected contract.

**Backlog cross-references (leads only, not reassessed):** #2019 and #1572 remain separate packages in [`performance-stability.md:20-42`](../../../_backlog/categories/performance-stability.md#L20-L42) and the security record; the Git-stream/error work package (#1338/#2246) aligns with C-01–C-04 and AR-01; #2023’s size/decoder ordering concerns align with AR-09/AR-10; #2121/#2271 align with the double `FindAllStringIndex` observation (§3, detection-order item 6). No backlog disposition was changed.

---

## 7. CP0 cross-check and corrections

| CP0 statement | CP1 result |
| --- | --- |
| Canonical code `5f6ec56`; HEAD separately recorded | Confirmed at HEAD `9fb7c4e`; source tree identical |
| Data flow: config bootstrap → `Translate` → `Detector` → `DetectSource` → sources → `DetectContext` → `detectRule` → findings → reporter | Confirmed with corrections: commands call `initConfig`, then `initDiagnostics`, then `Config`, then `Detector`; `protect` uses `DetectGit`; config path is overwritten after `Translate`; reporter selection happens inside `Detector`; report path is validated (and an existing file deleted) before scanning |
| “Up to five passes” for `MaxDecodeDepth` | Corrected: five decode iterations and up to six detection passes (raw + decoded) |
| Archive traversal in the source layer | Confirmed: `sources.File` identifies/extracts archives; Git binary archives use `git cat-file` |
| `MaxTargetMegaBytes` check occurs after decoding | Confirmed, with refinements: only when the flag is nonzero; dir source has an opt-in pre-read size check; Git has no pre-read cap; path-only rules run before the check |
| `DetectSource` logs fragment errors and returns nil | Confirmed; CLI exits 1 only if `DetectSource` itself returns an error, which source-level absorption usually prevents |
| Command layer logs scan errors and continues | Confirmed logging, but the summary exits 1 when the error survives; `pipe`/`stdin` use fatal immediately |
| Direct first-party HTTP client | No direct client identified in the inspected first-party code; localhost pprof exists and `extendURL` is a TODO. Dependency- or subprocess-driven network activity is not ruled out. |
| Windows symlink/CRLF/fixture-renaming caveats | Confirmed and carried forward; `detect` tests rename fixtures in place |

---

**Owner-review clarification (2026-11-09):** The execution-path headline now distinguishes direct from indirect `DetectSource` invocation; the comparison invariants distinguish established interfaces from defect candidates; and the absence of an HTTP client is scoped to inspected first-party code. These are documentation clarifications only, not additional runtime verification.

## 8. Known limitations and unresolved questions

1. **No runtime evidence.** Every behavioral consequence above is static; dynamic ordering, races, panics, and memory effects remain unverified. CP1 does not claim measured performance or security outcomes.
2. **No Go toolchain/dependencies.** CP0’s dependency barrier is unchanged; `build`/`test`/`vet`/`-race` remain NOT RUN.
3. **Semaphore semantics unverified.** `semgroup`’s cancellation/error behavior is inferred from usage; the module source was not available offline.
4. **Cobra root/no-subcommand behavior** and exact help exit status are unverified.
5. **RE2/WASM release behavior** is unverified; the default test path uses stdlib regex.
6. **Concurrent `yield` guarantees** depend on the source implementation; the `Source` interface does not document concurrency, and `DetectSource` is safe only because `AddFinding`/`addCommit` lock.
7. **Windows-specific path and symlink behavior** is inferred from code plus CP0’s environment; no Windows runtime test was run.
8. **`viper`/global-state behavior under repeated library use** is source-inferred; the CLI is a single-shot process and is unaffected in normal operation.
9. **No source/backlog changes** were made; `_backlog/` records are cited as leads only.

---

## 9. Working-tree integrity and handoff

- **Before CP1:** clean tree, HEAD `9fb7c4e6bcf6af0b05b281fc46c4a33da8314624`, branch `feat/issues-and-pr` synced with `origin/feat/issues-and-pr`.
- **After CP1:** only the CP1 report and the `_assessment/README.md`/roadmap status updates differ; no Go/test/config/build/CI/backlog file changed; HEAD unchanged.
- **Protected-file changes:** none. No test suite was run in the working checkout, so the CP0 fixture-renaming hazard was not triggered.
- **Commands:** read-only Git and filesystem inspection only; no build/test/generation/network.

**CP2 readiness:** **READY for CP2 with limitations.** CP2 can begin with source-level capability/regression mapping, but executable behavior requires the dependency/toolchain decision already recorded in CP0 (and, for tests, a writable disposable copy because of in-place fixture renames). The highest-value CP2 checks are AR-01/AR-02/AR-04/AR-05/AR-08 and the CLI contract matrix in §2.

---

## Definition-of-done review

- [x] Required investigations and all eight deliverables addressed (or explicitly marked unverified).
- [x] Claims labeled Verified (source/Git) / Inferred / Unverified; no measured runtime claims.
- [x] CP0 cross-check recorded, including corrected assumptions.
- [x] Source/backlog read-only boundary respected; before/after status recorded.
- [x] Runtime verification questions assigned to CP2/CP3/CP4.
- [x] Stopped pending owner authorization; CP2 not started.

**Recommendation:** Ready for owner review.

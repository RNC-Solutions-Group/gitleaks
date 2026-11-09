# CP4 — Security, integration, and developer-tooling assessment

**Status:** CP4A-corrected — awaiting owner review  
**Assessment date:** 2026-11-09 (America/New_York)  
**Repository / branch:** `RNC-Solutions-Group/gitleaks` / `feat/issues-and-pr`  
**Assessed code SHA:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (canonical code baseline)  
**Documentation HEAD at assessment:** `e1c04989b0c3a3566b7832f14f58b1a5c322d21f`  
**Assignment:** [`assignment-CP4.md`](assignment-CP4.md)  
**Review revision:** CP4A (2026-11-09) — corrected Makefile timing and fixture-filler interpretation; see [`review-CP4A.md`](review-CP4A.md).  
**Related:** [`report-CP2.md`](../CP2-capabilities-correctness/report-CP2.md), [`review-CP2A.md`](../CP2-capabilities-correctness/review-CP2A.md), [`report-CP3.md`](../CP3-performance-stability/report-CP3.md)

---

## 1. Executive security/integration assessment

The pinned implementation has a clear first-party security posture in several areas — report writers are single-purpose, template functions are restricted, Git arguments are passed without a shell, archive traversal is disabled by default, and CI actions are pinned by commit SHA. Those controls are real and should be preserved.

The most consequential confirmed weakness is **scan-integrity signaling**. A repository Git failure can produce a clean-looking successful run: a bounded synthetic repository with an invalid `git log` option produced `[git] fatal …` and `stderr is not empty` on stderr, **0 commits scanned**, `[]` findings, and **exit 0** — without any `--exit-code 0` override. The same class of issue exists for unreadable files (warning, partial scan, exit 0). An operator or CI job cannot reliably distinguish “complete, no findings” from “partial/failed” using exit status alone.

Two confirmed **filesystem/data-integrity** risks affect report destinations: report-path validation creates and then removes the target before scanning, so a pre-existing file is destroyed; and if that path is a symlink, validation truncates the symlink target and replaces the link with a regular report file. Both were reproduced with sentinel files in a disposable workspace. No remote exploit is implied.

**Sensitive-data exposure** is conditional on logging. `--redact` protects findings in JSON/CSV/JUnit/SARIF/template/verbose output (measured); it does **not** protect debug/trace logs, where the synthetic token appeared in stderr, nor path-only findings, whose `Match` is corrupted by redaction (CP2-004). Default (non-verbose, info-level) output is not affected by the debug leak, but operators who enable diagnostics can leak secrets into CI logs.

The **untrusted-repository configuration boundary** is broader than many users may expect: a repository-local `.gitleaks.toml` is auto-loaded by default and can define allowlists, rules, and `[extend] path`; a repository-local `.gitleaksignore` can suppress findings; malformed regexes panic instead of producing a controlled error. These are local/CI trust issues, not remote code execution. `.gitleaksignore` loading follows symlinks and has no regular-file/size/capability checks (the #1572 pathway); exploitability of a hostile device/symlink target was not tested and remains unverified. Decoder and archive processing have real resource costs (CP3), but no unbounded input was executed.

Tooling/release findings are moderate: the Makefile’s `VERSION` value contains a literal shell backtick expression; `git fetch --tags` executes only when a recipe that expands `$(LDFLAGS)` runs (`make build`, and `make profile` through its build prerequisite), not during Makefile parsing or `make -n`. Build/test/format/clean/profile targets still have mutating prerequisites or effects; CI uses an unpinned `gotestsum@latest` and mutates generated config in the checkout; the fork’s gitleaks scan workflow is disabled by its upstream-repository condition; the runtime image runs as root with no `.dockerignore`; and the `gore2regex` application binary builds while the full tagged package graph does not.

**Severity summary:** 2 High operational/integrity risks (fail-open scan status; destructive report destination), 4 Medium security/data-integrity risks (redaction/log leakage; untrusted config/ignore policy; baseline self-scan; timeout/partial handling), and a set of Low/Medium build/tooling and contract findings. No Critical remote vulnerability was established; no claim is made that dependency advisories are absent.

---

## 2. Environment and source-provenance manifest

| Item | Value |
| --- | --- |
| Canonical application SHA | `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` |
| Documentation HEAD | `e1c04989b0c3a3566b7832f14f58b1a5c322d21f` |
| Source equivalence | `git diff` over Go/build/config paths from baseline to HEAD is empty |
| Runtime binary | CP3 binary `c1fb9241c1937c65cb31381c767ac229cf56ab9cb559131389daf7f5fdede5eb`; built `CGO_ENABLED=0` from the canonical export; Go 1.27.3 linux/amd64 |
| Container | `golang:1.27.3`, digest `sha256:e519d3b0…`; 2 vCPU/1 GiB/no swap/`--pids-limit=256`/`--network none` |
| Native Go | Go 1.27.3 Windows is installed but was not used for CP4 runtime security checks; no claim of Windows parity |
| Disposable tree | Docker volume `cp3-work` (canonical export); CP4 fixtures in `/work/cp4`; outputs `/out/cp4` outside Git |
| Module/build cache | Reused `cp2-gomod`/`cp2-gocache`; `GOPROXY=off`, `-mod=readonly`; no new dependencies |
| Tests reused | CP2: 275/275 pass, `-race` pass, `go vet` fail on test `lostcancel`, coverage 30.9%; CP3 build/decode profiles |
| CP4 runtime checks | Bounded, synthetic, all under 60 s/1 GiB/2 vCPU; all completed; no OOM/timeout |

**Incomplete:** no Windows runtime security tests, no CI execution, no live advisory lookup, no device-symlink/decompression-bomb tests.

---

## 3. Threat model and trust-boundary map

### 3.1 Untrusted inputs and first-party paths

| Input | Entry point | First-party processing | Trust boundary / failure behavior |
| --- | --- | --- | --- |
| Git history/diff | `git`, `detect`, `protect` | `sources.Git.Fragments`, `git log/diff` subprocess, `gitdiff.Parse` | Repository-controlled data and Git configuration; subprocess stderr is logged; fragment errors are absorbed; Git exit status is debug-logged only |
| Directory contents/files | `dir`, `detect --no-git` | `sources.Files` → `sources.File.fileFragments` | Filesystem names/content; unreadable/oversized/binary files skipped with warnings; symlink following is opt-in |
| Standard input | `stdin`, `detect --pipe` | `sources.File` over `os.Stdin` | Stream is unbounded/blocking; context deadline is checked between reads, not during a blocked read |
| Paths/CLI args | all commands | `cmd`, `sources` | Operator-controlled paths; `--log-opts` is split into Git argv without a shell |
| Symlinks | directory scan | `sources.Files` `EvalSymlinks` | Default skip; opt-in follow can resolve outside the source tree; no containment check |
| Archives | directory/Git archive blobs | `archives.Identify/Extract/OpenReader`; `sources.File` recursion | Depth limit (default 0); no uncompressed-size limit; temp files for Zip/7z; errors mostly logged/skipped |
| Config file | `cmd.initConfig` | Viper/TOML → `config.Translate` | Local `.gitleaks.toml` is auto-loaded; `[extend] path` is a local file read; regex compile panics |
| Rules/regexes | config | `config.Rule`, `regexp.MustCompile` | Config-supplied patterns; malformed regex terminates via panic |
| Baseline | `--baseline-path` | JSON unmarshal into `report.Finding` | Explicit path; content trusted as report data; baseline file itself may be scanned if path comparison misses |
| `.gitleaksignore` | default/source path | `AddGitleaksIgnore` | Repository-controlled; follows symlinks; no regular-file/size/entry cap; `scanner.Err()` unchecked |
| Report destination | `--report-path` | `Detector` validation, `os.Create` in summary | Operator/repo-selected path; existing target is deleted/truncated before scan; symlinks are followed |
| Environment | `GITLEAKS_CONFIG`, `GITLEAKS_CONFIG_TOML`, `GITLEAKS_*` | `initConfig` | Process environment is trusted; content is parsed as config |

### 3.2 Protected assets and principal boundaries

| Asset | Boundary | Existing protection | Residual concern |
| --- | --- | --- | --- |
| Secret material | finding → console/report/log | `--redact` in `filter` before `AddFinding`; report writers consume redacted findings | Debug/trace logs occur before redaction; path-only finding corruption; default output is unredacted unless requested |
| Scan completeness | source → detector → CLI status | Source errors are logged; context cancellation returns error | Fragment/worker errors are absorbed; partial scans can exit 0 |
| Report correctness | finding model → writers | Typed writers; template funcs restricted; JSON/CSV/JUnit/SARIF validated in CP2 | Report target can be deleted/truncated before success; empty CSV has no header |
| Filesystem integrity | CLI/report path | Explicit path; `os.Create`/`os.Remove` | Validation deletes existing files; symlink target truncation; TOCTOU window |
| Availability | inputs/decoder/archive | `max-decode-depth`, `max-archive-depth`, optional `max-target-megabytes` | Default max-target 0; decoder regex cost; no uncompressed archive cap; `.gitleaksignore` blocking read |
| CI reliability | CLI exit/status | `--exit-code`, fatal config/report errors | Git/worker errors can look successful; gitleaks workflow disabled on forks |
| Build/release provenance | CI/Docker/goreleaser | Actions pinned by SHA; module checksums; CGO=0 release Docker build | Toolchain drift; tagged package-graph failure; no SBOM/signing observed; root runtime image |

---

## 4. Scan-completeness and error-propagation analysis

### 4.1 Propagation paths

| Stage | Error source | Handling | Caller-visible outcome |
| --- | --- | --- | --- |
| Git subprocess stderr | `listenForStdErr` | Logs each non-ignored line; sends one generic `stderr is not empty` error | `Git.Fragments` returns immediately; `DetectSource` callback logs and returns nil; CLI sees no error |
| Git process exit | `defer Cmd.Wait()` | Debug log only | Not surfaced |
| Git fragment/archive read | per-file worker / `file.Fragments` | Logged; `yield` error swallowed by detector callback | Partial scan; exit based only on findings |
| Directory worker | `s.Sema.Go` / local `wg` | Worker errors are not collected; only walk error returned | Partial scan; exit 0 unless a finding exists |
| File open/read | `sources.Files` | Permission/read errors skip the file with a warning; `yield` error swallowed | Partial scan; exit 0 |
| Archive extraction | `sources.File` | Logs and skips on errors; inner errors returned then swallowed | Partial scan |
| Cancellation | source `ctx.Done`/`fileFragments` | `Fragments` returns `ctx.Err()` | `DetectSource` returns error; CLI exits 1 |
| Config/report validation | `cmd`/`Detector` | `logging.Fatal` | Exit 1 (or panic 2 for regex) |
| Report write | `Reporter.Write` | `logging.Fatal` | Exit 1 |
| Interrupt | `main` goroutine | `logging.Fatal` | Immediate exit 1; no report finalization |

### 4.2 CP4 runtime confirmation (no `--exit-code 0`)

- **Valid Git scan:** synthetic repository, one finding → exit **1**; `leaks found: 1`.
- **Invalid Git option:** `gitleaks git … --log-opts "--not-a-real-option"` → stderr `[git] fatal: unrecognized argument`, `error="stderr is not empty"`, `0 commits scanned`, stdout `[]`, **exit 0**, “no leaks found.”
- **Clean directory:** exit 0.
- **Unreadable file (CP2 non-root):** warning, remaining file reported, exit 0.
- **Delayed stdin (CP2):** 1 s timeout observed after 2,009 ms, then fatal `context deadline exceeded`, exit 1.
- **Report failure:** missing template → exit 1 (after target deletion).

### 4.3 Operator/CI distinguishability

| State | Observable evidence | Distinguishable from “complete, clean”? |
| --- | --- | --- |
| Complete, no findings | exit 0; `no leaks found` | Baseline |
| Complete, findings | exit 1 (default) or configured `--exit-code`; `leaks found: N` | Yes, by exit/finding log |
| Partial (Git stderr) | Git error + `stderr is not empty` logged; `0 commits scanned`; exit **0** | **No** from exit alone; requires log inspection |
| Partial (unreadable file) | Warning `skipping file: permission denied`; exit 0 | **No** from exit alone |
| Failed config/report | fatal log; exit 1 (or panic 2) | Yes, but failure type not machine-readable |
| Timed out | delayed fatal `context deadline exceeded`; exit 1 | Only after the blocking operation returns; no partial-status code |

**Highest integrity risk:** CI gates that rely only on exit status can be green after a partial Git scan. This is a confirmed behavior, not a remote exploit.

---

## 5. Sensitive-data exposure matrix

| Output / boundary | Redaction behavior | CP4 evidence | Status |
| --- | --- | --- | --- |
| JSON report | Finding fields redacted by `filter` before write | `redact.json` has 1 finding; raw synthetic token absent | Protected when `--redact` |
| CSV report | Same | `redact.csv` row present; raw token absent | Protected when `--redact` |
| JUnit report | Same; failure data JSON uses redacted finding | `redact.junit` valid XML; raw token absent | Protected when `--redact` |
| SARIF report | Same; snippet uses redacted `Secret` | `redact.sarif` 1 result; raw token absent | Protected when `--redact` |
| Template report | Template receives redacted findings | `redact.txt` renders rule/file; raw token absent | Protected when `--redact` |
| Verbose stdout | `printFinding` receives findings after `filter` | `--verbose --redact` output contains no raw token | Protected when `--redact` |
| Debug/trace logs | Redaction happens after detection; decoder/allowlist logs occur earlier | Debug decoder log contains synthetic token; CP2 trace `gitleaks:allow` log contains it | **Not protected** |
| Git subprocess stderr | Logged with `[git] %s`; no redaction | Source [`sources/git.go:226-271`](../../../sources/git.go#L226-L271) | **Not protected**; content is subprocess-controlled |
| `gitleaks:allow` trace | Secret printed in trace before skip | CP2 measured; [`detect/detect.go:510-515`](../../../detect/detect.go#L510-L515) | **Not protected** at trace |
| Path-only finding redaction | `Secret` empty; 100% redaction corrupts `Match` | CP2-004 reproduced | Output integrity defect, no secret exposed |
| Diagnostics/profiles | Profiles contain stack/heap metadata; raw string exposure not verified | CP3 profiles were not inspected for plaintext secret strings | **Unverified** |
| Default non-verbose output | No findings printed unless `--verbose`; report still contains secrets unless `--redact` | By source/CP2 | Expected behavior; documentation should be explicit |

**Redaction boundary conclusion:** `--redact` protects finding/report/verbose output once findings exist, but it is applied too late to protect diagnostic logs emitted during detection/decoding. It also does not protect Git stderr. This matches the CLI flag description (“redact secrets from logs and stdout”) only partially: stdout/report findings are protected; diagnostic logs are not.

---

## 6. Filesystem, symlink, ignore-file, and archive analysis

### 6.1 Report destinations (confirmed)

- Validation at [`cmd/root.go:347-358`](../../../cmd/root.go#L347-L358) runs `os.Create(reportPath)` and then `os.Remove(reportPath)` for every non-`-` path, before the scan.
- **Measured:** a sentinel file `/out/cp4/target-delete.json` was absent after a failing template load (exit 1).
- **Measured symlink case:** `symlink-target.json` (sentinel) + `symlink-report.json -> symlink-target.json`; after `dir --report-path symlink-report.json`, the target was **0 bytes** and the path had become a regular 444-byte report file (the symlink was removed).
- **Impact:** local data-loss/truncation when the report path aliases an existing file; a pre-existing report is destroyed even if the scan later fails. No remote vector; requires operator/repo-selected path.
- **TOCTOU:** validation and final `os.Create` are separate operations; the path can change in between. Not reproduced.

### 6.2 Baseline and config path self-scan (confirmed)

- Baseline path is made relative to source ([`detect/baseline.go:52-79`](../../../detect/baseline.go#L52-L79)); skip is strict string equality ([`detect/detect.go:294-299`](../../../detect/detect.go#L294-L299)).
- **Measured:** baseline generated inside the source; second scan produced **2 findings from `baseline.json`** while suppressing the original token finding. The baseline JSON contains its own stored secrets and is re-scanned.
- Config `.gitleaks.toml` is partially masked by the default path allowlist (`gitleaks\.toml`); baseline JSON is not.

### 6.3 Symlink and filesystem boundaries

- Directory symlinks are skipped by default; `--follow-symlinks` calls `EvalSymlinks` and scans the real target ([`sources/files.go:103-121`](../../../sources/files.go#L103-L121)). There is **no containment check** to keep the target inside the source root; a repo-controlled symlink plus operator-enabled follow can read local files outside the tree. This is a local trust/configuration issue; no external scan was run.
- File permission failures are warnings and partial success ([`sources/files.go:155-162`](../../../sources/files.go#L155-L162)); CP2 non-root reproduction.
- Archive extraction uses in-memory readers and a temp file for Zip/7z ([`sources/file.go:90-147`](../../../sources/file.go#L90-L147)); inner paths are cleaned for scanning/allowlist display, not extracted to arbitrary disk paths. Depth default 0.
- No `.dockerignore`; Docker build context includes `.git` and assessment docs. The runtime image copies only the binary, but build context/provenance is broader than necessary.

### 6.4 `.gitleaksignore` (#1572 pathway)

- `fileExists` uses `os.Stat` (follows symlinks) ([`cmd/root.go:502-515`](../../../cmd/root.go#L502-L515)).
- `AddGitleaksIgnore` uses `os.Open` + `bufio.Scanner`, stores every nonblank/noncomment line, warns on malformed entries, and returns `nil` without checking `scanner.Err()` ([`detect/detect.go:153-193`](../../../detect/detect.go#L153-L193)).
- **Measured finite control:** `.gitleaksignore` as a symlink to a regular finite file suppressed the synthetic finding (0 findings); removing the symlink produced 1. This confirms symlink following and suppression behavior without a device target.
- **Missing controls:** no regular-file policy, no aggregate byte/entry cap, no blocking-read timeout, no containment; `bufio.Scanner`’s default token limit applies but its error is ignored. A hostile repository-controlled ignore path/symlink could block or consume resources; exploitability of the original #1572 device case was **not tested** and remains unverified.

### 6.5 Archive/decoding resource controls

- Archive depth defaults to 0; depth 1 extracts and scans bounded content. CP3 measured ~1.7× ordinary cost and 2,378 B compressed / 103,800 B uncompressed for 100 files.
- No uncompressed-size or per-entry cap is enforced by `sources.File`; the optional `max-target-megabytes` applies to file reads/detector content, not an explicit decompressed archive budget. Allocation occurs before any post-decode size check ([`detect/detect.go:431-440`](../../../detect/detect.go#L431-L440)); CP3 measured decoder allocation and CPU.
- The CP3 28.7× depth-0 result is a **coverage-changing control**, not a safe optimization. Any decoder gating must preserve percent, Unicode, hex, Base64, and nested-encoding coverage.

---

## 7. Configuration and rule-processing assessment

- **Precedence:** `--config` > `GITLEAKS_CONFIG` > `GITLEAKS_CONFIG_TOML` > `<source>/.gitleaks.toml` > embedded default ([`cmd/root.go:132-195`](../../../cmd/root.go#L132-L195)). CP2 measured all explicit/env/source-local paths.
- **Untrusted repository config:** when no explicit config is passed, a repository-local `.gitleaks.toml` is auto-loaded. It can add allowlists (suppressing findings), change rules, and set `[extend] path`. This is a scan-policy trust boundary for scanning untrusted repositories; `--config` is the available workaround. `.gitleaksignore` is also auto-loaded from the source directory.
- **Extend recursion/global state:** `extendDepth` is a package global ([`config/config.go:24`](../../../config/config.go#L24)) never reset in production; repeated `Translate` calls in one process skip validation/extension at depth ≥2. **Measured CP4 probe:** first translation added the default config (223 rules); a second translation of an invalid rule (no regex/path) returned **no error and accepted the rule**. CLI single-shot use is unaffected; library/embedded/repeated-init use is affected.
- **Malformed regex:** `regexp.MustCompile` panics ([`config/config.go:123-128`](../../../config/config.go#L123-L128)); CP2 measured exit 2 and a stack trace. A repository-local malformed config can crash a scan.
- **`extend.path`:** `viper.SetConfigFile` reads the supplied path; relative paths follow process/Viper CWD semantics. A repository config can point at local files; parsing only, but errors/paths may be exposed.
- **Allowlists/baseline:** global/rule path/commit/regex/stopword allowlists and baselines are applied at different stages; commit allowlist source logging is misleading but detector-level filtering prevents a false negative (CP2-014).
- **Rule selection:** `--enable-rule` replaces the rule map after prefilter construction; correctness preserved (superset prefilter), minor wasted work.

---

## 8. Git, hooks, and CI integration assessment

### 8.1 Exit/status behavior

| Scenario | Exit | CI signal |
| --- | --- | --- |
| Findings | 1 default / `--exit-code` | Red |
| Clean complete | 0 | Green |
| Git stderr/subprocess failure | **0** | **Green despite partial/failed Git scan** |
| Unreadable file | **0** | **Green despite skipped file** |
| Config/report fatal | 1 | Red (but no machine-readable class) |
| Invalid regex | 2 (panic) | Red |
| Timeout | 1 after blocked operation returns | Red, delayed |

### 8.2 Git and hooks

- `--log-opts` is split on spaces and appended to Git argv; no shell is used, so shell injection is not present. Arbitrary Git options remain operator-controlled and may have Git-defined filesystem/network effects; not exhaustively enumerated.
- `git --pre-commit` scans the **unstaged** diff; `--staged` scans the staged diff. CP2 reproduced 0 vs 1 finding for a staged-only change. The pre-commit hook entries pass `--staged` ([`.pre-commit-hooks.yaml:4`](../../../.pre-commit-hooks.yaml#L4)), which matches the intended staged workflow.
- `protect` is hidden/deprecated and routes through `DetectGit` with `context.Background()` ([`cmd/protect.go:61`](../../../cmd/protect.go#L61), [`detect/git.go:19-29`](../../../detect/git.go#L19-L29)); the Git subprocess is still created with the command context, but detector-level cancellation is absent.
- `scripts/pre-commit.py` runs a fixed `gitleaks protect -v --staged` command via `os.system`; no user input is interpolated. It checks only exit code 1; other nonzero statuses print nothing and `sys.exit` is not called (falls through), which can produce a misleading hook result for config/runtime failures. Source-verified.

### 8.3 CI workflows

- `test.yml` builds with Go 1.25, installs **`gotestsum@latest`** (unpinned/network), runs `go test -json ./... --race`, then `go generate ./... && git diff --exit-code`. It does not run `go vet`. `go generate` mutates the checkout (validates generated config) and depends on the generator toolchain.
- `gitleaks.yml` has `if: github.repository == 'gitleaks/gitleaks'`; on this fork it is a no-op, so the branch’s own repository does not run the Gitleaks scan workflow. A CI pipeline expecting that workflow to protect the fork would be misled.
- `release.yml` pins actions by SHA, uses Docker Hub/GHCR secrets, builds multi-arch images with build-push; it does not visibly run tests or produce an SBOM/signature/provenance manifest in this file.
- The Dockerfile runtime is Alpine with bash/git/openssh-client, runs as **root**, and sets `safe.directory '*'`; no `USER` or capability drop. The build stage uses Go 1.24 and `CGO_ENABLED=0`.
- Actions are pinned to commit SHAs (positive supply-chain control); no live workflow was triggered.

---

## 9. Makefile and script side-effect matrix

Sources: [`Makefile`](../../../Makefile#L1), [`scripts/profile.sh`](../../../scripts/profile.sh#L1), [`scripts/pre-commit.py`](../../../scripts/pre-commit.py#L1), [`.pre-commit-hooks.yaml`](../../../.pre-commit-hooks.yaml#L1).

| Target | Prerequisites | Commands | Filesystem/network/Git effects | Diagnostic suitability |
| --- | --- | --- | --- | --- |
| `VERSION` / `LDFLAGS` | — | `VERSION := \`git fetch --tags && git tag \| sort -V \| tail -1\``; `LDFLAGS` embeds `$(VERSION)` | **No execution at parse time.** Backticks are literal Make text; the shell executes them only when a recipe expands `$(LDFLAGS)` (`build`, `profile`). `make -n` prints without shell execution. Network + `.git` tag mutation applies to build/profile, not every Make invocation | Not safe for build/profile; parse-only inspection is not the risk |
| `config/gitleaks.toml` | wildcard generator sources | `go generate ./...` | Rewrites tracked generated config; toolchain/deps | Not safe |
| `format` | — | `go fmt ./...` | Rewrites all Go source | Not safe |
| `test-cover` | — | `go test -v ./... --race --cover --coverprofile=cover.out <pkg>`; `go tool cover -html` | Writes `cover.out`, opens browser; race needs CGO | Not safe as-written |
| `test` | `config/gitleaks.toml`, `format` | generate + fmt + `go test -v ./... --race <pkg>` | Rewrites source/config, writes test artifacts | Not safe |
| `failfast` | `format` | `go test -failfast ./...` | Rewrites source | Not safe |
| `build` | `config/gitleaks.toml`, `format` | `go mod tidy`, `go build $(LDFLAGS)` | Rewrites go.mod/go.sum, source/config, creates binary; **shell executes `git fetch --tags` via LDFLAGS**; `go mod tidy` may need network if cache is cold | Not safe |
| `lint` | — | `golangci-lint run` | Tool may write cache; no tag fetch; tool absent locally | Potentially observational only if tool/cache controlled |
| `clean` | — | `rm -rf profile`; delete `*.got.*`, `*.out` | Destructive deletion; no network/tag fetch | Not safe |
| `profile` | `build` | build + `scripts/profile.sh './gitleaks' '.'` | Inherits build network/mutation; writes `profile/<timestamp>` in CWD; runs scanner | Not safe |
| `test-cover`/`test` package argument | — | `$(PKG)` = `github.com/zricethezav/gitleaks` (no `/v8`) | May select the wrong/no package; unverified | Investigate before future use |

`profile.sh` sets `set -euo pipefail`, uses `--exit-code=0 --max-decode-depth 8`, writes `profile/<epoch>` and profiles under CWD ([`scripts/profile.sh:20-24`](../../../scripts/profile.sh#L20-L24), [`:38`](../../../scripts/profile.sh#L38), [`:66`](../../../scripts/profile.sh#L66)). It is a mutating diagnostic, not a read-only benchmark harness.

**Recommendation:** separate future read-only diagnostics (direct `go test`, explicit external outputs, no fetch) from maintenance targets; remove or make explicit the build/profile `git fetch --tags` dependency and require explicit version input. No Makefile change was made.

---

## 10. Build, dependency, and release reproducibility matrix

| Aspect | Observed | Assessment |
| --- | --- | --- |
| Module requirement | `go 1.24.11` ([`go.mod:3`](../../../go.mod#L3)) | Patch-level requirement |
| CI toolchain | Go 1.25 ([`test.yml:23`](../../../.github/workflows/test.yml#L23)) | Drift from module/Dockerfile |
| Docker build | `golang:1.24` ([`Dockerfile:1`](../../../Dockerfile#L1)); runtime Alpine 3.22 | Older build toolchain; CGO=0 |
| CP2/CP3 toolchain | Go 1.27.3 Linux, CGO 1/0 respectively | Not CI-equivalent |
| Application build | default and `gore2regex` app binaries build; `-tags gore2regex .` exit 0 | Working application binary |
| Complete tagged graph | `go build -tags gore2regex ./...` and `go test -tags gore2regex ./...` fail in `airtable.go`/`perplexity.go` | Generator packages incompatible with tag; not a broken release binary |
| Default tests | 275/275 pass; `-race` pass; shuffle pass | CP2 measured |
| `go vet` | Fails on two test `lostcancel` diagnostics | Static gate absent in CI; unresolved |
| Coverage | 30.9% instrumented; `cmd` 0%, `sources` 7.4% | Whole-program coverage lower |
| Release tag | `.goreleaser` `gore2regex`, `-s -w`, version ldflags, prerelease | Release packaging not executed |
| Version metadata | Makefile `build`/`profile` fetch tags via shell expansion of `LDFLAGS`; Dockerfile `git describe` | Network/Git-dependent version source for those targets |
| Runtime image | Root, bash/git/openssh, `safe.directory=*` | Broader attack surface than needed |
| Build context | No `.dockerignore` | `.git`/docs included in context |
| Dependency advisories | No live lookup performed | **Not claimed absent**; follow-up CP5/security process |
| SBOM/signing | None observed in release workflow | Provenance gap; CP5 input |

---

## 11. Confirmed findings and risk register

| ID | Finding | Class | Severity | Confidence | Evidence | Existing mitigation | Exploitability/limits |
| --- | --- | --- | --- | --- | --- | --- | --- |
| CP4-001 | Git subprocess/stderr failure produces exit 0 with `[]` findings | Scan-integrity | **High** | High (Measured) | `git-fail-open`; `0 commits`, Git fatal, exit 0 | Error logged; no exit-status mitigation | Local/CI scan with malformed Git option or Git stderr failure; not RCE |
| CP4-002 | Report target is deleted/truncated before scan; symlink target is truncated | Data-integrity | **High** | High (Measured) | `report-delete` target absent; `report-symlink` target 0 B, link replaced | Explicit path; no transaction/atomic rename | Local path control; no remote vector |
| CP4-003 | `--redact` does not cover debug/trace logs (decoder, `gitleaks:allow`, Git stderr) | Sensitive-data exposure | **Medium** | High (Measured/source) | `debug-leak` raw synthetic token in stderr; CP2 trace leak | Redaction for findings/reports/verbose | Requires elevated log level; logs may reach CI storage |
| CP4-004 | Repository-local `.gitleaks.toml` is auto-loaded and can suppress/alter scanning; malformed regex panics | Policy/availability | **Medium** | High (Measured/source) | CP2 precedence; CP2-008 panic; config trust path | `--config` override; validation for some errors | Scanning untrusted repo/local config; no remote code execution |
| CP4-005 | Baseline JSON inside source is rescanned and its stored secrets can be reported | Data exposure | **Medium** | High (Measured) | `baseline-second`: 2 findings in `baseline.json` | Explicit baseline path; default path allowlist for `gitleaks.toml` only | Requires baseline inside scanned source and path-form mismatch |
| CP4-006 | `.gitleaksignore` follows symlinks and has no regular-file/size/entry cap; `scanner.Err()` ignored | DoS/trust boundary (#1572) | **Medium–High** | Source High; exploit Unverified | Finite symlink suppression measured; source [`detect/detect.go:153-193`](../../../detect/detect.go#L153-L193) | None specific; ignore path operator-controlled | Hostile repo-controlled ignore path; device/blocking case not tested |
| CP4-007 | Archive processing has no explicit uncompressed-size budget; allocation precedes post-decode checks | Resource DoS | **Medium** | Source High; scale Unverified | CP3 archive/decoder measurements; source | `max-archive-depth` default 0; max target optional | Compressed-bomb path not tested; bounded evidence only |
| CP4-008 | Unreadable file yields warning + exit 0 with partial findings | Scan-integrity | **Medium** | High (Measured non-root) | CP2 `perm-nonroot` | Warning log | Local permission failure; no partial status |
| CP4-009 | Timeout does not interrupt a blocked stdin read promptly | Availability/integrity | **Medium** | High (Measured) | 2,009 ms for 1 s timeout; fatal after read | Context deadline eventually surfaces | Operator/CI delay; no unbounded test |
| CP4-010 | `extendDepth` global causes repeated translation to skip validation; invalid rule accepted | Library correctness | **Low–Medium** | High (Measured probe) | Second invalid config accepted; rules=223→1 | CLI single-shot; tests reset manually | Embedded/library/repeated initialization; not normal CLI |
| CP4-011 | Path-only `--redact` corrupts `Match` text | Output correctness | **Low–Medium** | High (Measured CP2) | CP2-004 repeated `REDACTED` | Redaction still removes secret | Cosmetic/data-quality; no secret exposure |
| CP4-012 | Makefile `VERSION` backticks execute `git fetch --tags` when `build`/`profile` expand `LDFLAGS`; build/test/format/clean/profile have mutating prerequisites/effects | Tooling/reproducibility | **Medium** | High (source) | Makefile matrix; no make execution | None; documentation warns | Developer/CI invoking build/profile (network/tag mutation); mutating test/format/build/clean targets |
| CP4-013 | Tagged full package graph fails; toolchain drift; vet failure | Build/release | **Medium** | High (Measured CP2) | CP2 `go test -tags gore2regex ./...` failure; `go vet` | App binary builds; CI actions pinned | Release/generator/CI workflows; not end-user binary |
| CP4-014 | Fork gitleaks CI workflow is disabled; test workflow uses unpinned `gotestsum@latest` and mutates generated config | CI integration | **Medium** | High (source) | `gitleaks.yml:9`; `test.yml:30,37` | Pinned actions; no gitleaks scan on fork | Misleading security gate; supply-chain drift |
| CP4-015 | Runtime container runs as root with broad tools and `safe.directory=*`; no `.dockerignore` | Container hardening | **Low–Medium** | Source | `Dockerfile:7-11`; no `.dockerignore` | Minimal binary copy; Alpine | Container compromise impact; not remotely reachable by scanner input |
| CP4-016 | Dependency advisories not checked | Supply chain | **Unverified** | — | No advisory lookup performed | go.sum checksums | Cannot claim absence; CP5/security process |
| CP4-017 | SARIF/report output includes secrets unless `--redact` | Contract | C4 | High (source/CP2) | CP2 SARIF snippet; redacted run | `--redact` | Expected reporting behavior; policy decision |
| CP4-018 | Template functions restricted; Git argv has no shell; archive default off; actions pinned | Positive controls | — | High (source/measured) | `report/template.go:33-37`; `sources/git.go`; `.goreleaser` | — | Preserve in future changes |

**Severity rationale:** High findings are confirmed, low-precondition operational integrity/confidentiality issues with clear CI/operator impact; no remote unauthenticated code execution was demonstrated. Medium findings require local/untrusted-input preconditions or elevated logging. Unverified items are not assigned severity as defects.

---

## 12. Mitigations and existing controls

- `--redact` (opt-in) for findings, reports, and verbose output; `--exit-code`; fatal report-write/config errors.
- Default archive depth 0; opt-in size/decode limits; binary detection; symlink skipping by default; archive depth bound.
- Template function map excludes `env`, `expandenv`, `getHostByName` ([`report/template.go:33-37`](../../../report/template.go#L33-L37)).
- Git subprocess uses `exec.CommandContext` with argv, no shell ([`sources/git.go:91-93`](../../../sources/git.go#L91-L93)); no command injection via `--log-opts`.
- `.gitleaksignore` finite file suppression; baseline suppression; allowlists.
- GitHub Actions pinned by commit SHA; module checksums in `go.sum`; CGO-disabled release Docker build.
- pprof diagnostics server binds `localhost` only ([`cmd/diagnostics.go:136-147`](../../../cmd/diagnostics.go#L136-L147)).

These controls should not be read as compensating for CP4-001/002/003/004/006.

---

## 13. Unverified hypotheses and limitations

1. `.gitleaksignore` device/blocking-read DoS (#1572) was not tested; only source behavior and a finite symlink control were established.
2. Archive decompression-bomb behavior and uncompressed-size exhaustion were not tested; CP3 used bounded archives only.
3. `--follow-symlinks` outside-source read was not reproduced; source shows no containment.
4. Debug/trace logs are confirmed to leak synthetic content; real-log exposure depends on logging configuration and retention.
5. Profiles/diagnostics may contain sensitive strings; not verified.
6. The report-target symlink behavior was measured locally; TOCTOU race and cross-platform behavior were not.
7. Windows-native behavior, CI Go 1.25, and release packaging were not executed; native Go is installed but unused.
8. Dependency vulnerabilities were not checked against live advisories; absence is not claimed.
9. `protect` end-to-end timeout and `git --pre-commit` behavior beyond CP2 were not fully explored.
10. Makefile targets were inspected only; no `make` or `make -n` was run.

---

## 14. CP0–CP3 crosswalk

| Prior finding | CP4 disposition |
| --- | --- |
| CP0 environment/toolchain, no native Go at CP0 | Confirmed historical; native Go now installed after CP0; CP4 used Linux container; no Windows tests |
| CP1 fail-open error flow (C-01/C-02, AR-01) | **Confirmed** by CP4 `git-fail-open` without exit-code override |
| CP1 report-target deletion (AR-13) | **Confirmed**, plus symlink target truncation |
| CP1 path-only redaction (AR-04) | **Confirmed** (CP2-004) |
| CP1 redaction logs (AR-03) | **Confirmed** debug/trace leak |
| CP1 baseline path (AR-05) | **Confirmed** baseline self-scan |
| CP1 `.gitleaksignore` (AR-11) | Source-confirmed; finite symlink control measured; device exploit unverified |
| CP1 `extendDepth` (AR-06) | **Measured** invalid rule accepted after repeated translation |
| CP1 archive/decoder resource (AR-09/10) | CP3 measured; no unbounded test |
| CP2-001/002/012/013/015 contract items | CP4 treats as contract decisions; timeout/exit behavior included |
| CP2-005/006/007/008/009/014 | Confirmed/refined; see findings register |
| CP2-010/011 | Confirmed static/build findings; app binary works |
| CP3 28.7× decoder control | Used as performance evidence; explicitly **not** a safe optimization (coverage change) |

---

## 15. Prioritized future verification/remediation candidates

| Priority | Candidate | Evidence | Recommendation |
| ---: | --- | --- | --- |
| 1 | Fail-closed/partial scan status | CP4-001/008 | Future design: nonzero or explicit partial status on Git/worker/file errors; add end-to-end tests |
| 2 | Atomic/symlink-safe report handling | CP4-002 | Future: validate without deleting; write temp + rename; `O_NOFOLLOW`-equivalent; preserve on failure |
| 3 | Redact or suppress diagnostic secrets | CP4-003 | Future: move redaction before logging or gate secret-bearing logs; test debug/trace/CI logs |
| 4 | Harden repository-local config/ignore trust | CP4-004/006 | Future: opt-out of local config/ignore for untrusted repos; regular-file/size/timeouts for ignore; recover regex panics |
| 5 | Explicit archive/decode resource budgets | CP4-007; CP3 | Future: uncompressed archive cap, decode gating with full coverage tests |
| 6 | Baseline self-scan exclusion | CP4-005 | Future: normalize/realpath-aware exclusion; test inside/outside source |
| 7 | Tooling/CI/release hygiene | CP4-012/013/014/015 | Future: make build/profile tag fetch explicit or avoid it; pin gotestsum; SBOM/signing; non-root image; tagged graph tests; vet |
| 8 | Dependency advisory process | CP4-016 | Future scheduled scan; do not claim absence |
| 9 | Report format contract decisions | CP4-001/017; CP2-001/002/015 | Document exit/report/empty-CSV/SARIF semantics before code changes |

No implementation was performed.

---

## 16. CP5 consolidation handoff

CP5 should consolidate:

- CP4-001/008/009 into a single **scan-completeness/status** finding with recommended evidence expectations for CI.
- CP4-002 and CP4-005 into **filesystem/data-integrity** controls (report destination and baseline self-scan).
- CP4-003/011 into a **sensitive-output/redaction boundary** matrix.
- CP4-004/006/010 into a **configuration/trust and global-state** finding.
- CP4-007 with CP3 decoder/archive data as a **resource-control** section.
- CP4-012–015/016 as a **build/CI/release reproducibility** section.

Unresolved: #1572 exploitability, archive-bomb behavior, profile secret content, Windows/CI parity, dependency advisories. These should remain explicit unknowns, not pass/fail claims.

---

## 17. Repository-integrity statement

- Tracked checkout: canonical source/test/fixture/config/build/CI/`_backlog` unchanged; only `_assessment/**` documentation added/edited.
- Disposable execution: `/work/cp4` + `/work/src` in `cp3-work`; CP3 binary and CP2/CP3 evidence reused; outputs/profiles external; no source edits in the tracked repository.
- Network: none for CP4 runtime checks; no image pull, dependency download, or external repository scan.
- Limits: 2 vCPU/1 GiB/no swap/60 s per workload; all checks completed; no OOM/timeout.
- Secrets: synthetic `CP4TOKEN-NNNN` fixtures only; no raw secret-like fixture values recorded.
- Final status and commit/push are reported in the CP4 handoff response after documentation validation.

**Recommendation:** Ready for owner review. CP5 not started.

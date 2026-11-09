# CP2 — Capabilities, correctness, and regression baseline

**Status:** CP2A review complete — accepted with corrections  
**Assessment date:** 2026-11-09 (America/New_York)  
**Repository / branch:** `RNC-Solutions-Group/gitleaks` / `feat/issues-and-pr`  
**Assessed code SHA:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (canonical code baseline)  
**Branch HEAD at assessment:** `c39f9a84339719c7c9d5602ca5caa910bd173e89` (documentation-only)  
**Assignment:** [`assignment-CP2.md`](assignment-CP2.md)  
**Review revision:** CP2A (2026-11-09) — classification corrections and evidence hardening; see [`review-CP2A.md`](review-CP2A.md).  
**Test environment:** disposable Linux tree in Docker volume `cp2-work`; `golang:1.27.3` (local image digest `sha256:e519d3b0…`), Go `go1.27.3 linux/amd64`, `CGO_ENABLED=1`, Git 2.47.3, gcc 14.2.0; module requirement `go 1.24.11`. Native Go 1.27.3 for Windows is now installed (`C:\Program Files\Go`, created 2026-11-09 14:50) but was not used for CP2 measurements and no Windows tests were run.

---

## Executive assessment

The pinned code builds, tests, and runs under the approved Linux container toolchain. The default suite passed **275/275 tests with 0 failures and 0 skips**; `-race` passed; `-shuffle=on` passed; `gofmt -l` was clean on the LF export. `go vet ./...` fails on two test-only `lostcancel` diagnostics. Coverage across instrumented packages is **30.9%** of statements, with `sources` at 7.4% and the `cmd` package at 0% (no tests).

The executable baseline is healthy in its core paths but contains a mix of confirmed defects, security/data-integrity risks, static/build failures, and contract questions. CP2A classifies them separately; this report’s measured results are unchanged.

- **Path-only findings are corrupted by `--redact`**: `Match` becomes `REDACTED` repeated between every character. Reproduced.
- **A Git subprocess failure returns success**: an invalid `git log` option produced `[git] fatal …`, `stderr is not empty`, 0 commits scanned, `[]` findings, and **exit 0** with “no leaks found.” Reproduced.
- **An unreadable file produces a warning and a successful partial scan**: reproduced as a non-root user; the remaining file was still reported and the CLI exited 0.
- **Baseline JSON inside the source is scanned**: the original finding was suppressed, but the baseline file itself produced two findings. Reproduced.
- **Report-target validation deletes an existing file before scanning**; reproduced by deleting a sentinel during a failing template load.
- **`--report-path -` requires an explicit `--report-format`**; otherwise it fatals with `Unknown report format:`.
- **`--report-template` does not imply template format**, despite its help text; reproduced fatal.
- **Invalid regexes panic** (`exit 2`) instead of returning a controlled error.
- **`--redact` does not redact debug/trace logs**; the synthetic token appeared in stderr in both modes.
- **`gore2regex` tagged builds/tests fail** in the generator rule package; the release binary itself builds and loads the default 222-rule config.
- **Finding order is nondeterministic** across runs because rules are traversed from a map; eight runs produced alternating rule orders.

Positive results: all scanning modes (`git`, `dir`, `stdin`, `detect`, `protect`) produced expected findings; config precedence (`--config`, `GITLEAKS_CONFIG`, `GITLEAKS_CONFIG_TOML`, source-local), allowlists, `gitleaks:allow`, `.gitleaksignore`, baseline suppression, entropy gates, generic dedupe, decode depth, archive depth, size limits, symlink following, report writers, and exit-code overrides all behaved as expected in controlled cases. Release-tag compilation and default-config loading also passed.

---

## Scope and methodology

**In scope:** CP2’s ten investigations — capability/CLI matrix, automated baseline, mode parity, error/completion integrity, configuration/path correctness, detection/filtering, reporting/redaction, platform/build variants, regression mapping, and handoff.

**Method:** the pinned tree was exported with container `git archive` to a writable Linux Docker volume; dependencies were acquired once with network enabled into an external module-cache volume and verified against `go.sum`; all builds, tests, vet, coverage, race, and CLI runs were then executed with `--network none`, `GOPROXY=off`, `GOFLAGS=-mod=readonly`, and external caches. The tracked assessment checkout was never used as an execution tree. All CLI fixtures are synthetic (`CP2TOKEN-NNNN`, obviously not credentials); no external repositories, real credentials, unbounded workloads, device symlinks, or #2019/#1572 reproductions were used.

**Evidence labels:** **Verified (source)** from CP1; **Measured** from executed commands; **Inferred** where runtime evidence is indirect; **Unverified/NOT RUN** where blocked.

**CP1/CP0 inputs:** CP0 environment findings and CP1’s AR register were read first. CP1’s source-only claims were tested rather than assumed; the CP1-004/AR-08 commit-allowlist log behavior and CP1-024 map-order nondeterminism were both reproduced.

---

## 1. Test-environment manifest

| Item | Value |
| --- | --- |
| Canonical code SHA | `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` |
| Branch HEAD | `c39f9a84339719c7c9d5602ca5caa910bd173e89`; source tree identical to baseline (`git diff` over Go/build paths empty) |
| Opening tracked status | Clean; `feat/issues-and-pr` synced with `origin/feat/issues-and-pr` |
| Disposable tree | Docker volume `cp2-work:/work/src`; 453 files; no `.git`; symlink fixture preserved as a real Linux symlink |
| Export provenance | Container `git -c safe.directory=/repo archive 5f6ec56` piped to `tar -x`; extracted `go.mod` blob `3a16143e…`, `go.sum` blob `30e18339…` match the pinned Git blobs exactly; tree manifest `3ff1f7da…` |
| Host-archive caveat | A first export made by Windows Git (`git archive` on the host) produced CRLF module files; it was discarded. All results use the clean LF export. This confirms CP0’s CRLF caveat. |
| Toolchain | `go1.27.3 linux/amd64`, `GOOS=linux`, `GOARCH=amd64`, `GOTOOLCHAIN=local`, `CGO_ENABLED=1`, `GOMODCACHE=/gomod`, `GOCACHE=/gocache` |
| Module requirement | `go 1.24.11` ([`go.mod:3`](../../../go.mod#L3)); CI 1.25 ([`test.yml:23`](../../../.github/workflows/test.yml#L23)); Dockerfile 1.24 ([`Dockerfile:1`](../../../Dockerfile#L1)) — **not directly comparable** |
| Dependency acquisition | `go mod download` with network enabled, exit 0, 24.1 s; external cache `/gomod` ≈517 MB; `go mod verify` → “all modules verified”; `go.mod`/`go.sum` blob IDs unchanged after all work |
| Build provenance | `go version -m /out/gitleaks`: Go 1.27.3, `CGO_ENABLED=1`, `GOOS=linux`, `GOARCH=amd64`, module `(devel)` because the export has no `.git`; dependencies carry `h1:` hashes |
| Release variant | `go build -tags gore2regex .` succeeded; build info records `-tags=gore2regex`, `CGO_ENABLED=1` |
| Test user | Root for the main suite; `uid 1000` for the unreadable-file test |
| Constraints | No Windows runtime/tests; native Go 1.27.3 Windows is available at CP2A review but was not used for CP2; no external scans; no #2019/#1572; no unbounded inputs; test artifacts outside the repository |

**Module provenance note:** `go list -m all` with `GOPROXY=off` failed because the offline cache lacked some module-graph metadata. It is not a test failure; `go mod verify`, the successful build/test runs, and embedded build info establish the dependency set.

---

## 2. Capability and CLI contract matrix

Status vocabulary: **PASS** = observed behavior matches expected; **FAIL** = observed defect/contract mismatch; **CONTRACT** = observed behavior requires a documented contract decision (not automatically a defect); **PARTIAL** = some behavior verified, some unverified/ambiguous; **NOT RUN**; **EXCLUDED**.

| Capability / flag | Code anchor | Expected | Observed | Status | Evidence |
| --- | --- | --- | --- | --- | --- |
| `version` | [`cmd/version.go:20-22`](../../../cmd/version.go#L20-L22) | Print version | `version is set by build process` (no ldflags in test build) | PASS | `version.out`, exit 0 |
| `--help` / visible command list | [`cmd/root.go:44-59`](../../../cmd/root.go#L44-L59) | Help, hidden modes absent | Shows `completion`, `dir`, `git`, `stdin`, `version`; hidden `detect`/`protect` omitted | PASS | `help.out`, exit 0 |
| Unknown flag | [`cmd/root.go:224-232`](../../../cmd/root.go#L224-L232) | Exit 126 | `Error: unknown flag…`; exit 126 | PASS | `badflag`, exit 126 |
| `dir` directory/file scan | [`cmd/directory.go:24-73`](../../../cmd/directory.go#L24-L73) | Finding, exit override | 1 finding; exit 0 with `--exit-code 0`; exit 1 by default | PASS | `dir-token`, `dir-token-exit1`, `dir-file` |
| `git` history | [`cmd/git.go:29-95`](../../../cmd/git.go#L29-L95) | Commit finding | 1 finding at commit/line; history works | PASS | `git-normal`, `detect-history` |
| `git --staged` / `--pre-commit` | [`cmd/git.go:64-77`](../../../cmd/git.go#L64-L77) | Diff modes | `--staged` finds staged token; `--pre-commit` scans unstaged diff and found 0 in this fixture | PASS (contract clarification) | `git-staged` 1; `git-precommit` 0 |
| `protect --staged` | [`cmd/protect.go:29-65`](../../../cmd/protect.go#L29-L65) | Staged finding | 1 finding | PASS | `protect-staged` |
| `stdin` / `detect --pipe` | [`cmd/stdin.go:23-53`](../../../cmd/stdin.go#L23-L53) | Finding from stdin | 1 finding each | PASS | `stdin-case`, `detect-pipe` |
| `detect --no-git` | [`cmd/detect.go:75-90`](../../../cmd/detect.go#L75-L90) | Directory finding | 1 finding | PASS | `detect-nogit` |
| `--config` / env / source-local / default | [`cmd/root.go:132-195`](../../../cmd/root.go#L132-L195) | Precedence | `source-local`, `external-rule`, `env-rule` selected correctly | PASS | `precedence-*`, `env-config-path` |
| Config `[extend]` | [`config/config.go:375-492`](../../../config/config.go#L375-L492) | Extend rules | Covered by existing config tests; not separately reproduced | PARTIAL | CP1 source; package tests pass |
| `--enable-rule` | [`cmd/root.go:332-345`](../../../cmd/root.go#L332-L345) | Restrict rules | Selected rule only; unknown rule fatal | PASS | `enable-specific`, `enable-generic`, `enable-unknown` |
| `--max-decode-depth` (default 5) | [`detect/detect.go:349-363`](../../../detect/detect.go#L349-L363) | N decodes, N+1 passes | depth 0→0 findings; 1→1; 2→2; default→2 on two encoded files | PASS | `decode-d0/d1/d2/default` |
| `--max-archive-depth` (default 0) | [`sources/file.go:52-84`](../../../sources/file.go#L52-L84) | 0 disables; 1 scans archive | 0→0; 1→1 finding (`archive.tar.gz!inner.txt`) | PASS | `archive-d0/d1` |
| `--max-target-megabytes` (default 0) | [`sources/files.go:92-99`](../../../sources/files.go#L92-L99), [`detect/detect.go:430-440`](../../../detect/detect.go#L430-L440) | 0 disables; 1MB skips larger file | 1.1MB file: flag 1→0 findings; unset→1 finding | PASS | `maxsize-1/0` |
| `--follow-symlinks` | [`sources/files.go:103-121`](../../../sources/files.go#L103-L121) | Off skips link; on scans target | off→1; on→2 (direct + via symlink) | PASS | `symlink-off/on` |
| `--redact` (100) | [`detect/utils.go:158-160`](../../../detect/utils.go#L158-L160), [`report/finding.go:78-86`](../../../report/finding.go#L78-L86) | Secret removed | Normal finding redacted; **path-only finding corrupted** (repeated `REDACTED` between every character) | **FAIL** | `pathonly-redact` |
| `--redact=50` | same | Partial mask | `Secret`/`Match` truncated to `CP2TOK...` | PASS | `redact50` |
| `gitleaks:allow` / `--ignore-gitleaks-allow` | [`detect/detect.go:510-515`](../../../detect/detect.go#L510-L515) | Inline suppression / override | allow→0; ignore-allow→1 | PASS | `inline-allow`, `inline-ignore` |
| Path allowlist | [`config/allowlist.go:136-146`](../../../config/allowlist.go#L136-L146) | Path filtered | ignored file omitted; kept file found | PASS | `allow-path` |
| `.gitleaksignore` finite file | [`detect/detect.go:153-193`](../../../detect/detect.go#L153-L193) | Fingerprint suppression | 0 findings; control 1 | PASS | `ignore-file`, `ignore-control` |
| Baseline | [`detect/baseline.go:12-35`](../../../detect/baseline.go#L12-L35) | Suppress matching findings | Original suppressed; baseline JSON itself re-scanned (2 findings) | **PARTIAL/FAIL** | `baseline-second` |
| Generic-rule dedupe | [`detect/utils.go:137-166`](../../../detect/utils.go#L137-L166) | Generic suppressed when specific matches | Only `specific-cp2` retained | PASS | `generic` |
| Entropy gate | [`detect/detect.go:542-554`](../../../detect/detect.go#L542-L554) | Low threshold keeps; high suppresses | entropy 2.0→finding; 4.0→0 | PASS | `entropy-low/high` |
| JSON report | [`report/json.go:13-16`](../../../report/json.go#L13-L16) | Valid JSON array | 1 finding parsed | PASS | `fmt-json` |
| CSV report | [`report/csv.go:15-77`](../../../report/csv.go#L15-L77) | Header+row | 1 row; empty findings → 0-byte file | PASS (normal) / CONTRACT (empty) | `fmt-csv`, `empty-csv` |
| JUnit report | [`report/junit.go:16-107`](../../../report/junit.go#L16-L107) | Valid XML | 1 testcase/failure | PASS | `fmt-junit` |
| SARIF report | [`report/sarif.go:17-217`](../../../report/sarif.go#L17-L217) | Valid SARIF 2.1.0 | Parsed; 1 result | PASS | `fmt-sarif` |
| Template report | [`report/template.go:19-53`](../../../report/template.go#L19-L53) | Render template | `cp2-token|path` | PASS | `fmt-template` |
| Report path `-` without format | [`cmd/root.go:366-379`](../../../cmd/root.go#L366-L379) | No implicit format documented; stdout likely expected | Fatal `Unknown report format:`; exit 1 | **CONTRACT** | `stdout-no-format` |
| Report path `.json` inference | same | JSON written | 448-byte JSON | PASS | `infer-json` |
| `--report-template` alone | [`cmd/root.go:399-402`](../../../cmd/root.go#L399-L402) | Help says it implies template | Fatal unless `--report-format=template` | **CONTRACT** | `tmpl-mismatch` |
| Existing report target | [`cmd/root.go:347-358`](../../../cmd/root.go#L347-L358) | Write/overwrite at end | Target deleted during validation; sentinel gone | **FAIL** | `report-delete` |
| `--exit-code` | [`cmd/root.go:493-499`](../../../cmd/root.go#L493-L499) | Findings→flag value; else 0/1 | `--exit-code 0`→0; default→1; clean→0 | PASS | `dir-token`, `dir-token-exit1`, `dir-clean` |
| `--diagnostics mem` | [`cmd/diagnostics.go:178-210`](../../../cmd/diagnostics.go#L178-L210) | Write profile | `mem.pprof` created | PASS | `diagnostics-mem` |
| `--diagnostics cpu/trace/http` | [`cmd/diagnostics.go:150-239`](../../../cmd/diagnostics.go#L150-L239) | Profiles/server | Not run | NOT RUN | — |
| `--verbose`, `--no-color`, explicit `--gitleaks-ignore-path` | [`cmd/root.go:81-90`](../../../cmd/root.go#L81-L90) | Output/log variants | Not exercised; default ignore path was tested | NOT RUN | — |
| `--platform github` links | [`detect/utils.go:23-110`](../../../detect/utils.go#L23-L110) | Commit link | Link generated under `--network none` | PASS | `git-link` |
| Invalid platform | [`cmd/git.go:74-76`](../../../cmd/git.go#L74-L76) | Fatal | `invalid scm platform value: bogus`; exit 1 | PASS | `bad-platform` |
| Git stderr/subprocess error | [`sources/git.go:414-420`](../../../sources/git.go#L414-L420), [`detect/detect.go:228-232`](../../../detect/detect.go#L228-L232) | Nonzero/failed scan | `[git] fatal…`, `stderr is not empty`, 0 commits, `[]`, exit 0, “no leaks found” | **FAIL** | `git-badflag` |
| Unreadable file (non-root) | [`sources/files.go:155-162`](../../../sources/files.go#L155-L162) | Skip/error | Warning, remaining finding, exit 0 | **CONTRACT** | `perm-nonroot` |
| Timeout on blocked stdin | [`cmd/root.go:48-58`](../../../cmd/root.go#L48-L58), [`sources/file.go:174-190`](../../../sources/file.go#L174-L190) | Abort near deadline | Waited 2,009 ms for a 1 s timeout; then fatal `context deadline exceeded` | **CONTRACT** | `timeout-stdin` |
| Invalid regex | [`config/config.go:123-128`](../../../config/config.go#L123-L128) | Controlled config error | Panic, exit 2 | **FAIL** | `bad-regex` |
| Finding order | [`detect/detect.go:327`](../../../detect/detect.go#L327) | Deterministic output | 8 runs alternated `alpha/beta` order | **CONTRACT** | `order-1..8` |
| `gore2regex` application binary | [`regexp/wasilibs_regex.go:1-14`](../../../regexp/wasilibs_regex.go#L1-L14) | Build and run | Binary built; 222-rule default config scanned successfully | PASS | `re2-default` |
| Tagged build/test graph | same | Build all packages | `go build -tags gore2regex ./...` and `go test -tags gore2regex ./...` fail in `cmd/generate/config/rules` | **FAIL** (C3) | tagged build/test logs |

---

## 3. Package-test and static-check evidence

All commands ran in `/work/src` (disposable volume) with `GOFLAGS=-mod=readonly`, `GOPROXY=off`, `GOMODCACHE=/gomod`, `GOCACHE=/gocache`, `GOTOOLCHAIN=local`, `CGO_ENABLED=1`, `--network none`.

| Command | Exit | Duration | Result |
| --- | --- | --- | --- |
| `go build ./...` | 0 | 20.08 s (with binary) | Success |
| `go build -o /out/gitleaks .` | 0 | (same run) | 24,570,239-byte Linux binary |
| `go test -mod=readonly -count=1 -timeout 600s -json ./...` | 0 | 6.79 s | **275 passed, 0 failed, 0 skipped** |
| `go test -race -mod=readonly -count=1 -timeout 600s ./...` | 0 | 34.39 s | Success on exercised paths only (not a proof for unexercised code) |
| `go test -shuffle=on -count=1 ./...` | 0 | 7.25 s (combined) | Success |
| `go vet ./...` | **1** | 6.11 s (with gofmt) | 2 `lostcancel` diagnostics, test file only |
| `gofmt -l .` | 0 | (same run) | 0 files listed (LF export) |
| `go test -coverprofile=/out/cover.out ./...` | 0 | 4.46 s | Total **30.9%** statements |
| `go build -tags gore2regex -o /out/gitleaks-re2 .` | 0 | 6.08 s | Success |
| `go build -tags gore2regex ./...` | **1** | (same run) | Two rule-generator compile errors |
| `go test -tags gore2regex -count=1 ./...` | **1** | 7.95 s | Rule-generator compile failure; other packages passed |

**Per-package test results (default tag):**

| Package | Status | Tests | Coverage |
| --- | --- | ---: | ---: |
| `cmd/generate/config/base` | pass | 14 | 100.0% |
| `cmd/generate/config/utils` | pass | 13 | 28.6% |
| `config` | pass | 45 | 79.1% |
| `detect` | pass | 113 | 73.9% |
| `detect/codec` | pass | 55 | 60.3% |
| `report` | pass | 28 | 76.8% |
| `sources` | pass | 7 | 7.4% |
| `main`, `cmd`, `cmd/generate/config`, `cmd/generate/config/rules`, `cmd/generate/secrets`, `cmd/scm`, `logging`, `regexp`, `version` | no test files | 0 | 0% / not instrumented |

**`go vet` failure (exact):**

```
detect/detect_test.go:2083:4: the cancel function is not used on all paths (possible context leak)
detect/detect_test.go:2122:3: this return statement may be reached without using the cancel var defined on line 2083
```

Source inspection shows `cancel()` is called in both the `expireContext` and normal branches ([`detect/detect_test.go:2083-2105`](../../../detect/detect_test.go#L2083-L2105)); the warning is likely conservative, but it is an existing static-check failure. `go vet -tests=false` still reported it; CI does not run vet ([`test.yml`](../../../.github/workflows/test.yml#L1)).

**Coverage interpretation:** 30.9% is the total over instrumented (tested) packages only; `cmd` and `main` are not instrumented by `-coverprofile` without `-coverpkg`, so effective whole-program coverage is lower. `sources` (7.4%) is the largest uncovered area relevant to the fail-open findings.

---

## 4. Correctness findings register

CP2A classification ([`review-CP2A.md`](review-CP2A.md)): **C1 correctness defects** — CP2-004, CP2-008, CP2-014 (logging/work; no demonstrated false negative); **C2 security/data-integrity risks** — CP2-003, CP2-005, CP2-007, CP2-009; **C3 build/static-analysis failures** — CP2-010, CP2-011 (tagged graph); **C4 contract decisions** — CP2-001, CP2-002, CP2-006, CP2-012, CP2-013, CP2-015; **C5 expected/intentional** — CP2-016; **C6 unverified** — CP2-017; **C7 positive verification** — CP2-018 and the `gore2regex` application-binary result. The table’s original measured statuses are preserved.

| ID | Finding | Status | Controlled reproduction (expected → observed) | Impact | Owner |
| --- | --- | --- | --- | --- | --- |
| CP2-001 | `--report-path -` without an explicit format is fatal | Measured | `dir --report-path -` (no `--report-format`) → expected stdout JSON; observed `FTL Unknown report format:` and exit 1 | CLI usability/automation break; `-` is documented but not self-describing | CP2 |
| CP2-002 | `--report-template` does not imply template format | Measured | `--report-format json --report-template t` → expected template or clear implication; observed fatal “Report format must be 'template'…” | Help/code contract mismatch | CP2 |
| CP2-003 | Existing report target is deleted during validation | Measured | Sentinel file + `--report-format template --report-template missing` → expected failure without touching sentinel; observed sentinel absent | Data-loss hazard before scan; symlink-target risk | CP2/CP4 |
| CP2-004 | Path-only finding corrupted by 100% redaction | Measured | Path-only rule + file + `--redact` → expected `file detected: …`; observed `Match` with `REDACTED` between every character, `Secret: REDACTED` | Corrupted output for `--redact` with path rules; no test coverage | CP2 |
| CP2-005 | Git subprocess/stderr failure reported as success | Measured | `git --log-opts --not-a-real-option` → expected failed/partial status; observed Git fatal + `stderr is not empty`, 0 commits, `[]`, exit 0, “no leaks found” | Operator cannot distinguish partial/failed scan from clean scan; CI false confidence | CP2 |
| CP2-006 | Unreadable file is a warning and scan succeeds | Measured (non-root) | Directory with unreadable + readable token files → expected explicit partial status; observed warning, 1 finding, exit 0 | Partial coverage not reflected in status | CP2 |
| CP2-007 | Baseline file inside source is scanned | Measured | Generate baseline in source, rescan with `--baseline-path` → expected original suppressed and baseline excluded; observed original suppressed but baseline JSON produced 2 findings (lines 9/10) | Baseline secrets can be re-reported; path-form issue | CP2 |
| CP2-008 | Invalid regex panics | Measured | Config with `regex = '['` → expected config error; observed panic at `regexp.MustCompile`, exit 2 | Crash on malformed config; poor operator experience | CP2 |
| CP2-009 | `--redact` leaves secrets in debug/trace logs | Measured | Encoded token + `--redact --log-level debug`; inline token + `--log-level trace` → synthetic token present in stderr | Secret exposure in logs despite redaction | CP4 |
| CP2-010 | `go vet` fails | Measured | `go vet ./...` → exit 1 on test lostcancel | Static-check gate failure; test-only | CP2 |
| CP2-011 | `gore2regex` full graph fails | Measured | `go test -tags gore2regex ./...` and `go build -tags gore2regex ./...` fail on `airtable.go:30`, `perplexity.go:16`; release binary `.` builds and default config loads | Tagged build/test parity defect; release path itself works | CP2/CP4 |
| CP2-012 | Finding order nondeterministic | Measured | 8 identical scans of two matching rules → sequences alternated `alpha,beta` / `beta,alpha` | Report/diff instability; consumers must sort | CP2 |
| CP2-013 | Blocking stdin delays timeout enforcement | Measured | 2 s delayed stdin, `--timeout 1` → observed 2,009 ms then fatal `context deadline exceeded` | Timeout is not a hard wall for blocked reads; run may exceed deadline | CP3 |
| CP2-014 | Git commit allowlist logs a skip but scans the diff | Measured | Trace run with commit allowlist → “skipping commit: global allowlist”, then detector “skipping file … allowed-commit=true”; 0 findings | Source-level log is misleading and work is redundant; detector guard prevents false negatives | CP2 |
| CP2-015 | Empty CSV report has no header | Measured | Clean scan + `--report-format csv` → 0-byte file | Downstream CSV parsers may fail | CP2 |
| CP2-016 | `git --pre-commit` vs `--staged` differ | Measured | Staged-only change: `--pre-commit`→0 findings; `--staged`→1 | Mode contract is easy to misuse; document explicitly | CP2 |
| CP2-017 | Truncated archive produced no diagnostic | Unverified | Truncated `.tar.gz` + depth 1 → no error, no findings | Whether skipped as binary or silently failed is unresolved | CP2/CP4 |
| CP2-018 | Positive capabilities verified | Measured | See §2: decode/archive/size/symlink/entropy/allowlists/reports/exit/preference | Core behavior is largely sound | — |

**Negative findings preserved:** CP1’s commit-allowlist bug did **not** produce a user-visible false negative in this test because the detector-level allowlist still suppressed the finding; it remains a source-level inefficiency/logging defect. The truncated-archive test did not reproduce an error path, so no archive-failure claim is made. Detailed CP2A reproduction evidence for CP2-003/004/005/006/007/008/009/013 is in [`review-CP2A.md` §3](review-CP2A.md); raw outputs remain outside the repository.

---

## 5. CP1 AR crosswalk

| AR | Disposition | Evidence |
| --- | --- | --- |
| AR-01 Error absorption / partial success | **CONFIRMED** | Git bad-flag exit 0 (`git-badflag`); non-root unreadable file warning + exit 0 (`perm-nonroot`); CP1 C-01/C-02 |
| AR-02 `protect` background context / timeout | **PARTIAL** | Source confirms `DetectGit` Background; delayed stdin timeout measured late (`timeout-stdin`); `protect` end-to-end timeout not run |
| AR-03 Redaction after detection logs | **CONFIRMED** | Synthetic token in debug/trace stderr (`redact-debug-leak`, `redact-trace-inline`) → CP4 |
| AR-04 Path-only redaction | **CONFIRMED** | `pathonly-redact` corruption |
| AR-05 Config/baseline path forms | **CONFIRMED** | Baseline JSON self-scan (`baseline-second`); config self-skip not separately reproduced |
| AR-06 Global `extendDepth` | **BLOCKED** | No multi-translation harness; source evidence only; existing tests reset `extendDepth` manually |
| AR-07 Invalid regex panic | **CONFIRMED** | `bad-regex`, exit 2 |
| AR-08 Git commit allowlist `continue` | **CONFIRMED (source log)** | `git-allow` trace; no user-visible false negative |
| AR-09 `max-target-megabytes` post-decode/opt-in | **DEFERRED to CP3** | Semantics verified (`maxsize-1/0`); no resource scaling measured |
| AR-10 Decoder allocation/per-fragment cache | **DEFERRED to CP3** | Depth semantics verified (`decode-*`); no memory measurement |
| AR-11 `.gitleaksignore` untrusted input | **DEFERRED to CP4** | Finite ignore file works (`ignore-file`); no device/unbounded tests |
| AR-12 `gore2regex` release parity | **PARTIAL/CONFIRMED defect** | Release binary builds and default config loads; full tagged build/test fails in rule-generator package |
| AR-13 Report target deletion | **CONFIRMED** | `report-delete` sentinel absent |
| AR-14 `--report-template` mismatch | **CONFIRMED** | `tmpl-mismatch`, exit 1 |
| AR-15 Generic dedupe semantics/perf | **PARTIAL** | `generic` selects specific rule; O(n²) performance deferred to CP3 |

---

## 6. Reporting, build, and platform observations

- **Report writers:** JSON, CSV, JUnit, SARIF, and template outputs were structurally validated (JSON parse, CSV row, XML parse, SARIF JSON, template text). CSV with no findings is a zero-byte file. SARIF embeds `Secret` in the snippet; redaction is inherited from findings.
- **Report path:** `-` writes stdout but requires an explicit format; extension inference works for `.json`. Validation deletes an existing target before scanning.
- **Exit behavior:** findings → `--exit-code` (default 1); clean/`--exit-code 0` → 0; invalid flag → 126; config/report fatal → 1; invalid regex panic → 2.
- **Redaction:** normal findings and partial redaction work; path-only findings are corrupted; debug/trace logs leak the raw synthetic token. Report-level redaction cannot compensate for log leakage.
- **Build tags:** default stdlib build passes; `gore2regex` release binary builds and loads all 222 default rules; the tagged all-package build/test graph fails because two rule-generator files import stdlib `regexp` ([`airtable.go:30`](../../../cmd/generate/config/rules/airtable.go#L30), [`perplexity.go:16`](../../../cmd/generate/config/rules/perplexity.go#L16)).
- **Platform:** all runtime results are Linux/amd64 under Go 1.27.3. CP0’s Windows caveats (CRLF, symlink materialization, two skipped tests) remain; no Windows execution occurred. The host-archive CRLF issue was encountered and avoided.
- **Toolchain comparability:** module 1.24.11, CI 1.25, Dockerfile 1.24, test toolchain 1.27.3. These results characterize 1.27.3 only and must not be presented as CI-equivalent.

---

## 7. Regression gap register

| Gap | Existing coverage | Proposed smallest future test | Priority |
| --- | --- | --- | --- |
| `--report-path -` format requirement | None in `cmd` | Table test: `-` without/with `--report-format`; assert error/output | High |
| `--report-template` implication | None | Run with template only; assert intended contract | Medium |
| Report target deletion | None | Pre-create sentinel + failing reporter; assert file survives/does not | High |
| Path-only + `--redact` | `report/finding_test.go` covers non-empty secret only | Path-only finding + `Redact(100)`; assert Match/Line | High |
| Git stderr/subprocess failure | `sources` 7.4% coverage; no fail-open test | Fake/real `git log` error; assert nonzero exit or explicit partial status | High |
| Unreadable file | No non-root test | Non-root container with `chmod 000`; assert status semantics | Medium |
| Baseline self-scan | No test | Baseline inside source vs outside; assert exclusion | High |
| Invalid regex | No panic test | Config with `[`; assert controlled error | Medium |
| Debug/trace redaction | No test | Synthetic secret + debug/trace + `--redact`; assert absent from stderr | High (CP4) |
| Finding order | Tests use `ElementsMatch` | Stable-sort contract test or documented nondeterminism | Medium |
| `gore2regex` all-package graph | CI runs default only | Tagged `go build ./...`/`go test ./...` in CI | High (CP4) |
| `go vet` lostcancel | No vet in CI | Fix test cancel handling or document false positive; run vet | Medium |
| Timeout/blocked I/O | No test | Bounded delayed stdin; assert timeout behavior | Medium (CP3) |
| Windows skips | Two known skips | Windows CI coverage or documented unsupported cases | Medium |
| CSV empty output | CSV test likely non-empty | Empty findings; assert header/contract | Low |
| Truncated archive | No test | Failing archive fixture; assert diagnostic/exit | Medium (CP4) |

---

## 8. Readiness and handoff

**CP2 status:** complete for Linux/amd64 under Go 1.27.3, with the documented exceptions above. **CP3 is READY to begin** subject to authorization to reuse the existing container image, module cache, and disposable-tree method; CP3 must not treat these results as performance measurements.

**CP3 handoff:** decode-depth memory and allocation growth (#2019 pathway), `max-target-megabytes` post-decode behavior, archive extraction, 40-worker concurrency/backpressure, finding accumulation, and generic dedupe O(n²). Keep #2019 unbounded workloads excluded.

**CP4 handoff:** debug/trace redaction leakage, `.gitleaksignore` symlink/size and #1572, SARIF secret snippets, report-path deletion/symlink targets, `gore2regex` release/test parity, and fail-open error status. The finite ignore test passed; no malicious/device fixtures were used.

**Unresolved prerequisites:** no Windows runtime; no native Go; no CI-Go 1.25 test run; `go vet` failure unresolved; `gore2regex` all-package graph unresolved; performance/security consequences deferred.

---

## 9. Integrity and definition-of-done

- **Tracked checkout:** clean before and after; only `_assessment/current-state/CP2-capabilities-correctness/report-CP2.md` plus authorized `_assessment/README.md`/roadmap status updates were added; HEAD `c39f9a8…` unchanged; no source/test/fixture/config/CI/backlog modification.
- **Disposable tree:** `go.mod`/`go.sum` blob IDs match the pinned blobs after all runs (`3a16143e…`, `30e18339…`); tests may have renamed fixtures inside the disposable volume only. Dependency cache is external (`/gomod`, ≈517 MB), verified by `go mod verify`.
- **Secrets:** only synthetic `CP2TOKEN-NNNN` values and artificial paths were used; no real credentials, upstream repository scans, or secret-like fixture values are recorded.
- **Commands:** module acquisition was the only network-enabled step; all build/test/vet/CLI runs used `--network none`, `GOPROXY=off`, `-mod=readonly`, and external caches.

**Definition-of-done review**

- [x] Every capability exercised in this report has a source anchor and an observed/unverified status; untested variants are marked NOT RUN.
- [x] Package tests, race, shuffle, coverage, vet, and build results include commands, exit codes, timing, and conditions.
- [x] All 15 CP1 AR items have a disposition; no unrun test is described as passing.
- [x] Correctness findings include controlled reproductions and expected-versus-actual behavior.
- [x] Protected files unchanged; before/after status recorded.
- [x] CP3/CP4 handoffs and unresolved prerequisites documented.
- [x] Stopped pending owner review; CP3 not started.

**Recommendation:** Ready for owner review.

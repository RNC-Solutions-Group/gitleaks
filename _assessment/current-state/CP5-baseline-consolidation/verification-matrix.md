# CP5 — Verification and regression matrix

**Status:** Complete — consolidated from CP0–CP4/CP4A  
**Canonical code SHA:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`  
**Tested runtime:** Linux/amd64, Go 1.27.3, `CGO_ENABLED=0` for performance and CP4 security cases; CP2 tests used `CGO_ENABLED=1`. No Windows, CI Go 1.25, or release-packaging execution.

**Evidence types:** Automated Go test; targeted CLI experiment; source verification; reused measured evidence; not run. “Verified” here means the stated capability was exercised or source-confirmed; it does not imply cross-platform parity.

---

## Capability, evidence, and regression gaps

| Capability | Expected / documented | Observed | Evidence (CP) | Environment | Evidence type | Regression gap | Future PR must demonstrate |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `git` history scan | Find added secrets; exit per findings | 1 finding, exit 1; clean exit 0 | CP2, CP4 `git-finding`/`dir-clean` | Linux Go 1.27.3 | CLI + source | CLI package has no direct tests | Git history/staged/unstaged parity; exit |
| `git` Git failure | Failed/partial scan signaled | Git fatal logged; 0 commits; `[]`; **exit 0** | CP4 `git-fail-open` (no `--exit-code 0`) | Linux Go 1.27.3 | CLI | No fail-open test | Nonzero/partial status; malformed stderr/diff |
| `dir`/file scan | Findings; clean success | 1 finding; clean 0; file path works | CP2/CP4 | Linux | CLI | CLI mode tests absent | Dir/file paths, permission failures |
| `stdin` / `detect --pipe` | Findings from stream | 1 finding each | CP2 | Linux | CLI | Blocking-stream test absent | Stream EOF/cancel/timeout |
| `detect --no-git` | Directory finding | 1 finding | CP2 | Linux | CLI | Legacy mode little coverage | Parity with `dir` |
| `protect --staged` | Staged finding | 1 finding; background detector context | CP2/CP4 source | Linux | CLI + source | Timeout/context test absent | Staged-only semantics; timeout |
| `git --pre-commit` / `--staged` | Documented diff mode | Unstaged scan 0 vs staged 1 in fixture | CP2 | Linux | CLI | Contract ambiguity | Explicit mode contract + tests |
| Exit contracts | Findings → flag; clean → 0; fatal → 1; unknown flag 126 | Confirmed; invalid regex panic 2 | CP2/CP4 | Linux | CLI | No CLI test suite | Full exit matrix incl. partial |
| Config precedence | `--config` > env > `GITLEAKS_CONFIG_TOML` > source-local > default | All four explicit/env/source paths selected expected rule | CP2 | Linux | CLI | No cmd tests; env matrix thin | Precedence + malformed config tests |
| Rules/allowlists | Rule matching, path/commit/regex/stopword allowlists | Path allowlist, inline allow, commit allowlist verified | CP2 | Linux | CLI + source | Commit-allowlist source log bug | Allowlist parity; no suppression bypass |
| `.gitleaksignore` | Fingerprint suppression | Finite file and finite symlink suppression work | CP2/CP4 | Linux | CLI + source | No device/size/cap test; #1572 unverified | Regular-file/size/blocking safeguards |
| Baseline | Suppress known findings | Original suppressed; baseline JSON inside source re-scanned | CP2/CP4 | Linux | CLI | No path-normalization test | Baseline exclusion/self-scan tests |
| Decode depth | Default 5; 0 disables; N+1 passes | 0→0, 1→100, 2→100 on encoded fixtures; 5 MiB depth 0/1/2 94/950/968 ms | CP2/CP3 | Linux | CLI + measurement | No full decoder matrix at scale | Coverage parity across all encodings |
| Archive depth | Default 0; depth 1 scans | 0→0; 1→100 findings; depth1 ~1.7× ordinary | CP2/CP3 | Linux | CLI + measurement | No uncompressed budget/bomb test | Depth/entry/size limits |
| Symlinks | Default skip; opt-in follow | Off→1; on→2 findings; finite ignore symlink followed | CP2/CP4 | Linux | CLI + source | No containment test | Contained follow policy |
| Size limits | `--max-target-megabytes` skips oversized | 1.1 MiB skipped at flag 1; 0.9 MiB unaffected | CP3 | Linux | Measurement | No Git-source pre-cap test | Threshold behavior both sources |
| Timeout/cancel | Abort near deadline | Delayed stdin: 2,009 ms for 1 s, fatal after read | CP2/CP4 | Linux | CLI | `protect` timeout not run | Deadline enforcement without blocking |
| Redaction | `--redact` protects findings/reports/logs | Reports/JSON/CSV/JUnit/SARIF/template/verbose clean; debug stderr leaks; path-only corrupt | CP2/CP4 | Linux | CLI | No regression test for debug/path-only | Output/log matrix with synthetic secret |
| Reports | JSON/CSV/JUnit/SARIF/template valid | All valid; empty CSV 0 bytes; `-` needs explicit format; SARIF includes secret unless redacted | CP2/CP4 | Linux | CLI | No cmd-level tests | Format/empty/redaction contracts |
| Finding order | Deterministic? | Nondeterministic across 8 runs | CP2 | Linux | CLI | Tests use order-insensitive match | Decided ordering contract |
| Concurrency | 40 workers; stable | 12–17 OS threads; ~1.9/2 cores on 32 MiB; no timeouts | CP3 | Linux 2 vCPU | Measurement | No goroutine/race-stress beyond CP2 | Race/soak at varied sizes |
| Default build | Builds | `go build ./...` pass; app binary pass | CP2/CP4 | Linux | Build | No native Windows build | Build across CI versions |
| `gore2regex` app | Tagged app builds/runs | App binary builds; default 222-rule scan pass | CP2/CP4 | Linux CGO0 | Build + CLI | Tagged full graph fails | App + full graph + tests with tag |
| `go vet` | Clean static analysis | Fails 2 test `lostcancel` diagnostics | CP2 | Linux | Static | CI doesn't run vet | Clean vet or documented suppression |
| Go tests | All pass | **275/275 pass**, 0 skip on Linux | CP2 | Linux | Automated | cmd/main untested; coverage gaps | Same/expanded suite per PR |
| Race tests | Pass in exercised scope | `-race` pass; no proof for unexercised paths | CP2 | Linux | Automated | Limited path coverage | Race tests for changed paths |
| Coverage | Measure | 30.9% instrumented; cmd 0%, sources 7.4% | CP2 | Linux | Automated | Whole-program lower | Coverage delta/justification |
| Windows | Supported CI matrix | **NOT RUN**; known skips/symlink/CRLF caveats | CP0/CP1 | — | Not run | Major platform gap | Windows CI or explicit unsupported |
| CI Go 1.25 | CI build/test | **NOT RUN**; local 1.27.3 only | CP2/CP4 | — | Not run | Toolchain parity gap | CI-equivalent run |
| Release packaging | Docker/goreleaser | **NOT RUN**; static inspect only | CP4 | — | Not run | Provenance/package gap | Release build/signing/SBOM review |
| Dependency advisories | Known-vulnerability status | **NOT RUN** | CP4/CP5 | — | Not run | Supply-chain gap | Advisory scan at pinned versions |

---

## Test-suite baseline (carried forward, unchanged)

| Metric | Result |
| --- | --- |
| `go test ./...` | 275/275 pass, 0 fail, 0 skip (Linux) |
| `go test -race ./...` | Pass on exercised paths |
| `go test -shuffle=on` | Pass |
| `gofmt -l` (LF export) | 0 files |
| Statement coverage | 30.9% instrumented; `cmd` 0%, `sources` 7.4%, `detect` 73.9%, `config` 79.1%, `codec` 60.3%, `report` 76.8% |
| `go vet` | Fail: `detect_test.go:2083`, `:2122` lostcancel; test-only |
| `gore2regex` app | Builds; default 222-rule config loads |
| `gore2regex` full graph | Build/test fail in `airtable.go:30`, `perplexity.go:16` |
| Platform | Linux/amd64 only; no Windows/CI Go 1.25/release execution |

---

## Regression-gap priorities

| Priority | Gap | Related finding | Required future test |
| ---: | --- | --- | --- |
| P0 | No fail-open/partial test | CP5-F01 | Induce Git stderr/exit error, unreadable file, worker error; assert partial/nonzero |
| P0 | No destructive report-path test | CP5-F02 | Sentinel + symlink target; failure must preserve target |
| P1 | No debug/trace/path-only redaction regression | CP5-F03 | Synthetic secret across all formats/log levels |
| P1 | No untrusted config/ignore safeguards test | CP5-F04/F05 | Allowlist evasion, panic, regular-file/size controls |
| P2 | Baseline self-scan untested | CP5-F06 | Baseline inside/outside source, path forms |
| P2 | Archive/decode resource limits untested at scale | CP5-F08/F09 | Bounded expansion/cancellation, coverage fixtures |
| P2 | Tagged full graph/vet not gated | CP5-F14 | Tagged `go build ./...`/`go test ./...`; vet |
| P2 | CI/release gaps | CP5-F15 | CI-equivalent run, pinned tools, SBOM/signing |
| P3 | Finding order/contracts | CP5-F12/F13 | Deterministic contract or sort; exit/report matrix |
| P3 | Windows/CI/release/advisories | CP5-F16/F18 | Platform matrix, advisory process |

This matrix is a **verification inventory**, not a claim of full coverage. Every “NOT RUN” row remains a gap.

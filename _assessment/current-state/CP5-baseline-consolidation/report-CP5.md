# CP5 — Final technical baseline consolidation

**Status:** Complete — awaiting owner review and baseline acceptance  
**Assessment date:** 2026-11-09 (America/New_York)  
**Repository / branch:** `RNC-Solutions-Group/gitleaks` / `feat/issues-and-pr`  
**Canonical application SHA:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`  
**Documentation HEAD at consolidation:** `f437f40059f541f7a6198a025afec08da35426e2` (CP4A/CP5 documentation commit follows)  
**Assignment:** [`assignment-CP5.md`](assignment-CP5.md)  
**Detailed evidence:** [`findings-register.md`](findings-register.md), [`verification-matrix.md`](verification-matrix.md), [`technical-roadmap.md`](technical-roadmap.md), and the CP0–CP4 reports.

---

## 1. Executive technical assessment

The fork is a working Go secret scanner with a mature core: the default test suite passes 275/275 on the tested Linux environment, the application binary builds and runs with both stdlib and `gore2regex` regex engines, and the principal CLI modes, configuration precedence, allowlists, ignore/baseline behavior, decode/archive controls, reports, and exit-code overrides have been exercised. These are real capabilities, not claims derived from upstream reputation.

The baseline is nevertheless **not clean enough to treat as a green security gate without qualification**. Two High-severity operational risks are confirmed:

1. **Incomplete/failed scans can look successful** — a Git subprocess failure produced `0 commits scanned`, `[]` findings, and **exit 0** without `--exit-code 0`; unreadable files and worker/fragment errors are similarly warning-only. An operator or CI job cannot reliably distinguish complete-clean from partial/failed by exit status.
2. **Report-path handling is destructive** — validation deletes/truncates an existing target before scanning; a symlink target was truncated to 0 bytes and replaced by a regular report file.

Additional confirmed issues include diagnostic-log secret leakage despite `--redact`, path-only redaction corruption, repository-local config/ignore trust issues, malformed-regex panics, baseline self-scanning, and decoder CPU/resource costs. The decoder profile is the dominant performance constraint: 8.20 s for a 32 MiB plain corpus at default depth versus 0.286 s with decoding disabled, with 92% of CPU in the encoding-detection regex path. **Depth 0 changes detection coverage; 28.7× is a configuration comparison, not a safe optimization.**

The baseline is internally consistent after CP4A: the Makefile `git fetch --tags` executes only when `build`/`profile` expand `$(LDFLAGS)`, not at Make parse time; the CP3 filler is repeated `a\n`, so the decoder cost is regex scanning, not long Base64-filler decoding. Those corrections are documented without altering measured results.

**Consolidated state:** 19 findings: **2 High**, **9 Medium**, **4 Low**, **1 Positive**, **3 Verification gaps**. Priorities: **2 P0**, **3 P1**, **6 P2**, **7 P3**, **1 Positive**. This is a mechanism count, not a sum of checkpoint findings; the crosswalk preserves all CP0–CP4 IDs.

**Recommendation:** **READY FOR OWNER BASELINE ACCEPTANCE**, subject to the explicit limitations in §9 and the owner decisions in §11. CP5 does not mark the baseline accepted.

---

## 2. Baseline identity and provenance

| Item | Value |
| --- | --- |
| Repository | `RNC-Solutions-Group/gitleaks` |
| Working branch | `feat/issues-and-pr` |
| Canonical application SHA | `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (fork `master`) |
| Documentation HEAD at consolidation | `f437f40059f541f7a6198a025afec08da35426e2`; CP4A/CP5 commit follows |
| Branch lineage | Assessment documentation commits on top of fork `master`; no application-source divergence |
| Source equivalence | `git diff` over Go/build/config paths from baseline to HEAD is empty |
| Module requirement | `go 1.24.11` ([`go.mod:3`](../../../go.mod#L3)); 16 direct/56 indirect modules; `go.sum` 423 lines |
| Tested toolchain | Go 1.27.3 linux/amd64; CP2 tests `CGO_ENABLED=1`; CP3/CP4 performance/security binary `CGO_ENABLED=0` |
| CI toolchain | Go 1.25 ([`test.yml:23`](../../../.github/workflows/test.yml#L23)); **not tested** |
| Docker build toolchain | `golang:1.24` ([`Dockerfile:1`](../../../Dockerfile#L1)); **not built in CP5** |
| Build tags | stdlib regex default; `gore2regex` release tag; app binary verified for both; full tagged graph fails |
| Docker/runtime | CP3/CP4 used local `golang:1.27.3` digest `sha256:e519d3b0…`, 2 vCPU/1 GiB/no swap/`--network none` |
| Module cache | External `cp2-gomod` (~517 MB); verified previously; no new downloads in CP4/CP5 |
| Native Windows Go | Installed but not used for tests/measurements |
| **Untested environments** | Windows native, CI Go 1.25, release packaging, multi-arch containers, external repositories |

**Baseline identity is the source SHA `5f6ec56…`, not the latest documentation commit.** Documentation commits are assessed as non-source-changing by the `git diff` evidence above.

---

## 3. Architecture and operationally significant contracts

Detailed maps remain in [`report-CP1.md`](../CP1-architecture/report-CP1.md); this is the consolidated reference.

### 3.1 CLI and execution modes

- Root `gitleaks` exposes `git`, `dir` (aliases `file`, `directory`), `stdin`, `version`; hidden legacy `detect` and `protect`.
- `git` history/diff, `dir` filesystem, `stdin` stream, `detect` legacy paths converge on `Detector.DetectSource`.
- **`protect` is the exception:** it calls deprecated `DetectGit`, which wraps `DetectSource` with `context.Background()` ([`cmd/protect.go:61`](../../../cmd/protect.go#L61), [`detect/git.go:19-29`](../../../detect/git.go#L19-L29)). This distinction must be preserved in any refactor.
- `--pre-commit` scans the unstaged diff; `--staged` scans staged changes; pre-commit hook entries pass `--staged`.

### 3.2 Initialization and configuration precedence

`initConfig` → `initDiagnostics` → `Config.Translate` → `Detector` construction → source scan → summary/report/exit. Config precedence: `--config` > `GITLEAKS_CONFIG` > `GITLEAKS_CONFIG_TOML` > `<source>/.gitleaks.toml` > embedded default ([`cmd/root.go:132-195`](../../../cmd/root.go#L132-L195)). Repository-local config and `.gitleaksignore` are auto-loaded when no explicit config overrides them; this is a trust boundary for untrusted repositories.

### 3.3 Source and detection contracts

- Sources implement `Fragments(ctx, yield)`; directory/Git sources submit per-file work to a semaphore (40-worker hard-coded).
- Archive depth default 0; depth 1 extracts bounded content; no explicit uncompressed-size budget.
- `DetectContext` applies global allowlists, keyword prefilter, per-rule evaluation, and up to `MaxDecodeDepth` decode iterations; default depth 5.
- Findings are filtered/redacted before storage; `.gitleaksignore` and baseline suppression apply at `AddFinding`.
- **Finding order is nondeterministic** (rule map traversal); report consumers must not assume order.

### 3.4 Error/exit and report contracts

- Git/file/worker errors are often logged and absorbed; config/report errors are fatal; invalid regex panics.
- Exit: findings → `--exit-code` (default 1); clean → 0; fatal → 1; unknown flag → 126; panic → 2.
- Report formats: JSON, CSV, JUnit, SARIF, template. `--report-path -` requires an explicit format; empty CSV is zero bytes; SARIF includes `Secret` unless redacted; template funcs exclude `env`/`expandenv`/`getHostByName`.
- External commands use argv, not a shell; `--log-opts` is split on spaces.

### 3.5 Build variants

Default build uses stdlib regex; release/goreleaser builds use `gore2regex`. The tagged application binary builds and loads the default 222-rule config, but `go build/test -tags gore2regex ./...` fails in generator rule files. Distinguish app binary from full package graph.

---

## 4. Functional and correctness state

**Verified/measured core capabilities:** CLI modes and hidden legacy modes; config precedence; rule/path/commit allowlists; inline `gitleaks:allow`; finite `.gitleaksignore`; baseline suppression; decode depths 0/1/2; archive depth 0/1; symlink skip/follow; size threshold; JSON/CSV/JUnit/SARIF/template reports; exit-code overrides; `gore2regex` app binary; default 275/275 tests.

**Confirmed correctness/integrity problems:**

- False-successful/partial scan status (F01).
- Destructive report paths and symlink truncation (F02).
- Redaction/log leakage and path-only redaction corruption (F03).
- Untrusted config/ignore trust and regex panic (F04/F05).
- Baseline/config self-scan (F06).
- Global `extendDepth` validation bypass in repeated library use (F07).
- Archive/decode resource constraints (F08).
- Decoder CPU and per-file buffer costs (F09/F10).
- Nondeterministic finding order and ambiguous CLI/report contracts (F11/F12/F13).
- Build/tag/vet and tooling/CI/release gaps (F14/F15).
- Unverified #1572/archive-bomb/profile/Windows/CI/release/advisory status (F16/F18/F19).

**Not defects / retained as expected or contract:** `git --pre-commit` vs `--staged` semantics; SARIF secret reporting unless redacted; empty CSV behavior pending decision; report stdout requiring explicit format pending contract; 40-worker hard code (design limitation, not verified optimum). CP4A also corrected the Makefile parse-time claim and CP3 filler explanation.

---

## 5. Test and verification results

| Item | Result | Scope |
| --- | --- | --- |
| `go test ./...` | **275/275 pass**, 0 fail, 0 skip | Linux Go 1.27.3, CGO=1 |
| `go test -race ./...` | Pass | Exercised paths only |
| `go test -shuffle=on` | Pass | Linux |
| `gofmt -l` | 0 files | LF export |
| Coverage | **30.9%** instrumented; `cmd` 0%, `sources` 7.4%, `detect` 73.9%, `config` 79.1%, `codec` 60.3%, `report` 76.8% | Instrumented packages only |
| `go vet ./...` | **Fail** — `detect_test.go:2083`, `:2122` lostcancel | Test-only static diagnostic |
| Tagged app | Builds; default 222-rule scan pass | `gore2regex`, Linux CGO=0 |
| Tagged full graph | Build/test fail in `airtable.go:30`, `perplexity.go:16` | Generator packages |
| Windows | **NOT RUN** | Known skips/symlink/CRLF caveats |
| CI Go 1.25 | **NOT RUN** | Toolchain parity gap |
| Release packaging | **NOT RUN** | Docker/goreleaser static only |
| Advisory lookup | **NOT RUN** | No vulnerability absence claim |

The detailed capability/evidence/regression matrix is [`verification-matrix.md`](verification-matrix.md). A passing test suite is not evidence of cross-platform parity or coverage completeness.

---

## 6. Performance and stability reference

Full measurements and corrected provenance: [`report-CP3.md`](../CP3-performance-stability/report-CP3.md), [`reproduction-manifest.md`](../CP3-performance-stability/reproduction-manifest.md).

| Observation | Result |
| --- | --- |
| 32 MiB / 2,000 files, default depth | **8,196 ms** median (7,374–9,014), ~4.0 MiB/s |
| Same corpus, depth 0 | **286 ms** median, ~114.6 MiB/s → **28.7×** |
| CPU profile attribution | **92%** in decoder encoding regex; ~1.9/2 cores |
| 5 MiB encoded depth 0/1/2 | 94 / 950 / 968 ms (1-level); 118 / 991 / 1,091 ms (2-level) |
| Per-file buffer | 100 KiB per file; 207.74 MiB of 292.97 MiB allocations on 2,000 files |
| Archive | depth 1 167 ms vs ordinary 95 ms; depth 0 16 ms; 2,378 B compressed/103,800 B uncompressed |
| Size threshold | 1.1 MiB file 95 ms skipped vs 399 ms scanned; 0.9 MiB unaffected |
| 1,000 findings | ~34 MiB peak RSS; reporting within noise; generic dedupe ~+29 ms median |
| Concurrency | 40 hard-coded workers; 12–17 OS threads; stable under limits |

The plain corpus is repeated `a\n` (CP4A-02); the decoder cost is regex scanning, not long Base64-filler decoding. The 28.7× figure disables decoding and changes coverage; it is **not** a validated optimization. All workloads were bounded by 2 vCPU/1 GiB/60 s; no OOM/timeout occurred.

---

## 7. Security and operational risk summary

1. **Scan completeness/exit-status integrity (High):** Git failure/unreadable files/worker errors can yield exit 0; timeout can be delayed. CP5-F01.
2. **Filesystem/report integrity (High):** report target deletion; symlink target truncation. CP5-F02.
3. **Sensitive-data handling (Medium):** debug/trace logs and Git stderr unredacted; path-only redaction corrupts. CP5-F03.
4. **Configuration/trust boundaries (Medium):** repository-local config/ignore can alter scanning; malformed regex panics; `extendDepth` global state. CP5-F04/F05/F07.
5. **Resource availability (Medium):** decoder CPU; archive expansion without uncompressed budget; blocking reads. CP5-F08/F09/F10.
6. **Build/CI/release reproducibility (Medium):** toolchain drift; tagged graph/vet; Makefile build/profile network/mutation; CI gaps; root image; advisory status unknown. CP5-F14/F15/F18.

**Highest confirmed operational risks — retain:** false-successful exit status after incomplete scanning; destructive report-path handling.

**Uncertainties — do not treat as validated vulnerabilities or passed checks:** #1572 device/blocking exploitability; archive/decompression bombs; #2019 workload; diagnostics/profile data exposure; Windows native; CI/toolchain parity; release packaging; dependency advisories.

---

## 8. Build, CI, and release posture

- **Application binary:** reproducible from the pinned export with cached modules; default and `gore2regex` app builds pass; binary is CGO=0 in CP3/CP4.
- **Package graph:** default `go build ./...` passes; tagged full graph fails; tests pass default only.
- **Toolchains:** module 1.24.11, CI 1.25, Dockerfile 1.24, tested 1.27.3 — no parity established.
- **Makefile:** build/profile execute `git fetch --tags` via shell expansion of `LDFLAGS`; test/format/config/clean mutate but do not fetch; `make -n` does not execute backticks; no Make execution was performed.
- **CI:** test workflow installs unpinned `gotestsum@latest` and runs `go generate` (checkout mutation); gitleaks scan workflow is conditioned to the upstream repository, so it is a no-op on this fork; actions are pinned by SHA.
- **Docker/release:** runtime Alpine root with bash/git/openssh and `safe.directory=*`; no `.dockerignore`; no SBOM/signing observed; release workflow builds/pushes multi-arch images.
- **Dependency advisories:** not reviewed; checksums only.

---

## 9. Limitations and unresolved questions

- Linux/Go 1.27.3 only; Windows, CI Go 1.25, and release packaging untested.
- Warm-cache, single-host performance; no cold-cache/cross-host extrapolation.
- #1572 device/blocking behavior, archive bombs, #2019, and profile secret content untested.
- Dependency vulnerability status unknown.
- Makefile execution semantics corrected from source; no `make` binary available for a live minimal reproduction.
- CP3 fixture filler corrected to `a\n`; measured timings remain valid, but input-level causal statements are narrowed.
- No external PR was evaluated; no remediation/implementation was performed.
- Owner baseline acceptance has **not** been granted.

---

## 10. Prioritized recommendations

See [`technical-roadmap.md`](technical-roadmap.md). Order:

1. Fail-closed scan completeness/status (P0).
2. Non-destructive, symlink-safe report destinations (P0).
3. Redaction/log boundary (P1).
4. Input/config/ignore trust (P1).
5. Baseline/path self-scan exclusion (P2).
6. Resource limits/cancellation (P2).
7. Build/CI/release reproducibility (P2).
8. Coverage-preserving decoder performance (P3).
9. Contract/determinism refinements (P3).

No initiative is implemented or authorized by this report.

---

## 11. Owner acceptance considerations

The baseline is suitable as a **reference point**, not as a claim that the scanner is secure or release-ready. Owner decisions requested:

1. Accept the baseline as the pinned technical reference with the explicit limitations in §9.
2. Decide the desired **scan-status contract** (nonzero/partial status vs warning-only) before any implementation.
3. Decide **report overwrite/symlink policy** and whether atomic writes are required.
4. Decide the **redaction guarantee** across debug/trace logs and Git stderr.
5. Decide whether untrusted-repository scans require explicit `--config` and whether `.gitleaksignore` needs hardening.
6. Decide report format/order/empty-CSV/SARIF contracts.
7. Authorize selected roadmap initiatives separately; none is authorized here.

---

## 12. External PR evaluation rubric

Do not evaluate an actual external PR during CP5. Future assessments must use:

| Dimension | Evaluation question |
| --- | --- |
| Applicability | Does the change address behavior present in this fork at `5f6ec56…`? |
| Technical correctness | Does it resolve the identified problem, with evidence? |
| Compatibility | Does it preserve CLI/report/config/`protect`/exit contracts and intentional behavior? |
| Security | Does it improve or weaken detection coverage and operational integrity? |
| Regression exposure | Which modes/encodings/archives/platforms are affected? |
| Performance | Is improvement measured under the CP3 conditions and comparable coverage? |
| Build/release | Does it work across required Go versions, tags, CGO, OSes, and packaging? |
| Maintenance | New dependencies, complexity, or long-term cost? |
| Evidence | Tests/experiments sufficient to establish the outcome? |
| Disposition | Adopt, adapt, defer, reject, or needs evidence. |

Required distinctions: correct-but-irrelevant; solves defect but changes intentional behavior; apparent performance gain by reducing detection coverage; applies cleanly but regresses; requires adaptation rather than cherry-pick. A clean patch or green unit test alone is insufficient.

---

## 13. Consolidation method, finding count, and integrity

- **Method:** reconciled CP0–CP4/CP4A reports; deduplicated by mechanism; preserved every original ID in the crosswalk; retained positive/expected/contract findings separately; no new broad discovery or test rerun.
- **Count:** 19 consolidated (2 High, 9 Medium, 4 Low, 1 Positive, 3 Verification gap); 2 P0, 3 P1, 6 P2, 7 P3, 1 Positive.
- **CP4A corrections:** Makefile fetch timing (build/profile, not parse); CP3 filler `a\n` and narrowed decoder causal explanation; 28.7× coverage caveat; status records.
- **Integrity:** canonical source SHA unchanged; only `_assessment/**` documentation changed; no source/test/fixture/config/build/CI/`_backlog` modifications; no external PRs evaluated; no fixes implemented.
- **Repository status:** final Git status, commit SHA, and push result are reported in the completion response after validation.

**Recommendation:** **READY FOR OWNER BASELINE ACCEPTANCE**, with §9 limitations and §11 decisions outstanding. The baseline is not marked accepted by this report.

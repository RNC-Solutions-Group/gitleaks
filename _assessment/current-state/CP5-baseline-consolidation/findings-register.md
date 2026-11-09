# CP5 — Canonical findings register

**Status:** Complete — consolidated for owner review  
**Canonical code SHA:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`  
**Documentation HEAD at consolidation:** `f437f40059f541f7a6198a025afec08da35426e2`  
**Basis:** CP0–CP4/CP4A reports and evidence. No new runtime discovery was performed for CP5.

---

## Classification and priority legend

| Field | Values |
| --- | --- |
| **Status** | **Measured** — executed and observed; **Verified** — source-confirmed; **Inferred** — reasoned consequence; **Unverified** — unresolved; **Expected** — intentional/contractual; **Decision required** — contract/policy choice; **Positive** — control/capability verified |
| **Severity** | Impact-oriented: **High**, **Medium**, **Low**, **Positive**, **Verification gap** |
| **Priority** | **P0** highest operational urgency, **P1** security/integrity, **P2** important but narrower, **P3** contract/performance/lower risk |

Deduplication groups findings by underlying mechanism, not by checkpoint count. One consolidated finding can carry several original IDs; the crosswalk preserves every original.

---

## Consolidated findings

| ID | Title | Category | Status | Severity | Priority | Evidence | Origin | Reproduction | Impact | Follow-up |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| **CP5-F01** | Incomplete/failed scans can appear successful | Scan integrity / exit contract | Measured + Verified | **High** | **P0** | CP4 `git-fail-open`: Git fatal + “stderr is not empty”, 0 commits, `[]`, **exit 0** without `--exit-code 0`; CP2 non-root unreadable file warning + exit 0; source [`sources/git.go:414-420`](../../../sources/git.go#L414-L420), [`detect/detect.go:228-232`](../../../detect/detect.go#L228-L232), [`sources/files.go:155-162`](../../../sources/files.go#L155-L162) | CP1-005, CP1-006, CP1-013; CP2-005, CP2-006, CP2-013; CP4-001, CP4-008, CP4-009 | CP4 cases `git-fail-open`, `git-finding`, `dir-clean`; CP2 `perm-nonroot`, `timeout-stdin`; CP3 not applicable | CI gates and operators can see green after partial/failed Git scans or skipped files; timeout can be delayed | Fail-closed/partial status design; end-to-end Git/file/timeout tests; distinct complete/partial/failed classes |
| **CP5-F02** | Report-path handling deletes/truncates existing targets, including symlink targets | Filesystem / data integrity | Measured + Verified | **High** | **P0** | CP4 `report-delete`: sentinel absent; `report-symlink`: target 0 bytes, link replaced by regular 444-byte report; source [`cmd/root.go:347-358`](../../../cmd/root.go#L347-L358), [`:463-491`](../../../cmd/root.go#L463-L491) | CP1-007; CP2-003; CP4-002 | CP4 cases `report-delete`, `report-symlink`; CP2 `report-delete` | Local data loss/truncation before/without a completed scan; report aliasing risk | Atomic write + non-destructive validation; symlink policy; failure-preserves-target tests |
| **CP5-F03** | Sensitive-data redaction does not cover diagnostic logs; path-only redaction corrupts output | Sensitive data / correctness | Measured + Verified | **Medium** | **P1** | CP4 `debug-leak`: synthetic token in debug stderr despite `--redact`; redacted JSON/CSV/JUnit/SARIF/template/verbose outputs contained no token; CP2 path-only `REDACTED`-between-characters; source [`detect/codec/decoder.go:93-101`](../../../detect/codec/decoder.go#L93-L101), [`detect/detect.go:510-515`](../../../detect/detect.go#L510-L515), [`report/finding.go:78-86`](../../../report/finding.go#L78-L86) | CP1-012; CP2-004, CP2-009; CP4-003, CP4-011 | CP4 `debug-leak`, five redacted reports, `redact-verbose`; CP2-004/009 | Secrets can reach debug/trace/CI logs; path-only findings have corrupted text (no secret exposure) | Redact/gate secret-bearing logs; fix empty-secret redaction; log-level regression tests |
| **CP5-F04** | Repository-local configuration is auto-loaded and can alter/suppress scanning; malformed regex panics | Configuration / trust | Measured + Verified | **Medium** | **P1** | CP2 config precedence (source-local selected); CP2 invalid regex panic exit 2; CP4 source [`cmd/root.go:132-195`](../../../cmd/root.go#L132-L195), [`config/config.go:123-128`](../../../config/config.go#L123-L128); commit-allowlist trace logs skip but detector filters (CP2-014) | CP1-010, CP1-011, CP1-018; CP2-008, CP2-014; CP4-004 | CP2 `precedence-*`, `bad-regex`, `git-allow` | Scanning an untrusted repo can let its config suppress/alter findings or crash the scanner; misleading commit-allowlist logging | `--config`-required mode for untrusted repos; recover/validate regex; clarify commit-allowlist logging |
| **CP5-F05** | `.gitleaksignore` follows symlinks and lacks regular-file/size/entry/blocking safeguards | Input trust / availability (#1572) | Verified; finite symlink measured; device exploit **Unverified** | **Medium** (High if exploit confirmed) | **P1** | Finite symlink to regular ignore file suppressed a finding (CP4 `ignore-symlink` 0 vs control 1); source `os.Stat` + `bufio.Scanner`, no `scanner.Err()` ([`cmd/root.go:502-515`](../../../cmd/root.go#L502-L515), [`detect/detect.go:153-193`](../../../detect/detect.go#L153-L193)) | CP1-017; CP4-006 | CP4 `ignore-symlink`/`ignore-control`; no device test | Hostile repository-controlled ignore path could block/consume resources; exploitability of original #1572 device/symlink case not demonstrated | Regular-file policy, size/entry caps, timeout; #1572-specific controlled validation |
| **CP5-F06** | Baseline/config files can be rescanned, re-reporting embedded secrets | Filesystem / data exposure | Measured + Verified | **Medium** | **P2** | CP4 `baseline-second`: 2 findings from `baseline.json`, original token suppressed; CP2 same; strict path comparison [`detect/detect.go:294-299`](../../../detect/detect.go#L294-L299), baseline relative path [`detect/baseline.go:52-79`](../../../detect/baseline.go#L52-L79) | CP1-009, CP1-023; CP2-007; CP4-005 | CP4 `baseline-create`/`baseline-second`; CP2 `baseline-second` | Baseline secrets can be re-reported when baseline/config path forms don't match the skip | Realpath/canonical exclusion; baseline inside/outside tests; config path normalization |
| **CP5-F07** | `extendDepth` package global causes repeated translation to skip validation | Library correctness / global state | Measured | **Low** | **P3** | CP4 probe: first translation 223 rules; second invalid rule (no regex/path) returned no error and was accepted; source [`config/config.go:24`](../../../config/config.go#L24), [`:375-414`](../../../config/config.go#L375-L414) | CP1-010; CP4-010 | CP4 `extenddepth-probe` (disposable Go probe) | Embedded/library/repeated-init consumers can accept invalid rules; normal CLI is single-shot | Reset/encapsulate extension state; repeated-translation tests |
| **CP5-F08** | Archive/decode resource controls are bounded by options but lack uncompressed-size/cancellation budgets | Resource / availability | Measured + Verified | **Medium** | **P2** | CP3 archive depth 0→16 ms, depth 1→167 ms vs ordinary 95 ms; no explicit uncompressed cap; decode default 5 vs depth 0 8,196 ms vs 286 ms; source [`sources/file.go:58-84`](../../../sources/file.go#L58-L84), [`detect/detect.go:349-363`](../../../detect/detect.go#L349-L363) | CP1-016; CP3-001, CP3-002, CP3-005, CP3-006; CP4-007 | CP3 workloads; CP2 archive/size tests; no bomb test | Bounded, opt-in controls exist, but decompressed expansion and blocking reads lack explicit budgets; DoS potential unquantified | Uncompressed-size/entry budgets; cancellation; bomb-safe finite tests |
| **CP5-F09** | Decoder encoding-detection scan dominates plain-text CPU | Performance | Measured + source | **Medium** (performance impact high) | **P3** | CP3: 32 MiB default 8.2 s vs depth 0 0.286 s; 92% CPU in decoder regex path; encoded 5 MiB depth 0/1/2 94/950/968 ms; profile [`detect/codec/decoder.go:23-40`](../../../detect/codec/decoder.go#L23-L40) | CP1-016; CP3-001, CP3-002, CP3-003, CP3-007, CP3-011 | CP3 `files2000`, `decoder_plain`, `decode_big`, profiles | Throughput ~4–5 MiB/s with decoder vs ~115 MiB/s depth 0; depth 0 changes detection coverage | Coverage-preserving gating; positive/negative encoding corpus; no depth-0 “fix” |
| **CP5-F10** | Per-file 100 KiB buffer is allocated regardless of file size | Performance / memory | Measured + source | **Low** | **P3** | CP3 allocation profiles: 207.74 MiB of 292.97 MiB on 2,000 files from `fileFragments`; source [`sources/file.go:167-171`](../../../sources/file.go#L167-L171) | CP3-004 | CP3 `files2000`, `findings1000`, profiles | Excess allocations/GC for small files; retained heap itself small | Size-aware/pooled buffer; chunk-boundary and archive tests |
| **CP5-F11** | Finding accumulation/dedupe/report costs are measurable but not dominant at tested limits | Performance / contract | Measured | **Low** | **P3** | CP3 `findings1000` 383 ms no report vs 391/375 ms report; dense1000 generic 122 ms vs specific 93 ms; in-use heap 4.86 MiB for 1,000 findings | CP3-008, CP3-010, CP3-012 | CP3 `findings1000`, `dense*` | Report writing within noise; generic dedupe cost modest at 1,000 findings; ordering/accumulation still matters at scale | Larger bounded finding-count study if ceilings change; output-contract decision |
| **CP5-F12** | Finding order is nondeterministic | Contract / determinism | Measured | **Low** | **P3** | 8 runs alternated rule order; map traversal [`detect/detect.go:327`](../../../detect/detect.go#L327); tests use `ElementsMatch` | CP1-024; CP2-012 | CP2 `order-1..8` | Report diffs/ordering vary across runs; consumers must sort | Decide deterministic order contract; stable-sort tests |
| **CP5-F13** | CLI/report contracts are inconsistent or undocumented | Contract / usability | Measured | **Medium** | **P3** | `--report-path -` without format fatal; `.json` inference works; `--report-template` help/code mismatch; empty CSV 0 bytes; SARIF includes secrets unless redacted; `git --pre-commit` vs `--staged` semantics | CP1-008, CP1-019, CP1-020; CP2-001, CP2-002, CP2-015, CP2-016; CP4-017 | CP2/CP4 cases | Automation and operator confusion; empty outputs; contract ambiguity | Document/decide; tests for stdout format, template implication, empty CSV, SARIF/redaction |
| **CP5-F14** | Build/toolchain drift, tagged graph failure, and vet failure | Build / verification | Measured + Verified | **Medium** | **P2** | CP2: 275/275 tests pass, race pass, coverage 30.9%, `go vet` lostcancel fail, `go test -tags gore2regex ./...` fails in generator rules, app binary builds; CP4 app tag build/smoke pass; source [`airtable.go:30`](../../../cmd/generate/config/rules/airtable.go#L30), [`perplexity.go:16`](../../../cmd/generate/config/rules/perplexity.go#L16) | CP1-022; CP2-010, CP2-011; CP4-013 | CP2 commands; CP4 `build-re2`, `re2-default` | Release app binary is viable; full tagged graph/CI variant and vet gate are not clean; toolchain versions differ | Tagged package builds/tests; vet fix; pin toolchains; separate app vs graph status |
| **CP5-F15** | Makefile/CI/release tooling has mutating, network, and provenance gaps | Build/CI/tooling | Verified + Measured CP4A | **Medium** | **P2** | CP4 matrix; CP4A correction: `VERSION` backticks execute only when `build`/`profile` expand `LDFLAGS`, not at parse time; CI `gotestsum@latest`, `go generate` mutation, gitleaks workflow disabled on fork; Docker root/no `.dockerignore`; no SBOM/signing observed | CP1-020; CP4-012, CP4-014, CP4-015; CP4A-01 | Source inspection; no Make execution | Build/profile network/tag mutation; CI/dependency drift; fork security scan no-op; container/release posture | Separate read-only diagnostics; pin tools; enable real scan; non-root image; SBOM/signing |
| **CP5-F16** | Platform, security, and resource hypotheses remain unverified | Verification gap | **Unverified / NOT RUN** | **Verification gap** | **P2** | #1572 device behavior, archive bomb, #2019, profile secret content, Windows native, CI Go 1.25, release packaging, dependency advisories, protect timeout | CP0 fixture absence; CP1-021; CP4-006/007/016; CP5 scope | Not run | Cannot claim parity/vulnerability absence | Define controlled future checks; keep as unknowns |
| **CP5-F17** | Positive controls and verified core capabilities | Positive | **Positive** | **Positive** | — | CP2 275/275 tests, redacted reports, template funcs restricted, Git argv no shell, archive default 0, pinned actions; CP4 tag app build; CP3 bounded stability | CP2-018; CP4-018; CP3-009, CP3-012 | CP2/CP3/CP4 evidence | Baseline has usable controls; future changes must preserve them | Preserve in contract tests |
| **CP5-F18** | Dependency advisory status is unknown | Verification gap | **Unverified / NOT RUN** | **Verification gap** | **P2** | No live advisory lookup performed; `go.sum` checksums only | CP4-016 | — | Absence of known vulnerabilities cannot be claimed | Schedule advisory review; record results |
| **CP5-F19** | Assessment-environment limitations affect reproducibility | Verification gap | **Verified limitations** | **Verification gap** | **P3** | CP0 sibling fixture repo absent; CRLF/symlink materialization; CP3 warm-cache, Linux Go 1.27.3 CGO=0, single host; CP4 native Go unused | CP0; CP1-021; CP3-012; CP3A; CP4A-02 | CP0/CP3/CP3A | Cross-platform/cold-cache comparisons not established; fixture content corrected to `a\n` | Reproduce on pinned CI/Windows/cold cache; use manifest |

---

## Consolidated counts

| Severity | Count |
| --- | ---: |
| High | 2 (F01, F02) |
| Medium | 9 (F03, F04, F05, F06, F08, F09, F13, F14, F15) |
| Low | 4 (F07, F10, F11, F12) |
| Positive | 1 (F17) |
| Verification gap | 3 (F16, F18, F19) |
| **Total consolidated** | **19** |

This total is a **mechanism count**, not a sum of checkpoint findings. One consolidated mechanism can subsume several CP IDs; conversely, not every original CP item is a defect (architecture facts, contract decisions, and positives are retained in the crosswalk).

---

## Crosswalk — CP0

| Original | Disposition |
| --- | --- |
| CP0 canonical SHA/source equality, environment, toolchain absence at CP0 | Baseline provenance; native Go later installed; no Windows tests — F19/F14 |
| CP0 sibling fixture repo absent, CRLF/symlink materialization | Verification limitation — F19 |
| CP0 safe-command matrix/diagnostics NOT RUN | Evidence constraint — CP5 verification matrix |

## Crosswalk — CP1 (24 findings)

| Original | CP5 disposition |
| --- | --- |
| CP1-001 baseline/source equality | Baseline identity |
| CP1-002 CLI convergence; protect wrapper | Architecture summary; F01 |
| CP1-003 empty-secret redaction | F03 |
| CP1-004 commit-allowlist `continue` | F04 |
| CP1-005 error absorption | F01 |
| CP1-006 protect Background context | F01 |
| CP1-007 report target deletion | F02 |
| CP1-008 template/help mismatch | F13 |
| CP1-009 config/baseline path forms | F06 |
| CP1-010 extendDepth global | F07 |
| CP1-011 MustCompile panic | F04 |
| CP1-012 redaction logs | F03 |
| CP1-013 worker errors discarded | F01 |
| CP1-014 40-worker hard code | F09/F17 |
| CP1-015 enable-rule prefilter | Architecture/contract retained |
| CP1-016 size guard post-decode/opt-in | F08/F09 |
| CP1-017 `.gitleaksignore` loader | F05 |
| CP1-018 unused `Fragment` field | Architecture/contract retained |
| CP1-019 reporting contracts | F13 |
| CP1-020 report fd/exit lifecycle | F13 |
| CP1-021 fixture-renaming tests | F19 |
| CP1-022 gore2regex variant | F14 |
| CP1-023 baseline/config allowlist | F06 |
| CP1-024 map-order nondeterminism | F12 |

## Crosswalk — CP2/CP2A (18 findings)

| Original | CP2A class | CP5 disposition |
| --- | --- | --- |
| CP2-001 `-` format | C4 | F13 |
| CP2-002 template implication | C4 | F13 |
| CP2-003 report deletion | C2 | F02 |
| CP2-004 path-only redaction | C1 | F03 |
| CP2-005 Git fail-open | C2 | F01 |
| CP2-006 unreadable file | C4 | F01 |
| CP2-007 baseline self-scan | C2 | F06 |
| CP2-008 invalid regex panic | C1 | F04 |
| CP2-009 log redaction leak | C2 | F03 |
| CP2-010 vet lostcancel | C3 | F14 |
| CP2-011 tagged graph/app binary | C3/C7 | F14/F17 |
| CP2-012 finding order | C4 | F12 |
| CP2-013 timeout | C4 | F01 |
| CP2-014 commit allowlist log | C1 | F04 |
| CP2-015 empty CSV | C4 | F13 |
| CP2-016 pre-commit/staged | C5 | F13 |
| CP2-017 truncated archive | C6 | F16 |
| CP2-018 positive capabilities | C7 | F17 |

## Crosswalk — CP3/CP3A (12 findings)

| Original | CP5 disposition |
| --- | --- |
| CP3-001 decoder dominates plain scans | F09 |
| CP3-002 first decode pass cost | F09 |
| CP3-003 encoding regex hotspot | F09 |
| CP3-004 per-file buffer allocation | F10 |
| CP3-005 archive depth cost | F08 |
| CP3-006 size-limit skip | F08 |
| CP3-007 throughput | F09 |
| CP3-008 finding/report/dedupe cost | F11 |
| CP3-009 worker/concurrency | F17/F09 |
| CP3-010 startup | F17 |
| CP3-011 profiling attribution | F09 |
| CP3-012 baseline variation | F19 |
| CP3A reproduction manifest | Evidence reference; corrected filler/`a\n` and 28.7× caveat (CP4A-02/03) |

## Crosswalk — CP4/CP4A (18 findings + corrections)

| Original | CP5 disposition |
| --- | --- |
| CP4-001 Git fail-open | F01 |
| CP4-002 report deletion/symlink | F02 |
| CP4-003 log redaction leak | F03 |
| CP4-004 config trust/regex panic | F04 |
| CP4-005 baseline self-scan | F06 |
| CP4-006 `.gitleaksignore` trust | F05 |
| CP4-007 archive resource | F08 |
| CP4-008 unreadable file | F01 |
| CP4-009 timeout | F01 |
| CP4-010 extendDepth | F07 |
| CP4-011 path-only redaction | F03 |
| CP4-012 Makefile | F15 (CP4A-01 corrected) |
| CP4-013 build/tag/vet | F14 |
| CP4-014 CI | F15 |
| CP4-015 Docker | F15 |
| CP4-016 advisories | F18 |
| CP4-017 SARIF/report contract | F13 |
| CP4-018 positive controls | F17 |
| CP4A-01 Make timing | Corrected in F15 |
| CP4A-02 fixture filler | Corrected in F09/F19 and CP3A manifest |
| CP4A-03 28.7× coverage caveat | F09/CP5 performance section |
| CP4A-04 status records | CP5 report/README/roadmap |

---

## Preservation statement

Every original CP finding ID is mapped above. Findings that were architecture facts, expected behavior, contract decisions, or positive controls are retained in their appropriate CP5 group rather than counted as defects. Original measured results were not rewritten; CP4A corrections are recorded as interpretation changes in [`review-CP4A.md`](../CP4-security-integrations-tooling/review-CP4A.md).

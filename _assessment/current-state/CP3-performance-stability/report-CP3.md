# CP3 — Performance and stability assessment

**Status:** CP3 accepted; CP3A reproduction manifest complete  
**Assessment date:** 2026-11-09 (America/New_York)  
**Repository / branch:** `RNC-Solutions-Group/gitleaks` / `feat/issues-and-pr`  
**Assessed code SHA:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (canonical code baseline)  
**Documentation HEAD at assessment:** `7d36d8266d52e96eecc76a3aa6d3cf3c80b7f8ee`  
**Assignment:** [`assignment-CP3.md`](assignment-CP3.md)  
**Reproduction manifest:** [`reproduction-manifest.md`](reproduction-manifest.md) (CP3A; verified evidence, no benchmark rerun)

---

## 1. Executive assessment

CP3 produced a bounded, reproducible Linux performance baseline for the pinned fork. The dominant measured cost is **unconditional encoding detection in the decoder**, not rule-regex matching. On a 32 MiB directory of plain text with no encodings, the default decode depth took **8.20 s** (median; range 7.37–9.01 s), while `--max-decode-depth 0` took **0.286 s** (median; 0.283–0.316 s) — about **28.7× faster**. CPU profiles attribute **92%** of CPU samples to `codec.Decoder.Decode → findEncodedSegments → findEncodingMatches → regexp.FindAllStringSubmatchIndex`; the depth-0 profile contains no decoder samples and finishes in ~26 ms.

On a 5 MiB encoded fixture, depth 0 was **94 ms** and depth 1 was **950 ms** (≈10× slower); depth 2 was 968 ms (≈+2%). Nested encoding was similar (118 ms → 991 ms → 1,091 ms). The first decode pass is therefore the dominant step in this workload; depth 2 adds a smaller incremental pass. This is consistent with the CP1 N+1 pass interpretation but shows that at the allowed scales the **first decoder encoding scan**, not the number of rule passes, dominates.

Other measured results: per-file 100 KiB buffer allocation accounts for **207.7 MiB of 293 MiB total allocations** on the 2,000-file workload; archive depth 1 costs ~1.7× the equivalent ordinary 100-file scan (167 ms vs 95 ms), depth 0 skips it (16 ms); `--max-target-megabytes 1` skips a 1.1 MiB file about 4.2× faster than scanning it (95 ms vs 399 ms) while leaving a 0.9 MiB file unaffected; 1,000 findings peak at ~34 MiB RSS and report writing is within noise; generic deduping adds ~29 ms median at 1,000 findings. Throughput with decoding is ~4–5 MiB/s on small files, versus ~115 MiB/s with decoding disabled for the 32 MiB corpus. All workloads completed under the 2-vCPU/1-GiB/60-s limits with no OOM or timeout.

The baseline is deliberately narrow: Linux/amd64, Go 1.27.3, `CGO_ENABLED=0`, warm cache, synthetic corpora, one rule configuration, decode depth ≤2, 2,000 files/1,000 findings ceilings. It is not CI-equivalent, not a Windows measurement, and not an #2019 OOM reproduction.

---

## 2. Test-environment and provenance manifest

| Item | Value |
| --- | --- |
| Canonical source SHA | `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` |
| Documentation HEAD | `7d36d8266d52e96eecc76a3aa6d3cf3c80b7f8ee` |
| Application source equivalence | `git diff` over Go/build/config paths from baseline to HEAD is empty |
| Disposable tree | Docker volume `cp3-work:/work/src`; container `git archive` of the canonical SHA; 453 files; symlink fixture preserved; no `.git` |
| Source integrity | Extracted `go.mod` blob `3a16143e…`, `go.sum` blob `30e18339…` match the pinned Git blobs |
| Build | `CGO_ENABLED=0 go build -o /out/gitleaks .`; exit 0; 11.85 s; binary 24,509,898 bytes |
| Toolchain | `go1.27.3 linux/amd64`, `GOTOOLCHAIN=local`; no native Windows execution |
| Module provenance | Reused external `cp2-gomod` cache; `GOPROXY=off`, `GOFLAGS=-mod=readonly`; prior `go mod verify` = all modules verified; no new downloads in CP3 |
| Runtime limits | `--cpus=2 --memory=1g --memory-swap=1g --pids-limit=256 --network none`; per-workload 60 s watchdog; inputs ≤64 MiB aggregate, ≤8 MiB/file, ≤2,000 files, ≤1,000 findings, decode ≤2, archive depth ≤1 |
| Host | Windows 11, Docker 29.9.2, Linux containers on WSL2; results characterize the container, not Windows |
| Native Go note | Native Go 1.27.3 for Windows is installed at `C:\Program Files\Go` but was **not** used for measurement; no Windows tests were run |
| Cache state | Fixtures created in the volume, warmups executed; subsequent repetitions read a warm page cache |
| Measurement harness | `/proc/<pid>/VmHWM` sampled every 10 ms; wall time from `date +%s%N`; cgroup `memory.peak` recorded as a container-level cross-check; stdout/stderr and metrics kept outside the repository |

**Not comparable to CI:** Go 1.27.3 vs CI 1.25/module 1.24.11/Dockerfile 1.24; CGO disabled; Linux container vs Windows CI matrix. These results must not be presented as CI-equivalent.

---

## 3. Workload and fixture matrix

| Fixture ID | Files | Bytes | Content | Findings | Used for |
| --- | ---: | ---: | --- | ---: | --- |
| `empty` | 0 | 0 | empty directory | 0 | startup/empty scan |
| `files100` | 100 | 102,400 | 1 KiB plain files | 0 | small file-count baseline |
| `files500` | 500 | 1,024,000 | 2 KiB plain files | 0 | medium baseline |
| `files2000` | 2,000 | 32,768,000 | 16 KiB plain files | 0 | large corpus |
| `files2000_1k` | 2,000 | 2,048,000 | 1 KiB plain files | 0 | file-count vs size isolation |
| `files500_16k` | 500 | 8,192,000 | 16 KiB plain files | 0 | size isolation |
| `findings1000` | 1,000 | 1,038,000 | 1 KiB files, one synthetic token each | 1,000 | finding density/reporting |
| `dense100/500/1000` | 1 each | 1.4/7.0/14.0 KB | 100/500/1,000 synthetic token lines in one file | 100/500/1,000 | accumulation/generic dedupe |
| `encoded1` | 100 | 2,000 | base64 token, 1 level | 0 at d0, 100 at d1+ | decode semantics |
| `encoded2` | 100 | 2,800 | base64 token, 2 levels | 0 at d0/d1, 100 at d2 | decode semantics |
| `encodedbig1` | 50 | 5,001,100 | 100 KiB filler + 1-level base64 token/file | 0 d0, 50 d1+ | decode cost |
| `encodedbig2` | 50 | 5,001,500 | 100 KiB filler + 2-level base64 token/file | 0 d0/d1, 50 d2 | decode cost |
| `ordinary` | 100 | 103,800 | 1 KiB token files | 100 | archive control |
| `archive1` | 1 | 2,378 compressed / 103,800 uncompressed | tar.gz of 100 token files | 0 d0, 100 d1 | archive cost |
| `sizesmall` | 1 | 900,014 | 0.9 MiB file with token | 1 | below-threshold size |
| `sizelarge` | 1 | 1,100,014 | 1.1 MiB file with token | 1 | above-threshold size |

All token values are synthetic placeholders matching `CP3TOKEN-[0-9]{4}`; no real credentials or fixture secret values are recorded.

---

## 4. Measurement methodology

- A disposable Linux source tree was built with `CGO_ENABLED=0`; all scans used the same binary and the `token.toml`/`generic.toml` custom configurations, never the 222-rule default, to isolate the measured machinery.
- Each workload had one or two warmups followed by **three measured repetitions** (five for `decode_big`), except single-run profiling and GC-trace runs. Medians and observed ranges are reported; a single execution is not treated as a baseline.
- Wall time is process start-to-exit. Peak RSS is the maximum sampled `VmHWM`; the sampler may miss extremely short peaks, so low values for fast runs are conservative. `cgroup memory.peak` is cumulative per container and is only a cross-check, not a per-rep metric.
- CPU/alloc profiles were captured with `--diagnostics cpu,mem` on representative workloads; profiles add overhead and are not used in the timing baseline.
- `GODEBUG=gctrace=1` runs are diagnostic only and excluded from timing medians.
- All scans used `--exit-code 0` to normalize success despite findings; unexpected nonzero exits/timeouts would be recorded as failures. All workloads exited 0 and none timed out.

---

## 5. Performance results and variation

### 5.1 Startup, empty scan, and throughput

| Workload | Median wall (range) | Peak RSS median | Median throughput |
| --- | ---: | ---: | ---: |
| `version` | 17 ms (17–20) | 1.6 MiB | — |
| `empty` scan | 17 ms (16–17) | 1.8 MiB | — |
| `files100` (100 KiB) | 33 ms (32–33) | 18.2 MiB | 3.1 MiB/s |
| `files500` (1 MiB) | 188 ms (129–198) | 23.8 MiB | 5.5 MiB/s |
| `files2000` (32 MiB) | **8,196 ms (7,374–9,014)** | 27.7 MiB | 4.0 MiB/s |
| `files2000_1k` (2 MiB, 2,000 files) | 473 ms (471–502) | 28.8 MiB | 4.3 MiB/s |
| `files500_16k` (8 MiB, 500 files) | 2,167 ms (2,010–2,404) | 23.9 MiB | 3.8 MiB/s |

The large-corpus result has ~22% run-to-run variation. At matched data volume, 2,000 small files (4.3 MiB/s) and 500 larger files (3.8 MiB/s) are comparable, indicating the dominant cost is per-byte decoding rather than a file-count cliff.

### 5.2 Decode-depth cost

| Fixture | Depth 0 | Depth 1 | Depth 2 |
| --- | ---: | ---: | ---: |
| `encodedbig1` (5.0 MiB, 1-level) | 94 ms (93–98), 53.2 MiB/s | **950 ms (907–1,088), 5.26 MiB/s** | 968 ms (883–1,102), 5.17 MiB/s |
| `encodedbig2` (5.0 MiB, 2-level) | 118 ms (94–143), 42.4 MiB/s | 991 ms (973–1,079), 5.05 MiB/s | **1,091 ms (1,049–1,301), 4.58 MiB/s** |

The first decode pass costs ~10× the no-decode scan; depth 2 adds roughly 2–10% on these fixtures. The small `encoded1`/`encoded2` timing set was dominated by cache/startup noise (medians 33–99 ms, non-monotonic), so it is not used for the timing conclusion; its finding counts were correct (0/100/100 at depths 0/1/2 as expected).

### 5.3 Decoder-disabled control

| Workload | Default depth (5) | Depth 0 | Speed-up |
| --- | ---: | ---: | ---: |
| `files2000` (32 MiB) | 8,196 ms | **286 ms (283–316)** | **≈28.7×** |
| `files500` (1 MiB) | 188 ms | 91 ms (90–111) | 2.1× |
| `findings1000` (1 MiB, 1,000 findings) | 383 ms | 209 ms (196–218) | 1.8× |
| `files100` (100 KiB) | 33 ms | 113 ms (74–116) | startup-dominated/noisy |

The 32 MiB control is the clearest evidence that the decoder, not rule evaluation, dominates plain-text scanning.

### 5.4 Archive and size-limit cases

| Workload | Median wall (range) | Peak RSS | Findings |
| --- | ---: | ---: | ---: |
| `ordinary` (100 files, 103.8 KB) | 95 ms (88–96) | 19.4 MiB | 100 |
| `archive1`, depth 0 | 16 ms (16–16) | 1.7 MiB | 0 |
| `archive1`, depth 1 | 167 ms (157–227) | 17.6 MiB | 100 |
| `sizesmall` 0.9 MiB, flag 0 | 341 ms (295–344) | 14.7 MiB | 1 |
| `sizesmall` 0.9 MiB, flag 1 | 307 ms (285–332) | 14.5 MiB | 1 |
| `sizelarge` 1.1 MiB, flag 0 | 399 ms (356–439) | 15.0 MiB | 1 |
| `sizelarge` 1.1 MiB, flag 1 | **95 ms (17–97)** | 11.8 MiB | 0 |

The large-file skip is real and saves roughly 4.2× median time; the small file is unaffected. The `17 ms` fast rep shows the skip can be near-startup cost; the 95 ms median includes startup/scheduling variance. Archive depth 1 costs ~1.7× the ordinary scan for the same 100 files; depth 0 is effectively a fast skip.

### 5.5 Finding density, reporting, and generic dedupe

| Workload | Median wall (range) | Peak RSS | Findings | Notes |
| --- | ---: | ---: | ---: | --- |
| `findings1000` no report | 383 ms (319–559) | 34.1 MiB | 1,000 | accumulation baseline |
| `findings1000` JSON file | 391 ms (303–395) | 34.3 MiB | 1,000 | reporting within noise |
| `findings1000` JSON stdout | 375 ms (374–394) | 38.3 MiB | 1,000 | within noise |
| `generic100` | 35 ms (19–44) | 1.6 MiB | 100 | warm/cold noise |
| `generic500` | 33 ms (32–33) | 14.9 MiB | 500 | |
| `generic1000` | 122 ms (115–179) | 16.9 MiB | 1,000 | generic+specific comparison |
| `specific1000` | 93 ms (91–97) | 16.9 MiB | 1,000 | same file, one rule |

Generic dedupe adds ~29 ms median at 1,000 findings (~31%), but absolute cost is small; variation is high and this does not by itself establish a dominant complexity class.

---

## 6. CPU, memory, and allocation observations

### 6.1 CPU profiles

| Workload | Duration/samples | Dominant CPU path | Share |
| --- | --- | --- | --- |
| `files2000` (32 MiB, default depth) | 8.03 s / 15.40 s (191.8%) | `Decoder.Decode → findEncodedSegments → findEncodingMatches → regexp.FindAllStringSubmatchIndex` | **92.3% cumulative** |
| `encodedbig1` depth 1 | 0.99 s / 1.99 s (200.5%) | same decoder path | **92.0% cumulative** |
| `encodedbig1` depth 0 | 25.8 ms / 80 ms | `readUntilSafeBoundary`, buffer growth, aho-corasick prefilter | decoder absent |
| `dense1000` | 28.4 ms / 20 ms sampled | `detect.filter` | 100% of tiny sample |
| `archive1` depth 1 | 33.1 ms / 30 ms | regexp/bitstate and runtime | small sample |

The profile shows the scan saturating ~1.9 of the 2 available vCPUs. `files2000` GC trace recorded 67 GC cycles (last: 7→7→3 MiB heap); `encodedbig1` depth 1 had 4 GCs with ~9% GC CPU; `findings1000` had 28 GCs. GC is secondary to the decoder regex path.

### 6.2 Allocation profiles

| Workload | Total alloc (sampled) | Top allocations |
| --- | ---: | --- |
| `files2000` | 292.97 MiB | `sources.(*File).fileFragments` 207.74 MiB (70.9%); aho-corasick 33.52 MiB; `bytes.Buffer.String` 32.50 MiB; `bufio.NewReaderSize` 9.54 MiB |
| `findings1000` | 153.83 MiB | `fileFragments` 98.90 MiB; `slices.AppendSeq` 28.29 MiB; `bufio` 6.02 MiB; `regexp.allStringIndex` 6.00 MiB |
| `dense1000` | 13.73 MiB | `shannonEntropy` 2.05 MiB; `AddFinding` 1.89 MiB; `filter` 1.20 MiB |
| `encodedbig1` depth 1 | 37.60 MiB | `bytes.growSlice` 10.51 MiB; aho-corasick 7.35 MiB; `bytes.Buffer.String` 6.22 MiB; `fileFragments` 4.53 MiB |
| `archive1` depth 1 | 19.05 MiB | `fileFragments` 9.62 MiB; `slices.AppendSeq` 2.59 MiB |

Retained heap (in-use) was small: 3.08 MiB for `files2000`, 4.86 MiB for `findings1000` (including 1.54 MiB `bytes.Buffer.String` and 758 KiB `AddFinding`), and 3.84 MiB for `dense1000`.

**Source explanation:** `source.File.fileFragments` allocates a **100 KiB buffer per file** when the caller did not supply one ([`sources/file.go:167-171`](../../../sources/file.go#L167-L171)); the 2,000-file run therefore allocates ~200 MiB of buffers regardless of file size. The decoder path always calls `findEncodedSegments`, which runs the combined encoding regex over the whole fragment ([`detect/codec/decoder.go:23-40`](../../../detect/codec/decoder.go#L23-L40), [`detect/codec/encodings.go:22-43`](../../../detect/codec/encodings.go#L22-L43), [`:120-135`](../../../detect/codec/encodings.go#L120-L135)); the base64 pattern is broad (`[\w\/+-]{16,}`) and matches long alphanumeric filler, which also drives `decodeBase64` attempts and allocation.

---

## 7. Concurrency and stability findings

- The detector creates one `semgroup.Group` with a hard-coded **40-worker** limit per detector ([`detect/detect.go:121-131`](../../../detect/detect.go#L121-L131)); sources submit per-file work through it ([`sources/files.go:142-182`](../../../sources/files.go#L142-L182)).
- Observed OS thread maxima were 12–17 for directory workloads; Goroutine counts were not directly captured (no runtime goroutine endpoint was enabled). The 2-vCPU profile consumed ~1.9–2.0 cores, indicating CPU saturation rather than idle-worker starvation at these sizes.
- Throughput at matched data volume was similar across 500 and 2,000 files (~3.8–4.3 MiB/s with decoding), with no file-count cliff. Peak process RSS stayed below ~35 MiB for 32 MiB of input; all workloads exited 0 with no timeout/OOM.
- No stability failures were observed across the measured repetitions. This is **not** a race/soak test; CP2’s `-race` pass covers only exercised paths.
- The worker count was not varied (no configuration exists); no claim is made that 40 is optimal.

---

## 8. Source-linked explanations and competing hypotheses

**Primary hypothesis (supported):** plain-text scanning is dominated by unconditional encoding detection. `DetectContext` calls `decoder.Decode` after each rule pass whenever depth remains ([`detect/detect.go:349-363`](../../../detect/detect.go#L349-L363)); `Decode` always runs `findEncodedSegments` ([`detect/codec/decoder.go:23-40`](../../../detect/codec/decoder.go#L23-L40)), which runs the combined encoding regex over the fragment even when the text contains no encoding. CPU profiles and the depth-0 control (28.7× faster) support this over rule-regex or GC as the primary cause.

**Competing explanation (rule matching):** keyword prefiltering should bypass most rules for plain filler; the profile attributes 92% to the codec path, not `detectRule` regex evaluation. This makes rule matching secondary for this corpus. The profile sample counts (15.4 s over 8.03 s wall) reflect a saturated 2-vCPU container; absolute CPU shares are more reliable than wall normalization.

**Memory hypothesis (supported):** per-file buffer allocation, not finding retention, drives allocation volume. `fileFragments` is 70.9% of allocations on the 2,000-file run and retained heap is only ~3 MiB. The 100 KiB buffer is created per `File` and is particularly wasteful for 16 KiB files.

**Archive hypothesis (supported):** depth 1 adds decompression/extraction and nested scanning overhead (~1.7× ordinary); depth 0 skips entirely. The archive is small (2.4 KiB compressed, 103.8 KiB uncompressed); no decompression-bomb behavior was tested.

**Finding/dedupe hypothesis (partially supported):** source `filter` is O(n²) over findings within a fragment ([`detect/utils.go:137-166`](../../../detect/utils.go#L137-L166)); at 1,000 findings the measured extra cost is ~29 ms median with high variation. This is consistent with extra dedupe work but does not establish a dominant complexity class at the allowed ceiling.

---

## 9. Performance findings register

| ID | Finding | Status | Evidence | Impact |
| --- | --- | --- | --- | --- |
| CP3-001 | Unconditional decoder encoding scan dominates plain-text scans | **Measured + source-linked** | 8,196 ms vs 286 ms at depth 0 (28.7×); 92.3% CPU in decoder path (`files2000`) | Highest measured CPU cost; default depth 5 pays it on every fragment. **The 28.7× figure disables decoding and therefore changes detection coverage; it is not a safe optimization gain.** |
| CP3-002 | First decode pass costs ~10× on encoded 5 MiB fixtures; depth 2 adds ~2–10% | Measured | 94→950→968 ms (`encodedbig1`); 118→991→1,091 ms (`encodedbig2`) | Default depth 5 necessarily pays at least the first pass; incremental depth cost is smaller at this scale |
| CP3-003 | Encoding regex/base64 candidate handling is the CPU hotspot | Measured | `findEncodingMatches` 92% cumulative; broad base64 pattern; `bytes.growSlice` allocations | Optimization target with coverage risk |
| CP3-004 | 100 KiB per-file buffer drives allocation volume | Measured + source | 207.74 MiB of 292.97 MiB on 2,000 files; source [`sources/file.go:167-171`](../../../sources/file.go#L167-L171) | GC/alloc pressure, especially small files |
| CP3-005 | Archive depth 1 costs ~1.7× ordinary; depth 0 skips | Measured | 167 ms vs 95 ms; depth 0 16 ms; 2,378 B compressed/103,800 B uncompressed | Bounded archive overhead; skip works |
| CP3-006 | `--max-target-megabytes` skip reduces large-file work | Measured | 1.1 MiB file: 95 ms (flag 1) vs 399 ms (flag 0); 0.9 MiB unaffected | Effective pre-filter for large files; small-file cost unchanged |
| CP3-007 | Throughput with decoder ~4–5 MiB/s; depth 0 ~115 MiB/s on 32 MiB | Measured | 4.0 MiB/s vs 114.6 MiB/s; ~1.9 cores used | Primary tuning lever |
| CP3-008 | 1,000 findings / reporting is not a dominant cost; generic dedupe is measurable but small | Measured | 383 ms no report vs 391/375 ms report; generic 122 ms vs specific 93 ms at 1,000 | Lower priority than decoder |
| CP3-009 | 40-worker bound; no file-count cliff; stable under limits | Measured + source | 12–17 threads; 4.3 vs 3.8 MiB/s at 2k/500 files; no timeouts; source 40 | No evidence to change worker count |
| CP3-010 | Startup/empty scan ~17 ms | Measured | version/empty medians 17 ms | Fixed overhead; minor |
| CP3-011 | Profiles identify decoder path; depth-0 control removes it | Measured | 92% vs absent decoder samples; `encodedbig` d0 profile | Confirms causal attribution |
| CP3-012 | Baseline variation is real but conclusions are stable | Measured | files2000 7.37–9.01 s; encodedbig 0.88–1.30 s; medians used | Do not over-read single runs |

---

## 10. Ranked optimization opportunities

| Rank | Candidate | Measured / expected impact | Confidence | Benefit | Correctness/compatibility risk | Complexity | Required regression tests |
| ---: | --- | --- | --- | --- | --- | --- | --- |
| 1 | Avoid running encoding detection on fragments with no plausible encoded content (cheap gate, decode only when needed) | 28.7× on 32 MiB plain text; ~10× on 5 MiB encoded depth-1 vs depth-0 | High | Very high CPU/throughput | **High**: may miss encodings if gate is too strict; must preserve all four decoders and nested cases | Medium | All decoder formats, nested depth 0/1/2, false-negative corpus, existing `detect` tests |
| 2 | Reduce base64 candidate scope / bound candidate length / reject long filler runs earlier | Decoder path is 92% CPU; base64 candidates drive allocations | Medium-high | High | High: can miss long real base64 secrets | Medium | Base64/URL-safe, long secrets, filler-only, nested |
| 3 | Reuse/size per-file buffer (pool or file-size-aware allocation) | 207.7 MiB allocations for 32 MiB corpus | High | Moderate memory/GC | Low: buffer semantics unchanged | Low–medium | Chunk-boundary, large/small files, archives, races |
| 4 | Optimize combined encoding regex / split scanning with early exits | CPU hotspot; unknown split-specific benefit | Medium | Moderate-high | Medium: overlap/precedence semantics | Medium | Encoding precedence/overlap tests, nested encodings |
| 5 | Replace linear generic-dedupe scan with a map/index | ~29 ms median at 1,000 findings | Medium-low | Low at current scale | Low | Low | Finding-set equivalence, ordering, generic naming |
| 6 | Reporting/JSON writer optimization | Within noise vs no report | Low | Low | Low | Low | Format parity |
| 7 | Worker-count tuning | No evidence; 40 is hard-coded and untested | Unverified | Unknown | Medium (concurrency behavior) | High | Soak/race at varied counts |

No optimization was implemented. Candidate 1 must not trade detection coverage for speed without a decoder-completeness test.

---

## 11. Risks and regression-test requirements

- **Decoder gating risk:** any optimization that skips `findEncodingMatches` can silently miss percent, unicode, hex, base64, URL-safe/nested encodings and change finding parity; require positive/negative fixtures for each decoder and decode depths 0/1/2.
- **Buffer changes:** ensure chunks/safe-boundary behavior and archive contents remain identical; test files at boundaries around 100 KiB and mixed sizes.
- **Archive changes:** preserve depth semantics, inner-path formatting (`archive!inner`), and bounded extraction; do not introduce decompression-bomb exposure.
- **Size-limit changes:** preserve the opt-in `0` behavior and the current dir pre-check; test below/above thresholds and Git source (which has no pre-read cap).
- **Dedupe changes:** preserve exact finding set and ordering expectations; tests currently use order-insensitive comparisons.
- **Worker changes:** require race/soak tests; 40 is unverified as optimal.

---

## 12. CP4 handoff and unresolved questions

**CP4:** decoder-gating security tradeoff (coverage guarantees), archive decompression controls, redaction/log leakage (CP2-009), baseline self-scan/report-target behavior, `.gitleaksignore` untrusted input/#1572, `gore2regex` release parity, and build/release reproducibility. CP3’s measurements are Linux/Go 1.27.3/CGO=0 only.

**Unresolved questions:** Why does the first decoder pass dominate even for plain text? (Source says it always scans; profile confirms.) What is the minimum safe gate for encoding detection? Does depth 2’s small increment hold at larger fragments? How do worker count/queueing behave under slower I/O or larger files? What is the cost on Windows/CI Go 1.25? These are CP3-follow-up/CP4 candidates, not conclusions.

---

## 13. Repository-integrity statement

- Tracked checkout: source/tests/fixtures/config/build/CI/`_backlog` unchanged; only authorized `_assessment/**` documentation was added/edited during CP2A/CP3.
- Disposable tree: fresh `cp3-work` volume from the canonical SHA; `go.mod`/`go.sum` matched pinned blobs; no application changes; raw profiles/binaries/caches outside the repository.
- All workloads used `--network none`, `GOPROXY=off`, `-mod=readonly`, external caches, 2 vCPU/1 GiB/60 s limits; no OOM, timeout, external repository, real credential, device stream, symlink cycle, decompression bomb, #2019, or #1572 reproduction.
- Synthetic tokens only; no raw secret-like fixture values recorded.

**Recommendation:** Ready for owner review. CP4 not started.

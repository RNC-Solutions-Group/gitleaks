# CP3A — CP3 reproduction manifest

**Status:** Complete — CP3A documentation hardening; no benchmark rerun  
**Date:** 2026-11-09 (America/New_York)  
**Canonical code SHA:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`  
**Documentation HEAD at CP3A:** `e1c04989b0c3a3566b7832f14f58b1a5c322d21f`  
**Parent report:** [`report-CP3.md`](report-CP3.md)  
**Purpose:** Provide enough verified information to reproduce or compare against the CP3 baseline without reconstructing that assessment conversation. **No benchmarks were rerun for this manifest.** Values below are taken from the retained CP3 artifacts (`%TEMP%\cp3-gitleaks\out`, external to Git) and the CP3 report; any item that cannot be recovered is explicitly marked.

---

## 1. Canonical code and build identity

| Item | Value |
| --- | --- |
| Canonical application SHA | `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` |
| Source export | Container `git archive` of the canonical SHA into Docker volume `cp3-work:/work/src`; 453 files; no `.git`; symlink fixture preserved |
| `go.mod` / `go.sum` blob IDs | `3a16143eab785122ab77865b331761b5cd1f5992` / `30e1833977a6b61ace3cf52fe1145240a7a82825` |
| Build command | `CGO_ENABLED=0 go build -o /out/gitleaks .` (run in `/work/src`, `GOFLAGS=-mod=readonly`, `GOPROXY=off`, external caches) |
| Build result | exit 0; 11.85 s; binary 24,509,898 bytes |
| Binary SHA-256 | `c1fb9241c1937c65cb31381c767ac229cf56ab9cb559131389daf7f5fdede5eb` |
| Embedded build info | `go version -m` recorded in `/out/buildinfo.txt`: Go 1.27.3, `CGO_ENABLED=0`, `GOOS=linux`, `GOARCH=amd64`, module `(devel)` (no `.git` in export), dependency `h1:` hashes retained |
| Build tags | none (stdlib regex engine; no `gore2regex`) |

**Unrecoverable:** the exact Git object IDs of dependency modules were not separately enumerated; `go.mod`/`go.sum` blob IDs plus `go mod verify` and the embedded `h1:` dependency hashes are the retained provenance evidence.

---

## 2. Toolchain, image, and runtime limits

| Item | Value |
| --- | --- |
| Toolchain | `go1.27.3 linux/amd64`, `GOTOOLCHAIN=local` |
| Container image | `golang:1.27.3`; ID/digest `sha256:e519d3b0a5b285e5aa53ab9eef02bc6bbb9cb04608389b743a136a9ce75f0d1e`; Linux/amd64 |
| Host | Windows 11; Docker 29.9.2; Linux containers on WSL2; Intel Core i7-11375H (4 cores / 8 logical, 3.30 GHz max) |
| Resource limits | `docker run --cpus=2 --memory=1g --memory-swap=1g --pids-limit=256 --network none` |
| Per-workload watchdog | 60 s wall; process terminated then marked `timedout=1` (all measured workloads completed; none timed out) |
| Module cache | external Docker volume `cp2-gomod`, ~517 MB; reused |
| Build cache | external Docker volume `cp2-gocache`, reused |
| Network | disabled for all CP3 scans; module download occurred only in the earlier CP2 acquisition |
| Native Go | Go 1.27.3 Windows is installed but **not used** for CP3; no Windows measurement |

**Unrecoverable:** CPU governor/frequency state, host load, WSL2 kernel scheduler details, and exact page-cache state before each run were not captured. Runs used warm caches after warmups; this is a warm-cache baseline, not cold-start storage I/O.

---

## 3. Module and cache provenance

- Dependencies were acquired once in CP2 with network enabled into `cp2-gomod`; CP2 recorded `go mod verify` = “all modules verified” and `go.mod`/`go.sum` unchanged.
- CP3 ran with `GOPROXY=off`, `GOFLAGS=-mod=readonly`, `GOMODCACHE=/gomod`, `GOCACHE=/gocache`; no new downloads.
- Build provenance is recorded in `/out/buildinfo.txt` and the binary hash above.

---

## 4. Rule configurations and hashes

Both configs are synthetic and use only `CP3TOKEN-[0-9]{4}`-shaped values.

`token.toml` — SHA-256 `592da648f83445ad7abb83b3da23104599e98c50d1a50a7f7f25917cebdd4c41`:

```toml
title = "cp3 token"
[[rules]]
id = "cp3-token"
description = "CP3 synthetic token"
regex = '''CP3TOKEN-[0-9]{4}'''
keywords = ["cp3token"]
```

`generic.toml` — SHA-256 `89d29c32f73556a3e736764517b6186fc6359bba4dd4bf26644a7d41858d98b8`:

```toml
title = "cp3 generic"
[[rules]]
id = "generic-cp3"
description = "CP3 generic rule"
regex = '''CP3TOKEN-[0-9]{4}'''
keywords = ["cp3token"]
[[rules]]
id = "specific-cp3"
description = "CP3 specific rule"
regex = '''CP3TOKEN-[0-9]{4}'''
keywords = ["cp3token"]
```

---

## 5. Deterministic fixture-generation procedure

All fixtures were generated in `/work/fixtures` on the disposable volume. The retained generator scripts are external to Git (`fixtures.sh`, `fixtures2.sh`, `fixtures3.sh`, `fixtures4.sh`); the verified procedure is:

```bash
F=/work/fixtures
mkdir -p "$F"

# Base filler (deterministic 'a' runs)
yes a | head -c 1024 > "$F/base1k"
yes a | head -c 2048 > "$F/base2k"
yes a | head -c 16384 > "$F/base16k"
yes a | head -c 100000 > "$F/base100k"

# Plain corpora
for i in $(seq 1 100);  do cp "$F/base1k"  "$F/files100/f$i.txt"; done
for i in $(seq 1 500);  do cp "$F/base2k"  "$F/files500/f$i.txt"; done
for i in $(seq 1 2000); do cp "$F/base16k" "$F/files2000/f$i.txt"; done
for i in $(seq 1 2000); do cp "$F/base1k"  "$F/files2000_1k/f$i.txt"; done
for i in $(seq 1 500);  do cp "$F/base16k" "$F/files500_16k/f$i.txt"; done

# Finding-density corpus: one synthetic token + 1KiB filler per file
for i in $(seq 1 1000); do
  { printf 'CP3TOKEN-%04d\n' "$i"; cat "$F/base1k"; } > "$F/findings1000/f$i.txt"
done

# Dense single-file corpora
for n in 100 500 1000; do
  seq 1 "$n" | awk '{printf "CP3TOKEN-%04d\n", $1}' > "$F/dense/dense$n.txt"
done

# One- and two-level Base64
printf 'CP3TOKEN-0001' | base64 -w0 > "$F/enc1.b64"
printf '%s' "$(cat "$F/enc1.b64")" | base64 -w0 > "$F/enc2.b64"
for i in $(seq 1 100); do cp "$F/enc1.b64" "$F/encoded1/f$i.txt"; cp "$F/enc2.b64" "$F/encoded2/f$i.txt"; done

# Large encoded corpora: 100KiB filler + encoded token
for i in $(seq 1 50); do
  { cat "$F/base100k"; printf '\n'; cat "$F/enc1.b64"; printf '\n'; } > "$F/encodedbig1/f$i.txt"
  { cat "$F/base100k"; printf '\n'; cat "$F/enc2.b64"; printf '\n'; } > "$F/encodedbig2/f$i.txt"
done

# Archive control: same 100 token files, ordinary and tar.gz
for i in $(seq 1 100); do
  { printf 'CP3TOKEN-%04d\n' "$i"; cat "$F/base1k"; } > "$F/archive_src/f$i.txt"
  cp "$F/archive_src/f$i.txt" "$F/ordinary/f$i.txt"
done
tar -czf "$F/archive1/archive.tar.gz" -C "$F/archive_src" .

# Size-threshold fixtures (0.9MiB and 1.1MiB, token appended)
{ yes a | head -c 900000;  printf 'CP3TOKEN-9001\n'; } > "$F/sizesmall/small_0_9mb.txt"
{ yes a | head -c 1100000; printf 'CP3TOKEN-9002\n'; } > "$F/sizelarge/large_1_1mb.txt"
```

**Reproducibility caveat:** `yes a | head -c ...` closes the pipe early; the generator used `|| true` in the scripts to tolerate the resulting `SIGPIPE` status. The resulting file contents are deterministic (`a` repeated).

### Verified fixture hashes

Manifest hashes are `sha256sum` over the deterministic sorted `sha256sum` output of the files in each directory; archive/config hashes are over file bytes.

| Fixture | Files | Bytes | Hash |
| --- | ---: | ---: | --- |
| `empty` | 0 | 0 | `abcfa6a9d4df344d1781bc2560b5e4cdcae08b39ed303063535e7e1e926a304a` |
| `files100` | 100 | 102,400 | `cc2532209e5d50af63b797157479de8837a5cfeb055e611a3e3c97d5015050eb` |
| `files500` | 500 | 1,024,000 | `b8b36bc2906a44ca264ea161089f79d8f442e105282ef4805d62992b32d5c7da` |
| `files2000` | 2,000 | 32,768,000 | `42350239cb127d9a55c3140d6859b27edca5665ffa242a66a4dae2a2785dece9` |
| `files2000_1k` | 2,000 | 2,048,000 | `289cfc095d84d9830bff40ba33712e724892488aab124207c8f19df26229003c` |
| `files500_16k` | 500 | 8,192,000 | `44819f150fd9ca0ec18cf379aefc997db9aca9238ef89990a283ecc5f4290868` |
| `findings1000` | 1,000 | 1,038,000 | `0280f6436d426f46cdadfb1745369ed0655dd586eb7c7acb601ac5fa6227a0af` |
| `encoded1` | 100 | 2,000 | `254411575477cafe88d5fcdaeb17981462d25b1a9cd21b190dd478b2767f739c` |
| `encoded2` | 100 | 2,800 | `480a1820093f1481fab7fa2797f51d7fb920bf4e89ba94243be702d17595f9ff` |
| `encodedbig1` | 50 | 5,001,100 | `035c4e38baf07685f28ff1e2171eded491a05db4bbb0ac473c4859eb683ba774` |
| `encodedbig2` | 50 | 5,001,500 | `b1fb3f3bf3450da5abde6394a750bc3effd9b3ec4704b868f72414b4b2dfbc2b` |
| `archive1` | 1 | 2,378 compressed / 103,800 uncompressed | `41077b8732bcdaf32adcee2360c7c9d51c4058a2c6f591bc2e4db7ca5d91a2f7` (directory manifest); archive bytes SHA-256 `3fd6a8c190e64a11f20ed6f262c976783edbf9491c2d07be4a22d9577ef73c23` |
| `ordinary` | 100 | 103,800 | `1d7fd202a68aec4a9e9c877e790f3343337dd2eeddea670197d1ef779cbde4b0` |
| `sizesmall` | 1 | 900,014 | `014254ca11ea2571c84dd3c6ce51c43e7cf8bf6179f6ff8562199b343ee599fb` |
| `sizelarge` | 1 | 1,100,014 | `173cd86af36a955bd0b50cf3aba08ff29b31a512969fd0ce4a9c0f82b18bda70` |
| `dense100` | 1 | 1,400 | `6b9f6d654197a02d96ace10dfcfac5e469ddbc77b2682a634c6e5957e201e8c0` |
| `dense500` | 1 | 7,000 | `0726368837847f153da699577cfad76f84632de37fe6c352152cc545cb3ece3a` |
| `dense1000` | 1 | 14,000 | `05cf7f69e538db4bfabd3346f6f55b9315f18a54c994e0976bdbf7334d85c009` |

---

## 6. Benchmark invocation patterns

Common environment in every measured container:

```bash
GOFLAGS=-mod=readonly GOPROXY=off GOMODCACHE=/gomod GOCACHE=/gocache \
GOTOOLCHAIN=local CGO_ENABLED=0
```

Scan template (all measured scans used `--no-banner --exit-code 0`; no report unless stated):

```bash
/out/gitleaks dir --config "$CONFIG" [OPTIONS] --no-banner --exit-code 0 "$FIXTURE"
```

| Group | Invocation | Options |
| --- | --- | --- |
| startup | `/out/gitleaks version` | — |
| throughput | `dir --config token.toml` | fixture `empty`, `files100/500/2000`, `findings1000` |
| reporting | `dir --config token.toml --report-format json --report-path <external>` | file or `-` (stdout) |
| decode | `dir --config token.toml --max-decode-depth {0,1,2}` | `encoded1/2`, `encodedbig1/2` |
| decoder control | `dir --config token.toml --max-decode-depth 0` | `files100/500/2000`, `findings1000` |
| archive | `dir --config token.toml --max-archive-depth {0,1}` | `archive1` |
| archive control | `dir --config token.toml` | `ordinary` |
| size | `dir --config token.toml --max-target-megabytes {0,1}` | `sizesmall`, `sizelarge` |
| dedupe | `dir --config generic.toml` | `dense100/500/1000` |
| dedupe control | `dir --config token.toml` | `dense1000` |
| concurrency | `dir --config token.toml` | `files2000_1k`, `files500_16k` |
| profiles | `dir --config ... --diagnostics cpu,mem --diagnostics-dir <external>` | representative workloads |

**Note:** `--exit-code 0` was used to normalize exit status for performance measurement. It does **not** imply that findings are absent; finding counts were recorded from stderr and JSON where applicable.

---

## 7. Critical comparisons and expected findings

| Comparison | Input | Options | Expected findings | CP3 result |
| --- | --- | --- | ---: | --- |
| Default vs decode-off | `files2000` (32 MiB, 2,000 files) | default depth 5 vs `--max-decode-depth 0` | 0 / 0 | 8,196 ms vs 286 ms median; **28.7×** |
| Encoded depth 0/1/2 | `encodedbig1` (5.0 MiB) | depth 0/1/2 | 0 / 50 / 50 | 94 / 950 / 968 ms |
| Encoded depth 0/1/2 | `encodedbig2` (5.0 MiB) | depth 0/1/2 | 0 / 0 / 50 | 118 / 991 / 1,091 ms |
| Ordinary vs archive | `ordinary` vs `archive1` | no archive / depth 1 | 100 / 100 | 95 ms vs 167 ms |
| Archive off | `archive1` | depth 0 | 0 | 16 ms |
| Size threshold | `sizesmall` 0.9 MiB | flag 0/1 | 1 / 1 | 341 / 307 ms |
| Size threshold | `sizelarge` 1.1 MiB | flag 0/1 | 1 / 0 | 399 / 95 ms |
| Generic dedupe | `dense1000` | generic vs specific | 1,000 (specific only) | 122 ms vs 93 ms |
| Finding density | `findings1000` | no report/file/stdout | 1,000 | 383 / 391 / 375 ms |

---

## 8. Warmup, repetition, timing, and memory sampling

- Each group used 1–2 warmups, then 3 measured repetitions (5 for `decode_big`); profiling/GC runs were single-run and excluded from timing.
- Wall time: `date +%s%N` around the process; milliseconds.
- Peak RSS: `/proc/<pid>/VmHWM` sampled every 10 ms by `measure.sh`; maximum value recorded. Fast workloads may under-sample their true peak; `cgroup memory.peak` (container-cumulative) was recorded as a cross-check, not a per-rep metric.
- OS thread maximum: `/proc/<pid>/status` `Threads` sampled in the same loop; Goroutine counts were not captured.
- `GODEBUG=gctrace=1` runs are diagnostic only; their wall times are not part of the medians.
- Profiles: `--diagnostics cpu,mem` produced `cpu.pprof`/`mem.pprof`; analyzed with `go tool pprof -top` and `-sample_index=alloc_space`.

Profile analysis commands:

```bash
go tool pprof -top -nodecount=10 /out/gitleaks /out/profiles/<case>/cpu.pprof
go tool pprof -top -cum -nodecount=16 /out/gitleaks /out/profiles/<case>/cpu.pprof
go tool pprof -sample_index=alloc_space -top -nodecount=12 /out/gitleaks /out/profiles/<case>/mem.pprof
go tool pprof -sample_index=inuse_space -top -nodecount=10 /out/gitleaks /out/profiles/<case>/mem.pprof
```

---

## 9. Known reproducibility limitations

1. **Warm cache / single host:** the baseline is warm-cache on one Windows/WSL2 host; it is not a cold-start or cross-host study.
2. **Toolchain mismatch:** Go 1.27.3, `CGO_ENABLED=0` — not CI Go 1.25, not the Dockerfile’s Go 1.24, not a Windows native build.
3. **Small decode fixtures noisy:** `encoded1/2` timings were dominated by cache/order effects; only the 5 MiB encoded corpora support the decode-cost conclusion.
4. **Sampling:** 10 ms `VmHWM` sampling can miss short peaks; `cgroup memory.peak` is container-cumulative.
5. **Profiles add overhead:** profile runs are diagnostic and not part of timing medians.
6. **Ceilings:** 2,000 files, 1,000 findings, 64 MiB aggregate, 5 MiB encoded corpora, decode ≤2, archive ≤1; no extrapolation beyond these.
7. **No #2019/#1572:** no OOM/device/unbounded behavior was tested; CP3 is a bounded performance baseline, not a memory-safety proof.
8. **No security conclusion:** the 28.7× figure comes from **disabling decoding**, which removes detection coverage. It is **not** a verified safe optimization gain; any future decoder gating must preserve percent, Unicode, hex, Base64, and nested-encoding coverage.
9. **CPU frequency/governor, host load, page-cache state, and exact scheduler behavior were not captured.**
10. **Raw profiles and binaries are outside Git** (`%TEMP%\cp3-gitleaks\out` and Docker volumes); the hashes and commands above identify them, but they are not committed.

---

## 10. Critical comparison note: 28.7× is not an optimization claim

The 32 MiB default-depth versus depth-0 comparison is a **coverage-changing control**: depth 0 disables recursive decoding entirely, so encoded secrets would not be found. The measurement identifies the cost of the decoder path; it does **not** show that depth 0 is an acceptable fix. A future optimization must be evaluated on both performance and detection coverage across all supported encodings and nested depths. This distinction is carried into CP4.

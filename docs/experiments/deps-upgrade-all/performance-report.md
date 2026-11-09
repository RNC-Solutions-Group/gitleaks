# Performance report — dependency and RE2 upgrade

**Status:** Exploratory, matched-resource experiment; not a production benchmark
**Environment:** Linux/amd64, Go 1.27.3, `CGO_ENABLED=0`, 2 vCPU, 1 GiB, no swap, `--network none`, warm module/build caches
**Matched control:** `GOMAXPROCS=2` on both stacks (see §2)
**Branch:** `test/deps-upgrade-all` (dependency commit `70aaf01`)

## 1. Inputs and method

Fixtures were generated in a disposable Docker volume. The plain corpus hash matches the CP3 manifest; the encoded/archive fixtures are this experiment’s own bounded equivalents.

| Fixture | Size / shape | SHA-256 |
| --- | --- | --- |
| `files2000` | 32,768,000 bytes, 2,000 × 16 KiB files; repeated `a\n` | `42350239cb127d9a55c3140d6859b27edca5665ffa242a66a4dae2a2785dece9` |
| `encodedbig1` | 50 × ~100 KiB files, each with one Base64 `UPGRADE-0001` | `6ed8b38230337bb36a8e39e16952aa5af817efb307b7d02c8bcfce38831cfa0e` |
| `archive1` | tar.gz of 100 token files; 103,800 uncompressed bytes | `7cef108b4eb98b8a33a98eedbefe01accfe7d3fa7a63c37534768402aa6d8877` |
| `findings1000` | 1,000 files, one synthetic token + 1 KiB filler | — (CP3 shape) |

Config: synthetic `UPGRADE-[0-9]{4}` rule, keyword `upgrade-`; `--no-banner --exit-code 0`; no reports during timing runs. Each workload had one warm-up then three measured repetitions. Startup (`version`) used five repetitions. Wall time was measured with `date +%s%N`; no profiler was enabled during timing runs. Peak RSS was sampled separately from `/proc/<pid>/status` `VmHWM` at ~10 ms, one run per case, because the image has no GNU `time`.

Binaries (all built from LF Git archives, `CGO_ENABLED=0`, Go 1.27.3):

| Binary | SHA-256 |
| --- | --- |
| Baseline stdlib `gitleaks` | `3ee5940e5a543de792a69e39885236abb4185e7e0aa0e0d0f14431ba01b43ec0` |
| Baseline RE2 v1.9 `gitleaks-re2` | `7a94c067c686c419fa4cb3bfcb79ee1ec114ba575c304265de9fed5d367e25a3` |
| Upgraded stdlib `gitleaks` | `ce50826bc92c8e06ed3c6a68ade16892d5c84ca516e74dee2c843e813a6c5a57` |
| Upgraded RE2 v1.13 `gitleaks-re2` | `a4b6cab12302937051856127feead54721d6c4e4d9297d1f697fa8329a41edcf` |

## 2. Critical methodology finding — `GOMAXPROCS` changed with the Go directive

The baseline `go.mod` says `go 1.24.11`; the upgraded one says `go 1.26.0`. In a 2-vCPU cgroup this changes the Go runtime default:

- Go 1.24-era default: `containermaxprocs=0`; the process used the host-visible CPU count, not the 2-CPU cgroup limit.
- Go 1.26-era default: container-aware `GOMAXPROCS` selects 2.

A separate three-repetition check on the **same** `files2000` workload, with no explicit `GOMAXPROCS`:

| Build / setting | Run 1 | Run 2 | Run 3 | Median |
| --- | ---: | ---: | ---: | ---: |
| Baseline, default runtime setting | 6,253 ms | 7,614 ms | 6,340 ms | **6,340 ms** |
| Baseline, `GODEBUG=containermaxprocs=1` | 3,442 ms | 3,401 ms | 3,408 ms | **3,408 ms** |
| Upgraded, default runtime setting | 3,335 ms | 3,268 ms | 3,336 ms | **3,336 ms** |

This is a **Go runtime default effect**, not a dependency effect. It also means the CP3 baseline’s 8,196 ms default-depth figure is not directly comparable to this experiment’s matched measurement; CP3 did not force `GOMAXPROCS`, and its 286 ms depth-0 figure was measured under the same runtime default. This experiment does **not** rewrite CP3 and does **not** claim the dependency upgrade alone produced the apparent ~2× improvement.

All main results below therefore set `GOMAXPROCS=2` for **both** stacks. Even so, small workload results remain noisy; these are descriptive medians, not significance tests.

## 3. Matched timing results

Medians of three repetitions; parenthesized values are min–max. Negative delta means the upgraded stack was faster.

### Standard Go regexp

| Workload | Baseline | Upgraded | Delta |
| --- | ---: | ---: | ---: |
| `files2000` default depth 5 | 3,447 ms (3,406–3,468) | 3,272 ms (3,221–3,391) | **−5.1%** |
| `files2000` depth 0 | 240 ms (220–257) | 226 ms (210–307) | **−5.8%** |
| `encodedbig1` depth 0 | 87 ms (76–91) | 99 ms (87–141) | +13.8% |
| `encodedbig1` depth 1 | 578 ms (532–579) | 549 ms (541–651) | **−5.0%** |
| `encodedbig1` depth 2 | 619 ms (479–630) | 563 ms (561–685) | **−9.0%** |
| `archive1` archive depth 1 | 89 ms (86–124) | 98 ms (97–103) | +10.1% |
| `findings1000` | 245 ms (202–275) | 197 ms (196–278) | **−19.6%** |
| Startup (`version`) | 52 ms (48–54) | 53 ms (52–54) | +1.9% |

Interpretation: the standard-engine results are close to neutral. The large plain/encoded workloads improve modestly; two sub-100 ms workloads show small absolute regressions (8–12 ms) that are within the observed variance. No coverage change was made to obtain these times.

### RE2 (`gore2regex`)

| Workload | Baseline RE2 v1.9 | Upgraded RE2 v1.13 | Delta |
| --- | ---: | ---: | ---: |
| `files2000` default depth 5 | 594 ms (560–741) | 299 ms (285–312) | **−49.7%** |
| `files2000` depth 0 | 478 ms (462–503) | 238 ms (226–295) | **−50.2%** |
| `encodedbig1` depth 0 | 433 ms (374–455) | 101 ms (89–110) | **−76.7%** |
| `encodedbig1` depth 1 | 361 ms (355–396) | 110 ms (98–119) | **−69.5%** |
| `encodedbig1` depth 2 | 381 ms (367–403) | 101 ms (98–148) | **−73.5%** |
| `archive1` archive depth 1 | 415 ms (411–431) | 107 ms (103–163) | **−74.2%** |
| `findings1000` | 766 ms (675–767) | 257 ms (239–340) | **−66.4%** |
| Startup (`version`) | 358 ms (316–414) | 68 ms (60–85) | **−81.0%** |

Interpretation: the RE2 v1.13 / `wasm2go` combination is materially faster than the v1.9 / wazero combination on every workload here. The improvement is largest on encoded, archive-filled, and high-finding-count inputs. On this corpus, upgraded RE2 is also faster than upgraded stdlib on the plain corpus (299 ms vs 3,272 ms), which is a real engine difference but should not be generalized to all rule sets or small inputs.

## 4. Peak memory

One run per case, matched `GOMAXPROCS=2`; `VmHWM` sampled at ~10 ms. Values are MiB.

| Workload | Baseline | Upgraded | Delta |
| --- | ---: | ---: | ---: |
| stdlib `files2000` default | 20.75 | 22.71 | +9.4% |
| stdlib `findings1000` | 21.34 | 19.82 | **−7.1%** |
| stdlib `encodedbig1` depth 2 | 25.57 | 28.86 | +12.9% |
| RE2 `files2000` default | 46.19 | 26.90 | **−41.8%** |
| RE2 `findings1000` | 56.25 | 31.93 | **−43.2%** |
| RE2 `encodedbig1` depth 2 | 48.02 | 30.92 | **−35.6%** |

The standard engine’s memory is roughly unchanged to slightly higher; the RE2 v1.13 path is substantially lower than v1.9 on these workloads. Sampling can under-read very fast peaks; these are comparative indicators, not allocator profiles.

## 5. Coverage and validity checks

- The detection corpus produced **40/40** normalized findings for stdlib and RE2, baseline and upgraded, at decode depth 2 and archive depth 1. At decode depth 0 it produced **38/38** (the two encoded findings are intentionally absent; archive remains).
- The matched timing runs used the same rule and config as the detection case and the same fixture hashes.
- `files2000` matched the CP3 fixture hash exactly; the CP3 8,196 ms number still cannot be compared without accounting for its runtime default.
- No run timed out or OOM-killed at 1 GiB.
- No raw profiles, large fixtures, or binaries are committed; they remain external evidence.

## 6. Limitations

- Three repetitions per workload; medians are descriptive, not statistically significant. Sub-100 ms results have high relative variance.
- Warm-cache, single-host, Linux/amd64 only. No cold-start storage, cross-host, Windows performance, or release-binary measurement.
- `VmHWM` sampling at ~10 ms may miss short peaks.
- The standard-engine upgrade effect is small and partly confounded by graph/toolchain changes even with `GOMAXPROCS` matched; this experiment cannot attribute it to one dependency.
- The RE2 result is engine-version-specific (`v1.9.0` → `v1.13.0`), not a generic claim about all RE2 builds or the optional `re2_cgo` variant.
- No claim is made that decode depth 0 is an optimization. It changes detection coverage and is used only as a control.

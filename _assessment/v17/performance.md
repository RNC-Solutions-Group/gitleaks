# CP-V17 performance dossier

This is local engineering evidence for the V17 candidate, not an independent benchmark
certification. The frozen harness and its records were not modified.

## 1. Provenance and method

| Field | Value |
| --- | --- |
| Accepted V16 base | `4112c3411f943aacfe6313fe9837ccd7078ae0e3` |
| V17 implementation checkpoint | `8446960` plus the documentation commit |
| Host | Windows 11 26300, 11th Gen Intel i7-11375H, 4 cores / 8 threads |
| Toolchain | Go 1.27.3 windows/amd64; Git 2.45.2 |
| Build | `go build -trimpath -ldflags "-s -w -X=.../version.Version=..."` |
| Repetition | Warm-up for each scanner, then alternating baseline/candidate order; medians and ranges reported |
| Variance note | The Windows host showed material run-to-run variance; absolute times differ between sessions. Direction and relative ordering were stable. |

Workloads were synthetic and bounded:

| Workload | Shape | Size |
| --- | --- | --- |
| `small-many` | 4,000 text files, ~1 KB each | 4.3 MB |
| `large-few` | 40 text files, 256 KB each | 10.5 MB |
| `finding-heavy` | 400 files with repeated generic-key-shaped lines | 1.65 MB, 835 findings |
| `decode-heavy` | 200 files of Base64 assignment lines | 1.25 MB, 14,641 findings |
| `archives` | 100 ZIPs with one text member | 31 KB compressed |
| `git-history` | 200 commits, one growing text file | ~377 KB objects |

All scans used the registered default configuration, `--exit-code 0`, JSON output, and the
same flags for both binaries. Report generation was part of every measurement.

## 2. Measured result

### Filesystem

| Workload | V16 median (s) | V17 median (s) | Delta | V16 range | V17 range |
| --- | ---: | ---: | ---: | --- | --- |
| `empty` | 0.220 | 0.241 | +9.5% (startup noise) | 0.090–0.364 | 0.088–0.601 |
| `small-many` | 1.823 | 1.152 | **−36.8%** | 1.540–3.472 | 0.969–1.617 |
| `large-few` | 1.277 | 0.199 | **−84.4%** | 1.003–1.529 | 0.177–0.227 |
| `finding-heavy` | 0.803 | 0.772 | −3.9% | 0.605–1.656 | 0.577–0.909 |
| `decode-heavy` | 1.426 | 1.485 | +4.1% (within noise) | 1.131–1.573 | 1.184–1.757 |
| `archives-off` | 0.114 | 0.119 | +4.4% | 0.106–0.218 | 0.105–0.247 |
| `archives-on` | 0.195 | 0.214 | +9.7% | 0.176–1.236 | 0.170–2.194 |

`small-many`, `large-few`, and Git used 9–12 alternating repetitions after a warm-up.
The smaller cases used 5–7 repetitions. Absolute values are host-specific; the
`large-few` and `small-many` direction was repeated across sessions.

The archive comparisons are within noise for this small corpus. The large speedups are
the small-file reader and the decoder prefilter; neither changes eligible content.

### Git history

Twelve alternating runs of the 200-commit repository:

| Scanner | Median (s) | Min (s) | Max (s) |
| --- | ---: | ---: | ---: |
| V16 | 1.102 | 0.809 | 2.384 |
| V17 | 0.932 | 0.631 | 1.956 |

The warm-cache median improvement is **15.4%**. An earlier 12-pair session measured
0.750 → 0.685; the direction was the same. The profile-level
startup cost fell from 0.33 s to 0.13 s after consolidating three `git rev-parse`
subprocesses into one normal-path call. The per-byte diff-section tracker was replaced
with newline-indexed scanning and the parser read buffer was raised to 64 KB.

### Scanner CPU and peak working set

| Workload | V16 median CPU (s) | V17 median CPU (s) | V16 peak MB | V17 peak MB |
| --- | ---: | ---: | ---: | ---: |
| `small-many` | 5 | 3 | 78.4 | 54.3 |
| `large-few` | 4 | <1 | 59.2 | 57.9 |
| `git-history` | <1 | <1 | 38.6 | 39.8 |

Peak working set was sampled while the process ran. CPU values are the process total
processor time, rounded to whole seconds by the sampler; they are directional, not
fine-grained.

## 3. Root causes and corrections

### 3.1 Per-file 100 KB buffer and archive trial

Baseline profile of `small-many`:

- 753 MB total allocated for 4.3 MB of source.
- 414.9 MB (57.1%) from `sources.(*File).fileFragments`, because every file allocated
  the full 100 KB read buffer.
- 205.4 MB (28.3%) from `archives.Brotli.Match` / `isValidBrotliStream`; the name-free
  content probe tried a speculative brotli decode on ordinary text.
- 17.6 MB from `readArchivePrefix` 4 KB reads.

Corrections:

- `fragmentReadBufferSize` fits the buffer to the known file size with a 4 KB floor and
  the same 100 KB ceiling for large or unknown streams.
- With `--max-archive-depth=0` (the default), no prefix is read for a file whose name
  does not identify an archive. A file with a misleading archive name is still
  stream-verified, so text named `.zip` remains scannable.
- With traversal enabled, the name-free content trial runs only when the prefix has a
  known archive signature or binary NUL structure. This removes the speculative brotli
  decode on text while preserving ZIP/gzip/tar/tar.gz recovery for #2010.

Result: total allocation on `small-many` fell from 753 MB to 91.8 MB (−87.8%); median
wall time fell 36.8% in the warm-cache session. Archive-enabled findings remained 56/56 with identical
fingerprints in the controlled corpus.

### 3.2 Encoding regexp on ordinary source

Baseline `large-few` CPU was 94.4% inside
`codec.(*Decoder).findEncodedSegments` / `findEncodingMatches`. The combined percent,
Unicode, hex, and base64 regexp scanned every fragment even when no encoding was
possible.

`mayContainEncoding` now performs a conservative byte scan first. It returns true only
when the data contains `%`, a Unicode escape introducer, or a run of at least 16
base64-alphabet bytes containing a character that `decodeBase64` requires anyway. A
run without such a character could never produce a decoded segment, so this is not a
detection narrowing. Every accepted candidate still runs the original regexp,
precedence, and decoder.

Result: `large-few` median fell from 1.277 s to 0.199 s (−84.4%) and CPU from 4 s to
under 1 s. `decode-heavy` remained within 4% of baseline and produced the same 14,641
fingerprints; the prefilter exits early on encoded input.

### 3.3 Newline lookup and byte accounting

`newLineRegexp.FindAllStringIndex(fragment.Raw, -1)` was executed once per rule that
produced a match, and again for each required rule. It is now computed once per
fragment and passed through rule and required-rule evaluation. The finding-heavy
synthetic case has only one matching rule, so its wall-time delta is small; the change
removes a known repeated scan for multi-rule fragments.

Directory fragments carry both `Raw` and `Bytes` views of the same content. The old
`TotalBytes` path added both, roughly doubling the reported scanned size. It now counts
one view. This is a reporting-accuracy correction, not a detection change.

### 3.4 Git startup and parser overhead

The baseline Git profile spent 0.33 s of a 0.80 s probe in
`NewGitLogCmdContext`: three sequential `git rev-parse` calls for shallow, bare, and
top-level state. The normal non-bare path now asks for all three facts in one command;
the bare path falls back to the two-state query. The Git command also uses a 64 KB
buffered reader and a newline-indexed section tracker rather than a byte-at-a-time Go
loop.

## 4. Correctness and parity

Default V16/V17 comparisons on identical fixtures:

| Surface | Result |
| --- | --- |
| `small-many`, `large-few`, `archives-off` | identical (empty) |
| `finding-heavy` | 835/835 identical fingerprints |
| `decode-heavy` | 14,641/14,641 identical fingerprints |
| `archives-on` | 56/56 identical fingerprints |
| `git-history` default | 2/2 identical fingerprints |
| C1 `compressed.other` archive | 1/1 identical |
| C1 text named `.zip` | 1/1 identical |
| C1 Airtable fixture | 1/1 identical |

The Git merge and context modes are opt-in and are tested separately. Default Git
history selection and fingerprints are unchanged.

## 5. Residuals and rejected work

- **Path-allowlist matching** remains 25–34% of `small-many` CPU. The source and
  detector both evaluate the same global path regexes. A safe shared cache or a
  path-only fast path is C3/C4 work; no speculative cache was added in V17.
- **Decode-heavy memory** is finding/segment dominated. At depth 5, the 1.25 MB corpus
  peaked at ~136 MB and produced 14,641 findings; depth 0 peaked at ~41 MB. This does
  not reproduce the reported 24 GB case, but it shows why a large encoded corpus can
  scale badly. No arbitrary segment cap was introduced.
- **Brotli with no filename and no magic** remains unsupported for name-free detection.
  The gate deliberately rejects it rather than trial-decoding every text prefix.
- **`--scan-merges`** adds one merge-side finding for a branch-introduced secret. It is
  opt-in because the duplicate and fingerprint change are not compatible with an
  invisible default change.
- **`--git-context-lines`** increases Git output and scan work only when enabled. The
  default remains addition-only.
- A one-process `git rev-parse` query was selected over three queries. A fallback for
  bare repositories remains because Git rejects `--show-toplevel` there.

## 6. Reproduce

```powershell
go build -trimpath -o gitleaks.exe .
# Alternate baseline and candidate builds, same flags:
.\gitleaks.exe directory <workload> --no-banner --exit-code 0 -f json -r out.json --log-level error
.\gitleaks.exe git <repo> --no-banner --exit-code 0 -f json -r out.json --log-level error
```

Profiles used `go tool pprof` over `-cpuprofile`/`-memprofile` runs of an equivalent
in-process probe. Disposable fixtures and raw reports were kept outside tracked source.

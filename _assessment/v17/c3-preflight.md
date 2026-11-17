# CP-V17 C3 preflight

C3 remains **V18 — source readers and path identity**. V17 does not implement the
cycle. This is the revised execution proposal after the V16 architecture, V17 reader
changes, and the independent V16 detection-gap portfolio.

## 1. V17 effects on the C3 scope

- The default archive path no longer probes ordinary files when traversal is disabled,
  and content probing is signature-gated when traversal is enabled. C3 must retain
  #2010 coverage tests for ZIP/gzip/tar/tar.gz and text-mislabel negatives.
- File read buffers now depend on known size. Boundary tests must cover a secret
  spanning the old 100 KB boundary, a small known-size file, a long line, CRLF, and an
  unknown-size stream.
- Git context mode introduces `AddedLines`, but it is opt-in and does not change
  default Git fingerprints. C3 must not make it a default while path identity changes.
- `TotalBytes` now counts one view of a fragment. Reports and tests should use the same
  unit for directory and Git scans.
- C2's `--scan-merges` creates a second occurrence at a merge SHA. C3's path/fingerprint
  work must not deduplicate that occurrence by content alone.

## 2. Package disposition for V18

| Package | Upstream references | V17 state | Planned C3 action |
| --- | --- | --- | --- |
| `SCAN-03` | #1651, #2324 | Current reader extends the current line to 1 MiB and errors rather than silently splitting beyond it; partial V16 fix | **C3A:** boundary matrix at 100 KB, 125 KB, 1 MiB, CRLF, long first line, secret crossing boundary; prove no false-clean |
| `SCAN-04` | #1757, #2348 | Invalid UTF-8 differs under `gore2regex`; no V17 change | **C3A:** choose lossless replacement/segmentation policy; both backends; binary/NUL controls |
| `SCAN-05` | #2217 | Test-only stale PR | **Reject/defer:** no product defect; reuse any useful test ideas |
| `SCAN-06` | #1951 | `filetype` MIME classifier remains; V17 gate uses signatures | **C3A design spike:** MIME vs signatures; do not replace without a negative corpus |
| `SCAN-07` | #1888, #1936 | Multi-root/glob still unsupported | **Defer:** feature, not a correctness package |
| `SCAN-08` | #1821 | Directory allowlist path lacks a trailing `/` distinction | **C3B:** directory/file path identity; negative fixture `goyaml.v2` |
| `SCAN-09` | #1641 | CRLF baseline/entropy parity mostly covered; reader boundary cases remain | **C3A:** CRLF/LF findings, columns, entropy, ignore parity |
| `SCAN-10` | #1638 | Path-only rule can repeat once per fragment; no dedup | **C3B:** one path-only finding per file/rule; archive-member analogue |
| `SCAN-11` | #1287 | Root-relative matching added, reported fingerprint unchanged | **C3B:** source-relative identity option; operator decision on reported path/fingerprint |
| `SCAN-12` | #1118 | Binary skip via `filetype`; benchmark only | **C3A:** measure binary/text classification cost before changing |
| `SCAN-13` | #1059 | Absolute vs relative ignore path gap remains | **C3B:** canonical ignore identity; migration note |
| `SCAN-14` | #1058 | Jupyter image false positives have a workaround/allowlist | **Accept:** do not build notebook image parsing without evidence |
| `SCAN-15` | #914 | Notebook JSON escaping/line shape can hide generic token | **C3A:** reproduce, then decide on notebook-aware matching |
| `SCAN-16` | #2107, #2114 | Git vs dir path identity differs; C2 merge identities added | **C3B:** combine with SCAN-11/13; explicit `--source-relative-paths` design |

## 3. Recommended waves

### C3A — Boundary, encoding, and notebook correctness (low compatibility)

Packages: `SCAN-03`, `SCAN-04`, `SCAN-09`, `SCAN-15`; design spike `SCAN-06`; benchmark
`SCAN-12`.

Acceptance:

- secret spanning 100 KB/1 MiB boundaries found once under both backends;
- invalid UTF-8 either preserved safely or explicitly bounded, with no backend-specific
  false-clean;
- LF/CRLF findings, entropy, columns and baseline suppression agree;
- notebook generic-token fixture either recovers or has a documented rule-level
  limitation;
- no binary file is newly scanned merely to pass a fixture.

### C3B — Path identity and duplicate suppression (high compatibility)

Packages: `SCAN-08`, `SCAN-10`, `SCAN-11`, `SCAN-13`, `SCAN-16`; `SCAN-14` accepted.

Acceptance:

- one canonical reported path and ignore identity per source mode, with an explicit
  operator-approved migration for existing fingerprints;
- path-only findings emitted once per file/rule;
- directory vs Git, absolute vs relative, symlink, archive-member and CRLF spellings
  have a documented matrix;
- no widening of allowlists; all old ignore spellings either supported or listed as a
  migration break.

### Deferred / not in C3

`SCAN-05` (stale test-only), `SCAN-07` (multi-root/glob feature), full `SCAN-06` MIME
replacement, and notebook image decoding. These need separate product or evidence gates.

## 4. Dependencies and decisions

- **OD-8 fingerprint v2 / path identity:** C3B cannot change reported path or
  fingerprint without an operator decision. C3A can proceed independently.
- **OD-9 source-local config trust:** affects directory ignore discovery; do not combine
  with C3B path changes.
- **C2 OD-1:** if merge scanning becomes default later, C3 path identity must account
  for merge-side duplicates; keep C3 identity work content-independent.
- **C4:** allowlist path semantics may need to land before or with C3B. Recommend C3A,
  C3B, then C4 unless C4's allowlist model is required to make C3B safe.

## 5. Test/instrumentation to reuse from V17

- `sources/file_perf_test.go` for no-probe/file-buffer invariants;
- `sources/file_test.go` long-line and CRLF tests;
- `detect/git_context_test.go` added-line overlap and unchanged-context controls;
- `detect/codec/prefilter_test.go` encoded/ordinary-text controls;
- temporary reader counters and bounded synthetic corpora (kept out of tracked source).

## 6. Priority effect of the V16 gap portfolio

The 447 risk-bearing misses are dominated by identifier/context/value-shape families
(WordPress 93, database passwords 76, application passwords 68), not by directory path
identity. The 26 `SOURCE_UNSUPPORTED` misses include 10 archive members and 12 database
records; archive-member behavior is C3-relevant, database/BSON/container support is not
a filesystem-reader fix.

Therefore C3 should not absorb the accuracy backlog. It should stabilize reader and
path contracts so C6/C7 can fix identifier and provider behavior without changing
identity underneath them. The V16 gap's 10 archive-member misses and 4 binary-file
misses justify retaining the #2010 archive tests and measuring the V17 signature gate.

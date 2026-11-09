# CP4A — Corrective review

**Status:** Complete — CP4 documentation corrected; CP5 consolidation complete, awaiting owner acceptance  
**Date:** 2026-11-09 (America/New_York)  
**Canonical code SHA:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`  
**Starting documentation HEAD:** `f437f40059f541f7a6198a025afec08da35426e2`  
**Parent report:** [`report-CP4.md`](report-CP4.md)  
**Purpose:** Correct two documentation inaccuracies, reconcile their implications, and preserve the original measurements.

---

## CP4A-01 — GNU Make execution semantics

### Original claim

CP4 stated that the Makefile’s `VERSION := \`git fetch --tags && git tag | sort -V | tail -1\`` assignment executes `git fetch --tags` during Makefile parsing, and that even `make -n` might evaluate it.

### Corrected claim

- GNU Make’s `:=` immediately expands **Make** variables/functions; backticks are ordinary characters in the value, not a Make function. Make does not execute them at parse time.
- `VERSION` is referenced only through `LDFLAGS` ([`Makefile:4-5`](../../../Makefile#L4-L5)); `LDFLAGS` is expanded only in the `build` recipe ([`Makefile:23`](../../../Makefile#L23)).
- The **shell** executes the backtick command substitution when `go build $(LDFLAGS)` runs. Therefore `git fetch --tags` runs for `make build`, and for `make profile` through its `build` prerequisite. It does **not** run for Makefile parsing, `make -n`, `make format`, `make test`, `make test-cover`, `make failfast`, `make clean`, `make config/gitleaks.toml`, or `make lint`.
- `make test`/`format`/`clean`/`config` still have their own mutating/local effects (formatting, generation, deletion), but not a tag fetch.
- `go mod tidy` in `build` can also require network if the module cache is cold; that is separate from the tag fetch.

### Evidence

- Source references above; `VERSION`/`LDFLAGS` have no other references in the Makefile.
- No `make`/`mingw32-make`/`gmake` executable is available on the current host, so no live minimal-Makefile execution was performed. The correction is based on GNU Make semantics and the source references; this is labeled **Verified from source**, with the execution path inferred from standard GNU Make behavior.
- The repository’s mutating Make targets were not executed.

### Impact on CP4 findings and documentation

- Executive tooling summary corrected.
- Makefile effect matrix corrected: `VERSION`/`LDFLAGS` row no longer claims parse-time execution; `build`/`profile` explicitly show the shell fetch; `clean`/`lint` distinguish no-fetch behavior.
- CP4-012 corrected to target-specific network/side effects.
- Version-metadata row and roadmap wording corrected.
- Original CP4 measurements are unaffected; no benchmark or security observation depended on this claim.

### Remaining uncertainty

- Exact shell execution behavior was not demonstrated with a local `make` binary. It follows directly from GNU Make’s documented variable/recipe expansion and the Makefile text; a future CP5/implementation task may reproduce it in a disposable environment with Make available. No repository target was run.

---

## CP4A-02 — CP3 fixture-generation accuracy

### Original claim

The CP3A manifest described the filler as uninterrupted runs of `a` characters and suggested that broad Base64 candidates in the filler could drive decoder work.

### Verified artifact evidence

`base16k` begins with the bytes `a \n a \n a \n a \n …`; `base1k` and `base100k` are the same two-byte repetition at their respective sizes. `yes a` emits `a\n`, and `head -c N` truncates that stream; all fixture sizes are even, so files end with `\n`. The filler is **not** a continuous `a` run.

### Corrected description

- `files100/500/2000`, `files2000_1k`, `files500_16k`, `findings1000` filler, `encodedbig1/2` filler (100 KiB), and the size fixtures all use repeated `a\n` bytes.
- The plain corpus therefore contains no long alphanumeric/Base64-like runs. The measured decoder CPU cost on the 32 MiB corpus is the **encoding-detection regex scanning the fragments**, not decoding long filler candidates.
- The encoded fixtures (`encoded1/2`, `encodedbig1/2`) contain actual Base64 token strings; their decode cost and `bytes.growSlice`/`Buffer.String` allocations are separate from the plain-corpus scan cost.
- Fixture counts, byte sizes, hashes, timings, and findings are unchanged.

### Procedure corrections

- The manifest now creates every fixture directory before writes (`files2000_1k`, `files500_16k`, `encodedbig1/2`, `sizesmall`, `sizelarge`, etc.).
- It documents `set -u` vs `pipefail` behavior: with `pipefail`, `yes` can exit 141 (SIGPIPE); `|| true` is shown for reproducible generation.
- The filler/decoder relationship is explicitly stated so future readers do not attribute the measured regex cost to long Base64-like filler.

### Impact on CP3 interpretation

- CP3 report §6.2, CP3-003, and the “reduce Base64 candidates” optimization wording were corrected.
- The 92% decoder-profile attribution and all timing/memory measurements remain valid; only the input-level causal explanation was narrowed.
- The 28.7× depth-0 comparison remains a configuration comparison and is not a validated safe optimization.

---

## CP4A-03 — Performance and detection-coverage interpretation

Every performance summary now preserves the distinction:

- default decode-depth scanning and `--max-decode-depth 0` are different detection configurations;
- disabling decoding changes detection coverage;
- the ~28.7× difference is a configuration comparison, not a security-equivalent optimization;
- any future gating must preserve percent, Unicode, hex, Base64, URL-safe, and nested-encoding behavior, verified with positive and negative fixtures.

This note was added to the CP3 report, the CP3A manifest, and the CP4/CP5 evidence.

---

## CP4A-04 — Assessment status records

Status semantics are now:

| Checkpoint | Execution | Owner acceptance |
| --- | --- | --- |
| CP0 | Complete | Pending/not declared |
| CP1 | Complete | Pending/not declared |
| CP2 | Complete | Accepted with CP2A corrections (documented) |
| CP2A | Complete | Corrections recorded |
| CP3 | Complete | Accepted as executed; CP3A manifest complete |
| CP3A | Complete | Pending/not declared |
| CP4 | Complete, with CP4A corrections | Pending/not declared |
| CP4A | Complete | Pending/not declared |
| CP5 | Complete | **Not owner-accepted** |

`_assessment/README.md` and `assessment-roadmap.md` were updated to reflect execution; CP5 explicitly records that owner baseline acceptance has not been granted.

---

## Remaining uncertainties

- No local `make` binary was available to demonstrate CP4A-01 execution; semantics are source-inferred from GNU Make and the Makefile.
- The exact page-cache/CPU state and raw CP3 artifacts remain external; fixture hashes and retained metrics support reconstruction.
- The corrected filler explanation does not change measured timings; it does change the input-level causal narrative.
- CP5 is documentation consolidation; it does not add runtime evidence.

## Preservation statement

No measured result, fixture hash, file count, size, timing, memory figure, finding ID, or original CP evidence was deleted or rewritten to make the corrected interpretation look consistent. Corrections are explicit, dated, and cross-linked.

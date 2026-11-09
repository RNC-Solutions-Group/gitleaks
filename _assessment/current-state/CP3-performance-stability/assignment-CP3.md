# CP3 — Performance and stability assessment

**Status:** Complete — CP3 executed; report in owner review (combined CP2A + CP3 authorization)  
**Authority:** This assignment was created and executed under the owner’s combined CP2A/CP3 instruction; no further confirmation is required for its defined scope.  
**Code baseline:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (record the documentation HEAD separately)  
**Scope:** read-only, bounded performance/stability measurement of the pinned application source. Write only under `_assessment/current-state/CP3-performance-stability/`, with necessary `_assessment/` index/roadmap updates.

## Objective

Establish a reproducible, evidence-backed performance and stability baseline for the pinned fork: scanning throughput, decode-depth and archive cost, size-limit behavior, worker concurrency/backpressure, finding accumulation and generic deduplication, and profiling-supported optimization hypotheses. Separate source-verified architecture, direct measurements, inferences, and unverified scaling hypotheses.

**Non-goals:** no source/test/benchmark changes, no optimization implementation, no security/PR integration work, no unbounded OOM/#2019/#1572 reproduction, no external repository scans, no new toolchain or dependency acquisition.

## Execution environment and limits

- Reuse the authorized local `golang:1.27.3` image, external `cp2-gomod`/`cp2-gocache` caches, and a disposable Linux source tree; no image pulls or new dependencies.
- Build with `CGO_ENABLED=0` for a reproducible application binary; record Go version, flags, image digest, and module verification. Do not compare across toolchains.
- All benchmark runs: `--network none`, `GOPROXY=off`, `GOFLAGS=-mod=readonly`, external outputs/profiles, disposable tree.
- Per-workload ceilings: 2 vCPUs, 1 GiB memory, no additional swap, 60 s timeout, ≤64 MiB aggregate synthetic input, ≤8 MiB per file, ≤2,000 files, ≤1,000 synthetic findings, decode depth ≤2, archive depth ≤1 with bounded uncompressed size. Reduce further where practical; never increase/retry a killed workload at larger scale.
- Exclusions: unbounded OOM/device/symlink-cycle/decompression-bomb workloads; `/dev/zero`/`/dev/random`; #2019/#1572 uncontrolled reproductions; external repositories or credentials.
- Inspect test/benchmark code before execution; no broad unknown targets.

## Required investigations

### A. Reproducible baseline
Record canonical SHA, documentation HEAD, Go version/flags, OS/arch, container limits, CPU/memory, fixture definitions, cache state, exact commands/exit codes, warmup and repetition methodology. Use warmup plus ≥3 measured repetitions for completed workloads; report median and variation, not a single run.

### B. Detection throughput
Deterministic synthetic workloads varying file count, aggregate size, size distribution, finding density, rule-match density, and fragment count. Separate startup, scan, and reporting overhead where measurable. Record elapsed time, throughput, peak memory, findings, and allocation/GC observations.

### C. Decode-depth cost
Compare depths 0/1/2 on small finite encoded fixtures: detection passes, allocation growth, CPU, temporary buffers, retention, repeated decoding. Evaluate CP1's N+1 detection-pass interpretation without extrapolating.

### D. Archive-processing cost
Measure bounded archive traversal at depths 0/1 versus ordinary files; record compressed/uncompressed size, time, memory, findings, errors; investigate size-check ordering and allocation-before-rejection without archive bombs.

### E. File-size limits
Evaluate `--max-target-megabytes` immediately below/above the threshold; determine whether limits reduce work and identify pre-filtering costs. Keep correctness, performance, and security separate.

### F. Worker concurrency and backpressure
Investigate the existing 40-worker architecture with bounded file/fragment distributions; measure CPU, peak memory, observable goroutine/thread behavior, throughput, queueing/backpressure, completion, and repetition stability. Do not modify worker counts or claim configurability.

### G. Finding accumulation and generic deduplication
Use finite workloads with increasing finding counts; investigate accumulation, allocation/retention, generic-rule dedupe cost, output formatting cost, and nonlinear growth. Do not assert a complexity class without adequate evidence.

### H. Profiling and optimization hypotheses
Capture CPU/heap profiles externally where safe; rank candidates by measured impact, confidence, benefit, correctness/compatibility risk, implementation complexity, and regression-test need. Do not implement optimizations.

## Measurement discipline

For every workload record: fixture ID and generation method; input size/composition; exact command; environment/limits; repetitions; duration distribution; memory/allocation observations; findings; exit status; errors/skips/timeouts; confidence/limitations. Separate source-verified architecture, direct measurements, inferences, and unverified scaling. Do not compare different toolchains/builds/configs/cache states/platforms as equivalent. Do not treat container OOM as an application memory ceiling.

## Deliverables

Create `_assessment/current-state/CP3-performance-stability/report-CP3.md` using the shared template, including all thirteen sections required by the owner instruction: executive assessment; environment/provenance; workload/fixture matrix; methodology; results/variation; CPU/memory/allocation observations; concurrency/stability; source-linked explanations/competing hypotheses; findings register; ranked optimization opportunities; risks/regression requirements; CP4 handoff/unresolved questions; repository-integrity statement.

Keep raw profiles, binaries, caches, and large measurements outside the repository; small sanitized summaries may live in the CP3 directory.

## Hard stops

Stop the affected action if the repository/branch/SHA differs, a protected path would change, a workload exceeds limits, a tool/image/dependency is unavailable, or results would expose real secrets. Mark unavailable investigations NOT RUN and continue safe work. Do not begin CP4.

## Definition of done

Meet the shared definition of done; every measured claim has conditions and exit status; every workload has limits and repetitions; source-only/inferred claims are labeled; protected files are unchanged; CP4 handoff and unresolved questions are recorded; stop for owner review after CP3.

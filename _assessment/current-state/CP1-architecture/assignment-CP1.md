# CP1 — Architecture and execution-flow assessment

**Status:** Ready for explicit authorization; not executed  
**Authority:** Execute only when Dave explicitly instructs `Execute CP1`  
**Code baseline:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (record the current documentation-branch HEAD separately)  
**Scope:** Read-only, source-linked architecture analysis of the existing Gitleaks fork. Write CP1 findings only inside `_assessment/current-state/CP1-architecture/`; update the `_assessment/` index/roadmap as needed.

## Objective

Establish how the current program is assembled and how data, control, state, errors, and findings move through it. Produce a defensible **as-implemented architecture** that later checkpoints can test and that future external pull-request assessments can compare against this fork.

Do not implement fixes, select upstream PRs for adoption, perform a complete capability/regression audit, or claim measured performance/security outcomes. CP1 answers **how the program works in source**; CP2–CP4 validate functional, resource, and security consequences.

## Required investigations

1. **Baseline and CP0 cross-check.** Reconfirm `RNC-Solutions-Group/gitleaks`, `feat/issues-and-pr`, the pinned code SHA, and opening Git status. Read the [CP0 report](../CP0-reconnaissance/report-CP0.md), its limitations, the [roadmap](../assessment-roadmap.md), and the shared [report template](../../templates/assessment-report.md). Treat CP0's initial data-flow map as a hypothesis to confirm or correct; record any contradiction rather than silently inheriting it.
2. **Component topology and interfaces.** Trace `main.go`, `cmd/`, `config/`, `sources/`, `detect/`, `detect/codec/`, `regexp/`, `report/`, `logging/`, `version/`, and the generated/embedded rule configuration. For each major component record responsibility, entry/exit symbols, exchanged types, state ownership, dependencies, and relevant source references. Separate architectural facts from assumed intent.
3. **CLI entry paths.** Trace startup, root flag registration and defaults, initialization/configuration, command dispatch, and termination. Map the actual paths and important divergences for legacy `detect`, `dir`, `git`, `stdin`, `protect`, and `version`. Identify which flags affect downstream components and where defaults or command-specific overrides are applied.
4. **Source-to-fragment lifecycle.** Trace file/directory traversal, Git command invocation/history/diff modes, stdin/stream ingestion, archive traversal, symlink handling, and source abstractions. Record where files, Git metadata, context, filenames, and scan options become fragments; where input is accepted, excluded, skipped, or fails; and which stages can call external processes or depend on filesystem semantics.
5. **Configuration and policy resolution.** Establish the precedence and lifetime of explicit config, environment-provided config, source-local config, embedded defaults, generated rule assets, `.gitleaksignore`, rule/commit/path allowlists, baselines, enabled-rule selectors, and associated validation. Distinguish configuration loading, translation, detector construction, and per-fragment/per-rule filtering.
6. **Detection and decoding pipeline.** Follow `DetectSource`, `DetectContext`, keyword prefilters, rule evaluation, recursive decoders, file-size/depth gates, archive handling, fingerprints, deduplication, redaction, and finding construction. Specify **the order** of transformations and checks, and where mutable buffers/caches are allocated and released. Verify or correct the preliminary #2019 and #1572 source-path descriptions without reproducing them.
7. **Concurrency, lifetime, and cancellation.** Inventory goroutines, worker groups, channels, mutexes, queues, contexts, timeouts, and shared/global state in the scanning paths. Identify producer/consumer relationships, backpressure, allocation ownership, shutdown conditions, cancellation propagation, and observable error-loss boundaries. If a bound or guarantee cannot be proven from source, label it unverified; do not infer throughput or memory metrics.
8. **Errors, results, and output contracts.** Trace errors from source readers, Git subprocesses, archive/decoder/rule operations, ignore/config loading, and reporters to logging, findings, CLI return values, exit codes, and output streams. Map reporter selection and JSON/CSV/JUnit/SARIF/template behavior at the interface level, including redaction and report-file creation. Distinguish errors that terminate, log-and-continue, are suppressed, or require runtime verification.
9. **Build and platform variants.** Document code-level consequences of build tags (stdlib regex versus `gore2regex`), generated/embedded configuration, and relevant Windows/Linux differences. Carry CP0's toolchain, CRLF, symlink, and fixture-renaming caveats forward; do not treat an available local binary as the pinned baseline.
10. **Assessment handoff.** Produce the dependency/coupling and architectural risk register, proposed runtime-verification questions for CP2/CP3/CP4, and an explicit list of interfaces/invariants that future external PR comparisons must check. Use the backlog only to locate research leads; do not reassess, reclassify, or implement backlog items.

## Big questions

- What exact control/data path does each scanning mode take from CLI input to output and exit status? Where do the paths converge or differ?
- Which components own the rules, fragment contents, decoder state, shared caches, cancellation signals, and findings at each stage?
- Where can an input be skipped, duplicated, transformed, or silently discarded? What evidence establishes the ordering of guards and transformations?
- When one fragment or subprocess fails, how does that error propagate through concurrent sources and what can the caller actually observe?
- Which assumptions differ for Git scans, directory scans, standard input, pre-commit protection, archives, and symlinks?
- Which architectural contracts depend on build tags, OS behavior, command defaults, or generated configuration?
- What would an external PR need to preserve or deliberately alter to improve this baseline without introducing a regression?

## Deliverables

Create `_assessment/current-state/CP1-architecture/report-CP1.md` from the shared report template. Include, either inline or in clearly linked supplementary documents **inside this CP1 directory**:

1. **Source-pinned component map:** ownership and public/internal interfaces, supported by real paths/symbols/line anchors. Show a compact topology diagram if useful.
2. **Execution-path matrix:** each CLI scanning mode, its source, configuration, detector entry, reporter/exit path, and meaningful differences.
3. **Data-flow and state-lifetime map:** fragment → validation/filtering → decoding/rules → findings → reporting; record where caches, buffers, and mutable state live.
4. **Configuration/precedence matrix:** input source, parser/translator, effective override/default, validation, and consumer.
5. **Concurrency/cancellation/error-flow register:** originating operation, propagation path, handling outcome, and evidence/uncertainty; distinguish verified source paths from unrun behavior.
6. **Reporting contracts and build-variant map:** output implementations, return/exit behavior, redaction boundaries, embedded assets, and regex engine variants.
7. **Coupling, invariants, and verification backlog:** source-linked architectural risks and targeted CP2/CP3/CP4 experiments. Cite relevant `_backlog/` records without changing them; preserve the separate #2019 and #1572 pathways.
8. **Integrity and handoff:** assessed SHA, exact investigative commands and exit status, before/after Git status, protected-file change statement, unresolved questions, and readiness recommendation for CP2.

A supplementary document should exist only when it materially improves navigation; the CP1 report remains the canonical entry point. Update `_assessment/README.md` to link the report when it exists and reflect the true CP1 state. Update the roadmap checkpoint status only when evidence supports the change.

## Evidence and verification rules

- Label significant claims **Verified (source)**, **Measured (executed)**, **Inferred**, or **Unverified**. A call graph established by source inspection is not runtime proof.
- Cite symbols with paths and line references **at the pinned code SHA** where practical; identify uncertainty about dynamic configuration or platform-dependent behavior.
- Cross-check CP0 findings, especially the absence of native Go/module cache, unverified executable baseline, fixture-mutating tests, Windows symlink materialization, build-tag drift, and the provenance of `gitleaks:rnc-base`.
- Do not overstate negative evidence: “not found in inspected first-party code” is different from “cannot happen,” and small fixtures are not automatically safe.
- Explicitly identify coverage questions; lack of direct `cmd/*_test.go` files does not by itself prove the CLI has no indirect test coverage.
- If a question requires a reproducible runtime check, record a bounded test proposal with prerequisites and assign it to CP2/CP3/CP4. Do not substitute estimates, stale images, or upstream issue claims for observations.

## Hard stops and continuation rules

- Follow root `AGENTS.md`, `_assessment/AGENTS.md`, and the common read-only boundary. **Do not modify source, tests, fixtures, config, build/release files, dependencies, CI, or `_backlog/`.** Do not alter Git refs, worktrees, remotes, or branch configuration.
- No toolchain/image installation, image pull, network dependency hydration, `make` targets, `go generate`, `go fmt`, `go mod tidy`, source regeneration, or external PR integration work is authorized by CP1.
- No unbounded workloads or reproduction of #2019, #1572, malicious symlink/device inputs, or other unsafe fixtures. Do not run the test suite in the working checkout: CP0 found in-place fixture renames/deletions. Disposable-copy testing requires separate authorization and controls.
- Do not require a working Go build for this source-level checkpoint. Where runtime evidence is unavailable, document the limitation and continue independent architecture analysis.
- Stop the affected action if the repository/branch differs, protected files would be modified, source mutations occur, or sensitive fixture values would enter a report. Do not reset/delete evidence to hide changes.
- **Stop after delivering CP1.** Do not create CP2 assignments, edit backlog dispositions, install tools, or begin a later checkpoint without a separate instruction.

## Definition of done

Meet the shared [`../../README.md`](../../README.md#shared-definition-of-done) checklist and address every required investigation and deliverable, or explicitly document a justified omission. Each major execution path and cross-component edge must have traceable source evidence; each important error/cancellation claim must state what is proven and what still needs runtime validation. Confirm the protected source is unchanged, report CP2 readiness/limitations, and stop for Dave's review.

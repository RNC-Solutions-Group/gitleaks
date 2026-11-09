# Architecture, capabilities, and technical baseline — assessment roadmap

**Status:** CP0 executed 2026-11-09; report awaiting owner review; CP1 assignment prepared, not executed  
**Repository:** `RNC-Solutions-Group/gitleaks`  
**Working branch:** `feat/issues-and-pr`  
**Starting branch SHA (pre-documentation):** `6cd8c033512ba0de86f895cbced6a8f13b7cc531`  
**Canonical code baseline:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (fork `master`; application source identical at `6cd8c03` and at the branch HEAD assessed in CP0)  
**Purpose:** Characterize the fork as it exists, then provide evidence for external PR-to-baseline compatibility assessments.

## Mission and non-goals

Produce a reliable architectural map, capability inventory, correctness and regression baseline, performance/stability profile, security/integration assessment, and prioritized technical findings. Examine implementation directly; supporting documentation and the upstream backlog are leads, not proof.

This is **not** a mandate to fix defects, refactor, adopt PRs, change scanner rules, or make the code pass tests. Document observed limitations, including failures. Source remains read-only. Assessment documents can be written under `_assessment/` only.

## Checkpoints and decision gates

| CP | Investigation | Required outcome | Gate |
| --- | --- | --- | --- |
| **[CP0](CP0-reconnaissance/assignment-CP0.md)** | Reconnaissance: identity, Git lineage, package/test/tooling inventory, safe command matrix, sibling test repo | Reproducible starting manifest and evidence report | Stop for review before CP1 |
| **[CP1](CP1-architecture/assignment-CP1.md)** | Architecture: entry points, component ownership, data flow, Git/dir/stdin sources, decoding, detector, config, rule evaluation, concurrency, outputs | Source-linked architecture and execution maps; coupling and risk hypotheses | Assignment ready; execute only upon explicit authorization; stop before CP2 |
| **CP2** | Capabilities/correctness: CLI behavior, feature matrix, expected vs observed behavior, failure modes, test/fixture coverage | Capability register, tested behavior, known defects and regression gaps | Stop for review before CP3 |
| **CP3** | Performance/stability: reproducible workloads, benchmarks, CPU/memory/allocations, concurrency and large-input behavior | Measurement conditions and baseline metrics; ranked, evidence-backed optimization hypotheses | Stop for review before CP4 |
| **CP4** | Security/integrations and developer tooling: sensitive-data exposure, config/allowlists, dependencies, Git/hooks, reports, Makefile/script side effects, build/release behavior and reproducibility | Security, integration, and build-tooling assessment with verified behaviors, risks, and reproducibility limitations | Stop for review before CP5 |
| **CP5** | Baseline consolidation: reconcile CP0–CP4, cross-check contradictions, index living architecture documentation | Dated baseline, findings register, verification matrix, prioritized investigation roadmap | Owner accepts baseline before PR comparisons |
| **Later** | External upstream/third-party PR comparison | Applicability, correctness, compatibility, risks, benefits, and disposition **relative to this fork** | Separately authorized per PR/package |

## CP4 emphasis — developer tooling and build reproducibility

In addition to the security and integration scope above, explicitly inspect the existing `Makefile`, related scripts, Go module/toolchain configuration, CI workflows, and release tooling **without modifying them**. CP4 must establish:

- Which targets perform purely observational testing versus formatting, code generation, dependency edits, cleanup, or other persistent writes. Trace prerequisite targets and shell expansion rather than assuming a target name describes all side effects.
- Which commands implicitly access external systems (for example, tag fetching in version resolution), require network or cached dependencies, or depend on remote services.
- Whether clean build/test procedures can run reproducibly from a pinned SHA with `-mod=readonly`, defined tool versions, build tags, and external output locations.
- Differences between local development, CI, supported operating systems, and release/version workflows; characterize failures and constraints without claiming untested platforms work.
- A proposed separation between **non-mutating diagnostics** and **explicitly mutating maintenance/release targets**, with tradeoffs and migration risks. Recommendations are future implementation candidates, not authority to edit source or build tooling on this branch.

Preserve the distinction between a demonstrable command side effect, a suspected reproducibility risk, and a measured defect. CP4 should assess the existing Makefile, **not remove, disable, or neutralize it** during the baseline investigation.

## Questions this assessment must answer

1. What does the fork actually implement today, and what behavior is intentional or contractual?
2. Which code paths carry input through detection, filtering, and reporting; where are the principal failure boundaries?
3. What is verified by existing automated tests, what can be reproduced against `../gitleaks-test-repo`, and what remains untested?
4. Which defects or bottlenecks are observed, which are only suspected, and under what conditions?
5. Which architectural constraints and invariants should an external PR comparison respect?
6. What evidence would prove that a proposed external fix improves, regresses, or simply does not apply to this baseline?
7. Which development, build, test, and release operations are reproducible and which have implicit writes or external dependencies?

## Evidence and baseline discipline

- Pin a **code commit SHA** separately from documentation commits. Subsequent assessment-document commits must not silently move the technical baseline.
- For each benchmark, record machine/OS, Go version, invocation, input corpus, repetitions, relevant settings, and baseline result; never substitute estimates for measurements.
- Distinguish verified code facts, observed measurements, architectural inferences, and unresolved questions. Source links should include paths/functions and, where feasible, line references at the assessed SHA.
- Record pre-existing failures without attempting source fixes. Do not silently discard failed or unavailable tests.
- Verify documentation against the actual code; later evidence may overturn the baseline. Preserve historical assessment snapshots instead of rewriting past measured results.
- Link back to `_backlog/` rather than copying its issue/PR text or changing its dispositions. External PRs may contain unsafe sample credentials; do not commit raw leaked values.

## Handoff and acceptance

Each checkpoint uses [`../templates/checkpoint-assignment.md`](../templates/checkpoint-assignment.md) and the common [definition of done](../README.md#shared-definition-of-done). Create its report using [`../templates/assessment-report.md`](../templates/assessment-report.md), summarize the evidence, identify unresolved questions, then stop. The next checkpoint requires explicit authorization.

**First authorized activity on a cold start:** none. Report your understanding and wait for a specific checkpoint instruction (e.g., `Execute CP1`).

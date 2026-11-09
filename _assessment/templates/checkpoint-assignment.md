# CPX — [Checkpoint title]

**Status:** Draft / Ready / In progress / Complete / Blocked  
**Authority:** Execute only upon explicit instruction for this checkpoint  
**Scope:** [Exactly what may be examined and the only authorized output paths]

## Objective

[What this checkpoint establishes and why it matters for the eventual baseline or external PR reconciliation. State non-goals.]

## Required investigations

1. [Specific code paths, systems, behaviors, or metrics to investigate.]
2. [Specific reproducible tests or commands and their safe execution conditions.]
3. [Required comparisons to earlier verified checkpoint evidence.]

## Big questions

- [What must be understood or disproved?]
- [What assumptions or competing explanations must be challenged?]
- [Which facts would change a later PR compatibility decision?]

## Deliverables

- [Exact report path, using `report-CPX.md` within this checkpoint directory.]
- [Required source-linked maps, matrices, measurements, or registers.]
- [Evidence that must be recorded, including negative and unverified results.]
- [Index/roadmap updates needed inside `_assessment/`.]

## Hard stops and continuation rules

- Stop actions that would violate the root `AGENTS.md` write/authority boundaries.
- Do not conceal a failed test, missing tool, unsupported hypothesis, or unknown baseline relationship.
- Continue independent safe investigations when one tool or dataset is unavailable.
- Never start the next checkpoint without explicit authorization.

## Definition of done

Meet the common [`../README.md`](../README.md#shared-definition-of-done) requirements and any additional checkpoint-specific acceptance criteria. Summarize the evidence, gaps, and next-checkpoint readiness, then stop for owner review.

# CP5 — Final technical baseline consolidation

**Status:** Complete — CP5 executed; baseline awaiting owner acceptance (combined CP4A + CP5 authorization)  
**Authority:** Owner assignment of 2026-11-09; no intermediate approval required.  
**Code baseline:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (canonical application SHA; documentation HEAD recorded separately)  
**Scope:** reconciliation and consolidation of CP0–CP4 evidence into an owner-reviewable baseline. No new broad discovery, no source changes, no remediation, no external PR evaluation.

## Objective

Produce one internally consistent technical characterization of the fork at the canonical SHA, suitable for evaluating future external PRs. Identify verified behavior, defects, risks, limitations, contracts, and the evidence required for change evaluation. Preserve historical measurements and correct interpretations explicitly rather than rewriting results.

## Required work

1. **Baseline identity/provenance** — repository/branch, canonical SHA, documentation HEAD, lineage, source/build/config equivalence, module/toolchain requirements, tested and untested environments, build tags, Docker/CI configs, cache provenance.
2. **Architecture/interfaces** — CLI modes; initialization/config; Git/dir/stdin/protect paths; fragments/readers; decoder/archive; rules/findings; allowlists/ignore/baseline; redaction/reporting; concurrency/global state/errors/cancellation; build variants. Summarize rather than duplicate CP1; preserve the `protect` distinction.
3. **Canonical findings register** — reconcile CP0–CP4/CP4A findings, deduplicate mechanisms, classify confirmed defects vs measured limitations vs intentional behavior vs contract decisions vs hardening opportunities vs unverified hypotheses vs positive controls. Crosswalk every original CP ID.
4. **Verification/regression matrix** — CLI modes/exit contracts, Git errors, file-read failures, config precedence, rules/allowlists, ignore/baseline, decode, archives, symlinks, size, timeout, redaction/logging, reports, ordering, concurrency, build variants, platforms/CI/release. Record expected/observed/evidence/environment/test type/gap/PR evidence.
5. **Performance/stability reference** — CP3/CP3A baseline with corrected filler/provenance; throughput, decode, CPU, buffers, archives, size, findings/dedupe, concurrency, variability. Rank optimization by impact/complexity/coverage risk/memory/tests/confidence; never present depth-0 as free.
6. **Security/operational baseline** — six risk groups: scan completeness/exit; filesystem/report; sensitive data; configuration/trust; resource availability; build/CI/release. Retain the two highest risks; preserve #1572, archive bomb, #2019, profile exposure, Windows, CI parity, release, advisory uncertainties.
7. **Prioritized roadmap** — problem, findings, severity, dependencies, implementation boundary, regression tests, compatibility risk, acceptance criteria, external-PR candidacy.
8. **External PR evaluation rubric** — applicability, correctness, compatibility, security, regression exposure, performance, build/release, maintenance, evidence, disposition. Distinguish irrelevant, behavior-changing, coverage-reducing, regression-introducing, and adaptation-requiring PRs. Do not evaluate any actual PR.

## Deliverables

- `assignment-CP5.md`
- `report-CP5.md`
- `findings-register.md`
- `verification-matrix.md`
- `technical-roadmap.md`
- `_assessment/README.md` and `assessment-roadmap.md` status updates.

## Evidence rules

All findings traceable to source or prior checkpoint evidence; no invented figures/vulnerability status/platform parity; distinguish source-verified vs measured vs inferred vs unverified; preserve original measurements and CP IDs; no silent omission; validate links; no source/protected-file changes.

## Hard stops

Stop on branch/SHA conflict, protected-file change, unexpected source divergence, unsafe workload, or missing evidence. Do not begin external PR comparison, implementation, remediation, backlog edits, or CP5 acceptance on the owner's behalf. Stop after CP5 for owner review.

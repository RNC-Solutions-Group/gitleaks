# CP-V15 — Course-of-action comparison

## 1. COAs

### COA 1 — Core-engineering-first (operator's approximate plan)

Complete the ~88 core packages, perform one architecture audit, then accuracy, then
provider modernization.

**Strengths:** clear sequencing; keeps detection changes away from a moving core; maps
to existing category ownership.

**Weaknesses:** treats 88 heterogeneous packages as one phase; delays confirmed high-risk
false negatives (#2010, #2000, #2257, #1028); puts the architecture audit after the
changes most likely to create duplication; leaves release identity and independent
accuracy until late.

### COA 2 — Risk-and-dependency-first (recommended, with gates)

Fix confirmed false negatives and shared correctness defects; consolidate the control
plane; finish Git/source/path and config/report compatibility; establish release
identity; run generic accuracy with an independent oracle; repair existing provider rules;
expand providers only after a rule-quality gate.

**Strengths:** removes known silent misses first; architecture work happens before four
waves of related edits; accuracy has a stable reader/decoder and testable scan setup;
release engineering is not held hostage to the provider backlog; low-value work is
rejected early.

**Weaknesses:** requires operator decisions up front (merge policy, module identity,
minVersion, harness access); reorders category ownership; some config/report features wait.

### COA 3 — Architecture-and-release-first

Refactor the scanner and settle `/v2`, release, SBOM and signing before detection work.

**Strengths:** stable API and release early; fewer integration churn points later.

**Weaknesses:** a large refactor without first fixing known defects risks preserving the
wrong behavior; delays security fixes; architecture can over-fit an unverified detection
model; release attestation before accuracy certification produces a well-signed scanner
with unknown detection quality.

### COA 4 — Accuracy-first

Implement the 51 accuracy packages and 61 provider rules immediately, then core and
release.

**Strengths:** fastest visible detection improvement; addresses user-facing misses.

**Weaknesses:** changes behavior on top of an unstable config/report/path model; without
an independent oracle, “improvement” is unproven; provider PRs conflict; increases
regression and maintenance burden; leaves fail-closed/release risks untouched.

### COA 5 — Minimal-maintenance / reject-most

Fix only confirmed high-severity defects, reject the rest, and maintain a small scanner.

**Strengths:** lowest cost and risk; good fit if RNC's product is a narrow internal
scanner.

**Weaknesses:** underuses the accepted platform work; leaves real provider/UX gaps; may be
too narrow for an independent product. It is a viable fallback if the operator does not
want a multi-cycle program.

## 2. Comparison

Scale: **++ strong, + acceptable, 0 neutral, − weak, −− unacceptable**.

| Criterion | COA 1 | COA 2+ | COA 3 | COA 4 | COA 5 |
|---|:---:|:---:|:---:|:---:|:---:|
| Correctness improvement (silent misses) | + | **++** | 0 | ++ | + |
| Security/operational risk reduction | + | **++** | + | − | + |
| Engineering complexity | 0 | **+** | − | − | **++** |
| Review/CI burden | 0 | + | − | − | **++** |
| Repeated work/integration churn | − | **+** | + | − | **++** |
| Compatibility impact | + | 0 | − | − | **++** |
| Architectural debt | − | **+** | ++ | − | ++ |
| Release readiness | − | **+** | **++** | − | + |
| Long-term maintainability | 0 | **++** | + | − | + |
| Detection accuracy confidence | − | **++** | − | 0 | − |
| Uses accepted platform investment | + | **++** | + | + | − |

COA 2+ wins because it is the only course that improves the three things that matter in
the right order: **correctness evidence, a stable internal model, and a verifiable release**.
It is not the lowest-effort option; COA 5 is. It is the best option if RNC-Mod is intended
to be an independently maintained security product rather than a small fork.

## 3. Why not the operator's COA 1

The operator's Phase II envelope (88 packages) is roughly the right amount of core work.
The problem is composition and order:

- It includes 27 configuration packages, many of which are feature requests or policy
  decisions, not core engineering.
- It defers confirmed false negatives until after a large batch of lower-risk work.
- It places one architecture audit after the waves that create the most coupling.
- It treats detection accuracy and provider rules as later phases even though some
  detection defects are core correctness.
- It does not resolve module/release identity early enough for a stable API.

A corrected COA 1 would be close to COA 2+, but then it is no longer “core first”; it is
risk-first with core compatibility waves.

## 4. Why not COA 3

A big-bang refactor is the wrong first move for a scanner with confirmed behavior defects
and no independent oracle. Refactoring before C1 would risk preserving the current
`withinColumns`, path and archive behavior behind cleaner abstractions. Refactoring before
C3 would move the reader/decoder while its false negatives are still being characterized.

COA 3's useful elements are adopted: AG1 after C1, AG2 after C6, and C9 release
engineering after AG1. The rejected element is “refactor everything before fixing known
defects.”

## 5. Why not COA 4

Accuracy-first is attractive because the product is a secret scanner. But a provider-rule
PR is not accuracy evidence, and a finding-count change is not a precision/recall result.
The current generic rule, path identity, extension merge and report columns are shared
control paths. Changing them in 50 independent PRs would create a conflict matrix worse
than the current backlog. C6 must be a campaign with one grammar revision and an oracle,
not 51 patches.

## 6. Recommendation

**Adopt COA 2+.** Sequence:

1. C1 — seven confirmed corrections.
2. AG1 — control-plane/config/report/Git seam.
3. C2 — merge coverage and runtime stability.
4. C3 — source/path identity.
5. C4 — configuration model.
6. C5 — reporting/fingerprint.
7. C6 — generic accuracy with independent oracle.
8. C7 — existing provider repair.
9. AG2 — detection/API/generator gate.
10. C8 — provider expansion (mostly triage).
11. C9 — release identity/artifacts.
12. C10 — elective triage.
13. VER — independent verification.

C9 can run in parallel after AG1; it does not need to wait for C6/C7. A release candidate
may be produced after C1–C5 + C9, but a **stable v2 claim** should wait for C6/C7, AG2 and
VER.

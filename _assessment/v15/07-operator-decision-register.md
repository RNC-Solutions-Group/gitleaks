# CP-V15 — Operator decision register

No V15 action implements any decision below. The register separates decisions that block
the next assignment from decisions that can safely remain deferred.

## 1. Immediate decisions (needed before or during C1/AG1)

| ID | Decision | Why now | Recommendation |
|---|---|---|---|
| OD-1 | **Merge coverage policy.** Keep default `git log -p`, or scan first-parent merge patches by default? | C2 cannot implement a default without a policy. V14 measured real merge-only recovery plus duplicates/new fingerprints. | Make first-parent coverage the default for secret scanning, with an explicit opt-out, duplicate-reporting note and baseline/ignore migration. If the operator rejects the default, approve a prominent comprehensive-history mode and stop calling the default complete. |
| OD-2 | **Approve AG1 refactor** of `cmd`/config/report-plan/Git seam with byte/fingerprint parity tests. | The next four waves otherwise edit the same 279-line function and global config. | Approve a semantics-preserving gate; stop on any parity failure. |
| OD-3 | **Provider-keyword exception for Airtable PAT.** Change `keywords = ["airtable"]` to the self-signaling prefix. | It is a default-rule behavior change, even though it fixes a confirmed miss. | Approve the narrow keyword/regeneration change; add bare-token and near-miss negatives. |
| OD-4 | **Independent harness access and adjudication.** Who runs the frozen oracle, and who owns new holdouts? | C6 cannot claim accuracy from Go fixtures. C1 can proceed with local reproductions, but accuracy campaigns cannot. | Name an independent operator/reviewer; keep the frozen oracle unmodified; require separate adjudication of any new corpus. |
| OD-5 | **Scope of the first program wave.** Approve C1 + AG1 before C2/C4/C6. | Prevents the old “next category” drift. | Approve C1 and AG1 as the first two assignments; do not authorize C6 yet. |

## 2. Decisions that block later cycles

| ID | Decision | Blocks | Recommendation |
|---|---|---|---|
| OD-6 | **`minVersion` policy.** Keep warning-only, add a fail-closed compatibility floor, or map RNC product versions to an embedded config-schema version? | C4 (config model), C9 (release tag/base), public config compatibility | Add an explicit `configCompatibilityVersion` mapping and make the policy testable. Do not compare an RNC `v2` product directly to upstream `8.25.0`. |
| OD-7 | **Module identity.** Retain upstream `github.com/zricethezav/gitleaks/v8`, or publish RNC-owned `github.com/RNC-Solutions-Group/gitleaks/v2`? | C9, AG2 public API freeze, `go install`/proxy proof | Choose COA2 `/v2` only if RNC will publish and support a Go module and verify proxy/checksum resolution. Otherwise choose COA1 and document “CLI/container distribution only” honestly. Do not claim both identities. |
| OD-8 | **Fingerprint v2.** Approve keyed versioned grammar, key source/rotation, dual-read vs v2-only, and migration. | C5, any path/merge identity migration | Approve a keyed v2, dual-read compatibility and explicit v1 limitation; otherwise keep v1 and stop promising content-safe suppression. |
| OD-9 | **Source-local config trust.** Keep compatible discovery, or make secure CI the default for untrusted targets? | C4 secure-default work; release guidance | Keep current default for compatibility, but require `--no-local-config` in RNC CI templates and document the residual. A secure-by-default change is a separate migration. |
| OD-10 | **Release attestation authorization.** Which artifacts, registries, signing identity, SBOM format and consumer verification? | C9 publication | Authorize only after C9 builds the pipeline; require a separate consumer verification before calling a release certified. |
| OD-11 | **Provider rule policy.** Default-on vs opt-in experimental rules; rule ID/rename/deprecation; severity model. | C7/C8 | Default-on only for exact, tested formats; opt-in for immature providers; no global severity invention; preserve rule-ID migration notes. |
| OD-12 | **Elective feature scope.** Which optional features, if any, are worth their code/maintenance burden? | C10 | Reject the unsafe/low-value list; defer multi-format/context until AG1; require a written benefit case for any survivor. |

## 3. Safely deferred decisions

| Decision | Why it can wait | Revisit |
|---|---|---|
| Permanent shell regression test for #1900 | Functional behavior is verified; no product defect remains | Before a release that advertises the profile helper, or in the VER gate |
| Keyed-grouping key lifecycle / Windows ACL enforcement | Grouping is opt-in and already documents external key storage; not a core scanner blocker | Before claiming grouped reporting as an operational product feature, or when a managed key service is proposed |
| Two-file crash atomicity | Known residual; handled failures are tested; no transaction design yet | If a release promises atomic primary+companion output, or an operator incident demonstrates the window |
| Report final-publication TOCTOU / parent-rename temp orphan | Narrow, owner-restricted temp residual; no reproduced substitution | With a directory-handle/transaction design, not an ordinary patch |
| RE2 constrained-runtime release behavior (#2017) | Needs a release artifact and RLIMIT_AS reproduction; not a source build issue | C9 release gate |
| Public Go API stability | Depends on OD-7 and AG2 | AG2, before stable v2 |
| First-parent duplicate suppression | Duplicate findings are acceptable for security coverage; automatic dedup changes identity | Only with a separate content/occurrence identity design |
| Provider severity model | Current SARIF level is opt-in and constant; no per-rule model | If a consumer requires it; otherwise do not invent one |

## 4. Architecture and release decisions already resolved by evidence

These are not open operator choices; they are conclusions from V15:

- **Do not retain the V14A 1,000-candidate aggregate cap.** V14B removed it; V13
  compatibility is restored.
- **Do not add a second report writer for multi-format output.** Use one validated report
  plan first.
- **Do not accept an unkeyed secret hash.** #1162 remains rejected; keyed grouping is the
  only implemented privacy-preserving alternative.
- **Do not adopt #2023 as a wholesale detect rewrite.** Its size/newline semantics are
  regression-prone; take only measured components if C6 needs them.
- **Do not adopt the npm/installer/Actions examples as-is.** They fail exit propagation,
  checksum verification or trust-boundary review.
- **Do not claim detection accuracy from the Go suite.** The independent harness remains
  a release gate.

## 5. Structural review gate

### When

- **AG1:** immediately after C1, before C2–C5.
- **AG2:** after C6/C7, before C8/C10 code and before a stable public API.

### Subsystems

AG1: `cmd` control plane, config loading, report plan/publication, Git command/parser seam.
AG2: `detect` rule/location/allowlist/required-rule pipeline, report model, generator
invariants, exported API surface.

### Measurements

- function length and branch proxy (not LOC targets);
- number of responsibilities per function;
- fatal/exit calls in testable paths;
- duplicate destination/validation paths;
- coverage of core packages;
- allocation/CPU profiles for scan/reader paths;
- number of filesystem identity calls in hot paths;
- parity: findings JSON, CSV, SARIF, console redaction, fingerprints, report bytes,
  permissions, exit codes, both regex backends.

### Parity proof

Before/after goldens and a differential run over the existing fixture corpus; a
behavioral difference is a failure unless explicitly in the gate's scope. A refactor that
makes tests faster but changes a fingerprint is not acceptable.

### Before stable public Go API

Resolve OD-7 (module), freeze `report.Finding` JSON fields, config export surface and
detect exported types; decide whether `sources.Fragment` is part of the API. If not,
remove it from the public report model or hide it.

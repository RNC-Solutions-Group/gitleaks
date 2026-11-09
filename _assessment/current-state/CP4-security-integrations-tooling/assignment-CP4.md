# CP4 — Security, integration, and developer-tooling assessment

**Status:** Complete — CP4 executed; report in owner review (combined CP3A + CP4 authorization)  
**Authority:** Owner instruction of 2026-11-09; no intermediate approval required for this scope.  
**Code baseline:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (record documentation HEAD separately)  
**Scope:** read-only source analysis plus bounded, synthetic runtime checks against a disposable Linux source tree. Write only under `_assessment/current-state/CP4-security-integrations-tooling/`, with authorized `_assessment/` index/roadmap updates.

## Objective

Produce a source-linked and, where safely possible, runtime-verified assessment of sensitive-data protection; scan completeness and failure signaling; filesystem/symlink/ignore-file/archive handling; configuration and allowlist trust; Git/hooks/CI integration; report safety; dependency/build/release reproducibility; and Makefile/script side effects. Identify existing protections, confirmed defects, realistic risks, architectural limitations, and unresolved questions. Do not repair, harden, refactor, optimize, or integrate fixes.

## Threat model and trust boundaries

**Untrusted inputs:** Git histories/diffs, directory contents, stdin, paths, symlinks, archives, configuration files, rule definitions/regexes, baseline files, `.gitleaksignore`, report destinations, CLI arguments/environment variables.

**Protected assets:** secret material; scan completeness/integrity; report correctness; filesystem integrity; CPU/memory availability; CI pipeline reliability; build/release provenance.

For each input, map the first-party processing path, trust boundary, failure behavior, and whether protection is first-party or delegated to a library/subprocess.

## Required investigations

- **A. Threat model/trust boundaries** — inputs, assets, boundaries; first-party vs library/subprocess behavior.
- **B. Scan integrity/fail-open** — revisit CP2-005/006/013/014/017; error propagation through CLI, `DetectSource`, Git/file fragment generation, archives, workers, cancellation, summary, exit status. Reproduce CP2-005 **without `--exit-code 0`**; classify complete/no-findings, complete/findings, partial, failed, timed-out.
- **C. Sensitive-data exposure/redaction** — normal findings, path-only findings, debug/trace logs, decoder diagnostics, `gitleaks:allow` logs, Git stderr, verbose output, JSON/CSV/JUnit/SARIF/template reports, diagnostics/profiles, terminal streams. Revisit CP2-004/009; determine whether `--redact` protects claimed outputs. Synthetic identifiers only.
- **D. Filesystem/report destinations** — creation/truncation/deletion, validation order, failure handling, symlink resolution, path normalization, relative/absolute paths, reports inside scanned dirs, baseline self-scan, permission failures, TOCTOU. Revisit CP2-003/007. Runtime checks only in uniquely named disposable files.
- **E. `.gitleaksignore`/untrusted filesystem input** — loading, size, symlink following, blocking reads, path containment, validation, errors, resource exhaustion (#1572). No device symlinks; finite regular-file controls and static analysis only; separate confirmed behavior/reproduced risk/inferred exploitability.
- **F. Archive/decoding/resource controls** — archive detection/traversal, compressed vs uncompressed checks, nested depth, allocation-before-rejection, cancellation/error propagation, encoded segments, buffer sizing, DoS pathways; relate to CP3; decoder-gating safety consequences and coverage guarantees.
- **G. Configuration/rule safety** — precedence, local vs explicit config, env config, extension/recursion, global `extendDepth`, malformed regex, allowlists, baseline, rule selection, repeated translation; revisit CP2-008/AR-06; distinguish malformed from adversarial config.
- **H. Git/hooks/CI** — local, pre-commit/staged, CI, GitHub Actions, security gates, containers; exit propagation, partial status, stderr/logs, cancellation, permissions, paths, Git subprocesses, network dependencies, log exposure; identify misleading-green-CI configurations. Do not trigger workflows or external scans.
- **I. Makefile/scripts/tooling** — target/effect matrix: name, prerequisites, commands, purpose, writes, network, toolchain assumptions, Git effects, cleanup, CI/release implications, diagnostic suitability. Trace prerequisite chains and shell expansion. Do not run mutating targets; do not assume `make -n` is side-effect-free. Do not run `make test/build/format/clean/profile`, `go generate`, `go mod tidy`.
- **J. Dependencies/build tags/release** — `go.mod`/`go.sum`, toolchain, provenance, Dockerfile, workflows, release scripts, version metadata, CGO, tags, platform behavior, `gore2regex` vs stdlib; revisit CP2-010/011 and CP3 toolchain/CGO differences. Separate application-binary build, complete package graph, tests, tagged tests, CI execution, release packaging. Bounded builds/smoke checks in disposable copies; no claim of untested parity or vulnerability absence without advisory lookup.

## Safety boundaries

Reuse `golang:1.27.3`, `cp2-gomod`/`cp2-gocache`, local images, disposable tree, synthetic fixtures, external outputs. No new dependencies/images/tools/network. Runtime checks: `--network none`, `GOPROXY=off`, `-mod=readonly`, external outputs, 2 vCPU/1 GiB/no swap/60 s unless smaller. Prohibited: unbounded memory/device/symlink-cycle/decompression-bomb/real-credential/malicious-external-repo/unknown Makefile/release-script/destructive operations; #2019/#1572 original inputs. Stop an action at its safety boundary and document it.

## Evidence and risk classification

Every material finding: stable CP4 ID; component/source anchor; entry point/trust boundary; preconditions; expected/observed; controlled reproduction or source evidence; impact; exploitability/reachability limits; confidence; severity rationale; existing mitigations; recommended future validation; relationship to CP0–CP3. Labels: **Measured**, **Verified from source**, **Inferred**, **Unverified/NOT RUN**. Do not conflate severity and confidence; a confirmed unsafe operation is not automatically remotely exploitable. Preserve findings that disprove/narrow earlier hypotheses.

## Deliverables

Create `_assessment/current-state/CP4-security-integrations-tooling/report-CP4.md` using the shared template with sections: executive assessment; environment/provenance; threat model/trust-boundary map; scan-completeness/error propagation; sensitive-data exposure matrix; filesystem/symlink/ignore/archive analysis; configuration/rule assessment; Git/hooks/CI; Makefile/script effect matrix; build/dependency/release matrix; confirmed findings/risk register; mitigations/existing controls; unverified hypotheses/limitations; CP0–CP3 crosswalk; prioritized future verification/remediation candidates; CP5 handoff; repository-integrity statement.

Include exact sanitized commands and exit codes for executable checks. Do not copy large raw logs/profiles/secrets/corpora into Git. Update `_assessment/README.md` and roadmap to reflect CP3A completion and CP4 actual status. Do not modify `_backlog/` dispositions.

## Definition of done

Every required investigation complete or explicitly marked incomplete; every material finding evidence-backed and classified; documentation links/cross-references valid; canonical source unchanged; final status recorded; only authorized `_assessment/**` documentation committed/pushed; stop for owner review. Do not begin CP5.

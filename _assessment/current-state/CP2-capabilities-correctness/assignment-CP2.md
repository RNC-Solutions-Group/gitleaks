# CP2 — Capabilities, correctness, and regression baseline

**Status:** Ready for explicit authorization; not executed  
**Authority:** Execute only when Dave explicitly instructs `Execute CP2`  
**Code baseline:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (record the current documentation-branch HEAD separately)  
**Scope:** Verify the fork's implemented capabilities and observable behavior against the CP0 environment inventory and CP1 source-linked architecture. Source, tests, fixtures, build/configuration, CI, Git metadata, and `_backlog/` remain read-only. Write only under `_assessment/current-state/CP2-capabilities-correctness/`, with necessary `_assessment/` index/roadmap updates.

## Objective

Establish an evidence-backed **expected-versus-observed functional baseline**: which CLI and library behaviors exist, which work, which fail, which differ by mode/platform/build tag, and which are currently untestable. Separate genuine regression gaps from undocumented-but-intentional behavior and source-only defect hypotheses. Later external PR reviews will test proposed changes against this baseline, not against assumptions copied from upstream discussions.

This checkpoint does **not** repair defects, change rules, adopt outside patches, run unbounded stability/security reproductions, or optimize performance. CP3 owns benchmark/resource measurement; CP4 owns detailed security, integrations, and release-tooling risk analysis.

## Preconditions and safe test environment

1. **Read first:** root `AGENTS.md`, `_assessment/AGENTS.md`, `_assessment/README.md`, [roadmap](../assessment-roadmap.md), [CP0 report](../CP0-reconnaissance/report-CP0.md), [CP1 report](../CP1-architecture/report-CP1.md), and this assignment. Start with repository identity, canonical SHA, Go/tool availability, branch state, and `git status --short`. CP1's 15-item AR register is the verification starting point, not proof of runtime defects.
2. **Toolchain capture:** record `go version`, `go env` settings (at least GOOS/GOARCH/GOVERSION/GOMODCACHE/GOCACHE/GOTOOLCHAIN/CGO_ENABLED), Go module requirement, and the exact binary build tags. A WinGet-installed Go toolchain may differ from CI Go 1.25; do not silently merge results across versions.
3. **Dependency acquisition gate:** the Go installation alone does not provide module dependencies. Do not install tools, fetch source branches, pull images, or download modules until Dave separately authorizes the required network/dependency operation. A missing dependency remains **NOT RUN**, not a test failure or evidence that the application is defective.
4. **Disposable execution tree:** never run the existing suite in the tracked assessment checkout. CP0 proved that tests rename and may delete `testdata/repos/*/dotGit` / `.git` paths. Export the pinned code tree to an explicitly identified temporary, writable, disposable directory. Keep the tracked checkout read-only, set external Go build/module caches and report outputs, and verify before/after Git status. For symlink-sensitive tests, extract and execute on a Linux filesystem/Docker volume rather than relying on Windows NTFS tar extraction or `core.symlinks=false`.
5. **Bounded, synthetic inputs:** use only local test fixtures and small, deliberately seeded synthetic repositories with documented expected outcomes. Never scan external/live repositories, use real credentials, print raw secret-like test values, or reproduce #2019/#1572 unbounded workloads or device symlinks. Test possible destructive report behavior only against uniquely named files inside the disposable workspace; never point the scanner at owner files.
6. **Execution controls:** only direct non-mutating Go commands with `-mod=readonly` against the disposable tree, when dependencies are ready. Do not run Makefile targets, `go fmt`, `go mod tidy`, `go generate`, or release scripts. Record any build, test, vet, CLI, and targeted runtime command with version, settings, exit code, duration, fixture ID, expected/actual result, and skipped cases. For Windows permissions tests, do not assume an unreadable fixture is actually unreadable; verify the precondition (or use a non-root Linux container).

## Required investigations

1. **Capability inventory and truth labels.** Build an itemized matrix of commands, flags/defaults, sources, config and rule behavior, allowlists/ignores/baselines, decoding, archive/symlink handling, findings, reporting, diagnostics, and exit conditions. Distinguish implemented/source-verified, executed-and-passed, executed-and-failed, skipped, and not run. Identify source locations and expected behavior independently of help text.
2. **Existing automated baseline.** Inventory direct and indirect package/CLI test coverage. Once the safe execution gate is satisfied, run bounded `go test -mod=readonly ./...` and `go vet -mod=readonly ./...` in the disposable tree; capture each package's results, failures, skips, and environmental conditions. Optional `-race` belongs in a compatible, verified Linux toolchain with C compiler and separately recorded results; missing race support must not invalidate otherwise useful tests.
3. **CLI contract and mode parity.** Verify `git`, staged/pre-commit, `dir`, `stdin`, hidden `detect` variants, hidden `protect`, `version`, help/invalid flags, and significant flag interactions. Map each mode to expected findings, stderr/stdout, report selection, and exit status. Confirm where `protect → DetectGit → DetectSource` loses the command context rather than calling it a separate detector endpoint.
4. **Error/completion integrity (AR-01, AR-02).** Exercise deterministic, bounded file read failures, Git subprocess/stderr or diff failure, archive processing failure, cancellation/timeout, and `protect` timeout behavior where safely possible. Assert not just findings but whether the CLI accurately distinguishes completed, partial, failed, and timed-out scans. Record exactly where an error was surfaced, logged only, suppressed, or converted into a successful return. Do not use indefinite waits or uncontrolled subprocesses.
5. **Configuration and source-path correctness (AR-05, AR-06, AR-07, AR-08).** Test explicit/env/source-local/default config precedence, rule selection, repeated config translation (including global `extendDepth`), invalid regex handling, baseline self-scan/path normalization, Git commit allowlists, path-only rules, and `gitleaksignore` behavior using finite files. Distinguish log messages from actual skips. Leave untrusted symlink-to-device and unbounded ignore-size analysis to CP4.
6. **Detection and filtering behavior.** Build a bounded correctness matrix for rule evaluation, `gitleaks:allow`, entropy, allowlists, fingerprint/ignore/baseline suppression, generic-rule deduplication, path rules, encoded findings, and archive depth/size options. Include decode depth 0/1/default with small finite inputs and examine N+1 pass semantics, but defer resource scaling/OOM conclusions to CP3.
7. **Reporting, redaction, and output contracts (AR-04, AR-13, AR-14).** Verify JSON, CSV, JUnit, SARIF, template, stdout versus file, empty findings, exit-code overrides, ordering, and the actual `--report-template`/`--report-format` interaction. Reproduce path-only redaction using artificial filenames and `--redact`, and the report-destination validation side effect using a disposable pre-existing file; preserve before/after evidence without exposing secrets. Redaction leakage at debug/trace and symlink/device security details belong to CP4, but record any directly observed leak as a finding.
8. **Platform and build-variant parity (AR-12).** Record real Windows test skips, CRLF effects, and symlink limitations. Where a pinned Linux toolchain and dependencies exist, run representative tests on Linux/symlink-preserving storage. Test `gore2regex` only if the necessary authorized dependencies/build environment are available; otherwise mark release parity NOT RUN and hand it to CP4. Never compare binaries from unknown or dirty revisions as equivalent baseline artifacts.
9. **Regression test mapping.** For every confirmed failure or significant unverified source hypothesis, identify the smallest repeatable reproduction, expected assertion, available existing test, and proposed future regression test location. Suggest fixes only as future implementation candidates; do not write or modify application tests on this branch.
10. **Cross-check and handoff.** Reconcile CP0 and CP1 claims with observed behavior, clearly preserve contradictions, classify findings by potential impact/confidence, and assign unresolved performance/security/release questions to CP3/CP4 without advancing those checkpoints.

## Big questions

- Which advertised/implemented capabilities are demonstrably correct in this pinned build, and under which modes, flags, and platforms?
- Can a scan terminate with success while omitting input, and how would an operator recognize partial coverage?
- Do findings and exit codes remain stable across paths, Git history/diffs, stdin, baseline/ignore rules, and output formats?
- Which defect hypotheses from CP1 are reproducible, which are disproved, and which remain untested due to the environment?
- Are secret redaction, report-target handling, and configuration behavior consistent with the CLI's documented contract?
- Which test assertions and invariants would be required to judge an external PR as beneficial, irrelevant, or regressive **relative to our code baseline**?

## Deliverables

Create `_assessment/current-state/CP2-capabilities-correctness/report-CP2.md` using the shared [report template](../../templates/assessment-report.md). Include:

1. **Test-environment manifest:** pinned code SHA, documentation HEAD, OS/toolchain/build tags, module/cache provenance, fixture IDs and hashes when appropriate, command inventory, constraints, and before/after status.
2. **Capability and CLI contract matrix:** feature/flag/source, code anchor, expected behavior, observed result, status, and coverage gaps. Include legacy and primary commands.
3. **Package-test and static-check evidence:** actual command lines, exit codes, timing, skipped tests, failure summaries, and explicit NOT RUN results. Keep raw test artifacts outside the tracked repository.
4. **Correctness findings register:** stable CP2 identifiers, verified/inferred/measured labels, exact controlled reproduction, expected versus actual behavior, user/CI impact, and evidence references. Preserve meaningful negative findings (a hypothesis that failed to reproduce).
5. **AR crosswalk:** disposition every CP1 AR-01 through AR-15 as confirmed, disproved, partially verified, deferred to CP3/CP4, or blocked; never call an unrun test passing.
6. **Regression gap register:** existing vs missing assertions, proposed safe fixtures/tests, OS/build-tag differences, and smallest follow-up verification necessary.
7. **Readiness and integrity:** unresolved prerequisites, CP3/CP4 handoffs, source/test/fixture modification statement, before/after status, and recommendation for Dave's review.

Supplementary evidence documents are optional and must remain within this CP2 directory. Update `_assessment/README.md` and the roadmap only to reflect actual work completed; leave `_backlog/` unchanged.

## Hard stops and continuation rules

- Stop the affected action if the repository/branch/SHA differs, any protected source/test/fixture/config path would be modified, a command would download unapproved dependencies, a test would touch real user data, or raw secret-like values would be recorded. Do not reset/delete unexpected changes to conceal them.
- Do **not** execute tests from the tracked checkout; all runtime work is in disposable source copies with external caches. Respect Windows symlink limits and preserve a Linux-native symlink test route where available.
- Do not install Go or GCC, pull Docker images, perform `go mod tidy`, execute Makefile targets, regenerate rule configuration, change Git metadata, make external requests beyond separately authorized module acquisition, or run unsafe OOM/device fixtures.
- When tooling, optional scanners, or corpus components are unavailable, record NOT RUN and continue independent safe checks. If a required runtime claim cannot be tested, mark CP2 conditional/incomplete for that claim rather than declaring the capability correct.
- **Stop after CP2.** Do not begin CP3, change backlog dispositions, modify production code/tests, or conduct external PR integration experiments without new authorization.

## Definition of done

Meet the common [`../../README.md`](../../README.md#shared-definition-of-done) checklist. Every advertised capability in scope must have a source reference and explicit observed/unverified status; every material executed test must include its conditions and exit code; all 15 CP1 verification leads must have dispositions; all protected files must remain unchanged. Deliver the report, identify exact verification gaps, and stop for Dave's review.

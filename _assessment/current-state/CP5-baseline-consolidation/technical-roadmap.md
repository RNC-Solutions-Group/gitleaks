# CP5 — Prioritized technical roadmap

**Status:** Complete — recommendations only; no implementation  
**Canonical code SHA:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`  
**Basis:** [`findings-register.md`](findings-register.md); no backlog disposition changes.

This roadmap orders future **investigation/remediation candidates** by operational risk. It is not an implementation plan and does not authorize code changes. Priority labels: P0 highest, P3 lower.

---

## CP5-R1 — Fail-closed scan completeness and exit status

- **Problem:** Git stderr/exit failures, unreadable files, and worker/fragment errors can produce partial scans that exit 0 and look clean.
- **Findings:** CP5-F01 (High, P0); CP4-001/008/009, CP2-005/006/013, CP1-005/006/013.
- **Dependencies:** None for design; requires a decision on exit-code compatibility.
- **Suggested boundary:** `sources` error propagation, `DetectSource` callback contract, CLI summary/exit mapping. Avoid changing finding semantics.
- **Required regression tests:** malformed `git log`/stderr, unreadable file (non-root), malformed diff tail, worker error, timeout, and clean/finding controls; assert complete/partial/failed classes.
- **Compatibility risk:** Existing scripts may rely on exit 0 for warning-only scans; define partial status carefully (new nonzero or explicit status).
- **Acceptance criteria:** An operator/CI can distinguish complete-clean, complete-findings, partial, failed, and timed-out scans without parsing logs alone.
- **External PR candidate:** Possible; evaluate applicability and compatibility against this fork, not upstream intent alone.

## CP5-R2 — Non-destructive, symlink-safe report destinations

- **Problem:** Report-path validation deletes/truncates existing files before scanning; symlink targets are truncated and links replaced.
- **Findings:** CP5-F02 (High, P0); CP4-002, CP2-003, CP1-007.
- **Dependencies:** Report contract decision (atomic write, overwrite/refuse).
- **Suggested boundary:** `cmd/root.go` report validation and `findingSummaryAndExit`; no reporter format changes.
- **Required tests:** sentinel file, symlink target, failed scan/template, `-` stdout, read-only directory, TOCTOU-safe behavior.
- **Compatibility risk:** Changing overwrite semantics may break automation; provide explicit policy.
- **Acceptance criteria:** Failed/aborted scans never destroy a pre-existing target; symlink behavior is explicit and tested.
- **External PR candidate:** Possible; must preserve report-path CLI contract.

## CP5-R3 — Sensitive-output and redaction boundary

- **Problem:** Debug/trace logs leak secrets despite `--redact`; Git stderr is unredacted; path-only redaction corrupts text.
- **Findings:** CP5-F03 (Medium, P1); CP4-003/011, CP2-004/009, CP1-012.
- **Dependencies:** Logging policy/level decision.
- **Suggested boundary:** decoder/detect logging, log redaction helper, `Finding.Redact`; no report format change.
- **Required tests:** synthetic secret at info/debug/trace, decoded/`gitleaks:allow` paths, JSON/CSV/JUnit/SARIF/template/verbose, path-only findings.
- **Compatibility risk:** Reducing debug detail may hinder troubleshooting; provide a safe secret-redaction mode.
- **Acceptance criteria:** `--redact` guarantees no raw secret in documented outputs/log levels, or docs explicitly narrow the guarantee.
- **External PR candidate:** Possible; security-sensitive, needs full format/log coverage.

## CP5-R4 — Input and configuration trust boundaries

- **Problem:** Repository-local `.gitleaks.toml`/`.gitleaksignore` are auto-loaded; malformed regex panics; ignore files follow symlinks without regular-file/size/blocking safeguards.
- **Findings:** CP5-F04 (Medium, P1), CP5-F05 (Medium, P1); CP4-004/006, CP2-008/014, CP1-010/011/017.
- **Dependencies:** Policy for untrusted-repo scanning; #1572 validation.
- **Suggested boundary:** config/ignore loaders and validation; add opt-out/required-explicit-config mode; no rule changes.
- **Required tests:** local config allowlist/evasion, invalid regex, extend path, finite/symlink/size/blocking ignore controls (no device tests).
- **Compatibility risk:** Existing repos rely on local config; changing defaults is behavior-breaking.
- **Acceptance criteria:** Untrusted repo cannot silently disable detection without an explicit policy; malformed config yields controlled error; ignore reader has bounded behavior.
- **External PR candidate:** Mixed; some upstream hardening PRs may be applicable, others assume different behavior.

## CP5-R5 — Baseline and path self-scan exclusion

- **Problem:** Baseline/config files can be scanned and re-report embedded secrets when path comparisons miss.
- **Findings:** CP5-F06 (Medium, P2); CP4-005, CP2-007, CP1-009/023.
- **Dependencies:** Path canonicalization policy.
- **Suggested boundary:** `DetectContext` skip logic, baseline path resolution.
- **Required tests:** baseline inside/outside source, absolute/relative source, config in repo, symlinked source.
- **Compatibility risk:** Canonicalization may change allowlist matching on Windows.
- **Acceptance criteria:** Baseline/config file never re-reported for supported source/path forms.
- **External PR candidate:** Possible; must handle cross-platform paths.

## CP5-R6 — Resource limits, cancellation, and archive/ignore availability

- **Problem:** Archive expansion lacks uncompressed budget; decoder cost is high; blocking stdin/ignore reads delay cancellation.
- **Findings:** CP5-F08 (Medium, P2), CP5-F09/F10 (P3); CP4-007/009, CP3-001/004, CP1-016.
- **Dependencies:** Resource-policy decision.
- **Suggested boundary:** archive/file/decoder limits, context-aware readers; no detection-rule changes.
- **Required tests:** bounded archive entries/uncompressed size, decode depth 0/1/2 coverage, blocked-read cancellation, memory ceilings.
- **Compatibility risk:** New limits can skip content; must be opt-in/default-documented.
- **Acceptance criteria:** Bounded inputs terminate predictably and preserve detection for allowed content.
- **External PR candidate:** Possible for archive limits; performance changes need coverage proof.

## CP5-R7 — Build, CI, and release reproducibility

- **Problem:** Toolchain drift; tagged full-graph/vet failures; Makefile build/profile network/mutation; CI unpinned `gotestsum`, generated-config mutation, fork gitleaks workflow disabled; root image/no SBOM.
- **Findings:** CP5-F14/F15/F18 (Medium, P2); CP4-012/013/014/015/016; CP4A-01.
- **Dependencies:** CI ownership/platform decisions.
- **Suggested boundary:** workflows, Makefile, Dockerfile, release metadata; no scanner logic.
- **Required tests:** CI-equivalent Go 1.25 build/test, tagged app + full graph, vet, non-root container, SBOM/signing check, pinned tool versions.
- **Compatibility risk:** Workflow/tool changes affect contributors/releases.
- **Acceptance criteria:** Reproducible build/test/release with documented toolchains and no hidden network/mutation in diagnostics.
- **External PR candidate:** Mixed; tooling PRs must be evaluated against fork CI, not upstream assumptions.

## CP5-R8 — Coverage-preserving decoder performance work

- **Problem:** Decoder encoding scan dominates plain-text CPU (~28.7× vs depth 0); per-file buffer allocation adds pressure.
- **Findings:** CP5-F09/F10 (Medium/Low, P3); CP3-001/003/004/007.
- **Dependencies:** Decoder coverage test corpus; CP3 manifest.
- **Suggested boundary:** `detect/codec` and `sources/file` buffer handling; no report/CLI changes.
- **Required tests:** percent/Unicode/hex/Base64/URL-safe/nested fixtures, depth 0/1/2, positive/negative parity, allocations/throughput on pinned fixtures.
- **Compatibility risk:** Any gating can miss secrets; depth-0 is **not** an acceptable equivalent.
- **Acceptance criteria:** Measured improvement with identical findings across the decoder corpus; no coverage regression.
- **External PR candidate:** Possible only with full coverage evidence; do not accept a faster scan that omits work.

## CP5-R9 — Contract and determinism refinements

- **Problem:** Finding order is nondeterministic; `--report-path -`/template/empty-CSV/SARIF semantics are ambiguous; empty CSV lacks header.
- **Findings:** CP5-F11/F12/F13 (Low/Medium, P3); CP2-001/002/012/015/016, CP1-008/019/024.
- **Dependencies:** Owner/contract decisions.
- **Suggested boundary:** documentation and output ordering/format behavior.
- **Required tests:** deterministic order or documented sort, stdout format, template implication, empty CSV, SARIF redaction.
- **Compatibility risk:** Output changes can break parsers.
- **Acceptance criteria:** Documented, tested contracts; no silent behavior change.
- **External PR candidate:** Mixed; contract PRs need owner decision.

---

## Ordering summary

1. CP5-R1 scan completeness (P0).
2. CP5-R2 report filesystem (P0).
3. CP5-R3 sensitive output (P1).
4. CP5-R4 config/ignore trust (P1).
5. CP5-R5 baseline/path self-scan (P2).
6. CP5-R6 resource/cancellation (P2).
7. CP5-R7 build/CI/release (P2).
8. CP5-R8 performance with coverage (P3).
9. CP5-R9 contracts/determinism (P3).

No initiative is implemented; none is authorized by this roadmap. External PRs may be candidates but must be evaluated independently against this pinned baseline.

# V13 Reporting Architecture Options (Assessment)

**Status:** design/assessment only. No multi-format, contextual-excerpt or
repeated-finding feature is enabled in V13.

## V13-03 — multiple report formats from one scan

### Current architecture

`Detector` returns one finding slice. `cmd.reportOutput` validates one primary
destination, renders a temporary file, optionally renders one informative
companion, and publishes according to overwrite/aggregate/reject policy. A scan
is performed once; reporting is already decoupled from detection.

### Safe proposed contract

A future `--report-plan` or repeatable `--report` option should use a parsed,
ordered plan such as `json=findings.json,csv=findings.csv,sarif=findings.sarif`.
Before scanning:

1. Parse every `format=path` entry.
2. Reject duplicate formats/paths and unknown formats.
3. Validate every destination independently without mutation.
4. Reject collisions with every input, configuration, baseline, ignore file,
   template and other output — including symlink/hardlink identity.
5. Render every output to a restrictive temporary file after one scan.
6. Publish in a deterministic order.
7. If a later publication fails, report which outputs were published; do not
   claim rollback or cross-file atomicity.

The existing single `--report-path` interface remains the compatible default.
The phase-2 implementation would factor a `reportPlan` from `reportOutput`
rather than adding arrays to today's `reportOutput` struct. It would need new
tests for duplicate paths, cross-output hardlinks, informative companion pairing,
aggregate numbering across a plan, partial failure and permissions.

**Disposition:** `DEFERRED — OUTSIDE AUTHORIZED SCOPE`. The required validation
and failure semantics are larger than a safe V13 patch, and no output plan is
enabled without them.

## V13-04 — bounded finding context

### Current contract

`Finding.Line` contains the matched source line but is excluded from JSON
(`json:"-"`), CSV and SARIF. `Finding.Match` contains the rule match. V12 did not
add surrounding context to machine schemas.

### Proposed opt-in contract

If authorized, introduce a separate, explicitly selected report mode rather than
silently extending existing JSON/SARIF/CSV schemas:

- `--context-lines N` (default 0, maximum bounded by a documented byte budget);
- a `Context` object with `before`, `line`, `after` and a truncation marker;
- no raw context in the default report;
- context inclusion only after all known sensitive spans are redacted, including
  other findings in the window;
- no Git historical blob fetch: context comes only from a currently available
  fragment;
- byte/line caps and CRLF/Unicode-aware slicing;
- context omitted from shared report formats unless the operator explicitly
  selects the context mode.

The design must decide whether context is a console-only view, a sidecar report,
or a new opt-in schema field. It must not change allowlist/detection semantics.

**Disposition:** `DEFERRED — OUTSIDE AUTHORIZED SCOPE`. Implementing this safely
requires a new report schema and a redaction window over all findings; an
unbounded `Finding.Line` export would be unsafe.

## V13-06 — repeated-finding reporting

### Finding identity

A raw SHA-256 of a secret is an unkeyed, enumerable identifier. It is not safe as
a default public grouping token, particularly for low-entropy secrets. Primary
findings, fingerprints, commit attribution and ignore semantics must remain
untouched.

### Safe proposed policy

A future grouping report would be a separate, opt-in sidecar keyed by an
operator-supplied secret (HMAC-SHA-256, displayed only as a scoped grouping
identifier). The key would never be persisted in the report. Without a key, the
UI may group by rule ID plus location/commit without exposing a secret-derived
token, or show counts only.

**Disposition:** `POLICY DECISION REQUIRED`. No grouping or deduplication is
enabled in V13. Unkeyed secret hashes are explicitly rejected.

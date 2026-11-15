# CP-V15 — Complete disposition matrix

The complete matrix is machine-readable in
[`disposition-matrix.csv`](disposition-matrix.csv), with the same records in
[`disposition-matrix.json`](disposition-matrix.json) and a readable
[`disposition-matrix.md`](disposition-matrix.md). It contains **220 active packages**,
one record per package, after the V15 retirements.

Each record carries:

- stable package ID, category and upstream issue/PR references;
- current technical status and accepted source status;
- evidence confidence, separate from severity and implementation confidence;
- recommended action, severity, composite priority and complexity;
- planned cycle, acceptance requirement and all remaining acceptance conditions.

The matrix is a **planning disposition**, not a canonical claim that an upstream item is
fixed. Canonical removal is limited to the six V15 retirements in the reconciliation
record.

## 1. Planning counts

| Action | Packages |
|---|---:|
| `FIX` — confirmed defect, implement | 19 |
| `ADAPT` — use upstream work with a corrected design | 15 |
| `COMBINE` — merge into a shared correction | 16 |
| `IMPLEMENT` — direct RNC implementation | 5 |
| `POLICY+IMPLEMENT` — decision plus code | 1 |
| `IMPROVE` — bounded residual improvement | 1 |
| `REPRO` — targeted reproduction before disposition | 52 |
| `DEFER` — valid but not now | 43 |
| `REJECT` — unsafe, redundant or unjustified | 35 |
| `RELEASE` — release/artifact gate | 10 |
| `ACCEPT` — current behavior/partial feature accepted | 5 |
| `OPERATOR` — policy decision required | 7 |
| `ARCH` — architecture/design gate | 3 |
| `TEST` — test-only evidence | 3 |
| `BENCH` — measurement first | 1 |
| `POLICY` — detection tradeoff decision | 1 |
| `ESCALATE` — external incident owner | 3 |

Severity: **25 S1, 113 S2, 65 S3, 17 S4**. Priority: **3 P0, 32 P1, 111 P2,
71 P3, 3 external P0\***. Complexity: **80 S, 118 M, 22 L**. Evidence confidence:
**40 High, 83 Medium, 97 Low**. These dimensions are intentionally separate: a
high-severity issue can have low evidence; a low-complexity fix can be high value.

`P0*` is used only for the three upstream Slack incident reports. Those are not RNC
scanner defects and do not consume RNC implementation priority.

## 2. Immediate confirmed corrections

These are the only packages with both (a) a reproduced/current-source defect and
(b) a sufficiently bounded correction to start without a broader design:

| Package | Issue / PR | Confirmed behavior | Correction |
|---|---|---|---|
| `SCAN-01` | #2010 | A ZIP named `.other` is skipped as `application/zip`; the member secret is missed | Content-first archive identification using the existing bounded prefix, with text-mislabel and corrupt-archive negatives |
| `SCAN-02` | #2000 / #2172, #2208 | Base64 payload decoding to printable non-ASCII is rejected; synthetic GitLab token missed | Accept printable UTF-8 under a control-character envelope; prefer #2172's narrower filter |
| `ACC-03` | #2122 / #2150, #2255 | `withinColumns` compares columns across lines; primary at line 1 joins context at line 6 | Require the same line before column proximity; add cross-line negative |
| `ACC-09` | #2327 / #2335 | Generic rule does not match a high-entropy value beginning `+` or `/` | Extend the value class without accepting arbitrary delimiters |
| `ACC-10` | #2257 / #2260 | Airtable PAT is not detected without the literal `airtable`; the regex itself matches | Change the keyword to the self-signaling prefix and add a bare-token fixture |
| `CONF-08` | #2024 | Extending a rule drops `RequiredRules`/`SkipReport` | Merge those fields with explicit precedence; add extend tests |
| `REP-02` | #1424 / #2215 | Later-line CRLF columns are one byte high because byte columns use `prevNewLine` | Base byte columns on `lineStart`; retain SARIF codepoint behavior |

`GIT-02` (#1028 / #2301) is equally important but is a **policy + implementation**
package: first-parent coverage recovers merge-only secrets but adds duplicates, new
fingerprints and baseline/ignore migration. It belongs at the front of C2 after the
operator's decision, not in C1's unconditional fix set.

## 3. Packages deliberately not implemented

The matrix rejects 35 packages. The main groups are:

- **Unsafe network/remote surfaces:** #1992 URL/JSON input, #2056 MCP, #2154 SDK,
  #1013 live credential validation, #1557/#2186 remote config, #1874 CEL.
- **Unsafe identity/IO:** #1162 unkeyed secret hash, #2101 secret-value ignore,
  #1182 pre-scan text replacement, #1581 broad alpha/numeric allowlist.
- **Unsafe distribution:** #1661 npm wrapper, #1485 installer as written, #2084
  Actions example with bad semantics.
- **Overbroad/unsupported provider rules:** #1335 Chargebee, #1086 Plivo ID,
  #774 Segment public token, #2077 Spring placeholder, #1089 New Relic browser token,
  #2001 mixed 18-file provider bundle.
- **Duplicate/non-secret requests:** #1234 selected files (covered by `gitleaks files`),
  #739 cookie SQLite, #1525 password-manager exports, #1082 wallet seed phrases,
  #2094 Claude Code session URL, #1829 Python JWT generic.

`DEFER` is not a rejection. It records a real request that is not justified before the
core/release gates, often because it needs a trust model, a schema, or measured demand.

## 4. Plan-level caveats

- 52 packages require reproduction before any implementation decision. This is the
  honest state of a backlog assembled from upstream titles and partial source review.
- 35 `REJECT` rows are planning recommendations. The operator may choose to keep an
  upstream item open in the canonical ledger; V15 does not close it.
- The matrix does not assign 220 code changes to 220 developers. It assigns every
  **review obligation** to a decision path and, where justified, a cycle.
- A finding-count change is never acceptance for C6. A green Go suite is never an
  accuracy oracle.

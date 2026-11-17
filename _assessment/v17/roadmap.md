# CP-V17 roadmap and decision update

## 1. Sequence

The V15 sequence remains valid in outline:

```text
V18 C3  Source readers and path identity
V19 C4  Configuration and allowlist model
V20 C5  Reporting and fingerprint compatibility
        Formal structural architecture checkpoint
V21 C6  Generic detection accuracy
V22 C7  Existing provider-rule repair
        AG2 Detection model, generator invariants, public API
V23 C8  Provider expansion and rule governance
V24 C9  Release identity and artifact integrity
V25 C10 Elective feature triage
    VER Independent harness and release-candidate verification
```

Changes justified by V16/V17 evidence:

1. **C3 is split into C3A (boundary/encoding/notebook) and C3B (path identity).** V17
   changed the reader buffer and archive probe; C3A can stabilize those contracts before
   any fingerprint migration. C3B needs OD-8 and should not be bundled into one batch.
2. **C3 does not absorb the accuracy backlog.** The V16 gap is concentrated in
   identifier/context/value-shape misses, not path identity. C3 should make readers
   deterministic; C6/C7 should change detection semantics.
3. **C5 is promoted in importance, not moved earlier.** C2's opt-in merge mode and the
   context mode create explicit new occurrence/fingerprint cases. C5 must settle
   duplicate, ignore and fingerprint policy before merge/context defaults or path
   identity change.
4. **C6's first wave should follow the V16 gap order:** WordPress authentication
   identifiers (93), database-password context (76), application-password context (68),
   then positional/encoded and URI/connection-string formats. Do not start with broad
   low-confidence provider expansion.
5. **C7 remains after C6.** The V16 gap's `RULE_ABSENT` and provider-specific misses need
   the generic identifier/context model first; otherwise each fix becomes a one-off
   regex.
6. **C9 is not blocked by C6/C7.** Release identity and artifact integrity remain a
   parallel gate after AG1, as V15 proposed.
7. **AG2 remains after C6/C7.** V17 changed internal detector signatures but not the
   exported model; the public-API freeze still needs the accuracy/model work.

## 2. V16 gap evidence impact

The independent V16 report recorded:

- 1,125 positive occurrences; 533 missed; 447 risk-bearing.
- 165 missed occurrences recovered by another scanner.
- 26 development opportunities.
- Main risk concentrations: WordPress authentication material (93), database password
  identifiers (76), application password identifiers (68), positional/encoded
  passwords (61), URI/connection strings (27).
- Root causes: identifier mismatch 167, format unsupported 112, context insufficient
  72, rule absent 34, source unsupported 26, unresolved 31.

This does not change C3's identity scope. It does change C6's priority order and makes
scanner-side candidate instrumentation a prerequisite for value-validator and
context-insufficient cases.

## 3. Decision register update

| ID | Decision | Status / recommendation |
| --- | --- | --- |
| OD-1 | Merge coverage default | **Open.** Keep default off in V17; ship `--scan-merges`; require C5 migration before changing the default. |
| OD-2 | AG1 refactor | **Accepted** in V16. V17 builds on it. |
| OD-3 | Airtable keyword | **Accepted** in V16. |
| OD-4 | Independent harness ownership | **Open.** V17 did not run or modify the harness. C6 remains gated. |
| OD-5 | First program wave | **Accepted.** C1/AG1 complete; C2/V17 under review. |
| OD-6 | `minVersion` policy | **Open**, C4/C9. |
| OD-7 | Module identity | **Open**, C9/AG2. |
| OD-8 | Fingerprint v2 / path identity | **Open**, C3B/C5. |
| OD-9 | Source-local config trust | **Open**, C4. |
| OD-10 | Release attestation | **Open**, C9. |
| OD-11 | Provider rule policy | **Open**, C7/C8. |
| OD-12 | Elective features | **Open**, C10. |
| OD-13 (new) | Git context default | **Recommend opt-in only.** `--git-context-lines` is available; default remains 0 until C5/C3 prove the performance and duplicate envelope. |
| OD-14 (new) | Git subdirectory scope | **Recommend explicit pathspec.** Do not silently narrow repository-wide history; document the warning and migration. |

## 4. Performance follow-up

- **C3A:** add path-allowlist match caching only if it can share a decision with the
  source without weakening AND-allowlists. The V17 profile shows 25–34% of small-file
  CPU in duplicate global-path regex evaluation.
- **C5:** decide whether merge/context occurrence identity needs a separate occurrence
  key rather than content dedup.
- **C6:** use scanner-side candidate instrumentation before changing generic entropy or
  identifier behavior; the V16 gap explicitly warns against a global entropy reduction.
- **C9:** keep the V16 Docker `go mod download` proxy failure classified as
  infrastructure; the rerun passed with no product change.

## 5. Not proposed

- A wholesale detect rewrite from #2023.
- Making merge or context scanning default by implication.
- A broad provider-rule expansion solely to improve a benchmark score.
- A second report writer or unkeyed secret hash.
- Backfilling missing shallow history.

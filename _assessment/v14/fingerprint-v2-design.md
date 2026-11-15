# Fingerprint v2 Design (OP-06)

**Decision state:** DESIGN ONLY. Existing v1 `.gitleaksignore` behavior, baselines and filters are unchanged in V14.

## Problem

V1 fingerprints identify a location, not content:

```text
<commit>:<file>:<rule-id>:<line>
<file>:<rule-id>:<line>          # commitless scans
```

Suppressing a finding at a location can therefore suppress a different secret later placed at the same rule/path/line. The design must prevent silent location-only suppression without breaking existing ignores.

## Proposed opt-in v2 grammar

```text
v2:<key-id>:<hmac-sha256-hex>
```

- `v2` is a literal version marker.
- `key-id` is a non-secret operator identifier for the key epoch (for example a short fingerprint of the public key context, never the key itself).
- `hmac-sha256-hex` is HMAC-SHA-256 over a canonical, domain-separated content identity.
- The fingerprint does not expose the secret or an unkeyed enumerable hash.
- The exact canonical input includes a versioned domain separator, rule ID, normalized path identity, normalized secret bytes and a scope discriminator for commit versus commitless scans. The design must specify whether commit is included; it should be optional/encoded in the domain so a v2 ignore can be location+content or commit+location+content.

## Why keyed

A bare SHA-256 over a secret is enumerable for low-entropy credentials and should not be a default public identifier. A keyed HMAC:
- resists offline dictionary confirmation without the key;
- allows cross-run stable IDs only when the same key is used;
- makes key rotation an explicit migration, not an accidental identity change.

## Key handling

- Key source: protected operator file/environment outside the repository and outside reports; never a command-line value visible in process listings/logs.
- Minimum 32 bytes of key material; reject missing/weak/invalid keys.
- Never print, log, profile, report, or commit the key.
- Key identity (`key-id`) may appear; the key may not.
- Rotation: v2 IDs change when the key changes. Existing v2 entries for an old key remain accepted only if a previous key is explicitly supplied for migration; otherwise they are reported as unmatched, never silently reinterpreted.
- Cross-run grouping: deterministic only within the same key/domain. Cross-key grouping is intentionally not stable.

## Compatibility and migration

1. **Default remains v1.** No current behavior changes without opt-in.
2. **Dual acceptance:** in compatibility mode, v1 and v2 entries are both accepted. This explicitly preserves the old location-based suppression behavior; operators must understand that a v1 entry can still suppress a different secret at the same location.
3. **Strict v2-only mode:** optional later policy; rejects ambiguous v1 entries but requires a migration report and operator approval.
4. **Baseline comparison:** v1 baseline content/identity remains unchanged by default. A future v2 baseline must carry its version/key metadata and dual-read rules.
5. **Reports:** primary fingerprints remain v1 unless a v2 reporting mode is explicitly selected. Adding v2 must not silently change JSON/CSV/SARIF fingerprints.
6. **Windows/Linux paths:** normalize separators/root identity consistently; do not use absolute host paths as canonical content identity. CRLF secret normalization must match baseline semantics but not erase meaningful differences.
7. **Repeated secrets across commits:** content-only v2 groups/suppresses all occurrences; commit-scoped v2 suppresses only the specific occurrence. The operator chooses the scope.

## Threat model

- Untrusted repository cannot choose the key or derive HMAC values without it.
- An attacker with the key can confirm candidate secrets; protect the key accordingly.
- A leaked v2 fingerprint is not reversible but is stable within the key epoch.
- v1 entries remain location-based and are not equivalent to content assurance.
- Baseline and ignore files may contain sensitive paths/metadata; store them as confidential operational data.

## Upstream references

The design compares the problem space described by upstream issue #2239 and PRs #2286/#1566. Neither patch is assumed suitable or adopted; any future implementation must be independently written and evaluated against this design.

## Recommendation

Do not implement or activate v2 in V14. Recommended next step: approve a versioned v2 grammar, key-source mechanism, dual-acceptance policy and migration report before any code change. Explicitly document that accepting v1 keeps the old suppression semantics.

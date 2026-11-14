# CP-V14 Operator Policy Matrix

**Authoritative policy:** `_backlog/operator-decisions/2026-11-14-v13-policy.md` on `feat/issues-and-pr`
**Starting product baseline:** V13 candidate `audit/v13-integration@71be6f3` (V13 not merged when V14 work began)
**V14 branch:** `feat/v14-policy-implementation`
**Scope:** six operator decisions; no default scan-history, fingerprint, module-path or release change.

| Decision | Current behavior | V14 change/experiment | Authorization | Compatibility/security risk | Acceptance evidence | Deliverable |
|---|---|---|---|---|---|---|
| OP-01 keyed HMAC grouping | No grouping; primary findings only | New opt-in `--grouped-report` plus `--grouping-key-file`; separate JSON sidecar; HMAC-SHA-256 with domain separator; no suppression | **APPROVED FOR IMPLEMENTATION** | New optional output can expose paths/fingerprints; key file trust; destination collision; no primary change | Unit/CLI tests for repeated values, low entropy, redaction, key strength, destination collision, missing key; standard/RE2/race | `report/grouped.go`, `cmd/root.go`, tests; docs |
| OP-05 strict config | Unknown fields ignored; malformed/empty blocks fail | Opt-in `--strict-config` recursive schema validation for root/extend/rules/allowlists/required; default unchanged | **APPROVED FOR IMPLEMENTATION** | Stricter mode can reject configs that previously loaded; not default | Unit/CLI/extension tests; default non-strict compatibility | `config/strict.go`, `cmd/root.go`, tests |
| OP-04 trusted CI config | Local config trusted by default; `--no-local-config` available | Document explicit trusted config/ignore/baseline contract and synthetic hostile-fixture test | **APPROVED** | Guidance only; no default change | `TestTrustedCIConfigurationContract`; precedence probes | `CI-SECURITY.md`, test |
| OP-02 module identity | Upstream module path `github.com/zricethezav/gitleaks/v8` | COA comparison; no `go.mod` change | **EVALUATION ONLY** | Module path changes break imports/install; COA1 lowest risk | Source/module/manifest analysis; controlled resolution reasoning | `_assessment/v14/module-identity-comparison.md` |
| OP-03 first-parent default | Merge patches omitted by default; explicit first-parent opt-in | Synthetic evil/normal merge comparison; no default change | **EVALUATION ONLY** | Default change duplicates findings/fingerprints | Synthetic repo results | `_assessment/v14/first-parent-default-evaluation.md` |
| OP-06 fingerprint v2 | v1 location-based fingerprints | Versioned keyed-content fingerprint design, migration/threat model; no rollout | **DESIGN ONLY** | v1 acceptance preserves old suppression; key lifecycle | Design review | `_assessment/v14/fingerprint-v2-design.md` |

## Implementation boundary

OP-01/OP-05 are implemented opt-in features. OP-04 is documentation plus tests. OP-02/OP-03/OP-06 do **not** change product defaults, module identity or `.gitleaksignore` behavior. No release, tag, image, package or attestation is published.

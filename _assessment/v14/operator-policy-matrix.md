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
| OP-02 module identity | Upstream module path `github.com/zricethezav/gitleaks/v8`; `minVersion` comparison warns rather than rejects | COA comparison plus isolated `/v2` experiment; no production `go.mod` change | **EVALUATION ONLY** | COA2 breaks downstream imports and would require a separate `minVersion` compatibility floor; public install unverified | Local `/v2` build/test/RE2/vet/install/linker/cross-build | `_assessment/v14/module-identity-comparison.md` |
| OP-03 first-parent default | Merge patches omitted by default; explicit first-parent opt-in | Expanded synthetic ordinary/evil/conflict/octopus/repeated/pathspec/ignore/large-history matrix; no default change | **EVALUATION ONLY** | Default change duplicates findings/fingerprints | Measured scenario table | `_assessment/v14/first-parent-default-evaluation.md` |
| OP-06 fingerprint v2 | v1 location-based fingerprints | Versioned keyed-content fingerprint design, migration/threat model; no rollout | **DESIGN ONLY** | v1 acceptance preserves old suppression; key lifecycle | Design review | `_assessment/v14/fingerprint-v2-design.md` |

## Implementation boundary

OP-01/OP-05 are implemented opt-in features, hardened by V14A. OP-04 is documentation plus tests. V14A-01 through V14A-09 are recorded in `v14a-corrective-review.md`; V14A-07 removes overbroad grouped-name exclusion. **V14B supersedes V14A-08:** the 1,000-candidate cap is removed, preflight is structural, and allocation remains unbounded in numbering. OP-02/OP-03/OP-06 do **not** change product defaults, module identity or `.gitleaksignore` behavior. No release, tag, image, package or attestation is published.

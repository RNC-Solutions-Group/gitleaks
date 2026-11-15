# CP-V14 Implementation and Evaluation Report

**Policy authority:** `_backlog/operator-decisions/2026-11-14-v13-policy.md`
**Starting product baseline:** `audit/v13-integration@71be6f3fdf99ac711e7d1dbb5ee1683147ff2883` (V13 candidate, not merged)
**V14 branch:** `feat/v14-policy-implementation`
**Final PR:** created from this branch against the V13 integration branch (stacked review)
**Product baseline:** V13 + V14 opt-in features; no default scan/fingerprint/module change

## Six policy outcomes

| Decision | Outcome | Evidence |
|---|---|---|
| OP-01 keyed HMAC grouping | **APPROVED + IMPLEMENTED (opt-in), V14A-hardened** | `--grouped-report` + `--grouping-key-file`; separate sidecar; HMAC-SHA-256 domain-separated; no primary suppression; key strength/permissions/location; redaction-safe; destination-set collision and output-policy checks; tests. |
| OP-05 strict config | **APPROVED + IMPLEMENTED (opt-in)** | `--strict-config` recursive schema validation for root/extend/rules/allowlists/required; errors name fields not values; default permissive loading unchanged; extension-chain tests. |
| OP-04 trusted CI config | **APPROVED + DOCUMENTED/TESTED** | `CI-SECURITY.md`; hostile-repository precedence test; `--no-local-config` + explicit trusted config/env; no workflow activation. |
| OP-02 module identity | **EVALUATED — COA2 conditional recommendation after V14A `/v2` proof** | `module-identity-comparison.md`; isolated `/v2` build/test/install experiment passes locally; current upstream module path unchanged; public proxy/install unverified. |
| OP-03 first-parent default | **EVALUATED — retain current opt-in** | `first-parent-default-evaluation.md`; expanded ordinary/evil/conflict/octopus/repeated/pathspec/ignore/large-history matrix; no default change. |
| OP-06 fingerprint v2 | **DESIGN ONLY** | `fingerprint-v2-design.md`; keyed versioned grammar, dual-acceptance/v1 caveat, key rotation and migration; no code/default change. |

## Implementation details

### V14A correction summary

V14A-01 through V14A-05 are documented in `v14a-corrective-review.md`. Key/destination protection, grouped output policy, cross-report collision checks, inline strict validation, and key lifecycle guidance/tests were corrected on the existing branch. No new PR or retarget was performed.

### OP-01

- `report/grouped.go`: HMAC-SHA-256, domain `rnc-mod/report-grouping/v1`, normalized secret CRLF→LF, group by exact secret (case/whitespace preserved), sorted deterministic occurrence/group output.
- Sidecar excludes Secret/Match/Line; includes rule/path/line/commit/v1 fingerprint.
- `report.Finding.SetGroupingSecret` stores a pre-redaction unexported value so `--redact` does not collapse all groups into `REDACTED`.
- `cmd/root.go`: both flags required; sidecar path cannot be stdout, primary path, informative companion, or unsafe destination; key file must be regular, ≥32 bytes, and not group/world readable on POSIX.
- Writes through the existing validated temporary-file/report publication path after a completed scan. Primary report and exit semantics unchanged.
- Key file is never logged; key rotation changes group IDs (documented).

### OP-05

- `config/strict.go`: parses TOML into a map, validates known top-level keys, `extend`, rules, legacy/current allowlists and required-rule tables; malformed TOML and unknown fields are distinct errors.
- `--strict-config` validates every loaded config file, including extension-chain files, or inline `GITLEAKS_CONFIG_TOML`.
- Default non-strict parsing and all existing aliases/deprecated fields remain.

### OP-04

Trusted CI contract: explicit config/environment, `--no-local-config`, trusted ignore/baseline outside the checkout, protected report destination, documented precedence and v1 caveats. Test uses a hostile local config that suppresses a synthetic secret and demonstrates recovery with the trusted contract.

### OP-02

The isolated `/v2` experiment mechanically migrated 207 files and passed local build, standard/RE2 tests, vet, run, install, linker-stamp and cross-build checks. Public module proxy/`go install` was not verified because nothing was published. COA2 is conditionally recommended for the first independent RNC release, contingent on operator approval and a `minVersion` configuration-compatibility floor; COA1 remains the conservative fallback. `go.mod` is unchanged.

### OP-03

V14A expanded the matrix. Ordinary merges: default 1 vs first-parent 2 (one duplicate). Evil merge: 1 vs 3 (one duplicate + one merge-only). Conflict-resolution and octopus merge-only secrets: 0 vs 1. Repeated values: 2 vs 3 (one duplicate). Pathspec: 1 vs 2. A feature-fingerprint ignore leaves the merge duplicate at 1. Extra findings are not accuracy. The default remains unchanged because fingerprints, ignores, baselines and Git-version compatibility would need migration.

### OP-06

v2 design uses keyed content identity with a versioned grammar, key ID, rotation and dual-acceptance. Accepting v1 explicitly preserves old location suppression. No v2 implementation or default switch.

## Validation

| Check | Result |
|---|---|
| `go build ./...` | PASS |
| `go test -count=1 ./...` | PASS |
| `go test -count=1 -tags gore2regex ./...` | PASS |
| `go vet ./...` | PASS |
| `go mod verify` | PASS |
| `go generate ./...` and clean tracked tree | PASS (generated config line endings restored) |
| Non-root Linux `go test -race -count=1 ./...` at corrective head `75e9c0a` | PASS |
| Grouped report unit/CLI tests | PASS: repetition, pre-redaction grouping, key strength/location/identity, output policy, primary/informative/aggregate/parent-symlink/hardlink collisions, no secret output |
| Strict config tests | PASS: valid, unknown root/rule/allowlist, malformed, extension chain, inline root+extension, deprecated alias, default non-strict compatibility |
| Trusted CI hostile-fixture precedence | PASS |

No detection-rule, benchmark-harness, release, tag, image or production-module change was made.

## Residuals/operator decisions

- HMAC key lifecycle/key-management service is operator-owned; v1 fingerprints remain location-based.
- v2 fingerprint design requires a future operator decision and dual-acceptance/migration implementation.
- Module identity: COA2 conditionally recommended after the `/v2` experiment; COA1 remains fallback pending operator decision and `minVersion` policy.
- First-parent default remains unchanged pending duplicate/fingerprint/baseline migration approval.
- No release/tag/image/package/attestation was published.

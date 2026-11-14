# CP-V14 Implementation and Evaluation Report

**Policy authority:** `_backlog/operator-decisions/2026-11-14-v13-policy.md`
**Starting product baseline:** `audit/v13-integration@71be6f3fdf99ac711e7d1dbb5ee1683147ff2883` (V13 candidate, not merged)
**V14 branch:** `feat/v14-policy-implementation`
**Final PR:** created from this branch against the V13 integration branch (stacked review)
**Product baseline:** V13 + V14 opt-in features; no default scan/fingerprint/module change

## Six policy outcomes

| Decision | Outcome | Evidence |
|---|---|---|
| OP-01 keyed HMAC grouping | **APPROVED + IMPLEMENTED (opt-in)** | `--grouped-report` + `--grouping-key-file`; separate sidecar; HMAC-SHA-256 domain-separated; no primary suppression; key strength/permissions; redaction-safe; destination checks; tests. |
| OP-05 strict config | **APPROVED + IMPLEMENTED (opt-in)** | `--strict-config` recursive schema validation for root/extend/rules/allowlists/required; errors name fields not values; default permissive loading unchanged; extension-chain tests. |
| OP-04 trusted CI config | **APPROVED + DOCUMENTED/TESTED** | `CI-SECURITY.md`; hostile-repository precedence test; `--no-local-config` + explicit trusted config/env; no workflow activation. |
| OP-02 module identity | **EVALUATED — COA1 recommended** | `module-identity-comparison.md`; current upstream module path unchanged; binaries/containers are the supported RNC distribution. |
| OP-03 first-parent default | **EVALUATED — retain current default** | `first-parent-default-evaluation.md`; synthetic default 1 finding, first-parent 3 (1 merge-only, 1 duplicate); no default change. |
| OP-06 fingerprint v2 | **DESIGN ONLY** | `fingerprint-v2-design.md`; keyed versioned grammar, dual-acceptance/v1 caveat, key rotation and migration; no code/default change. |

## Implementation details

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

COA1 retains `github.com/zricethezav/gitleaks/v8`, avoiding import/linker churn. COA2 is a breaking module migration; COA3 is consumer-side replace behavior, not a transparent dual identity. No `go.mod` change.

### OP-03

Synthetic evil-merge results:

- default: 1 finding, feature commit only; merge-only secret absent;
- first-parent: 3 findings (feature duplicate at merge, feature commit, merge-only secret);
- extra count is not accuracy; 1 duplicate + 1 genuine merge-only finding;
- default change would create fingerprints/ignore/baseline duplication and requires migration.

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
| `go generate ./...` and clean tracked tree | PASS |
| Non-root Linux `go test -race -count=1 ./...` at `b41ebe9` | PASS |
| Grouped report unit/CLI tests | PASS: repetition, pre-redaction grouping, key strength, missing flags, primary collision, no secret output |
| Strict config tests | PASS: valid, unknown root, unknown rule, extension chain, default non-strict compatibility |
| Trusted CI hostile-fixture precedence | PASS |

No detection-rule, benchmark-harness, release, tag, image or production-module change was made.

## Residuals/operator decisions

- HMAC key lifecycle/key-management service is operator-owned; v1 fingerprints remain location-based.
- v2 fingerprint design requires a future operator decision and dual-acceptance/migration implementation.
- Module identity remains COA1 pending separate operator decision.
- First-parent default remains unchanged pending duplicate/fingerprint/baseline migration approval.
- No release/tag/image/package/attestation was published.

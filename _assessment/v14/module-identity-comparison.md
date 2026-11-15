# Go Module Identity Comparison (OP-02, revised V14A)

**Decision state:** EVALUATION ONLY. No production `go.mod`, import path, linker path, repository name or tag was changed. The `/v2` experiment lives outside PR #7.

## Starting facts

- Repository: `RNC-Solutions-Group/gitleaks`.
- Current declared module: `github.com/zricethezav/gitleaks/v8` (`go.mod`).
- Current product baseline: V13/V14 candidate, no stable independent release.
- Current `minVersion = "v8.25.0"` comparisons use the scanner's stamped product version. `config.validateMinVersion` logs a warning when the scanner is older; it does not fail the load or scan solely on that comparison.
- Builds stamp `github.com/zricethezav/gitleaks/v8/version.Version`.
- `go install github.com/zricethezav/gitleaks/v8@...` resolves upstream, not this fork.
- `go install github.com/RNC-Solutions-Group/gitleaks@...` is not a supported module identity today.

## Controlled `/v2` experiment

An isolated disposable worktree was created at `f4327bc` and mechanically migrated to:

```go
module github.com/RNC-Solutions-Group/gitleaks/v2
```

The experiment replaced the upstream import/linker prefix in source, `go.mod`, `.goreleaser.yml`, `Dockerfile`, `Makefile` and related text. It was not committed to PR #7.

### Actual local results

| Check | Result |
|---|---|
| Import/module/linker replacements | 207 files mechanically changed |
| `go list -m` | `github.com/RNC-Solutions-Group/gitleaks/v2` |
| `go build ./...` | PASS |
| `go test -count=1 ./...` | PASS |
| `go test -count=1 -tags gore2regex ./...` | PASS |
| `go vet ./...` | PASS |
| `go run . version` | PASS (`version is set by build process`) |
| `go install .` to temporary `GOBIN` | PASS; binary runs |
| Linker `-X .../v2/version.Version=2.0.0-rnc.test` | PASS; stamped version reported |
| Cross-build linux/amd64, windows/amd64, darwin/arm64 | PASS |
| `go mod verify` | PASS |
| `go generate ./...` | PASS (experiment worktree; no production generation) |
| Public module-proxy `go install` | **NOT VERIFIED** — no module/tag was published |

The experiment proves local source/build/test/install mechanics. It does **not** prove public `go install` resolution, proxy caching, checksum-database acceptance, release immutability or artifact signing.

## Separating three version questions

1. **RNC product version:** proposed first independent release `v2.0.0`.
2. **Go module major:** `/v2`, required by semantic import versioning for a v2 product line.
3. **Upstream configuration compatibility:** existing configs use `minVersion = "v8.25.0"`. An RNC 2.0.0 binary compared directly to 8.25.0 would be considered older and would emit a compatibility warning for otherwise compatible upstream configs. Current behavior does **not** reject the configuration solely because of that version comparison; invalid SemVer or an unparsable scanner version produces an error, and those are separate cases.

These must not be conflated. A `/v2` module does not mean upstream configuration compatibility moved to v2.

## `minVersion` proposal (not implemented)

Recommended future policy:

- Treat an upstream-style `minVersion` with major **>= 8** as a **configuration schema compatibility floor**, compared against an embedded RNC `configCompatibilityVersion` (for example `8.26.0`), not against the RNC product version.
- Add an explicit, tested mapping for RNC releases: `v2.0.0` product ↔ upstream configuration compatibility `8.26.0`.
- Permit RNC-specific `minVersion` major 2 only if a future config schema is genuinely RNC-owned.
- Document both namespaces and add compatibility tests before any module migration.
- Do not silently reinterpret `minVersion`; a release note and migration table are required.

This preserves existing upstream configs without forcing RNC product releases to reuse upstream `v8.x` numbers.

## Revised COA comparison

| Criterion | COA1: retain upstream `/v8` | COA2: RNC-owned `/v2`, product `v2.0.0` | COA3: transitional compatibility |
|---|---|---|---|
| RNC module identity | Misleading: source repo ≠ module identity | Clear RNC identity | No true dual identity; consumer-local replace only |
| Existing internal imports | No change | 207-file mechanical change; all local tests pass | Internal imports remain upstream |
| Downstream Go imports | Receive upstream path unless replaced | Breaking; consumers must migrate to RNC path | Replace directives are per-consumer and fragile |
| `go install`/`go get` | Upstream path installs upstream | RNC path becomes canonical after publication | Not transparent |
| Module proxy | Upstream namespace | New RNC `/v2` namespace after publication | Proxy cannot alias one module to two paths |
| Versioning | RNC product version and `/v8` mismatch | Product `v2.0.0` and module `/v2` align | Mixed versions persist |
| Config `minVersion` | Upstream comparison already works | Requires compatibility-floor policy above | Same issue unless policy added |
| Release artifacts | Binaries/containers independent of module path | Binaries/containers independent; Go module identity improves | Ambiguous |
| Migration cost | Lowest | Mechanical source change + release/module publication + downstream migration | High consumer coordination, low transparency |
| Risk | Independent release remains confusing for Go users | Breaking Go API consumers; proxy/public install unproven | Easy to misrepresent as compatibility |

COA3 is not recommended as a primary identity: Go modules have one canonical path; `replace` is a consumer build workaround, not a second published identity.

## Recommendation

**COA2 is the recommended target for the first independent RNC release, conditional on operator approval and a completed migration plan.** The `/v2` experiment demonstrates that the source/build/test/install mechanics are feasible. COA1 remains the conservative fallback if RNC will not publish a Go module and values upstream import compatibility above independent module identity.

No migration is performed now. Before COA2:
- implement and test the `minVersion` compatibility floor;
- publish a versioned module with `/v2` tags and verify proxy/checksum resolution;
- update GoReleaser, Dockerfile, Makefile, docs and all linker paths;
- provide a downstream import migration note;
- keep binaries/containers as the primary supported distribution during the transition.

## Reopen/approval criteria

Operator must explicitly select COA1 or COA2. COA2 requires a separate product/module/release authorization; this report does not grant it.

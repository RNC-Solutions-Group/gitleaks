# CP-V13 Source-Diff and Verification Report

**V13 development HEAD:** `4b77ac90d1d4b37fcdd04e02c9fb21bd233580f1` (tree `7fcaaef841c08d578bab26c8b875429d7b8c6ea3`)
**Integration product commits:** `74a042e522c3d8646d39fd41e8a0ed88d00e5e83` and `d13e14ef7f7eb1e2684472d3c61d74b1c75b8427` (final tree `7fcaaef841c08d578bab26c8b875429d7b8c6ea3`)
**Integration base:** `rnc-mod` at `3a4635f82dcd0a216f9d2b88ebbd72341b3b60cc`
**Final PR:** `audit/v13-integration` → `rnc-mod`

## Product equality

`git diff --exit-code fix/v13-operational-hardening d13e14e` returned zero before the documentation-only provenance addendum. The product tree is object-identical:

```text
development tree: 7fcaaef841c08d578bab26c8b875429d7b8c6ea3
integration tree: 7fcaaef841c08d578bab26c8b875429d7b8c6ea3
```

Only `_assessment/contributor-provenance/v13-*.md` is added after the equality point.

## Difference from V12

```text
30 files changed, 1615 insertions(+), 221 deletions(-)
```

No provider detection rule, generic-secret pattern, entropy policy or benchmark-harness file was changed.

## Package dispositions

| ID | Disposition | Summary |
|---|---|---|
| V13-01 | FIXED | Unused-ignore tracking and opt-in completed-scan diagnostic. |
| V13-02 | FIXED | `gitleaks config` safe effective-config summary. |
| V13-03 | DEFERRED | Multi-format report-plan design; no unsafe implementation. |
| V13-04 | DEFERRED | Bounded context/redaction schema design. |
| V13-05 | FIXED | Opt-in constant SARIF result level. |
| V13-06 | POLICY | Keyed grouping design; unkeyed secret hashes rejected. |
| V13-07 | PARTIAL | Opt-in local file links; decoded highlighting residual. |
| V13-08 | FIXED | Shell-free log-opts parser and parsed-vector validation. |
| V13-09 | PARTIAL | Shallow-history warning; no identity rewrite. |
| V13-10 | PARTIAL | Bare/subdirectory provenance warnings; default scope retained. |
| V13-11 | FIXED | Azure SSH normalization, credential stripping, URL escaping. |
| V13-12 | PARTIAL | Empty/malformed config fail closed; unknown fields remain tolerated. |
| V13-13 | NOT REPRODUCED | Local pre-commit passes; Docker/macOS/PTY root cause unproven. |
| V13-14 | FIXED | Reliable profiling workflow; Linux smoke run. |
| V13-15 | NOT REPRODUCED | Bounded decode probe no depth-correlated RSS growth; defaults retained. |
| V13-16 | FIXED | Bounded deterministic explicit-file concurrency; measured improvement. |
| V13-17 | FIXED | GHCR-only RNC release identity, tag/ancestry/privacy controls. |
| V13-18 | POLICY | RNC SemVer/install matrix; module path retained. |
| V13-19 | PARTIAL | Release-binary portability/packaging plan; no artifact acceptance. |
| V13-20 | DEFERRED | SBOM/provenance/signature verification design; nothing published. |

The full 20-row evidence matrix is `_assessment/v13-remediation-matrix.md`.

## Verification

Executed on the V13 development tree and/or frozen commit:

```text
go build ./...
go test -count=1 ./...
go test -count=1 -tags gore2regex ./...
go vet ./...
go mod verify
go generate ./...
git diff --check
```

All passed. Additional evidence:

- Real pre-commit `--files`/`--all-files` integration: PASS.
- Non-root Linux `go test -race -count=1 ./...` at the pre-fix tree: PASS. The cross-platform alias fix was validated by the Windows/macOS CI jobs that originally exposed `TestGitRepositoryScopeDetection`.
- Profile workflow Linux smoke: `dir` and `git`, CPU/memory/trace artifacts produced: PASS.
- Release workflow YAML parse/static scope: PASS; no workflow executed with registry credentials.
- Release tag/base probe: `v8.26.0-rnc.1` valid; `v8.25.0-rnc.1` rejected; malformed/`-rnc.0` rejected.
- Explicit-file synthetic benchmark (200 × 50 KB, Windows): sequential V12 median ~0.577 s; V13 median ~0.266 s.
- Decode peak-RSS probe (6.67 MB, line-bounded, depths 0/1/3/5): ~18.0–18.8 MB, no depth-correlated growth.

Platform limits: Docker image builds and macOS-native pre-commit scenarios were not executed locally in this cycle; no pass is claimed for them. No release artifact, image, tag, SBOM or signature was published.

## Compatibility and security notes

- No finding, fingerprint, filter or detection-rule change was introduced by V13.
- New CLI controls are opt-in; default JSON/CSV/console/SARIF output is unchanged except deterministic ordering for the explicit-file command.
- Local file links and unused-ignore diagnostics can expose local paths/fingerprints; both are opt-in and console/stderr-only.
- `--log-opts` validation now uses the exact executed vector, closing quote-based bypasses of V12 fail-closed restrictions.
- Release workflow no longer references upstream Docker Hub or its credentials.
- No production release credentials were used.

## Residuals

See `_assessment/v13-remediation-matrix.md` and `RNC-MOD.md` §2.8. Release remains operator-gated; V13 is release-ready in source only.

## Merge recommendation

Use a normal merge commit to preserve the audited RNC integration commit. Rebase and merge is acceptable if linear history is required and SHA rewriting is accepted; do not Squash and merge. Operator approval is mandatory.

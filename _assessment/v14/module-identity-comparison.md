# Go Module Identity Comparison (OP-02)

**Decision state:** EVALUATION ONLY. No production `go.mod`, import path, linker path, repository name or tag was changed.

## Current facts

- Repository: `RNC-Solutions-Group/gitleaks`.
- Declared module: `github.com/zricethezav/gitleaks/v8` (`go.mod`).
- Declared Go minimum: `1.26.0`; CI/container toolchain: `1.27.4`.
- Builds stamp `github.com/zricethezav/gitleaks/v8/version.Version`.
- `go install github.com/zricethezav/gitleaks/v8@...` resolves the upstream module/repository, not this fork.
- `go install github.com/RNC-Solutions-Group/gitleaks@...` is not a supported path because the module declares the upstream identity; a module proxy/repository path mismatch is expected.
- `go run` from an RNC checkout works with the local module; import consumers see the upstream module path.

## COA comparison

| Criterion | COA1: retain upstream module path | COA2: RNC-owned module path | COA3: transitional compatibility |
|---|---|---|---|
| Existing Go imports | No change | Every internal and downstream import changes | Requires replacement/aliasing; Go module identity constraints limit practical support |
| `go install` | Upstream path installs upstream code; RNC uses binaries/source | RNC path could install RNC after a versioned module release | Two identities cannot be one module simultaneously; replace directives are consumer-local |
| `go get`/API consumers | Consumers receive upstream unless they use a replace | Consumers import RNC path explicitly | Replace can redirect a build but is not a transparent distribution identity |
| GitHub repository identity | Repository differs from module path, which is common for forks | Repository and module align | Partial alignment only in consumer build files |
| Module proxy | Upstream proxy path remains upstream | New proxy namespace/cache required | Proxy cannot serve one module under two canonical paths |
| Internal imports | No churn | Repository-wide churn; high risk of missed linker paths | Internal imports remain upstream; external compatibility is consumer-specific |
| GoReleaser/linker | Current paths work | All `-X` paths and docs change | Mixed paths invite stale stamps |
| Future API/library consumers | Clear but surprising: module source is upstream path | Clear RNC identity | Ambiguous and easy to misuse |
| Fork ancestry | Preserved | Preserved | Preserved |
| Migration risk | Low | High; breaking API/release identity | High; not reliably transparent |
| Maintenance cost | Lowest | Highest | Medium-high; every consumer needs special handling |

## Controlled observations

- `go list -m` in the RNC worktree reports the upstream module path, confirming the declared identity.
- All internal imports resolve locally under that path; no RNC-owned module path appears in source.
- The release/container policy already distinguishes RNC artifacts from the Go module path (see `RELEASING.md`).
- A module proxy cannot publish two canonical identities for one module; COA3 would therefore be a consumer-side replace strategy, not a true dual identity.

## Recommendation

**COA1: retain the upstream Go module path for V14.** Distribute RNC-Mod as source, binaries and the RNC GHCR image under the RNC repository identity. This preserves import/build compatibility and avoids a sweeping rename with no release benefit.

COA3 should be documented only as a consumer-side `replace` pattern for controlled internal builds; it should not be advertised as transparent RNC module installation. COA2 remains a future operator decision requiring a versioned migration, rewritten imports/linker paths, module proxy publication and explicit downstream coordination.

## Reopen criteria

Reopen if RNC publishes a Go library/API with an RNC-owned import contract, or if a proxy/registry requirement cannot be met while retaining the upstream module path.

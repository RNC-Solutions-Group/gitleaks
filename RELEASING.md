# RNC-Mod Release and Installation Policy

This document defines the **candidate** RNC-Mod distribution contract. It is not
a published release, and no tag, image, package, SBOM or signature is created by
this document or by the V13 development cycle.

## Release identity

- Repository: `RNC-Solutions-Group/gitleaks`
- Protected integration/release branch: `rnc-mod`
- Container registry: `ghcr.io/rnc-solutions-group/gitleaks` only
- No Docker Hub publication and no upstream `zricethezav/gitleaks` image is
  presented as an RNC-Mod artifact.
- Go module path remains `github.com/zricethezav/gitleaks/v8`; this preserves
  import compatibility and avoids a sweeping module-path rename. The RNC
  repository is not the canonical source for that module path.

## Version convention

RNC release tags use:

```text
vMAJOR.MINOR.PATCH-rnc.N
```

Examples: `v8.26.0-rnc.1`, `v8.26.1-rnc.2`.

The numeric base MUST be **strictly newer** than the embedded default
configuration's `minVersion`. A tag such as `v8.25.0-rnc.1` has lower SemVer
precedence than `v8.25.0` because the prerelease suffix is compared after the
base version; it must not be used. The release workflow validates this rule and
requires the release commit to be reachable from protected `rnc-mod`.

Container tags are derived from the same version without the leading `v`:

- immutable: `8.26.0-rnc.1`
- moving aliases: `8.26`, `8`
- no `latest` alias is created

Docker tag grammar permits these values; Git tag and container tag validation
are deliberately separate.

## Supported installation matrix

| Method | Supported | Command/contract | Consequences |
|---|---|---|---|
| Build from RNC source | Yes | `git clone https://github.com/RNC-Solutions-Group/gitleaks.git`; `make build` or `go build` | Requires Go `1.26.0` minimum; CI builds with Go `1.27.4`. The binary reports the RNC source build version. |
| RNC release binary | Yes, after an actual release | Download the versioned archive and verify the published checksum | No release artifact exists until the operator publishes one. |
| RNC GHCR image | Yes, after an actual release | `docker pull ghcr.io/rnc-solutions-group/gitleaks:<version>` | Default image is root-based; `runtime-nonroot` is an opt-in build target. |
| `go install github.com/zricethezav/gitleaks/v8@...` | Not an RNC install | Installs the upstream module, not RNC-Mod | Do not use this command to claim RNC distribution. |
| `go install github.com/RNC-Solutions-Group/gitleaks@...` | Not supported | Declared module path is upstream | Use a source checkout or a release artifact. |
| `go run` from RNC checkout | Developer use | `go run . <command>` | Version may be a development stamp; not a release artifact. |
| Docker Hub `zricethezav/gitleaks` | Not an RNC image | None | Do not label it RNC-Mod. |

The Go module path change remains a separate operator decision. Renaming the
repository or detaching the module path merely for branding is not authorized.

## Toolchain and platform baseline

- Minimum Go directive: `1.26.0`.
- CI/container build toolchain: Go `1.27.4` (includes current net/http security
  fixes).
- Minimum Git for the explicit `--diff-merges=first-parent` opt-in: Git 2.31.
- Standard regexp and `gore2regex` builds are both supported; the latter retains
  the constrained-runtime address-space limitation tracked as issue #2017.
- Release artifacts should cover linux/amd64, linux/arm64, darwin/amd64,
  darwin/arm64, windows/amd64 and windows/arm64. Other architectures require a
  tested justification; they are not implied by the source build.

## Binary packaging assessment

The checked-in `.goreleaser.yml` is a **candidate**, not an invoked release
pipeline. It currently builds `gore2regex` artifacts across legacy architectures
and has not been validated against the constrained-runtime RE2 issue. A future
release must either:

- publish a standard-regexp binary set plus clearly identified RE2 variants, or
- prove the RE2 binaries meet the documented runtime memory/address-space
  contract; and
- smoke-test every OS/architecture artifact, record checksums, and verify the
  version stamp.

No GoReleaser run or artifact publication is claimed by V13.

## Container publication controls

`.github/workflows/release.yml`:

- is triggered only by a published GitHub Release;
- validates the RNC tag convention and numeric base against `minVersion`;
- verifies the release commit is an ancestor of `rnc-mod`;
- requests only `contents: read` and `packages: write`;
- logs in only to GHCR with `GITHUB_TOKEN`;
- publishes multiarch `linux/amd64,linux/arm64` with version/minor/major tags;
- does not use Docker Hub credentials or an upstream image name.

No production registry credentials are used during development. No image is
published by V13.

## Proposed attestations (not yet enabled)

The following are the acceptance targets for a later operator-authorized release
workflow. A workflow merely containing an `sbom`/`provenance` option is not proof
that an attestation exists.

1. Build the multiarch image and record its digest.
2. Generate a CycloneDX or SPDX SBOM and attach it to the release.
3. Generate GitHub/Sigstore provenance for the image digest.
4. Generate checksums for every downloadable archive.
5. Optionally sign checksums/SBOM with keyless Cosign using OIDC, with minimal
   `id-token: write` permission isolated to the signing job.
6. Independently verify the digest, SBOM and signature as a separate consumer.
7. Define retention, revocation and rollback procedures for a bad artifact.

Verification must be recorded against exact digests and release assets, not
against a successful build log alone.

## Release checklist

- [ ] Operator approves the release version and tag base.
- [ ] `rnc-mod` is the source commit and all required checks pass.
- [ ] Version base is strictly newer than `minVersion`.
- [ ] Local source, standard/RE2 tests, race tests and Docker smoke tests pass.
- [ ] Release notes list behavior changes and residual limitations.
- [ ] Checksums, SBOM, provenance and optional signatures exist.
- [ ] A separate consumer verifies the published digest and assets.
- [ ] Rollback/retention plan is recorded.

Until every item is true, RNC-Mod is **release-ready in source**, not released.

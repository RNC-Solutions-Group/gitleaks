# Round 2 — Go 1.27.3 toolchain, Alpine 3.24 runtime, and Lip Gloss v2

**Status:** Experimental implementation checkpoint complete — **not merged, not released, not production-approved.**
**Branch:** `test/deps-upgrade-all`
**Round 2 base commit:** `c5911d4c0decb848a99397a306079760c9f9d4b2` (round 1 documentation HEAD)
**Previous round:** preserved unchanged in [`README.md`](README.md), [`dependency-matrix.md`](dependency-matrix.md), [`compatibility-report.md`](compatibility-report.md), [`regression-register.md`](regression-register.md), [`performance-report.md`](performance-report.md), and [`upgrade-assessment.md`](upgrade-assessment.md).

Round 1 upgraded module versions **within existing module paths** and deliberately changed nothing outside `go.mod`, `go.sum`, and documentation. Round 2 continues from that point and is the first round that intentionally changes application source, the container build, and CI. It does not revise or delete round 1 results; where the two differ, both are recorded.

## What changed in round 2

| Area | Change |
| --- | --- |
| Go toolchain | Standardized build/test on **Go 1.27.3** (minimum remains `go 1.26.0`) |
| Container build | Reproducible multi-stage build on `golang:1.27.3-alpine3.24` and patched `alpine:3.24` |
| Runtime hardening | Package-pruned runtime, new `.dockerignore`, documented privilege decision |
| Lip Gloss | `github.com/charmbracelet/lipgloss` v1.1.0 → **`charm.land/lipgloss/v2` v2.0.6** |
| Display output | Finding output now passes through a color-downsampling writer |
| CI | Rewritten workflow with 8 jobs, slash-branch trigger, pinned tooling |
| Tests | One focused regression test added for ANSI-free redirected output |

## 1. Final dependency inventory

**Verified** with `go list -m all`, `go list -deps`, and `go mod why` on the round 2 working tree.

| Metric | Round 1 | Round 2 |
| --- | --- | --- |
| `go.mod` `require` entries | 69 | 70 |
| Modules in the module graph | 145 | 145 |
| Modules linked into the default binary | 66 | 65 |
| Direct modules | 16 | 17 |
| `go mod verify` | pass | pass |
| `go mod tidy -diff` | clean | clean |

### Direct-module changes

| Module | Before | After | Reason |
| --- | --- | --- | --- |
| `charm.land/lipgloss/v2` | — | v2.0.6 | Major-version module-path migration (replaces v1) |
| `github.com/charmbracelet/lipgloss` | v1.1.0 | removed | Superseded by `charm.land/lipgloss/v2` |
| `github.com/charmbracelet/colorprofile` | v0.4.3 (indirect) | v0.4.3 (**direct**) | Now imported directly for output downsampling |

All other direct modules were already at their latest stable versions and were left unchanged:

`aho-corasick` v1.0.3, `sprig/v3` v3.3.0, `semgroup` v1.3.0, `gitleaks/go-gitdiff` v0.9.1, `go-cmp` v0.7.0, `filetype` v1.1.3, `go-version` v1.9.0, `archives` v0.1.5, `zerolog` v1.35.1, `cobra` v1.10.2, `viper` v1.21.0, `testify` v1.12.1, `x/exp` (pinned pseudo-version), `lucasjones/reggen` (pinned pseudo-version), `wasilibs/go-re2` v1.13.0.

### Transitive changes caused by the Lip Gloss v2 move

| Module | Before | After | Note |
| --- | --- | --- | --- |
| `github.com/muesli/termenv` | v0.16.0 | removed | v1-only terminal layer |
| `github.com/charmbracelet/x/cellbuf` | v0.0.15 | removed | v1-only renderer cell buffer |
| `github.com/aymanbagabas/go-osc52/v2` | v2.0.1 | removed | v1-only |
| `github.com/charmbracelet/ultraviolet` | — | added (indirect) | v2 terminal query/cancel-reader layer |
| `github.com/muesli/cancelreader` | — | added (indirect) | v2 dependency |
| `github.com/charmbracelet/x/termios`, `x/windows` | — | added (indirect) | v2 platform support |

Net effect on the **linked** binary: one fewer module (`66 → 65`). The terminal stack is still present, but through Lip Gloss v2's own dependencies rather than the v1 renderer stack.

### Graph-only dependencies (not linked)

The 53 indirect modules previously identified as parent-pinned remain parent-pinned. They were **not** forced to their newest versions. The Google Cloud / OpenTelemetry sub-graph that round 1 observed is still present in the module graph and still **not linked into the scanner binary**:

- **Verified:** `go list -deps` for the default build returns no `cloud.google.com`, `google.golang.org`, `go.opentelemetry.io`, `grpc`, `spiffe`, or `envoy` module.

## 2. Go / toolchain compatibility matrix

The minimum supported Go version and the pinned build/test toolchain are deliberately different things.

| Item | Value | Status |
| --- | --- | --- |
| `go.mod` minimum (`go` directive) | **1.26.0** | Verified (unchanged by this round) |
| Pinned build/test toolchain | **1.27.3** | Verified locally and in the container build stage |
| CI toolchain | **1.27.3** (asserted at runtime) | Inferred (workflow not executed on GitHub in this round) |
| Container build image | `golang:1.27.3-alpine3.24` | Verified (`go version` → `go1.27.3 linux/amd64`) |
| Go 1.26.0 exact build | — | **Not executed.** Inferred compatible from the `go 1.26.0` directive |
| Go 1.25 / 1.24 | — | Verified failing at the `go.mod` gate in round 1 |

Toolchain evidence:

- Local host: `go version go1.27.3 windows/amd64`, `go env GOVERSION` → `go1.27.3`.
- Builder image `golang:1.27.3-alpine3.24`, digest `sha256:fa3d596898d9f126af7c01d14ee7597bab3f39cccf392ef36b69fe77cd681f18`, reports `go version go1.27.3 linux/amd64`.

The `go` directive was **not** raised. It was checked that the chosen Lip Gloss release does not force a raise:

| Lip Gloss release | Module Go requirement | Decision |
| --- | --- | --- |
| `charm.land/lipgloss/v2` v2.0.6 | `go 1.25.0` | **Selected** — compatible with `go 1.26.0` |
| `charm.land/lipgloss/v2` v2.1.0 | `go 1.26.7` | Deliberately not taken — would force a `go` directive raise |

No `toolchain` directive was added to `go.mod`: a `toolchain go1.27.3` line would raise the effective requirement for consumers running with `GOTOOLCHAIN=local`, contradicting the stated minimum of 1.26.0. The pin is expressed in CI and in the Dockerfile instead.

## 3. Alpine image provenance and patch methodology

### Candidate image, as provided

The locally available candidate was accessible and matches the stated identity:

| Property | Value |
| --- | --- |
| Tag | `alpine:3.24-patched-test` |
| Image ID | `sha256:57910b1b435b97201eeb141f4e0f5e1c1e3ed87f4c8a2be22e14322bc8e059c0` |
| Architecture | `amd64/linux` |
| Created | 2026-11-07 (approximately 46 hours before inspection) |
| Size | 13.2 MB |

### Provenance (verified)

Image history, newest first:

1. `RUN /bin/sh -c apk upgrade --no-cache`
2. `CMD ["/bin/sh"]`
3. `ADD alpine-minirootfs-3.24.4-x86_64.tar.gz /`

So the candidate is **stock official `alpine:3.24` (minirootfs 3.24.4) with a single `apk upgrade --no-cache` layer** — nothing bespoke, nothing private.

### Patch state (verified by package diff)

Comparing `apk list --installed` between `alpine:3.24` and `alpine:3.24-patched-test`, the **only** difference is:

| Package | Stock `alpine:3.24` | `alpine:3.24-patched-test` |
| --- | --- | --- |
| `zlib` | 1.3.2-r1 | **1.4-r0** |

Every other package is identical (16 packages in the bare image). The "patch" is therefore a single Alpine point release of zlib delivered through the public repository.

### Reproducibility (verified)

`docker run --rm alpine:3.24 sh -c 'apk --no-cache upgrade'` today reports:

```
(1/1) Upgrading zlib (1.3.2-r1 -> 1.4-r0)
OK: 8212 KiB in 16 packages
```

So the candidate's patch level is reproducible **from the public Alpine v3.24 APK repositories**, with no local artifacts:

- The local image ID **must not** be assumed to exist on GitHub runners; CI never references it.
- The repository `Dockerfile` reconstructs the same patch level with `apk upgrade --no-cache` during image construction.
- Because `apk upgrade` resolves against the live repository, a rebuild can only move **forward** to newer patch levels, never backward.

### Implemented Dockerfile

```dockerfile
FROM golang:1.27.3-alpine3.24 AS build    # pinned build toolchain
...
FROM alpine:3.24
RUN apk upgrade --no-cache && \
    apk add --no-cache bash git openssh-client
```

Runtime image as built (`linux/amd64`): Alpine **3.24.4**, 38 installed packages, 66.5 MB, `zlib-1.4-r0` present.

Retained runtime packages and why:

| Package | Retained | Reason |
| --- | --- | --- |
| `bash` | yes | Existing image contract — shell-based invocations and wrappers use it |
| `git` (+ `git-init-template`) | yes | Git repository scanning, including commit history |
| `openssh-client` | yes | Scanning repositories cloned over SSH |
| `ca-certificates-bundle` | yes (base) | TLS for git/SSH and HTTPS remotes |

Dropped: none were added beyond the above; the previous image also installed `bash git openssh-client`. Build-only tooling (`git` in the builder stage) never reaches the runtime image.

### Privileges (evaluated, deliberately unchanged)

The image still runs as **root** and still sets `git config --global --add safe.directory '*'`.

- Changing to a non-root user would break scanning of bind-mounted workspaces owned by the caller's uid, which is the primary container use case for this tool.
- Narrowing `safe.directory` would cause git to refuse mounted repositories as "dubious ownership".

This is recorded as a **remaining limitation and an adoption blocker for hardening**, not silently accepted as ideal.

### `.dockerignore`

Added. It excludes documentation/assessment trees, fixtures, report samples, editor metadata, and build/test artifacts. It deliberately **does not** exclude `.git` (used for version stamping) or `config/gitleaks.toml` (embedded via `go:embed`) — see the verification below.

### Version stamping (verified)

- With `--build-arg VERSION=8.99.0-test` → `gitleaks version 8.99.0-test`.
- Without the build argument on this tagless fork clone → `gitleaks version unknown` (graceful; the previous Dockerfile would have produced an empty stamp because `git` was absent from the Alpine builder).
- CI passes `VERSION=ci-${GITHUB_SHA}` so images are always traceable.

## 4. Lip Gloss v2 migration outcome

**Outcome: migrated successfully to `charm.land/lipgloss/v2` v2.0.6.**

### What the migration required

The Style API surface used by gitleaks — `NewStyle`, `SetString`, `Bold`, `Italic`, `Foreground`, `Color`, and `Style.String()` — is unchanged in v2. The mechanically required changes were small:

| File | Change |
| --- | --- |
| `detect/utils.go` | Import path; write finding output through a downsampling writer |
| `report/finding.go` | Import path; write required-findings output through a downsampling writer; new `ColorWriter()` helper |
| `detect/utils_test.go` | New regression test |

### The behavioral trap (and why it is handled explicitly)

This is the part of the migration that is **not** a mechanical import swap.

| Behavior | Lip Gloss v1 | Lip Gloss v2 |
| --- | --- | --- |
| Where downsampling happens | Inside `Style.Render()` via the global renderer | At the output layer only |
| `Style.Render()` / `Style.String()` when stdout is **not** a TTY | Produces plain text (profile detected as Ascii) | Produces **full-fidelity truecolor ANSI** |
| Piped / redirected output | No escape sequences | Escape sequences unless written through a downsampling writer |

Measured with an isolated probe:

- v1, stdout is a pipe → profile `Ascii`, `Style.String()` contains **0** ESC bytes.
- v2, stdout is a pipe → `Style.String()` contains **2** ESC bytes (`\x1b[1;3;38;2;240;92;7m…`).
- v2 written through `colorprofile.NewWriter` → **0** ESC bytes for a `NoTTY` profile, truecolor preserved for a truecolor profile, 256-colour downsampling for ANSI256.

A naive import-only migration would therefore have injected ANSI escape sequences into piped output, redirected files, and CI logs. That regression was **observed directly** (see the mutation check below) and fixed by routing all finding display output through `report.ColorWriter()`:

```go
// ColorWriter returns a writer for standard output that down-samples ANSI
// color sequences to the color profile supported by the current terminal.
func ColorWriter() io.Writer {
	return colorprofile.NewWriter(os.Stdout, os.Environ())
}
```

The writer is created per call, so it binds the **current** `os.Stdout`. This matters because the existing test harness (`captureStdout`) swaps `os.Stdout` for a pipe; a package-level writer bound at init would have written to the real terminal and broken that capture.

### Verified rendering behavior

Compared against a binary built from the round 2 base commit (v1) on the same fixtures:

| Scenario | v1 (baseline) | v2 (migrated) | Result |
| --- | --- | --- | --- |
| Piped/redirected, default env | 0 ESC, 449 bytes | 0 ESC, 449 bytes | **Byte-identical** |
| Piped, `TERM=xterm-256color`, `CLICOLOR_FORCE=1` | 0 ESC (v1 ignores the force when not a TTY) | 8 ESC, downsampled to 256-colour (`38;5;221`, `38;5;202`) | Colour now honoured and correctly downsampled; no truecolor leak |
| `--no-color`, forced colour env | plain | plain (0 ESC) | `--no-color` still wins |
| `--redact`, redirected | `REDACTED` | `REDACTED` without escapes | Unchanged |
| JSON / SARIF report files | no ANSI | no ANSI | Unchanged (report writers never used Lip Gloss) |
| `dir` scan and `git` scan | identical | identical | Unchanged |

Note the only divergence is the **explicit** `CLICOLOR_FORCE` override on a non-TTY, where v2 now emits downstream-correct downsampled colour instead of silently ignoring the request. Default redirected output is unchanged.

### Regression test (mutation-verified)

`Test_printFinding_colorOutput` in `detect/utils_test.go`:

1. Asserts that a v2 style renders ANSI sequences (so the test cannot silently become vacuous).
2. Asserts that `printFinding` output captured through a pipe contains the rule ID and secret but **no** `\x1b[` sequence.
3. Asserts the same for the `noColor` path.

The test was validated by temporarily writing directly to `os.Stdout` instead of the downsampling writer: it failed with truecolor escapes (`\x1b[38;2;245;212;69m…`) in the captured output, then passed again once restored. The mutation was reverted before commit.

### Migration cost

- One fewer linked module, four transitive modules removed, six added.
- No public CLI behavior change beyond the forced-colour case above.
- No theme/color redesign was needed; the two hex colours are unchanged.

## 5. Other major-version migration candidates

**Audit result: there are no other direct modules with a newer major-version module path.**

Every direct module was checked for available updates. All are already at their latest stable version; the only direct major-version move available in the ecosystem is Lip Gloss v2, which this round performs.

| Candidate | Finding | Action |
| --- | --- | --- |
| `github.com/charmbracelet/lipgloss` → `charm.land/lipgloss/v2` | Major module-path change, official v2 migration guide | **Migrated** (v2.0.6) |
| `github.com/charmbracelet/lipgloss/v2` v2.1.0 | Requires Go 1.26.7, forcing a `go` directive raise | **Deferred** (would violate the 1.26.0 minimum) |
| `wasilibs/go-re2` | v1.13.0 is latest, already migrated in round 1 | No change |
| `mholt/archives` | v0.1.5 is latest | No change |
| `spf13/viper`, `spf13/cobra`, `rs/zerolog`, `stretchr/testify`, `Masterminds/sprig/v3` | Latest majors already in use | No change |

### Indirect-module candidates evaluated and deferred

Indirect modules were not bulk-upgraded. Each was classified by linkage before being judged:

| Module | Installed → available | Linked? | Decision | Reason |
| --- | --- | --- | --- | --- |
| `golang.org/x/crypto` | v0.58.0 → `v0.58.1-0.20261009181718-…` | yes | **Deferred** | Fix exists only as a pseudo-version; advisory (`GO-2026-5932`) concerns `openpgp`, which gitleaks does not import and cannot reach |
| `github.com/andybalholm/brotli` | v1.2.8 → v1.3.0 | yes | **Deferred** | Linked decompressor, but no advisory and no functional need; round 1 already took the security-relevant bump to v1.2.8 |
| `github.com/klauspost/compress` | v1.20.4 (latest) | yes | No change | Already current |
| `golang.org/x/net`, `x/oauth2`, `x/term`, `x/time`, `x/mod` | behind latest | **no** | **Deferred** | Graph-only; pulled by the cloud/OTel sub-graph, not linked into the binary |
| ~40 cloud/OTel/`google.golang.org` modules | behind latest | **no** | **Deferred** | Graph-only serialization surface; upgrading them would inflate the diff without changing the binary |

**Verified:** no cloud/OTel module is linked into the default binary, so none of these deferrals affect runtime behavior or binary content.

### Unused direct dependencies

**Verified:** every direct module has a real import path (`go mod why -m` resolves each one to a gitleaks package). No unused direct dependency was found; none was removed.

## 6. Test and build results by platform

All commands were executed directly; no Python scanner harness, Levain, Leaky Repo, or external fixture suite was used.

| Check | Platform | Result |
| --- | --- | --- |
| `go build ./...` | Windows amd64, Go 1.27.3 | **Pass** |
| `go test ./...` | Windows amd64, Go 1.27.3 | **Pass** — all packages ok |
| `go test -race ./...` (CGO) | Linux amd64, Go 1.27.3 container | **No data races.** 2 environment-limited test failures, reproduced identically at the base commit (see §6.1) |
| `go test -race ./...` | Windows amd64 | **Not supported here** — `-race requires cgo`, no C compiler on host |
| `go vet ./...` | Windows amd64 | 2 pre-existing `lostcancel` findings only (unchanged) |
| Standard application build | Windows amd64 | **Pass** |
| RE2 (`gore2regex`) application build | Windows amd64 | **Pass**; engine confirmed as `github.com/wasilibs/go-re2` |
| `go build -tags gore2regex ./...` | Windows amd64 | **Fail** — known baseline failure (see §7) |
| `go test -tags gore2regex ./...` | Windows amd64 | **Fail** — same known baseline failure |
| `go mod verify` | Windows amd64 | **Pass** — all modules verified |
| `go mod tidy -diff` | Windows amd64 | **Clean** |
| CLI `--version` / `--help` | Windows binary | **Pass** |
| Docker image build | linux/amd64 | **Pass** |
| Docker image build | linux/arm64 | **Pass** (QEMU emulation) |
| Docker runtime smoke | linux/amd64 | **Pass** — below |
| Docker runtime smoke | linux/arm64 | **Pass** (QEMU emulation) |
| Config generation drift | — | Not run locally (would rewrite a tracked generated file); validated in CI |

### 6.1 Race tests

Race tests were run on Linux/amd64 inside a `golang:1.27.3` container with `CGO_ENABLED=1`. Windows-native race testing is unavailable on this host because no C toolchain (`gcc`/`clang`) is installed and the environment has `CGO_ENABLED=0`; this is recorded as a tooling limitation.

Result: **the race detector reported no data races** (`DATA RACE` occurrences: 0). Two tests fail inside the Linux container:

- `TestFromGit` — the `testdata/repos/*` fixture repositories are empty in this checkout ("0 commits scanned").
- `TestDetectWithSymlinks` — symlinks are not materialized in this working tree.

Both are properties of **this Windows checkout**, not of the code:

| Evidence | Value |
| --- | --- |
| `git config core.symlinks` | `false` → symlinks checked out as plain files |
| `git config core.autocrlf` | `true` |
| Fixture repos | empty; nested/submodule repositories not initialized |
| Same tests at base commit `c5911d4` in the same container | **identical failures**, same `Test:` names |
| Same tests on the Windows host (`go test ./...`) | **pass** |

Because the identical failures occur at the pre-change commit in the same container, they are classified as environment-limited, **not** introduced regressions. They ran green on the Windows host for the full suite, and on a normal Linux checkout (where symlinks and fixture repositories are present) they are expected to pass — but that specific Linux pass is **not verified in this round**.

### 6.2 Docker runtime smoke (verified)

Image `gitleaks:deps-upgrade-amd64` (`linux/amd64`, Alpine 3.24.4, 66.5 MB), built with `VERSION=8.99.0-test`:

| Check | Result |
| --- | --- |
| `--version` | `gitleaks version 8.99.0-test` |
| `--help` | renders usage |
| `bash --version` | GNU bash 5.3.9 |
| `git --version` | git 2.54.0 |
| `ssh -V` | OpenSSH 10.3p1, OpenSSL 3.5.9 |
| Mounted workspace scan, clean target | exit `0` |
| Mounted workspace scan, synthetic finding | exit `1` |
| Mounted git repository scan | exit `1` (finding detected) |
| `--report-format json` | valid JSON with `RuleID` |
| `--report-format sarif` | valid SARIF |
| `--verbose --redact` on redirected output | `REDACTED`, no ANSI |

The arm64 image was built and smoke-tested under QEMU emulation; see §6.3.

### 6.3 Docker runtime smoke, linux/arm64 (verified)

Image `gitleaks:deps-upgrade-arm64` (`arm64/linux`), built with `VERSION=8.99.0-test` under QEMU:

| Check | Result |
| --- | --- |
| `uname -m` in container | `aarch64` |
| `/etc/alpine-release` | `3.24.4` |
| `zlib` package | `1.4-r0` (patched) |
| `--version` | `gitleaks version 8.99.0-test` |
| `--help` | renders usage |
| `bash --version` | GNU bash 5.3.9 (aarch64) |
| `git --version` | git 2.54.0 |
| Mounted workspace scan, finding | exit `1` |
| Mounted workspace scan, clean | exit `0` |

Both architectures therefore build reproducibly and exercise the same detection, report, and CLI paths.

## 7. Existing failures vs introduced regressions

### Pre-existing failures (reproduced, not fixed, not suppressed)

| # | Failure | Evidence | Classification |
| --- | --- | --- | --- |
| 1 | `go vet ./...` reports two test-only `lostcancel` diagnostics | `detect/detect_test.go:2083:4` and `detect/detect_test.go:2122:3` | Pre-existing (round 1 / CP5 baseline), unchanged |
| 2 | `go build -tags gore2regex ./...` and `go test -tags gore2regex ./...` fail | `cmd/generate/config/rules/airtable.go:30` and `perplexity.go:16` — stdlib `*regexp.Regexp` where the tagged wrapper type is expected | Pre-existing, unchanged; the RE2 **application binary still builds** |
| 3 | `TestTranslateExtend` is order-sensitive | Full `config` package with `-shuffle=1794264498871797264` fails on `valid/extend`; the same test alone passes | Pre-existing (`CP5-F07` global-state/test-order issue), unchanged |

### Introduced regressions

**None retained.** One regression was introduced *and fixed within this round*, and is documented rather than hidden: the initial Lip Gloss v2 import swap emitted truecolor ANSI into redirected output. It was detected by direct byte comparison and by the new regression test (which was itself mutation-verified), then fixed by routing output through a downsampling writer.

### Not silently suppressed

- CI does not add `--tags gore2regex` to the full package graph, because that graph is a known baseline failure; it builds the RE2 **application binary** instead and the workflow comment states why.
- `go vet` findings remain visible; nothing was excluded to make the suite appear green.
- No test was modified to change an asserted result. The only test change is the **addition** of `Test_printFinding_colorOutput`.

## 8. Container vulnerability results

### Scanner provenance

| Item | Value |
| --- | --- |
| Scanner | Trivy **0.76.2** (`aquasec/trivy:0.76.2`) |
| Advisory source | `mirror.gcr.io/aquasec/trivy-db:2` |
| Database fetched | 2026-11-10T01:55:38Z |
| Image scanned | `gitleaks:deps-upgrade-amd64` (Alpine 3.24.4, `linux/amd64`) |
| Scanners | OS packages (`alpine 3.24`, 36 packages) and Go binary |

### Results by severity

| Target | Type | CRITICAL | HIGH | MEDIUM | LOW | UNKNOWN | Total |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Base image `alpine 3.24.4` | OS packages | 0 | 0 | 0 | 0 | 0 | **0** |
| `/usr/bin/gitleaks` | Go binary | 0 | 0 | 0 | 0 | 1 | **1** |

The single application finding is `GO-2026-5932` against `golang.org/x/crypto` v0.58.0 ("`openpgp` package is unmaintained, unsafe by design"). It is:

- **OS-package vs application separation:** this is an application/module finding, not an OS-package finding.
- **Not reachable:** `govulncheck` v1.9.0 (source mode, DB updated 2026-11-05) reports **0 reachable vulnerabilities** and **0 imported-package vulnerabilities**; the advisory appears only in the module-level list.
- **Pre-existing:** `x/crypto` v0.58.0 is unchanged from the round 2 base commit.

### Comparison with the previously patched image

All three images were scanned with the same Trivy 0.76.2 build and the same advisory database snapshot:

| Image | OS packages scanned | OS-package findings | Detail |
| --- | --- | --- | --- |
| `alpine:3.24` (stock) | 16 | **1 MEDIUM** | `CVE-2026-4738`, `zlib` 1.3.2-r1, fixed in 1.4-r0 |
| `alpine:3.24-patched-test` (candidate) | 16 | **0** | `zlib` 1.4-r0 |
| `gitleaks:deps-upgrade-amd64` (this round) | 36 | **0** | `zlib` 1.4-r0 + scan-time packages |

**The reproducible image reaches the same OS-package result as the locally patched candidate**, because it applies the same `apk upgrade` and therefore the same `zlib` 1.4-r0 that clears `CVE-2026-4738`. The stock image demonstrates that the patch step is load-bearing: without `apk upgrade` the unpatched `zlib` is reported. This is a measured result for this scan and this database snapshot only.

### Explicit caveat

**No claim of "zero vulnerabilities overall" is made.** The OS layer is clean; the binary carries one unreachable module advisory. Any future assertion of "zero CVEs" must be re-measured against a fresh advisory database.

## 9. Remaining limitations

**Verified limitations**

1. Windows-native `-race` testing is unavailable in this environment (no C compiler, `CGO_ENABLED=0`); race coverage was obtained on Linux, where **no data races** were reported.
2. This Windows checkout has `core.symlinks=false` and uninitialized fixture repositories, so `TestDetectWithSymlinks` and `TestFromGit` fail inside a Linux container. The identical failures reproduce at the base commit, so they are environment-limited rather than regressions; a clean Linux checkout was not exercised in this round.
3. This fork clone has **no git tags**, so a container built without `VERSION` reports `unknown`. CI supplies `VERSION`.
4. The image still runs as **root** with `safe.directory '*'` — required for mounted-workspace scanning, but not a hardened posture.
5. `bash` is retained for backward compatibility with existing image usage; gitleaks itself does not need it.
6. Container scanning depends on the public Trivy advisory database; a database outage degrades the assessment.
7. `.git` is included in the Docker build context so version stamping can fall back to `git describe`; this makes the context larger than strictly necessary for compilation.

**Inferred, not verified**

1. Go 1.26.0 (the declared minimum) builds this tree — inferred from the `go 1.26.0` directive; only 1.27.3 was executed.
2. CI workflow behavior on GitHub-hosted runners (QEMU arm64 execution, `load: true` multi-arch builds, Trivy via mounted Docker socket) — the jobs were validated for YAML and locally where possible, but not executed on GitHub in this round.
3. Windows console colour rendering on a real VT-capable console — the downsampling logic is platform-agnostic and Windows env detection lives in `colorprofile`, but no Windows PTY was available to observe rendered colour.
4. Visual appearance on a real TTY. ANSI presence/absence and downsampling level were verified programmatically; human-visual rendering was not.

**Not executed this round**

1. Go 1.26.0 toolchain build and test.
2. `go generate ./...` drift check locally (it rewrites a tracked, generated config file; delegated to CI).
3. Python scanner harness, Levain, Leaky Repo, and external fixture suites — **explicitly out of scope**.
4. Publishing any image, package, or release.

**Deliberately deferred**

1. Lip Gloss v2.1.0 (Go 1.26.7 requirement).
2. `golang.org/x/crypto` pseudo-version bump (unreachable `openpgp` advisory).
3. `andybalholm/brotli` v1.3.0 (no advisory, no functional need).
4. The cloud/OpenTelemetry module sub-graph (not linked).
5. Non-root runtime user (would break mounted-workspace scanning as currently used).
6. The three pre-existing failures in §7 — they belong to separate remediation work and were intentionally not "fixed" by the dependency upgrade.

## 10. Adoption recommendation

**Recommendation: adopt for continued experimentation; do not merge or release as-is.**

What is now demonstrably better than the round 2 base commit:

1. **The toolchain gate is resolved.** CI and the container build no longer use Go 1.24/1.25, which fail immediately at `go.mod` (`requires go >= 1.26.0`). Both are pinned to the tested 1.27.3.
2. **The runtime image is reproducible and patched.** Alpine 3.24.4 with `apk upgrade` reproduces the previously hand-built patched image from public sources, and scans clean at the OS layer.
3. **The Lip Gloss v2 migration is complete and behaviorally safe.** Piped output is byte-identical to v1, `--no-color` still wins, reports are unaffected, and the one real regression risk (ANSI leakage into redirected output) is fixed and covered by a mutation-verified test.
4. **CI actually triggers on this branch.** `branches: ["*"]` could never match `test/deps-upgrade-all`; the workflow now uses `**` and asserts the Go version at runtime.

Conditions before any merge or release:

1. Execute the workflow on GitHub and confirm all eight jobs, including arm64 under QEMU and the Trivy job.
2. Decide the runtime posture deliberately: keep root + `bash` for mounted-workspace compatibility, or adopt a non-root variant with a documented invocation change.
3. Keep the three pre-existing failures (vet `lostcancel`, `gore2regex` full graph, `TestTranslateExtend` ordering) as separate remediation items. This round neither fixes nor hides them.
4. Re-measure container and module advisories against a fresh database before making any "zero CVE" statement.
5. Treat the round 1 performance and detection-parity results as still applicable to the standard engine; the Lip Gloss change affects display output only and does not touch detection paths.

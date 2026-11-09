# CP0 — Repository and environment reconnaissance

**Status:** Draft — ready for owner review  
**Assessment date:** 2026-11-09 (America/New_York)  
**Repository / branch:** `RNC-Solutions-Group/gitleaks` / `feat/issues-and-pr`  
**Assessed code SHA:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` (canonical code baseline)  
**Branch HEAD at assessment:** `a89c5f7653ba1037f747df285312eda5831fbbe5` (documentation-only commits)  
**Environment:** Windows 11 Home 10.0.26300, AMD64, 8 logical CPUs, 31.7 GB RAM; Git 2.45.2.windows.1; Docker 29.9.2 (Linux containers, WSL2 kernel 6.18.54.1); container Go 1.27.3; **no native Go toolchain**  
**Assignment:** [`assignment-CP0.md`](assignment-CP0.md)

---

## Executive assessment

CP0 established a reproducible starting manifest for the fork. The canonical code baseline is `5f6ec56` (fork `master`); the assessment branch adds 60 commits that touch only documentation (`AGENTS.md`, `README.md`, `_assessment/**`, `_backlog/**`) and no application source. The working tree was clean before and after all work, and no protected path was modified.

The decisive environmental fact is that **no native Go toolchain, module cache, C compiler, or `make` exists on this host**. A usable Go 1.27.3 toolchain does exist in the locally cached `golang:1.27.3` image and works with network disabled, but the module cache is empty and all 16 direct / 56 indirect module requirements are unresolved. An offline `go build ./...` therefore fails at dependency resolution; `go test` and `go vet` are recorded as **NOT RUN** for that reason. A dependency-free `gofmt` probe parsed all 214 Go files successfully and found zero formatting deviations after normalizing the checkout's CRLF line endings.

Other baseline-relevant findings: the CLI command layer has no tests; several `detect` tests rename `testdata/repos/*/dotGit` ↔ `.git` in place and therefore are unsafe on a read-only mount; Windows disables symlink materialization and skips two tests; `../gitleaks-test-repo` is absent; and the locally cached `gitleaks:rnc-base` image was built from `ec7d126` with a dirty tree and `v8.30.0` ldflags, so it is **not** baseline evidence. Toolchain versions differ across the module (1.24.11), CI (1.25), Dockerfile (1.24), and the available container (1.27.3); results from these are not directly comparable.

CP1 can begin on the canonical source SHA without any executable Go diagnostics. Unlocking native or containerized `build`/`test`/`vet` requires an explicit dependency-acquisition authorization described in the closing proposal.

---

## Scope and methodology

**In scope:** repository identity and lineage; package/entry-point inventory; toolchain and dependency inventory; Docker environment assessment; safe-command matrix; non-mutating baseline diagnostics; sibling fixture-repository readiness; working-tree integrity; CP1 readiness.

**Method:** read-only Git and filesystem inspection; source reading at the assessed SHA; local Docker metadata and offline container probes. Every container run used `--network none` and `--pull never`; the repository was mounted `:ro`; caches were directed to container-local `/tmp`; one write probe verified the mount is read-only. No `make` target, native build, test, scanner run, benchmark, dependency download, image pull, or network request occurred. No #2019 or #1572 reproduction was attempted.

**Evidence labels:** **Verified** = directly observed in source/Git metadata; **Measured** = produced by an executed command; **Inferred** = reasoned from code but not executed; **Unverified** = not established. All source references are at canonical SHA `5f6ec56`; the source tree is identical at branch HEAD.

**Correction carried from pre-flight review:** documentation-only divergence since `5f6ec56` includes `AGENTS.md` and `_assessment/**` in addition to `README.md` and `_backlog/**`. The conclusion that application source is unchanged is preserved and re-verified below.

---

## Findings register

| ID | Component / claim | Status | Evidence | Impact / risk | Follow-up |
| --- | --- | --- | --- | --- | --- |
| CP0-001 | Canonical code baseline is `5f6ec56`; branch HEAD `a89c5f7` differs only in documentation | Verified | `git merge-base HEAD master` = `5f6ec56`; `git diff --name-only 5f6ec56..HEAD` returns only `.md`, `.json`, `.csv` paths, including `AGENTS.md` and `_assessment/**` | Future findings must cite `5f6ec56`; HEAD must not be used as a code baseline | Use `5f6ec56` in every CP report |
| CP0-002 | 60 documentation-only commits separate the canonical baseline from HEAD | Verified | `git rev-list --count 5f6ec56..HEAD` = 60; no non-documentation path in the diff | Branch history is not a source change | Cite canonical SHA plus HEAD separately |
| CP0-003 | No native Go toolchain is installed | Verified | `Get-Command go` → not found; `go version` → `CommandNotFoundException`; no Go under Program Files, WinGet, chocolatey, `%USERPROFILE%\go`, or `C:\Go`; no `GO*` environment variables | Native `build`/`test`/`vet`/`-race` impossible | Provision toolchain only with separate authorization |
| CP0-004 | Locally cached `golang:1.27.3` runs offline and reports `go1.27.3 linux/amd64` | Measured | `docker run --network none --pull never golang:1.27.3 go version` → `go version go1.27.3 linux/amd64`, exit 0 | A container Go path exists, but it is not the CI/module toolchain | Pin a digest and align version before measurements |
| CP0-005 | Dependencies cannot be resolved offline; `go build ./...` fails | Measured | Offline probe with `-mod=readonly`, `GOPROXY=off`, empty module cache, `:ro` repo mount → exit 1 after 3.68 s; 14 `module lookup disabled by GOPROXY=off` errors | `go build`, `go test`, `go vet` are blocked without network or a pre-populated cache | Record NOT RUN; see closing authorization proposal |
| CP0-006 | The container image contains no module cache | Measured | `docker run … golang:1.27.3 ls -A /go/pkg/mod` → `No such file or directory` | Every build starts cold; no offline build is possible today | Populate an external cache or accept NOT RUN |
| CP0-007 | Toolchain versions are inconsistent across the project | Verified | `go.mod:3` = `go 1.24.11`; [`.github/workflows/test.yml:23`](../../../.github/workflows/test.yml#L23) = `1.25`; [`Dockerfile:1`](../../../Dockerfile#L1) = `golang:1.24`; local image = `1.27.3` | Results from different toolchains are not directly comparable | Pin one version for CP2–CP4 evidence |
| CP0-008 | The Windows checkout applies CRLF to all 214 Go files | Measured | `git ls-files --eol -- '*.go'` → 214 `w/crlf`, 0 `w/lf`; `core.autocrlf=true`; no `.gitattributes` | In-place `gofmt -l` lists every file; file hashes/diffs must state LF vs CRLF | Normalize copies for formatting checks; record hashes with line-ending context |
| CP0-009 | All Go files parse and are gofmt-clean after LF normalization | Measured | Container probe copied 214 files, stripped CR, `gofmt -l /tmp/s` → no output, exit 0 | No syntax or formatting defect found; CRLF is environmental | Re-run with the pinned toolchain once available |
| CP0-010 | The CLI command layer (`package cmd`) has no tests | Verified | 19 test files: `config` (2), `detect` (5), `detect/codec` (1), `report` (7), `sources` (2), `cmd/generate` (2); none under `cmd/` | Flag parsing, command routing, report selection, and exit behavior are unverified by tests | CP2 must treat CLI behavior as a coverage gap |
| CP0-011 | Some `detect` tests rename fixture directories in place | Verified | [`detect/detect_test.go:1282-1283`](../../../detect/detect_test.go#L1282-L1283), [`:1368-1369`](../../../detect/detect_test.go#L1368-L1369), implementation [`:2432-2455`](../../../detect/detect_test.go#L2432-L2455) | A test run can mutate `testdata/repos/*/dotGit`; interruption can leave `.git` behind; read-only test mounts fail | Run tests against a disposable copy; verify status after |
| CP0-012 | Windows disables the only tracked symlink and skips two tests | Verified | `core.symlinks=false`; fixture is a 25-byte regular file with no `LinkType`; skips at [`detect/detect_test.go:852`](../../../detect/detect_test.go#L852) and [`:2127-2132`](../../../detect/detect_test.go#L2127-L2132) | Symlink scanning and one Git-ref scenario lack Windows coverage | Label skipped tests; do not use the fixture as #1572 evidence |
| CP0-013 | Expected sibling fixture repository is absent | Verified | `../gitleaks-test-repo` does not exist; adjacent `gitleaks-harness` has only a disabled dataset manifest | External functional/accuracy corpus unavailable | Report as limitation; do not substitute the harness |
| CP0-014 | The local `gitleaks:rnc-base` image is not baseline-pinned | Measured | Extracted binary: `go version -m` → `vcs.revision=ec7d12622164674ed2b9bcbeb25f0fdabdefd138`, `vcs.modified=true`, `go1.27.3`, `CGO_ENABLED=0`, ldflags `v8.30.0`, no `-tags`; `gitleaks version` → `v8.30.0` | Image predates the assessment and contains dirty-tree content; using it would misstate the baseline | Do not cite as baseline; rebuild from `5f6ec56` when authorized |
| CP0-015 | `ec7d126` is a Dockerfile-only descendant of the canonical SHA | Verified | `git log -1 ec7d126` = “build(docker): refresh Go and Alpine base images”; `git diff --stat 5f6ec56..ec7d126` = `Dockerfile` only; image binary reports `vcs.modified=true` | Explains image provenance while still disqualifying it as exact evidence | CP4 build-tooling topic |
| CP0-016 | Every build/test/format/profile Makefile entry point mutates, deletes, generates, or fetches | Verified | [`Makefile:4`](../../../Makefile#L4) `git fetch --tags`; [`:8-10`](../../../Makefile#L8-L10) coverage output; [`:12-16`](../../../Makefile#L12-L16) `go fmt` + `go generate` prerequisites; [`:21-23`](../../../Makefile#L21-L23) `go mod tidy`; [`:28-31`](../../../Makefile#L28-L31) deletion; [`:33-34`](../../../Makefile#L33-L34) build+profile. `lint` ([`:25-26`](../../../Makefile#L25-L26)) is the only non-mutating target, but `golangci-lint` is absent | No usable Makefile target is a safe diagnostic entry point | Continue to prohibit build/test targets; record `make lint` as unavailable |
| CP0-017 | Diagnostics and profiling default to the current directory | Verified | [`cmd/diagnostics.go:48-61`](../../../cmd/diagnostics.go#L48-L61); [`cmd/root.go:97`](../../../cmd/root.go#L97); [`scripts/profile.sh:24`](../../../scripts/profile.sh#L24) | Running without external output paths writes into the repository | Always pass external `--diagnostics-dir`; never run `profile.sh` in-repo |
| CP0-018 | Configuration and ignore files are source-controlled and auto-loaded | Verified | [`cmd/root.go:132-194`](../../../cmd/root.go#L132-L194) loads `--config`, `GITLEAKS_CONFIG`, `GITLEAKS_CONFIG_TOML`, `<source>/.gitleaks.toml`, embedded default; [`:300-321`](../../../cmd/root.go#L300-L321) loads `.gitleaksignore`; repository has `.gitleaks.toml` (allowlists) and `.gitleaksignore` (944 lines) | Rules, allowlists, and suppressions are part of scan behavior and a security boundary | CP4 topic; do not alter |
| CP0-019 | No application HTTP/network client exists | Verified | `rg`-equivalent search: only `net/http/pprof` in [`cmd/diagnostics.go:6-7`](../../../cmd/diagnostics.go#L6-L7); [`config/config.go:416`](../../../config/config.go#L416) `extendURL` is a TODO stub; Git access is via `exec.CommandContext` | Network exposure is confined to Git subprocesses and the local pprof server | No external scan performed |
| CP0-020 | Remote-URL resolution runs `git ls-remote --get-url` | Verified command / Inferred network effect | [`sources/git.go:485`](../../../sources/git.go#L485); Git documents `--get-url` as returning the URL without contacting the remote | Scanning a repo with remotes spawns a Git subprocess; no traffic was captured | Confirm offline behavior before CP2 CLI tests |
| CP0-021 | Two regex engines are selected by build tag | Verified | [`regexp/stdlib_regex.go:1`](../../../regexp/stdlib_regex.go#L1) `!gore2regex`; [`regexp/wasilibs_regex.go:1`](../../../regexp/wasilibs_regex.go#L1) `gore2regex`; [`.goreleaser.yml:18-19`](../../../.goreleaser.yml#L18-L19) releases with `gore2regex` | Default/CI builds use stdlib; release binaries use go-re2; behavior and resources may differ | CP2/CP4 must state build tags |
| CP0-022 | The embedded default config contains 222 rules and is generated | Verified | [`config/config.go:19-20`](../../../config/config.go#L19-L20) `go:embed gitleaks.toml`; [`config/gitleaks.toml:1`](../../../config/gitleaks.toml#L1) auto-generated header; 222 `[[rules]]`; generator at [`cmd/generate/config/main.go:18`](../../../cmd/generate/config/main.go#L18) | Rule inventory is deterministic but regeneration mutates a tracked file | Never run `go generate` on this branch |
| CP0-023 | Test fixtures are small and bounded; no unbounded-device fixture exists | Verified | `testdata` = 215 files / 92,552 bytes; largest archive 10,240-byte tar; `nested.tar.gz` 3,675 bytes; no `/dev/random` or `/dev/urandom` reference outside `_backlog` | Fixtures are safe for unit tests; they are not a performance corpus | CP3 needs synthetic corpora |

---

## Execution evidence

All container commands used `--network none` and `--pull never`; repository mounts were `:ro`; caches were container-local `/tmp`.

| Command or experiment | Input and conditions | Result / exit code | Duration or metric | Notes |
| --- | --- | --- | --- | --- |
| `git rev-parse --show-toplevel`, `rev-parse HEAD`, `status --short --branch` | Repository root | Root, `a89c5f765…`, clean branch | — | Before-state; repeated after |
| `git merge-base HEAD master`; `git rev-list --count 5f6ec56..HEAD`; `git diff --name-only 5f6ec56..HEAD` | Canonical SHA vs HEAD | Merge base `5f6ec56`; 60 commits; documentation-only paths | — | Corrected lineage evidence |
| `go version` (native) | Host PATH | `CommandNotFoundException` | — | Go not installed |
| `Get-Command go,make,gcc,golangci-lint,gotestsum,hyperfine` | Host | All not found | — | Tooling inventory |
| `docker version`; `docker info` | Local daemon | Client/server 29.9.2; linux/amd64; overlayfs; cgroup v2 | — | Docker usable |
| `docker image ls --digests golang` / `gitleaks` | Local image cache | `golang:1.27.3` (digest `sha256:e519d3b0…`); `gitleaks:rnc-base`, `pr2308mod`, `release-latest` | — | No pulls performed |
| `docker run --network none --pull never golang:1.27.3 go version` | Local image | `go version go1.27.3 linux/amd64`; exit 0 | — | Container toolchain verified |
| `docker run … golang:1.27.3 go env …` | Local image | `linux amd64 go1.27.3 local /go /go/pkg/mod /root/.cache/go-build 1` | — | `CGO_ENABLED=1` |
| `docker run … golang:1.27.3 ls -A /go/pkg/mod` | Local image | `No such file or directory`; exit 2 | — | No module cache |
| `docker run … -v <repo>:/src:ro … touch /src/.cp0-write-probe` | Read-only mount | `Read-only file system`; exit 1 | — | Host probe file absent |
| `docker run … -v <repo>:/src:ro -w /src -e GOFLAGS=-mod=readonly -e GOPROXY=off … go build ./...` | Canonical tree, empty cache | Exit 1; 14 `module lookup disabled by GOPROXY=off` errors | 3.68 s | Dependencies block build |
| `docker run … gofmt -l -e .` | CRLF working tree | 214 files listed; exit 0 | — | Environmental CRLF artifact |
| Container copy + CR-strip + `gofmt -l /tmp/s` | 214 normalized files | No output; exit 0 | — | Syntax/format clean |
| `docker run … gitleaks:rnc-base version` | Local image | `v8.30.0`; exit 0 | — | Version string only |
| `go version -m /tmp/gitleaks.bin` (extracted image binary) | Local image | `vcs.revision=ec7d126…`, `vcs.modified=true`, `go1.27.3`, `CGO_ENABLED=0`, no build tags | — | Provenance established, not baseline |
| `go test ./...` | Native and container | **NOT RUN** | — | Dependencies unresolved; fixture-rename hazard |
| `go vet ./...` | Native and container | **NOT RUN** | — | Same dependency barrier |

---

## Architectural and functional observations

### 1. Repository and lineage manifest

| Item | Value |
| --- | --- |
| Canonical code SHA | `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c` |
| Branch HEAD | `a89c5f7653ba1037f747df285312eda5831fbbe5` |
| Branch upstream | `origin/feat/issues-and-pr` (0 ahead / 0 behind at assessment) |
| Default branch | `origin/master` = `5f6ec56` |
| Commits between baseline and HEAD | 60, all documentation-only |
| Tracked files / commits / pack size | 484 / 1,353 / 6.54 MiB |
| Go files / test files / total Go lines (non-blank) | 214 / 19 / 19,482 |
| Remotes | `origin` RNC fork; `upstream` gitleaks upstream (not fetched) |

### 2. Package and entry-point inventory

| Package / path | Files / src lines | Role | Key symbols (source references) | Tests |
| --- | --- | --- | --- | --- |
| `main.go` | 1 / 19 | Process entry; signal handler | `main` → `cmd.Execute()` ([`main.go:11-19`](../../../main.go#L11-L19)) | none |
| `cmd/` | 8 / 1,047 | Cobra CLI, flags, config bootstrap, detector assembly, report selection, exit behavior | `rootCmd` ([`root.go:44`](../../../cmd/root.go#L44)); `Execute` ([`:224`](../../../cmd/root.go#L224)); `Detector` ([`:249`](../../../cmd/root.go#L249)); `runDetect` ([`detect.go:50`](../../../cmd/detect.go#L50)); `runDirectory` ([`directory.go:24`](../../../cmd/directory.go#L24)); `runGit` ([`git.go:29`](../../../cmd/git.go#L29)); `runProtect` ([`protect.go:29`](../../../cmd/protect.go#L29)); `runStdIn` ([`stdin.go:23`](../../../cmd/stdin.go#L23)); `runVersion` ([`version.go:20`](../../../cmd/version.go#L20)) | **none** |
| `config/` | 6 / 715 | Viper config translation, rule/allowlist model, embedded default config, extension chain | `ViperConfig` ([`config.go:32`](../../../config/config.go#L32)); `Config` ([`:87`](../../../config/config.go#L87)); `Translate` ([`:109`](../../../config/config.go#L109)); `DefaultConfig` embed ([`:19-20`](../../../config/config.go#L19-L20)); `extendDefault` ([`:375`](../../../config/config.go#L375)); `extendPath` ([`:395`](../../../config/config.go#L395)); `extendURL` TODO ([`:416`](../../../config/config.go#L416)); `Rule` ([`rule.go:12`](../../../config/rule.go#L12)); `Allowlist` ([`allowlist.go:29`](../../../config/allowlist.go#L29)) | 2 |
| `detect/` | 12 / 1,358 | Detection loop, rule evaluation, decoding orchestration, baseline suppression, reading, location/format helpers | `Detector` ([`detect.go:39`](../../../detect/detect.go#L39)); `NewDetector` ([`:115`](../../../detect/detect.go#L115)); `AddGitleaksIgnore` ([`:153`](../../../detect/detect.go#L153)); `DetectSource` ([`:208`](../../../detect/detect.go#L208)); `DetectContext` ([`:276`](../../../detect/detect.go#L276)); `detectRule` ([`:371`](../../../detect/detect.go#L371)); `AddBaseline` ([`baseline.go`](../../../detect/baseline.go)) | 5 |
| `detect/codec/` | 10 / 772 | Recursive encoding detection/decoding: percent, unicode, hex, base64, ASCII checks | `Decoder` ([`decoder.go:10`](../../../detect/codec/decoder.go#L10)); `Decode` ([`:23`](../../../detect/codec/decoder.go#L23)); `findEncodedSegments` ([`:43`](../../../detect/codec/decoder.go#L43)); registry ([`encodings.go:56`](../../../detect/codec/encodings.go#L56)) | 1 |
| `sources/` | 8 / 998 | Fragment sources: single file/stream, directory walk, Git history/diff | `Source` interface ([`source.go:12`](../../../sources/source.go#L12)); `File` ([`file.go:30`](../../../sources/file.go#L30)); `Files` ([`files.go:50`](../../../sources/files.go#L50)); `ScanTarget` ([`files.go:17`](../../../sources/files.go#L17)); `Git` ([`git.go:280`](../../../sources/git.go#L280)); `Fragments` ([`:299`](../../../sources/git.go#L299)) | 2 |
| `report/` | 15 / 521 | Finding model and output writers | `Reporter` ([`report.go:14`](../../../report/report.go#L14)); `CsvReporter`, `JsonReporter`, `JunitReporter`, `SarifReporter`, `TemplateReporter`; selection at [`cmd/root.go:380-397`](../../../cmd/root.go#L380-L397) | 7 |
| `regexp/` | 2 / 20 | Regex engine indirection by build tag | `Version`, `MustCompile` ([`stdlib_regex.go:1`](../../../regexp/stdlib_regex.go#L1), [`wasilibs_regex.go:1`](../../../regexp/wasilibs_regex.go#L1)) | none |
| `logging/` | 1 / 39 | zerolog configuration wrapper | package-level logger | none |
| `version/` | 1 / 4 | Build-time version variable (ldflags) | `Version` | none |
| `cmd/scm/` | 1 / 47 | SCM platform detection for finding links | platform parsing | none |
| `cmd/generate/` | 139 / 8,171 | Rule/config code generation; 131 rule definition files | `//go:generate` ([`cmd/generate/config/main.go:18`](../../../cmd/generate/config/main.go#L18)) | 2 |
| `testdata/` | 215 files / 92,552 B | Synthetic fixtures used by tests | repositories, archives, configs, reports, ignore file | n/a |

### 3. CLI surface (initial map)

| Command | Use | Flags defined in |
| --- | --- | --- |
| `gitleaks` (root) | Global flags: `--config`, `--exit-code` (default 1), `--report-path`, `--report-format`, `--report-template`, `--baseline-path`, `--log-level`, `--verbose`, `--no-color`, `--max-target-megabytes` (0), `--ignore-gitleaks-allow`, `--redact` (default 100), `--no-banner`, `--enable-rule`, `--gitleaks-ignore-path` (`.`), `--max-decode-depth` (5), `--max-archive-depth` (0), `--timeout` (0), `--diagnostics`, `--diagnostics-dir` | [`cmd/root.go:75-97`](../../../cmd/root.go#L75-L97) |
| `detect` | Legacy detector: `--no-git`, `--pipe`, `--follow-symlinks`, `--source`, `--log-opts`, `--platform` | [`cmd/detect.go:34-45`](../../../cmd/detect.go#L34-L45) |
| `dir` | Directory/file scan: `--follow-symlinks` | [`cmd/directory.go:13-20`](../../../cmd/directory.go#L13-L20) |
| `git` | Git history/diff scan: `--platform`, `--staged`, `--pre-commit`, `--log-opts` | [`cmd/git.go:15-25`](../../../cmd/git.go#L15-L25) |
| `protect` | Pre-commit protection: `--staged`, `--log-opts`, `--source` | [`cmd/protect.go:16-24`](../../../cmd/protect.go#L16-L24) |
| `stdin` | Read from standard input | [`cmd/stdin.go:14-19`](../../../cmd/stdin.go#L14-L19) |
| `version` | Print version | [`cmd/version.go:11-16`](../../../cmd/version.go#L11-L16) |

**Initial data-flow map (to be deepened in CP1):** `main` → `cmd.Execute` → config bootstrap (`--config` / env / `<source>/.gitleaks.toml` / embedded default) → `Config.Translate` → `Detector` assembly (flags, ignore files, baseline, reporter, diagnostics) → `DetectSource` → one of `Files`/`File`/`Git` yields `Fragment`s → `DetectContext` applies path/commit allowlists, keyword prefilter, per-rule `detectRule`, up to `MaxDecodeDepth` decode passes, archive traversal in the source layer → findings are filtered/redacted → `Reporter.Write` and exit code. This is a provisional map; CP1 must verify each boundary.

### 4. Dependencies and tooling

| Item | Value |
| --- | --- |
| Module path | `github.com/zricethezav/gitleaks/v8` |
| Module Go requirement | `go 1.24.11` ([`go.mod:3`](../../../go.mod#L3)) |
| Direct / indirect / total requires | 16 / 56 / 72 |
| `go.sum` | 423 lines; no `vendor/`; no module cache |
| CI toolchain | Go 1.25 ([`test.yml:23`](../../../.github/workflows/test.yml#L23)) |
| Dockerfile build image | `golang:1.24` ([`Dockerfile:1`](../../../Dockerfile#L1)); runtime `alpine:3.22` ([`:7`](../../../Dockerfile#L7)) |
| Local container toolchain | `go1.27.3 linux/amd64`, `CGO_ENABLED=1`, `GOTOOLCHAIN=local` |
| Native toolchain | absent |
| Release build tags | `gore2regex` ([`.goreleaser.yml:18-19`](../../../.goreleaser.yml#L18-L19)) |
| Default build engine | stdlib `regexp` (`!gore2regex`) |
| Lint/tooling | `golangci-lint` absent; config at [`.golangci.yaml`](../../../.golangci.yaml#L1); `gotestsum` absent |
| Profiling | in-binary pprof via `--diagnostics`; `hyperfine`/`yq` absent |

**Not comparable:** the module’s required 1.24.11, CI’s 1.25, Dockerfile’s 1.24, and the local image’s 1.27.3. Any build/test/benchmark result produced in this environment must name its toolchain and cannot be treated as interchangeable with CI results.

### 5. Docker environment assessment

- Docker client/server 29.9.2, Linux containers on WSL2, overlayfs, cgroup v2. Local images include `golang:1.27.3` and three `gitleaks` images.
- `golang:1.27.3` works with `--network none`; the repository can be mounted read-only (write probe correctly failed).
- The image has no module cache (`/go/pkg/mod` absent). An offline `go build ./...` fails at dependency resolution, so no build, test, or vet can be produced from this image alone.
- `gitleaks:rnc-base` was built from `ec7d126` (a Dockerfile-only descendant of the canonical SHA) with `vcs.modified=true`, Go 1.27.3, `CGO_ENABLED=0`, ldflags `v8.30.0`, and **no** `gore2regex` tag. Its `version` output is `v8.30.0`. It is a useful local smoke artifact but **not** a reproducible baseline.
- Docker reports 23.92 GB of build cache. It may contain module layers from prior builds, but cache contents are unverified and cache-based resolution is not reproducible; it was not used.

### 6. Safe command matrix

| Operation | Command (exact form) | Status | Reason / controls |
| --- | --- | --- | --- |
| Git inspection | `git status/log/diff/show/ls-files/rev-list/grep/merge-base` | **ALLOWED — executed** | Read-only; no fetch |
| Filesystem/source inspection | PowerShell `Get-ChildItem`, `Select-String`, read | **ALLOWED — executed** | Read-only |
| Docker metadata | `docker version/info/image ls/inspect/history` | **ALLOWED — executed** | Local daemon; no pull |
| Offline container shell | `docker run --rm --network none --pull never golang:1.27.3 <cmd>` | **ALLOWED — executed** | No network; ephemeral `/tmp`; `:ro` repo |
| Go syntax/format check | `gofmt -l` on LF-normalized copy | **ALLOWED — executed** | No repo writes; no dependencies |
| Go build | `go build -mod=readonly ./...` | **BLOCKED** | Dependencies unresolved offline; network not authorized |
| Go test | `go test -mod=readonly ./...` | **NOT RUN** | Same dependency barrier; also in-place fixture renames |
| Go race test | `go test -mod=readonly -race ./...` | **NOT RUN** | Native: no Go and no C compiler; container: dependencies blocked |
| Go vet | `go vet -mod=readonly ./...` | **NOT RUN** | Same dependency barrier |
| Benchmarks | `go test -mod=readonly -run '^$' -bench … -cpuprofile …` | **NOT RUN** | Requires build/test capability; profiles must be external |
| `make test/build/format/profile/clean` | `make <target>` | **FORBIDDEN** | `format`, `go generate`, `go mod tidy`, `git fetch --tags`, deletions (see CP0-016); `make` absent |
| `make lint` | `make lint` | **NOT RUN** | Target is non-mutating, but `golangci-lint` is absent and `make` itself is not installed |
| Code generation | `go generate ./...` | **FORBIDDEN** | Rewrites tracked `config/gitleaks.toml` |
| Formatting mutation | `go fmt ./...` | **FORBIDDEN** | Rewrites all source files |
| Dependency mutation | `go mod tidy` | **FORBIDDEN** | Rewrites `go.mod`/`go.sum` |
| Profiling script | `./scripts/profile.sh …` | **FORBIDDEN** | Writes `profile/` under CWD; runs scanner; `set -euo pipefail` with default exit code 1 |
| Pre-commit helper | `python scripts/pre-commit.py` | **NOT AUTHORIZED** | Requires installed `gitleaks`; runs `protect --staged` |
| Image build / pull | `docker build`, `docker pull` | **NOT AUTHORIZED** | Network and image-state changes |
| Git fetch/tags | `git fetch --tags` | **FORBIDDEN** | Network and `.git` mutation |
| Release build | `goreleaser` | **NOT AUTHORIZED** | Writes `dist/`; network/release semantics |
| Scanner runs on live/remote data | `gitleaks git <external repo>` | **FORBIDDEN** | External systems and untrusted input |
| #2019 / #1572 reproductions | any unbounded memory/device scenario | **EXCLUDED** | Explicitly out of CP0; unsafe |

**Commands that may establish external connections:** `git fetch`/`ls-remote` against remotes, `go build/test/vet` without a populated cache (module proxy), `docker pull/build`, `goreleaser`, and the GitHub workflows. **Commands that can modify source/config/deps/generated files:** `make test`, `make build`, `make format`, `make clean`, `make profile`, `go generate`, `go fmt`, `go mod tidy`, `make test-cover` (writes `cover.out`), `scripts/profile.sh` (writes `profile/`), and `go test ./detect` (temporarily renames `testdata/repos/*/dotGit`).

### 7. Fixture readiness

- **Expected external corpus:** `../gitleaks-test-repo` is absent. This is recorded as an unavailable check, not a successful test.
- **Adjacent resources (outside this assessment’s authorization):** `C:\Dev\gitleaks\gitleaks-harness` (orphan `test/scan-harness` branch; only `datasets/manifest.yaml`, Levain corpus disabled with `revision: null`), `pr2308`, `pr2336`, and a dirty `review` working copy of the security-integration branch. None were used as evidence.
- **Internal fixtures:** 215 files / 92,552 bytes, dominated by 136 files in `testdata/repos` (Git fixtures stored as `dotGit` and renamed by tests), 50 config files, and 13 files under `testdata/archives`. Archives are bounded (largest 10,240 bytes); `nested.tar.gz` is 3,675 bytes. No `/dev/random`, `/dev/urandom`, or unbounded stream reference exists in source/scripts/config outside `_backlog`.
- **Hazards:** `detect` tests rename fixture directories in place ([`detect_test.go:2432-2455`](../../../detect/detect_test.go#L2432-L2455)); the tracked symlink is a 25-byte text file under this checkout; two tests skip on Windows. Any test run must use a disposable copy or be followed by a status check and fixture-integrity verification.

### 8. Working-tree integrity

- **Before diagnostics:** `## feat/issues-and-pr...origin/feat/issues-and-pr`; no staged, unstaged, untracked, or ignored changes; HEAD `a89c5f7653ba1037f747df285312eda5831fbbe5`.
- **After diagnostics:** identical output; HEAD unchanged.
- **Protected-file changes:** none. The container repository mount was read-only; the write probe failed with `Read-only file system`; no probe file appeared on the host. Temporary binaries/scripts were removed from `%TEMP%`. No container or volume created by this checkpoint remains (`cp0-extract` removed).

### 9. Known limitations and unresolved questions

1. **No executable Go baseline.** Dependencies are unresolved and network acquisition is unauthorized. `build`, `test`, `vet`, and benchmarks are NOT RUN. This does not block CP1’s source-level architecture work, but it blocks CP2/CP3 runtime evidence until resolved.
2. **Toolchain choice unresolved.** Should measurements use Go 1.24.11 (module), 1.25 (CI), or another pinned version? Results must not mix toolchains.
3. **Cross-platform limits.** Windows checkout has CRLF, no symlinks, two skipped tests, and no `-race` capability without a C compiler. Container Linux results will not describe native Windows behavior.
4. **Fixture-mutating tests.** The `detect` suite writes to `testdata/repos/*`; a read-only mount cannot run it, and an interrupted run can leave the tree changed. A disposable copy at the pinned SHA is the safe route.
5. **No external corpus.** `../gitleaks-test-repo` is absent and the harness corpus is disabled; performance and accuracy work will need separately authorized synthetic corpora.
6. **Local image provenance.** `gitleaks:rnc-base` is `v8.30.0` from `ec7d126` with `vcs.modified=true`; it must not be cited as the baseline.
7. **Git/CI parity unknown.** Git 2.45.2 locally versus CI runner versions; no fetch was performed, so upstream refs and the current upstream backlog state are unverified.

### 10. CP1 readiness

**Recommendation: READY for CP1, with the limitations above.** CP1 is a source-linked architecture investigation and requires no executable Go baseline. It should cite `5f6ec56`, not HEAD, and should treat the CLI/data-flow map above as provisional. The dependency-network decision affects CP2/CP3, not CP1.

---

## Definition-of-done review and handoff

- [x] Required deliverables complete or omissions explained.
- [x] Significant claims linked to evidence and labeled Verified / Measured / Inferred / Unverified.
- [x] Failures, unavailable checks, and unknowns recorded; no unrun test described as passing.
- [x] Source-read-only boundary respected; before/after status recorded; no protected-file changes.
- [x] Readiness and risks for the next CP documented.
- [x] Stopped pending explicit owner authorization; CP1 not started.

**Recommendation:** Ready for owner review.

---

## Proposed additional authorization (not executed)

To make the executable baseline real without touching the repository, the following three-step operation would materially improve CP0/CP2/CP3 evidence:

**Operation:** one-time dependency hydration plus a pinned, CI-matched Go toolchain.
1. Pull `golang:1.25` (or `golang:1.24.11`) and pin the resulting image by digest.
2. Populate an external module cache in a named Docker volume:
   `docker run --rm --network bridge -v <repo>:/src:ro -v gitleaks-cp0-gomod:/gomod -w /src -e GOFLAGS=-mod=readonly -e GOMODCACHE=/gomod golang:<pinned> go mod download`
3. Run diagnostics offline against a **writable disposable copy** of the canonical tree (the `detect` tests rename fixtures, so the copy must be writable; it is discarded afterward):
   `git archive --format=tar -o <temp>/cp0-tree.tar 5f6ec56`, then
   `tar -xf <temp>/cp0-tree.tar -C <temp>/tree`, then
   `docker run --rm --network none -v <temp>/tree:/src -v gitleaks-cp0-gomod:/gomod -w /src -e GOFLAGS=-mod=readonly -e GOPROXY=off -e GOMODCACHE=/gomod -e CGO_ENABLED=1 golang:<pinned> go test ./...`
   Add `-race` to the final `go test` invocation when race coverage is required; the image includes `gcc`.

**Benefits:** enables `go build`, `go test`, `go vet`, `-race` (image includes gcc), coverage, and scoped benchmarks on the exact canonical tree; produces a defensible executable baseline.

**Controls:** digest-pinned image; network only during the `go mod download` step; `go.sum` checksum verification; repository never writable (tests run on an extracted copy); caches/profiles outside the repository; `-mod=readonly`; status and fixture integrity checked before/after; no credentials or external scanner runs; #2019/#1572 remain excluded.

**Alternative:** install native Go plus MinGW/gcc on the host. This is more invasive, changes the host toolchain, and still needs module downloads; it is not recommended for reproducibility.

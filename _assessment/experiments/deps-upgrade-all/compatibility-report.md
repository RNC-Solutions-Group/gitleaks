# Compatibility report — dependency and RE2 upgrade

**Status:** Experimental; Go 1.27.3 local evidence, no CI/release execution
**Branch:** `test/deps-upgrade-all`
**Canonical source:** `5f6ec56cfa6f8c67648d863582c7ae2f206d5c1c`
**Baseline documentation HEAD:** `afa2c3c8280b7dd581c6bf3b2bc4c7d7b224a20b`

## 1. Method and source control

Both stacks were exported as LF archives from the same Git object database:

- Baseline: `git archive` of `afa2c3c` (tree equivalent to canonical `5f6ec56` for application code).
- Upgraded: `git archive` of `test/deps-upgrade-all` at `70aaf01` (dependency commit).
- Archives preserve Git symlink entries and normalize fixture line endings; extracting the Windows working tree directly is not valid for the test fixtures (`core.symlinks=false` produced regular files and false failures).
- Baseline archive SHA-256: `5cc7ff685263b83dfcee3421b178e7fc5dd7bc8eec072f6b5a3a77ce2b6a871c`.
- Upgraded archive SHA-256: `2d3cd225f47cc94a39d9acfd5528f4391e0c2733b759e8a4335b705ebc7925c1`.

Primary container: `golang:1.27.3` linux/amd64, `GOTOOLCHAIN=local`, `GOFLAGS=-mod=readonly`, `GOMODCACHE=/gomod`, `GOCACHE=/gocache`. Tests used `CGO_ENABLED=1`; application binaries used `CGO_ENABLED=0`. Toolchain: Go 1.27.3.

## 2. Build and test matrix

| Check | Baseline | Upgraded | Result |
| --- | --- | --- | --- |
| `go mod verify` | Pass | Pass | all modules verified |
| `go build ./...` | Pass | Pass | no compile errors |
| `go test ./... -count=1` | Pass, 275/275 | Pass, 275/275 | 0 fail, 0 skip |
| `go test -race ./... -count=1` | Pass | Pass, 275/275 | exercised paths |
| `go test -shuffle=on ./...` | Pass (random seed) | Fail (seed `1794264498871797264`) | order-sensitive; see §3 |
| `go vet ./...` | Fail | Fail | same two `lostcancel` diagnostics |
| Standard app build (`CGO_ENABLED=0`) | Pass | Pass | `gitleaks` |
| `gore2regex` app build (`CGO_ENABLED=0`) | Pass | Pass | `gitleaks-re2` |
| `go build -tags gore2regex ./...` | Fail | Fail | same generator-rule type errors |
| `go test -tags gore2regex ./...` | Fail | Fail | same generator-rule type errors |
| RE2-tagged targeted packages | Pass | Pass | `detect`, `detect/codec`, `report`, `sources`, `config` |
| CLI config/scan smoke | Pass | Pass | one synthetic finding, both engines |
| Report formats + redaction | Pass | Pass | JSON/CSV/JUnit/SARIF/template valid |
| Config precedence | Pass | Pass | source-local, env, explicit all selected correctly |
| Windows native build | Pass | Pass | stdlib and RE2 |
| Windows native CLI parity | Pass | Pass | 40/40 normalized findings |
| Go 1.25 build | — | **Fail** | `go.mod requires go >= 1.26.0` |
| Go 1.24 build | — | **Fail** | `go.mod requires go >= 1.26.0` |

### Exact pre-existing static/build failures

`go vet ./...`:

```text
detect/detect_test.go:2083:4: the cancel function is not used on all paths (possible context leak)
detect/detect_test.go:2122:3: this return statement may be reached without using the cancel var defined on line 2083
```

Tagged full graph (`go build -tags gore2regex ./...` and `go test -tags gore2regex ./...`):

```text
cmd/generate/config/rules/airtable.go:30:16: cannot use regexp.MustCompile(...) (value of type *"regexp".Regexp)
  as *"github.com/zricethezav/gitleaks/v8/regexp".Regexp value in struct literal
cmd/generate/config/rules/perplexity.go:16:16: cannot use regexp.MustCompile(...) (value of type *"regexp".Regexp)
  as *"github.com/zricethezav/gitleaks/v8/regexp".Regexp value in struct literal
```

These are **pre-existing** baseline failures (CP5-F14), not introduced by the upgrade. The tagged application binary is viable because it does not compile the generator packages.

### RE2-tagged targeted tests

Because the full tagged graph is blocked by the generator packages, the packages that do compile under `gore2regex` were tested directly:

```text
go test -tags gore2regex ./detect ./detect/codec ./report ./sources ./config -count=1
```

Both baseline and upgraded stacks pass. Upgraded durations: `detect` 4.457 s, `detect/codec` 0.036 s, `report` 0.018 s, `sources` 0.026 s, `config` 0.354 s. Baseline durations were longer (for example `detect` 8.625 s), consistent with the RE2 v1.13 improvement; those are not a controlled benchmark because package scheduling and cgroup `GOMAXPROCS` differ.

## 3. Shuffle-order finding

The upgraded random `-shuffle=on` run failed:

```text
--- FAIL: TestTranslateExtend/valid/extend
    config_test.go:574: valid/extend diff: (-want +got)
```

Both stacks were then run with the same seed:

```text
go test -shuffle=1794264498871797264 ./config -count=1 -v
```

Both baseline and upgraded stacks fail identically. This is the pre-existing `extendDepth`/global-config test-order sensitivity described in CP5-F07. It is **not** a dependency regression, but it means a single passing `-shuffle=on` run is not a reliable claim of order independence.

## 4. Go toolchain compatibility

`go-re2` v1.13.0 declares `go 1.26.0`, so `go get` raised the module directive from `go 1.24.11` to `go 1.26.0`. No `toolchain` directive was added.

| Toolchain | Result | Exact failure |
| --- | --- | --- |
| Go 1.27.3 | Pass | all local build/test checks |
| Go 1.25.14 (`golang:1.25`, CI version) | **Fail** | `go: go.mod requires go >= 1.26.0 (running go 1.25.14; GOTOOLCHAIN=local)` |
| Go 1.24.13 (`golang:1.24`, Dockerfile version) | **Fail** | `go: go.mod requires go >= 1.26.0 (running go 1.24.13; GOTOOLCHAIN=local)` |

This is the primary migration blocker. It requires a CI/Dockerfile toolchain change, not an application-code change. The previous CP5 CI gap was “Go 1.25 not run”; this experiment establishes that Go 1.25 cannot run the upgraded stack at all.

## 5. RE2 implementation details

### Default WebAssembly-based RE2

- `github.com/wasilibs/go-re2` `v1.9.0` → `v1.13.0`.
- v1.13.0 uses `wasm2go`; `wazero` and `wazero-helpers` are removed from the module graph.
- The Gitleaks wrapper still compiles against the same `Compile`/`MustCompile`/`Regexp` surface; no source change was required.
- The application binary builds and runs; default config loading and rules work.
- Full tagged graph still fails only because generator packages import stdlib `regexp` directly.

### Optional native `re2_cgo`

The native variant is still gated by the `re2_cgo` build tag in addition to `gore2regex`. `CGO_ENABLED=1` alone does not select it. With the stock `golang:1.27.3` image, the build fails at system dependency resolution:

```text
# github.com/wasilibs/go-re2/internal/cre2
# [pkg-config --cflags -- re2]
Package re2 was not found in the pkg-config search path.
```

Requirement: install `re2` development headers/pkg-config metadata (for Debian/Ubuntu, `libre2-dev`) and set `PKG_CONFIG_PATH` if necessary. This variant was **not** built or tested further; it is not part of the default WebAssembly-based result.

## 6. Windows native checks

Windows/amd64, Go 1.27.3, `CGO_ENABLED=0`:

- Standard application build: Pass.
- `gore2regex` application build: Pass.
- `go build ./...`: Pass.
- CLI detection on the synthetic corpus: stdlib and RE2, baseline and upgraded, all produce **40/40** normalized findings; encoded depth-2 and archive depth-1 findings are present in all four.
- Normalized comparison across baseline/upgraded and stdlib/RE2: equivalent.

Limitation: the full Windows test suite was **not run**. The test fixtures include Git symlink entries that a default Windows checkout materializes as regular files; that environment issue would produce false failures unrelated to dependencies. Linux tests were run from an archive with symlinks preserved.

## 7. Configuration, CLI, and reporting surface

These checks were run with the upgraded application binary, standard and RE2 engines where relevant:

- Custom config with rule allowlist, path allowlist, entropy threshold, Unicode `\p{L}`, generic/specific duplicate rules: loaded and behaved as expected.
- Embedded default config: loaded; three expected default-rule findings on the synthetic corpus (AWS access token, GitHub PAT, private key).
- Config precedence:
  - source-local `.gitleaks.toml`: selected `local` rule only;
  - `GITLEAKS_CONFIG`: selected `env` rule only;
  - explicit `--config` overrides `GITLEAKS_CONFIG`: selected `explicit` rule only.
- Reports: JSON/CSV/JUnit/SARIF/template all produced valid output; CSV had a header; SARIF was valid JSON; JUnit parsed as XML; template rendered.
- `--redact`: 39/39 JSON findings had `Secret: "REDACTED"` and no raw synthetic token appeared.
- `--max-decode-depth 2` found both one-level and two-level Base64 secrets and an archive-contained secret; depth 0 omitted the encoded findings while retaining the archive finding. This confirms decoding/archive behavior was exercised rather than bypassed.

## 8. Conclusion

The upgraded stack is compatible with the Go 1.27.3 application source and test suite, with no source changes. Its only **introduced mandatory compatibility change** is the **Go 1.26 minimum**; the existing `go vet`, tagged full-graph, and shuffle-order failures remain. The optional native `re2_cgo` path has an additional system-library prerequisite and is not covered by this experiment.

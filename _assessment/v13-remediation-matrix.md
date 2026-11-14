# CP-V13 Remediation Matrix

**Starting baseline:** `3a4635f82dcd0a216f9d2b88ebbd72341b3b60cc` (merged V12)
**Development branch:** `fix/v13-operational-hardening`
**Development HEAD:** final V13 development commit (see branch tip)
**Integration:** `audit/v13-integration` → PR to `rnc-mod`
**Scope:** 20 designated operational, reporting, configuration, performance, stability and release-engineering packages. No detection-rule changes and no benchmark-harness changes.

## Status summary

| ID | Package | Disposition | Evidence |
|---|---|---|---|
| V13-01 | Unused `.gitleaksignore` diagnostics | **FIXED — LOCALLY VALIDATED** | Usage tracking + opt-in completed-scan stderr diagnostic; no machine-schema change; detector/CLI tests. |
| V13-02 | Effective configuration inspection | **FIXED — LOCALLY VALIDATED** | `gitleaks config` safe summary; precedence/local/invalid tests; no regex export. |
| V13-03 | Multiple report formats | **DEFERRED — OUTSIDE AUTHORIZED SCOPE** | Safe report-plan architecture documented; no partial unsafe implementation. |
| V13-04 | Bounded finding context | **DEFERRED — OUTSIDE AUTHORIZED SCOPE** | Separate opt-in context schema/redaction design documented; no schema extension. |
| V13-05 | SARIF result severity | **FIXED — LOCALLY VALIDATED** | Opt-in constant `--sarif-level`; default omission preserved; report/CLI tests. |
| V13-06 | Repeated-finding grouping | **POLICY DECISION REQUIRED** | Unkeyed secret SHA rejected; keyed-HMAC grouping design only. |
| V13-07 | Console location / decoded UX | **PARTIALLY CORRECTED — RESIDUAL CONFIRMED** | Opt-in local file links implemented; decoded highlighting remains separate. |
| V13-08 | Git `--log-opts` parsing | **FIXED — LOCALLY VALIDATED** | Shell-free tokenizer, parsed-vector validation, quote/escape/Windows/pathspec tests. |
| V13-09 | Shallow-clone provenance | **PARTIALLY CORRECTED — RESIDUAL CONFIRMED** | Shallow warning added; fingerprints/SHAs unchanged; missing history cannot be reconstructed. |
| V13-10 | Git source scope / bare repos | **PARTIALLY CORRECTED — RESIDUAL CONFIRMED** | Subdirectory and bare-repository warnings; no default scope change; explicit config/ignore path. |
| V13-11 | Git remote normalization / links | **FIXED — LOCALLY VALIDATED** | Azure v3 SSH mapping, credential stripping, URL path escaping; remote/link tests. |
| V13-12 | Configuration validation/diagnostics | **PARTIALLY CORRECTED — RESIDUAL CONFIRMED** | Empty/malformed blocks fail closed; unknown fields still ignored; safe summary. |
| V13-13 | Pre-commit environment reliability | **NOT REPRODUCED — EVIDENCE INSUFFICIENT** | Local real pre-commit passes; Docker/macOS/PTY evidence unavailable or not a shared root cause. |
| V13-14 | Profiling/benchmark tooling | **FIXED — LOCALLY VALIDATED** | Command parity, quoting, exit-code handling, metadata, artifact checks; Linux dir/git CPU/mem/trace smoke run. |
| V13-15 | Recursive decoding / hang stability | **NOT REPRODUCED — EVIDENCE INSUFFICIENT** | 6.67 MB line-bounded decode probe at depths 0/1/3/5: peak RSS ~18.0–18.8 MB, no depth-correlated growth; hangs not reproduced. Defaults retained. |
| V13-16 | Explicit-file scalability | **FIXED — LOCALLY VALIDATED** | Bounded concurrency, deterministic sort, dedupe/error/cancellation invariants; 200×50 KB median ~0.577 s → ~0.266 s. |
| V13-17 | RNC release workflow / image identity | **FIXED — LOCALLY VALIDATED** | GHCR-only, RNC image, tag/base/ancestry checks, least privilege; static workflow/tag probes. No publication. |
| V13-18 | Versioning / install policy | **POLICY DECISION REQUIRED** | `vX.Y.Z-rnc.N`, base strictly > minVersion, container tags and module-path consequences documented; no module rename/release. |
| V13-19 | Release-binary portability | **PARTIALLY CORRECTED — RESIDUAL CONFIRMED** | Artifact/RE2/architecture plan documented; goreleaser unavailable and no artifact-level test. |
| V13-20 | Supply-chain attestations | **DEFERRED — OUTSIDE AUTHORIZED SCOPE** | Digest/SBOM/provenance/checksum/signature verification design and checklist; nothing claimed or published. |

## Important implementation notes

### V13-01

`Detector` records ignore entries that actually suppress a finding. `--report-unused-ignore` prints sorted entries after a **completed** scan only. It explicitly says “not triggered” is not proof of obsolescence because a path/commit/revision may not have been examined. The diagnostic is stderr-only and opt-in; no JSON schema field was added. Tests cover used/unused entries and CLI behavior.

### V13-02

`gitleaks config [path]` resolves the same configuration precedence/extension path as a scan and prints source, loaded files, active/disabled/ordered rule counts, global/rule allowlist counts, minVersion/title, report policy and limits. It cannot print regexes/patterns. Tests cover explicit config, target-local discovery, hidden-pattern safety and invalid configuration.

### V13-03 / V13-04 / V13-06

Designs are recorded in `_assessment/v13-reporting-options.md`:

- Multi-format requires a report-plan abstraction, all-destination validation and explicit non-atomic partial-failure semantics.
- Context requires a separate opt-in schema/mode, byte/line budgets and redaction of every known sensitive span in the window.
- Grouping requires a keyed HMAC and key lifecycle; an unkeyed SHA-256 of a secret is unsafe and is rejected.

No unsafe approximation was shipped.

### V13-05

`--sarif-level` accepts `none`, `note`, `warning`, `error`, requires SARIF output and leaves `level` absent by default. No per-rule severity classification is inferred.

### V13-07

`--file-links` adds a console-only `File-Link:` URI with `#Lline` for non-Git findings. `--no-color` suppresses OSC 8 escapes; machine reports are unchanged. Decoded-value highlighting remains a residual not represented as fixed.

### V13-08

The parser supports single/double quotes, escaped whitespace/quotes, literal Windows backslashes and pathspec `--` boundaries; unmatched quotes, empty arguments and unsafe formats fail. `validateGitLogArgs` operates on the exact vector later executed, so quoting cannot bypass V12 fail-closed restrictions. No shell is used.

### V13-09/V13-10

Warnings are emitted for shallow repositories, bare repositories and repository-wide history from a subdirectory. They do not rewrite commit IDs, paths, fingerprints or default history scope. A bare repository must receive explicit `--config`/`--gitleaks-ignore-path` when those inputs are needed.

### V13-11

SCP remotes normalize to HTTPS; Azure `ssh.dev.azure.com:v3/org/project/repo` becomes `https://dev.azure.com/org/project/_git/repo`; userinfo/tokens are stripped; link paths use `url.URL` escaping. Existing GitHub/GitLab/Gitea/Bitbucket link fixtures still pass.

### V13-12

Empty `[[rules]]` and malformed TOML fail closed with contextual errors. Unknown fields remain accepted/ignored because strict rejection would break existing configurations; this is documented rather than changed abruptly. The safe config summary gives operators an effective-view diagnostic.

### V13-13

The local real pre-commit runner passes for the staged and explicit-file Go/system hooks. Docker/macOS/PTY reports were not reproduced and are not assigned an invented cause. Hook scope and `safe.directory` policy were not broadened.

### V13-14

`scripts/profile.sh` now:

- uses one command shape for hyperfine and profiles;
- defaults to decode depth 5, overridable explicitly;
- uses `--exit-code=0` so findings do not abort, while operational failures still abort;
- quotes paths and arrays;
- records UTC metadata, binary/toolchain/Git state;
- writes CPU/memory/trace artifacts and checks each exists;
- documents optional hyperfine/yq dependencies.

Linux smoke run produced all six profile artifacts for `dir` and `git`.

### V13-15

A finite 6.67 MB encoded, line-bounded corpus was measured with GNU time in a Linux container (empty rule set, depths 0/1/3/5): maximum RSS ~18.0, 18.8, 18.6, 18.5 MB; wall time ~0.11–0.15 s. No depth-correlated growth or hang reproduced. This is not a universal memory bound; encoded-input amplification remains a residual.

### V13-16

Explicit-file scanning uses the detector semaphore (bounded at 40) while retaining exact-file scope, dedupe, symlink policy, close/read errors and cancellation. Findings are sorted by file, line, column and rule. 200 synthetic 50 KB files on Windows: sequential V12 median ~0.577 s; V13 median ~0.266 s (4 timed runs, warm-up discarded). This is a local synthetic measurement, not a universal throughput claim.

### V13-17/V13-18/V13-19/V13-20

Release policy is documented in `RELEASING.md`:

- GHCR-only `ghcr.io/rnc-solutions-group/gitleaks`; no Docker Hub.
- `vMAJOR.MINOR.PATCH-rnc.N`; numeric base strictly newer than `minVersion`; `v8.25.0-rnc.1` is explicitly invalid for that reason.
- Release commit must be reachable from protected `rnc-mod`.
- Container tags use semver version/minor/major, no `latest`.
- Module path remains upstream; direct `go install` of the RNC repo path is not supported.
- SBOM/provenance/checksum/signature verification is a future release gate, not a current claim.

No tag, image, package, release, SBOM or signature was published by V13.

## Validation gates

Executed locally on the V13 development tree:

```text
go build ./...
go test -count=1 ./...
go test -count=1 -tags gore2regex ./...
go vet ./...
go mod verify
go generate ./...
git diff --check
```

Additional evidence:

- Real pre-commit integration: PASS.
- Linux non-root race: to be run on the frozen integration candidate.
- Docker build/smoke: local/CI dependent; no V13 publication.
- Release workflow parse/scope and tag-policy probes: PASS in source.
- Profile script Linux smoke: PASS.
- Explicit-file benchmark: measured, synthetic.
- Decode memory probe: finite, Linux container.

## Residuals

- Multi-format report plan, bounded context, keyed grouping: not enabled.
- Decoded-finding highlighting: not corrected.
- Shallow history cannot become complete; bare/subdirectory scope remains explicit/default-safe.
- Unknown configuration fields and allowlist-placement diagnostics remain partial.
- Docker/macOS/PTY pre-commit failures remain unreproduced.
- Recursive decoding/hang amplification remains a resource residual.
- GoReleaser/artifact portability and attestations are release-gated; nothing published.
- No independent scanner-harness certification.

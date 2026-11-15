# CP-V15 — Engineering evidence appendix

## 1. Environment

| Item | Value |
|---|---|
| Assessment date | 2026-11-15, America/New_York |
| Product worktree | `C:\Dev\gitleaks\gitleaks-upgrade` |
| Product HEAD | `7d74921841114be7677a423cf60620210e832bf7` |
| Product tree | `8eb327dd5860a4665daf38ce2f2fbc0aa5c07963` |
| Product branch | `rnc-mod`, clean and synchronized with `origin/rnc-mod` |
| Backlog worktree | `C:\Dev\gitleaks\gitleaks-assessment` |
| Backlog starting HEAD | `7f7e4eada3bdf3b594a3c0ec1fd101b89f2f08e2` |
| OS | Windows 11, NTFS |
| Go | `go1.27.3 windows/amd64` |
| Git | `2.45.2.windows.1` |
| Python | 3.13.5 (analysis/scripts only) |
| Docker | Not run locally; post-merge GitHub Actions Docker jobs passed |
| Independent harness/oracle | Not present in the repository or local workspaces; not run |

## 2. Baseline/ancestry evidence

Commands:

```text
git status --short --branch
git rev-parse HEAD
git rev-parse '71be6f3^{tree}'
git rev-parse '1e88ac8^{tree}'
git rev-parse 'cb7b058^{tree}'
git rev-parse '7d74921^{tree}'
git diff --quiet 7d74921 cb7b058
git diff --quiet 71be6f3 1e88ac8
```

Observed:

```text
71be6f3 tree = c1686124229b7751be58851211cb379341d45e27
1e88ac8 tree = c1686124229b7751be58851211cb379341d45e27
cb7b058 tree = 8eb327dd5860a4665daf38ce2f2fbc0aa5c07963
7d74921 tree = 8eb327dd5860a4665daf38ce2f2fbc0aa5c07963
tree diff exit = 0
```

`1e88ac8` parent is `3a4635f`; `7d74921` parent is `1e88ac8`. The squash merges did
not alter the reviewed product trees.

CI evidence:

```text
gh api repos/RNC-Solutions-Group/gitleaks/commits/7d74921.../check-runs
```

Result: 14 check runs — 13 `success`, one `skipped` (`gitleaks` action), zero failures.
Run `43763741641`; last job completed 2026-11-15T21:41:01Z.

## 3. Local validation results

| Command | Result | Notes |
|---|---|---|
| `go build ./...` | PASS | No output |
| `go vet ./...` | PASS | No output |
| `go test -count=1 ./...` | PASS | `cmd` 89.498 s; exit 0 |
| `go test -count=1 -tags gore2regex ./...` | PASS | `cmd` 89.010 s; exit 0 |
| `go test -count=1 -coverprofile=... ./...` | PASS | 61.8% statements total; `cmd` 98.214 s |
| `go generate ./...` | NOT RE-RUN in V15 | V14 validated and restored generated line endings; no V15 source edit |
| `go test -race ./...` | NOT RUN locally | CI Linux/Windows/macOS race jobs passed; no local C toolchain claim |
| `staticcheck` / `golangci-lint` | NOT RUN | No installed tool; `go vet` only |
| `deadcode` | NOT RUN | No dead-code claim |
| Docker build/smoke | NOT RUN locally | Post-merge CI amd64/arm64 root/non-root jobs passed |

Package coverage:

| Package | Coverage |
|---|---:|
| `cmd` | 69.9% |
| `cmd/generate/config` | 67.7% |
| `cmd/generate/config/base` | 100.0% |
| `cmd/generate/config/rules` | 2.2% |
| `cmd/generate/config/utils` | 44.8% |
| `cmd/generate/secrets` | 77.8% |
| `config` | 64.2% |
| `detect` | 84.1% |
| `detect/codec` | 66.2% |
| `internal/informative` | 86.8% |
| `internal/scanerr` | 95.7% |
| `report` | 79.8% |
| `sources` | 76.3% |
| **Total** | **61.8%** |

Coverage is statement coverage, not detection accuracy. `cmd/generate/config/rules`
coverage is low because most of its files are declarative rule constructors; the Airtable
keyword defect shows why low rule-definition coverage matters.

## 4. Architecture measurements

Method: a temporary Go AST tool outside the repository counted production and test lines,
function spans, and a branch proxy (`if`, `for`, `range`, `case`/`comm`, `&&`, `||`).
It skipped `_test.go` and `cmd/generate` for control-plane metrics. A separate line count
included generator and test code.

Totals:

```text
Non-generator production: 62 files, 10,075 lines
Generator:                137 files, 8,380 lines
Tests:                     85 files, 15,859 lines
Total Go:                 284 files, 34,314 lines
```

Largest production functions:

```text
279 lines / 73 branch proxy  cmd/root.go:270 Detector
214 / 38                     detect/detect.go:632 detectRuleWithPaths
188 / 40                     config/config.go:134 translate
160 / 24                     cmd/report_pair.go:96 writePaired
138 / 23                     sources/git.go:759 Fragments
111 / 17                     sources/files.go:78 scanTargets
109 / 20                     detect/utils.go:224 printFinding
109 / 18                     sources/file.go:299 extractorFragments
 95 / 26                     cmd/root.go:730 findingSummaryAndExit
 93 / 12                     sources/git.go:898 processDiffFile
```

Filesystem call census in non-generator production:

```text
os.Stat 26; os.Lstat 11; filepath.EvalSymlinks 5; filepath.Clean 18
os.CreateTemp 4; os.Rename 1; os.Open 6; os.ReadFile 3
viper. 22; strings.ToLower 17
```

No import cycles were found among the 18 packages. The `report` package imports
`sources`; `detect` imports `report` and `sources`; `sources` imports `config`.

## 5. Manual reproductions

All fixtures are synthetic and live under `C:\Temp`; no fixture was committed. Commands
below are abridged; `C:\Temp\gitleaks-v15.exe` was built from the accepted HEAD.

### 5.1 `withinColumns` compares across lines (#2122)

Config: `primary` rule `P_[A-Z]{4}`, required `context` rule `C_[A-Z]{4}` with
`withinColumns = 1`. Input: `P_AAAA` on line 1, `C_BBBB` on line 6.

Observed:

```text
Finding: P_AAAA
Required: context:6:C_BBBB
```

The primary finding was emitted with a required rule five lines away because
`withinProximity` checks `abs(primary.StartColumn - required.StartColumn)` without a
same-line guard (`detect/detect.go:943-968`). This is a confirmed false-positive /
composite-semantics defect.

### 5.2 Airtable PAT keyword blind spot (#2257)

Input A: the synthetic PAT alone. Input B: `airtable_pat = "<PAT>"`.

Observed: A produced no finding; B produced `airtable-personnal-access-token`.
Source: `cmd/generate/config/rules/airtable.go` uses a self-signaling `pat...` regex but
`Keywords: []string{"airtable"}`. The fragment prefilter never evaluates the regex on
the bare token. Confirmed.

### 5.3 Archive with an unrecognized extension (#2010)

Created a real ZIP containing a synthetic GitLab token, copied it to `compressed.other`,
and scanned with `--max-archive-depth 2`.

Observed:

```text
File: .../compressed.zip!secret.txt   (zip scanned)
DBG skipping binary file mime_type=application/zip path="...\compressed.other"
```

The `.other` archive was skipped as binary. Source: `archives.Identify(ctx, s.Path, nil)`
uses the filename and nil content; a `NoMatch`/nil format falls through to plain content,
then MIME classifies it as `application/zip` and skips it (`sources/file.go:102-149`,
`487-502`). Confirmed.

### 5.4 Base64 text decoding to printable Unicode (#2000)

Input: base64 of `说明 token=glpat-<synthetic> 结束`.

Observed: no finding. Source: `detect/codec/base64.go` accepts only
`isPrintableASCII(decodedValue)`, so valid printable UTF-8 is rejected. Confirmed at the
codec boundary; the exact renderer/finding expectation still needs a fixture in C1.

### 5.5 Generic rule leading `+`/`/` (#2327)

Input: `api_token = "+Q2x9vLm8Zp4Rt7Yw1Nc6Bd3Fg5Hj0Ks"`.

Observed: no finding. Accepted generated regex value alternatives are
`[\w.=-]{10,150}` and `[a-z0-9][a-z0-9+/]{11,}={0,3}`; neither can start with `+` or
`/`. Confirmed source-level false negative; whether every such string is a credential is
a detection-policy question for C1/C6.

### 5.6 Later-line byte column (#1424)

Input (LF): `first line\nxx ALPHA_1234ABCD`.

Observed JSON: `StartLine: 2`, `StartColumn: 5`, `EndColumn: 18`.

The `ALPHA` token begins at byte column 4 on line 2. `location.go` computes byte columns
from `prevNewLine` (the newline byte) while `lineStart` is already available and is used
only for codepoint columns. Confirmed later-line +1 drift; CRLF shows the same class.
SARIF's codepoint path is separate and should remain unchanged.

### 5.7 Extended rule drops RequiredRules (#2024)

Base config defines rule `alpha`; child config extends it, redefines `alpha`, and adds
`[[rules.required]] id="context" withinLines=1`, plus a `context` rule. Input contains
both tokens on one line.

Observed: `alpha` finding has no `Required:` line; `context` appears separately.
Source: `config.extend` merges description/entropy/secret group/regex/path/tags/keywords/
allowlists but not `RequiredRules` or `SkipReport` (`config/config.go:496-529`). Confirmed.

### 5.8 No-final-newline redaction (#1243)

Input: `some.datasource.password=password12345\nasana=1234567890123456` with no trailing
newline, scanned with `--redact`.

Observed:

```text
Finding: asana=REDACTED
Secret:  REDACTED
Line:    2
```

No preceding-line password appeared. The V11 line-start correction (`f6c111c8`,
`detect/location.go:59,73`) plus known-exposure grouping (`detect/utils.go:166-217`)
corrects the reported case. Retired.

### 5.9 `targetRules` on a global allowlist (#1919)

Config defines two rules (`alpha`, `beta`) and a path-targeted allowlist
`targetRules=["alpha"]` for `input.txt`, which contains both tokens.

Observed: only `beta` was reported; `alpha` was suppressed. A separate regex-targeted
fixture behaved the same. The accepted loader attaches target rules after extension
(`config/config.go:273-282`). Retired.

### 5.10 Profile script (#1900)

Executed with Git Bash, a path containing a space, a synthetic finding, `DECODE_DEPTH=1`,
and a Git-initialized target. `hyperfine` and `yq` were not installed.

Observed: exit 0; `dir` and `git` each produced `cpu.pprof`, `mem.pprof`, `trace.out`;
metadata recorded target, Git HEAD, shallow=false, Go version and decode depth. The script
uses the same command array for benchmark and profile and includes `--exit-code=0`
(`scripts/profile.sh:66,72,86`). Functional behavior verified; permanent shell test
remains absent.

## 6. CI and artifact limitations

- Post-merge CI is the strongest available operational evidence. It does not run the
  independent scanner harness/oracle, measure detection precision/recall, or establish
  release-artifact signatures.
- Docker CI built and smoke-tested amd64/arm64 root/non-root images, but no image was
  published and no registry digest was verified.
- The `gitleaks` action check is skipped by policy. It is not a security pass.
- No GoReleaser run, SBOM, Cosign signature, checksum manifest or release binary exists.
- #2017's RE2 constrained-address-space behavior was not tested against a release binary.
- `go generate ./...` was not re-run in V15 because no product source was edited; V14's
  generated-config evidence remains the accepted baseline.
- No real provider API was contacted. Provider-format claims remain source/patch review
  unless a synthetic format fixture was executed.

## 7. Confidence statement

**Verified by execution:** baseline trees/CI; Go build/test/RE2; coverage; the manual
reproductions above; profile script smoke; targetRules/redaction behavior.

**Verified by source inspection:** `Detector` complexity, `location.go` byte/codepoint
split, `config.extend` merge omissions, Airtable keyword, archive `Identify`/MIME path,
base64 ASCII filter, generator keyword semantics.

**Not verified:** 24 GB decode OOM, macOS PTY hang, Docker Desktop pre-commit failure,
RE2 RLIMIT_AS release behavior, full provider token formats, independent detection
accuracy, Windows ACL/key lifecycle, release signatures/SBOM.

No finding-count or benchmark result in this appendix is presented as detection accuracy.

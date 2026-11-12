# Contribution guidelines

## Build toolchain

The minimum supported Go version is declared by the `go` directive in
[`go.mod`](go.mod) (currently Go 1.26.0). The build and test toolchain used by
CI and by the container image is pinned to **Go 1.27.4**:

- CI: [`.github/workflows/test.yml`](.github/workflows/test.yml)
- Container build stage: [`Dockerfile`](Dockerfile) (`golang:1.27.4-alpine3.24`)

Building with a Go release older than the `go.mod` requirement fails at the
module version gate before compilation.

## Build and test commands

Routine commands are non-mutating. They do not format, generate, tidy, fetch
remote tags, or rewrite tracked configuration:

| Command | Purpose |
| --- | --- |
| `make build` | Compile the scanner with the resolved version stamp. |
| `make test` | Run the full test suite without the race detector. |
| `make test-race` | Run the full test suite with the race detector. |
| `make test-cover` | Run race tests with coverage and open the HTML report. |
| `make failfast` | Stop at the first test failure. |
| `make verify` | Run `go vet`, compare generated configuration, and check module consistency on Linux/macOS. |
| `make lint` | Run `golangci-lint`. |
| `make version` | Print the resolved version and its source. |
| `make help` | Show the command contract. |

`make test` intentionally does not imply `-race`; use `make test-race` for the
racy build. The same commands work from a cached module set with network access
disabled. On a Windows CRLF checkout, `make verify` skips `go mod tidy -diff`
because Git's line-ending conversion makes the diff a false positive; the Linux
CI job performs that check.

Mutating maintenance is explicit:

| Command | Purpose |
| --- | --- |
| `make format` | Rewrite Go sources with `go fmt`. |
| `make tidy` | Rewrite `go.mod`/`go.sum` with `go mod tidy`. |
| `make generate` | Regenerate `config/gitleaks.toml`. |
| `make clean` | Remove profiles and test artifacts. |

Do not make routine builds or tests depend on these commands. `make verify`
checks generated configuration by running the generator into a temporary file;
it fails rather than repairing a stale tracked configuration.

### Version resolution

The build version is resolved in this order:

1. `VERSION` supplied on the command line or in the environment.
2. An exact local Git tag, with `+dirty` when tracked files are modified.
3. The nearest local tag plus `+dev.<commit>` build metadata.
4. `v8.25.0+dev.source` when Git metadata is unavailable.

No remote fetch is performed. The fallback base is the minimum version accepted
by the embedded default configuration, so tagless source builds remain valid
SemVer. Pass an explicit release version when producing a release artifact:

```bash
make build VERSION=v8.30.0
```

### Regex backends

The application uses the standard library regexp package by default. The
`gore2regex` build tag selects the RE2 backend:

```bash
go build ./...
go test -count=1 ./...

go build -tags gore2regex ./...
go test -count=1 -tags gore2regex ./...
```

Both the generator packages and the application must compile and test under
both backends. The generated configuration is identical; run `make verify` (or
`go test ./cmd/generate/config -run TestGeneratedConfigMatchesTracked`) after
changing generator code.

### Pre-commit integration

The Linux development CI installs a pinned pre-commit version and runs
`TestPreCommitHookIntegration` with `REQUIRE_PRE_COMMIT=1`. In that mode a
missing Python or pre-commit installation fails the job instead of skipping it.
The test builds the current scanner and exercises the real `try-repo` flow for
zero, one, multiple, and space-containing staged files, `--all-files`, clean
changes, synthetic findings, and scanner configuration errors. The manifest
structure test runs on every platform even when pre-commit is unavailable.

### CI coverage and intentional skips

The Linux job runs the full race suite and real pre-commit integration. Windows
and both macOS jobs run the full race suite; pre-commit is not provisioned
there, so that integration test skips. Windows also skips tests that
require POSIX permission semantics, FIFO or special filename support, or
unprivileged symlink creation. The Windows relative-root path fixtures now
derive a relative root from the fixture's parent when the temporary directory is
on another volume, so those F16 cases execute instead of skipping. The separate
`gore2regex` job builds the application and runs the complete tagged test
suite. Docker jobs cover the default root image and the opt-in non-root target
on `linux/amd64` and `linux/arm64`. The upstream Gitleaks Action workflow is not
part of this fork's CI.

## General

### Issues

If you have a feature or bug fix you would like to contribute please check if
there are any open issues describing your proposed addition. If there are open
issues, make a comment stating you are working on fixing or implementing said
issue. If not, then please open an issue describing your addition. Make sure to
link your PR to an issue.

### Pull Requests

Fill out the template as best you can. Make sure your tests pass. If you see a
PR that isn't one you opened and want it introduced in the next release,
give it a :thumbsup: on the PR description.

## Adding new Gitleaks rules

If you want to add a new rule to the [default Gitleaks configuration](https://github.com/zricethezav/gitleaks/blob/master/config/gitleaks.toml) then follow these steps.

1. Create a `cmd/generate/config/rules/{provider}.go` file.
   This file is used to generate a new Gitleaks rule.
   Let's look at `beamer.go` for example. Comments have been added for context.

   ```golang
   func Beamer() *config.Rule {
       // Define Rule
       r := config.Rule{
           // Human readable description of the rule
           Description: "Beamer API token",

           // Unique ID for the rule
           RuleID:      "beamer-api-token",

           // Regex capture group for the actual secret



           // Regex used for detecting secrets. See regex section below for more details
           Regex: GenerateSemiGenericRegex([]string{"beamer"}, `b_[a-z0-9=_\-]{44}`, true)

           // Keywords used for string matching on fragments (think of this as a prefilter)
           Keywords: []string{"beamer"},
       }

       // validate
       tps := []string{
           generateSampleSecret("beamer", "b_"+secrets.NewSecret(alphaNumericExtended("44"))),
       }
       fps := []string{
           `R21A-A-V010SP13RC181024R16900-CN-B_250K-Release-OTA-97B6C6C59241976086FABDC41472150C.bfu`,
       }
       return validate(r, tps, fps)
   }
   ```

   Feel free to use this example as a template when writing new rules.
   This file should be fairly self-explanatory except for a few items;
   regex and secret generation. To help with maintenance, _most_ rules should
   be uniform. The functions,
   [`GenerateSemiGenericRegex`](https://github.com/zricethezav/gitleaks/blob/master/cmd/generate/config/rules/rule.go#L31) and [`GenerateUniqueTokenRegex`](https://github.com/zricethezav/gitleaks/blob/master/cmd/generate/config/rules/rule.go#L44) will generate rules
   that follow defined patterns.

   The function signatures look like this:

   ```golang
   func GenerateSemiGenericRegex(identifiers []string, secretRegex string, isCaseInsensitive bool) *regexp.Regexp

   func GenerateUniqueTokenRegex(secretRegex string, isCaseInsensitive bool) *regexp.Regexp
   ```

   `GenerateSemiGenericRegex` accepts a list of identifiers, a regex, and a boolean indicating whether the pattern should be case-insensitive.
   The list of identifiers _should_ match the list of `Keywords` in the rule
   definition above. Both `identifiers` in the `GenerateSemiGenericRegex`
   function _and_ `Keywords` act as filters for Gitleaks telling the program
   "_at least one of these strings must be present to be considered a leak_"

   `GenerateUniqueTokenRegex` just accepts a regex and a boolean indicating whether the pattern should be case-insensitive. If you are writing a rule for a
   token that is unique enough not to require an identifier then you can use
   this function. For example, Pulumi's API Token has the prefix `pul-` which is
   unique enough to use `GenerateUniqueTokenRegex`. But something like Beamer's API
   token that has a `b_` prefix is not unique enough to use `GenerateUniqueTokenRegex`,
   so instead we use `GenerateSemiGenericRegex` and require a `beamer`
   identifier is part of the rule.
   If a token's prefix has more than `3` characters then you could
   probably get away with using `GenerateUniqueTokenRegex`.

   Last thing you'll want to hit before we move on from this file is the
   validation part. You can use `generateSampleSecret` to create a secret for the
   true positives (`tps` in the example above) used in `validate`.

1. Update `cmd/generate/config/main.go`. Extend `configRules` slice with
   the `rules.Beamer(),` in `main()`. Try and keep
   this alphabetically pretty please.

1. Run `make generate`

1. Check out your new rules in `config/gitleaks.toml` and see if everything looks good.

1. Open a PR

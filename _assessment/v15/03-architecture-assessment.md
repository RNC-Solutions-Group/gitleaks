# CP-V15 — Architecture assessment

**Baseline:** `rnc-mod@7d74921841114be7677a423cf60620210e832bf7`

This is a current-state reconnaissance, not the post-Phase-II structural refactor. No
production source was changed during V15.

## 1. Package map

| Package | Production files / lines | Responsibility | Architecture note |
|---|---:|---|---|
| `cmd` | 17 / 2,813 | Cobra commands, Viper loading, detector construction, output policy, publication, diagnostics | **Control-plane concentration.** `Detector` alone is 279 lines / 73 branch proxy. |
| `detect` | 18 / 2,741 | Rule evaluation, allowlists, baseline, decoding, location, console formatting, SCM links | **Mixed layers.** Core engine also owns presentation and links. `detectRuleWithPaths` is 214/38. |
| `sources` | 7 / 2,120 | Directory/Git/explicit-file readers, fragments, archives, path identity | **Git monolith.** `sources/git.go` is 1,010 lines; reader and parser share the package. |
| `config` | 5 / 930 | TOML schema, Viper translation, extension, allowlists, strict validation | **Global state.** 22 `viper.` references; package-global configuration is a concurrency/API residual. |
| `report` | 9 / 876 | Finding model, redaction, JSON/CSV/JUnit/SARIF/template, grouped HMAC | **Source coupling.** Imports `sources` for `Finding.Fragment`; serialization model is part of the public API. |
| `internal/informative` | 1 / 463 | Coverage collector, schema, bounded detail, console/JSON document | Self-contained; 86.8% statement coverage. |
| `internal/scanerr` | 1 / 59 | Bounded path-aware error collection | Small, well-tested (95.7%). |
| `regexp` | 2 / 30 | stdlib / `gore2regex` abstraction | Small and clean; release runtime behavior is the residual, not the wrapper. |
| `logging` | 1 / 39 | Global zerolog configuration | Global logger is practical for CLI, limiting for embedders. |
| `version` | 1 / 4 | Build-stamped product version | Product version vs configuration compatibility is unresolved. |
| `cmd/generate` | 137 / 8,380 | Rule generator, provider definitions, stopwords, secrets tooling | **Large generated-rule surface.** 131 files in `rules`; generator validation is the weak point. |

Totals: **10,075 non-generator production lines, 8,380 generator lines, 15,859 test
lines**, 284 Go files. Test/source ratio is healthy (1.57 non-generator; 0.92 including
generator code). It does not establish detection accuracy.

## 2. Dependency direction

```text
main
  -> cmd
       -> config -> regexp, version, logging, viper
       -> sources -> config, internal/informative, internal/scanerr, scm
       -> detect -> config, sources, report, codec, scanerr, informative
       -> report -> config, sources          <-- presentation model reaches back to sources
       -> internal/informative
       -> internal/scanerr
       -> regexp
```

There are **no import cycles** among the 18 packages. The package graph is small enough
to reason about. The problem is not package count; it is responsibility concentration and
a few cross-layer reversals:

1. `report` depends on `sources` solely because `Finding.Fragment` is a
   `*sources.Fragment`. A public report type should not require the reader package.
2. `detect` owns console rendering and SCM URL construction (`detect/utils.go`), although
   `report` is the presentation package and `cmd/scm` owns platform classification.
3. `sources` knows allowlist path policy through `config`, so a reader makes
   suppression decisions before detection.
4. `cmd` owns report publication, while `report` owns serialization; the security-critical
   transaction sits outside the package that defines the output contract.
5. `config` uses package-global Viper and logging, so two configurations cannot be
   loaded independently in one process without careful sequencing.

## 3. Control flow

Ordinary CLI path:

1. `cmd` loads Viper, applies flags/environment/config precedence, translates into
   `config.Config`.
2. `cmd.Detector` constructs `detect.Detector`, loads ignore/baseline, resolves report
   mode, constructs primary/informative/grouped outputs and validates the destination
   set.
3. A command-specific runner (`runDetect`, `runGit`, `runFiles`, `runStdin`) builds
   `sources` and calls detector methods.
4. `sources` emits `Fragment` values; `detect` evaluates rules, allowlists, baselines,
   required rules and decoding; `report.Finding` values are accumulated.
5. `cmd.reportOutput` renders to a restrictive temporary file and publishes primary /
   companion / grouped outputs.
6. `cmd` prints findings and maps scan/finding state to exit codes.

The control plane is the only place that knows the complete destination policy. That is
both a strength (one fail-closed path) and a debt (a 279-line function with many
`logging.Fatal` exits, hard to test in-process).

## 4. Quantitative observations

| Measure | Result |
|---|---|
| Packages with import cycles | 0 |
| Non-generator production lines | 10,075 |
| Generator lines | 8,380 |
| Test lines | 15,859 |
| Statement coverage | 61.8% overall; `detect` 84.1%, `report` 79.8%, `sources` 76.3%, `cmd` 69.9%, `config` 64.2% |
| `cmd` test duration | 98.2 s for the coverage run; roughly 89 s for ordinary runs |
| Functions > 100 lines | 8 in non-generator code |
| Highest branch proxy | `cmd.Detector` 73; `config.translate` 40; `detect.detectRuleWithPaths` 38; `detect.createScmLink` 29; `detect.DetectContext` 27 |
| Filesystem calls in production | 26 `os.Stat`, 11 `os.Lstat`, 5 `filepath.EvalSymlinks`, 18 `filepath.Clean` |
| Viper references in production | 22 |

### Largest functions

| Function | Lines | Branch proxy | Why it matters |
|---|---:|---:|---|
| `cmd.Detector` | 279 | 73 | Combines flags, config, ignore/baseline, report policy, key handling, output set and fatal exits |
| `detect.detectRuleWithPaths` | 214 | 38 | Combines path-only finding, regex matching, decoding provenance, location, entropy, allowlists and required rules |
| `config.translate` | 188 | 40 | Combines rule validation, allowlists, extension, version checks, file tracking and target-rule attachment |
| `cmd.report_pair.writePaired` | 160 | 24 | Paired publication and rollback, security-critical |
| `sources.git.Fragments` | 138 | 23 | Command construction, stream parsing, metadata and section checks |
| `sources.files.scanTargets` | 111 | 17 | Walk, size/skip policy, archive routing, error aggregation |
| `detect.printFinding` | 109 | 20 | Presentation and escaping in the detector |
| `sources.file.extractorFragments` | 109 | 18 | Archive recursion, member policy, size limits |
| `cmd.findingSummaryAndExit` | 95 | 26 | Output and exit policy |
| `sources.git.processDiffFile` | 93 | 12 | Diff parser and finding assembly |

Branch proxy counts `if`, `for`, `range`, `case`/`comm` and `&&`/`||`; it is a ranking
proxy, not formal cyclomatic complexity.

### Coverage interpretation

Coverage is strongest where correctness/cost is highest (`detect`, `report`, `sources`,
`internal`). `cmd` is lower because much of its logic is only reachable through
`logging.Fatal` subprocess tests. The `cmd/generate/config/rules` package is 2.2% covered;
that is the rule-definition surface where the Airtable keyword defect survives. Coverage
does not measure rule-format correctness.

## 5. Technical-debt register

| ID | Debt | Severity | Evidence | Disposition |
|---|---|---|---|---|
| A1 | `cmd.Detector` control-plane god function; fatal exits prevent in-process testing | High | 279 lines / 73 branch proxy; 79 `logging.Fatal` calls in `cmd`, 64 in `root.go` | AG1: extract option loading, policy validation, report plan; return errors |
| A2 | Global Viper/config state; three concerns in `config.translate` | High | 22 `viper.` references; 188/40 function | AG1: instance-scoped loader; split translation phases |
| A3 | `sources/git.go` mixes command runner, parser, metadata, remote normalization | High | 1,010 lines; 138/23 `Fragments` | AG1/C2: split runner/parser/normalizer without semantic change |
| A4 | `detect.detectRuleWithPaths` mixes matching, location, metadata and suppression | High | 214/38; withinColumns defect lives here | AG2: extract matcher, location builder, finding annotation, required-rule join |
| A5 | Detection owns console/SCM presentation | Medium | `detect/utils.go` `printFinding`, `createScmLink` | AG2: move presentation to report/console adapter |
| A6 | Report publication split across `cmd` and `report`; pair logic complex | High (security-sensitive) | `writePaired` 160/24; 468-line file | AG1: introduce a report plan/transaction value; preserve tests byte-for-byte |
| A7 | Generator keyword/regex relationship not validated | High for provider work | Airtable PAT: keyword not implied by regex; 2.2% coverage | C4/C7: add generator lint and self-signaling token rule |
| A8 | Archive/MIME routing and size policy spread across reader | Medium | `sources/file.go` 512 lines; #2010 inverse | C1/C3: content-first archive detection; keep bounded prefix |
| A9 | `report` depends on `sources` for `Finding.Fragment` | Medium | report import graph | AG2: decouple public report model from reader type |
| A10 | Module identity/API stability unresolved | High for release | Upstream `/v8` path for an RNC product; no `/v2` decision | OD-7 / C9; freeze only after AG2 |

### What is not a problem

- Package count and file count are not excessive for a scanner of this scope.
- A high test/source ratio is not waste; the tests are the main evidence for operational
  behavior.
- Long functions in security-critical publication code are not automatically improved by
  splitting; they need parity tests first.
- Some duplication in report modes is deliberate defense-in-depth. The target is one
  validated destination plan, not fewer lines.

## 6. Consolidation opportunities

1. **Report plan object.** One value that owns primary, informative and grouped paths,
   mode, protected inputs, reservation, publication order and rollback. `Detector` should
   build it, not perform each step inline.
2. **Config loader instance.** `Load(viperInstance, options) -> Config, error` with no
   package-global state. Split rule translation, allowlist translation, extension merge
   and version-policy validation.
3. **Git acquisition split.** `gitCommand` (argv/env/stderr), `gitParser` (stream/section
   state), `remoteNormalizer` (SCM identity). C2's merge and shallow work then has a
   testable seam.
4. **Detection pipeline seam.** `Match` → `Locate` → `Annotate` → `Allowlist` →
   `RequiredJoin`. The composite `withinColumns` fix belongs at the join; generic rule
   work belongs at matching; column work belongs at location.
5. **Rule-generation invariant.** At minimum: every self-signaling token rule's keyword
   must be a literal implied by its regex, and every generated rule must have positive
   and negative fixtures. This is a lint, not a new framework.
6. **Report model boundary.** Keep `Finding` JSON-compatible; remove the public
   dependency on `sources.Fragment` before declaring a stable Go API.

## 6.1 What was not measured

- **Dead code:** no `deadcode`/`staticcheck` analysis was run. Unused symbols are not
  claimed.
- **Allocation/CPU hot spots:** the V14 benchmark and V15 profile-script smoke provide
  artifacts but not a new controlled production profile. No allocation claim is inferred
  from source size.
- **Duplication:** the report-mode duplication and layers above are source-observed, not
  token-clone measurements.
- **Public API count:** exported surface is described qualitatively; an API compatibility
  inventory belongs in AG2.
- **Build/test durations:** the `cmd` suite is ~89–98 s and dominates; no per-test timing
  study was performed.

## 6.2 Upcoming-work architecture risks

| Upcoming work | Architectural risk | Control |
|---|---|---|
| C2 merge coverage | First-parent merges add duplicate findings and new fingerprints; low-level `git` changes can bypass section/drain checks | Differential merge matrix; duplicate/fingerprint counts; no parser bypass |
| C3 path identity | Three independent path rewrites can create incompatible ignores/baselines | One canonical identity; migration/dual-read; no per-call `EvalSymlinks` sprawl |
| C4 config model | A global Viper loader and overlay features can silently suppress rules or race | Instance-scoped loader; precedence matrix; negative allowlist controls |
| C5 report/fingerprint | A second writer or unkeyed hash creates a new confidentiality surface | One report plan; keyed v2 only; redaction of every excerpt |
| C6 generic accuracy | A single regex change can alter hundreds of rules and hide false negatives | Oracle/holdout; per-credential expectations; staged waves |
| C7 provider repair | Rule IDs, keyword prefilters and generated config drift | Generator invariant; exact format fixtures; no config-only edits |
| C9 release | Module/version/attestation changes can invalidate docs, installs and API | OD-2/OD-6/OD-10; consumer verification; no release claim from build success |

## 7. Gate placement

A single post-Phase-II audit is too late. The first gate must occur after C1 because the
next three waves all touch the same control plane:

- **AG1 after C1:** `cmd`/`config`/report-plan/Git seam. No detection semantics; parity
  tests and byte-identical reports/fingerprints.
- **AG2 after C6:** detection pipeline, generator invariants, public report/API boundary.
  This is the point to freeze a public Go API or decide not to have one.

Do not use line-count reduction as an objective. Measure: function size/branch proxy,
number of responsibilities per function, testability without subprocess, duplicate
destination validation paths, and parity failures. A refactor that reduces lines but
changes fingerprints, report bytes or exit behavior is a regression.

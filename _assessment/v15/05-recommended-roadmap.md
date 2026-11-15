# CP-V15 — Recommended development roadmap

## 1. Design rule

A development cycle is a **shared-root-cause engineering package**, not a quota of
15–25 upstream items. The operator's proposed 92/88-slot Phase II is the right rough
envelope for core work, but its composition is wrong: it mixes confirmed false negatives,
security policy, elective configuration features, release engineering and 61 provider
rules into one queue.

This roadmap assigns every one of the 220 active packages to exactly one planning cycle.
The exact roster is in [cycle-rosters.md](cycle-rosters.md); the complete evidence fields
are in [disposition-matrix.csv](disposition-matrix.csv).

```text
V16  C1   Critical corrections (7)
     AG1  Control-plane/configuration/report/Git seam
V17  C2   Git/merge coverage and runtime stability (12)
V18  C3   Source readers and path identity (14)
V19  C4   Configuration/allowlist model (27)
V20  C5   Reporting and fingerprint compatibility (7)
V21  C6   Generic accuracy and independent oracle (48)
V22  C7   Existing provider-rule repair (20)
     AG2  Detection model, generator invariants, public API
V23  C8   Provider expansion and rule governance (41)
V24  C9   Release identity and artifact integrity (19)
V25  C10  Elective feature triage (22, mostly decisions)
     VER  Independent harness and release-candidate verification
OPS       External upstream incident response (3; not RNC code)
```

This is **nine implementation/release cycles, one elective-triage gate and two
architecture gates**. It is fewer than the operator's cycle count but larger in scope per
cycle; within C4/C6/C8 the work is split into waves of 8–12 changes so review does not
collapse under a 48-item batch.

## 2. Why this order

1. **C1 reduces silent misses immediately** and fixes shared code before it is refactored.
2. **AG1 prevents the next four waves from editing a 279-line control plane and a global
   config loader independently.** It is cheaper than four local workarounds.
3. **C2/C3/C4/C5 are the core compatibility envelope.** They can be sequenced; C2 and C3
   share Git/path identity, so C2 goes first.
4. **C6 needs a stable reader/decoder and a testable scan setup.** It cannot produce
   trustworthy accuracy evidence while C1/C3 are moving.
5. **C7 repairs existing rules using C6's generic model**; C8 expands only where a token
   format and a negative corpus survive the rule-quality gate.
6. **C9 starts after AG1, not after C6.** Release engineering is a parallel gate; a
   product can produce a release candidate before every provider rule is modernized.
7. **C10 is a decision gate first.** Most elective features should be rejected/deferred
   before engineering time is spent.

## 3. Cycle specifications

### C1 — V16 Critical detector/reader corrections

**Mission:** remove seven confirmed, bounded false-negative/false-positive/correctness
defects without changing defaults beyond the specific corrected behavior.

**Packages (7):** `SCAN-01` #2010; `SCAN-02` #2000/#2172/#2208; `ACC-03`
#2122/#2150/#2255; `ACC-09` #2327/#2335; `ACC-10` #2257/#2260; `CONF-08` #2024;
`REP-02` #1424/#2215.

**Expected changes:**
- bounded content-first archive identification, retaining text-mislabel negatives;
- printable UTF-8 base64 envelope (prefer #2172's control filter, not #2208's NUL/invisible
  acceptance);
- same-line guard for `withinColumns`;
- generic value class allowing leading `+`/`/` without accepting arbitrary punctuation;
- Airtable PAT keyword changed to its self-signaling prefix, with a bare-token fixture;
- extension merge preserves `RequiredRules` and `SkipReport`;
- byte columns from `lineStart`, keeping SARIF codepoints unchanged.

**Touched packages:** `sources/file.go`, `detect/codec/base64.go`, `detect/detect.go`,
`config/config.go`, `detect/location.go`, `cmd/generate/config/rules/`, `config/gitleaks.toml`.

**Dependencies:** none for six packages. `ACC-10` needs operator acceptance of a provider
keyword change (OD-3); the change is narrow and does not alter default rule IDs.

**Effort/risk:** 1 wave, mostly S/M; medium integration risk because these are core
paths. Evidence confidence high for all seven; this is the lowest-uncertainty code wave.

**Acceptance:** each defect has a failing-before/passing-after fixture; standard and
`gore2regex`; no unrelated finding changes; generated config is clean; report/fingerprint
goldens unchanged where not targeted.

**Completion:** all seven dispositions recorded; #2010/#2000/#2122/#2257/#2024/#1424/
#2327 review slots stay active only for residuals. No default merge/path/config policy
change.

**Residuals:** #2010 archive formats outside the supported extractor set; #2000
non-printable decoded content remains excluded; #2327 character classes outside the
tested envelope; #1424 console/CSV historical byte contract is intentionally retained.

### AG1 — Control-plane, configuration, report-plan and Git seam

**Mission:** make the next four waves testable without changing scanner behavior. This is
an engineering refactor, not a feature cycle.

**Scope:**
- Extract `cmd.Detector` option loading and policy validation into small functions; return
  errors instead of 79 fatal exits across `cmd`.
- Introduce a report-plan value that owns primary/informative/grouped paths, mode,
  protected inputs, reservation, publication order and rollback. Preserve current bytes,
  permissions, fingerprints and exit codes.
- Make configuration loading instance-scoped (no package-global Viper dependency in the
  loader) and split `config.translate` into rule/allowlist/extension/version phases.
- Split `sources/git.go` command runner, parser and remote normalizer; no parser change.
- Add golden parity tests: findings JSON, SARIF, CSV, console redaction, fingerprints,
  report bytes, exit status and destination failures.

**Dependencies:** C1 complete; operator approval for a semantics-preserving refactor.

**Effort/risk:** 1–2 waves, medium-large; high review burden, low intended behavior risk.
If parity fails, stop and preserve the old path.

**Acceptance:** no diff in golden reports/fingerprints/exit codes; `go test`, race,
RE2, CI matrix; function-size/branch proxy reduced or responsibilities clearly separated;
new tests run in-process rather than subprocess where practical.

**Residuals:** global logger remains; public API is not frozen yet; report publication
still uses filesystem transactions with the known TOCTOU residuals.

### C2 — V17 Git/merge coverage and runtime stability

**Mission:** close the highest-risk Git false-negative and separate unverified hang/OOM
reports from real defects.

**Packages (12):** `GIT-01` #2003/#2128; `GIT-02` #1028/#2301; `GIT-03` #1993;
`GIT-04` #1647; `GIT-05` #1364; `GIT-06` #1310; `GIT-07` #821; `GIT-08` #1293;
`GIT-09` #1163; `PERF-01` #2019; `PERF-03` #1437; `PERF-04` #2023.

**Expected changes:**
- `GIT-02` begins with OD-1: choose a coverage policy. Recommended: **make
  `--diff-merges=first-parent` the default for secret scanning**, but only with a
  duplicate-reporting rule, fingerprint/baseline migration note, and an opt-out. If the
  operator rejects the default, ship a prominent “comprehensive history” advisory and
  make the opt-in discoverable; do not claim default completeness.
- Complete quoted `--log-opts` behavior; do not accept #2128 as proof of the issue's
  notes/CSV claim.
- Build bounded reproducers for macOS PTY hang, diff-context/K8s, Docker hook, shallow
  provenance and bare-repo ignore expectations.
- Measure decode memory at depths 0/1/3/5 with peak RSS, not post-GC heap; use a finite
  synthetic corpus. Reject #2023 as a wholesale detect rewrite unless its path-only,
  newline-cache and size-limit changes pass differential tests.

**Touched packages:** `sources/git.go`, `detect/detect.go`, `detect/codec`,
`cmd/git.go`, `cmd/root.go`, `scripts`/benchmarks.

**Dependencies:** AG1; OD-1 merge policy; C1 for any shared location/size changes.

**Effort/risk:** 2 waves, medium-large; false-negative impact high, policy/migration risk
high. Evidence for #1028 is strong (V14 measured); #1993/#2019/#1437 remain unverified.

**Acceptance:** ordinary merge, evil merge, conflict resolution, octopus, repeated value,
pathspec, ignore and baseline fixtures; duplicate/fingerprint counts explicit; peak RSS
and stack traces for hang/OOM; no false-clean.

**Residuals:** first-parent duplicates and new fingerprints are a permanent semantic
cost; missing shallow history cannot be reconstructed; blocking OS reads remain.

### C3 — V18 Source readers and path identity

**Mission:** make identical content produce identical findings, locations, paths and
ignore identity across dir/Git/explicit-file/archive and both regex backends.

**Packages (14):** `SCAN-03` #1651/#2324; `SCAN-04` #1757/#2348; `SCAN-05` #2217;
`SCAN-06` #1951; `SCAN-07` #1888/#1936; `SCAN-08` #1821; `SCAN-09` #1641;
`SCAN-10` #1638; `SCAN-11` #1287; `SCAN-12` #1118; `SCAN-13` #1059;
`SCAN-14` #1058; `SCAN-15` #914; `SCAN-16` #2107/#2114.

**Expected changes:** long-line/cross-boundary matching; invalid UTF-8 policy; content
MIME/archive interaction; relative/absolute path and fingerprint normalization; CRLF
parity; path-only duplicate suppression; notebook escaping; multi-root/glob decision;
binary-skip measurement.

**Dependencies:** AG1, C1 archive/decoder, C2 Git identity seam. OD-9 source-local config
trust is related but not a prerequisite.

**Effort/risk:** 2 waves, medium; medium-high compatibility risk because path changes
touch ignores/baselines. Do not combine independent path rewrites.

**Acceptance:** one canonical path identity; positive/negative fixtures for exact,
directory, symlink, absolute/relative, Windows/Unix, CRLF/LF, archive member and
notebook; both regex backends; no widened allowlist.

**Residuals:** rule-level `path` remains reported-path-based unless separately designed;
multi-root/glob remains a feature decision; path identity migration may not preserve every
historical ignore.

### C4 — V19 Configuration/allowlist model

**Mission:** make configuration precedence, extension, allowlist targeting, strictness
and trust explicit and testable; stop accumulating one-off flags.

**Packages (27):** all active `CONF-*` except `CONF-08` (C1), `CONF-10` (retired) and
`CONF-30` (C5). Key groups: overlays #1325/#1251/#1323/#947; strict #1279; minVersion
#1754/#764; deprecation #2120; target/report options #2296; ignore glob #1870; decoder
depth #1732; rule schema #1998/#1344; generator lint #1675/#1932/#1499/#1500; gitignore
patterns #2060/#1195.

**Expected changes:** one resolved-config model; explicit precedence; optional overlay
inputs; strict-mode residual for #1279; minVersion compatibility floor (OD-6); generator
invariants; deprecation of old allowlist syntax; rule-ID migration policy.

**Dependencies:** AG1; OD-6 `minVersion`; OD-9 source-local config trust. #2024 is already
fixed in C1.

**Effort/risk:** 2–3 waves, medium; configuration changes can silently suppress findings,
so every change needs negative controls. Low-value feature rows should be deferred rather
than implemented.

**Acceptance:** precedence matrix; extension chain; targeted/global AND/OR; strict vs
permissive; minVersion warn/reject policy; no default rule loss; no new broad allowlist.

**Residuals:** global Viper state should be gone after AG1 but third-party consumers may
retain assumptions; schema renames need a migration period; `minVersion` semantics remain
an operator policy.

### C5 — V20 Reporting and fingerprint compatibility

**Mission:** settle the report/identity contract before accuracy and provider work add
new finding types.

**Packages (7):** `REP-01` #1866/#2326; `REP-04` #1766; `REP-06` #1708;
`REP-07` #1357; `REP-08` #1074; `REP-09` #1162; `CONF-30` #2239/#2286/#1566.

**Expected changes:** stdout documentation and warning behavior; optional bounded
context; decoded highlighting or explicit non-support; long-line coordinate proof;
optional full-line excerpt; reject unkeyed `SecretSHA`; keyed fingerprint-v2 design
(OD-8) if approved.

**Dependencies:** AG1, C1 column fix, OD-8 fingerprint v2. Reporting context and
multi-format (#1048) remain optional/architecture.

**Effort/risk:** 1–2 waves, medium; schema/redaction/privacy risk. Do not add a second
writer.

**Acceptance:** schema compatibility; report bytes for disabled/default paths; every
excerpt/context redacts all known secrets; fingerprint dual-read/migration; no unkeyed
secret digest.

**Residuals:** two-file crash atomicity remains; fingerprint v2 may remain opt-in;
spreadsheet apostrophe durability remains.

### C6 — V21 Generic accuracy and independent oracle

**Mission:** replace one-off generic-rule patches with a documented credential grammar
and evidence-driven arbitration, measured against an independent oracle.

**Packages (48):** all active `ACC-*` except `ACC-03`, `ACC-09`, `ACC-10` (C1). This
includes generic false negatives (#2223, #1785, #1564, #1772/#1230, #1047/#1308/#1833,
#1239/#1436/#1958, #1331, #1319), false positives (#1897, #1837, #1857, #1830, #1828,
#1728, #1321, #1052/#1073/#1053), private-key behavior (#1043, #1594, #1475, #1085,
#1034), entropy/Unicode (#2230), XML/K8s/Terraform/JSON, and specificity arbitration
(#1997).

**Expected changes:** a generic matcher with explicit value grammar, keyword/prefilter
interaction, entropy normalization and negative allowlist envelope; a frozen oracle
holdout; false-positive budget; provider-specific precedence only where tested.

**Dependencies:** C1, C3, AG1; independent harness/frozen oracle (OD-4); OD-11 provider
policy for overlaps. This is the largest cycle because the backlog is large, not because
all 48 will be implemented.

**Execution:** three waves:
1. reproduce and classify all 48 against the current scanner and oracle;
2. implement 8–12 confirmed high-value changes as one grammar revision;
3. re-run oracle/holdout; defer or reject the rest with evidence.

**Effort/risk:** high; detection behavior changes are the most dangerous in the program.
No finding-count-only acceptance. Do not tune the oracle to make the scanner look better.

**Acceptance:** expected findings by credential, not just count; positive/negative
controls; both regex backends; no broad suppression; performance envelope; holdout
unchanged and independently adjudicated.

**Residuals:** generic detection will remain probabilistic; low-entropy passwords and
provider-specific formats will still need explicit rules; some title-only issues will be
closed as not reproduced rather than fixed.

### C7 — V22 Existing provider-rule repair

**Mission:** repair existing/default rules whose format, keyword, false-positive or
compatibility behavior is demonstrably wrong. Do not add a provider merely because a PR
exists.

**Packages (20):** `RULE-01` OpenRouter; `RULE-02` Supabase PAT; `RULE-03` Sourcegraph
FP; `RULE-04` Basic Auth URL; `RULE-08` SQL Server; `RULE-09` GCP service account;
`RULE-13` LangSmith/Pinecone; `RULE-15` Stripe webhook; `RULE-16` Supabase/DB URI;
`RULE-18` OpenAI; `RULE-19` HashiCorp; `RULE-20` GitHub App; `RULE-23` Anthropic;
`RULE-29` MongoDB; `RULE-41` LaunchDarkly; `RULE-49` New Relic browser; `RULE-52` AWS
FP; `RULE-57` AWS key format; `RULE-59` Azure; `RULE-60` GitLab glpat.

**Expected changes:** choose one canonical format per family, exact rule IDs, generator
and bundled config parity, positive and public/non-secret negatives, no duplicate
findings.

**Dependencies:** C6 generic model; AG2 after C6/C7 for generation/API; OD-11 provider
policy. AG2 can run after the repair wave, before expansion.

**Effort/risk:** 1–2 waves, medium; provider formats change, so fixtures can become stale.
Require an authoritative format source per family.

**Acceptance:** exact token-format fixture; public/anon/ID negative; generated config
clean; no duplicate classification; existing rule ID compatibility or migration note.

**Residuals:** provider formats evolve; New Relic browser token may be removed rather
than “fixed”; Sourcegraph bare-SHA mode may remain context-only.

### C8 — V23 Provider expansion and rule governance

**Mission:** add only provider coverage with a verifiable format, a real secret class,
and a negative corpus. Most of this roster is triage, not implementation.

**Packages (41):** all active `RULE-*` not in C7. The 15 `REJECT` rows include #643,
#2333, #2119, #2118, #2094, #2077, #1996, #1829, #1778, #1335, #1228, #1114, #1086,
#774, #2001. The 18 `DEFER` rows include low-demand/new providers. Six are `COMBINE`
with Azure/K8s/Mongo/SQL families.

**Execution:** a rule-quality gate, then 1–2 expansion waves of at most 10 rules each.
Do not port a config-only fragment; do not accept a PR that also changes CI/Docker/report
code.

**Dependencies:** C7, AG2, OD-11, a maintained provider-format source.

**Effort/risk:** medium, low security urgency; risk is false positives and maintenance
burden. Rejecting low-value rules is a success condition.

**Acceptance:** token format, public/non-secret control, generator/config parity, no
duplicate IDs, rule documentation, update/rotation notes.

**Residuals:** many provider requests remain deferred; a deferred rule is not a claim
that the provider is safe.

### C9 — V24 Release identity and artifact integrity

**Mission:** make an RNC-owned, verifiable release possible; do not publish a release in
this cycle without separate operator authorization.

**Packages (19):** `SEC-05`; `DEP-01`…`DEP-08`; `INT-01`…`INT-06`; `DOC-01`…`DOC-04`.

**Expected changes:** module identity decision (OD-7); `minVersion` compatibility
mapping (OD-6); Go/Git support floors; `go install`/`go run` documentation; release
binary provenance; RE2 constrained-runtime proof or a separate backend artifact; Docker
version stamp; SBOM; Cosign signing rewrite; stable-tag policy; platform docs; reject the
unsafe Actions example.

**Dependencies:** AG1; OD-7, OD-6, OD-10; C6/C7 not strictly required for a release
candidate, but no stable `v2` claim before AG2/VER.

**Effort/risk:** 2 waves, medium-large; release work is cross-cutting but independent.
The main risk is claiming attestation before a consumer verifies it.

**Acceptance:** exact artifact digest, SBOM, checksums, provenance/signature verification
by a separate consumer; version stamp inside pushed images; RE2 constrained run or
documented standard-backend release; no upstream module ambiguity.

**Residuals:** key lifecycle/Windows ACL remains operator-owned; release retention and
revocation policy may remain manual.

### C10 — V25 Elective feature triage

**Mission:** stop optional scope from quietly becoming engineering debt. This is a
decision gate; only survivors get code.

**Packages (22):** all `FEAT-*`. Recommended: 10 reject, 8 defer, 4 accept current
partial behavior. No new network, remote-config, CEL, SDK/MCP, npm, installer-as-is or
source-rewriting feature.

**Dependencies:** C7/AG2 for any integration feature; operator scope decision (OD-12).

**Acceptance:** written benefit/cost/security case; opt-in; bounded; fail-closed; no
change to default scan/report/fingerprint behavior. Otherwise the package remains
deferred/rejected.

### OPS — External upstream incident response

`SEC-01`/`SEC-02`/`SEC-03` are upstream Slack token exposure reports, not RNC scanner
defects. They require the credential owner/upstream maintainer to validate, revoke and
rotate. Do not use or copy reported values. RNC code work is limited to ensuring the
relevant token rules remain correct if separately authorized.

## 4. Architecture and verification gates

### AG2 — Detection model, generator invariants, public API

After C6 and C7, before C8/C10 code:

- split `detectRuleWithPaths` into matching, locating, annotation, allowlist and
  required-rule join;
- add generator lint: keywords implied by self-signaling regex, valid rule IDs, positive
  and negative fixtures, both regex backends;
- decide whether `report.Finding`/`config`/`detect` are a supported Go API;
- freeze only what the operator will support; otherwise document “CLI/binary only”.

### VER — Independent verification and release candidate

- Run the independent RNC scanner harness against a frozen oracle; do not modify/tune the
  oracle.
- Add the permanent `profile.sh` shell regression test and verify #1900's command parity.
- Run standard/RE2, race, Docker root/nonroot, release-artifact and constrained-memory
  checks.
- Verify SBOM/signature/checksum as a consumer, not from build logs.
- Produce operator acceptance evidence; do not publish a release without separate
  authorization.

## 5. First three engineering assignments

1. **V16/C1 — Critical corrections.** Seven packages, one wave. This is the highest
   value per review hour and does not require a broad policy change.
2. **AG1 — Control-plane consolidation.** A semantics-preserving refactor with golden
   parity. It is a prerequisite for doing C2–C5 without repeating the same edits.
3. **V17/C2 — Git/merge coverage.** Begins with OD-1, then implements the chosen policy
   and the bounded Git/runtime reproducers.

The operator should not authorize C4/C6 before C1/AG1 simply because they are “next” in a
category table. C9 can start in parallel after AG1 if release identity is a priority.

## 6. Workload realism

| Cycle | Packages | Relative effort | Review burden |
|---|---:|---|---|
| C1 | 7 | Low-medium | 1 focused PR |
| C2 | 12 | Medium | 1–2 PRs |
| C3 | 14 | Medium | 2 PRs |
| C4 | 27 | Medium-high | 2–3 PRs |
| C5 | 7 | Medium | 1–2 PRs |
| C6 | 48 | High | 2–3 waves, oracle-gated |
| C7 | 20 | Medium | 1–2 PRs |
| C8 | 41 | Low-medium (mostly triage) | rule-quality gate + 2 waves |
| C9 | 19 | Medium-high | 2 PRs + release verification |
| C10 | 22 | Low (decisions) | 1 decision record |

The package count is not the effort. C1 has seven high-confidence fixes; C8 has 41
packages but recommends 15 rejections and 18 deferrals. This is why the 15–25-item quota
should be abandoned.

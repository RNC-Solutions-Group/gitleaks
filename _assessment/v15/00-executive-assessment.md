# CP-V15 — Executive assessment

## 1. Bottom line

RNC-Mod at `7d74921841114be7677a423cf60620210e832bf7` is a credible, independently
maintained scanner with substantially better operational behavior than the upstream
baseline it forked. Its fail-closed execution, report publication, path handling,
redaction, configuration diagnostics and release-source controls have been exercised by
a serious local suite and by the full GitHub Actions matrix on the accepted merge.

It is **not** a release-certified detection product. The independent scanner harness and
frozen oracle were not run, no release artifact exists, and the accepted code still has
confirmed false-negative classes. The most important of these are not exotic: merge-only
changes, archives with ordinary but unrecognized extensions, decoded Unicode text, and a
shipped provider rule whose keyword prefilter cannot fire on its own token format.

The program should stop treating “Phase II,” “architecture,” “accuracy,” and “providers”
as four sequential worlds. Detection correctness is core correctness. Architecture is a
gate, not a phase. Provider rules need a rule-quality gate before expansion. Release
identity is a prerequisite for a credible `v2`, not an afterthought.

## 2. Accepted product condition

| Dimension | Assessment |
|---|---|
| Operational correctness | **Good, with residuals.** Source/Git/worker errors are collected and fail closed; report writes are non-destructive and published by temp/rename; V14B removed the 1,000-candidate aggregate regression and the quadratic collision path. |
| Security engineering | **Materially improved, not complete.** Report destinations, key/config protection, redaction and Git trust boundaries are tested. Source-local config trust remains compatible-by-default; `--no-local-config` is the CI contract. |
| Detection accuracy | **Unverified at product level.** Go fixtures are not an independent oracle. Three confirmed false-negative classes and several provider-rule defects are present. |
| Architecture | **Maintainable but concentrated.** 18 packages, no import cycles, 10,075 production lines, 15,859 test lines. The control plane and four core functions carry disproportionate complexity. |
| Release readiness | **Source-ready only.** Post-merge CI is green; no tag, image, binary, SBOM, signature, checksum or consumer verification exists. |
| Backlog authority | **Reconciled to 180 / 142 / 220.** The V14A arithmetic was incomplete; #1243 and #1919 are already corrected. |

## 3. Highest-value opportunities

1. **Fix confirmed false negatives before broad refactoring.** Merge-only history,
   unknown-extension archives, decoded Unicode and Airtable PAT keyword gating.
2. **Fix two confirmed correctness defects in shared code.** `withinColumns` compares
   columns across different lines; extension merging drops `RequiredRules` and
   `SkipReport`.
3. **Fix byte-column drift on later lines.** CRLF later-line findings report a column one
   byte too high; existing `location.go` mixes `prevNewLine` and `lineStart`.
4. **Consolidate the command/control plane.** `cmd.Detector` is 279 lines with a branch
   proxy of 73 and mixes flag loading, policy, report planning, grouping-key handling and
   fatal exits.
5. **Make release identity a real gate.** Decide module identity and `minVersion`
   semantics; produce and verify an actual signed/SBOM artifact.
6. **Build the independent accuracy lane before claiming accuracy.** A frozen oracle,
   holdouts and new stress workloads must remain separate and independently adjudicated.
7. **Stop low-value provider and feature expansion.** More than one third of the active
   queue is better deferred or rejected than implemented.

## 4. Recommended strategy

**COA 2+ — risk-and-dependency-first, with two architecture gates.**

- **V16 / C1:** seven confirmed corrections.
- **AG1:** control-plane/configuration/report-plan consolidation after C1.
- **V17 / C2:** Git/merge coverage and runtime stability.
- **V18 / C3:** source readers and path identity.
- **V19 / C4:** configuration/allowlist model.
- **V20 / C5:** reporting and fingerprint compatibility.
- **V21 / C6:** generic accuracy campaign with independent oracle.
- **V22 / C7:** existing provider-rule repair.
- **AG2:** detection model/public API freeze.
- **V23 / C8:** provider expansion and rule governance.
- **V24 / C9:** release identity, artifacts and documentation.
- **V25 / C10:** elective-feature triage; decision-only unless a case survives.
- **VER:** independent harness, release-candidate verification, operator acceptance.

C1 and AG1 are the critical path. C9 can start after AG1 and should not wait for the
provider backlog. C6 starts after C1 and C3, not after every core package.

## 5. What the program should not do

Do not implement live credential validation, unkeyed secret hashes, remote
configuration, CEL allowlists, a mutable npm wrapper with broken exit propagation,
the current installer, source-file rewriting, or broad speculative provider patterns.
Do not make a single post-Phase-II refactor the only architecture gate. Do not treat
finding counts, upstream PR merges, or a green Go suite as accuracy evidence.

## 6. Final briefing — direct answers

1. **Actual condition?** A strong operational scanner with a substantial test suite,
   green post-merge CI, and unverified detection quality. Not release-certified.
2. **How much backlog is genuinely actionable?** 220 active packages. Planning triage
   finds 19 clear fixes, 15 adaptations, 16 combinations, 5 direct implementations,
   52 targeted reproductions, 43 deferrals, 35 rejections and 10 release actions, plus
   policy/incident items. Only about 55–70 are immediate code candidates.
3. **Which reported defects are already corrected?** At minimum #1243, #1858, #1919,
   #1714, #1900 and #1809; #811/#827 are corrected within the still-active strict-config
   package. #1858/#1714 were the V14A proposal; #1243/#1919 are V15 additions.
4. **Greatest current risk?** Default Git history omits merge-only patches, so a
   merge-resolution secret can be silently missed. The source-local config default is a
   deliberate compatibility risk. Decode memory and very-long-line/coordinate handling
   remain unresolved.
5. **Which upstream proposals are inferior or unnecessary?** #1992 URL/JSON input,
   #2056 MCP, #2154 SDK, #1013 live validation, #1874 CEL, #1557/#2186 remote config,
   #1661 npm, #1485 installer, #1162 unkeyed hash, #2101 secret-value ignore,
   #2001 mixed provider bundle, #2023 detect rewrite, and many title-only provider rules.
6. **Which changes should combine?** Generic rule clusters, composite `withinColumns`,
   path identity, config overlays, fingerprint/dedup, archive sniffing, Supabase,
   Anthropic/AI providers, payment providers, Azure, Kubernetes and MongoDB.
7. **Is the proposed Phase II structure appropriate?** No. It is a category inventory,
   not a risk-ordered program. It delays confirmed false negatives and includes dozens of
   elective configuration features.
8. **When should architecture review occur?** Twice: AG1 after the seven critical fixes
   and before the core waves; AG2 after generic accuracy and before stable API/provider
   expansion. Not one large audit after Phase II.
9. **What should Phase III contain?** C6: generic detection semantics, false-negative
   recovery, false-positive reduction, entropy/allowlist arbitration and provider-rule
   evidence. It should not be a re-run of 51 title-only issues.
10. **When should accuracy/provider modernization begin?** Accuracy begins in V16 with
    confirmed defects, then C6 after C1/C3/AG1. Existing provider-rule repair begins after
    C6; new provider expansion only after C7 and AG2.
11. **What would you deliberately not implement?** The rejected list in §5 and the
    `REJECT` rows in the disposition matrix.
12. **How many cycles?** Nine implementation/release cycles plus one elective-triage gate,
    with two architecture gates and one verification gate. The cycles are defined by
    shared root cause, not by a 15–25 package quota.
13. **First three assignments?** (1) V16 critical corrections; (2) AG1 control-plane
    consolidation; (3) V17 Git/merge coverage after the operator’s merge-policy decision.
14. **Decisions before those?** Merge-default policy, approval to run architecture
    refactor with parity tests, independent-harness access, and provider-rule policy for
    the early Airtable fix. Module/`minVersion`/release decisions are needed before C9.
15. **Recommended COA?** COA 2+, hybrid risk-first with architecture gates. It beats
    COA 1 by moving confirmed false negatives forward, beats COA 3 by not delaying
    security fixes for a big-bang refactor, and beats COA 4 by stabilizing the model
    before accuracy claims.

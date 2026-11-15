# CP-V15 — Dependency and conflict analysis

## 1. Program dependency graph

```text
C1  Critical corrections (7)
 |
AG1 Control-plane / config / report-plan / Git seam
 |\
 | +--> C2  Git + merge + runtime (12) ----+
 | +--> C3  Readers + path identity (14) --+
 | +--> C4  Configuration model (27) ------+--> C6 Generic accuracy (48)
 | +--> C5  Reporting/fingerprint (7) -----+       |
 |                                                AG2 Detection/API/generator seam
 |                                                   |
 |                                        C7 Provider repair (20)
 |                                                   |
 |                                        C8 Provider expansion (41, mostly triage)
 |
 +--> C9  Release identity / artifacts (19) -- independent of C6/C7 except final release
 |
 +--> C10 Elective triage (22, mostly decisions)

C6/C7 -> AG2 -> C8/C10 code -> VER independent harness + release candidate
```

C1 and AG1 are the critical path. C9 can run in parallel after AG1. C6 cannot start in
earnest until C1/C3 have stabilized the reader/decoder and AG1 has made the scan setup
testable. C7 needs C6's generic model and rule-generation invariants.

## 2. Shared root causes

| Root cause | Packages / proposals | One correction, not many |
|---|---|---|
| Generic rule is a single brittle regex with a hardcoded allowlist envelope | #2327, #2223, #1785, #1564, #1772/#1230, #1047/#1308/#1833, #908/#1637, #1239/#1436/#1958, #1331, #1319, #1997 and related ACC rows | A documented generic-credential grammar with an explicit negative corpus and entropy/allowlist arbitration |
| Keyword prefilter omits a token's self-signaling prefix | #2257 Airtable; LaunchDarkly #1366; other provider rules | Generator invariant: keyword must be a literal implied by the regex, or the rule must be explicitly context-gated |
| Composite-rule proximity is semantically wrong | #2122, #2150, #2255 | Same-line guard plus one proximity matrix; do not merge both PRs |
| Later-line byte columns use the newline byte, not line start | #1424, #2215; overlaps #2238 | One `location.go` correction; retain SARIF codepoint tests |
| Root/path identity differs between Git/dir/absolute/explicit-file | #1287, #1059, #2107/#2114, #1821; #2248 already partial | One canonical path-identity policy with baseline/ignore migration; do not add independent rewrites |
| Config overlay/extension merge is incomplete | #2024, #1325, #1251, #1323, #947, #1344, #2296, #1998 | One resolved-config model with explicit precedence, then optional overlays |
| Ignore fingerprint is location-only | #2239, #2286/#1566, #1162, #1284, #1205 | One keyed v2 design; reject unkeyed hash; preserve v1 dual-read |
| Archive recognition is filename-first | #2010, #2316 (fixed), #1951 | Content-first bounded sniff; keep text-mislabel negatives |
| Git merge patches are omitted by default | #1028/#2301, #1647, #1933 (context) | One merge-coverage policy and duplicate/fingerprint migration; do not adopt #2301 blindly |
| Provider rule families conflict | Supabase, Anthropic/AI, payments, Azure, K8s, MongoDB, Sourcegraph | One canonical rule ID/format matrix per family; preserve all contributors but select one implementation |
| Reporting output has multiple single-output assumptions | #1048, #1766, #1074, #1162, #1048/#1232 | Report-plan architecture before multi-format/context; no unsafe partial writer |

## 3. Competing-PR / conflict matrix

| Area | Packages | Conflict | Recommended resolution |
|---|---|---|---|
| Composite columns | #2122 + #2150, #2255 | Two patches, same bug; different test depth | RNC-authored minimal same-line guard; credit both; do not merge both |
| Generic value syntax | #2327, #2223, #1785, #1564, #1772/#1230, #1047/#1308/#1833 | Overlapping character classes and word-boundary changes; independent PRs conflict | One grammar + oracle corpus; port only tested components |
| XML generic matching | #908 / #1265, #1637 | PR #1265 edits an obsolete generator API; issue #1637 is broader | Independently implement XML tag/attribute cases on current generator |
| Private-key rules | #1043, #1594, #1475, #1085/#2341, #1034 | Broaden/narrow/duplicate behavior; some are test-only | One private-key format matrix; regression test first; avoid broad regex |
| Path identity | #1287, #1059, #2107/#2114, #1821 | `--source-relative-paths` changes fingerprints; root-relative ignore behavior differs | One canonical identity + migration; reject independent display rewrite |
| Config overlays | #1325, #1251, #1323, #947 | Four distinct workflows; overlapping “extra config” ideas | One precedence model; separate rules/allowlist inputs only if proven |
| Strict config | #811/#827, #1279 | Default permissive vs strict; shared package | Keep strict opt-in; retire issue records #811/#827; keep #1279 active |
| Fingerprint/dedup | #2239, #2286/#1566, #1162, #1054 | Unkeyed SHA, keyed HMAC, content-only, commit-scoped, primary dedup | Keyed v2 + dual-read; reject unkeyed; keep primary findings |
| Merge coverage | #1028/#2301 | First-parent recovers merge-only but creates duplicates/fingerprints; `--cc` breaks parser | Policy decision + duplicate/migration design; no blind adoption |
| Archive/MIME | #2010, #1951, #2316 | Content sniff vs MIME allowlist vs extension | Content-first bounded sniff; MIME remains a second gate |
| Supabase | #2225, #2278/#2279, #2280 | New PAT/secret formats vs service_role JWT/DB URI | Separate rule IDs by credential class; test `anon` public control |
| Anthropic/AI | #1467, #2158/#2104, #2094, #2180, #2314/#2177, #2242 | OpenAI legacy, Claude OAuth, session URL, Groq/xAI, LangSmith/LangGraph, ActionGate | One provider-format matrix; reject session URL; resolve LangSmith length collision |
| Payments | #2078, #2117, #2118, #2119, #2303, #1335 | Same Mercado Pago ID; config-only Bundles; broad Chargebee | One canonical generator rule set; reject config-only fragments and broad prefixes |
| Azure | #539, #1687, #1327 | Four PRs, broad CredScan port vs specific storage/SHIR rules | Build one provider matrix; select exact formats; no broad legacy port |
| Kubernetes | #1933, #1513, #1568 | Multi-entry context, additional secret types, secretdata | One K8s rule family; #1933 is the highest-confidence behavior |
| MongoDB | #2007, #2085, #2280 | URI rule vs Atlas tokens; general DB URI | One generator family with separate token/URI rules |
| Sourcegraph | #1898, #2083, #2109, #1697 | Bare SHA false positive; two fixes; test-only fixture | One exact `sgp_` rule; bare 40-hex only in explicitly contextual mode |
| Report context | #1766, #1074, #1048 | Context, full line, multi-format all touch report schema | One bounded report-plan design; defer until AG1 |
| Distribution | #1661, #1485, #1217, #2220, #885 | npm/installer signing/SBOM/tags all release-pipeline changes | One release gate; reject unsafe npm/installer patches; implement only verified pipeline steps |

## 4. Integration hazards

1. **Fingerprint changes invalidate `.gitleaksignore` and baselines.** Any path, merge
   or v2 identity change needs a migration table and dual-read tests. Do not let report,
   merge and provider changes land in one release without a compatibility report.
2. **Two regex backends.** Generator code must import the project regexp abstraction;
   test-only PRs importing stdlib `regexp` can fail the `gore2regex` build.
3. **Global Viper state.** Parallel config loads and SDK/MCP experiments can race. Any
   new embedding feature must be rejected until config is instance-scoped.
4. **Report publication set.** Primary, informative and grouped outputs share protected
   inputs and aggregate families. A new reporter or multi-format change must reuse the
   destination plan; it must not add a second ad hoc writer.
5. **Regex/keyword ordering.** Adding a keyword for a self-signaling rule increases
   candidate fragments but must not change detection of ordinary tokens. Add both a bare
   positive and a near-miss negative.
6. **Generated config drift.** Adding a rule only to `config/gitleaks.toml` or only to
   the generator is incomplete. `go generate ./...` must be clean; a generated config
   diff is part of acceptance.
7. **Default changes.** First-parent, strict config, fingerprint v2 and new provider
   rules can change baseline/ignore semantics. They are release-note and migration
   events, not ordinary fixes.
8. **Exit-code wrappers.** npm/installer/MCP examples must propagate nonzero findings
   and operational errors. A wrapper that returns 0 on a leak is a release blocker.

## 5. Contributor preservation

No V15 retirement or rejection changes upstream attribution. The accepted
`RNC-MOD.md` §3 register remains the authority for adopted/adapted commits. This analysis
adds documentary credit only:

- #1900 — @vignesh1507, independent profile design influence, no adopted shell code.
- #1809 — @JoostVoskuil, Azure SSH normalization research, independent implementation.
- #1919 — @maksym-dilanian, targetRules report; upstream implementation present.
- #1243 — @SyedMonzAli, redaction defect research; V11 correction.
- #1858 — @bdovaz, SARIF level request; opt-in RNC implementation.
- #1714 — @rgmz, local file-link request; opt-in RNC implementation.
- #811/#827 — @svyotov/@sergiomarotco, strict-configuration requests; shared package
  remains with #1279.
- All competing provider PRs retain their authors even when only one design is selected.
  A later RNC implementation is not a claim that an upstream PR was merged or copied.

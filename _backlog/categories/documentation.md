# Documentation — RNC review queue

Snapshot: 2026-11-08 · RNC fork `master` checked at `5f6ec56cfa6f` · Triage branch: `feat/issues-and-pr`.

**8 issues · 13 PRs · 13 work packages.**

Documentation, setup guidance, labels, and examples. Compare proposed text with current code and README behavior before consolidating patches.

**Review states:** `SOURCE-CHECKED` records a targeted source/patch comparison, **not** a built or tested correction. `PENDING` records an unresolved issue or competing proposal. No PR is adopted or upstream item closed by categorization. See [RNC exclusions](exclusions.md) for retired review slots.

## Candidates

| Open issue(s) | Open PR(s) and contributor(s) | Scope and proposed change | RNC status | Source evidence / remaining gate | RNC implementation PR/commit |
| --- | --- | --- | --- | --- | --- |
| [#2173](https://github.com/gitleaks/gitleaks/issues/2173) | [#2174](https://github.com/gitleaks/gitleaks/pull/2174) (@moutons), [#2213](https://github.com/gitleaks/gitleaks/pull/2213) (@rylena), [#2344](https://github.com/gitleaks/gitleaks/pull/2344) (@delaysystem) | Document supported `--platform` choices in CLI help and README. | **SOURCE-CHECKED — TEST PENDING** | #2174 derives values from SCM enum but misses `detect` help; #2213 documents five values; #2344 adds `none` in help but omits README. Consolidate semantics and test both commands before adoption. | — |
| [#2218](https://github.com/gitleaks/gitleaks/issues/2218) | [#2236](https://github.com/gitleaks/gitleaks/pull/2236) (@RajeshRajendiran) | Clarify that range-restricted Git scans only inspect added diff lines; pre-existing unchanged secrets are outside the scan scope. | **DOCS BUNDLE — EXACT-PATCH CANDIDATE** | `sources/git.go` consumes added fragments from `git log -p`, and the current README does not explain pre-receive range limitations. PR adds documentation only; it does not fix pre-receive hook revision selection (#821), nor does it implement rescanning unchanged content. | — |
| — | [#2212](https://github.com/gitleaks/gitleaks/pull/2212) (@RajeshRajendiran) | Correct documented `--max-decode-depth` default from 0/disabled to 5/enabled. | **DOCS BUNDLE — EXACT-PATCH CANDIDATE** | `master/cmd/root.go:91` sets default 5; README still contradicts it twice. Patch fixes both references. No memory/OOM remediation: issue #2019 remains in [performance](performance-stability.md). | — |
| [#1862](https://github.com/gitleaks/gitleaks/issues/1862) | [#2298](https://github.com/gitleaks/gitleaks/pull/2298) (@bharathvemula) | Update baseline prose to `git` command and fix the article/grammar. | **DOCS BUNDLE — EXACT-PATCH CANDIDATE** | Master still says `Once as baseline` and refers to retired `detect`; PR changes the two prose locations. | — |
| — | [#2287](https://github.com/gitleaks/gitleaks/pull/2287) (@ump45nose) | Fix pre-commit Docker-image section link. | **DOCS BUNDLE — ADAPT REQUIRED** | `#docker` is broken, but PR replaces it with `#docker-ghcrio`, which is only a Bash comment inside a code fence in current README, **not** a heading. Use a real heading/anchor or direct URL; credit #2287 without copying the broken anchor. | — |
| — | [#2282](https://github.com/gitleaks/gitleaks/pull/2282) (@haimingZZ) | Fix `maintence` in CONTRIBUTING and `built ontop` in README. | **DOCS BUNDLE — EXACT-PATCH CANDIDATE** | Both misspellings remain in current master; PR changes only these two phrases in two files. | — |
| [#2256](https://github.com/gitleaks/gitleaks/issues/2256) | [#2258](https://github.com/gitleaks/gitleaks/pull/2258) (@jankesec) | `identiifer` → `identifier` in README config example. | **DOCS BUNDLE — EXACT-PATCH CANDIDATE** | Misspelling remains in current master, same line and context. | — |
| [#1855](https://github.com/gitleaks/gitleaks/issues/1855) | [#2247](https://github.com/gitleaks/gitleaks/pull/2247) (@CAOShurong) | Replace stale blog.gitleaks.io config-guide URL with Medium link. | **DOCS BUNDLE — EXACT-PATCH CANDIDATE** | Old blog URL remains; the proposed [Medium guide](https://medium.com/gitleaks/stop-leaking-secrets-configuration-2-3-aeed293b1fbf) was independently found with the matching title/content on 2026-11-08. Recheck at integration time; old example versions inside the article are historical. | — |
| — | [#2084](https://github.com/gitleaks/gitleaks/pull/2084) (@vincent067) | Add `examples/README.md` and a new reusable GitHub Actions example. | **HOLD — NOT SAFE AS-IS** | The proposed `gitleaks git --baseline-path ...` claims to **generate** a baseline (it actually consumes one); default workflow runs an Action rather than the shown manual SARIF upload; suggested `--exit-code 0` conflicts with blocking deployments; workflow summary has malformed shell quoting around `.gitleaksignore`. Requires functional/security review, not a typo bundle. | — |
| — | [#2008](https://github.com/gitleaks/gitleaks/pull/2008) (@jylenhof) | Document `mise use -g gitleaks@latest` as an optional installation. | **DOCS BUNDLE — EXACT-PATCH CANDIDATE** | One additive README snippet, absent from master; confirm current mise registry availability as an integration smoke check before publishing. | — |
| [#1067](https://github.com/gitleaks/gitleaks/issues/1067) | — | Legacy `examples/leaky-repo.toml` benchmark config removed. | **ISSUE — REVIEW PENDING** | Original request refers to a removed 2022-era example; current master does not include it. Whether a maintained benchmark configuration is useful remains a separate decision, not solved by the docs typo bundle. | — |
| [#2295](https://github.com/gitleaks/gitleaks/issues/2295) | [#2306](https://github.com/gitleaks/gitleaks/pull/2306) (@rotirotirafa) | Document Firebase example-key provenance in generated GCP API-key allowlist. | **PATCH-REVIEWED — DOC/GENERATOR GATE** | PR #2306 adds explanatory comments in GCP rule generator/bundled output; verify both remain synchronized when regenerated. No claim that detection or allowlisting semantics change. | — |
| [#1263](https://github.com/gitleaks/gitleaks/issues/1263) | — | Explain safe storage and lifecycle for secret-bearing baseline reports in centralized CI/pre-commit workflows. | **SOURCE-CHECKED — SECURITY DOCS REVIEW** | Upstream asks where a baseline should live and when to prune it. Current README's `Creating a baseline` section describes reuse but does not address that JSON findings can retain `Secret`/`Match` values. Draft security-sensitive guidance on protected access, retention, redacted artifact handling and periodic remediation; do not publish real secrets or assert a product defect. | — |

## Documentation consolidation recommendation — separate owner branch, no adoption here

**Bundle A — eight PRs / seven credited authors:** [#2298](https://github.com/gitleaks/gitleaks/pull/2298) @bharathvemula, [#2282](https://github.com/gitleaks/gitleaks/pull/2282) @haimingZZ, [#2258](https://github.com/gitleaks/gitleaks/pull/2258) @jankesec, [#2212](https://github.com/gitleaks/gitleaks/pull/2212) @RajeshRajendiran, [#2236](https://github.com/gitleaks/gitleaks/pull/2236) @RajeshRajendiran, [#2008](https://github.com/gitleaks/gitleaks/pull/2008) @jylenhof, [#2287](https://github.com/gitleaks/gitleaks/pull/2287) @ump45nose, and [#2247](https://github.com/gitleaks/gitleaks/pull/2247) @CAOShurong.

| Upstream PR | Proposed head commit (pinned 2026-11-08) | Integration approach |
| --- | --- | --- |
| [#2298](https://github.com/gitleaks/gitleaks/pull/2298) | `e63bf8fa4fa06b7c79cd482fbaec49af54f91723` | Preserve exact prose patch |
| [#2282](https://github.com/gitleaks/gitleaks/pull/2282) | `4c25e35e24483e2ebe5da1dba14e64ed6912eb0c` | Preserve both README and CONTRIBUTING corrections |
| [#2258](https://github.com/gitleaks/gitleaks/pull/2258) | `f5b3892e7d30cb5ecdcfc47b12a3ac7a30fc6856` | Preserve exact typo fix |
| [#2212](https://github.com/gitleaks/gitleaks/pull/2212) | `2380e927f3961b686b7194dd3a59f36c94e95bdf` | Preserve exact default/decoding descriptions |
| [#2236](https://github.com/gitleaks/gitleaks/pull/2236) | `0ef7c0a02c81bf794039c4fe09153a2904cf399e` | Preserve explanatory Git/pre-receive limitations |
| [#2247](https://github.com/gitleaks/gitleaks/pull/2247) | `718896a43f42cdce3e0733ccf8f67dab815a6a24` | Preserve replacement link; verify at merge time |
| [#2008](https://github.com/gitleaks/gitleaks/pull/2008) | `f4ff4532d350cf36e923ad36384e88898882a131` | Preserve snippet only after mise registry smoke check |
| [#2287](https://github.com/gitleaks/gitleaks/pull/2287) | `298f58e8320e09a7325ccde723ecfac6d42502c8` | Adapt target anchor; credit original intent |

- **Copy only the exact documentation corrections after applying them to a fresh branch from `master`:** #2298, #2282, #2258, #2212, #2236, and #2247 (replacement Medium article source independently checked). #2008 is a one-line additive candidate, conditional on verifying the current mise installation source. None of these needs scanner-code changes.
- **Adapt with explicit credit:** #2287 needs a real navigable anchor, not `#docker-ghcrio` pointing into a code fence. Reconfirm #2247's external destination at merge time; a current link check cannot guarantee future availability.
- **Keep separate:** #2084 is an executable integration recipe with inaccurate baseline instructions; it is **not safe to copy** into the documentation bundle without a functional audit.
- **History and attribution:** preserve each author's upstream PR URL and source commit SHA in the owner consolidation PR. Prefer retaining original commits where they apply cleanly; adapted fixes should be a distinct owner-authored commit that explicitly credits the originating contributor (`Co-authored-by` when appropriate). Never claim that this triage branch adopted, tested, or merged any source PR.
- **Integration boundary:** the **owner's second development branch** may incorporate this bundle, run a Markdown-link and command/usage check, and open one RNC PR. The upstream-submission branch and this triage branch receive only explanatory README/PR attribution notes as required; no unreviewed docs source changes are cherry-picked there.


## Decision gate

1. Reproduce against a pinned `master` commit (or record why not possible); inspect the PR diff and its test coverage.
2. Compare against the RNC security PR, shared-security branch, and hardened image where relevant; a Docker bump does not change Go module pins or scanner behavior.
3. Record a disposition and evidence. **Adopted** only after upstream code is integrated, with author/PR/SHA and RNC PR/commit. **Corrected** only for RNC-origin fixes. **Superseded** requires an identified replacement; **Closed/Obsolete** requires a verified reason.
4. Update the category ledger and strikethrough without editing immutable `_backlog/data/` exports.

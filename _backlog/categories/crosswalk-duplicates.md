# Crosswalk duplicates — initial review

Source snapshot: [2026-11-08 Gitleaks backlog](../README.md) · [Candidate-link CSV](../data/crosswalk.csv) · Branch: `feat/issues-and-pr`

## Scope and counting

This is a **candidate duplicate register**, not a list of merged or implemented fixes. Items promoted to other active categories are removed here; the source snapshot counts below remain unchanged. It uses the 174 initial candidate relationships and distinguishes high-signal `Fixes/Closes/Resolves` references or explicit issue URLs from low-signal bare number mentions.

| Check | Count | Meaning |
| --- | ---: | --- |
| Open issues in snapshot | 285 | All open upstream issues from the collection |
| Open PRs in snapshot | 252 | All open upstream PRs |
| All initial candidate relationships | 174 | Includes weak textual references and out-of-scope issues |
| Relationships with closing keyword or direct issue URL | 119 | Stronger candidate evidence, **not** proof the code works |
| Issue targets with 2+ strong-signal open PRs | 18 | 17 currently in the open-issue snapshot, 1 outside it |
| Open PRs with strong-signal references to 2+ issue targets | 5 | Sometimes duplicate issues, sometimes distinct problems addressed together |

**Important:** A shared issue number means *competing proposals for one issue*, not necessarily identical patches. Two issues with similar wording can also be distinct. No PR diffs were integrated, scanner builds performed, or technical fixes validated in this review.

## A. Candidate duplicate issues

These are distinct issue records that may describe the same root cause.

| Issue cluster | Assessment | Evidence and preliminary disposition |
| --- | --- | --- |
| [#2121](https://github.com/gitleaks/gitleaks/issues/2121) + [#2160](https://github.com/gitleaks/gitleaks/issues/2160) | **High-confidence duplicate reports** | Both describe `detectRule` invoking `FindAllStringIndex` twice on the same fragment. [#2199](https://github.com/gitleaks/gitleaks/pull/2199) explicitly addresses both. Consolidate into one performance work package; retain both upstream links. |
| [#2048](https://github.com/gitleaks/gitleaks/issues/2048) + [#2120](https://github.com/gitleaks/gitleaks/issues/2120) | **High-confidence overlapping defect** | Both ask to replace deprecated `[allowlist]` with `[[allowlists]]` in bundled/generated config. [#2343](https://github.com/gitleaks/gitleaks/pull/2343) addresses both. Confirm generator and checked-in config scope before selecting one PR. |
| [#1895](https://github.com/gitleaks/gitleaks/issues/1895) + [#1859](https://github.com/gitleaks/gitleaks/issues/1859) | **Likely duplicate operational symptom** | Both report failures of the `gitleaks-system` pre-commit hook; [#1943](https://github.com/gitleaks/gitleaks/pull/1943) references both with a `pass_filenames: false` fix. Verify both reproduction paths. |
| [#2110](https://github.com/gitleaks/gitleaks/issues/2110) + [#2111](https://github.com/gitleaks/gitleaks/issues/2111) + [#2112](https://github.com/gitleaks/gitleaks/issues/2112) | **Same title; do not collapse automatically** | Each is titled “Possible leaked API key in this repository.” Identical titles are not evidence that the underlying exposed values or locations match. Handle sensitively, with no secret values copied here. |
| [#1203](https://github.com/gitleaks/gitleaks/issues/1203) + [#2214](https://github.com/gitleaks/gitleaks/issues/2214) | **Related detection gap, not confirmed duplicates** | Both describe Basic Auth detection, but the examples and encoding contexts differ. A shared test matrix may uncover whether one rule change addresses both. |

## B. Competing PRs for the same issue

This is the complete set of issue targets from the initial crosswalk that have **at least two high-signal open PR references**. Every contributor remains individually credited.

| Issue | Problem | Strong PR candidates | PR authors | Triage |
| --- | --- | ---: | --- | --- |
| [#2122](https://github.com/gitleaks/gitleaks/issues/2122) | bug: withinColumns compares column numbers across different lines in composite rules, producing incorrect prox | 5 | [#2123](https://github.com/gitleaks/gitleaks/pull/2123) (@c-tonneslan), [#2175](https://github.com/gitleaks/gitleaks/pull/2175) (@ardittirana), [#2187](https://github.com/gitleaks/gitleaks/pull/2187) (@driphtyio), [#2201](https://github.com/gitleaks/gitleaks/pull/2201) (@RajeshRajendiran), [#2255](https://github.com/gitleaks/gitleaks/pull/2255) (@AARUSHGAUR) | Review first |
| [#2121](https://github.com/gitleaks/gitleaks/issues/2121) | perf: detectRule() calls r.Regex.FindAllStringIndex twice per match cycle, doubling regex work on every scan | 5 | [#2151](https://github.com/gitleaks/gitleaks/pull/2151) (@wahajahmed010), [#2162](https://github.com/gitleaks/gitleaks/pull/2162) (@driphtyio), [#2163](https://github.com/gitleaks/gitleaks/pull/2163) (@driphtyio), [#2199](https://github.com/gitleaks/gitleaks/pull/2199) (@RajeshRajendiran), [#2271](https://github.com/gitleaks/gitleaks/pull/2271) (@dhanushs11) | Review first |
| [#2048](https://github.com/gitleaks/gitleaks/issues/2048) | Use [[allowlists]] in generated gitleaks.toml | 4 | [#2049](https://github.com/gitleaks/gitleaks/pull/2049) (@dariozachow), [#2053](https://github.com/gitleaks/gitleaks/pull/2053) (@akk85), [#2171](https://github.com/gitleaks/gitleaks/pull/2171) (@driphtyio), [#2343](https://github.com/gitleaks/gitleaks/pull/2343) (@delaysystem) | Review |
| [#2173](https://github.com/gitleaks/gitleaks/issues/2173) | --platform enumeration insufficiently documented | 3 | [#2174](https://github.com/gitleaks/gitleaks/pull/2174) (@moutons), [#2213](https://github.com/gitleaks/gitleaks/pull/2213) (@rylena), [#2344](https://github.com/gitleaks/gitleaks/pull/2344) (@delaysystem) | Review |
| [#2232](https://github.com/gitleaks/gitleaks/issues/2232) | Unreadable files are silently skipped and the run still reports "no leaks found" with exit 0 | 3 | [#2235](https://github.com/gitleaks/gitleaks/pull/2235) (@hamodywe), [#2249](https://github.com/gitleaks/gitleaks/pull/2249) (@vaibhav8a), [#2313](https://github.com/gitleaks/gitleaks/pull/2313) (@Varshith-Kali) | Review first |
| [#1981](https://github.com/gitleaks/gitleaks/issues/1981) | git errors are being incorrectly ignored | 2 | [#2080](https://github.com/gitleaks/gitleaks/pull/2080) (@Smeet23), [#2339](https://github.com/gitleaks/gitleaks/pull/2339) (@delaysystem) | Review |
| [#1846](https://github.com/gitleaks/gitleaks/issues/1846) (outside open snapshot) | Outside current issue snapshot | 2 | [#2082](https://github.com/gitleaks/gitleaks/pull/2082) (@Smeet23), [#2241](https://github.com/gitleaks/gitleaks/pull/2241) (@harriiinnii) | Review |
| [#1898](https://github.com/gitleaks/gitleaks/issues/1898) | False positive - sourcegraph-access-token matches Git SHA | 2 | [#2083](https://github.com/gitleaks/gitleaks/pull/2083) (@Smeet23), [#2109](https://github.com/gitleaks/gitleaks/pull/2109) (@obchain) | Review |
| [#2120](https://github.com/gitleaks/gitleaks/issues/2120) | fix: update default gitleaks.toml to use [[allowlists]] instead of deprecated [allowlist] | 2 | [#2124](https://github.com/gitleaks/gitleaks/pull/2124) (@c-tonneslan), [#2343](https://github.com/gitleaks/gitleaks/pull/2343) (@delaysystem) | Review |
| [#1897](https://github.com/gitleaks/gitleaks/issues/1897) | False Positive: (Some) Vscode keybindings - "key": "ctrl+shift+down" | 2 | [#2125](https://github.com/gitleaks/gitleaks/pull/2125) (@obchain), [#2195](https://github.com/gitleaks/gitleaks/pull/2195) (@alexchen-sys) | Review |
| [#2003](https://github.com/gitleaks/gitleaks/issues/2003) | Passing "--no-notes" via --log-opts option doesn't work with CSV output | 2 | [#2128](https://github.com/gitleaks/gitleaks/pull/2128) (@puneetdixit200), [#2205](https://github.com/gitleaks/gitleaks/pull/2205) (@trinathone) | Review |
| [#2160](https://github.com/gitleaks/gitleaks/issues/2160) | detectRule runs the rule regex over each fragment twice | 2 | [#2161](https://github.com/gitleaks/gitleaks/pull/2161) (@mazen-salah), [#2199](https://github.com/gitleaks/gitleaks/pull/2199) (@RajeshRajendiran) | Review |
| [#2165](https://github.com/gitleaks/gitleaks/issues/2165) | fix: global commit allowlist silently bypassed due to misscoped continue | 2 | [#2166](https://github.com/gitleaks/gitleaks/pull/2166) (@ashvinctrl), [#2200](https://github.com/gitleaks/gitleaks/pull/2200) (@RajeshRajendiran) | Review |
| [#2000](https://github.com/gitleaks/gitleaks/issues/2000) | Decoding skips values with Unicode characters | 2 | [#2172](https://github.com/gitleaks/gitleaks/pull/2172) (@arpitjain099), [#2208](https://github.com/gitleaks/gitleaks/pull/2208) (@trinathone) | Review |
| [#2182](https://github.com/gitleaks/gitleaks/issues/2182) | GitLeaks overreach flagging variable called "privateKey" | 2 | [#2183](https://github.com/gitleaks/gitleaks/pull/2183) (@sanmaxdev), [#2207](https://github.com/gitleaks/gitleaks/pull/2207) (@trinathone) | Review |
| [#2225](https://github.com/gitleaks/gitleaks/issues/2225) | Add detection rules for Supabase secret API keys and personal access tokens | 2 | [#2226](https://github.com/gitleaks/gitleaks/pull/2226) (@Yusuf-Gadelrab), [#2259](https://github.com/gitleaks/gitleaks/pull/2259) (@jankesec) | Review |
| [#2299](https://github.com/gitleaks/gitleaks/issues/2299) | Report file containing discovered secrets is written with default, often world-readable permissions | 2 | [#2300](https://github.com/gitleaks/gitleaks/pull/2300) (@ToniAdreal), [#2312](https://github.com/gitleaks/gitleaks/pull/2312) (@Varshith-Kali) | Review first |
| [#2309](https://github.com/gitleaks/gitleaks/issues/2309) | Add OpenRouter API key rule (sk-or-v1-) | 2 | [#2310](https://github.com/gitleaks/gitleaks/pull/2310) (@Varshith-Kali), [#2325](https://github.com/gitleaks/gitleaks/pull/2325) (@Tarikdxd) | Review |

**Strong candidates for early comparison:**

1. [#2299](https://github.com/gitleaks/gitleaks/issues/2299) — sensitive report permissions. Compare [#2300](https://github.com/gitleaks/gitleaks/pull/2300) and [#2312](https://github.com/gitleaks/gitleaks/pull/2312) for file creation, existing-file behavior, cross-platform effects, and tests. This is an actual security-hardening candidate.
2. [#2232](https://github.com/gitleaks/gitleaks/issues/2232) — unreadable files must not silently produce a clean result. Compare [#2235](https://github.com/gitleaks/gitleaks/pull/2235), [#2249](https://github.com/gitleaks/gitleaks/pull/2249), and [#2313](https://github.com/gitleaks/gitleaks/pull/2313) for errors, exit codes, and partial-scan semantics.
3. [#2122](https://github.com/gitleaks/gitleaks/issues/2122) — `withinColumns` cross-line bug has five strong candidate fixes. Check whether approaches agree on correct matching semantics before judging regression.
4. [#2121](https://github.com/gitleaks/gitleaks/issues/2121) / [#2160](https://github.com/gitleaks/gitleaks/issues/2160) — duplicate regex execution has several proposals. [#2162](https://github.com/gitleaks/gitleaks/pull/2162) and [#2163](https://github.com/gitleaks/gitleaks/pull/2163) even have the same title and author, but distinct head SHAs. Evaluate behavior and benchmarks, not filename/title alone.
5. [#2048](https://github.com/gitleaks/gitleaks/issues/2048) / [#2120](https://github.com/gitleaks/gitleaks/issues/2120) — multiple generator/config changes target one deprecation. Prefer a consolidated implementation with both outputs covered.
6. [#2309](https://github.com/gitleaks/gitleaks/issues/2309) — OpenRouter detection has [#2310](https://github.com/gitleaks/gitleaks/pull/2310) and [#2325](https://github.com/gitleaks/gitleaks/pull/2325) formally proposing the same feature, plus earlier [#2032](https://github.com/gitleaks/gitleaks/pull/2032) without a matching explicit issue reference. Review three implementations, avoid merging duplicate rules.

## C. One PR references multiple issues

These 5 open PRs have strong textual references to more than one issue. This is useful evidence for issue consolidation, but a multi-fix PR does not prove its issues are duplicates.

| PR | Issue references | Assessment |
| --- | --- | --- |
| [#1943](https://github.com/gitleaks/gitleaks/pull/1943) | [#1895](https://github.com/gitleaks/gitleaks/issues/1895) and [#1859](https://github.com/gitleaks/gitleaks/issues/1859) | Potential duplicate reports of the same pre-commit hook fault |
| [#2154](https://github.com/gitleaks/gitleaks/pull/2154) | [#2152](https://github.com/gitleaks/gitleaks/issues/2152) and [#2153](https://github.com/gitleaks/gitleaks/issues/2153) | Two related feature requests, not duplicates |
| [#2199](https://github.com/gitleaks/gitleaks/pull/2199) | [#2160](https://github.com/gitleaks/gitleaks/issues/2160) and [#2121](https://github.com/gitleaks/gitleaks/issues/2121) | Same duplicate regex-work defect described by two issues |
| [#2343](https://github.com/gitleaks/gitleaks/pull/2343) | [#2120](https://github.com/gitleaks/gitleaks/issues/2120) and [#2048](https://github.com/gitleaks/gitleaks/issues/2048) | Same deprecated allowlist syntax, overlapping config scopes |

## D. Related issues that must remain distinct for now

- **Chunking and scanning false negatives:** [#1651](https://github.com/gitleaks/gitleaks/issues/1651) reports 60 Regex101 matches versus 50 from Gitleaks, increasing to 55/58/59 with larger chunks; [#2328](https://github.com/gitleaks/gitleaks/issues/2328) describes binary signature misclassification on later chunks of long first lines. These are related investigation paths, **not** proven duplicate bugs. Track their distinct reproductions and fixed expectations separately.
- **Archive behavior:** [#2010](https://github.com/gitleaks/gitleaks/issues/2010) (archive recognized by MIME but skipped under unknown extension) and [#2316](https://github.com/gitleaks/gitleaks/issues/2316) (plain text mistakenly classified as archive because of a filename substring) are inverse failure modes. They may share filename/type-dispatch code but are not interchangeable tests.
- **Staged, unreadable, and skipped files:** [#2232](https://github.com/gitleaks/gitleaks/issues/2232), [#2337](https://github.com/gitleaks/gitleaks/issues/2337), [#2316](https://github.com/gitleaks/gitleaks/issues/2316), and [#2328](https://github.com/gitleaks/gitleaks/issues/2328) all involve silent clean results after incomplete scanning, yet each arises from a different path. A common incomplete-scan result contract could be useful, but keep bug-specific fixes/tests separate.
- **Dependency upgrades:** dependency-only PRs live exclusively in [dependencies](dependencies.md), not this register; similar titles are not evidence of duplicates.

## E. Crosswalk false-positive safeguards

- The first-pass CSV parser recorded bare `#N` occurrences from foreign package changelogs and release notes. **Do not treat third-party repository issue numbers as Gitleaks duplicate relationships or fixes.** A confirmed case exists in a Cobra release-note excerpt; its PR is owned by the dependency register.
- A reference within a code block, upstream dependency changelog, or third-party URL is lower confidence than an explicit `Fixes/Closes/Resolves` statement aimed at `gitleaks/gitleaks`.
- PR references outside the open-issue snapshot have not been treated as confirmed closed issues; they require direct status verification.
- The 119 strong-signal relationships still require individual code and CI review. This initial review does not assert that all are genuine fixes.
- Recommended next crosswalk improvement: capture `evidence_source`, `confidence`, `duplicate_group_id`, and `confirmed_issue_relation` separately, without rewriting the raw source snapshots. Preserve a review override when refreshing.

## F. Disposition and future implementation tracking

All candidates remain **UNREVIEWED / NOT IMPLEMENTED** in this register. For each selected group, create a dedicated category review with:

- Reproduced behavior and independently derived expected result.
- Candidate PRs, authors, commit SHAs, and whether code changes overlap or conflict.
- Regression results on pinned images and fixtures, with clear exclusions.
- A decision: selected / duplicate implementation / deferred / rejected / blocked / failed experiment, with rationale.
- Once an actual implementation is accepted: record **RNC branch, RNC PR URL, implementation commit SHA, date, and covered upstream issue/PR IDs** in the durable tracking scheme.

The tracking field design was deliberately **not added to the source JSON/CSV/Markdown inventories yet**. That is a separate follow-up decision so source snapshots and human adoption decisions stay distinguishable.

No upstream PR was cherry-picked, no RNC security branch was modified, and no merge was performed for this assessment.

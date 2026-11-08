# 02 — Open pull requests

Snapshot: 2026-11-08T07:14:05.557Z | Source: [gitleaks/gitleaks PRs](https://github.com/gitleaks/gitleaks/pulls?q=is%3Apr+is%3Aopen)

- **252** open PRs collected; **112** are flagged by GitHub's `linked:issue` search; **140** are not so flagged.
- **144** PRs have at least one parsed issue reference (including mentions); **108** have none detected.
- These sets differ: a PR can mention an issue without having a formal link, or have GitHub linkage not visible in the parsed text.
- All PRs are **unreviewed**; source code and CI statuses have not been tested.

## First-pass categories

| Suggested category | PRs |
| --- | ---: |
| Dependencies and CI | 21 |
| Detection accuracy | 105 |
| Git and sources | 10 |
| Integrations | 9 |
| Other / needs triage | 85 |
| Reporting and output | 19 |
| Security hardening | 3 |

## Complete register

| PR | Title | Author | Created (UTC) | Changes | Formal issue link | Suggested category |
| --- | --- | --- | --- | ---: | --- | --- |
| [#2348](https://github.com/gitleaks/gitleaks/pull/2348) | Continue scanning after invalid UTF-8 bytes | blaze0089 | 2026-11-08 | 3 | Yes | Other / needs triage |
| [#2347](https://github.com/gitleaks/gitleaks/pull/2347) | feat(rules): detect SQL Server connection strings with passwords | delaysystem | 2026-11-08 | 3 | Yes | Detection accuracy |
| [#2346](https://github.com/gitleaks/gitleaks/pull/2346) | fix(detect): exclude the previous line's newline from Finding.Line | delaysystem | 2026-11-08 | 4 | Yes | Detection accuracy |
| [#2345](https://github.com/gitleaks/gitleaks/pull/2345) | fix(rules): drop dead GCPServiceAccount rule, cover GCP keys in private-key tests | delaysystem | 2026-11-08 | 3 | No | Detection accuracy |
| [#2344](https://github.com/gitleaks/gitleaks/pull/2344) | docs(cli): list all supported --platform values in help text | delaysystem | 2026-11-08 | 2 | Yes | Other / needs triage |
| [#2343](https://github.com/gitleaks/gitleaks/pull/2343) | fix(config): emit [[allowlists]] instead of deprecated [allowlist] | delaysystem | 2026-11-08 | 2 | Yes | Detection accuracy |
| [#2342](https://github.com/gitleaks/gitleaks/pull/2342) | test(rules): enable two dormant generic-api-key true positives | delaysystem | 2026-11-08 | 1 | No | Integrations |
| [#2341](https://github.com/gitleaks/gitleaks/pull/2341) | test(detect): lock in two separate findings for two private keys | delaysystem | 2026-11-08 | 1 | Yes | Detection accuracy |
| [#2340](https://github.com/gitleaks/gitleaks/pull/2340) | feat(rules): detect telegram bot tokens in Telegraf constructor calls | delaysystem | 2026-11-08 | 2 | No | Detection accuracy |
| [#2339](https://github.com/gitleaks/gitleaks/pull/2339) | fix(sources): return git exit error instead of reporting a clean scan | delaysystem | 2026-11-07 | 2 | Yes | Other / needs triage |
| [#2338](https://github.com/gitleaks/gitleaks/pull/2338) | fix(git): scan staged text marked binary by git config | VISHNU0906 | 2026-11-06 | 2 | Yes | Git and sources |
| [#2336](https://github.com/gitleaks/gitleaks/pull/2336) | build(deps): bump klauspost/compress to v1.18.7 | VISHNU0906 | 2026-11-05 | 2 | No | Dependencies and CI |
| [#2335](https://github.com/gitleaks/gitleaks/pull/2335) | fix: detect generic keys with + or / prefix | dajiaohuang | 2026-11-05 | 2 | Yes | Detection accuracy |
| [#2334](https://github.com/gitleaks/gitleaks/pull/2334) | fix(dir): respect the given path instead of silently scanning CWD (#1979) | Shanmukh-l33t | 2026-11-04 | 2 | Yes | Other / needs triage |
| [#2333](https://github.com/gitleaks/gitleaks/pull/2333) | feat: add Nuntiq API key detection | EwoutVo | 2026-10-31 | 3 | No | Detection accuracy |
| [#2332](https://github.com/gitleaks/gitleaks/pull/2332) | test: cover multiline SourceGraph token from #1697 | Tarikdxd | 2026-10-31 | 1 | No | Detection accuracy |
| [#2329](https://github.com/gitleaks/gitleaks/pull/2329) | fix(sources): check the filetype only on the first chunk | cardene777 | 2026-10-29 | 2 | Yes | Other / needs triage |
| [#2326](https://github.com/gitleaks/gitleaks/pull/2326) | docs: document JSON report to stdout and report-format/report-path dependency | bharathvemula | 2026-10-29 | 1 | Yes | Reporting and output |
| [#2325](https://github.com/gitleaks/gitleaks/pull/2325) | feat(rules): add OpenRouter API key detection rule | Tarikdxd | 2026-10-28 | 3 | Yes | Detection accuracy |
| [#2324](https://github.com/gitleaks/gitleaks/pull/2324) | fix(sources): finish the current line when no blank line is found at a chunk boundary | cardene777 | 2026-10-28 | 2 | No | Other / needs triage |
| [#2323](https://github.com/gitleaks/gitleaks/pull/2323) | build(deps): bump docker/login-action from 4.6.0 to 4.7.0 in the actions group | dependabot[bot] | 2026-10-28 | 1 | No | Dependencies and CI |
| [#2319](https://github.com/gitleaks/gitleaks/pull/2319) | fix: disable filename forwarding for the system pre-commit hook | abhinavrastogi2627 | 2026-10-24 | 3 | Yes | Git and sources |
| [#2318](https://github.com/gitleaks/gitleaks/pull/2318) | fix: skip allowlisted commits in Git.Fragments | MAHIR-BABBAR | 2026-10-23 | 2 | No | Other / needs triage |
| [#2317](https://github.com/gitleaks/gitleaks/pull/2317) | fix(sources): verify archive content before trusting filename extension | Varshith-Kali | 2026-10-22 | 2 | Yes | Other / needs triage |
| [#2315](https://github.com/gitleaks/gitleaks/pull/2315) | add LangSmith API keys | ramonpetgrave64 | 2026-10-20 | 3 | Yes | Integrations |
| [#2313](https://github.com/gitleaks/gitleaks/pull/2313) | Surface unreadable files as a partial scan instead of a clean run | Varshith-Kali | 2026-10-20 | 2 | Yes | Other / needs triage |
| [#2312](https://github.com/gitleaks/gitleaks/pull/2312) | Write report files with restrictive permissions (0600) | Varshith-Kali | 2026-10-20 | 2 | Yes | Reporting and output |
| [#2311](https://github.com/gitleaks/gitleaks/pull/2311) | enforce standard git log format to avoid custom format.pretty breaks | tonymet | 2026-10-19 | 1 | No | Reporting and output |
| [#2310](https://github.com/gitleaks/gitleaks/pull/2310) | Add new OpenRouter API key rule | Varshith-Kali | 2026-10-19 | 3 | Yes | Detection accuracy |
| [#2308](https://github.com/gitleaks/gitleaks/pull/2308) | build(deps): update dependencies to resolve CVEs (#2305, #2291) | sanjitpawar2007-hue | 2026-10-19 | 2 | Yes | Dependencies and CI |
| [#2307](https://github.com/gitleaks/gitleaks/pull/2307) | fix: preserve existing reports during report validation | jdotc1 | 2026-10-17 | 2 | Yes | Reporting and output |
| [#2306](https://github.com/gitleaks/gitleaks/pull/2306) | docs: preserve Firebase example key provenance in generated config | rotirotirafa | 2026-10-17 | 3 | Yes | Other / needs triage |
| [#2303](https://github.com/gitleaks/gitleaks/pull/2303) | feat(rules): add Mercado Pago access token rule | JoacoB-dev | 2026-10-14 | 3 | No | Detection accuracy |
| [#2301](https://github.com/gitleaks/gitleaks/pull/2301) | fix(git): detect secrets introduced in merge commits | bhaskargurram-ai | 2026-10-11 | 3 | Yes | Detection accuracy |
| [#2300](https://github.com/gitleaks/gitleaks/pull/2300) | fix: write report files with owner-only (0600) permissions | ToniAdreal | 2026-10-11 | 2 | Yes | Reporting and output |
| [#2298](https://github.com/gitleaks/gitleaks/pull/2298) | docs: replace deprecated detect wording in baseline guide | bharathvemula | 2026-10-08 | 1 | Yes | Detection accuracy |
| [#2297](https://github.com/gitleaks/gitleaks/pull/2297) | feat: support report options and redact in config file | AdamMagued | 2026-10-05 | 6 | Yes | Reporting and output |
| [#2294](https://github.com/gitleaks/gitleaks/pull/2294) | fix(config): group boolean and null alternation to anchor entire allowlist regex (fixes #2293) | i-am-paradox | 2026-10-03 | 3 | Yes | Detection accuracy |
| [#2292](https://github.com/gitleaks/gitleaks/pull/2292) | fix(git): return an error when the git command fails | Saadanjum0 | 2026-10-01 | 2 | No | Other / needs triage |
| [#2290](https://github.com/gitleaks/gitleaks/pull/2290) | chore(rules/stripe.go): add Stripe Web Hooks secrets to rule | arthurzenika | 2026-09-30 | 1 | No | Detection accuracy |
| [#2289](https://github.com/gitleaks/gitleaks/pull/2289) | build(deps): bump github.com/wasilibs/go-re2 from 1.9.0 to 1.13.0 | dependabot[bot] | 2026-09-30 | 2 | No | Dependencies and CI |
| [#2287](https://github.com/gitleaks/gitleaks/pull/2287) | docs: fix anchor for official Docker images link in README | ump45nose | 2026-09-29 | 1 | No | Dependencies and CI |
| [#2286](https://github.com/gitleaks/gitleaks/pull/2286) | Add opt-in content-specific fingerprints | breken-ai | 2026-09-28 | 13 | Yes | Other / needs triage |
| [#2284](https://github.com/gitleaks/gitleaks/pull/2284) | fix: force --no-color on git diff/log so forced color config doesn't break parsing | lekanjava | 2026-09-24 | 2 | Yes | Git and sources |
| [#2282](https://github.com/gitleaks/gitleaks/pull/2282) | docs: fix typos in CONTRIBUTING.md and README.md | haimingZZ | 2026-09-22 | 2 | No | Other / needs triage |
| [#2280](https://github.com/gitleaks/gitleaks/pull/2280) | add supabase-service-role-key and database-connection-string rules | MohamedQiqa | 2026-09-19 | 4 | No | Other / needs triage |
| [#2271](https://github.com/gitleaks/gitleaks/pull/2271) | perf(detect): reuse regex matches | dhanushs11 | 2026-09-11 | 2 | Yes | Detection accuracy |
| [#2261](https://github.com/gitleaks/gitleaks/pull/2261) | Mask secrets by rune so redaction stays valid UTF-8 | arpitjain099 | 2026-09-02 | 2 | No | Other / needs triage |
| [#2260](https://github.com/gitleaks/gitleaks/pull/2260) | fix(rules): use 'pat' keyword for airtable PAT prefilter to eliminate blind spot | jankesec | 2026-09-01 | 2 | Yes | Detection accuracy |
| [#2259](https://github.com/gitleaks/gitleaks/pull/2259) | feat(rules): add detection rules for Supabase secret API keys and personal access tokens | jankesec | 2026-09-01 | 3 | Yes | Detection accuracy |
| [#2258](https://github.com/gitleaks/gitleaks/pull/2258) | docs: fix identifier typo in README | jankesec | 2026-09-01 | 1 | Yes | Other / needs triage |
| [#2255](https://github.com/gitleaks/gitleaks/pull/2255) | fix(detect): require same line for within_columns proximity | AARUSHGAUR | 2026-08-29 | 2 | Yes | Detection accuracy |
| [#2254](https://github.com/gitleaks/gitleaks/pull/2254) | feat(rules): add Microsoft/Azure secret detection rules (#539) | AARUSHGAUR | 2026-08-29 | 3 | No | Detection accuracy |
| [#2253](https://github.com/gitleaks/gitleaks/pull/2253) | Return source read errors from DetectSource | arpitjain099 | 2026-08-27 | 2 | Yes | Other / needs triage |
| [#2252](https://github.com/gitleaks/gitleaks/pull/2252) | fix: propagate git-scan errors so a failed scan exits non-zero (gitleaks#1464) | CAOShurong | 2026-08-27 | 5 | No | Other / needs triage |
| [#2250](https://github.com/gitleaks/gitleaks/pull/2250) | build(deps): bump github.com/stretchr/testify from 1.10.0 to 1.12.1 | dependabot[bot] | 2026-08-26 | 2 | No | Dependencies and CI |
| [#2249](https://github.com/gitleaks/gitleaks/pull/2249) | fix(sources): report files that could not be read instead of exiting 0 | vaibhav8a | 2026-08-26 | 2 | Yes | Reporting and output |
| [#2248](https://github.com/gitleaks/gitleaks/pull/2248) | fix(dir): match .gitleaksignore, allowlist paths and baselines relative to the scan root | hamodywe | 2026-08-26 | 7 | No | Detection accuracy |
| [#2247](https://github.com/gitleaks/gitleaks/pull/2247) | docs: replace dead blog.gitleaks.io link in README | CAOShurong | 2026-08-25 | 1 | No | Other / needs triage |
| [#2246](https://github.com/gitleaks/gitleaks/pull/2246) | fix(sources): stop a failed diff parse from hanging or silently truncating a scan | aron-intframe | 2026-08-25 | 2 | Yes | Other / needs triage |
| [#2244](https://github.com/gitleaks/gitleaks/pull/2244) | feat: report unused .gitleaksignore entries | Tens1des | 2026-08-21 | 4 | Yes | Reporting and output |
| [#2242](https://github.com/gitleaks/gitleaks/pull/2242) | feat(rules): add ActionGate Prove Token and LangGraph API Key secret detection rules | AAH20 | 2026-08-17 | 4 | No | Detection accuracy |
| [#2241](https://github.com/gitleaks/gitleaks/pull/2241) | Fix panic on invalid regex in user config | harriiinnii | 2026-08-17 | 8 | Yes | Detection accuracy |
| [#2240](https://github.com/gitleaks/gitleaks/pull/2240) | fix(rules): detect legacy and service-account OpenAI API keys | buiducnhat | 2026-08-16 | 2 | Yes | Detection accuracy |
| [#2238](https://github.com/gitleaks/gitleaks/pull/2238) | fix(detect): count columns in codepoints, not bytes | hamodywe | 2026-08-11 | 2 | Yes | Detection accuracy |
| [#2237](https://github.com/gitleaks/gitleaks/pull/2237) | fix: prevent CSV formula injection in csv report output | bunlongheng | 2026-08-11 | 1 | No | Reporting and output |
| [#2236](https://github.com/gitleaks/gitleaks/pull/2236) | docs: clarify that git log -p scans only added lines, not pre-existing content | RajeshRajendiran | 2026-08-10 | 1 | Yes | Reporting and output |
| [#2235](https://github.com/gitleaks/gitleaks/pull/2235) | Report unreadable files as a partial scan instead of silent success | hamodywe | 2026-08-08 | 3 | Yes | Reporting and output |
| [#2233](https://github.com/gitleaks/gitleaks/pull/2233) | fix: escape SARIF artifact location URIs | CBOSSX | 2026-08-07 | 2 | Yes | Reporting and output |
| [#2230](https://github.com/gitleaks/gitleaks/pull/2230) | Count runes, not bytes, when computing entropy | arpitjain099 | 2026-08-03 | 2 | No | Detection accuracy |
| [#2229](https://github.com/gitleaks/gitleaks/pull/2229) | fix(rules): generic-api-key misses secrets after an unrelated comma-separated field | RajeshRajendiran | 2026-08-02 | 3 | Yes | Integrations |
| [#2227](https://github.com/gitleaks/gitleaks/pull/2227) | fix(sources): don't skip files on path match when allowlist condition is AND | BasitS-hash | 2026-08-01 | 2 | No | Detection accuracy |
| [#2226](https://github.com/gitleaks/gitleaks/pull/2226) | feat(rules): add Supabase secret key and personal access token rules | Yusuf-Gadelrab | 2026-07-31 | 3 | Yes | Detection accuracy |
| [#2224](https://github.com/gitleaks/gitleaks/pull/2224) | Detect basic-auth credentials in URLs | rylena | 2026-07-30 | 3 | Yes | Detection accuracy |
| [#2220](https://github.com/gitleaks/gitleaks/pull/2220) | ci: attach provenance and SBOM attestations to the published images | kobihikri | 2026-07-28 | 1 | No | Dependencies and CI |
| [#2219](https://github.com/gitleaks/gitleaks/pull/2219) | fix(rules): require the identifier to end with the keyword in hashicorp-tf-password | Eljees | 2026-07-28 | 2 | Yes | Detection accuracy |
| [#2217](https://github.com/gitleaks/gitleaks/pull/2217) | test(sources): cover file, directory walking and platform detection | Qalipso | 2026-07-28 | 3 | No | Detection accuracy |
| [#2216](https://github.com/gitleaks/gitleaks/pull/2216) | Add line-offset gitleaks:allow directives | hantaee22 | 2026-07-28 | 5 | Yes | Other / needs triage |
| [#2215](https://github.com/gitleaks/gitleaks/pull/2215) | fix: calculate columns from line starts | Eljees | 2026-07-28 | 3 | Yes | Other / needs triage |
| [#2213](https://github.com/gitleaks/gitleaks/pull/2213) | docs: list supported platform values | rylena | 2026-07-27 | 3 | Yes | Other / needs triage |
| [#2212](https://github.com/gitleaks/gitleaks/pull/2212) | docs: fix stale --max-decode-depth default in README | RajeshRajendiran | 2026-07-26 | 1 | Yes | Other / needs triage |
| [#2209](https://github.com/gitleaks/gitleaks/pull/2209) | fix: Org. repo with licence + PR from fork = missing license key error | trinathone | 2026-07-25 | 1 | Yes | Other / needs triage |
| [#2208](https://github.com/gitleaks/gitleaks/pull/2208) | fix: Decoding skips values with Unicode characters | trinathone | 2026-07-25 | 2 | Yes | Other / needs triage |
| [#2207](https://github.com/gitleaks/gitleaks/pull/2207) | fix: GitLeaks overreach flagging variable called privateKey | trinathone | 2026-07-25 | 2 | Yes | Other / needs triage |
| [#2206](https://github.com/gitleaks/gitleaks/pull/2206) | fix: Run as a non-root user in Docker | trinathone | 2026-07-25 | 1 | Yes | Dependencies and CI |
| [#2205](https://github.com/gitleaks/gitleaks/pull/2205) | fix: Passing --no-notes via --log-opts option doesn't work with CSV output | trinathone | 2026-07-25 | 3 | Yes | Reporting and output |
| [#2204](https://github.com/gitleaks/gitleaks/pull/2204) | perf(detect): reuse regex matches instead of scanning the fragment twice | rayair250-droid | 2026-07-24 | 1 | No | Detection accuracy |
| [#2203](https://github.com/gitleaks/gitleaks/pull/2203) | fix: stderr from git log no longer silently drops remaining commits | RajeshRajendiran | 2026-07-23 | 2 | Yes | Reporting and output |
| [#2202](https://github.com/gitleaks/gitleaks/pull/2202) | docs: fix baseline typo | rylena | 2026-07-23 | 1 | No | Detection accuracy |
| [#2201](https://github.com/gitleaks/gitleaks/pull/2201) | fix: withinColumns compares column numbers across different lines in composite rules | RajeshRajendiran | 2026-07-22 | 2 | Yes | Other / needs triage |
| [#2200](https://github.com/gitleaks/gitleaks/pull/2200) | fix: global commit allowlist silently bypassed due to misscoped continue | RajeshRajendiran | 2026-07-22 | 1 | Yes | Detection accuracy |
| [#2199](https://github.com/gitleaks/gitleaks/pull/2199) | fix: avoid running rule regex twice per matched fragment | RajeshRajendiran | 2026-07-22 | 1 | Yes | Detection accuracy |
| [#2195](https://github.com/gitleaks/gitleaks/pull/2195) | fix(rules): allowlist editor keybinding chords in generic-api-key | alexchen-sys | 2026-07-20 | 2 | Yes | Detection accuracy |
| [#2193](https://github.com/gitleaks/gitleaks/pull/2193) | fix(config): detect new GitHub App installation tokens (fixes #2192) | tianrking | 2026-07-16 | 2 | Yes | Detection accuracy |
| [#2187](https://github.com/gitleaks/gitleaks/pull/2187) | fix: withinColumns should not compare columns across different lines | driphtyio | 2026-07-05 | 1 | Yes | Other / needs triage |
| [#2185](https://github.com/gitleaks/gitleaks/pull/2185) | feat(rules): add Resend API key detection rule | mwoss | 2026-07-01 | 3 | No | Detection accuracy |
| [#2183](https://github.com/gitleaks/gitleaks/pull/2183) | fix: allowlist Swift private key type declarations | sanmaxdev | 2026-06-26 | 2 | Yes | Detection accuracy |
| [#2181](https://github.com/gitleaks/gitleaks/pull/2181) | feat: add detection rules for Groq and xAI API keys | PhinehasNarh | 2026-06-25 | 4 | Yes | Detection accuracy |
| [#2178](https://github.com/gitleaks/gitleaks/pull/2178) | feat(rules): add Anthropic OAuth refresh token rule (sk-ant-ort01-) | DevamShah | 2026-06-22 | 3 | No | Detection accuracy |
| [#2177](https://github.com/gitleaks/gitleaks/pull/2177) | feat(rules): add Pinecone and LangSmith API key detection rules | DevamShah | 2026-06-22 | 4 | No | Detection accuracy |
| [#2175](https://github.com/gitleaks/gitleaks/pull/2175) | fix: only apply withinColumns proximity within the same line | ardittirana | 2026-06-19 | 2 | Yes | Other / needs triage |
| [#2174](https://github.com/gitleaks/gitleaks/pull/2174) | doc: document --platform flag in helptext & readme | moutons | 2026-06-17 | 3 | Yes | Other / needs triage |
| [#2172](https://github.com/gitleaks/gitleaks/pull/2172) | Decode base64 values containing valid UTF-8 characters | arpitjain099 | 2026-06-15 | 3 | Yes | Other / needs triage |
| [#2171](https://github.com/gitleaks/gitleaks/pull/2171) | fix: use [[allowlists]] instead of deprecated [allowlist] in default config | driphtyio | 2026-06-14 | 1 | Yes | Detection accuracy |
| [#2168](https://github.com/gitleaks/gitleaks/pull/2168) | Add rule for Zhipu AI (Z.ai) API keys | moizxsec | 2026-06-11 | 3 | No | Detection accuracy |
| [#2166](https://github.com/gitleaks/gitleaks/pull/2166) | fix: global commit allowlist bypassed due to misscoped continue | ashvinctrl | 2026-06-10 | 4 | Yes | Detection accuracy |
| [#2163](https://github.com/gitleaks/gitleaks/pull/2163) | perf: reuse regex match results in detectRule() to avoid double regex execution | driphtyio | 2026-06-08 | 1 | Yes | Detection accuracy |
| [#2162](https://github.com/gitleaks/gitleaks/pull/2162) | perf: reuse regex match results in detectRule() to avoid double regex execution | driphtyio | 2026-06-08 | 1 | Yes | Detection accuracy |
| [#2161](https://github.com/gitleaks/gitleaks/pull/2161) | perf(detect): reuse match results instead of scanning the fragment twice | mazen-salah | 2026-06-04 | 1 | Yes | Detection accuracy |
| [#2159](https://github.com/gitleaks/gitleaks/pull/2159) | feat(rules): add Anthropic OAuth access and refresh token rules | mazen-salah | 2026-06-04 | 3 | Yes | Detection accuracy |
| [#2156](https://github.com/gitleaks/gitleaks/pull/2156) | chore(github): add SDK embedding issue template | MChorfa | 2026-06-03 | 1 | Yes | Integrations |
| [#2154](https://github.com/gitleaks/gitleaks/pull/2154) | feat(sdk): add embeddable scanner, canary, and vector state APIs | MChorfa | 2026-06-03 | 7 | Yes | Other / needs triage |
| [#2151](https://github.com/gitleaks/gitleaks/pull/2151) | fix(detect): remove duplicate regex call in detectRule (#2121) | wahajahmed010 | 2026-06-01 | 1 | Yes | Detection accuracy |
| [#2150](https://github.com/gitleaks/gitleaks/pull/2150) | Fix composite withinColumns cross-line matching | AlbusZhuo | 2026-06-01 | 3 | No | Other / needs triage |
| [#2149](https://github.com/gitleaks/gitleaks/pull/2149) | build(deps): bump github.com/spf13/cobra from 1.9.1 to 1.10.2 | dependabot[bot] | 2026-06-01 | 2 | No | Dependencies and CI |
| [#2147](https://github.com/gitleaks/gitleaks/pull/2147) | build(deps): bump github.com/hashicorp/go-version from 1.7.0 to 1.9.0 | dependabot[bot] | 2026-06-01 | 2 | No | Dependencies and CI |
| [#2145](https://github.com/gitleaks/gitleaks/pull/2145) | build(deps): bump github.com/charmbracelet/lipgloss from 0.5.0 to 1.1.0 | dependabot[bot] | 2026-06-01 | 2 | No | Dependencies and CI |
| [#2144](https://github.com/gitleaks/gitleaks/pull/2144) | build(deps): bump github.com/rs/zerolog from 1.33.0 to 1.35.1 | dependabot[bot] | 2026-06-01 | 2 | No | Dependencies and CI |
| [#2141](https://github.com/gitleaks/gitleaks/pull/2141) | build(deps): bump github.com/fatih/semgroup from 1.2.0 to 1.3.0 | dependabot[bot] | 2026-06-01 | 2 | No | Dependencies and CI |
| [#2139](https://github.com/gitleaks/gitleaks/pull/2139) | build(deps): bump github.com/spf13/viper from 1.19.0 to 1.21.0 | dependabot[bot] | 2026-06-01 | 2 | No | Dependencies and CI |
| [#2136](https://github.com/gitleaks/gitleaks/pull/2136) | fix: reduce generic api key false positives in JS properties | Hackerchen716 | 2026-05-30 | 2 | Yes | Integrations |
| [#2134](https://github.com/gitleaks/gitleaks/pull/2134) | chore: update logo | santoshyadavdev | 2026-05-27 | 1 | No | Other / needs triage |
| [#2132](https://github.com/gitleaks/gitleaks/pull/2132) | feat(rules): add sealed-env detection rules | davidalmeidac | 2026-05-22 | 3 | Yes | Detection accuracy |
| [#2128](https://github.com/gitleaks/gitleaks/pull/2128) | Fix quoted --log-opts parsing | puneetdixit200 | 2026-05-21 | 2 | Yes | Reporting and output |
| [#2127](https://github.com/gitleaks/gitleaks/pull/2127) | docs: improve README phrasing and fix typographical errors | StephenJarso | 2026-05-19 | 1 | No | Other / needs triage |
| [#2126](https://github.com/gitleaks/gitleaks/pull/2126) | feat(detect): add --log-ignored to surface suppressed findings | obchain | 2026-05-19 | 3 | Yes | Detection accuracy |
| [#2125](https://github.com/gitleaks/gitleaks/pull/2125) | fix(rules): allowlist keyboard chord secrets in generic-api-key | obchain | 2026-05-19 | 2 | Yes | Detection accuracy |
| [#2124](https://github.com/gitleaks/gitleaks/pull/2124) | config: switch the bundled gitleaks.toml to [[allowlists]] | c-tonneslan | 2026-05-18 | 2 | Yes | Other / needs triage |
| [#2123](https://github.com/gitleaks/gitleaks/pull/2123) | withinProximity: require same line for withinColumns to apply | c-tonneslan | 2026-05-18 | 2 | Yes | Other / needs triage |
| [#2119](https://github.com/gitleaks/gitleaks/pull/2119) | feat: add pagseguro-token and pagseguro-api-key detection rules | gugacyber | 2026-05-15 | 1 | No | Detection accuracy |
| [#2118](https://github.com/gitleaks/gitleaks/pull/2118) | feat: add Cielo payment gateway secret detection rules | gugacyber | 2026-05-15 | 1 | No | Detection accuracy |
| [#2117](https://github.com/gitleaks/gitleaks/pull/2117) | Add Mercado Pago access token and public key detection rules | gugacyber | 2026-05-14 | 1 | No | Detection accuracy |
| [#2116](https://github.com/gitleaks/gitleaks/pull/2116) | chore(gitignore): drop stale config paths | obchain | 2026-05-14 | 1 | Yes | Git and sources |
| [#2115](https://github.com/gitleaks/gitleaks/pull/2115) | chore: bump Go to 1.26.3 to fix vulnerabilities | AlonBilman | 2026-05-14 | 3 | No | Other / needs triage |
| [#2114](https://github.com/gitleaks/gitleaks/pull/2114) | feat(dir): add --source-relative-paths for git/dir fingerprint parity | obchain | 2026-05-14 | 3 | Yes | Detection accuracy |
| [#2109](https://github.com/gitleaks/gitleaks/pull/2109) | fix(rules): require sourcegraph context for bare 40-char hex tokens | obchain | 2026-05-13 | 2 | Yes | Other / needs triage |
| [#2108](https://github.com/gitleaks/gitleaks/pull/2108) | fix(rules): tighten square-access-token to cut base64 FPs | obchain | 2026-05-13 | 2 | Yes | Detection accuracy |
| [#2106](https://github.com/gitleaks/gitleaks/pull/2106) | feat: add Azure Storage account connection string detection | 54ry | 2026-05-09 | 3 | No | Detection accuracy |
| [#2105](https://github.com/gitleaks/gitleaks/pull/2105) | feat(rules): detect Claude Code OAuth tokens and MCP server credentials | ipunithgowda | 2026-05-09 | 4 | Yes | Detection accuracy |
| [#2103](https://github.com/gitleaks/gitleaks/pull/2103) | build: update Go toolchain for CVE-2025-68121 | dkarter | 2026-05-08 | 1 | No | Security hardening |
| [#2101](https://github.com/gitleaks/gitleaks/pull/2101) | feat: Add secret-value support to .gitleaksignore | rossguilbert-gov | 2026-05-05 | 4 | No | Detection accuracy |
| [#2100](https://github.com/gitleaks/gitleaks/pull/2100) | chore: Upgrade lipgloss to v2.0.3 | corani | 2026-05-05 | 4 | No | Other / needs triage |
| [#2099](https://github.com/gitleaks/gitleaks/pull/2099) | Honor max target size inside archives | alephnull-sh | 2026-05-04 | 10 | No | Other / needs triage |
| [#2098](https://github.com/gitleaks/gitleaks/pull/2098) | Emit composite rule requirements in generated config | alephnull-sh | 2026-05-04 | 4 | Yes | Detection accuracy |
| [#2097](https://github.com/gitleaks/gitleaks/pull/2097) | feat(config): read [tool.gitleaks] from pyproject.toml (#2066) | ChrisJr404 | 2026-05-03 | 4 | Yes | Other / needs triage |
| [#2096](https://github.com/gitleaks/gitleaks/pull/2096) | feat(report): add gitlab-code-quality report format (#2068) | ChrisJr404 | 2026-05-03 | 5 | Yes | Reporting and output |
| [#2095](https://github.com/gitleaks/gitleaks/pull/2095) | feat: add Claude Code session URL detection | eFAILution | 2026-05-02 | 3 | Yes | Detection accuracy |
| [#2085](https://github.com/gitleaks/gitleaks/pull/2085) | Add MongoDB Atlas and connection string detection rules | alexbevi | 2026-04-15 | 3 | No | Detection accuracy |
| [#2084](https://github.com/gitleaks/gitleaks/pull/2084) | docs: add enhanced GitHub Actions example with security best practices | vincent067 | 2026-04-15 | 2 | No | Security hardening |
| [#2083](https://github.com/gitleaks/gitleaks/pull/2083) | fix(rules): remove bare-hex alternative from sourcegraph-access-token | Smeet23 | 2026-04-10 | 2 | Yes | Detection accuracy |
| [#2082](https://github.com/gitleaks/gitleaks/pull/2082) | fix(config): return error on invalid regex instead of panicking | Smeet23 | 2026-04-10 | 5 | Yes | Detection accuracy |
| [#2081](https://github.com/gitleaks/gitleaks/pull/2081) | fix(docker): inject version via build-arg instead of git describe | Smeet23 | 2026-04-10 | 2 | Yes | Dependencies and CI |
| [#2080](https://github.com/gitleaks/gitleaks/pull/2080) | fix(sources): propagate git stderr errors instead of swallowing them | Smeet23 | 2026-04-10 | 2 | Yes | Other / needs triage |
| [#2079](https://github.com/gitleaks/gitleaks/pull/2079) | fix(config): treat empty [[rules]] blocks as no-ops | Smeet23 | 2026-04-10 | 3 | Yes | Other / needs triage |
| [#2078](https://github.com/gitleaks/gitleaks/pull/2078) | Add payment processor detection rules: Adyen, MercadoPago, Razorpay, Braintree, Midtrans, Paystack | momenbasel | 2026-04-08 | 7 | No | Detection accuracy |
| [#2077](https://github.com/gitleaks/gitleaks/pull/2077) | Add default Spring Boot YAML placeholder secret detection rule | ravikyada | 2026-04-08 | 3 | No | Detection accuracy |
| [#2072](https://github.com/gitleaks/gitleaks/pull/2072) | Add Azure Storage Account Key detection rule | lumir-sliva | 2026-03-31 | 3 | No | Detection accuracy |
| [#2061](https://github.com/gitleaks/gitleaks/pull/2061) | feat(dir): add optional --gitignore support for dir scans | shivsegv | 2026-03-19 | 6 | Yes | Git and sources |
| [#2056](https://github.com/gitleaks/gitleaks/pull/2056) | feat: add MCP server command | raajheshkannaa | 2026-03-14 | 3 | Yes | Integrations |
| [#2053](https://github.com/gitleaks/gitleaks/pull/2053) | fix: use [[allowlists]] in generated gitleaks.toml | akk85 | 2026-03-10 | 2 | Yes | Other / needs triage |
| [#2049](https://github.com/gitleaks/gitleaks/pull/2049) | feat: update gitleaks.toml to use new allowlists | dariozachow | 2026-03-05 | 2 | Yes | Other / needs triage |
| [#2041](https://github.com/gitleaks/gitleaks/pull/2041) | feat: add Replicate API token detection rule | clarabennettdev | 2026-02-22 | 3 | No | Detection accuracy |
| [#2039](https://github.com/gitleaks/gitleaks/pull/2039) | fix: Add pass_filenames: false to gitleaks-system | maledb404 | 2026-02-11 | 1 | No | Other / needs triage |
| [#2036](https://github.com/gitleaks/gitleaks/pull/2036) | Fix false-positive SAST finding by replacing realistic Facebook access token in test config | Tanmaykaturi | 2026-02-05 | 1 | No | Detection accuracy |
| [#2035](https://github.com/gitleaks/gitleaks/pull/2035) | chore(go.mod): update go version and dependencies | sgaist | 2026-02-05 | 2 | Yes | Dependencies and CI |
| [#2032](https://github.com/gitleaks/gitleaks/pull/2032) | feat: add OpenRouter API key detection | AdrianDeAnda | 2026-01-30 | 3 | No | Detection accuracy |
| [#2030](https://github.com/gitleaks/gitleaks/pull/2030) | refactor(config): adjust to new glpat routable token | max-wittig | 2026-01-27 | 3 | No | Detection accuracy |
| [#2025](https://github.com/gitleaks/gitleaks/pull/2025) | fix: stop passing filenames for gitleaks-system | JanKleine | 2026-01-12 | 1 | No | Other / needs triage |
| [#2024](https://github.com/gitleaks/gitleaks/pull/2024) | fix RequiredRules and SkipReport not properly applied when extending configuration | tralph3 | 2026-01-10 | 1 | No | Other / needs triage |
| [#2023](https://github.com/gitleaks/gitleaks/pull/2023) | Reducing regex operations | sissl0 | 2026-01-10 | 2 | No | Detection accuracy |
| [#2008](https://github.com/gitleaks/gitleaks/pull/2008) | docs(README): add mise alternative method installation | jylenhof | 2025-12-19 | 1 | No | Other / needs triage |
| [#2007](https://github.com/gitleaks/gitleaks/pull/2007) | add: mongodb connection string rule | sravan-pv | 2025-12-16 | 3 | No | Detection accuracy |
| [#2001](https://github.com/gitleaks/gitleaks/pull/2001) | feat(config): add Nigerian fintech & betting secret rules by @Lloydcoder | LloydCoder | 2025-12-06 | 18 | No | Detection accuracy |
| [#1996](https://github.com/gitleaks/gitleaks/pull/1996) | feat: add Django insecure keys detection | k0ral | 2025-11-28 | 3 | No | Detection accuracy |
| [#1992](https://github.com/gitleaks/gitleaks/pull/1992) | Add URL & JSON sources | bplaxco | 2025-11-19 | 20 | No | Reporting and output |
| [#1991](https://github.com/gitleaks/gitleaks/pull/1991) | Update Dockerfile | mtcolman | 2025-11-18 | 1 | No | Other / needs triage |
| [#1960](https://github.com/gitleaks/gitleaks/pull/1960) | Upgrade Go and Alpine versions in Dockerfile | Jiri-Stary | 2025-10-17 | 1 | No | Dependencies and CI |
| [#1959](https://github.com/gitleaks/gitleaks/pull/1959) | Update go to 1.24.9 in go.mod | micheltuu | 2025-10-15 | 1 | No | Dependencies and CI |
| [#1954](https://github.com/gitleaks/gitleaks/pull/1954) | Added `&` in secrets suffix boundary when key is in URL | skj-skj | 2025-10-14 | 3 | No | Other / needs triage |
| [#1953](https://github.com/gitleaks/gitleaks/pull/1953) | Fixed k8s secrets detection when kube yml have multiple secrets | skj-skj | 2025-10-09 | 2 | No | Detection accuracy |
| [#1951](https://github.com/gitleaks/gitleaks/pull/1951) | sources: replace filetype with mimetype and add explicit human-readable filter | mikelolasagasti | 2025-10-02 | 4 | No | Other / needs triage |
| [#1943](https://github.com/gitleaks/gitleaks/pull/1943) | fix(pre-commit): don't pass filenames to system hook | pcrockett | 2025-09-17 | 1 | Yes | Git and sources |
| [#1940](https://github.com/gitleaks/gitleaks/pull/1940) | create composite rule error | moamlrh | 2025-09-05 | 130 | No | Detection accuracy |
| [#1936](https://github.com/gitleaks/gitleaks/pull/1936) | Allow multiple sources for `dir` command | maxvt | 2025-08-27 | 2 | No | Other / needs triage |
| [#1930](https://github.com/gitleaks/gitleaks/pull/1930) | Update and re-add `GCPServiceAccount` rule | jrdnbradford | 2025-08-11 | 4 | Yes | Detection accuracy |
| [#1918](https://github.com/gitleaks/gitleaks/pull/1918) | Fix pre-commit hook by use git-agnostic file-based checking | Holzhaus | 2025-07-18 | 1 | Yes | Git and sources |
| [#1915](https://github.com/gitleaks/gitleaks/pull/1915) | ci: add GitHub Actions workflow to run Gitleaks on push/PR | SGreeshma | 2025-07-17 | 1 | No | Dependencies and CI |
| [#1900](https://github.com/gitleaks/gitleaks/pull/1900) | Update profile.sh | vignesh1507 | 2025-06-27 | 1 | No | Other / needs triage |
| [#1899](https://github.com/gitleaks/gitleaks/pull/1899) | Update pre-commit.py | vignesh1507 | 2025-06-27 | 1 | No | Git and sources |
| [#1896](https://github.com/gitleaks/gitleaks/pull/1896) | Added functionality for using remote configs hosted through HTTP/-s during the runtime | TheArqsz | 2025-06-21 | 2 | No | Other / needs triage |
| [#1874](https://github.com/gitleaks/gitleaks/pull/1874) | Experiment: allowlist expressions using CEL | rgmz | 2025-05-25 | 18 | No | Detection accuracy |
| [#1871](https://github.com/gitleaks/gitleaks/pull/1871) | feat(config): added gitleaksGlobIgnore for *very* basic glob ignore pattern matching | Scriptornix | 2025-05-20 | 3 | Yes | Other / needs triage |
| [#1867](https://github.com/gitleaks/gitleaks/pull/1867) | Tweak `generic-api-key` false-positives from aws-sdk-go-v2 | rgmz | 2025-05-19 | 3 | No | Integrations |
| [#1852](https://github.com/gitleaks/gitleaks/pull/1852) | docs(security): add initial security policy | janderssonse | 2025-05-11 | 1 | No | Security hardening |
| [#1841](https://github.com/gitleaks/gitleaks/pull/1841) | simple archive handling | zricethezav | 2025-04-28 | 15 | No | Other / needs triage |
| [#1829](https://github.com/gitleaks/gitleaks/pull/1829) | feat(rules): Add new detection rules for Python Jose Jwt Hardcoded Key | AlexandreEXFO | 2025-04-15 | 3 | No | Detection accuracy |
| [#1824](https://github.com/gitleaks/gitleaks/pull/1824) | Introduce context propagation for git operations | sirakav | 2025-04-13 | 6 | No | Other / needs triage |
| [#1823](https://github.com/gitleaks/gitleaks/pull/1823) | fix(adobe-client-secret): adobe client secret regex | hd-deman | 2025-04-11 | 2 | No | Detection accuracy |
| [#1809](https://github.com/gitleaks/gitleaks/pull/1809) | Add Azure DevOps SSH platform detection | JoostVoskuil | 2025-04-01 | 2 | No | Detection accuracy |
| [#1745](https://github.com/gitleaks/gitleaks/pull/1745) | Fix typos | jessp01 | 2025-02-04 | 4 | No | Other / needs triage |
| [#1724](https://github.com/gitleaks/gitleaks/pull/1724) | Experiment: create an 'experimental' config | rgmz | 2025-01-23 | 6 | No | Other / needs triage |
| [#1683](https://github.com/gitleaks/gitleaks/pull/1683) | Add warning if rule doesn't have keywords (#1675) | vk-NEU7 | 2024-12-31 | 1 | No | Detection accuracy |
| [#1661](https://github.com/gitleaks/gitleaks/pull/1661) | chore: release npm | ianzone | 2024-11-28 | 8 | No | Dependencies and CI |
| [#1654](https://github.com/gitleaks/gitleaks/pull/1654) | Create command to print config | rgmz | 2024-11-23 | 5 | No | Other / needs triage |
| [#1639](https://github.com/gitleaks/gitleaks/pull/1639) | Fix duplicate findings when scanning directory with path only rules | elchocarrero | 2024-11-14 | 4 | Yes | Other / needs triage |
| [#1630](https://github.com/gitleaks/gitleaks/pull/1630) | cmd: adding config and fetch commands to load custom config from remote URL | cduggn | 2024-11-08 | 6 | Yes | Other / needs triage |
| [#1594](https://github.com/gitleaks/gitleaks/pull/1594) | Tweak private key regex to reduce false positives | rgmz | 2024-10-22 | 3 | No | Detection accuracy |
| [#1581](https://github.com/gitleaks/gitleaks/pull/1581) | Update global allowlist to ignore alpha or numeric sequences | rgmz | 2024-10-19 | 10 | No | Detection accuracy |
| [#1568](https://github.com/gitleaks/gitleaks/pull/1568) | Detect Kuberenetes secret secretdata | rgmz | 2024-10-14 | 3 | No | Detection accuracy |
| [#1566](https://github.com/gitleaks/gitleaks/pull/1566) | WIP: use secret hash as fingerprint | rgmz | 2024-10-14 | 9 | No | Detection accuracy |
| [#1539](https://github.com/gitleaks/gitleaks/pull/1539) | Warn when a file cannot be scanned due to Git LFS | rgmz | 2024-09-30 | 1 | No | Git and sources |
| [#1534](https://github.com/gitleaks/gitleaks/pull/1534) | Ignore leading newlines | rgmz | 2024-09-27 | 2 | No | Other / needs triage |
| [#1500](https://github.com/gitleaks/gitleaks/pull/1500) | Catch common pattern errors with `\b` | rgmz | 2024-09-09 | 13 | No | Other / needs triage |
| [#1499](https://github.com/gitleaks/gitleaks/pull/1499) | Experiment: run validation as tests | rgmz | 2024-09-08 | 8 | No | Other / needs triage |
| [#1497](https://github.com/gitleaks/gitleaks/pull/1497) | Add new SaladCloud API keys rule | seniorquico | 2024-09-07 | 3 | No | Detection accuracy |
| [#1495](https://github.com/gitleaks/gitleaks/pull/1495) | Add Onfido API Token detector to recognize this type of secrets. | lucasan1 | 2024-09-02 | 3 | No | Detection accuracy |
| [#1485](https://github.com/gitleaks/gitleaks/pull/1485) | Add install.sh bash script installation option | suhovius | 2024-08-25 | 2 | No | Other / needs triage |
| [#1435](https://github.com/gitleaks/gitleaks/pull/1435) | fix#1432: made changes to the reading logic | sapienhwaker | 2024-06-30 | 2 | No | Other / needs triage |
| [#1394](https://github.com/gitleaks/gitleaks/pull/1394) | Create mssql.go | fhverga | 2024-04-18 | 1 | No | Other / needs triage |
| [#1373](https://github.com/gitleaks/gitleaks/pull/1373) | Add PuttyPrivateKey rule and allowlist for raw content | baruchiro | 2024-03-14 | 9 | No | Detection accuracy |
| [#1358](https://github.com/gitleaks/gitleaks/pull/1358) | fix: reduce false positives by hashicorp | baruchiro | 2024-02-29 | 4 | No | Other / needs triage |
| [#1356](https://github.com/gitleaks/gitleaks/pull/1356) | feat: add aws-secret-key rule | baruchiro | 2024-02-28 | 3 | No | Detection accuracy |
| [#1350](https://github.com/gitleaks/gitleaks/pull/1350) | fix: secrets found with quotas | baruchiro | 2024-02-20 | 10 | No | Other / needs triage |
| [#1335](https://github.com/gitleaks/gitleaks/pull/1335) | feat: Detect Chargebee Access Tokens | bharathvaj-ganesan | 2024-02-05 | 4 | No | Detection accuracy |
| [#1327](https://github.com/gitleaks/gitleaks/pull/1327) | add support for Azure Data Factory SHIR Key | RachidAZ | 2024-01-11 | 3 | No | Other / needs triage |
| [#1314](https://github.com/gitleaks/gitleaks/pull/1314) | chore: Fix Contributing Link in README | riley-kohler | 2023-12-09 | 1 | No | Other / needs triage |
| [#1296](https://github.com/gitleaks/gitleaks/pull/1296) | zricethezav -> gitleaks in README | jdanyow | 2023-11-21 | 1 | No | Other / needs triage |
| [#1291](https://github.com/gitleaks/gitleaks/pull/1291) | 1079 | StephenHogg | 2023-11-16 | 4 | Yes | Other / needs triage |
| [#1266](https://github.com/gitleaks/gitleaks/pull/1266) | chore: refactor regex for redundant pipe | baruchiro | 2023-09-11 | 2 | No | Detection accuracy |
| [#1265](https://github.com/gitleaks/gitleaks/pull/1265) | feat: catch secret in XML child | baruchiro | 2023-09-11 | 3 | No | Detection accuracy |
| [#1264](https://github.com/gitleaks/gitleaks/pull/1264) | Multiple files source | rajatmehra2307 | 2023-09-06 | 5 | No | Other / needs triage |
| [#1249](https://github.com/gitleaks/gitleaks/pull/1249) | Support reading .gitleaksignore using git show | savely-krasovsky | 2023-08-11 | 75 | No | Other / needs triage |
| [#1232](https://github.com/gitleaks/gitleaks/pull/1232) | Support multiple report formats in a single execution | Coruscant11 | 2023-07-15 | 4 | Yes | Reporting and output |
| [#1217](https://github.com/gitleaks/gitleaks/pull/1217) | Sign container using Cosign | becojo | 2023-07-01 | 1 | No | Other / needs triage |
| [#1162](https://github.com/gitleaks/gitleaks/pull/1162) | Add hash of a secret support | mn3m0nic | 2023-04-25 | 14 | No | Detection accuracy |
| [#1130](https://github.com/gitleaks/gitleaks/pull/1130) | Update Generic API Key rule to detect secrets with special characters | SonaliRathore | 2023-03-24 | 2 | No | Detection accuracy |
| [#1103](https://github.com/gitleaks/gitleaks/pull/1103) | fix(sec): upgrade gopkg.in/yaml.v3 to 3.0.0 | chncaption | 2023-02-19 | 2 | No | Other / needs triage |
| [#1099](https://github.com/gitleaks/gitleaks/pull/1099) | Fix: go.mod version changed back to 1.18 | Harshetjain666 | 2023-02-12 | 1 | No | Dependencies and CI |
| [#1086](https://github.com/gitleaks/gitleaks/pull/1086) | Add Plivo Auth ID rule | michaelricordeau | 2023-01-30 | 5 | No | Detection accuracy |
| [#1071](https://github.com/gitleaks/gitleaks/pull/1071) | add rule for checking nuget.config passwords | jankonikola93 | 2023-01-11 | 3 | No | Detection accuracy |
| [#1061](https://github.com/gitleaks/gitleaks/pull/1061) | feat(config): extends generic-api-key to support more special characters | cheesi | 2022-12-19 | 2 | No | Integrations |
| [#1060](https://github.com/gitleaks/gitleaks/pull/1060) | Add expansion of the SAST abbreviation | it-praktyk | 2022-12-19 | 1 | No | Other / needs triage |
| [#1046](https://github.com/gitleaks/gitleaks/pull/1046) | Add Zapier webhook regex | karankohli-cf | 2022-11-28 | 1 | No | Detection accuracy |
| [#1043](https://github.com/gitleaks/gitleaks/pull/1043) | + non-encrypted prv key in non-pem format(no -----begin header) check pt2 | icecream-monster | 2022-11-24 | 3 | No | Reporting and output |
| [#968](https://github.com/gitleaks/gitleaks/pull/968) | New rule sshpass detect | ivaldivieso | 2022-08-26 | 1 | No | Detection accuracy |
| [#958](https://github.com/gitleaks/gitleaks/pull/958) | better explanation for [extend] feature | YogevBokobza | 2022-08-18 | 1 | No | Other / needs triage |
| [#955](https://github.com/gitleaks/gitleaks/pull/955) | feat: adding support for detecting checked in OpenStack RC files. | srelf-eschercloud | 2022-08-17 | 3 | No | Other / needs triage |
| [#774](https://github.com/gitleaks/gitleaks/pull/774) | Add segment public api token rule | salolivares | 2022-01-19 | 3 | No | Detection accuracy |
| [#652](https://github.com/gitleaks/gitleaks/pull/652) | Minimum go version should be enforced when installing gitleaks with go get command | weineran | 2021-10-01 | 4 | No | Other / needs triage |
| [#643](https://github.com/gitleaks/gitleaks/pull/643) | Pull request to detect API keys,bearer token and JWT token hardcoded in gitlab/github ci file | ds50421 | 2021-09-24 | 1 | No | Detection accuracy |
| [#637](https://github.com/gitleaks/gitleaks/pull/637) | Add specific permissions to workflows under .github/workflows | varunsh-coder | 2021-09-11 | 1 | No | Integrations |
| [#555](https://github.com/gitleaks/gitleaks/pull/555) | Add file_hash to commit | grogancolin | 2021-04-25 | 4 | No | Git and sources |

## Exclusions

**108 PRs without a detected issue number or URL:** [#2342](https://github.com/gitleaks/gitleaks/pull/2342), [#2340](https://github.com/gitleaks/gitleaks/pull/2340), [#2333](https://github.com/gitleaks/gitleaks/pull/2333), [#2318](https://github.com/gitleaks/gitleaks/pull/2318), [#2311](https://github.com/gitleaks/gitleaks/pull/2311), [#2303](https://github.com/gitleaks/gitleaks/pull/2303), [#2290](https://github.com/gitleaks/gitleaks/pull/2290), [#2289](https://github.com/gitleaks/gitleaks/pull/2289), [#2287](https://github.com/gitleaks/gitleaks/pull/2287), [#2282](https://github.com/gitleaks/gitleaks/pull/2282), [#2261](https://github.com/gitleaks/gitleaks/pull/2261), [#2248](https://github.com/gitleaks/gitleaks/pull/2248), [#2242](https://github.com/gitleaks/gitleaks/pull/2242), [#2237](https://github.com/gitleaks/gitleaks/pull/2237), [#2230](https://github.com/gitleaks/gitleaks/pull/2230), [#2227](https://github.com/gitleaks/gitleaks/pull/2227), [#2220](https://github.com/gitleaks/gitleaks/pull/2220), [#2217](https://github.com/gitleaks/gitleaks/pull/2217), [#2202](https://github.com/gitleaks/gitleaks/pull/2202), [#2185](https://github.com/gitleaks/gitleaks/pull/2185), [#2177](https://github.com/gitleaks/gitleaks/pull/2177), [#2168](https://github.com/gitleaks/gitleaks/pull/2168), [#2150](https://github.com/gitleaks/gitleaks/pull/2150), [#2147](https://github.com/gitleaks/gitleaks/pull/2147), [#2145](https://github.com/gitleaks/gitleaks/pull/2145), [#2141](https://github.com/gitleaks/gitleaks/pull/2141), [#2134](https://github.com/gitleaks/gitleaks/pull/2134), [#2127](https://github.com/gitleaks/gitleaks/pull/2127), [#2119](https://github.com/gitleaks/gitleaks/pull/2119), [#2118](https://github.com/gitleaks/gitleaks/pull/2118), [#2117](https://github.com/gitleaks/gitleaks/pull/2117), [#2115](https://github.com/gitleaks/gitleaks/pull/2115), [#2103](https://github.com/gitleaks/gitleaks/pull/2103), [#2101](https://github.com/gitleaks/gitleaks/pull/2101), [#2100](https://github.com/gitleaks/gitleaks/pull/2100), [#2099](https://github.com/gitleaks/gitleaks/pull/2099), [#2085](https://github.com/gitleaks/gitleaks/pull/2085), [#2084](https://github.com/gitleaks/gitleaks/pull/2084), [#2078](https://github.com/gitleaks/gitleaks/pull/2078), [#2041](https://github.com/gitleaks/gitleaks/pull/2041), [#2039](https://github.com/gitleaks/gitleaks/pull/2039), [#2036](https://github.com/gitleaks/gitleaks/pull/2036), [#2032](https://github.com/gitleaks/gitleaks/pull/2032), [#2030](https://github.com/gitleaks/gitleaks/pull/2030), [#2025](https://github.com/gitleaks/gitleaks/pull/2025), [#2024](https://github.com/gitleaks/gitleaks/pull/2024), [#2023](https://github.com/gitleaks/gitleaks/pull/2023), [#2008](https://github.com/gitleaks/gitleaks/pull/2008), [#2007](https://github.com/gitleaks/gitleaks/pull/2007), [#2001](https://github.com/gitleaks/gitleaks/pull/2001), [#1996](https://github.com/gitleaks/gitleaks/pull/1996), [#1992](https://github.com/gitleaks/gitleaks/pull/1992), [#1991](https://github.com/gitleaks/gitleaks/pull/1991), [#1960](https://github.com/gitleaks/gitleaks/pull/1960), [#1959](https://github.com/gitleaks/gitleaks/pull/1959), [#1951](https://github.com/gitleaks/gitleaks/pull/1951), [#1915](https://github.com/gitleaks/gitleaks/pull/1915), [#1900](https://github.com/gitleaks/gitleaks/pull/1900), [#1899](https://github.com/gitleaks/gitleaks/pull/1899), [#1896](https://github.com/gitleaks/gitleaks/pull/1896), [#1874](https://github.com/gitleaks/gitleaks/pull/1874), [#1867](https://github.com/gitleaks/gitleaks/pull/1867), [#1852](https://github.com/gitleaks/gitleaks/pull/1852), [#1841](https://github.com/gitleaks/gitleaks/pull/1841), [#1829](https://github.com/gitleaks/gitleaks/pull/1829), [#1824](https://github.com/gitleaks/gitleaks/pull/1824), [#1823](https://github.com/gitleaks/gitleaks/pull/1823), [#1809](https://github.com/gitleaks/gitleaks/pull/1809), [#1745](https://github.com/gitleaks/gitleaks/pull/1745), [#1661](https://github.com/gitleaks/gitleaks/pull/1661), [#1654](https://github.com/gitleaks/gitleaks/pull/1654), [#1594](https://github.com/gitleaks/gitleaks/pull/1594), [#1581](https://github.com/gitleaks/gitleaks/pull/1581), [#1534](https://github.com/gitleaks/gitleaks/pull/1534), [#1500](https://github.com/gitleaks/gitleaks/pull/1500), [#1499](https://github.com/gitleaks/gitleaks/pull/1499), [#1497](https://github.com/gitleaks/gitleaks/pull/1497), [#1495](https://github.com/gitleaks/gitleaks/pull/1495), [#1485](https://github.com/gitleaks/gitleaks/pull/1485), [#1394](https://github.com/gitleaks/gitleaks/pull/1394), [#1373](https://github.com/gitleaks/gitleaks/pull/1373), [#1358](https://github.com/gitleaks/gitleaks/pull/1358), [#1350](https://github.com/gitleaks/gitleaks/pull/1350), [#1335](https://github.com/gitleaks/gitleaks/pull/1335), [#1327](https://github.com/gitleaks/gitleaks/pull/1327), [#1314](https://github.com/gitleaks/gitleaks/pull/1314), [#1296](https://github.com/gitleaks/gitleaks/pull/1296), [#1266](https://github.com/gitleaks/gitleaks/pull/1266), [#1265](https://github.com/gitleaks/gitleaks/pull/1265), [#1249](https://github.com/gitleaks/gitleaks/pull/1249), [#1217](https://github.com/gitleaks/gitleaks/pull/1217), [#1162](https://github.com/gitleaks/gitleaks/pull/1162), [#1130](https://github.com/gitleaks/gitleaks/pull/1130), [#1103](https://github.com/gitleaks/gitleaks/pull/1103), [#1099](https://github.com/gitleaks/gitleaks/pull/1099), [#1086](https://github.com/gitleaks/gitleaks/pull/1086), [#1071](https://github.com/gitleaks/gitleaks/pull/1071), [#1060](https://github.com/gitleaks/gitleaks/pull/1060), [#1046](https://github.com/gitleaks/gitleaks/pull/1046), [#1043](https://github.com/gitleaks/gitleaks/pull/1043), [#968](https://github.com/gitleaks/gitleaks/pull/968), [#958](https://github.com/gitleaks/gitleaks/pull/958), [#955](https://github.com/gitleaks/gitleaks/pull/955), [#774](https://github.com/gitleaks/gitleaks/pull/774), [#652](https://github.com/gitleaks/gitleaks/pull/652), [#643](https://github.com/gitleaks/gitleaks/pull/643), [#637](https://github.com/gitleaks/gitleaks/pull/637), [#555](https://github.com/gitleaks/gitleaks/pull/555).

These are discovery candidates, including one-off features, maintenance, dependencies, and fixes without explicit issue references. They are **not** automatically deemed unrelated to an issue.

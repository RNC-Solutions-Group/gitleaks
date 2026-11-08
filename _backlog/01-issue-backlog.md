# 01 — Open issue backlog

Snapshot: 2026-11-08T07:14:05.557Z | Source: [gitleaks/gitleaks issues](https://github.com/gitleaks/gitleaks/issues?q=is%3Aissue+is%3Aopen)

**RNC triage overlay — full routing pass (2026-11-08):** A ~~struck-through upstream link~~ means that issue or PR has been routed to exactly one [category register](categories/README.md), including the RNC-closed ledger. This changes **display only**: the source snapshot counts and upstream GitHub state are unchanged. Raw JSON/CSV exports remain untouched.

- **285** open issues collected; **114** have at least one detected relation to an open PR; **171** have no detected open-PR relation.
- All issue records were **unreviewed at snapshot collection**; subsequent RNC dispositions are in the category registers. “No detected relation” does not establish that an issue is unfixed.
- Descriptions were used for cross-reference extraction but excluded from the public metadata export.

## First-pass categories

These counts use **title-only keyword rules**, with one suggested category per issue. They are routing hints, not security severity or confirmed status.

| Suggested category | Issues |
| --- | ---: |
| Dependencies and CI | 9 |
| Detection accuracy | 97 |
| Git and sources | 13 |
| Integrations | 23 |
| Other / needs triage | 109 |
| Performance and stability | 3 |
| Reporting and output | 24 |
| Security hardening | 7 |

## Complete register

| Issue | Title | Open-PR references | Category |
| --- | --- | ---: | --- |
| ~~[#2337](https://github.com/gitleaks/gitleaks/issues/2337)~~ | `gitleaks git --pre-commit --staged` silently skips staged text files that git treats as binary (`core.bigFileThreshold`, `-diff` attribute) | 4 | Git and sources |
| ~~[#2330](https://github.com/gitleaks/gitleaks/issues/2330)~~ | fix below CVEs | 1 | Other / needs triage |
| ~~[#2328](https://github.com/gitleaks/gitleaks/issues/2328)~~ | Secrets are missed when a chunk of a long first line starts with a binary signature | 2 | Other / needs triage |
| ~~[#2327](https://github.com/gitleaks/gitleaks/issues/2327)~~ | generic-api-key misses base64 secrets that start with + or / | 2 | Integrations |
| ~~[#2322](https://github.com/gitleaks/gitleaks/issues/2322)~~ | ghghj | 0 | Other / needs triage |
| ~~[#2320](https://github.com/gitleaks/gitleaks/issues/2320)~~ | regexTarget = "line" allowlists see the previous line's newline, so ^-anchored regexes never match after the first line | 1 | Other / needs triage |
| ~~[#2316](https://github.com/gitleaks/gitleaks/issues/2316)~~ | A plain-text file is silently skipped ("scanned ~0 bytes", exit 0) when its name contains an archive extension such as .7z, .xz, .zip or .br | 2 | Other / needs triage |
| ~~[#2314](https://github.com/gitleaks/gitleaks/issues/2314)~~ | feat: add LangSmith API keys detection | 1 | Detection accuracy |
| ~~[#2309](https://github.com/gitleaks/gitleaks/issues/2309)~~ | Add OpenRouter API key rule (sk-or-v1-) | 2 | Detection accuracy |
| ~~[#2305](https://github.com/gitleaks/gitleaks/issues/2305)~~ | Please Resolve CVEs in 8.30.1 | 2 | Other / needs triage |
| ~~[#2304](https://github.com/gitleaks/gitleaks/issues/2304)~~ | Existing report is removed before invalid report format is rejected | 2 | Reporting and output |
| ~~[#2302](https://github.com/gitleaks/gitleaks/issues/2302)~~ | gitleaks-system pre-commit hook missing pass_filenames: false | 1 | Git and sources |
| ~~[#2299](https://github.com/gitleaks/gitleaks/issues/2299)~~ | Report file containing discovered secrets is written with default, often world-readable permissions | 3 | Security hardening |
| ~~[#2296](https://github.com/gitleaks/gitleaks/issues/2296)~~ | Feature request: allow `--redact` and other report options in `.gitleaks.toml` | 1 | Reporting and output |
| ~~[#2295](https://github.com/gitleaks/gitleaks/issues/2295)~~ | Preserve Firebase example-key provenance in the generated gcp-api-key allowlist | 1 | Detection accuracy |
| ~~[#2293](https://github.com/gitleaks/gitleaks/issues/2293)~~ | Default global boolean/null allowlist overmatches substrings, including custom rules | 2 | Detection accuracy |
| ~~[#2291](https://github.com/gitleaks/gitleaks/issues/2291)~~ | Go CVE Vulnerabilities | 1 | Security hardening |
| ~~[#2281](https://github.com/gitleaks/gitleaks/issues/2281)~~ | gitleaks git --pre-commit --staged does not detect secrets in staged files | 1 | Detection accuracy |
| ~~[#2279](https://github.com/gitleaks/gitleaks/issues/2279)~~ | Rule request: detect database connection strings with inline credentials (postgres/mysql/mongodb) | 2 | Detection accuracy |
| ~~[#2278](https://github.com/gitleaks/gitleaks/issues/2278)~~ | Rule request: detect Supabase `service_role` JWTs (distinct from safe `anon` keys) | 1 | Detection accuracy |
| ~~[#2257](https://github.com/gitleaks/gitleaks/issues/2257)~~ | airtable-personnal-access-token never matches a real token unless the word "airtable" is also nearby (keyword prefilter blind spot) | 1 | Detection accuracy |
| ~~[#2256](https://github.com/gitleaks/gitleaks/issues/2256)~~ | Typo in README: identiifer should be identifier | 1 | Other / needs triage |
| ~~[#2239](https://github.com/gitleaks/gitleaks/issues/2239)~~ | .gitleaksignore fingerprint (file:rule:line) matches by position, not secret content | 1 | Detection accuracy |
| ~~[#2232](https://github.com/gitleaks/gitleaks/issues/2232)~~ | Unreadable files are silently skipped and the run still reports "no leaks found" with exit 0 | 3 | Other / needs triage |
| ~~[#2225](https://github.com/gitleaks/gitleaks/issues/2225)~~ | Add detection rules for Supabase secret API keys and personal access tokens | 3 | Detection accuracy |
| ~~[#2223](https://github.com/gitleaks/gitleaks/issues/2223)~~ | generic-api-key: false negative when an unrelated key precedes the credential in single-line JSON | 1 | Detection accuracy |
| ~~[#2218](https://github.com/gitleaks/gitleaks/issues/2218)~~ | Pre-receive hook docs: --log-opts misses pre-existing secrets in changed files | 1 | Reporting and output |
| ~~[#2214](https://github.com/gitleaks/gitleaks/issues/2214)~~ | http Basic-Auth secret not detected | 1 | Detection accuracy |
| ~~[#2196](https://github.com/gitleaks/gitleaks/issues/2196)~~ | anthropic-api-key default rule silently fails to fire on macOS ARM64 (Homebrew build), works on Linux/aarch64 | 0 | Detection accuracy |
| ~~[#2192](https://github.com/gitleaks/gitleaks/issues/2192)~~ | GitHub App format changing | 1 | Reporting and output |
| ~~[#2188](https://github.com/gitleaks/gitleaks/issues/2188)~~ | [Notice] FYI I submitted a schema for gitleaks to schema store | 0 | Other / needs triage |
| ~~[#2186](https://github.com/gitleaks/gitleaks/issues/2186)~~ | Config file over web | 0 | Other / needs triage |
| ~~[#2182](https://github.com/gitleaks/gitleaks/issues/2182)~~ | GitLeaks overreach flagging variable called "privateKey" | 2 | Other / needs triage |
| ~~[#2180](https://github.com/gitleaks/gitleaks/issues/2180)~~ | feat: add detection rules for Groq and xAI API keys | 1 | Detection accuracy |
| ~~[#2176](https://github.com/gitleaks/gitleaks/issues/2176)~~ | Run as a non-root user in Docker | 1 | Dependencies and CI |
| ~~[#2173](https://github.com/gitleaks/gitleaks/issues/2173)~~ | --platform enumeration insufficiently documented | 3 | Other / needs triage |
| ~~[#2165](https://github.com/gitleaks/gitleaks/issues/2165)~~ | fix: global commit allowlist silently bypassed due to misscoped continue | 2 | Detection accuracy |
| ~~[#2160](https://github.com/gitleaks/gitleaks/issues/2160)~~ | detectRule runs the rule regex over each fragment twice | 2 | Detection accuracy |
| ~~[#2158](https://github.com/gitleaks/gitleaks/issues/2158)~~ | Add detection for Anthropic OAuth tokens (sk-ant-oat01-, sk-ant-ort01-) | 1 | Detection accuracy |
| ~~[#2155](https://github.com/gitleaks/gitleaks/issues/2155)~~ | Add dedicated SDK embedding issue template | 1 | Other / needs triage |
| ~~[#2153](https://github.com/gitleaks/gitleaks/issues/2153)~~ | SDK canary exemptions and vector-state evaluation | 1 | Other / needs triage |
| ~~[#2152](https://github.com/gitleaks/gitleaks/issues/2152)~~ | SDK wrapper package for programmatic embedding | 1 | Other / needs triage |
| ~~[#2131](https://github.com/gitleaks/gitleaks/issues/2131)~~ | Add detection rules for sealed-env (npm + Maven Central) | 1 | Detection accuracy |
| ~~[#2130](https://github.com/gitleaks/gitleaks/issues/2130)~~ | Report unused gitleaks ignore entries | 1 | Reporting and output |
| ~~[#2129](https://github.com/gitleaks/gitleaks/issues/2129)~~ | `gitleaks git` silently reports "no leaks found" with exit 0 when `git log` writes anything to stderr (git-crypt repos affected by default) | 1 | Reporting and output |
| ~~[#2122](https://github.com/gitleaks/gitleaks/issues/2122)~~ | bug: withinColumns compares column numbers across different lines in composite rules, producing incorrect proximity results | 5 | Other / needs triage |
| ~~[#2121](https://github.com/gitleaks/gitleaks/issues/2121)~~ | perf: detectRule() calls r.Regex.FindAllStringIndex twice per match cycle, doubling regex work on every scan | 6 | Detection accuracy |
| ~~[#2120](https://github.com/gitleaks/gitleaks/issues/2120)~~ | fix: update default gitleaks.toml to use [[allowlists]] instead of deprecated [allowlist] | 2 | Detection accuracy |
| ~~[#2112](https://github.com/gitleaks/gitleaks/issues/2112)~~ | Possible leaked API key in this repository | 0 | Integrations |
| ~~[#2111](https://github.com/gitleaks/gitleaks/issues/2111)~~ | Possible leaked API key in this repository | 0 | Integrations |
| ~~[#2110](https://github.com/gitleaks/gitleaks/issues/2110)~~ | Possible leaked API key in this repository | 0 | Integrations |
| ~~[#2107](https://github.com/gitleaks/gitleaks/issues/2107)~~ | consistent ignore paths for git vs dir | 1 | Other / needs triage |
| ~~[#2104](https://github.com/gitleaks/gitleaks/issues/2104)~~ | Request: detection rules for Claude Code OAuth tokens and MCP server credentials | 2 | Detection accuracy |
| ~~[#2094](https://github.com/gitleaks/gitleaks/issues/2094)~~ | Add detection rule for Claude Code Remote Control session URLs | 1 | Detection accuracy |
| ~~[#2093](https://github.com/gitleaks/gitleaks/issues/2093)~~ | False positive: square-access-token matches EAAA… inside base64-encoded MIME / archive payloads | 1 | Detection accuracy |
| ~~[#2086](https://github.com/gitleaks/gitleaks/issues/2086)~~ | Tag v8.30.1 not an ancestor of master, breaks pre-commit autoupdate | 0 | Git and sources |
| ~~[#2074](https://github.com/gitleaks/gitleaks/issues/2074)~~ | Untrusted impersonation of the project at gitleaks.org | 0 | Security hardening |
| ~~[#2068](https://github.com/gitleaks/gitleaks/issues/2068)~~ | Support --report-format gitlab-code-quality | 1 | Reporting and output |
| ~~[#2066](https://github.com/gitleaks/gitleaks/issues/2066)~~ | Support configuration from `pyproject.toml` | 1 | Other / needs triage |
| ~~[#2063](https://github.com/gitleaks/gitleaks/issues/2063)~~ | Org. repo with licence + PR from fork = missing license key error | 1 | Other / needs triage |
| ~~[#2060](https://github.com/gitleaks/gitleaks/issues/2060)~~ | Allow .gitignore syntax in .gitleaksignore or add a flag to use .gitignore directly? | 1 | Git and sources |
| ~~[#2048](https://github.com/gitleaks/gitleaks/issues/2048)~~ | Use [[allowlists]] in generated gitleaks.toml | 5 | Other / needs triage |
| ~~[#2031](https://github.com/gitleaks/gitleaks/issues/2031)~~ | Rename in the go.mod file? | 0 | Dependencies and CI |
| ~~[#2019](https://github.com/gitleaks/gitleaks/issues/2019)~~ | Memory issue in GitLab when scanning large repos with Gitleaks (--max-decode-depth behavior) | 1 | Performance and stability |
| ~~[#2017](https://github.com/gitleaks/gitleaks/issues/2017)~~ | Request: Switch released builds to use stdlib for RE? | 1 | Other / needs triage |
| ~~[#2010](https://github.com/gitleaks/gitleaks/issues/2010)~~ | Can't scan compressed files with unknown extension | 0 | Other / needs triage |
| ~~[#2003](https://github.com/gitleaks/gitleaks/issues/2003)~~ | Passing "--no-notes" via --log-opts option doesn't work with CSV output | 3 | Reporting and output |
| ~~[#2000](https://github.com/gitleaks/gitleaks/issues/2000)~~ | Decoding skips values with Unicode characters | 2 | Other / needs triage |
| ~~[#1999](https://github.com/gitleaks/gitleaks/issues/1999)~~ | Unable to detect simple username and password | 0 | Detection accuracy |
| ~~[#1998](https://github.com/gitleaks/gitleaks/issues/1998)~~ | Rule version in the gitleaks.toml | 0 | Detection accuracy |
| ~~[#1997](https://github.com/gitleaks/gitleaks/issues/1997)~~ | Rules priority | 0 | Other / needs triage |
| ~~[#1993](https://github.com/gitleaks/gitleaks/issues/1993)~~ | gitleaks protect --staged hangs in non-standard terminal environments (PTY/pipe deadlock) | 1 | Git and sources |
| ~~[#1986](https://github.com/gitleaks/gitleaks/issues/1986)~~ | rule: hashicorp-tf-password \| But I dont set any password | 0 | Detection accuracy |
| ~~[#1985](https://github.com/gitleaks/gitleaks/issues/1985)~~ | Configuration with empty [[rules]] array fails to load in v8.29.0 | 1 | Other / needs triage |
| ~~[#1984](https://github.com/gitleaks/gitleaks/issues/1984)~~ | --version flag returns placeholder text in v8.29.0 instead of actual version number | 1 | Other / needs triage |
| ~~[#1981](https://github.com/gitleaks/gitleaks/issues/1981)~~ | git errors are being incorrectly ignored | 3 | Other / needs triage |
| ~~[#1979](https://github.com/gitleaks/gitleaks/issues/1979)~~ | gitleaks dir detect &lt;filename&gt; not respecting the file | 1 | Detection accuracy |
| ~~[#1962](https://github.com/gitleaks/gitleaks/issues/1962)~~ | Column numbers should count unicode codepoints, not bytes | 1 | Other / needs triage |
| ~~[#1958](https://github.com/gitleaks/gitleaks/issues/1958)~~ | Passwords with certain special characters are not detected | 0 | Other / needs triage |
| ~~[#1946](https://github.com/gitleaks/gitleaks/issues/1946)~~ | Add Looker client ID and secret rules | 0 | Detection accuracy |
| ~~[#1944](https://github.com/gitleaks/gitleaks/issues/1944)~~ | Build release binary in GitHub Actions | 0 | Dependencies and CI |
| ~~[#1942](https://github.com/gitleaks/gitleaks/issues/1942)~~ | GCP API key not detected in Google Chat API | 1 | Integrations |
| ~~[#1939](https://github.com/gitleaks/gitleaks/issues/1939)~~ | detect vulnerabilities in gitleaks binary | 0 | Detection accuracy |
| ~~[#1933](https://github.com/gitleaks/gitleaks/issues/1933)~~ | Kubernetes secrets not found if `kind: Secret` comes after multiple base64 encoded secret data | 1 | Detection accuracy |
| ~~[#1932](https://github.com/gitleaks/gitleaks/issues/1932)~~ | Creating Composite Rule Errors | 2 | Detection accuracy |
| ~~[#1927](https://github.com/gitleaks/gitleaks/issues/1927)~~ | Should Gitleaks scan .git/ and dotfiles？ So we can detect secrets in hidden files or folders | 0 | Detection accuracy |
| ~~[#1925](https://github.com/gitleaks/gitleaks/issues/1925)~~ | Config Request: Logging ignored secrets? | 1 | Other / needs triage |
| ~~[#1924](https://github.com/gitleaks/gitleaks/issues/1924)~~ | Using tiktoken token density to estimate randomness | 0 | Detection accuracy |
| ~~[#1923](https://github.com/gitleaks/gitleaks/issues/1923)~~ | Support loading gitleaks configs from git refs | 0 | Other / needs triage |
| ~~[#1919](https://github.com/gitleaks/gitleaks/issues/1919)~~ | targetRules field in [[allowLists]] is ignored | 0 | Other / needs triage |
| ~~[#1913](https://github.com/gitleaks/gitleaks/issues/1913)~~ | GitLeaks to support on merge_group | 0 | Other / needs triage |
| ~~[#1898](https://github.com/gitleaks/gitleaks/issues/1898)~~ | False positive - sourcegraph-access-token matches Git SHA | 2 | Detection accuracy |
| ~~[#1897](https://github.com/gitleaks/gitleaks/issues/1897)~~ | False Positive: (Some) Vscode keybindings - "key": "ctrl+shift+down" | 2 | Detection accuracy |
| ~~[#1895](https://github.com/gitleaks/gitleaks/issues/1895)~~ | `gitleaks-system` fails due to missing `pass_filenames: false` | 1 | Other / needs triage |
| ~~[#1893](https://github.com/gitleaks/gitleaks/issues/1893)~~ | The results generated on the Windows platform and the Linux platform are inconsistent. | 0 | Other / needs triage |
| ~~[#1888](https://github.com/gitleaks/gitleaks/issues/1888)~~ | gitleaks dir [path] argument does not allow multiple values or globs | 2 | Other / needs triage |
| ~~[#1870](https://github.com/gitleaks/gitleaks/issues/1870)~~ | Allow wildcards in .gitleaksignore | 1 | Other / needs triage |
| ~~[#1869](https://github.com/gitleaks/gitleaks/issues/1869)~~ | Provide MCP support to allow command running from coding agents | 1 | Integrations |
| ~~[#1866](https://github.com/gitleaks/gitleaks/issues/1866)~~ | Reporting: Include explicitly example for writing output as JSON to stdout | 1 | Reporting and output |
| ~~[#1864](https://github.com/gitleaks/gitleaks/issues/1864)~~ | Sample AKIA credentials cause SAST/Security Scanners to false flag | 0 | Security hardening |
| ~~[#1862](https://github.com/gitleaks/gitleaks/issues/1862)~~ | The READMe.md mentions to run a deprecated command | 1 | Other / needs triage |
| ~~[#1859](https://github.com/gitleaks/gitleaks/issues/1859)~~ | pre-commit gitleaks-system hook fails with '.pre-commit-config.yaml': Not a directory | 1 | Git and sources |
| ~~[#1858](https://github.com/gitleaks/gitleaks/issues/1858)~~ | The “level” field is not displayed in the sarif report | 0 | Reporting and output |
| ~~[#1857](https://github.com/gitleaks/gitleaks/issues/1857)~~ | False positive: ways for environment variable specification | 0 | Detection accuracy |
| ~~[#1855](https://github.com/gitleaks/gitleaks/issues/1855)~~ | Blog unavailable blocked by cloudflair | 1 | Other / needs triage |
| ~~[#1839](https://github.com/gitleaks/gitleaks/issues/1839)~~ | Module path discrepancy after repository migration | 0 | Other / needs triage |
| ~~[#1838](https://github.com/gitleaks/gitleaks/issues/1838)~~ | `.gitleaks.toml` in `.gitignore` | 1 | Git and sources |
| ~~[#1837](https://github.com/gitleaks/gitleaks/issues/1837)~~ | False positive: e.unk_token,this.vocab=new | 1 | Detection accuracy |
| ~~[#1835](https://github.com/gitleaks/gitleaks/issues/1835)~~ | False positive for `authentication.",     "approve_with_id":` | 0 | Detection accuracy |
| ~~[#1833](https://github.com/gitleaks/gitleaks/issues/1833)~~ | Gitleaks does not detect passwords if they contain special chars? | 0 | Detection accuracy |
| ~~[#1832](https://github.com/gitleaks/gitleaks/issues/1832)~~ | hashicorp-tf-password generates excessive false positives | 1 | Other / needs triage |
| ~~[#1830](https://github.com/gitleaks/gitleaks/issues/1830)~~ | Gitleaks Entropy detection is including all sorts of plaintext variable names and placeholders | 0 | Detection accuracy |
| ~~[#1828](https://github.com/gitleaks/gitleaks/issues/1828)~~ | Password environment variable example with empty value causes false-positive | 0 | Other / needs triage |
| ~~[#1827](https://github.com/gitleaks/gitleaks/issues/1827)~~ | Invalid URI in SARIF file report | 1 | Reporting and output |
| ~~[#1821](https://github.com/gitleaks/gitleaks/issues/1821)~~ | Idiosyncrasy with `gitleaks dir` and allowlist paths | 0 | Detection accuracy |
| ~~[#1818](https://github.com/gitleaks/gitleaks/issues/1818)~~ | `GCPServiceAccount` rule | 2 | Detection accuracy |
| ~~[#1814](https://github.com/gitleaks/gitleaks/issues/1814)~~ | Support custom gitconfig | 0 | Other / needs triage |
| ~~[#1785](https://github.com/gitleaks/gitleaks/issues/1785)~~ | generic secrets with python type hints | 0 | Other / needs triage |
| ~~[#1778](https://github.com/gitleaks/gitleaks/issues/1778)~~ | Add Tailsale keys | 0 | Other / needs triage |
| ~~[#1773](https://github.com/gitleaks/gitleaks/issues/1773)~~ | Add more prefixes to generic-api-key | 0 | Integrations |
| ~~[#1772](https://github.com/gitleaks/gitleaks/issues/1772)~~ | Generic API Key rule should contain additional whitespace before the assignment operator | 1 | Detection accuracy |
| ~~[#1766](https://github.com/gitleaks/gitleaks/issues/1766)~~ | Add "context" to findings | 0 | Other / needs triage |
| ~~[#1762](https://github.com/gitleaks/gitleaks/issues/1762)~~ | The "generic-api-key" rule may raise false positive in vuejs files | 0 | Detection accuracy |
| ~~[#1757](https://github.com/gitleaks/gitleaks/issues/1757)~~ | Handle non-unicode bytes in input | 1 | Other / needs triage |
| ~~[#1754](https://github.com/gitleaks/gitleaks/issues/1754)~~ | Setting a minimum version of GitLeaks when useDefault = true | 0 | Other / needs triage |
| ~~[#1753](https://github.com/gitleaks/gitleaks/issues/1753)~~ | Unable to receive license from email for github organization account | 0 | Integrations |
| ~~[#1732](https://github.com/gitleaks/gitleaks/issues/1732)~~ | Support setting `max-decode-depth` in config file | 0 | Other / needs triage |
| ~~[#1728](https://github.com/gitleaks/gitleaks/issues/1728)~~ | K8s Sealed Secrets detected as Generic-Api-Key in v8.23.2 | 0 | Integrations |
| ~~[#1727](https://github.com/gitleaks/gitleaks/issues/1727)~~ | Passing multiple arguments to `dir` causes a fill scan of the current directory | 1 | Other / needs triage |
| ~~[#1723](https://github.com/gitleaks/gitleaks/issues/1723)~~ | Create an 'experimental' config | 1 | Other / needs triage |
| ~~[#1715](https://github.com/gitleaks/gitleaks/issues/1715)~~ | Create rule for SQL Server connection strings | 2 | Detection accuracy |
| ~~[#1714](https://github.com/gitleaks/gitleaks/issues/1714)~~ | Generate links for directory scans? | 1 | Other / needs triage |
| ~~[#1708](https://github.com/gitleaks/gitleaks/issues/1708)~~ | Base64-decoded findings aren't printed with colour | 1 | Other / needs triage |
| ~~[#1697](https://github.com/gitleaks/gitleaks/issues/1697)~~ | Create rule for SourceGraph tokens | 1 | Detection accuracy |
| ~~[#1689](https://github.com/gitleaks/gitleaks/issues/1689)~~ | Missing XML report format file | 0 | Reporting and output |
| ~~[#1687](https://github.com/gitleaks/gitleaks/issues/1687)~~ | Add support for Azure Development Account Key | 0 | Other / needs triage |
| ~~[#1675](https://github.com/gitleaks/gitleaks/issues/1675)~~ | Add warning if rule doesn't have keywords | 1 | Detection accuracy |
| ~~[#1668](https://github.com/gitleaks/gitleaks/issues/1668)~~ | Add errormessage when specifying non-existant .gitleaksignore | 0 | Other / needs triage |
| ~~[#1651](https://github.com/gitleaks/gitleaks/issues/1651)~~ | Chunking can split secrets, causing them to be undetected. | 2 | Other / needs triage |
| ~~[#1647](https://github.com/gitleaks/gitleaks/issues/1647)~~ | Kubernetes Secrets Can Slip Through Detection in Git Mode | 0 | Detection accuracy |
| ~~[#1641](https://github.com/gitleaks/gitleaks/issues/1641)~~ | Normalize line-endings | 0 | Other / needs triage |
| ~~[#1638](https://github.com/gitleaks/gitleaks/issues/1638)~~ | Duplicate Path Only Rule Findings When Scanning Directory | 1 | Detection accuracy |
| ~~[#1637](https://github.com/gitleaks/gitleaks/issues/1637)~~ | XML Regex is not working in Gitleaks | 0 | Detection accuracy |
| ~~[#1626](https://github.com/gitleaks/gitleaks/issues/1626)~~ | output git command to show context for finding | 0 | Reporting and output |
| ~~[#1624](https://github.com/gitleaks/gitleaks/issues/1624)~~ | Organisation key email never received | 0 | Other / needs triage |
| ~~[#1588](https://github.com/gitleaks/gitleaks/issues/1588)~~ | Request: Normalize build process | 0 | Dependencies and CI |
| ~~[#1572](https://github.com/gitleaks/gitleaks/issues/1572)~~ | .gitleaksignore should be treated as untrusted input | 0 | Security hardening |
| ~~[#1567](https://github.com/gitleaks/gitleaks/issues/1567)~~ | Supress `generic-api-key` findings that were ignored by a specific rule | 0 | Detection accuracy |
| ~~[#1564](https://github.com/gitleaks/gitleaks/issues/1564)~~ | Support generic-api-key with ternary operator | 0 | Integrations |
| ~~[#1557](https://github.com/gitleaks/gitleaks/issues/1557)~~ | Managing custom config file | 1 | Other / needs triage |
| ~~[#1544](https://github.com/gitleaks/gitleaks/issues/1544)~~ | invalid go version '1.22.0': must match format 1.23 | 0 | Reporting and output |
| ~~[#1525](https://github.com/gitleaks/gitleaks/issues/1525)~~ | Scanning for popular password manager export formats | 0 | Other / needs triage |
| ~~[#1513](https://github.com/gitleaks/gitleaks/issues/1513)~~ | Additional Kubernetes Secret Detections | 1 | Detection accuracy |
| ~~[#1511](https://github.com/gitleaks/gitleaks/issues/1511)~~ | Possibility to use a custom replacement string | 0 | Other / needs triage |
| ~~[#1475](https://github.com/gitleaks/gitleaks/issues/1475)~~ | Gitleaks Regex Overreach: Private Key Detection Captures Excessive Unrelated Content | 0 | Detection accuracy |
| ~~[#1467](https://github.com/gitleaks/gitleaks/issues/1467)~~ | OpenAI API Key rule is missing keys | 1 | Detection accuracy |
| ~~[#1465](https://github.com/gitleaks/gitleaks/issues/1465)~~ | ubi8/python 3.12 emits RuntimeWarning for  tarfile | 0 | Other / needs triage |
| ~~[#1464](https://github.com/gitleaks/gitleaks/issues/1464)~~ | Gitleaks returns exit code 0 when a sequence of invalid flags is passed | 2 | Other / needs triage |
| ~~[#1456](https://github.com/gitleaks/gitleaks/issues/1456)~~ | Declare API keys with checksums | 0 | Integrations |
| ~~[#1448](https://github.com/gitleaks/gitleaks/issues/1448)~~ | I'm integrating gitleaks as a precommit commit hooks but it is not working | 0 | Git and sources |
| ~~[#1447](https://github.com/gitleaks/gitleaks/issues/1447)~~ | GitLeaks does not detect Discord Bot Tokens in source code | 0 | Detection accuracy |
| ~~[#1444](https://github.com/gitleaks/gitleaks/issues/1444)~~ | Support Hubspot Private App Access Tokens and Developer API keys | 0 | Integrations |
| ~~[#1438](https://github.com/gitleaks/gitleaks/issues/1438)~~ | Gitleaks on kubernetes | 0 | Other / needs triage |
| ~~[#1437](https://github.com/gitleaks/gitleaks/issues/1437)~~ | hang on futex | 0 | Performance and stability |
| ~~[#1436](https://github.com/gitleaks/gitleaks/issues/1436)~~ | Generic-api-key has a lot of FN | 0 | Integrations |
| ~~[#1424](https://github.com/gitleaks/gitleaks/issues/1424)~~ | Off-by-one error in startColumn and endColumn calculation | 1 | Other / needs triage |
| ~~[#1409](https://github.com/gitleaks/gitleaks/issues/1409)~~ | gitleaks does not scan correct files with pre-commit run --files | 1 | Git and sources |
| ~~[#1405](https://github.com/gitleaks/gitleaks/issues/1405)~~ | Please update the deprecated way how to set output | 0 | Reporting and output |
| ~~[#1397](https://github.com/gitleaks/gitleaks/issues/1397)~~ | Enhance gitleaks:allow to be usable outside the same line | 1 | Other / needs triage |
| ~~[#1385](https://github.com/gitleaks/gitleaks/issues/1385)~~ | lsetxattr: permission denied with podman | 0 | Security hardening |
| ~~[#1381](https://github.com/gitleaks/gitleaks/issues/1381)~~ | gitleaks with the --amend command | 0 | Other / needs triage |
| ~~[#1366](https://github.com/gitleaks/gitleaks/issues/1366)~~ | launchdarkly-access-token is not detecting LaunchDarkly API Access Tokens | 0 | Detection accuracy |
| ~~[#1365](https://github.com/gitleaks/gitleaks/issues/1365)~~ | Installing gitleaks not working | 0 | Other / needs triage |
| ~~[#1364](https://github.com/gitleaks/gitleaks/issues/1364)~~ | Gitleaks-docker crashes due to Git error: unknown option `staged' | 0 | Git and sources |
| ~~[#1361](https://github.com/gitleaks/gitleaks/issues/1361)~~ | Consider yaml/yml files in scan by default | 1 | Other / needs triage |
| ~~[#1357](https://github.com/gitleaks/gitleaks/issues/1357)~~ | Incorrect Start and End Column Values in Gitleaks Report for Findings | 0 | Reporting and output |
| ~~[#1348](https://github.com/gitleaks/gitleaks/issues/1348)~~ | Alibaba secret in CSV format | 0 | Detection accuracy |
| ~~[#1345](https://github.com/gitleaks/gitleaks/issues/1345)~~ | Try to find config file and ignore file in current directory | 0 | Other / needs triage |
| ~~[#1344](https://github.com/gitleaks/gitleaks/issues/1344)~~ | Update the use of allowlist and add ignorelist for gitleaks.toml files | 0 | Detection accuracy |
| ~~[#1343](https://github.com/gitleaks/gitleaks/issues/1343)~~ | Resource not accessible by integration | 0 | Other / needs triage |
| ~~[#1338](https://github.com/gitleaks/gitleaks/issues/1338)~~ | gitdiff report "parser" can cause an infinite wait on parser error | 1 | Reporting and output |
| ~~[#1331](https://github.com/gitleaks/gitleaks/issues/1331)~~ | Some simple API key are not detected in configuration file | 0 | Integrations |
| ~~[#1329](https://github.com/gitleaks/gitleaks/issues/1329)~~ | Facing cannot spawn astextplain in gitleak scan | 0 | Other / needs triage |
| ~~[#1328](https://github.com/gitleaks/gitleaks/issues/1328)~~ | Simple example needed on how to ignore or whitelist a directory | 0 | Other / needs triage |
| ~~[#1326](https://github.com/gitleaks/gitleaks/issues/1326)~~ | Detect a npmrc auth being checked in | 0 | Detection accuracy |
| ~~[#1325](https://github.com/gitleaks/gitleaks/issues/1325)~~ | Enhance function of .gitleaksignore | 0 | Other / needs triage |
| ~~[#1324](https://github.com/gitleaks/gitleaks/issues/1324)~~ | Change AWS regex to match only allowed characters | 1 | Detection accuracy |
| ~~[#1323](https://github.com/gitleaks/gitleaks/issues/1323)~~ | Cannot whitelist a line at repo level in Gitleaks 8.18.1 with regex | 0 | Detection accuracy |
| ~~[#1321](https://github.com/gitleaks/gitleaks/issues/1321)~~ | False positive "generic-api-key" for HashiCorp Vault Plugin configuration in Jenkins pipeline | 0 | Detection accuracy |
| ~~[#1319](https://github.com/gitleaks/gitleaks/issues/1319)~~ | Does not detect leaks in appsettings.json | 0 | Detection accuracy |
| ~~[#1316](https://github.com/gitleaks/gitleaks/issues/1316)~~ | gitconfig breaks commit info for detect command | 0 | Detection accuracy |
| ~~[#1315](https://github.com/gitleaks/gitleaks/issues/1315)~~ | Git leaks keep on scanning previous commit how to mention in the workflow to scan last two commits or latest commit | 0 | Git and sources |
| ~~[#1312](https://github.com/gitleaks/gitleaks/issues/1312)~~ | adding support for writing formatted reports to stdout when --report-path = "" and --report-format is != "" | 0 | Reporting and output |
| ~~[#1311](https://github.com/gitleaks/gitleaks/issues/1311)~~ | FTL unable to load gitleaks config, err: Config File ".gitleaks" Not Found in "[/repo]" | 0 | Other / needs triage |
| ~~[#1310](https://github.com/gitleaks/gitleaks/issues/1310)~~ | detect in a repo that was cloned with `--depth 1` reports wrong commit IDs | 0 | Detection accuracy |
| ~~[#1308](https://github.com/gitleaks/gitleaks/issues/1308)~~ | Generic api key not find special symbols in passwords | 0 | Integrations |
| ~~[#1306](https://github.com/gitleaks/gitleaks/issues/1306)~~ | Gitleaks github actions throwing intermittent error  "ERR [git] fatal: Invalid revision range" | 0 | Integrations |
| ~~[#1305](https://github.com/gitleaks/gitleaks/issues/1305)~~ | git version 1.8.3.1 does not support git -C which the latest version of gitleaks uses. date(2023-11-30). | 0 | Other / needs triage |
| ~~[#1303](https://github.com/gitleaks/gitleaks/issues/1303)~~ | stopwords conflict with regex | 0 | Detection accuracy |
| ~~[#1301](https://github.com/gitleaks/gitleaks/issues/1301)~~ | Add detection of Docker registry credentials | 0 | Detection accuracy |
| ~~[#1293](https://github.com/gitleaks/gitleaks/issues/1293)~~ | can gitleaks understand .gitleaksignore while scanning a bare git repo? | 0 | Other / needs triage |
| ~~[#1287](https://github.com/gitleaks/gitleaks/issues/1287)~~ | Detect NoGit scan appends source folder in the Fingerprint | 0 | Detection accuracy |
| ~~[#1285](https://github.com/gitleaks/gitleaks/issues/1285)~~ | Question: several config files in v8 | 0 | Other / needs triage |
| ~~[#1284](https://github.com/gitleaks/gitleaks/issues/1284)~~ | Baseline entries require more than the fingerprint | 0 | Detection accuracy |
| ~~[#1279](https://github.com/gitleaks/gitleaks/issues/1279)~~ | Allowlist may not be working as expected, per Discord | 0 | Detection accuracy |
| ~~[#1277](https://github.com/gitleaks/gitleaks/issues/1277)~~ | [Security] Vulnerabilities reported in docker image | 0 | Security hardening |
| ~~[#1275](https://github.com/gitleaks/gitleaks/issues/1275)~~ | Help on finding right flag to export report in json | 0 | Reporting and output |
| ~~[#1268](https://github.com/gitleaks/gitleaks/issues/1268)~~ | detect: ability to scan diff files and generate directory-aware reports | 0 | Detection accuracy |
| ~~[#1263](https://github.com/gitleaks/gitleaks/issues/1263)~~ | Question: where to properly store the baseline? | 0 | Detection accuracy |
| ~~[#1261](https://github.com/gitleaks/gitleaks/issues/1261)~~ | gitleaks not correctly showing all result  when no ID is added in rule | 0 | Detection accuracy |
| ~~[#1257](https://github.com/gitleaks/gitleaks/issues/1257)~~ | Failed with Self hosted GitHub Runner | 0 | Integrations |
| ~~[#1251](https://github.com/gitleaks/gitleaks/issues/1251)~~ | Support for Separate Files for Rules and Allowlist in gitleaks | 0 | Detection accuracy |
| ~~[#1246](https://github.com/gitleaks/gitleaks/issues/1246)~~ | Suppress "gitleaks:allow" inline annotation | 0 | Other / needs triage |
| ~~[#1243](https://github.com/gitleaks/gitleaks/issues/1243)~~ | Output shows redacted secret if final new line not present | 0 | Detection accuracy |
| ~~[#1239](https://github.com/gitleaks/gitleaks/issues/1239)~~ | Generic-api-key - checking for the presence of a digit or a stop-word? | 0 | Integrations |
| ~~[#1234](https://github.com/gitleaks/gitleaks/issues/1234)~~ | Adding the ability to scan selected files from the source folder | 2 | Other / needs triage |
| ~~[#1230](https://github.com/gitleaks/gitleaks/issues/1230)~~ | Generic rule false negative when there are many spaces | 0 | Detection accuracy |
| ~~[#1228](https://github.com/gitleaks/gitleaks/issues/1228)~~ | [NEW RULE] [SQL] Cleartext passwords in SQL | 0 | Detection accuracy |
| ~~[#1226](https://github.com/gitleaks/gitleaks/issues/1226)~~ | `go run ...` module declares its path as: github.com/zricethezav/gitleaks | 0 | Integrations |
| ~~[#1210](https://github.com/gitleaks/gitleaks/issues/1210)~~ | Debian and GitLeak build failures | 0 | Dependencies and CI |
| ~~[#1209](https://github.com/gitleaks/gitleaks/issues/1209)~~ | Cannot build from source: incorrect go.mod | 0 | Dependencies and CI |
| ~~[#1206](https://github.com/gitleaks/gitleaks/issues/1206)~~ | gitleaks protect does not detect commits if an external diff tool is used | 0 | Detection accuracy |
| ~~[#1205](https://github.com/gitleaks/gitleaks/issues/1205)~~ | RuleID name change from 'slack-web-hook' to 'slack-webhook-url' in 8.17.0 is not compatible with previous .gitleaksignore Fingerprints | 0 | Integrations |
| ~~[#1203](https://github.com/gitleaks/gitleaks/issues/1203)~~ | basic auth not detected | 0 | Other / needs triage |
| ~~[#1195](https://github.com/gitleaks/gitleaks/issues/1195)~~ | Feature - use the .gitignore as allowlist | 0 | Detection accuracy |
| ~~[#1182](https://github.com/gitleaks/gitleaks/issues/1182)~~ | Feature - Ignore false positives by specifying text to replace before testing line | 0 | Other / needs triage |
| ~~[#1181](https://github.com/gitleaks/gitleaks/issues/1181)~~ | Feature - Ability to ignore files for rule | 0 | Detection accuracy |
| ~~[#1163](https://github.com/gitleaks/gitleaks/issues/1163)~~ | detect --source folder is not honored as a folder in git mode | 0 | Detection accuracy |
| ~~[#1147](https://github.com/gitleaks/gitleaks/issues/1147)~~ | Detect nuget passwords | 0 | Detection accuracy |
| ~~[#1137](https://github.com/gitleaks/gitleaks/issues/1137)~~ | Central gitleaks allowlist stopwords is not working as expected | 0 | Detection accuracy |
| ~~[#1128](https://github.com/gitleaks/gitleaks/issues/1128)~~ | Secrets still being detected in deletions | 0 | Other / needs triage |
| ~~[#1127](https://github.com/gitleaks/gitleaks/issues/1127)~~ | Gitleaks 8.14.1 Strange bug!? | 0 | Other / needs triage |
| ~~[#1122](https://github.com/gitleaks/gitleaks/issues/1122)~~ | Add option to completely disable color codes | 0 | Other / needs triage |
| ~~[#1120](https://github.com/gitleaks/gitleaks/issues/1120)~~ | Please bump golang.org/x/sys dependency | 0 | Dependencies and CI |
| ~~[#1118](https://github.com/gitleaks/gitleaks/issues/1118)~~ | Skip binary files by default without the need to add the extensions in the allow list | 0 | Other / needs triage |
| ~~[#1114](https://github.com/gitleaks/gitleaks/issues/1114)~~ | Add Commvault rule | 0 | Detection accuracy |
| ~~[#1110](https://github.com/gitleaks/gitleaks/issues/1110)~~ | Unclear Sarif Short Description | 0 | Reporting and output |
| ~~[#1091](https://github.com/gitleaks/gitleaks/issues/1091)~~ | Improve gitleaksignore to handle uncommitted secrets | 1 | Other / needs triage |
| ~~[#1089](https://github.com/gitleaks/gitleaks/issues/1089)~~ | New Relic ingest browser API token incorrectly identified as secret | 1 | Detection accuracy |
| ~~[#1085](https://github.com/gitleaks/gitleaks/issues/1085)~~ | Private key regex is greedy | 2 | Detection accuracy |
| ~~[#1082](https://github.com/gitleaks/gitleaks/issues/1082)~~ | Detecting private keys and seed phrases for cryptocurrency wallets | 0 | Other / needs triage |
| ~~[#1080](https://github.com/gitleaks/gitleaks/issues/1080)~~ | Mistake in the description for aws-access-token | 0 | Detection accuracy |
| ~~[#1078](https://github.com/gitleaks/gitleaks/issues/1078)~~ | How to create custom rules to scan secerts in .yaml or .properties file | 1 | Other / needs triage |
| ~~[#1075](https://github.com/gitleaks/gitleaks/issues/1075)~~ | False positives | 0 | Other / needs triage |
| ~~[#1074](https://github.com/gitleaks/gitleaks/issues/1074)~~ | Need more context surrounding secret | 0 | Detection accuracy |
| ~~[#1073](https://github.com/gitleaks/gitleaks/issues/1073)~~ | Public Keys being identified as secrets | 0 | Other / needs triage |
| ~~[#1070](https://github.com/gitleaks/gitleaks/issues/1070)~~ | Leak is not detected | 0 | Other / needs triage |
| ~~[#1067](https://github.com/gitleaks/gitleaks/issues/1067)~~ | Configuration examples gone | 0 | Other / needs triage |
| ~~[#1065](https://github.com/gitleaks/gitleaks/issues/1065)~~ | Validate report-path before scan | 0 | Reporting and output |
| ~~[#1059](https://github.com/gitleaks/gitleaks/issues/1059)~~ | gitleaksignore fingerprint doesn't support relative paths | 0 | Detection accuracy |
| ~~[#1058](https://github.com/gitleaks/gitleaks/issues/1058)~~ | Ignoring embedded images in Jupyter notebooks. | 0 | Other / needs triage |
| ~~[#1057](https://github.com/gitleaks/gitleaks/issues/1057)~~ | Go version 1.19 is unsuported by tidy ( requires 1.19 ) | 0 | Other / needs triage |
| ~~[#1054](https://github.com/gitleaks/gitleaks/issues/1054)~~ | Feature request: findings deduplication | 0 | Other / needs triage |
| ~~[#1053](https://github.com/gitleaks/gitleaks/issues/1053)~~ | RuleID generic-api-key gets confused by key_algorithm | 0 | Integrations |
| ~~[#1052](https://github.com/gitleaks/gitleaks/issues/1052)~~ | In generic-api-key rule, ignore PublicKey and PublicKeyToken | 0 | Detection accuracy |
| ~~[#1049](https://github.com/gitleaks/gitleaks/issues/1049)~~ | Possible false positives with aws-access-token rule | 1 | Detection accuracy |
| ~~[#1048](https://github.com/gitleaks/gitleaks/issues/1048)~~ | Support multiple report formats in a single execution | 1 | Reporting and output |
| ~~[#1047](https://github.com/gitleaks/gitleaks/issues/1047)~~ | Extend generic-api-key rule | 1 | Detection accuracy |
| ~~[#1041](https://github.com/gitleaks/gitleaks/issues/1041)~~ | CentOS 7 - Not able to compile Gitleaks + Binary error | 0 | Other / needs triage |
| ~~[#1036](https://github.com/gitleaks/gitleaks/issues/1036)~~ | Better detection of API tokens | 0 | Detection accuracy |
| ~~[#1034](https://github.com/gitleaks/gitleaks/issues/1034)~~ | Exclude encrypted cosign private keys | 0 | Other / needs triage |
| ~~[#1029](https://github.com/gitleaks/gitleaks/issues/1029)~~ | Documentation: installation of dependencies when installing from source. | 0 | Dependencies and CI |
| ~~[#1028](https://github.com/gitleaks/gitleaks/issues/1028)~~ | Unable to detect secrets introduced during merge | 1 | Detection accuracy |
| ~~[#1027](https://github.com/gitleaks/gitleaks/issues/1027)~~ | Additional documentation for "extending a central config" | 0 | Other / needs triage |
| ~~[#1022](https://github.com/gitleaks/gitleaks/issues/1022)~~ | Error: unknown option "staged" when running pre-commit with docker_image language | 0 | Git and sources |
| ~~[#1013](https://github.com/gitleaks/gitleaks/issues/1013)~~ | adding secret validation feature | 0 | Detection accuracy |
| ~~[#1004](https://github.com/gitleaks/gitleaks/issues/1004)~~ | Add logging around failure to write report | 0 | Reporting and output |
| ~~[#999](https://github.com/gitleaks/gitleaks/issues/999)~~ | How to configure gitleaks | 0 | Other / needs triage |
| ~~[#988](https://github.com/gitleaks/gitleaks/issues/988)~~ | No json report file generated in docker version v8.13.0 | 0 | Reporting and output |
| ~~[#983](https://github.com/gitleaks/gitleaks/issues/983)~~ | print resulting configuration for debug purpose if extend directive is used | 0 | Other / needs triage |
| ~~[#949](https://github.com/gitleaks/gitleaks/issues/949)~~ | can we configure gitleaks as a global level hook/template | 0 | Integrations |
| ~~[#947](https://github.com/gitleaks/gitleaks/issues/947)~~ | Problems with current implementation - adding the ability to extend configuration files | 0 | Other / needs triage |
| ~~[#914](https://github.com/gitleaks/gitleaks/issues/914)~~ | gitleaks fails to detect generic API token in Juypter notebook | 0 | Detection accuracy |
| ~~[#908](https://github.com/gitleaks/gitleaks/issues/908)~~ | Default Generic API rule - does not scan secrets in XML tags | 0 | Detection accuracy |
| ~~[#885](https://github.com/gitleaks/gitleaks/issues/885)~~ | Create docker tags also for major and minor versions | 0 | Dependencies and CI |
| ~~[#833](https://github.com/gitleaks/gitleaks/issues/833)~~ | Undocumented golang and git version requirements | 0 | Other / needs triage |
| ~~[#829](https://github.com/gitleaks/gitleaks/issues/829)~~ | Flag -add_full_row | 0 | Other / needs triage |
| ~~[#827](https://github.com/gitleaks/gitleaks/issues/827)~~ | Add validation config parameter | 0 | Other / needs triage |
| ~~[#821](https://github.com/gitleaks/gitleaks/issues/821)~~ | Incoming commits are not scanned as part of pre-receive hook | 0 | Integrations |
| ~~[#811](https://github.com/gitleaks/gitleaks/issues/811)~~ | gitleaks ignores bad fields in the configuration silently | 0 | Other / needs triage |
| ~~[#800](https://github.com/gitleaks/gitleaks/issues/800)~~ | Add a flag to enable scanning both the git commit message as well as the commit content | 0 | Git and sources |
| ~~[#764](https://github.com/gitleaks/gitleaks/issues/764)~~ | configuration version and verification | 0 | Other / needs triage |
| ~~[#752](https://github.com/gitleaks/gitleaks/issues/752)~~ | invalid memory address or nil pointer dereference | 1 | Performance and stability |
| ~~[#739](https://github.com/gitleaks/gitleaks/issues/739)~~ | Include rules for browser cookie databases in default config | 0 | Other / needs triage |
| ~~[#539](https://github.com/gitleaks/gitleaks/issues/539)~~ | Gitleaks misses quite a few Azure rules. | 4 | Other / needs triage |

## Exclusions

**171 open issues without a detected relation to an open PR:** ~~[#2322](https://github.com/gitleaks/gitleaks/issues/2322)~~, ~~[#2196](https://github.com/gitleaks/gitleaks/issues/2196)~~, ~~[#2188](https://github.com/gitleaks/gitleaks/issues/2188)~~, ~~[#2186](https://github.com/gitleaks/gitleaks/issues/2186)~~, ~~[#2112](https://github.com/gitleaks/gitleaks/issues/2112)~~, ~~[#2111](https://github.com/gitleaks/gitleaks/issues/2111)~~, ~~[#2110](https://github.com/gitleaks/gitleaks/issues/2110)~~, ~~[#2086](https://github.com/gitleaks/gitleaks/issues/2086)~~, ~~[#2074](https://github.com/gitleaks/gitleaks/issues/2074)~~, ~~[#2031](https://github.com/gitleaks/gitleaks/issues/2031)~~, ~~[#2010](https://github.com/gitleaks/gitleaks/issues/2010)~~, ~~[#1999](https://github.com/gitleaks/gitleaks/issues/1999)~~, ~~[#1998](https://github.com/gitleaks/gitleaks/issues/1998)~~, ~~[#1997](https://github.com/gitleaks/gitleaks/issues/1997)~~, ~~[#1986](https://github.com/gitleaks/gitleaks/issues/1986)~~, ~~[#1958](https://github.com/gitleaks/gitleaks/issues/1958)~~, ~~[#1946](https://github.com/gitleaks/gitleaks/issues/1946)~~, ~~[#1944](https://github.com/gitleaks/gitleaks/issues/1944)~~, ~~[#1939](https://github.com/gitleaks/gitleaks/issues/1939)~~, ~~[#1927](https://github.com/gitleaks/gitleaks/issues/1927)~~, ~~[#1924](https://github.com/gitleaks/gitleaks/issues/1924)~~, ~~[#1923](https://github.com/gitleaks/gitleaks/issues/1923)~~, ~~[#1919](https://github.com/gitleaks/gitleaks/issues/1919)~~, ~~[#1913](https://github.com/gitleaks/gitleaks/issues/1913)~~, ~~[#1893](https://github.com/gitleaks/gitleaks/issues/1893)~~, ~~[#1864](https://github.com/gitleaks/gitleaks/issues/1864)~~, ~~[#1858](https://github.com/gitleaks/gitleaks/issues/1858)~~, ~~[#1857](https://github.com/gitleaks/gitleaks/issues/1857)~~, ~~[#1839](https://github.com/gitleaks/gitleaks/issues/1839)~~, ~~[#1835](https://github.com/gitleaks/gitleaks/issues/1835)~~, ~~[#1833](https://github.com/gitleaks/gitleaks/issues/1833)~~, ~~[#1830](https://github.com/gitleaks/gitleaks/issues/1830)~~, ~~[#1828](https://github.com/gitleaks/gitleaks/issues/1828)~~, ~~[#1821](https://github.com/gitleaks/gitleaks/issues/1821)~~, ~~[#1814](https://github.com/gitleaks/gitleaks/issues/1814)~~, ~~[#1785](https://github.com/gitleaks/gitleaks/issues/1785)~~, ~~[#1778](https://github.com/gitleaks/gitleaks/issues/1778)~~, ~~[#1773](https://github.com/gitleaks/gitleaks/issues/1773)~~, ~~[#1766](https://github.com/gitleaks/gitleaks/issues/1766)~~, ~~[#1762](https://github.com/gitleaks/gitleaks/issues/1762)~~, ~~[#1754](https://github.com/gitleaks/gitleaks/issues/1754)~~, ~~[#1753](https://github.com/gitleaks/gitleaks/issues/1753)~~, ~~[#1732](https://github.com/gitleaks/gitleaks/issues/1732)~~, ~~[#1728](https://github.com/gitleaks/gitleaks/issues/1728)~~, ~~[#1689](https://github.com/gitleaks/gitleaks/issues/1689)~~, ~~[#1687](https://github.com/gitleaks/gitleaks/issues/1687)~~, ~~[#1668](https://github.com/gitleaks/gitleaks/issues/1668)~~, ~~[#1647](https://github.com/gitleaks/gitleaks/issues/1647)~~, ~~[#1641](https://github.com/gitleaks/gitleaks/issues/1641)~~, ~~[#1637](https://github.com/gitleaks/gitleaks/issues/1637)~~, ~~[#1626](https://github.com/gitleaks/gitleaks/issues/1626)~~, ~~[#1624](https://github.com/gitleaks/gitleaks/issues/1624)~~, ~~[#1588](https://github.com/gitleaks/gitleaks/issues/1588)~~, ~~[#1572](https://github.com/gitleaks/gitleaks/issues/1572)~~, ~~[#1567](https://github.com/gitleaks/gitleaks/issues/1567)~~, ~~[#1564](https://github.com/gitleaks/gitleaks/issues/1564)~~, ~~[#1544](https://github.com/gitleaks/gitleaks/issues/1544)~~, ~~[#1525](https://github.com/gitleaks/gitleaks/issues/1525)~~, ~~[#1511](https://github.com/gitleaks/gitleaks/issues/1511)~~, ~~[#1475](https://github.com/gitleaks/gitleaks/issues/1475)~~, ~~[#1465](https://github.com/gitleaks/gitleaks/issues/1465)~~, ~~[#1456](https://github.com/gitleaks/gitleaks/issues/1456)~~, ~~[#1448](https://github.com/gitleaks/gitleaks/issues/1448)~~, ~~[#1447](https://github.com/gitleaks/gitleaks/issues/1447)~~, ~~[#1444](https://github.com/gitleaks/gitleaks/issues/1444)~~, ~~[#1438](https://github.com/gitleaks/gitleaks/issues/1438)~~, ~~[#1437](https://github.com/gitleaks/gitleaks/issues/1437)~~, ~~[#1436](https://github.com/gitleaks/gitleaks/issues/1436)~~, ~~[#1405](https://github.com/gitleaks/gitleaks/issues/1405)~~, ~~[#1385](https://github.com/gitleaks/gitleaks/issues/1385)~~, ~~[#1381](https://github.com/gitleaks/gitleaks/issues/1381)~~, ~~[#1366](https://github.com/gitleaks/gitleaks/issues/1366)~~, ~~[#1365](https://github.com/gitleaks/gitleaks/issues/1365)~~, ~~[#1364](https://github.com/gitleaks/gitleaks/issues/1364)~~, ~~[#1357](https://github.com/gitleaks/gitleaks/issues/1357)~~, ~~[#1348](https://github.com/gitleaks/gitleaks/issues/1348)~~, ~~[#1345](https://github.com/gitleaks/gitleaks/issues/1345)~~, ~~[#1344](https://github.com/gitleaks/gitleaks/issues/1344)~~, ~~[#1343](https://github.com/gitleaks/gitleaks/issues/1343)~~, ~~[#1331](https://github.com/gitleaks/gitleaks/issues/1331)~~, ~~[#1329](https://github.com/gitleaks/gitleaks/issues/1329)~~, ~~[#1328](https://github.com/gitleaks/gitleaks/issues/1328)~~, ~~[#1326](https://github.com/gitleaks/gitleaks/issues/1326)~~, ~~[#1325](https://github.com/gitleaks/gitleaks/issues/1325)~~, ~~[#1323](https://github.com/gitleaks/gitleaks/issues/1323)~~, ~~[#1321](https://github.com/gitleaks/gitleaks/issues/1321)~~, ~~[#1319](https://github.com/gitleaks/gitleaks/issues/1319)~~, ~~[#1316](https://github.com/gitleaks/gitleaks/issues/1316)~~, ~~[#1315](https://github.com/gitleaks/gitleaks/issues/1315)~~, ~~[#1312](https://github.com/gitleaks/gitleaks/issues/1312)~~, ~~[#1311](https://github.com/gitleaks/gitleaks/issues/1311)~~, ~~[#1310](https://github.com/gitleaks/gitleaks/issues/1310)~~, ~~[#1308](https://github.com/gitleaks/gitleaks/issues/1308)~~, ~~[#1306](https://github.com/gitleaks/gitleaks/issues/1306)~~, ~~[#1305](https://github.com/gitleaks/gitleaks/issues/1305)~~, ~~[#1303](https://github.com/gitleaks/gitleaks/issues/1303)~~, ~~[#1301](https://github.com/gitleaks/gitleaks/issues/1301)~~, ~~[#1293](https://github.com/gitleaks/gitleaks/issues/1293)~~, ~~[#1287](https://github.com/gitleaks/gitleaks/issues/1287)~~, ~~[#1285](https://github.com/gitleaks/gitleaks/issues/1285)~~, ~~[#1284](https://github.com/gitleaks/gitleaks/issues/1284)~~, ~~[#1279](https://github.com/gitleaks/gitleaks/issues/1279)~~, ~~[#1277](https://github.com/gitleaks/gitleaks/issues/1277)~~, ~~[#1275](https://github.com/gitleaks/gitleaks/issues/1275)~~, ~~[#1268](https://github.com/gitleaks/gitleaks/issues/1268)~~, ~~[#1263](https://github.com/gitleaks/gitleaks/issues/1263)~~, ~~[#1261](https://github.com/gitleaks/gitleaks/issues/1261)~~, ~~[#1257](https://github.com/gitleaks/gitleaks/issues/1257)~~, ~~[#1251](https://github.com/gitleaks/gitleaks/issues/1251)~~, ~~[#1246](https://github.com/gitleaks/gitleaks/issues/1246)~~, ~~[#1243](https://github.com/gitleaks/gitleaks/issues/1243)~~, ~~[#1239](https://github.com/gitleaks/gitleaks/issues/1239)~~, ~~[#1230](https://github.com/gitleaks/gitleaks/issues/1230)~~, ~~[#1228](https://github.com/gitleaks/gitleaks/issues/1228)~~, ~~[#1226](https://github.com/gitleaks/gitleaks/issues/1226)~~, ~~[#1210](https://github.com/gitleaks/gitleaks/issues/1210)~~, ~~[#1209](https://github.com/gitleaks/gitleaks/issues/1209)~~, ~~[#1206](https://github.com/gitleaks/gitleaks/issues/1206)~~, ~~[#1205](https://github.com/gitleaks/gitleaks/issues/1205)~~, ~~[#1203](https://github.com/gitleaks/gitleaks/issues/1203)~~, ~~[#1195](https://github.com/gitleaks/gitleaks/issues/1195)~~, ~~[#1182](https://github.com/gitleaks/gitleaks/issues/1182)~~, ~~[#1181](https://github.com/gitleaks/gitleaks/issues/1181)~~, ~~[#1163](https://github.com/gitleaks/gitleaks/issues/1163)~~, ~~[#1147](https://github.com/gitleaks/gitleaks/issues/1147)~~, ~~[#1137](https://github.com/gitleaks/gitleaks/issues/1137)~~, ~~[#1128](https://github.com/gitleaks/gitleaks/issues/1128)~~, ~~[#1127](https://github.com/gitleaks/gitleaks/issues/1127)~~, ~~[#1122](https://github.com/gitleaks/gitleaks/issues/1122)~~, ~~[#1120](https://github.com/gitleaks/gitleaks/issues/1120)~~, ~~[#1118](https://github.com/gitleaks/gitleaks/issues/1118)~~, ~~[#1114](https://github.com/gitleaks/gitleaks/issues/1114)~~, ~~[#1110](https://github.com/gitleaks/gitleaks/issues/1110)~~, ~~[#1082](https://github.com/gitleaks/gitleaks/issues/1082)~~, ~~[#1080](https://github.com/gitleaks/gitleaks/issues/1080)~~, ~~[#1075](https://github.com/gitleaks/gitleaks/issues/1075)~~, ~~[#1074](https://github.com/gitleaks/gitleaks/issues/1074)~~, ~~[#1073](https://github.com/gitleaks/gitleaks/issues/1073)~~, ~~[#1070](https://github.com/gitleaks/gitleaks/issues/1070)~~, ~~[#1067](https://github.com/gitleaks/gitleaks/issues/1067)~~, ~~[#1065](https://github.com/gitleaks/gitleaks/issues/1065)~~, ~~[#1059](https://github.com/gitleaks/gitleaks/issues/1059)~~, ~~[#1058](https://github.com/gitleaks/gitleaks/issues/1058)~~, ~~[#1057](https://github.com/gitleaks/gitleaks/issues/1057)~~, ~~[#1054](https://github.com/gitleaks/gitleaks/issues/1054)~~, ~~[#1053](https://github.com/gitleaks/gitleaks/issues/1053)~~, ~~[#1052](https://github.com/gitleaks/gitleaks/issues/1052)~~, ~~[#1041](https://github.com/gitleaks/gitleaks/issues/1041)~~, ~~[#1036](https://github.com/gitleaks/gitleaks/issues/1036)~~, ~~[#1034](https://github.com/gitleaks/gitleaks/issues/1034)~~, ~~[#1029](https://github.com/gitleaks/gitleaks/issues/1029)~~, ~~[#1027](https://github.com/gitleaks/gitleaks/issues/1027)~~, ~~[#1022](https://github.com/gitleaks/gitleaks/issues/1022)~~, ~~[#1013](https://github.com/gitleaks/gitleaks/issues/1013)~~, ~~[#1004](https://github.com/gitleaks/gitleaks/issues/1004)~~, ~~[#999](https://github.com/gitleaks/gitleaks/issues/999)~~, ~~[#988](https://github.com/gitleaks/gitleaks/issues/988)~~, ~~[#983](https://github.com/gitleaks/gitleaks/issues/983)~~, ~~[#949](https://github.com/gitleaks/gitleaks/issues/949)~~, ~~[#947](https://github.com/gitleaks/gitleaks/issues/947)~~, ~~[#914](https://github.com/gitleaks/gitleaks/issues/914)~~, ~~[#908](https://github.com/gitleaks/gitleaks/issues/908)~~, ~~[#885](https://github.com/gitleaks/gitleaks/issues/885)~~, ~~[#833](https://github.com/gitleaks/gitleaks/issues/833)~~, ~~[#829](https://github.com/gitleaks/gitleaks/issues/829)~~, ~~[#827](https://github.com/gitleaks/gitleaks/issues/827)~~, ~~[#821](https://github.com/gitleaks/gitleaks/issues/821)~~, ~~[#811](https://github.com/gitleaks/gitleaks/issues/811)~~, ~~[#800](https://github.com/gitleaks/gitleaks/issues/800)~~, ~~[#764](https://github.com/gitleaks/gitleaks/issues/764)~~, ~~[#739](https://github.com/gitleaks/gitleaks/issues/739)~~.

See [crosswalk](03-pr-issue-crosswalk.md) for evidence classifications, link-only out-of-scope issues, and the CSV review register.

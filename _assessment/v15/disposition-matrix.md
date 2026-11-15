# RNC-Mod V15 disposition matrix

All 220 active packages after post-V14 reconciliation. This is a planning disposition, not an upstream closure. Evidence confidence is not implementation confidence. `P0*` means external incident severity, not RNC scanner priority.

## C1 — Critical detector/reader corrections (7 packages)

| Package | Upstream issue(s) | Upstream PR(s) | Status | Action | Severity | Priority | Confidence | Complexity | Acceptance |
|---|---|---|---|---|---|---|---|---|---|
| `SCAN-01` | #2010 | - | CONFIRMED_DEFECT | **FIX** | S1 | P0 | High | M | Reproduce, fix the shared root cause, add positive/negative regression, both regex backends; no default weakening. |
| `SCAN-02` | #2000 | #2172 #2208 | CONFIRMED_DEFECT | **FIX** | S1 | P1 | High | S | Reproduce, fix the shared root cause, add positive/negative regression, both regex backends; no default weakening. |
| `ACC-03` | #2122 | #2150 #2255 | CONFIRMED_DEFECT | **FIX** | S2 | P1 | High | S | Reproduce, fix the shared root cause, add positive/negative regression, both regex backends; no default weakening. |
| `ACC-09` | #2327 | #2335 | CONFIRMED_SOURCE_DEFECT | **FIX** | S1 | P1 | High | S | Reproduce, fix the shared root cause, add positive/negative regression, both regex backends; no default weakening. |
| `ACC-10` | #2257 | #2260 | CONFIRMED_DEFECT | **FIX** | S1 | P0 | High | S | Reproduce, fix the shared root cause, add positive/negative regression, both regex backends; no default weakening. |
| `CONF-08` | - | #2024 | CONFIRMED_DEFECT | **FIX** | S2 | P1 | High | S | Reproduce, fix the shared root cause, add positive/negative regression, both regex backends; no default weakening. |
| `REP-02` | #1424 | #2215 | CONFIRMED_DEFECT | **FIX** | S2 | P1 | High | S | Reproduce, fix the shared root cause, add positive/negative regression, both regex backends; no default weakening. |

## C2 — Git/merge coverage and runtime stability (12 packages)

| Package | Upstream issue(s) | Upstream PR(s) | Status | Action | Severity | Priority | Confidence | Complexity | Acceptance |
|---|---|---|---|---|---|---|---|---|---|
| `GIT-01` | #2003 | #2128 | PARTIAL_FIX | **REPRO** | S2 | P2 | Medium | S | Controlled Git fixture and exact finding/fingerprint/exit parity; policy/default decision recorded. |
| `GIT-02` | #1028 | #2301 | CONFIRMED_FALSE_NEGATIVE | **POLICY+IMPLEMENT** | S1 | P0 | High | L | Controlled Git fixture and exact finding/fingerprint/exit parity; policy/default decision recorded. |
| `GIT-03` | #1993 | - | UNVERIFIED_HANG | **REPRO** | S2 | P2 | Low | M | Controlled Git fixture and exact finding/fingerprint/exit parity; policy/default decision recorded. |
| `GIT-04` | #1647 | - | UNVERIFIED_CONTEXT_GAP | **REPRO** | S2 | P2 | Low | M | Controlled Git fixture and exact finding/fingerprint/exit parity; policy/default decision recorded. |
| `GIT-05` | #1364 | - | ENVIRONMENT_BLOCKED | **REPRO** | S3 | P3 | Low | M | Controlled Git fixture and exact finding/fingerprint/exit parity; policy/default decision recorded. |
| `GIT-06` | #1310 | - | PARTIAL_WARNING_ONLY | **IMPROVE** | S2 | P2 | Medium | S | Controlled Git fixture and exact finding/fingerprint/exit parity; policy/default decision recorded. |
| `GIT-07` | #821 | - | FEATURE_DESIGN | **DEFER** | S2 | P3 | Low | M | Controlled Git fixture and exact finding/fingerprint/exit parity; policy/default decision recorded. |
| `GIT-08` | #1293 | - | POLICY_DESIGN | **DEFER** | S3 | P3 | Low | S | Controlled Git fixture and exact finding/fingerprint/exit parity; policy/default decision recorded. |
| `GIT-09` | #1163 | - | POLICY_DESIGN | **DEFER** | S3 | P3 | Low | S | Controlled Git fixture and exact finding/fingerprint/exit parity; policy/default decision recorded. |
| `PERF-01` | #2019 | - | UNVERIFIED_OOM | **REPRO** | S1 | P1 | Low | L | Controlled Git fixture and exact finding/fingerprint/exit parity; policy/default decision recorded. |
| `PERF-03` | #1437 | - | UNVERIFIED_HANG | **REPRO** | S1 | P1 | Low | M | Controlled Git fixture and exact finding/fingerprint/exit parity; policy/default decision recorded. |
| `PERF-04` | - | #2023 | UNSAFE_PATCH | **REJECT** | S2 | P2 | Medium | M | Controlled Git fixture and exact finding/fingerprint/exit parity; policy/default decision recorded. |

## C3 — Source readers and path identity (14 packages)

| Package | Upstream issue(s) | Upstream PR(s) | Status | Action | Severity | Priority | Confidence | Complexity | Acceptance |
|---|---|---|---|---|---|---|---|---|---|
| `SCAN-03` | #1651 | #2324 | PARTIAL_BOUNDARY_FIX | **FIX** | S1 | P1 | Medium | M | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |
| `SCAN-04` | #1757 | #2348 | BACKEND_COVERAGE_GAP | **ADAPT** | S1 | P2 | Medium | M | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |
| `SCAN-05` | - | #2217 | TEST_ONLY_STALE | **REJECT** | S4 | P3 | High | S | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |
| `SCAN-06` | - | #1951 | MIME_DESIGN_OPEN | **DEFER** | S2 | P2 | Low | M | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |
| `SCAN-07` | #1888 | #1936 | FEATURE_DEFERRED | **DEFER** | S3 | P3 | High | M | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |
| `SCAN-08` | #1821 | - | UNVERIFIED_PATH_GAP | **REPRO** | S2 | P1 | Low | S | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |
| `SCAN-09` | #1641 | - | UNVERIFIED_CRLF_GAP | **REPRO** | S2 | P1 | Low | M | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |
| `SCAN-10` | #1638 | - | UNVERIFIED_DUPLICATE | **REPRO** | S2 | P2 | Low | M | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |
| `SCAN-11` | #1287 | - | PATH_IDENTITY_GAP | **COMBINE** | S2 | P1 | Medium | M | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |
| `SCAN-12` | #1118 | - | PARTIAL_BEHAVIOR | **BENCH** | S3 | P2 | Medium | M | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |
| `SCAN-13` | #1059 | - | PATH_IDENTITY_GAP | **COMBINE** | S2 | P1 | Medium | M | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |
| `SCAN-14` | #1058 | - | WORKAROUND_PRESENT | **ACCEPT** | S3 | P3 | Medium | S | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |
| `SCAN-15` | #914 | - | UNVERIFIED_ENCODING_GAP | **REPRO** | S2 | P2 | Low | M | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |
| `SCAN-16` | #2107 | #2114 | PATH_IDENTITY_GAP | **COMBINE** | S2 | P1 | Medium | M | Synthetic boundary/path fixture; no false-clean, no widened suppression; absolute and codepoint locations. |

## C4 — Configuration/allowlist model (27 packages)

| Package | Upstream issue(s) | Upstream PR(s) | Status | Action | Severity | Priority | Confidence | Complexity | Acceptance |
|---|---|---|---|---|---|---|---|---|---|
| `CONF-01` | #1397 | #2216 | RISKY_FEATURE | **DEFER** | S3 | P3 | Low | L | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-02` | #2120 | #2343 | DOC_CONFIG_GAP | **IMPLEMENT** | S3 | P2 | High | S | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-03` | #2296 | #2297 | FEATURE_POLICY | **DEFER** | S3 | P3 | Low | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-04` | #1925 | #2126 | PRIVACY_RISK | **REJECT** | S2 | P2 | Medium | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-05` | #1838 | #2116 | DOC_GAP | **IMPLEMENT** | S4 | P3 | High | S | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-06` | - | #2101 | UNSAFE_FEATURE | **REJECT** | S2 | P1 | Medium | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-07` | #1985 | #2079 | POLICY_DECISION | **OPERATOR** | S3 | P3 | Medium | S | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-09` | #1998 | - | FEATURE_POLICY | **COMBINE** | S3 | P3 | Low | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-11` | #1870 | #1871 | FEATURE_DEFERRED | **DEFER** | S3 | P3 | Low | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-12` | #1732 | - | FEATURE_DEFERRED | **DEFER** | S4 | P3 | Low | S | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-13` | - | #1581 | UNSAFE_ALLOWLIST | **REJECT** | S2 | P2 | Medium | S | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-14` | #1284 | - | IDENTITY_DESIGN | **COMBINE** | S2 | P2 | Medium | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-15` | #1205 | - | COMPATIBILITY_POLICY | **OPERATOR** | S2 | P2 | Medium | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-16` | #1182 | - | UNSAFE_FEATURE | **REJECT** | S3 | P2 | Medium | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-17` | #1137 | - | UNVERIFIED_STOPWORD_BUG | **REPRO** | S2 | P2 | Low | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-18` | #1344 | - | SCHEMA_POLICY | **DEFER** | S3 | P3 | Low | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-19` | #1311 | - | ENV_REPRO | **REPRO** | S2 | P2 | Low | S | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-20` | #1078 | - | SUPPORT_GAP | **DEFER** | S4 | P3 | Low | S | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-21` | #1675 | #1683 | TOOLING_GAP | **IMPLEMENT** | S3 | P2 | High | S | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-22` | #1567 | - | ACCURACY_OVERLAP | **COMBINE** | S2 | P2 | Medium | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-23` | - | #1499 | TOOLING_DEFERRED | **DEFER** | S3 | P3 | Low | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-24` | - | #1500 | TOOLING_GAP | **ADAPT** | S3 | P2 | Medium | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-25` | #1932 | #1940 #2098 | GENERATOR_GAP | **COMBINE** | S3 | P2 | Medium | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-26` | #1195 #2060 | #2061 | FEATURE_DEFERRED | **DEFER** | S3 | P3 | Low | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-27` | #947 #1251 #1323 #1325 | - | CONFIG_ARCHITECTURE_GAP | **ARCH** | S2 | P2 | Medium | L | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-28` | #764 #1754 | - | POLICY_FLOOR_GAP | **OPERATOR** | S1 | P1 | High | M | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |
| `CONF-29` | #1279 | - | STRICT_RESIDUAL | **REPRO** | S2 | P2 | Medium | S | Config precedence and negative fixtures; never broaden suppression; schema/CLI compatibility documented. |

## C5 — Reporting and fingerprint compatibility (7 packages)

| Package | Upstream issue(s) | Upstream PR(s) | Status | Action | Severity | Priority | Confidence | Complexity | Acceptance |
|---|---|---|---|---|---|---|---|---|---|
| `CONF-30` | #2239 | #1566 #2286 | FINGERPRINT_DESIGN | **ARCH** | S2 | P2 | High | L | Byte/schema/fingerprint parity; bounded redaction of all known secrets; no unkeyed public hash. |
| `REP-01` | #1866 | #2326 | DOC_GAP | **IMPLEMENT** | S4 | P3 | High | S | Byte/schema/fingerprint parity; bounded redaction of all known secrets; no unkeyed public hash. |
| `REP-04` | #1766 | - | SCHEMA_DESIGN | **DEFER** | S2 | P2 | Low | L | Byte/schema/fingerprint parity; bounded redaction of all known secrets; no unkeyed public hash. |
| `REP-06` | #1708 | - | UX_GAP | **DEFER** | S3 | P3 | Medium | M | Byte/schema/fingerprint parity; bounded redaction of all known secrets; no unkeyed public hash. |
| `REP-07` | #1357 | - | COORDINATE_GAP | **REPRO** | S1 | P1 | Medium | L | Byte/schema/fingerprint parity; bounded redaction of all known secrets; no unkeyed public hash. |
| `REP-08` | #1074 | - | PRIVACY_SCHEMA_GAP | **DEFER** | S2 | P2 | Medium | M | Byte/schema/fingerprint parity; bounded redaction of all known secrets; no unkeyed public hash. |
| `REP-09` | - | #1162 | UNSAFE_HASH | **REJECT** | S1 | P1 | High | M | Byte/schema/fingerprint parity; bounded redaction of all known secrets; no unkeyed public hash. |

## C6 — Generic accuracy and independent oracle (48 packages)

| Package | Upstream issue(s) | Upstream PR(s) | Status | Action | Severity | Priority | Confidence | Complexity | Acceptance |
|---|---|---|---|---|---|---|---|---|---|
| `ACC-01` | #1897 | #2125 #2195 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-02` | #2182 | #2183 #2207 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-04` | #1203 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-05` | - | #1043 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-06` | - | #1745 | SEMANTIC_CHANGE | **REJECT** | S3 | P3 | Medium | S | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-07` | - | #2342 | TEST_ONLY | **TEST** | S4 | P3 | High | S | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-08` | #1085 | #2341 | TEST_ONLY | **TEST** | S4 | P3 | High | S | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-11` | - | #2230 | UNICODE_ENTROPY_BUG | **FIX** | S3 | P3 | High | S | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-12` | #2223 | #2229 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-13` | #2196 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-14` | #1837 | #2136 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-15` | #2093 | #2108 | PROVIDER_FALSE_POSITIVE | **ADAPT** | S2 | P2 | Medium | S | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-16` | #1999 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-17` | #1986 | - | PROVIDER_OVERLAP | **COMBINE** | S2 | P2 | Medium | S | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-18` | #1924 | - | UNSUPPORTED_DESIGN | **REJECT** | S4 | P3 | Medium | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-19` | - | #1867 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-20` | #1857 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-21` | #1835 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-22` | #1830 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-23` | #1828 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-24` | - | #1823 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-25` | #1773 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-26` | #1762 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-27` | #1728 | - | GENERIC_FALSE_POSITIVE | **FIX** | S2 | P2 | Medium | S | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-28` | #1637 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-29` | - | #1594 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-30` | #1564 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-31` | #1475 | - | PRIVATE_KEY_OVERREACH | **FIX** | S2 | P2 | Medium | S | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-32` | #1447 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-33` | - | #1358 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-34` | - | #1350 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-35` | #1331 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-36` | #1321 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-37` | #1319 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-38` | #1075 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-39` | #1070 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-40` | #1034 | - | PRIVATE_KEY_FALSE_POSITIVE | **FIX** | S2 | P2 | Medium | S | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-41` | #1361 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-42` | #1785 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-43` | #1036 | - | REQUIREMENTS_MISSING | **ARCH** | S2 | P2 | Low | L | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-44` | #908 | #1265 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-45` | #1933 | #1953 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-46` | #1230 #1772 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-47` | #1047 #1308 #1833 | #1061 #1130 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-48` | #1942 | #1954 | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-49` | #1239 #1436 #1958 | - | DETECTION_TRADEOFF | **POLICY** | S1 | P1 | Medium | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-50` | #1052 #1053 #1073 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |
| `ACC-51` | #1997 | - | UNVERIFIED_DETECTION | **REPRO** | S2 | P2 | Low | M | Independent oracle positive AND negative case; no finding-count-only acceptance; holdout unchanged. |

## C7 — Existing provider-rule repair (20 packages)

| Package | Upstream issue(s) | Upstream PR(s) | Status | Action | Severity | Priority | Confidence | Complexity | Acceptance |
|---|---|---|---|---|---|---|---|---|---|
| `RULE-01` | #2309 | #2032 #2310 #2325 | PARTIAL_OR_BROKEN_RULE | **ADAPT** | S2 | P2 | Medium | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-02` | #2225 | #2226 #2259 | PARTIAL_OR_BROKEN_RULE | **ADAPT** | S2 | P2 | Medium | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-03` | #1898 | #2083 #2109 | CONFIRMED_FALSE_POSITIVE | **FIX** | S2 | P1 | High | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-04` | #2214 | #2224 | PARTIAL_OR_BROKEN_RULE | **ADAPT** | S2 | P2 | Medium | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-08` | #1715 | #2347 | PARTIAL_OR_BROKEN_RULE | **ADAPT** | S2 | P2 | Medium | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-09` | #1818 | #1930 #2345 | PARTIAL_OR_BROKEN_RULE | **ADAPT** | S2 | P2 | Medium | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-13` | #2314 | #2177 #2315 | PARTIAL_OR_BROKEN_RULE | **ADAPT** | S2 | P2 | Medium | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-15` | - | #2290 | CONFIRMED_FALSE_NEGATIVE | **FIX** | S1 | P1 | High | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-16` | #2278 #2279 | #2280 | PARTIAL_OR_BROKEN_RULE | **ADAPT** | S2 | P2 | Medium | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-18` | #1467 | #2240 | CONFIRMED_FALSE_NEGATIVE | **FIX** | S1 | P1 | High | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-19` | #1832 | #2219 | CONFIRMED_FALSE_POSITIVE | **FIX** | S2 | P1 | Medium | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-20` | #2192 | #2193 | CONFIRMED_FALSE_NEGATIVE | **FIX** | S1 | P1 | High | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-23` | #2104 #2158 | #2105 #2159 #2178 | OVERLAPPING_FORMATS | **COMBINE** | S1 | P1 | Medium | M | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-29` | - | #2007 #2085 | PARTIAL_OR_BROKEN_RULE | **ADAPT** | S2 | P2 | Medium | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-41` | #1366 | - | PARTIAL_OR_BROKEN_RULE | **ADAPT** | S2 | P2 | Medium | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-49` | #1089 | - | PUBLIC_TOKEN_MISCLASSIFICATION | **FIX** | S2 | P1 | High | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-52` | #1049 | - | PARTIAL_OR_BROKEN_RULE | **ADAPT** | S2 | P2 | Medium | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-57` | #1324 | #1356 | COMPATIBILITY_GAP | **ADAPT** | S2 | P2 | Medium | M | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-59` | #539 | #1291 #2072 #2106 #2254 | PROVIDER_MATRIX | **COMBINE** | S2 | P2 | Medium | L | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |
| `RULE-60` | - | #2030 | FORMAT_GAP | **FIX** | S2 | P1 | Medium | S | Format-valid positive, public/nonsecret negative, generator/bundle parity, no duplicate rule ID. |

## C8 — Provider expansion and rule governance (41 packages)

| Package | Upstream issue(s) | Upstream PR(s) | Status | Action | Severity | Priority | Confidence | Complexity | Acceptance |
|---|---|---|---|---|---|---|---|---|---|
| `RULE-05` | - | #955 | OVERBROAD_PATTERN | **ADAPT** | S2 | P2 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-06` | - | #643 | LEGACY_OVERBROAD | **REJECT** | S3 | P3 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-07` | - | #1373 | MIXED_SCOPE | **COMBINE** | S2 | P2 | Medium | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-10` | - | #2340 | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-11` | - | #2333 | LOW_VALUE_PROVIDER | **REJECT** | S4 | P3 | Low | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-12` | #1697 | #2332 | TEST_ONLY | **TEST** | S4 | P3 | High | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-14` | - | #2078 #2117 #2303 | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-17` | - | #2242 | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-21` | - | #2185 | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-22` | #2180 | #2181 | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-24` | - | #2168 | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-25` | #2131 | #2132 | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-26` | - | #2119 | CONFIG_ONLY_CONFLICT | **REJECT** | S3 | P3 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-27` | - | #2118 | CONFIG_ONLY_CONFLICT | **REJECT** | S3 | P3 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-28` | #2094 | #2095 | NON_SECRET_URL | **REJECT** | S3 | P3 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-30` | - | #2077 | OVERBROAD_PATTERN | **REJECT** | S2 | P2 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-31` | - | #2041 | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-32` | - | #1996 | GENERIC_DUPLICATE | **REJECT** | S3 | P3 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-33` | #1946 | - | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-34` | - | #1829 | GENERIC_DUPLICATE | **REJECT** | S3 | P3 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-35` | #1778 | - | OBSOLETE_PROVIDER | **REJECT** | S4 | P3 | Low | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-36` | #1687 | - | PROVIDER_OVERLAP | **COMBINE** | S3 | P2 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-37` | - | #1568 | PROVIDER_OVERLAP | **COMBINE** | S3 | P2 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-38` | - | #1497 | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-39` | - | #1495 | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-40` | #1444 | - | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-42` | - | #1335 | OVERBROAD_PATTERN | **REJECT** | S2 | P2 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-43` | - | #1327 | PROVIDER_OVERLAP | **COMBINE** | S3 | P2 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-44` | #1326 | - | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-45` | #1301 | - | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-46` | #1228 | - | GENERIC_DUPLICATE | **REJECT** | S3 | P3 | Medium | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-47` | #1147 | - | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-48` | #1114 | - | LOW_VALUE_PROVIDER | **REJECT** | S4 | P3 | Low | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-50` | - | #1086 | PUBLIC_IDENTIFIER | **REJECT** | S2 | P2 | High | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-51` | #1082 | - | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-53` | - | #968 | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-54` | - | #774 | PUBLIC_TOKEN | **REJECT** | S2 | P2 | High | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-55` | #1513 | - | PROVIDER_OVERLAP | **COMBINE** | S3 | P2 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-56` | #1348 | - | PROVIDER_OVERLAP | **COMBINE** | S3 | P2 | Medium | S | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-58` | - | #1046 | FEATURE_PROPOSAL | **DEFER** | S3 | P3 | Low | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |
| `RULE-61` | - | #2001 | MIXED_SCOPE | **REJECT** | S2 | P2 | Medium | M | Provider-format evidence and user value; otherwise defer/reject; no broad unbounded patterns. |

## C9 — Release identity and artifact integrity (19 packages)

| Package | Upstream issue(s) | Upstream PR(s) | Status | Action | Severity | Priority | Confidence | Complexity | Acceptance |
|---|---|---|---|---|---|---|---|---|---|
| `SEC-05` | - | #2036 | FIXTURE_HYGIENE | **IMPLEMENT** | S3 | P2 | High | S | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `DEP-01` | - | #2323 | RELEASE_GAP | **RELEASE** | S2 | P2 | Medium | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `DEP-02` | - | #2115 | TOOLCHAIN_POLICY | **OPERATOR** | S2 | P2 | Medium | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `DEP-03` | #1839 | - | MODULE_IDENTITY | **OPERATOR** | S2 | P1 | High | L | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `DEP-04` | #1305 | - | SUPPORT_FLOOR | **OPERATOR** | S3 | P2 | Medium | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `DEP-05` | #1226 | - | RELEASE_GAP | **RELEASE** | S2 | P2 | Medium | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `DEP-06` | #1209 | - | RELEASE_GAP | **RELEASE** | S2 | P2 | Medium | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `DEP-07` | #833 | - | RELEASE_GAP | **RELEASE** | S2 | P2 | Medium | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `DEP-08` | #1944 | - | RELEASE_GAP | **RELEASE** | S2 | P2 | Medium | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `INT-01` | #2017 | - | RELEASE_RUNTIME_GAP | **REPRO** | S1 | P1 | Medium | L | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `INT-02` | #1984 | #2081 | RELEASE_GAP | **RELEASE** | S2 | P2 | Medium | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `INT-03` | - | #2220 | RELEASE_GAP | **RELEASE** | S2 | P2 | Medium | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `INT-04` | #1127 | - | UNVERIFIED_REPORT | **REPRO** | S3 | P3 | Low | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `INT-05` | - | #1217 | RELEASE_SECURITY_GAP | **RELEASE** | S1 | P1 | Medium | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `INT-06` | #885 | - | RELEASE_POLICY | **OPERATOR** | S3 | P2 | Medium | S | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `DOC-01` | #2173 | #2174 #2213 #2344 | RELEASE_GAP | **RELEASE** | S2 | P2 | Medium | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `DOC-02` | - | #2084 | UNSAFE_EXAMPLE | **REJECT** | S2 | P2 | High | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `DOC-03` | #1067 | - | LOW_VALUE_DOC | **DEFER** | S4 | P3 | Medium | S | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |
| `DOC-04` | #2295 | #2306 | RELEASE_GAP | **RELEASE** | S2 | P2 | Medium | M | Artifact/digest/CI evidence; consumer verification; no release claim from build success alone. |

## C10 — Elective feature triage (22 packages)

| Package | Upstream issue(s) | Upstream PR(s) | Status | Action | Severity | Priority | Confidence | Complexity | Acceptance |
|---|---|---|---|---|---|---|---|---|---|
| `FEAT-01` | - | #1992 | UNSAFE_NETWORK_FEATURE | **REJECT** | S1 | P1 | High | L | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-02` | #2152 #2153 | #2154 | SDK_COUPLING | **DEFER** | S2 | P3 | Low | L | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-03` | #1869 | #2056 | AGENT_SURFACE | **DEFER** | S2 | P3 | Low | L | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-04` | #1054 | - | PARTIAL_IMPLEMENTATION | **ACCEPT** | S3 | P3 | High | S | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-05` | #1013 | - | NETWORK_EGRESS_FEATURE | **REJECT** | S1 | P1 | High | L | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-06` | #1456 | - | SCHEMA_COMPLEXITY | **DEFER** | S3 | P3 | Low | L | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-07` | #1923 | - | TRUST_BOUNDARY_FEATURE | **REJECT** | S2 | P2 | Medium | M | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-08` | #2130 | #2244 | PARTIAL_IMPLEMENTATION | **ACCEPT** | S4 | P3 | High | S | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-09` | #2066 | #2097 | LOW_VALUE_FEATURE | **DEFER** | S4 | P3 | Low | S | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-10` | - | #1874 | NEW_LANGUAGE_FEATURE | **REJECT** | S3 | P3 | Medium | L | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-11` | #1557 #2186 | #1630 #1896 | REMOTE_TRUST_FEATURE | **REJECT** | S1 | P1 | High | M | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-12` | #1723 | #1724 | RULESET_POLICY | **REJECT** | S3 | P3 | Low | M | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-13` | #983 | #1654 | PARTIAL_IMPLEMENTATION | **ACCEPT** | S4 | P3 | High | S | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-14` | - | #1661 | UNSAFE_PACKAGING | **REJECT** | S2 | P2 | High | L | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-15` | - | #1485 | UNSAFE_INSTALLER | **REJECT** | S2 | P2 | High | M | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-16` | #800 | - | OPTIONAL_SOURCE_MODE | **DEFER** | S3 | P3 | Low | M | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-17` | #2068 | #2096 | INTEGRATION_POLICY | **DEFER** | S3 | P3 | Low | L | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-18` | #1048 | - | REPORT_ARCHITECTURE | **DEFER** | S2 | P2 | Medium | L | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-19` | #1525 | - | ENCRYPTED_FORMAT_FEATURE | **REJECT** | S3 | P3 | Low | L | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-20` | #1268 | - | NEW_SOURCE_MODE | **DEFER** | S3 | P3 | Low | L | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-21` | #1234 | - | DUPLICATE_FEATURE | **ACCEPT** | S4 | P3 | High | S | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |
| `FEAT-22` | #739 | - | LOW_VALUE_BINARY_FEATURE | **REJECT** | S3 | P3 | Low | M | Opt-in, bounded, fail-closed and independently justified; otherwise reject as unnecessary. |

## OPS — External upstream incident response (3 packages)

| Package | Upstream issue(s) | Upstream PR(s) | Status | Action | Severity | Priority | Confidence | Complexity | Acceptance |
|---|---|---|---|---|---|---|---|---|---|
| `SEC-01` | #2110 | - | EXTERNAL_INCIDENT | **ESCALATE** | S1 | P0* | Low | S | Upstream incident owner/credential owner only; never use or copy reported values. |
| `SEC-02` | #2111 | - | EXTERNAL_INCIDENT | **ESCALATE** | S1 | P0* | Low | S | Upstream incident owner/credential owner only; never use or copy reported values. |
| `SEC-03` | #2112 | - | EXTERNAL_INCIDENT | **ESCALATE** | S1 | P0* | Low | S | Upstream incident owner/credential owner only; never use or copy reported values. |

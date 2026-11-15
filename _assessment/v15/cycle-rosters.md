# V15 cycle rosters

Every active package is assigned to exactly one planning cycle. The upstream references are package references from the canonical row; context-only links are not counted as open PRs. `OPS` is not a development cycle.

## C1 — Critical detector/reader corrections (7 packages)

| Package | Upstream issues | Upstream PRs | Status | Action |
|---|---|---|---|---|
| `SCAN-01` | #2010 | - | CONFIRMED_DEFECT | **FIX** |
| `SCAN-02` | #2000 | #2172 #2208 | CONFIRMED_DEFECT | **FIX** |
| `ACC-03` | #2122 | #2150 #2255 | CONFIRMED_DEFECT | **FIX** |
| `ACC-09` | #2327 | #2335 | CONFIRMED_SOURCE_DEFECT | **FIX** |
| `ACC-10` | #2257 | #2260 | CONFIRMED_DEFECT | **FIX** |
| `CONF-08` | - | #2024 | CONFIRMED_DEFECT | **FIX** |
| `REP-02` | #1424 | #2215 | CONFIRMED_DEFECT | **FIX** |

## C2 — Git/merge coverage and runtime stability (12 packages)

| Package | Upstream issues | Upstream PRs | Status | Action |
|---|---|---|---|---|
| `GIT-01` | #2003 | #2128 | PARTIAL_FIX | **REPRO** |
| `GIT-02` | #1028 | #2301 | CONFIRMED_FALSE_NEGATIVE | **POLICY+IMPLEMENT** |
| `GIT-03` | #1993 | - | UNVERIFIED_HANG | **REPRO** |
| `GIT-04` | #1647 | - | UNVERIFIED_CONTEXT_GAP | **REPRO** |
| `GIT-05` | #1364 | - | ENVIRONMENT_BLOCKED | **REPRO** |
| `GIT-06` | #1310 | - | PARTIAL_WARNING_ONLY | **IMPROVE** |
| `GIT-07` | #821 | - | FEATURE_DESIGN | **DEFER** |
| `GIT-08` | #1293 | - | POLICY_DESIGN | **DEFER** |
| `GIT-09` | #1163 | - | POLICY_DESIGN | **DEFER** |
| `PERF-01` | #2019 | - | UNVERIFIED_OOM | **REPRO** |
| `PERF-03` | #1437 | - | UNVERIFIED_HANG | **REPRO** |
| `PERF-04` | - | #2023 | UNSAFE_PATCH | **REJECT** |

## C3 — Source readers and path identity (14 packages)

| Package | Upstream issues | Upstream PRs | Status | Action |
|---|---|---|---|---|
| `SCAN-03` | #1651 | #2324 | PARTIAL_BOUNDARY_FIX | **FIX** |
| `SCAN-04` | #1757 | #2348 | BACKEND_COVERAGE_GAP | **ADAPT** |
| `SCAN-05` | - | #2217 | TEST_ONLY_STALE | **REJECT** |
| `SCAN-06` | - | #1951 | MIME_DESIGN_OPEN | **DEFER** |
| `SCAN-07` | #1888 | #1936 | FEATURE_DEFERRED | **DEFER** |
| `SCAN-08` | #1821 | - | UNVERIFIED_PATH_GAP | **REPRO** |
| `SCAN-09` | #1641 | - | UNVERIFIED_CRLF_GAP | **REPRO** |
| `SCAN-10` | #1638 | - | UNVERIFIED_DUPLICATE | **REPRO** |
| `SCAN-11` | #1287 | - | PATH_IDENTITY_GAP | **COMBINE** |
| `SCAN-12` | #1118 | - | PARTIAL_BEHAVIOR | **BENCH** |
| `SCAN-13` | #1059 | - | PATH_IDENTITY_GAP | **COMBINE** |
| `SCAN-14` | #1058 | - | WORKAROUND_PRESENT | **ACCEPT** |
| `SCAN-15` | #914 | - | UNVERIFIED_ENCODING_GAP | **REPRO** |
| `SCAN-16` | #2107 | #2114 | PATH_IDENTITY_GAP | **COMBINE** |

## C4 — Configuration/allowlist model (27 packages)

| Package | Upstream issues | Upstream PRs | Status | Action |
|---|---|---|---|---|
| `CONF-01` | #1397 | #2216 | RISKY_FEATURE | **DEFER** |
| `CONF-02` | #2120 | #2343 | DOC_CONFIG_GAP | **IMPLEMENT** |
| `CONF-03` | #2296 | #2297 | FEATURE_POLICY | **DEFER** |
| `CONF-04` | #1925 | #2126 | PRIVACY_RISK | **REJECT** |
| `CONF-05` | #1838 | #2116 | DOC_GAP | **IMPLEMENT** |
| `CONF-06` | - | #2101 | UNSAFE_FEATURE | **REJECT** |
| `CONF-07` | #1985 | #2079 | POLICY_DECISION | **OPERATOR** |
| `CONF-09` | #1998 | - | FEATURE_POLICY | **COMBINE** |
| `CONF-11` | #1870 | #1871 | FEATURE_DEFERRED | **DEFER** |
| `CONF-12` | #1732 | - | FEATURE_DEFERRED | **DEFER** |
| `CONF-13` | - | #1581 | UNSAFE_ALLOWLIST | **REJECT** |
| `CONF-14` | #1284 | - | IDENTITY_DESIGN | **COMBINE** |
| `CONF-15` | #1205 | - | COMPATIBILITY_POLICY | **OPERATOR** |
| `CONF-16` | #1182 | - | UNSAFE_FEATURE | **REJECT** |
| `CONF-17` | #1137 | - | UNVERIFIED_STOPWORD_BUG | **REPRO** |
| `CONF-18` | #1344 | - | SCHEMA_POLICY | **DEFER** |
| `CONF-19` | #1311 | - | ENV_REPRO | **REPRO** |
| `CONF-20` | #1078 | - | SUPPORT_GAP | **DEFER** |
| `CONF-21` | #1675 | #1683 | TOOLING_GAP | **IMPLEMENT** |
| `CONF-22` | #1567 | - | ACCURACY_OVERLAP | **COMBINE** |
| `CONF-23` | - | #1499 | TOOLING_DEFERRED | **DEFER** |
| `CONF-24` | - | #1500 | TOOLING_GAP | **ADAPT** |
| `CONF-25` | #1932 | #1940 #2098 | GENERATOR_GAP | **COMBINE** |
| `CONF-26` | #1195 #2060 | #2061 | FEATURE_DEFERRED | **DEFER** |
| `CONF-27` | #947 #1251 #1323 #1325 | - | CONFIG_ARCHITECTURE_GAP | **ARCH** |
| `CONF-28` | #764 #1754 | - | POLICY_FLOOR_GAP | **OPERATOR** |
| `CONF-29` | #1279 | - | STRICT_RESIDUAL | **REPRO** |

## C5 — Reporting and fingerprint compatibility (7 packages)

| Package | Upstream issues | Upstream PRs | Status | Action |
|---|---|---|---|---|
| `CONF-30` | #2239 | #1566 #2286 | FINGERPRINT_DESIGN | **ARCH** |
| `REP-01` | #1866 | #2326 | DOC_GAP | **IMPLEMENT** |
| `REP-04` | #1766 | - | SCHEMA_DESIGN | **DEFER** |
| `REP-06` | #1708 | - | UX_GAP | **DEFER** |
| `REP-07` | #1357 | - | COORDINATE_GAP | **REPRO** |
| `REP-08` | #1074 | - | PRIVACY_SCHEMA_GAP | **DEFER** |
| `REP-09` | - | #1162 | UNSAFE_HASH | **REJECT** |

## C6 — Generic accuracy and independent oracle (48 packages)

| Package | Upstream issues | Upstream PRs | Status | Action |
|---|---|---|---|---|
| `ACC-01` | #1897 | #2125 #2195 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-02` | #2182 | #2183 #2207 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-04` | #1203 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-05` | - | #1043 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-06` | - | #1745 | SEMANTIC_CHANGE | **REJECT** |
| `ACC-07` | - | #2342 | TEST_ONLY | **TEST** |
| `ACC-08` | #1085 | #2341 | TEST_ONLY | **TEST** |
| `ACC-11` | - | #2230 | UNICODE_ENTROPY_BUG | **FIX** |
| `ACC-12` | #2223 | #2229 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-13` | #2196 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-14` | #1837 | #2136 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-15` | #2093 | #2108 | PROVIDER_FALSE_POSITIVE | **ADAPT** |
| `ACC-16` | #1999 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-17` | #1986 | - | PROVIDER_OVERLAP | **COMBINE** |
| `ACC-18` | #1924 | - | UNSUPPORTED_DESIGN | **REJECT** |
| `ACC-19` | - | #1867 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-20` | #1857 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-21` | #1835 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-22` | #1830 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-23` | #1828 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-24` | - | #1823 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-25` | #1773 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-26` | #1762 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-27` | #1728 | - | GENERIC_FALSE_POSITIVE | **FIX** |
| `ACC-28` | #1637 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-29` | - | #1594 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-30` | #1564 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-31` | #1475 | - | PRIVATE_KEY_OVERREACH | **FIX** |
| `ACC-32` | #1447 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-33` | - | #1358 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-34` | - | #1350 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-35` | #1331 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-36` | #1321 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-37` | #1319 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-38` | #1075 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-39` | #1070 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-40` | #1034 | - | PRIVATE_KEY_FALSE_POSITIVE | **FIX** |
| `ACC-41` | #1361 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-42` | #1785 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-43` | #1036 | - | REQUIREMENTS_MISSING | **ARCH** |
| `ACC-44` | #908 | #1265 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-45` | #1933 | #1953 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-46` | #1230 #1772 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-47` | #1047 #1308 #1833 | #1061 #1130 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-48` | #1942 | #1954 | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-49` | #1239 #1436 #1958 | - | DETECTION_TRADEOFF | **POLICY** |
| `ACC-50` | #1052 #1053 #1073 | - | UNVERIFIED_DETECTION | **REPRO** |
| `ACC-51` | #1997 | - | UNVERIFIED_DETECTION | **REPRO** |

## C7 — Existing provider-rule repair (20 packages)

| Package | Upstream issues | Upstream PRs | Status | Action |
|---|---|---|---|---|
| `RULE-01` | #2309 | #2032 #2310 #2325 | PARTIAL_OR_BROKEN_RULE | **ADAPT** |
| `RULE-02` | #2225 | #2226 #2259 | PARTIAL_OR_BROKEN_RULE | **ADAPT** |
| `RULE-03` | #1898 | #2083 #2109 | CONFIRMED_FALSE_POSITIVE | **FIX** |
| `RULE-04` | #2214 | #2224 | PARTIAL_OR_BROKEN_RULE | **ADAPT** |
| `RULE-08` | #1715 | #2347 | PARTIAL_OR_BROKEN_RULE | **ADAPT** |
| `RULE-09` | #1818 | #1930 #2345 | PARTIAL_OR_BROKEN_RULE | **ADAPT** |
| `RULE-13` | #2314 | #2177 #2315 | PARTIAL_OR_BROKEN_RULE | **ADAPT** |
| `RULE-15` | - | #2290 | CONFIRMED_FALSE_NEGATIVE | **FIX** |
| `RULE-16` | #2278 #2279 | #2280 | PARTIAL_OR_BROKEN_RULE | **ADAPT** |
| `RULE-18` | #1467 | #2240 | CONFIRMED_FALSE_NEGATIVE | **FIX** |
| `RULE-19` | #1832 | #2219 | CONFIRMED_FALSE_POSITIVE | **FIX** |
| `RULE-20` | #2192 | #2193 | CONFIRMED_FALSE_NEGATIVE | **FIX** |
| `RULE-23` | #2104 #2158 | #2105 #2159 #2178 | OVERLAPPING_FORMATS | **COMBINE** |
| `RULE-29` | - | #2007 #2085 | PARTIAL_OR_BROKEN_RULE | **ADAPT** |
| `RULE-41` | #1366 | - | PARTIAL_OR_BROKEN_RULE | **ADAPT** |
| `RULE-49` | #1089 | - | PUBLIC_TOKEN_MISCLASSIFICATION | **FIX** |
| `RULE-52` | #1049 | - | PARTIAL_OR_BROKEN_RULE | **ADAPT** |
| `RULE-57` | #1324 | #1356 | COMPATIBILITY_GAP | **ADAPT** |
| `RULE-59` | #539 | #1291 #2072 #2106 #2254 | PROVIDER_MATRIX | **COMBINE** |
| `RULE-60` | - | #2030 | FORMAT_GAP | **FIX** |

## C8 — Provider expansion and rule governance (41 packages)

| Package | Upstream issues | Upstream PRs | Status | Action |
|---|---|---|---|---|
| `RULE-05` | - | #955 | OVERBROAD_PATTERN | **ADAPT** |
| `RULE-06` | - | #643 | LEGACY_OVERBROAD | **REJECT** |
| `RULE-07` | - | #1373 | MIXED_SCOPE | **COMBINE** |
| `RULE-10` | - | #2340 | FEATURE_PROPOSAL | **DEFER** |
| `RULE-11` | - | #2333 | LOW_VALUE_PROVIDER | **REJECT** |
| `RULE-12` | #1697 | #2332 | TEST_ONLY | **TEST** |
| `RULE-14` | - | #2078 #2117 #2303 | FEATURE_PROPOSAL | **DEFER** |
| `RULE-17` | - | #2242 | FEATURE_PROPOSAL | **DEFER** |
| `RULE-21` | - | #2185 | FEATURE_PROPOSAL | **DEFER** |
| `RULE-22` | #2180 | #2181 | FEATURE_PROPOSAL | **DEFER** |
| `RULE-24` | - | #2168 | FEATURE_PROPOSAL | **DEFER** |
| `RULE-25` | #2131 | #2132 | FEATURE_PROPOSAL | **DEFER** |
| `RULE-26` | - | #2119 | CONFIG_ONLY_CONFLICT | **REJECT** |
| `RULE-27` | - | #2118 | CONFIG_ONLY_CONFLICT | **REJECT** |
| `RULE-28` | #2094 | #2095 | NON_SECRET_URL | **REJECT** |
| `RULE-30` | - | #2077 | OVERBROAD_PATTERN | **REJECT** |
| `RULE-31` | - | #2041 | FEATURE_PROPOSAL | **DEFER** |
| `RULE-32` | - | #1996 | GENERIC_DUPLICATE | **REJECT** |
| `RULE-33` | #1946 | - | FEATURE_PROPOSAL | **DEFER** |
| `RULE-34` | - | #1829 | GENERIC_DUPLICATE | **REJECT** |
| `RULE-35` | #1778 | - | OBSOLETE_PROVIDER | **REJECT** |
| `RULE-36` | #1687 | - | PROVIDER_OVERLAP | **COMBINE** |
| `RULE-37` | - | #1568 | PROVIDER_OVERLAP | **COMBINE** |
| `RULE-38` | - | #1497 | FEATURE_PROPOSAL | **DEFER** |
| `RULE-39` | - | #1495 | FEATURE_PROPOSAL | **DEFER** |
| `RULE-40` | #1444 | - | FEATURE_PROPOSAL | **DEFER** |
| `RULE-42` | - | #1335 | OVERBROAD_PATTERN | **REJECT** |
| `RULE-43` | - | #1327 | PROVIDER_OVERLAP | **COMBINE** |
| `RULE-44` | #1326 | - | FEATURE_PROPOSAL | **DEFER** |
| `RULE-45` | #1301 | - | FEATURE_PROPOSAL | **DEFER** |
| `RULE-46` | #1228 | - | GENERIC_DUPLICATE | **REJECT** |
| `RULE-47` | #1147 | - | FEATURE_PROPOSAL | **DEFER** |
| `RULE-48` | #1114 | - | LOW_VALUE_PROVIDER | **REJECT** |
| `RULE-50` | - | #1086 | PUBLIC_IDENTIFIER | **REJECT** |
| `RULE-51` | #1082 | - | FEATURE_PROPOSAL | **DEFER** |
| `RULE-53` | - | #968 | FEATURE_PROPOSAL | **DEFER** |
| `RULE-54` | - | #774 | PUBLIC_TOKEN | **REJECT** |
| `RULE-55` | #1513 | - | PROVIDER_OVERLAP | **COMBINE** |
| `RULE-56` | #1348 | - | PROVIDER_OVERLAP | **COMBINE** |
| `RULE-58` | - | #1046 | FEATURE_PROPOSAL | **DEFER** |
| `RULE-61` | - | #2001 | MIXED_SCOPE | **REJECT** |

## C9 — Release identity and artifact integrity (19 packages)

| Package | Upstream issues | Upstream PRs | Status | Action |
|---|---|---|---|---|
| `SEC-05` | - | #2036 | FIXTURE_HYGIENE | **IMPLEMENT** |
| `DEP-01` | - | #2323 | RELEASE_GAP | **RELEASE** |
| `DEP-02` | - | #2115 | TOOLCHAIN_POLICY | **OPERATOR** |
| `DEP-03` | #1839 | - | MODULE_IDENTITY | **OPERATOR** |
| `DEP-04` | #1305 | - | SUPPORT_FLOOR | **OPERATOR** |
| `DEP-05` | #1226 | - | RELEASE_GAP | **RELEASE** |
| `DEP-06` | #1209 | - | RELEASE_GAP | **RELEASE** |
| `DEP-07` | #833 | - | RELEASE_GAP | **RELEASE** |
| `DEP-08` | #1944 | - | RELEASE_GAP | **RELEASE** |
| `INT-01` | #2017 | - | RELEASE_RUNTIME_GAP | **REPRO** |
| `INT-02` | #1984 | #2081 | RELEASE_GAP | **RELEASE** |
| `INT-03` | - | #2220 | RELEASE_GAP | **RELEASE** |
| `INT-04` | #1127 | - | UNVERIFIED_REPORT | **REPRO** |
| `INT-05` | - | #1217 | RELEASE_SECURITY_GAP | **RELEASE** |
| `INT-06` | #885 | - | RELEASE_POLICY | **OPERATOR** |
| `DOC-01` | #2173 | #2174 #2213 #2344 | RELEASE_GAP | **RELEASE** |
| `DOC-02` | - | #2084 | UNSAFE_EXAMPLE | **REJECT** |
| `DOC-03` | #1067 | - | LOW_VALUE_DOC | **DEFER** |
| `DOC-04` | #2295 | #2306 | RELEASE_GAP | **RELEASE** |

## C10 — Elective feature triage (22 packages)

| Package | Upstream issues | Upstream PRs | Status | Action |
|---|---|---|---|---|
| `FEAT-01` | - | #1992 | UNSAFE_NETWORK_FEATURE | **REJECT** |
| `FEAT-02` | #2152 #2153 | #2154 | SDK_COUPLING | **DEFER** |
| `FEAT-03` | #1869 | #2056 | AGENT_SURFACE | **DEFER** |
| `FEAT-04` | #1054 | - | PARTIAL_IMPLEMENTATION | **ACCEPT** |
| `FEAT-05` | #1013 | - | NETWORK_EGRESS_FEATURE | **REJECT** |
| `FEAT-06` | #1456 | - | SCHEMA_COMPLEXITY | **DEFER** |
| `FEAT-07` | #1923 | - | TRUST_BOUNDARY_FEATURE | **REJECT** |
| `FEAT-08` | #2130 | #2244 | PARTIAL_IMPLEMENTATION | **ACCEPT** |
| `FEAT-09` | #2066 | #2097 | LOW_VALUE_FEATURE | **DEFER** |
| `FEAT-10` | - | #1874 | NEW_LANGUAGE_FEATURE | **REJECT** |
| `FEAT-11` | #1557 #2186 | #1630 #1896 | REMOTE_TRUST_FEATURE | **REJECT** |
| `FEAT-12` | #1723 | #1724 | RULESET_POLICY | **REJECT** |
| `FEAT-13` | #983 | #1654 | PARTIAL_IMPLEMENTATION | **ACCEPT** |
| `FEAT-14` | - | #1661 | UNSAFE_PACKAGING | **REJECT** |
| `FEAT-15` | - | #1485 | UNSAFE_INSTALLER | **REJECT** |
| `FEAT-16` | #800 | - | OPTIONAL_SOURCE_MODE | **DEFER** |
| `FEAT-17` | #2068 | #2096 | INTEGRATION_POLICY | **DEFER** |
| `FEAT-18` | #1048 | - | REPORT_ARCHITECTURE | **DEFER** |
| `FEAT-19` | #1525 | - | ENCRYPTED_FORMAT_FEATURE | **REJECT** |
| `FEAT-20` | #1268 | - | NEW_SOURCE_MODE | **DEFER** |
| `FEAT-21` | #1234 | - | DUPLICATE_FEATURE | **ACCEPT** |
| `FEAT-22` | #739 | - | LOW_VALUE_BINARY_FEATURE | **REJECT** |

## OPS — External upstream incident response (3 packages)

| Package | Upstream issues | Upstream PRs | Status | Action |
|---|---|---|---|---|
| `SEC-01` | #2110 | - | EXTERNAL_INCIDENT | **ESCALATE** |
| `SEC-02` | #2111 | - | EXTERNAL_INCIDENT | **ESCALATE** |
| `SEC-03` | #2112 | - | EXTERNAL_INCIDENT | **ESCALATE** |

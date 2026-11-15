# CP-V13 Contributor and Provenance Manifest

**Integration branch:** `audit/v13-integration`
**Integration base:** `rnc-mod` at `3a4635f82dcd0a216f9d2b88ebbd72341b3b60cc`
**V13 development HEAD (post-V13A):** `77376f101a4e36a7c47897b09584d175a8825f39`
**RNC product commits:** `74a042e522c3d8646d39fd41e8a0ed88d00e5e83`, `d13e14ef7f7eb1e2684472d3c61d74b1c75b8427`, `b8131048add7a20f49ee615cc90e0e6d75641c24`
**Final product tree:** `89844b1b6cd63d76f77cedb2b9584844ee092b66`

## Upstream commit adoption

**None.** V13 introduced no new upstream-authored commit and no cherry-pick. All new implementation is RNC-authored on the development branch and appears in the single RNC product tree.

This is not a claim that the work was uninformed by upstream discussion. The following design/research influences are acknowledged in `RNC-MOD.md` Section 3 and are not commit-level attribution:

| Source | Contributor | Use |
|---|---|---|
| [#2128](https://github.com/gitleaks/gitleaks/pull/2128) | puneetdixit200 | Quoted `--log-opts` problem/regression research; RNC wrote an independent tokenizer and parsed-vector validation. |
| [#1809](https://github.com/gitleaks/gitleaks/pull/1809) | JoostVoskuil | Azure DevOps SSH remote identity research; original patch targets removed files and was not transplanted. |
| [#1162](https://github.com/gitleaks/gitleaks/pull/1162) | mn3m0nic | Repeated-finding/secret-hash proposal; unkeyed hash rejected, keyed grouping remains policy-gated. |
| [#1232](https://github.com/gitleaks/gitleaks/pull/1232) | Coruscant11 | Multiple-report-writer proposal; safe report-plan architecture documented, no legacy patch adopted. |

No co-author trailer was fabricated. No contributor commit was amended, squashed or replayed.

## Reconciliation

The integration product commit was created by restoring the frozen V13 development tree and verifying object equality:

```text
development tree: 89844b1b6cd63d76f77cedb2b9584844ee092b66
integration tree: 89844b1b6cd63d76f77cedb2b9584844ee092b66 (before documentation-only provenance additions)
```

The provenance reports are documentation-only additions after the product equality point. The final PR target is `rnc-mod`; no merge or release publication is performed.

## Licensing

The upstream MIT license and copyright notice remain in place. Git attribution, design-influence recognition and this manifest do not replace that license.

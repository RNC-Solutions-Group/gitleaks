# Closed — RNC disposition ledger

This ledger is the **single exclusion authority** for future `_backlog/categories/` active candidate registers. “Closed” means **no further RNC work is currently planned for the stated scope**, not that a maintainer closed the upstream issue. Never edit GitHub issue state automatically.

## Closed for RNC

| Upstream item | Source status | RNC status | Reason | Evidence | Reopen condition |
| --- | --- | --- | --- | --- | --- |
| [#1939](https://github.com/gitleaks/gitleaks/issues/1939) — Vulnerabilities in Gitleaks binary | **OPEN upstream** (checked 2026-11-08) | **CLOSED_RNC — owner confirmed** | The report concerned stale Go stdlib plus `ulikunitz/xz`. Current RNC hardened builder uses Go **1.27.3**, while fork `master/go.mod` already pins `ulikunitz/xz v0.5.15` (the fixed version cited by the issue, rather than the report's v0.5.12). | [Hardened Dockerfile](https://github.com/RNC-Solutions-Group/gitleaks/blob/fix/docker-security-bases/Dockerfile) · [Pinned go.mod](https://github.com/RNC-Solutions-Group/gitleaks/blob/master/go.mod) · Owner's latest image assessment. **Current image SBOM/Trivy result not independently attached to this ledger.** | Reopen if a pinned RNC image digest still reports either issue-specific CVE or the runtime build provenance does not match the hardened builder. |

| [PR #1435](https://github.com/gitleaks/gitleaks/pull/1435) — legacy chunk-read change by @sapienhwaker | **OPEN upstream** (checked 2026-11-08); related [issue #1432](https://github.com/gitleaks/gitleaks/issues/1432) **CLOSED upstream** (2025-01-13) | **SUPERSEDED_RNC — reviewed** | The 2024 PR proposes a 100-byte chunk overlap for the legacy 10 KB file reader. It was never merged; the proposed slice retains an alias into a reused read buffer, and no automated tests or formal reviews were found for its head. Later merged PRs [#1652](https://github.com/gitleaks/gitleaks/pull/1652), [#1665](https://github.com/gitleaks/gitleaks/pull/1665), and [#1676](https://github.com/gitleaks/gitleaks/pull/1676) addressed chunking and added tests. [#1872](https://github.com/gitleaks/gitleaks/pull/1872) subsequently replaced the legacy `detect/directory.go` path. | [Original PR](https://github.com/gitleaks/gitleaks/pull/1435) · [issue resolution](https://github.com/gitleaks/gitleaks/issues/1432) · [current fork file reader](https://github.com/RNC-Solutions-Group/gitleaks/blob/master/sources/file.go) · [boundary tests](https://github.com/RNC-Solutions-Group/gitleaks/blob/master/sources/common_test.go). **Review and source comparison only; no fresh runtime reproduction performed.** | Only if current-master testing reproduces a missed detection at chunk boundaries; investigate against the current reader rather than cherry-picking #1435. |

### Exclusion contract

- Canonical IDs: `upstream_issue:gitleaks/gitleaks#1939` and `upstream_pr:gitleaks/gitleaks#1435`.
- Linked issue #1432 was already **closed upstream** and is not present in the open-issue inventory; it is contextual evidence for PR #1435, **not** a second RNC-closed work item.
- Exclude these items from all active category tables; reference this ledger instead. PR #1435 is **superseded**, not adopted or corrected by RNC.
- Do not use the upstream GitHub `state=open` as a reason to reinstate an RNC-closed issue automatically.
- A closure needs: an RNC disposition, precise scope, artifact/commit or image digest, evidence, and an explicit reopen trigger. “Resolved in RNC” does not mean “fixed upstream.”
- When refreshing the CSV or JSON snapshots, preserve the closure decision in a **separate reviewed-dispositions layer**. The raw upstream metadata remains immutable.

**Follow-up evidence improvement:** capture the exact current binary/image digest and a reproducible scan showing the two reported issue-specific CVEs absent. That will upgrade this from owner-confirmed closure to scan-attested closure without changing its current RNC disposition.

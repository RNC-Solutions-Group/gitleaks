# Closed — RNC disposition ledger

This ledger is the **single exclusion authority** for future `_backlog/categories/` active candidate registers. “Closed” means **no further RNC work is currently planned for the stated scope**, not that a maintainer closed the upstream issue. Never edit GitHub issue state automatically.

## Closed for RNC

| Upstream issue | Source status | RNC status | Reason | Evidence | Reopen condition |
| --- | --- | --- | --- | --- | --- |
| [#1939](https://github.com/gitleaks/gitleaks/issues/1939) — Vulnerabilities in Gitleaks binary | **OPEN upstream** (checked 2026-11-08) | **CLOSED_RNC — owner confirmed** | The report concerned stale Go stdlib plus `ulikunitz/xz`. Current RNC hardened builder uses Go **1.27.3**, while fork `master/go.mod` already pins `ulikunitz/xz v0.5.15` (the fixed version cited by the issue, rather than the report's v0.5.12). | [Hardened Dockerfile](https://github.com/RNC-Solutions-Group/gitleaks/blob/fix/docker-security-bases/Dockerfile) · [Pinned go.mod](https://github.com/RNC-Solutions-Group/gitleaks/blob/master/go.mod) · Owner's latest image assessment. **Current image SBOM/Trivy result not independently attached to this ledger.** | Reopen if a pinned RNC image digest still reports either issue-specific CVE or the runtime build provenance does not match the hardened builder. |

### Exclusion contract

- Canonical ID: `upstream_issue:gitleaks/gitleaks#1939`.
- Exclude this issue from active dependencies, security-hardening, and future category tables; reference this ledger instead.
- Do not use the upstream GitHub `state=open` as a reason to reinstate an RNC-closed issue automatically.
- A closure needs: an RNC disposition, precise scope, artifact/commit or image digest, evidence, and an explicit reopen trigger. “Resolved in RNC” does not mean “fixed upstream.”
- When refreshing the CSV or JSON snapshots, preserve the closure decision in a **separate reviewed-dispositions layer**. The raw upstream metadata remains immutable.

**Follow-up evidence improvement:** capture the exact current binary/image digest and a reproducible scan showing the two reported issue-specific CVEs absent. That will upgrade this from owner-confirmed closure to scan-attested closure without changing its current RNC disposition.

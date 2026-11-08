# Category review notes

This directory is reserved for human-reviewed analyses split from the [crosswalk](../03-pr-issue-crosswalk.md). The original upstream snapshot was unreviewed; the category documents below record subsequent **RNC-specific** dispositions. Those assessments do not change upstream issue or PR state.

Current reviews:

- [Crosswalk duplicates](crosswalk-duplicates.md): competing PRs for the same issue, overlapping issue reports, and references that should **not** be treated as duplicates. Candidate matches are not implementations.
- [Dependencies](dependencies.md): Go modules, compiler/builder versions, Docker and CI dependency candidates, including local supersession checks.
- [Features](features.md): functional and developer-experience changes, starting with legacy Go-install compatibility PR #652.
- [Closed](closed.md): the authoritative **RNC-only exclusion register**, currently #1939; upstream remains open.
- [Security hardening](security-hardening.md): PRs #2308 and #2336 with their linked issues; **reviewed — testing**, not yet integrated.

## Single-owner category routing

Each upstream issue or PR has **exactly one canonical RNC category** at a time. A related issue and PR are **one work package**, and must reside together in that same category; a category may contain both issues and PRs. Moving an item removes its active listing from the old category; the [closed ledger](closed.md) overrides active categories. Incidental references to another category's candidates are kept to a minimum and do **not** establish a second ownership assignment.

| Canonical category | Issues | PRs | Scope |
| --- | ---: | ---: | --- |
| [Crosswalk duplicates](crosswalk-duplicates.md) | 33 | 48 | Competing proposals and nearby-but-distinct reports awaiting review |
| [Dependencies](dependencies.md) | 11 | 18 | Build/toolchain, pinned modules, CI and related issues |
| [Features](features.md) | 0 | 1 | Legacy installation and developer experience; pending review |
| [Security hardening](security-hardening.md) | 3 | 2 | Two reviewed security PRs and their directly linked issues; testing |
| [Closed](closed.md) | 1 | 0 | No more current RNC work; upstream may remain open |
| **Unique assigned IDs** | **48** | **69** | Two duplicate-category issues fall outside the open-issue snapshot |

The three top-level Markdown reports ([issues](../01-issue-backlog.md), [PRs](../02-open-pull-requests.md), and [crosswalk](../03-pr-issue-crosswalk.md)) ~~strike through~~ items once assigned to a category. Strikethrough is **routing**, not upstream closure or implemented code. Their original counts remain historical snapshot counts. JSON and CSV source exports remain unchanged.

Use one Markdown file per meaningful topic (for example `scanner-correctness.md` or `output-safety.md`). Record for each considered issue/PR:

| Field | Required evidence |
| --- | --- |
| Issue(s) and PR(s) | Original upstream links and contributor attribution |
| Scope | Security defect, scanning behavior, build dependency, feature, or unrelated |
| Root cause | Reproduction and why the behavior occurs |
| Proposed fix | Existing PR SHA or our independent patch commit |
| Baseline | Base SHA, container image digest, test fixture revision |
| Tests | Commands, assertions, complete fixture coverage, failures and exclusions |
| Integration | Branch, cherry-pick SHA, conflicts, follow-up changes |
| Disposition | Adopt / defer / reject / duplicate / superseded / blocked / failed test |
| Rationale | Explicit *why*, including side effects or unsuccessful results |
| Promotion | User decision and destination branch; no automatic merges |

Security-only contributions from Dave remain separate from combined community security contributions; general enhancements are kept downstream. This backlog does not change branch policy and does not itself promote code.

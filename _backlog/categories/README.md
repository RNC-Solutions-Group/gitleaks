# Category review notes

This directory is reserved for human-reviewed analyses split from the [crosswalk](../03-pr-issue-crosswalk.md). No candidate has been accepted, rejected, tested, or promoted by this initial inventory.

Use one Markdown file per meaningful topic (for example `security-hardening.md`, `scanner-correctness.md`, `dependency-updates.md`, or `output-safety.md`). Record for each considered issue/PR:

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
